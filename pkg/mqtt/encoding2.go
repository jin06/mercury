package mqtt

import (
	"errors"
	"io"

	"github.com/jin06/mercury/pkg/utils"
)

const (
	maxString int = 65535
)

func readUint64(reader io.Reader) (uint64, error) {
	res, err := read(reader, 8)
	if err != nil {
		return 0, err
	}
	return utils.ToUint64(res)
}

func readUint64Ptr(reader io.Reader) (*uint64, error) {
	res, err := readUint64(reader)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func readUint32(reader io.Reader) (uint32, error) {
	res, err := read(reader, 4)
	if err != nil {
		return 0, err
	}
	return utils.ToUint32(res)
}

func decodeUint32(data []byte) (uint32, error) {
	return utils.ToUint32(data)
}

func readUint32Ptr(reader io.Reader) (*uint32, error) {
	res, err := readUint32(reader)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func decodeUint32Ptr(data []byte) (*uint32, error) {
	res, err := decodeUint32(data)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func decodeUint16(data []byte) (uint16, error) {
	return utils.ToUint16(data)
}

func decodeUint16Ptr(data []byte) (*uint16, error) {
	res, err := decodeUint16(data)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func readUint16(reader io.Reader) (uint16, error) {
	res, err := read(reader, 2)
	if err != nil {
		return 0, err
	}
	return utils.ToUint16(res)
}

func readUint16Ptr(reader io.Reader) (*uint16, error) {
	res, err := readUint16(reader)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func readBool(reader io.Reader) (bool, error) {
	b, err := readByte(reader)
	if err != nil {
		return false, err
	}
	if b == 0 {
		return false, nil
	}
	if b == 1 {
		return true, nil
	}
	return false, ErrProtocol
}

func readBoolPtr(reader io.Reader) (*bool, error) {
	res, err := readBool(reader)
	if err != nil {
		return &res, err
	}
	return nil, err
}

func readUint8(reader io.Reader) (uint8, error) {
	res, err := readByte(reader)
	if err != nil {
		return 0, err
	}
	return uint8(res), nil
}

func readUint8Ptr(reader io.Reader) (*uint8, error) {
	res, err := readUint8(reader)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func readBytes(reader io.Reader, n int) ([]byte, error) {
	buf := make([]byte, n)
	returnNum, err := io.ReadFull(reader, buf)
	if err != nil {
		return nil, err
	}
	if returnNum != n {
		return nil, errors.New("no sufficient bytes")
	}
	return buf, nil
}

func readByte(reader io.Reader) (byte, error) {
	res, err := read(reader, 1)
	if err != nil {
		return 0, err
	}
	return res[0], nil
}

func readBytePtr(reader io.Reader) (*byte, error) {
	res, err := readByte(reader)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func readStr(reader io.Reader) (string, error) {
	str, _, err := readStrN(reader)
	return str, err
}

func readStrPtr(reader io.Reader) (*string, error) {
	res, err := readStr(reader)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func readStrN(reader io.Reader) (str string, n int, err error) {
	var l uint16
	n = 2
	if l, err = readUint16(reader); err != nil {
		return
	}
	if l != 0 {
		var res []byte
		if res, err = read(reader, int(l)); err != nil {
			return
		}
		str = string(res)
		n = n + int(l)
	}
	return
}

func read(reader io.Reader, n int) ([]byte, error) {
	res := make([]byte, n)
	_, err := io.ReadFull(reader, res)
	return res, err
}

func uint32ToBytes(source uint32) []byte {
	return []byte{
		byte(source >> 24),
		byte(source >> 16),
		byte(source >> 8),
		byte(source),
	}
}

func packetIDToBytes(id PacketID) []byte {
	return []byte{
		byte(id >> 8),
		byte(id),
	}
}

func uint16ToBytes(source uint16) []byte {
	return []byte{
		byte(source >> 8),
		byte(source),
	}
}

func boolTobyte(source bool) byte {
	if source {
		return 1
	}
	return 0
}

func strToBytes(s string) ([]byte, error) {
	if len(s) > maxString {
		return nil, errors.New("")
	}
	l := uint16(len(s))
	result := make([]byte, 0, l+2)
	result = append(result, byte(l>>8), byte(l))
	result = append(result, []byte(s)...)
	return result, nil
}

func readVariableInt(r io.Reader) (uint32, error) {
	var res uint32
	for i := 0; i < 4; i++ {
		b, err := readByte(r)
		if err != nil {
			return res, err
		}
		res = res + uint32(0b01111111&b<<(7*i))
		if 0b10000000&b != 0b10000000 {
			break
		}
	}
	return res, nil
}

func encodeVariableInt(val uint32) ([]byte, error) {
	res := []byte{}
	for i := 0; i < 4; i++ {
		b := byte(val % 128)
		val = val / 128
		if val > 0 {
			b = 0b10000000 | b
		}
		res = append(res, b)
		if val <= 0 {
			break
		}
	}
	return res, nil
}

func decodeKeepAlive(l []byte) uint16 {
	res, _ := utils.ToUint16(l)
	return res
}

func encodeLength(l int) ([]byte, error) {
	res := make([]byte, 2)
	res[0] = byte(l >> 8)
	res[1] = byte(l)
	return res, nil
}

func decodeLength(b []byte) (int, error) {
	l, err := utils.ToUint16(b)
	if err != nil {
		return 0, err
	}
	return int(l), nil
}

func readLength(reader io.Reader) (l int, err error) {
	b := make([]byte, 2)
	if _, err = reader.Read(b); err != nil {
		return
	}
	return decodeLength(b)
}

func readProtocolName(reader io.Reader) (res []byte, err error) {
	return readUTF8(reader)
}

// func decodeUTF8(reader io.Reader) (res []byte, err error) {
// 	var l uint16
// 	if l, err = decodeLength(reader); err != nil {
// 		return
// 	}
// 	res = make([]byte, l)
// 	_, err = reader.Read(res)
// 	return
// }

func readUTF8(reader io.Reader) (res []byte, err error) {
	var l int
	if l, err = readLength(reader); err != nil {
		return
	}
	res = make([]byte, l)
	_, err = reader.Read(res)
	return
}

func readUTF8Str(reader io.Reader) (res string, err error) {
	b, err := readUTF8(reader)
	if err != nil {
		return
	}
	return string(b), err
}
