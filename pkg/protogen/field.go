package protogen

import (
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

var _ ProtoTyper = (*Field)(nil)

// Field represents a field on a Message
type Field struct {
	msg *Message
	dp  *descriptorpb.FieldDescriptorProto
}

// Request returns the [pluginpb.CodeGeneratorRequest] received by
// the [Plugin]
func (p *Field) Request() *pluginpb.CodeGeneratorRequest {
	return p.File().Request()
}

// Proto returns the underlying protobuf structure
func (p *Field) Proto() *descriptorpb.FieldDescriptorProto {
	return p.dp
}

// File returns the [File] that defines this type
func (p *Field) File() *File {
	return p.msg.File()
}

// Package returns the package name associated to this type
func (p *Field) Package() string {
	return p.File().Package()
}

// Name returns the relative name of this field
func (p *Field) Name() string {
	return optional(p.dp.Name, "")
}

// FullName returns the fully qualified name of this type
func (p *Field) FullName() string {
	s0 := p.msg.FullName()
	s1 := p.Name()

	switch {
	case s0 != "" && s1 != "":
		return s0 + "." + s1
	default:
		return s1
	}
}

// Fields returns the fields of this [Message]
func (p *Message) Fields() []*Field {
	if p.fields == nil {
		p.loadFields()
	}
	return p.fields
}

func (p *Message) loadFields() {
	out := make([]*Field, 0, len(p.dp.Field))
	for _, dp := range p.dp.Field {
		if dp == nil {
			continue
		}

		q := &Field{
			msg: p,
			dp:  dp,
		}

		out = append(out, q)
	}

	// sort fields by name
	Sort(out, func(a, b *Field) bool {
		return a.Name() < b.Name()
	})
	p.fields = out
}
