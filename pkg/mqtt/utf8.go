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
