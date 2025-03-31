package mqtt

type BinaryData struct {
	Length uint16
	Data   []byte
}

func (b *BinaryData) Encode() []byte {
	data, _ := encodeBinaryData(b.Data)
	return data
}

func (b *BinaryData) Decode(data []byte) (int, error) {
	res, n, err := decodeBinaryData(data)
	if err != nil {
		return n, err
	}
	b.Length = uint16(len(res))
	b.Data = res
	return n, nil
}

func encodeBinaryData(data []byte) ([]byte, error) {
	l := len(data)
	lbytes, err := encodeLength(l)
	if err != nil {
		return nil, err
	}
	res := append(lbytes, data...)
	return res, nil
}

func decodeBinaryData(data []byte) (res []byte, n int, err error) {
	if len(data) < 2 {
		return nil, 0, ErrBytesShorter
	}
	l, err := decodeLength(data[:2])
	if err != nil {
		return nil, 0, err
	}
	return data[2 : 2+l], 2 + int(l), nil
}
