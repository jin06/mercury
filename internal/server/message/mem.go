package message

import (
	"sync"
	"time"

	"github.com/jin06/mercury/internal/model"
	"github.com/jin06/mercury/internal/utils"
	"github.com/jin06/mercury/pkg/mqtt"
)

func NewMemStore(delivery chan *model.Record) *memStore {
	s := &memStore{
		used:           make(map[mqtt.PacketID]*model.Record, mqtt.MAXPACKETID),
		nextFreeID:     1,
		max:            mqtt.MAXPACKETID,
		expiry:         time.Hour * 24,
		delivery:       delivery,
		resendDuration: time.Second * 5,
	}
	go s.run()
	return s
}

type memStore struct {
	used           map[mqtt.PacketID]*model.Record
	nextFreeID     mqtt.PacketID
	max            mqtt.PacketID
	mu             sync.Mutex
	expiry         time.Duration
	delivery       chan *model.Record
	closing        chan struct{}
	resendDuration time.Duration
}

func (s *memStore) Receive(p *mqtt.Pubrel) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.used[p.PacketID] != nil {
		s.used[p.PacketID].Content = p
	}
	return nil
}

func (s *memStore) Save(p *mqtt.Publish, source string, dest string) (*model.Record, error) {
	if p.Qos.Zero() {
		return model.NewRecord(p.Clone(), source, dest), nil
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
		r := model.NewRecord(np, source, dest)
		s.used[id] = r
		return r, nil
	}
	return nil, utils.ErrPacketIDUsed
}

func (s *memStore) Delete(pid mqtt.PacketID) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var has bool
	if _, has = s.used[pid]; has {
		delete(s.used, pid)
	}
	return has, nil
}

func (s *memStore) Change(id mqtt.PacketID, state model.State) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.used[id] == nil {
		return utils.ErrPacketIDNotExist
	}
	s.used[id].State = state
	return nil
}

func (s *memStore) run() error {
	ticker := time.NewTicker(s.resendDuration)
	defer ticker.Stop()
	for {
		select {
		case <-s.closing:
			return nil
		case <-ticker.C:
			s.resend()
		}
	}
}

func (s *memStore) resend() {
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
			if record.State == model.ReadyState {
				publish, ok := record.Content.(*mqtt.Publish)
				if ok {
					publish.Dup = true
					record.Content = publish
					s.delivery <- record
				}
			}
			if record.State == model.ReceivedState {
				s.delivery <- record
			}
		}
		record.Times++
	}
}

func (s *memStore) Close() error {
	close(s.closing)
	return nil
}

func (s *memStore) IsStop() bool {
	_, ok := <-s.closing
	return ok
}
