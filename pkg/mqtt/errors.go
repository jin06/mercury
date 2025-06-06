package mqtt

import "errors"

var (
	ErrProtocol           = errors.New("protocol error")
	ErrMalformedPacket    = errors.New("malformed packet")
	ErrUnsupportVersion   = errors.New("unsupported version")
	ErrNullClientID       = errors.New("null client id")
	ErrMaximumPacketSize  = errors.New("maximum packet size can't by zero")
	ErrBytesShorter       = errors.New("the length of []byte is shorter than expected")
	ErrUTFLengthShoter    = errors.New("the length of []byte is shorter for UTF string")
	ErrPacketEncoding     = errors.New("packet encoding error")
	ErrPacketDecoding     = errors.New("packet decoding error")
	ErrProtocolViolation  = errors.New("protocol violation")
	ErrUTFLengthTooLong   = errors.New("max length that can be encoded in 2 bytes")
	ErrReadNotEnoughBytes = errors.New("not enough bytes read from reader")
	ErrInvalidQoS         = errors.New("invalid QoS")
	ErrInsufficientData   = errors.New("insufficient data")
	ErrNotUint16          = errors.New("value out of uint16 range")
	ErrTopicIsEmpty       = errors.New("topic is empty")
	ErrNotValidTopic      = errors.New("not valid topic")
)

type ProtocolError struct{}

func (p *ProtocolError) Error() string {
	return ""
}
