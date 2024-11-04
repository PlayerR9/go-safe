package queue

import "errors"

var (
	// ErrEmptyQueue occurs when the queue is empty.
	//
	// Format:
	//   "queue is empty"
	ErrEmptyQueue error
)

func init() {
	ErrEmptyQueue = errors.New("queue is empty")
}
