package main

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/staticcheck"
	"strings"
)

func main() {
	var analyzers []*analysis.Analyzer

	// добавляем стандартные анализаторы
	analyzers = append(analyzers, printf.Analyzer, shadow.Analyzer, structtag.Analyzer)

	for _, v := range staticcheck.Analyzers {
		// добавляем SA-анализаторы из staticcheck
		if strings.HasPrefix(v.Analyzer.Name, "SA") {
			analyzers = append(analyzers, v.Analyzer)
		}

		// анализаторы из остальных классов staticcheck
		if v.Analyzer.Name == "S1000" || v.Analyzer.Name == "ST1000" {
			analyzers = append(analyzers, v.Analyzer)
		}
	}

	// добавляем анализатор для поиска os.Exit
	analyzers = append(analyzers, NoOsExitInMainAnalyzer)

	multichecker.Main(analyzers...)
}

// NoOsExitInMainAnalyzer is an analyzer that checks if os.Exit is called in the main function.
var NoOsExitInMainAnalyzer = &analysis.Analyzer{
	Name: "exitguard",
	Doc:  "checks if os.Exit is called in main function",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	if pass.Pkg.Name() != "main" {
		return nil, nil
	}

	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			callExpr, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}

			ident, ok := selExpr.X.(*ast.Ident)
			if !ok {
				return true
			}

			if ident.Name == "os" && selExpr.Sel.Name == "Exit" {
				pass.Reportf(callExpr.Pos(), "direct call of os.Exit in main function is not allowed")
			}

			return true
		})
	}
	return nil, nil
}
