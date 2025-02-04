package mqtt

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
