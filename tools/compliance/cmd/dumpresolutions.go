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
	conditions      = newMultiString("c", "conditions to resolve")
	graphViz        = flag.Bool("dot", false, "Whether to output graphviz (i.e. dot) format.")
	labelConditions = flag.Bool("label_conditions", false, "Whether to label target nodes with conditions.")
	stripPrefix     = flag.String("strip_prefix", "", "Prefix to remove from paths. i.e. path to root")
)

type context struct {
	conditions      []string
	graphViz        bool
	labelConditions bool
	stripPrefix     string
}

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

func (l byTargetName) Len() int      { return len(l) }
func (l byTargetName) Swap(i, j int) { l[i], l[j] = l[j], l[i] }
func (l byTargetName) Less(i, j int) bool {
	return l[i].Name() < l[j].Name()
}

// byOrding orders license conditions by originating target name.
type byOrigin []compliance.LicenseCondition

func (l byOrigin) Len() int      { return len(l) }
func (l byOrigin) Swap(i, j int) { l[i], l[j] = l[j], l[i] }
func (l byOrigin) Less(i, j int) bool {
	if l[i].Origin().Name() == l[j].Origin().Name() {
		return l[i].Name() < l[j].Name()
	}
	return l[i].Origin().Name() < l[j].Origin().Name()
}

func main() {
	flag.Parse()

	// Must specify at least one root target.
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(2)
	}

	ctx := &context{
		conditions:      append([]string{}, *conditions...),
		graphViz:        *graphViz,
		labelConditions: *labelConditions,
		stripPrefix:     *stripPrefix,
	}
	err := dumpResolutions(ctx, os.Stdout, os.Stderr, flag.Args()...)
	if err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}

func dumpResolutions(ctx *context, stdout, stderr io.Writer, files ...string) error {
	if len(files) < 1 {
		return fmt.Errorf("no license metadata files requested")
	}

	// Read the license graph from the license metadata files (*.meta_lic).
	licenseGraph, err := compliance.ReadLicenseGraph(files)
	if err != nil {
		fmt.Fprintf(stderr, "Unable to read license metadata file(s) %q: %v\n", files, err)
		return err
	}
	if licenseGraph == nil {
		fmt.Fprintf(stderr, "No licenses\n")
		return fmt.Errorf("no licenses")
	}

	// resolutions will contain the requested set of resolutions.
	var resolutions compliance.ResolutionSet

	if len(ctx.conditions) == 0 {
		resolutions = compliance.ResolveBottomUpConditions(licenseGraph)
	} else {
		rlist := make([]compliance.ResolutionSet, 0, len(ctx.conditions))
		for _, c := range ctx.conditions {
			rlist = append(rlist, compliance.ResolveTopDownForCondition(licenseGraph, c))
		}
		if len(rlist) == 1 {
			resolutions = rlist[0]
		} else {
			resolutions = compliance.JoinResolutions(rlist...)
		}
	}

	nodes := make(map[string]string)
	n := 0

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

	makeNode := func(target compliance.TargetNode) {
		tName := target.Name()
		if _, ok := nodes[tName]; !ok {
			nodeName := fmt.Sprintf("n%d", n)
			nodes[tName] = nodeName
			fmt.Fprintf(stdout, "\t%s [label=\"%s\"];\n", nodeName, targetOut(target, "\\n"))
			n++
		}
	}

	outputEdge := func(tname, oname string, cnames []string) {
		if ctx.graphViz {
			// ... one edge per line labelled with \\n-separated annotations.
			tNode := nodes[tname]
			oNode := nodes[oname]
			fmt.Fprintf(stdout, "\t%s -> %s [label=\"%s\"];\n", tNode, oNode, strings.Join(cnames, "\\n"))
		} else {
			// ... one edge per line with names in a colon-separated tuple.
			fmt.Fprintf(stdout, "%s %s %s\n", tname, oname, strings.Join(cnames, ":"))
		}
	}

	outputSingleton := func(tname string) {
		if !ctx.graphViz {
			fmt.Fprintf(stdout, "%s\n", tname)
		}
	}

	// Sort the resolutions by targetname for repeatability/stability.
	targets := resolutions.AppliesTo()
	sort.Sort(byTargetName(targets))

	// if graphviz output, start the directed graph
	if ctx.graphViz {
		fmt.Fprintf(stdout, "strict digraph {\n\trankdir=LR;\n")
		for _, target := range targets {
			makeNode(target)
			for _, lc := range resolutions.Conditions(target).AsList() {
				makeNode(lc.Origin())
			}
		}
	}

	// Output the sorted targets.
	for _, target := range targets {
		var tname string
		if ctx.graphViz {
			tname = target.Name()
		} else {
			tname = targetOut(target, ":")
		}

		// Sort the conditions that apply to `target` for repeatability/stability.
		conditions := resolutions.Conditions(target).AsList()
		sort.Sort(byOrigin(conditions))

		poname := ""
		cnames := make([]string, 0, len(conditions))
		if len(conditions) == 0 {
			outputSingleton(tname)
		}
		// Output 1 line for each target+origin combination.
		for _, condition := range conditions {
			var oname string
			if ctx.graphViz {
				oname = condition.Origin().Name()
			} else {
				oname = targetOut(condition.Origin(), ":")
			}
			if poname != oname && poname != "" {
				outputEdge(tname, poname, cnames)
				cnames = cnames[:0]
			}
			poname = oname
			cnames = append(cnames, condition.Name())
		}
		if poname == "" {
			outputSingleton(tname)
		} else {
			outputEdge(tname, poname, cnames)
		}
	}
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
