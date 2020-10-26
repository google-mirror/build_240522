package main

import (
	"flag"
	"fmt"
	"go.starlark.net/starlark"
	"os"
	"product_config_go/context"
	"strings"
)

var (
	execprog = flag.String("c", "", "execute program `prog`")
	rootdir  = flag.String("d", ".", "the value of // for load paths")
)

func main() {
	os.Exit(doMain())
}

func doMain() int {
	flag.Parse()
	var filename string
	var src interface{}
	var env []string
	for _, arg := range flag.Args() {
		if strings.Contains(arg, "=") {
			env = append(env, arg)
		} else if filename == "" {
			filename = arg
		} else {
			fmt.Fprintln(os.Stderr, "only one file can be executed")
			return 2
		}
	}
	if *execprog != "" {
		if filename != "" {
			fmt.Fprintf(os.Stderr, "either -c or file name should be present")
			return 2
		}
		filename = "cmdline"
		src = *execprog
	}
	if filename == "" {
		flag.Usage()
		return 1
	}
	if stat, err := os.Stat(*rootdir); os.IsNotExist(err) || !stat.IsDir() {
		fmt.Fprintf(os.Stderr, "%s is not a directory\n", *rootdir)
		return 2
	}
	context.LoadPathRoot = *rootdir
	if err := context.Run(filename, src, env); err != nil {
		if evalErr, ok := err.(*starlark.EvalError); ok {
			fmt.Fprintln(os.Stderr, evalErr.Backtrace())
		} else {
			fmt.Fprintln(os.Stderr, err)
		}
		return 2
	}
	context.PrintConfig(os.Stdout)
	return 0
}
