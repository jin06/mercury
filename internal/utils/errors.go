package utils

import "errors"

var (
	ErrNotConnectPacket = errors.New("not connect packet error")
	ErrClosedChannel    = errors.New("closed channel")
)
