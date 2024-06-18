package mqtt

type Pubrecl struct {
	PacketID PacketID
	QoS      QoS
	Dup      bool
}
