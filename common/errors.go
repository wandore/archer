package common

import "errors"

var (
	ErrCyclicGraph    = errors.New("Cyclic graph, toposort failed!")
	ErrKeyExists      = errors.New("key already exists")
	ErrWaitMismatch   = errors.New("unexpected wait result")
	ErrTooManyClients = errors.New("too many clients")
	ErrNoWatcher      = errors.New("no watcher channel")
)