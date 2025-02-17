package clients

import (
	"context"
	"net"

	"github.com/jin06/mercury/internal/utils"
	"github.com/jin06/mercury/logs"
	"github.com/jin06/mercury/pkg/mqtt"
)

func NewClient(handler Handler, conn net.Conn) *Client {
	c := Client{
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

type Client struct {
	ID        string
	Conn      net.Conn
	handler   Handler
	reader    *mqtt.Reader
	writer    *mqtt.Writer
	options   Options
	connected chan struct{}
	closing   chan struct{}
	closed    chan struct{}
}

func (c *Client) Run(ctx context.Context) (err error) {
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

func (c *Client) runloop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		}
	}
}

func (c *Client) HandlePacket(p mqtt.Packet) {
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

func (c *Client) connect() {
}

func (c *Client) Connect() {

}

func (c *Client) ReadStream() {

}

func (c *Client) WriteStream() {

}
