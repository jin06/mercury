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

func readUint32(reader io.Reader) (uint32, error) {
	res, err := read(reader, 4)
	if err != nil {
		return 0, err
	}
	return utils.ToUint32(res)
}

func readUint16(reader io.Reader) (uint16, error) {
	res, err := read(reader, 2)
	if err != nil {
		return 0, err
	}
	return utils.ToUint16(res)
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

func readByte(reader io.Reader) (byte, error) {
	res, err := read(reader, 1)
	if err != nil {
		return 0, err
	}
	return res[0], nil
}

func readStr(reader io.Reader) (str string, err error) {
	str, _, err = readStrN(reader)
	return
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

func uint32ToBytes(source uint32) []byte {
	return []byte{
		byte(source >> 24),
		byte(source >> 16),
		byte(source >> 8),
		byte(source),
	}
}

func uint16ToBytes(source uint16) []byte {
	return []byte{
		byte(source >> 8),
		byte(source),
	}
}

func boolTobyte(source bool) byte {
	if source {
		return 1
	}
	return 0
}

func encodeString(s string) ([]byte, error) {
	if len(s) > maxString {
		return nil, errors.New("")
	}
	l := uint16(len(s))
	result := make([]byte, 0, l+2)
	result = append(result, byte(l>>8), byte(l))
	result = append(result, []byte(s)...)
	return result, nil
}

func encodeUint16(i uint16) []byte {
	return []byte{
		byte(i >> 8),
		byte(i),
	}
}
