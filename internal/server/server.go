package server

import (
	"context"

	"github.com/jin06/mercury/internal/server/clients"
)

type Server interface {
	Run(context.Context) error
	Reg(*clients.Client) error
}

func NewServer() Server {
	return newGeneric()
}
