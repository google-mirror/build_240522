package main

import (
	"flag"
	"fmt"
	"go.starlark.net/starlark"
	"os"
	"product_config_go/context"
)

var (
	execprog   = flag.String("c", "", "execute program `prog`")
	rootdir = flag.String("d", ".", "the value of // for load paths")
)

func main() {
	os.Exit(doMain())
}

func doMain() int {
	flag.Parse()
	var filename string
	var src interface{}
	if flag.NArg() == 1 {
		filename = flag.Arg(0)
	} else if *execprog != "" {
		filename = "cmdline"
		src = *execprog
	} else {
		flag.Usage()
		return 1
	}
	if stat, err := os.Stat(*rootdir); os.IsNotExist(err) || !stat.IsDir() {
		fmt.Fprintf(os.Stderr, "%s is not a directory\n", *rootdir)
		return 2
	}
	context.LoadPathRoot = *rootdir
	if err := context.Run(filename, src); err != nil {
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