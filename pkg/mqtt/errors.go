package mqtt

import "errors"

var (
	ProtocolError    = errors.New("protocol error")
	UnsupportVersion = errors.New("unsupported version")
)
