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
	Qos      mqtt.QoS
	Receive  time.Time
	Send     time.Time
	Expiry   time.Duration
	ClientID string
	State    State
	Times    uint64
	Content  mqtt.Packet
}

func NewRecord(cid string, p mqtt.Packet) *Record {
	r := &Record{
		Content:  p,
		Receive:  time.Now(),
		Send:     time.Now(),
		State:    ReadyState,
		Times:    0,
		ClientID: cid,
	}
	return r
}
