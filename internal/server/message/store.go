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
	Receive(p *mqtt.Pubrel) error
}

func NewRingBufferStore(delivery chan *Record) *ringBufferStore {
	s := &ringBufferStore{
		used:           make([]*Record, mqtt.MAXPACKETID),
		nextFreeID:     1,
		max:            mqtt.MAXPACKETID,
		expiry:         time.Hour * 24,
		delivery:       delivery,
		resendDuration: time.Second * 5,
	}
	go s.run()
	return s
}

type ringBufferStore struct {
	used           []*Record
	nextFreeID     mqtt.PacketID
	max            mqtt.PacketID
	mu             sync.Mutex
	expiry         time.Duration
	delivery       chan *Record
	closing        chan struct{}
	resendDuration time.Duration
}

func (s *ringBufferStore) Receive(p *mqtt.Pubrel) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.used[p.PacketID] != nil {
		s.used[p.PacketID].Content = p
	}
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

func (s *ringBufferStore) run() error {
	ticker := time.NewTicker(s.resendDuration)
	defer ticker.Stop()
	for {
		select {
		case <-s.closing:
			return nil
		case <-ticker.C:
			s.resend()
		default:
		}
	}
}

func (s *ringBufferStore) resend() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, record := range s.used {
		if record.Qos == mqtt.QoS1 {
			publish, ok := record.Content.(*mqtt.Publish)
			if ok {
				publish.Dup = true
				record.Content = publish
				s.delivery <- record
			}
		}
		if record.Qos == mqtt.QoS2 {
			if record.State == ReadyState {
				publish, ok := record.Content.(*mqtt.Publish)
				if ok {
					publish.Dup = true
					record.Content = publish
					s.delivery <- record
				}
			}
			if record.State == ReceivedState {
				s.delivery <- record
			}
		}
		record.Times++
	}
}

func (s *ringBufferStore) Close() error {
	close(s.closing)
	return nil
}

func (s *ringBufferStore) IsStop() bool {
	_, ok := <-s.closing
	return ok
}
