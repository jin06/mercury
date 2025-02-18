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
		closing:   make(chan struct{}),
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
	closing   chan struct{}
	closed    chan struct{}
	closeOnce sync.Once
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
		if r := recover(); r != nil {
		}
		if err != nil {
			logs.Logger.Err(err)
			c.setError(err)
		}
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
		err = utils.ErrNotConnectPacket
		return
	}
	c.runloop(ctx)
	return nil
}

func (c *generic) runloop(ctx context.Context) {
	wg := sync.WaitGroup{}
	select {
	case <-ctx.Done():
		return
	case <-c.closing:
		return
	case <-c.closed:
		return
	}
	wg.Done()
}

func (c *generic) readloop(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-c.closing:
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
		case <-c.closing:
			return nil
		case p := <-c.output:
			if err := c.Write(*p); err != nil {
				logs.Logger.Err(err)
				return err
			}
		}
	}
	return nil
}

func (c *generic) setError(err error) {
	atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&c.err)), nil, unsafe.Pointer(&err))
}

func (c *generic) handlePacket(p mqtt.Packet) {

	switch p.(type) {
	case *mqtt.Connect:
		{
			_ = &mqtt.Connack{
				FixHeader: &mqtt.FixedHeader{
					PacketType:      mqtt.CONNACK,
					Flags:           0,
					RemainingLength: 2,
				},
			}
		}
	}
}

func (c *generic) Write(p mqtt.Packet) error {
	return nil
}

func (c *generic) Close(ctx context.Context) error {
	var err error
	c.closeOnce.Do(func() {
		close(c.closing)
	})
	return err
}
