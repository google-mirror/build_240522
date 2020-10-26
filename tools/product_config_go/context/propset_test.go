package context

import (
	"fmt"
	"go.starlark.net/resolve"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
	"go.starlark.net/starlarktest"
	"path/filepath"
	"runtime"
	"testing"
)

func init() {
	// The tests make extensive use of these not-yet-standard features.
	resolve.AllowLambda = true
	starlarktest.DataFile = func(pkgdir, filename string) string {
		// TODO(handle Bazel case
		// We know that the caller of this methods is in starlarktest.go. Ughh.
		_, starlarktestSrcFile, _, _ := runtime.Caller(1)
		if filepath.Base(starlarktestSrcFile) != "starlarktest.go" {
			panic(fmt.Errorf("this function should be called from starlarktest.go, got %s",
				starlarktestSrcFile))
		}
		return filepath.Join(filepath.Dir(starlarktestSrcFile), filename)
	}
}

func Test(t *testing.T) {
	thread := &starlark.Thread{Load: load}
	starlarktest.SetReporter(thread, t)
	_, thisSrcFile, _, _ := runtime.Caller(0)
	filename := filepath.Join(filepath.Dir(filepath.Dir(thisSrcFile)), "data/propset.star")
	predeclared := starlark.StringDict{
		"propset": starlark.NewBuiltin("struct", MakePropset),
		"module":  starlark.NewBuiltin("module", starlarkstruct.MakeModule),
	}
	if _, err := starlark.ExecFile(thread, filename, nil, predeclared); err != nil {
		if err, ok := err.(*starlark.EvalError); ok {
			t.Fatal(err.Backtrace())
		}
		t.Fatal(err)
	}
}

// load implements the 'load' operation as used in the evaluator tests.
func load(_ *starlark.Thread, module string) (starlark.StringDict, error) {
	if module == "assert.star" {
		return starlarktest.LoadAssertModule()
	}
	return nil, fmt.Errorf("load not implemented")
}
