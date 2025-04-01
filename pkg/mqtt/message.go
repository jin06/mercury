package mqtt

import "time"

type Message interface {
	Packet
	ID() PacketID
	Time() time.Time
	Expiry() time.Duration
}
