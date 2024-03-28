package clients

import (
	"context"
	"fmt"
	"net"

	"github.com/jin06/mercury/logs"
	"github.com/jin06/mercury/pkg/encoder"
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
		c.HandlePacket(p)
		// os.Exit(1)
		// logs.Logger.Info().Msgf("%b", b)
	}
	return
}

func (c *Client) HandlePacket(p mqtt.Packet) {
	val, ok := p.(*mqtt.Connect)
	fmt.Println(123)
	if ok {
		ack := mqtt.Connack{
			Properties: &mqtt.ConnackProperties{
				MaximumPacketSize:               1000,
				RetainAvailable:                 true,
				SharedSubscriptionAvailable:     true,
				SubscriptionIdentifierAvailable: true,
				TopicAliasMaximum:               99,
				WildcardSubscriptionAvailable:   true,
				ReceiveMaximum:                  111,
			},
		}
		if val.Version == mqtt.MQTT5 {
			ack.ReasonCode = mqtt.V5_CONACK_SUCCESS
		} else {
			ack.ReasonCode = mqtt.V3_CONNACK_ACCEPT
		}
		fmt.Println(val.Version)
		bytes, err := ack.Encode(val.Version)
		fmt.Println(bytes)
		if err != nil {
			panic(err)
		}
		c.Conn.Write(bytes)
	}
}

func (c *Client) Connect() {

}
