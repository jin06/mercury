package message

import (
	"github.com/dgraph-io/badger/v4"
	"github.com/jin06/mercury/internal/config"
	"github.com/jin06/mercury/internal/model"
	"github.com/jin06/mercury/pkg/mqtt"
)

func NewBadgerStore(options config.BadgerConfig) (*badgerStore, error) {
	s := &badgerStore{
		options: options,
	}
	if db, err := badger.Open(badger.DefaultOptions(options.Dir)); err != nil {
		return nil, err
	} else {
		s.db = db
	}
	return s, nil
}

type badgerStore struct {
	options config.BadgerConfig
	db      *badger.DB
}

func (store *badgerStore) Save(p *mqtt.Publish, source string, dest string) (record *model.Record, err error) {
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
