package orbitdb

import (
	"encoding/json"
	"sync"

	"github.com/keks/go-ipfs-colog"
)

type EventPayload struct {
	Op   string          `json:"op"`
	Data json.RawMessage `json:"data"`
}

func (pl EventPayload) DataString() string {
	var s string
	json.Unmarshal(pl.Data, &s)
	return s
}

type EventResult func() (EventPayload, error)

type eventIndex struct {
	l sync.Mutex

	added map[colog.Hash]struct{}
}

func (idx *eventIndex) handleAdd(e *colog.Entry) error {
	_, err := eventCast(e)
	if err != nil {
		return err
	}

	idx.l.Lock()
	idx.added[e.Hash] = struct{}{}
	idx.l.Unlock()

	return nil
}

func (idx *eventIndex) has(hash colog.Hash) bool {
	idx.l.Lock()
	defer idx.l.Unlock()

	_, has := idx.added[hash]
	return has
}

type EventStore struct {
	idx *eventIndex
	db  *OrbitDB
}

func NewEventStore(db *OrbitDB) *EventStore {
	evs := &EventStore{
		db: db,
		idx: &eventIndex{
			added: make(map[colog.Hash]struct{}),
		},
	}

	mux := NewHandlerMux()
	mux.AddHandler(OpAdd, evs.idx.handleAdd)

	go mux.Serve(db)

	return evs
}

func (evs *EventStore) Add(data interface{}) (*colog.Entry, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	payload := EventPayload{
		Op:   "ADD",
		Data: jsonData,
	}

	return evs.db.Add(&payload)
}

func (evs *EventStore) Query(qry colog.Query) EventResult {
	res := evs.db.colog.Query(qry)

	return func() (EventPayload, error) {
		for {
			e, err := res()
			if err != nil {
				return EventPayload{}, err
			}

			if evs.idx.has(e.Hash) {
				return eventCast(e)
			}
		}

	}
}
