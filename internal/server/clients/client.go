package clients

import (
	"context"
	"net"
	"os"

	"github.com/jin06/mercury/logs"
	"github.com/jin06/mercury/pkg/encoder"
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
	reader := encoder.NewReader(c.Conn)
	for {
		// b := make([]byte, 1000)
		p, err := reader.ReadPacket()
		if err != nil {
			panic(err)
		}
		logs.Logger.Info().Msgf("%v", p)
		os.Exit(1)
		// logs.Logger.Info().Msgf("%b", b)
	}
	return
}
