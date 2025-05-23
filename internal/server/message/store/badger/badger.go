package badgerStore

import (
	"context"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/jin06/mercury/internal/config"
	"github.com/jin06/mercury/internal/logger"
	"github.com/jin06/mercury/internal/model"
	"github.com/jin06/mercury/pkg/mqtt"
)

var (
	def             *badger.DB
	packetIDKey     = "packetid:%s"  // -> packetid:{clientID}
	recordKey       = "record:%s:%d" // -> record:{clientID}:{PacketID}
	recordPrefixKey = "record:%s"
)

func Init(options config.BadgerConfig) (err error) {
	def, err = badger.Open(badger.DefaultOptions(options.Dir))
	return
}

func New(cid string) *badgerStore {
	s := &badgerStore{
		options:        config.Def.MessageStore.BadgerConfig,
		db:             def,
		cid:            cid,
		resendDuration: time.Second * 5,
		expiry:         config.Def.MQTTConfig.MessageExpiryInterval,
		closing:        make(chan struct{}),
	}
	return s
}

// func NewBadgerStore(cid string, delivery chan *model.Record) *badgerStore {
// 	s := &badgerStore{
// 		options:        config.Def.MessageStore.BadgerConfig,
// 		db:             def,
// 		cid:            cid,
// 		delivery:       delivery,
// 		resendDuration: time.Second * 5,
// 		expiry:         config.Def.MQTTConfig.MessageExpiryInterval,
// 		closing:        make(chan struct{}),
// 	}
// 	return s
// }

type badgerStore struct {
	options        config.BadgerConfig
	db             *badger.DB
	cid            string
	expiry         time.Duration
	resendDuration time.Duration
	// delivery       chan *model.Record
	closing chan struct{}
}

func (s *badgerStore) Run(ctx context.Context, ch chan mqtt.Packet) error {
	ticker := time.NewTicker(s.resendDuration)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-s.closing:
			return nil
		case <-ticker.C:
			s.resend(ch)
		}
	}
}

func (store *badgerStore) resend(ch chan mqtt.Packet) {
	store.db.View(func(txn *badger.Txn) error {
		prefix := store.getRecordPrefix()
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = true
		opts.Prefix = prefix

		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			v, err := item.ValueCopy(nil)
			if err != nil {
				logger.Error(err)
				continue
			}
			record, err := decodeRecord(v)
			if err != nil {
				logger.Error(err)
				continue
			}
			ch <- record.Content
			record.Times++
		}
		return nil
	})
}

func (store *badgerStore) Clean() (err error) {
	if err = store.db.DropPrefix([]byte(store.getRecordPrefix())); err != nil {
		return
	}
	if err = store.db.DropPrefix([]byte(store.getPacketIDKey())); err != nil {
		return
	}
	return
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
	return store.delete(p.PacketID)
	// return nil
	// return store.update(p)
}

func (store *badgerStore) Complete(pid mqtt.PacketID) error {
	return store.delete(pid)
}

func (store *badgerStore) create(p *mqtt.Publish) (record *model.Record, err error) {
	if p.Qos.Zero() {
		return model.NewRecord(store.cid, p.Clone(), store.expiry), nil
	}
	err = store.db.Update(func(txn *badger.Txn) (err error) {
		var currentID mqtt.PacketID

		item, err := txn.Get([]byte(store.getPacketIDKey()))

		switch err {
		case badger.ErrKeyNotFound:
			currentID = 1
			err = nil
		case nil:
			if val, err := item.ValueCopy(nil); err != nil {
				return err
			} else {
				if err = currentID.Decode(val); err != nil {
					return err
				}
			}
		default:
			return err
		}

		if currentID == 0 || currentID > mqtt.MAX_PACKET_ID {
			currentID = 1
		}
		nextID := currentID + 1
		np := p.Clone()
		np.PacketID = currentID
		record = model.NewRecord(store.cid, np, config.Def.MQTTConfig.MessageExpiryInterval)
		buf, err := encodeRecord(record)
		if err != nil {
			return err
		}
		if err = txn.Set([]byte(store.getRecordKey(currentID)), buf); err != nil {
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
	err := store.db.Update(func(txn *badger.Txn) (err error) {
		currentID := p.PID()
		r := model.NewRecord(store.cid, p, store.expiry)
		buf, err := encodeRecord(r)
		if err != nil {
			return err
		}
		return txn.Set([]byte(store.getRecordKey(currentID)), buf)
	})
	return err
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

func (store *badgerStore) getRecordPrefix() []byte {
	return []byte(fmt.Sprintf(recordPrefixKey, store.cid))
}

func encodeRecord(r *model.Record) ([]byte, error) {
	data := []byte{}
	data = append(data, byte(r.Version))
	data = binary.BigEndian.AppendUint64(data, r.Times)
	data = binary.BigEndian.AppendUint64(data, uint64(r.Expiry))
	if raw, err := mqtt.EncodeUTF8(r.ClientID); err != nil {
		return nil, err
	} else {
		data = append(data, raw...)
	}
	if raw, err := r.Receive.MarshalBinary(); err != nil {
		return nil, err
	} else {
		data = append(data, raw...)
	}
	if raw, err := r.Send.MarshalBinary(); err != nil {
		return nil, err
	} else {
		data = append(data, raw...)
	}
	if raw, err := r.Content.Encode(); err != nil {
		return nil, err
	} else {
		data = append(data, raw...)
	}
	return data, nil
}

func decodeRecord(data []byte) (*model.Record, error) {
	r := &model.Record{}
	i := 0

	r.Version = mqtt.ProtocolVersion(data[i])
	i++

	r.Times = binary.BigEndian.Uint64(data[i:])
	i += 8

	r.Expiry = time.Duration(binary.BigEndian.Uint64(data[i:]))
	i += 8

	if clientID, n, err := mqtt.DecodeUTF8(data[i:]); err != nil {
		return nil, err
	} else {
		r.ClientID = clientID
		i += n
	}

	// Receive time.Time (binary, variable length, so use UnmarshalBinary)
	if err := r.Receive.UnmarshalBinary(data[i : i+15]); err != nil { // 15 is typical size, adjust as needed
		return nil, err
	}
	i += 15

	// Send time.Time
	if err := r.Send.UnmarshalBinary(data[i : i+15]); err != nil {
		return nil, err
	}
	i += 15

	if packet, err := mqtt.Decode(r.Version, data[i:]); err != nil {
		return nil, err
	} else {
		r.Content = packet
	}
	return r, nil
}

func (b *badgerStore) Close() error {
	close(b.closing)
	return nil
}
