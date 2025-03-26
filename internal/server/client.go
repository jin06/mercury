package server

import (
	"context"

	"github.com/jin06/mercury/pkg/mqtt"
)

type Client interface {
	Run(ctx context.Context) error
	Close(ctx context.Context) error
	ClientID() string
	UUID() string
	Write(p mqtt.Packet) (err error)
	Read() (mqtt.Packet, error)
	KeepAlive()
}
