package servers

import (
	"context"
	"errors"
	"fmt"

	"github.com/jin06/mercury/internal/model"
	"github.com/jin06/mercury/internal/server"
	"github.com/jin06/mercury/internal/server/subscriptions"
	"github.com/jin06/mercury/pkg/mqtt"
)

func newGeneric() *generic {
	server := &generic{
		manager:    server.NewManager(),
		subManager: subscriptions.NewTrie(),
	}
	// server.delivery = newSingle(server.subManager)
	return server
}

type generic struct {
	manager    *server.Manager
	subManager subscriptions.SubManager
	// delivery   Delivery // remove delivery or not
}

func (g *generic) Run(ctx context.Context) error {
	select {}
}

func (g *generic) Register(c server.Client) error {
	if c == nil {
		return errors.New("client is nil")
	}
	g.manager.Set(c)
	return nil
}

func (g *generic) Deregister(c server.Client) error {
	if c == nil {
		return errors.New("client is nil")
	}
	g.manager.RemoveClient(c)
	return nil
}

func (g *generic) HandlePacket(packet mqtt.Packet, cid string) (resp mqtt.Packet, err error) {
	fmt.Println(packet)
	switch p := packet.(type) {
	case *mqtt.Connect:
		return g.HandleConnect(p)
	case *mqtt.Publish:
		return g.HandlePublish(p, cid)
	case *mqtt.Pingreq:
		return g.HandlePingreq(p, cid)
	case *mqtt.Pubrec:
		return g.HandlePubrec(p, cid)
	case *mqtt.Pubrel:
		return g.HandlePubrel(p, cid)
	case *mqtt.Pubcomp:
		return g.HandlePubcomp(p, cid)
	case *mqtt.Subscribe:
		return g.HandleSubscribe(p, cid)
	case *mqtt.Unsubscribe:
		return g.HandleUnsubscribe(p, cid)
	case *mqtt.Disconnect:
		return nil, g.HandleDisconnect(p, cid)
	case *mqtt.Auth:
		return
	}
	return
}

func (g *generic) HandleConnect(p *mqtt.Connect) (resp *mqtt.Connack, err error) {
	resp = p.Response()
	return
}

func (g *generic) HandleConnack(p *mqtt.Connack) error {
	panic("implement me")
}

func (g *generic) HandlePublish(p *mqtt.Publish, cid string) (resp mqtt.Packet, err error) {
	subers := g.subManager.GetSubers(p.Topic)
	for _, s := range subers {
		msg := &model.Message{}
		msg.FromPublish(p)
		// todo
		msg.PacketID = g.manager.GetPacketID()
		g.Delivery(s.ClientID, msg)
	}
	resp, err = p.Response()
	return
}

func (g *generic) HandlePuback(p *mqtt.Puback, cid string) (resp mqtt.Packet, err error) {
	return
}

func (g *generic) HandlePubrec(p *mqtt.Pubrec, cid string) (resp mqtt.Packet, err error) {
	return
}

func (g *generic) HandlePubrel(p *mqtt.Pubrel, cid string) (resp mqtt.Packet, err error) {
	resp = p.Response()
	return
}

func (g *generic) HandlePubcomp(p *mqtt.Pubcomp, cid string) (resp mqtt.Packet, err error) {
	return
}

func (g *generic) HandleSubscribe(p *mqtt.Subscribe, cid string) (resp *mqtt.Suback, err error) {
	for _, sub := range p.Subscriptions {
		if err = g.subManager.Sub(sub.TopicFilter, cid); err != nil {
			return nil, err
		}
	}
	resp = p.Response()
	return
}

func (g *generic) HandleSuback(p *mqtt.Suback) error {
	panic("implement me")
}

func (g *generic) HandleUnsubscribe(p *mqtt.Unsubscribe, cid string) (resp *mqtt.Unsuback, err error) {
	fmt.Println(123)
	resp = p.Response()
	return
}

func (g *generic) HandleUnsuback(p *mqtt.Unsuback, cid string) error {
	return nil
}

func (g *generic) HandlePingreq(p *mqtt.Pingreq, cid string) (resp *mqtt.Pingresp, err error) {
	resp = p.Response()
	return
}

func (g *generic) HandlePingresp(p *mqtt.Pingresp) error {
	panic("implement me")
}

func (g *generic) HandleDisconnect(p *mqtt.Disconnect, cid string) error {
	g.manager.Remove(cid)
	return nil
}

func (g *generic) HandleAuth(p *mqtt.Auth) error {
	panic("implement me")
}

//	func (g *generic) HandlePacket(p mqtt.Packet) error {
//		panic("implement me")
//	}

func (g *generic) Delivery(cid string, msg *model.Message) error {
	if client := g.manager.Get(cid); client != nil {
		return client.Write(msg.ToPublish())
	}
	return nil
}
