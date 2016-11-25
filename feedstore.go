package orbitdb

import (
	"encoding/json"
	"fmt"

	"github.com/keks/go-ipfs-colog"
	"github.com/keks/go-orbitdb/handler"
)

// ErrMalformedEntry is returned when a colog Entry does not have the expected
// format.
var ErrMalformedEntry = fmt.Errorf("read malformed colog entry")

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
	delete(idx.added, colog.Hash(pl.Event().GetString()))
	idx.l.Unlock()

	return nil
}

func (idx *feedIndex) has(hash colog.Hash) bool {
	idx.l.Lock()
	defer idx.l.Unlock()

	_, has := idx.added[hash]

	return has
}

// FeedStore is similar to an EventStore but also allows deleting events.
type FeedStore struct {
	db  *OrbitDB
	idx *feedIndex
}

// NewFeedStore returns a FeedStore for the given OrbitDB.
func NewFeedStore(db *OrbitDB) *FeedStore {
	fs := &FeedStore{
		db: db,
		idx: &feedIndex{
			added: make(map[colog.Hash]struct{}),
		},
	}

	mux := handler.NewHandlerMux()
	mux.AddHandler(OpAdd, fs.idx.handleAdd)
	mux.AddHandler(OpDel, fs.idx.handleDel)

	go db.Notify(mux)

	return fs
}

// Delete deletes the event at the given hash.
func (fs *FeedStore) Delete(hash colog.Hash) (*colog.Entry, error) {
	jsonHash, err := json.Marshal(hash)
	if err != nil {
		return nil, err
	}

	payload := eventPayload{
		Op:   OpDel,
		Data: jsonHash,
	}

	return fs.db.Add(&payload)
}

// Add adds a new Event.
func (fs *FeedStore) Add(data interface{}) (*colog.Entry, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	payload := eventPayload{
		Op:   OpAdd,
		Data: jsonData,
	}

	return fs.db.Add(&payload)
}

// Query queries the events using a given query qry. It omits deleted events.
func (fs *FeedStore) Query(qry colog.Query) EventResult {
	res := fs.db.colog.Query(qry)

	return func() (Event, error) {
		for {
			e, err := res()
			if err != nil {
				return Event{}, err
			}

			if fs.idx.has(e.Hash) {
				pl, err := eventCast(e)
				if err != nil {
					return Event{}, err
				}

				return pl.Event(), nil
			}
		}
	}
}
