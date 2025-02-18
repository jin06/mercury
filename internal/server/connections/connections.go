package connections

import (
	"github.com/jin06/mercury/pkg/mqtt"
)

// need to implement net.Conn
type Connection interface {
	Read() (mqtt.Packet, error)
	Write(p mqtt.Packet) error
	Close() error
}
