package mqtt

import (
	"errors"
	"io"

	"github.com/jin06/mercury/pkg/utils"
)

const (
	maxString int = 65535
)

func readUint64(reader io.Reader) (uint64, error) {
	res, err := read(reader, 8)
	if err != nil {
		return 0, err
	}
	return utils.ToUint64(res)
}

func readUint64Ptr(reader io.Reader) (*uint64, error) {
	res, err := readUint64(reader)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func readUint32(reader io.Reader) (uint32, error) {
	res, err := read(reader, 4)
	if err != nil {
		return 0, err
	}
	return utils.ToUint32(res)
}

func readUint32Ptr(reader io.Reader) (*uint32, error) {
	res, err := readUint32(reader)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func readUint16(reader io.Reader) (uint16, error) {
	res, err := read(reader, 2)
	if err != nil {
		return 0, err
	}
	return utils.ToUint16(res)
}

func readUint16Ptr(reader io.Reader) (*uint16, error) {
	res, err := readUint16(reader)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func readBool(reader io.Reader) (bool, error) {
	b, err := readByte(reader)
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

func readBoolPtr(reader io.Reader) (*bool, error) {
	res, err := readBool(reader)
	if err != nil {
		return &res, err
	}
	return nil, err
}

func readUint8(reader io.Reader) (uint8, error) {
	res, err := readByte(reader)
	if err != nil {
		return 0, err
	}
	return uint8(res), nil
}

func readUint8Ptr(reader io.Reader) (*uint8, error) {
	res, err := readUint8(reader)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func readBytes(reader io.Reader, n int) ([]byte, error) {
	buf := make([]byte, n)
	returnNum, err := io.ReadFull(reader, buf)
	if err != nil {
		return nil, err
	}
	if returnNum != n {
		return nil, errors.New("no sufficient bytes")
	}
	return buf, nil
}

func readByte(reader io.Reader) (byte, error) {
	res, err := read(reader, 1)
	if err != nil {
		return 0, err
	}
	return res[0], nil
}

func readBytePtr(reader io.Reader) (*byte, error) {
	res, err := readByte(reader)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

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
