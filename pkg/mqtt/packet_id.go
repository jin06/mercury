package mqtt

const (
	MIN_PACKET_ID PacketID = 0
	MAX_PACKET_ID PacketID = 65535
)

type PacketID uint16

// ToBytes return MSB and LSB
// byte 1 is MSB
// byte 2 is LSB
func (p PacketID) ToBytes() []byte {
	return []byte{byte(p >> 8), byte(p)}
}
