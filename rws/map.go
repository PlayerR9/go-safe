package rws

import (
	"iter"
	"sync"

	"github.com/PlayerR9/go-safe/common"
)

// Map is a thread-safe map. To create an empty map, use the `sm := new(rws.Map[T, U])`
// constructor.
type Map[T comparable, U any] struct {
	// m is the underlying map.
	m map[T]U

	// mu is the mutex to synchronize map access.
	mu sync.RWMutex
}

// Copy is a method that returns a copy of the Map.
//
// Returns:
//   - *Map[T, U]: A copy of the Map. An empty Map if the receiver is nil.
func (sm *Map[T, U]) Copy() *Map[T, U] {
	if sm == nil {
		return new(Map[T, U])
	}

	sm.mu.RLock()
	defer sm.mu.RUnlock()

	new_map := make(map[T]U, len(sm.m))
	for key, value := range sm.m {
		new_map[key] = value
	}

	return &Map[T, U]{
		m: new_map,
	}
}

// Entry is a method that returns an iterator over the entries in the Map.
//
// Returns:
//   - iter.Seq2[T, U]: An iterator over the entries in the Map. Never returns nil.
func (sm *Map[T, U]) Entry() iter.Seq2[T, U] {
	if sm == nil {
		return func(yield func(T, U) bool) {}
	}

	sm.mu.RLock()
	defer sm.mu.RUnlock()

	fn := func(yield func(T, U) bool) {
		for key, value := range sm.m {
			if !yield(key, value) {
				return
			}
		}
	}

	return fn
}

// Get retrieves a value from the map.
//
// Parameters:
//   - key: The key to retrieve the value.
//
// Returns:
//   - U: The value associated with the key.
//   - bool: A boolean indicating if the key exists in the map.
func (sm *Map[T, U]) Get(key T) (U, bool) {
	if sm == nil || len(sm.m) == 0 {
		return *new(U), false
	}

	sm.mu.RLock()
	defer sm.mu.RUnlock()

	val, ok := sm.m[key]
	return val, ok
}

// Set sets a value in the map. Does nothing if the receiver is nil.
//
// Parameters:
//   - key: The key to set the value.
//   - val: The value to set.
//
// Returns:
//   - error: An error if the receiver is nil.
func (sm *Map[T, U]) Set(key T, val U) error {
	if sm == nil {
		return common.ErrNilReceiver
	}

	sm.mu.Lock()
	defer sm.mu.Unlock()

	if sm.m == nil {
		sm.m = make(map[T]U)
	}

	sm.m[key] = val

	return nil
}

// Delete removes a key from the map. Does nothing if the key does not exist in the map.
//
// Parameters:
//   - key: The key to remove.
func (sm *Map[T, U]) Delete(key T) {
	if sm == nil {
		return
	}

	sm.mu.Lock()
	defer sm.mu.Unlock()

	delete(sm.m, key)
}

// Len returns the number of elements in the map.
//
// Returns:
//   - int: The number of elements in the map. 0 if the receiver is nil.
func (sm *Map[T, U]) Len() int {
	if sm == nil {
		return 0
	}

	sm.mu.RLock()
	defer sm.mu.RUnlock()

	return len(sm.m)
}

// Reset removes all elements from the map.
func (sm *Map[T, U]) Reset() {
	if sm == nil {
		return
	}

	sm.mu.Lock()
	defer sm.mu.Unlock()

	if len(sm.m) == 0 {
		return
	}

	clear(sm.m)
	sm.m = nil
}

// GetMap returns the underlying map.
//
// Returns:
//   - map[T]U: The underlying map. Nil if the receiver is nil.
func (sm *Map[T, U]) GetMap() map[T]U {
	if sm == nil {
		return nil
	}

	sm.mu.RLock()
	defer sm.mu.RUnlock()

	map_copy := make(map[T]U, len(sm.m))
	for key, value := range sm.m {
		map_copy[key] = value
	}

	return map_copy
}
