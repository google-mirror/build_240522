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

// LicenseCondition describes an individual license condition or requirement
// attached or originating at a specific target node. (immutable)
type LicenseCondition interface {
	// Name identifies the type of condition. e.g. "notice" or "restricted"
	//
	// Policy prescribes the set of condition names and how each must be
	// resolved.
	Name() string

	// Origin identifies the target node that originated the condition.
	//
	// e.g. If a library with a notice policy is linked into a binary, the
	// notice condition originates at the library, and gets resolved by the
	// binary when the binary provides the necessary notice.
	Origin() TargetNode
}


// ConditionList implements introspection methods to arrays of LicenseCondition.
type ConditionList []LicenseCondition

// HasByName returns true if the list contains any condition matching `name`.
func (cl ConditionList) HasByName(name ConditionNames) bool {
	for _, c := range cl {
		if name.Contains(c.Name()) {
			return true
		}
	}
	return false
}

// ByName returns the sublist of conditions that match `name`.
func (cl ConditionList) ByName(name ConditionNames) ConditionList {
	result := make(ConditionList, 0, cl.CountByName(name))
	for _, c := range cl {
		if name.Contains(c.Name()) {
			result = append(result, c)
		}
	}
	return result
}

// CountByName returns the size of the sublist of conditions that match `name`.
func (cl ConditionList) CountByName(name ConditionNames) int {
	size := 0
	for _, c := range cl {
		if name.Contains(c.Name()) {
			size++
		}
	}
	return size
}

// HasByOrigin returns true if the list contains any condition originating at `origin`.
func (cl ConditionList) HasByOrigin(origin TargetNode) bool {
	for _, c := range cl {
		if c.Origin().Name() == origin.Name() {
			return true
		}
	}
	return false
}

// ByOrigin returns the sublist of conditions that originate at `origin`.
func (cl ConditionList) ByOrigin(origin TargetNode) ConditionList {
	result := make(ConditionList, 0, cl.CountByOrigin(origin))
	for _, c := range cl {
		if c.Origin().Name() == origin.Name() {
			result = append(result, c)
		}
	}
	return result
}

// CountByOrigin returns the size of the sublist of conditions that originate at `origin`.
func (cl ConditionList) CountByOrigin(origin TargetNode) int {
	size := 0
	for _, c := range cl {
		if c.Origin().Name() == origin.Name() {
			size++
		}
	}
	return size
}


// LicenseConditionSet describes a mutable set of immutable license conditions.
//
// Use:
//   func NewLicenseConditionSet(conditions... LicenseCondition) LicenseConditionSet
// to construct LicenseConditionSet variables for local use.
type LicenseConditionSet interface {
	// Add makes each `condition` a member of the set.
	Add(condition... LicenseCondition)

	// ByName returns all of the conditions in the set matching `names`.
	ByName(names... ConditionNames) ConditionList
	// HasAnyByName returns true if the set contains any conditions
	// matching `names`.
	HasAnyByName(names... ConditionNames) bool
	// CountByName returns the number of conditions in the set matching `names`.
	CountByName(names... ConditionNames) int

	// ByOrigin returns all of the conditions that originate at `origin`.
	ByOrigin(origin TargetNode) ConditionList
	// HasAnyByOrigin returns true if the set contains any conditions
	// originating at `origin`.
	HasAnyByOrigin(origin TargetNode) bool
	// CountByOrigin returns the number of conditions in the set that originate at `origin`.
	CountByOrigin(origin TargetNode) int

	// AsList returns all of the license conditions in the set.
	// (unordered)
	AsList() ConditionList
	// Count returns the number of license conditions in the set.
	Count() int

	// Copy makes a duplicate of the set that can be changed independently.
	Copy() LicenseConditionSet

	// IsEmpty returns true if the set contains zero elements.
	IsEmpty() bool
}


// ConditionNames implements the Contains predicate for slices of condition
// name strings.
type ConditionNames []string

// Contains returns true if the name matches one of the ConditionNames.
func (cn ConditionNames) Contains(name string) bool {
        for _, c := range cn {
                if c == name {
                        return true
                }
        }
        return false
}
