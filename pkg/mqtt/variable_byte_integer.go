package mqtt

import (
	"errors"
	"fmt"
	"math"
)

type VariableByteInteger int

func (v VariableByteInteger) Int() int {
	return int(v)
}

func (v *VariableByteInteger) Clone() *VariableByteInteger {
	if v == nil {
		return nil
	}
	clone := *v
	return &clone
}

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
