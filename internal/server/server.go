package server

import (
	"github.com/jin06/mercury/pkg/mqtt"
)

type Server interface {
	Register(client Client) error
	Deregister(client Client) error
	HandlePacket(packet mqtt.Packet, cid string) (response mqtt.Packet, err error)
	HandleConnect(p *mqtt.Connect, c Client) (resp *mqtt.Connack, err error)
	Dispatch(cid string, p *mqtt.Publish) error
	Delivery(cid string, msg *mqtt.Publish) error
}
