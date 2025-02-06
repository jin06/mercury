package mqtt

const (
	QoS0 QoS = iota
	QoS1
	QoS2
)

type QoS byte

type VariableHeader struct{}

type Payload struct{}

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

// FixedHeader
// 3.1 link: https://public.dhe.ibm.com/software/dw/webservices/ws-mqtt/mqtt-v3r1.html#fixed-header
// 3.1.1 link: https://docs.oasis-open.org/mqtt/mqtt/v3.1.1/os/mqtt-v3.1.1-os.html#_Toc398718020
// 5.0 link: https://docs.oasis-open.org/mqtt/mqtt/v5.0/os/mqtt-v5.0-os.html#_Toc3901021
type FixedHeader struct {
	PacketType      PacketType
	Flags           byte
	RemainingLength int
}

func (f *FixedHeader) Encode() ([]byte, error) {
	panic("todo")
}

// Decode decodes a fixed header from the given byte slice ([]byte).
// It extracts the packet type, flags, and remaining length (decoded using decodeVariableByteInteger).
// The function returns the number of bytes consumed and any error encountered during decoding.
//
// Parameters:
//   - data []byte: The byte slice containing the encoded fixed header.
//
// Returns:
//   - int: The number of bytes consumed during decoding.
//   - error: Any error encountered (e.g., byte slice too short or decoding issue).
func (f *FixedHeader) Decode(data []byte) (int, error) {
	if len(data) < 2 {
		return 0, ErrBytesShorter
	}
	f.PacketType = PacketType(data[0] >> 4)
	f.Flags = 0b0001111 & data[0]
	length, n, err := decodeVariableByteInteger(data[1:])
	if err != nil {
		return 0, err
	}
	total := 1 + n + length
	if len(data) != total {
		return 0, ErrPacketEncoding
	}
	f.RemainingLength = length
	return n + 1, nil
}

func (f *FixedHeader) Read(reader *Reader) error {
	b, err := reader.ReadByte()
	if err != nil {
		return err
	}
	f.PacketType = PacketType(b >> 4)
	f.Flags = 0b00001111 & b
	l, _, err := readVariableByteInteger(reader)
	if err != nil {
		return err
	}
	f.RemainingLength = l
	return nil
}

func (f *FixedHeader) Write(writer *Writer) error {
	b := byte(f.PacketType<<4) | (f.Flags)
	if err := writer.WriteByte(b); err != nil {
		return err
	}
	return writeVariableByteInteger(writer, f.RemainingLength)
}
