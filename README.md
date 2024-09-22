## Union Type in Go and Static Check Tool of switch-case handling  

## gounioncheck

### Installation 

```
$ go install github.com/avoronkov/gounion/cmd/gounioncheck@latest
```

### Example

```
// main.go
package main

type Foo interface {
	foo()
}

type Bar struct{}

func (b Bar) foo() {}

type Baz struct{}

func (b *Baz) foo() {}

func main() {
	var foo Foo = &Bar{}
	switch foo.(type) {
	case *Bar:
		// Bar
	}
}
```

```console
$ gounioncheck .
main.go:17:2: uncovered cases for main.Foo type switch: *main.Baz, main.Bar
```

### Additional information

It cheks only so-called "private" interfaces i.e. interfaces with private methods which can be implemented only whithin the package they are declared.
