package utils

import "errors"

func ToUint16(s []byte) (ret uint16, err error) {
	if len(s) > 2 {
		return ret, errors.New("input slice length must <= 4")
	}
	for i := 0; i < len(s); i++ {
		ret = ret << 8
		ret = ret + uint16(s[i])
	}
	return
}

func ToUint32(s []byte) (ret uint32, err error) {
	if len(s) > 4 {
		return ret, errors.New("input slice length must <= 4")
	}
	for i := 0; i < len(s); i++ {
		ret = ret << 8
		ret = ret + uint32(s[i])
	}
	return
}

func ToUint64(s []byte) (ret uint64, err error) {
	if len(s) > 8 {
		return ret, errors.New("input slice length must <= 4")
	}
	for i := 0; i < len(s); i++ {
		ret = ret << 8
		ret = ret + uint64(s[i])
	}
	return
}
