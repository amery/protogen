package nanogrpc

import (
	"text/template"

	"github.com/amery/protogen/pkg/protogen"
)

// FileSuffix indicates the suffix used for generated files
const FileSuffix = "_nanogrpc.go"

// Generator is our nanopb code generator for C
type Generator struct {
	*protogen.Plugin

	opts *Options
	t    *template.Template
}

// NewGenerator creates a new nanopb [Generator]
func NewGenerator(opts *Options, gen *protogen.Plugin) (*Generator, error) {
	if opts == nil {
		opts = new(Options)
	}

	if err := opts.SetDefaults(); err != nil {
		return nil, err
	}

	g := &Generator{
		Plugin: gen,
		opts:   opts,
	}

	if err := g.withTemplates(); err != nil {
		return nil, err
	}

	return g, nil
}

// Run handles the protoc plugin protocol for the nanopb generator
func (gen *Generator) Run() error {
	var err error

	gen.ForEachFile(func(f *protogen.File) {
		switch {
		case err != nil:
			// aborting
		case !f.Generate():
			// skip
		default:
			// generate
			err = gen.generateFile(f)
		}
	})

	return err
}

func (gen *Generator) generateFile(src *protogen.File) error {
	// generate data
	data, err := gen.newGeneratedFileData(src)
	if err != nil {
		return err
	}

	// output file
	out, err := gen.NewGeneratedFile("%s%s", src.Base(), FileSuffix)
	if err != nil {
		return err
	}

	// apply template
	err = gen.T("file", out, data)
	if err != nil {
		_ = out.Discard()
		return err
	}

	// success
	return out.Close()
}
