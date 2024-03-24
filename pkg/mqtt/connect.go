package mqtt

import (
	"fmt"
	"io"

	"github.com/jin06/mercury/pkg/utils"
)

type Connect struct {
	Version           ProtocolVersion
	Clean             bool // v4
	KeepAlive         uint16
	ClientID          []byte
	Username          string
	Password          string
	WillQoS           bool
	WillFlag          bool
	WillRetain        bool // v4
	WillDelay         uint32
	WillProperties    map[string]string
	WillTopic         []byte
	WillMessage       []byte
	willDelayInterval uint64 // only mqtt5
	MessageExpiry     uint64
	// SessionExpiryInterval      uint64 // mqtt5
	// RequestResponseInformation byte
	// RequestProblemInformation  byte
	PayloadFormat   bool
	ContentType     string
	ResponseTopic   string
	CorrelationData []byte
	UserProperty    map[string]string
	Properties      Properties
}

// mqtt5
type Properties struct {
	RequestProblemInformation  byte
	RequestResponseInformation byte
	SessionExpiryInterval      uint32
	//ReceiveMaximum The Client uses this value to limit the number of QoS 1 and QoS 2 publications that it is willing to process concurrently.
	ReceiveMaximum uint16
	// MaximumPacketSize The packet size is the total number of bytes in an MQTT Control Packet
	MaximumPacketSize uint32
	TopicAliasMax     uint16
	UserProperties    UserProperties
}

type UserProperties map[string]string

func (c *Connect) String() string {
	return ""
}

func (c *Connect) Decode(reader io.Reader) (err error) {
	// buf := make([]byte, 16)

	protocolName, err := decodeProtocolName(reader)
	if err != nil {
		return
	}
	fmt.Println(protocolName)
	buf := make([]byte, 4)
	if _, err = reader.Read(buf); err != nil {
		return
	}

	c.Version = ProtocolVersion(buf[0])
	// var usernameFlag, passFlag, willRetain, willQoS, willFlag, cleanStart, reserved bool
	usernameFlag := buf[1]&0b10000000 == 0b10000000
	passFlag := buf[1]&0b01000000 == 0b01000000
	// willRetain := buf[8] & 0b00100000
	// willQoS := QoS(buf[8] & 0b00011000)
	willFlag := buf[1]&0b00000100 == 0b00000100
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
	if c.ClientID, err = decodeUTF8(reader); err != nil {
		return
	}
	switch c.Version {
	case MQTT3, MQTT4:
		{
			if willFlag {
				if c.WillTopic, err = decodeUTF8(reader); err != nil {
					return
				}
				if c.WillMessage, err = decodeUTF8(reader); err != nil {
					return
				}
			}
			if usernameFlag {
				if c.Username, err = decodeUTF8Str(reader); err != nil {
					return
				}
			}
			if passFlag {
				if c.Password, err = decodeUTF8Str(reader); err != nil {
					return
				}

			}
		}
	case MQTT5:
		{
			// read property length
			proLen := make([]byte, 1)
			if _, err = reader.Read(proLen); err != nil {
				return
			}
			willProperties := make([]byte, int(proLen[0]))
			if _, err = reader.Read(willProperties); err != nil {
				return
			}
			for {
				if len(willProperties) <= 0 {
					break
				}
				idetifier := willProperties[0]
				switch idetifier {
				case 0x18:
					{
						if c.willDelayInterval, err = bytesToUint64(willProperties[1:5]); err != nil {
							return
						}
					}
					// payload format
				case 0x01:
					{

					}
					// publication expiry interval
				case 0x02:
					{

					}
					// content type
				case 0x03:
					{

					}
					// reponse topic
				case 0x08:
					{

					}
					// correlation data
				case 0x09:
					{

					}
				}

			}

			if c.willDelayInterval, err = bytesToUint64(willProperties[0:4]); err != nil {
				return
			}
		}
	}
	return
}

func decodeProperties(reader io.Reader) (result Properties, err error) {
	result = Properties{}
	var total int
	if res, err := readByte(reader); err != nil {
		return result, err
	} else {
		total = int(res)
	}

	for i := 0; i < total; {
		var identifier byte
		identifier, err = readByte(reader)
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
				userProperties := UserProperties{}
				list := []string{}
				for j := 0; j < total-i; {
					var val string
					var n int
					if val, n, err = readStr(reader); err != nil {
						return
					}
					j = j + n
					list = append(list, val)
				}
				if len(list)%2 == 1 {
					err = ProtocolError
					return
				}
				for i := 0; i < len(list); i += 2 {
					userProperties[list[i]] = userProperties[list[i+1]]
				}
				break
			}
		}
	}
	return
}

func decodeUserProperties(s []byte) (err error) {
	return
}

func readUint64(reader io.Reader) (uint64, error) {
	res, err := read(reader, 8)
	if err != nil {
		return 0, err
	}
	return utils.ToUint64(res)
}

func readUint32(reader io.Reader) (uint32, error) {
	res, err := read(reader, 4)
	if err != nil {
		return 0, err
	}
	return utils.ToUint32(res)
}

func readUint16(reader io.Reader) (uint16, error) {
	res, err := read(reader, 2)
	if err != nil {
		return 0, err
	}
	return utils.ToUint16(res)
}

func readByte(reader io.Reader) (byte, error) {
	res, err := read(reader, 1)
	if err != nil {
		return 0, err
	}
	return res[0], nil
}

func readStr(reader io.Reader) (str string, n int, err error) {
	var l uint16
	if l, err = readUint16(reader); err != nil {
		return
	}
	var res []byte
	if res, err = read(reader, int(l)); err != nil {
		return
	}
	str = string(res)
	n = int(l + 2)
	return
}

func read(reader io.Reader, n int) ([]byte, error) {
	res := make([]byte, n)
	_, err := io.ReadFull(reader, res)
	return res, err
}
