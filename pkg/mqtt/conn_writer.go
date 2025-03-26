package mqtt

import (
	"bufio"
	"io"
)

func newWriter(c io.Writer) *Writer {
	return &Writer{
		raw:    c,
		Writer: bufio.NewWriter(c),
	}
}

type Writer struct {
	raw io.Writer
	*bufio.Writer
}

func (w *Writer) WritePacket(p Packet) error {
	data, err := p.Encode()
	if err != nil {
		return err
	}
	// _, err = w.Writer.Write(data)
	//todo: buffer io writer will not write data to the underlying writer until the buffer is full
	_, err = w.raw.Write(data)
	return err
}

func (w *Writer) Write(data []byte) (int, error) {
	return w.Writer.Write(data)
}

func (w *Writer) WriteBool(b bool) (int, error) {
	return w.Writer.Write([]byte{encodeBool(b)})
}

func (w *Writer) WriteUint8(u uint8) (int, error) {
	return w.Writer.Write([]byte{byte(u)})
}

func (w *Writer) WriteUint16(u uint16) (int, error) {
	return w.Writer.Write(encodeUint16(u))
}

func (w *Writer) WriteUint32(u uint32) (int, error) {
	return w.Writer.Write(encodeUint32(u))
}

func (w *Writer) WriteUint64(u uint64) (int, error) {
	return w.Writer.Write(encodeUint64(u))
}

func (w *Writer) WriteUTF8Str(str string) (int, error) {
	data, err := encodeUTF8Str(str)
	if err != nil {
		return 0, err
	}
	return w.Writer.Write(data)
}

func (w *Writer) WriteVariableByteInteger(l VariableByteInteger) (int, error) {
	data, err := encodeVariableByteInteger(l)
	if err != nil {
		return 0, err
	}
	return w.Writer.Write(data)
}
