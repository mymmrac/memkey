package memkey

import "sync"

type Store[K comparable] struct {
	data map[K]any
	init sync.Once
	lock sync.RWMutex
}

type Entry[K comparable, V any] struct {
	key   K
	value V
}

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

func GetRaw[K comparable](store *Store[K], key K) (any, bool) {
	store.lock.RLock()
	defer store.lock.RUnlock()

	rawValue, ok := store.data[key]
	if !ok {
		return nil, false
	}

	return rawValue, true
}

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

func Has[V any, K comparable](store *Store[K], key K) bool {
	store.lock.RLock()
	defer store.lock.RUnlock()

	data, ok := store.data[key]
	if !ok {
		return false
	}

	_, ok = data.(V)
	if !ok {
		return false
	}

	return true
}

func HasKey[K comparable](store *Store[K], key K) bool {
	store.lock.RLock()
	defer store.lock.RUnlock()

	_, ok := store.data[key]
	if !ok {
		return false
	}

	return true
}

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

func DeleteRaw[K comparable](store *Store[K], key K) {
	store.lock.Lock()
	defer store.lock.Unlock()

	delete(store.data, key)
}

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

func Len[K comparable](store *Store[K]) int {
	store.lock.RLock()
	defer store.lock.RUnlock()

	return len(store.data)
}

func Keys[K comparable](store *Store[K]) []K {
	store.lock.RLock()
	defer store.lock.RUnlock()

	keys := make([]K, 0, len(store.data))
	for key := range store.data {
		keys = append(keys, key)
	}

	return keys
}

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

func Values[K comparable](store *Store[K]) []any {
	store.lock.RLock()
	defer store.lock.RUnlock()

	values := make([]any, 0, len(store.data))
	for _, rawValue := range store.data {
		values = append(values, rawValue)
	}

	return values
}

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

func ForEach[K comparable](store *Store[K], f func(key K, value any)) {
	for key, rawValue := range store.data {
		f(key, rawValue)
	}
}

func ForEachOf[V any, K comparable](store *Store[K], f func(key K, value V)) {
	for key, rawValue := range store.data {
		if value, ok := rawValue.(V); ok {
			f(key, value)
		}
	}
}
