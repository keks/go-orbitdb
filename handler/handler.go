package handler

import (
	"fmt"
	"sync"

	"github.com/keks/go-ipfs-colog"
)

// Op is the string that defines an operation
type Op string

var (
	// WrongOp is returned if the called handler does not support the operation
	// specified in the given colog Entry.
	WrongOp = fmt.Errorf("operation not supported")
)

// HandlerFunc is a function that can be called to handle an incoming colog
// Entry.
type HandlerFunc func(*colog.Entry) error

// Handler is a type that has a HandlerFunc Handle.
type Handler interface {
	Handle(*colog.Entry) error
}

type opPayload struct {
	Op `json:"op"`
}

// HandlerMux manages several handlers and calles them based on the operation.
type HandlerMux struct {
	l        sync.Mutex
	handlers map[Op]HandlerFunc
}

// NewHandlerMux returns a new HandlerMux.
func NewHandlerMux() *HandlerMux {
	return &HandlerMux{
		handlers: make(map[Op]HandlerFunc),
	}
}

// AddHandler adds a handler h that is to be called when a Entry with operation
// op arrives.
func (hm *HandlerMux) AddHandler(op Op, h HandlerFunc) {
	hm.l.Lock()
	hm.handlers[op] = h
	hm.l.Unlock()
}

// Handle examines the operation specified in e and calles the appropriate HandlerFunc.
// Returns WrongOp if the given op is not supported.
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
