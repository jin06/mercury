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

func NewBadgerStore(options config.BadgerConfig) (*badgerStore, error) {
	s := &badgerStore{
		options: options,
		db:      def,
	}
	return s, nil
}

type badgerStore struct {
	options config.BadgerConfig
	db      *badger.DB
}

func (store *badgerStore) Save(p *mqtt.Publish, source string, dest string) (record *model.Record, err error) {
	if p.Qos.Zero() {
		return model.NewRecord(p.Clone(), source, dest), nil
	}
	// store.db.NewStream().Send
	err = store.db.Update(func(txn *badger.Txn) (err error) {
		item, err := txn.Get([]byte(getPacketIDKey(dest)))
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
		r := model.NewRecord(np, source, dest)
		if err = encoder.Encode(r); err != nil {
			return
		}
		if err = txn.Set([]byte(getRecordKey(dest, currentID)), buf.Bytes()); err != nil {
			return
		}
		if err = txn.Set([]byte(getPacketIDKey(dest)), nextID.Encode()); err != nil {
			return
		}
		return
	})
	return
}

func (store *badgerStore) Delete(pid mqtt.PacketID) (ok bool, err error) {
	return
}
func (store *badgerStore) Change(id mqtt.PacketID, state model.State) (err error) {
	return
}
func (store *badgerStore) Receive(p *mqtt.Pubrel) (err error) {
	return
}
func getPacketIDKey(clientID string) string {
	return fmt.Sprintf(packetIDKey, clientID)
}

func getRecordKey(clientID string, packetID mqtt.PacketID) string {
	return fmt.Sprintf(recordKey, clientID, packetID)
}
