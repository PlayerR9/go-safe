package queue

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/PlayerR9/go-safe/common"
	sbj "github.com/PlayerR9/go-safe/subject"
)

// Queue is a generic type that represents a thread-safe queue data
// structure without a limited capacity, implemented using a linked list.
//
// An empty queue is created by using the `q := new(Queue[T])` constructor.
type Queue[T any] struct {
	// front and back are pointers to the first and last nodes in the safe queue,
	// respectively.
	front, back *queue_node[T]

	// frontMutex and backMutex are sync.RWMutexes, which are used to ensure that
	// concurrent reads and writes to the front and back nodes are thread-safe.
	mu sync.RWMutex

	// size is the size that observers observe.
	size *sbj.Subject[int]
}

// GoString implements the fmt.GoStringer interface.
func (queue *Queue[T]) GoString() string {
	if queue == nil {
		return "Queue[size=0, values=[]]"
	}

	queue.mu.RLock()
	defer queue.mu.RUnlock()

	size := queue.size.MustGet()
	if size == 0 {
		return "Queue[size=0, values=[]]"
	}

	values := make([]string, 0, size)
	for node := queue.front; node != nil; node = node.next {
		values = append(values, fmt.Sprint(node.value))
	}

	var builder strings.Builder

	builder.WriteString("Queue[size=")
	builder.WriteString(strconv.Itoa(size))
	builder.WriteString(", values=[")
	builder.WriteString(strings.Join(values, ", "))
	builder.WriteString("]]")

	return builder.String()
}

// Enqueue enqueues a value in the queue in a safe way.
//
// Parameters:
//   - value: The value to be enqueued.
//
// Returns:
//   - error: An error if the receiver is nil.
func (queue *Queue[T]) Enqueue(value T) error {
	if queue == nil {
		return common.ErrNilReceiver
	}

	queue.mu.Lock()
	defer queue.mu.Unlock()

	node := &queue_node[T]{
		value: value,
	}

	if queue.back == nil {
		queue.front = node
	} else {
		queue.back.next = node
	}

	queue.back = node

	if queue.size == nil {
		queue.size = new(sbj.Subject[int])
	}

	_ = queue.size.Edit(func(size *int) {
		*size = *size + 1
	})

	return nil
}

// EnqueueMany enqueues multiple values in the queue in a safe way.
//
// Parameters:
//   - values: The values to be enqueued.
//
// Returns:
//   - error: An error if the receiver is nil.
func (queue *Queue[T]) EnqueueMany(values []T) error {
	if len(values) == 0 {
		return nil
	} else if queue == nil {
		return common.ErrNilReceiver
	}

	for _, value := range values {
		_ = queue.Enqueue(value)
	}

	return nil
}

// Dequeue removes and returns the first element in the queue in a safe way.
//
// Returns:
//   - T: The first element in the queue.
//   - error: An error if the dequeue operation fails.
//
// Errors:
//   - ErrEmptyQueue: If the queue is empty.
//   - common.ErrNilReceiver: If the receiver is nil.
func (queue *Queue[T]) Dequeue() (T, error) {
	if queue == nil {
		return *new(T), common.ErrNilReceiver
	}

	queue.mu.Lock()
	defer queue.mu.Unlock()

	if queue.front == nil {
		return *new(T), ErrEmptyQueue
	}

	toRemove := queue.front

	if queue.front.next == nil {
		queue.front = nil
		queue.back = nil
	} else {
		queue.front = queue.front.next
	}

	_ = queue.size.Edit(func(size *int) {
		*size = *size - 1
	})

	return toRemove.value, nil
}

// Peek returns the first element in the queue in a safe way.
//
// Returns:
//   - T: The first element in the queue.
//   - error: An error if the peek operation fails.
//
// Errors:
//   - ErrEmptyQueue: If the queue is empty.
//   - common.ErrNilReceiver: If the receiver is nil.
func (queue *Queue[T]) Peek() (T, error) {
	if queue == nil {
		return *new(T), common.ErrNilReceiver
	}

	queue.mu.RLock()
	defer queue.mu.RUnlock()

	if queue.front == nil {
		return *new(T), ErrEmptyQueue
	}

	return queue.front.value, nil
}

// IsEmpty checks whether the queue is empty.
//
// Returns:
//   - bool: True if the queue is empty, false otherwise.
func (queue *Queue[T]) IsEmpty() bool {
	if queue == nil {
		return true
	}

	queue.mu.RLock()
	defer queue.mu.RUnlock()

	return queue.front == nil
}

// Size returns the size of the queue in a safe way.
//
// Returns:
//   - int: The size of the queue.
func (queue *Queue[T]) Size() int {
	if queue == nil {
		return 0
	}

	queue.mu.RLock()
	defer queue.mu.RUnlock()

	return queue.size.MustGet()
}

// Reset removes all elements from the queue in a safe way.
func (queue *Queue[T]) Reset() {
	if queue == nil {
		return
	}

	queue.mu.Lock()
	defer queue.mu.Unlock()

	if queue.front == nil {
		return // Queue is already empty
	}

	queue.front = nil
	queue.back = nil
	queue.size = nil
}

// Slice returns a copy of the elements in the queue.
//
// Returns:
//   - []T: A copy of the elements in the queue.
func (queue *Queue[T]) Slice() []T {
	if queue == nil {
		return nil
	}

	queue.mu.RLock()
	defer queue.mu.RUnlock()

	size := queue.size.MustGet()

	if size == 0 {
		return nil
	}

	slice := make([]T, 0, size)

	for node := queue.front; node != nil; node = node.next {
		slice = append(slice, node.value)
	}

	return slice
}

// Copy copies the queue in a safe way. It does not copy
// the observers.
//
// Returns:
//   - *Queue[T]: A copy of the queue. Nil if the receiver is nil.
//
// Deprecated: There is no need for this function.
func (queue *Queue[T]) Copy() *Queue[T] {
	if queue == nil {
		return nil
	}

	queue.mu.RLock()
	defer queue.mu.RUnlock()

	q_copy := &Queue[T]{
		size: queue.size.Copy(),
	}

	if queue.front == nil {
		return q_copy
	}

	// First node
	node := &queue_node[T]{
		value: queue.front.value,
	}

	q_copy.front = node
	q_copy.back = node

	// Subsequent nodes
	for n := queue.front.next; n != nil; n = n.next {
		node = &queue_node[T]{
			value: n.value,
		}

		q_copy.back.next = node
		q_copy.back = node
	}

	return q_copy
}

// ObserveSize adds an observer to the size of the queue. Does nothing if the
// function is nil.
//
// Parameters:
//   - fn: The function to be called when the size changes.
//
// Returns:
//   - error: An error if the receiver is nil.
func (queue *Queue[T]) ObserveSize(fn sbj.Action[int]) error {
	if fn == nil {
		return nil
	} else if queue == nil {
		return common.ErrNilReceiver
	}

	o := sbj.FromAction(fn)

	queue.mu.Lock()
	defer queue.mu.Unlock()

	if queue.size == nil {
		queue.size = new(sbj.Subject[int])
	}

	_ = queue.size.Attach(o)

	return nil
}
