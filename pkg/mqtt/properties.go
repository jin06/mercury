package mqtt

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
