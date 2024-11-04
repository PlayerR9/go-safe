package internal

import (
	"fmt"
	"sync"

	"github.com/PlayerR9/go-safe/common"
	lls "github.com/PlayerR9/go-safe/queue"
	sbj "github.com/PlayerR9/go-safe/subject"
)

// Buffer is a thread-safe, generic data structure that allows multiple
// goroutines to produce and consume elements in a synchronized manner.
// It is implemented as a queue and uses channels to synchronize the
// goroutines.
//
// Information: To close the buffer, just close the send-only channel.
// Once that is done, a cascade of events will happen:
//   - The goroutine that listens for incoming messages will stop listening
//     and exit.
//   - The goroutine that sends messages from the Buffer to the receive
//     channel will stop sending messages once the Buffer is empty, and then exit.
//   - The Buffer will be cleaned up.
//
// Of course, a Close method is also provided to manually close the Buffer but
// it is not necessary to call it if the send-only channel is closed.
//
// To create an empty Buffer, use the `b := new(Buffer[T])` constructor.
type Buffer[T any] struct {
	// q is a pointer to the SafeQueue that stores the elements of the Buffer.
	q *lls.Queue[T]

	// sendTo is a channel that receives messages and sends them to the Buffer.
	sendTo chan T

	// receiveFrom is a channel that receives messages from the Buffer and
	// sends them to the consumer.
	receiveFrom chan T

	// wg is a WaitGroup that is used to wait for the goroutines to finish.
	wg sync.WaitGroup

	// locker is a pointer to the RWSafe that synchronizes the Buffer.
	locker *sbj.Locker[BufferCondition]
}

// listenForIncomingMessages is a method of the Buffer type that listens for
// incoming messages from the receiveChannel and enqueues them in the Buffer.
//
// It must be run in a separate goroutine to avoid blocking the main thread.
func (b *Buffer[T]) listenForIncomingMessages() {
	defer b.wg.Done()

	for msg := range b.sendTo {
		_ = b.q.Enqueue(msg)
	}

	_ = b.locker.ChangeValue(IsRunning, false)
}

// sendMessagesFromBuffer is a method of the Buffer type that sends
// messages from the Buffer to the sendChannel.
//
// It must be run in a separate goroutine to avoid blocking the main thread.
func (b *Buffer[T]) sendMessagesFromBuffer() {
	defer b.wg.Done()

	for {
		value, err := b.locker.Get(IsRunning)
		if err != nil {
			panic(fmt.Errorf("unable to get whether the buffer is running or not: %w", err))
		}

		if !value {
			break
		}

		fn := func(m map[BufferCondition]bool) error {
			for {
				isEmpty, ok := b.sendSingleMessage()
				if !ok || isEmpty {
					break
				}
			}

			if m[IsRunning] {
				return nil
			} else {
				return sbj.ErrStop
			}
		}

		err = b.locker.DoFunc(fn)
		if err == nil {
			continue
		} else if err == sbj.ErrStop {
			break
		} else {
			panic(err)
		}
	}

	for {
		isEmpty, _ := b.sendSingleMessage()
		if isEmpty {
			break
		}
	}

	b.locker = nil
	b.q = nil
}

// sendSingleMessage is a method of the Buffer type that sends a single message
// from the Buffer to the send channel.
//
// Returns:
//   - bool: A boolean indicating if the queue is empty.
//   - bool: A boolean indicating if a message was sent successfully.
func (b *Buffer[T]) sendSingleMessage() (bool, bool) {
	msg, err := b.q.Peek()
	if err != nil {
		return true, true
	}

	select {
	case b.receiveFrom <- msg:
		_, err := b.q.Dequeue()
		if err != nil {
			return true, false
		}

		return false, true
	default:
		return false, false
	}
}

// Start implements the Runner interface.
func (b *Buffer[T]) Start() error {
	if b == nil {
		return common.ErrNilReceiver
	} else if b.locker != nil {
		// already started
		return nil
	}

	b.locker = sbj.NewLocker[BufferCondition]()
	b.locker.SetSubject(IsEmpty, true, true)
	b.locker.SetSubject(IsRunning, true, true)

	b.q = new(lls.Queue[T])

	err := b.q.ObserveSize(func(val int) error {
		err := b.locker.ChangeValue(IsEmpty, val == 0)
		return err
	})
	if err != nil {
		return err
	}

	b.sendTo = make(chan T)
	b.receiveFrom = make(chan T)

	b.wg.Add(2)

	go b.listenForIncomingMessages()
	go b.sendMessagesFromBuffer()

	return nil
}

// Close implements the Runner interface.
func (b *Buffer[T]) Close() {
	if b == nil || b.sendTo == nil {
		return
	}

	close(b.sendTo)
	b.sendTo = nil

	b.wg.Wait()

	close(b.receiveFrom)
	b.receiveFrom = nil
}

// Reset removes all elements from the Buffer, effectively resetting
// it to an empty state. Precalculated elements are kept as they are no longer
// in the buffer but in the channel. It locks the firstMutex to ensure
// thread-safety during the operation.
//
// This method is safe for concurrent use by multiple goroutines.
func (b *Buffer[T]) Reset() {
	if b == nil || b.q == nil {
		return
	}

	b.q.Reset()
}

// Send implements the Sender interface.
func (b *Buffer[T]) Send(msg T) error {
	if b == nil {
		return common.ErrNilReceiver
	}

	if b.sendTo == nil {
		return ErrAlreadyClosed
	}

	b.sendTo <- msg

	return nil
}

// Receive implements the Receiver interface.
func (b *Buffer[T]) Receive() (T, error) {
	if b == nil {
		return *new(T), common.ErrNilReceiver
	} else if b.receiveFrom == nil {
		return *new(T), ErrAlreadyClosed
	}

	msg, ok := <-b.receiveFrom
	if !ok {
		return *new(T), ErrAlreadyClosed
	}

	return msg, nil
}
