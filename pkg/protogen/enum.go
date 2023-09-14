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
	msg  *Message
	dp   *descriptorpb.EnumDescriptorProto

	values   []*EnumValue
	min, max int32
}

// Request returns the [pluginpb.CodeGeneratorRequest] received by
// the [Plugin]
func (p *Enum) Request() *pluginpb.CodeGeneratorRequest {
	return p.File().Request()
}

// Proto returns the underlying protobuf structure
func (p *Enum) Proto() *descriptorpb.EnumDescriptorProto {
	return p.dp
}

// File returns the [File] that defines this type
func (p *Enum) File() *File {
	switch {
	case p.file != nil:
		return p.file
	case p.msg != nil:
		return p.msg.File()
	default:
		panic("unreachable")
	}
}

// Package returns the package name associated to this type
func (p *Enum) Package() string {
	return p.File().Package()
}

// Name returns the relative name of this type
func (p *Enum) Name() string {
	return optional(p.dp.Name, "")
}

// FullName returns the fully qualified name of this type
func (p *Enum) FullName() string {
	var s0 string
	switch {
	case p.msg != nil:
		s0 = p.msg.FullName()
	case p.file != nil:
		s0 = p.file.Package()
	}

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

// Minimum returns the minimum numeric value used on this Enum
func (p *Enum) Minimum() int {
	return int(p.min)
}

// Maximum returns the maximum numeric value used on this Enum
func (p *Enum) Maximum() int {
	return int(p.max)
}

func (p *Enum) init() {
	var next int32

	p.values = make([]*EnumValue, 0, len(p.dp.Value))

	for _, dp := range p.dp.Value {
		next = p.newValue(dp, next)
	}

	// sort values
	Sort(p.values, func(a, b *EnumValue) bool {
		return a.Number() < b.Number()
	})
}

// EnumValue represents a possible value of a [Enum]
type EnumValue struct {
	enum *Enum
	dp   *descriptorpb.EnumValueDescriptorProto

	number int32
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

// Number returns the integer representation of the EnumValue
func (p *EnumValue) Number() int {
	return int(p.number)
}

func (p *Enum) newValue(dp *descriptorpb.EnumValueDescriptorProto, next int32) int32 {
	cur := optional(dp.Number, next)

	v := &EnumValue{
		enum:   p,
		dp:     dp,
		number: cur,
	}

	switch {
	case len(p.values) == 0:
		// first
		p.min, p.max = cur, cur
	case cur > p.max:
		// new max
		p.max = cur
	case cur < p.min:
		// new min
		p.min = cur
	}

	p.values = append(p.values, v)
	return cur + 1
}

// Enums returns all the [Enum] types local to this Message
func (p *Message) Enums() []*Enum {
	if p.enums == nil {
		p.enums = loadEnums(Enum{msg: p}, p.dp.EnumType)
	}
	return p.enums
}

// Enums returns all the [Enum] types defined on this file
func (f *File) Enums() []*Enum {
	if f.enums == nil {
		f.enums = loadEnums(Enum{file: f}, f.dp.EnumType)
	}
	return f.enums
}

// EnumByName finds a [Enum] by name
func (f *File) EnumByName(name string) *Enum {
	pkgName, name, _ := SplitName(name)

	switch {
	case pkgName != "" && f.Package() != pkgName:
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

func loadEnums(ref Enum, enums []*descriptorpb.EnumDescriptorProto) []*Enum {
	out := make([]*Enum, 0, len(enums))
	for _, dp := range enums {
		if dp == nil {
			// TODO: log error
			continue
		}

		q := ref
		q.dp = dp

		q.init()
		out = append(out, &q)
	}

	// sort enums by name
	Sort(out, func(a, b *Enum) bool {
		return a.Name() < b.Name()
	})

	return out
}
