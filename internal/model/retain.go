package model

import (
	"time"

	"github.com/jin06/mercury/pkg/mqtt"
)

func NewRetain(publish mqtt.Publish) Retain {
	return Retain{
		Publish: publish,
		Time:    time.Now(),
		Expiry:  time.Minute * 60 * 24,
	}
}

type Retain struct {
	Publish mqtt.Publish
	Time    time.Time
	Expiry  time.Duration
}
