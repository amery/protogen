package nanopb

import (
	"os"
	"path/filepath"

	"github.com/amery/protogen/pkg/protogen"
)

// GeneratedData contains the data used by the templates
// to generate a header file
type GeneratedData struct {
	opts *Options
	src  *protogen.File

	BaseName string

	ProtoCGenCmd    string
	ProtoFile       string
	ProtoParameters string

	HeaderSentinel string
	HeaderInclude  []string
	SourceInclude  []string
}

func (f *GeneratedData) init() error {
	f.BaseName = f.src.Base()
	f.ProtoCGenCmd = filepath.Base(os.Args[0])
	f.ProtoFile = f.src.Name()

	return nil
}

func (gen *Generator) newGeneratedData(src *protogen.File) (*GeneratedData, error) {
	f := &GeneratedData{
		opts: gen.opts,
		src:  src,
	}

	if err := f.init(); err != nil {
		return nil, err
	}

	return f, nil
}
