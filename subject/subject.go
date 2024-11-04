package subject

import (
	"errors"
	"sync"

	common "github.com/PlayerR9/go-safe/common"
)

// Subject is a subject that can be observed.
type Subject[T any] struct {
	// state is the current state of the subject.
	state T

	// observers are the observers of the subject.
	observers []Observer[T]

	// mu is the mutex for the subject.
	mu sync.RWMutex
}

// Notify implements the Observer interface.
func (s *Subject[T]) Notify(change T) error {
	if s == nil {
		return common.ErrNilReceiver
	}

	s.mu.Lock()
	s.state = change
	s.mu.Unlock()

	err := s.NotifyAll()
	return err
}

// Cleanup implements the Observer interface.
func (s *Subject[T]) Cleanup() {
	if s == nil {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.state = *new(T)

	if len(s.observers) > 0 {
		for _, o := range s.observers {
			o.Cleanup()
		}

		clear(s.observers)
		s.observers = nil
	}
}

// New creates a new Subject with the given initial value.
//
// Parameters:
//   - value: The initial value of the subject.
//
// Returns:
//   - *Subject[T]: The new Subject. Never returns nil.
func New[T any](value T) *Subject[T] {
	return &Subject[T]{
		state: value,
	}
}

// Set sets the value of the subject.
//
// Parameters:
//   - value: The value to set.
//
// Returns:
//   - error: An error if the receiver is nil or if any of the observers fail to notify.
//
// Errors:
//   - common.ErrNilReceiver: If the receiver is nil.
//   - any other error: If any of the observers fail to notify.
func (s *Subject[T]) Set(value T) error {
	if s == nil {
		return common.ErrNilReceiver
	}

	s.mu.Lock()
	s.state = value
	s.mu.Unlock()

	err := s.NotifyAll()
	return err
}

// Get returns the value of the subject.
//
// Returns:
//   - T: The value of the subject. The zero value if the receiver is nil.
//   - error: An error if the receiver is nil.
//
// Errors:
//   - common.ErrNilReceiver: If the receiver is nil.
func (s *Subject[T]) Get() (T, error) {
	if s == nil {
		return *new(T), common.ErrNilReceiver
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.state, nil
}

// MustGet returns the value of the subject.
//
// Returns:
//   - T: The value of the subject. The zero value if the receiver is nil.
func (s *Subject[T]) MustGet() T {
	if s == nil {
		return *new(T)
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.state
}

// Edit edits the value of the subject.
//
// Parameters:
//   - fn: The function to edit the value.
//
// Returns:
//   - error: An error if the receiver is nil or the function is nil.
//
// Errors:
//   - common.ErrNilReceiver: If the receiver is nil.
//   - any other error: If any of the observers fail to notify.
func (s *Subject[T]) Edit(fn func(elem *T)) error {
	if fn == nil {
		return nil
	} else if s == nil {
		return common.ErrNilReceiver
	}

	s.mu.Lock()
	fn(&s.state)
	s.mu.Unlock()

	err := s.NotifyAll()
	return err
}

// Attach adds an observer to the subject. Does nothing if the observer is nil.
//
// Parameters:
//   - o: The observer to add.
//
// Returns:
//   - error: An error if the receiver is nil.
//
// Errors:
//   - common.ErrNilReceiver: If the receiver is nil.
func (s *Subject[T]) Attach(o Observer[T]) error {
	if o == nil {
		return nil
	} else if s == nil {
		return common.ErrNilReceiver
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.observers = append(s.observers, o)

	return nil
}

// NotifyAll notifies all observers of the subject.
//
// Returns:
//   - error: An error if the receiver is nil or if any of the observers fail to notify.
//
// Errors:
//   - common.ErrNilReceiver: If the receiver is nil.
//   - any other error: If any of the observers fail to notify.
func (s *Subject[T]) NotifyAll() error {
	if s == nil {
		return common.ErrNilReceiver
	}

	// Critical section.
	s.mu.RLock()
	state := s.state
	observers := make([]Observer[T], len(s.observers))
	copy(observers, s.observers)
	s.mu.RUnlock()

	// Prepare to notify.
	var wg sync.WaitGroup

	errs := make([]error, 0, len(observers))
	var mu sync.Mutex
	fns := make([]func(), 0, len(observers))

	for _, o := range observers {
		fn := func() {
			defer wg.Done()

			err := o.Notify(state)
			if err == nil {
				return
			}

			mu.Lock()
			errs = append(errs, err)
			mu.Unlock()
		}

		fns = append(fns, fn)
	}

	// Notify.
	wg.Add(len(fns))

	for _, fn := range fns {
		go fn()
	}

	wg.Wait()

	return errors.Join(errs...)
}

// Copy returns a deep copy of the Subject. If the receiver is nil, a new
// empty Subject is returned.
//
// Returns:
//   - *Subject[T]: The copy of the Subject. Never returns nil.
//
// Does not copy the observers.
func (s *Subject[T]) Copy() *Subject[T] {
	if s == nil {
		return new(Subject[T])
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	return &Subject[T]{
		state: s.state,
	}
}
