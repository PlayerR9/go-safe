package subject

import (
	"fmt"
	"sync"

	"github.com/PlayerR9/go-safe/common"
)

// Conditioner is an interface that represents a condition.
type Conditioner interface {
	~int

	fmt.Stringer
}

// Locker is a thread-Subject locker that allows multiple goroutines to wait for a condition.
type Locker[T Conditioner] struct {
	// cond is the condition variable.
	cond *sync.Cond

	// elems is the list of elements.
	subjects map[T]*Subject[bool]

	// mu is the mutex to synchronize map access.
	mu sync.RWMutex
}

// NewLocker creates a new Locker.
//
// Use Locker.Set for observer boolean predicates.
//
// Parameters:
//   - keys: The keys to initialize the locker.
//
// Returns:
//   - *Locker[T]: A new Locker.
//
// Behaviors:
//   - All the predicates are initialized to true.
func NewLocker[T Conditioner]() *Locker[T] {
	l := &Locker[T]{
		cond:     sync.NewCond(&sync.Mutex{}),
		subjects: make(map[T]*Subject[bool]),
	}

	return l
}

// SetSubject adds a new subject to the locker.
//
// Parameters:
//   - key: The key to add.
//   - subject: The subject to add.
//   - broadcast: A flag indicating whether the subject should broadcast or signal.
//
// Behaviors:
//   - If the subject is nil, it will not be added.
//   - It overwrites the existing subject if the key already exists.
func (l *Locker[T]) SetSubject(key T, value bool, broadcast bool) {
	subject := New(value)

	var fn Action[bool]

	if broadcast {
		fn = func(b bool) error {
			l.cond.L.Lock()
			defer l.cond.L.Unlock()

			l.cond.Broadcast()

			return nil
		}
	} else {
		fn = func(b bool) error {
			l.cond.L.Lock()
			defer l.cond.L.Unlock()

			l.cond.Signal()

			return nil
		}
	}

	o := FromAction(fn)
	_ = subject.Attach(o)

	l.mu.Lock()
	defer l.mu.Unlock()

	l.subjects[key] = subject
}

// ChangeValue changes the value of a subject.
//
// Parameters:
//   - key: The key to change the value.
//   - value: The new value.
//
// Returns:
//   - error: An error if the receiver is nil or if the key does not exist.
//
// Errors:
//   - common.ErrNilReceiver: If the receiver is nil.
//   - ErrKeyNotExist: If the key does not exist.
func (l *Locker[T]) ChangeValue(key T, value bool) error {
	if l == nil {
		return common.ErrNilReceiver
	}

	l.mu.Lock()
	subject, ok := l.subjects[key]
	l.mu.Unlock()

	if !ok {
		return ErrKeyNotExist
	}

	_ = subject.Set(value)
	return nil
}

// hasFalse is a private method that checks if at least one of the conditions is false.
//
// Returns:
//   - map[T]bool: A copy of the map of conditions.
//   - bool: True if at least one of the conditions is false, false otherwise.
func (l *Locker[T]) hasFalse() (map[T]bool, bool) {
	l.mu.RLock()

	map_copy := make(map[T]bool, len(l.subjects))

	for key, value := range l.subjects {
		map_copy[key] = value.MustGet()
	}
	l.mu.RUnlock()

	for _, value := range map_copy {
		if !value {
			return map_copy, true
		}
	}

	return map_copy, false
}

// Get returns the value of a predicate.
//
// Parameters:
//   - key: The key to get the value.
//
// Returns:
//   - bool: The value of the predicate.
//   - error: An error if the receiver is nil or if the key does not exist.
//
// Errors:
//   - common.ErrNilReceiver: If the receiver is nil.
//   - ErrKeyNotExist: If the key does not exist.
func (l *Locker[T]) Get(key T) (bool, error) {
	if l == nil {
		return false, common.ErrNilReceiver
	}

	l.mu.RLock()
	val, ok := l.subjects[key]
	l.mu.RUnlock()

	if ok {
		return val.MustGet(), nil
	} else {
		return false, ErrKeyNotExist
	}
}

// DoFunc is a function that executes a function while waiting for the condition to be false.
//
// Parameters:
//   - f: The function to execute that takes a map of conditions as a parameter and returns
//     an error if something goes wrong. ErrStop will be returned if the function should stop.
//
// Returns:
//   - error: An error if the do function fails.
//
// Errors:
//   - common.ErrNilReceiver: If the receiver is nil.
//   - ErrStop: If the function should stop.
func (l *Locker[T]) DoFunc(f func(map[T]bool) error) error {
	if l == nil {
		return common.ErrNilReceiver
	}

	l.cond.L.Lock()

	var map_copy map[T]bool
	var ok bool

	for {
		map_copy, ok = l.hasFalse()
		if ok {
			l.cond.L.Unlock()
			break
		}

		l.cond.Wait()
	}

	err := f(map_copy)
	return err
}
