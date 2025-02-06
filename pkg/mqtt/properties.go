package mqtt

import (
	"errors"
	"io"
)

const (
	ID_PayloadFormat                   byte = 0x01
	ID_MessageExpiryInterval           byte = 0x02
	ID_ContentType                     byte = 0x03
	ID_ResponseTopic                   byte = 0x08
	ID_CorrelationData                 byte = 0x09
	ID_SubscriptionIdentifier          byte = 0x0B
	ID_SessionExpiryInterval           byte = 0x11
	ID_AssignedClientID                byte = 0x12
	ID_ServerKeepAlive                 byte = 0x13
	ID_AuthenticationMethod            byte = 0x15
	ID_AuthenticationData              byte = 0x16
	ID_RequestProblemInfo              byte = 0x17
	ID_WillDelayInterval               byte = 0x18
	ID_RequestResponseInformation      byte = 0x19
	ID_ResponseInformation             byte = 0x1A
	ID_ServerReference                 byte = 0x1C
	ID_ReasonString                    byte = 0x1F
	ID_ReceiveMaximum                  byte = 0x21
	ID_TopicAliasMaximum               byte = 0x22
	ID_TopicAlias                      byte = 0x23
	ID_MaximumQoS                      byte = 0x24
	ID_RetainAvailable                 byte = 0x25
	ID_UserProperties                  byte = 0x26
	ID_MaximumPacketSize               byte = 0x27
	ID_WildcardSubscriptionAvailable   byte = 0x28
	ID_SubscriptionIdentifierAvailable byte = 0x29
	ID_SharedSubscriptionAvailable     byte = 0x2A
)

// mqtt5
type Properties struct {
	PayloadFormat          *byte
	MessageExpiryInterval  *uint32
	ContentType            *string
	ResponseTopic          *string
	CorrelationData        []byte
	SubscriptionIdentifier []byte
	// SessionExpiryInterval second
	SessionExpiryInterval      *uint32
	AssignedClientID           *string
	ServerKeepAlive            *uint16
	AuthenticationMethod       *string
	AuthenticationData         *string
	RequestProblemInformation  *byte
	WillDelayInterval          *uint32
	RequestResponseInformation *byte
	ResponseInformation        *string
	ServerReference            *string
	ReasonString               *string
	//ReceiveMaximum The Client uses this value to limit the number of QoS 1 and QoS 2 publications that it is willing to process concurrently.
	ReceiveMaximum    *uint16
	TopicAliasMaximum *uint16
	TopicAlias        *uint16
	MaximumQoS        *byte
	RetainAvailable   *byte
	UserProperties    *UserProperties
	// MaximumPacketSize The packet size is the total number of bytes in an MQTT Control Packet
	MaximumPacketSize               *uint32
	WildcardSubscriptionAvailable   *byte
	SubscriptionIdentifierAvailable *byte
	SharedSubscriptionAvailable     *byte
}

func (p *Properties) Len() uint64 {
	// unimplemented
	return 0
}

func (p *Properties) Encode() ([]byte, error) {
	result := []byte{0}

	if p.SessionExpiryInterval != nil {
		result = append(result, ID_SessionExpiryInterval)
		result = append(result, encodeUint32(*p.SessionExpiryInterval)...)
	}

	if p.RequestResponseInformation != nil {
		result = append(result, ID_RequestResponseInformation, *p.RequestResponseInformation)
	}
	if p.RequestProblemInformation != nil {
		result = append(result, 0x17)
		result = append(result, *p.RequestProblemInformation)
	}
	if p.ReceiveMaximum != nil {
		result = append(result, 0x21)
		result = append(result, encodeUint16(*p.ReceiveMaximum)...)
	}
	if p.MaximumPacketSize != nil {
		result = append(result, 0x27)
		result = append(result, encodeUint32(*p.MaximumPacketSize)...)
	}
	if p.TopicAliasMaximum != nil {
		result = append(result, 0x22)
		result = append(result, encodeUint16(*p.TopicAliasMaximum)...)
	}
	// todo userPropertis decode and encode
	if p.UserProperties != nil {
		for key, val := range *p.UserProperties {
			result = append(result, ID_UserProperties)
			if buf, err := encodeUTF8Str(key); err != nil {
				return nil, err
			} else {
				result = append(result, buf...)
			}
			if buf, err := encodeUTF8Str(val); err != nil {
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
		case ID_PayloadFormat:
			p.PayloadFormat = new(byte)
			*p.PayloadFormat = data[i]
			i++
		case ID_MessageExpiryInterval:
			if p.MessageExpiryInterval, err = decodeUint32Ptr(data[i : i+4]); err != nil {
				return i + total, err
			}
			i = i + 4
		case ID_SessionExpiryInterval:
			if p.SessionExpiryInterval, err = decodeUint32Ptr(data[i : i+4]); err != nil {
				return i + total, err
			}
			i = i + 4
		case ID_ContentType:
			if p.ContentType, n, err = decodeUTF8Ptr(data[i:]); err != nil {
				return i + total, err
			} else {
				i = i + n
			}
		case ID_ResponseTopic:
			if p.ResponseTopic, n, err = decodeUTF8Ptr(data[i:]); err != nil {
				return i + total, err
			} else {
				i = i + n
			}
		case ID_CorrelationData:
			if p.CorrelationData, n, err = decodeBinaryData(data[i:]); err != nil {
				return i + total, err
			} else {
				i = i + n
			}
		case ID_ReceiveMaximum:
			if p.ReceiveMaximum, err = decodeUint16Ptr(data[i : i+2]); err != nil {
				return i + total, err
			}
			i = i + 2
		case ID_MaximumPacketSize:
			if p.MaximumPacketSize, err = decodeUint32Ptr(data[i : i+4]); err != nil {
				return i + total, err
			}
			i = i + 4
		case ID_TopicAliasMaximum:
			if p.TopicAliasMaximum, err = decodeUint16Ptr(data[i : i+2]); err != nil {
				return i + total, err
			}
			i = i + 2
		case ID_RequestResponseInformation:
		default:
			return 0, ErrProtocolViolation
		}
	}
	return n + total, nil
}

func (p *Properties) Read(r *Reader) error {
	// result = &Properties{}
	var total int

	total, _, err := r.ReadVariableByteInteger()
	if err != nil {
		return err
	}

	for i := 0; i < total; {
		var identifier byte
		if identifier, err = r.ReadByte(); err != nil {
			return err
		}
		i++

		switch identifier {
		case 0x11:
			{
				i += 4
				if p.SessionExpiryInterval, err = r.ReadUint32Ptr(); err != nil {
					return err
				}
			}
		case 0x19:
			{
				i++
				if p.RequestResponseInformation, err = r.ReadBytePtr(); err != nil {
					return err
				}
			}
		case 0x17:
			{
				i++
				if p.RequestProblemInformation, err = r.ReadBytePtr(); err != nil {
					return err
				}
			}
			// receive max
		case 0x21:
			{
				i += 2
				if p.ReceiveMaximum, err = r.ReadUint16Ptr(); err != nil {
					return err
				}
			}
			// Max packet size
		case 0x27:
			{
				i += 4
				if max, err := r.ReadUint32Ptr(); err != nil {
					return err
				} else if max == nil {
					p.MaximumPacketSize = max
				}
			}
			//  Topic Alias Max
		case 0x22:
			{
				i += 2
				if p.TopicAliasMaximum, err = r.ReadUint16Ptr(); err != nil {
					return err
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
					if val, n, err = r.ReadUTF8Str(); err != nil {
						return err
					}
					ul += n
					j = j + n
					list = append(list, val)
				}
				if len(list)%2 == 1 {
					return ErrProtocol
				}
				userProperties := UserProperties{}

				for i := 0; i < len(list); i += 2 {
					userProperties[list[i]] = userProperties[list[i+1]]
				}
				p.UserProperties = &userProperties
				i += ul
				return nil
			}
		}
	}
	return nil
}

// unimplemented
func writeProperties(writer io.Writer, p *Properties) error {
	return nil
}

type UserProperties map[string]string

func (u *UserProperties) toBytes() (result []byte, err error) {
	result = []byte{}
	for key, val := range *u {
		if bytes, err := encodeUTF8Str(key); err != nil {
			return result, err
		} else {
			result = append(result, bytes...)
		}
		if bytes, err := encodeUTF8Str(val); err != nil {
			return result, err
		} else {
			result = append(result, bytes...)
		}
	}
	return
}

func (u *UserProperties) Read(r *Reader) error {
	propertyLength, err := r.ReadUint8()
	if err != nil {
		return err
	}
	// buf, err := readBytes(reader, int(propertyLength))
	arr := []string{}
	for i := 0; i <= int(propertyLength); {
		str, n, err := r.ReadUTF8Str()
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
