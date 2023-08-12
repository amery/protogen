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
