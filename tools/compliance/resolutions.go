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

// ResolutionSet describes an immutable set of targets and the license
// conditions each target must satisfy or "resolve" in a specific context.
//
// Ultimately, the purpose of recording the license metadata and building a
// license graph is to identify, describe, and verify the necessary actions or
// operations for compliance policy.
//
// i.e. What is the source-sharing policy? Has it been met? Meet it.
//
// i.e. Are there incompatible policy requirements? Such as a source-sharing
// policy applied to code that policy also says may not be shared? If so, stop
// and remove the dependencies that create the situation.
//
// The ResolutionSet is the base unit for mapping license conditions to the
// targets triggering some necessary action per policy. Different ResolutionSet
// values may be calculated for different contexts.
//
// e.g. Suppose an unencumbered binary links in a notice .a library.
//
// An "unencumbered" condition would originate from the binary, and a "notice"
// condition would originate from the .a library. A ResolutionSet for the
// context of the Notice policy might apply both conditions to the binary while
// preserving the origin of each condition. By applying the notice condition to
// the binary, the ResolutionSet stipulates the policy that the release of the
// unencumbered binary must provide suitable notice for the .a library.
//
// The resulting ResolutionSet could be used for building a notice file, for
// validating that a suitable notice has been built into the distribution, or
// for reporting what notices need to be given.
//
// See also: resolve.go for:
//  * ResolveBottomUpConditions(...)
//  * ResolveTopDownCondition(...)
//  * ResolveTopDownRestricted(...)
type ResolutionSet interface {
	// AppliesTo identifies the list of targets requiring action to resolve
	// conditions (unordered).
	AppliesTo() []TargetNode

	// Conditions returns the set of license conditions that `appliesTo`
	// target must resolve, and whether the `appliesTo` target appears in
	// the set.
	Conditions(appliesTo TargetNode) (LicenseConditionSet, bool)

	// AppliesToTarget returns true if the set contains conditions that
	// apply to `appliesTo`.
	AppliesToTarget(appliesTo TargetNode) bool

	// AnyByNameApplyToTarget returns true if the set contains conditions
	// matching `names` that apply to `appliesTo`.
	AnyByNameApplyToTarget(appliesTo TargetNode, names... ConditionNames) bool

	// AllByNameApplyTo returns true if the set contains at least one
	// condition matching each element of `names` for `appliesTo`.
	AllByNameApplyToTarget(appliesTo TargetNode, names... ConditionNames) bool
}


// LicenseCondition describes an individual license condition or requirement.
//
// Origin identifies the origin of the requirement and gets .
type LicenseCondition interface {
	// Name identifies the type of condition. e.g. "notice" or "restricted"
	//
	// Policy prescribes the set of condition names and how each must be
	// resolved.
	Name() string

	// Origin identifies the target node that originated the condition.
	Origin() TargetNode
}

// LicenseConditionSet describes a set of license conditions.
//
// Use:
//   func NewLicenseConditionSet(conditions... LicenseCondition) LicenseConditionSet
// to construct LicenseConditionSet variables for local use.
type LicenseConditionSet interface {
	// Add makes each `condition` a member of the set.
	Add(condition... LicenseCondition)

	// ByName returns all of the conditions in the set matching `names`.
	ByName(names... ConditionNames) []LicenseCondition
	// HasAnyByName returns true if the set contains any conditions
	// matching `names`.
	HasAnyByName(name... ConditionNames) bool

	// ByOrigin returns all of the conditions that originate at `origin`.
	ByOrigin(origin TargetNode) []LicenseCondition
	// HasAnyByOrigin returns true if the set contains any conditions
	// originating at `origin`.
	HasAnyByOrigin(origin TargetNode) bool

	// Conditions returns all of the license conditions in the set.
	// (unordered)
	Conditions() []LicenseCondition

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
