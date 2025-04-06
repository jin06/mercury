package message

import (
	"github.com/jin06/mercury/internal/model"
	"github.com/jin06/mercury/pkg/mqtt"
)

type Store interface {
	Save(p *mqtt.Publish, source string, dest string) (*model.Record, error)
	Delete(pid mqtt.PacketID) (bool, error)
	Change(id mqtt.PacketID, state model.State) error
	Receive(p *mqtt.Pubrel) error
}
