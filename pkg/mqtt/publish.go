package mqtt

type Publish struct {
	PacketID               PacketID
	Dup                    bool
	Qos                    QoS
	Retain                 bool
	Topic                  string
	Payload                []byte
	PayloadFormatIndicator bool   // from mqtt5
	MessageExpiryInterval  uint16 // from mqtt5 (seconds)
	TopicAlias             uint16 // from mqtt5
	ResponseTopic          string // from mqtt5
	CorrelationData        []byte // from mqtt5
	Properties             *Properties
	SubscriptionIdentifier uint32
	ContentType            string
	Version                ProtocolVersion
}

func (p *Publish) Encode() ([]byte, error) {
	// write header
	result := toHeader(PUBLISH)
	if p.Dup {
		result[0] |= 0b00001000
	}
	result[0] |= (byte(p.Qos) << 1)
	// write topic name
	if bytes, err := strToBytes(p.Topic); err != nil {
		return nil, err
	} else {
		result = append(result, bytes...)
	}
	// write message ID
	result = append(result, packetIDToBytes(p.PacketID)...)

	if p.Version == MQTT5 {
		bytes, err := encodeProperties(p.Properties)
		if err != nil {
			return nil, err
		}
		result = append(result, bytes...)
	}
	result = append(result, p.Payload...)

	return result, nil
}
