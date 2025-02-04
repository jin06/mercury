package mqtt

import (
	"io"
)

func NewConnect(header *FixedHeader) *Connect {
	return &Connect{FixHeader: header}
}

type Connect struct {
	Version      ProtocolVersion
	FixHeader    *FixedHeader
	ProtocolName string
	//Clean Clean Session(v3,v4) or Clean Start(v5)
	UserNameFlag bool
	PasswordFlag bool
	WillFlag     bool
	Will         *Will
	Clean        bool
	KeepAlive    uint16
	ClientID     string
	Username     string
	Password     string
	Properties   *Properties
}

func (c *Connect) Encode() (result []byte, err error) {
	//Fixed header
	result = toHeader(CONNECT)
	//Variable header
	var buf []byte
	if buf, err = strToBytes(c.protocolName()); err != nil {
		return
	}
	result = append(result, buf...)
	result = append(result, byte(c.Version))
	result = append(result, c.encodeFlag())
	result = append(result, uint16ToBytes(c.KeepAlive)...)
	if buf, err = c.Properties.Encode(); err != nil {
		result = append(result, buf...)
	}
	if c.ClientID == "" {
		return nil, ErrNullClientID
	}
	if buf, err = strToBytes(c.ClientID); err != nil {
		return
	}
	result = append(result, buf...)
	if c.Will != nil {
		if c.Version == MQTT5 {
			// if buf, err = encodeWillProperties(c.Will.Properties); err != nil {
			// 	return
			// }
			result = append(result, buf...)
		}
		if buf, err = strToBytes(c.Will.Topic); err != nil {
			return
		}
		result = append(result, buf...)
		if buf, err = strToBytes(c.Will.Message); err != nil {
			return
		}
		result = append(result, buf...)
	}
	if buf, err = strToBytes(c.Username); err != nil {
		return
	}
	result = append(result, buf...)
	if buf, err = strToBytes(c.Password); err != nil {
		return
	}
	result = append(result, buf...)

	return
}

func (c *Connect) EncodeBody() ([]byte, error) {
	return nil, nil
}

func (c *Connect) Decode(data []byte) (int, error) {
	l, err := c.FixHeader.Decode(data)
	if err != nil {
		return 0, err
	}
	return c.DecodeBody(data[l:])
}

func (c *Connect) DecodeBody(data []byte) (int, error) {
	// total := len(data)
	var start int

	if n, err := c.FixHeader.Decode(data); err != nil {
		return start, err
	} else {
		start = n
	}

	if protocol, n, err := decodeUTF8Str(data[start:]); err != nil {
		return start, err
	} else {
		c.ProtocolName = protocol
		start = start + n
	}

	if version, err := decodeProtocolVersion(data[start]); err != nil {
		return start, err
	} else {
		c.Version = version
		start = start + 1
	}

	{
		flag := data[start]
		c.UserNameFlag = (flag&0b1000000 == 0b1000000)
		c.PasswordFlag = (flag&0b01000000 == 0b01000000)
		if flag&0b00000100 == 0b00000100 {
			c.Will = &Will{}
			c.Will.Retain = (flag&0b00100000 == 0b00100000)
			c.Will.QoS = QoS(flag & 0b00011000)
		}
		c.Clean = (flag&0b00000010 == 0b00000010)
		c.WillFlag = (flag&0b00000100 == 0b00000100)
		if c.WillFlag {
			c.Will = &Will{
				Retain:     flag&0b00100000 == 0b00100000,
				QoS:        QoS(flag & 0b00011000),
				Properties: new(Properties),
			}
		}

		start = start + 1
	}

	{
		c.KeepAlive = decodeKeepAlive(data[start : start+2])
		start = start + 2
	}

	if c.Version.IsMQTT5() {
		c.Properties = new(Properties)
		n, err := c.Properties.Decode(data[start:])
		if err != nil {
			return start, err
		}
		start = start + n
	}
	{
		clientID, n, err := decodeUTF8Str(data[start:])
		if err != nil {
			return start, err
		}
		c.ClientID = clientID
		start = start + n
	}
	if c.WillFlag {
		if c.Version.IsMQTT5() {
			c.Will.Properties.Decode(data[start:])
		}
		if topic, n, err := decodeUTF8Str(data[start:]); err != nil {
			return start, err
		} else {
			c.Will.Topic = topic
			start = start + n
		}
		if message, n, err := decodeUTF8Str(data[start:]); err != nil {
			return start, err
		} else {
			c.Will.Message = message
			start = start + n
		}
	}
	if c.UserNameFlag {
		if user, n, err := decodeUTF8Str(data[start:]); err != nil {
			return start, err
		} else {
			c.Username = user
			start = start + n
		}
	}
	if c.PasswordFlag {
		if pass, n, err := decodeUTF8Str(data[start:]); err != nil {
			return start, err
		} else {
			c.Password = pass
			start = start + n
		}
	}

	return start, nil
}

func (c *Connect) Read(reader io.Reader) (err error) {
	// buf := make([]byte, 16)
	if _, err := readProtocolName(reader); err != nil {
		return err
	}
	buf := make([]byte, 4)
	if _, err = io.ReadFull(reader, buf); err != nil {
		return
	}

	c.Version = ProtocolVersion(buf[0])
	// var usernameFlag, passFlag, willRetain, willQoS, willFlag, cleanStart, reserved bool
	usernameFlag := buf[1]&0b10000000 == 0b10000000
	passFlag := buf[1]&0b01000000 == 0b01000000
	if buf[1]&0b00000100 == 0b00000100 {
		c.Will = &Will{Properties: new(Properties)}
		c.Will.Retain = (buf[1]&0b00100000 == 0b00100000)
		c.Will.QoS = QoS(buf[1] & 0b00011000)
	}
	c.Clean = buf[1]&0b00000010 == 0b00000010
	// reserved := buf[8] & 0b00000001
	// measured in seconds
	c.KeepAlive = decodeKeepAlive(buf[2:])
	if c.Version == MQTT5 {
		properties, err := readProperties(reader)
		if err != nil {
			return err
		}
		c.Properties = properties
	}
	if c.ClientID, err = readStr(reader); err != nil {
		return
	}
	if c.Will != nil {
		if c.Version == MQTT5 {
			// if c.Will.Properties, err = decodeWillProperties(reader); err != nil {
			// return
			// }
		}
		if c.Will.Topic, err = readStr(reader); err != nil {
			return
		}
		if c.Will.Message, err = readStr(reader); err != nil {
			return
		}
	}
	if usernameFlag {
		if c.Username, err = readStr(reader); err != nil {
			return
		}
	}
	if passFlag {
		if c.Password, err = readStr(reader); err != nil {
			return
		}
	}

	return
}

func (c *Connect) ReadBody(r io.Reader) error {
	return nil
}
func (c *Connect) Write(reader io.Writer) error {
	return nil
}

func (c *Connect) WriteBody(w io.Writer) error {
	return nil
}

func (c *Connect) protocolName() string {
	switch c.Version {
	case MQTT3:
		return "MQIsdp"
	case MQTT4:
		return "MQTT"
	case MQTT5:
		return "MQTT"
	}
	return ""
}

func (c *Connect) encodeFlag() byte {
	var flags byte
	if c.Username != "" {
		flags = flags | 0b1
	}
	if c.Password != "" {
		flags = flags | 0b01
	}
	if c.Will != nil {
		flags = flags | (byte(c.Will.QoS) << 3)
		flags = flags | 0b00000010
	}
	if c.Clean {
		flags = flags | 0b00000001
	}
	return flags
}

func (c *Connect) String() string {
	return ""
}
