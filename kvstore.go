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

type kvIndex struct {
	sync.Mutex

	watchCh <-chan *colog.Entry
	kv      map[string]string
}

func (idx *kvIndex) work() {
	for {
		e := <-idx.watchCh

		if e == nil {
			break
		}

		p := kvPayload{}

		err := e.Get(&p)
		if err != nil {
			continue
		}

		if p.Op == "PUT" {
			idx.Lock()
			idx.kv[p.Key] = p.Value
			idx.Unlock()
		}

		if p.Op == "DEL" {
			idx.Lock()
			delete(idx.kv, p.Key)
			idx.Unlock()
		}
	}
}

func (idx *kvIndex) Get(key string) (value string, err error) {
	var ok bool

	idx.Lock()
	value, ok = idx.kv[key]
	idx.Unlock()

	if !ok {
		err = ErrNotFound
	}

	return value, err
}

type KVStore struct {
	store *Store
	idx   *kvIndex
}

func NewKVStore(store *Store) *KVStore {
	kvs := &KVStore{
		store: store,
		idx: &kvIndex{
			watchCh: store.Watch(),
			kv:      make(map[string]string),
		},
	}

	go kvs.idx.work()

	return kvs

}

func (kv *KVStore) Put(key, value string) error {
	payload := kvPayload{
		Key:   key,
		Value: value,
		Op:    "PUT",
	}

	_, err := kv.store.Add(payload)
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

	_, err := kv.store.Add(payload)
	return err
}
