package main

import (
	"github.com/avoronkov/gounion"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(gounion.Analyzer)
}
