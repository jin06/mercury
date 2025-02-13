package server

import (
	"context"
	"errors"

	"github.com/jin06/mercury/internal/server/clients"
)

func newGeneric() *generic {
	return &generic{
		manager: NewManager(),
	}
}

type generic struct {
	manager *Manager
}

func (g *generic) Run(ctx context.Context) error {
	select {}
}

func (g *generic) Reg(c *clients.Client) error {
	if c == nil {
		return errors.New("client is nil")
	}
	g.manager.Set(c)
	return nil
}
