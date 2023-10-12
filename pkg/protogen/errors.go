package protogen

import (
	"bytes"
	"errors"
	"fmt"

	"google.golang.org/protobuf/types/descriptorpb"
)

var (
	// ErrNotImplemented tells certain logic hasn't been implemented yet
	ErrNotImplemented = errors.New("not implemented")
	// ErrInvalidName tells the requested file name isn't acceptable
	ErrInvalidName = errors.New("invalid name")
	// ErrInvalidUTF8Content tells the plugin generated unacceptable content
	ErrInvalidUTF8Content = errors.New("invalid UTF-8 content generated")

	// ErrUnknownParam tells the plug-in parameter isn't recognized
	ErrUnknownParam = errors.New("unknown protoc option")
	// ErrInvalidParam tells the plug-in parameter is known but the value is not
	// acceptable.
	ErrInvalidParam = errors.New("invalid protoc option value")
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

// PluginError is a wrapped error referencing a .proto file
type PluginError struct {
	Path string
	Hint string
	Err  error
}

func (e PluginError) Error() string {
	s0 := e.Path

	s1 := e.Hint
	if s1 == "" && e.Err != nil {
		s1 = e.Err.Error()
	}

	if s1 == "" {
		s1 = "unspecified error"
	}

	switch {
	case s0 == "":
		return s1
	default:
		return s0 + ": " + s1
	}
}

func (e PluginError) Unwrap() error {
	return e.Err
}

// NewPluginError creates a new PluginError taking the Path from
// [descriptorpb.FileDescriptorProto#Name]
func NewPluginError(file *descriptorpb.FileDescriptorProto, err error, hint string) *PluginError {
	var name string
	if file != nil && file.Name != nil {
		name = *file.Name
	}

	return &PluginError{Path: name, Err: err, Hint: hint}
}

// ErrAggregation bundles multiple errors
type ErrAggregation struct {
	Errs []error
}

// Errors returns the included errors
func (e *ErrAggregation) Errors() []error {
	if e == nil || len(e.Errs) == 0 {
		return nil
	}

	return e.Errs
}

func (e *ErrAggregation) Error() string {
	var buf bytes.Buffer

	if len(e.Errs) == 0 {
		return "OK"
	}

	for _, err := range e.Errs {
		_, _ = fmt.Fprint(&buf, "* ", err.Error(), "\n")
	}

	return buf.String()
}

// AsError returns nil if there are no errors stored,
// or itself if there are.
func (e *ErrAggregation) AsError() error {
	if len(e.Errs) == 0 {
		return nil
	}
	return e
}

// Append stores an error on the aggregation
func (e *ErrAggregation) Append(err error) {
	if err != nil {
		e.Errs = append(e.Errs, err)
	}
}

// AppendWrapped stores an error on the aggregation wrapped with a hint
func (e *ErrAggregation) AppendWrapped(err error, hint string, args ...any) {
	if err != nil {
		err = Wrap(err, hint, args...)

		e.Errs = append(e.Errs, err)
	}
}
