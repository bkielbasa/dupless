package analyzer

import (
	"os"
	"path/filepath"
	"testing"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/analysistest"
)

func resetGlobals() {
	forbiddenFuncNamesArgs = []string{}
	forbiddenPackageNamesArgs = []string{}
	forbiddenVariableNamesArgs = []string{}
}

func TestFuncNames(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get wd: %s", err)
	}

	resetGlobals()
	forbiddenFuncNamesArgs = []string{"dupa"}

	testdata := filepath.Join(filepath.Dir(filepath.Dir(wd)), "testdata")
	analysistest.Run(t, testdata, NewAnalyzer(), "funcs")
}

func TestPackageName(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get wd: %s", err)
	}

	resetGlobals()
	forbiddenPackageNamesArgs = []string{"dup"}

	testdata := filepath.Join(filepath.Dir(filepath.Dir(wd)), "testdata")
	analysistest.Run(t, testdata, NewAnalyzer(), "dupaInPkgName")
}

func TestInvalidRegexp(t *testing.T) {
	forbiddenPackageNamesArgs = []string{"["}
	anal := NewAnalyzer()
	_, err := anal.Run(&analysis.Pass{})
	if err == nil {
		t.Fatalf("expected to get an error but got nil")
	}
	resetGlobals()
	forbiddenFuncNamesArgs = []string{"["}
	_, err = anal.Run(&analysis.Pass{})
	if err == nil {
		t.Fatalf("expected to get an error but got nil")
	}
}

func TestVariableName(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get wd: %s", err)
	}

	resetGlobals()
	forbiddenVariableNamesArgs = []string{"xyz", "xxx"}

	testdata := filepath.Join(filepath.Dir(filepath.Dir(wd)), "testdata")
	analysistest.Run(t, testdata, NewAnalyzer(), "varnames")
}
