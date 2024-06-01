package mqtt

type Suback struct {
	PacketID   uint16
	GrantedQoS QoS
}
