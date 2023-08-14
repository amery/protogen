package protogen

import (
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

var (
	_ Message = (*MessageDescriptor)(nil)
)

// MessageDescriptor implements Message
type MessageDescriptor struct {
	file *FileDescriptor
	dp   *descriptorpb.DescriptorProto
}

// Request returns the [pluginpb.CodeGeneratorRequest] received by
// the [Plugin]
func (p *MessageDescriptor) Request() *pluginpb.CodeGeneratorRequest {
	return p.file.Request()
}

// Proto returns the underlying protobuf structure
func (p *MessageDescriptor) Proto() *descriptorpb.DescriptorProto {
	return p.dp
}

// File returns the [File] that defines this type
func (p *MessageDescriptor) File() File {
	return p.file
}

// Package returns the package name associated to this type
func (p *MessageDescriptor) Package() string {
	return p.file.Package()
}

// Name returns the relative name of this type
func (p *MessageDescriptor) Name() string {
	return optional(p.dp.Name, "")
}

// FullName returns the fully qualified name of this type
func (p *MessageDescriptor) FullName() string {
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
func (f *FileDescriptor) Messages() []Message {
	return f.messages
}

// MessageByName find a [Message] by name
func (f *FileDescriptor) MessageByName(name string) Message {
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

func (f *FileDescriptor) loadMessages() {
	out := make([]Message, 0, len(f.dp.MessageType))
	for _, dp := range f.dp.MessageType {
		if dp == nil {
			continue
		}

		p := &MessageDescriptor{
			file: f,
			dp:   dp,
		}

		out = append(out, p)
	}
	f.messages = out
}
