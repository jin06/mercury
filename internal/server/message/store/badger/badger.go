package badgerStore

import (
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

func NewBadgerStore(cid string, delivery chan *model.Record) *badgerStore {
	s := &badgerStore{
		options:        config.Def.MessageStore.BadgerConfig,
		db:             def,
		cid:            cid,
		delivery:       delivery,
		resendDuration: time.Second * 5,
		expiry:         time.Hour * 24,
		closing:        make(chan struct{}),
	}
	go s.run()
	return s
}

type badgerStore struct {
	options        config.BadgerConfig
	db             *badger.DB
	cid            string
	expiry         time.Duration
	resendDuration time.Duration
	delivery       chan *model.Record
	closing        chan struct{}
}

func (s *badgerStore) run() error {
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

func (store *badgerStore) resend() {
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
			if record.Qos == mqtt.QoS1 {
				if publish, ok := record.Content.(*mqtt.Publish); ok {
					publish.Dup = true
					record.Content = publish
					store.delivery <- record
				}
				if record.Qos == mqtt.QoS2 {
					if publish, ok := record.Content.(*mqtt.Publish); ok {
						publish.Dup = true
						record.Content = publish
						store.delivery <- record
					}
				}
				record.Times++
			}
		}
		return nil
	})
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
		return model.NewRecord(store.cid, p.Clone()), nil
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
		r := model.NewRecord(store.cid, np)
		buf, err := encodeRecord(r)
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
		r := model.NewRecord(store.cid, p)
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
	data = append(data, byte(r.Qos))
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
	data = binary.BigEndian.AppendUint64(data, uint64(r.Expiry.Nanoseconds()))
	data = append(data, byte(r.Version))
	data = append(data, []byte(r.ClientID)...)
	data = binary.BigEndian.AppendUint64(data, r.Times)
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

	// Qos
	r.Qos = mqtt.QoS(data[i])
	i++

	// Receive time.Time (binary, variable length, so use UnmarshalBinary)
	var t time.Time
	if err := t.UnmarshalBinary(data[i : i+15]); err != nil { // 15 is typical size, adjust as needed
		return nil, err
	}
	r.Receive = t
	i += len(t.AppendFormat(nil, time.RFC3339Nano)) // or exact MarshalBinary size

	// Send time.Time
	var t2 time.Time
	if err := t2.UnmarshalBinary(data[i : i+15]); err != nil {
		return nil, err
	}
	r.Send = t2
	i += len(t2.AppendFormat(nil, time.RFC3339Nano))

	// Expiry Duration (8 bytes)
	if len(data[i:]) < 8 {
		return nil, fmt.Errorf("invalid expiry bytes")
	}
	r.Expiry = time.Duration(int64(binary.BigEndian.Uint64(data[i:])))
	i += 8

	r.Version = mqtt.ProtocolVersion(data[i])
	i++

	// ClientID (must store its length before writing in encode!)
	clientIDLen := binary.BigEndian.Uint16(data[i : i+2])
	i += 2
	if len(data[i:]) < int(clientIDLen) {
		return nil, fmt.Errorf("invalid clientID")
	}
	r.ClientID = string(data[i : i+int(clientIDLen)])
	i += int(clientIDLen)

	// Times (8 bytes)
	if len(data[i:]) < 8 {
		return nil, fmt.Errorf("invalid Times bytes")
	}
	r.Times = binary.BigEndian.Uint64(data[i:])
	i += 8
	packet, err := mqtt.Decode(r.Version, data[i:])
	if err != nil {
		return nil, err
	}
	r.Content = packet
	return r, nil
}
