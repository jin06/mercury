package mqtt

type Publish struct {
	Dup                    bool
	Qos                    QoS
	Retain                 bool
	Topic                  string
	PacketID               uint16
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
