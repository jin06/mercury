package clients

import (
	"sync"
	"time"

	"github.com/jin06/mercury/pkg/mqtt"
)

func newRecordDB() *recordDB {
	return &recordDB{
		records: make(map[mqtt.PacketID]*pubRecord),
		closing: make(chan struct{}),
	}
}

type recordDB struct {
	mu      sync.RWMutex
	records map[mqtt.PacketID]*pubRecord
	closing chan struct{}
}

func (db *recordDB) iter(f func(mqtt.Packet) error) error {
	db.mu.Lock()
	defer db.mu.Unlock()
	for _, record := range db.records {
		if err := f(record.response); err != nil {
			return err
		}
	}
	return nil
}

func (db *recordDB) close() {
	close(db.closing)
}

func (db *recordDB) save(p *mqtt.Publish, response mqtt.Packet) {
	if p.Qos != mqtt.QoS2 {
		return
	}
	db.mu.Lock()
	defer db.mu.Unlock()
	_, ok := db.records[p.ID()]
	if !ok {
		db.records[p.ID()] = newRecord(p, response)
	}
}
func (db *recordDB) delivery(id mqtt.PacketID, f func(*mqtt.Publish) error) (err error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	if record, ok := db.records[id]; ok {
		if err = f(record.content); err != nil {
			return
		}
	}
	delete(db.records, id)
	return
}

func newRecord(p *mqtt.Publish, response mqtt.Packet) *pubRecord {
	return &pubRecord{
		content:  p,
		response: response,
		receive:  time.Now(),
		send:     time.Now(),
		times:    1,
	}
}

type pubRecord struct {
	content  *mqtt.Publish
	response mqtt.Packet
	receive  time.Time
	send     time.Time
	times    int
}
