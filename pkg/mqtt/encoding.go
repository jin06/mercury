package mqtt

import (
	"errors"
	"fmt"
	"math"
)

func encodeVariableByteInteger(length int) ([]byte, error) {
	if length < 0 {
		return nil, errors.New("length cannot be negative")
	}

	var result []byte
	for length > 0 {
		byteValue := byte(length & 0x7F)
		if len(result) > 0 {
			byteValue |= 0x80
		}

		result = append(result, byteValue)

		length >>= 7
	}

	if len(result) == 0 {
		return []byte{0}, nil
	}

	return result, nil
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
func decodeVariableByteInteger(data []byte) (int, int, error) {
	var multiplier int = 1 // Multiplier for each byte (1, 128, 16384, ...)
	var length int = 0     // The length being built
	var byteValue byte     // Single byte to read
	var idx int = 0        // Index to read from the byte slice
	var n int = 0

	for {
		if idx >= len(data) {
			return 0, n, ErrBytesShorter
		}

		// Read one byte from the slice
		byteValue = data[idx]
		idx++
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

	return length, n, nil
}

func readVariableByteInteger(reader *Reader) (uint64, error) {
	var multiplier uint64 = 1 // Multiplier for each byte (1, 128, 16384, ...)
	var length uint64 = 0     // The length being built
	var byteValue byte        // Single byte to read

	for {
		// Read one byte from the reader
		b, err := reader.ReadByte()
		if err != nil {
			return 0, err
		}
		byteValue = b

		// Add the 7 bits to the length value
		length += uint64(byteValue&0x7F) * multiplier

		// If the MSB (most significant bit) is 0, it's the last byte
		if byteValue&0x80 == 0 {
			break
		}

		// Update multiplier for next byte (128, 16384, 2097152, etc.)
		multiplier *= 128
	}

	return length, nil
}

func writeVariableByteInteger(writer *Writer, length uint64) error {
	bytes, err := encodeVariableByteInteger(length)
	if err != nil {
		return err
	}
	_, err = writer.Write(bytes)
	return err
}

func bytesToUint64(l []byte) (ret uint64, err error) {
	if len(l) > 4 {
		return ret, errors.New("input slice length must <= 4")
	}
	for i := 0; i < len(l); i++ {
		ret = ret << 8
		ret = ret + uint64(l[i])
	}
	return
}
