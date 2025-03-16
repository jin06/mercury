package server

import (
	"math/rand"
	"sync"
	"time"

	"github.com/jin06/mercury/pkg/mqtt"
)

func NewManager() *Manager {
	return &Manager{
		clients: map[string]Client{},
		// todo persistent store
	}
}

type Manager struct {
	clients map[string]Client
	mu      sync.Mutex
}

func (m *Manager) Set(c Client) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.clients[c.ClientID()] != nil {
		return
	}
	m.clients[c.ClientID()] = c
}

func (m *Manager) Remove(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.clients, id)
}

func (m *Manager) RemoveClient(c Client) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.clients[c.ClientID()]; ok {
		if m.clients[c.ClientID()].UUID() != c.UUID() {
			delete(m.clients, c.ClientID())
		}
	}
	return
}

func (m *Manager) Get(id string) Client {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.clients[id]
}

func (m *Manager) All() map[string]Client {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.clients
}

func (m *Manager) Len() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.clients)
}

func (m *Manager) Iterator(f func(c Client)) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, c := range m.clients {
		f(c)
	}
}

// todo
func (m *Manager) GetPacketID() mqtt.PacketID {
	rand.Seed(time.Now().UnixNano())
	return mqtt.PacketID(rand.Intn(65535) + 1) // Packet ID is a 16-bit unsigned integer
}
