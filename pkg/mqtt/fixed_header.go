package mqtt

// FixedHeader
// 3.1 link: https://public.dhe.ibm.com/software/dw/webservices/ws-mqtt/mqtt-v3r1.html#fixed-header
// 3.1.1 link: https://docs.oasis-open.org/mqtt/mqtt/v3.1.1/os/mqtt-v3.1.1-os.html#_Toc398718020
// 5.0 link: https://docs.oasis-open.org/mqtt/mqtt/v5.0/os/mqtt-v5.0-os.html#_Toc3901021
type FixedHeader struct {
	PacketType      PacketType
	Flags           byte
	RemainingLength VariableByteInteger
}

func (f *FixedHeader) Length() int {
	return f.RemainingLength.Int()
}

func (f *FixedHeader) Encode() ([]byte, error) {
	var data []byte
	// data = append(data, byte(f.PacketType<<4))
	// data = append(data, (f.Flags & 0b0001111))
	bit1 := byte(f.PacketType<<4) | (f.Flags & 0b1111)
	data = append(data, bit1)
	if length, err := encodeVariableByteInteger(f.RemainingLength); err != nil {
		return nil, err
	} else {
		data = append(data, length...)
	}
	return data, nil
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
	total := 1 + n + int(length)
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

func newBasePacket(t PacketType, v ProtocolVersion) *BasePacket {
	return &BasePacket{
		FixedHeader: &FixedHeader{PacketType: t},
		Version:     v,
	}
}

type BasePacket struct {
	*FixedHeader
	Version ProtocolVersion
}
