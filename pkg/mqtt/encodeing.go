package mqtt

import (
	"errors"
	"io"
)

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

func strLength(l []byte) uint16 {
	ret := uint16(l[1]) + (uint16(l[0]) << 8)
	return ret
}

func decodeKeepAlive(l []byte) uint16 {
	return strLength(l)
}

func decodeLength(reader io.Reader) (l uint16, err error) {
	b := make([]byte, 2)
	if _, err = reader.Read(b); err != nil {
		return
	}
	l = strLength(b)
	return
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
