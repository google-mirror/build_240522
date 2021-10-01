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
	return l[i].AppliesTo().Name() < l[j].AppliesTo().Name() || (
		l[i].AppliesTo().Name() == l[j].AppliesTo().Name() && l[i].Name() < l[j].Name())
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

	// resolutions maps target names to the conditions they must satisfy.
	resolutions := make(map[string]compliance.LicenseConditionSet)

	// Even when a target includes code from a reciprocal project, only the source
	// for the reciprocal project needs to be shared.
	recip := licenseGraph.WalkDepsForCondition("reciprocal")
	for _, t := range recip.Targets() {
		conditions, _ := recip.Conditions(t)
		for _, c := range conditions.Conditions() {
			appliesTo := c.AppliesTo()  // only need to record the AppliesTo targets.
			if _, ok := resolutions[appliesTo.Name()]; !ok {
				resolutions[appliesTo.Name()] = compliance.NewLicenseConditionSet(c)
			}
		}
	}

	// When a target includes code from a restricted project, the restricted conditon
	// requires sharing both the restricted project and the including project(s).
	restrict := licenseGraph.WalkRestricted(*fullWalk /* revisit nodes to identify all resolutions? */)
	for _, t := range restrict.Targets() {
		conditions, _ := restrict.Conditions(t)
		if _, ok := resolutions[t.Name()]; !ok {
			resolutions[t.Name()] = compliance.NewLicenseConditionSet(conditions.Conditions()...)
			continue
		}
		resolutions[t.Name()].Add(conditions.Conditions()...)
		// no need to do iterate the conditions and record `AppliesTo` because `restrct`, the walk result,
		// will include separate resolutions for them.
	}

	// group the target resolutions by project
	presolution := make(map[string]compliance.LicenseConditionSet)
	for target, rs := range resolutions {
		for _, p := range licenseGraph.TargetNode(target).Projects() {
			if _, ok := presolution[p]; !ok {
				presolution[p] = compliance.NewLicenseConditionSet(rs.Conditions()...)
				continue
			}
			presolution[p].Add(rs.Conditions()...)
		}
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
			fmt.Fprintf(os.Stdout, ",%s:%s", c.AppliesTo().Name(), c.Name())
		}
		fmt.Fprintf(os.Stdout, "\n")
	}
	os.Exit(0)
}
