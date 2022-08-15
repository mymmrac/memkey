/*
Package memkey provides very simple type-safe, thread-safe in memory key-value store with zero dependencies.
*/
package memkey

import (
	"fmt"
	"sync"
)

// Store represents key-value storage that is type-safe and thread-safe to use
type Store[K comparable] struct {
	data map[K]any
	init sync.Once
	lock sync.RWMutex
}

// Entry represents a pair of key and value that can be retrieved from Store
type Entry[K comparable, V any] struct {
	Key   K
	Value V
}

// Get returns a value stored in the store if it exists, or zero value for the type and false
func Get[V any, K comparable](store *Store[K], key K) (V, bool) {
	store.lock.RLock()
	defer store.lock.RUnlock()

	rawValue, ok := store.data[key]
	if !ok {
		return zero[V](), false
	}

	value, ok := rawValue.(V)
	if !ok {
		return zero[V](), false
	}

	return value, true
}

// Get returns raw value stored in the store if it exists, or nil and false
func (s *Store[K]) Get(key K) (any, bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	rawValue, ok := s.data[key]
	if !ok {
		return nil, false
	}

	return rawValue, true
}

// Set stores value with the specified type in the store
func Set[V any, K comparable](store *Store[K], key K, value V) {
	store.lock.Lock()
	defer store.lock.Unlock()

	store.init.Do(func() {
		if store.data == nil {
			store.data = make(map[K]any)
		}
	})

	store.data[key] = value
}

// Set stores value in the store
func (s *Store[K]) Set(key K, value any) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.init.Do(func() {
		if s.data == nil {
			s.data = make(map[K]any)
		}
	})

	s.data[key] = value
}

// Type returns type name of value that is stored, if not found returns false and empty string
func Type[K comparable](store *Store[K], key K) (string, bool) {
	return store.Type(key)
}

// Type returns type name of value that is stored, if not found returns false and empty string
func (s *Store[K]) Type(key K) (string, bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	data, ok := s.data[key]
	if !ok {
		return "", false
	}

	return fmt.Sprintf("%T", data), true
}

// Has returns true if value with the specified key and type exist in the store
func Has[V any, K comparable](store *Store[K], key K) bool {
	store.lock.RLock()
	defer store.lock.RUnlock()

	data, ok := store.data[key]
	if !ok {
		return false
	}

	_, ok = data.(V)
	return ok
}

// Has returns a true if value with the specified key exists in the store with any type
func (s *Store[K]) Has(key K) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()

	_, ok := s.data[key]
	return ok
}

// Delete deletes value from the store if it exists with a specified type and returns true, if not found returns false
func Delete[V any, K comparable](store *Store[K], key K) bool {
	store.lock.Lock()
	defer store.lock.Unlock()

	data, ok := store.data[key]
	if !ok {
		return false
	}

	_, ok = data.(V)
	if !ok {
		return false
	}

	delete(store.data, key)
	return true
}

// Delete deletes value from the store and returns true or if not found reruns false
func (s *Store[K]) Delete(key K) bool {
	s.lock.Lock()
	defer s.lock.Unlock()

	_, ok := s.data[key]
	if !ok {
		return false
	}

	delete(s.data, key)
	return true
}

// Len returns number of values with a specified type that are stored
func Len[V any, K comparable](store *Store[K]) int {
	store.lock.RLock()
	defer store.lock.RUnlock()

	count := 0
	for _, rawValue := range store.data {
		if _, ok := rawValue.(V); !ok {
			continue
		}

		count++
	}

	return count
}

// Len returns number of values that are stored
func (s *Store[K]) Len() int {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return len(s.data)
}

// Keys returns keys of all values with a specified type that are stored
func Keys[V any, K comparable](store *Store[K]) []K {
	store.lock.RLock()
	defer store.lock.RUnlock()

	keys := make([]K, 0, len(store.data))
	for key, rawValue := range store.data {
		if _, ok := rawValue.(V); !ok {
			continue
		}

		keys = append(keys, key)
	}

	return keys
}

// Keys returns keys of all values that are stored
func (s *Store[K]) Keys() []K {
	s.lock.RLock()
	defer s.lock.RUnlock()

	keys := make([]K, 0, len(s.data))
	for key := range s.data {
		keys = append(keys, key)
	}

	return keys
}

// Values returns all values with a specified type that are stored
func Values[V any, K comparable](store *Store[K]) []V {
	store.lock.RLock()
	defer store.lock.RUnlock()

	values := make([]V, 0, len(store.data))
	for _, rawValue := range store.data {
		if value, ok := rawValue.(V); ok {
			values = append(values, value)
		}
	}

	return values
}

// Values returns all values that are stored
func (s *Store[K]) Values() []any {
	s.lock.RLock()
	defer s.lock.RUnlock()

	values := make([]any, 0, len(s.data))
	for _, rawValue := range s.data {
		values = append(values, rawValue)
	}

	return values
}

// Entries returns entries (key-value pairs) where value is of a specified type that are stored
func Entries[V any, K comparable](store *Store[K]) []Entry[K, V] {
	store.lock.RLock()
	defer store.lock.RUnlock()

	entries := make([]Entry[K, V], 0, len(store.data))
	for key, rawValue := range store.data {
		if value, ok := rawValue.(V); ok {
			entries = append(entries, Entry[K, V]{
				Key:   key,
				Value: value,
			})
		}
	}

	return entries
}

// Entries returns entries (key-value pairs) that are stored
func (s *Store[K]) Entries() []Entry[K, any] {
	s.lock.RLock()
	defer s.lock.RUnlock()

	entries := make([]Entry[K, any], 0, len(s.data))
	for key, rawValue := range s.data {
		entries = append(entries, Entry[K, any]{
			Key:   key,
			Value: rawValue,
		})
	}

	return entries
}

// ForEach goes in loop through all values of a specified type and calls f with a key and value
// Warning: May be not thread-safe depending on your usage
func ForEach[V any, K comparable](store *Store[K], f func(key K, value V)) {
	for key, rawValue := range store.data {
		if value, ok := rawValue.(V); ok {
			f(key, value)
		}
	}
}

// ForEach goes in loop through all values and calls f with a key and value
// Warning: May be not thread-safe depending on your usage
func (s *Store[K]) ForEach(f func(key K, value any)) {
	for key, rawValue := range s.data {
		f(key, rawValue)
	}
}
