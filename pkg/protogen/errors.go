package protogen

import (
	"errors"
	"fmt"
)

var (
	// ErrNotImplemented tells certain logic hasn't been implemented yet
	ErrNotImplemented = errors.New("not implemented")
	// ErrInvalidName tells the requested file name isn't acceptable
	ErrInvalidName = errors.New("invalid name")
	// ErrInvalidUTF8Content tells the plugin generated unacceptable content
	ErrInvalidUTF8Content = errors.New("invalid UTF-8 content generated")
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
func Wrap(err error, hint string, args ...any) error {
	if err == nil {
		return nil
	}

	if len(args) > 0 {
		hint = fmt.Sprintf(hint, args...)
	}

	if hint == "" {
		return err
	}

	return &WrappedError{
		Err:  err,
		Hint: hint,
	}
}
