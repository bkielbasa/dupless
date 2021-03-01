package analyzer

import (
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestFuncNames(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get wd: %s", err)
	}

	forbiddenFuncNames = []string{"dupa"}

	testdata := filepath.Join(filepath.Dir(filepath.Dir(wd)), "testdata")
	analysistest.Run(t, testdata, NewAnalyzer(), "funcs")
}

func TestPackageName(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get wd: %s", err)
	}

	forbiddenPackageNames = []*regexp.Regexp{regexp.MustCompile("dup")}

	testdata := filepath.Join(filepath.Dir(filepath.Dir(wd)), "testdata")
	analysistest.Run(t, testdata, NewAnalyzer(), "dupaInPkgName")
}
