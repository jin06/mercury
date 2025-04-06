package message

import (
	"github.com/dgraph-io/badger/v4"
	"github.com/jin06/mercury/internal/config"
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
