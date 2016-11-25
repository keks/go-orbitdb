package handler

import (
	"fmt"
	"sync"

	"github.com/keks/go-ipfs-colog"
)

// Op is the string that defines an operation
type Op string

var (
	// ErrWrongOp is returned if the called handler does not support the operation
	// specified in the given colog Entry.
	ErrWrongOp = fmt.Errorf("operation not supported")
)

// Func is a function that can be called to handle an incoming colog
// Entry.
type Func func(*colog.Entry) error

// Handler is a type that has a HandlerFunc Handle.
type Handler interface {
	Handle(*colog.Entry) error
}

type opPayload struct {
	Op `json:"op"`
}

// Mux manages several handlers and calles them based on the operation.
type Mux struct {
	l        sync.Mutex
	handlers map[Op]Func
}

// NewMux returns a new HandlerMux.
func NewMux() *Mux {
	return &Mux{
		handlers: make(map[Op]Func),
	}
}

// AddHandler adds a handler h that is to be called when a Entry with operation
// op arrives.
func (hm *Mux) AddHandler(op Op, h Func) {
	hm.l.Lock()
	hm.handlers[op] = h
	hm.l.Unlock()
}

// Handle examines the operation specified in e and calles the appropriate HandlerFunc.
// Returns ErrWrongOp if the given op is not supported.
func (hm *Mux) Handle(e *colog.Entry) error {
	var opPl opPayload

	err := e.Get(&opPl)
	if err != nil {
		return ErrWrongOp
	}

	hm.l.Lock()
	h, ok := hm.handlers[opPl.Op]
	hm.l.Unlock()

	if !ok {
		return ErrWrongOp
	}

	return h(e)
}
