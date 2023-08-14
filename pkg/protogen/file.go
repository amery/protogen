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
	_ File = (*FileDescriptor)(nil)
)

// FileDescriptor implements a File
type FileDescriptor struct {
	gen *Plugin
	dp  *descriptorpb.FileDescriptorProto

	generate bool

	enums    []Enum
	messages []Message
}

// Request returns the [pluginpb.CodeGeneratorRequest] received by
// the [Plugin]
func (f *FileDescriptor) Request() *pluginpb.CodeGeneratorRequest {
	return f.gen.Request()
}

// Proto returns the underlying protobuf structure
func (f *FileDescriptor) Proto() *descriptorpb.FileDescriptorProto {
	return f.dp
}

// Generate indicates the file was directly specified when
// calling protoc
func (f *FileDescriptor) Generate() bool {
	return f.generate
}

func (gen *Plugin) setFilesGenerate(files ...string) error {
	for _, filename := range files {
		if !gen.setFileGenerate(filename) {
			return &fs.PathError{
				Op:   "generate",
				Path: filename,
				Err:  fs.ErrNotExist,
			}
		}
	}

	return nil
}

func (gen *Plugin) setFileGenerate(filename string) bool {
	file := gen.getFileByName(filename)
	if file != nil {
		file.generate = true
		return true
	}
	return false
}

// Name returns the full file name of proto file
func (f *FileDescriptor) Name() string {
	return optional(f.dp.Name, "")
}

// Base returns the name of the proto file including directory
// but excluding extensions
func (f *FileDescriptor) Base() string {
	name := f.Name()
	if ext := filepath.Ext(name); ext != "" {
		return strings.TrimSuffix(name, ext)
	}
	return name
}

// Package returns the package name associated to this file
func (f *FileDescriptor) Package() string {
	return optional(f.dp.Package, "")
}

// PackageDirectory returns the package name associated to this file
// converted to a directory path
func (f *FileDescriptor) PackageDirectory() string {
	s := optional(f.dp.Package, "")
	if s == "" {
		return "."
	}

	sep := string([]rune{os.PathSeparator})
	return strings.ReplaceAll(s, ".", sep)
}

// Dependencies returns the source proto files this one depends on
func (f *FileDescriptor) Dependencies() []File {
	out := make([]File, len(f.dp.Dependency))
	for i, fn := range f.dp.Dependency {
		out[i] = f.gen.getFileByName(fn)
	}
	return out
}

// Files returns a slice of all source proto files
func (gen *Plugin) Files() []File {
	return gen.files
}

// ForEachFile calls a function for each source proto file
func (gen *Plugin) ForEachFile(fn func(File)) {
	for _, f := range gen.files {
		fn(f)
	}
}

// FileByName returns a source proto file by name
func (gen *Plugin) FileByName(filename string) File {
	return gen.getFileByName(filename)
}

func (gen *Plugin) getFileByName(filename string) *FileDescriptor {
	for _, f := range gen.files {
		if f.Name() == filename {
			return f.(*FileDescriptor)
		}
	}

	return nil
}

func (gen *Plugin) loadFiles(files ...*descriptorpb.FileDescriptorProto) {
	for _, dp := range files {
		f := &FileDescriptor{
			dp:  dp,
			gen: gen,
		}

		gen.files = append(gen.files, f)
	}
}
