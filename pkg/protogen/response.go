package protogen

import (
	"bytes"
	"io"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

// Response returns a reference to the assembled [pluginpb.CodeGeneratorResponse]
func (gen *Plugin) Response() *pluginpb.CodeGeneratorResponse {
	return &gen.resp
}

// MarshalCodeGeneratorResponse writes the proto encoded representation of the
// given [pluginpb.CodeGeneratorResponse]
func MarshalCodeGeneratorResponse(resp *pluginpb.CodeGeneratorResponse,
	w io.Writer) (int64, error) {
	// encode
	b, err := proto.Marshal(resp)
	if err != nil {
		return 0, err
	}

	// write
	buf := bytes.NewBuffer(b)
	return buf.WriteTo(w)
}

// WriteTo writes the generated [pluginpb.CodeGeneratorResponse] to the
// provided [io.Writer]
func (gen *Plugin) WriteTo(w io.Writer) (int64, error) {
	return MarshalCodeGeneratorResponse(&gen.resp, w)
}

func (gen *Plugin) Write() (int64, error) {
	return gen.WriteTo(gen.options.Stdout)
}

// WriteErrorTo generates an error response and writes it to the
// provided [io.Writer]
func (*Plugin) WriteErrorTo(err error, w io.Writer) (int64, error) {
	s := err.Error()

	resp := &pluginpb.CodeGeneratorResponse{
		Error: &s,
	}

	return MarshalCodeGeneratorResponse(resp, w)
}

// WriteError generates an error response and writes it to Stdout
func (gen *Plugin) WriteError(err error) (int64, error) {
	return gen.WriteErrorTo(err, gen.options.Stdout)
}
