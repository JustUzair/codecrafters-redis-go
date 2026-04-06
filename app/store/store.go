package store

import (
	"fmt"
	"sync"
)

type Storage[T any] struct {
	mu    sync.RWMutex
	store map[string]T
}

var Cache = &Storage[string]{
	store: make(map[string]string),
}

func (s *Storage[T]) Set(key string, value T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.store[key] = value
}

func (s *Storage[T]) Get(key string) (T, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.store[key]
	if !ok {
		var zero T // return default value of generic
		err := fmt.Errorf("Key not found")
		return zero, err
	}
	return val, nil
}
