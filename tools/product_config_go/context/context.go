package context

import (
	"fmt"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
	"io"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

type configValue struct {
	value interface{}
	frozen bool
	list bool
}

const callerDirKey = "callerDir"
var configVariables map[string]configValue

var LoadPathRoot = "."

type modentry struct {
	globals starlark.StringDict
	err error
}

var moduleCache = make(map[string]*modentry)

var builtins starlark.StringDict

func moduleName2AbsPath(moduleName string, callerDir string) (string, error) {
	if strings.HasPrefix(moduleName, "//") {
		return filepath.Abs(filepath.Join(LoadPathRoot, moduleName[2:]))
	} else if strings.HasPrefix(moduleName, ":") {
		return filepath.Abs(filepath.Join(callerDir, moduleName[1:]))
	} else {
		return filepath.Abs(moduleName)
	}
}
// loader can be passed to starlark.ExecFile to read starlark file and execute it.
// Executed file can use the builtins.
func loader(thread *starlark.Thread, module string) (starlark.StringDict, error) {
	// TODO(asmundak): use thread.Local to maintain caller module's directory,
	// so that load(":foo.star", ...) will look for "foo" in the caller's directory
	modulePath, err := moduleName2AbsPath(module, thread.Local(callerDirKey).(string))
	if err != nil {
		return nil, err
	}
	e, ok := moduleCache[modulePath]
	if e == nil {
		if ok {
			return nil, fmt.Errorf("cycle in load graph")
		}

		// Add a placeholder to indicate "load in progress".
		moduleCache[modulePath] = nil

		// Load it.

		thread := &starlark.Thread{Name: "exec " + module, Load: thread.Load}
		thread.SetLocal(callerDirKey, filepath.Dir(modulePath))
		globals, err := starlark.ExecFile(thread, modulePath, nil, builtins)
		e = &modentry{globals, err}

		// Update the cache.
		moduleCache[modulePath] = e
	}
	return e.globals, e.err
}

func unpackAndSet(b *starlark.Builtin, args starlark.Tuple,
	kwargs []starlark.Tuple, freeze bool) error {
	var cvName, cvValue string
	if err := starlark.UnpackPositionalArgs(b.Name(), args, kwargs, 2, &cvName, &cvValue); err != nil {
		return err
	}
	// TODO(asmundak): make it thread-safe
	if v, ok := configVariables[cvName]; ok {
		if v.list {
			return fmt.Errorf("%s is a list variable", cvName)
		}
		if v.frozen {
			return fmt.Errorf("%s has been already set and cannot be changed", cvName)
		}
	}
	configVariables[cvName]=configValue{value: cvValue, frozen: freeze}
	return nil
}

// set("X", "value")` creates variable `X` and sets its value to `"value"`. It will fail if
// `X` already exists and has been frozen (see `setFinal`)
func set(_ *starlark.Thread, b *starlark.Builtin, args starlark.Tuple,
		kwargs []starlark.Tuple) (starlark.Value, error) {
	return starlark.None, unpackAndSet(b, args, kwargs, false)
}

// setFinal("X", "value") works as `set`, only that `X` value cannot be changed after the call.
func setFinal(_ *starlark.Thread, b *starlark.Builtin, args starlark.Tuple,
		kwargs []starlark.Tuple) (starlark.Value, error) {
	return starlark.None, unpackAndSet(b, args, kwargs, true)
}

// appendTo("X", "value")` creates a list variable if it does not exist and appends `"value"` to
// its value list.
func appendTo(_ *starlark.Thread, b *starlark.Builtin, args starlark.Tuple,
		kwargs []starlark.Tuple) (starlark.Value, error) {
	var cvName, cvValue string
	if err := starlark.UnpackPositionalArgs(b.Name(), args, kwargs, 2, &cvName, &cvValue); err != nil {
		return starlark.None, err
	}

	var sl []string
	// TODO(asmundak): make it thread-safe
	if v, ok := configVariables[cvName]; ok {
		if v.list {
			sl = append(v.value.([]string), cvValue)
			configVariables[cvName] = configValue{value: sl, list: true}
		} else {
			return starlark.None, fmt.Errorf("%s is not a list variable", cvName)
		}
	} else {
		sl = []string{cvValue}
	}
	configVariables[cvName] = configValue{value: sl, list:true}
	return starlark.None, nil
}

// loadGenerated("cmd", [args]) runs command which generates Starlark script on stdout
// and then executes this script.
func loadGenerated(_ *starlark.Thread, b *starlark.Builtin, args starlark.Tuple,
		kwargs []starlark.Tuple) (starlark.Value, error) {
	var command string
	var stringList *starlark.List
	if err := starlark.UnpackArgs(b.Name(), args, kwargs, "c", &command, "a", &stringList); err != nil {
		return starlark.None, err
	}

	cmdArgs := make([]string, stringList.Len())
	for i := 0; i < stringList.Len(); i++ {
		v := stringList.Index(i)
		if strVal, ok := v.(starlark.String); ok {
			cmdArgs[i] = string(strVal)
		} else {
			return starlark.None, fmt.Errorf("command argument list can contain only strings, got %s", v.Type())
		}
	}
	cmd := exec.Command(command,  cmdArgs...)
	bytes, err := cmd.Output()
	if err != nil {
		return starlark.None, fmt.Errorf("'%s' failed: %s", cmd , err)
	}
	loadGenThread := &starlark.Thread{
		Name: "loadGenerated " + command,
		Load: loader,
	}
	_, err = starlark.ExecFile(loadGenThread, fmt.Sprintf("$(%s)",cmd),
			string(bytes), builtins)
	if err != nil {
		return starlark.None, err
	}
	return starlark.None, nil
}

// PrintConfig writes out makefile-style variable assignments sorted by variable name
func PrintConfig(w io.Writer) {
	keys := make([]string, len(configVariables))
	i := 0
	for k := range configVariables {
		keys[i] = k
		i++
	}
	sort.Strings(keys)

	for _, k:= range keys {
		v := configVariables[k]
		if v.list {
			sl := v.value.([]string)
			sorted := make([]string, len(sl))
			copy(sorted, sl)
			sort.Strings(sorted)
			fmt.Fprintf(w, "%s:=%s\n", k, strings.Join(sorted, " "))
		} else {
			fmt.Fprintf(w, "%s:=%s\n", k, v.value.(string))
		}
	}
}

// Parses, resolves, and executes a Starlark file.
// filename and src parameters are as for starlark.ExecFile:
// filename is the name of the file to execute,
// and the name that appears in error messages;
// src is an optional source of bytes to use instead of filename
// (it can be a string, or a byte array, or an io.Reader instance)
func Run(filename string, src interface{}) error {
	configVariables = make(map[string]configValue)
	builtins = starlark.StringDict	{
		"setFinal": starlark.NewBuiltin("setFinal", setFinal),
		"set": starlark.NewBuiltin("set", set),
		"appendTo": starlark.NewBuiltin("appendTo", appendTo),
		"loadGenerated": starlark.NewBuiltin("loadGenerated", loadGenerated),
		"module":  starlark.NewBuiltin("module", starlarkstruct.MakeModule),
	}
	mainThread := &starlark.Thread{Name: "main",Load: loader}
	absPath, err := filepath.Abs(filename)
	if err == nil {
		mainThread.SetLocal(callerDirKey, filepath.Dir(absPath))
		_, err = starlark.ExecFile(mainThread, absPath, src, builtins)
	}
	return err
}
