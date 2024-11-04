package common

import (
	"context"
)

// RejectNilAction takes a pointer to an array of Actions and rejects all
// nil Actions. If the resulting slice is empty, the pointer is set to nil.
// Otherwise, the resulting slice is truncated.
//
// Parameters:
//   - slice: A pointer to an array of Actions.
func RejectNilAction(slice *[]Action) {
	if slice == nil || len(*slice) == 0 {
		return
	}

	var top int

	for _, act := range *slice {
		if act != nil {
			(*slice)[top] = act
			top++
		}
	}

	if top == 0 {
		clear(*slice)
		*slice = nil
	} else {
		clear((*slice)[top:])
		*slice = (*slice)[:top]
	}
}

// Action is a function that takes a context and returns an error.
type Action interface {
	// Run runs the action.
	//
	// Parameters:
	//   - ctx: The context to run the action in.
	//
	// Returns:
	//   - error: An error if the action fails to run.
	Run(ctx context.Context) error
}

// Run runs a list of actions in the given context. If any of the actions fails
// to run, Run will immediately return the error. If the context is canceled or
// times out, Run will return the context's error. Run will not wait for the
// context to be closed to return.
//
// Parameters:
//   - ctx: The context to run the actions in.
//   - acts: The list of actions to run.
//
// Returns:
//   - error: An error if any of the actions fails to run or if the context is
//     canceled or times out.
func Run(ctx context.Context, acts ...Action) error {
	RejectNilAction(&acts)
	if len(acts) == 0 {
		return nil
	} else if ctx == nil {
		return NewErrNilParam("ctx")
	}

	done := ctx.Done()

	for _, act := range acts {
		if act == nil {
			continue
		}

		select {
		case <-done:
			return ctx.Err()
		default:
			err := act.Run(ctx)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
