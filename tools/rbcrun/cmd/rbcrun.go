package main

import (
	"flag"
	"fmt"
	"go.starlark.net/starlark"
	"os"
	"rbcrun"
	"strings"
)

var (
	execprog = flag.String("c", "", "execute program `prog`")
	rootdir  = flag.String("d", ".", "the value of // for load paths")
	perfFile = flag.String("perf", "", "save performance data")
)

func main() {
	flag.Parse()
	var filename string
	var src interface{}
	var env []string

	rc := 0
	for _, arg := range flag.Args() {
		if strings.Contains(arg, "=") {
			env = append(env, arg)
		} else if filename == "" {
			filename = arg
		} else {
			quit("only one file can be executed\n")
		}
	}
	if *execprog != "" {
		if filename != "" {
			quit("either -c or file name should be present\n")
		}
		filename = "cmdline"
		src = *execprog
	}
	if filename == "" {
		flag.Usage()
		os.Exit(1)
	}
	if stat, err := os.Stat(*rootdir); os.IsNotExist(err) || !stat.IsDir() {
		quit("%s is not a directory\n", *rootdir)
	}
	if *perfFile != "" {
		pprof, err := os.Create(*perfFile)
		if err != nil {
			quit("%s: err", *perfFile)
		}
		defer pprof.Close()
		if err := starlark.StartProfile(pprof); err != nil {
			quit("%s\n", err)
		}
	}
	rbcrun.LoadPathRoot = *rootdir
	err := rbcrun.Run(filename, src, env)
	if *perfFile != "" {
		if err2 := starlark.StopProfile(); err2 != nil {
			fmt.Fprintln(os.Stderr, err2)
			rc = 1
		}
	}
	if err != nil {
		if evalErr, ok := err.(*starlark.EvalError); ok {
			quit("%s\n", evalErr.Backtrace())
		} else {
			quit("%s\n", err)
		}
	}
	os.Exit(rc)
}

func quit(format string, s ...interface{}) {
	fmt.Fprintln(os.Stderr, format, s)
	os.Exit(2)
}
