package subject

import "errors"

var (
	// ErrKeyNotExist occurs when the key does not exist.
	//
	// Format:
	//   "key does not exist"
	ErrKeyNotExist error

	// ErrStop occurs when the do function of the locker should stop.
	//
	// Format:
	//   "should stop"
	ErrStop error
)

func init() {
	ErrKeyNotExist = errors.New("key does not exist")

	ErrStop = errors.New("should stop")
}
