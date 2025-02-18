package clients

import (
	"context"
	"net"
	"sync"

	"github.com/jin06/mercury/internal/server"
	"github.com/jin06/mercury/internal/utils"
	"github.com/jin06/mercury/logs"
	"github.com/jin06/mercury/pkg/mqtt"
)

func NewClient(handler server.Server, conn net.Conn) *generic {
	c := generic{
		handler:   handler,
		Conn:      conn,
		connected: make(chan struct{}),
		closing:   make(chan struct{}),
		closed:    make(chan struct{}),
		options:   Options{},
		reader:    mqtt.NewReader(conn),
		writer:    mqtt.NewWriter(conn),
	}
	return &c
}

type generic struct {
	ID        string
	Conn      net.Conn
	handler   server.Server
	reader    *mqtt.Reader
	writer    *mqtt.Writer
	options   Options
	connected chan struct{}
	closing   chan struct{}
	closed    chan struct{}
	closeOnce sync.Once
}

func (c *generic) ClientID() string {
	return c.ID
}

func (c *generic) Run(ctx context.Context) (err error) {
	c.handler.Reg(c)
	p, err := mqtt.ReadPacket(c.reader)
	if err != nil {
		logs.Logger.Err(err)
		return err
	}
	if packet, ok := p.(*mqtt.Connect); ok {
		if err := c.handler.HandleConnect(packet); err != nil {
			logs.Logger.Err(err)
			return err
		}
	} else {
		logs.Logger.Err(utils.ErrNotConnectPacket)
		return
	}
	c.runloop(ctx)
	logs.Logger.Info().Msgf("%v", p)
	return nil
}

func (c *generic) runloop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		}
	}
}

func (c *generic) HandlePacket(p mqtt.Packet) {
	var response mqtt.Packet

	switch p.(type) {
	case *mqtt.Connect:
		{
			response = &mqtt.Connack{
				FixHeader: &mqtt.FixedHeader{
					PacketType:      mqtt.CONNACK,
					Flags:           0,
					RemainingLength: 2,
				},
			}
		}
	}
	response.Write(c.writer)
}

func (c *generic) connect() {
}

func (c *generic) Connect() {

}

func (c *generic) ReadStream() {

}

func (c *generic) WriteStream() {

}

func (c *generic) Close(ctx context.Context) error {
	var err error
	c.closeOnce.Do(func() {
		close(c.closing)
	})
	return err
}
