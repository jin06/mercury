package message

import (
	"errors"
	"sync"

	"github.com/jin06/mercury/internal/model"
	"github.com/jin06/mercury/pkg/mqtt"
)

type Manager struct {
	clients  map[string]Store
	mu       sync.RWMutex
	delivery chan *model.Record
}

func NewManager(delivery chan *model.Record) *Manager {
	m := &Manager{
		clients:  map[string]Store{},
		delivery: delivery,
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

func (m *Manager) Get(cid string) Store {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.clients[cid]
}

func (m *Manager) Set(cid string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.clients[cid]; ok {
		return errors.New("exist cid")
	}
	m.clients[cid] = NewMemStore(cid, m.delivery)
	return nil
}
