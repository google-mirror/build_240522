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

// NewLicenseConditionSet creates a new instance or variable of *LicenseConditionSet.
func NewLicenseConditionSet(conditions ...LicenseCondition) *LicenseConditionSet {
	cs := &LicenseConditionSet{nil, make(map[string]map[string]interface{})}
	cs.Add(conditions...)
	return cs
}

// LicenseConditionSet describes a mutable set of immutable license conditions.
type LicenseConditionSet struct {
	// lg identifies the graph to which the condition applies
	lg *LicenseGraph

	// conditions describes the set of license conditions i.e. (condition name, origin target name) pairs
	// by mapping condition name -> origin target name -> nil.
	//
	// In very limited contexts during resolve walks, the final mapping changes to target name -> true|false
	// for whether to treat the origin target as a purely aggregating container.
	conditions map[string]map[string]interface{}
}

// Add makes all `conditions` members of the set if they were not previously.
func (cs *LicenseConditionSet) Add(conditions ...LicenseCondition) {
	if len(conditions) == 0 {
		return
	}
	for _, lc := range conditions {
		if cs.lg == nil {
			cs.lg = lc.origin.lg
		} else if lc.origin.lg != cs.lg {
			panic(fmt.Errorf("attempting to combine license conditions from different graphs"))
		}
		if _, ok := cs.conditions[lc.name]; !ok {
			cs.conditions[lc.name] = make(map[string]interface{})
		}
		cs.conditions[lc.name][lc.origin.file] = nil
	}
}

// AddSet makes all elements of `conditions` members of the set if they were not previously.
func (cs *LicenseConditionSet) AddSet(other *LicenseConditionSet) {
	if len(other.conditions) == 0 {
		return
	}
	if cs.lg == nil {
		cs.lg = other.lg
	} else if other.lg == nil {
		panic(fmt.Errorf("attempting to add %d conditions from nil graph: %s", len(other.conditions), strings.Join(other.asStringList(":"), ", ")))
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

// ByName returns a list of the conditions in the set matching `names`.
func (cs *LicenseConditionSet) ByName(names ...ConditionNames) *LicenseConditionSet {
	other := &LicenseConditionSet{cs.lg, make(map[string]map[string]interface{})}
	for _, cn := range names {
		for _, name := range cn {
			if origins, ok := cs.conditions[name]; ok {
				other.conditions[name] = make(map[string]interface{})
				for origin := range origins {
					other.conditions[name][origin] = nil
				}
			}
		}
	}
	return other
}

// HasAnyByName returns true if the set contains any conditions matching `names` originating at any target.
func (cs *LicenseConditionSet) HasAnyByName(names ...ConditionNames) bool {
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
func (cs *LicenseConditionSet) CountByName(names ...ConditionNames) int {
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
func (cs *LicenseConditionSet) ByOrigin(origin TargetNode) *LicenseConditionSet {
	other := &LicenseConditionSet{cs.lg, make(map[string]map[string]interface{})}
	for name, origins := range cs.conditions {
		if _, ok := origins[origin.file]; ok {
			other.conditions[name] = make(map[string]interface{})
			other.conditions[name][origin.file] = nil
		}
	}
	return other
}

// HasAnyByOrigin returns true if the set contains any conditions originating at `origin` regardless of condition name.
func (cs *LicenseConditionSet) HasAnyByOrigin(origin TargetNode) bool {
	return cs.hasByOrigin(origin.file)
}

// CountByOrigin returns the number of conditions originating at `origin` regardless of condition name.
func (cs *LicenseConditionSet) CountByOrigin(origin TargetNode) int {
	return cs.countByOrigin(origin.file)
}

// AsList returns a list of all the conditions in the set.
func (cs *LicenseConditionSet) AsList() ConditionList {
	result := make([]LicenseCondition, 0, cs.Count())
	for name, origins := range cs.conditions {
		for origin := range origins {
			result = append(result, LicenseCondition{name, TargetNode{cs.lg, origin}})
		}
	}
	return result
}

// Count returns the number of conditions in the set.
func (cs *LicenseConditionSet) Count() int {
	size := 0
	for _, origins := range cs.conditions {
		size += len(origins)
	}
	return size
}

// Copy creates a new LicenseCondition variable with the same value.
func (cs *LicenseConditionSet) Copy() *LicenseConditionSet {
	other := &LicenseConditionSet{cs.lg, make(map[string]map[string]interface{})}
	for name := range cs.conditions {
		other.conditions[name] = make(map[string]interface{})
		for origin := range cs.conditions[name] {
			other.conditions[name][origin] = cs.conditions[name][origin]
		}
	}
	return other
}

// HasCondition returns true if the set contains any condition matching both `names` and `origin`.
func (cs *LicenseConditionSet) HasCondition(names ConditionNames, origin TargetNode) bool {
	if cs.lg == nil {
		return false
	} else if origin.lg == nil {
		origin.lg = cs.lg
	} else if cs.lg != origin.lg {
		panic(fmt.Errorf("attempt to query license conditions from different graph"))
	}
	for _, name := range names {
		if origins, ok := cs.conditions[name]; ok {
			_, isPresent := origins[origin.file]
			if isPresent {
				return true
			}
		}
	}
	return false
}

// IsEmpty returns true when the set conditions contains zero elements.
func (cs *LicenseConditionSet) IsEmpty() bool {
	for _, origins := range cs.conditions {
		if len(origins) > 0 {
			return false
		}
	}
	return true
}

// RemoveAllByName changes the set to delete all conditions matching `names`.
func (cs *LicenseConditionSet) RemoveAllByName(names ...ConditionNames) {
	for _, cn := range names {
		for _, name := range cn {
			delete(cs.conditions, name)
		}
	}
}

// Remove changes the set to delete `conditions`.
func (cs *LicenseConditionSet) Remove(conditions ...LicenseCondition) {
	for _, lc := range conditions {
		if _, isPresent := cs.conditions[lc.name]; !isPresent {
			panic(fmt.Errorf("attempt to remove non-existent condition: %q", lc.asString(":")))
		}
		if _, isPresent := cs.conditions[lc.name][lc.origin.file]; !isPresent {
			panic(fmt.Errorf("attempt to remove non-existent origin: %q", lc.asString(":")))
		}
		delete(cs.conditions[lc.name], lc.origin.file)
	}
}

// removeSet changes the set to delete all conditions also present in `other`.
func (cs *LicenseConditionSet) RemoveSet(other *LicenseConditionSet) {
	for name, origins := range other.conditions {
		if _, isPresent := cs.conditions[name]; !isPresent {
			continue
		}
		for origin := range origins {
			delete(cs.conditions[name], origin)
		}
	}
}


// compliance-only LicenseConditionSet methods

// newLicenseConditionSet constructs a set of `conditions`.
func newLicenseConditionSet(lg *LicenseGraph) *LicenseConditionSet {
	return &LicenseConditionSet{lg, make(map[string]map[string]interface{})}
}

// addAll changes the set to include each element of `conditions` originating at `origin`.
func (cs *LicenseConditionSet) addAll(origin string, conditions ...string) {
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
func (cs *LicenseConditionSet) add(condition string, origins ...TargetNode) {
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

// byOrigin returns all of the conditions that originate at a target named `origin` regardless of condition name.
func (cs *LicenseConditionSet) byOrigin(origin string) []LicenseCondition {
	size := cs.countByOrigin(origin)
	result := make([]LicenseCondition, 0, size)
	for name, origins := range cs.conditions {
		if _, ok := origins[origin]; ok {
			result = append(result, LicenseCondition{name, TargetNode{cs.lg, origin}})
		}
	}
	return result
}

// hasByOrigin returns true if any of the conditions originate at a target named `origin` regardless of condition name.
func (cs *LicenseConditionSet) hasByOrigin(origin string) bool {
	for _, origins := range cs.conditions {
		if _, ok := origins[origin]; ok {
			return true
		}
	}
	return false
}

// countByOrigin returns the number of conditions that originate at a target named `origin` regardless of condition name.
func (cs *LicenseConditionSet) countByOrigin(origin string) int {
	size := 0
	for _, origins := range cs.conditions {
		if _, ok := origins[origin]; ok {
			size++
		}
	}
	return size
}

// removeAllByOrigin changes the set to delete all conditions that originate at target `file`.
func (cs *LicenseConditionSet) removeAllByOrigin(file string) {
	for name := range cs.conditions {
		delete(cs.conditions[name], file)
	}
}

// asStringList returns the conditions in the set as `separator`-separated pair strings.
func (cs *LicenseConditionSet) asStringList(separator string) []string {
	result := make([]string, 0, cs.Count())
	for name, origins := range cs.conditions {
		for origin := range origins {
			result = append(result, origin+separator+name)
		}
	}
	return result
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
