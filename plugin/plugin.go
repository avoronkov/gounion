package main

import (
	"github.com/avoronkov/gounion/checker"
	"golang.org/x/tools/go/analysis"
)

func New(conf any) ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{checker.Analyzer}, nil
}
