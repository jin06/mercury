package mqtt

type Message interface {
	Packet
	PID() PacketID
}
