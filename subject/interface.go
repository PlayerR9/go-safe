package subject

// Observer is an observer of a subject.
type Observer[T any] interface {
	// Notify notifies the observer of a change in the subject.
	//
	// Parameters:
	//   - change: The change in the subject.
	//
	// Returns:
	//   - error: An error if the observer fails to notify.
	//
	// Errors:
	// 	- common.ErrNilReceiver: If the receiver is nil.
	// 	- any other error: If the observer fails to notify.
	Notify(change T) error

	// Cleanup cleans up the observer.
	Cleanup()
}
