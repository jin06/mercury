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
	ErrProtocolViolation  = errors.New("protocol violation")
	ErrUTFLengthTooLong   = errors.New("max length that can be encoded in 2 bytes")
	ErrReadNotEnoughBytes = errors.New("not enough bytes read from reader")
	ErrInvalidQoS         = errors.New("invalid QoS")
)
