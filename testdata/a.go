package testdata

type A interface {
	a()
}

type One struct{}

func (*One) a() {}

type Two struct{}

func (Two) a() {}

func test() {
	var a A = &One{}

	switch a.(type) { // want `uncovered cases for testdata.A type switch: \*testdata.One, testdata.Two`
	case *Two:
	}
}
