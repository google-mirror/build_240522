package main

import (
	"flag"
	"fmt"
	"os"
	"product_config_go/context"
)

var (
	execprog   = flag.String("c", "", "execute program `prog`")
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

	if err := context.Run(filename, src); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		return 2
	}
	context.PrintConfig(os.Stdout)
	return 0
}