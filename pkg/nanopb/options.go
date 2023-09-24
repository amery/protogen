package nanopb

import (
	"github.com/amery/protogen/pkg/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

// Options determines how the generator operates
type Options struct {
	protogen.Options
}

// SetDefaults fills any gap in the Options object
func (opts *Options) SetDefaults() error {
	opts.ParamFunc = opts.SetParam
	opts.Features = pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL
	opts.Options.SetDefaults()

	return nil
}

// SetParam sets parameters passed by protoc
func (opts *Options) SetParam(key, _ string) error {
	opts.Logger.Printf("%s:%q: %s", "warning", key, "unknown parameter")
	return nil
}

// NewOptions creates a new [Options] optionally using the given [protogen.Options]
// as reference
func NewOptions(pgo *protogen.Options) (*Options, error) {
	return prepareOptions(nil, pgo)
}

func prepareOptions(opts *Options, pgo *protogen.Options) (*Options, error) {
	if opts == nil {
		opts = new(Options)
	}

	if pgo != nil {
		opts.Options = *pgo
	}

	if err := opts.SetDefaults(); err != nil {
		return nil, err
	}

	return opts, nil
}
