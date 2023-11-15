package plugin

import "fmt"

// ExitError is an error that is expected to be
// handled directly via os.Exit(e.Code)
type ExitError struct {
	Err  error
	Code int
}

func (e *ExitError) Error() string {
	var s string

	code := e.ExitCode()

	switch {
	case e.Err != nil:
		s = fmt.Sprintf("%s (exit:%v)", e.Err.Error(), code)
	case code == 0:
		s = "OK"
	default:
		s = fmt.Sprintf("Exit Code %v", code)
	}

	return s
}

func (e *ExitError) Unwrap() error {
	return e.Err
}

// ExitCode returns the value to be used on [os.Exit] calls
func (e *ExitError) ExitCode() int {
	return e.Code & 0x7f
}

// WithExitCode wraps a fatal error in [ExitError] so callers
// know how to [os.Exit]
func WithExitCode(err error, code int) *ExitError {
	return &ExitError{
		Err:  err,
		Code: code & 0x7f,
	}
}
