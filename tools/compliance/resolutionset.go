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

// JoinResolutionSets returns a new ResolutionSet combining the resolutions from
// multiple resolution sets. All sets must be derived from the same license
// graph.
//
// e.g. combine "restricted", "reciprocal", and "proprietary" resolutions.
func JoinResolutionSets(resolutions ...*ResolutionSet) *ResolutionSet {
	if len(resolutions) < 1 {
		panic(fmt.Errorf("attempt to join 0 resolution sets"))
	}
	rmap := make(map[*TargetNode]actionSet)
	for _, r := range resolutions {
		if len(r.resolutions) < 1 {
			continue
		}
		for attachesTo, as := range r.resolutions {
			if as.isEmpty() {
				continue
			}
			if _, ok := rmap[attachesTo]; !ok {
				rmap[attachesTo] = as.copy()
				continue
			}
			rmap[attachesTo].addSet(as)
		}
	}
	return &ResolutionSet{rmap}
}

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
// Resolutions for different contexts may be combined in a new ResolutionSet
// using JoinResolutions(...).
//
// See: resolve.go for:
//  * ResolveBottomUpConditions(...)
//  * ResolveTopDownForCondition(...)
// See also: policy.go for:
//  * ResolveSourceSharing(...)
//  * ResolveSourcePrivacy(...)
type ResolutionSet struct {
	// resolutions maps names of target with applicable conditions to the set of conditions that apply.
	resolutions map[*TargetNode]actionSet
}

// String returns a string representation of the set.
func (rs *ResolutionSet) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "{")
	sep := ""
	for attachesTo, as := range rs.resolutions {
		fmt.Fprintf(&sb, "%s%s -> %s", sep, attachesTo.name, as.String())
		sep = ", "
	}
	fmt.Fprintf(&sb, "}")
	return sb.String()
}

// AttachesTo identifies the list of targets triggering action to resolve
// conditions. (unordered)
func (rs *ResolutionSet) AttachesTo() TargetNodeList {
	targets := make(TargetNodeList, 0, len(rs.resolutions))
	for attachesTo := range rs.resolutions {
		targets = append(targets, attachesTo)
	}
	return targets
}

// ActsOn identifies the list of targets to act on (share, give notice etc.)
// to resolve conditions. (unordered)
func (rs *ResolutionSet) ActsOn() TargetNodeList {
	tset := make(map[*TargetNode]struct{})
	for _, as := range rs.resolutions {
		as.VisitAll(func(a resolutionAction) { tset[a.actsOn] = struct{}{} })
	}
	targets := make(TargetNodeList, 0, len(tset))
	for target := range tset {
		targets = append(targets, target)
	}
	return targets
}

// Resolutions returns the list of resolutions that `attachedTo`
// target must resolve. Returns empty list if no conditions apply.
//
// Panics if `attachedTo` does not appear in the set.
func (rs *ResolutionSet) Resolutions(attachedTo *TargetNode) ResolutionList {
	as, ok := rs.resolutions[attachedTo]
	if !ok {
		return ResolutionList{}
	}
	result := make(ResolutionList, 0, as.Len())
	as.VisitAll(func(a resolutionAction) {
		result = append(result, Resolution{attachedTo, a.actsOn, a.cs})
	})
	return result
}

// ResolutionsByActsOn returns the list of resolutions that must `actOn` to
// resolvee. Returns empty list if no conditions apply.
//
// Panics if `actOn` does not appear in the set.
func (rs *ResolutionSet) ResolutionsByActsOn(actOn *TargetNode) ResolutionList {
	c := 0
	for _, as := range rs.resolutions {
		as.VisitAll(func(a resolutionAction) {
			if a.actsOn == actOn {
				c++
			}
		})
	}
	result := make(ResolutionList, 0, c)
	for attachedTo, as := range rs.resolutions {
		as.VisitAll(func(a resolutionAction) {
			if a.actsOn == actOn {
				result = append(result, Resolution{attachedTo, actOn, a.cs})
			}
		})
	}
	return result
}

// AttachesToTarget returns true if the set contains conditions that
// are `attachedTo`.
func (rs *ResolutionSet) AttachesToTarget(attachedTo *TargetNode) bool {
	_, isPresent := rs.resolutions[attachedTo]
	return isPresent
}

// AnyMatchingAttachToTarget returns true if the set contains conditions matching
// `names` that attach to `attachedTo`.
func (rs *ResolutionSet) AnyMatchingAttachToTarget(attachedTo *TargetNode, conditions ...LicenseConditionSet) bool {
	as, isPresent := rs.resolutions[attachedTo]
	if !isPresent {
		return false
	}
	return as.matchesAnySet(conditions...)
}

// IsEmpty returns true if the set contains no conditions to resolve.
func (rs *ResolutionSet) IsEmpty() bool {
	for _, as := range rs.resolutions {
		if !as.isEmpty() {
			return false
		}
	}
	return true
}

// compliance-only ResolutionSet methods

// newResolutionSet constructs a new, empty instance of resolutionSetImp for graph `lg`.
func newResolutionSet() *ResolutionSet {
	return &ResolutionSet{make(map[*TargetNode]actionSet)}
}

// addConditions attaches all of the license conditions in `as` to `attachTo` to act on the original node if not already applied.
func (rs *ResolutionSet) addConditions(attachTo *TargetNode, as actionSet) {
	_, ok := rs.resolutions[attachTo]
	if !ok {
		rs.resolutions[attachTo] = as.copy()
		return
	}
	rs.resolutions[attachTo].addSet(as)
}

// add attaches all of the license conditions in `as` to `attachTo` to act on `attachTo` if not already applied.
func (rs *ResolutionSet) addSelf(attachTo *TargetNode, c chan resolutionAction) {
	for a := range c {
		if a.cs.IsEmpty() {
			continue
		}
		_, ok := rs.resolutions[attachTo]
		if !ok {
			rs.resolutions[attachTo] = actionSet{attachTo.lg, &IntervalSet{}}
		}
		rs.resolutions[attachTo].add(attachTo, a.cs)
	}
}
