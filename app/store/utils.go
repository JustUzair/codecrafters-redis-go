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
	if workers, ok := s.notifiers[list_key]; ok && len(workers) > 0 {
		ringBell := workers[0]
		s.notifiers[list_key] = workers[1:]
		ringBell <- struct{}{}
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

func (s *Storage[T]) LPop(list_key string, n_pop int) ([]any, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	entry, ok := s.store[list_key]
	if !ok {
		return nil, fmt.Errorf("key not found")
	}
	var list []any
	list = any(entry.Value).([]any) // for popping value must be any[] list
	if n_pop > len(list) {
		n_pop = len(list)
		emptyList := make([]any, 0)
		delete(s.store, list_key)
		return emptyList, nil
	}

	var elements []any
	elements = list[:n_pop]
	// 0 1 2 3 4 5, n_pop = 2, list [2 : ] ===> correct
	newList := list[n_pop:]

	s.store[list_key] = Value[T]{
		Value:            any(newList).(T),
		Deadline:         entry.Deadline,
		IsDeadlineMillis: entry.IsDeadlineMillis,
	}
	return elements, nil

}

func (s *Storage[T]) BLPop(list_key string, timeout float64) ([]any, error) {
	s.mu.Lock() // acquire lock
	// ----  Step 1. If KV is present, check if list is non empty, if so pop immediately ----
	entry, exists := s.store[list_key]

	if exists {
		list, ok := any(entry.Value).([]any)
		if ok && len(list) > 0 {
			// Found data! Pop it using your internal (non-locking) logic
			element := list[0]
			newList := list[1:]

			// Update store
			if len(newList) == 0 {
				delete(s.store, list_key)
			} else {
				s.store[list_key] = Value[T]{Value: any(newList).(T), Deadline: entry.Deadline, IsDeadlineMillis: entry.IsDeadlineMillis}
			}

			s.mu.Unlock()
			return []any{list_key, element}, nil
		}
	}
	// ----  Step 1. End ----

	// ---- Step 2. Entry doesn't exist, create a bell, that rings channel for updates in data----
	bell := make(chan struct{}, 1)
	s.notifiers[list_key] = append(s.notifiers[list_key], bell)
	s.mu.Unlock()

	var timeoutChannel <-chan time.Time
	if timeout != 0 {
		timeoutChannel = time.After(time.Duration(timeout * float64(time.Second)))
	}
	select {
	case <-bell: // Bell rung, there were updates in the list_key kv pair, consume bell, and pop
		fmt.Println("Channel changes occured")
		elements, err := s.LPop(list_key, 1)
		if err != nil {
			return nil, err
		}
		workers := s.notifiers[list_key]
		for i, v := range workers {
			if v == bell {
				s.notifiers[list_key] = append(workers[:i], workers[i+1:]...)
				break
			}
		}
		return []any{list_key, elements[0]}, nil
	case <-timeoutChannel: // Timeout has occured
		return nil, fmt.Errorf("Timeout")
	}
	// ---- Step 2. ----

}

func (s *Storage[T]) Type(key string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	entry, ok := s.store[key]
	if !ok {
		return "none"
	}

	switch any(entry.Value).(type) {
	case string:
		return "string"
	case []any:
		return "list"
	case Stream:
		return "stream"
	default:
		return "none"
	}

}

func (s *Storage[T]) XAdd(list_key string, id string, fields []Field) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	entry, exists := s.store[list_key]
	if !exists {
		// inner most entry
		var streamEntry StreamEntry = StreamEntry{
			ID:     id,
			Fields: fields,
		}
		// create 1 entry space for the above entry

		var streamEntries []StreamEntry = make([]StreamEntry, 1)
		streamEntries[0] = streamEntry

		var stream Stream = Stream{
			StreamEntries: streamEntries,
		}
		// create stream
		s.store[list_key] = Value[T]{
			Value:            any(stream).(T),
			Deadline:         -1,
			IsDeadlineMillis: false,
		}
		return 1, nil
	} else {
		existingStream := any(entry.Value).(Stream)
		existingStream.StreamEntries = append(existingStream.StreamEntries, StreamEntry{
			ID:     id,
			Fields: fields,
		})
		s.store[list_key] = Value[T]{
			Value:            any(existingStream).(T),
			Deadline:         entry.Deadline,
			IsDeadlineMillis: entry.IsDeadlineMillis,
		}
	}

	return 0, nil
}
