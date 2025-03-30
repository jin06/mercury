package model

import (
	"time"

	"github.com/jin06/mercury/pkg/mqtt"
)

type Message struct {
	PacketID               mqtt.PacketID             `json:"packet_id"`
	Version                mqtt.ProtocolVersion      `json:"version"`
	Dup                    bool                      `json:"dup"`
	QoS                    mqtt.QoS                  `json:"qos"`
	Retained               bool                      `json:"retained"`
	Topic                  string                    `json:"topic"`
	Payload                []byte                    `json:"payload"`
	ContentType            *string                   `json:"content_type"`
	CorrelationData        *mqtt.BinaryData          `json:"correlation_data"`
	MessageExpiry          *uint32                   `json:"message_expiry"`
	PayloadFormat          *mqtt.PayloadFormat       `json:"payload_format"`
	ResponseTopic          *string                   `json:"response_topic"`
	SubscriptionIdentifier *mqtt.VariableByteInteger `json:"subscription_identifier"`
	UserProperties         []mqtt.UserProperty       `json:"user_properties"`
	Time                   time.Time                 `json:"time"`
}

func (m *Message) FromPublish(p *mqtt.Publish) {
	m.Dup = p.Dup
	m.QoS = p.Qos
	m.Retained = p.Retain
	m.Topic = p.Topic.String()
	m.Payload = p.Payload
	m.Version = p.Version
	m.PacketID = p.PacketID
	if p.Properties != nil {
		m.PayloadFormat = p.Properties.PayloadFormat
		m.ContentType = p.Properties.ContentType
		m.CorrelationData = p.Properties.CorrelationData
		m.MessageExpiry = p.Properties.MessageExpiryInterval
		m.ResponseTopic = p.Properties.ResponseTopic
		m.SubscriptionIdentifier = p.Properties.SubscriptionIdentifier
		m.UserProperties = p.Properties.UserProperties
	}
	m.Time = time.Now()

}

// todo
func (m *Message) ToPublish() *mqtt.Publish {
	header := mqtt.FixedHeader{
		PacketType: mqtt.PUBLISH,
	}
	p := mqtt.NewPublish(&header, m.Version)
	p.Dup = m.Dup
	p.Qos = m.QoS
	p.Retain = m.Retained
	p.Topic = mqtt.UTF8String(m.Topic)
	p.Payload = m.Payload
	p.Version = m.Version
	p.PacketID = m.PacketID
	p.Properties = &mqtt.Properties{
		PayloadFormat:          m.PayloadFormat,
		UserProperties:         m.UserProperties,
		ContentType:            m.ContentType,
		CorrelationData:        m.CorrelationData,
		MessageExpiryInterval:  m.MessageExpiry,
		ResponseTopic:          m.ResponseTopic,
		SubscriptionIdentifier: m.SubscriptionIdentifier,
	}
	return p

}
