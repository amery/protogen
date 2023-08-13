package protogen

import (
	"io"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

// Request returns the received [pluginpb.CodeGeneratorRequest]
func (gen *Plugin) Request() *pluginpb.CodeGeneratorRequest {
	return gen.req
}

// UnmarshalCodeGeneratorRequest reads the proto encoded representation of the
// [pluginpb.CodeGeneratorRequest] from a [io.Reader]
func UnmarshalCodeGeneratorRequest(r io.Reader) (*pluginpb.CodeGeneratorRequest, error) {
	in, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	req := &pluginpb.CodeGeneratorRequest{}
	if err := proto.Unmarshal(in, req); err != nil {
		return nil, err
	}

	return req, nil
}

func (gen *Plugin) loadRequest(req *pluginpb.CodeGeneratorRequest) error {
	// TODO: populate Plugin in a useful way
	gen.req = req
	return nil
}
