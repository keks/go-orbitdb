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

type evIndex struct {
	sync.Mutex

	watchCh <-chan *colog.Entry
	events  map[colog.Hash]*colog.Entry
}

func (idx *evIndex) work() {
	for {
		e := <-idx.watchCh

		if e == nil {
			break
		}

		p := EventPayload{}

		err := e.Get(&p)
		if err != nil || p.Op != "ADD" {
			continue
		}

		idx.Lock()
		idx.events[e.Hash] = e
		idx.Unlock()
	}
}

type EventStore struct {
	idx   *evIndex
	store *Store
}

func NewEventStore(store *Store) *EventStore {
	evs := &EventStore{
		store: store,
		idx: &evIndex{
			watchCh: store.Watch(),
			events:  make(map[colog.Hash]*colog.Entry),
		},
	}

	go evs.idx.work()

	return evs
}

func (evs *EventStore) Add(data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	payload := EventPayload{
		Op:   "ADD",
		Data: jsonData,
	}

	_, err = evs.store.Add(&payload)

	return err
}

func (evs *EventStore) Query(qry colog.Query) EventResult {
	res := evs.store.colog.Query(qry)

	return func() (EventPayload, error) {
		var payload EventPayload
		for {
			e, err := res()
			if err != nil {
				return EventPayload{}, err
			}

			if err = e.Get(&payload); err != nil {
				continue
			}

			if payload.Op == "ADD" {
				return payload, err
			}
		}

	}
}
