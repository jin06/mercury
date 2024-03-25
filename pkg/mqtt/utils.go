package mqtt

import (
	"io"

	"github.com/jin06/mercury/pkg/utils"
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

func readByte(reader io.Reader) (byte, error) {
	res, err := read(reader, 1)
	if err != nil {
		return 0, err
	}
	return res[0], nil
}

func readStr(reader io.Reader) (str string, n int, err error) {
	var l uint16
	if l, err = readUint16(reader); err != nil {
		return
	}
	var res []byte
	if res, err = read(reader, int(l)); err != nil {
		return
	}
	str = string(res)
	n = int(l + 2)
	return
}

func read(reader io.Reader, n int) ([]byte, error) {
	res := make([]byte, n)
	_, err := io.ReadFull(reader, res)
	return res, err
}
