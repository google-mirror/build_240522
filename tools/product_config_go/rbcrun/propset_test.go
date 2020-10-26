package rbcrun

import (
	"fmt"
	"go.starlark.net/resolve"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarktest"
	"path/filepath"
	"runtime"
	"testing"
)

func init() {
	resolve.AllowLambda = true
	starlarktest.DataFile = func(pkgdir, filename string) string {
		// The caller expects this function to return the path to the
		// data file. The implementation assumes that the source file
		// containing the caller and the data file are in the same
		// directory. It's ugly. Not sure what's the better way.
		// TODO(asmundak): handle Bazel case
		_, starlarktestSrcFile, _, _ := runtime.Caller(1)
		if filepath.Base(starlarktestSrcFile) != "starlarktest.go" {
			panic(fmt.Errorf("this function should be called from starlarktest.go, got %s",
				starlarktestSrcFile))
		}
		return filepath.Join(filepath.Dir(starlarktestSrcFile), filename)
	}
}

func Test(t *testing.T) {
	// In order to use "assert.star" from go/starlark.net/starlarktest in the tests, provide:
	//  * load function that handles "assert.star"
	//  * starlarktest.DataFile function that finds its location
	thread := &starlark.Thread{
		Load: func(thread *starlark.Thread, module string) (starlark.StringDict, error) {
			if module == "assert.star" {
				return starlarktest.LoadAssertModule()
			}
			return nil, fmt.Errorf("load not implemented")
		}}
	starlarktest.SetReporter(thread, t)
	_, thisSrcFile, _, _ := runtime.Caller(0)
	filename := filepath.Join(filepath.Dir(filepath.Dir(thisSrcFile)), "testdata/propset.star")
	if _, err := starlark.ExecFile(thread, filename, nil, predeclared(nil)); err != nil {
		if err, ok := err.(*starlark.EvalError); ok {
			t.Fatal(err.Backtrace())
		}
		t.Fatal(err)
	}
}
