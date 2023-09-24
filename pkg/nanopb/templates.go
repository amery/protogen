package nanopb

import (
	"embed"
	"fmt"
	"io"
	"io/fs"

	"text/template"
)

//go:embed templates/**.gotmpl
var content embed.FS

func withTemplates(root *template.Template) (*template.Template, error) {
	if root == nil {
		root = template.New("")
	}

	dir, err := fs.Sub(content, "templates")
	if err != nil {
		return nil, err
	}

	return root.ParseFS(dir, "*.gotmpl")
}

func (gen *Generator) withTemplates() error {
	root := template.New("").Funcs(
		// TODO: add template functions
		template.FuncMap{},
	)

	t, err := withTemplates(root)
	if err != nil {
		return err
	}

	gen.t = t
	return nil
}

// T executes a template
func (gen *Generator) T(name string, out io.Writer, data any) error {
	var t *template.Template

	if gen.t != nil {
		t = gen.t.Lookup(name + ".gotmpl")
	}

	if t == nil {
		return &template.ExecError{
			Name: name,
			Err:  fmt.Errorf("template %q not found", name),
		}
	}

	return t.Execute(out, data)
}
