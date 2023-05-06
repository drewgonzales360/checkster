// docwrap checks that all the godoc comments in a file are
// the correct length.
package docwrap

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	lineLength = 100
	errMsg = "the length of the comment is too long for one line"
)

var Analyzer = &analysis.Analyzer{
	Name:     "docwrap",
	Doc:      "Checks the line length of Godoc comments.",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return nil, fmt.Errorf("analyzer is not type *inspector.Inspector")
	}

	nodeFilter := []ast.Node{
		(*ast.Comment)(nil),
		(*ast.CommentGroup)(nil),
	}

	inspect.Preorder(nodeFilter, func(node ast.Node) {
		switch stmt := node.(type) {
		case *ast.Comment:
			if len(stmt.Text) > lineLength{
				pass.Reportf(stmt.Pos(), errMsg)
			}
		case *ast.CommentGroup:
			lines := stmt.List
			for _, line := range lines {
				if len(line.Text) > lineLength {
					pass.Reportf(line.Pos(), errMsg)
				}
			}
		}
	})

	return nil, nil
}
