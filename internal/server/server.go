package server

import (
	"github.com/jin06/mercury/internal/server/clients"
	"github.com/jin06/mercury/pkg/mqtt"
)

type Server interface {
	Reg(*clients.Client) error
	HandleConnect(packet *mqtt.Connect) error
}
