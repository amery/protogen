package protogen

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

var (
	_ FileDescriptor = (*File)(nil)
)

// File is the base foundation of a FileDescriptor
type File struct {
	gen *Plugin
	fdp *descriptorpb.FileDescriptorProto

	generate bool
}

// Request returns the [pluginpb.CodeGeneratorRequest] received by
// the [Plugin]
func (f *File) Request() *pluginpb.CodeGeneratorRequest {
	return f.gen.Request()
}

// Proto returns the underlying protobuf structure
func (f *File) Proto() *descriptorpb.FileDescriptorProto {
	return f.fdp
}

// Generate indicates the file was directly specified when
// calling protoc
func (f *File) Generate() bool {
	return f.generate
}

func (gen *Plugin) setFilesGenerate(files ...string) error {
	for _, fn1 := range files {
		if !gen.setFileGenerate(fn1) {
			return &fs.PathError{
				Op:   "generate",
				Path: fn1,
				Err:  fs.ErrNotExist,
			}
		}
	}

	return nil
}

func (gen *Plugin) setFileGenerate(fn1 string) bool {
	for _, f := range gen.files {
		file, ok := f.(*File)
		if !ok {
			continue
		}

		fn2, ok := optional2(file.fdp.Name, "")
		if !ok {
			continue
		}

		if fn1 == fn2 {
			file.generate = true
			return true
		}
	}

	return false
}

// Name returns the full file name of proto file
func (f *File) Name() string {
	return optional(f.fdp.Name, "")
}

// Base returns the name of the proto file including directory
// but excluding extensions
func (f *File) Base() string {
	name := f.Name()
	if ext := filepath.Ext(name); ext != "" {
		return strings.TrimSuffix(name, ext)
	}
	return name
}

// Package returns the package name associated to this file
func (f *File) Package() string {
	return optional(f.fdp.Package, "")
}

// PackageDirectory returns the package name associated to this file
// converted to a directory path
func (f *File) PackageDirectory() string {
	s := optional(f.fdp.Package, "")
	if s == "" {
		return "."
	}

	sep := string([]rune{os.PathSeparator})
	return strings.ReplaceAll(s, ".", sep)
}

// Files returns a slice of all source proto files
func (gen *Plugin) Files() []FileDescriptor {
	return gen.files
}

// ForEachFile calls a function for each source proto file
func (gen *Plugin) ForEachFile(fn func(FileDescriptor)) {
	for _, f := range gen.files {
		fn(f)
	}
}

func (gen *Plugin) loadFiles(files ...*descriptorpb.FileDescriptorProto) {
	for _, fdp := range files {
		f := &File{
			fdp: fdp,
			gen: gen,
		}

		gen.files = append(gen.files, f)
	}
}
