package mqtt

type Puback struct {
	PacketID       uint16
	ReasonCode     ReasonCode
	ReasonString   string
	UserProperties UserProperties
}
