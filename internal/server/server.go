package server

import (
	"github.com/jin06/mercury/pkg/mqtt"
)

type Server interface {
	Reg(client Client) error
	HandlePacket(packet mqtt.Packet) (response mqtt.Packet, err error)
}
