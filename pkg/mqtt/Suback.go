package mqtt

type Suback struct {
	MessageID  uint16
	GrantedQoS QoS
}
