package main

import (
	"os"

	"github.com/avoronkov/gounion/checker"
	"honnef.co/go/lint/lintutil"
)

func main() {
	lintutil.ProcessArgs("gosumcheck", checker.NewChecker(), os.Args[1:])
}
