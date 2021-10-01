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
		fmt.Fprintf(os.Stderr, `Usage: %s {-full_walk} file.meta_lic {file.meta_lic...}

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

var (
	fullWalk = flag.Bool("full_walk", true, "Whether to re-visit nodes to identify all resolutions.")
)

// ConditionList implements introspection methods to arrays of LicenseCondition.
type ConditionList []compliance.LicenseCondition

// HasByName returns true if the list contains any condition matching `name`.
func (cl ConditionList) HasByName(name compliance.ConditionNames) bool {
	for _, c := range cl {
		if name.Contains(c.Name()) {
			return true
		}
	}
	return false
}

// ByName returns the sublist of conditions that match `name`.
func (cl ConditionList) ByName(name compliance.ConditionNames) ConditionList {
	result := make(ConditionList, 0)
	for _, c := range cl {
		if name.Contains(c.Name()) {
			result = append(result, c)
		}
	}
	return result
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

	// recip is the set of resolutions for the Reciprocal condition.
	recip := compliance.ResolveTopDownCondition(licenseGraph, "reciprocal")

	// proprietary is the set of resolutions for the Proprietary condition.
	proprietary := compliance.ResolveTopDownCondition(licenseGraph, "proprietary")

	// restrict is the set of resolutions for the Restricted condition.
	restrict := compliance.ResolveTopDownRestricted(licenseGraph, *fullWalk /* revisit nodes to identify all resolutions? */)

	// shareSource is the set of all source-sharing and proprietary resolutions.
	shareSource := compliance.JoinResolutions(recip, restrict, proprietary)

	// report outputs an error message to stderr
	report := func(target compliance.TargetNode, p, ss compliance.LicenseCondition) {
		fmt.Fprintf(os.Stderr, "%s %s from %s and must share from %s %s\n", target.Name(), p.Name(), p.Origin().Name(), ss.Name(), ss.Origin().Name())
	}

	// status is the exit status for the program. Assumes success until a conflict found.
	status := 0
	for _, t := range shareSource.Targets() {
		cs, _ := shareSource.Conditions(t)

		// allconditions includes all of the conditions for the target
		allconditions := ConditionList(cs.Conditions())

		// ownconditions includes only the conditions originating at the target
		ownconditions := ConditionList(cs.ByOrigin(t))

		// pconditions identifies the conditions to examine for Proprietary: own for containers, all for everything else
		var pconditions ConditionList
		if t.IsContainer() {
			// A pure aggregate can contain both proprietary targets and targets with source-sharing requirements.
			// However, if the source-code for the aggregate itself is proprietary, it cannot inherit source-sharing requirements.
			// For the proprietary check, only look at license conditions originating with the container itself.
			pconditions = ownconditions
		} else {
			// Works cannot be derivative of proprietary code and have source-sharing requirements
			pconditions = allconditions
		}
		if pconditions.HasByName(compliance.Proprietary) {
			// reciprocal only matters when the condition originates at the target, but restricted can originate anywhere.
			if allconditions.HasByName(compliance.Restricted) || ownconditions.HasByName(compliance.Reciprocal) {
				// one or more conflicts detected
				status = 1
				// report all conflicting condition combinations
				for _, p := range pconditions.ByName(compliance.Proprietary) {
					for _, ss := range allconditions.ByName(compliance.Restricted) {
						report(t, p, ss)
					}
					for _, ss := range ownconditions.ByName(compliance.Reciprocal) {
						report(t, p, ss)
					}
				}
			}
		}
	}

	// indicate pass or fail on stdout
	if status == 0 {
		fmt.Fprintln(os.Stdout, "PASS")
	} else {
		fmt.Fprintln(os.Stdout, "FAIL")
	}
	os.Exit(status)
}
