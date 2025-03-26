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
