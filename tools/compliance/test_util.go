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


// byTargetName orders TargetNode elements by name.
type byTargetName []TargetNode

func (l byTargetName) Len() int      { return len(l) }
func (l byTargetName) Swap(i, j int) { l[i], l[j] = l[j], l[i] }
func (l byTargetName) Less(i, j int) bool {
	return l[i].Name() < l[j].Name()
}

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


// res describes test data resolutions to define test resolution sets.
type res struct {
	attachesTo, actsOn, origin, condition string
}

// toResolutionSet converts a list of res test data into a test resolution set.
func toResolutionSet(lg *licenseGraphImp, data []res) ResolutionSet {
	rmap := make(map[string]actionSet)
	for _, r := range data {
		_ = newTestNode(lg, r.attachesTo)
		_ = newTestNode(lg, r.actsOn)
		onode := newTestNode(lg, r.origin).(targetNodeImp)
		if _, ok := rmap[r.attachesTo]; !ok {
			rmap[r.attachesTo] = make(actionSet)
		}
		if _, ok := rmap[r.attachesTo][r.actsOn]; !ok {
			rmap[r.attachesTo][r.actsOn] = newLicenseConditionSet(&onode)
		}
		rmap[r.attachesTo][r.actsOn].addAll(r.origin, r.condition)
	}
	return &resolutionSetImp{lg, rmap}
}

// checkSameActions compares an actual action set to an expected action set for a test.
func checkSameActions(lg *licenseGraphImp, asActual, asExpected actionSet, t *testing.T) {
	rsActual := resolutionSetImp{lg, make(map[string]actionSet)}
	rsExpected := resolutionSetImp{lg, make(map[string]actionSet)}
	_ = newTestNode(lg, "test")
	rsActual.resolutions["test"] = asActual
	rsExpected.resolutions["test"] = asExpected
	checkSame(&rsActual, &rsExpected, t)
}

// checkSame compares an actual resolution set to an expected resolution set for a test.
func checkSame(rsActual, rsExpected ResolutionSet, t *testing.T) {
	expectedTargets := rsExpected.AttachesTo()
	sort.Sort(byTargetName(expectedTargets))
	for _, target := range expectedTargets {
		if !rsActual.AttachesToTarget(target) {
			t.Errorf("unexpected missing target: got AttachesToTarget(%q) is false in %s, want true in %s", target.Name(), rsActual, rsExpected)
			continue
		}
		expectedRl := ResolutionList(rsExpected.Resolutions(target))
		sort.Sort(expectedRl)
		actualRl := ResolutionList(rsActual.Resolutions(target))
		sort.Sort(actualRl)
		if len(expectedRl) != len(actualRl) {
			t.Errorf("unexpected number of resolutions attach to %q: got %s with %d elements, want %s with %d elements",
				target.Name(), actualRl, len(actualRl), expectedRl, len(expectedRl))
			continue
		}
		for i := 0; i < len(expectedRl); i++ {
			if expectedRl[i].AttachesTo().Name() != actualRl[i].AttachesTo().Name() || expectedRl[i].ActsOn().Name() != actualRl[i].ActsOn().Name() {
				t.Errorf("unexpected resolution attaches to %q at index %d: got %s, want %s",
					target.Name(), i, actualRl[i].(resolutionImp).asString(),
					expectedRl[i].(resolutionImp).asString())
				continue
			}
			expectedConditions := expectedRl[i].Resolves().AsList()
			actualConditions := actualRl[i].Resolves().AsList()
			sort.Sort(expectedConditions)
			sort.Sort(actualConditions)
			if len(expectedConditions) != len(actualConditions) {
				t.Errorf("unexpected number of conditions apply to %q acting on %q: got %v with %d elements, want %v with %d elements",
					target.Name(), expectedRl[i].ActsOn().Name(),
					actualConditions, len(actualConditions),
					expectedConditions, len(expectedConditions))
				continue
			}
			for j := 0; j < len(expectedConditions); j++ {
				if expectedConditions[j] != actualConditions[j] {
					t.Errorf("unexpected condition attached to %q acting on %q at index %d: got %s at index %d in %v, want %s in %v",
						target.Name(), expectedRl[i].ActsOn().Name(), i,
						actualConditions[j].(licenseConditionImp).asString(":"), j, actualConditions,
						expectedConditions[j].(licenseConditionImp).asString(":"), expectedConditions)
				}
			}
		}

	}
	actualTargets := rsActual.AttachesTo()
	sort.Sort(byTargetName(actualTargets))
	for i, target := range actualTargets {
		if !rsExpected.AttachesToTarget(target) {
			t.Errorf("unexpected target: got %q element %d in AttachesTo() %v with %d elements in %s, want %v with %d elements in %s",
				target.Name(), i, actualTargets, len(actualTargets), rsActual, expectedTargets, len(expectedTargets), rsExpected)
		}
	}
}
