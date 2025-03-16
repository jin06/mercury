package servers

import (
	"context"
	"errors"

	"github.com/jin06/mercury/internal/model"
	"github.com/jin06/mercury/internal/server"
	"github.com/jin06/mercury/internal/server/subscriptions"
	"github.com/jin06/mercury/pkg/mqtt"
)

func newGeneric() *generic {
	return &generic{
		manager:    server.NewManager(),
		subManager: subscriptions.NewTrie(),
	}
}

type generic struct {
	manager    *server.Manager
	subManager subscriptions.SubManager
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
	switch p := packet.(type) {
	case *mqtt.Connect:
		return g.handleConnect(p)
	case *mqtt.Pingreq:
		return nil, nil
	case *mqtt.Subscribe:
		return g.HandleSubscribe(p, cid)
	case *mqtt.Disconnect:
	}
	return
}

func (g *generic) handleConnect(p *mqtt.Connect) (resp *mqtt.Connack, err error) {
	resp = p.Response()
	return
}

func (g *generic) HandleConnack(p *mqtt.Connack) error {
	panic("implement me")
}

func (g *generic) HandlePublish(p *mqtt.Publish) error {
	subers := g.subManager.GetSubers(p.Topic)
	for _, s := range subers {
		msg := &model.Message{}
		msg.FromPublish(p)
		// todo
		msg.PacketID = g.manager.GetPacketID()
		g.Delivery(s.ClientID, msg)
	}
	return nil
}

func (g *generic) HandlePuback(p *mqtt.Puback) error {
	panic("implement me")
}

func (g *generic) HandlePubrec(p *mqtt.Pubrec) error {
	panic("implement me")
}

func (g *generic) HandlePubrel(p *mqtt.Pubrel) error {
	panic("implement me")
}

func (g *generic) HandlePubcomp(p *mqtt.Pubcomp) error {
	panic("implement me")
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

func (g *generic) HandleUnsubscribe(p *mqtt.Unsubscribe) error {
	panic("implement me")
}

func (g *generic) HandleUnsuback(p *mqtt.Unsuback) error {
	panic("implement me")
}

func (g *generic) HandlePingreq(p *mqtt.Pingreq) error {
	panic("implement me")
}

func (g *generic) HandlePingresp(p *mqtt.Pingresp) error {
	panic("implement me")
}

func (g *generic) HandleDisconnect(p *mqtt.Disconnect) error {
	panic("implement me")
}

func (g *generic) HandleAuth(p *mqtt.Auth) error {
	panic("implement me")
}

//	func (g *generic) HandlePacket(p mqtt.Packet) error {
//		panic("implement me")
//	}
func (g *generic) Delivery(cid string, msg *model.Message) error {
	client := g.manager.Get(cid)
	if client != nil {
	}
	return nil
}
