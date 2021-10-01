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
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var (
	graphViz        = flag.Bool("dot", false, "Whether to output graphviz (i.e. dot) format.")
	labelConditions = flag.Bool("label_conditions", false, "Whether to label target nodes with conditions.")
	stripPrefix     = flag.String("strip_prefix", "", "Prefix to remove from paths. i.e. path to root")
)

type context struct {
	graphViz        bool
	labelConditions bool
	stripPrefix     string
}

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: %s {options} file.meta_lic {file.meta_lic...}

Outputs space-separated Target Dependency Annotations tuples for each
edge in the license graph.

In plain text mode, multiple values within a field are colon-separated.
e.g. multiple annotations appear as annotation1:annotation2:annotation3
or when -label_conditions is requested, Target and Dependency become
target:condition1:condition2 etc.

Options:
`, filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
}

// byTargetDep orders target edges by Target then Dependency then sorted Annotations.
type byTargetDep []compliance.TargetEdge

func (l byTargetDep) Len() int      { return len(l) }
func (l byTargetDep) Swap(i, j int) { l[i], l[j] = l[j], l[i] }
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

// byName orders target nodes by Name
type byName []compliance.TargetNode

func (l byName) Len() int           { return len(l) }
func (l byName) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
func (l byName) Less(i, j int) bool { return l[i].Name() < l[j].Name() }

func main() {
	flag.Parse()

	// Must specify at least one root target.
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(2)
	}

	ctx := &context{*graphViz, *labelConditions, *stripPrefix}

	err := dumpGraph(ctx, os.Stdout, os.Stderr, flag.Args()...)
	if err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}

// dumpGraph implements the dumpgraph utility.
func dumpGraph(ctx *context, stdout, stderr io.Writer, files ...string) error {
	if len(files) < 1 {
		return fmt.Errorf("no license metadata files requested")
	}

	// Read the license graph from the license metadata files (*.meta_lic).
	licenseGraph, err := compliance.ReadLicenseGraph(files)
	if err != nil {
		fmt.Fprintf(stderr, "Unable to read license metadata file(s) %q: %w\n", files, err)
		return err
	}
	if licenseGraph == nil {
		fmt.Fprintf(stderr, "No licenses\n")
		return fmt.Errorf("no licenses")
	}

	// Sort the edges of the graph.
	edges := licenseGraph.Edges()
	sort.Sort(byTargetDep(edges))

	// nodes maps license metadata file names to graphViz node names when ctx.graphViz is true.
	var nodes map[string]string
	n := 0

	// targetOut calculates the string to output for `target` separating conditions as needed using `sep`.
	targetOut := func(target compliance.TargetNode, sep string) string {
		tOut := strings.TrimPrefix(target.Name(), ctx.stripPrefix)
		if ctx.labelConditions {
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

	// makeNode maps `target` to a graphViz node name.
	makeNode := func(target compliance.TargetNode) {
		tName := target.Name()
		if _, ok := nodes[tName]; !ok {
			nodeName := fmt.Sprintf("n%d", n)
			nodes[tName] = nodeName
			fmt.Fprintf(stdout, "\t%s [label=\"%s\"];\n", nodeName, targetOut(target, "\\n"))
			n++
		}
	}

	// If graphviz output, map targets to node names, and start the directed graph.
	if ctx.graphViz {
		nodes = make(map[string]string)
		targets := licenseGraph.Targets()
		sort.Sort(byName(targets))

		fmt.Fprintf(stdout, "strict digraph {\n\trankdir=RL;\n")
		for _, target := range targets {
			makeNode(target)
		}
	}

	// Print the sorted edges to stdout ...
	for _, e := range edges {
		// sort the annotations for repeatability/stability
		annotations := e.Annotations().AsList()
		sort.Strings(annotations)

		tName := e.Target().Name()
		dName := e.Dependency().Name()

		if ctx.graphViz {
			// ... one edge per line labelled with \\n-separated annotations.
			tNode := nodes[tName]
			dNode := nodes[dName]
			fmt.Fprintf(stdout, "\t%s -> %s [label=\"%s\"];\n", dNode, tNode, strings.Join(annotations, "\\n"))
		} else {
			// ... one edge per line with annotations in a colon-separated tuple.
			fmt.Fprintf(stdout, "%s %s %s\n", targetOut(e.Target(), ":"), targetOut(e.Dependency(), ":"), strings.Join(annotations, ":"))
		}
	}

	// If graphViz output, rank the root nodes together, and complete the directed graph.
	if ctx.graphViz {
		fmt.Fprintf(stdout, "\t{rank=same;")
		for _, f := range files {
			fName := f
			if !strings.HasSuffix(fName, ".meta_lic") {
				fName += ".meta_lic"
			}
			if fNode, ok := nodes[fName]; ok {
				fmt.Fprintf(stdout, " %s", fNode)
			}
		}
		fmt.Fprintf(stdout, "}\n}\n")
	}
	return nil
}
