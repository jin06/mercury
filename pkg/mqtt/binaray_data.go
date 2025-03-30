package mqtt

type BinaryData struct {
	Length uint16
	Data   []byte
}

func (b *BinaryData) Encode() []byte {
	lengthData := encodeUint16(b.Length)
	data := append(lengthData, b.Data...)
	return data
}

func (b *BinaryData) Decode(data []byte) (int, error) {
	if len(data) < 2 {
		return 0, ErrBytesShorter
	}
	length := len(data)
	if length > MAXUINT16 {
		return 0, ErrNotUint16
	}
	b.Length = uint16(length)
	b.Data = data[2:]
	return 2 + length, nil
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
	l, err := decodeLength(data[:2])
	if err != nil {
		return nil, 0, err
	}
	return data[2 : 2+l], 2 + int(l), nil
}
