package model

import (
	"time"

	"github.com/jin06/mercury/pkg/mqtt"
)

type State byte

const (
	ReadyState State = 0
	// PublishedState State = 1
	ReceivedState State = 2
	ReleasedState State = 3
	// AckState State = 4
)

type Record struct {
	Receive  time.Time
	Send     time.Time
	Expiry   time.Duration
	Version  mqtt.ProtocolVersion
	ClientID string
	Times    uint64
	Content  mqtt.Packet
}

func NewRecord(cid string, p mqtt.Message, expiry time.Duration) *Record {
	r := &Record{
		Content:  p,
		Expiry:   expiry,
		Receive:  time.Now(),
		Send:     time.Now(),
		Version:  p.GetVersion(),
		Times:    0,
		ClientID: cid,
	}
	return r
}
