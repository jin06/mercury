package server

import (
	"sync"

	"github.com/jin06/mercury/internal/server/clients"
)

func NewManager() *Manager {
	return &Manager{
		clients: map[string]*clients.Client{},
	}
}

type Manager struct {
	clients map[string]*clients.Client
	mu      sync.Mutex
}

func (m *Manager) Set(c *clients.Client) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.clients[c.ID] != nil {
		return
	}
	m.clients[c.ID] = c
}

func (m *Manager) Remove(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.clients, id)
}

func (m *Manager) Get(id string) *clients.Client {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.clients[id]
}

func (m *Manager) All() map[string]*clients.Client {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.clients
}

func (m *Manager) Len() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.clients)
}

func (m *Manager) Iterator(f func(c *clients.Client)) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, c := range m.clients {
		f(c)
	}
}
