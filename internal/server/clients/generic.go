package clients

import (
	"context"
	"net"
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/jin06/mercury/internal/logs"
	"github.com/jin06/mercury/internal/server"
	"github.com/jin06/mercury/internal/server/connections"
	"github.com/jin06/mercury/internal/utils"
	"github.com/jin06/mercury/pkg/mqtt"
)

func NewClient(handler server.Server, conn net.Conn) *generic {
	c := generic{
		handler:   handler,
		conn:      connections.NewTCP(conn),
		connected: make(chan struct{}),
		stopping:  make(chan struct{}),
		closed:    make(chan struct{}),
		options:   Options{},
	}
	return &c
}

type generic struct {
	id        string
	conn      connections.Connection
	handler   server.Server
	options   Options
	connected chan struct{}
	stopping  chan struct{}
	stopOnce  sync.Once
	closed    chan struct{}
	err       error // first error that occurs exits the client
	// packet channels
	input  chan *mqtt.Packet
	output chan *mqtt.Packet

	// client info
	ip string // client ip
}

func (c *generic) ClientID() string {
	return c.id
}

func (c *generic) Run(ctx context.Context) (err error) {
	defer func() {
		if err != nil {
			c.setError(err)
			c.stop(err)
		}
		c.Close(ctx)
	}()

	defer close(c.closed)

	c.handler.Reg(c)

	var p mqtt.Packet

	if p, err = c.conn.Read(); err != nil {
		return err
	}

	if packet, ok := p.(*mqtt.Connect); ok {
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
	select {
	case <-ctx.Done():
		return nil
	case <-c.stopping:
		return nil
	}
	wg.Done()
	return nil
}

func (c *generic) readloop(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-c.stopping:
			return nil
		default:
			p, err := c.conn.Read()
			if err != nil {
				return err
			}
			c.input <- &p
		}
	}
}

func (c *generic) writeloop(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-c.stopping:
			return nil
		case p := <-c.output:
			if err := c.Write(*p); err != nil {
				logs.Logger.Err(err)
				return err
			}
		}
	}
}

func (c *generic) setError(err error) {
	atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&c.err)), nil, unsafe.Pointer(&err))
}

// func (c *generic) handlePacket(p mqtt.Packet) {

// 	switch p.(type) {
// 	case *mqtt.Connect:
// 		{
// 			_ = &mqtt.Connack{
// 				FixHeader: &mqtt.FixedHeader{
// 					PacketType:      mqtt.CONNACK,
// 					Flags:           0,
// 					RemainingLength: 2,
// 				},
// 			}
// 		}
// 	}
// }

func (c *generic) Write(p mqtt.Packet) error {
	return nil
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
	if c.conn != nil {
		err = c.conn.Close()
	}
	return
}
