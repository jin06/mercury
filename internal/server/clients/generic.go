package clients

import (
	"context"
	"fmt"
	"io"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/google/uuid"
	"github.com/jin06/mercury/internal/server"
	"github.com/jin06/mercury/internal/utils"
	"github.com/jin06/mercury/pkg/mqtt"
)

func NewClient(handler server.Server, conn io.ReadWriteCloser) *generic {
	c := generic{
		handler:    handler,
		Connection: mqtt.NewConnection(conn),
		stopping:   make(chan struct{}),
		closed:     make(chan struct{}),
		options:    Options{},
		input:      make(chan mqtt.Packet, 2000),
		output:     make(chan mqtt.Packet, 2000),
		uuid:       uuid.New().String(),
		keep:       time.Now(),
	}
	return &c
}

type generic struct {
	id string
	*mqtt.Connection
	handler   server.Server
	connected bool
	options   Options
	stopping  chan struct{}
	stopOnce  sync.Once
	closed    chan struct{}
	disOnce   sync.Once
	err       error // first error that occurs exits the client
	// packet channels
	input  chan mqtt.Packet
	output chan mqtt.Packet
	uuid   string
	keep   time.Time
}

func (c *generic) ClientID() string {
	return c.id
}

func (c *generic) UUID() string {
	return c.uuid
}

func (c *generic) Run(ctx context.Context) (err error) {
	defer close(c.closed)
	defer c.Close(ctx)
	defer c.stop(err)

	if err = c.connect(); err != nil {
		return
	}

	if err = c.runloop(ctx); err != nil {
		return
	}

	return
}

func (c *generic) connect() (err error) {
	var p mqtt.Packet
	var response mqtt.Packet

	if p, err = c.ReadPacket(); err != nil {
		return
	}
	cp, ok := p.(*mqtt.Connect)
	if !ok {
		return utils.ErrMalformedPacket
	}

	if response, err = c.handler.HandleConnect(cp); err != nil {
		return
	}

	if err = c.Write(response); err != nil {
		return
	}
	if err = c.handler.Register(c); err != nil {
		return
	}
	c.connected = true
	c.id = cp.ClientID
	c.Reader.Version = cp.Version
	return nil
}

func (c *generic) disconnect(p *mqtt.Disconnect) (err error) {
	c.disOnce.Do(func() {
		c.Write(p)
	})
	return
}

func (c *generic) runloop(ctx context.Context) error {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := c.inputLoop(ctx); err != nil {
			c.stop(err)
			return
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := c.outputLoop(ctx); err != nil {
			c.stop(err)
			return
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := c.handleLoop(ctx); err != nil {
			c.stop(err)
			return
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := c.keepLoop(ctx); err != nil {
			c.stop(err)
			return
		}
	}()
	wg.Wait()
	return nil
}

func (c *generic) inputLoop(ctx context.Context) error {
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
		default:
			if p != nil {
				c.input <- p
			}
		}
	}
}

func (c *generic) outputLoop(ctx context.Context) error {
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

func (c *generic) handleLoop(ctx context.Context) error {
	for {
		var resp mqtt.Packet
		var err error
		select {
		case <-ctx.Done():
			return nil
		case <-c.stopping:
			return nil
		case p, ok := <-c.input:
			if !ok {
				return nil
			}
			switch val := p.(type) {
			case *mqtt.Pingreq:
				// resp = val.Response()
				resp, err = c.handler.HandlePacket(val, c.id)
			case *mqtt.Publish:
				resp, err = c.handler.HandlePacket(val, c.id)
			case *mqtt.Pubrec:
			case *mqtt.Pubrel:
			case *mqtt.Pubcomp:
			case *mqtt.Subscribe:
				resp, err = c.handler.HandlePacket(val, c.id)
			case *mqtt.Unsubscribe:
				fmt.Println("unsubscribe", val)
				resp, err = c.handler.HandlePacket(val, c.id)
			case *mqtt.Disconnect:
			case *mqtt.Auth:
			}
		}
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("%v\n", resp)
		c.KeepAlive()
		if resp != nil {
			c.Write(resp)
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
	// defer func(e *error) {
	// 	if r := recover(); r != nil {
	// 		*e = utils.ErrClosedChannel
	// 	}
	// }(&err)
	c.output <- p
	return
}

func (c *generic) stop(err error) {
	c.stopOnce.Do(func() {
		c.setError(err)
		close(c.stopping)
	})
}

func (c *generic) keepLoop(ctx context.Context) error {
	// ticker := time.NewTicker(time.Minute)
	ch := make(chan time.Time, 1)
	defer close(ch)
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-c.stopping:
			return nil
		case <-ch:
			duration := time.Since(c.keep)
			if duration > time.Hour {
				panic("keep alive")
			}
		}
	}
}

func (c *generic) KeepAlive() {
	c.keep = time.Now()
}

func (c *generic) Close(ctx context.Context) (err error) {
	if c.connected {
		if c.err != nil {
			c.disconnect(mqtt.NewDisconnect(&mqtt.FixedHeader{}))
		}
	}
	if c.Connection != nil {
		c.Connection.Close()
	}
	err = c.handler.Deregister(c)
	return
}
