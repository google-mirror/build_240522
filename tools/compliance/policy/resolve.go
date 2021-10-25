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

	// short-cut if already walked and cached
	lg.mu.Lock()
	rs := lg.rsBU
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
				for _, cs := range preresolved {
					result.addSet(cs)
				}
				return result
			}
			if _, asAggregate := cmap[f]; !asAggregate {
				for _, cs := range preresolved {
					result.addSet(cs)
				}
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
	if lg.rsBU == nil {
		lg.rsBU = rs
	}
	lg.mu.Unlock()

	return rs
}

// ResolveTopDownCondtions performs a top-down walk of the LicenseGraph
// resolving all reachable nodes for `condition`. Policy establishes the rules
// for transforming and propagating resolutions down the graph.
//
// e.g. For current policy, none of the conditions propagate from target to
// dependency except restricted. For restricted, the policy is to share the
// source of any libraries linked to restricted code and to provide notice.
func ResolveTopDownConditions(graph LicenseGraph) ResolutionSet {
	lg := graph.(*licenseGraphImp)

	// short-cut if already walked and cached
	lg.mu.Lock()
	rs := lg.rsTD
	lg.mu.Unlock()

	if rs != nil {
		return rs
	}

	// start with the conditions propagated up the graph
	rs = ResolveBottomUpConditions(lg).(*resolutionSetImp)

	// rmap maps 'appliesTo' targets to their applicable conditions
	//
	// rmap + lg is the resulting ResolutionSet
	rmap := make(map[string]actionSet)
	for attachesTo, as := range rs.resolutions {
		rmap[attachesTo] = as.copy()
	}

	path := make([]*dependencyEdge, 0, 10)

	firstNonAggregate := func(f string) string {
		if len(path) < 1 {
			return f
		}
		for i := 0; i < len(path); i++ {
			if !lg.targets[path[i].target].GetIsContainer() {
				return path[i].target
			}
		}
		return path[len(path)-1].target
	}

	var walk func(f string, cs *licenseConditionSetImp, treatAsAggregate bool)

	walk = func(f string, cs *licenseConditionSetImp, treatAsAggregate bool) {
		attachTo := firstNonAggregate(f)
		if _, ok := rmap[attachTo]; !ok {
			rmap[attachTo] = make(actionSet)
		}
		if _, ok := rmap[attachTo][f]; !ok {
			rmap[attachTo][f] = cs.Copy().(*licenseConditionSetImp)
		} else {
			rmap[attachTo][f].addSet(cs)
		}
		// add the bottom-up conditions for `f` to the top-down conditions before pushing into deps
		if as, ok := rs.resolutions[f]; ok {
			cs = cs.Copy().(*licenseConditionSetImp)
			for _, fcs := range as {
				cs.addSet(fcs)
			}
		}
		// for each dependency
		for _, edge := range lg.index[f] {
			path = append(path, edge)
			// add any dependency conditions that come down from the parent target
			// dcs holds the dependency conditions inheritd from the target
			dcs := targetConditionsApplicableToDep(targetEdgeImp{lg, edge}, cs, treatAsAggregate)
			if dcs.IsEmpty() {
				path = path[:len(path)-1]
				continue
			}
			// add the conditions to the dependency
			walk(edge.dependency, dcs, treatAsAggregate && lg.targets[edge.dependency].GetIsContainer())
			path = path[:len(path)-1]
		}
	}

	// walk each of the roots
	for _, r := range lg.rootFiles {
		as, ok := rs.resolutions[r]
		if !ok {
			// no conditions in root or transitive closure of dependencies
			continue
		}
		if as.isEmpty() {
			continue
		}

		path = path[:0]
		// add the conditions to the root and its transitive closure
		walk(r, newLicenseConditionSet(&targetNodeImp{lg, r}), lg.targets[r].GetIsContainer())
	}

	return &resolutionSetImp{lg, rmap}
}

// WalkResolutionsForCondition performs a top-down walk of the LicenseGraph
// resolving all included works for condition `names`.
func WalkResolutionsForCondition(rs ResolutionSet, names ConditionNames) ResolutionSet {
	lg := rs.(*resolutionSetImp).lg

	// rmap maps 'attachesTo' targets to the `actsOn` targets and applicable conditions
	//
	// rmap + lg is the resulting ResolutionSet
	rmap := make(map[string]actionSet)

	var walk func(f string, as actionSet)

	walk = func(f string, as actionSet) {
		if _, ok := rmap[f]; ok {
			rmap[f].addSet(as)
		} else {
			rmap[f] = as.copy()
		}
		// For non-containers, all dependencies handled by applicable conditions.
		if !lg.targets[f].GetIsContainer() {
			return
		}
		// Containers have visible internal structure.

		// for each dependency
		for _, edge := range lg.index[f] {
			// Only need to include what is shipped or derived from.
			if !edgeIsDerivation(edge) {
				continue
			}
			if _, alreadyWalked := rmap[edge.dependency]; alreadyWalked {
				continue
			}
			// get dependency's actions
			das, ok := rs.(*resolutionSetImp).resolutions[edge.dependency]
			if !ok {
				continue
			}
			// filter the actionable conditions by name
			das = das.byName(names)
			if das.isEmpty() {
				continue
			}
			// add the conditions to the dependency
			walk(edge.dependency, das)
		}
	}

	// walk each of the roots
	for _, r := range lg.rootFiles {
		as, ok := rs.(*resolutionSetImp).resolutions[r]
		if !ok {
			// no conditions in root or transitive closure of dependencies
			continue
		}

		// restrict to the requested condition names
		as = as.byName(names)
		if as.isEmpty() {
			continue
		}

		// add the conditions to the root and its transitive closure
		walk(r, as)
	}

	return &resolutionSetImp{lg, rmap}
}
