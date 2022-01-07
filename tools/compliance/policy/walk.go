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

type EdgeContext interface {
	Context(lg *LicenseGraph, path TargetEdgePath, edge *TargetEdge) interface{}
}

type NoEdgeContext struct{}

func (ctx NoEdgeContext) Context(lg *LicenseGraph, path TargetEdgePath, edge *TargetEdge) interface{} {
	return nil
}

type ApplicableConditionsContext struct {
	universe LicenseConditionSet
}

func (ctx ApplicableConditionsContext) Context(lg *LicenseGraph, path TargetEdgePath, edge *TargetEdge) interface{} {
	universe := ctx.universe
	if len(path) > 0 {
		universe = path[len(path)-1].ctx.(LicenseConditionSet)
	}
	return applicableConditions(lg, edge, universe)
}

// VisitNode is called for each root and for each walked dependency node by
// WalkTopDown. When VisitNode returns true, WalkTopDown will proceed to walk
// down the dependences of the node
type VisitNode func(lg *LicenseGraph, target *TargetNode, path TargetEdgePath) bool

// WalkTopDown does a top-down walk of `lg` calling `visit` and descending
// into depenencies when `visit` returns true.
func WalkTopDown(ctx EdgeContext, lg *LicenseGraph, visit VisitNode) {
	path := NewTargetEdgePath(32)

	var walk func(fnode *TargetNode)
	walk = func(fnode *TargetNode) {
		visitChildren := visit(lg, fnode, *path)
		if !visitChildren {
			return
		}
		for _, edge := range fnode.edges {
			var edgeContext interface{}
			if ctx == nil {
				edgeContext = nil
			} else {
				edgeContext = ctx.Context(lg, *path, edge)
			}
			path.Push(edge, edgeContext)
			walk(edge.dependency)
			path.Pop()
		}
	}

	for _, r := range lg.rootNodes {
		path.Clear()
		walk(r)
	}
}

// WalkResolutionsForCondition performs a top-down walk of the LicenseGraph
// resolving all distributed works for `conditions`.
func WalkResolutionsForCondition(lg *LicenseGraph, conditions LicenseConditionSet) ResolutionSet {
	shipped := ShippedNodes(lg)

	// rmap maps 'attachesTo' targets to the `actsOn` targets and applicable conditions
	//
	// rmap is the resulting ResolutionSet
	rmap := make(ResolutionSet)

	cmap := make(map[*TargetNode]LicenseConditionSet)

	result := make(ResolutionSet)
	WalkTopDown(ApplicableConditionsContext{conditions}, lg, func(lg *LicenseGraph, tn *TargetNode, path TargetEdgePath) bool {
		universe := conditions
		if len(path) > 0 {
			universe = path[len(path)-1].ctx.(LicenseConditionSet)
		}

		if universe.IsEmpty() {
			return false
		}
		priorUniverse, alreadyWalked := cmap[tn]
		if alreadyWalked && universe != priorUniverse.Intersection(universe) {
			alreadyWalked = false
			universe = universe.Union(priorUniverse)
		}
		if alreadyWalked {
			pure := true
			for _, p := range path {
				target := p.Target()
				for actsOn, cs := range rmap[tn] {
					rmap[target][actsOn] = cs
				}
				if _, ok := result[tn]; ok && pure {
					if _, ok := result[target]; !ok {
						result[target] = make(ActionSet)
					}
					for actsOn, cs := range result[tn] {
						result[target][actsOn] = cs
					}
					pure = target.IsContainer()
				}
			}
			if pure {
				match := rmap[tn][tn].Intersection(universe)
				if !match.IsEmpty() {
					if _, ok := result[tn]; !ok {
						result[tn] = make(ActionSet)
					}
					result[tn][tn] = match
				}
			}
			return false
		}
		if !shipped.Contains(tn) {
			return false
		}
		if _, ok := rmap[tn]; !ok {
			rmap[tn] = make(ActionSet)
		}
		rmap[tn][tn] = tn.resolution
		cmap[tn] = universe
		cs := tn.resolution
		if !cs.IsEmpty() {
			cs = cs.Intersection(universe)
			pure := true
			for _, p := range path {
				target := p.Target()
				rmap[target][tn] = tn.resolution
				if pure && !cs.IsEmpty() {
					if _, ok := result[target]; !ok {
						result[target] = make(ActionSet)
					}
					result[target][tn] = cs
					pure = target.IsContainer()
				}
			}
			if pure && !cs.IsEmpty() {
				if _, ok := result[tn]; !ok {
					result[tn] = make(ActionSet)
				}
				result[tn][tn] = cs
			}
		}
		return true
	})

	return result
}

// WalkActionsForCondition performs a top-down walk of the LicenseGraph
// resolving all distributed works for `conditions`.
func WalkActionsForCondition(lg *LicenseGraph, conditions LicenseConditionSet) ActionSet {
	shipped := ShippedNodes(lg)

	// rmap maps 'actsOn' targets to the applicable conditions
	//
	// rmap is the resulting ActionSet
	amap := make(ActionSet)
	WalkTopDown(ApplicableConditionsContext{conditions}, lg, func(lg *LicenseGraph, tn *TargetNode, path TargetEdgePath) bool {
		universe := conditions
		if len(path) > 0 {
			universe = path[len(path)-1].ctx.(LicenseConditionSet)
		}
		if universe.IsEmpty() {
			return false
		}
		if _, ok := amap[tn]; ok {
			return false
		}
		if !shipped.Contains(tn) {
			return false
		}
		cs := universe.Intersection(tn.resolution)
		if !cs.IsEmpty() {
			amap[tn] = cs
		}
		return true
	})

	return amap
}
