package protogen

import (
	"github.com/pkg/errors"

	"google.golang.org/protobuf/types/pluginpb"
)

// Plugin is the protoc code generator engine
type Plugin struct {
	options Options
	req     *pluginpb.CodeGeneratorRequest
	resp    pluginpb.CodeGeneratorResponse
}

func (gen *Plugin) init(req *pluginpb.CodeGeneratorRequest) error {
	gen.options.SetDefaults()

	if req == nil {
		var err error

		// read encoded request from stdin
		req, err = UnmarshalCodeGeneratorRequest(gen.options.Stdin)
		if err != nil {
			return errors.Wrap(err, "UnmarshalCodeGeneratorRequest")
		}
	}

	return gen.unsafeLoadRequest(req)
}

// Print logs an error in the manner of fmt.Print
func (gen *Plugin) Print(v ...any) {
	gen.options.Logger.Print(v...)
}

// Println logs an error in the manner of fmt.Println
func (gen *Plugin) Println(v ...any) {
	gen.options.Logger.Println(v...)
}

// Printf logs an error in the manner of fmt.Printf
func (gen *Plugin) Printf(format string, v ...any) {
	gen.options.Logger.Printf(format, v...)
}

// NewPlugin creates and initialises a new protoc Plugin handler.
// If a CodeGeneratorRequest isn't provided, Options.Stdin will
// be unmarshalled instead.
func NewPlugin(opts *Options, req *pluginpb.CodeGeneratorRequest) (*Plugin, error) {
	if opts == nil {
		opts = &Options{}
	}

	gen := &Plugin{
		options: *opts,
	}

	if err := gen.init(req); err != nil {
		return nil, err
	}

	return gen, nil
}
