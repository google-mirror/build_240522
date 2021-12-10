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
	"os"
)

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
func ResolveBottomUpConditions(lg *LicenseGraph) *ResolutionSet {

	// short-cut if already walked and cached
	lg.mu.Lock()
	rs := lg.rsBU
	lg.mu.Unlock()

	if rs != nil {
		return rs
	}

	rs = resolveBottomUp(lg, make(map[*TargetNode]actionSet) /* empty map; no prior resolves */)

	// if not yet cached, save the result
	lg.mu.Lock()
	if lg.rsBU == nil {
		lg.rsBU = rs
	} else {
		// if we end up with 2, release the later for garbage collection
		rs = lg.rsBU
	}
	lg.mu.Unlock()

	fmt.Fprintf(os.Stderr, "bottom-up %d actions\n", len(lg.actions))

	return rs
}

// ResolveTopDownCondtions performs a top-down walk of the LicenseGraph
// resolving all reachable nodes for `condition`. Policy establishes the rules
// for transforming and propagating resolutions down the graph.
//
// e.g. For current policy, none of the conditions propagate from target to
// dependency except restricted. For restricted, the policy is to share the
// source of any libraries linked to restricted code and to provide notice.
func ResolveTopDownConditions(lg *LicenseGraph) *ResolutionSet {

	// short-cut if already walked and cached
	lg.mu.Lock()
	rs := lg.rsTD
	lg.mu.Unlock()

	if rs != nil {
		return rs
	}

	// start with the conditions propagated up the graph
	rs = ResolveBottomUpConditions(lg)

	// rmap maps 'appliesTo' targets to their applicable conditions
	//
	// rmap is the resulting ResolutionSet
	rmap := make(map[*TargetNode]actionSet)

	// cmap contains the set of targets walked as pure aggregates. i.e. containers
	cmap := make(map[*TargetNode]struct{})

	progress := 0

	var walk func(fnode *TargetNode, cs LicenseConditionSet, treatAsAggregate bool)

	walk = func(fnode *TargetNode, cs LicenseConditionSet, treatAsAggregate bool) {
		progress++
		if _, ok := rmap[fnode]; !ok {
			rmap[fnode] = actionSet{lg, &IntervalSet{}}
		}
		rmap[fnode].add(fnode, cs)
		if treatAsAggregate {
			cmap[fnode] = struct{}{}
		}
		// add conditions attached to `fnode`
		rs.resolutions[fnode].VisitAll(func(a resolutionAction) {
			cs = cs.Union(a.cs)
		})
		// for each dependency
		for _, edge := range fnode.edges {
			// dcs holds the dpendency conditions inherited from the target
			dcs := targetConditionsApplicableToDep(lg, edge, cs, treatAsAggregate)
			if dcs.IsEmpty() && !treatAsAggregate {
				continue
			}
			dnode := edge.dependency
			if as, alreadyWalked := rmap[dnode]; alreadyWalked {
				if dcs.Difference(as.conditions()).IsEmpty() {
					// no new conditions

					// pure aggregates never need walking a 2nd time with same conditions
					if treatAsAggregate {
						continue
					}
					// non-aggregates don't need walking as non-aggregate a 2nd time
					if _, asAggregate := cmap[dnode]; !asAggregate {
						continue
					}
					// previously walked as pure aggregate; need to re-walk as non-aggregate
					delete(cmap, dnode)
				}
			}
			// add the conditions to the dependency
			walk(dnode, dcs, treatAsAggregate && dnode.IsContainer())
		}
	}

	// walk each of the roots
	for _, rnode := range lg.rootNodes {
		as, ok := rs.resolutions[rnode]
		if !ok {
			// no conditions in root or transitive closure of dependencies
			continue
		}
		if as.isEmpty() {
			continue
		}

		// add the conditions to the root and its transitive closure
		walk(rnode, NewLicenseConditionSet(), rnode.IsContainer())
	}

	// back-fill any bottom-up conditions on targets missed by top-down walk
	for attachesTo, as := range rs.resolutions {
		if _, ok := rmap[attachesTo]; !ok {
			rmap[attachesTo] = as.copy()
		} else {
			rmap[attachesTo].addSet(as)
		}
	}

	// propagate any new conditions back up the graph
	rs = resolveBottomUp(lg, rmap)

	// if not yet cached, save the result
	lg.mu.Lock()
	if lg.rsTD == nil {
		lg.rsTD = rs
	} else {
		// if we end up with 2, release the later for garbage collection
		rs = lg.rsTD
	}
	lg.mu.Unlock()

	fmt.Fprintf(os.Stderr, "top-down %d actions\n", len(lg.actions))

	return rs
}

// resolveBottomUp implements a bottom-up resolve propagating conditions both
// from the graph, and from a `priors` map of resolutions.
func resolveBottomUp(lg *LicenseGraph, priors map[*TargetNode]actionSet) *ResolutionSet {
	rs := newResolutionSet()

	// cmap contains an entry for every target that was previously walked as a pure aggregate only.
	cmap := make(map[*TargetNode]struct{})

	var walk func(target *TargetNode, treatAsAggregate bool) actionSet

	progress := 0
	walk = func(target *TargetNode, treatAsAggregate bool) actionSet {
		result := actionSet{lg, &IntervalSet{}}
		result.add(target, target.licenseConditions)
		progress++
		if pas, ok := priors[target]; ok {
			result.addSet(pas)
		}
		if preresolved, ok := rs.resolutions[target]; ok {
			if treatAsAggregate {
				result.addSet(preresolved)
				return result
			}
			if _, asAggregate := cmap[target]; !asAggregate {
				result.addSet(preresolved)
				return result
			}
			// previously walked in a pure aggregate context,
			// needs to walk again in non-aggregate context
			delete(cmap, target)
		}
		if treatAsAggregate {
			cmap[target] = struct{}{}
		}

		// add all the conditions from all the dependencies
		for _, edge := range target.edges {
			// walk dependency to get its conditions
			as := walk(edge.dependency, treatAsAggregate && edge.dependency.IsContainer())

			// turn those into the conditions that apply to the target
			as = depActionsApplicableToTarget(lg, edge, as, treatAsAggregate)

			// add them to the result
			result.addSet(as)
		}

		// record these conditions as applicable to the target
		rs.addConditions(target, result)
		if len(priors) == 0 {
			// on the first bottom-up resolve, parents have their own sharing and notice needs
			// on the later resolve, if priors is empty, there will be nothing new to add
			rs.addSelf(target, result.matchingAnySet(ImpliesRestricted))
		}

		// return this up the tree
		return result
	}

	// walk each of the roots
	for _, rnode := range lg.rootNodes {
		_ = walk(rnode, rnode.IsContainer())
	}

	return rs
}
