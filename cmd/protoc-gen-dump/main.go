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

var cmdName = plugin.CmdName()

func generate(gen *protogen.Plugin) error {
	return saveRawRequest(gen)
}

func run(in io.ReadCloser, out io.WriteCloser) error {
	opts := protogen.Options{
		Logger:   log.New(os.Stderr, cmdName+": ", 0),
		Stdin:    in,
		Stdout:   out,
		Features: pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL,
	}

	return opts.Run(generate)
}

func main() {
	pc := &plugin.Config{
		Name: cmdName,
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
