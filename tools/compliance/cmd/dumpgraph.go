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
	graphViz = flag.Bool("dot", false, "Whether to output graphviz (i.e. dot) format.")
	labelConditions = flag.Bool("label_conditions", false, "Whether to label target nodes with conditions.")
	stripPrefix = flag.String("strip_prefix", "", "Prefix to remove from paths. i.e. path to root")
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

	var nodes map[string]string
	n := 0

	targetOut := func(target compliance.TargetNode, sep string) string {
		tOut := strings.TrimPrefix(target.Name(), *stripPrefix)
		if *labelConditions {
			conditions := make([]string, 0, target.LicenseConditions().Count())
			for _, lc := range target.LicenseConditions().AsList() {
				conditions = append(conditions, lc.Name())
			}
			sort.Strings(conditions)
			if len(conditions) > 0 {
				tOut += sep + strings.Join(conditions, sep)
			}
		}
		return tOut
	}

	makeNode := func(target compliance.TargetNode) {
		tName := target.Name()
		if _, ok := nodes[tName]; !ok {
			nodeName := fmt.Sprintf("n%d", n)
			nodes[tName] = nodeName
			fmt.Fprintf(os.Stdout, "\t%s [label=\"%s\"];\n", nodeName, targetOut(target, "\\n"))
			n++
		}
	}

	// if graphviz output, start the directed graph
	if *graphViz {
		fmt.Fprintf(os.Stdout, "strict digraph {\n\trankdir=RL;\n")
		nodes = make(map[string]string)
		for _, e := range edges {
			makeNode(e.Target())
			makeNode(e.Dependency())
		}
	}

	// Print the sorted edges to stdout ...
	for _, e := range edges {
		// sort the annotations for repeatability/stability
		annotations := e.Annotations().AsList()
		sort.Strings(annotations)

		tName := e.Target().Name()
		dName := e.Dependency().Name()

		if *graphViz {
			// ... one edge per line labelled with \\n-separated annotations.
			tNode := nodes[tName]
			dNode := nodes[dName]
			fmt.Fprintf(os.Stdout, "\t%s -> %s [label=\"%s\"];\n", dNode, tNode, strings.Join(annotations, "\\n"))
		} else {
			// ... one edge per line with annotations in a colon-separated tuple.
			fmt.Fprintf(os.Stdout, "%s %s %s\n", targetOut(e.Target(), ":"), targetOut(e.Dependency(), ":"), strings.Join(annotations, ":"))
		}
	}

	if *graphViz {
		fmt.Fprintf(os.Stdout, "\t{rank=same;")
		for _, f := range files {
			fName := f
			if !strings.HasSuffix(fName, ".meta_lic") {
				fName += ".meta_lic"
			}
			if fNode, ok := nodes[fName]; ok {
				fmt.Fprintf(os.Stdout, " %s", fNode)
			}
		}
		fmt.Fprintf(os.Stdout, "}\n}\n")
	}
	os.Exit(0)
}
