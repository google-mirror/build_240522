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
	"sort"
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: %s file.meta_lic {file.meta_lic...}

Outputs a csv file with 1 project per line in the first field followed
by target:condition pairs describing why the project must be shared.

Each target is the path to a generated license metadata file for a
Soong module or Make target, and the license condition is either
restricted (e.g. GPL) or reciprocal (e.g. MPL).

%s always output all of the projects. When the -full_walk flag
is given, %s also includes a complete list of triggering targets
and conditions. This takes longer.

Options:
`, filepath.Base(os.Args[0]), filepath.Base(os.Args[0]), filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
}

type byTargetName []compliance.TargetNode

func (l byTargetName) Len() int           { return len(l) }
func (l byTargetName) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
func (l byTargetName) Less(i, j int) bool {
	return l[i].Name() < l[j].Name()
}

type byAppliesTo []compliance.LicenseCondition

func (l byAppliesTo) Len() int           { return len(l) }
func (l byAppliesTo) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
func (l byAppliesTo) Less(i, j int) bool {
	if l[i].AppliesTo().Name() == l[j].AppliesTo().Name() {
		return l[i].Name() < l[j].Name()
	}
	return l[i].AppliesTo().Name() < l[j].AppliesTo().Name()
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

	resolutions := compliance.ResolveGraphConditions(licenseGraph)
	targets := resolutions.Targets()
	sort.Sort(byTargetName(targets))
	for _, t := range targets {
		cs, _ := resolutions.Conditions(t)
		conditions := cs.Conditions()
		sort.Sort(byAppliesTo(conditions))
		for _, c := range conditions {
			fmt.Fprintf(os.Stdout, "%s %s %s\n", t.Name(), c.AppliesTo().Name(), c.Name())
		}
	}

	os.Exit(0)
}
