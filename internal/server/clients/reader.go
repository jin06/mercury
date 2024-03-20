package clients

import (
	"bufio"
	"errors"
	"io"
	"net"

	"github.com/jin06/mercury/pkg/mqtt"
)

func NewReader(conn net.Conn) *Reader {
	r := &Reader{
		conn: conn,
	}
	r.rbuf = bufio.NewReader(conn)
	return r
}

type Reader struct {
	conn net.Conn
	rbuf *bufio.Reader
}

func packetType(b byte) (t string, err error) {
	switch b {
	case b & 0b00010000:
		{
			t = "connect"
			return
		}
	case b & 0b00100000:
		{
			t = "connack"
			return
		}
	default:
		{
			err = errors.New("unknown packet type")
			return
		}
	}
	return
}

func (r *Reader) ReadPacket() (p mqtt.Packet, err error) {
	fh := make([]byte, 2)
	_, err = io.ReadFull(r.rbuf, fh)
	if err != nil {
		return
	}
	t, err := packetType(fh[0])
	panic(t)

	return
}

type Writer struct {
	net.Conn
}
