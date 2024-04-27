package mqtt

type Publish struct {
	Dup       bool
	Qos       QoS
	Retain    bool
	Topic     string
	MessageID uint16
	Payload   []byte
}
