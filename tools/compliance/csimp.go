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
	cs.addAll(conditions)
	return cs
}

// conditionNamesArray implements a `contains` predicate for arrays of ConditionNames
type conditionNamesArray []ConditionNames

func (cn conditionNamesArray) contains(name string) bool {
	for _, c := range cn {
		if c.Contains(name) {
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
	// by mapping condition name -> origin target name -> nil
	conditions map[string]map[string]interface{}
}

// Add makes all `conditions` members of the set if they were not previously.
func (cs *licenseConditionSetImp) Add(conditions... LicenseCondition) {
	cs.addAll(conditions)
}

// ByName returns a list of the conditions in the set matching `name`.
func (cs *licenseConditionSetImp) ByName(name... ConditionNames) []LicenseCondition {
	size := 0
	for _, cn := range name {
		for _, n := range cn {
			if origins, ok := cs.conditions[n]; ok {
				size += len(origins)
			}
		}
	}
	l := make([]LicenseCondition, 0, size)
	for _, cn := range name {
		for _, n := range cn {
			if origins, ok := cs.conditions[n]; ok {
				for t := range origins {
					l = append(l, licenseConditionImp{n, targetNodeImp{cs.lg, t}})
				}
			}
		}
	}
	return l
}

// ByOrigin returns all of the conditions that originate at `origin` regardless of name.
func (cs *licenseConditionSetImp) ByOrigin(origin TargetNode) []LicenseCondition {
	oimp := origin.(targetNodeImp)
	size := 0
	for _, origins := range cs.conditions {
		if _, ok := origins[oimp.file]; ok {
			size++
		}
	}
	l := make([]LicenseCondition, 0, size)
	for name, origins := range cs.conditions {
		if _, ok := origins[oimp.file]; ok {
			l = append(l, licenseConditionImp{name, targetNodeImp{cs.lg, oimp.file}})
		}
	}
	return l
}

// HasAnyByOrigin returns true if the set contains any conditions originating at `origin` regarless of name.
func (cs *licenseConditionSetImp) HasAnyByOrigin(origin TargetNode) bool {
	oimp := origin.(targetNodeImp)
	for _, origins := range cs.conditions {
		if _, ok := origins[oimp.file]; ok {
			return true
		}
	}
	return false
}

// Conditions returns a list of all the conditions in the set.
func (cs *licenseConditionSetImp) Conditions() []LicenseCondition {
	size := 0
	for _, origins := range cs.conditions {
		size += len(origins)
	}
	l := make([]LicenseCondition, 0, size)
	for c, origins := range cs.conditions {
		for t := range origins {
			l = append(l, licenseConditionImp{c, targetNodeImp{cs.lg, t}})
		}
	}
	return l
}

// Copy creates a new LicenseCondition variable with the same value.
func (cs *licenseConditionSetImp) Copy() LicenseConditionSet {
	other := licenseConditionSetImp{cs.lg, make(map[string]map[string]interface{})}
	for name, origins := range cs.conditions {
		other.conditions[name] = make(map[string]interface{})
		for origin := range origins {
			other.conditions[name][origin] = nil
		}
	}
	return &other
}

// HasAnyByName returns true if the set contains any conditions matching `name` originating at any target.
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

// IsEmpty returns true when the sent contains zero conditions.
func (cs *licenseConditionSetImp) IsEmpty() bool {
	for _, origins := range cs.conditions {
		if len(origins) > 0 {
			return false
		}
	}
	return true
}

// compliance-only licenseConditionSetImp methods

// newLicenseConditionSet constructs a set of conditions
func newLicenseConditionSet(origin *targetNodeImp, condition... string) *licenseConditionSetImp {
	cs := &licenseConditionSetImp{nil, make(map[string]map[string]interface{})}
	if origin != nil {
		cs.lg = origin.lg
	} else if len(condition) > 0 {
		panic(fmt.Errorf("attempt to add conditions to nil target"))
	}
	for _, name := range condition {
		cs.conditions[name][origin.file] = nil
	}
	return cs
}

// add changes the set to include `condition` if it does not already
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

// addAll modifies `cs` to include all of the `conditions` if they were not previously members.
func (cs *licenseConditionSetImp) addAll(conditions []LicenseCondition) {
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

// removeAllByName changes the set to delete all conditions matching `names`.
func (cs *licenseConditionSetImp) removeAllByName(names... ConditionNames) {
	for _, cn := range names {
		for _, name := range cn {
			delete(cs.conditions, name)
		}
	}
}

// removeAllByTarget changes the set to delete all conditions that originate at target `file`.
func (cs *licenseConditionSetImp) removeAllByOrigin(file string) {
	for name := range cs.conditions {
		delete(cs.conditions[name], file)
	}
}

// remove changes the set to delete all conditions also present in `other`.
func (cs *licenseConditionSetImp) remove(other *licenseConditionSetImp) {
	for name, origins := range other.conditions {
		if _, isPresent := cs.conditions[name]; !isPresent {
			continue
		}
		for origin := range origins {
			delete(cs.conditions[name], origin)
		}
	}
}

// rename behaves similar to filter except it limits the search to a single ConditionNames, and it
// changes the name of each condition in the output to `newName`.
func (cs *licenseConditionSetImp) rename(names ConditionNames, newName string) *licenseConditionSetImp {
	result := &licenseConditionSetImp{cs.lg, make(map[string]map[string]interface{})}
	for name, origins := range cs.conditions {
		if !names.Contains(name) {
			continue
		}
		for origin := range origins {
			if _, ok := result.conditions[name]; !ok {
				result.conditions[newName] = make(map[string]interface{})
			}
			result.conditions[newName][origin] = nil
		}
	}
	return result
}

// union returns a new set calculated as the union of `cs` with some `other` set.
func (cs *licenseConditionSetImp) union(other *licenseConditionSetImp) *licenseConditionSetImp {
	if cs.lg == nil {
		cs.lg = other.lg
	} else if other.lg == nil {
		other.lg = cs.lg
	} else if cs.lg != other.lg {
		panic(fmt.Errorf("attempt to union condition sets from different graphs"))
	}
	result := &licenseConditionSetImp{cs.lg, make(map[string]map[string]interface{})}
	for name, origins := range cs.conditions {
		for origin := range origins {
			if _, ok := result.conditions[name]; !ok {
				result.conditions[name] = make(map[string]interface{})
			}
			result.conditions[name][origin] = nil
		}
	}
	for name, origins := range other.conditions {
		for origin := range origins {
			if _, ok := result.conditions[name]; !ok {
				result.conditions[name] = make(map[string]interface{})
			}
			result.conditions[name][origin] = nil
		}
	}
	return result
}


// filter returns a new liceneseConditionSetImp containing the subset of conditions matching `names`.
func (cs *licenseConditionSetImp) filter(names... ConditionNames) *licenseConditionSetImp {
	result := &licenseConditionSetImp{cs.lg, make(map[string]map[string]interface{})}
	for name, origins := range cs.conditions {
		if !conditionNamesArray(names).contains(name) {
			continue
		}
		for origin := range origins {
			if _, ok := result.conditions[name]; !ok {
				result.conditions[name] = make(map[string]interface{})
			}
			result.conditions[name][origin] = nil
		}
	}
	return result
}
