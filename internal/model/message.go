package model

import (
	"time"

	"github.com/jin06/mercury/pkg/mqtt"
)

type Message struct {
	PacketID               mqtt.PacketID       `json:"packet_id"`
	Dup                    bool                `json:"dup"`
	QoS                    uint8               `json:"qos"`
	Retained               bool                `rjson:"retained"`
	Topic                  string              `json:"topic"`
	Payload                []byte              `json:"payload"`
	ContentType            string              `json:"content_type"`
	CorrelationData        []byte              `json:"correlation_data"`
	MessageExpiry          uint32              `json:"message_expiry"`
	PayloadFormat          mqtt.PayloadFormat  `json:"payload_format"`
	ResponseTopic          string              `json:"response_topic"`
	SubscriptionIdentifier []uint32            `json:"subscription_identifier"`
	UserProperties         []mqtt.UserProperty `json:"user_properties"`
	Time                   time.Time           `json:"time"`
}
