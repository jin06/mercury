package broker

import (
	"context"
	"net"

	"github.com/jin06/mercury/internal/server"
	"github.com/jin06/mercury/internal/server/clients"
)

func NewBroker() *Broker {
	b := &Broker{
		Server: server.NewServer(),
	}
	return b
}

type Broker struct {
	Server *server.Server
}

func (b *Broker) Run(ctx context.Context) (err error) {
	// go listen
	// listener create connection
	// create client -> conn -> read bytes -> tranlate packet -> server -> if publish packet -> search topics tree -> send client
	// server

	b.listen()
	return
}
func (b *Broker) listen() {
	listener, err := net.Listen("tcp", ":1883")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		client := clients.NewClient(b.Server, conn)
		go client.Run(context.Background())
	}
}
