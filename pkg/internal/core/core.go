package core

import (
	"context"
	"plugin"
	"rack/pkg/model"
	"rack/pkg/template"
)

type ExecutionContext struct {
	Ctx     context.Context
	Plugins map[model.Plugin]*PluginContext
}

var execCtx *ExecutionContext

type PluginContext struct {
	Ctx context.Context
	model.Plugin
	Templates interface{}
	Plug      *plugin.Plugin
	Quit      context.CancelFunc
}

func GetExecutionContext() *ExecutionContext {
	return execCtx
}

func init() {
	if execCtx == nil {
		execCtx = &ExecutionContext{
			Ctx:     context.Background(),
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

func Deregister(p model.Plugin) error {
	for k, v := range execCtx.Plugins {
		if k.Concrete == p.Concrete && k.Version == p.Version {
			v.Quit()
			delete(execCtx.Plugins, k)
			return nil
		}
	}
	return nil
}

func Register(p model.Plugin) (err error) {
	plug, err := loader(p.Path)
	if err != nil {
		return err
	}

	pc := PluginContext{
		Ctx:  context.Background(),
		Plug: plug,
	}
	execCtx.Plugins[p] = &pc
	if err := exec(&pc); err != nil {
		return err
	}
	return nil
}

func exec(plugCtx *PluginContext) error {
	ConcreteSym, err := plugCtx.Plug.Lookup(plugCtx.Concrete)
	if err != nil {
		return err
	}
	switch v := ConcreteSym.(type) {
	case template.Hulu:
		go func() {
			var ctx context.Context
			ctx, plugCtx.Quit = context.WithCancel(plugCtx.Ctx)
			template.NewHuluCtx(ctx).Run(v)
		}()
	}

	return nil
}
