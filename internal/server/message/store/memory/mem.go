package memStore

import (
	"sync"
	"time"

	"github.com/jin06/mercury/internal/config"
	"github.com/jin06/mercury/internal/model"
	"github.com/jin06/mercury/internal/utils"
	"github.com/jin06/mercury/pkg/mqtt"
)

func NewMemStore(cid string, delivery chan *model.Record) *memStore {
	s := &memStore{
		cid:        cid,
		used:       make(map[mqtt.PacketID]*model.Record, mqtt.MAX_PACKET_ID),
		nextFreeID: 1,
		// max:            mqtt.MAX_PACKET_ID,
		expiry:         config.Def.MQTTConfig.MessageExpiryInterval,
		delivery:       delivery,
		resendDuration: time.Second * 5,
		closing:        make(chan struct{}),
	}
	go s.run()
	return s
}

type memStore struct {
	cid        string
	used       map[mqtt.PacketID]*model.Record
	nextFreeID mqtt.PacketID
	// max            mqtt.PacketID
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

func (s *memStore) Release(p *mqtt.Pubcomp) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.used[p.PacketID] != nil {
		s.used[p.PacketID].Content = p
	}
	return nil
}

func (s *memStore) Ack(pid mqtt.PacketID) (err error) {
	_, err = s.delete(pid)
	return
}

func (s *memStore) Complete(pid mqtt.PacketID) (err error) {
	_, err = s.delete(pid)
	return
}

func (s *memStore) Publish(p *mqtt.Publish) (*model.Record, error) {
	return s.save(p)
}

func (s *memStore) save(p *mqtt.Publish) (*model.Record, error) {
	if p.Qos.Zero() {
		return model.NewRecord(s.cid, p.Clone(), s.expiry), nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.used[s.nextFreeID] == nil {
		id := s.nextFreeID
		if s.nextFreeID++; s.nextFreeID > mqtt.MAX_PACKET_ID {
			s.nextFreeID = 1
		}
		np := p.Clone()
		np.PacketID = id
		r := model.NewRecord(s.cid, np, s.expiry)
		s.used[id] = r
		return r, nil
	}
	return nil, utils.ErrPacketIDUsed
}

func (s *memStore) delete(pid mqtt.PacketID) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var has bool
	if _, has = s.used[pid]; has {
		delete(s.used, pid)
	}
	return has, nil
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
		if publish, ok := record.Content.(*mqtt.Publish); ok {
			publish.Dup = true
			record.Content = publish
			s.delivery <- record
		}
		s.delivery <- record
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
