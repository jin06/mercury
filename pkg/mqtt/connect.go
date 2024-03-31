package mqtt

import (
	"fmt"
	"io"
)

type Connect struct {
	Version    ProtocolVersion
	Clean      bool // v4
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

func (c *Connect) Decode(reader io.Reader) (err error) {
	// buf := make([]byte, 16)

	if protocolName, err := decodeProtocolName(reader); err != nil {
		return err
	} else {
		fmt.Println(protocolName)
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

func decodeProperties(reader io.Reader) (result *Properties, err error) {
	result = &Properties{
		UserProperties: make(UserProperties),
	}
	var total int
	if res, err := readByte(reader); err != nil {
		return result, err
	} else {
		total = int(res)
	}

	for i := 0; i < total; {
		var identifier byte
		if identifier, err = readByte(reader); err != nil {
			return
		}
		i++

		switch identifier {
		case 0x11:
			{
				i += 4
				if result.SessionExpiryInterval, err = readUint32(reader); err != nil {
					return
				}
			}
		case 0x19:
			{
				i++
				result.RequestResponseInformation, err = readByte(reader)
				if err != nil {
					return
				}
			}
		case 0x17:
			{
				i++
				result.RequestProblemInformation, err = readByte(reader)
				if err != nil {
					return
				}
			}
			// receive max
		case 0x21:
			{
				i += 2
				if result.ReceiveMaximum, err = readUint16(reader); err != nil {
					return
				}
			}
			// Max packet size
		case 0x27:
			{
				i += 4
				if result.MaximumPacketSize, err = readUint32(reader); err != nil {
					return
				}
			}
			//  Topic Alias Max
		case 0x22:
			{
				i += 2
				if result.TopicAliasMax, err = readUint16(reader); err != nil {
					return
				}
			}
			// User properties
		case 0x26:
			{
				var ul int
				list := []string{}
				for j := 0; j < total-i; {
					var val string
					var n int
					if val, n, err = readStrN(reader); err != nil {
						return
					}
					ul += n
					j = j + n
					list = append(list, val)
				}
				if len(list)%2 == 1 {
					return result, ErrProtocol
				}
				for i := 0; i < len(list); i += 2 {
					result.UserProperties[list[i]] = result.UserProperties[list[i+1]]
				}
				i += ul
				return
			}
		}
	}
	return
}

func (c *Connect) Encode() (result []byte, err error) {
	//Fixed header
	result = toHeader(CONNECT)
	//Variable header
	var buf []byte
	if buf, err = encodeString(c.protocolName()); err != nil {
		return
	}
	result = append(result, buf...)
	result = append(result, byte(c.Version))
	result = append(result, c.encodeFlag())
	result = append(result, encodeUint16(c.KeepAlive)...)
	if c.ClientID == "" {
		return nil, ErrNullClientID
	}
	if buf, err = encodeString(c.ClientID); err != nil {
		return
	}
	result = append(result, buf...)
	if c.Will != nil {
		if buf, err = encodeString(c.Will.Topic); err != nil {
			return
		}
		result = append(result, buf...)
		if buf, err = encodeString(c.Will.Message); err != nil {
			return
		}
		result = append(result, buf...)
	}
	if buf, err = encodeString(c.Username); err != nil {
		return
	}
	result = append(result, buf...)
	if buf, err = encodeString(c.Password); err != nil {
		return
	}
	result = append(result, buf...)
	return
}

type UserProperties map[string]string

// mqtt5
type Properties struct {
	RequestProblemInformation  byte
	RequestResponseInformation byte
	// SessionExpiryInterval second
	SessionExpiryInterval uint32
	//ReceiveMaximum The Client uses this value to limit the number of QoS 1 and QoS 2 publications that it is willing to process concurrently.
	ReceiveMaximum uint16
	// MaximumPacketSize The packet size is the total number of bytes in an MQTT Control Packet
	MaximumPacketSize uint32
	TopicAliasMax     uint16
	UserProperties    UserProperties
}
