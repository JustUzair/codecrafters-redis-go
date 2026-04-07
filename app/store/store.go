package store

import (
	"sync"
)

type Value[T any] struct {
	Value            T
	Deadline         int64
	IsDeadlineMillis bool
}
type Storage[T any] struct {
	mu    sync.RWMutex
	store map[string]Value[T]
}

var Cache = &Storage[any]{
	store: make(map[string]Value[any]),
}
