// Package main provides a protoc generator for C using nanopb
package main

import (
	"io"
	"log"
	"os"

	"github.com/amery/protogen/pkg/nanopb"
	"github.com/amery/protogen/pkg/protogen"
	"github.com/amery/protogen/pkg/protogen/plugin"
)

var cmdName = plugin.CmdName()

func run(in io.ReadCloser, out io.WriteCloser) error {
	opts := &protogen.Options{
		Name:   cmdName,
		Stdin:  in,
		Stdout: out,
	}

	return nanopb.RunPlugin(opts)
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
