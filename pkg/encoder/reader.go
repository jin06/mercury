package encoder

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/jin06/mercury/pkg/mqtt"
)

func NewReader(conn net.Conn) *Reader {
	r := &Reader{
		conn: conn,
	}
	r.readbuf = bufio.NewReader(conn)
	return r
}

type Reader struct {
	conn    net.Conn
	readbuf *bufio.Reader
}

func packetType(b byte) (t mqtt.PacketType, err error) {
	t = mqtt.PacketType(b >> 4)
	return
}

func parsePacket(fiexedBytes []byte) (mqtt.Packet, error) {
	if len(fiexedBytes) != 2 {
		return nil, errors.New(fmt.Sprintf("fixedBytes length expected 2, found %d", len(fiexedBytes)))
	}
	t := mqtt.PacketType(fiexedBytes[0] & 0b11110000)

	t = t >> 4
	switch t {
	case mqtt.CONNECT:
		{
			return &mqtt.Connect{}, nil
		}
	case mqtt.CONNACK:
		{
			return &mqtt.Connect{}, nil

		}
	}
	return nil, errors.New("unsupported type")
}

func (r *Reader) ReadPacket() (p mqtt.Packet, err error) {
	fh := make([]byte, 2)

	if _, err = io.ReadFull(r.readbuf, fh); err != nil {
		return
	}
	if p, err = parsePacket(fh); err != nil {
		return
	}
	p.Decode(r.readbuf)
	return
}

type Writer struct {
	net.Conn
}
