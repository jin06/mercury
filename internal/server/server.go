package server

import (
	"github.com/jin06/mercury/pkg/mqtt"
)

type Server interface {
	Reg(c Client) error
	HandleConnect(packet *mqtt.Connect) error
}
