package connections

import (
	"net"

	"github.com/jin06/mercury/pkg/mqtt"
)

func NewTCP(c net.Conn) Connection {
	return &TCP{
		conn:   c,
		reader: mqtt.NewReader(c),
		writer: mqtt.NewWriter(c),
	}
}

type TCP struct {
	conn   net.Conn
	reader *mqtt.Reader
	writer *mqtt.Writer
}

func (t *TCP) Read() (mqtt.Packet, error) {
	return mqtt.ReadPacket(t.reader)
}

func (t *TCP) Write(p mqtt.Packet) error {
	return mqtt.WritePacket(t.writer, p)
}

func (t *TCP) Close() error {
	if t.conn == nil {
		t.conn.Close()
	}
	return nil
}
