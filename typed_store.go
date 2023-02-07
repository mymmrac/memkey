package memkey

import (
	"sync"
	"time"
)

// TypedStore represents key-value storage with defined keys and values that is type-safe and thread-safe to use
type TypedStore[K comparable, V any] struct {
	data map[K]V
	init sync.Once
	lock sync.RWMutex

	initTTL sync.Once
	ttl     map[K]time.Time
}

// Get return value stored in the store if it exists, or zero value and false
func (s *TypedStore[K, V]) Get(key K) (V, bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	value, ok := s.data[key]
	return value, ok
}

// Set stores value in the store
func (s *TypedStore[K, V]) Set(key K, value V) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.init.Do(func() {
		if s.data == nil {
			s.data = make(map[K]V)
		}
	})

	s.data[key] = value
}

// SetWithTTL stores value in the store with TTL, expiration happens only if ExpireTTL was called
func (s *TypedStore[K, V]) SetWithTTL(key K, value V, ttl time.Duration) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.init.Do(func() {
		if s.data == nil {
			s.data = make(map[K]V)
		}
	})

	s.initTTL.Do(func() {
		if s.ttl == nil {
			s.ttl = make(map[K]time.Time)
		}
	})

	s.data[key] = value
	s.ttl[key] = time.Now().Add(ttl)
}

// ExpireTTL run check for TTL in specified time, and if expired func not nil it will be called with removed item
func (s *TypedStore[K, V]) ExpireTTL(check time.Duration, expired func(key K, value V)) {
	s.initTTL.Do(func() {
		if s.ttl == nil {
			s.ttl = make(map[K]time.Time)
		}
	})

	for now := range time.Tick(check) {
		s.lock.Lock()

		for k, v := range s.ttl {
			if v.Before(now) {
				if expired != nil {
					expired(k, s.data[k])
				}

				delete(s.data, k)
				delete(s.ttl, k)
			}
		}

		s.lock.Unlock()
	}
}

// Has returns a true if value with the specified key exists in the store
func (s *TypedStore[K, V]) Has(key K) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()

	_, ok := s.data[key]
	return ok
}

// Delete deletes value from the store and returns true or if not found reruns false
func (s *TypedStore[K, V]) Delete(key K) bool {
	s.lock.Lock()
	defer s.lock.Unlock()

	_, ok := s.data[key]
	if !ok {
		return false
	}

	delete(s.data, key)
	return true
}

// Len returns number of values that are stored
func (s *TypedStore[K, V]) Len() int {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return len(s.data)
}

// Keys returns keys of all values that are stored, no order is expected
func (s *TypedStore[K, V]) Keys() []K {
	s.lock.RLock()
	defer s.lock.RUnlock()

	keys := make([]K, 0, len(s.data))
	for key := range s.data {
		keys = append(keys, key)
	}

	return keys
}

// Values returns all values that are stored, no order is expected
func (s *TypedStore[K, V]) Values() []V {
	s.lock.RLock()
	defer s.lock.RUnlock()

	values := make([]V, 0, len(s.data))
	for _, rawValue := range s.data {
		values = append(values, rawValue)
	}

	return values
}

// Entries returns entries (key-value pairs) that are stored
func (s *TypedStore[K, V]) Entries() []Entry[K, V] {
	s.lock.RLock()
	defer s.lock.RUnlock()

	entries := make([]Entry[K, V], 0, len(s.data))
	for key, rawValue := range s.data {
		entries = append(entries, Entry[K, V]{
			Key:   key,
			Value: rawValue,
		})
	}

	return entries
}

// ForEach goes in loop through all values and calls f with a key and value
// Warning: May be not thread-safe depending on your usage
func (s *TypedStore[K, V]) ForEach(f func(key K, value V) (stop bool)) {
	for key, rawValue := range s.data {
		if f(key, rawValue) {
			return
		}
	}
}
