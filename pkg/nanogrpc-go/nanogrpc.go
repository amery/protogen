// Package nanogrpc generates Go code for nano-grpc
package nanogrpc

import (
	"github.com/amery/protogen/pkg/protogen"
)

// RunPlugin handles the protoc plugin protocol for the nanogrpc for Go generator
func RunPlugin(pgo *protogen.Options) error {
	opts, err := NewOptions(pgo)
	if err != nil {
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
