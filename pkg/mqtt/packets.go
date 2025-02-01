package mqtt

import (
	"io"
)

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
	Decode([]byte) error
	// Decode body only
	DecodeBody([]byte)
	// Read all with header
	Read(io.Reader) error
	// Only body
	ReadBody(io.Reader) error
	// Write all to writer with header
	Write(io.Writer) error
	// Only body
	WriteBody(io.Writer) error
	// Remaining length
	Len() uint64
}

// FixedHeader
// 3.1 link: https://public.dhe.ibm.com/software/dw/webservices/ws-mqtt/mqtt-v3r1.html#fixed-header
// 3.1.1 link: https://docs.oasis-open.org/mqtt/mqtt/v3.1.1/os/mqtt-v3.1.1-os.html#_Toc398718020
// 5.0 link: https://docs.oasis-open.org/mqtt/mqtt/v5.0/os/mqtt-v5.0-os.html#_Toc3901021
type FixedHeader struct {
	PacketType      PacketType
	Flags           byte
	RemainingLength uint64
}

func (f *FixedHeader) Encode() ([]byte, error) {
	panic("todo")
}

func (f *FixedHeader) Decode([]byte) error {
	panic("todo")
}

func (f *FixedHeader) Read(reader *Reader) error {
	b, err := reader.ReadByte()
	if err != nil {
		return err
	}
	f.PacketType = PacketType(b >> 4)
	f.Flags = 0b00001111 & b
	l, err := readVariableByteInteger(reader)
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
