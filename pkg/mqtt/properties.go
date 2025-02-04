package mqtt

import (
	"errors"
	"io"
)

const (
	IDPayloadFormat              byte = 0x01
	IDMessageExpiry              byte = 0x02
	IDContentType                byte = 0x03
	IDResponseTopic              byte = 0x08
	IDCorrelationData            byte = 0x09
	IDSubscriptionIdentifier     byte = 0x0B
	IDSessionExpiryInterval      byte = 0x11
	IDAssignedClientID           byte = 0x12
	IDServerKeepAlive            byte = 0x13
	IDAuthMethod                 byte = 0x15
	IDAuthData                   byte = 0x16
	IDRequestProblemInfo         byte = 0x17
	IDWillDelayInterval          byte = 0x18
	IDRequestResponseInformation byte = 0x19
	IDResponseInfo               byte = 0x1A
	IDServerReference            byte = 0x1C
	IDReasonString               byte = 0x1F
	IDReceiveMaximum             byte = 0x21
	IDTopicAliasMaximum          byte = 0x22
	IDTopicAlias                 byte = 0x23
	IDMaximumQoS                 byte = 0x24
	IDRetainAvailable            byte = 0x25
	IDUserProperties             byte = 0x26
	IDMaximumPacketSize          byte = 0x27
	IDWildcardSubAvailable       byte = 0x28
	IDSubIDAvailable             byte = 0x29
	IDSharedSubAvailable         byte = 0x2A
)

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

func (p *Properties) Len() uint64 {
	// unimplemented
	return 0
}

func (p *Properties) Encode() ([]byte, error) {
	result := []byte{0}

	if p.SessionExpiryInterval != nil {
		result = append(result, IDSessionExpiryInterval)
		result = append(result, uint32ToBytes(*p.SessionExpiryInterval)...)
	}

	if p.RequestResponseInformation != nil {
		result = append(result, IDRequestResponseInformation, *p.RequestResponseInformation)
	}
	if p.RequestProblemInformation != nil {
		result = append(result, 0x17)
		result = append(result, *p.RequestProblemInformation)
	}
	if p.ReceiveMaximum != nil {
		result = append(result, 0x21)
		result = append(result, uint16ToBytes(*p.ReceiveMaximum)...)
	}
	if p.MaximumPacketSize != nil {
		result = append(result, 0x27)
		result = append(result, uint32ToBytes(*p.MaximumPacketSize)...)
	}
	if p.TopicAliasMax != nil {
		result = append(result, 0x22)
		result = append(result, uint16ToBytes(*p.TopicAliasMax)...)
	}
	if p.UserProperties != nil {
		for key, val := range *p.UserProperties {
			result = append(result, IDUserProperties)
			if buf, err := strToBytes(key); err != nil {
				return nil, err
			} else {
				result = append(result, buf...)
			}
			if buf, err := strToBytes(val); err != nil {
				return nil, err
			} else {
				result = append(result, buf...)
			}
		}
	}
	lengthBytes, err := encodeVariableByteInteger((len(result)))
	if err != nil {
		return nil, err
	}
	return append(lengthBytes, result...), nil
}

func (p *Properties) Decode(data []byte) (int, error) {
	total, n, err := decodeVariableByteInteger(data)
	if err != nil {
		return n, err
	}
	// data[n : n+l]
	for i := n; i < total; {
		identifier := data[i]
		i++
		switch identifier {
		case IDSessionExpiryInterval:
			if p.SessionExpiryInterval, err = decodeUint32Ptr(data[i : i+4]); err != nil {
				return 0, err
			}
			i = i + 4
		case IDRequestResponseInformation:
			p.RequestProblemInformation = new(byte)
			*p.RequestProblemInformation = data[i]
			i++
		default:
			return 0, ErrProtocolViolation
		}
	}
	return n + total, nil
}

func readProperties(reader io.Reader) (result *Properties, err error) {
	result = &Properties{}
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
				if result.SessionExpiryInterval, err = readUint32Ptr(reader); err != nil {
					return
				}
			}
		case 0x19:
			{
				i++
				result.RequestResponseInformation, err = readBytePtr(reader)
				if err != nil {
					return
				}
			}
		case 0x17:
			{
				i++
				result.RequestProblemInformation, err = readBytePtr(reader)
				if err != nil {
					return
				}
			}
			// receive max
		case 0x21:
			{
				i += 2
				if result.ReceiveMaximum, err = readUint16Ptr(reader); err != nil {
					return
				}
			}
			// Max packet size
		case 0x27:
			{
				i += 4
				if max, err := readUint32Ptr(reader); err != nil {
					return nil, err
				} else if max == nil {
					result.MaximumPacketSize = max
				}
			}
			//  Topic Alias Max
		case 0x22:
			{
				i += 2
				if result.TopicAliasMax, err = readUint16Ptr(reader); err != nil {
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
				userProperties := UserProperties{}

				for i := 0; i < len(list); i += 2 {
					userProperties[list[i]] = userProperties[list[i+1]]
				}
				result.UserProperties = &userProperties
				i += ul
				return
			}
		}
	}
	return
}

// unimplemented
func writeProperties(writer io.Writer, p *Properties) error {
	return nil
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
