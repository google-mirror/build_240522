package main

/*
   Canoninja reads a Ninja file and changes the rule names to be the digest of the rule contents.
   Feed  it to a filter that extracts only build statements, sort them, and you will have a crude
   but effective tool to find small differences between two Ninja files.
*/

import (
	"canoninja"
	"flag"
	"fmt"
	"os"
)

func main() {
	flag.Parse()
	files := flag.Args()
	if len(files) == 0 {
		files = []string{"/dev/stdin"}
	}
	rc := 0
	for _, f := range files {
		if err := canoninja.Generate(f); err != nil {
			//goland:noinspection GoUnhandledErrorResult
			fmt.Fprintln(os.Stderr, err)
			rc = 1
		}
	}
	os.Exit(rc)
}
