package model

import (
	"time"

	"github.com/jin06/mercury/pkg/mqtt"
)

type Message struct {
	Publish *mqtt.Publish
	Time    time.Time
	Source  string
	Dest    string
}

func NewMessage(p *mqtt.Publish, source string, dest string) *Message {
	return &Message{
		Publish: p,
		Time:    time.Now(),
		Source:  source,
		Dest:    dest,
	}
}
