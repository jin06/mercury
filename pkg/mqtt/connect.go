package mqtt

import "fmt"

func NewConnect(header *FixedHeader, v ProtocolVersion) *Connect {
	return &Connect{
		BasePacket: &BasePacket{header, v},
		// Version:    v,
		Properties: new(Properties),
	}
}

type Connect struct {
	*BasePacket
	// Version      ProtocolVersion
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

func (c *Connect) String() string {
	return fmt.Sprintf("Connect - Version: %s, ProtocolName: %s, Clean: %t, KeepAlive: %d, ClientID: %s, Username: %s, Password: %s",
		c.Version, c.ProtocolName, c.Clean, c.KeepAlive, c.ClientID, c.Username, c.Password)
}

func (c *Connect) Response() *Connack {
	resp := &Connack{
		BasePacket: newBasePacket(CONNACK, c.Version),
		ReasonCode: RET_CONNACK_ACCEPT,
	}
	resp.ReasonCode = V5_SUCCESS

	return resp
}

func (c *Connect) Encode() ([]byte, error) {
	body, err := c.EncodeBody()
	if err != nil {
		return nil, err
	}
	c.FixedHeader.RemainingLength = VariableByteInteger(len(body))
	header, err := c.FixedHeader.Encode()
	if err != nil {
		return nil, err
	}
	return append(header, body...), nil
}

func (c *Connect) EncodeBody() ([]byte, error) {
	var data []byte

	// Encode Protocol Name
	if protocolData, err := encodeUTF8Str(c.ProtocolName); err != nil {
		return nil, err
	} else {
		data = append(data, protocolData...)
	}

	// Encode Protocol Version
	if version, err := encodeProtocolVersion(c.Version); err != nil {
		return nil, err
	} else {
		data = append(data, version)
	}

	// Encode Connect Flags
	if flagData, err := c.encodeFlag(); err != nil {
		return nil, err
	} else {
		data = append(data, flagData)
	}

	// Encode KeepAlive
	data = append(data, encodeKeepAlive(c.KeepAlive)...)

	// Encode Properties (MQTT 5.0 only)
	if c.Version.IsMQTT5() && c.Properties != nil {
		propertiesData, err := c.Properties.Encode()
		if err != nil {
			return nil, err
		}
		data = append(data, propertiesData...)
	}

	// Encode Client ID
	if clientIDData, err := encodeUTF8Str(c.ClientID); err != nil {
		return nil, err
	} else {
		data = append(data, clientIDData...)
	}

	// Encode Will Message (if Will Flag is set)
	if c.WillFlag {
		if c.Version.IsMQTT5() {
			propertiesData, err := c.Will.Properties.Encode()
			if err != nil {
				return nil, err
			}
			data = append(data, propertiesData...)
		}
		if willTopicData, err := encodeUTF8Str(c.Will.Topic); err != nil {
			return nil, err
		} else {
			data = append(data, willTopicData...)
		}
		if willMessageData, err := encodeUTF8Str(c.Will.Message); err != nil {
			return nil, err
		} else {
			data = append(data, willMessageData...)
		}
	}

	// Encode Username (if Username Flag is set)
	if c.UserNameFlag {
		if usernameData, err := encodeUTF8Str(c.Username); err != nil {
			return nil, err
		} else {
			data = append(data, usernameData...)
		}
	}

	// Encode Password (if Password Flag is set)
	if c.PasswordFlag {
		if passwordData, err := encodeUTF8Str(c.Password); err != nil {
			return nil, err
		} else {
			data = append(data, passwordData...)
		}
	}

	return data, nil
}

func (c *Connect) Decode(data []byte) (int, error) {
	n, err := c.FixedHeader.Decode(data)
	if err != nil {
		return 0, err
	}
	bodyLen, err := c.DecodeBody(data[n:])
	return bodyLen + n, err
}

func (c *Connect) DecodeBody(data []byte) (int, error) {
	var start int

	// Decode Protocol Name
	if protocol, n, err := decodeUTF8Str(data[start:]); err != nil {
		return start, err
	} else {
		c.ProtocolName = protocol
		start += n
	}

	// Decode Protocol Version
	if version, err := decodeProtocolVersion(data[start]); err != nil {
		return start, err
	} else {
		c.Version = version
		start++
	}
	// Decode Connect Flags
	c.decodeFlag(data[start])
	start++

	// Decode KeepAlive
	c.KeepAlive = decodeKeepAlive(data[start : start+2])
	start += 2

	// Decode Properties (MQTT 5.0 only)
	if c.Version.IsMQTT5() {
		c.Properties = new(Properties)
		n, err := c.Properties.Decode(data[start:])
		if err != nil {
			return start, err
		}
		start += n
	}
	{
		// Decode Client ID
		clientID, n, err := decodeUTF8Str(data[start:])
		if err != nil {
			return start, err
		}
		c.ClientID = clientID
		start += n
	}
	// Decode Will Message (if Will Flag is set)
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
		// Decode Username (if Username Flag is set)
		if user, n, err := decodeUTF8Str(data[start:]); err != nil {
			return start, err
		} else {
			c.Username = user
			start = start + n
		}
	}
	if c.PasswordFlag {
		// Decode Password (if Password Flag is set)
		if pass, n, err := decodeUTF8Str(data[start:]); err != nil {
			return start, err
		} else {
			c.Password = pass
			start = start + n
		}
	}

	return start, nil
}

func (c *Connect) Read(r *Reader) (err error) {
	c.FixedHeader = new(FixedHeader)
	if err := c.FixedHeader.Read(r); err != nil {
		return err
	}
	return c.ReadBody(r)
}

func (c *Connect) ReadBody(r *Reader) error {
	data, err := r.Read(c.Length())
	if err != nil {
		return err
	}
	_, err = c.DecodeBody(data)
	return err
}

func (c *Connect) Write(writer *Writer) error {
	data, err := c.Encode()
	if err != nil {
		return err
	}
	_, err = writer.Write(data)
	return err
}

func (c *Connect) WriteBody(w *Writer) error {
	data, err := c.EncodeBody()
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

func (c *Connect) encodeFlag() (byte, error) {
	var flag byte
	if c.UserNameFlag {
		flag = flag | 0b1000000
	}
	if c.PasswordFlag {
		flag = flag | 0b01000000
	}
	if c.WillFlag && c.Will != nil {
		flag = flag | 0b00100100
		flag = flag | byte((c.Will.QoS&0b00011000)<<3)
	}
	if c.Clean {
		flag = flag | 0b00000010
	}
	return flag, nil
}

func (c *Connect) decodeFlag(flag byte) {
	c.UserNameFlag = (flag&0b1000000 == 0b1000000)
	c.PasswordFlag = (flag&0b01000000 == 0b01000000)
	if flag&0b00000100 == 0b00000100 {
		c.Will = &Will{}
		c.Will.Retain = (flag&0b00100000 == 0b00100000)
		c.Will.QoS = QoS((flag & 0b00011000) >> 3)
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
}
