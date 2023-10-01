package protogen

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"google.golang.org/protobuf/types/pluginpb"
)

// Options specifies callbacks and streams to be used by the protogen [Plugin]
type Options struct {
	// Name indicates the name of the Plugin
	Name string

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
	//   protogen.Run(opts, func(p protogen.Generator) error {
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

	// Features indicates what extra features the plugin supports.
	// 0: None
	// 1: Proto3 Optional
	Features pluginpb.CodeGeneratorResponse_Feature
}

// SetDefaults fills any gap in the Options object
func (opts *Options) SetDefaults() {
	if opts.Name == "" {
		opts.Name = filepath.Base(os.Args[0])
	}

	if IsNil(opts.Stdin) {
		opts.Stdin = os.Stdin
	}

	if IsNil(opts.Stdout) {
		opts.Stdout = os.Stdout
	}

	if IsNil(opts.Stderr) {
		opts.Stderr = os.Stderr
	}

	if IsNil(opts.Logger) {
		prefix := opts.Name + ": "
		opts.Logger = log.New(opts.Stderr, prefix, log.Lmsgprefix)
	}
}

// New allocates a Generator using the Options values
func (opts *Options) New() (*Plugin, error) {
	return NewPlugin(opts, nil)
}

// Run handles the protoc plugin protocol using the provided
// handler and Options values
func (opts *Options) Run(h Handler) error {
	return Run(opts, h)
}
