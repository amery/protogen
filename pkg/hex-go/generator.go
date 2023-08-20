package hex

import (
	"text/template"

	plugin "github.com/amery/protogen/pkg/lang/go"
	"github.com/amery/protogen/pkg/protogen"
)

const (
	// FileSuffix indicates the suffix used for generated files
	FileSuffix = "_hex.go"
)

// Generator is our hexagonal code generator for go
type Generator struct {
	*plugin.Plugin

	opts *Options
	t    *template.Template
}

func prepareGenerator(opts *Options, gen *protogen.Plugin) (*Options, *plugin.Plugin, error) {
	opts, err := prepareOptions(opts, nil)
	if err != nil {
		return nil, nil, err
	}

	lp, err := plugin.NewPlugin(gen)
	if err != nil {
		return nil, nil, err
	}

	return opts, lp, nil
}

// NewGenerator creates a new go-hex [Generator]
func NewGenerator(opts *Options, gen *protogen.Plugin) (*Generator, error) {
	opts, gp, err := prepareGenerator(opts, gen)
	if err != nil {
		return nil, err
	}

	g := &Generator{
		Plugin: gp,
		opts:   opts,
	}

	if err := g.withTemplates(); err != nil {
		return nil, err
	}

	return g, nil
}

// Run handles the protoc plugin protocol for the go-hex generator
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
	data, err := gen.newGeneratedData(src)
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
