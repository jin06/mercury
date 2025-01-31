package server

import (
	"github.com/jin06/mercury/internal/server/clients"
	"github.com/jin06/mercury/logs"
)

func NewServer() *Server {
	s := Server{
		Clients: map[string]*clients.Client{},
	}
	return &s
}

type Server struct {
	Clients map[string]*clients.Client
}

func (s *Server) Run() (err error) {
	select {}
}

func (s *Server) On(c *clients.Client) {
	s.Clients[c.ID] = c
	logs.Logger.Debug().Msgf("new client %v", c.ID)
}
