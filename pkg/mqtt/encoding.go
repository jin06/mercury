package mqtt

import (
	"errors"
)

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
