package pubsub

import (
	shell "github.com/keks/go-ipfs-api"
	pubsub "github.com/keks/go-orbitdb/pubsub"
)

var sh *shell.Shell = shell.NewShell("http://localhost:5001")

type Record struct {
	rec *shell.PubSubRecord
}

func (r Record) From() string {
	return string(r.rec.From)
}

func (r Record) Data() string {
	return string(r.rec.Data)
}

type Subscription struct {
	sub *shell.PubSubSubscription
}

func (s *Subscription) Next() (pubsub.Record, error) {
	r, err := s.sub.Next()

	return Record{r}, err
}

func (s *Subscription) Cancel() error {
	return s.sub.Cancel()
}

type PubSub struct {
	sh *shell.Shell
}

func (ps PubSub) Subscribe(topic string) (pubsub.Subscription, error) {
	sub, err := ps.sh.PubSubSubscribe(topic)

	return &Subscription{
		sub,
	}, err
}

func (ps PubSub) Publish(topic, data string) error {
	return ps.sh.PubSubPublish(topic, data)
}

func New() PubSub {
	return PubSub{sh}
}
