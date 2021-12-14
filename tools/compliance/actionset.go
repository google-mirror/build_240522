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
)

type resolutionAction struct {
	actsOn *TargetNode
	cs LicenseConditionSet
}

func (a resolutionAction) IsEqual(other resolutionAction) bool {
	return a.actsOn == other.actsOn && a.cs == other.cs
}

func (a resolutionAction) IsLess(other resolutionAction) bool {
	if a.actsOn == other.actsOn {
		return a.cs < other.cs
	}
	aName := a.actsOn.Name()
	oName := other.actsOn.Name()
	if aName == oName {
		panic(fmt.Errorf("identically named targets at different locations: %s and %s", a.actsOn, other.actsOn))
	}
	return aName < oName
}

func (a resolutionAction) String() string {
	return fmt.Sprintf("%s{%s}", a.actsOn.Name(), strings.Join(a.cs.Names(), ", "))
}

type resolution struct {
	attachesTo *TargetNode
	action resolutionAction
}

func (r resolution) IsEqual(other resolution) bool {
	return r.attachesTo == other.attachesTo && r.action.IsEqual(other.action)
}

func (r resolution) IsLess(other resolution) bool {
	if r.attachesTo == other.attachesTo {
		return r.action.IsLess(other.action)
	}
	rName := r.attachesTo.Name()
	oName := other.attachesTo.Name()
	if rName == oName {
		panic(fmt.Errorf("identically named targets at different locations: %s and %s", r.attachesTo, other.attachesTo))
	}
	return rName < oName
}

func (r resolution) String() string {
	return fmt.Sprintf("%s:%s", r.attachesTo.Name(), r.action.String())
}


// actionSet maps `actOn` target nodes to the license conditions the action on each target resolves.
type actionSet struct {
	lg *LicenseGraph
	indexes *IntervalSet
}

// String returns a string representation of the set.
func (as actionSet) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "{")

	sep := ""
	as.indexes.VisitAll(func(index int) {
		fmt.Fprintf(&sb, "%s%s", sep, as.lg.actions[index].String())
		sep = ", "
	})
	fmt.Fprintf(&sb, "}")
	return sb.String()
}

func (as actionSet) FindFirst(a *resolutionAction, matches func(resolutionAction) bool) bool {
	var index int
	if !as.indexes.FindFirst(&index, func(i int) bool { return matches(as.lg.actions[i]) }) {
		return false
	}
	*a = as.lg.actions[index]
	return true
}

func (as actionSet) VisitAll(visit func(a resolutionAction)) {
	as.indexes.VisitAll(func(index int) { visit(as.lg.actions[index]) })
}

func (as actionSet) Len() int {
	return as.indexes.Len()
}

func (as actionSet) matchingAny(condition ...LicenseCondition) chan resolutionAction {
	result := make(chan resolutionAction)

	go func() {
		as.indexes.VisitAll(func(index int) {
			a := as.lg.actions[index]
			cs := a.cs
			matching := cs.MatchingAny(condition...)
			if !matching.IsEmpty() {
				result <- resolutionAction{a.actsOn, matching}
			}
		})
		close(result)
	}()

	return result
}

func (as actionSet) matchingAnySet(conditionSet ...LicenseConditionSet) chan resolutionAction {
	result := make(chan resolutionAction)

	go func() {
		as.indexes.VisitAll(func(index int) {
			a := as.lg.actions[index]
			cs := a.cs
			matching := cs.MatchingAnySet(conditionSet...)
			if !matching.IsEmpty() {
				result <- resolutionAction{a.actsOn, matching}
			}
		})
		close(result)
	}()

	return result
}

func (as actionSet) hasAny(condition ...LicenseCondition) bool {
	var index int
	return as.indexes.FindFirst(&index, func(i int) bool {
		return as.lg.actions[i].cs.HasAny(condition...)
	})
}

func (as actionSet) matchesAnySet(conditions ...LicenseConditionSet) bool {
	var index int
	return as.indexes.FindFirst(&index, func(i int) bool {
		return as.lg.actions[i].cs.MatchesAnySet(conditions...)
	})
}

// byActsOn returns the subset of `as` where `actsOn` is in the `reachable` target node set.
func (as actionSet) byActsOn(reachable *TargetNodeSet) actionSet {
	result := actionSet{as.lg, &IntervalSet{}}

	as.indexes.VisitAll(func(index int) {
		if reachable.Contains(as.lg.actions[index].actsOn) {
			result.indexes.Insert(index)
		}
	})
	return result
}

// copy returns another actionSet with the same value as `as`
func (as actionSet) copy() actionSet {
	var indexes IntervalSet
	as.indexes.Copy(&indexes)
	return actionSet{as.lg, &indexes}
}

// addSet adds all of the actions of `other` if not already present.
func (as actionSet) addSet(other actionSet) {
	other.indexes.VisitAll(func(index int) {
		as.add(other.lg.actions[index].actsOn, other.lg.actions[index].cs)
	})
}

// addAll adds all of the actions fro channel `c` if not already present.
func (as actionSet) addAll(c chan resolutionAction) {
	for a := range c {
		as.add(a.actsOn, a.cs)
	}
}

// add makes the action on `actsOn` to resolve the conditions in `cs` a member of the set.
func (as actionSet) add(actsOn *TargetNode, cs LicenseConditionSet) {
	index := as.actionIndex(as.lg, actsOn, cs)
	as.indexes.Insert(index)
}

// addCondition makes the action on `actsOn` to resolve `lc` a member of the set.
func (as actionSet) addCondition(actsOn *TargetNode, lc LicenseCondition) {
	index := as.actionIndex(as.lg, actsOn, LicenseConditionSet(lc))
	as.indexes.Insert(index)
}

// isEmpty returns true if no action to resolve a condition exists.
func (as actionSet) isEmpty() bool {
	return as.indexes.IsEmpty()
}

// isEqual returns true if `other` contains the same set as `as`.
func (as actionSet) isEqual(other actionSet) bool {
	return as.lg == other.lg && as.indexes.IsEqual(other.indexes)
}

// conditions returns the set of conditions resolved by the action set.
func (as actionSet) conditions() LicenseConditionSet {
	var result LicenseConditionSet
	as.indexes.VisitAll(func(index int) {
		result = result.Union(as.lg.actions[index].cs)
	})
	return result
}

func (as actionSet) actionIndex(lg *LicenseGraph, actsOn *TargetNode, cs LicenseConditionSet) int {
	var index int
	if as.indexes.FindFirst(&index, func(i int) bool { return lg.actions[i].actsOn == actsOn }) {
		if cs == (cs & lg.actions[index].cs) {
			return index
		}
		cs |= lg.actions[index].cs
		as.indexes.Remove(index)
	}

	if !actsOn.actions.FindFirst(&index, func(i int) bool {
		if lg.actions[i].actsOn != actsOn {
			panic(fmt.Errorf("action for wrong target: got %s, want %s", lg.actions[i].actsOn, actsOn))
		}
		return lg.actions[i].cs == cs
	}) {
		lg.mu.Lock()
		index = len(lg.actions)
		lg.actions = append(lg.actions, resolutionAction{actsOn, cs})
		actsOn.actions.Insert(index)
		lg.mu.Unlock()
	}
	return index
}
