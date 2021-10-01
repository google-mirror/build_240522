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
		fmt.Fprintf(os.Stderr, `Usage: %s {-full_walk} file.meta_lic {file.meta_lic...}

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

var (
	fullWalk = flag.Bool("full_walk", false, "Whether to re-visit nodes to identify all resolutions.")
)

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

	// Even when a target includes code from a reciprocal project, only the source
	// for the reciprocal project, i.e. the origin of the condition, needs to be shared.
	recip := compliance.ResolveTopDownCondition(licenseGraph, "reciprocal")

	// When a target includes code from a restricted project, the restricted conditon
	// requires sharing both the restricted project and the including project(s).
	restrict := compliance.ResolveTopDownRestricted(licenseGraph, *fullWalk /* revisit nodes to identify all resolutions? */)

	// shareSource is the set of all source-sharing resolutions.
	shareSource := compliance.JoinResolutions(recip, restrict)

	// group the resolutions by project
	presolution := make(map[string]compliance.LicenseConditionSet)
	for _, t := range shareSource.Targets() {
		conditions, _ := shareSource.Conditions(t)
		for _, c := range conditions.Conditions() {
			// Always share the origin of the condition
			for _, p := range c.Origin().Projects() {
				if _, ok := presolution[p]; !ok {
					presolution[p] = compliance.NewLicenseConditionSet(c)
					continue
				}
				presolution[p].Add(c)
			}
			// For reciprocal, only the origin source needs to be shared.
			if compliance.Reciprocal.Contains(c.Name()) {
				continue
			}
			// no need to tag the same target twice
			if t.Name() == c.Origin().Name() {
				continue
			}
			// For restricted, share the project the resolution applies to as well as the origin.
			for _, p := range t.Projects() {
				if _, ok := presolution[p]; !ok {
					presolution[p] = compliance.NewLicenseConditionSet(c)
					continue
				}
				presolution[p].Add(c)
			}
		}
	}

	// sort the projects to make the results more usable
	projects := make([]string, 0, len(presolution))
	for p := range presolution {
		projects = append(projects, p)
	}
	sort.Strings(projects)

	// output the sorted projects and the license conditions sharing the project resolves.
	for _, p := range projects {
		fmt.Fprintf(os.Stdout, "%s", p)
		cs := presolution[p]
		conditions := cs.Conditions()
		sort.Sort(byTarget(conditions))
		for _, c := range conditions {
			fmt.Fprintf(os.Stdout, ",%s:%s", c.Origin().Name(), c.Name())
		}
		fmt.Fprintf(os.Stdout, "\n")
	}
	os.Exit(0)
}
