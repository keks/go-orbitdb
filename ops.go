package orbitdb

import (
	"fmt"
	"github.com/keks/go-ipfs-colog"
	"sync"
)

type Op string

const (
	OpAdd     Op = "ADD"
	OpDel     Op = "DEL"
	OpPut     Op = "PUT"
	OpCounter Op = "COUNTER"
)

type HandlerFunc func(*colog.Entry) error

type Handler interface {
	Handle(*colog.Entry) error
}

var (
	WrongOp = fmt.Errorf("operation not supported")
)

type opPayload struct {
	Op `json:"op"`
}

type HandlerMux struct {
	l        sync.Mutex
	handlers map[Op]HandlerFunc
}

func NewHandlerMux() *HandlerMux {
	return &HandlerMux{
		handlers: make(map[Op]HandlerFunc),
	}
}

func (hm *HandlerMux) AddHandler(op Op, h HandlerFunc) {
	hm.l.Lock()
	hm.handlers[op] = h
	hm.l.Unlock()
}

func (hm *HandlerMux) Handle(e *colog.Entry) error {
	var opPl opPayload

	err := e.Get(&opPl)
	if err != nil {
		return WrongOp
	}

	hm.l.Lock()
	h, ok := hm.handlers[opPl.Op]
	hm.l.Unlock()

	if !ok {
		return WrongOp
	}

	return h(e)
}

func (hm *HandlerMux) Serve(db *OrbitDB) {
	for e := range db.Watch() {
		err := hm.Handle(e)
		if err != nil && err != WrongOp {
			// ignore WrongOp errors
			//hm.logger.Log(err)
		}
	}
}
