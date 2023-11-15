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

	switch {
	case e.Err != nil:
		s = fmt.Sprintf("%s (exit:%v)", e.Err.Error(), e.Code)
	case e.Code == 0:
		s = "OK"
	default:
		s = fmt.Sprintf("Exit Code %v", e.Code)
	}

	return s
}

func (e *ExitError) Unwrap() error {
	return e.Err
}

// WithExitCode wraps a fatal error in [ExitError] so callers
// know how to [os.Exit]
func WithExitCode(err error, code int) *ExitError {
	return &ExitError{
		Err:  err,
		Code: code,
	}
}
