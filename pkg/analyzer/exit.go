package analyzer

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// ExitAnalyzer inspects the main func of package main and reports any os.Exit calls.
var ExitAnalyzer = &analysis.Analyzer{
	Name: "exit",
	Doc:  "prohibits os.Exit calls in main function of package main",
	Run:  runExitAnalyzer,
}

func runExitAnalyzer(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		if file.Name.Name != "main" {
			continue
		}

		for _, decl := range file.Decls {
			fnDecl, ok := decl.(*ast.FuncDecl)
			if !ok || fnDecl.Name.Name != "main" {
				continue
			}

			ast.Inspect(fnDecl.Body, func(n ast.Node) bool {
				if sel, ok := getSelector(n); ok && sel == "Exit" {
					pass.Reportf(n.Pos(), "os.Exit calls are prohibited in main function of package main")
					return false
				}
				return true
			})
		}
	}
	return nil, nil
}

func getSelector(n ast.Node) (string, bool) {
	c, ok := n.(*ast.CallExpr)
	if !ok {
		return "", false
	}

	s, ok := c.Fun.(*ast.SelectorExpr)
	if !ok {
		return "", false
	}

	return s.Sel.Name, true
}
