package orbitdb

import (
	"github.com/keks/go-ipfs-colog"
	db "github.com/keks/go-ipfs-colog/immutabledb/ipfs-api"
	"github.com/keks/go-orbitdb/pubsub"
	ippubsub "github.com/keks/go-orbitdb/pubsub/ipfs-api"

	"log"
	"os"
)

type Store struct {
	topic string

	logger *log.Logger
	colog  *colog.CoLog
	pubsub pubsub.PubSub
}

func NewStore(topic string) (*Store, error) {
	s := &Store{
		topic:  topic,
		logger: log.New(os.Stderr, "orbit.Store ", log.Ltime|log.Lshortfile),
		colog:  colog.New(db.New()),
		pubsub: ippubsub.New(),
	}

	go s.handleSubscription(topic)

	return s, nil
}

func (s *Store) Add(data interface{}) (*colog.Entry, error) {
	e, err := s.colog.Add(data)
	if err != nil {
		return e, err
	}

	err = s.pubsub.Publish(s.topic, string(e.Hash))

	return e, err
}

func (s *Store) handleSubscription(topic string) {
	sub, err := s.pubsub.Subscribe(topic)
	if err != nil {
		s.logger.Println("subscribe error:", err, "; aborting")
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

			err := s.colog.FetchFromHead(colog.Hash(rec.Data()))
			if err != nil {
				s.logger.Println("fetch error:", err, "; continuing")
			}

		case err := <-errCh:
			s.logger.Println("pubsub error:", err, "; cancelling")
			break L
		}
	}
}

func (s *Store) Watch() <-chan *colog.Entry {
	return s.colog.Watch()
}
