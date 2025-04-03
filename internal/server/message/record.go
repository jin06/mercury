package message

import (
	"time"

	"github.com/jin06/mercury/pkg/mqtt"
)

type State byte

const (
	ReadyState     State = 0
	PublishedState State = 1
	ReceivedState  State = 2
	ReleasedState  State = 3
	AckState       State = 4
)

type Record struct {
	Content *mqtt.Publish
	Receive time.Time
	Send    time.Time
	State   State
	Source  string
	Dest    string
	Times   int
}

func NewRecord(p *mqtt.Publish, source string, dest string) *Record {
	r := &Record{
		Content: p,
		Receive: time.Now(),
		Send:    time.Now(),
		State:   ReadyState,
		Source:  source,
		Dest:    dest,
		Times:   0,
	}
	return r
}
