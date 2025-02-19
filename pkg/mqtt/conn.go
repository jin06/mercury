package mqtt

import "io"

type Connection struct {
	conn io.ReadWriteCloser
	*Reader
	*Writer
}

func NewConnection(c io.ReadWriteCloser) *Connection {
	return &Connection{
		conn:   c,
		Reader: newReader(c),
		Writer: newWriter(c),
	}
}

func (c *Connection) Close() error {
	return c.conn.Close()
}
