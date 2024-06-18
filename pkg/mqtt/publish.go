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
	UserProperties         UserProperties
	SubscriptionIdentifier uint32
	ContentType            string
}

func (p *Publish) Encode([]byte, error) {
	result := toHeader(PUBLISH)
	if p.Dup {
		result[0] |= 0b00001000
	}
	result[0] |= (byte(p.Qos) << 1)

	if p.Version = MQTT5 {

	}
}
