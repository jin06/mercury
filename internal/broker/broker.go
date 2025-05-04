package broker

import (
	"context"
	"net"
	"sync"

	"github.com/rs/zerolog/log"

	"github.com/jin06/mercury/internal/config"
	"github.com/jin06/mercury/internal/server"
	"github.com/jin06/mercury/internal/server/clients"
	badgerStore "github.com/jin06/mercury/internal/server/message/store/badger"
	"github.com/jin06/mercury/internal/server/servers"
)

func NewBroker() *Broker {
	b := &Broker{
		Server:  servers.NewServer(config.Def.Mode),
		closing: make(chan struct{}),
		closed:  make(chan struct{}),
		options: &Options{},
	}
	return b
}

type Broker struct {
	Server    server.Server
	options   *Options
	closeOnce sync.Once
	closing   chan struct{}
	closed    chan struct{}
}

func (b *Broker) Run(ctx context.Context) (err error) {
	defer close(b.closed)
	defer b.close()
	if err = badgerStore.Init(config.Def.MessageStore.BadgerConfig); err != nil {
		return
	}
	return b.listen(ctx)
}

func (b *Broker) listen(ctx context.Context) error {
	wg := sync.WaitGroup{}
	for _, l := range config.Def.Listeners {
		switch l.Type {
		case "tcp":
			wg.Add(1)
			go func() {
				if err := b.listenTCP(ctx, l.Addr); err != nil {
					log.Error().Err(err).Msg("listen tcp error")
				}
				b.close()
				wg.Done()
			}()
		case "mqtt":
		}
	}
	wg.Wait()
	return nil
}

func (b *Broker) listenTCP(ctx context.Context, addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		client := clients.NewClient(b.Server, conn)
		go func() {
			if err := client.Run(ctx); err != nil {
				log.Error().Err(err).Msg("client run error")
			}
		}()
	}
}

func (b *Broker) close() error {
	b.closeOnce.Do(func() {
		close(b.closing)
	})
	return nil
}
