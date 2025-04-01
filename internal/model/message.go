package model

import (
	"time"

	"github.com/jin06/mercury/pkg/mqtt"
)

type MessageState byte

const (
	ReadyState MessageState = 0
	PubState   MessageState = 1
	RecState   MessageState = 1
	AckState   MessageState = 1
)

type Message struct {
	Publish *mqtt.Publish
	Time    time.Time
	State   MessageState
	Source  string
	Dest    string
}

func NewMessage(p *mqtt.Publish, source string, dest string) *Message {
	return &Message{
		Publish: p,
		Time:    time.Now(),
		Source:  source,
		Dest:    dest,
		State:   ReadyState,
	}
}
