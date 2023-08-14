package protogen

import (
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

var (
	_ ProtoTyper = (*Message)(nil)
)

// Message represents a type
type Message struct {
	file *File
	dp   *descriptorpb.DescriptorProto
}

// Request returns the [pluginpb.CodeGeneratorRequest] received by
// the [Plugin]
func (p *Message) Request() *pluginpb.CodeGeneratorRequest {
	return p.file.Request()
}

// Proto returns the underlying protobuf structure
func (p *Message) Proto() *descriptorpb.DescriptorProto {
	return p.dp
}

// File returns the [File] that defines this type
func (p *Message) File() *File {
	return p.file
}

// Package returns the package name associated to this type
func (p *Message) Package() string {
	return p.file.Package()
}

// Name returns the relative name of this type
func (p *Message) Name() string {
	return optional(p.dp.Name, "")
}

// FullName returns the fully qualified name of this type
func (p *Message) FullName() string {
	s0 := p.file.Package()
	s1 := p.Name()

	switch {
	case s0 != "" && s1 != "":
		return s0 + "." + s1
	default:
		return s1
	}
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
	f.messages = out
}
