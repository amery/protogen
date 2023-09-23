package plugin

import "fmt"

// ExitError is an error that is expected to be
// handled directly via os.Exit(e.Code)
type ExitError struct {
	Code int
}

func (e *ExitError) Error() string {
	var s string

	switch {
	case e.Code == 0:
		s = "OK"
	default:
		s = fmt.Sprintf("Exit Code %v", e.Code)
	}

	return s
}
