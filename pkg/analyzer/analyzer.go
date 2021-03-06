package analyzer

import (
	"flag"
	"fmt"
	"go/ast"
	"regexp"
	"strings"

	"golang.org/x/tools/go/analysis"
)

//nolint:gochecknoglobals
var flagSet flag.FlagSet

type arrayFlags []string

func (f *arrayFlags) String() string {
	s := []string{}
	for _, a := range *f {
		s = append(s, a)
	}

	return fmt.Sprintf("%v", s)
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

//nolint:gochecknoglobals
var (
	forbiddenFuncNamesArgs    arrayFlags
	forbiddenFuncNames        []*regexp.Regexp
	forbiddenPackageNames     []*regexp.Regexp
	forbiddenPackageNamesArgs arrayFlags
)

//nolint:gochecknoglobals
var defaultPackageNames = arrayFlags{"^util", "^helper", "^base"}

//nolint:gochecknoinits
func init() {
	flagSet.Var(&forbiddenFuncNamesArgs, "functionNames", "list of regexps that are forbidden to use in function names")
	flagSet.Var(&forbiddenPackageNamesArgs, "packageNames", "list of regexps that are forbidden to use in package names")
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
	for _, pattern := range forbiddenFuncNamesArgs {
		rxp, err := regexp.Compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("cannot parse function pattern: %w", err)
		}

		forbiddenFuncNames = append(forbiddenFuncNames, rxp)
	}

	if len(forbiddenPackageNamesArgs) == 0 {
		forbiddenPackageNamesArgs = defaultPackageNames
	}
	for _, pattern := range forbiddenPackageNamesArgs {
		rxp, err := regexp.Compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("cannot parse package pattern: %w", err)
		}

		forbiddenPackageNames = append(forbiddenPackageNames, rxp)
	}

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

	for _, word := range forbiddenFuncNamesArgs {
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
