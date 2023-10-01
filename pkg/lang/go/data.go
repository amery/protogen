package gogen

import (
	"github.com/amery/protogen/pkg/protogen"
)

// FileData ...
type FileData struct {
	protogen.FileData

	gen *Plugin

	goPkgName string
}

func (fd *FileData) init() error {
	s0, _, err := fd.gen.GoPackage(fd.File())
	if err != nil {
		return err
	}

	fd.goPkgName = s0

	return nil
}

// GoPackage returns the package name for the file
func (fd *FileData) GoPackage() string {
	return fd.goPkgName
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

	return fd, fd.init()
}
