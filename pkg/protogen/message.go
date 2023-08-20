package protogen

import (
	"sort"

	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

var (
	_ ProtoTyper = (*Message)(nil)
)

// Message represents a type
type Message struct {
	file *File
	msg  *Message
	dp   *descriptorpb.DescriptorProto

	enums    []*Enum
	messages []*Message
	fields   []*Field
}

// Request returns the [pluginpb.CodeGeneratorRequest] received by
// the [Plugin]
func (p *Message) Request() *pluginpb.CodeGeneratorRequest {
	return p.File().Request()
}

// Proto returns the underlying protobuf structure
func (p *Message) Proto() *descriptorpb.DescriptorProto {
	return p.dp
}

// File returns the [File] that defines this type
func (p *Message) File() *File {
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
func (p *Message) Package() string {
	return p.File().Package()
}

// Name returns the relative name of this type
func (p *Message) Name() string {
	return optional(p.dp.Name, "")
}

// FullName returns the fully qualified name of this type
func (p *Message) FullName() string {
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

// Messages returns all the [Message] subtypes defined on this message
func (p *Message) Messages() []*Message {
	if p.messages == nil {
		p.loadMessages()
	}
	return p.messages
}

func (p *Message) loadMessages() {
	out := make([]*Message, 0, len(p.dp.NestedType))
	for _, dp := range p.dp.NestedType {
		if dp == nil {
			continue
		}

		q := &Message{
			msg: p,
			dp:  dp,
		}

		out = append(out, q)
	}

	// sort enums by name
	sort.SliceStable(out, func(i, j int) bool {
		a := out[i].Name()
		b := out[j].Name()
		return a < b
	})

	p.messages = out
}

// Messages returns all the [Message] types defined on this file
func (f *File) Messages() []*Message {
	if f.messages == nil {
		f.loadMessages()
	}
	return f.messages
}

// MessageByName find a [Message] by name
func (f *File) MessageByName(name string) *Message {
	pkgname, name, _ := SplitName(name)

	switch {
	case pkgname != "" && f.Package() != pkgname:
		// wrong package
		return nil
	case name == "":
		// no name
		return nil
	default:
		for _, p := range f.Messages() {
			if name == p.Name() {
				// match
				return p
			}
		}
		return nil
	}
}

func (f *File) loadMessages() {
	out := make([]*Message, 0, len(f.dp.MessageType))
	for _, dp := range f.dp.MessageType {
		if dp == nil {
			continue
		}

		p := &Message{
			file: f,
			dp:   dp,
		}

		out = append(out, p)
	}

	// sort enums by name
	sort.SliceStable(out, func(i, j int) bool {
		a := out[i].Name()
		b := out[j].Name()
		return a < b
	})

	f.messages = out
}
