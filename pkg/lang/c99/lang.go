// Package c99gen implements C99 specific details for generators
package c99gen

import (
	"os"
	"path/filepath"

	"github.com/amery/protogen/pkg/protogen"
)

// Plugin is a protogen.Plugin wrapper for C99 generators
type Plugin struct {
	*protogen.Plugin

	Name string
}

func (p *Plugin) setDefaults() {
	if p.Name == "" {
		p.Name = filepath.Base(os.Args[0])
	}
}

func (p *Plugin) init() error {
	p.setDefaults()
	return nil
}

// NewPlugin ...
func NewPlugin(gen *protogen.Plugin) (*Plugin, error) {
	if gen == nil {
		return nil, nil
	}

	p := &Plugin{Plugin: gen}
	if err := p.init(); err != nil {
		return nil, err
	}

	return p, nil
}
