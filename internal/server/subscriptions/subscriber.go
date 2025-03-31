package subscriptions

import (
	"time"

	"github.com/jin06/mercury/pkg/mqtt"
)

type Subscriber struct {
	Type     Type
	ClientID string
	Group    string
	Time     time.Time

	RetainAsPublished bool
	NoLocal           bool

	RetainHandling byte
	Qos            mqtt.QoS
}
