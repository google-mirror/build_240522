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
	"strings"
	"testing"
)

func TestResolutionAction(t *testing.T) {
	tna := &TargetNode{name: "testNodeA"}
	tnb := &TargetNode{name: "testNodeB"}
	r1 := resolutionAction{tna, LicenseConditionSet(NoticeCondition)}
	r2 := resolutionAction{tnb, ImpliesShared}
	r3 := resolutionAction{tnb, ImpliesShared}
	r4 := resolutionAction{tnb, LicenseConditionSet(NoticeCondition)}

	t.Logf("A:Notice=%s", r1.String())
	t.Logf("B:Shared(1st instance)=%s", r2.String())
	t.Logf("B:Shared(2nd instance)=%s", r3.String())
	t.Logf("B:Notice=%s", r4.String())

	if !r1.IsEqual(r1) {
		t.Errorf("A:Notice.IsEqual(self), got false, want true")
	}
	if r1.IsEqual(r2) {
		t.Errorf("A:Notice.IsEqual(B:Shared), got true, want false")
	}
        if r1.IsEqual(r4) {
		t.Errorf("A:Notice.IsEqual(B:Notice), got true, want false")
	}
	if r2.IsEqual(r1) {
		t.Errorf("B:Shared(1st).IsEqual(A:Notice), got true, want false")
	}
	if !r2.IsEqual(r3) {
		t.Errorf("B:Shared(1st).IsEqual(B:Shared(2nd)), got false, want true")
	}
	if r2.IsEqual(r4) {
		t.Errorf("B:Shared(1st).IsEqual(B:Notice), got true, want false")
	}
	if r3.IsEqual(r1) {
		t.Errorf("B:Shared(2nd).IsEqual(A:Notice), got true, want false")
	}
	if !r3.IsEqual(r2) {
		t.Errorf("B:Shared(2nd).IsEqual(B:Shared(1st)), got false, want true")
	}
	if r3.IsEqual(r4) {
		t.Errorf("B:Shared(2nd).IsEqual(B:Notice), got true, want false")
	}
	if r4.IsEqual(r1) {
		t.Errorf("B:Notice.IsEqual(A:Notice), got true, want false")
	}
	if r4.IsEqual(r2) {
		t.Errorf("B:Notice.IsEqual(B:Shared), got true, want false")
	}
	if !r4.IsEqual(r4) {
		t.Errorf("B:Notce.IsEqual(self), got fales, want true")
	}

	if r1.IsLess(r1) {
		t.Errorf("A:Notice.IsLess(self), got true, want false")
	}
	if !r1.IsLess(r2) {
		t.Errorf("A:Notice.IsLess(B:Shared), got false, want true")
	}
	if !r1.IsLess(r4) {
		t.Errorf("A:Notice.IsLess(B:Notice), got false, want true")
	}
	if r2.IsLess(r1) {
		t.Errorf("B:Shared.IsLess(A:Notice), got true, want false")
	}
	if r2.IsLess(r2) {
		t.Errorf("B:Shared.IsLess(self), got true, want false")
	}
	if r2.IsLess(r3) {
		t.Errorf("B:Shared.IsLess(B:Shared), got true, want false")
	}
	if r2.IsLess(r4) {
		t.Errorf("B:Shared.IsLess(B:Notice), got true, want false")
	}
	if r4.IsLess(r1) {
		t.Errorf("B:Notice.IsLess(A:Notice), got true, want false")
	}
	if !r4.IsLess(r2) {
		t.Errorf("B:Notice.IsLess(B:Shared), got false, want true")
	}
	if r4.IsLess(r4) {
		t.Errorf("B:Notice.IsLess(self), got true, want false")
	}

	if r1.String() != "testNodeA{notice}" {
		t.Errorf("A:Notice.String(), got %q, want %q", r1.String(), "testNodeA{notice}")
	}
	if !strings.HasPrefix(r2.String(), "testNodeB{reciprocal, restricted") {
		t.Errorf("B:Shared.String(), got %q, want string starting with %q", r1.String(), "testNodeB{reciprocal, restricted")
	}
	if r4.String() != "testNodeB{notice}" {
		t.Errorf("B:Notice.String(), got %q, want %q", r4.String(), "testNodeB{notice}")
	}
}

func TestActionSet(t *testing.T) {
	tests := []struct {
		name        string
		actions     []act
		plus        *[]act
		matchingAny map[string][]act
		expected    []act
	}{
		{
			name:       "empty",
			actions: []act{},
			plus:       &[]act{},
			matchingAny: map[string][]act{
				"notice":     []act{},
				"restricted": []act{},
				"restricted|reciprocal": []act{},
			},
			expected:   []act{},
		},
		{
			name: "noticeonly",
			actions: []act{{"image", "bin1", "notice"}},
			matchingAny: map[string][]act{
				"notice":     []act{{"image", "bin1", "notice"}},
				"notice|proprietary":     []act{{"image", "bin1", "notice"}},
				"restricted": []act{},
			},
			expected: []act{{"image", "bin1", "notice"}},
		},
		{
			name: "allnoticeonly",
			actions: []act{{"image", "bin1", "notice"}},
			plus: &[]act{{"image", "bin2", "notice"}},
			matchingAny: map[string][]act{
				"notice":     []act{{"image", "bin1", "notice"}, {"image", "bin2", "notice"}},
				"notice|proprietary":     []act{{"image", "bin1", "notice"}, {"image", "bin2", "notice"}},
				"restricted": []act{},
			},
			expected: []act{{"image", "bin1", "notice"}, {"image", "bin2", "notice"}},
		},
		{
			name: "selfnoticeonly",
			actions: []act{{"image", "bin1", "notice"}},
			plus: &[]act{{"image", "bin1", "notice"}},
			matchingAny: map[string][]act{
				"notice":     []act{{"image", "bin1", "notice"}},
				"notice|proprietary":     []act{{"image", "bin1", "notice"}},
				"restricted": []act{},
			},
			expected: []act{{"image", "bin1", "notice"}},
		},
		{
			name: "emptyplusnotice",
			actions: []act{},
			plus: &[]act{{"bin1", "bin1", "notice"}},
			matchingAny: map[string][]act{
				"notice":     []act{{"bin1", "bin1", "notice"}},
				"notice|proprietary":     []act{{"bin1", "bin1", "notice"}},
				"restricted": []act{},
			},
			expected: []act{{"bin1", "bin1", "notice"}},
		},
		{
			name: "everything",
			actions: []act{
				{"image", "lib1", "unencumbered"},
				{"image", "lib2", "permissive"},
				{"image", "bin1", "notice"},
				{"image", "bin2", "reciprocal"},
				{"image", "bin3", "restricted"},
				{"image", "bin4", "proprietary"},
			},
			plus: &[]act{
				{"bin3", "lib3", "restricted_with_classpath_exception"},
				{"bin3", "lib4", "restricted_allows_dynamic_linking"},
				{"bin4", "bin4", "by_exception_only"},
				{"image", "bin5", "not_allowed"},
			},
			matchingAny: map[string][]act{
				"unencumbered": []act{{"image", "lib1", "unencumbered"}},
				"permissive":       []act{{"image", "lib2", "permissive"}},
				"notice":     []act{{"image", "bin1", "notice"}},
				"reciprocal":     []act{{"image", "bin2", "reciprocal"}},
				"restricted":     []act{{"image", "bin3", "restricted"}},
				"restricted_with_classpath_exception":     []act{{"bin3", "lib3", "restricted_with_classpath_exception"}},
				"restricted_allows_dynamic_linking":     []act{{"bin3", "lib4", "restricted_allows_dynamic_linking"}},
				"proprietary":     []act{{"image", "bin4", "proprietary"}},
				"by_exception_only":     []act{{"bin4", "bin4", "by_exception_only"}},
				"not_allowed":     []act{{"image", "bin5", "not_allowed"}},
				"notice|proprietary":     []act{{"image", "bin1", "notice"}, {"image", "bin4", "proprietary"}},
			},
			expected: []act{
				{"image", "lib1", "unencumbered"},
				{"image", "lib2", "permissive"},
				{"image", "bin1", "notice"},
				{"image", "bin2", "reciprocal"},
				{"image", "bin3", "restricted"},
				{"bin3", "lib3", "restricted_with_classpath_exception"},
				{"bin3", "lib4", "restricted_allows_dynamic_linking"},
				{"image", "bin4", "proprietary"},
				{"bin4", "bin4", "by_exception_only"},
				{"image", "bin5", "not_allowed"},
			},
		},
		{
			name: "everythingonone",
			actions: []act{
				{"image", "lib1", "unencumbered|permissive|notice|reciprocal|restricted" +
						"|proprietary|restricted_with_classpath_exception" +
						"|restricted_allows_dynamic_linking|by_exception_only|not_allowed"},
			},
			plus: &[]act{
				{"image", "lib1", "restricted_with_classpath_exception"},
				{"image", "lib1", "restricted_allows_dynamic_linking"},
				{"image", "lib1", "by_exception_only"},
				{"image", "lib1", "not_allowed"},
			},
			matchingAny: map[string][]act{
				"unencumbered": []act{{"image", "lib1", "unencumbered"}},
				"permissive": []act{{"image", "lib1", "permissive"}},
				"notice": []act{{"image", "lib1", "notice"}},
				"reciprocal": []act{{"image", "lib1", "reciprocal"}},
				"restricted": []act{{"image", "lib1", "restricted"}},
				"restricted_with_classpath_exception": []act{{"image", "lib1", "restricted_with_classpath_exception"}},
				"restricted_allows_dynamic_linking": []act{{"image", "lib1", "restricted_allows_dynamic_linking"}},
				"proprietary": []act{{"image", "lib1", "proprietary"}},
				"by_exception_only": []act{{"image", "lib1", "by_exception_only"}},
				"not_allowed": []act{{"image", "lib1", "not_allowed"}},
				"notice|proprietary": []act{{"image", "lib1", "notice|proprietary"}},
			},
			expected: []act{
				{"image", "lib1", "unencumbered|permissive|notice|reciprocal|restricted" +
						"|proprietary|restricted_with_classpath_exception" +
						"|restricted_allows_dynamic_linking|by_exception_only|not_allowed"},
			},
		},
		{
			name: "everythingplusone",
			actions: []act{
				{"image", "lib1", "unencumbered|permissive|notice|reciprocal|restricted|proprietary"},
			},
			plus: &[]act{
				{"image", "lib1", "restricted_with_classpath_exception"},
				{"image", "lib1", "restricted_allows_dynamic_linking"},
				{"image", "lib1", "by_exception_only"},
				{"image", "lib1", "not_allowed"},
			},
			matchingAny: map[string][]act{
				"unencumbered": []act{{"image", "lib1", "unencumbered"}},
				"permissive": []act{{"image", "lib1", "permissive"}},
				"notice": []act{{"image", "lib1", "notice"}},
				"reciprocal": []act{{"image", "lib1", "reciprocal"}},
				"restricted": []act{{"image", "lib1", "restricted"}},
				"restricted_with_classpath_exception": []act{{"image", "lib1", "restricted_with_classpath_exception"}},
				"restricted_allows_dynamic_linking": []act{{"image", "lib1", "restricted_allows_dynamic_linking"}},
				"proprietary": []act{{"image", "lib1", "proprietary"}},
				"by_exception_only": []act{{"image", "lib1", "by_exception_only"}},
				"not_allowed": []act{{"image", "lib1", "not_allowed"}},
				"notice|proprietary": []act{{"image", "lib1", "notice|proprietary"}},
			},
			expected: []act{
				{"image", "lib1", "unencumbered|permissive|notice|reciprocal|restricted" +
						"|proprietary|restricted_with_classpath_exception" +
						"|restricted_allows_dynamic_linking|by_exception_only|not_allowed"},
			},
		},
		{
			name: "allbutone",
			actions: []act{
				{"image", "lib1", "unencumbered|permissive|notice|reciprocal|restricted|proprietary"},
			},
			plus: &[]act{
				{"image", "lib1", "restricted_allows_dynamic_linking"},
				{"image", "lib1", "by_exception_only"},
				{"image", "lib1", "not_allowed"},
			},
			matchingAny: map[string][]act{
				"unencumbered": []act{{"image", "lib1", "unencumbered"}},
				"permissive": []act{{"image", "lib1", "permissive"}},
				"notice": []act{{"image", "lib1", "notice"}},
				"reciprocal": []act{{"image", "lib1", "reciprocal"}},
				"restricted": []act{{"image", "lib1", "restricted"}},
				"restricted_with_classpath_exception": []act{},
				"restricted_allows_dynamic_linking": []act{{"image", "lib1", "restricted_allows_dynamic_linking"}},
				"proprietary": []act{{"image", "lib1", "proprietary"}},
				"by_exception_only": []act{{"image", "lib1", "by_exception_only"}},
				"not_allowed": []act{{"image", "lib1", "not_allowed"}},
				"notice|proprietary": []act{{"image", "lib1", "notice|proprietary"}},
			},
			expected: []act{
				{"image", "lib1", "unencumbered|permissive|notice|reciprocal|restricted" +
						"|proprietary|restricted_allows_dynamic_linking|by_exception_only|not_allowed"},
			},
		},
		{
			name: "restrictedplus",
			actions: []act{
				{"bin1", "bin1", "restricted"},
				{"bin1", "bin1", "restricted_with_classpath_exception"},
				{"bin1", "bin1", "restricted_allows_dynamic_linking"},
			},
			plus: &[]act{
				{"bin1", "bin1", "permissive|notice"},
				{"bin1", "bin1", "restricted|proprietary"},
			},
			matchingAny: map[string][]act{
				"unencumbered":     []act{},
				"permissive":     []act{{"bin1", "bin1", "permissive"}},
				"notice":     []act{{"bin1", "bin1", "notice"}},
				"restricted":     []act{{"bin1", "bin1", "restricted"}},
				"restricted_with_classpath_exception":     []act{{"bin1", "bin1", "restricted_with_classpath_exception"}},
				"restricted_allows_dynamic_linking":     []act{{"bin1", "bin1", "restricted_allows_dynamic_linking"}},
				"proprietary":     []act{{"bin1", "bin1", "proprietary"}},
				"restricted|proprietary":     []act{{"bin1", "bin1", "restricted|proprietary"}},
				"by_exception_only": []act{},
				"proprietary|by_exception_only":     []act{{"bin1", "bin1", "proprietary"}},
			},
			expected: []act{{"bin1", "bin1", "permissive|notice|restricted|restricted_with_classpath_exception|restricted_allows_dynamic_linking|proprietary"}},
		},
	}
	for _, tt := range tests {
		toConditionSet := func(names []string) LicenseConditionSet {
			result := LicenseConditionSet(0x0000)
			for _, name := range names {
				result |= LicenseConditionSet(RecognizedConditionNames[name])
			}
			return result
		}
		populate := func(lg *LicenseGraph, t *testing.T) actionSet {
			as := toActionSet(lg, tt.actions)
			t.Logf("action set(%v): %s", tt.actions, as.String())
			if tt.plus != nil {
				for _, a := range *tt.plus {
					for _, lc := range strings.Split(a.condition, "|") {
						as.addCondition(newTestNode(lg, a.actsOn), newTestCondition(lg, a.origin, lc))
					}
				}
			}
			return as
		}
		populateSet := func(lg *LicenseGraph, t *testing.T) actionSet {
			result := actionSet{lg, &IntervalSet{}}
			result.addSet(toActionSet(lg, tt.actions))
			t.Logf("action set(%v): %s", tt.actions, result.String())
			if tt.plus != nil {
				result.addSet(toActionSet(lg, *tt.plus))
			}
			return result
		}
		checkMatching := func(as actionSet, t *testing.T) {
			for data, expectedActions := range tt.matchingAny {
				expected := toActionSet(as.lg, expectedActions)
				actual := actionSet{as.lg, &IntervalSet{}}
				actual.addAll(as.matchingAny(toConditionSet(strings.Split(data, "|")).AsList()...))

				t.Logf("matchingAny(%s): actual set %s", data, actual.String())
				t.Logf("matchingAny(%s): expected set %s", data, expected.String())
				actual.VisitAll(func(actualAction resolutionAction) {
					var expectedAction resolutionAction
					if !expected.FindFirst(&expectedAction, func(expectedAction resolutionAction) bool { return actualAction.actsOn == expectedAction.actsOn}) {
						t.Errorf("matchingAny(%s): unexpected node: found %s, want missing", data, actualAction.actsOn.Name())
						return
					}
					if expectedAction.cs != (actualAction.cs & expectedAction.cs) {
						t.Errorf("matchingAny(%s): unexpected conditions: got %q, want %q", data, (actualAction.cs & expectedAction.cs).Names(), expectedAction.cs.Names())
					}
				})
				expected.VisitAll(func(expectedAction resolutionAction) {
					var actualAction resolutionAction
					if !actual.FindFirst(&actualAction, func(actualAction resolutionAction) bool { return actualAction.actsOn == expectedAction.actsOn}) {
						t.Errorf("matchigAny(%s): missing node: found missing, want %s", data, expectedAction.actsOn.Name())
					}
				})
			}
		}
		checkMatchingSet := func(as actionSet, t *testing.T) {
			for data, expectedActions := range tt.matchingAny {
				expected := toActionSet(as.lg, expectedActions)
				actual := actionSet{as.lg, &IntervalSet{}}
				actual.addAll(as.matchingAnySet(toConditionSet(strings.Split(data, "|"))))

				t.Logf("matchingAnySet(%s): actual set %s", data, actual.String())
				t.Logf("matchingAnySet(%s): expected set %s", data, expected.String())

				actual.VisitAll(func(actualAction resolutionAction) {
					var expectedAction resolutionAction
					if !expected.FindFirst(&expectedAction, func(expectedAction resolutionAction) bool { return actualAction.actsOn == expectedAction.actsOn}) {
						t.Errorf("matchingAnySet(%s): unexpected node: found %s, want missing", data, actualAction.actsOn.Name())
						return
					}
					if expectedAction.cs != (actualAction.cs & expectedAction.cs) {
						t.Errorf("matchingAnySet(%s): unexpected conditions: got %q, want %q", data, (actualAction.cs & expectedAction.cs).Names(), expectedAction.cs.Names())
					}
				})
				expected.VisitAll(func(expectedAction resolutionAction) {
					var actualAction resolutionAction
					if !actual.FindFirst(&actualAction, func(actualAction resolutionAction) bool { return actualAction.actsOn == expectedAction.actsOn}) {
						t.Errorf("matchigAnySet(%s): missing node: found missing, want %s", data, expectedAction.actsOn.Name())
					}
				})
			}
		}
		actListString := func(l []act) string {
			var sb strings.Builder
			sep := ""
			for _, a := range l {
				fmt.Fprintf(&sb, "%s%s", sep, a.String())
				sep = ", "
			}
			return sb.String()
		}

		checkExpected := func(actual actionSet, t *testing.T) bool {
			t.Logf("checkExpected{%s}", actListString(tt.expected))

			expected := toActionSet(actual.lg, tt.expected)

			t.Logf("actual action set: %s", actual)
			t.Logf("expected action set: %s", expected)

			if !actual.isEqual(expected) {
				t.Errorf("checkExpected: got %s, want %s", actual, expected)
				return false
			}

			expectedStrings := make([]string, 0)
			actualStrings := make([]string, 0)
			var first resolutionAction
			expected.FindFirst(&first, func(a resolutionAction) bool {
				expectedStrings = append(expectedStrings, a.String())
				return false
			})

			actual.VisitAll(func(a resolutionAction) {
				actualStrings = append(actualStrings, a.String())
			})

			t.Logf("actual strings: %v", actualStrings)
			t.Logf("expected strings: %v", expectedStrings)

			if len(actualStrings) != len(expectedStrings) {
				t.Errorf("len(actualStrings): got %d, want %d", len(actualStrings), len(expectedStrings))
			} else {
				for i := 0; i < len(actualStrings); i++ {
					if actualStrings[i] != expectedStrings[i] {
						t.Errorf("actualStrings[%d]: got %s, want %s", i, actualStrings[i], expectedStrings[i])
						break
					}
				}
			}

			found := actual.FindFirst(&first, func(a resolutionAction) bool { return true })
			if len(tt.expected) == 0 {
				if found {
					t.Errorf("actual.FindFirst(&first, true): got true with found %s, want false", first.String())
				}
			} else if len(expectedStrings) == 0 {
				t.Errorf("error in test len(tt.expected)==%d, got len(expectedStrings)==0, want > 0", len(tt.expected))
			} else if !found {
				t.Errorf("actual.FindFirst(&first, true): got false, want true with found %s", expectedStrings[0])
			} else if first.String() != expectedStrings[0] {
				t.Errorf("actual.FindFirst(&first, true): got true found %s, want true found %s", first.String(), expectedStrings[0])
			}

			expectedConditions := expected.conditions().AsList()

			if len(tt.expected) == 0 {
				if !actual.isEmpty() {
					t.Errorf("actual.isEmpty(): got false, want true")
				}
				if actual.hasAny(expectedConditions...) {
					t.Errorf("actual.hasAny(): got true, want false")
				}
			} else {
				if actual.isEmpty() {
					t.Errorf("actual.isEmpty(): got true, want false")
				}
				if !actual.hasAny(expectedConditions...) {
					t.Errorf("actual.hasAny(all expected): got false, want true")
				}
			}
			for _, expectedCondition := range expectedConditions {
				if !actual.hasAny(expectedCondition) {
					t.Errorf("actual.hasAny(%q): got false, want true", expectedCondition.Name())
				}
			}

			notExpected := (AllLicenseConditions &^ expected.conditions())
			notExpectedList := notExpected.AsList()
			t.Logf("not expected license condition set: %04x {%s}", notExpected, strings.Join(notExpected.Names(), ", "))

			if len(tt.expected) == 0 {
				if actual.hasAny(append(expectedConditions, notExpectedList...)...) {
					t.Errorf("actual.hasAny(all conditions): want false, got true")
				}
			} else {
				if !actual.hasAny(append(expectedConditions, notExpectedList...)...) {
					t.Errorf("actual.hasAny(all conditions): want true, got false")
				}
			}
			for _, unexpectedCondition := range notExpectedList {
				if actual.hasAny(unexpectedCondition) {
					t.Errorf("actual.hasAny(%q): got true, want false", unexpectedCondition.Name())
				}
			}
			return true
		}

		checkExpectedSet := func(actual actionSet, t *testing.T) bool {
			t.Logf("checkExpectedSet{%s}", actListString(tt.expected))

			expected := toActionSet(actual.lg, tt.expected)

			t.Logf("actual action set: %s", actual)
			t.Logf("expected action set: %s", expected)

			if !actual.isEqual(expected) {
				t.Errorf("checkExpectedSet: got %s, want %s", actual, expected)
				return false
			}

			expectedConditions := expected.conditions()

			if len(tt.expected) == 0 {
				if !actual.isEmpty() {
					t.Errorf("actual.isEmpty(): got false, want true")
				}
				if actual.matchesAnySet(expectedConditions) {
					t.Errorf("actual.matchesAnySet({}): got true, want false")
				}
			} else {
				if actual.isEmpty() {
					t.Errorf("actual.isEmpty(): got true, want false")
				}
				if !actual.matchesAnySet(expectedConditions) {
					t.Errorf("actual.matchesAnySet({all expected}): want true, got false")
				}
			}

			notExpected := (AllLicenseConditions &^ expectedConditions)
			t.Logf("not expected license condition set: %04x {%s}", notExpected, strings.Join(notExpected.Names(), ", "))

			if len(tt.expected) == 0 {
				if actual.matchesAnySet(expectedConditions, notExpected) {
					t.Errorf("empty actual.matchesAnySet({expected}, {not expected}): want false, got true")
				}
			} else {
				if !actual.matchesAnySet(expectedConditions, notExpected) {
					t.Errorf("actual.matchesAnySet({expected}, {not expected}): want true, got false")
				}
			}
			if actual.matchesAnySet(notExpected) {
				t.Errorf("actual.matchesAnySet({not expected}): want false, got true")
			}
			return true
		}

		t.Run(tt.name, func(t *testing.T) {
			as := populate(newLicenseGraph(), t)
			if checkExpected(as, t) {
				checkMatching(as, t)
			}
			if checkExpectedSet(as, t) {
				checkMatchingSet(as, t)
			}
		})

		t.Run(tt.name+"_sets", func(t *testing.T) {
			as := populateSet(newLicenseGraph(), t)
			if checkExpected(as, t) {
				checkMatching(as, t)
			}
			if checkExpectedSet(as, t){
				checkMatchingSet(as, t)
			}
		})
	}
}
