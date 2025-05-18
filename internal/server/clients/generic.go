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
	"github.com/jin06/mercury/internal/config"
	"github.com/jin06/mercury/internal/server"
	"github.com/jin06/mercury/internal/server/message/store"
	"github.com/jin06/mercury/internal/utils"
	"github.com/jin06/mercury/pkg/mqtt"
)

func NewClient(handler server.Server, conn io.ReadWriteCloser) *generic {
	c := generic{
		handler:    handler,
		Connection: mqtt.NewConnection(conn),
		stopping:   make(chan struct{}),
		closed:     make(chan struct{}),
		options:    DefaultOptions(),
		input:      make(chan mqtt.Packet, 2000),
		output:     make(chan mqtt.Packet, 2000),
		uuid:       uuid.New().String(),
		keep:       time.Now(),
		db:         newRecordDB(),
	}
	return &c
}

type generic struct {
	id string
	*mqtt.Connection
	handler   server.Server
	connected bool
	options   *Options
	stopping  chan struct{}
	stopOnce  sync.Once
	closed    chan struct{}
	disOnce   sync.Once
	err       error // first error that occurs exits the client
	// packet channels
	input         chan mqtt.Packet
	output        chan mqtt.Packet
	uuid          string
	keep          time.Time
	db            *recordDB
	connectedTime time.Time
	msgStore      store.Store
	cleanSession  bool
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

	c.Reader.Version = cp.Version
	c.id = cp.ClientID
	c.cleanSession = cp.Clean

	c.msgStore = store.NewStore(config.Def.MessageStore.Mode, c.id)

	fmt.Printf("[IN] - [%s] | %v \n", cp.ClientID, cp)

	if response, err = c.handler.HandleConnect(cp, c); err != nil {
		return
	}

	if err = c.Write(response); err != nil {
		return
	}

	c.connected = true
	c.connectedTime = time.Now()

	return nil
}

func (c *generic) disconnect(p *mqtt.Disconnect) (err error) {
	c.disOnce.Do(func() {
		c.Write(p)
	})
	return
}

func (c *generic) runloop(ctx context.Context) error {
	go func() {
		if err := c.inputLoop(ctx); err != nil {
			c.stop(err)
			return
		}
	}()
	go func() {
		if err := c.outputLoop(ctx); err != nil {
			c.stop(err)
			return
		}
	}()
	go func() {
		if err := c.handleLoop(ctx); err != nil {
			c.stop(err)
			return
		}
	}()
	go func() {
		if err := c.keepLoop(ctx); err != nil {
			c.stop(err)
			return
		}
	}()
	go func() {
		if err := c.recordLoop(ctx); err != nil {
			c.stop(err)
			return
		}
	}()
	go func() {
		if err := c.msgStore.Run(ctx, c.output); err != nil {
			return
		}
	}()
	<-c.stopping
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
			fmt.Printf("[OUT] - [%s] | %v \n", c.id, p)
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
			fmt.Printf("[IN] - [ClientID: %s] | %v \n", c.id, p)
			switch val := p.(type) {
			case *mqtt.Pingreq:
				resp, err = c.handler.HandlePacket(val, c.id)
			case *mqtt.Publish:
				resp, err = c.handler.HandlePacket(val, c.id)
				if err == nil {
					if val.Qos == mqtt.QoS2 {
						c.db.save(val, resp)
					}
				}
			case *mqtt.Puback:
				resp, err = c.handler.HandlePacket(val, c.id)
			case *mqtt.Pubrec:
				resp, err = c.handler.HandlePacket(val, c.id)
			case *mqtt.Pubrel:
				if resp, err = c.handler.HandlePacket(val, c.id); err != nil {
					break
				}
				err = c.db.dispatch(val.PacketID, func(p *mqtt.Publish) error {
					return c.handler.Dispatch(c.id, p)
				})
			case *mqtt.Pubcomp:
				resp, err = c.handler.HandlePacket(val, c.id)
			case *mqtt.Subscribe:
				resp, err = c.handler.HandlePacket(val, c.id)
			case *mqtt.Unsubscribe:
				resp, err = c.handler.HandlePacket(val, c.id)
			case *mqtt.Disconnect:
				resp, err = c.handler.HandlePacket(val, c.id)
			case *mqtt.Auth:
			}
		}
		if err != nil {
			fmt.Println(err)
		}
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
			c.disconnect(mqtt.NewDisconnect(&mqtt.FixedHeader{}, c.Version))
		}
	}
	if c.Connection != nil {
		c.Connection.Close()
	}
	if c.msgStore != nil {
		c.msgStore.Close()
		c.msgStore.Clean()
	}
	err = c.handler.Deregister(c)
	return
}

func (c *generic) recordLoop(ctx context.Context) error {
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			c.db.iter(c.Write)
		case <-c.stopping:
			return nil
		case <-c.closed:
			return nil
		}
	}
}
