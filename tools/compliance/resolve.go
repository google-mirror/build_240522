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

// ConditionNames implements the Contains predicate for lists of condition names.
type ConditionNames []string

var (
	// Notice lists the condition names implying a notice or attribution policy.
	Notice = ConditionNames{"unencumbered", "permissive", "notice", "reciprocal", "restricted"}

	// Proprietary lists the condition names implying a confidentiality or secrecy policy.
	Proprietary = ConditionNames{"proprietary"}

	// Reciprocal lists the condition names implying a local source-sharing policy.
	Reciprocal = ConditionNames{"reciprocal"}

	// Restricted lists the condition names implying an infectious source-sharing policy.
	Restricted = ConditionNames{"restricted"}
)

// Contains returns true if the name matches one of the ConditionNames.
func (cn ConditionNames) Contains(name string) bool {
	for _, c := range cn {
		if c == name {
			return true
	        }
	}
	return false
}


// resolve returns the conditions which propagate up an edge from dependency to target.
func resolve(e targetEdgeImp, cs LicenseConditionSet) *licenseConditionSetImp {
	result := cs.(*licenseConditionSetImp).Copy().(*licenseConditionSetImp)
	if !e.isDerivation() {
		result.removeAllByName(ConditionNames{"unencumbered", "permissive", "notice", "reciprocal", "proprietary"})
	}
	return result
}

// edgeIsDynamicLink returns true for edges representing shared libraries etc. linked dynamically at runtime.
func edgeIsDynamicLink(e *dependencyEdge) bool {
	_, isPresent := e.annotations["dynamic"]
	return isPresent
}

// edgeIsDerivation returns true for edges where the target is a derivative work of dependency.
func edgeIsDerivation(e *dependencyEdge) bool {
	_, isDynamic := e.annotations["dynamic"]
	_, isToolchain := e.annotations["toolchain"]
	return !isDynamic && !isToolchain
}


// ResolveBottomUpConditions performs a bottom-up walk of the LicenseGraph
// propagating conditions up the graph as necessary according to the properties
// of each edge and according to each license condition in question.
//
// Subsequent top-down walks of the graph will filter some resolutions and may
// introduce new resolutions.
//
// e.g. if a "restricted" condition applies to a binary, it also applies to all
// of the statically-linked libraries and the transitive closure of their static
// dependencies; even if neither they nor the transitive closure of their
// dependencies originate any "restricted" conditions. The bottom-up walk will
// not resolve the library and its transitive closure, but the later top-down
// walk will.
func ResolveBottomUpConditions(graph LicenseGraph) ResolutionSet {
	lg := graph.(*licenseGraphImp)

	lg.mu.Lock()
	rs := lg.rs
	lg.mu.Unlock()

	if rs != nil {
		return rs
	}

	lg.indexForward()

	rs = newResolutionSet(lg)

	var walk func(f string) *licenseConditionSetImp

	walk = func(f string) *licenseConditionSetImp {
		result := newLicenseConditionSet(&targetNodeImp{lg, f})
		if preresolved, ok := rs.resolutions[f]; ok {
			return result.union(preresolved)
		}
		for _, e := range lg.index[f] {
			cs := walk(e.dependency)
			cs = resolve(targetEdgeImp{lg, e}, cs)
			result = result.union(cs)
		}
		result = lg.targets[f].licenseConditions.union(result)
		rs.add(f, result)
		return result
	}

	for _, r := range lg.rootFiles {
		cs := walk(r)
		rs.add(r, lg.targets[r].licenseConditions.union(cs))
	}

	lg.mu.Lock()
	if lg.rs == nil{
		lg.rs = rs
	}
	lg.mu.Unlock()

        return rs
}

// ResolveTopDownRestricted performs a top-down walk of the LicenseGraph
// resolving all reachable "restricted" conditions and the transitive closure
// of their derivation subgraph.
//
// For fullWalk == true, the walk includes every resolved condition on every
// target. If only the list of targets with any restricted condition applied
// matters, the walk can run faster by not revisiting target nodes.
//
// e.g. If library L is linked into binary B and shared library S, recording
// the resolutions for both B and S on L and its transive dependencies requires
// walking the sub-tree rooted at L twice. In applications where only the fact
// that L has a source-sharing requirement matters, the walk can complete
// faster by only walking the sub-tree rooted at L once: fullWalk == false.
func ResolveTopDownRestricted(graph LicenseGraph, fullWalk bool) ResolutionSet {
	lg := graph.(*licenseGraphImp)

	rs := ResolveBottomUpConditions(lg)

	rmap := make(map[string]*licenseConditionSetImp)
	cmap := make(map[string]interface{})

	var walkContainer, walkNonContainer func(string, *licenseConditionSetImp)

	walkNonContainer = func(f string, cs *licenseConditionSetImp) {
		if _, ok := rmap[f]; ok {
			rmap[f] = rmap[f].union(cs)
		} else {
			rmap[f] = cs.Copy().(*licenseConditionSetImp)
		}
		for _, e := range lg.index[f] {
			if !edgeIsDerivation(e) {
				continue
			}
			dcs, ok := rs.(*resolutionSetImp).resolutions[e.dependency]
			if !ok {
				dcs = rmap[f].Copy().(*licenseConditionSetImp)
			} else {
				dcs = dcs.filter(Restricted).union(rmap[f])
			}
			if pcs, alreadyWalked := rmap[e.dependency]; alreadyWalked {
				if _, asContainer := cmap[e.dependency]; asContainer {
					delete(cmap, e.dependency)
					walkNonContainer(e.dependency, dcs)
					continue
				}
				if !fullWalk {
					continue
				}
				pcs = pcs.Copy().(*licenseConditionSetImp)
				pcs.remove(dcs)
				if pcs.IsEmpty() {
					continue
				}
			}
			walkNonContainer(e.dependency, dcs)
		}
	}

	walkContainer = func(f string, cs *licenseConditionSetImp) {
		if _, ok := rmap[f]; ok {
			rmap[f] = rmap[f].union(cs)
		} else {
			rmap[f] = cs.Copy().(*licenseConditionSetImp)
			cmap[f] = nil
		}
		for _, e := range lg.index[f] {
			if !edgeIsDerivation(e) {
				continue
			}
			dcs, ok := rs.(*resolutionSetImp).resolutions[e.dependency]
			if !ok {
				continue
			}
			dcs = dcs.filter(Restricted)
			if dcs.IsEmpty() {
				continue
			}
			if pcs, alreadyWalked := rmap[e.dependency]; alreadyWalked {
				if !fullWalk {
					continue
				}
				pcs = pcs.Copy().(*licenseConditionSetImp)
				pcs.remove(dcs)
				if pcs.IsEmpty() {
					continue
				}
			}
			if lg.targets[e.dependency].isContainer {
				walkContainer(e.dependency, dcs)
			} else {
				walkNonContainer(e.dependency, dcs)
			}
		}
	}

	for _, r := range lg.rootFiles {
		cs, ok := rs.(*resolutionSetImp).resolutions[r]
		if !ok {
			continue
		}
		if !cs.HasAnyByName(Restricted) {
			continue
		}
		rcs := newLicenseConditionSet(&targetNodeImp{lg, r})
		rcs.addAll(cs.ByName(Restricted))
		if lg.targets[r].isContainer {
			walkContainer(r, rcs)
		} else {
			walkNonContainer(r, rcs)
		}
	}

	return &resolutionSetImp{lg, rmap}
}

// ResolveTopDownNotice performs a top-down walk of the LicenseGraph
// resolving all reachable conditions where policy is to require notice
// or attribution.
//
// Has the effect of turning "restricted" and "permissive" into "notice" for
// this specific resolution set. The conditions on the graph remain unchanged.
func ResolveTopDownNotice(graph LicenseGraph) ResolutionSet {
	lg := graph.(*licenseGraphImp)

	rs := ResolveBottomUpConditions(lg)

	rmap := make(map[string]*licenseConditionSetImp)

	var walk func(string)

	walk = func(f string) {
		if _, ok := rmap[f]; ok {
			return
		}
		for _, e := range lg.index[f] {
			if !edgeIsDerivation(e) {
				continue
			}
			dcs, ok := rs.(*resolutionSetImp).resolutions[e.dependency]
			if !ok {
				continue
			} else if !dcs.HasAnyByName(Notice) {
				continue
			}
			rcs := newLicenseConditionSet(nil)
			rcs.addAll(dcs.rename(Notice, "notice"))
			rmap[e.dependency] = rcs
			if lg.targets[e.dependency].isContainer {
				walk(e.dependency)
			}
		}

	}

	for _, r := range lg.rootFiles {
		cs, ok := rs.(*resolutionSetImp).resolutions[r]
		if !ok {
			continue
		}
		if !cs.HasAnyByName(Notice) {
			continue
		}
		rcs := newLicenseConditionSet(&targetNodeImp{lg, r})
		rcs.addAll(cs.rename(Notice, "notice"))
		rmap[r] = rcs
		if lg.targets[r].isContainer {
			walk(r)
		}
	}
	return &resolutionSetImp{lg, rmap}
}

// ResolveTopDownCondition performs a top-down walk of the LicenseGraph
// resolving all reachable nodes with `condition` attached.
//
// Diverts "notice" and "restricted", if requested, to their special walks.
func ResolveTopDownCondition(graph LicenseGraph, condition string) ResolutionSet {
	// policy requires a special walk for restricted to make source-sharing infectious
	if Restricted.Contains(condition) {
		return ResolveTopDownRestricted(graph, true /* full walk */)
	}

	// policy is to display notices for several license types
	if condition == "notice" {
		return ResolveTopDownNotice(graph)
	}

	conditionName := ConditionNames{condition}

	lg := graph.(*licenseGraphImp)

	rs := ResolveBottomUpConditions(lg)

	rmap := make(map[string]*licenseConditionSetImp)

	var walk func(string)

	walk = func(f string) {
		if _, ok := rmap[f]; ok {
			return
		}
		for _, e := range lg.index[f] {
			if !edgeIsDerivation(e) {
				continue
			}
			dcs, ok := rs.(*resolutionSetImp).resolutions[e.dependency]
			if !ok {
				continue
			} else if !dcs.HasAnyByName(conditionName) {
				continue
			}
			rcs := newLicenseConditionSet(nil)
 			rcs.addAll(dcs.ByName(conditionName))
			rmap[e.dependency] = rcs
			if lg.targets[e.dependency].isContainer {
				walk(e.dependency)
			}
		}

	}

	for _, r := range lg.rootFiles {
		cs, ok := rs.(*resolutionSetImp).resolutions[r]
		if !ok {
			continue
		}
		if !cs.HasAnyByName(conditionName) {
			continue
		}
		rcs := newLicenseConditionSet(&targetNodeImp{lg, r})
		rcs.addAll(cs.ByName(conditionName))
		rmap[r] = rcs
		if lg.targets[r].isContainer {
			walk(r)
		}
	}
	return &resolutionSetImp{lg, rmap}
}

// JoinResolutions returns a new ResolutionSet combining the resolutions from
// multiple resolution sets derived from the same license graph.
//
// e.g. combine "restricted", "reciprocal", and "proprietary" resolutions.
func JoinResolutions(resolutions... ResolutionSet) ResolutionSet {
	rmap := make(map[string]*licenseConditionSetImp)
	var lg *licenseGraphImp
	for _, r := range resolutions {
		if lg == nil {
			lg = r.(*resolutionSetImp).lg
		} else if lg != r.(*resolutionSetImp).lg {
			panic(fmt.Errorf("attempt to join resolutions from multiple graphs"))
		}
		for t, cs := range r.(*resolutionSetImp).resolutions {
			if _, ok := rmap[t]; !ok {
				rmap[t] = cs.Copy().(*licenseConditionSetImp)
				continue
			}
			rmap[t].addAll(cs.Conditions())
		}
	}
	return &resolutionSetImp{lg, rmap}
}
