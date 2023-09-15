// Package main implements a pluginpb.CodeGeneratorRequest dumper
package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/amery/protogen/pkg/protogen"
)

func generate(gen *protogen.Plugin) error {
	var errs protogen.ErrAggregation

	if name, ok := gen.Param("req"); ok {
		if err := saveRawRequest(gen, name); err != nil {
			errs.AppendWrapped(err, "req=%s", name)
		}
	}

	gen.ForEachFile(func(f *protogen.File) {
		if err := saveRawFile(gen, f); err != nil {
			errs.AppendWrapped(err, f.Name())
		}
	})

	return errs.AsError()
}

func main() {
	arg0 := filepath.Base(os.Args[0])
	opts := protogen.Options{
		Logger: log.New(os.Stderr, arg0+": ", 0),
	}

	err := opts.Run(generate)
	if err != nil {
		log.Fatal(err)
	}
}
