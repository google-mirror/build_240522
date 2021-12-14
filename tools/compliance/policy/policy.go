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
	"regexp"
	"strings"
)

var (
	// ImpliesUnencumbered lists the condition names representing an author attempt to disclaim copyright.
	ImpliesUnencumbered = LicenseConditionSet(UnencumberedCondition)

	// ImpliesPermissive lists the condition names representing copyrighted but "licensed without policy requirements".
	ImpliesPermissive = LicenseConditionSet(PermissiveCondition)

	// ImpliesNotice lists the condition names implying a notice or attribution policy.
	ImpliesNotice = LicenseConditionSet(UnencumberedCondition | PermissiveCondition | NoticeCondition | ReciprocalCondition |
			RestrictedCondition | RestrictedClasspathExceptionCondition | WeaklyRestrictedCondition |
			ProprietaryCondition | ByExceptionOnlyCondition)

	// ImpliesReciprocal lists the condition names implying a local source-sharing policy.
	ImpliesReciprocal = LicenseConditionSet(ReciprocalCondition)

	// Restricted lists the condition names implying an infectious source-sharing policy.
	ImpliesRestricted = LicenseConditionSet(RestrictedCondition | RestrictedClasspathExceptionCondition | WeaklyRestrictedCondition)

	// ImpliesProprietary lists the condition names implying a confidentiality policy.
	ImpliesProprietary = LicenseConditionSet(ProprietaryCondition)

	// ImpliesByExceptionOnly lists the condition names implying a policy for "license review and approval before use".
	ImpliesByExceptionOnly = LicenseConditionSet(ProprietaryCondition | ByExceptionOnlyCondition)

	// ImpliesPrivate lists the condition names implying a source-code privacy policy.
	ImpliesPrivate = LicenseConditionSet(ProprietaryCondition)

	// ImpliesShared lists the condition names implying a source-code sharing policy.
	ImpliesShared = LicenseConditionSet(ReciprocalCondition | RestrictedCondition | RestrictedClasspathExceptionCondition | WeaklyRestrictedCondition)
)

var (
	anyLgpl      = regexp.MustCompile(`^SPDX-license-identifier-LGPL.*`)
	versionedGpl = regexp.MustCompile(`^SPDX-license-identifier-GPL-\p{N}.*`)
	genericGpl   = regexp.MustCompile(`^SPDX-license-identifier-GPL$`)
	ccBySa       = regexp.MustCompile(`^SPDX-license-identifier-CC-BY.*-SA.*`)
)


// newLicenseConditionSet returns a set containing the recognized `names` and
// silently ignoring or discarding the unrecognized `names`.
func LicenseConditionSetFromNames(tn *TargetNode, names ...string) LicenseConditionSet {
	cs := LicenseConditionSet(0x00)
	for _, name := range names {
		if name == "restricted" {
			if 0 == len(tn.LicenseKinds()) {
				cs = cs.Plus(RestrictedCondition)
				continue
			}
			hasLgpl := false
			hasClasspath := false
			hasGeneric := false
			for _, kind := range tn.LicenseKinds() {
				if strings.HasSuffix(kind, "-with-classpath-exception") {
					cs = cs.Plus(RestrictedClasspathExceptionCondition)
					hasClasspath = true
				} else if anyLgpl.MatchString(kind) {
					cs = cs.Plus(WeaklyRestrictedCondition)
					hasLgpl = true
				} else if versionedGpl.MatchString(kind) {
					cs = cs.Plus(RestrictedCondition)
				} else if genericGpl.MatchString(kind) {
					hasGeneric = true
				} else if kind == "legacy_restricted" || ccBySa.MatchString(kind) {
					cs = cs.Plus(RestrictedCondition)
				} else {
					cs = cs.Plus(RestrictedCondition)
				}
			}
			if hasGeneric && !hasLgpl && !hasClasspath {
				cs = cs.Plus(RestrictedCondition)
			}
			continue
		}
		if lc, ok := RecognizedConditionNames[name]; ok {
			cs |= LicenseConditionSet(lc)
		}
	}
	return cs
}


// Resolution happens in two passes:
//
// 1. A bottom-up traversal propagates license conditions up to targets from
// dendencies as needed.
//
// 2. For each condition of interest, a top-down traversal adjusts the attached
// conditions pushing restricted down from targets into linked dependencies.
//
// The behavior of the 2 passes gets controlled by the 2 functions below.
//
// The first function controls what happens during the bottom-up traversal. In
// general conditions flow up through static links but not other dependencies;
// except, restricted sometimes flows up through dynamic links.
//
// In general, too, the originating target gets acted on to resolve the
// condition (e.g. providing notice), but again restricted is special in that
// it requires acting on (i.e. sharing source of) both the originating module
// and the target using the module.
//
// The latter function controls what happens during the top-down traversal. In
// general, only restricted conditions flow down at all, and only through
// static links.
//
// Not all restricted licenses are create equal. Some have special rules or
// exceptions. e.g. LGPL or "with classpath excption".

// depActionsApplicableToTarget returns the actions which propagate up an
// edge from dependency to target.
//
// This function sets the policy for the bottom-up traversal and how conditions
// flow up the graph from dependencies to targets.
//
// If a pure aggregation is built into a derivative work that is not a pure
// aggregation, per policy it ceases to be a pure aggregation in the context of
// that derivative work. The `treatAsAggregate` parameter will be false for
// non-aggregates and for aggregates in non-aggregate contexts.
func depActionsApplicableToTarget(lg *LicenseGraph, e *TargetEdge, depActions actionSet, treatAsAggregate bool) actionSet {
	result := actionSet{lg, &IntervalSet{}}
	if edgeIsDerivation(e) {
		result.addSet(depActions)
		for a := range depActions.matchingAnySet(ImpliesRestricted){
			result.add(e.Target(), a.cs & ImpliesRestricted)
		}
		return result
	}
	if !edgeIsDynamicLink(e) {
		return result
	}

	for a := range depActions.matchingAny(RestrictedCondition, RestrictedClasspathExceptionCondition) {
		if a.cs.HasAny(RestrictedCondition) {
			result.addCondition(a.actsOn, RestrictedCondition)
			result.addCondition(e.Target(), RestrictedCondition)
		}
		if a.cs.HasAny(RestrictedClasspathExceptionCondition) && !edgeNodesAreIndependentModules(e) {
			result.addCondition(a.actsOn, RestrictedClasspathExceptionCondition)
			result.addCondition(e.Target(), RestrictedClasspathExceptionCondition)
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
func targetConditionsApplicableToDep(lg *LicenseGraph, e *TargetEdge, targetConditions LicenseConditionSet, treatAsAggregate bool) LicenseConditionSet {
	result := targetConditions

	// reverse direction -- none of these apply to things depended-on, only to targets depending-on.
	result = result.Minus(UnencumberedCondition, PermissiveCondition, NoticeCondition, ReciprocalCondition, ProprietaryCondition, ByExceptionOnlyCondition)

	if !edgeIsDerivation(e) && !edgeIsDynamicLink(e) {
		// target is not a derivative work of dependency and is not linked to dependency
		result = result.Difference(ImpliesRestricted)
		return result
	}
	if treatAsAggregate {
		// If the author of a pure aggregate licenses it restricted, apply restricted to immediate dependencies.
		// Otherwise, restricted does not propagate back down to dependencies.
		if !LicenseConditionSetFromNames(e.target, e.target.proto.LicenseConditions...).MatchesAnySet(ImpliesRestricted) {
			result = result.Difference(ImpliesRestricted)
		}
		return result
	}
	if edgeIsDerivation(e) {
		return result
	}
	result = result.Minus(WeaklyRestrictedCondition)
	if edgeNodesAreIndependentModules(e) {
		result = result.Minus(RestrictedClasspathExceptionCondition)
	}
	return result
}

// edgeIsDynamicLink returns true for edges representing shared libraries
// linked dynamically at runtime.
func edgeIsDynamicLink(e *TargetEdge) bool {
	return e.annotations.HasAnnotation("dynamic")
}

// edgeIsDerivation returns true for edges where the target is a derivative
// work of dependency.
func edgeIsDerivation(e *TargetEdge) bool {
	isDynamic := e.annotations.HasAnnotation("dynamic")
	isToolchain := e.annotations.HasAnnotation("toolchain")
	return !isDynamic && !isToolchain
}

// edgeNodesAreIndependentModules returns true for edges where the target and
// dependency are independent modules.
func edgeNodesAreIndependentModules(e *TargetEdge) bool {
	return e.target.PackageName() != e.dependency.PackageName()
}
