package message

import (
	"errors"
	"sync"

	"github.com/jin06/mercury/pkg/mqtt"
)

type Manager struct {
	clients map[string]Store
	mu      sync.RWMutex
}

func NewManager() *Manager {
	m := &Manager{
		clients: map[string]Store{},
	}
	return m
}

func (m *Manager) Save(p *mqtt.Publish, source string, dest string) (*Record, error) {
	if s := m.Get(dest); s == nil {
		if err := m.Set(dest); err != nil {
			return nil, err
		}
	}
	return m.Get(dest).Save(p, source, dest)
}

func (m *Manager) Delete(cid string, packetID mqtt.PacketID) error {
	if s := m.Get(cid); s != nil {
		_, err := s.Delete(packetID)
		return err
	}
	return nil
}

func (m *Manager) Change(cid string, packetID mqtt.PacketID, state State) error {
	if s := m.Get(cid); s != nil {
		return s.Change(packetID, state)
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
	m.clients[cid] = NewRingBufferStore()
	return nil
}

func (m *Manager) Remove(cid string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.clients, cid)
}
