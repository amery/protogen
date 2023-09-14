package protogen

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"path"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

var (
	_ io.Closer = (*GeneratedFile)(nil)
	_ io.Writer = (*GeneratedFile)(nil)
)

// GeneratedFile implements GeneratedFile
type GeneratedFile struct {
	gen  *Plugin
	buf  *bytes.Buffer
	name string
}

// Name returns the output name associated to this file
func (f *GeneratedFile) Name() string {
	return f.name
}

// Write writes content to the file
func (f *GeneratedFile) Write(b []byte) (int, error) {
	switch {
	case f == nil:
		return 0, fs.ErrInvalid
	case f.buf == nil:
		return 0, fs.ErrClosed
	default:
		return f.buf.Write(b)
	}
}

// Close tells the plugin to emit the file
func (f *GeneratedFile) Close() error {
	switch {
	case f == nil:
		return fs.ErrInvalid
	case f.buf == nil:
		return fs.ErrClosed
	default:
		return f.gen.saveGenerated(f)
	}
}

func (gen *Plugin) saveGenerated(f *GeneratedFile) error {
	var err error

	// double check we are the right instance
	f0, ok := gen.generated[f.name]
	if !ok || f0 != f {
		return fs.ErrInvalid
	}

	b := f.buf.Bytes()
	switch {
	case utf8.Valid(b):
		// append to response
		f := &pluginpb.CodeGeneratorResponse_File{
			Name:    proto.String(f.name),
			Content: proto.String(string(b)),
		}

		gen.resp.File = append(gen.resp.File, f)
	default:
		// invalid UTF-8
		err = Wrap(ErrInvalidUTF8Content, f.name)
	}

	// remove
	delete(f.gen.generated, f.name)
	f.buf = nil
	return err
}

// Discard tells the plugin not to emit this file
func (f *GeneratedFile) Discard() error {
	switch {
	case f == nil:
		return fs.ErrInvalid
	case f.buf == nil:
		return fs.ErrClosed
	default:
		return f.gen.discardGenerated(f)
	}
}

func (gen *Plugin) discardGenerated(f *GeneratedFile) error {
	// double check we are the right instance
	f0, ok := gen.generated[f.name]
	if !ok || f0 != f {
		return fs.ErrInvalid
	}

	// purge buffer and discard
	delete(gen.generated, f.name)
	f.buf = nil
	return nil
}

// NewGeneratedFile creates a new output file
func (gen *Plugin) NewGeneratedFile(format string, args ...any) (*GeneratedFile, error) {
	var err error

	name, ok := getGeneratedName(format, args...)
	if !ok {
		err = ErrInvalidName
	} else if _, ok = gen.generated[name]; ok {
		err = fs.ErrExist
	}

	if err != nil {
		return nil, &fs.PathError{
			Op:   "create",
			Path: name,
			Err:  err,
		}
	}

	f := &GeneratedFile{
		gen:  gen,
		buf:  new(bytes.Buffer),
		name: name,
	}

	gen.generated[name] = f

	return f, nil
}

func getGeneratedName(s string, args ...any) (string, bool) {
	if len(args) > 0 {
		s = fmt.Sprintf(s, args...)
	}

	switch {
	case s != path.Clean(s):
		return s, false
	case filepath.IsAbs(s):
		return s, false
	case strings.ContainsRune(s, '\\'):
		return s, false
	default:
		return s, true
	}
}
