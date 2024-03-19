package broker

import "context"

type Broker struct {
}

func (b *Broker) Run(ctx context.Context) (err error) {
	// go listen
	// listener create connection
	// create client -> conn -> read bytes -> tranlate packet -> server -> if publish packet -> search topics tree -> send client
	// server
	return
}
