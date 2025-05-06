package mqtt

const (
	MAXUINT16 = 65535
)

func encodeLength(l int) ([]byte, error) {
	res := make([]byte, 2)
	res[0] = byte(l >> 8)
	res[1] = byte(l)
	return res, nil
}

func decodeLength(b []byte) (int, error) {
	l, err := decodeUint16(b)
	if err != nil {
		return 0, err
	}
	return int(l), nil
}

func encodeKeepAlive(u uint16) []byte {
	return encodeUint16(u)
}

func decodeKeepAlive(data []byte) uint16 {
	res, _ := decodeUint16(data)
	return res
}

func encodeBool(source bool) byte {
	if source {
		return 1
	}
	return 0
}

func decodeBool(b byte) (bool, error) {
	if b == 0 {
		return false, nil
	}
	if b == 1 {
		return true, nil
	}
	return false, ErrProtocol
}

func decodeBoolPtr(b byte) (*bool, error) {
	res, err := decodeBool(b)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func encodeUint16(source uint16) []byte {
	return []byte{
		byte(source >> 8),
		byte(source),
	}
}

func decodeUint16(data []byte) (uint16, error) {
	var ret uint16
	if len(data) > 2 {
		return 0, ErrBytesShorter
	}
	for i := 0; i < len(data); i++ {
		ret = ret << 8
		ret = ret + uint16(data[i])
	}
	return ret, nil
}

func decodeUint16Ptr(data []byte) (*uint16, error) {
	res, err := decodeUint16(data)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func encodeUint32(source uint32) []byte {
	return []byte{
		byte(source >> 24),
		byte(source >> 16),
		byte(source >> 8),
		byte(source),
	}
}

func decodeUint32(data []byte) (uint32, error) {
	var ret uint32
	if len(data) > 4 {
		return 0, ErrBytesShorter
	}
	for i := range data {
		ret = ret << 8
		ret = ret + uint32(data[i])
	}
	return ret, nil
}

func decodeUint32Ptr(data []byte) (*uint32, error) {
	res, err := decodeUint32(data)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// encodeUint64 encodes a uint64 value into a byte slice.
func encodeUint64(value uint64) []byte {
	var data []byte
	// Extract each byte by shifting the uint64 value and appending it to the slice
	for i := range 8 {
		// Shift and mask to get each byte from most significant to least significant
		data = append([]byte{byte(value >> (8 * (7 - i)))}, data...)
	}
	// Remove leading zero bytes, keeping only necessary bytes to represent the uint64
	for len(data) > 1 && data[0] == 0 {
		data = data[1:]
	}
	return data
}

func decodeUint64(data []byte) (uint64, error) {
	if len(data) > 8 {
		return 0, ErrBytesShorter
	}
	var ret uint64
	for i := range data {
		ret = ret << 8
		ret = ret + uint64(data[i])
	}
	return ret, nil
}

func encodePacketID(id PacketID) []byte {
	return []byte{
		byte(id >> 8),
		byte(id),
	}
}

func decodeBytePrt(b byte) *QoS {
	return (*QoS)(&b)
}

// func encodeStringPair(pair map[string]string) ([]byte, error) {
// 	var data []byte
// 	for key, val := range pair {

// 	}
// 	return
// }

func encodeStringPair(k, v string) ([]byte, error) {
	key, err := encodeUTF8Str(k)
	if err != nil {
		return nil, err
	}
	value, err := encodeUTF8Str(v)
	if err != nil {
		return nil, err
	}
	return append(key, value...), nil
}

func decodeStringPair(data []byte) (string, string, int, error) {
	k, nk, err := decodeUTF8Str(data)
	if err != nil {
		return "", "", nk, err
	}
	v, nv, err := decodeUTF8Str(data[nk:])
	if err != nil {
		return "", "", nk + nv, err
	}
	return k, v, nk + nv, nil
}
