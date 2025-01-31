package mqtt

import (
	"errors"
	"io"

	"github.com/jin06/mercury/pkg/utils"
)

func Encoding(souce []byte) (dst Packet) {
	return
}

func Decoding(souce Packet) (dst []byte) {
	return
}

func variableByteInteger(length int) ([]byte, error) {
	if length < 0 {
		return nil, errors.New("length cannot be negative")
	}

	var result []byte
	for length > 0 {
		byteValue := byte(length & 0x7F)
		if len(result) > 0 {
			byteValue |= 0x80
		}

		result = append(result, byteValue)

		length >>= 7
	}

	if len(result) == 0 {
		return []byte{0}, nil
	}

	return result, nil
}

func bytesToUint64(l []byte) (ret uint64, err error) {
	if len(l) > 4 {
		return ret, errors.New("input slice length must <= 4")
	}
	for i := 0; i < len(l); i++ {
		ret = ret << 8
		ret = ret + uint64(l[i])
	}
	return

}

func decodeKeepAlive(l []byte) uint16 {
	res, _ := utils.ToUint16(l)
	return res
}

func decodeLength(reader io.Reader) (l uint16, err error) {
	b := make([]byte, 2)
	if _, err = reader.Read(b); err != nil {
		return
	}
	return utils.ToUint16(b)
}

func decodeProtocolName(reader io.Reader) (res []byte, err error) {
	// var l uint16
	// if l, err = decodeLength(reader); err != nil {
	// return
	// }
	// res = make([]byte, l)
	// _, err = reader.Read(res)
	// return
	return decodeUTF8(reader)
}

func decodeUTF8(reader io.Reader) (res []byte, err error) {
	var l uint16
	if l, err = decodeLength(reader); err != nil {
		return
	}
	res = make([]byte, l)
	_, err = reader.Read(res)
	return
}

func decodeUTF8Str(reader io.Reader) (res string, err error) {
	b, err := decodeUTF8(reader)
	if err != nil {
		return
	}
	return string(b), err
}
