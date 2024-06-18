package mqtt

type Subscribe struct {
	PacketID      PacketID
	TopicWildcard TopicWildcard
	Payload       []Subscription
}

type Subscription struct {
	TopicWildcard TopicWildcard
	QoS           QoS
}
