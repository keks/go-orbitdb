package orbitdb

import (
	"encoding/json"
	"sync"

	"github.com/keks/go-ipfs-colog"
	"github.com/keks/go-orbitdb/handler"
)

type eventPayload struct {
	handler.Op `json:"op"`
	Data       json.RawMessage `json:"data"`
}

func eventCast(e *colog.Entry) (eventPayload, error) {
	var pl eventPayload

	err := e.Get(&pl)
	if err != nil {
		return pl, err
	}

	if len(pl.Data) == 0 {
		return pl, ErrMalformedEntry
	}

	return pl, nil
}

func (pl eventPayload) Event() Event {
	return Event{data: pl.Data}
}

// Event is an event stored in an EventStore or FeedStore.
type Event struct {
	data json.RawMessage
}

// GetString returns the string value of the contained data.
// If the contained data is not a string, it returns "".
func (e Event) GetString() string {
	var s string
	json.Unmarshal(e.data, &s)
	return s
}

// Get parses the contained data into v. v Needs to be a pointer.
func (e Event) Get(v interface{}) error {
	return json.Unmarshal(e.data, v)
}

// EventResult is the result of a query to an EventStore or FeedStore.
type EventResult func() (Event, error)

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

// EventStore stores events in an OrbitDB
type EventStore struct {
	idx *eventIndex
	db  *OrbitDB
}

// NewEventStore returns an EventStore for the given OrbitDB.
func NewEventStore(db *OrbitDB) *EventStore {
	evs := &EventStore{
		db: db,
		idx: &eventIndex{
			added: make(map[colog.Hash]struct{}),
		},
	}

	mux := handler.NewMux()
	mux.AddHandler(OpAdd, evs.idx.handleAdd)

	go db.Notify(mux)

	return evs
}

// Add adds an event to the store.
func (evs *EventStore) Add(data interface{}) (*colog.Entry, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	payload := eventPayload{
		Op:   OpAdd,
		Data: jsonData,
	}

	return evs.db.Add(&payload)
}

// Query queries the events using a given query qry.
func (evs *EventStore) Query(qry colog.Query) EventResult {
	res := evs.db.colog.Query(qry)

	return func() (Event, error) {
		for {
			e, err := res()
			if err != nil {
				return Event{}, err
			}

			if evs.idx.has(e.Hash) {
				pl, err := eventCast(e)
				if err != nil {
					return Event{}, err
				}

				return pl.Event(), nil

			}
		}

	}
}
