package store

import (
	"fmt"
	"time"
)

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

func (s *Storage[T]) Push(list_key string, values []string, isLPUSH bool) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	var list []any
	entry, ok := s.store[list_key]
	if ok {
		var temp any = entry.Value
		list, ok = temp.([]any)
		if !ok {
			// Not a list type
			return -1
		}
	}
	/*
		list.append(vals,existing)
		existing[]
		val[]
	*/
	if isLPUSH {
		// new vals --> a b c
		// old --> 1 2 3
		// lpush --> c b a
		// final --> c b a 1 2 3
		reversed := make([]any, len(values))
		for i, v := range values {
			reversed[len(values)-1-i] = v
		}
		list = append(reversed, list...)
	} else {
		for _, v := range values {
			list = append(list, v)
		}
	}
	s.store[list_key] = Value[T]{
		Value:            any(list).(T),
		Deadline:         -1,
		IsDeadlineMillis: false,
	}
	return len(list)
}

func (s *Storage[T]) LRange(list_key string, start int64, stop int64) []any {
	// Locks are set in Get()
	temp, err := Cache.Get(list_key)

	defaultValue := make([]any, 0)
	if err != nil {
		// KV pair doesnt exist
		return defaultValue
	}
	var vals []any
	vals, ok := temp.([]any)
	if !ok {
		// wrong type
		return defaultValue
	}
	var valueLen64 int64 = int64(len(vals))

	// a,b,c,d,e
	if start < 0 {
		start = valueLen64 + start // ex: 5 + (-2) = 3
	}
	if stop < 0 {
		stop = valueLen64 + stop // ex: 5 + (-1) = 4
	}

	//  Clamp the indices if deduction introduced out of bounds index
	if start < 0 {
		start = 0
	}
	if stop < 0 {
		stop = 0
	}

	/*
		The LRANGE command has several behaviors to keep in mind:

		If the list doesn't exist, an empty array is returned.
		If the start index is greater than or equal to the list's length, an empty array is returned.
		If the stop index is greater than or equal to the list's length, the stop index is treated as the last element.
		If the start index is greater than the stop index, an empty array is returned.

	*/
	if start >= valueLen64 || valueLen64 == 0 {
		return defaultValue
	} else if stop >= valueLen64 {
		stop = valueLen64 - 1 // We return by incrementing  stop + 1 for inclusiveness we reduce by 1 so the return doesn't go out of bounds
	} else if start > stop {
		return defaultValue
	}
	return vals[start : stop+1]
}

func (s *Storage[T]) LLen(list_key string) int {
	temp, err := Cache.Get(list_key)

	if err != nil {
		// KV pair doesnt exist
		return 0
	}

	var list []any
	list = temp.([]any)
	return len(list)

}

func (s *Storage[T]) LPop(list_key string) (any, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry, ok := s.store[list_key]
	if !ok {
		return nil, fmt.Errorf("key not found")
	}
	var list []any
	list = any(entry.Value).([]any) // for popping value must be any[] list
	elememt := list[0]
	newList := list[1:]

	s.store[list_key] = Value[T]{
		Value:            any(newList).(T),
		Deadline:         entry.Deadline,
		IsDeadlineMillis: entry.IsDeadlineMillis,
	}
	return elememt, nil

}
