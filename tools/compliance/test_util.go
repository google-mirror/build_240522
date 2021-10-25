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
	lg.targets[targetName] = &targetNode{name: targetName}
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

// Stat not implemented to obviate implmenting fs.FileInfo.
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
				t.Errorf("unexpected condition applies to %q at index %d: got %s, want %s",
					target.Name(), i, actualConditions[i].(licenseConditionImp).asString(":"),
					expectedConditions[i].(licenseConditionImp).asString(":"))
			}
		}
	}
	for _, target := range rsActual.AppliesTo() {
		if !rsExpected.AppliesToTarget(target) {
			t.Errorf("unexpected target: %q in AppliesTo(), want not in AppliesTo()", target.Name())
		}
	}
}
