// Package protogen assists on making Go protoc plugins for any language
package protogen

import "google.golang.org/protobuf/types/pluginpb"

// Generator is the interface implemented by our Plugin for the Handler
type Generator interface {
	// Print logs an error in the manner of fmt.Print
	Print(...any)
	// Println logs an error in the manner of fmt.Println
	Println(...any)
	// Printf logs an error in the manner of fmt.Printf
	Printf(string, ...any)

	// Request returns the received [pluginpb.CodeGeneratorRequest]
	Request() *pluginpb.CodeGeneratorRequest
}

// Run handles the protoc plugin protocol using the provided
// Options and handler.
// if Options is nil, a new one will be created with
// default values.
func Run(opts *Options, fn func(Generator) error) error {
	gen, err := NewPlugin(opts, nil)
	if err != nil {
		gen.Print(err)
		_, _ = gen.WriteError(err)
		return err
	}

	err = fn(gen)
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
