package mqtt

import (
	"bufio"
	"net"
)

func NewWriter(c net.Conn) *Writer {
	return &Writer{
		conn:   c,
		Writer: bufio.NewWriter(c),
	}
}

type Writer struct {
	conn net.Conn
	*bufio.Writer
}
