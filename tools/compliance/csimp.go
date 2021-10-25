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
)

// NewLicenseConditionSet creates a new instance or variable of LicenseConditionSet.
func NewLicenseConditionSet(conditions... LicenseCondition) LicenseConditionSet {
	cs := &licenseConditionSetImp{nil, make(map[string]map[string]interface{})}
	cs.addList(conditions)
	return cs
}

// conditionNamesArray implements a `contains` predicate for arrays of ConditionNames
type conditionNamesArray []ConditionNames

func (cn conditionNamesArray) contains(name string) bool {
	for _, names := range cn {
		if names.Contains(name) {
			return true
		}
	}
	return false
}


// licenseConditionSetImp implements publicly mutable LicenseConditionSet.
type licenseConditionSetImp struct {
	// lg identifies the graph to which the condition applies
	lg *licenseGraphImp

	// conditions describes the set of license conditions i.e. (condition name, origin target name) pairs
	// by mapping condition name -> origin target name -> nil.
	//
	// In very limited contexts during resolve walks, the final mapping changes to target name -> true|false
	// for whether to treat the origin target as a purely aggregating container.
	conditions map[string]map[string]interface{}
}

// Add makes all `conditions` members of the set if they were not previously.
func (cs *licenseConditionSetImp) Add(conditions... LicenseCondition) {
	cs.addList(conditions)
}

// ByName returns a list of the conditions in the set matching `names`.
func (cs *licenseConditionSetImp) ByName(names... ConditionNames) ConditionList {
	l := make(ConditionList, 0, cs.CountByName(names...))
	for _, cn := range names {
		for _, name := range cn {
			if origins, ok := cs.conditions[name]; ok {
				for t := range origins {
					l = append(l, licenseConditionImp{name, targetNodeImp{cs.lg, t}})
				}
			}
		}
	}
	return l
}

// HasAnyByName returns true if the set contains any conditions matching `names` originating at any target.
func (cs *licenseConditionSetImp) HasAnyByName(names... ConditionNames) bool {
	for _, cn := range names {
		for _, name := range cn {
			if origins, ok := cs.conditions[name]; ok {
				if len(origins) > 0 {
					return true
				}
			}
		}
	}
	return false
}

// CountByName returns the number of conditions matching `names` originating at any target.
func (cs *licenseConditionSetImp) CountByName(names... ConditionNames) int {
	size := 0
	for _, cn := range names {
		for _, name := range cn {
			if origins, ok := cs.conditions[name]; ok {
				size += len(origins)
			}
		}
	}
	return size
}

// ByOrigin returns all of the conditions that originate at `origin` regardless of name.
func (cs *licenseConditionSetImp) ByOrigin(origin TargetNode) ConditionList {
	oimp := origin.(targetNodeImp)
	return cs.byOrigin(oimp.file)
}

// HasAnyByOrigin returns true if the set contains any conditions originating at `origin` regardless of condition name.
func (cs *licenseConditionSetImp) HasAnyByOrigin(origin TargetNode) bool {
	oimp := origin.(targetNodeImp)
	return cs.hasByOrigin(oimp.file)
	return false
}

// CountByOrigin returns the number of conditions originating at `origin` regardless of condition name.
func (cs *licenseConditionSetImp) CountByOrigin(origin TargetNode) int {
	oimp := origin.(targetNodeImp)
	return cs.countByOrigin(oimp.file)
}

// AsList returns a list of all the conditions in the set.
func (cs *licenseConditionSetImp) AsList() ConditionList {
	l := make([]LicenseCondition, 0, cs.Count())
	for c, origins := range cs.conditions {
		for t := range origins {
			l = append(l, licenseConditionImp{c, targetNodeImp{cs.lg, t}})
		}
	}
	return l
}

// Count returns the number of conditions in the set.
func (cs *licenseConditionSetImp) Count() int {
	size := 0
	for _, origins := range cs.conditions {
		size += len(origins)
	}
	return size
}

// Copy creates a new LicenseCondition variable with the same value.
func (cs *licenseConditionSetImp) Copy() LicenseConditionSet {
	other := licenseConditionSetImp{cs.lg, make(map[string]map[string]interface{})}
	for name := range cs.conditions {
		other.conditions[name] = make(map[string]interface{})
		for origin := range cs.conditions[name] {
			other.conditions[name][origin] = cs.conditions[name][origin]
		}
	}
	return &other
}

// HasCondition returns true if the set contains any condition matching both `names` and `origin`.
func (cs *licenseConditionSetImp) HasCondition(names ConditionNames, origin TargetNode) bool {
	oimp := origin.(targetNodeImp)
	if cs.lg == nil {
		return false
	} else if oimp.lg == nil {
		oimp.lg = cs.lg
	} else if cs.lg != oimp.lg {
		panic(fmt.Errorf("attempt to query license conditions from different graph"))
	}
	for _, name := range names {
		if origins, ok := cs.conditions[name]; ok {
			_, isPresent := origins[oimp.file]
			if isPresent {
				return true
			}
		}
	}
	return false
}

// IsEmpty returns true when the set conditions contains zero elements.
func (cs *licenseConditionSetImp) IsEmpty() bool {
	for _, origins := range cs.conditions {
		if len(origins) > 0 {
			return false
		}
	}
	return true
}

// compliance-only licenseConditionSetImp methods

// newLicenseConditionSet constructs a set of `conditions`.
func newLicenseConditionSet(origin *targetNodeImp, conditions... string) *licenseConditionSetImp {
	cs := &licenseConditionSetImp{nil, make(map[string]map[string]interface{})}
	if origin != nil {
		cs.lg = origin.lg
	}
	cs.addAll(origin.file, conditions...)
	return cs
}

// addAll changes the set to include each element of `conditions` originating at `origin`.
func (cs *licenseConditionSetImp) addAll(origin string, conditions... string) {
	if cs.lg == nil && len(conditions) > 0 {
		panic(fmt.Errorf("attempt to add conditions to nil target"))
	}
	for _, name := range conditions {
		if _, ok := cs.conditions[name]; !ok {
			cs.conditions[name] = make(map[string]interface{})
		}
		cs.conditions[name][origin] = nil
	}
}

// add changes the set to include `condition` for each element of `origins`.
func (cs *licenseConditionSetImp) add(condition string, origins... targetNodeImp) {
	if len(origins) == 0 {
		return
	}
	if _, ok := cs.conditions[condition]; !ok {
		cs.conditions[condition] = make(map[string]interface{})
	}

	for _, origin := range origins {
		if cs.lg == nil {
			cs.lg = origin.lg
		} else if origin.lg == nil {
			origin.lg = cs.lg
		} else if origin.lg != cs.lg {
			panic(fmt.Errorf("attempting to combine license conditions from different graphs"))
		}
		found := false
		for otherorigin := range cs.conditions[condition] {
			if origin.file == otherorigin {
				found = true
				break
			}
		}
		if !found {
			cs.conditions[condition][origin.file] = nil
		}
	}
}

// addList modifies `cs` to include all of the `conditions` if they were not previously members.
func (cs *licenseConditionSetImp) addList(conditions []LicenseCondition) {
	if len(conditions) == 0 {
		return
	}
	for _, c := range conditions {
		cimp := c.(licenseConditionImp)
		if cs.lg == nil {
			cs.lg = cimp.origin.lg
		} else if cimp.origin.lg != cs.lg {
			panic(fmt.Errorf("attempting to combine license conditions from different graphs"))
		}
		if _, ok := cs.conditions[cimp.name]; !ok {
			cs.conditions[cimp.name] = make(map[string]interface{})
		}
		cs.conditions[cimp.name][cimp.origin.file] = nil
	}
}

// addSet modifies `cs` to include all of the elements of `other` if they were not previously members.
func (cs *licenseConditionSetImp) addSet(other *licenseConditionSetImp) {
	if len(other.conditions) == 0 {
		return
	}
	if cs.lg == nil {
		cs.lg = other.lg
	} else if other.lg != cs.lg {
		panic(fmt.Errorf("attempting to add set of license conditions from different graph"))
	}
	for name, origins := range other.conditions {
		if len(origins) == 0 {
			continue
		}
		if _, ok := cs.conditions[name]; !ok {
			cs.conditions[name] = make(map[string]interface{})
		}
		for origin := range origins {
			cs.conditions[name][origin] = other.conditions[name][origin]
		}
	}
}

// removeAllByName changes the set to delete all conditions matching `names`.
func (cs *licenseConditionSetImp) removeAllByName(names... ConditionNames) {
	for _, cn := range names {
		for _, name := range cn {
			delete(cs.conditions, name)
		}
	}
}

// removeAllByOrigin changes the set to delete all conditions that originate at target `file`.
func (cs *licenseConditionSetImp) removeAllByOrigin(file string) {
	for name := range cs.conditions {
		delete(cs.conditions[name], file)
	}
}

// removeSet changes the set to delete all conditions also present in `other`.
func (cs *licenseConditionSetImp) removeSet(other *licenseConditionSetImp) {
	for name, origins := range other.conditions {
		if _, isPresent := cs.conditions[name]; !isPresent {
			continue
		}
		for origin := range origins {
			delete(cs.conditions[name], origin)
		}
	}
}

// byOrigin returns all of the conditions that originate at a target named `origin` regardless of condition name.
func (cs *licenseConditionSetImp) byOrigin(origin string) []LicenseCondition {
	size := cs.countByOrigin(origin)
	l := make([]LicenseCondition, 0, size)
	for name, origins := range cs.conditions {
		if _, ok := origins[origin]; ok {
			l = append(l, licenseConditionImp{name, targetNodeImp{cs.lg, origin}})
		}
	}
	return l
}

// hasByOrigin returns true if any of the conditions originate at a target named `origin` regardless of condition name.
func (cs *licenseConditionSetImp) hasByOrigin(origin string) bool {
	for _, origins := range cs.conditions {
		if _, ok := origins[origin]; ok {
			return true
		}
	}
	return false
}

// countByOrigin returns the number of conditions that originate at a target named `origin` regardless of condition name.
func (cs *licenseConditionSetImp) countByOrigin(origin string) int {
	size := 0
	for _, origins := range cs.conditions {
		if _, ok := origins[origin]; ok {
			size++
		}
	}
	return size
}

// filter returns a new liceneseConditionSetImp containing the subset of conditions matching `names`.
//
// `treatAsAggregate` is used as the map existence payload instead of nil
func (cs *licenseConditionSetImp) filter(treatAsAggregate bool, names... ConditionNames) *licenseConditionSetImp {
	result := &licenseConditionSetImp{cs.lg, make(map[string]map[string]interface{})}
	for name, origins := range cs.conditions {
		if !conditionNamesArray(names).contains(name) {
			continue
		}
		for origin := range origins {
			if _, ok := result.conditions[name]; !ok {
				result.conditions[name] = make(map[string]interface{})
			}
			result.conditions[name][origin] = treatAsAggregate
		}
	}
	return result
}

// rename behaves similar to filter except it changes the name of each condition in the output to `newName`.
//
// Used by policy.go for policies like the notice policy, which applies to most conditions--not just to the "notice" condition.
func (cs *licenseConditionSetImp) rename(treatAsAggregate bool, newName string, names... ConditionNames) *licenseConditionSetImp {
	result := &licenseConditionSetImp{cs.lg, make(map[string]map[string]interface{})}
	for name, origins := range cs.conditions {
		if !conditionNamesArray(names).contains(name) {
			continue
		}
		for origin := range origins {
			if _, ok := result.conditions[newName]; !ok {
				result.conditions[newName] = make(map[string]interface{})
			}
			result.conditions[newName][origin] = treatAsAggregate
		}
	}
	return result
}

// copy returns a new liceneseConditionSetImp containing all of the conditions in `cs` except
// using `treatAsAggregate` instead of nil as the existence payload.
func (cs *licenseConditionSetImp) copy(treatAsAggregate bool) *licenseConditionSetImp {
	result := &licenseConditionSetImp{cs.lg, make(map[string]map[string]interface{})}
	for name, origins := range cs.conditions {
		for origin := range origins {
			if _, ok := result.conditions[name]; !ok {
				result.conditions[name] = make(map[string]interface{})
			}
			result.conditions[name][origin] = treatAsAggregate
		}
	}
	return result
}
