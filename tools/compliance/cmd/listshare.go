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

List the projects that must be shared for OS license compliance when releasing
targets whose license metadata is represented by the argument filenames.

Options:
`, filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
}

var (
	fullWalk = flag.Bool("full_walk", false, "Whether to re-visit nodes to identify all resolutions.")
)

func main() {
	flag.Parse()

	// FIXME: remove debug output
	fmt.Fprintf(os.Stderr, "Environment:\n")
	for _, ev := range os.Environ() {
		fmt.Fprintf(os.Stderr, "%s\n", ev)
	}
	fmt.Fprintf(os.Stderr, "\n")

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
	fmt.Fprintln(os.Stderr, "graph read")
	resolutions := make(map[string]compliance.LicenseConditionSet)
	recip := licenseGraph.WalkDepsForCondition("reciprocal")
	fmt.Fprintln(os.Stderr, "reciprocal walked")
	for _, t := range recip.Targets() {
		conditions, _ := recip.Conditions(t)
		for _, c := range conditions.Conditions() {
			appliesTo := c.AppliesTo()
			if _, ok := resolutions[appliesTo.Name()]; !ok {
				resolutions[appliesTo.Name()] = compliance.NewLicenseConditionSet(c)
			}
		}
	}
	fmt.Fprintln(os.Stderr, "reciprocal resolved")
	restrict := licenseGraph.WalkRestricted(*fullWalk /* revisit nodes to identify all resolutions? */)
	fmt.Fprintln(os.Stderr, "restricted walked")
	for _, t := range restrict.Targets() {
		conditions, _ := restrict.Conditions(t)
		if _, ok := resolutions[t.Name()]; !ok {
			resolutions[t.Name()] = compliance.NewLicenseConditionSet(conditions.Conditions()...)
			continue
		}
		resolutions[t.Name()].Add(conditions.Conditions()...)
	}
	fmt.Fprintln(os.Stderr, "restricted resolved")
	presolution := make(map[string]compliance.LicenseConditionSet)
	for _, rs := range resolutions {
		for _, c := range rs.Conditions() {
			for _, p := range c.AppliesTo().Projects() {
				if _, ok := presolution[p]; !ok {
					presolution[p] = compliance.NewLicenseConditionSet(c)
					continue
				}
				presolution[p].Add(c)
			}
		}
	}
	fmt.Fprintln(os.Stderr, "projects resolved")
	projects := make([]string, 0, len(presolution))
	for p := range presolution {
		projects = append(projects, p)
	}
	sort.Strings(projects)
	fmt.Fprintln(os.Stderr, "projects sorted")
	for _, p := range projects {
		fmt.Fprintf(os.Stdout, "%s", p)
		cs := presolution[p]
		for _, c := range cs.Conditions() {
			fmt.Fprintf(os.Stdout, ",%s:%s", c.AppliesTo().Name(), c.Name())
		}
		fmt.Fprintf(os.Stdout, "\n")
	}
	os.Exit(0)
}

