package msgBadger

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"

	"github.com/dgraph-io/badger/v4"
	"github.com/jin06/mercury/internal/config"
	"github.com/jin06/mercury/internal/model"
	"github.com/jin06/mercury/pkg/mqtt"
)

var (
	def         *badger.DB
	packetIDKey = "packetid:%s"  // -> packetid:{clientID}
	recordKey   = "record:%s:%d" // -> record:{clientID}:{PacketID}
)

func Init(options config.BadgerConfig) (err error) {
	def, err = badger.Open(badger.DefaultOptions(options.Dir))
	return
}

func NewBadgerStore(cid string, options config.BadgerConfig) (*badgerStore, error) {
	s := &badgerStore{
		options: options,
		db:      def,
		cid:     cid,
	}
	return s, nil
}

type badgerStore struct {
	options config.BadgerConfig
	db      *badger.DB
	cid     string
}

func (store *badgerStore) Publish(p *mqtt.Publish) (*model.Record, error) {
	return store.create(p)
}

func (store *badgerStore) Ack(pid mqtt.PacketID) error {
	return store.delete(pid)
}

func (store *badgerStore) Receive(p *mqtt.Pubrel) error {
	return store.update(p)
}

func (store *badgerStore) Release(p *mqtt.Pubcomp) error {
	return store.update(p)
}

func (store *badgerStore) Complete(pid mqtt.PacketID) error {
	return store.delete(pid)
}

func (store *badgerStore) create(p *mqtt.Publish) (record *model.Record, err error) {
	if p.Qos.Zero() {
		return model.NewRecord(p.Clone()), nil
	}
	err = store.db.Update(func(txn *badger.Txn) (err error) {
		item, err := txn.Get([]byte(store.getPacketIDKey()))
		if err != nil {
			return
		}
		val, err := item.ValueCopy(nil)
		if err != nil {
			return
		}
		currentID := mqtt.PacketID(binary.BigEndian.Uint16(val))

		if currentID == 0 {
			currentID = 1
		}
		switch {
		case currentID == 0:
			currentID = 1
		case currentID > mqtt.MAXPACKETID:
			currentID = 1
		}
		nextID := currentID + 1

		var buf bytes.Buffer
		encoder := gob.NewEncoder(&buf)
		np := p.Clone()
		np.PacketID = currentID
		r := model.NewRecord(np)
		if err = encoder.Encode(r); err != nil {
			return
		}
		if err = txn.Set([]byte(store.getRecordKey(currentID)), buf.Bytes()); err != nil {
			return
		}
		if err = txn.Set([]byte(store.getPacketIDKey()), nextID.Encode()); err != nil {
			return
		}
		return
	})
	return
}

func (store *badgerStore) update(p mqtt.Message) error {
	return nil
}

func (store *badgerStore) delete(pid mqtt.PacketID) error {
	return store.db.DropPrefix([]byte(store.getRecordKey(pid)))
}

func (store *badgerStore) getPacketIDKey() string {
	return fmt.Sprintf(packetIDKey, store.cid)
}

func (store *badgerStore) getRecordKey(packetID mqtt.PacketID) string {
	return fmt.Sprintf(recordKey, store.cid, packetID)
}
