package store

import (
	"context"

	"github.com/jin06/mercury/internal/config"
	"github.com/jin06/mercury/internal/model"
	badgerStore "github.com/jin06/mercury/internal/server/message/store/badger"
	memStore "github.com/jin06/mercury/internal/server/message/store/memory"
	"github.com/jin06/mercury/pkg/mqtt"
)

func NewStore(mode string, cid string, clean bool) Store {
	var s Store
	switch mode {
	case "memory":
		s = memStore.New(cid, clean)
	case "badger":
		s = badgerStore.New(cid, clean)
	case "redis":
	default:
		panic("memory store: " + config.Def.MessageStore.Mode)
	}
	return s
}

type Store interface {
	Publish(p *mqtt.Publish) (*model.Record, error)
	Ack(mqtt.PacketID) error
	Receive(*mqtt.Pubrel) error
	Complete(mqtt.PacketID) error
	Release(*mqtt.Pubcomp) error
	Run(ctx context.Context, ch chan mqtt.Packet) error
	Clean() error
	Close() error
}
