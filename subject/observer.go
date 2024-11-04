package subject

// Action is a change action.
//
// Parameters:
//   - change: The change in the subject.
//
// Returns:
//   - error: An error if the observer fails to notify.
type Action[T any] func(change T) error

// reactiveObserver is a reactive observer of a subject.
type reactiveObserver[T any] struct {
	// action is the action to perform when the observer is notified.
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
	action Action[T]
}

// Notify implements the Observer interface.
func (o reactiveObserver[T]) Notify(change T) error {
	err := o.action(change)
	return err
}

// Cleanup implements the Observer interface.
func (o *reactiveObserver[T]) Cleanup() {
	if o == nil {
		return
	}

	o.action = nil
}

// FromAction returns a new reactive observer of a subject given a change action.
//
// Parameters:
//   - action: The action to perform when the observer is notified.
//
// Returns:
//   - Observer[T]: The new observer. Returns nil if the action is nil.
func FromAction[T any](action Action[T]) Observer[T] {
	if action == nil {
		return nil
	}

	return &reactiveObserver[T]{
		action: action,
	}
}
