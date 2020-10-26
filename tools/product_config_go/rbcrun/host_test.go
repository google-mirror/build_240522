package rbcrun

import (
	"fmt"
	"go.starlark.net/resolve"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarktest"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// In order to use "assert.star" from go/starlark.net/starlarktest in the tests,
// provide:
//  * load function that handles "assert.star"
//  * starlarktest.DataFile function that finds its location

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

// Common setup for the tests: create thread, change to the test directory
func testSetup(t *testing.T) (*starlark.Thread, string) {
	thread := &starlark.Thread{
		Load: func(thread *starlark.Thread, module string) (starlark.StringDict, error) {
			if module == "assert.star" {
				return starlarktest.LoadAssertModule()
			}
			return nil, fmt.Errorf("load not implemented")
		}}
	starlarktest.SetReporter(thread, t)
	_, thisSrcFile, _, _ := runtime.Caller(0)
	dataDir := filepath.Join(filepath.Dir(filepath.Dir(thisSrcFile)), "testdata")
	if err := os.Chdir(dataDir); err != nil {
		t.Fatal(err)
	}
	return thread, dataDir
}

func TestFileOps(t *testing.T) {
	thread, dataDir := testSetup(t)
	// see files_op.star:
	if err := os.Setenv("TEST_DATA_DIR", dataDir); err != nil {
		t.Fatal(err)
	}
	if _, err := starlark.ExecFile(thread, "file_ops.star", nil, predeclared(nil)); err != nil {
		if err, ok := err.(*starlark.EvalError); ok {
			t.Fatal(err.Backtrace())
		}
		t.Fatal(err)
	}
}

func TestCliAndEnv(t *testing.T) {
	thread, _ := testSetup(t)
	if err := os.Setenv("TEST_ENVIRONMENT_FOO", "test_environment_foo"); err != nil {
		t.Fatal(err)
	}
	if _, err := starlark.ExecFile(thread, "cli_and_env.star", nil, predeclared([]string{"CLI_FOO=foo"})); err != nil {
		if err, ok := err.(*starlark.EvalError); ok {
			t.Fatal(err.Backtrace())
		}
		t.Fatal(err)
	}
}

func TestRegex(t *testing.T) {
	thread, _ := testSetup(t)
	if _, err := starlark.ExecFile(thread, "regex.star", nil, predeclared(nil)); err != nil {
		if err, ok := err.(*starlark.EvalError); ok {
			t.Fatal(err.Backtrace())
		}
		t.Fatal(err)
	}
}

func TestLoad(t *testing.T) {
	thread, dataDir := testSetup(t)
	thread.Load = func(thread *starlark.Thread, module string) (starlark.StringDict, error) {
		if module == "assert.star" {
			return starlarktest.LoadAssertModule()
		} else {
			return loader(thread, module)
		}
	}
	thread.SetLocal(callerDirKey, dataDir)
	LoadPathRoot = filepath.Dir(dataDir)
	if _, err := starlark.ExecFile(thread, "load.star", nil, predeclared(nil)); err != nil {
		if err, ok := err.(*starlark.EvalError); ok {
			t.Fatal(err.Backtrace())
		}
		t.Fatal(err)
	}

}
