package mqtt

import "errors"

var (
	ErrProtocol         = errors.New("protocol error")
	ErrUnsupportVersion = errors.New("unsupported version")
	ErrNullClientID     = errors.New("null client id")
)
