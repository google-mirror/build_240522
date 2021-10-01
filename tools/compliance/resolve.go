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
	cmap := make(map[string]interface{})

	var walk func(f string, treatAsAggregate bool) *licenseConditionSetImp

	walk = func(f string, treatAsAggregate bool) *licenseConditionSetImp {
		result := newLicenseConditionSet(&targetNodeImp{lg, f})
		if preresolved, ok := rs.resolutions[f]; ok {
			if treatAsAggregate {
				return result.union(preresolved)
			}
			if _, asAggregate := cmap[f]; !asAggregate {
				return result.union(preresolved)
			}
			delete(cmap, f)
		}
		if treatAsAggregate {
			cmap[f] = nil
		}
		for _, e := range lg.index[f] {
			cs := walk(e.dependency, treatAsAggregate && lg.targets[e.dependency].isContainer)
			cs = depConditionsApplicableToTarget(targetEdgeImp{lg, e}, cs, treatAsAggregate)
			result = result.union(cs)
		}
		result = lg.targets[f].licenseConditions.union(result)
		rs.add(f, result)
		return result
	}

	for _, r := range lg.rootFiles {
		cs := walk(r, lg.targets[r].isContainer)
		rs.add(r, lg.targets[r].licenseConditions.union(cs))
	}

	lg.mu.Lock()
	if lg.rs == nil{
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

	rs := ResolveBottomUpConditions(lg)

	rmap := make(map[string]*licenseConditionSetImp)
	cmap := make(map[string]interface{})

	var walk func(f string, cs *licenseConditionSetImp, treatAsAggregate bool)

	walk = func(f string, cs *licenseConditionSetImp, treatAsAggregate bool) {
		if _, ok := rmap[f]; ok {
			rmap[f] = rmap[f].union(cs)
		} else {
			rmap[f] = cs.Copy().(*licenseConditionSetImp)
			if treatAsAggregate {
				cmap[f] = nil
			}
		}
		for _, e := range lg.index[f] {
			dcs, ok := rs.(*resolutionSetImp).resolutions[e.dependency]
			if !ok {
				continue
			}
			dcs = selfConditionsApplicableForConditionName(condition, dcs, treatAsAggregate)
			dcs = dcs.union(targetConditionsApplicableToDep(targetEdgeImp{lg, e}, rmap[f], treatAsAggregate))
			if dcs.IsEmpty() {
				continue
			}
			if pcs, alreadyWalked := rmap[e.dependency]; alreadyWalked {
				if !treatAsAggregate {
					if _, asAggregate := cmap[e.dependency]; asAggregate {
						delete(cmap, e.dependency)
						walk(e.dependency, dcs, false /* treat as aggregate */)
						continue
				        }
				}
				pcs = pcs.Copy().(*licenseConditionSetImp)
				pcs.remove(dcs)
				if pcs.IsEmpty() {
					continue
				}
			}
			walk(e.dependency, dcs, treatAsAggregate && lg.targets[e.dependency].isContainer)
		}
	}

	for _, r := range lg.rootFiles {
		cs, ok := rs.(*resolutionSetImp).resolutions[r]
		if !ok {
			continue
		}
		cs = selfConditionsApplicableForConditionName(condition, cs, lg.targets[r].isContainer)
		if cs.IsEmpty() {
			continue
		}
		walk(r, cs, lg.targets[r].isContainer)
	}

	return &resolutionSetImp{lg, rmap}
}
