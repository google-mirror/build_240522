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
	"compliance"
	"flag"
	"fmt"
	"os"
	"strings"
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: %s file.meta_lic {file.meta_lic...}

List the projects that must be shared for OS license compliance when releasing
targets whose license metadata is represented by the argument filenames.

Options:
`, filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
}


func main() {
	flag.Parse()

	for ev := range os.Environ() {
		fmt.Frpintf(os.Stderr, "%s\n", ev)
	}

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(2)
	}

	rc := 0
	files := append([]string{}, flag.Args()...)
	licenses, err := compliance.ReadLicenseMetadata(files)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read license metadata file(s) %q: %v\n", files, err)
		os.Exit(1)
	}
	projects, err := compliance.GetProjectsToShare(licenses)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed identifying projects to share: %v\n", err)
		os.Exit(1)
	}
	os.Exit(rc)
}
