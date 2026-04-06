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
	deadline = int64(time.Now().Unix())

	if isDeadlineMillis {
		deadline = int64((time.Now().Unix() * 1000) + expiry) //expiry is in ms
	} else {
		deadline = int64(time.Now().Unix() + expiry) // expiry is in secs
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
	var currentTime int64 = time.Now().Unix()
	if ok && (isDeadlineMillis && currentTime*1000 > deadline) || (!isDeadlineMillis && currentTime > deadline) {
		return zero, fmt.Errorf("Key expired")
	}

	return val.Value, nil
}

var Cache = &Storage[string]{
	store: make(map[string]Value[string]),
}
