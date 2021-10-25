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
	"strings"
	"testing"
)

func toConditionList(lg *licenseGraphImp, conditions map[string][]string) ConditionList {
	cl := make(ConditionList, 0)
	for name, origins := range conditions {
		for _, origin := range origins {
			cl = append(cl, licenseConditionImp{name, newTestNode(lg, origin).(targetNodeImp)})
		}
	}
	return cl
}

func TestConditionNames(t *testing.T) {
	impliesShare := ConditionNames([]string{"restricted", "reciprocal"})

	if impliesShare.Contains("notice") {
		t.Errorf("impliesShare.Contains(\"notice\") got true, want false")
	}

	if !impliesShare.Contains("restricted") {
		t.Errorf("impliesShare.Contains(\"restricted\") got false, want true")
	}

	if !impliesShare.Contains("reciprocal") {
		t.Errorf("impliesShare.Contains(\"reciprocal\") got false, want true")
	}

	if impliesShare.Contains("") {
		t.Errorf("impliesShare.Contains(\"\") got true, want false")
	}
}

func newTestNode(lg *licenseGraphImp, targetName string) TargetNode {
	lg.targets[targetName] = &targetNode{name: targetName}
	return targetNodeImp{lg, targetName}
}

func TestConditionList(t *testing.T) {
	tests := []struct {
		name       string
		conditions map[string][]string
		byName     map[string][]string
		byOrigin   map[string][]string
	}{
		{
			name: "noticeonly",
			conditions: map[string][]string{
				"notice": []string{"bin1", "lib1"},
			},
			byName: map[string][]string{
				"notice":     []string{"bin1", "lib1"},
				"restricted": []string{},
			},
			byOrigin: map[string][]string{
				"bin1": []string{"notice"},
				"lib1": []string{"notice"},
				"bin2": []string{},
				"lib2": []string{},
			},
		},
		{
			name:       "empty",
			conditions: map[string][]string{},
			byName: map[string][]string{
				"notice":     []string{},
				"restricted": []string{},
			},
			byOrigin: map[string][]string{
				"bin1": []string{},
				"lib1": []string{},
				"bin2": []string{},
				"lib2": []string{},
			},
		},
		{
			name: "everything",
			conditions: map[string][]string{
				"notice":            []string{"bin1", "bin2", "lib1", "lib2"},
				"reciprocal":        []string{"bin1", "bin2", "lib1", "lib2"},
				"restricted":        []string{"bin1", "bin2", "lib1", "lib2"},
				"by_exception_only": []string{"bin1", "bin2", "lib1", "lib2"},
			},
			byName: map[string][]string{
				"permissive":        []string{},
				"notice":            []string{"bin1", "bin2", "lib1", "lib2"},
				"reciprocal":        []string{"bin1", "bin2", "lib1", "lib2"},
				"restricted":        []string{"bin1", "bin2", "lib1", "lib2"},
				"by_exception_only": []string{"bin1", "bin2", "lib1", "lib2"},
			},
			byOrigin: map[string][]string{
				"bin1":  []string{"notice", "reciprocal", "restricted", "by_exception_only"},
				"bin2":  []string{"notice", "reciprocal", "restricted", "by_exception_only"},
				"lib1":  []string{"notice", "reciprocal", "restricted", "by_exception_only"},
				"lib2":  []string{"notice", "reciprocal", "restricted", "by_exception_only"},
				"other": []string{},
			},
		},
		{
			name: "allbutoneeach",
			conditions: map[string][]string{
				"notice":            []string{"bin2", "lib1", "lib2"},
				"reciprocal":        []string{"bin1", "lib1", "lib2"},
				"restricted":        []string{"bin1", "bin2", "lib2"},
				"by_exception_only": []string{"bin1", "bin2", "lib1"},
			},
			byName: map[string][]string{
				"permissive":        []string{},
				"notice":            []string{"bin2", "lib1", "lib2"},
				"reciprocal":        []string{"bin1", "lib1", "lib2"},
				"restricted":        []string{"bin1", "bin2", "lib2"},
				"by_exception_only": []string{"bin1", "bin2", "lib1"},
			},
			byOrigin: map[string][]string{
				"bin1":  []string{"reciprocal", "restricted", "by_exception_only"},
				"bin2":  []string{"notice", "restricted", "by_exception_only"},
				"lib1":  []string{"notice", "reciprocal", "by_exception_only"},
				"lib2":  []string{"notice", "reciprocal", "restricted"},
				"other": []string{},
			},
		},
		{
			name: "oneeach",
			conditions: map[string][]string{
				"notice":            []string{"bin1"},
				"reciprocal":        []string{"bin2"},
				"restricted":        []string{"lib1"},
				"by_exception_only": []string{"lib2"},
			},
			byName: map[string][]string{
				"permissive":        []string{},
				"notice":            []string{"bin1"},
				"reciprocal":        []string{"bin2"},
				"restricted":        []string{"lib1"},
				"by_exception_only": []string{"lib2"},
			},
			byOrigin: map[string][]string{
				"bin1":  []string{"notice"},
				"bin2":  []string{"reciprocal"},
				"lib1":  []string{"restricted"},
				"lib2":  []string{"by_exception_only"},
				"other": []string{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lg := newLicenseGraphImp()
			cl := toConditionList(lg, tt.conditions)
			for names, expected := range tt.byName {
				name := ConditionNames(strings.Split(names, ":"))
				if cl.HasByName(name) {
					if len(expected) == 0 {
						t.Errorf("unexpected ConditionList.HasByName(%q): got true, want false", name)
					}
				} else {
					if len(expected) != 0 {
						t.Errorf("unexpected ConditionList.HasByName(%q): got false, want true", name)
					}
				}
				if len(expected) != cl.CountByName(name) {
					t.Errorf("unexpected ConditionList.CountByName(%q): got %d, want %d", name, cl.CountByName(name), len(expected))
				}
				byName := cl.ByName(name)
				if len(expected) != len(byName) {
					t.Errorf("unexpected ConditionList.ByName(%q): got %v, want %v", name, byName, expected)
				} else {
					sort.Strings(expected)
					actual := make([]string, 0, len(byName))
					for _, lc := range byName {
						actual = append(actual, lc.Origin().Name())
					}
					sort.Strings(actual)
					for i := 0; i < len(expected); i++ {
						if expected[i] != actual[i] {
							t.Errorf("unexpected ConditionList.ByName(%q) index %d in %v: got %s, want %s", name, i, actual, actual[i], expected[i])
						}
					}
				}
			}
			for origin, expected := range tt.byOrigin {
				onode := newTestNode(lg, origin)
				if cl.HasByOrigin(onode) {
					if len(expected) == 0 {
						t.Errorf("unexpected ConditionList.HasByOrigin(%q): got true, want false", origin)
					}
				} else {
					if len(expected) != 0 {
						t.Errorf("unexpected ConditionList.HasByOrigin(%q): got false, want true", origin)
					}
				}
				if len(expected) != cl.CountByOrigin(onode) {
					t.Errorf("unexpected ConditionList.CountByOrigin(%q): got %d, want %d", origin, cl.CountByOrigin(onode), len(expected))
				}
				byOrigin := cl.ByOrigin(onode)
				if len(expected) != len(byOrigin) {
					t.Errorf("unexpected ConditionList.ByOrigin(%q): got %v, want %v", origin, byOrigin, expected)
				} else {
					sort.Strings(expected)
					actual := make([]string, 0, len(byOrigin))
					for _, lc := range byOrigin {
						actual = append(actual, lc.Name())
					}
					sort.Strings(actual)
					for i := 0; i < len(expected); i++ {
						if expected[i] != actual[i] {
							t.Errorf("unexpected ConditionList.ByOrigin(%q) index %d in %v: got %s, want %s", origin, i, actual, actual[i], expected[i])
						}
					}
				}
			}
		})
	}
}

type byName map[string][]string

func (bn byName) checkPublic(ls LicenseConditionSet, t *testing.T) {
	for names, expected := range bn {
		name := ConditionNames(strings.Split(names, ":"))
		if ls.HasAnyByName(name) {
			if len(expected) == 0 {
				t.Errorf("unexpected LicenseConditionSet.HasAnyByName(%q): got true, want false", name)
			}
		} else {
			if len(expected) != 0 {
				t.Errorf("unexpected LicenseConditionSet.HasAnyByName(%q): got false, want true", name)
			}
		}
		if len(expected) != ls.CountByName(name) {
			t.Errorf("unexpected LicenseConditionSet.CountByName(%q): got %d, want %d", name, ls.CountByName(name), len(expected))
		}
		byName := ls.ByName(name)
		if len(expected) != len(byName) {
			t.Errorf("unexpected LicenseConditionSet.ByName(%q): got %v, want %v", name, byName, expected)
		} else {
			sort.Strings(expected)
			actual := make([]string, 0, len(byName))
			for _, lc := range byName {
				actual = append(actual, lc.Origin().Name())
			}
			sort.Strings(actual)
			for i := 0; i < len(expected); i++ {
				if expected[i] != actual[i] {
					t.Errorf("unexpected LicenseConditionSet.ByName(%q) index %d in %v: got %s, want %s", name, i, actual, actual[i], expected[i])
				}
			}
		}
	}
}

type byOrigin map[string][]string

func (bo byOrigin) checkPublic(lg *licenseGraphImp, ls LicenseConditionSet, t *testing.T) {
	expectedCount := 0
	for origin, expected := range bo {
		expectedCount += len(expected)
		onode := newTestNode(lg, origin)
		if ls.HasAnyByOrigin(onode) {
			if len(expected) == 0 {
				t.Errorf("unexpected LicenseConditionSet.HasAnyByOrigin(%q): got true, want false", origin)
			}
		} else {
			if len(expected) != 0 {
				t.Errorf("unexpected LicenseConditionSet.HasAnyByOrigin(%q): got false, want true", origin)
			}
		}
		if len(expected) != ls.CountByOrigin(onode) {
			t.Errorf("unexpected LicenseConditionSet.CountByOrigin(%q): got %d, want %d", origin, ls.CountByOrigin(onode), len(expected))
		}
		byOrigin := ls.ByOrigin(onode)
		if len(expected) != len(byOrigin) {
			t.Errorf("unexpected LicenseConditionSet.ByOrigin(%q): got %v, want %v", origin, byOrigin, expected)
		} else {
			sort.Strings(expected)
			actual := make([]string, 0, len(byOrigin))
			for _, lc := range byOrigin {
				actual = append(actual, lc.Name())
			}
			sort.Strings(actual)
			for i := 0; i < len(expected); i++ {
				if expected[i] != actual[i] {
					t.Errorf("unexpected LicenseConditionSet.ByOrigin(%q) index %d in %v: got %s, want %s", origin, i, actual, actual[i], expected[i])
				}
			}
		}
	}
	if expectedCount != ls.Count() {
		t.Errorf("unexpected LicenseConditionSet.Count(): got %d, want %d", ls.Count(), expectedCount)
	}
	if ls.IsEmpty() {
		if expectedCount != 0 {
			t.Errorf("unexpected LicenseConditionSet.IsEmpty(): got true, want false")
		}
	} else {
		if expectedCount == 0 {
			t.Errorf("unexpected LicenseConditionSet.IsEmpty(): got false, want true")
		}
	}
}

func TestConditionSet(t *testing.T) {
	tests := []struct {
		name       string
		conditions map[string][]string
		add        map[string][]string
		byName     map[string][]string
		byOrigin   map[string][]string
	}{
		{
			name:       "empty",
			conditions: map[string][]string{},
			add:        map[string][]string{},
			byName: map[string][]string{
				"notice":     []string{},
				"restricted": []string{},
			},
			byOrigin: map[string][]string{
				"bin1": []string{},
				"lib1": []string{},
				"bin2": []string{},
				"lib2": []string{},
			},
		},
		{
			name: "noticeonly",
			conditions: map[string][]string{
				"notice": []string{"bin1", "lib1"},
			},
			byName: map[string][]string{
				"notice":     []string{"bin1", "lib1"},
				"restricted": []string{},
			},
			byOrigin: map[string][]string{
				"bin1": []string{"notice"},
				"lib1": []string{"notice"},
				"bin2": []string{},
				"lib2": []string{},
			},
		},
		{
			name: "noticeonlyadded",
			conditions: map[string][]string{
				"notice": []string{"bin1", "lib1"},
			},
			add: map[string][]string{
				"notice": []string{"bin1", "bin2"},
			},
			byName: map[string][]string{
				"notice":     []string{"bin1", "bin2", "lib1"},
				"restricted": []string{},
			},
			byOrigin: map[string][]string{
				"bin1": []string{"notice"},
				"lib1": []string{"notice"},
				"bin2": []string{"notice"},
				"lib2": []string{},
			},
		},
		{
			name: "everything",
			conditions: map[string][]string{
				"notice":            []string{"bin1", "bin2", "lib1", "lib2"},
				"reciprocal":        []string{"bin1", "bin2", "lib1", "lib2"},
				"restricted":        []string{"bin1", "bin2", "lib1", "lib2"},
				"by_exception_only": []string{"bin1", "bin2", "lib1", "lib2"},
			},
			add: map[string][]string{
				"notice":            []string{"bin1", "bin2", "lib1", "lib2"},
				"reciprocal":        []string{"bin1", "bin2", "lib1", "lib2"},
				"restricted":        []string{"bin1", "bin2", "lib1", "lib2"},
				"by_exception_only": []string{"bin1", "bin2", "lib1", "lib2"},
			},
			byName: map[string][]string{
				"permissive":        []string{},
				"notice":            []string{"bin1", "bin2", "lib1", "lib2"},
				"reciprocal":        []string{"bin1", "bin2", "lib1", "lib2"},
				"restricted":        []string{"bin1", "bin2", "lib1", "lib2"},
				"by_exception_only": []string{"bin1", "bin2", "lib1", "lib2"},
			},
			byOrigin: map[string][]string{
				"bin1":  []string{"notice", "reciprocal", "restricted", "by_exception_only"},
				"bin2":  []string{"notice", "reciprocal", "restricted", "by_exception_only"},
				"lib1":  []string{"notice", "reciprocal", "restricted", "by_exception_only"},
				"lib2":  []string{"notice", "reciprocal", "restricted", "by_exception_only"},
				"other": []string{},
			},
		},
		{
			name: "allbutoneeach",
			conditions: map[string][]string{
				"notice":            []string{"bin2", "lib1", "lib2"},
				"reciprocal":        []string{"bin1", "lib1", "lib2"},
				"restricted":        []string{"bin1", "bin2", "lib2"},
				"by_exception_only": []string{"bin1", "bin2", "lib1"},
			},
			byName: map[string][]string{
				"permissive":        []string{},
				"notice":            []string{"bin2", "lib1", "lib2"},
				"reciprocal":        []string{"bin1", "lib1", "lib2"},
				"restricted":        []string{"bin1", "bin2", "lib2"},
				"by_exception_only": []string{"bin1", "bin2", "lib1"},
			},
			byOrigin: map[string][]string{
				"bin1":  []string{"reciprocal", "restricted", "by_exception_only"},
				"bin2":  []string{"notice", "restricted", "by_exception_only"},
				"lib1":  []string{"notice", "reciprocal", "by_exception_only"},
				"lib2":  []string{"notice", "reciprocal", "restricted"},
				"other": []string{},
			},
		},
		{
			name: "allbutoneeachadded",
			conditions: map[string][]string{
				"notice":            []string{"bin2", "lib1", "lib2"},
				"reciprocal":        []string{"bin1", "lib1", "lib2"},
				"restricted":        []string{"bin1", "bin2", "lib2"},
				"by_exception_only": []string{"bin1", "bin2", "lib1"},
			},
			add: map[string][]string{
				"notice":            []string{"bin2", "lib1", "lib2"},
				"reciprocal":        []string{"bin1", "lib1", "lib2"},
				"restricted":        []string{"bin1", "bin2", "lib2"},
				"by_exception_only": []string{"bin1", "bin2", "lib1"},
			},
			byName: map[string][]string{
				"permissive":        []string{},
				"notice":            []string{"bin2", "lib1", "lib2"},
				"reciprocal":        []string{"bin1", "lib1", "lib2"},
				"restricted":        []string{"bin1", "bin2", "lib2"},
				"by_exception_only": []string{"bin1", "bin2", "lib1"},
			},
			byOrigin: map[string][]string{
				"bin1":  []string{"reciprocal", "restricted", "by_exception_only"},
				"bin2":  []string{"notice", "restricted", "by_exception_only"},
				"lib1":  []string{"notice", "reciprocal", "by_exception_only"},
				"lib2":  []string{"notice", "reciprocal", "restricted"},
				"other": []string{},
			},
		},
		{
			name: "allbutoneeachfilled",
			conditions: map[string][]string{
				"notice":            []string{"bin2", "lib1", "lib2"},
				"reciprocal":        []string{"bin1", "lib1", "lib2"},
				"restricted":        []string{"bin1", "bin2", "lib2"},
				"by_exception_only": []string{"bin1", "bin2", "lib1"},
			},
			add: map[string][]string{
				"notice":            []string{"bin1", "bin2", "lib1"},
				"reciprocal":        []string{"bin1", "bin2", "lib2"},
				"restricted":        []string{"bin1", "lib1", "lib2"},
				"by_exception_only": []string{"bin2", "lib1", "lib2"},
			},
			byName: map[string][]string{
				"permissive":        []string{},
				"notice":            []string{"bin1", "bin2", "lib1", "lib2"},
				"reciprocal":        []string{"bin1", "bin2", "lib1", "lib2"},
				"restricted":        []string{"bin1", "bin2", "lib1", "lib2"},
				"by_exception_only": []string{"bin1", "bin2", "lib1", "lib2"},
			},
			byOrigin: map[string][]string{
				"bin1":  []string{"notice", "reciprocal", "restricted", "by_exception_only"},
				"bin2":  []string{"notice", "reciprocal", "restricted", "by_exception_only"},
				"lib1":  []string{"notice", "reciprocal", "restricted", "by_exception_only"},
				"lib2":  []string{"notice", "reciprocal", "restricted", "by_exception_only"},
				"other": []string{},
			},
		},
		{
			name: "oneeach",
			conditions: map[string][]string{
				"notice":            []string{"bin1"},
				"reciprocal":        []string{"bin2"},
				"restricted":        []string{"lib1"},
				"by_exception_only": []string{"lib2"},
			},
			byName: map[string][]string{
				"permissive":        []string{},
				"notice":            []string{"bin1"},
				"reciprocal":        []string{"bin2"},
				"restricted":        []string{"lib1"},
				"by_exception_only": []string{"lib2"},
			},
			byOrigin: map[string][]string{
				"bin1":  []string{"notice"},
				"bin2":  []string{"reciprocal"},
				"lib1":  []string{"restricted"},
				"lib2":  []string{"by_exception_only"},
				"other": []string{},
			},
		},
		{
			name: "oneeachoverlap",
			conditions: map[string][]string{
				"notice":            []string{"bin1"},
				"reciprocal":        []string{"bin2"},
				"restricted":        []string{"lib1"},
				"by_exception_only": []string{"lib2"},
			},
			add: map[string][]string{
				"notice":            []string{"lib2"},
				"reciprocal":        []string{"lib1"},
				"restricted":        []string{"bin2"},
				"by_exception_only": []string{"bin1"},
			},
			byName: map[string][]string{
				"permissive":        []string{},
				"notice":            []string{"bin1", "lib2"},
				"reciprocal":        []string{"bin2", "lib1"},
				"restricted":        []string{"bin2", "lib1"},
				"by_exception_only": []string{"bin1", "lib2"},
			},
			byOrigin: map[string][]string{
				"bin1":  []string{"by_exception_only", "notice"},
				"bin2":  []string{"reciprocal", "restricted"},
				"lib1":  []string{"reciprocal", "restricted"},
				"lib2":  []string{"by_exception_only", "notice"},
				"other": []string{},
			},
		},
		{
			name: "oneeachadded",
			conditions: map[string][]string{
				"notice":            []string{"bin1"},
				"reciprocal":        []string{"bin2"},
				"restricted":        []string{"lib1"},
				"by_exception_only": []string{"lib2"},
			},
			add: map[string][]string{
				"notice":            []string{"bin1"},
				"reciprocal":        []string{"bin2"},
				"restricted":        []string{"lib1"},
				"by_exception_only": []string{"lib2"},
			},
			byName: map[string][]string{
				"permissive":        []string{},
				"notice":            []string{"bin1"},
				"reciprocal":        []string{"bin2"},
				"restricted":        []string{"lib1"},
				"by_exception_only": []string{"lib2"},
			},
			byOrigin: map[string][]string{
				"bin1":  []string{"notice"},
				"bin2":  []string{"reciprocal"},
				"lib1":  []string{"restricted"},
				"lib2":  []string{"by_exception_only"},
				"other": []string{},
			},
		},
	}
	for _, tt := range tests {
		testPublicInterface := func(lg *licenseGraphImp, ls LicenseConditionSet, t *testing.T) {
			byName(tt.byName).checkPublic(ls, t)
			byOrigin(tt.byOrigin).checkPublic(lg, ls, t)
		}
		t.Run(tt.name+"_public_interface", func(t *testing.T) {
			lg := newLicenseGraphImp()
			ls := NewLicenseConditionSet(toConditionList(lg, tt.conditions)...)
			if tt.add != nil {
				ls.Add(toConditionList(lg, tt.add)...)
			}
			testPublicInterface(lg, ls, t)
		})

		t.Run("Copy() of "+tt.name+"_public_interface", func(t *testing.T) {
			lg := newLicenseGraphImp()
			ls := NewLicenseConditionSet(toConditionList(lg, tt.conditions)...)
			if tt.add != nil {
				ls.Add(toConditionList(lg, tt.add)...)
			}
			testPublicInterface(lg, ls.Copy(), t)
		})

		testPrivateInterface := func(lgimp *licenseGraphImp, lsimp *licenseConditionSetImp, t *testing.T) {
			slist := make([]string, 0, lsimp.Count())
			for origin, expected := range tt.byOrigin {
				for _, name := range expected {
					slist = append(slist, origin+";"+name)
				}
				if lsimp.hasByOrigin(origin) {
					if len(expected) == 0 {
						t.Errorf("unexpected licenseConditionSetImp.hasByOrigin(%q): got true, want false", origin)
					}
				} else {
					if len(expected) != 0 {
						t.Errorf("unexpected licenseConditionSetImp.hasByOrigin(%q): got false, want true", origin)
					}
				}
				if len(expected) != lsimp.countByOrigin(origin) {
					t.Errorf("unexpected licenseConditionSetImp.countByOrigin(%q): got %d, want %d", origin, lsimp.countByOrigin(origin), len(expected))
				}
				byOrigin := lsimp.byOrigin(origin)
				if len(expected) != len(byOrigin) {
					t.Errorf("unexpected licenseConditionSetImp.byOrigin(%q): got %v, want %v", origin, byOrigin, expected)
				} else {
					sort.Strings(expected)
					actual := make([]string, 0, len(byOrigin))
					for _, lc := range byOrigin {
						actual = append(actual, lc.Name())
					}
					sort.Strings(actual)
					for i := 0; i < len(expected); i++ {
						if expected[i] != actual[i] {
							t.Errorf("unexpected licenseConditionSetImp.byOrigin(%q) index %d in %v: got %s, want %s", origin, i, actual, actual[i], expected[i])
						}
					}
				}
			}
			actualSlist := lsimp.asStringList(";")
			if len(slist) != len(actualSlist) {
				t.Errorf("unexpected licenseConditionSetImp.asStringList(\";\"): got %v, want %v", actualSlist, slist)
			} else {
				sort.Strings(slist)
				sort.Strings(actualSlist)
				for i := 0; i < len(slist); i++ {
					if slist[i] != actualSlist[i] {
						t.Errorf("unexpected licenseConditionSetImp.asStringList(\";\") index %d in %v: got %s, want %s", i, actualSlist, actualSlist[i], slist[i])
					}
				}
			}
		}

		t.Run(tt.name+"_private_list_interface", func(t *testing.T) {
			lg := newLicenseGraphImp()
			lsimp := newLicenseConditionSet(&targetNodeImp{lg, ""})
			for name, origins := range tt.conditions {
				for _, origin := range origins {
					_ = newTestNode(lg, origin)
					lsimp.addAll(origin, name)
				}
			}
			if tt.add != nil {
				lsimp.addList(toConditionList(lg, tt.add))
			}
			testPrivateInterface(lg, lsimp, t)
		})

		t.Run("copy() of "+tt.name+"_private_list_interface", func(t *testing.T) {
			lg := newLicenseGraphImp()
			lsimp := newLicenseConditionSet(&targetNodeImp{lg, ""})
			for name, origins := range tt.conditions {
				for _, origin := range origins {
					_ = newTestNode(lg, origin)
					lsimp.addAll(origin, name)
				}
			}
			if tt.add != nil {
				lsimp.addList(toConditionList(lg, tt.add))
			}
			testPrivateInterface(lg, lsimp.copy(false), t)
		})

		t.Run(tt.name+"_private_set_interface", func(t *testing.T) {
			lg := newLicenseGraphImp()
			lsimp := newLicenseConditionSet(&targetNodeImp{lg, ""})
			for name, origins := range tt.conditions {
				for _, origin := range origins {
					lsimp.add(name, newTestNode(lg, origin).(targetNodeImp))
				}
			}
			if tt.add != nil {
				other := newLicenseConditionSet(nil)
				for name, origins := range tt.add {
					for _, origin := range origins {
						other.add(name, newTestNode(lg, origin).(targetNodeImp))
					}
				}
				lsimp.addSet(other)
			}
			testPrivateInterface(lg, lsimp, t)
		})

		t.Run("copy() of "+tt.name+"_private_set_interface", func(t *testing.T) {
			lg := newLicenseGraphImp()
			lsimp := newLicenseConditionSet(&targetNodeImp{lg, ""})
			for name, origins := range tt.conditions {
				for _, origin := range origins {
					lsimp.add(name, newTestNode(lg, origin).(targetNodeImp))
				}
			}
			if tt.add != nil {
				other := newLicenseConditionSet(nil)
				for name, origins := range tt.add {
					for _, origin := range origins {
						other.add(name, newTestNode(lg, origin).(targetNodeImp))
					}
				}
				lsimp.addSet(other)
			}
			testPrivateInterface(lg, lsimp.copy(true), t)
		})
	}
}

func TestConditionSet_Removals(t *testing.T) {
	tests := []struct {
		name           string
		conditions     map[string][]string
		removeByName   []ConditionNames
		removeByOrigin []string
		removeSet      map[string][]string
		byName         map[string][]string
		byOrigin       map[string][]string
	}{
		{
			name:         "emptybyname",
			conditions:   map[string][]string{},
			removeByName: []ConditionNames{{"reciprocal", "restricted"}},
			byName: map[string][]string{
				"notice":     []string{},
				"restricted": []string{},
			},
			byOrigin: map[string][]string{
				"bin1": []string{},
				"lib1": []string{},
				"bin2": []string{},
				"lib2": []string{},
			},
		},
		{
			name:           "emptybyorigin",
			conditions:     map[string][]string{},
			removeByOrigin: []string{"bin1", "bin2", "lib1", "lib2"},
			byName: map[string][]string{
				"notice":     []string{},
				"restricted": []string{},
			},
			byOrigin: map[string][]string{
				"bin1": []string{},
				"lib1": []string{},
				"bin2": []string{},
				"lib2": []string{},
			},
		},
		{
			name:       "emptybyset",
			conditions: map[string][]string{},
			removeSet: map[string][]string{
				"notice":     []string{"bin1", "bin2"},
				"restricted": []string{"lib1", "lib2"},
			},
			byName: map[string][]string{
				"notice":     []string{},
				"restricted": []string{},
			},
			byOrigin: map[string][]string{
				"bin1": []string{},
				"lib1": []string{},
				"bin2": []string{},
				"lib2": []string{},
			},
		},
		{
			name: "everythingremovenone",
			conditions: map[string][]string{
				"notice":            []string{"bin1", "bin2", "lib1", "lib2"},
				"reciprocal":        []string{"bin1", "bin2", "lib1", "lib2"},
				"restricted":        []string{"bin1", "bin2", "lib1", "lib2"},
				"by_exception_only": []string{"bin1", "bin2", "lib1", "lib2"},
			},
			removeByName:   []ConditionNames{{"permissive", "unencumbered"}},
			removeByOrigin: []string{"apex", "other"},
			removeSet: map[string][]string{
				"notice": []string{"apk1"},
			},
			byName: map[string][]string{
				"permissive":        []string{},
				"notice":            []string{"bin1", "bin2", "lib1", "lib2"},
				"reciprocal":        []string{"bin1", "bin2", "lib1", "lib2"},
				"restricted":        []string{"bin1", "bin2", "lib1", "lib2"},
				"by_exception_only": []string{"bin1", "bin2", "lib1", "lib2"},
			},
			byOrigin: map[string][]string{
				"bin1":  []string{"notice", "reciprocal", "restricted", "by_exception_only"},
				"bin2":  []string{"notice", "reciprocal", "restricted", "by_exception_only"},
				"lib1":  []string{"notice", "reciprocal", "restricted", "by_exception_only"},
				"lib2":  []string{"notice", "reciprocal", "restricted", "by_exception_only"},
				"other": []string{},
			},
		},
		{
			name: "everythingremovesome",
			conditions: map[string][]string{
				"notice":            []string{"bin1", "bin2", "lib1", "lib2"},
				"reciprocal":        []string{"bin1", "bin2", "lib1", "lib2"},
				"restricted":        []string{"bin1", "bin2", "lib1", "lib2"},
				"by_exception_only": []string{"bin1", "bin2", "lib1", "lib2"},
			},
			removeByName:   []ConditionNames{{"restricted", "by_exception_only"}},
			removeByOrigin: []string{"bin1", "bin2"},
			removeSet: map[string][]string{
				"notice": []string{"lib1"},
			},
			byName: map[string][]string{
				"permissive":        []string{},
				"notice":            []string{"lib2"},
				"reciprocal":        []string{"lib1", "lib2"},
				"restricted":        []string{},
				"by_exception_only": []string{},
			},
			byOrigin: map[string][]string{
				"bin1":  []string{},
				"bin2":  []string{},
				"lib1":  []string{"reciprocal"},
				"lib2":  []string{"notice", "reciprocal"},
				"other": []string{},
			},
		},
		{
			name: "everythingremoveall",
			conditions: map[string][]string{
				"notice":            []string{"bin1", "bin2", "lib1", "lib2"},
				"reciprocal":        []string{"bin1", "bin2", "lib1", "lib2"},
				"restricted":        []string{"bin1", "bin2", "lib1", "lib2"},
				"by_exception_only": []string{"bin1", "bin2", "lib1", "lib2"},
			},
			removeByName:   []ConditionNames{{"restricted", "by_exception_only"}},
			removeByOrigin: []string{"bin1", "bin2", "lib2"},
			removeSet: map[string][]string{
				"notice":     []string{"bin1", "lib1"},
				"reciprocal": []string{"bin2", "lib1"},
				"restricted": []string{"lib1", "lib2"},
			},
			byName: map[string][]string{
				"permissive":        []string{},
				"notice":            []string{},
				"reciprocal":        []string{},
				"restricted":        []string{},
				"by_exception_only": []string{},
			},
			byOrigin: map[string][]string{
				"bin1":  []string{},
				"bin2":  []string{},
				"lib1":  []string{},
				"lib2":  []string{},
				"other": []string{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			lg := newLicenseGraphImp()
			lsimp := newLicenseConditionSet(&targetNodeImp{lg, ""})
			for name, origins := range tt.conditions {
				for _, origin := range origins {
					_ = newTestNode(lg, origin)
					lsimp.addAll(origin, name)
				}
			}
			if tt.removeByName != nil {
				lsimp.removeAllByName(tt.removeByName...)
			}
			if tt.removeByOrigin != nil {
				for _, origin := range tt.removeByOrigin {
					lsimp.removeAllByOrigin(origin)
				}
			}
			if tt.removeSet != nil {
				other := newLicenseConditionSet(nil)
				for name, origins := range tt.removeSet {
					for _, origin := range origins {
						other.add(name, newTestNode(lg, origin).(targetNodeImp))
					}
				}
				lsimp.removeSet(other)
			}
			byName(tt.byName).checkPublic(lsimp, t)
			byOrigin(tt.byOrigin).checkPublic(lg, lsimp, t)
		})
	}
}

func checkAsAggregate(treatAsAggregate bool, lsimp *licenseConditionSetImp, t *testing.T) {
	for _, origins := range lsimp.conditions {
		for _, actual := range origins {
			if actual != treatAsAggregate {
				t.Errorf("unexpected treat as aggregate payload: got %v, want %v", actual, treatAsAggregate)
			}
		}
	}
}

func TestConditionSet_filter(t *testing.T) {
	tests := []struct {
		name             string
		treatAsAggregate bool
		conditions       map[string][]string
		filterNames      []ConditionNames
		byName           map[string][]string
		byOrigin         map[string][]string
	}{
		{
			name:             "empty",
			treatAsAggregate: false,
			conditions:       map[string][]string{},
			filterNames:      []ConditionNames{{"notice", "reciprocal", "restricted"}},
			byName: map[string][]string{
				"notice":     []string{},
				"reciprocal": []string{},
				"restricted": []string{},
			},
			byOrigin: map[string][]string{
				"bin1": []string{},
				"lib1": []string{},
				"bin2": []string{},
				"lib2": []string{},
			},
		},
		{
			name:             "emptycontainer",
			treatAsAggregate: true,
			conditions:       map[string][]string{},
			filterNames:      []ConditionNames{{"notice", "reciprocal", "restricted"}},
			byName: map[string][]string{
				"notice":     []string{},
				"reciprocal": []string{},
				"restricted": []string{},
			},
			byOrigin: map[string][]string{
				"bin1": []string{},
				"lib1": []string{},
				"bin2": []string{},
				"lib2": []string{},
			},
		},
		{
			name:             "everything",
			treatAsAggregate: false,
			conditions: map[string][]string{
				"notice":            []string{"bin1", "bin2", "lib1", "lib2"},
				"reciprocal":        []string{"bin1", "bin2", "lib1", "lib2"},
				"restricted":        []string{"bin1", "bin2", "lib1", "lib2"},
				"by_exception_only": []string{"bin1", "bin2", "lib1", "lib2"},
			},
			filterNames: []ConditionNames{{"permissive", "notice", "reciprocal", "restricted", "by_exception_only", "proprietary"}},
			byName: map[string][]string{
				"permissive":        []string{},
				"notice":            []string{"bin1", "bin2", "lib1", "lib2"},
				"reciprocal":        []string{"bin1", "bin2", "lib1", "lib2"},
				"restricted":        []string{"bin1", "bin2", "lib1", "lib2"},
				"by_exception_only": []string{"bin1", "bin2", "lib1", "lib2"},
			},
			byOrigin: map[string][]string{
				"bin1":  []string{"notice", "reciprocal", "restricted", "by_exception_only"},
				"bin2":  []string{"notice", "reciprocal", "restricted", "by_exception_only"},
				"lib1":  []string{"notice", "reciprocal", "restricted", "by_exception_only"},
				"lib2":  []string{"notice", "reciprocal", "restricted", "by_exception_only"},
				"other": []string{},
			},
		},
		{
			name:             "everythingcontainer",
			treatAsAggregate: true,
			conditions: map[string][]string{
				"notice":            []string{"bin1", "bin2", "lib1", "lib2"},
				"reciprocal":        []string{"bin1", "bin2", "lib1", "lib2"},
				"restricted":        []string{"bin1", "bin2", "lib1", "lib2"},
				"by_exception_only": []string{"bin1", "bin2", "lib1", "lib2"},
			},
			filterNames: []ConditionNames{{"permissive", "notice", "reciprocal", "restricted", "by_exception_only", "proprietary"}},
			byName: map[string][]string{
				"permissive":        []string{},
				"notice":            []string{"bin1", "bin2", "lib1", "lib2"},
				"reciprocal":        []string{"bin1", "bin2", "lib1", "lib2"},
				"restricted":        []string{"bin1", "bin2", "lib1", "lib2"},
				"by_exception_only": []string{"bin1", "bin2", "lib1", "lib2"},
			},
			byOrigin: map[string][]string{
				"bin1":  []string{"notice", "reciprocal", "restricted", "by_exception_only"},
				"bin2":  []string{"notice", "reciprocal", "restricted", "by_exception_only"},
				"lib1":  []string{"notice", "reciprocal", "restricted", "by_exception_only"},
				"lib2":  []string{"notice", "reciprocal", "restricted", "by_exception_only"},
				"other": []string{},
			},
		},
		{
			name:             "everythingshared",
			treatAsAggregate: false,
			conditions: map[string][]string{
				"notice":            []string{"bin1", "bin2", "lib1", "lib2"},
				"reciprocal":        []string{"bin1", "bin2", "lib1", "lib2"},
				"restricted":        []string{"bin1", "bin2", "lib1", "lib2"},
				"by_exception_only": []string{"bin1", "bin2", "lib1", "lib2"},
			},
			filterNames: []ConditionNames{{"reciprocal", "restricted"}},
			byName: map[string][]string{
				"permissive":        []string{},
				"notice":            []string{},
				"reciprocal":        []string{"bin1", "bin2", "lib1", "lib2"},
				"restricted":        []string{"bin1", "bin2", "lib1", "lib2"},
				"by_exception_only": []string{},
			},
			byOrigin: map[string][]string{
				"bin1":  []string{"reciprocal", "restricted"},
				"bin2":  []string{"reciprocal", "restricted"},
				"lib1":  []string{"reciprocal", "restricted"},
				"lib2":  []string{"reciprocal", "restricted"},
				"other": []string{},
			},
		},
		{
			name:             "oneeach",
			treatAsAggregate: true,
			conditions: map[string][]string{
				"notice":            []string{"bin1"},
				"reciprocal":        []string{"bin2"},
				"restricted":        []string{"lib1"},
				"by_exception_only": []string{"lib2"},
			},
			filterNames: []ConditionNames{{"notice", "by_exception_only"}},
			byName: map[string][]string{
				"permissive":        []string{},
				"notice":            []string{"bin1"},
				"reciprocal":        []string{},
				"restricted":        []string{},
				"by_exception_only": []string{"lib2"},
			},
			byOrigin: map[string][]string{
				"bin1":  []string{"notice"},
				"bin2":  []string{},
				"lib1":  []string{},
				"lib2":  []string{"by_exception_only"},
				"other": []string{},
			},
		},
		{
			name:             "oneeachnoneleft",
			treatAsAggregate: true,
			conditions: map[string][]string{
				"notice":            []string{"bin1"},
				"reciprocal":        []string{"bin2"},
				"restricted":        []string{"lib1"},
				"by_exception_only": []string{"lib2"},
			},
			filterNames: []ConditionNames{{"unencumbered", "permissive"}},
			byName: map[string][]string{
				"permissive":        []string{},
				"notice":            []string{},
				"reciprocal":        []string{},
				"restricted":        []string{},
				"by_exception_only": []string{},
			},
			byOrigin: map[string][]string{
				"bin1":  []string{},
				"bin2":  []string{},
				"lib1":  []string{},
				"lib2":  []string{},
				"other": []string{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			lg := newLicenseGraphImp()
			lsimp := newLicenseConditionSet(&targetNodeImp{lg, ""})
			for name, origins := range tt.conditions {
				for _, origin := range origins {
					_ = newTestNode(lg, origin)
					lsimp.addAll(origin, name)
				}
			}
			lsimp = lsimp.filter(tt.treatAsAggregate, tt.filterNames...)
			checkAsAggregate(tt.treatAsAggregate, lsimp, t)
			byName(tt.byName).checkPublic(lsimp, t)
			byOrigin(tt.byOrigin).checkPublic(lg, lsimp, t)
		})
	}
}

func TestConditionSet_rename(t *testing.T) {
	tests := []struct {
		name             string
		treatAsAggregate bool
		conditions       map[string][]string
		newName          string
		filterNames      []ConditionNames
		byName           map[string][]string
		byOrigin         map[string][]string
	}{
		{
			name:             "empty",
			treatAsAggregate: false,
			conditions:       map[string][]string{},
			newName:          "notice",
			filterNames:      []ConditionNames{{"notice", "reciprocal", "restricted"}},
			byName: map[string][]string{
				"notice":     []string{},
				"reciprocal": []string{},
				"restricted": []string{},
			},
			byOrigin: map[string][]string{
				"bin1": []string{},
				"lib1": []string{},
				"bin2": []string{},
				"lib2": []string{},
			},
		},
		{
			name:             "emptycontainer",
			treatAsAggregate: true,
			conditions:       map[string][]string{},
			newName:          "notice",
			filterNames:      []ConditionNames{{"notice", "reciprocal", "restricted"}},
			byName: map[string][]string{
				"notice":     []string{},
				"reciprocal": []string{},
				"restricted": []string{},
			},
			byOrigin: map[string][]string{
				"bin1": []string{},
				"lib1": []string{},
				"bin2": []string{},
				"lib2": []string{},
			},
		},
		{
			name:             "everything",
			treatAsAggregate: false,
			conditions: map[string][]string{
				"notice":            []string{"bin1", "bin2", "lib1", "lib2"},
				"reciprocal":        []string{"bin1", "bin2", "lib1", "lib2"},
				"restricted":        []string{"bin1", "bin2", "lib1", "lib2"},
				"by_exception_only": []string{"bin1", "bin2", "lib1", "lib2"},
			},
			newName:     "bob",
			filterNames: []ConditionNames{{"permissive", "notice", "reciprocal", "restricted", "by_exception_only", "proprietary"}},
			byName: map[string][]string{
				"permissive":        []string{},
				"notice":            []string{},
				"reciprocal":        []string{},
				"restricted":        []string{},
				"by_exception_only": []string{},
				"bob":               []string{"bin1", "bin2", "lib1", "lib2"},
			},
			byOrigin: map[string][]string{
				"bin1":  []string{"bob"},
				"bin2":  []string{"bob"},
				"lib1":  []string{"bob"},
				"lib2":  []string{"bob"},
				"other": []string{},
			},
		},
		{
			name:             "everythingcontainer",
			treatAsAggregate: true,
			conditions: map[string][]string{
				"notice":            []string{"bin1", "bin2", "lib1", "lib2"},
				"reciprocal":        []string{"bin1", "bin2", "lib1", "lib2"},
				"restricted":        []string{"bin1", "bin2", "lib1", "lib2"},
				"by_exception_only": []string{"bin1", "bin2", "lib1", "lib2"},
			},
			newName:     "notice",
			filterNames: []ConditionNames{{"permissive", "notice", "reciprocal", "restricted", "by_exception_only", "proprietary"}},
			byName: map[string][]string{
				"permissive":        []string{},
				"notice":            []string{"bin1", "bin2", "lib1", "lib2"},
				"reciprocal":        []string{},
				"restricted":        []string{},
				"by_exception_only": []string{},
			},
			byOrigin: map[string][]string{
				"bin1":  []string{"notice"},
				"bin2":  []string{"notice"},
				"lib1":  []string{"notice"},
				"lib2":  []string{"notice"},
				"other": []string{},
			},
		},
		{
			name:             "everythingshared",
			treatAsAggregate: false,
			conditions: map[string][]string{
				"notice":            []string{"bin1", "bin2", "lib1", "lib2"},
				"reciprocal":        []string{"bin1", "bin2", "lib1", "lib2"},
				"restricted":        []string{"bin1", "bin2", "lib1", "lib2"},
				"by_exception_only": []string{"bin1", "bin2", "lib1", "lib2"},
			},
			newName:     "shared",
			filterNames: []ConditionNames{{"reciprocal", "restricted"}},
			byName: map[string][]string{
				"permissive":        []string{},
				"notice":            []string{},
				"reciprocal":        []string{},
				"restricted":        []string{},
				"by_exception_only": []string{},
				"shared":            []string{"bin1", "bin2", "lib1", "lib2"},
			},
			byOrigin: map[string][]string{
				"bin1":  []string{"shared"},
				"bin2":  []string{"shared"},
				"lib1":  []string{"shared"},
				"lib2":  []string{"shared"},
				"other": []string{},
			},
		},
		{
			name:             "oneeach",
			treatAsAggregate: true,
			conditions: map[string][]string{
				"notice":            []string{"bin1"},
				"reciprocal":        []string{"bin2"},
				"restricted":        []string{"lib1"},
				"by_exception_only": []string{"lib2"},
			},
			newName:     "other",
			filterNames: []ConditionNames{{"notice", "by_exception_only"}},
			byName: map[string][]string{
				"permissive":        []string{},
				"notice":            []string{},
				"reciprocal":        []string{},
				"restricted":        []string{},
				"by_exception_only": []string{},
				"other":             []string{"bin1", "lib2"},
			},
			byOrigin: map[string][]string{
				"bin1":  []string{"other"},
				"bin2":  []string{},
				"lib1":  []string{},
				"lib2":  []string{"other"},
				"other": []string{},
			},
		},
		{
			name:             "oneeachnoneleft",
			treatAsAggregate: true,
			conditions: map[string][]string{
				"notice":            []string{"bin1"},
				"reciprocal":        []string{"bin2"},
				"restricted":        []string{"lib1"},
				"by_exception_only": []string{"lib2"},
			},
			newName:     "notice",
			filterNames: []ConditionNames{{"unencumbered", "permissive"}},
			byName: map[string][]string{
				"permissive":        []string{},
				"notice":            []string{},
				"reciprocal":        []string{},
				"restricted":        []string{},
				"by_exception_only": []string{},
			},
			byOrigin: map[string][]string{
				"bin1":  []string{},
				"bin2":  []string{},
				"lib1":  []string{},
				"lib2":  []string{},
				"other": []string{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			lg := newLicenseGraphImp()
			lsimp := newLicenseConditionSet(&targetNodeImp{lg, ""})
			for name, origins := range tt.conditions {
				for _, origin := range origins {
					_ = newTestNode(lg, origin)
					lsimp.addAll(origin, name)
				}
			}
			lsimp = lsimp.rename(tt.treatAsAggregate, tt.newName, tt.filterNames...)
			checkAsAggregate(tt.treatAsAggregate, lsimp, t)
			byName(tt.byName).checkPublic(lsimp, t)
			byOrigin(tt.byOrigin).checkPublic(lg, lsimp, t)
		})
	}
}
