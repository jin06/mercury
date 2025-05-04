package store

import (
	"github.com/jin06/mercury/internal/model"
	"github.com/jin06/mercury/pkg/mqtt"
)

type Store interface {
	Publish(p *mqtt.Publish) (*model.Record, error)
	Ack(mqtt.PacketID) error
	Receive(*mqtt.Pubrel) error
	Complete(mqtt.PacketID) error
	Release(*mqtt.Pubcomp) error
}
