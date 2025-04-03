package message

import (
	"sync"
	"time"

	"github.com/jin06/mercury/internal/utils"
	"github.com/jin06/mercury/pkg/mqtt"
)

type Store interface {
	Save(p *mqtt.Publish, source string, dest string) (*Record, error)
	Delete(pid mqtt.PacketID) (bool, error)
	Change(id mqtt.PacketID, state State) error
}

func NewRingBufferStore() *ringBufferStore {
	return &ringBufferStore{
		used:       make([]*Record, mqtt.MAXPACKETID),
		nextFreeID: 1,
		max:        mqtt.MAXPACKETID,
		expiry:     time.Hour * 24,
	}
}

type ringBufferStore struct {
	used       []*Record
	nextFreeID mqtt.PacketID
	max        mqtt.PacketID
	mu         sync.Mutex
	expiry     time.Duration
}

func (s *ringBufferStore) run() error {
	defer s.close()
	return nil
}

func (s *ringBufferStore) close() error {
	return nil
}

func (s *ringBufferStore) Save(p *mqtt.Publish, source string, dest string) (*Record, error) {
	if p.Qos.Zero() {
		return nil, nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.used[s.nextFreeID] == nil {
		id := s.nextFreeID
		if s.nextFreeID++; s.nextFreeID > s.max {
			s.nextFreeID = 1
		}
		np := p.Clone()
		np.PacketID = id
		r := NewRecord(np, source, dest)
		s.used[id] = NewRecord(np, source, dest)
		return r, nil
	}
	return nil, utils.ErrPacketIDUsed
}

func (s *ringBufferStore) Delete(pid mqtt.PacketID) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var has bool
	if s.used[pid] != nil {
		s.used[pid] = nil
		has = true
	}
	return has, nil
}

func (s *ringBufferStore) Change(id mqtt.PacketID, state State) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.used[id] == nil {
		return utils.ErrPacketIDNotExist
	}
	s.used[id].State = state
	return nil
}
