package analyzer

import (
	"flag"
	"go/ast"
	"regexp"
	"strings"

	"golang.org/x/tools/go/analysis"
)

//nolint:gochecknoglobals
var flagSet flag.FlagSet

//nolint:gochecknoglobals
var (
	maxComplexity         int
	forbiddenFuncNames    []string
	forbiddenPackageNames []*regexp.Regexp
)

//nolint:gochecknoinits
func init() {
	// flagSet.IntVar(&maxComplexity, "maxComplexity", 10, "max complexity the function can have")
}

func NewAnalyzer() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name:  "dupless",
		Doc:   "checks whatever function, variables or packages contain forbidden words",
		Run:   run,
		Flags: flagSet,
	}
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			if f, ok := node.(*ast.FuncDecl); ok {
				checkFunctionNames(pass, f)
			}
			if pkg, ok := node.(*ast.File); ok {
				checkPkgNames(pass, pkg)
			}

			return true
		})
	}

	return nil, nil
}

func checkFunctionNames(pass *analysis.Pass, f *ast.FuncDecl) {
	funcName := strings.ToLower(f.Name.Name)

	for _, word := range forbiddenFuncNames {
		if strings.Contains(funcName, word) {
			pass.Reportf(f.Pos(), "the function name contains a forbidden word: %s", word)
		}
	}
}

func checkPkgNames(pass *analysis.Pass, f *ast.File) {
	name := strings.ToLower(f.Name.Name)

	for _, reg := range forbiddenPackageNames {
		if reg.MatchString(name) {
			pass.Reportf(f.Pos(), "the package name matches the pattern: %s", reg)
		}
	}
}
