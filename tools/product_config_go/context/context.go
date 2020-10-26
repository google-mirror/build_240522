package context

import (
	"fmt"
	"go.starlark.net/starlark"
	"io"
	"os/exec"
	"sort"
	"strings"
)

type configValue struct {
	value interface{}
	frozen bool
	list bool
}

var configVariables map[string]configValue

type modentry struct {
	globals starlark.StringDict
	err error
}

var moduleCache = make(map[string]*modentry)

var builtins starlark.StringDict

// loader can be passed to starlark.ExecFile to read starlark file and and execute it.
// Executed file can use the builtins.
func loader(thread *starlark.Thread, module string) (starlark.StringDict, error) {
	e, ok := moduleCache[module]
	if e == nil {
		if ok {
			return nil, fmt.Errorf("cycle in load graph")
		}

		// Add a placeholder to indicate "load in progress".
		moduleCache[module] = nil

		// Load it.
		thread := &starlark.Thread{Name: "exec " + module, Load: thread.Load}
		globals, err := starlark.ExecFile(thread, module, nil, builtins)
		e = &modentry{globals, err}

		// Update the cache.
		moduleCache[module] = e
	}
	return e.globals, e.err
}

func unpackAndSet(b *starlark.Builtin, args starlark.Tuple,
	kwargs []starlark.Tuple, freeze bool) error {
	var cvName, cvValue string
	if err := starlark.UnpackArgs(b.Name(), args, kwargs,
		"n", &cvName, "v", &cvValue); err != nil {
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
	if err := starlark.UnpackArgs(b.Name(), args, kwargs, "n", &cvName, "v", &cvValue); err != nil {
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
			fmt.Fprintf(w, "%s:=%s\n", k, strings.Join(v.value.([]string), " "))
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
	}

	_, err := starlark.ExecFile(&starlark.Thread{Name: "main",Load: loader},
		filename, src, builtins)
	return err
}
