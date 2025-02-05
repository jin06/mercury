package mqtt

import (
	"fmt"
)

// Read FixedHeader and remaining length
func ReadPacket(reader *Reader) (packet Packet, err error) {
	header := &FixedHeader{}
	if err := header.Read(reader); err != nil {
		return nil, err
	}
	switch header.PacketType {
	case CONNECT:
		packet = NewConnect(header)
		if err := packet.ReadBody(reader); err != nil {
			return nil, err
		}
	default:
		fmt.Println(header)
	}
	return
}

func WritePacket(writer *Writer, p Packet) error {
	return p.Write(writer)
}
