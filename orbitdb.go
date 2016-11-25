// OrbitDB is a Go implementation of https://github.com/haadcode/orbit-db.
package orbitdb

import (
	"github.com/keks/go-ipfs-colog"
	db "github.com/keks/go-ipfs-colog/immutabledb/ipfs-api"
	"github.com/keks/go-orbitdb/pubsub"
	ippubsub "github.com/keks/go-orbitdb/pubsub/ipfs-api"

	"log"
	"os"
)

// OrbitDB manages a colog and keeps it up-to-date using ipfs pubsub.
type OrbitDB struct {
	topic string

	logger *log.Logger
	colog  *colog.CoLog
	pubsub pubsub.PubSub
}

// NewOrbitDB returns a new OrbitDB and subscribes to the given topic string.
func NewOrbitDB(topic string) (*OrbitDB, error) {
	db := &OrbitDB{
		topic:  topic,
		logger: log.New(os.Stderr, "orbit.OrbitDB ", log.Ltime|log.Lshortfile),
		colog:  colog.New(db.New()),
		pubsub: ippubsub.New(),
	}

	go db.handleSubscription(topic)

	return db, nil
}

// Add adds a given value to the database.
func (db *OrbitDB) Add(data interface{}) (*colog.Entry, error) {
	e, err := db.colog.Add(data)
	if err != nil {
		return e, err
	}

	err = db.pubsub.Publish(db.topic, string(e.Hash))

	return e, err
}

func (db *OrbitDB) handleSubscription(topic string) {
	sub, err := db.pubsub.Subscribe(topic)
	if err != nil {
		db.logger.Println("subscribe error:", err, "; aborting")
		return
	}

	defer sub.Cancel()

	recCh := make(chan pubsub.Record)
	errCh := make(chan error)

	next := func() {
		rec, err := sub.Next()
		if err != nil {
			errCh <- err
		} else {
			recCh <- rec
		}
	}

	go next()

L:
	for {
		select {
		case rec := <-recCh:
			go next()

			err := db.colog.FetchFromHead(colog.Hash(rec.Data()))
			if err != nil {
				db.logger.Println("fetch error:", err, "; continuing")
			}

		case err := <-errCh:
			db.logger.Println("pubsub error:", err, "; cancelling")
			break L
		}
	}
}

// Watch returns a new channel that returns new Entries as they arrive.
func (db *OrbitDB) Watch() <-chan *colog.Entry {
	return db.colog.Watch()
}
