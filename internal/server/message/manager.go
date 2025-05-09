package message

import (
	"errors"
	"sync"

	"github.com/jin06/mercury/internal/config"
	"github.com/jin06/mercury/internal/model"
	"github.com/jin06/mercury/internal/server/message/store"
	badgerStore "github.com/jin06/mercury/internal/server/message/store/badger"
	memStore "github.com/jin06/mercury/internal/server/message/store/memory"
	"github.com/jin06/mercury/pkg/mqtt"
)

type Manager struct {
	clients  map[string]store.Store
	mu       sync.RWMutex
	delivery chan *model.Record
	newStore func(cid string, delivery chan *model.Record) store.Store
}

func NewManager(delivery chan *model.Record) *Manager {
	m := &Manager{
		clients:  map[string]store.Store{},
		delivery: delivery,
	}
	switch config.Def.MessageStore.Mode {
	case "memory":
		m.newStore = func(cid string, delivery chan *model.Record) store.Store {
			return memStore.NewMemStore(cid, delivery)
		}
	case "badger":
		m.newStore = func(cid string, delivery chan *model.Record) store.Store {
			return badgerStore.NewBadgerStore(cid, delivery)
		}
	default:
		panic("memory store: " + config.Def.MessageStore.Mode)
	}
	return m
}

func (m *Manager) Publish(p *mqtt.Publish, cid string) (*model.Record, error) {
	if s := m.Get(cid); s == nil {
		if err := m.Set(cid); err != nil {
			return nil, err
		}
	}
	return m.Get(cid).Publish(p)
}

func (m *Manager) Receive(cid string, p *mqtt.Pubrel) error {
	if s := m.Get(cid); s != nil {
		err := s.Receive(p)
		return err
	}
	return nil
}

func (m *Manager) Ack(cid string, packetID mqtt.PacketID) error {
	if s := m.Get(cid); s != nil {
		return s.Ack(packetID)
	}
	return nil
}

func (m *Manager) Release(cid string, p *mqtt.Pubcomp) error {
	if s := m.Get(cid); s != nil {
		return s.Release(p)
	}
	return nil
}

func (m *Manager) Complete(cid string, pid mqtt.PacketID) error {
	if s := m.Get(cid); s != nil {
		return s.Complete(pid)
	}
	return nil
}

func (m *Manager) Get(cid string) store.Store {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.clients[cid]
}

func (m *Manager) Set(cid string) (err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.clients[cid]; ok {
		return errors.New("exist cid")
	}
	m.clients[cid] = m.newStore(cid, m.delivery)
	return
}

func (m *Manager) Del(cid string) (err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.clients, cid)
	return nil
}
