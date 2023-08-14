// Package protogen assists on making Go protoc plugins for any language
package protogen

import (
	"google.golang.org/protobuf/types/pluginpb"
)

// Handler uses [Plugin] to generate code
type Handler func(*Plugin) error

// ProtoTyper is the common abstraction for types defined on a proto file
type ProtoTyper interface {
	// Request returns the received [pluginpb.CodeGeneratorRequest]
	Request() *pluginpb.CodeGeneratorRequest

	// File returns the [File] that defines this type
	File() *File
	// Package returns the package name associated to this type
	Package() string
	// Name returns the relative name of this type
	Name() string
	// FullName returns the fully qualified name of this type
	FullName() string
}

// Run handles the protoc plugin protocol using the provided
// Options and handler.
// if Options is nil, a new one will be created with
// default values.
func Run(opts *Options, h Handler) error {
	gen, err := NewPlugin(opts, nil)
	if err != nil {
		gen.Print(err)
		_, _ = gen.WriteError(err)
		return err
	}

	err = h(gen)
	if err != nil {
		_, _ = gen.WriteError(err)
		return err
	}

	_, err = gen.Write()
	if err != nil {
		gen.Print(err)
		return err
	}

	return nil
}
