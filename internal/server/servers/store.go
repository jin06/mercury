package servers

import (
	"sync"
	"time"

	"github.com/jin06/mercury/internal/model"
	"github.com/jin06/mercury/internal/utils"
	"github.com/jin06/mercury/pkg/mqtt"
)

type Store interface {
}

func newRingBufferStore() *ringBufferStore {
	return &ringBufferStore{
		used:       make([]*model.Message, mqtt.MAXPACKETID),
		nextFreeID: 1,
		max:        mqtt.MAXPACKETID,
		expiry:     time.Hour * 24,
	}
}

type ringBufferStore struct {
	used       []*model.Message
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

func (s *ringBufferStore) PopPacketID(p *model.Message) (mqtt.PacketID, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.used[s.nextFreeID] == nil { // 说明该ID尚未使用
		id := s.nextFreeID
		s.nextFreeID++
		if s.nextFreeID > s.max {
			s.nextFreeID = 1
		}
		s.used[id] = p
		return id, nil
	}
	return 0, utils.ErrPacketIDUsed
}

func (s *ringBufferStore) DiscardExpiry() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()

	for id, msg := range s.used {
		if msg == nil {
			continue
		}
		expiry := s.expiry
		if msg.Publish.Properties != nil && msg.Publish.Properties.MessageExpiryInterval != nil {
			if *msg.Publish.Properties.MessageExpiryInterval != 0 {
				expiry = time.Duration(*msg.Publish.Properties.MessageExpiryInterval) * time.Second
			}
		}
		if now.Unix()-msg.Time.Unix() > int64(expiry.Seconds()) {
			s.used[id] = nil
		}
	}
	return nil
}
