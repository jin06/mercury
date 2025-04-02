package model

import (
	"time"

	"github.com/jin06/mercury/pkg/mqtt"
)

type MessageState byte

const (
	ReadyState    MessageState = 0
	PubState      MessageState = 1
	ReceivedState MessageState = 2
	ReleaseState  MessageState = 3
	AckState      MessageState = 4
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
