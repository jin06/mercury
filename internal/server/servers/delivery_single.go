package servers

import (
	"github.com/jin06/mercury/internal/model"
	"github.com/jin06/mercury/internal/server"
	"github.com/jin06/mercury/internal/server/subscriptions"
)

func newSingle(s subscriptions.SubManager) *singleDelivery {
	return &singleDelivery{
		sub: s,
	}
}

type singleDelivery struct {
	sub    subscriptions.SubManager
	server server.Server
}

func (c *singleDelivery) Handle(msg *model.Message) (int, error) {
	subers := c.sub.GetSubers(msg.Topic)
	if len(subers) == 0 {
		return 0, nil
	}
	for _, s := range subers {
		err := c.server.Delivery(s.ClientID, msg)
		if err != nil {
		}
	}
	return 0, nil
}
