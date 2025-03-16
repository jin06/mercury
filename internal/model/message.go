package model

import (
	"time"

	"github.com/jin06/mercury/pkg/mqtt"
)

type Message struct {
	PacketID               mqtt.PacketID       `json:"packet_id"`
	Dup                    bool                `json:"dup"`
	QoS                    mqtt.QoS            `json:"qos"`
	Retained               bool                `json:"retained"`
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

func (m *Message) FromPublish(p *mqtt.Publish) {
	m.Dup = p.Dup
	m.QoS = p.Qos
	m.Retained = p.Retain
	m.Topic = p.Topic
	m.Payload = p.Payload
	m.ContentType = p.ContentType
	m.CorrelationData = p.CorrelationData
	m.MessageExpiry = uint32(p.MessageExpiryInterval)
	m.PayloadFormat = *p.Properties.PayloadFormat
	m.ResponseTopic = p.ResponseTopic
	m.SubscriptionIdentifier = p.SubscriptionIdentifier
	m.UserProperties = p.Properties.UserProperties
	m.Time = time.Now()

}

// todo
func (m *Message) ToPublish() *mqtt.Publish {
	return &mqtt.Publish{}
}
