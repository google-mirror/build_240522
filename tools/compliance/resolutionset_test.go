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
	"testing"
)

var (
	// bottomUp describes the bottom-up resolve of a hypothetical graph
	// the graph has a container image, a couple binaries, and a couple
	// libraries. bin1 statically links lib1 and dynamically links lib2;
	// bin2 dynamically links lib1 and statically links lib2.
	// binc represents a compiler or other toolchain binary used for
	// building the other binaries.
	bottomUp = []res{
		{"image", "image", "image", "notice"},
		{"image", "image", "bin2", "restricted"},
		{"image", "bin1", "bin1", "reciprocal"},
		{"image", "bin2", "bin2", "restricted"},
		{"image", "lib1", "lib1", "notice"},
		{"image", "lib2", "lib2", "notice"},
		{"binc", "binc", "binc", "proprietary"},
		{"bin1", "bin1", "bin1", "reciprocal"},
		{"bin1", "lib1", "lib1", "notice"},
		{"bin2", "bin2", "bin2", "restricted"},
		{"bin2", "lib2", "lib2", "notice"},
		{"lib1", "lib1", "lib1", "notice"},
		{"lib2", "lib2", "lib2", "notice"},
	}

	// notice describes bottomUp after a top-down notice resolve.
	notice = []res{
		{"image", "image", "image", "notice"},
		{"image", "image", "bin2", "restricted"},
		{"image", "bin1", "bin1", "reciprocal"},
		{"image", "bin2", "bin2", "restricted"},
		{"image", "lib1", "lib1", "notice"},
		{"image", "lib2", "bin2", "restricted"},
		{"image", "lib2", "lib2", "notice"},
		{"bin1", "bin1", "bin1", "reciprocal"},
		{"bin1", "lib1", "lib1", "notice"},
		{"bin2", "bin2", "bin2", "restricted"},
		{"bin2", "lib2", "bin2", "restricted"},
		{"bin2", "lib2", "lib2", "notice"},
		{"lib1", "lib1", "lib1", "notice"},
		{"lib2", "lib2", "lib2", "notice"},
	}

	// share describes bottomUp after a top-down share resolve.
	share = []res{
		{"image", "image", "bin2", "restricted"},
		{"image", "bin1", "bin1", "reciprocal"},
		{"image", "bin2", "bin2", "restricted"},
		{"image", "lib2", "bin2", "restricted"},
		{"bin1", "bin1", "bin1", "reciprocal"},
		{"bin2", "bin2", "bin2", "restricted"},
		{"bin2", "lib2", "bin2", "restricted"},
	}

	// proprietary describes bottomUp after a top-down proprietary resolve.
	// Note that the proprietary binc is not reachable through the toolchain
	// dependency.
	proprietary = []res{}
)

func TestResolutionSet_JoinResolutionSets(t *testing.T) {
	lg := newLicenseGraph()

	rsNotice := toResolutionSet(lg, notice)
	rsShare := toResolutionSet(lg, share)
	rsExpected := toResolutionSet(lg, append(notice, share...))

	rsActual := JoinResolutionSets(rsNotice, rsShare)
	t.Logf("Joining %s and %s yields %s, want %s", rsNotice.String(), rsShare.String(), rsActual.String(), rsExpected.String())
	checkSame(rsActual, rsExpected, t)
}

func TestResolutionSet_JoinResolutionsEmpty(t *testing.T) {
	lg := newLicenseGraph()

	rsShare := toResolutionSet(lg, share)
	rsProprietary := toResolutionSet(lg, proprietary)
	rsExpected := toResolutionSet(lg, append(share, proprietary...))

	rsActual := JoinResolutionSets(rsShare, rsProprietary)
	t.Logf("Joining %s and %s yields %s, want %s", rsShare.String(), rsProprietary.String(), rsActual.String(), rsExpected.String())
	checkSame(rsActual, rsExpected, t)
}

func TestResolutionSet_AttachedToTarget(t *testing.T) {
	lg := newLicenseGraph()

	rsShare := toResolutionSet(lg, share)

	t.Logf("checking resolution set %s", rsShare.String())

	if rsShare.AttachesToTarget(newTestNode(lg, "binc")) {
		t.Errorf("actual.AttachedToTarget(\"binc\"): got true, want false")
	}
	if !rsShare.AttachesToTarget(newTestNode(lg, "image")) {
		t.Errorf("actual.AttachedToTarget(\"image\"): got false want true")
	}
}

func TestResolutionSet_AnyMatchingAttachToTarget(t *testing.T) {
	lg := newLicenseGraph()

	rs := toResolutionSet(lg, bottomUp)

	t.Logf("checking resolution set %s", rs.String())

	pandp := LicenseConditionSet(PermissiveCondition | ProprietaryCondition)
	pandn := LicenseConditionSet(PermissiveCondition | NoticeCondition)
	p := LicenseConditionSet(ProprietaryCondition)
	r := LicenseConditionSet(RestrictedCondition)

	if rs.AnyMatchingAttachToTarget(newTestNode(lg, "image"), pandp, p) {
		t.Errorf("actual.AnyMatchingAttachToTarget(\"image\", \"proprietary\", \"permissive\") in %s: want false, got true", rs.String())
	}
	if !rs.AnyMatchingAttachToTarget(newTestNode(lg, "binc"), p) {
		t.Errorf("actual.AnyMatchingAttachToTarget(\"binc\", \"proprietary\"): want true, got false")
	}
	if !rs.AnyMatchingAttachToTarget(newTestNode(lg, "image"), pandn) {
		t.Errorf("actual.AnyMatchingAttachToTarget(\"image\", \"permissive\", \"notice\"): want true, got false")
	}
	if !rs.AnyMatchingAttachToTarget(newTestNode(lg, "image"), r, pandn) {
		t.Errorf("actual.AnyMatchingAttachToTarget(\"image\", \"restricted\", \"notice\"): want true, got false")
	}
	if !rs.AnyMatchingAttachToTarget(newTestNode(lg, "image"), r, p) {
		t.Errorf("actual.AnyMatchingAttachToTarget(\"image\", \"restricted\", \"proprietary\"): want true, got false")
	}
}

func TestResolutionSet_AttachesToTarget(t *testing.T) {
	lg := newLicenseGraph()

	rsShare := toResolutionSet(lg, share)

	t.Logf("checking resolution set %s", rsShare.String())

	if rsShare.AttachesToTarget(newTestNode(lg, "binc")) {
		t.Errorf("actual.AttachesToTarget(\"binc\"): got true, want false")
	}
	if !rsShare.AttachesToTarget(newTestNode(lg, "image")) {
		t.Errorf("actual.AttachesToTarget(\"image\"): got false want true")
	}
}
