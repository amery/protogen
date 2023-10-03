// Package main implements a pluginpb.CodeGeneratorRequest dumper
package main

import (
	"errors"
	"io"
	"log"
	"os"

	"google.golang.org/protobuf/types/pluginpb"

	"github.com/amery/protogen/pkg/protogen"
	"github.com/amery/protogen/pkg/protogen/plugin"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	cmdName = plugin.CmdName()

	rawRequestFlag    *pflag.Flag
	rawRequestValue   *string
	skipRawWriteValue *bool
)

func setExtraRootFlags(cmd *cobra.Command) {
	flags := cmd.LocalFlags()

	rawRequestValue = flags.StringP("raw-request", "w", "",
		"override where the raw request is saved ('none' or empty to disable)")
	rawRequestFlag = flags.Lookup("raw-request")

	skipRawWriteValue = flags.BoolP("skip-raw-write", "W", false,
		"do not save raw request")
}

func getRawRequestName(gen *protogen.Plugin) (string, bool) {
	switch {
	case *skipRawWriteValue:
		// skip
		return "", true
	case rawRequestFlag.Changed:
		// given via command line
		return *rawRequestValue, true
	default:
		// protoc option
		return gen.Param("raw_request")
	}
}

func rawOutputName(gen *protogen.Plugin) (string, bool) {
	name, ok := getRawRequestName(gen)
	if ok {
		switch name {
		case "", "none", "false":
			// disable
			return "", true
		default:
			// use given
			return name, true
		}
	}

	// find name
	gen.ForEachFile(func(f *protogen.File) {
		if f.Generate() && name == "" {
			// use
			name = f.Base()
		}
	})

	return name, name != ""
}

func generate(gen *protogen.Plugin) error {
	// get name
	name, ok := rawOutputName(gen)
	switch {
	case !ok:
		return errors.New("couldn't determine name for raw request")
	case name != "":
		// save as
		if err := saveRawRequest(gen, name); err != nil {
			return err
		}
	}

	return nil
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
	var err error
	pc := &plugin.Config{
		Name: cmdName,
		RunE: run,
	}

	rootCmd, err := plugin.NewRoot(pc)
	if err == nil {
		setExtraRootFlags(rootCmd)

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
