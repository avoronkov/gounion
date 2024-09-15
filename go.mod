module github.com/avoronkov/gounion

go 1.23.1

replace honnef.co/go/lint => ./honnef.co/go/lint

require (
	golang.org/x/tools v0.25.0
	honnef.co/go/lint v0.0.0-00010101000000-000000000000
)

require github.com/kisielk/gotool v1.0.0 // indirect
