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

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: %s file.meta_lic {file.meta_lic...}

Outputs space-separated Target Dependency Annotations tuples for each edge in
the license graph.

Options:
`, filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
}

// byTargetDep orders target edges by Target then Dependency then sorted Annotations.
type byTargetDep []compliance.TargetEdge

func (l byTargetDep) Len() int           { return len(l) }
func (l byTargetDep) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
func (l byTargetDep) Less(i, j int) bool {
	if l[i].Target().Name() == l[j].Target().Name() {
		if l[i].Dependency().Name() == l[j].Dependency().Name() {
			a1 := l[i].Annotations().AsList()
			a2 := l[j].Annotations().AsList()
			sort.Strings(a1)
			sort.Strings(a2)
			for k := 0; k < len(a1) && k < len(a2); k++ {
				if a1[k] < a2[k] {
					return true
				}
				if a1[k] > a2[k] {
					return false
				}
			}
			return len(a1) < len(a2)
		}
		return l[i].Dependency().Name() < l[j].Dependency().Name()
	}
	return l[i].Target().Name() < l[j].Target().Name()
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

	// Sort the edges of the graph.
	edges := licenseGraph.Edges()
	sort.Sort(byTargetDep(edges))

	// Print the sorted edges to stdout ...
	for _, e := range edges {
		// sort the annotations for repeatability/stability
		annotations := e.Annotations().AsList()
		sort.Strings(annotations)

		// ... one edge per line with annotations in a colon-separated tuple.
		fmt.Fprintf(os.Stdout, "%s %s %s\n", e.Target().Name(), e.Dependency().Name(), strings.Join(annotations, ":"))
	}

	os.Exit(0)
}
