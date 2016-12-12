package pubsub

import (
	shell "github.com/ipfs/go-ipfs-api"
	pubsub "github.com/keks/go-orbitdb/pubsub"
)

var sh *shell.Shell = shell.NewShell("http://localhost:5001")

type Subscription struct {
	sub *shell.PubSubSubscription
}

func (s *Subscription) Next() (pubsub.Record, error) {
	return s.sub.Next()
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
