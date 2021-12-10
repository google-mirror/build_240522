// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"compliance"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

var (
	outputFile      = flag.String("o", "-", "Where to write the NOTICE text file. (default stdout)")
	includeHeader   = flag.Bool("header", false, "Whether to print a table of contents.")
	stripPrefix     = flag.String("strip_prefix", "", "Prefix to remove from paths. i.e. path to root")

	failNoneRequested = fmt.Errorf("\nNo license metadata files requested")
	failNoLicenses = fmt.Errorf("No licenses found")
)

type context struct {
	stdout io.Writer
	stderr io.Writer
	rootFS fs.FS
	includeHeader bool
	stripPrefix string
}

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: %s {options} file.meta_lic {file.meta_lic...}

Outputs a text NOTICE file.

Options:
`, filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	// Must specify at least one root target.
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(2)
	}

	if len(*outputFile) == 0 {
		flag.Usage()
		fmt.Fprintf(os.Stderr, "must specify file for -o; use - for stdout\n")
		os.Exit(2)
	} else {
		dir, err := filepath.Abs(filepath.Dir(*outputFile))
		if err != nil {
			fmt.Fprintf(os.Stderr, "cannot determine path to %q: %w\n", *outputFile, err)
			os.Exit(1)
		}
		fi, err := os.Stat(dir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "cannot read directory %q of %q: %w\n", dir, *outputFile, err)
			os.Exit(1)
		}
		if !fi.IsDir() {
			fmt.Fprintf(os.Stderr, "parent %q of %q is not a directory\n", dir, *outputFile)
			os.Exit(1)
		}
	}

	var ofile io.Writer
	ofile = os.Stdout
	if *outputFile != "-" {
		ofile = &bytes.Buffer{}
	}

	ctx := &context{ofile, os.Stderr, os.DirFS("."), *includeHeader, *stripPrefix}

	err := textNotice(ctx, flag.Args()...)
	if err != nil {
		if err == failNoneRequested {
			flag.Usage()
		}
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
	if *outputFile != "-" {
		err := os.WriteFile(*outputFile, ofile.(*bytes.Buffer).Bytes(), 0666)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not write output to %q: %w\n", *outputFile, err)
			os.Exit(1)
		}
	}
	os.Exit(0)
}

// textNotice implements the textNotice utility.
func textNotice(ctx *context, files ...string) error {
	// Must be at least one root file.
	if len(files) < 1 {
		return failNoneRequested
	}

	fmt.Fprintf(os.Stderr, "reading graph\n")
	// Read the license graph from the license metadata files (*.meta_lic).
	licenseGraph, err := compliance.ReadLicenseGraph(ctx.rootFS, ctx.stderr, files)
	if err != nil {
		return fmt.Errorf("Unable to read license metadata file(s) %q: %v\n", files, err)
	}
	if licenseGraph == nil {
		return failNoLicenses
	}

	fmt.Fprintf(os.Stderr, "resolving notices\n")
	// rs contains all notice resolutions.
	rs := compliance.ResolveNotices(licenseGraph)

	fmt.Fprintf(os.Stderr, "indexing licenes texts\n")
	ni, err := compliance.IndexLicenseTexts(ctx.rootFS, licenseGraph, rs)
	if err != nil {
		return fmt.Errorf("Unable to read license text file(s) for %q: %v\n", files, err)
	}

	for _, r := range ni.Resolutions() {
		for _, p := range r.InstallPaths() {
			fmt.Fprintf(ctx.stdout, "%s %s\n", strings.TrimPrefix(p, ctx.stripPrefix), ni.LibraryName(r))
		}
	}
	return nil
}
