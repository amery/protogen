package nanogrpc

import "github.com/amery/protogen/pkg/protogen"

// GeneratedFile contains the data used by the templates
// to generate a header file
type GeneratedFile struct {
	opts *Options
	src  *protogen.File
}

func (*GeneratedFile) init() error {
	return nil
}

func (gen *Generator) newGeneratedFileData(src *protogen.File) (*GeneratedFile, error) {
	f := &GeneratedFile{
		opts: gen.opts,
		src:  src,
	}

	if err := f.init(); err != nil {
		return nil, err
	}

	return f, nil
}
