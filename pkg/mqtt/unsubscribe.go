package mqtt

type Unsubscribe struct {
	QoS      QoS
	Dup      bool
	PacketID PacketID
	Payload  []TopicWildcard
}
