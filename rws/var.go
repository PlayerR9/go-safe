package rws

import (
	"sync"

	"github.com/PlayerR9/go-safe/common"
)

// Var is a variable that can be read and written safely. To create a Var with the zero
// value, use the `v := new(rws.Var[T])` constructor.
type Var[T any] struct {
	// data is the value of the variable.
	data T

	// mu is the mutex for the variable.
	mu sync.RWMutex
}

// New creates a new Var.
//
// Parameters:
//   - data: The initial value of the variable.
//
// Returns:
//   - *Var[T]: The new Var. Never returns nil.
func New[T any](data T) *Var[T] {
	return &Var[T]{
		data: data,
	}
}

// Copy returns a copy of the Var.
//
// Returns:
//   - *Var[T]: A copy of the Var. Never returns nil.
//
// If the receiver is nil, a new Var is returned.
func (s *Var[T]) Copy() *Var[T] {
	if s == nil {
		return new(Var[T])
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	return &Var[T]{
		data: s.data,
	}
}

// Get returns the value of the variable.
//
// Returns:
//   - T: The value.
//   - error: An error if the receiver is nil.
//
// Errors:
//   - common.ErrNilReceiver: If the receiver is nil.
func (s *Var[T]) Get() (T, error) {
	if s == nil {
		return *new(T), common.ErrNilReceiver
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.data, nil
}

// MustGet returns the value of the variable.
//
// Returns:
//   - T: The value. The zero value if the receiver is nil.
func (s *Var[T]) MustGet() T {
	if s == nil {
		return *new(T)
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.data
}

// Set sets the value of the variable.
//
// Parameters:
//   - data: The value to set.
//
// Returns:
//   - error: An error if the receiver is nil.
func (s *Var[T]) Set(data T) error {
	if s == nil {
		return common.ErrNilReceiver
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.data = data

	return nil
}

// Edit edits the value of the variable. Does nothing if the function is nil.
//
// Parameters:
//   - fn: The function to edit the value.
//
// Returns:
//   - error: An error if the receiver is nil.
func (s *Var[T]) Edit(fn func(elem *T)) error {
	if fn == nil {
		return nil
	} else if s == nil {
		return common.ErrNilReceiver
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	fn(&s.data)

	return nil
}
