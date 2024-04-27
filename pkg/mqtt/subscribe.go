package mqtt

type Subscribe struct {
	MessageID     uint16
	TopicWildcard TopicWildcard
	Payload       []Subscription
}

type Subscription struct {
	TopicWildcard TopicWildcard
	QoS           QoS
}
