package clients

import (
	"context"
	"fmt"
	"io"
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/jin06/mercury/internal/server"
	"github.com/jin06/mercury/internal/utils"
	"github.com/jin06/mercury/pkg/mqtt"
)

func NewClient(handler server.Server, conn io.ReadWriteCloser) *generic {
	c := generic{
		handler:    handler,
		Connection: mqtt.NewConnection(conn),
		connected:  make(chan struct{}),
		stopping:   make(chan struct{}),
		closed:     make(chan struct{}),
		options:    Options{},
		input:      make(chan mqtt.Packet, 2000),
		output:     make(chan mqtt.Packet, 2000),
	}
	return &c
}

type generic struct {
	id string
	*mqtt.Connection
	handler   server.Server
	options   Options
	connected chan struct{}
	stopping  chan struct{}
	stopOnce  sync.Once
	closed    chan struct{}
	err       error // first error that occurs exits the client
	// packet channels
	input  chan mqtt.Packet
	output chan mqtt.Packet
}

func (c *generic) ClientID() string {
	return c.id
}

func (c *generic) Run(ctx context.Context) (err error) {
	defer close(c.closed)
	defer func() {
		if err != nil {
			c.setError(err)
			c.stop(err)
		}
		c.Close(ctx)
	}()

	var p mqtt.Packet

	if p, err = c.ReadPacket(); err != nil {
		return err
	}

	if packet, ok := p.(*mqtt.Connect); ok {
		fmt.Println(packet)
		if err = c.handler.HandleConnect(packet); err != nil {
			return
		}
	} else {
		return utils.ErrNotConnectPacket
	}
	if err = c.runloop(ctx); err != nil {
		return err
	}
	return nil
}

func (c *generic) runloop(ctx context.Context) error {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := c.readLoop(ctx); err != nil {
			c.setError(err)
			return
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := c.writeLoop(ctx); err != nil {
			c.setError(err)
			return
		}
	}()
	select {
	case <-ctx.Done():
		break
	case <-c.stopping:
		break
	}
	wg.Done()
	return nil
}

func (c *generic) readLoop(ctx context.Context) error {
	for {
		p, err := c.ReadPacket()
		if err != nil {
			return err
		}
		select {
		case <-ctx.Done():
			return nil
		case <-c.stopping:
			return nil
		case c.input <- p:
		}
	}
}

func (c *generic) writeLoop(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-c.stopping:
			return nil
		case p, ok := <-c.output:
			if !ok {
				return nil
			}
			if err := c.WritePacket(p); err != nil {
				return err
			}
		}
	}
}

func (c *generic) setError(err error) {
	atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&c.err)), nil, unsafe.Pointer(&err))
}

func (c *generic) Read() (mqtt.Packet, error) {
	p, ok := <-c.input
	if !ok {
		return nil, utils.ErrClosedChannel
	}
	return p, nil
}

func (c *generic) Write(p mqtt.Packet) (err error) {
	defer func(e *error) {
		if r := recover(); r != nil {
			*e = utils.ErrClosedChannel
		}
	}(&err)
	c.output <- p
	return
}

func (c *generic) stop(err error) {
	c.stopOnce.Do(func() {
		c.setError(err)
		close(c.stopping)
	})
}

func (c *generic) Close(ctx context.Context) (err error) {
	if c.err != nil {
	}
	if c.Connection != nil {
		err = c.Connection.Close()
	}
	return
}
