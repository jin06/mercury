package mqtt

import "io"

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

func encodeProperties(props *Properties) ([]byte, error) {
	if props == nil {
		return []byte{0}, nil
	}
	result := []byte{0}

	if props.SessionExpiryInterval != nil {
		result = append(result, IDSessionExpiryInterval)
		result = append(result, uint32ToBytes(*props.SessionExpiryInterval)...)
	}

	if props.RequestResponseInformation != nil {
		result = append(result, IDRequestResponseInformation, *props.RequestResponseInformation)
	}
	if props.RequestProblemInformation != nil {
		result = append(result, 0x17)
		result = append(result, *props.RequestProblemInformation)
	}
	if props.ReceiveMaximum != nil {
		result = append(result, 0x21)
		result = append(result, uint16ToBytes(*props.ReceiveMaximum)...)
	}
	if props.MaximumPacketSize != nil {
		result = append(result, 0x27)
		result = append(result, uint32ToBytes(*props.MaximumPacketSize)...)
	}
	if props.TopicAliasMax != nil {
		result = append(result, 0x22)
		result = append(result, uint16ToBytes(*props.TopicAliasMax)...)
	}
	if props.UserProperties != nil {
		for key, val := range *props.UserProperties {
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
	lengthBytes, err := variableByteInteger(len(result))
	if err != nil {
		return nil, err
	}
	return append(lengthBytes, result...), nil
}

func decodeProperties(data []byte) (*Properties, error) {
	prop := &Properties{}
	return prop, nil
}

func readProperties(reader io.Reader) (result *Properties, err error) {
	result = &Properties{
		UserProperties:    make(UserProperties),
		MaximumPacketSize: -1,
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
				var max uint32
				if max, err = readUint32(reader); err != nil {
					return
				} else if max == 0 {
					err = ErrMaximumPacketSize
					return
				}
				result.MaximumPacketSize = int64(max)
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
