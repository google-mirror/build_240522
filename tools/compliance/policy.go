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

var (
	// ImpliesUnencumbered lists the condition names representing an author attempt to disclaim copyright.
	ImpliesUnencumbered = ConditionNames{"unencumbered"}

	// ImpliesPermissive lists the condition names representing copyrighted but "licensed without policy requirements".
	ImpliesPermissive = ConditionNames{"permissive"}

	// ImpliesNotice lists the condition names implying a notice or attribution policy.
	ImpliesNotice = ConditionNames{"unencumbered", "permissive", "notice", "reciprocal", "restricted"}

	// ImpliesReciprocal lists the condition names implying a local source-sharing policy.
	ImpliesReciprocal = ConditionNames{"reciprocal"}

	// Restricted lists the condition names implying an infectious source-sharing policy.
	ImpliesRestricted = ConditionNames{"restricted"}

	// ImpliesProprietary lists the condition names implying a confidentiality policy.
	ImpliesProprietary = ConditionNames{"proprietary"}

	// ImpliesByExceptionOnly lists the condition names implying a policy for "license review and approval before use".
	ImpliesByExceptionOnly = ConditionNames{"proprietary", "by_exception_only"}

	// ImpliesPrivate lists the condition names implying a source-code privacy policy.
	ImpliesPrivate = ConditionNames{"proprietary"}

	// ImpliesShared lists the condition names implying a source-code sharing policy.
	ImpliesShared = ConditionNames{"reciprocal", "restricted"}
)

// Resolution happens in two passes:
//
// 1. A bottom-up traversal propagates license conditions up to targets from
// dendencies as needed.
//
// 2. For each condition of interest, a top-down traversal adjusts the attached
// conditions pushing restricted down from targets into linked dependencies, or
// turning all manner of conditions into "notice" for the notice policy.
//
// The behavior of the 2 passes gets controlled by the 3 functions below.
//
// The first function controls what happens during the bottom-up traversal. In
// general conditions flow up through static links but not other dependencies;
// except, restricted sometimes flows up through dynamic links.
//
// The latter two functions control what happens during the top-down traversal.
// In general, only restricted conditions flow down at all, and only through
// static links. Because top-down traversals are context-specific or policy-
// specific, they generally filter the conditions to the one of interest. In
// the case of the "notice" context, it renames all of the relevant conditions
// to "notice" before filtering.


// depConditionsApplicableToTarget returns the conditions which propagate up an
// edge from dependency to target.
//
// This function sets the policy for the bottom-up traversal and how conditions
// flow up the graph from dependencies to targets.
//
// If a pure aggregation is built into a derivative work that is not a pure
// aggregation, per policy it ceases to be a pure aggregation in the context of
// that derivative work. The `treatAsAggregate` parameter will be false for
// non-aggregates and for aggregates in non-aggregate contexts.
func depConditionsApplicableToTarget(e targetEdgeImp, depConditions LicenseConditionSet, treatAsAggregate bool) *licenseConditionSetImp {
	result := depConditions.(*licenseConditionSetImp).Copy().(*licenseConditionSetImp)
	if !e.isDerivation() {
		// target is not a derivative work of dependency
		result.removeAllByName(ConditionNames{"unencumbered", "permissive", "notice", "reciprocal", "proprietary"})
		// FIXME: need code here for e.areIndependentModules() and certain kinds of GPL
	}
	return result
}

// targetConditionsApplicableToDep returns the conditions which propagate down
// an edge from target to dependency.
//
// This function sets the policy for the top-down traversal and how conditions
// flow down the graph from targets to dependencies.
//
// If a pure aggregation is built into a derivative work that is not a pure
// aggregation, per policy it ceases to be a pure aggregation in the context of
// that derivative work. The `treatAsAggregate` parameter will be false for
// non-aggregates and for aggregates in non-aggregate contexts.
func targetConditionsApplicableToDep(e targetEdgeImp, targetConditions LicenseConditionSet, treatAsAggregate bool) *licenseConditionSetImp {
	result := targetConditions.(*licenseConditionSetImp).Copy().(*licenseConditionSetImp)

	// reverse direction -- none of these apply to things depended-on, only to targets depending-on.
	result.removeAllByName(ConditionNames{"unencumbered", "permissive", "notice", "reciprocal", "proprietary", "by_exception_only"})

	if treatAsAggregate || !e.isDerivation() {
		// FIXME: probably need to propagate restricted if the aggregate itself originates the condition
		result.removeAllByName(ImpliesRestricted)
	}
	return result
}

// selfConditionsApplicableForConditionName adjusts the conditions per the
// condition name.
//
// This function sets the policy for top-down traversal and which conditions
// are relevant for a given context.
func selfConditionsApplicableForConditionName(conditionName string, selfConditions LicenseConditionSet, treatAsAggregate bool) *licenseConditionSetImp {
	result := selfConditions.(*licenseConditionSetImp) // no copy required -- copied in filter and rename
	switch conditionName {
	case "unencumbered":
		return result.filter(treatAsAggregate, ImpliesUnencumbered)
	case "permissive":
		return result.filter(treatAsAggregate, ImpliesPermissive)
	case "notice":
		return result.rename(treatAsAggregate, "notice", ImpliesNotice)
	case "reciprocal":
		return result.filter(treatAsAggregate, ImpliesReciprocal)
	case "restricted":
		return result.filter(treatAsAggregate, ImpliesRestricted)
	case "proprietary":
		return result.filter(treatAsAggregate, ImpliesProprietary)
	case "by_exception_only":
		return result.filter(treatAsAggregate, ImpliesByExceptionOnly)
	case "all":
		return result.copy(treatAsAggregate)
	default:
		panic(fmt.Errorf("resolve requested for unknown license condition: %q", conditionName))
	}
}

// SourceSharePrivacyConflict describes an individual conflict between a source-sharing
// condition and a source privacy condition
type SourceSharePrivacyConflict struct {
	AppliesTo TargetNode
	ShareCondition LicenseCondition
	PrivacyCondition LicenseCondition
}

// Error returns a string describing the conflict.
func (conflict SourceSharePrivacyConflict) Error() string {
	return fmt.Sprintf("%s %s from %s and must share from %s %s\n",
		conflict.AppliesTo.Name(),
		conflict.PrivacyCondition.Name(), conflict.PrivacyCondition.Origin().Name(),
		conflict.ShareCondition.Name(), conflict.ShareCondition.Origin().Name())
}

// ConflictingSharedPrivateSource lists all of the targets where conflicting conditions to
// share the source and to keep the source private apply to the target.
func ConflictingSharedPrivateSource(licenseGraph LicenseGraph) []SourceSharePrivacyConflict {
	// shareSource is the set of all source-sharing resolutions.
	shareSource := ResolveSourceSharing(licenseGraph)

	// privateSource is the set of all source privacy resolutions.
	privateSource := ResolveSourcePrivacy(licenseGraph)

	// combined is the combination of source-sharing and source privacy.
	combined := JoinResolutions(shareSource, privateSource)

	// size is the size of the result
	size := 0
	for _, appliesTo := range combined.AppliesTo() {
		cs := combined.Conditions(appliesTo)
		size += cs.CountByName(ImpliesShared) * cs.CountByName(ImpliesPrivate)
	}
	result := make([]SourceSharePrivacyConflict, 0, size)
	for _, appliesTo := range combined.AppliesTo() {
		cs := combined.Conditions(appliesTo)

		pconditions := cs.ByName(ImpliesPrivate)
		ssconditions := cs.ByName(ImpliesShared)

		// report all conflicting condition combinations
		for _, p := range pconditions {
			for _, ss := range ssconditions {
				result = append(result, SourceSharePrivacyConflict{appliesTo, ss, p})
			}
		}
	}
	return result
}

// ResolveSourceSharing implements the policy for source-sharing.
//
// Reciprocal and Restricted top-down walks are joined except for reciprocal
// only the conditions where the AppliesTo and Origin are the same are preserved.
//
// Per policy, the source for reciprocal triggering target does not need to be
// shared unless it is also the originating target.
func ResolveSourceSharing(licenseGraph LicenseGraph) ResolutionSet {
	lg := licenseGraph.(*licenseGraphImp)

	// When a target includes code from a reciprocal project, policy is only the source
	// for the reciprocal project, i.e. the origin of the condition, needs to be shared.
	recip := ResolveTopDownForCondition(licenseGraph, "reciprocal")

	// When a target includes code from a restricted project, the policy
	// requires sharing the restricted project, the including project(s), and the
	// transitive closure of the including project(s) derivation dependencies.
	restrict := ResolveTopDownForCondition(licenseGraph, "restricted")

	// shareSource is the set of all source-sharing resolutions.
	shareSource := make(map[string]*licenseConditionSetImp)

	// For reciprocal, policy is to share the originating code only.
	rsimp := recip.(*resolutionSetImp)
	for appliesTo, cs := range rsimp.resolutions {
		conditions := cs.byOrigin(appliesTo)
		if len(conditions) == 0 {
			continue
		}
		shareSource[appliesTo] = newLicenseConditionSet(&targetNodeImp{lg, appliesTo})
		shareSource[appliesTo].addSet(rsimp.resolutions[appliesTo])
	}

	// For restricted, policy is to share the originating code and any code linked to it.
	rsimp = restrict.(*resolutionSetImp)
	for appliesTo, cs := range rsimp.resolutions {
		conditions := cs.AsList()
		if len(conditions) == 0 {
			continue
		}
		shareSource[appliesTo] = newLicenseConditionSet(&targetNodeImp{lg, appliesTo})
		shareSource[appliesTo].addSet(rsimp.resolutions[appliesTo])
	}

	return &resolutionSetImp{lg, shareSource}
}


// ResolveSourcePrivacy implements the policy for source privacy.
//
// Based on the "proprietary" walk, selectively prunes inherited conditions
// from pure aggregates.
func ResolveSourcePrivacy(licenseGraph LicenseGraph) ResolutionSet {
	lg := licenseGraph.(*licenseGraphImp)

	// For pure aggregates, a container may contain both binaries with private source
	// and binaries with shared source so the pure aggregate does not need to inherit
	// the "proprietary" condition. Policy requires non-aggregates and
	// not-pure-aggregates with source sharing requirements to share all of the source
	// affecting any private source contained therein.
	proprietary := ResolveTopDownForCondition(licenseGraph, "proprietary")

	// privateSource is the set of all source privacy resolutions pruned to account for
	// pure aggregations.
	privateSource := make(map[string]*licenseConditionSetImp)

	rsimp := proprietary.(*resolutionSetImp)
	for appliesTo, cs := range rsimp.resolutions {
		found := false
		for _, origins := range cs.conditions {
			// selfConditionsApplicableForConditionName(...) above causes all
			// conditions in a ResolutionSet to record whether they apply to
			// pure aggregates.
			for origin, treatAsAggregate := range origins {
				if treatAsAggregate == nil || treatAsAggregate == false || origin == appliesTo {
					found = true
				}
			}
		}
		if !found {
			// all conditions for `appliesTo` pruned so skip rest
			continue
		}
		for name, origins := range cs.conditions {
			for origin, treatAsAggregate := range origins {
				if treatAsAggregate == nil || treatAsAggregate == false || origin == appliesTo {
					if _, ok := privateSource[appliesTo]; !ok {
						privateSource[appliesTo] = newLicenseConditionSet(nil)
					}
					if _, ok := privateSource[appliesTo].conditions[name]; !ok {
						privateSource[appliesTo].conditions[name] = make(map[string]interface{})
					}
					privateSource[appliesTo].conditions[name][origin] = rsimp.resolutions[appliesTo].conditions[name][origin]
				}
			}
		}
	}

	return &resolutionSetImp{lg, privateSource}
}


// edgeIsDynamicLink returns true for edges representing shared libraries
// linked dynamically at runtime.
func edgeIsDynamicLink(e *dependencyEdge) bool {
	_, isPresent := e.annotations["dynamic"]
	return isPresent
}

// edgeIsDerivation returns true for edges where the target is a derivative
// work of dependency.
func edgeIsDerivation(e *dependencyEdge) bool {
	// FIXME: is this the correct definition/policy?
	_, isDynamic := e.annotations["dynamic"]
	_, isToolchain := e.annotations["toolchain"]
	return !isDynamic && !isToolchain
}

// edgeNodesAreIndependentModules returns true for edges where the target and
// dependency are independent modules.
func edgeNodesAreIndependentModules(e TargetEdge) bool {
	// FIXME: is this the correct definition/policy?
	return e.Target().PackageName() != e.Dependency().PackageName()
}
