package mqtt

import (
	"bufio"
	"fmt"
	"net"
)

func NewReader(c net.Conn) *Reader {
	return &Reader{
		conn:   c,
		Reader: bufio.NewReader(c),
	}
}

type Reader struct {
	conn net.Conn
	*bufio.Reader
}

type Writer struct {
	conn net.Conn
	*bufio.Writer
}

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
