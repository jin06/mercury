package clients

import (
	"context"
	"net"

	"github.com/jin06/mercury/logs"
	"github.com/jin06/mercury/pkg/mqtt"
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
	reader *mqtt.Reader
	writer *mqtt.Writer
}

func (c *Client) Run(ctx context.Context) (err error) {
	c.server.On(c)
	c.reader = mqtt.NewReader(c.Conn)
	for {
		// b := make([]byte, 1000)
		p, err := mqtt.ReadPacket(c.reader)
		if err != nil {
			logs.Logger.Err(err)
			// panic(err)
		}
		logs.Logger.Info().Msgf("%v", p)
		// c.HandlePacket(p)
		// os.Exit(1)
		// logs.Logger.Info().Msgf("%b", b)
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

func (c *Client) Connect() {

}
