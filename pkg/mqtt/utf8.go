package mqtt

type UTF8String string

func (u UTF8String) Encode() ([]byte, error) {
	return encodeUTF8Str(u)
}

func (u *UTF8String) Decode(data []byte) (int, error) {
	str, n, err := decodeUTF8Str(data)
	if err != nil {
		return n, err
	}
	*u = UTF8String(str)
	return n, nil
}

func (u UTF8String) String() string {
	return string(u)
}

func encodeUTF8(data []byte) ([]byte, error) {
	// Calculate the length of the data to encode
	l := len(data)

	// Ensure the length fits within the 2-byte limit
	if l > MAXUINT16 { // Max length that can be encoded in 2 bytes
		return nil, ErrUTFLengthTooLong
	}

	// Create a new byte slice to hold the length and the data
	result := make([]byte, 2+l)

	// Encode the length into the first two bytes
	length, err := encodeLength(l)
	if err != nil {
		return nil, err
	}
	copy(result[0:2], length)

	// Copy the data into the result slice after the length
	copy(result[2:], data)

	// Return the encoded byte slice without the length
	return result, nil
}

func decodeUTF8(data []byte) (res []byte, n int, err error) {
	if len(data) < 2 {
		return nil, 0, ErrBytesShorter
	}
	l, err := decodeLength(data[:2])
	if err != nil {
		return nil, 0, err
	}
	total := 2 + int(l)
	if len(data) < total {
		return nil, 0, ErrUTFLengthShoter
	}
	return data[2:total], total, nil
}

func encodeUTF8Str[T string | UTF8String](s T) ([]byte, error) {
	return encodeUTF8([]byte(s))
}

func decodeUTF8Str(data []byte) (string, int, error) {
	r, n, err := decodeUTF8(data)
	return string(r), n, err
}

func decodeUTF8Ptr(data []byte) (*string, int, error) {
	s, n, err := decodeUTF8Str(data)
	if err != nil {
		return nil, n, err
	}
	return &s, n, err
}
