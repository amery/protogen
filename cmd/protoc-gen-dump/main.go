// Package main implements a pluginpb.CodeGeneratorRequest dumper
package main

import (
	"io"
	"log"
	"os"

	"google.golang.org/protobuf/types/pluginpb"

	"github.com/amery/protogen/pkg/protogen"
	"github.com/amery/protogen/pkg/protogen/plugin"
)

func generate(gen *protogen.Plugin) error {
	var errs protogen.ErrAggregation

	if name, ok := gen.Param("req"); ok {
		if err := saveRawRequest(gen, name); err != nil {
			errs.AppendWrapped(err, "req=%s", name)
		}
	}

	return errs.AsError()
}

func run(in io.ReadCloser, out io.WriteCloser) error {
	opts := protogen.Options{
		Logger:   log.New(os.Stderr, plugin.CmdName()+": ", 0),
		Stdin:    in,
		Stdout:   out,
		Features: pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL,
	}

	return opts.Run(generate)
}

func main() {
	pc := &plugin.Config{
		RunE: run,
	}

	rootCmd, err := plugin.NewRoot(pc)
	if err == nil {
		err = rootCmd.Execute()
	}

	switch e := err.(type) {
	case *plugin.ExitError:
		os.Exit(e.Code)
	case nil:
		os.Exit(0)
	default:
		log.Fatal(err)
	}
}
