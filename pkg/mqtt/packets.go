package mqtt

import "io"

type Message struct {
	packet      Packet
	fixedHeader FixedHeader
}

// Packet for mqtt packet
type Packet interface {
	// String() string
	Encode() ([]byte, error)
	Decode([]byte) error
	Write(io.Writer) error
	Read(io.Reader) error
}

// FixedHeader
// 3.1 link: https://public.dhe.ibm.com/software/dw/webservices/ws-mqtt/mqtt-v3r1.html#fixed-header
// 3.1.1 link: https://docs.oasis-open.org/mqtt/mqtt/v3.1.1/os/mqtt-v3.1.1-os.html#_Toc398718020
// 5.0 link: https://docs.oasis-open.org/mqtt/mqtt/v5.0/os/mqtt-v5.0-os.html#_Toc3901021
type FixedHeader struct {
	PacketType PacketType
	Reserved   byte
	Remaining  []byte
}

type QoS byte

const (
	QoS0 QoS = iota
	QoS1
	QoS2
)

// type FiexedHeader struct {
// 	Raw [2]byte
// 	// Type  byte
// 	// Flags byte
// }

type VariableHeader struct{}

type Payload struct{}
