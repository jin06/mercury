package server

import (
	"github.com/jin06/mercury/pkg/mqtt"
)

type Server interface {
	Register(client Client) error
	Deregister(client Client) error
	HandlePacket(packet mqtt.Packet, cid string) (response mqtt.Packet, err error)
}
