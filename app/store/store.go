package store

import (
	"fmt"
	"sync"
	"time"
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

func (s *Storage[T]) Set(key string, value T, expiry int64, isDeadlineMillis bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	// Never expires
	if expiry == -1 {
		s.store[key] = Value[T]{
			Value:            value,
			Deadline:         -1,
			IsDeadlineMillis: false,
		}
		return
	}

	var deadline int64
	now := time.Now().UnixMilli()

	if isDeadlineMillis {
		deadline = now + expiry //expiry is in ms
	} else {
		deadline = now + (expiry * 1000) // expiry is in secs
	}
	s.store[key] = Value[T]{
		Value:            value,
		Deadline:         deadline,
		IsDeadlineMillis: isDeadlineMillis,
	}
}

func (s *Storage[T]) Get(key string) (T, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.store[key]
	var zero T // return default value of generic
	if !ok {
		err := fmt.Errorf("Key not found")
		return zero, err
	}

	var deadline int64 = val.Deadline
	if deadline == -1 {
		return val.Value, nil
	}
	var isDeadlineMillis bool = val.IsDeadlineMillis
	var currentTime int64 = time.Now().UnixMilli()
	if ok && (isDeadlineMillis && currentTime >= deadline) || (!isDeadlineMillis && currentTime >= deadline*1000) {
		return zero, fmt.Errorf("Key expired")
	}

	return val.Value, nil
}

var Cache = &Storage[any]{
	store: make(map[string]Value[any]),
}
