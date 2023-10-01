package gogen

import (
	"errors"
)

var (
	// ErrInvalidGoPackageName ...
	ErrInvalidGoPackageName = errors.New("invalid go package name")
	// ErrNoGoPackageName ...
	ErrNoGoPackageName = errors.New("failed to infer go package name")
	// ErrGoPackageNotExist ...
	ErrGoPackageNotExist = errors.New("unknown package")
)
