package main

import (
	"github.com/avoronkov/gounion/checker"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(checker.Analyzer)
}
