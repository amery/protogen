package protogen

import (
	"errors"
	"fmt"
)

var (
	// ErrNotImplemented tells certain logic hasn't been implemented yet
	ErrNotImplemented = errors.New("not implemented")
)

// WrappedError is a simple wrapped error container
type WrappedError struct {
	Hint string
	Err  error
}

func (e WrappedError) Error() string {
	switch {
	case e.Hint == "":
		return e.Err.Error()
	case e.Err == nil:
		return e.Hint
	default:
		return e.Hint + ": " + e.Err.Error()
	}
}

func (e WrappedError) Unwrap() error {
	return e.Err
}

// Wrap wraps an error with a formatted hint message
func Wrap(err error, hint string, args ...any) *WrappedError {
	if err == nil {
		return nil
	}

	if len(args) > 0 {
		hint = fmt.Sprintf(hint, args...)
	}

	return &WrappedError{
		Err:  err,
		Hint: hint,
	}
}
