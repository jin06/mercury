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

func (w *Writer) Write(data []byte) (int, error) {
	return w.conn.Write(data)
}

func (w *Writer) WriteBool(b bool) (int, error) {
	return w.conn.Write([]byte{encodeBool(b)})
}

func (w *Writer) WriteUint8(u uint8) (int, error) {
	return w.conn.Write([]byte{byte(u)})
}

func (w *Writer) WriteUint16(u uint16) (int, error) {
	return w.conn.Write(encodeUint16(u))
}

func (w *Writer) WriteUint32(u uint32) (int, error) {
	return w.conn.Write(encodeUint32(u))
}

func (w *Writer) WriteUint64(u uint64) (int, error) {
	return w.conn.Write(encodeUint64(u))
}

func (w *Writer) WriteUTF8Str(str string) (int, error) {
	data, err := encodeUTF8Str(str)
	if err != nil {
		return 0, err
	}
	return w.conn.Write(data)
}

func (w *Writer) WriteVariableByteInteger(l int) (int, error) {
	data, err := encodeVariableByteInteger(l)
	if err != nil {
		return 0, err
	}
	return w.conn.Write(data)
}
