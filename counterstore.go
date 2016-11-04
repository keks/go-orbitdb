package orbitdb

import (
	"fmt"
	"sync"

	"github.com/keks/go-ipfs-colog"
)

type ctrIndex struct {
	sync.Mutex

	watchCh <-chan *colog.Entry
	value   int
}

func (idx *ctrIndex) work() {
	for e := range idx.watchCh {
		if e == nil {
			break
		}

		var p CtrPayload

		err := e.Get(&p)
		if err != nil {
			continue
		}

		if p.Op == "COUNTER" {
			idx.Lock()
			idx.value += p.Value
			idx.Unlock()
		}
	}
	fmt.Println("ctridx.work stop.")
}

func (idx *ctrIndex) Value() int {
	idx.Lock()
	defer idx.Unlock()

	return idx.value
}

type CtrPayload struct {
	Op    string `json:"op"`
	Value int    `json:"value"`
}

type CtrStore struct {
	s   *Store
	idx *ctrIndex
}

func NewCtrStore(store *Store) *CtrStore {
	s := &CtrStore{
		s: store,
		idx: &ctrIndex{
			watchCh: store.Watch(),
		},
	}

	go s.idx.work()

	return s
}

func (cs *CtrStore) Increment(by int) (*colog.Entry, error) {
	payload := CtrPayload{
		Op:    "COUNTER",
		Value: by,
	}

	return cs.s.Add(&payload)
}

func (cs *CtrStore) Value() int {
	return cs.idx.Value()
}
