package common

import "errors"

var (
	// ErrNilReceiver occurs when a method is called on a receiver who was not
	// expected to be nil. This error can be checked with the == operator.
	//
	// Format:
	// 	"receiver must not be nil"
	ErrNilReceiver error
)

func init() {
	ErrNilReceiver = errors.New("receiver must not be nil")
}

// ErrBadParam occurs when a parameter is bad. (i.e., not a valid value).
type ErrBadParam struct {
	// ParamName is the name of the parameter causing the error.
	ParamName string

	// Msg is the error message describing why the parameter is bad.
	Msg string
}

// Error implements the error interface.
func (e ErrBadParam) Error() string {
	var msg string

	if e.Msg == "" {
		msg = "is invalid"
	} else {
		msg = e.Msg
	}

	if e.ParamName == "" {
		return "parameter " + msg
	} else {
		return "parameter (" + e.ParamName + ") " + msg
	}
}

// NewErrBadParam creates a new ErrBadParam error with the specified parameter name and message.
//
// Parameters:
//   - param_name: The name of the parameter causing the error.
//   - msg: The error message describing why the parameter is bad.
//
// Returns:
//   - error: An instance of ErrBadParam. Never returns nil.
//
// Format:
//
//	"parameter (<param_name>) <msg>"
//
// where:
//   - (<param_name>): The name of the parameter. If empty, it is omitted.
//   - <msg>: The error message describing why the parameter is bad. If empty, "is invalid" is used.
func NewErrBadParam(param_name, msg string) error {
	return &ErrBadParam{
		ParamName: param_name,
		Msg:       msg,
	}
}

// NewErrNilParam is a convenience function that creates a new ErrBadParam error with the specified
// parameter name and the message "must not be nil".
//
// Parameters:
//   - param_name: The name of the parameter causing the error.
//
// Returns:
//   - error: An instance of ErrBadParam. Never returns nil.
//
// Format:
//
//	"parameter (<param_name>) must not be nil"
//
// where:
//   - (<param_name>): The name of the parameter. If empty, it is omitted.
func NewErrNilParam(param_name string) error {
	return &ErrBadParam{
		ParamName: param_name,
		Msg:       "must not be nil",
	}
}
