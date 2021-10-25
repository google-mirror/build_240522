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
	"regexp"
	"strings"
)

var (
	// ImpliesUnencumbered lists the condition names representing an author attempt to disclaim copyright.
	ImpliesUnencumbered = ConditionNames{"unencumbered"}

	// ImpliesPermissive lists the condition names representing copyrighted but "licensed without policy requirements".
	ImpliesPermissive = ConditionNames{"permissive"}

	// ImpliesNotice lists the condition names implying a notice or attribution policy.
	ImpliesNotice = ConditionNames{"unencumbered", "permissive", "notice", "reciprocal", "restricted", "proprietary", "by_exception_only"}

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

var (
	anyLgpl = regexp.MustCompile(`^SPDX-license-identifier-LGPL.*`)
	versionedGpl = regexp.MustCompile(`^SPDX-license-identifier-GPL-\p{N}.*`)
	genericGpl = regexp.MustCompile(`^SPDX-license-identifier-GPL$`)
	ccBySa = regexp.MustCompile(`^SPDX-license-identifier-CC-BY.*-SA.*`)
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
		if !e.isDynamicLink() {
			// target is not a derivative work of dependency and is not linked to dependency
			result.removeAllByName(ImpliesRestricted)
			return result
		}
		restricted := result.ByName(ImpliesRestricted)
		for _, lc := range restricted {
			hasGpl := false
			hasLgpl := false
			hasClasspath := false
			hasGeneric := false
			hasOther := false
			for _, kind := range lc.Origin().LicenseKinds() {
				if strings.HasSuffix(kind, "-with-classpath-exception") {
					hasClasspath = true
				} else if anyLgpl.MatchString(kind) {
					hasLgpl = true
				} else if versionedGpl.MatchString(kind) {
					hasGpl = true
				} else if genericGpl.MatchString(kind) {
					hasGeneric = true
				} else if kind == "legacy_restricted" || ccBySa.MatchString(kind) {
					hasOther = true
				}
			}
			if hasOther || hasGpl {
				continue
			}
			if hasClasspath && !e.areIndependentModules() {
				continue
			}
			if hasGeneric && !hasLgpl && !hasClasspath {
				continue
			}
			result.remove(lc)
		}
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

	if !e.isDerivation() && !e.isDynamicLink() {
		// target is not a derivative work of dependency and is not linked to dependency
		result.removeAllByName(ImpliesRestricted)
		return result
	}
	if treatAsAggregate {
		// If the author of a pure aggregate licenses it restricted, apply restricted to immediate dependencies.
		// Otherwise, restricted does not propagate back down to dependencies.
		restricted := result.ByName(ImpliesRestricted)
		for _, lc := range restricted {
			if lc.Origin().Name() != e.e.target {
				result.remove(lc)
			}
		}
		return result
	}
	if e.isDerivation() {
		return result
	}
	restricted := result.ByName(ImpliesRestricted)
	for _, lc := range restricted {
		hasGpl := false
		hasLgpl := false
		hasClasspath := false
		hasGeneric := false
		hasOther := false
		for _, kind := range lc.Origin().LicenseKinds() {
			if strings.HasSuffix(kind, "-with-classpath-exception") {
				hasClasspath = true
			} else if anyLgpl.MatchString(kind) {
				hasLgpl = true
			} else if versionedGpl.MatchString(kind) {
				hasGpl = true
			} else if genericGpl.MatchString(kind) {
				hasGeneric = true
			} else if kind == "legacy_restricted" || ccBySa.MatchString(kind) {
				hasOther = true
			}
		}
		if hasOther || hasGpl {
			continue
		}
		if hasClasspath && !e.areIndependentModules() {
			continue
		}
		if hasGeneric && !hasLgpl && !hasClasspath {
			continue
		}
		result.remove(lc)
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
