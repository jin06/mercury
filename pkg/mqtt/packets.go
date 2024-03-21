package mqtt

import "io"

// Packet for mqtt packet
type Packet interface {
	String() string
	Decode(reader io.Reader) error
}

// FixedHeader
// 3.1 link: https://public.dhe.ibm.com/software/dw/webservices/ws-mqtt/mqtt-v3r1.html#fixed-header
// 3.1.1 link: https://docs.oasis-open.org/mqtt/mqtt/v3.1.1/os/mqtt-v3.1.1-os.html#_Toc398718020
// 5.0 link: https://docs.oasis-open.org/mqtt/mqtt/v5.0/os/mqtt-v5.0-os.html#_Toc3901021
// type FixedHeader struct {
// 	PacketType byte
// 	Reserved   byte
// 	Remaining  byte
// }

type PacketID uint16

type QoS byte

const (
	QoS0 QoS = iota
	QoS1
	QoS2
)

type PacketType byte

const (
	CONNECT PacketType = iota + 1
	CONNACK
	PUBLISH
	PUBACK
	PUBREC
	PUBREL
	PUBCOMP
	SUBSCRIBE
	SUBACK
	UNSUBSCRIBE
	UNSUBACK
	PINGREQ
	PINGRESP
	DISCONNECT
	AUTH
)

type FiexedHeader struct {
	Raw [2]byte
	// Type  byte
	// Flags byte
}

type VariableHeader struct{}

type Payload struct{}
