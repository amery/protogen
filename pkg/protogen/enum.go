package protogen

import (
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

var (
	_ EnumDescriptor = (*Enum)(nil)
)

// Enum is the implementation of EnumDescriptor
type Enum struct {
	file *File
	dp   *descriptorpb.EnumDescriptorProto
}

// Request returns the [pluginpb.CodeGeneratorRequest] received by
// the [Plugin]
func (p *Enum) Request() *pluginpb.CodeGeneratorRequest {
	return p.file.Request()
}

// Proto returns the underlying protobuf structure
func (p *Enum) Proto() *descriptorpb.EnumDescriptorProto {
	return p.dp
}

// File returns the [File] that defines this type
func (p *Enum) File() FileDescriptor {
	return p.file
}

// Package returns the package name associated to this type
func (p *Enum) Package() string {
	return p.file.Package()
}

// Name returns the relative name of this type
func (p *Enum) Name() string {
	return optional(p.dp.Name, "")
}

// FullName returns the fully qualified name of this type
func (p *Enum) FullName() string {
	s0 := p.file.Package()
	s1 := p.Name()

	switch {
	case s0 != "" && s1 != "":
		return s0 + "." + s1
	default:
		return s1
	}
}

// Enums returns all the [Enum] types defined on this file
func (f *File) Enums() []EnumDescriptor {
	return f.enums
}

// EnumByName finds a [Enum] by name
func (f *File) EnumByName(name string) EnumDescriptor {
	pkgname, name, _ := SplitName(name)

	switch {
	case pkgname != "" && f.Package() != pkgname:
		// wrong package
		return nil
	case name == "":
		// no name
		return nil
	default:
		for _, p := range f.Enums() {
			if name == p.Name() {
				// match
				return p
			}
		}

		return nil
	}
}

func (f *File) loadEnums() {
	out := make([]EnumDescriptor, 0, len(f.dp.EnumType))
	for _, dp := range f.dp.EnumType {
		if dp == nil {
			continue
		}

		p := &Enum{
			file: f,
			dp:   dp,
		}

		out = append(out, p)
	}
	f.enums = out
}
