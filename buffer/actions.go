package buffer

import (
	"context"

	"github.com/PlayerR9/go-safe/common"
)

// resetAct is an action that resets the Buffer.
type resetAct[T any] struct{}

// Run implements the common.Action interface.
func (act *resetAct[T]) Run(ctx context.Context) error {
	c, err := fromContext[T](ctx)
	if err != nil {
		return err
	}

	c.buffer.Reset()

	return nil
}

// Reset removes all elements from the Buffer, effectively resetting
// it to an empty state. Precalculated elements are kept as they are no longer
// in the buffer but in the channel. It locks the firstMutex to ensure
// thread-safety during the operation.
//
// This method is safe for concurrent use by multiple goroutines.
//
// Returns:
//   - common.Action: The reset action. Never returns nil.
func Reset[T any]() common.Action {
	return &resetAct[T]{}
}

/* // IsClosed implements the Runner interface.
func (b *Buffer[T]) IsClosed() bool {
	return b == nil || b.locker == nil
} */

// sendAct is an action that sends a message to the Buffer.
type sendAct[T any] struct {
	// msg is the message to send.
	msg T
}

// Run implements the common.Action interface.
func (act *sendAct[T]) Run(ctx context.Context) error {
	c, err := fromContext[T](ctx)
	if err != nil {
		return err
	}

	return c.buffer.Send(act.msg)
}

// Send sends a message to the Buffer.
//
// Parameters:
//   - msg: The message to send.
//
// Returns:
//   - common.Action: The send action. Never returns nil.
func Send[T any](msg T) common.Action {
	return &sendAct[T]{
		msg: msg,
	}
}

// receiveAct is an action that receives a message from the Buffer.
type receiveAct[T any] struct {
	// msg is the destination to receive the message.
	msg *T
}

// Run implements the common.Action interface.
func (act *receiveAct[T]) Run(ctx context.Context) error {
	c, err := fromContext[T](ctx)
	if err != nil {
		return err
	}

	msg, err := c.buffer.Receive()
	if err != nil {
		return err
	}

	*act.msg = msg

	return nil
}

// Receive receives a message from the Buffer.
//
// Parameters:
//   - dest: The destination to receive the message.
//
// Returns:
//   - common.Action: The receive action. Nil if dest is nil.
func Receive[T any](dest *T) common.Action {
	if dest == nil {
		return nil
	}

	return &receiveAct[T]{
		msg: dest,
	}
}
