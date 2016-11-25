package orbitdb

import (
	"github.com/keks/go-orbitdb/handler"
)

const (
	// OpAdd is the add operation.
	OpAdd handler.Op = "ADD"
	// OpDel is the delete operation.
	OpDel handler.Op = "DEL"
	// OpPut is the put operation
	OpPut handler.Op = "PUT"
	// OpCounter is the counter operation
	OpCounter handler.Op = "COUNTER"
)
