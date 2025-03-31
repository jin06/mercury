package mqtt

import (
	"errors"
	"fmt"
	"math"
)

const (
	MAXUINT16            = 65535
	MAXPACKETID PacketID = 65535
)

func encodeVariableByteInteger[T VariableByteInteger | int](length T) ([]byte, error) {
	if length < 0 {
		return nil, errors.New("length cannot be negative")
	}
	encoded := make([]byte, 0)
	for {
		byteVal := byte(length & 0x7F)
		length >>= 7
		if length > 0 {
			byteVal |= 0x80
		}
		encoded = append(encoded, byteVal)
		if length == 0 {
			break
		}
	}

	if len(encoded) == 0 {
		return []byte{0x00}, nil
	}
	return encoded, nil
}

// decodeVariableByteInteger decodes a variable-length integer from a byte slice ([]byte).
// The integer is encoded in 7-bit chunks, where the MSB indicates if more bytes follow.
// It returns the decoded integer, the number of bytes used, and an error if something goes wrong.
//
// Parameters:
//   - data []byte: A byte slice containing the encoded integer.
//
// Returns:
//   - int: The decoded integer value.
//   - int: The number of bytes consumed.
//   - error: Any error encountered (e.g., byte slice too short or value too large).
func decodeVariableByteInteger(data []byte) (VariableByteInteger, int, error) {
	var multiplier int = 1 // Multiplier for each byte (1, 128, 16384, ...)
	var length int = 0     // The length being built
	var byteValue byte     // Single byte to read
	// var idx int = 0        // Index to read from the byte slice
	var n int = 0

	for {
		if n >= len(data) {
			return 0, n, ErrBytesShorter
		}

		// Read one byte from the slice
		byteValue = data[n]
		// idx++
		n++

		// Add the 7 bits to the length value
		length += int(byteValue&0x7F) * multiplier

		// If the MSB (most significant bit) is 0, it's the last byte
		if byteValue&0x80 == 0 {
			break
		}

		// Update multiplier for next byte (128, 16384, 2097152, etc.)
		multiplier *= 128
	}

	if length > math.MaxInt {
		return 0, n, fmt.Errorf("value exceeds maximum int size")
	}

	return VariableByteInteger(length), n, nil
}

func decodeVariableByteIntegerPtr(data []byte) (*VariableByteInteger, int, error) {
	res, n, err := decodeVariableByteInteger(data)
	if err != nil {
		return nil, n, err
	}
	return &res, n, nil
}

func readVariableByteInteger(reader *Reader) (VariableByteInteger, int, error) {
	var multiplier int = 1 // Multiplier for each byte (1, 128, 16384, ...)
	var length int = 0     // The length being built
	var byteValue byte     // Single byte to read
	var n int = 0

	for {
		n++
		// Read one byte from the reader
		b, err := reader.ReadByte()
		if err != nil {
			return 0, 0, err
		}
		byteValue = b

		// Add the 7 bits to the length value
		length += int(byteValue&0x7F) * multiplier

		// If the MSB (most significant bit) is 0, it's the last byte
		if byteValue&0x80 == 0 {
			break
		}

		// Update multiplier for next byte (128, 16384, 2097152, etc.)
		multiplier *= 128
	}

	return VariableByteInteger(length), n, nil
}

func writeVariableByteInteger(writer *Writer, length VariableByteInteger) error {
	bytes, err := encodeVariableByteInteger(length)
	if err != nil {
		return err
	}
	_, err = writer.Write(bytes)
	return err
}

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

func encodeProtocolVersion(v ProtocolVersion) (byte, error) {
	if err := v.Check(); err != nil {
		return 0, err
	}
	return byte(v), nil
}

func decodeProtocolVersion(b byte) (ProtocolVersion, error) {
	switch b {
	case 3:
		return MQTT3, nil
	case 4:
		return MQTT4, nil
	case 5:
		return MQTT5, nil
	}
	return 0, ErrUnsupportVersion
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
