package mqtt

type Unsubscribe struct {
	QoS       QoS
	Dup       bool
	MessageID uint16
	Payload   []TopicWildcard
}
