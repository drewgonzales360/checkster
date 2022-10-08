package funclen

import (
	"fmt"
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	statementLimit = 25
)

// Analyzer runs static analysis.
var Analyzer = &analysis.Analyzer{
	Name:     "funclen",
	Doc:      "Checks the length of functions",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return nil, fmt.Errorf("analyzer is not type *inspector.Inspector")
	}

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
		(*ast.File)(nil),
	}

	pkg := "main"
	inspect.Preorder(nodeFilter, func(node ast.Node) {
		switch stmt := node.(type) {
		case *ast.File:
			pkg = stmt.Name.Name

		case *ast.FuncDecl:
			statements := stmt.Body.List
			funcStatements := len(statements)
			recvType := getRecvType(pass, stmt.Recv)

			if funcStatements > statementLimit {
				pass.Reportf(
					node.Pos(),
					"%s.%s%s has more than %d statments(%d)",
					pkg,
					recvType,
					stmt.Name.String(),
					statementLimit,
					funcStatements,
				)
			}
		}
	})

	return nil, nil
}

func getRecvType(pass *analysis.Pass, recv *ast.FieldList) string {
	if recv == nil {
		return ""
	}
	if recv.List == nil {
		return ""
	}
	if len(recv.List) != 1 {
		return ""
	}

	fullType := strings.Split(pass.TypesInfo.Types[recv.List[0].Type].Type.String(), "/")
	pkgType := fullType[len(fullType)-1]
	return fmt.Sprintf("(%s)", pkgType)
}
