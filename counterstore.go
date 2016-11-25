package orbitdb

import (
	"sync"

	"github.com/keks/go-ipfs-colog"
	"github.com/keks/go-orbitdb/handler"
)

type ctrPayload struct {
	handler.Op `json:"op"`
	Value      int `json:"value"`
}

func ctrCast(e *colog.Entry) (pl ctrPayload, err error) {
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

// CtrStore manages a counter stored in an OrbitDB.
type CtrStore struct {
	db  *OrbitDB
	idx *ctrIndex
}

// NewCtrStore returns a CtrStore for the given OrbitDB.
func NewCtrStore(db *OrbitDB) *CtrStore {
	s := &CtrStore{
		db:  db,
		idx: &ctrIndex{},
	}

	mux := handler.NewMux()
	mux.AddHandler(OpCounter, s.idx.handleCounter)

	go db.Notify(mux)

	return s
}

// Increment increments the counter by n.
func (cs *CtrStore) Increment(n int) (*colog.Entry, error) {
	payload := ctrPayload{
		Op:    OpCounter,
		Value: n,
	}

	return cs.db.Add(&payload)
}

// Value returns the current value of the counter.
func (cs *CtrStore) Value() int {
	return cs.idx.Value()
}
