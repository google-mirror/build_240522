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
	"sync"
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
func ResolveBottomUpConditions(lg *LicenseGraph) *ResolutionStep {

	// short-cut if already walked and cached
	lg.mu.Lock()
	rs := lg.rsBU

	if rs != nil {
		lg.mu.Unlock()
		rs.wg.Wait()
		return rs
	}
	rs, err := newResolutionStep(lg)
	if err == nil {
		lg.rsBU = rs
	}
	lg.mu.Unlock()
	if err != nil {
		panic(err)
	}

	parent := 0

	resolveBottomUp(rs, parent)

	rs.wg.Done()
	return rs
}

// ResolveTopDownCondtions performs a top-down walk of the LicenseGraph
// resolving all reachable nodes for `condition`. Policy establishes the rules
// for transforming and propagating resolutions down the graph.
//
// e.g. For current policy, none of the conditions propagate from target to
// dependency except restricted. For restricted, the policy is to share the
// source of any libraries linked to restricted code and to provide notice.
func ResolveTopDownConditions(lg *LicenseGraph) *ResolutionStep {

	// short-cut if already walked and cached
	lg.mu.Lock()
	rs := lg.rsTD

	if rs != nil {
		lg.mu.Unlock()
		rs.wg.Wait()
		return rs
	}
	rs, err := newResolutionStep(lg)
	if err == nil {
		lg.rsTD = rs
	}
	lg.mu.Unlock()
	if err != nil {
		panic(err)
	}

	// start with the conditions propagated up the graph
	parent := ResolveBottomUpConditions(lg).index

	// cmap contains the set of targets walked as pure aggregates. i.e. containers
	// (guarded by mu)
	cmap := make(map[*TargetNode]struct{})

	// mu guards concurrent access to cmap
	var mu sync.Mutex

	// wg signals when the first walk is complete
	wg := sync.WaitGroup{}

	var walk func(fnode *TargetNode, cs LicenseConditionSet, treatAsAggregate bool)

	walk = func(fnode *TargetNode, cs LicenseConditionSet, treatAsAggregate bool) {
		defer wg.Done()
		mu.Lock()
		fnode.resolutions[rs.index] |= fnode.resolutions[parent]
		fnode.resolutions[rs.index] |= cs
		if treatAsAggregate {
			cmap[fnode] = struct{}{}
		}
		cs = fnode.resolutions[rs.index]
		mu.Unlock()
		// for each dependency
		for _, edge := range fnode.edges {
			func(edge *TargetEdge) {
				// dcs holds the dpendency conditions inherited from the target
				dcs := targetConditionsApplicableToDep(lg, edge, cs, treatAsAggregate)
				dnode := edge.dependency
				mu.Lock()
				defer mu.Unlock()
				depcs := dnode.resolutions[rs.index]
				if !dcs.IsEmpty() && !depcs.IsEmpty() {
					if dcs.Difference(depcs).IsEmpty() {
						// no new conditions

						// pure aggregates never need walking a 2nd time with same conditions
						if treatAsAggregate {
							return
						}
						// non-aggregates don't need walking as non-aggregate a 2nd time
						if _, asAggregate := cmap[dnode]; !asAggregate {
							return
						}
						// previously walked as pure aggregate; need to re-walk as non-aggregate
						delete(cmap, dnode)
					}
				}
				// add the conditions to the dependency
				wg.Add(1)
				go walk(dnode, dcs, treatAsAggregate && dnode.IsContainer())
			}(edge)
		}
	}

	// walk each of the roots
	for _, rnode := range lg.rootNodes {
		wg.Add(1)
		// add the conditions to the root and its transitive closure
		walk(rnode, NewLicenseConditionSet(), rnode.IsContainer())
	}
	wg.Wait()

	// propagate any new conditions back up the graph
	resolveBottomUp(rs, rs.index /* use own top-down results as parent */)

	rs.wg.Done()
	return rs
}

// resolveBottomUp implements a bottom-up resolve propagating conditions both
// from the graph, and from a `priors` map of resolutions.
func resolveBottomUp(rs *ResolutionStep, parent int) {

	lg := rs.lg

	// cmap indentifies targets previously walked as pure aggregates. i.e. as containers
	// (guarded by mu)
	cmap := make(map[*TargetNode]struct{})
	var mu sync.Mutex

	var walk func(target *TargetNode, treatAsAggregate bool) LicenseConditionSet

	walk = func(target *TargetNode, treatAsAggregate bool) LicenseConditionSet {
		priorWalkResults := func() LicenseConditionSet {
			mu.Lock()
			defer mu.Unlock()

			if !target.resolutions[rs.index].IsEmpty() {
				if treatAsAggregate {
					return target.resolutions[rs.index]
				}
				if _, asAggregate := cmap[target]; !asAggregate {
					return target.resolutions[rs.index]
				}
				// previously walked in a pure aggregate context,
				// needs to walk again in non-aggregate context
				delete(cmap, target)
			} else {
				target.resolutions[rs.index] = target.resolutions[parent]
			}
			if treatAsAggregate {
				cmap[target] = struct{}{}
			}
			return NewLicenseConditionSet()
		}
		if cs := priorWalkResults(); !cs.IsEmpty() {
			return cs
		}

		c := make(chan LicenseConditionSet, len(target.edges))
		// add all the conditions from all the dependencies
		for _, edge := range target.edges {
			go func(edge *TargetEdge) {
				// walk dependency to get its conditions
				cs := walk(edge.dependency, treatAsAggregate && edge.dependency.IsContainer())

				// turn those into the conditions that apply to the target
				cs = depConditionsApplicableToTarget(lg, edge, cs, treatAsAggregate)

				c <- cs
			}(edge)
		}
		mu.Lock()
		cs := target.resolutions[rs.index]
		mu.Unlock()
		for i := 0; i < len(target.edges); i++ {
			cs |= <-c
		}
		mu.Lock()
		target.resolutions[rs.index] |= cs
		mu.Unlock()

		// return conditions up the tree
		return cs
	}

	// walk each of the roots
	for _, rnode := range lg.rootNodes {
		_ = walk(rnode, rnode.IsContainer())
	}
}
