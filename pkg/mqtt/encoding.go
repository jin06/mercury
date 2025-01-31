package mqtt

import (
	"errors"
)

func variableByteInteger(length int) ([]byte, error) {
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
