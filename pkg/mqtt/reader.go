package mqtt

import (
	"bufio"
	"io"
	"net"

	"github.com/jin06/mercury/pkg/utils"
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

func (r *Reader) Read(n int) ([]byte, error) {
	p := make([]byte, n)
	if rn, err := r.Reader.Read(p); err != nil {
		return nil, err
	} else if rn != n {
		return nil, ErrReadNotEnoughBytes
	}
	return p, nil
}

func (r *Reader) ReadBytePtr() (*byte, error) {
	b, err := r.ReadByte()
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *Reader) ReadBool() (bool, error) {
	b, err := r.ReadByte()
	if err != nil {
		return false, err
	}
	if b == 0 {
		return false, nil
	}
	if b == 1 {
		return true, nil
	}
	return false, ErrProtocol
}

func (r *Reader) ReadBoolPtr() (*bool, error) {
	b, err := r.ReadBool()
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *Reader) ReadUint8() (uint8, error) {
	b, err := r.ReadByte()
	if err != nil {
		return 0, err
	}
	return uint8(b), nil
}

func (r *Reader) ReadUint8Ptr(reader io.Reader) (*uint8, error) {
	b, err := r.ReadUint8()
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *Reader) ReadUint16() (uint16, error) {
	p, err := r.Read(2)
	if err != nil {
		return 0, err
	}
	return utils.ToUint16(p)
}

func (r *Reader) ReadUint16Ptr() (*uint16, error) {
	p, err := r.ReadUint16()
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *Reader) ReadUint32() (uint32, error) {
	p, err := r.Read(4)
	if err != nil {
		return 0, err
	}
	return utils.ToUint32(p)
}

func (r *Reader) ReadUint32Ptr() (*uint32, error) {
	p, err := r.ReadUint32()
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *Reader) ReadUint64() (uint64, error) {
	p, err := r.Read(8)
	if err != nil {
		return 0, err
	}
	return utils.ToUint64(p)
}

func (r *Reader) ReadUint64Ptr() (*uint64, error) {
	p, err := r.ReadUint64()
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *Reader) ReadLength() (int, error) {
	l, err := r.ReadUint16()
	if err != nil {
		return 0, err
	}
	return int(l), err
}

func (r *Reader) ReadUTF8Str() (string, int, error) {
	var s string
	n, err := r.ReadLength()
	if err != nil {
		return s, 0, err
	}
	if n != 0 {
		if p, err := r.Read(n); err != nil {
			return "", 0, nil
		} else {
			s = string(p)
		}
	}
	return s, 2 + n, nil
}

func (r *Reader) ReadUTF8Ptr() (*string, int, error) {
	s, n, err := r.ReadUTF8Str()
	if err != nil {
		return nil, 0, err
	}
	return &s, n, nil
}

func (r *Reader) ReadVariableByteInteger() (int, int, error) {
	var multiplier int = 1 // Multiplier for each byte (1, 128, 16384, ...)
	var length int = 0     // The length being built
	var byteValue byte     // Single byte to read
	var n int = 0

	for {
		n++
		// Read one byte from the reader
		b, err := r.ReadByte()
		if err != nil {
			return 0, 0, err
		}
		byteValue = b

		// Add the 7 bits to the length value
		length += int(byteValue&0x7F) * multiplier

		// If the MSB (most significant bit) is 0, it's the last byte
		if byteValue&0x80 == 0 {
			break
		}

		// Update multiplier for next byte (128, 16384, 2097152, etc.)
		multiplier *= 128
	}

	return length, n, nil
}
