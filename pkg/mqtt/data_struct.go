package mqtt

type VariableByteInteger int

func (v VariableByteInteger) Encode() ([]byte, error) {
	return encodeVariableByteInteger(v)
}

func (v *VariableByteInteger) Decode(data []byte) (int, error) {
	length, n, err := decodeVariableByteInteger(data)
	if err != nil {
		return n, err
	}
	*v = length
	return n, nil
}

func (v VariableByteInteger) Int() int {
	return int(v)
}

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
