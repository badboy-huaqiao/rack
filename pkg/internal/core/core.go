package core

import (
	"context"
	"plugin"
	"rack/pkg/model"
	"rack/pkg/template"
)

type Executor interface {
	Run()
}

type ExecutionContext struct {
	Plugins map[model.Plugin]*PluginContext
}

var executionContext *ExecutionContext

type PluginContext struct {
	Ctx  context.Context
	Plug *plugin.Plugin
}

func init() {
	if executionContext == nil {
		executionContext = &ExecutionContext{
			Plugins: make(map[model.Plugin]*PluginContext),
		}
	}
}

func loader(path string) (p *plugin.Plugin, err error) {
	if p, err = plugin.Open(path); err != nil {
		return nil, err
	}
	return p, nil
}

func Load(p model.Plugin) error {
	plug, err := loader(p.Path)
	if err != nil {
		return err
	}
	pc := PluginContext{
		Ctx:  context.Background(),
		Plug: plug,
	}
	executionContext.Plugins[p] = &pc
	go executor(&pc)
	return nil
}

func executor(plugContext *PluginContext) error {
	RackContextSym, err := plugContext.Plug.Lookup("RackContext")
	if err != nil {
		return err
	}

	switch RackContextSym.(type) {
	case template.Hulu:
		go executeHulu(RackContextSym)
	}

	return nil
}
