// Package nanopb generates code using nanopb
package nanopb

import "github.com/amery/protogen/pkg/protogen"

// RunPlugin handles the protoc plugin protocol for the nanopb generator
func RunPlugin(pgo *protogen.Options) error {
	opts, err := NewOptions(pgo)
	if err != nil {
		return err
	}

	return Run(opts)
}

// Run handles the protoc plugin protocol for the nanopb generator
func (opts *Options) Run() error {
	return Run(opts)
}

// Run handles the protoc plugin protocol for the nanopb generator
func Run(opts *Options) error {
	var err error

	if opts, err = prepareOptions(opts, nil); err != nil {
		return err
	}

	return opts.Options.Run(func(gen *protogen.Plugin) error {
		g, err := NewGenerator(opts, gen)
		if err != nil {
			return err
		}

		return g.Run()
	})
}
