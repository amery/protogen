// Package protogen assists on making Go protoc plugins for any language
package protogen

import (
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

// Handler uses Generator to generate code
type Handler func(Generator) error

// RequestDescriptor is the abstraction of an element on the Request
type RequestDescriptor interface {
	// Request returns the received [pluginpb.CodeGeneratorRequest]
	Request() *pluginpb.CodeGeneratorRequest
}

// Generator is the interface implemented by our Plugin for the Handler
type Generator interface {
	RequestDescriptor

	// Print logs an error in the manner of fmt.Print
	Print(...any)
	// Println logs an error in the manner of fmt.Println
	Println(...any)
	// Printf logs an error in the manner of fmt.Printf
	Printf(string, ...any)

	// Param returns the value of a parameter if specified
	Param(string) (string, bool)
	// Params returns all specified parameters
	Params() map[string]string

	// Files returns a slice of all source proto files
	Files() []FileDescriptor
	// ForEachFile calls a function for each source proto file
	ForEachFile(func(FileDescriptor))
	// FileByName returns a source proto file by name
	FileByName(string) FileDescriptor
}

// FileDescriptor represents a source proto file
type FileDescriptor interface {
	RequestDescriptor

	// Proto returns the underlying protobuf structure
	Proto() *descriptorpb.FileDescriptorProto

	// Generate indicates the file was directly specified when
	// calling protoc
	Generate() bool

	// Name returns the full file name of proto file
	Name() string
	// Base returns the name of the proto file including directory
	// but excluding extensions
	Base() string
	// Package returns the package name associated to this file
	Package() string
	// PackageDirectory returns the package name associated to this file
	// converted to a directory path
	PackageDirectory() string
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
