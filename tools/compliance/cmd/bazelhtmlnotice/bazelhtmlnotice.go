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
	"compress/gzip"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	inputFile  string
	outputFile = flag.String("o", "", "output file")
)

type LicenseKind struct {
	Target     string   `json:"target"`
	Name       string   `json:"name"`
	Conditions []string `json:"conditions"`
}

type License struct {
	Rule            string        `json:"rule"`
	CopyrightNotice string        `json:"copyright_notice"`
	PackageName     string        `json:"package_name"`
	PackageUrl      string        `json:"package_url"`
	PackageVersion  string        `json:"package_version"`
	LicenseFile     string        `json:"license_text"`
	LicenseKinds    []LicenseKind `json:"license_kinds"`
}

func maybeQuit(err error) {
	if err == nil {
		return
	}

	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func doit(in string, sink io.Writer) {
	decoder := json.NewDecoder(strings.NewReader(in))
	decoder.DisallowUnknownFields() //useful to detect typos, e.g. in unit tests
	var licenses []License
	err := decoder.Decode(&licenses)
	maybeQuit(err)
	fmt.Fprint(sink, `<!DOCTYPE html>
<html>
  <head>
	<style type="text/css"
      body { padding: 2px; margin: 0; }
      ul { list-style-type: none; margin: 0; padding: 0; }
      li { padding-left: 1em; }
      .file-list { margin-left: 1em; }
    </style>
  </head>
  <body>
  </body>
</html>
`)
}

func processArgs() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, `usage: bazelhtmlnotice -o <output> <input>`)
		flag.PrintDefaults()
		os.Exit(2)
	}
	flag.Parse()
	if len(flag.Args()) != 1 {
		flag.Usage()
	}
	inputFile = flag.Arg(0)
}

func setupWriting() (io.Writer, io.Closer, *os.File) {
	if *outputFile == "" {
		return os.Stdout, nil, nil
	}
	ofile, err := os.Create(*outputFile)
	maybeQuit(err)
	if !strings.HasSuffix(*outputFile, ".gz") {
		return ofile, nil, ofile
	}
	gz, err := gzip.NewWriterLevel(ofile, gzip.BestCompression)
	maybeQuit(err)
	return gz, gz, ofile
}

func main() {
	processArgs()
	data, err := os.ReadFile(inputFile)
	maybeQuit(err)
	sink, closer, ofile := setupWriting()
	doit(string(data), sink)
	if closer != nil {
		maybeQuit(closer.Close())
	}
	if ofile != nil {
		maybeQuit(ofile.Close())
	}
}
