package server

import (
	"github.com/jin06/mercury/internal/model"
	"github.com/jin06/mercury/pkg/mqtt"
)

type Server interface {
	Register(client Client) error
	Deregister(client Client) error
	HandlePacket(packet mqtt.Packet, cid string) (response mqtt.Packet, err error)
	HandleConnect(p *mqtt.Connect, c Client) (resp *mqtt.Connack, err error)
	DeliveryPublish(cid string, p *mqtt.Publish) error
	DeliveryOne(cid string, msg *model.Message) error
}
