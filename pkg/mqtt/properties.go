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
	ID_RequestProblemInfomation        byte = 0x17
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

// Properties represents the various properties in an MQTT packet.
// These are used for controlling the behavior of an MQTT connection,
// message, or subscription and include options like message expiry,
// authentication, and session details.
type Properties struct {
	// PayloadFormat specifies the format of the payload.
	// If this is not provided, the default format is assumed.
	// A value of 0 indicates "UTF-8", and 1 indicates "binary".
	PayloadFormat *byte

	// MessageExpiryInterval specifies the expiry interval of the message.
	// This is the maximum time (in seconds) that the message is valid for.
	// After this interval, the message is discarded by the broker.
	MessageExpiryInterval *uint32

	// ContentType is used to specify the content type of the payload.
	// This helps identify the type of data in the payload (e.g., JSON, XML, etc.).
	ContentType *string

	// ResponseTopic is the topic to which the response should be sent.
	// This can be used in request-response messaging patterns.
	ResponseTopic *string

	// CorrelationData is used to correlate a response message with a request.
	// It is typically included in the response to match a previously sent request.
	CorrelationData []byte

	// SubscriptionIdentifier is an identifier for the subscription.
	// It can be used to relate a subscription to a specific client or purpose.
	SubscriptionIdentifier *int

	// SessionExpiryInterval specifies the session expiry time in seconds.
	// This defines how long the broker should keep the session alive after the client disconnects.
	// A value of 0 indicates that the session is valid indefinitely.
	SessionExpiryInterval *uint32

	// AssignedClientID is the ID assigned to a client when it connects.
	// This is the unique identifier used by the broker to identify a client.
	AssignedClientID *string

	// ServerKeepAlive is the maximum time interval between two consecutive MQTT messages
	// (e.g., PINGREQ and PINGRESP). If no messages are sent in this time, the server will disconnect the client.
	ServerKeepAlive *uint16

	// AuthenticationMethod is used to specify the authentication method
	// to be used by the client for connection. This is an optional field.
	AuthenticationMethod *string

	// AuthenticationData contains the data required for the client’s authentication method.
	// This can be a password, certificate, etc., depending on the authentication method.
	AuthenticationData []byte

	// RequestProblemInformation specifies whether the client wants the broker to include problem information in responses.
	// This can be used to indicate issues with message processing.
	RequestProblemInformation *bool

	// WillDelayInterval specifies the delay in seconds before the last will message is sent after disconnection.
	// This can be used to delay the sending of the will message to allow for certain conditions to be met.
	WillDelayInterval *uint32

	// RequestResponseInformation indicates whether the client wants the broker to provide response information
	// (i.e., if the message expects a response).
	RequestResponseInformation *bool

	// ResponseInformation is the information provided by the broker in response to a message,
	// typically used in request-response scenarios.
	ResponseInformation *string

	// ServerReference provides a reference to the server that can be used in the response.
	// This is often used for identifying which server the message was processed by.
	ServerReference *string

	// ReasonString is a human-readable string that provides additional information
	// about the reason for a particular message or state.
	ReasonString *string

	// ReceiveMaximum indicates the maximum number of QoS 1 and QoS 2 publications
	// the client is willing to process concurrently.
	ReceiveMaximum *uint16

	// TopicAliasMaximum specifies the maximum number of topic aliases the client is willing to accept.
	// Topic aliases are used to reduce the size of the topic in messages.
	TopicAliasMaximum *uint16

	// TopicAlias is a topic alias used to reduce the size of the topic in messages.
	// It is an identifier for a specific topic.
	TopicAlias *uint16

	// MaximumQoS specifies the maximum QoS level supported by the client for message delivery.
	// This can be 0 (At most once), 1 (At least once), or 2 (Exactly once).
	MaximumQoS *QoS

	// RetainAvailable indicates whether the retain flag is supported by the server.
	// If true, the client can request retained messages for topics.
	RetainAvailable *bool

	// UserProperties allow the client and broker to exchange arbitrary key-value pairs.
	// This can be used for custom metadata, such as application-specific properties.
	UserProperties *UserProperties

	// MaximumPacketSize specifies the maximum size of an MQTT control packet in bytes.
	// This is used to ensure the control packets do not exceed the specified size limit.
	MaximumPacketSize *uint32

	// WildcardSubscriptionAvailable indicates whether the broker supports wildcard subscriptions.
	// If true, the client can subscribe to topics using wildcard characters like '#' or '+'.
	WildcardSubscriptionAvailable *bool

	// SubscriptionIdentifierAvailable indicates whether the broker supports subscription identifiers.
	// These identifiers are used to track and manage subscriptions.
	SubscriptionIdentifierAvailable *bool

	// SharedSubscriptionAvailable indicates whether the broker supports shared subscriptions.
	// Shared subscriptions allow multiple clients to share a single subscription to a topic.
	SharedSubscriptionAvailable *bool
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
		result = append(result, ID_RequestResponseInformation, encodeBool(*p.RequestResponseInformation))
	}
	if p.RequestProblemInformation != nil {
		result = append(result, 0x17)
		result = append(result, encodeBool(*p.RequestProblemInformation))
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
		var vl int
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
		case ID_ContentType:
			if p.ContentType, vl, err = decodeUTF8Ptr(data[i:]); err != nil {
				return i + total, err
			}
			i = i + vl
		case ID_ResponseTopic:
			if p.ResponseTopic, vl, err = decodeUTF8Ptr(data[i:]); err != nil {
				return i + total, err
			}
			i = i + vl
		case ID_CorrelationData:
			if p.CorrelationData, vl, err = decodeBinaryData(data[i:]); err != nil {
				return i + total, err
			}
			i = i + vl
		case ID_SubscriptionIdentifier:
			if p.SubscriptionIdentifier, vl, err = decodeVariableByteIntegerPtr(data[i:]); err != nil {
				return i + total, err
			}
			i = i + vl
		case ID_SessionExpiryInterval:
			if p.SessionExpiryInterval, err = decodeUint32Ptr(data[i : i+4]); err != nil {
				return i + total, err
			}
			i = i + 4
		case ID_AssignedClientID:
			if p.AssignedClientID, vl, err = decodeUTF8Ptr(data[i:]); err != nil {
				return i + total, err
			}
			i += vl
		case ID_ServerKeepAlive:
			if p.ServerKeepAlive, err = decodeUint16Ptr(data[i : i+2]); err != nil {
				return i + total, err
			}
			i += 2
		case ID_AuthenticationMethod:
			if p.AuthenticationMethod, vl, err = decodeUTF8Ptr(data[i:]); err != nil {
				return i + total, err
			}
			i += vl
		case ID_AuthenticationData:
			if p.AuthenticationData, vl, err = decodeBinaryData(data[i:]); err != nil {
				return i + total, err
			}
			i += vl
		case ID_RequestProblemInfomation:
			if p.RequestProblemInformation, err = decodeBoolPtr(data[i]); err != nil {
				return i + total, err
			}
			i += vl
		case ID_WillDelayInterval:
			if p.WillDelayInterval, err = decodeUint32Ptr(data[i:]); err != nil {
				return i + total, err
			}
			i += 4
		case ID_RequestResponseInformation:
			if p.RequestResponseInformation, err = decodeBoolPtr(data[i]); err != nil {
				return i + total, err
			}
			i++
		case ID_ResponseInformation:
			if p.ResponseInformation, vl, err = decodeUTF8Ptr(data[i:]); err != nil {
				return i + total, err
			}
			i = i + vl
		case ID_ServerReference:
			if p.ServerReference, vl, err = decodeUTF8Ptr(data[i:]); err != nil {
				return i + total, err
			}
			i += vl
		case ID_ReasonString:
			if p.ReasonString, vl, err = decodeUTF8Ptr(data[i:]); err != nil {
				return i + total, err
			}
			i += vl
		case ID_ReceiveMaximum:
			if p.ReceiveMaximum, err = decodeUint16Ptr(data[i : i+2]); err != nil {
				return i + total, err
			}
			i = i + 2
		case ID_TopicAliasMaximum:
			if p.TopicAliasMaximum, err = decodeUint16Ptr(data[i : i+2]); err != nil {
				return i + total, err
			}
			i = i + 2
		case ID_TopicAlias:
			if p.TopicAlias, err = decodeUint16Ptr(data[i : i+2]); err != nil {
				return i + total, err
			}
			i += 2
		case ID_MaximumQoS:
			p.MaximumQoS = decodeBytePrt(data[i])
			i++
		case ID_RetainAvailable:
			if p.RetainAvailable, err = decodeBoolPtr(data[i]); err != nil {
				return i + total, err
			}
			i++
		case ID_UserProperties:
			panic("not imeplement")
		case ID_MaximumPacketSize:
			if p.MaximumPacketSize, err = decodeUint32Ptr(data[i : i+4]); err != nil {
				return i + total, err
			}
			i = i + 4
		case ID_WildcardSubscriptionAvailable:
			if p.WildcardSubscriptionAvailable, err = decodeBoolPtr(data[i]); err != nil {
				return i + total, err
			}
			i++
		case ID_SubscriptionIdentifierAvailable:
			if p.SubscriptionIdentifierAvailable, err = decodeBoolPtr(data[i]); err != nil {
				return i + total, err
			}
			i++
		case ID_SharedSubscriptionAvailable:
			if p.SharedSubscriptionAvailable, err = decodeBoolPtr(data[i]); err != nil {
				return i + total, err
			}
			i++
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
				if p.RequestResponseInformation, err = r.ReadBoolPtr(); err != nil {
					return err
				}
			}
		case 0x17:
			{
				if p.RequestProblemInformation, err = r.ReadBoolPtr(); err != nil {
					return err
				}
				i++
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

func (u *UserProperties) Encode() (result []byte, err error) {
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
