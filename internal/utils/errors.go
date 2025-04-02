package utils

import (
	"errors"

	"github.com/jin06/mercury/pkg/mqtt"
)

var (
	ErrNotConnectPacket = errors.New("not connect packet error")
	ErrClosedChannel    = errors.New("closed channel")
	ErrMalformedPacket  = errors.New("malformed packet")
	ErrNotValidTopic    = errors.New("not valid topic")
	ErrPacketIDUsed     = errors.New("packet ID is already used")
	ErrPacketIDNotExist = errors.New("packet ID is not exist")
)

func PacketError(p mqtt.Packet, err error) {
	if p == nil || err == nil {
		return
	}
}
