/*
Package memkey provides very simple type-safe, thread-safe in memory key-value store with zero dependencies.
*/
package memkey

import "sync"

// Store represents key-value storage that is type-safe and thread-safe to use
type Store[K comparable] struct {
	data map[K]any
	init sync.Once
	lock sync.RWMutex
}

// Entry represents a pair of key and value that can be retrieved from Store
type Entry[K comparable, V any] struct {
	key   K
	value V
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

// GetRaw returns raw value stored in the store if it exists, or nil and false
func GetRaw[K comparable](store *Store[K], key K) (any, bool) {
	store.lock.RLock()
	defer store.lock.RUnlock()

	rawValue, ok := store.data[key]
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

// HasKey returns true if value with the specified key exist in the store with any type
func HasKey[K comparable](store *Store[K], key K) bool {
	store.lock.RLock()
	defer store.lock.RUnlock()

	_, ok := store.data[key]
	return ok
}

// Delete deletes value from the store if it exists with a specified type, if not found this is no-op
func Delete[V any, K comparable](store *Store[K], key K) {
	store.lock.Lock()
	defer store.lock.Unlock()

	data, ok := store.data[key]
	if !ok {
		return
	}

	_, ok = data.(V)
	if !ok {
		return
	}

	delete(store.data, key)
}

// DeleteOk deletes value from the store if it exists with a specified type and returns true, if not found returns false
func DeleteOk[V any, K comparable](store *Store[K], key K) bool {
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

// DeleteRaw deletes value from the store, if not found this is no-op
func DeleteRaw[K comparable](store *Store[K], key K) {
	store.lock.Lock()
	defer store.lock.Unlock()

	delete(store.data, key)
}

// DeleteRawOk deletes value from the store and returns true or if not found reruns false
func DeleteRawOk[K comparable](store *Store[K], key K) bool {
	store.lock.Lock()
	defer store.lock.Unlock()

	_, ok := store.data[key]
	if !ok {
		return false
	}

	delete(store.data, key)
	return true
}

// Len returns number of elements stored
func Len[K comparable](store *Store[K]) int {
	store.lock.RLock()
	defer store.lock.RUnlock()

	return len(store.data)
}

// Keys returns keys of all values that are stored
func Keys[K comparable](store *Store[K]) []K {
	store.lock.RLock()
	defer store.lock.RUnlock()

	keys := make([]K, 0, len(store.data))
	for key := range store.data {
		keys = append(keys, key)
	}

	return keys
}

// KeysOf returns keys of all values with a specified type that are stored
func KeysOf[V any, K comparable](store *Store[K]) []K {
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

// Values returns all values that are stored
func Values[K comparable](store *Store[K]) []any {
	store.lock.RLock()
	defer store.lock.RUnlock()

	values := make([]any, 0, len(store.data))
	for _, rawValue := range store.data {
		values = append(values, rawValue)
	}

	return values
}

// ValuesOf returns all values with a specified that are stored
func ValuesOf[V any, K comparable](store *Store[K]) []V {
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

// Entries returns entries (key-value pairs) that are stored
func Entries[K comparable](store *Store[K]) []Entry[K, any] {
	store.lock.RLock()
	defer store.lock.RUnlock()

	entries := make([]Entry[K, any], 0, len(store.data))
	for key, rawValue := range store.data {
		entries = append(entries, Entry[K, any]{
			key:   key,
			value: rawValue,
		})
	}

	return entries
}

// EntriesOf returns entries (key-value pairs) where value is of a specified type that are stored
func EntriesOf[V any, K comparable](store *Store[K]) []Entry[K, V] {
	store.lock.RLock()
	defer store.lock.RUnlock()

	entries := make([]Entry[K, V], 0, len(store.data))
	for key, rawValue := range store.data {
		if value, ok := rawValue.(V); ok {
			entries = append(entries, Entry[K, V]{
				key:   key,
				value: value,
			})
		}
	}

	return entries
}

// ForEach goes in loop through all values and calls f with a key and value
// Warning: Not thread-safe
func ForEach[K comparable](store *Store[K], f func(key K, value any)) {
	for key, rawValue := range store.data {
		f(key, rawValue)
	}
}

// ForEachOf goes in loop through all values of a specified type and calls f with a key and value
// Warning: Not thread-safe
func ForEachOf[V any, K comparable](store *Store[K], f func(key K, value V)) {
	for key, rawValue := range store.data {
		if value, ok := rawValue.(V); ok {
			f(key, value)
		}
	}
}
