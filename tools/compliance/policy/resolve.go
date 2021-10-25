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
	"sort"
	"strings"
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
func ResolveBottomUpConditions(graph LicenseGraph) ResolutionSet {
	lg := graph.(*licenseGraphImp)

	// short-cut if already walked and cached
	lg.mu.Lock()
	rs := lg.rs
	lg.mu.Unlock()

	if rs != nil {
		return rs
	}

	// must be indexed for fast lookup
	lg.indexForward()

	rs = newResolutionSet(lg)

	// cmap contains an entry for every target that was previously walked as a pure aggregate only.
	cmap := make(map[string]interface{})

	var walk func(f string, treatAsAggregate bool) *licenseConditionSetImp

	walk = func(f string, treatAsAggregate bool) *licenseConditionSetImp {
		result := newLicenseConditionSet(&targetNodeImp{lg, f}, lg.targets[f].LicenseConditions...)
		if preresolved, ok := rs.resolutions[f]; ok {
			if treatAsAggregate {
				result.addSet(preresolved)
				return result
			}
			if _, asAggregate := cmap[f]; !asAggregate {
				result.addSet(preresolved)
				return result
			}
			// previously walked in a pure aggregate context,
			// needs to walk again in non-aggregate context
			delete(cmap, f)
		}
		if treatAsAggregate {
			cmap[f] = nil
		}

		// add all the conditions from all the dependencies
		for _, edge := range lg.index[f] {
			// walk dependency to get its conditions
			cs := walk(edge.dependency, treatAsAggregate && lg.targets[edge.dependency].GetIsContainer())

			// turn those into the conditions that apply to the target
			cs = depConditionsApplicableToTarget(targetEdgeImp{lg, edge}, cs, treatAsAggregate)

			// add them to the result
			result.addSet(cs)
		}

		// record these conditions as applicable to the target
		rs.add(f, result)

		// return this up the tree
		return result
	}

	// walk each of the roots
	for _, r := range lg.rootFiles {
		_ = walk(r, lg.targets[r].GetIsContainer())
	}

	// if not yet cached, save the result
	lg.mu.Lock()
	if lg.rs == nil {
		lg.rs = rs
	}
	lg.mu.Unlock()

	return rs
}

// ResolveTopDownForCondition performs a top-down walk of the LicenseGraph
// resolving all reachable nodes for `condition`. Policy establishes the rules
// for transforming and propagating resolutions down the graph.
//
// e.g. Current policy is to provide notices for all license types.
//
// The top-down resolve for notice transforms all the permissive and restricted
// etc. conditions to "notice" conditions to meet this policy.
//
// For current policy, none of the conditions propagate from target to
// dependency except restricted. For restricted, the policy is to share the
// source of any libraries linked to restricted code.
func ResolveTopDownForCondition(graph LicenseGraph, condition string) ResolutionSet {
	lg := graph.(*licenseGraphImp)

	// start with the conditions propagated up the graph
	rs := ResolveBottomUpConditions(lg)

	// rmap maps 'appliesTo' targets to their applicable conditions
	//
	// rmap + lg is the resulting ResolutionSet
	rmap := make(map[string]*licenseConditionSetImp)

	// wmap contains an entry for every target+conditions that was previously walked.
	wmap := make(map[string]interface{})
	// cmap contains an entry for every target+conditions that was previously walked as a pure aggregate only.
	cmap := make(map[string]interface{})

	walkKey := func(t string, cs *licenseConditionSetImp) string {
		conditions := make([]string, 0, cs.Count())
		for cname := range cs.conditions {
			conditions = append(conditions, cname)
		}
		// Sort conditions for consistent, reproducible order.
		sort.Strings(conditions)
		return fmt.Sprintf("%s:%s", t, strings.Join(conditions, ":"))
	}

	var walk func(f string, cs *licenseConditionSetImp, treatAsAggregate bool)

	walk = func(f string, cs *licenseConditionSetImp, treatAsAggregate bool) {
		wstring := walkKey(f, cs)

		// add the conditions from above to the target
		if _, walkedBefore := wmap[wstring]; walkedBefore {
			rmap[f].addSet(cs)
		} else {
			if _, ok := rmap[f]; ok {
				rmap[f].addSet(cs)
			} else {
				rmap[f] = cs.Copy().(*licenseConditionSetImp)
			}
			if treatAsAggregate {
				cmap[wstring] = nil
			}
		}
		// for each dependency
		for _, edge := range lg.index[f] {
			// get dependency's conditions
			dcs, ok := rs.(*resolutionSetImp).resolutions[edge.dependency]
			if !ok {
				dcs = newLicenseConditionSet(&targetNodeImp{lg, edge.dependency})
			}
			// adjust the dependency's own conditions per policy
			dcs = selfConditionsApplicableForConditionName(condition, dcs, treatAsAggregate)
			// add any dependencies that come down from the parent target
			dcs.addSet(targetConditionsApplicableToDep(targetEdgeImp{lg, edge}, rmap[f], treatAsAggregate))
			if dcs.IsEmpty() {
				continue
			}
			wstring = walkKey(edge.dependency, dcs)
			// get conditions from prior walks (only need to process if new conditions)
			if _, alreadyWalked := wmap[wstring]; alreadyWalked {
				if !treatAsAggregate {
					if _, asAggregate := cmap[wstring]; asAggregate {
						// was previously walked as pure aggregate,
						// walk again in non-aggregate context
						delete(cmap, wstring)
						walk(edge.dependency, dcs, false /* treat as aggregate */)
					}
				}
				continue
			}
			// add the conditions to the dependency
			walk(edge.dependency, dcs, treatAsAggregate && lg.targets[edge.dependency].GetIsContainer())
		}
	}

	// walk each of the roots
	for _, r := range lg.rootFiles {
		cs, ok := rs.(*resolutionSetImp).resolutions[r]
		if !ok {
			// no conditions in root or transitive closure of dependencies
			continue
		}

		// adjust the root's own conditions per policy
		cs = selfConditionsApplicableForConditionName(condition, cs, lg.targets[r].GetIsContainer())
		if cs.IsEmpty() {
			continue
		}

		// add the conditions to the root and its transitive closure
		walk(r, cs, lg.targets[r].GetIsContainer())
	}

	return &resolutionSetImp{lg, rmap}
}
