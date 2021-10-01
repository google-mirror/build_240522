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
	"strings"
)

var (
	conditions = newMultiString("c", "conditions to resolve")
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: %s {-c condition}... file.meta_lic {file.meta_lic...}

Outputs a space-separated Target Origin Condition tuple for each resolution in
the graph.

If one or more '-c condition' conditions are given, outputs the joined set of
resolutions for all of the conditions. Otherwise, outputs the result of the
bottom-up resolve only.

Options:
`, filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
}

// newMultiString creates a flag that allows multiple values in an array.
func newMultiString(name, usage string) *multiString {
	var f multiString
	flag.Var(&f, name, usage)
	return &f
}

// multiString implements the flag `Value` interface for multiple strings.
type multiString []string

func (ms *multiString) String() string     { return strings.Join(*ms, ", ") }
func (ms *multiString) Set(s string) error { *ms = append(*ms, s); return nil }


// byTargetName orders target nodes by name.
type byTargetName []compliance.TargetNode

func (l byTargetName) Len() int           { return len(l) }
func (l byTargetName) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
func (l byTargetName) Less(i, j int) bool {
	return l[i].Name() < l[j].Name()
}


// byOrding orders license conditions by originating target name.
type byOrigin []compliance.LicenseCondition

func (l byOrigin) Len() int           { return len(l) }
func (l byOrigin) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
func (l byOrigin) Less(i, j int) bool {
	if l[i].Origin().Name() == l[j].Origin().Name() {
		return l[i].Name() < l[j].Name()
	}
	return l[i].Origin().Name() < l[j].Origin().Name()
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

	// files identifies the roots of the license graph.
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

	// resolutions will contain the requested set of resolutions.
	var resolutions compliance.ResolutionSet

	if len(*conditions) == 0 {
		resolutions = compliance.ResolveBottomUpConditions(licenseGraph)
	} else {
		rlist := make([]compliance.ResolutionSet, 0, len(*conditions))
		for _, c := range *conditions {
			rlist = append(rlist, compliance.ResolveTopDownForCondition(licenseGraph, c))
		}
		if len(rlist) == 1 {
			resolutions = rlist[0]
		} else {
			resolutions = compliance.JoinResolutions(rlist...)
		}
	}

	// Sort the resolutions by targetname for repeatability/stability.
	targets := resolutions.AppliesTo()
	sort.Sort(byTargetName(targets))

	// Output the sorted targets.
	for _, target := range targets {
		// Sort the conditions that apply to `target` for repeatability/stability.
		conditions := resolutions.Conditions(target).AsList()
		sort.Sort(byOrigin(conditions))

		// Output 1 line for each target+condition combination.
		for _, condition := range conditions {
			fmt.Fprintf(os.Stdout, "%s %s %s\n",
				target.Name(),
				condition.Origin().Name(),
				condition.Name())
		}
	}

	os.Exit(0)
}
