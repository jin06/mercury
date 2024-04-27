package mqtt

type Pubrecl struct {
	QoS       QoS
	Dup       bool
	MessageID uint16
}
