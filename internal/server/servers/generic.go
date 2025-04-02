package servers

import (
	"context"
	"errors"

	"github.com/jin06/mercury/internal/model"
	"github.com/jin06/mercury/internal/server"
	"github.com/jin06/mercury/internal/server/message"
	"github.com/jin06/mercury/internal/server/subscriptions"
	"github.com/jin06/mercury/pkg/mqtt"
)

func newGeneric() *generic {
	server := &generic{
		manager:    server.NewManager(),
		subManager: subscriptions.NewTrie(),
		msgManager: *message.NewManager(),
	}
	return server
}

type generic struct {
	manager    *server.Manager
	subManager subscriptions.SubManager
	msgManager message.Manager
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

func (g *generic) HandleConnect(p *mqtt.Connect, c server.Client) (resp *mqtt.Connack, err error) {
	if err = g.Register(c); err != nil {
		return
	}
	resp = p.Response()
	return
}

func (g *generic) HandleConnack(p *mqtt.Connack) error {
	panic("implement me")
}

func (g *generic) HandlePublish(p *mqtt.Publish, cid string) (resp mqtt.Packet, err error) {
	if p.Qos == mqtt.QoS0 || p.Qos == mqtt.QoS1 {
		if err = g.DeliveryPublish(cid, p); err != nil {
			return
		}
	}
	resp, err = p.Response()
	return
}

func (g *generic) HandlePuback(p *mqtt.Puback, cid string) (resp mqtt.Packet, err error) {
	_, err = g.msgManager.Ack(cid, p.ID())
	return
}

func (g *generic) HandlePubrec(p *mqtt.Pubrec, cid string) (resp mqtt.Packet, err error) {
	_, err = g.msgManager.Rec(cid, p.PacketID)
	resp = p.Response()
	return
}

func (g *generic) HandlePubrel(p *mqtt.Pubrel, cid string) (resp mqtt.Packet, err error) {
	resp = p.Response()
	return
}

func (g *generic) HandlePubcomp(p *mqtt.Pubcomp, cid string) (resp mqtt.Packet, err error) {
	_, err = g.msgManager.Ack(cid, p.ID())
	return
}

func (g *generic) HandleSubscribe(p *mqtt.Subscribe, cid string) (resp *mqtt.Suback, err error) {
	for _, sub := range p.Subscriptions {
		if _, err = g.subManager.Sub(sub.TopicFilter, cid); err != nil {
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
	for _, v := range p.TopicFilters {
		g.subManager.Unsub(v, cid)
	}
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
	return nil
}

func (g *generic) HandleDisconnect(p *mqtt.Disconnect, cid string) error {
	g.manager.Remove(cid)
	return nil
}

func (g *generic) HandleAuth(p *mqtt.Auth) error {
	panic("implement me")
}

func (g *generic) DeliveryPublish(cid string, p *mqtt.Publish) error {
	subers := g.subManager.GetSubers(p.Topic.String())
	for _, s := range subers {
		msg := model.NewMessage(p, cid, s.ClientID)
		msg.Publish.PacketID = g.manager.GetPacketID()
		if _, err := g.msgManager.Pop(msg); err != nil {
			return err
		}
		go g.DeliveryOne(s.ClientID, msg)
	}
	return nil
}

func (g *generic) DeliveryOne(cid string, msg *model.Message) error {
	if client := g.manager.Get(cid); client != nil {
		return client.Write(msg.Publish)
	}
	return nil
}
