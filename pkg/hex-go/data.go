package hex

import (
	plugin "github.com/amery/protogen/pkg/lang/go"
	"github.com/amery/protogen/pkg/protogen"
)

// FileData contains the data used by the templates
// to generate our Go file
type FileData struct {
	plugin.FileData

	opts *Options
}

func (*FileData) init() error {
	return nil
}

func (gen *Generator) newGeneratedData(src *protogen.File) (*FileData, error) {
	fd, err := gen.Plugin.NewFileData(src)
	if err != nil {
		return nil, err
	}

	f := &FileData{
		FileData: fd,

		opts: gen.opts,
	}

	if err := f.init(); err != nil {
		return nil, err
	}

	return f, nil
}
