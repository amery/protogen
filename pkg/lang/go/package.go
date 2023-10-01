package gogen

import (
	"strings"

	"github.com/amery/protogen/pkg/protogen"
	"github.com/amery/protogen/pkg/text"
	"google.golang.org/protobuf/types/descriptorpb"
)

// Package ...
type Package struct {
	src *protogen.File

	GoPackageName string
	GoPackagePath string

	ProtoFile string
}

func (p *Package) loadFile(src *protogen.File) error {
	p.src = src
	if p.GoPackageName == "" {
		name, path, err := getGoPackage(src.Proto())
		if err != nil {
			return err
		}

		p.GoPackageName = name
		p.GoPackagePath = path
	}

	return nil
}

func (gen *Plugin) mustPackage(protoFile string, src *protogen.File) (*Package, error) {
	if protoFile == "" && src != nil {
		protoFile = src.Name()
	}

	p, ok := gen.pkgMap[protoFile]
	switch {
	case !ok:
		// new
		return gen.newPackage(protoFile, src)
	case p.src == src:
		// done
	case p.src == nil:
		// update
		if err := p.loadFile(src); err != nil {
			return nil, err
		}
	default:
		// *protogen.File is singleton
		panic("unreachable")
	}

	return p, nil
}

func (gen *Plugin) newPackage(protoFile string, src *protogen.File) (*Package, error) {
	p := &Package{
		ProtoFile: protoFile,
	}

	if src != nil {
		if err := p.loadFile(src); err != nil {
			return nil, err
		}
	}

	gen.pkgMap[protoFile] = p
	return p, nil
}

// GetPackageByProtoFile ...
func (gen *Plugin) GetPackageByProtoFile(protoName string) (*Package, bool) {
	p, ok := gen.pkgMap[protoName]
	return p, ok
}

// GetPackageByFile ...
func (gen *Plugin) GetPackageByFile(src *protogen.File) (*Package, bool) {
	if src == nil {
		return nil, false
	}

	return gen.GetPackageByProtoFile(src.Name())
}

// SetGoPackage ...
func (gen *Plugin) SetGoPackage(protoName, goPackageOption string) error {
	name, path, ok := goPackageFromOption(goPackageOption)
	if !ok {
		err := &protogen.PluginError{
			Path: protoName,
			Err:  ErrInvalidGoPackageName,
			Hint: goPackageOption,
		}
		return err
	}

	p, err := gen.mustPackage(protoName, nil)
	if err != nil {
		return err
	}

	p.GoPackageName = name
	p.GoPackagePath = path
	return nil
}

// GoPackage returns the Go package associated with a given .proto.
// derived from the following:
//   - The MprotoFile=goPackage protoc parameters, if provided
//   - The `go_package` identified option, if provided.
//   - The basename of the package import path, if provided.
//   - The package statement in the .proto file, if present.
//   - The basename of the .proto file, without extension.
//
// `;` can be used after the go import path to specify the package name
func (gen *Plugin) GoPackage(src *protogen.File) (name string, path string, err error) {
	p, err := gen.mustPackage("", src)
	if err != nil {
		return "", "", err
	}

	return p.GoPackageName, p.GoPackagePath, nil
}

func getGoPackage(file *descriptorpb.FileDescriptorProto) (name string, path string, err error) {
	for _, fn := range []func(*descriptorpb.FileDescriptorProto) (string, string, error){
		getGoPackageFromOption,  // identified `go_package`
		getGoPackageFromPath,    // path to .proto
		getGoPackageFromPackage, // proto package import path
		getGoPackageFromFile,    // basename of .proto

	} {
		name, path, err = fn(file)
		switch {
		case err != nil:
			return "", "", err
		case name != "":
			return name, path, nil
		}
	}

	// unreachable
	return "", "", protogen.NewPluginError(file, ErrNoGoPackageName, "")
}

func getGoPackageFromOption(file *descriptorpb.FileDescriptorProto) (name string, path string, err error) {
	if p := file.GetOptions().GoPackage; p != nil && *p != "" {
		var ok bool

		name, path, ok = goPackageFromOption(*p)
		if !ok {
			err = protogen.NewPluginError(file, ErrInvalidGoPackageName, *p)
		}
	}

	return name, path, err
}

func getGoPackageFromPath(file *descriptorpb.FileDescriptorProto) (name string, path string, err error) {
	if p := file.Name; p != nil && *p != "" {
		var ok bool

		// TODO: ToSlash()?
		path, _, ok = text.CutLastRune(*p, '/')
		if ok {
			name, path, ok = goPackageValidate("", path)
			if !ok {
				err = protogen.NewPluginError(file, ErrInvalidGoPackageName, *p)
			}
			return name, path, err
		}

		// dir-less
	}
	return "", "", nil
}

func getGoPackageFromFile(file *descriptorpb.FileDescriptorProto) (name string, path string, err error) {
	if p := file.Name; p != nil && *p != "" {
		var ok bool

		name, path, ok = goPackageFromFile(*p)
		if !ok {
			err = protogen.NewPluginError(file, ErrInvalidGoPackageName, *p)
		}
	}
	return name, path, err
}

// getGoPackageFromPackage attempts to get go package from the proto package
func getGoPackageFromPackage(file *descriptorpb.FileDescriptorProto) (name string, path string, err error) {
	if p := file.Package; p != nil && *p != "" {
		var ok bool

		name, path, ok = goPackageFromPackage(*p)
		if !ok {
			err = protogen.NewPluginError(file, ErrInvalidGoPackageName, *p)
		}
	}
	return name, path, err
}

// goPackageFromOption splits a `path;name` option entry
func goPackageFromOption(s string) (name string, path string, ok bool) {
	path, name, ok = text.CutLastRune(s, ';')
	if !ok {
		name = ""
		path = s
	}

	return goPackageValidate(name, path)
}

func goPackageFromFile(s string) (name string, path string, ok bool) {
	path, ok = text.CutSuffix(s, ".proto")
	if ok {
		return goPackageValidate("", path)
	}
	return "", "", false
}

// goPackageFromPackage extracts go package name and import path from
// the .proto's package name
func goPackageFromPackage(s string) (name string, path string, ok bool) {
	return goPackageValidate("", strings.ReplaceAll(s, ".", "/"))
}

// goPackageValidate confirms the chosen name and import path are valid.
// if name is missing, the last element of the path will be chosen
func goPackageValidate(sName string, sPath string) (name string, path string, ok bool) {
	switch {
	case sName == "" && sPath == "":
		return "", "", false
	case sName == "":
		var ok bool

		_, name, ok = text.CutLastRune(sPath, '/')
		if !ok {
			name = sPath
		}

		path = sPath
	case sPath == "":
		name, path = sPath, sPath
	}

	// TODO: validate further
	return name, path, true
}
