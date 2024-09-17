package checker

import (
	"fmt"
	"go/ast"
	"go/types"
	"log"
	"strings"

	"github.com/avoronkov/gounion"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/types/typeutil"
)

var Analyzer = &analysis.Analyzer{
	Name: "gounioncheck",
	Doc:  "Static Check Tool of type switch-case handling",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	all := typeutil.Dependencies(pass.Pkg)
	info := pass.TypesInfo
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			typeswitch, ok := node.(*ast.TypeSwitchStmt)
			if !ok {
				return true
			}
			named := switchInterfaceType(typeswitch, *info)
			if named == nil {
				return true
			}
			iface := gounion.NewSumInterface(named, all)
			if iface == nil {
				return true
			}
			covered := make(map[types.Type]bool, len(iface.Implements.Pointers))
			for _, ptr := range iface.Implements.Pointers {
				covered[ptr.Elem()] = false
			}
			coveredTypes := make(map[types.Type]bool, len(iface.Implements.Types))
			for _, typ := range iface.Implements.Types {
				coveredTypes[*typ] = false
			}

			coveredTypsByInterface := make(map[*types.Named][]types.Type)
			for _, i := range iface.Implements.Interfaces {
				name := i.NamedInterface
				coveredTypsByInterface[name] = make([]types.Type, len(i.Implements.Pointers))
				for j, ptr := range i.Implements.Pointers {
					coveredTypsByInterface[name][j] = ptr.Elem()
				}
			}
			for _, caseClause := range typeswitch.Body.List {
				c, ok := caseClause.(*ast.CaseClause)
				if !ok {
					log.Printf("got unexpected node: %v", caseClause)
					continue
				}
				if c.List == nil {
					// TODO: handle default clause
				} else {
					for _, expr := range c.List {
						tv, ok := info.Types[expr]
						if !ok {
							// Just ignore this case and continue.
							//   log.Printf("fail to got type: %v", expr)
							// You can see sample cases when you run gosumcheck to docker.
							// $ gosumcheck github.com/docker/docker/...
							//   2016/12/04 23:57:42 fail to got type: &{client ErrRepoNotInitialized}
							//   2016/12/04 23:57:42 fail to got type: &{client ErrRepositoryNotExist}
							//   2016/12/04 23:57:42 fail to got type: &{signed ErrExpired}
							// ...
							continue
						}
						typ := tv.Type
						if types.IsInterface(typ) {
							for _, typ := range coveredTypsByInterface[typ.(*types.Named)] {
								covered[typ] = true
							}
						} else {
							if ptr, ok := typ.(*types.Pointer); ok {
								covered[ptr.Elem()] = true
							} else {
								coveredTypes[typ] = true
							}
						}
					}
				}
			}

			var uncovered []string
			for elem, b := range covered {
				if !b {
					uncovered = append(uncovered, types.NewPointer(elem).String())
				}
			}
			for elem, b := range coveredTypes {
				if !b {
					uncovered = append(uncovered, elem.String())
				}
			}

			if len(uncovered) > 0 {
				pass.Report(analysis.Diagnostic{
					Pos:            typeswitch.Pos(),
					End:            0,
					Category:       "",
					Message:        fmt.Sprintf("uncovered cases for %v type switch\n\t- %v", named.String(), strings.Join(uncovered, "\n\t- ")),
					SuggestedFixes: nil,
				})
			}

			return true
		})
	}
	return nil, nil
}

// switchInterfaceType returns interface type of type switch statement. It may
// return nil.
func switchInterfaceType(node *ast.TypeSwitchStmt, info types.Info) *types.Named {
	ae := assertExpr(node)
	if ae == nil {
		return nil
	}
	tv, ok := info.Types[ae.X]
	if !ok {
		return nil
	}

	if named, ok := tv.Type.(*types.Named); ok && types.IsInterface(named) {
		return named
	}
	return nil
}

func assertExpr(x *ast.TypeSwitchStmt) *ast.TypeAssertExpr {
	switch a := x.Assign.(type) {
	case *ast.AssignStmt: // x := y.(type)
		for _, expr := range a.Rhs {
			ae, ok := expr.(*ast.TypeAssertExpr)
			if !ok {
				continue
			}
			return ae
		}
		return nil
	case *ast.ExprStmt: // y.(type)
		ae, ok := a.X.(*ast.TypeAssertExpr)
		if !ok {
			return nil
		}
		return ae
	}
	return nil
}
