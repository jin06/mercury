package clients

import "github.com/jin06/mercury/pkg/mqtt"

type Handler interface {
	Reg(*Client) error
	HandleConnect(packet *mqtt.Connect) error
}
