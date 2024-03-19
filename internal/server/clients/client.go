package clients

import (
	"context"
	"net"

	"github.com/jin06/mercury/logs"
)

func NewClient(server Server, conn net.Conn) *Client {
	c := Client{
		server: server,
		Conn:   conn,
	}
	return &c
}

type Client struct {
	ID     string
	Conn   net.Conn
	server Server
}

func (c *Client) Run(ctx context.Context) (err error) {
	c.server.On(c)
	for {
		b := make([]byte, 1000)
		_, rerr := c.Conn.Read(b)
		if rerr != nil {
			logs.Logger.Err(rerr)
		}
		logs.Logger.Info().Msgf("%b", b)
	}
	return
}
