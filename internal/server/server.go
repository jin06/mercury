package server

import (
	"github.com/jin06/mercury/pkg/mqtt"
)

type Server interface {
	HandleConnect(packet *mqtt.Connect) error
}
