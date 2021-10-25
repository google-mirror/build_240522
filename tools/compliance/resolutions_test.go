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
	"sort"
	"testing"
)

type res struct {
	appliesTo, origin, condition string
}

var (
	// bottomUp describes the bottom-up resolve of a hypothetical graph
	// the graph has a container image, a couple binaries, and a couple
	// libraries. bin1 statically links lib1 and dynamically links lib2;
	// bin2 dynamically links lib1 and statically links lib2.
	// binc represents a compiler or other toolchain binary used for
	// building the other binaries.
	bottomUp = []res{
		{"image", "image", "notice"},
		{"image", "bin1", "reciprocal"},
		{"image", "bin2", "restricted"},
		{"image", "lib1", "notice"},
		{"image", "lib2", "notice"},
		{"binc", "binc", "proprietary"},
		{"bin1", "bin1", "reciprocal"},
		{"bin1", "lib1", "notice"},
		{"bin2", "bin2", "restricted"},
		{"bin2", "lib2", "notice"},
		{"lib1", "lib1", "notice"},
		{"lib2", "lib2", "notice"},
	}

	// notice describes bottomUp after a top-down notice resolve.
	notice = []res{
		{"image", "image", "notice"},
		{"image", "bin1", "notice"},
		{"image", "bin2", "notice"},
		{"image", "lib1", "notice"},
		{"image", "lib2", "notice"},
		{"bin1", "bin1", "notice"},
		{"bin1", "lib1", "notice"},
		{"bin2", "bin2", "notice"},
		{"bin2", "lib2", "notice"},
		{"lib1", "lib1", "notice"},
		{"lib2", "lib2", "notice"},
	}

	// share describes bottomUp after a top-down share resolve.
	share = []res{
		{"image", "bin1", "reciprocal"},
		{"image", "bin2", "restricted"},
		{"bin1", "bin1", "reciprocal"},
		{"bin2", "bin2", "restricted"},
		{"lib2", "bin1", "restricted"},
	}

	// proprietary describes bottomUp after a top-down proprietary resolve.
	// Note that the proprietary binc is not reachable through the toolchain
	// dependency.
	proprietary = []res{}
)

func testNode(lg *licenseGraphImp, targetName string) targetNodeImp {
	lg.targets[targetName] = &targetNode{name: targetName}
	return targetNodeImp{lg, targetName}
}

func toResolutionSet(lg *licenseGraphImp, data []res) ResolutionSet {
	rmap := make(map[string]*licenseConditionSetImp)
	for _, r := range data {
		_ = testNode(lg, r.appliesTo)
		onode := testNode(lg, r.origin)
		if _, ok := rmap[r.appliesTo]; !ok {
			rmap[r.appliesTo] = newLicenseConditionSet(&onode)
		}
		rmap[r.appliesTo].addAll(r.origin, r.condition)
	}
	return &resolutionSetImp{lg, rmap}
}

// byOriginName orders license conditions by originating target then condition name.
type byOriginName []LicenseCondition

func (l byOriginName) Len() int      { return len(l) }
func (l byOriginName) Swap(i, j int) { l[i], l[j] = l[j], l[i] }
func (l byOriginName) Less(i, j int) bool {
	return l[i].Origin().Name() < l[j].Origin().Name() || (l[i].Origin().Name() == l[j].Origin().Name() && l[i].Name() < l[j].Name())
}

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

func TestResolutionSet_JoinResolutions(t *testing.T) {
	lg := newLicenseGraphImp()

	rsNotice := toResolutionSet(lg, notice)
	rsShare := toResolutionSet(lg, share)
	rsExpected := toResolutionSet(lg, append(notice, share...))

	rsActual := JoinResolutions(rsNotice, rsShare)
	checkSame(rsActual, rsExpected, t)
}

func TestResolutionSet_JoinResolutionsEmpty(t *testing.T) {
	lg := newLicenseGraphImp()

	rsShare := toResolutionSet(lg, share)
	rsProprietary := toResolutionSet(lg, proprietary)
	rsExpected := toResolutionSet(lg, append(share, proprietary...))

	rsActual := JoinResolutions(rsShare, rsProprietary)
	checkSame(rsActual, rsExpected, t)
}

func TestResolutionSet_Origins(t *testing.T) {
	lg := newLicenseGraphImp()

	rsShare := toResolutionSet(lg, share)

	origins := make([]string, 0)
	for _, target := range rsShare.Origins() {
		origins = append(origins, target.Name())
	}
	sort.Strings(origins)
	if len(origins) != 2 {
		t.Errorf("unexpected number of origins: got %v with %d elements, want [\"bin1\", \"bin2\"] with 2 elements", origins, len(origins))
	}
	if origins[0] != "bin1" {
		t.Errorf("unexpected origin at element 0: got %s, want \"bin1\"", origins[0])
	}
	if origins[1] != "bin2" {
		t.Errorf("unexpected origin at element 0: got %s, want \"bin2\"", origins[0])
	}
}

func TestResolutionSet_AppliesToTarget(t *testing.T) {
	lg := newLicenseGraphImp()

	rsShare := toResolutionSet(lg, share)

	if rsShare.AppliesToTarget(testNode(lg, "binc")) {
		t.Errorf("unexpected AppliesToTarget(\"binc\"): got true, want false")
	}
	if !rsShare.AppliesToTarget(testNode(lg, "image")) {
		t.Errorf("unexpected AppliesToTarget(\"image\"): got false want true")
	}
}

func TestResolutionSet_AnyByNameApplyToTarget(t *testing.T) {
	lg := newLicenseGraphImp()

	rs := toResolutionSet(lg, bottomUp)

	pandp := ConditionNames{"permissive", "proprietary"}
	pandn := ConditionNames{"permissive", "notice"}
	p := ConditionNames{"proprietary"}
	r := ConditionNames{"restricted"}

	if rs.AnyByNameApplyToTarget(testNode(lg, "image"), pandp, p) {
		t.Errorf("unexpected AnyByNameApplyToTarget(\"image\", \"proprietary\", \"permissive\"): want false, got true")
	}
	if !rs.AnyByNameApplyToTarget(testNode(lg, "binc"), p) {
		t.Errorf("unexpected AnyByNameApplyToTarget(\"binc\", \"proprietary\"): want true, got false")
	}
	if !rs.AnyByNameApplyToTarget(testNode(lg, "image"), pandn) {
		t.Errorf("unexpected AnyByNameApplyToTarget(\"image\", \"permissive\", \"notice\"): want true, got false")
	}
	if !rs.AnyByNameApplyToTarget(testNode(lg, "image"), r, pandn) {
		t.Errorf("unexpected AnyByNameApplyToTarget(\"image\", \"restricted\", \"notice\"): want true, got false")
	}
	if !rs.AnyByNameApplyToTarget(testNode(lg, "image"), r, p) {
		t.Errorf("unexpected AnyByNameApplyToTarget(\"image\", \"restricted\", \"proprietary\"): want true, got false")
	}
}

func TestResolutionSet_AllByNameApplyToTarget(t *testing.T) {
	lg := newLicenseGraphImp()

	rs := toResolutionSet(lg, bottomUp)

	pandp := ConditionNames{"permissive", "proprietary"}
	pandn := ConditionNames{"permissive", "notice"}
	p := ConditionNames{"proprietary"}
	r := ConditionNames{"restricted"}

	if rs.AllByNameApplyToTarget(testNode(lg, "image"), pandp, p) {
		t.Errorf("unexpected AllByNameApplyToTarget(\"image\", \"proprietary\", \"permissive\"): want false, got true")
	}
	if !rs.AllByNameApplyToTarget(testNode(lg, "binc"), p) {
		t.Errorf("unexpected AllByNameApplyToTarget(\"binc\", \"proprietary\"): want true, got false")
	}
	if !rs.AllByNameApplyToTarget(testNode(lg, "image"), pandn) {
		t.Errorf("unexpected AllByNameApplyToTarget(\"image\", \"notice\"): want true, got false")
	}
	if !rs.AllByNameApplyToTarget(testNode(lg, "image"), r, pandn) {
		t.Errorf("unexpected AllByNameApplyToTarget(\"image\", \"restricted\", \"notice\"): want true, got false")
	}
	if rs.AllByNameApplyToTarget(testNode(lg, "image"), r, p) {
		t.Errorf("unexpected AllByNameApplyToTarget(\"image\", \"restricted\", \"proprietary\"): want false, got true")
	}
}

func TestResolutionSet_hasTarget(t *testing.T) {
	lg := newLicenseGraphImp()

	rsShare := toResolutionSet(lg, share).(*resolutionSetImp)

	if rsShare.hasTarget("binc") {
		t.Errorf("unexpected hasTarget(\"binc\"): got true, want false")
	}
	if !rsShare.hasTarget("image") {
		t.Errorf("unexpected AppliesToTarget(\"image\"): got false want true")
	}
}

func TestResolutionSet_hasAnyByName(t *testing.T) {
	lg := newLicenseGraphImp()

	rs := toResolutionSet(lg, bottomUp).(*resolutionSetImp)

	pandp := ConditionNames{"permissive", "proprietary"}
	pandn := ConditionNames{"permissive", "notice"}
	p := ConditionNames{"proprietary"}
	r := ConditionNames{"restricted"}

	if rs.hasAnyByName("image", pandp, p) {
		t.Errorf("unexpected hasAnyByName(\"image\", \"proprietary\", \"permissive\"): want false, got true")
	}
	if !rs.hasAnyByName("binc", p) {
		t.Errorf("unexpected hasAnyByName(\"binc\", \"proprietary\"): want true, got false")
	}
	if !rs.hasAnyByName("image", pandn) {
		t.Errorf("unexpected hasAnyByName(\"image\", \"permissive\", \"notice\"): want true, got false")
	}
	if !rs.hasAnyByName("image", r, pandn) {
		t.Errorf("unexpected hasAnyByName(\"image\", \"restricted\", \"notice\"): want true, got false")
	}
	if !rs.hasAnyByName("image", r, p) {
		t.Errorf("unexpected AnyByNameApplyToTarget(\"image\", \"restricted\", \"proprietary\"): want true, got false")
	}
}

func TestResolutionSet_hasAllByName(t *testing.T) {
	lg := newLicenseGraphImp()

	rs := toResolutionSet(lg, bottomUp).(*resolutionSetImp)

	pandp := ConditionNames{"permissive", "proprietary"}
	pandn := ConditionNames{"permissive", "notice"}
	p := ConditionNames{"proprietary"}
	r := ConditionNames{"restricted"}

	if rs.hasAllByName("image", pandp, p) {
		t.Errorf("unexpected hasAllByName(\"image\", \"proprietary\", \"permissive\"): want false, got true")
	}
	if !rs.hasAllByName("binc", p) {
		t.Errorf("unexpected hasAllByName(\"binc\", \"proprietary\"): want true, got false")
	}
	if !rs.hasAllByName("image", pandn) {
		t.Errorf("unexpected hasAllByName(\"image\", \"notice\"): want true, got false")
	}
	if !rs.hasAllByName("image", r, pandn) {
		t.Errorf("unexpected hasAllByName(\"image\", \"restricted\", \"notice\"): want true, got false")
	}
	if rs.hasAllByName("image", r, p) {
		t.Errorf("unexpected hasAllByName(\"image\", \"restricted\", \"proprietary\"): want false, got true")
	}
}
