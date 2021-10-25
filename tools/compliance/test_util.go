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

package compliance

import (
	"fmt"
	"io"
	"io/fs"
	"sort"
	"strings"
	"testing"
)

const (
	// AOSP starts a test metadata file for Android Apache-2.0 licensing.
	AOSP = `` +
`package_name: "Android"
license_kinds: "SPDX-license-identifier-Apache-2.0"
license_conditions: "notice"
`

	// GPL starts a test metadata file for GPL 2.0 licensing.
	GPL = `` +
`package_name: "Free Software"
license_kinds: "SPDX-license-identifier-GPL-2.0"
license_conditions: "restricted"
`

	// Classpath starts a test metadata file for GPL 2.0 with classpath exception licensing.
	Classpath = `` +
`package_name: "Free Software"
license_kinds: "SPDX-license-identifier-GPL-2.0-with-classpath-exception"
license_conditions: "restricted"
`

	// DependentModule starts a test metadata file for a module in the same package as `Classpath`.
	DependentModule = `` +
`package_name: "Free Software"
license_kinds: "SPDX-license-identifier-MIT"
license_conditions: "notice"
`

	// LGPL starts a test metadata file for a module with LGPL 2.0 licensing.
	LGPL = `` +
`package_name: "Free Library"
license_kinds: "SPDX-license-identifier-LGPL-2.0"
license_conditions: "restricted"
`

	// MPL starts a test metadata file for a module with MPL 2.0 reciprical licensing.
	MPL = `` +
`package_name: "Reciprocal"
license_kinds: "SPDX-license-identifier-MPL-2.0"
license_conditions: "reciprocal"
`

	// MIT starts a test metadata file for a module with generic notice (MIT) licensing.
	MIT = `` +
`package_name: "Android"
license_kinds: "SPDX-license-identifier-MIT"
license_conditions: "notice"
`

	// Proprietary starts a test metadata file for a module with proprietary licensing.
	Proprietary = `` +
`package_name: "Android"
license_kinds: "legacy_proprietary"
license_conditions: "proprietary"
`

	// ByException starts a test metadata file for a module with by_exception_only licensing.
	ByException = `` +
`package_name: "Special"
license_kinds: "legacy_by_exception_only"
license_conditions: "by_exception_only"
`

)

var (
	// meta maps test file names to metadata file content without dependencies.
	meta = map[string]string{
		"apacheBin.meta_lic": AOSP,
		"apacheLib.meta_lic": AOSP,
		"apacheContainer.meta_lic": AOSP + "is_container: true\n",
		"dependentModule.meta_lic": DependentModule,
		"gplWithClasspathException.meta_lic": Classpath,
		"gplBin.meta_lic": GPL,
		"gplLib.meta_lic": GPL,
		"gplContainer.meta_lic": GPL + "is_container: true\n",
		"lgplBin.meta_lic": LGPL,
		"lgplLib.meta_lic": LGPL,
		"mitBin.meta_lic": MIT,
		"mitLib.meta_lic": MIT,
		"mplBin.meta_lic": MPL,
		"mplLib.meta_lic": MPL,
		"proprietary.meta_lic": Proprietary,
		"by_exception.meta_lic": ByException,
	}
)



// toConditionList converts a test data map of condition name to origin names into a ConditionList.
func toConditionList(lg *licenseGraphImp, conditions map[string][]string) ConditionList {
	cl := make(ConditionList, 0)
	for name, origins := range conditions {
		for _, origin := range origins {
			cl = append(cl, licenseConditionImp{name, newTestNode(lg, origin).(targetNodeImp)})
		}
	}
	return cl
}


// newTestNode constructs a test node in the license graph.
func newTestNode(lg *licenseGraphImp, targetName string) TargetNode {
	if _, ok := lg.targets[targetName]; !ok {
		lg.targets[targetName] = &targetNode{name: targetName}
	}
	return targetNodeImp{lg, targetName}
}


// testFS implements a test file system (fs.FS) simulated by a map from filename to []byte content.
type testFS map[string][]byte

// Open implements fs.FS.Open() to open a file based on the filename.
func (fs *testFS) Open(name string) (fs.File, error) {
	if _, ok := (*fs)[name]; !ok {
		return nil, fmt.Errorf("unknown file %q", name)
	}
	return &testFile{fs, name, 0}, nil
}

// testFile implements a test file (fs.File) based on testFS above.
type testFile struct {
	fs   *testFS
	name string
	posn int
}

// Stat not implemented to obviate implementing fs.FileInfo.
func (f *testFile) Stat() (fs.FileInfo, error) {
	return nil, fmt.Errorf("unimplemented")
}

// Read copies bytes from the testFS map.
func (f *testFile) Read(b []byte) (int, error) {
	if f.posn < 0 {
		return 0, fmt.Errorf("file not open: %q", f.name)
	}
	if f.posn >= len((*f.fs)[f.name]) {
		return 0, io.EOF
	}
	n := copy(b, (*f.fs)[f.name][f.posn:])
	f.posn += n
	return n, nil
}

// Close marks the testFile as no longer in use.
func (f *testFile) Close() error {
	if f.posn < 0 {
		return fmt.Errorf("file already closed: %q", f.name)
	}
	f.posn = -1
	return nil
}


// edge describes test data edges to define test graphs.
type edge struct {
	target, dep string
}

func (e edge) String() string {
	return e.target + " -> " + e.dep
}

// byEdge orders edges by target then dep name then annotations.
type byEdge []edge

func (l byEdge) Len() int      { return len(l) }
func (l byEdge) Swap(i, j int) { l[i], l[j] = l[j], l[i] }
func (l byEdge) Less(i, j int) bool {
	if l[i].target == l[j].target {
		return l[i].dep < l[j].dep
	}
	return l[i].target < l[j].target
}


// annotated describes annotated test data edges to define test graphs.
type annotated struct {
	target, dep string
	annotations []string
}

func (e annotated) String() string {
	if e.annotations != nil {
		return e.target + " -> " + e.dep + " [" + strings.Join(e.annotations, ", ") + "]"
	}
	return e.target + " -> " + e.dep
}

func (e annotated) IsEqualTo(other annotated) bool {
	if e.target != other.target {
		return false
	}
	if e.dep != other.dep {
		return false
	}
        if len(e.annotations) != len(other.annotations) {
		return false
	}
	a1 := append([]string{}, e.annotations...)
	a2 := append([]string{}, other.annotations...)
	for i := 0; i < len(a1); i++ {
		if a1[i] != a2[i] {
			return false
		}
	}
	return true
}

// toGraph converts a list of roots and a list of annotated edges into a test license graph.
func toGraph(stderr io.Writer, roots []string, edges []annotated) (LicenseGraph, error) {
	deps := make(map[string][]annotated)
	for _, root := range roots {
		deps[root] = []annotated{}
	}
	for _, edge := range edges {
		if prev, ok := deps[edge.target]; ok {
			deps[edge.target] = append(prev, edge)
		} else {
			deps[edge.target] = []annotated{edge}
		}
		if _, ok := deps[edge.dep]; !ok {
			deps[edge.dep] = []annotated{}
		}
	}
	fs := make(testFS)
	for file, edges := range deps {
		body := meta[file]
		for _, edge := range edges {
			body += fmt.Sprintf("deps: {\n  file: %q\n", edge.dep)
			for _, ann := range edge.annotations {
				body += fmt.Sprintf("  annotations: %q\n", ann)
			}
			body += "}\n"
		}
		fs[file] = []byte(body)
	}

	return ReadLicenseGraph(&fs, stderr, roots)
}


// byAnnotatedEdge orders edges by target then dep name then annotations.
type byAnnotatedEdge []annotated

func (l byAnnotatedEdge) Len() int      { return len(l) }
func (l byAnnotatedEdge) Swap(i, j int) { l[i], l[j] = l[j], l[i] }
func (l byAnnotatedEdge) Less(i, j int) bool {
	if l[i].target == l[j].target {
		if l[i].dep == l[j].dep {
			ai := append([]string{}, l[i].annotations...)
			aj := append([]string{}, l[j].annotations...)
			sort.Strings(ai)
			sort.Strings(aj)
			for k := 0; k < len(ai) && k < len(aj); k++ {
				if ai[k] == aj[k] {
					continue
				}
				return ai[k] < aj[k]
			}
			return len(ai) < len(aj)
		}
		return l[i].dep < l[j].dep
	}
	return l[i].target < l[j].target
}

// res describes test data resolutions to define test resolution sets.
type res struct {
	appliesTo, origin, condition string
}

// toResolutionSet converts a list of res test data into a test resolution set.
func toResolutionSet(lg *licenseGraphImp, data []res) ResolutionSet {
	rmap := make(map[string]*licenseConditionSetImp)
	for _, r := range data {
		_ = newTestNode(lg, r.appliesTo)
		onode := newTestNode(lg, r.origin).(targetNodeImp)
		if _, ok := rmap[r.appliesTo]; !ok {
			rmap[r.appliesTo] = newLicenseConditionSet(&onode)
		}
		rmap[r.appliesTo].addAll(r.origin, r.condition)
	}
	return &resolutionSetImp{lg, rmap}
}

// byOriginName orders license conditions by originating target then by condition name.
type byOriginName []LicenseCondition

func (l byOriginName) Len() int      { return len(l) }
func (l byOriginName) Swap(i, j int) { l[i], l[j] = l[j], l[i] }
func (l byOriginName) Less(i, j int) bool {
	return l[i].Origin().Name() < l[j].Origin().Name() || (l[i].Origin().Name() == l[j].Origin().Name() && l[i].Name() < l[j].Name())
}

// checkSame compares an actual resolution set to an expected resolution set for a test.
func checkSame(rsActual, rsExpected ResolutionSet, t *testing.T) {
	for _, target := range rsExpected.AppliesTo() {
		if !rsActual.AppliesToTarget(target) {
			t.Errorf("unexpected missing target: got AppliesToTarget(%q) is false, want true", target.Name())
			continue
		}
		expectedConditions := rsExpected.Conditions(target).AsList()
		actualConditions := rsActual.Conditions(target).AsList()
		sort.Sort(byOriginName(expectedConditions))
		sort.Sort(byOriginName(actualConditions))
		if len(expectedConditions) != len(actualConditions) {
			t.Errorf("unexpected number of conditions apply to %q: got %v with %d elements, want %v with %d elements",
				target.Name(), actualConditions, len(actualConditions), expectedConditions, len(expectedConditions))
			continue
		}
		for i := 0; i < len(expectedConditions); i++ {
			if expectedConditions[i] != actualConditions[i] {
				t.Errorf("unexpected condition applies to %q at index %d: got %s in %v, want %s in %v",
					target.Name(), i, actualConditions[i].(licenseConditionImp).asString(":"), actualConditions,
					expectedConditions[i].(licenseConditionImp).asString(":"), expectedConditions)
			}
		}
	}
	for i, target := range rsActual.AppliesTo() {
		if !rsExpected.AppliesToTarget(target) {
			t.Errorf("unexpected target: got %q element %d in AppliesTo() %v with %d elements, want %v with %d elements",
				target.Name(), i, rsActual.AppliesTo(), len(rsActual.AppliesTo()), rsExpected.AppliesTo(), len(rsExpected.AppliesTo()))
		}
	}
}
