package orbitdb

import (
	"sync"

	"github.com/keks/go-ipfs-colog"
)

type CtrPayload struct {
	Op    string `json:"op"`
	Value int    `json:"value"`
}

func ctrCast(e *colog.Entry) (pl CtrPayload, err error) {
	err = e.Get(&pl)
	return pl, err
}

type ctrIndex struct {
	l sync.Mutex

	value int
}

func (idx *ctrIndex) handleCounter(e *colog.Entry) error {
	pl, err := ctrCast(e)
	if err != nil {
		return err
	}

	idx.l.Lock()
	idx.value += pl.Value
	idx.l.Unlock()

	return nil
}

func (idx *ctrIndex) Value() int {
	idx.l.Lock()
	defer idx.l.Unlock()

	return idx.value
}

type CtrStore struct {
	db  *OrbitDB
	idx *ctrIndex
}

func NewCtrStore(db *OrbitDB) *CtrStore {
	s := &CtrStore{
		db:  db,
		idx: &ctrIndex{},
	}

	mux := NewHandlerMux()
	mux.AddHandler(OpCounter, s.idx.handleCounter)

	go mux.Serve(db)

	return s
}

func (cs *CtrStore) Increment(by int) (*colog.Entry, error) {
	payload := CtrPayload{
		Op:    "COUNTER",
		Value: by,
	}

	return cs.db.Add(&payload)
}

func (cs *CtrStore) Value() int {
	return cs.idx.Value()
}
