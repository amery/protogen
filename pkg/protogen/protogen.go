// Package protogen assists on making Go protoc plugins for any language
package protogen

import (
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

// Handler uses Generator to generate code
type Handler func(Generator) error

// ProtoRequester is the abstraction of an element on the Request
type ProtoRequester interface {
	// Request returns the received [pluginpb.CodeGeneratorRequest]
	Request() *pluginpb.CodeGeneratorRequest
}

// Generator is the interface implemented by our Plugin for the Handler
type Generator interface {
	ProtoRequester

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
	Files() []File
	// ForEachFile calls a function for each source proto file
	ForEachFile(func(File))
	// FileByName returns a source proto file by name
	FileByName(string) File
}

// File represents a source proto file
type File interface {
	ProtoRequester

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

	// Dependencies returns the source proto files this one depends on
	Dependencies() []File

	// Enums returns all the [Enum] types defined on this file
	Enums() []Enum
	// EnumByName finds a [Enum] by name
	EnumByName(string) Enum
}

// ProtoTyper is the common abstraction for types defined on a proto file
type ProtoTyper interface {
	ProtoRequester

	// File returns the [File] that defines this type
	File() File
	// Package returns the package name associated to this type
	Package() string
	// Name returns the relative name of this type
	Name() string
	// FullName returns the fully qualified name of this type
	FullName() string
}

// Enum represents an Enum type
type Enum interface {
	ProtoTyper

	// Proto returns the underlying protobuf structure
	Proto() *descriptorpb.EnumDescriptorProto
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
