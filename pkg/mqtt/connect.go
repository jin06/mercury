package mqtt

import (
	"errors"
	"io"
)

type Connect struct {
	Version ProtocolVersion
	//Clean Clean Session(v3,v4) or Clean Start(v5)
	Clean      bool
	KeepAlive  uint16
	ClientID   string
	Username   string
	Password   string
	Properties *Properties
	Will       *Will
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

func (c *Connect) Decode(data []byte) (err error) {
	return
}

func (c *Connect) Read(reader io.Reader) error {
	return nil
}

func (c *Connect) Write(reader io.Reader) (err error) {
	// buf := make([]byte, 16)

	if _, err := decodeProtocolName(reader); err != nil {
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
		c.Will = &Will{}
		c.Will.Retain = (buf[1]&0b00100000 == 0b00100000)
		c.Will.QoS = QoS(buf[1] & 0b00011000)
	}
	c.Clean = buf[1]&0b00000010 == 0b00000010
	// reserved := buf[8] & 0b00000001
	// measured in seconds
	c.KeepAlive = decodeKeepAlive(buf[2:])
	if c.Version == MQTT5 {
		properties, err := decodeProperties(reader)
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
			if c.Will.Properties, err = decodeWillProperties(reader); err != nil {
				return
			}
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
	if buf, err = encodeProperties(c.Properties); err != nil {
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
			if buf, err = encodeWillProperties(c.Will.Properties); err != nil {
				return
			}
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

type UserProperties map[string]string

func (u *UserProperties) toBytes() (result []byte, err error) {
	result = []byte{}
	for key, val := range *u {
		if bytes, err := strToBytes(key); err != nil {
			return result, err
		} else {
			result = append(result, bytes...)
		}
		if bytes, err := strToBytes(val); err != nil {
			return result, err
		} else {
			result = append(result, bytes...)
		}
	}
	return
}

func (u *UserProperties) fromReader(reader io.Reader) error {
	propertyLength, err := readUint8(reader)
	if err != nil {
		return err
	}
	// buf, err := readBytes(reader, int(propertyLength))
	arr := []string{}
	for i := 0; i <= int(propertyLength); {
		str, n, err := readStrN(reader)
		if err != nil {
			return err
		}
		arr = append(arr, str)
		i += n
	}
	if len(arr)%2 == 1 {
		return errors.New("count of string can't fulfil requirements of key-value pair")
	}
	for i := 0; i < len(arr); i += 2 {
		(*u)[arr[i]] = arr[i+1]
	}
	return nil
}

// mqtt5
type Properties struct {
	RequestProblemInformation  *byte
	RequestResponseInformation *byte
	// SessionExpiryInterval second
	SessionExpiryInterval *uint32
	//ReceiveMaximum The Client uses this value to limit the number of QoS 1 and QoS 2 publications that it is willing to process concurrently.
	ReceiveMaximum *uint16
	// MaximumPacketSize The packet size is the total number of bytes in an MQTT Control Packet
	MaximumPacketSize *uint32
	TopicAliasMax     *uint16
	UserProperties    *UserProperties
}
