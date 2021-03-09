package analyzer

import (
	"flag"
	"fmt"
	"go/ast"
	"go/token"
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

func (f *arrayFlags) Set(value string) error {
	*f = append(*f, value)
	return nil
}

//nolint:gochecknoglobals
var (
	forbiddenFuncNamesArgs     arrayFlags
	forbiddenFuncNames         []*regexp.Regexp
	forbiddenPackageNamesArgs  arrayFlags
	forbiddenPackageNames      []*regexp.Regexp
	forbiddenVariableNames     []*regexp.Regexp
	forbiddenVariableNamesArgs arrayFlags
)

//nolint:gochecknoglobals
var defaultPackageNames = arrayFlags{"^util[s]$", "^helper[s]$", "^base$", "^interfaces$"}

//nolint:gochecknoinits
func init() {
	flagSet.Var(&forbiddenFuncNamesArgs, "functionNames", "list of regexps that are forbidden to use in function names")
	flagSet.Var(&forbiddenPackageNamesArgs, "packageNames", "list of regexps that are forbidden to use in package names")
	flagSet.Var(&forbiddenFuncNamesArgs, "variableNames", "list of regexps that are forbidden to use in variable names")
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

	for _, pattern := range forbiddenVariableNamesArgs {
		rxp, err := regexp.Compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("cannot parse variable pattern: %w", err)
		}

		forbiddenVariableNames = append(forbiddenVariableNames, rxp)
	}

	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			if f, ok := node.(*ast.FuncDecl); ok {
				checkFunctionNames(pass, f)
			}

			if pkg, ok := node.(*ast.File); ok {
				checkPkgNames(pass, pkg)
			}

			if ident, ok := node.(*ast.AssignStmt); ok {
				checkVarNames(pass, ident.Lhs)
			}
			if ident, ok := node.(*ast.GenDecl); ok {
				if ident.Tok != token.VAR {
					return true
				}
				checkVarNamesInValueSpec(pass, ident.Specs)
			}

			return true
		})
	}

	return nil, nil
}

func checkVarNamesInValueSpec(pass *analysis.Pass, specs []ast.Spec) {
	for _, spec := range specs {
		var value *ast.ValueSpec
		var ok bool
		if value, ok = spec.(*ast.ValueSpec); !ok {
			return
		}

		for _, ident := range value.Names {
			varName := ident.Name

			for _, reg := range forbiddenVariableNames {
				if reg.MatchString(varName) {
					pass.Reportf(spec.Pos(), "the variable name contains the forbidden pattern: %s", reg)
				}
			}
		}

	}
}

func checkVarNames(pass *analysis.Pass, expressions []ast.Expr) {
	for _, expr := range expressions {
		if ident, ok := expr.(*ast.Ident); ok {
			varName := strings.ToLower(ident.Name)

			for _, reg := range forbiddenVariableNames {
				if reg.MatchString(varName) {
					pass.Reportf(ident.Pos(), "the variable name contains the forbidden pattern: %s", reg)
				}
			}
		}
	}
}
func checkFunctionNames(pass *analysis.Pass, f *ast.FuncDecl) {
	funcName := strings.ToLower(f.Name.Name)

	for _, reg := range forbiddenFuncNames {
		if reg.MatchString(funcName) {
			pass.Reportf(f.Pos(), "the function name contains the forbidden pattern: %s", reg)
		}
	}
}

func checkPkgNames(pass *analysis.Pass, f *ast.File) {
	name := strings.ToLower(f.Name.Name)

	for _, reg := range forbiddenPackageNames {
		if reg.MatchString(name) {
			pass.Reportf(f.Pos(), "the package name matches the forbidden pattern: %s", reg)
		}
	}
}
