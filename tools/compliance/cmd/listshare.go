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

Options:
`, filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
}

// byTarget orders license conditions by originating target then condition name.
type byTarget []compliance.LicenseCondition

func (l byTarget) Len() int           { return len(l) }
func (l byTarget) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
func (l byTarget) Less(i, j int) bool {
	return l[i].Origin().Name() < l[j].Origin().Name() || (
		l[i].Origin().Name() == l[j].Origin().Name() && l[i].Name() < l[j].Name())
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

	// files identifies the roots of the licens graph.
	files := append([]string{}, flag.Args()...)

	// Read the license graph from the license metadata files (*.meta_lic).
	licenseGraph, err := compliance.ReadLicenseGraph(files)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read license metadata file(s) %q: %v\n", files, err)
		os.Exit(1)
	}
	if licenseGraph == nil {
		fmt.Fprintf(os.Stderr, "No licenses\n")
		os.Exit(1)
	}

	// shareSource contains all source-sharing resolutions.
	shareSource := compliance.ResolveSourceSharing(licenseGraph)

	// Group the resolutions by project.
	presolution := make(map[string]compliance.LicenseConditionSet)
	for _, target := range shareSource.AppliesTo() {
		for _, c := range shareSource.Conditions(target).AsList() {
			for _, p := range target.Projects() {
				if _, ok := presolution[p]; !ok {
					presolution[p] = compliance.NewLicenseConditionSet(c)
					continue
				}
				presolution[p].Add(c)
			}
		}
	}

	// Sort the projects for repeatability/stability.
	projects := make([]string, 0, len(presolution))
	for p := range presolution {
		projects = append(projects, p)
	}
	sort.Strings(projects)

	// Output the sorted projects and the source-sharing license conditions that each project resolves.
	for _, p := range projects {
		fmt.Fprintf(os.Stdout, "%s", p)

		// Sort the conditions for repeatability/stability.
		conditions := presolution[p].AsList()
		sort.Sort(byTarget(conditions))

		// Output the sorted origin:condition pairs.
		for _, c := range conditions {
			fmt.Fprintf(os.Stdout, ",%s:%s", c.Origin().Name(), c.Name())
		}
		fmt.Fprintf(os.Stdout, "\n")
	}
	os.Exit(0)
}
