package mqtt

const (
	QoS0 QoS = iota
	QoS1
	QoS2
)

type QoS byte

func (q QoS) Zero() bool {
	return q == QoS0
}

// Packet for mqtt packet
type Packet interface {
	// Encode all packet with header
	Encode() ([]byte, error)
	// Encode body without header, encode body(payload)
	EncodeBody() ([]byte, error)
	// Decode with header to bytes
	Decode([]byte) (int, error)
	// Decode body only
	DecodeBody([]byte) (int, error)
	// Read all with header
	Read(*Reader) error
	// Only body
	ReadBody(*Reader) error
	// Write all to writer with header
	Write(*Writer) error
	// Only body
	WriteBody(*Writer) error
}
