// Package gogen implements Go specific details for generators
package gogen

import (
	"github.com/amery/protogen/pkg/protogen"
	"github.com/amery/protogen/pkg/text"
)

// Plugin is a protogen.Plugin wrapper for Go generators
type Plugin struct {
	*protogen.Plugin

	pkgMap map[string]*Package
}

func (*Plugin) init() error {
	return nil
}

// SetParam ...
func (gen *Plugin) SetParam(key, value string) error {
	if fn, ok := text.CutPrefix(key, "M"); ok {
		return gen.SetGoPackage(fn, value)
	}

	return protogen.ErrUnknownParam
}

// NewPlugin ...
func NewPlugin(gen *protogen.Plugin) (*Plugin, error) {
	if gen == nil {
		return nil, nil
	}

	p := &Plugin{
		Plugin: gen,

		pkgMap: make(map[string]*Package),
	}

	if err := p.init(); err != nil {
		return nil, err
	}

	return p, nil
}
