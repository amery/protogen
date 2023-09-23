package protogen

import (
	"io"
	"strings"

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
	// []ProtoFile
	gen.loadFiles(req.ProtoFile...)
	if err := gen.setFilesGenerate(req.FileToGenerate...); err != nil {
		return err
	}

	// Parameter
	if p := req.Parameter; p != nil {
		err := gen.loadParams(*p)
		if err != nil {
			// bad parameters
			return err
		}
	}

	gen.req = req
	return nil
}

// Param returns the value of a parameter if specified
func (gen *Plugin) Param(key string) (string, bool) {
	value, found := gen.params[key]
	return value, found
}

// Params returns all specified parameters
func (gen *Plugin) Params() map[string]string {
	return gen.params
}

func (gen *Plugin) loadParams(params string) error {
	for _, s := range strings.Split(params, ",") {
		s = strings.TrimSpace(s)
		if s == "" {
			// skip empty
			continue
		}

		k, v, found := strings.Cut(s, "=")
		if !found {
			v = "true"
		}

		err := gen.setParam(k, v)
		if err != nil {
			return err
		}
	}

	return nil
}

func (gen *Plugin) setParam(k, v string) error {
	var err error

	gen.params[k] = v

	if gen.options.ParamFunc != nil {
		err = gen.options.ParamFunc(k, v)
	}

	return err
}
