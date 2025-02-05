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

func (r *Reader) ReadUTF8Str() (string, error) {
	var s string
	l, err := r.ReadUint16()
	if err != nil {
		return s, err
	}
	if l != 0 {
		if p, err := r.Read(int(l)); err != nil {
			return "", nil
		} else {
			s = string(p)
		}
	}
	return s, nil
}

func (r *Reader) ReadUTF8Ptr() (*string, error) {
	s, err := r.ReadUTF8Str()
	if err != nil {
		return nil, err
	}
	return &s, nil
}

const (
	maxString int = 65535
)

// func readUint64(reader io.Reader) (uint64, error) {
// 	res, err := read(reader, 8)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return utils.ToUint64(res)
// }

// func readUint64Ptr(reader io.Reader) (*uint64, error) {
// 	res, err := readUint64(reader)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &res, nil
// }

// func readUint32(reader io.Reader) (uint32, error) {
// 	res, err := read(reader, 4)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return utils.ToUint32(res)
// }

// func readUint32Ptr(reader io.Reader) (*uint32, error) {
// 	res, err := readUint32(reader)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &res, nil
// }

// func readUint16(reader io.Reader) (uint16, error) {
// 	res, err := read(reader, 2)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return utils.ToUint16(res)
// }

// func readUint16Ptr(reader io.Reader) (*uint16, error) {
// 	res, err := readUint16(reader)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &res, nil
// }

// func readBool(reader io.Reader) (bool, error) {
// 	b, err := readByte(reader)
// 	if err != nil {
// 		return false, err
// 	}
// 	if b == 0 {
// 		return false, nil
// 	}
// 	if b == 1 {
// 		return true, nil
// 	}
// 	return false, ErrProtocol
// }

// func readBoolPtr(reader io.Reader) (*bool, error) {
// 	res, err := readBool(reader)
// 	if err != nil {
// 		return &res, err
// 	}
// 	return nil, err
// }

// func readUint8(reader io.Reader) (uint8, error) {
// 	res, err := readByte(reader)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return uint8(res), nil
// }

// func readUint8Ptr(reader io.Reader) (*uint8, error) {
// 	res, err := readUint8(reader)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &res, nil
// }

// func readBytes(reader io.Reader, n int) ([]byte, error) {
// 	buf := make([]byte, n)
// 	returnNum, err := io.ReadFull(reader, buf)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if returnNum != n {
// 		return nil, errors.New("no sufficient bytes")
// 	}
// 	return buf, nil
// }

// func readByte(reader io.Reader) (byte, error) {
// 	res, err := read(reader, 1)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return res[0], nil
// }

// func readBytePtr(reader io.Reader) (*byte, error) {
// 	res, err := readByte(reader)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &res, nil
// }

func readStr(reader io.Reader) (string, error) {
	str, _, err := readStrN(reader)
	return str, err
}

func readStrPtr(reader io.Reader) (*string, error) {
	res, err := readStr(reader)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func readStrN(reader io.Reader) (str string, n int, err error) {
	var l uint16
	n = 2
	if l, err = readUint16(reader); err != nil {
		return
	}
	if l != 0 {
		var res []byte
		if res, err = read(reader, int(l)); err != nil {
			return
		}
		str = string(res)
		n = n + int(l)
	}
	return
}

func read(reader io.Reader, n int) ([]byte, error) {
	res := make([]byte, n)
	_, err := io.ReadFull(reader, res)
	return res, err
}

func decodeKeepAlive(l []byte) uint16 {
	res, _ := utils.ToUint16(l)
	return res
}

func readLength(reader io.Reader) (l int, err error) {
	b := make([]byte, 2)
	if _, err = reader.Read(b); err != nil {
		return
	}
	return decodeLength(b)
}

func readProtocolName(reader io.Reader) (res []byte, err error) {
	return readUTF8(reader)
}

// func decodeUTF8(reader io.Reader) (res []byte, err error) {
// 	var l uint16
// 	if l, err = decodeLength(reader); err != nil {
// 		return
// 	}
// 	res = make([]byte, l)
// 	_, err = reader.Read(res)
// 	return
// }

func readUTF8(reader io.Reader) (res []byte, err error) {
	var l int
	if l, err = readLength(reader); err != nil {
		return
	}
	res = make([]byte, l)
	_, err = reader.Read(res)
	return
}

func readUTF8Str(reader io.Reader) (res string, err error) {
	b, err := readUTF8(reader)
	if err != nil {
		return
	}
	return string(b), err
}
