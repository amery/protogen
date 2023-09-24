package nanopb

import (
	"text/template"

	plugin "github.com/amery/protogen/pkg/lang/c99"
	"github.com/amery/protogen/pkg/protogen"
)

const (
	// HeaderFileSuffix indicates the suffix used for generated .h files
	HeaderFileSuffix = "_pb.h"
	// FileSuffix indicates the suffix used for generated .c files
	FileSuffix = "_pb.c"
)

// Generator is our nanopb code generator for C
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

// NewGenerator creates a new nanopb [Generator]
func NewGenerator(opts *Options, gen *protogen.Plugin) (*Generator, error) {
	opts, lp, err := prepareGenerator(opts, gen)
	if err != nil {
		return nil, err
	}

	g := &Generator{
		Plugin: lp,
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
			err = gen.generateFiles(f)
		}
	})

	return err
}

func (gen *Generator) generateFiles(src *protogen.File) error {
	// generate data
	data, err := gen.newGeneratedData(src)
	if err != nil {
		return err
	}

	// output files
	files := generatedFiles([]generatedFile{
		{
			suffix:   HeaderFileSuffix,
			template: "file_h",
		},
		{
			suffix:   FileSuffix,
			template: "file_c",
		},
	})

	// discard all on error
	defer files.DiscardAll()

	// open files
	for i := range files {
		if err := files[i].Open(gen, data); err != nil {
			return err
		}
	}

	// apply templates
	for i := range files {
		if err := files[i].T(gen, data); err != nil {
			return err
		}
	}

	// success
	return files.CloseAll()
}

type generatedFiles []generatedFile

func (m generatedFiles) DiscardAll() {
	for i := range m {
		f := m[i].out
		if f != nil {
			_ = f.Discard()
		}
	}
}

func (m generatedFiles) CloseAll() error {
	var err error

	for i := range m {
		f := m[i].out
		if f != nil {
			e := f.Close()
			if e != nil && err == nil {
				err = e
			}
		}
	}

	return err
}

type generatedFile struct {
	suffix   string
	template string
	out      *protogen.GeneratedFile
	enabled  func(*Options) bool
}

func (g *generatedFile) Enabled(opts *Options) bool {
	switch {
	case g == nil:
		return false
	case g.enabled == nil:
		return true
	default:
		return g.enabled(opts)
	}
}

func (g *generatedFile) Open(gen *Generator, data *GeneratedData) error {
	if !g.Enabled(gen.opts) {
		return nil
	}

	out, err := gen.NewGeneratedFile(data.BaseName + g.suffix)
	if err != nil {
		return err
	}

	g.out = out
	return nil
}

func (g *generatedFile) T(gen *Generator, data any) error {
	if g.out != nil {
		return gen.T(g.template, g.out, data)
	}
	return nil
}
