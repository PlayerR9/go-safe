package internal

import "errors"

var (
	// ErrAlreadyClosed occurs when the buffer is already closed.
	//
	// Format:
	//   "buffer is already closed"
	ErrAlreadyClosed error
)

func init() {
	ErrAlreadyClosed = errors.New("buffer is already closed")
}
