package mqtt

const (
	MIN_PACKET_ID PacketID = 0
	MAX_PACKET_ID PacketID = 65535
)

type PacketID uint16
