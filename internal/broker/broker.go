package broker

import (
	"context"
	"net"
)

type Broker struct {
}

func (b *Broker) Run(ctx context.Context) (err error) {
	// go listen
	// listener create connection
	// create client -> conn -> read bytes -> tranlate packet -> server -> if publish packet -> search topics tree -> send client
	// server
	return
}

func listen() {
	listener, err := net.Listen("tcp", ":1883")
	if err != nil {
		panic(err)
	}

	for conn, err := range listener.Accept() {

	}
}
