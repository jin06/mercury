package message

import (
	"errors"
	"sync"

	"github.com/jin06/mercury/internal/model"
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

func (m *Manager) Pop(p *model.Message) (mqtt.PacketID, error) {
	s := m.Get(p.Dest)
	if s == nil {
		if err := m.Set(p.Dest); err != nil {
			return 0, err
		}
	}
	return m.Get(p.Dest).Pop(p)
}

func (m *Manager) Rec(cid string, pid mqtt.PacketID) (bool, error) {
	return m.Get(cid).Rec(pid)
}

func (m *Manager) Ack(cid string, pid mqtt.PacketID) (bool, error) {
	return m.Get(cid).Ack(pid)
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
