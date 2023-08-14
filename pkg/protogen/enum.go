package protogen

import (
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

var (
	_ ProtoTyper = (*Enum)(nil)
	_ ProtoTyper = (*EnumValue)(nil)
)

// Enum represents an enumeration type
type Enum struct {
	file *File
	dp   *descriptorpb.EnumDescriptorProto

	values []*EnumValue
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
func (p *Enum) File() *File {
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

// Values returns the possible values for this type
func (p *Enum) Values() []*EnumValue {
	return p.values
}

func (p *Enum) init() {
	out := make([]*EnumValue, 0, len(p.dp.Value))
	for _, dp := range p.dp.Value {
		p := &EnumValue{
			enum: p,
			dp:   dp,
		}

		out = append(out, p)
	}
	p.values = out
}

// EnumValue represents a possible value of a [Enum]
type EnumValue struct {
	enum *Enum
	dp   *descriptorpb.EnumValueDescriptorProto
}

// Request returns the [pluginpb.CodeGeneratorRequest] received by
// the [Plugin]
func (p *EnumValue) Request() *pluginpb.CodeGeneratorRequest {
	return p.enum.Request()
}

// Proto returns the underlying protobuf structure
func (p *EnumValue) Proto() *descriptorpb.EnumValueDescriptorProto {
	return p.dp
}

// File returns the [File] associates to this value type
func (p *EnumValue) File() *File {
	return p.enum.File()
}

// Package returns the package name associated to this type
func (p *EnumValue) Package() string {
	return p.enum.Package()
}

// Enum returns the [Enum] associates to this value type
func (p *EnumValue) Enum() *Enum {
	return p.enum
}

// Name returns the relative name of this value type
func (p *EnumValue) Name() string {
	return optional(p.dp.Name, "")
}

// FullName returns the fully qualified name of this value type
func (p *EnumValue) FullName() string {
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
func (f *File) Enums() []*Enum {
	if f.enums == nil {
		f.loadEnums()
	}
	return f.enums
}

// EnumByName finds a [Enum] by name
func (f *File) EnumByName(name string) *Enum {
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
	out := make([]*Enum, 0, len(f.dp.EnumType))
	for _, dp := range f.dp.EnumType {
		if dp == nil {
			continue
		}

		p := &Enum{
			file: f,
			dp:   dp,
		}

		p.init()
		out = append(out, p)
	}
	f.enums = out
}
