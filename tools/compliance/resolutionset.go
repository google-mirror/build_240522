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
	rmap := make(map[string]actionSet)
	lg := resolutions[0].lg
	for _, r := range resolutions {
		if len(r.resolutions) < 1 {
			continue
		}
		if lg == nil {
			lg = r.lg
		} else if r.lg == nil {
			panic(fmt.Errorf("attempt to join nil resolution set with %d resolutions", len(r.resolutions)))
		} else if lg != r.lg {
			panic(fmt.Errorf("attempt to join resolutions from multiple graphs"))
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
	return &ResolutionSet{lg, rmap}
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
	// lg defines the scope of the resolutions to one LicenseGraph
	lg *LicenseGraph

	// resolutions maps names of target with applicable conditions to the set of conditions that apply.
	resolutions map[string]actionSet
}

func (rs *ResolutionSet) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "{")
	sep := ""
	for attachesTo, as := range rs.resolutions {
		fmt.Fprintf(&sb, "%s%s -> %s", sep, attachesTo, as.String())
		sep = ", "
	}
	fmt.Fprintf(&sb, "}")
	return sb.String()
}

// AttachesTo identifies the list of targets triggering action to resolve
// conditions. (unordered)
func (rs *ResolutionSet) AttachesTo() []TargetNode {
	targets := make([]TargetNode, 0, len(rs.resolutions))
	for target := range rs.resolutions {
		targets = append(targets, TargetNode{rs.lg, target})
	}
	return targets
}

// ActsOn identifies the list of targets to act on (share, give notice etc.)
// to resolve conditions. (unordered)
func (rs *ResolutionSet) ActsOn() []TargetNode {
	tset := make(map[string]interface{})
	for _, as := range rs.resolutions {
		for actsOn := range as {
			tset[actsOn] = nil
		}
	}
	targets := make([]TargetNode, 0, len(tset))
	for target := range tset {
		targets = append(targets, TargetNode{rs.lg, target})
	}
	return targets
}

// Origins identifies the list of targets originating conditions to resolve.
// (unordered)
func (rs *ResolutionSet) Origins() []TargetNode {
	tset := make(map[string]interface{})
	for _, as := range rs.resolutions {
		for _, cs := range as {
			for _, origins := range cs.conditions {
				for origin := range origins {
					tset[origin] = nil
				}
			}
		}
	}
	targets := make([]TargetNode, 0, len(tset))
	for target := range tset {
		targets = append(targets, TargetNode{rs.lg, target})
	}
	return targets
}

// Resolutions returns the list of resolutions that `attachedTo`
// target must resolve. Returns empty list if no conditions apply.
//
// Panics if `attachedTo` does not appear in the set.
func (rs *ResolutionSet) Resolutions(attachedTo TargetNode) ResolutionList {
	if rs.lg == nil {
		rs.lg = attachedTo.lg
	} else if attachedTo.lg == nil {
		attachedTo.lg = rs.lg
	} else if rs.lg != attachedTo.lg {
		panic(fmt.Errorf("attempt to query target resolutions for wrong graph"))
	}
	as, ok := rs.resolutions[attachedTo.file]
	if !ok {
		return ResolutionList{}
	}
	result := make(ResolutionList, 0, len(as))
	for actsOn, cs := range as {
		result = append(result, Resolution{attachedTo.file, actsOn, cs.Copy()})
	}
	return result
}

// ResolutionsByActsOn returns the list of resolutions that must `actOn` to
// resolvee. Returns empty list if no conditions apply.
//
// Panics if `actOn` does not appear in the set.
func (rs *ResolutionSet) ResolutionsByActsOn(actOn TargetNode) ResolutionList {
	if rs.lg == nil {
		rs.lg = actOn.lg
	} else if actOn.lg == nil {
		actOn.lg = rs.lg
	} else if rs.lg != actOn.lg {
		panic(fmt.Errorf("attempt to query target resolutions for wrong graph"))
	}
	c := 0
	for _, as := range rs.resolutions {
		if _, ok := as[actOn.file]; ok {
			c++
		}
	}
	result := make(ResolutionList, 0, c)
	for attachedTo, as := range rs.resolutions {
		if cs, ok := as[actOn.file]; ok {
			result = append(result, Resolution{attachedTo, actOn.file, cs.Copy()})
		}
	}
	return result
}

// AttachesToByOrigin identifies the list of targets requiring action to
// resolve conditions originating at `origin`. (unordered)
func (rs *ResolutionSet) AttachesToByOrigin(origin TargetNode) []TargetNode {
	if rs.lg == nil {
		rs.lg = origin.lg
	} else if origin.lg == nil {
		origin.lg = rs.lg
	} else if rs.lg != origin.lg {
		panic(fmt.Errorf("attempt to query targets by origin for wrong graph"))
	}
	tset := make(map[string]interface{})
	for target, as := range rs.resolutions {
		for _, cs := range as {
			if cs.hasByOrigin(origin.file) {
				tset[target] = nil
				break
			}
		}
	}
	targets := make([]TargetNode, 0, len(tset))
	for target := range tset {
		targets = append(targets, TargetNode{rs.lg, target})
	}
	return targets
}

// AttachesToTarget returns true if the set contains conditions that
// are `attachedTo`.
func (rs *ResolutionSet) AttachesToTarget(attachedTo TargetNode) bool {
	if rs.lg == nil {
		rs.lg = attachedTo.lg
	} else if attachedTo.lg == nil {
		attachedTo.lg = rs.lg
	} else if rs.lg != attachedTo.lg {
		panic(fmt.Errorf("attempt to query resolved targets for wrong graph"))
	}
	return rs.hasTarget(attachedTo.file)
}

// AnyByNameAttachToTarget returns true if the set contains conditions matching
// `names` that attach to `attachedTo`.
func (rs *ResolutionSet) AnyByNameAttachToTarget(attachedTo TargetNode, names ...ConditionNames) bool {
	if rs.lg == nil {
		rs.lg = attachedTo.lg
	} else if attachedTo.lg == nil {
		attachedTo.lg = rs.lg
	} else if rs.lg != attachedTo.lg {
		panic(fmt.Errorf("attempt to query target resolutions for wrong graph"))
	}
	return rs.hasAnyByName(attachedTo.file, names...)
}

// AllByNameAttachTo returns true if the set contains at least one condition
// matching each element of `names` for `attachedTo`.
func (rs *ResolutionSet) AllByNameAttachToTarget(attachedTo TargetNode, names ...ConditionNames) bool {
	if rs.lg == nil {
		rs.lg = attachedTo.lg
	} else if attachedTo.lg == nil {
		attachedTo.lg = rs.lg
	} else if rs.lg != attachedTo.lg {
		panic(fmt.Errorf("attempt to query target resolutions for wrong graph"))
	}
	return rs.hasAllByName(attachedTo.file, names...)
}

// compliance-only ResolutionSet methods

// newResolutionSet constructs a new, empty instance of resolutionSetImp for graph `lg`.
func newResolutionSet(lg *LicenseGraph) *ResolutionSet {
	return &ResolutionSet{lg, make(map[string]actionSet)}
}

// add attaches all of the license conditions in `cs` to `file` to act on the originating node if not already applied.
func (rs *ResolutionSet) add(file string, cs *LicenseConditionSet) {
	if cs.IsEmpty() {
		return
	}
	_, ok := rs.resolutions[file]
	if !ok {
		rs.resolutions[file] = make(actionSet)
	}

	for name, origins := range cs.conditions {
		for origin := range origins {
			if _, ok := rs.resolutions[file][origin]; !ok {
				rs.resolutions[file][origin] = newLicenseConditionSet(rs.lg)
			}
			rs.resolutions[file][origin].addAll(origin, name)
		}
	}
}

// addConditions attaches all of the license conditions in `as` to `file` to act on the originating node if not already applied.
func (rs *ResolutionSet) addConditions(file string, as actionSet) {
	_, ok := rs.resolutions[file]
	if !ok {
		rs.resolutions[file] = as.copy()
		return
	}
	rs.resolutions[file].addSet(as)
}

// add attaches all of the license conditions in `as` to `file` to act on `file` if not already applied.
func (rs *ResolutionSet) addSelf(file string, as actionSet) {
	for _, cs := range as {
		if cs.IsEmpty() {
			return
		}
		_, ok := rs.resolutions[file]
		if !ok {
			rs.resolutions[file] = make(actionSet)
		}
		_, ok = rs.resolutions[file][file]
		if !ok {
			rs.resolutions[file][file] = newLicenseConditionSet(rs.lg)
		}
		rs.resolutions[file][file].AddSet(cs)
	}
}

// hasTarget returns true if `file` appears as one of the targets in the set.
func (rs *ResolutionSet) hasTarget(file string) bool {
	_, isPresent := rs.resolutions[file]
	return isPresent
}

// hasAnyByName returns true if the target for `file` has at least 1 condition matching `names`.
func (rs *ResolutionSet) hasAnyByName(file string, names ...ConditionNames) bool {
	objects, isPresent := rs.resolutions[file]
	if !isPresent {
		return false
	}
	for _, cs := range objects {
		for _, cn := range names {
			for _, name := range cn {
				_, isPresent = cs.conditions[name]
				if isPresent {
					return true
				}
			}
		}
	}
	return false
}

// hasAllByName returns true if the target for `file` has at least 1 condition for each element of `names`.
func (rs *ResolutionSet) hasAllByName(file string, names ...ConditionNames) bool {
	as, isPresent := rs.resolutions[file]
	if !isPresent {
		return false
	}
	for _, cn := range names {
		found := false
	asloop:
		for _, cs := range as {
			for _, name := range cn {
				_, isPresent = cs.conditions[name]
				if isPresent {
					found = true
					break asloop
				}
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// hasResolution returns true if the exact `condition` applies to `file` acting on `actsOn`.
func (rs *ResolutionSet) hasResolution(file, actsOn string, condition LicenseCondition) bool {
	as, isPresent := rs.resolutions[file]
	if !isPresent {
		return false
	}
	cs, isPresent := as[actsOn]
	if !isPresent {
		return false
	}
	origins, isPresent := cs.conditions[condition.Name()]
	if !isPresent {
		return false
	}
	_, isPresent = origins[condition.Origin().Name()]
	return isPresent
}
