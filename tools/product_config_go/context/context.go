package context

import (
	"fmt"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

type configValue struct {
	value  interface{}
	frozen bool
	list   bool
}

const callerDirKey = "callerDir"
const cliModuleName = "cli"

var configVariables map[string]configValue

var LoadPathRoot = "."

type modentry struct {
	globals starlark.StringDict
	err     error
}

var moduleCache = make(map[string]*modentry)

var builtins starlark.StringDict

func moduleName2AbsPath(moduleName string, callerDir string) (string, error) {
	path := moduleName
	if ix := strings.LastIndex(path, ":"); ix >= 0 {
		path = path[0:ix] + string(os.PathSeparator) + path[ix+1:]
	}
	if strings.HasPrefix(path, "//") {
		return filepath.Abs(filepath.Join(LoadPathRoot, path[2:]))
	} else if strings.HasPrefix(moduleName, ":") {
		return filepath.Abs(filepath.Join(callerDir, path[1:]))
	} else {
		return filepath.Abs(path)
	}
}

// loader can be passed to starlark.ExecFile to read starlark file and execute it.
// Executed file can use the builtins.
func loader(thread *starlark.Thread, module string) (starlark.StringDict, error) {
	// The module name's format is [//path]:base[|symbol]
	// The file path is $ROOT/path/base if path is present, <caller_dir>/base otherwise.
	// The presence of `|symbol` indicates that the loader should return a single symbol
	// bound to None if file is missing.
	pipePos := strings.LastIndex(module, "|")
	mustLoad := pipePos < 0
	var defaultSymbol string
	if !mustLoad {
		defaultSymbol = module[pipePos+1:]
		module = module[:pipePos]
	}
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

		// Decide if we should load.
		if !mustLoad {
			if _, err := os.Stat(modulePath); err == nil {
				mustLoad = true
			}
		}

		// Load or return default
		if mustLoad {
			thread := &starlark.Thread{Name: "exec " + module, Load: thread.Load}
			thread.SetLocal(callerDirKey, filepath.Dir(modulePath))
			globals, err := starlark.ExecFile(thread, modulePath, nil, builtins)
			e = &modentry{globals, err}
		} else {
			e = &modentry{starlark.StringDict{defaultSymbol: starlark.None}, nil}
		}

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
	configVariables[cvName] = configValue{value: cvValue, frozen: freeze}
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
	configVariables[cvName] = configValue{value: sl, list: true}
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
	cmd := exec.Command(command, cmdArgs...)
	bytes, err := cmd.Output()
	if err != nil {
		return starlark.None, fmt.Errorf("'%s' failed: %s", cmd, err)
	}
	loadGenThread := &starlark.Thread{
		Name: "loadGenerated " + command,
		Load: loader,
	}
	_, err = starlark.ExecFile(loadGenThread, fmt.Sprintf("$(%s)", cmd),
		string(bytes), builtins)
	if err != nil {
		return starlark.None, err
	}
	return starlark.None, nil
}

// fileExists returns True if file with given name exists.
func fileExists(_ *starlark.Thread, b *starlark.Builtin, args starlark.Tuple,
	kwargs []starlark.Tuple) (starlark.Value, error) {
	var path string
	if err := starlark.UnpackPositionalArgs(b.Name(), args, kwargs, 1, &path); err != nil {
		return starlark.None, err
	}
	if stat, err := os.Stat(path); err != nil || stat.IsDir() {
		return starlark.False, nil
	}
	return starlark.True, nil
}

// wildcard(pattern, top) finds all the files matching pattern under top
// (in other words, it's like running `cd top && find -name pattern` and removing
// "./" at the beginning of each line).
func wildcard(_ *starlark.Thread, b *starlark.Builtin, args starlark.Tuple,
	kwargs []starlark.Tuple) (starlark.Value, error) {
	var pattern, top string
	if err := starlark.UnpackPositionalArgs(b.Name(), args, kwargs, 2, &pattern, &top); err != nil {
		return starlark.None, err
	}
	var files []string
	err := filepath.Walk(top, func(path string, info os.FileInfo, err error) error {
		if match, err2 := filepath.Match(pattern, info.Name()); err2 == nil && match {
			if relPath, err3 := filepath.Rel(top, path); err3 == nil {
				files = append(files, relPath)
			}
		}
		return nil
	})
	if err != nil {
		return starlark.None, err
	}
	return makeStringList(files), nil
}

func makeStringList(items []string) *starlark.List {
	elems := make([]starlark.Value, len(items))
	for i, item := range items {
		elems[i] = starlark.String(item)
	}
	return starlark.NewList(elems)
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

	for _, k := range keys {
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
// env is an array of "VAR=value" items. They are accessible from
// the starlark script as members of the `cli` module.
func Run(filename string, src interface{}, env []string) error {
	configVariables = make(map[string]configValue)
	cliVariables := make(starlark.StringDict, len(env))
	for _, s := range env {
		pair := strings.SplitN(s, "=", 2)
		cliVariables[pair[0]] = starlark.String(pair[1])
	}
	builtins = starlark.StringDict{
		"setFinal":         starlark.NewBuiltin("setFinal", setFinal),
		"set":              starlark.NewBuiltin("set", set),
		"appendTo":         starlark.NewBuiltin("appendTo", appendTo),
		"loadGenerated":    starlark.NewBuiltin("loadGenerated", loadGenerated),
		"module":           starlark.NewBuiltin("module", starlarkstruct.MakeModule),
		"rblf_file_exists": starlark.NewBuiltin("rblf_file_exists", fileExists),
		"rblf_wildcard":    starlark.NewBuiltin("rblf_wildcard", wildcard),
		cliModuleName:      &starlarkstruct.Module{cliModuleName, cliVariables},
	}
	mainThread := &starlark.Thread{
		Name:  "main",
		Print: func(_ *starlark.Thread, msg string) { fmt.Println(msg) },
		Load:  loader,
	}
	absPath, err := filepath.Abs(filename)
	if err == nil {
		mainThread.SetLocal(callerDirKey, filepath.Dir(absPath))
		_, err = starlark.ExecFile(mainThread, absPath, src, builtins)
	}
	return err
}
