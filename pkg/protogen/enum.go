package protogen

import (
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

var (
	_ Enum      = (*EnumDescriptor)(nil)
	_ EnumValue = (*EnumValueDescriptor)(nil)
)

// EnumDescriptor is the implementation of Enum
type EnumDescriptor struct {
	file *FileDescriptor
	dp   *descriptorpb.EnumDescriptorProto

	values []EnumValue
}

// Request returns the [pluginpb.CodeGeneratorRequest] received by
// the [Plugin]
func (p *EnumDescriptor) Request() *pluginpb.CodeGeneratorRequest {
	return p.file.Request()
}

// Proto returns the underlying protobuf structure
func (p *EnumDescriptor) Proto() *descriptorpb.EnumDescriptorProto {
	return p.dp
}

// File returns the [File] that defines this type
func (p *EnumDescriptor) File() File {
	return p.file
}

// Package returns the package name associated to this type
func (p *EnumDescriptor) Package() string {
	return p.file.Package()
}

// Name returns the relative name of this type
func (p *EnumDescriptor) Name() string {
	return optional(p.dp.Name, "")
}

// FullName returns the fully qualified name of this type
func (p *EnumDescriptor) FullName() string {
	s0 := p.file.Package()
	s1 := p.Name()

	switch {
	case s0 != "" && s1 != "":
		return s0 + "." + s1
	default:
		return s1
	}
}

// Values returns the possible values for this type
func (p *EnumDescriptor) Values() []EnumValue {
	return p.values
}

func (p *EnumDescriptor) init() {
	out := make([]EnumValue, 0, len(p.dp.Value))
	for _, dp := range p.dp.Value {
		p := &EnumValueDescriptor{
			enum: p,
			dp:   dp,
		}

		out = append(out, p)
	}
	p.values = out
}

// EnumValueDescriptor is the implementation of EnumValue
type EnumValueDescriptor struct {
	enum *EnumDescriptor
	dp   *descriptorpb.EnumValueDescriptorProto
}

// Request returns the [pluginpb.CodeGeneratorRequest] received by
// the [Plugin]
func (p *EnumValueDescriptor) Request() *pluginpb.CodeGeneratorRequest {
	return p.enum.Request()
}

// Proto returns the underlying protobuf structure
func (p *EnumValueDescriptor) Proto() *descriptorpb.EnumValueDescriptorProto {
	return p.dp
}

// File returns the [File] associates to this value type
func (p *EnumValueDescriptor) File() File {
	return p.enum.File()
}

// Package returns the package name associated to this type
func (p *EnumValueDescriptor) Package() string {
	return p.enum.Package()
}

// Enum returns the [Enum] associates to this value type
func (p *EnumValueDescriptor) Enum() Enum {
	return p.enum
}

// Name returns the relative name of this value type
func (p *EnumValueDescriptor) Name() string {
	return optional(p.dp.Name, "")
}

// FullName returns the fully qualified name of this value type
func (p *EnumValueDescriptor) FullName() string {
	s0 := p.enum.FullName()
	s1 := p.Name()

	switch {
	case s0 != "" && s1 != "":
		return s0 + "." + s1
	default:
		return s1
	}
}

// Enums returns all the [Enum] types defined on this file
func (f *FileDescriptor) Enums() []Enum {
	if f.enums == nil {
		f.loadEnums()
	}
	return f.enums
}

// EnumByName finds a [Enum] by name
func (f *FileDescriptor) EnumByName(name string) Enum {
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

func (f *FileDescriptor) loadEnums() {
	out := make([]Enum, 0, len(f.dp.EnumType))
	for _, dp := range f.dp.EnumType {
		if dp == nil {
			continue
		}

		p := &EnumDescriptor{
			file: f,
			dp:   dp,
		}

		p.init()
		out = append(out, p)
	}
	f.enums = out
}
