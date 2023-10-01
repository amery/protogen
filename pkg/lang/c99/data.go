package c99gen

import (
	"github.com/amery/protogen/pkg/protogen"
)

// FileData ...
type FileData struct {
	protogen.FileData

	gen *Plugin
}

// NewFileData ...
func (gen *Plugin) NewFileData(src *protogen.File) (FileData, error) {
	fd0, err := gen.Plugin.NewFileData(src)
	if err != nil {
		return FileData{}, err
	}

	fd := FileData{
		FileData: fd0,

		gen: gen,
	}

	return fd, nil
}
