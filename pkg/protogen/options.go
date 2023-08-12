package protogen

import (
	"io"
	"log"
	"os"
)

// Options specifies callbacks and streams to be used by the protogen [Plugin]
type Options struct {
	// If ParamFunc is non-nil, it will be called with each unknown
	// generator parameter.
	//
	// Plugins for protoc can accept parameters from the command line,
	// passed in the --<generator>_out protoc, separated from the output
	// directory with a colon; e.g.,
	//
	//   --foo_out=<param1>=<value1>,<param2>=<value2>:<output_directory>
	//
	// Parameters passed in this fashion as a comma-separated list of
	// key=value pairs will be passed to the ParamFunc.
	//
	// The (flag.FlagSet).Set method matches this function signature,
	// so parameters can be converted into flags as in the following:
	//
	//   var flags flag.FlagSet
	//   value := flags.Bool("param", false, "")
	//   opts := &protogen.Options{
	//     ParamFunc: flags.Set,
	//   }
	//   protogen.Run(opts, func(p *protogen.Plugin) error {
	//     if *value { ... }
	//   })
	ParamFunc func(name, value string) error

	// Stdin is the source of the encoded [pluginpb.CodeGeneratorRequest]
	Stdin io.Reader
	// Stdout is where we write the encoded [pluginpb.CodeGeneratorResponse]
	Stdout io.Writer
	// Stderr is where we write the logs if Logger isn't specified
	Stderr io.Writer
	// Logger is optional [log.Logger] to use for errors. If not specified
	// one will be built using Stderr
	Logger *log.Logger
}

// SetDefaults fills any gap in the Options object
func (opts *Options) SetDefaults() {
	if opts.Stdin == nil {
		opts.Stdin = os.Stdin
	}

	if opts.Stdout == nil {
		opts.Stdout = os.Stdout
	}

	if opts.Stderr == nil {
		opts.Stderr = os.Stderr
	}

	if opts.Logger == nil {
		opts.Logger = log.New(opts.Stderr, "protogen: ", 0)
	}
}

// NewPlugin allocates a Plugin using the Options values
func (opts *Options) NewPlugin() (*Plugin, error) {
	return NewPlugin(opts, nil)
}

// Run handles the protoc plugin protocol using the provided
// handler and Options values
func (opts *Options) Run(fn func(Generator) error) error {
	return Run(opts, fn)
}
