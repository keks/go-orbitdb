package pubsub

type PubSub interface {
	Subscribe(topic string) (Subscription, error)
	Publish(topic string, data string) error
}

type Subscription interface {
	Next() (Record, error)
	Cancel() error
}

type Record interface {
	From() string
	Data() string
}
