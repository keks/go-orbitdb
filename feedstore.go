package orbitdb

import (
	"encoding/json"
	"fmt"
	"github.com/keks/go-ipfs-colog"
)

var ErrMalformedEntry = fmt.Errorf("read malformed colog entry")

func eventCast(e *colog.Entry) (EventPayload, error) {
	var pl EventPayload

	err := e.Get(&pl)
	if err != nil {
		return pl, err
	}

	if len(pl.Data) == 0 {
		return pl, ErrMalformedEntry
	}

	return pl, nil
}

type feedIndex eventIndex

func (idx *feedIndex) handleAdd(e *colog.Entry) error {
	_, err := eventCast(e)
	if err != nil {
		return err
	}

	idx.l.Lock()
	idx.added[e.Hash] = struct{}{}
	idx.l.Unlock()

	return nil
}

func (idx *feedIndex) handleDel(e *colog.Entry) error {
	pl, err := eventCast(e)
	if err != nil {
		return err
	}

	idx.l.Lock()
	delete(idx.added, colog.Hash(pl.DataString()))
	idx.l.Unlock()

	return nil
}

func (idx *feedIndex) has(hash colog.Hash) bool {
	idx.l.Lock()
	defer idx.l.Unlock()

	_, has := idx.added[hash]

	return has
}

type FeedStore struct {
	db  *OrbitDB
	idx *feedIndex
}

func NewFeedStore(db *OrbitDB) *FeedStore {
	fs := &FeedStore{
		db: db,
		idx: &feedIndex{
			added: make(map[colog.Hash]struct{}),
		},
	}

	mux := NewHandlerMux()
	mux.AddHandler(OpAdd, fs.idx.handleAdd)
	mux.AddHandler(OpDel, fs.idx.handleDel)

	go mux.Serve(db)

	return fs
}

func (fs *FeedStore) Get(hash colog.Hash) (*colog.Entry, error) {
	if added := fs.idx.has(hash); !added {
		return nil, ErrNotFound
	}

	return fs.db.colog.Get(hash)
}

func (fs *FeedStore) Delete(hash colog.Hash) (*colog.Entry, error) {
	jsonHash, err := json.Marshal(hash)
	if err != nil {
		return nil, err
	}

	payload := EventPayload{
		Op:   OpDel,
		Data: jsonHash,
	}

	return fs.db.Add(&payload)
}

func (fs *FeedStore) Add(data interface{}) (*colog.Entry, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	payload := EventPayload{
		Op:   OpAdd,
		Data: jsonData,
	}

	return fs.db.Add(&payload)
}

func (fs *FeedStore) Query(qry colog.Query) EventResult {
	res := fs.db.colog.Query(qry)

	return func() (EventPayload, error) {
		for {
			e, err := res()
			if err != nil {
				return EventPayload{}, err
			}

			if fs.idx.has(e.Hash) {
				return eventCast(e)
			}
		}
	}
}
