package linter

import (
	"github.com/avoronkov/gounion/checker"
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
)

func init() {
	register.Plugin("gounioncheck", New)
}

func New(conf any) (register.LinterPlugin, error) {
	return &PluginUnionChech{}, nil
}

type PluginUnionChech struct{}

func (p *PluginUnionChech) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{
		checker.Analyzer,
	}, nil
}

func (p *PluginUnionChech) GetLoadMode() string {
	return register.LoadModeSyntax
}
