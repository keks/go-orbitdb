package pubsub

import (
	"github.com/libp2p/go-libp2p-peer"
)

type PubSub interface {
	Subscribe(topic string) (Subscription, error)
	Publish(topic string, data string) error
}

type Subscription interface {
	Next() (Record, error)
	Cancel() error
}

type Record interface {
	From() peer.ID
	Data() []byte
}
