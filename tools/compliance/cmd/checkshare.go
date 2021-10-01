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
	"path/filepath"
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: %s file.meta_lic {file.meta_lic...}

Reports on stderr any targets where policy dictates the source both
must and must not be shared. The error report indicates the target, the
license condition with origin that blocks code sharing, and the license
condition with origin that requires code sharing per policy.

Any given target may appear multiple times with different combinatiions
of conflicting license conditions.

If none of the source code that policy dictates be shared is blocked by
policy from being shared, outputs "PASS" to stdout and exits with
status 0.

If policy dictates any source must both be shared and not be shared,
outputs "FAIL" to stdout and exits with status 1.

%s always output all of the targets with conflicts. When the -full_walk
flag is given, %s also includes a complete list of triggering targets
and conditions. This takes longer, and is the default behavior.

Options:
`, filepath.Base(os.Args[0]), filepath.Base(os.Args[0]), filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	// FIXME: remove debug output
	fmt.Fprintf(os.Stderr, "Environment:\n")
	for _, ev := range os.Environ() {
		fmt.Fprintf(os.Stderr, "%s\n", ev)
	}
	fmt.Fprintf(os.Stderr, "\n")

	// Must specify at least one root target.
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(2)
	}

	files := append([]string{}, flag.Args()...)
	licenseGraph, err := compliance.ReadLicenseGraph(files)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read license metadata file(s) %q: %v\n", files, err)
		os.Exit(1)
	}
	if licenseGraph == nil {
		fmt.Fprintf(os.Stderr, "No licenses\n")
		os.Exit(1)
	}

	// conflicts contains the targets and conflicting conditions per policy
	conflicts := compliance.ConflictingSharedPrivateSource(licenseGraph)
	for _, conflict := range conflicts {
		fmt.Fprintln(os.Stderr, conflict.Error())
	}

	// indicate pass or fail on stdout
	if len(conflicts) > 0 {
		fmt.Fprintln(os.Stdout, "FAIL")
		os.Exit(1)
	}
	fmt.Fprintln(os.Stdout, "PASS")
	os.Exit(0)
}
