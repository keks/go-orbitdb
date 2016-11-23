package orbitdb

import (
	"fmt"
	"sync"

	"github.com/keks/go-ipfs-colog"
)

var (
	ErrNotFound = fmt.Errorf("key not found in index")
)

type kvPayload struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Op    string `json:"op"`
}

func kvCast(e *colog.Entry) (pl kvPayload, err error) {
	err = e.Get(&pl)
	return pl, err
}

type kvIndex struct {
	l sync.Mutex

	kv map[string]string
}

func (idx *kvIndex) handlePut(e *colog.Entry) error {
	pl, err := kvCast(e)
	if err != nil {
		return err
	}

	idx.l.Lock()
	idx.kv[pl.Key] = pl.Value
	idx.l.Unlock()

	return nil
}

func (idx *kvIndex) handleDel(e *colog.Entry) error {
	pl, err := kvCast(e)
	if err != nil {
		return err
	}

	idx.l.Lock()
	delete(idx.kv, pl.Key)
	idx.l.Unlock()

	return nil
}

func (idx *kvIndex) Get(key string) (value string, err error) {
	var ok bool

	idx.l.Lock()
	value, ok = idx.kv[key]
	idx.l.Unlock()

	if !ok {
		err = ErrNotFound
	}

	return value, err
}

type KVStore struct {
	db  *OrbitDB
	idx *kvIndex
}

func NewKVStore(db *OrbitDB) *KVStore {
	kvs := &KVStore{
		db: db,
		idx: &kvIndex{
			kv: make(map[string]string),
		},
	}

	mux := NewHandlerMux()
	mux.AddHandler(OpPut, kvs.idx.handlePut)
	mux.AddHandler(OpDel, kvs.idx.handleDel)

	go mux.Serve(db)

	return kvs
}

func (kv *KVStore) Put(key, value string) error {
	payload := kvPayload{
		Key:   key,
		Value: value,
		Op:    "PUT",
	}

	_, err := kv.db.Add(&payload)
	return err
}

func (kv *KVStore) Get(key string) (string, error) {
	return kv.idx.Get(key)
}

func (kv *KVStore) Delete(key string) error {
	payload := kvPayload{
		Key: key,
		Op:  "DEL",
	}

	_, err := kv.db.Add(&payload)
	return err
}
