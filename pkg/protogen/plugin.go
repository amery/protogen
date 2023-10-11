package protogen

import (
	"google.golang.org/protobuf/types/pluginpb"
)

// Plugin is the protoc code generator engine
type Plugin struct {
	options Options
	req     *pluginpb.CodeGeneratorRequest
	resp    pluginpb.CodeGeneratorResponse

	params    map[string]string
	files     []*File
	generated map[string]*GeneratedFile
}

func (gen *Plugin) init(req *pluginpb.CodeGeneratorRequest) error {
	gen.options.SetDefaults()

	// extra features supported by the plugin
	gen.resp.SupportedFeatures = PointerOrNil(uint64(gen.options.Features))

	if req == nil {
		var err error

		// read encoded request from stdin
		req, err = UnmarshalCodeGeneratorRequest(gen.options.Stdin)
		if err != nil {
			return Wrap(err, "UnmarshalCodeGeneratorRequest")
		}
	}

	return gen.loadRequest(req)
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

// NewPlugin creates and initializes a new protoc Plugin handler.
// If a CodeGeneratorRequest isn't provided, Options.Stdin will
// be deserialized instead.
func NewPlugin(opts *Options, req *pluginpb.CodeGeneratorRequest) (*Plugin, error) {
	if opts == nil {
		opts = &Options{}
	}

	gen := &Plugin{
		options:   *opts,
		params:    make(map[string]string),
		generated: make(map[string]*GeneratedFile),
	}

	// always return the *Plugin so it can be used to respond
	return gen, gen.init(req)
}
