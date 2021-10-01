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
	"sync"
)

// NewLicenseConditionSet creates a new instance of LicenseConditionSet.
func NewLicenseConditionSet(conditions... LicenseCondition) LicenseConditionSet {
	cs := &licenseConditionSetImp{nil, make(map[string]map[string]interface{})}
	cs.addAll(conditions)
	return cs
}

// licenseConditionImp implements LicenseCondition.
type licenseConditionImp struct {
	name string
	origin targetNodeImp
}

// Name returns the name of the condition. e.g. "restricted" or "notice"
func (c licenseConditionImp) Name() string {
	return c.name
}

// Origin identifies the TargetNode where the condition originates.
func (c licenseConditionImp) Origin() TargetNode {
	return c.origin
}


// licenseConditionSetImp implements LicenseConditionSet.
type licenseConditionSetImp struct {
	lg *licenseGraphImp
	conditions map[string]map[string]interface{}
}

// Add makes all `conditions` members of the set if they were not previously.
func (cs *licenseConditionSetImp) Add(conditions... LicenseCondition) {
	cs.addAll(conditions)
}

// ByName returns a list of the conditions in the set matching `name`.
func (cs *licenseConditionSetImp) ByName(name... ConditionNames) []LicenseCondition {
	l := make([]LicenseCondition, 0)
	for _, cn := range name {
		for _, n := range cn {
			if origins, ok := cs.conditions[n]; ok {
				for t := range origins {
					l = append(l, licenseConditionImp{n, targetNodeImp{cs.lg, t}})
				}
			}
		}
	}
	return l
}

// ByOrigin returns all of the conditions that originate at `target` regardless of name.
func (cs *licenseConditionSetImp) ByOrigin(target TargetNode) []LicenseCondition {
	l := make([]LicenseCondition, 0)
	for name, origins := range cs.conditions {
		if _, ok := origins[target.Name()]; ok {
			l = append(l, licenseConditionImp{name, targetNodeImp{cs.lg, target.Name()}})
		}
	}
	return l
}

// HasAnyByOrigin returns true if the set contains any conditions originating at `target` regarless of name.
func (cs *licenseConditionSetImp) HasAnyByOrigin(target TargetNode) bool {
	for _, origins := range cs.conditions {
		if _, ok := origins[target.Name()]; ok {
			return true
		}
	}
	return false
}

// Conditions returns a list of all the conditions in the set.
func (cs *licenseConditionSetImp) Conditions() []LicenseCondition {
	l := make([]LicenseCondition, 0, len(cs.conditions))
	for c, origins := range cs.conditions {
		for t := range origins {
			l = append(l, licenseConditionImp{c, targetNodeImp{cs.lg, t}})
		}
	}
	return l
}

// Copy creates a new LicenseCondition variable with the same value.
func (cs *licenseConditionSetImp) Copy() LicenseConditionSet {
	other := licenseConditionSetImp{cs.lg, make(map[string]map[string]interface{})}
	for name, origins := range cs.conditions {
		for t := range origins {
			other.add(name, targetNodeImp{cs.lg, t})
		}
	}
	return &other
}

// HasAnyByName returns true if the set contains any conditions matching `name` originating at any target.
func (cs *licenseConditionSetImp) HasAnyByName(name... ConditionNames) bool {
	for _, cn := range name {
		for _, n := range cn {
			if origins, ok := cs.conditions[n]; ok {
				if len(origins) > 0 {
					return true
				}
			}
		}
	}
	return false
}

// HasCondition returns true if the set contains any condition matching both `name` and `origin`.
func (cs *licenseConditionSetImp) HasCondition(name ConditionNames, origin TargetNode) bool {
	timp := origin.(targetNodeImp)
	if cs.lg == nil {
		return false
	} else if timp.lg == nil {
		timp.lg = cs.lg
	} else if cs.lg != timp.lg {
		panic(fmt.Errorf("attempt to query license conditions from different graph"))
	}
	for _, n := range name {
		if targets, ok := cs.conditions[n]; ok {
			_, isPresent := targets[timp.file]
			if isPresent {
				return true
			}
		}
	}
	return false
}

// IsEmpty returns true when the sent contains zero conditions.
func (cs *licenseConditionSetImp) IsEmpty() bool {
	for _, origins := range cs.conditions {
		if len(origins) > 0 {
			return false
		}
	}
	return true
}

// compliance-only licenseConditionSetImp methods

// newLicenseConditionSet constructs a set of conditions
func newLicenseConditionSet(origin *targetNodeImp, condition... string) *licenseConditionSetImp {
	cs := &licenseConditionSetImp{nil, make(map[string]map[string]interface{})}
	if origin != nil {
		cs.lg = origin.lg
	} else if len(condition) > 0 {
		panic(fmt.Errorf("attempt to add conditions to nil target"))
	}
	for _, c := range condition {
		cs.conditions[c][origin.file] = nil
	}
	return cs
}

// add changes the set to include `condition` if it does not already
func (cs *licenseConditionSetImp) add(condition string, origins... targetNodeImp) {
	if len(origins) == 0 {
		return
	}
	if _, ok := cs.conditions[condition]; !ok {
		cs.conditions[condition] = make(map[string]interface{})
	}

	for _, t := range origins {
		if cs.lg == nil {
			cs.lg = t.lg
		} else if t.lg == nil {
			t.lg = cs.lg
		} else if t.lg != cs.lg {
			panic(fmt.Errorf("attempting to combine license conditions from different graphs"))
		}
		found := false
		for othert := range cs.conditions[condition] {
			if t.file == othert {
				found = true
				break
			}
		}
		if !found {
			cs.conditions[condition][t.file] = nil
		}
	}
}

// addAll modifies `cs` to include all of the `conditions` if they were not previously members.
func (cs *licenseConditionSetImp) addAll(conditions []LicenseCondition) {
	if len(conditions) == 0 {
		return
	}
	for _, c := range conditions {
		cimp := c.(licenseConditionImp)
		if cs.lg == nil {
			cs.lg = cimp.origin.lg
		} else if cimp.origin.lg != cs.lg {
			panic(fmt.Errorf("attempting to combine license conditions from different graphs"))
		}
		if _, ok := cs.conditions[cimp.name]; !ok {
			cs.conditions[cimp.name] = make(map[string]interface{})
		}
		cs.conditions[cimp.name][cimp.origin.file] = nil
	}
}

// removeAllByName changes the set to delete all conditions matching `name`.
func (cs *licenseConditionSetImp) removeAllByName(name... ConditionNames) {
	for _, cn := range name {
		for _, n := range cn {
			delete(cs.conditions, n)
		}
	}
}

// removeAllByTarget changes the set to delete all conditions that apply to target `file`.
func (cs *licenseConditionSetImp) removeAllByTarget(file string) {
	for c := range cs.conditions {
		delete(cs.conditions[c], file)
	}
}

// remove changes the set to delete all conditions also present in `other`.
func (cs *licenseConditionSetImp) remove(other *licenseConditionSetImp) {
	for c, targets := range other.conditions {
		if _, isPresent := cs.conditions[c]; !isPresent {
			continue
		}
		for t := range targets {
			delete(cs.conditions[c], t)
		}
	}
}

// rename behaves similar to ByName except it limits the search to a single ConditionNames, and it
// changes the name of each condition in the output to `newName`.
func (cs *licenseConditionSetImp) rename(name ConditionNames, newName string) []LicenseCondition {
	l := make([]LicenseCondition, 0)
	for _, n := range name {
		if origins, ok := cs.conditions[n]; ok {
			for t := range origins {
				l = append(l, licenseConditionImp{newName, targetNodeImp{cs.lg, t}})
			}
		}
	}
	return l
}

// union returns a new set calculated as the union of `cs` with some `other` set.
func (cs *licenseConditionSetImp) union(other *licenseConditionSetImp) *licenseConditionSetImp {
	if cs.lg == nil {
		cs.lg = other.lg
	} else if other.lg == nil {
		other.lg = cs.lg
	} else if cs.lg != other.lg {
		panic(fmt.Errorf("attempt to union condition sets from different graphs"))
	}
	result := &licenseConditionSetImp{cs.lg, make(map[string]map[string]interface{})}
	for c, origins := range cs.conditions {
		for t := range origins {
			if _, ok := result.conditions[c]; !ok {
				result.conditions[c] = make(map[string]interface{})
			}
			result.conditions[c][t] = nil
		}
	}
	for c, origins := range other.conditions {
		for t := range origins {
			if _, ok := result.conditions[c]; !ok {
				result.conditions[c] = make(map[string]interface{})
			}
			result.conditions[c][t] = nil
		}
	}
	return result
}

// conditionNamesArray implements a `contains` predicate for arrays of ConditionNames
type conditionNamesArray []ConditionNames

func (cn conditionNamesArray) contains(name string) bool {
	for _, c := range cn {
		if c.Contains(name) {
			return true
		}
	}
	return false
}


// filter returns a new liceneseConditionSetImp containing the subset of conditions matching `names`.
func (cs *licenseConditionSetImp) filter(names... ConditionNames) *licenseConditionSetImp {
	result := &licenseConditionSetImp{cs.lg, make(map[string]map[string]interface{})}
	for c, origins := range cs.conditions {
		if !conditionNamesArray(names).contains(c) {
			continue
		}
		for t := range origins {
			if _, ok := result.conditions[c]; !ok {
				result.conditions[c] = make(map[string]interface{})
			}
			result.conditions[c][t] = nil
		}
	}
	return result
}

// resolutionSetImp implements ResolutionSet.
type resolutionSetImp struct {
	lg *licenseGraphImp
	resolutions map[string]*licenseConditionSetImp
}

// Targets returns the list of targets with applicable licens conditions attached.
func (rs *resolutionSetImp) Targets() []TargetNode {
	targets := make([]TargetNode, 0, len(rs.resolutions))
	for t := range rs.resolutions {
		targets = append(targets, targetNodeImp{rs.lg, t})
	}
	return targets
}

// Conditions returns the set of conditions applied to `target` or false if `target` not in set.
func (rs *resolutionSetImp) Conditions(target TargetNode) (LicenseConditionSet, bool) {
	timp := target.(targetNodeImp)
	if rs.lg == nil {
		rs.lg = timp.lg
	} else if timp.lg == nil {
		timp.lg = rs.lg
	} else if rs.lg != timp.lg {
		panic(fmt.Errorf("attempt to query target resolutions for wrong graph"))
	}
	cs, ok := rs.resolutions[timp.file]
	return cs.Copy(), ok
}

// HasTarget returns true if `target` appears in the set.
func (rs *resolutionSetImp) HasTarget(target TargetNode) bool {
	timp := target.(targetNodeImp)
	if rs.lg == nil {
		rs.lg = timp.lg
	} else if timp.lg == nil {
		timp.lg = rs.lg
	} else if rs.lg != timp.lg {
		panic(fmt.Errorf("attempt to query resolved targets for wrong graph"))
	}
	return rs.hasTarget(timp.file)
}

// HasAnyByName returns true if `target` appears in the set with any conditions matching `name`.
func (rs resolutionSetImp) HasAnyByName(target TargetNode, name... ConditionNames) bool {
	timp := target.(targetNodeImp)
	if rs.lg == nil {
		rs.lg = timp.lg
	} else if timp.lg == nil {
		timp.lg = rs.lg
	} else if rs.lg != timp.lg {
		panic(fmt.Errorf("attempt to query target resolutions for wrong graph"))
	}
	return rs.hasAnyByName(timp.file, name...)
}

// HasAllByName returns true if `target` appears in the set with conditions matching every element of `name`.
func (rs resolutionSetImp) HasAllByName(target TargetNode, name... ConditionNames) bool {
	timp := target.(targetNodeImp)
	if rs.lg == nil {
		rs.lg = timp.lg
	} else if timp.lg == nil {
		timp.lg = rs.lg
	} else if rs.lg != timp.lg {
		panic(fmt.Errorf("attempt to query target resolutions for wrong graph"))
	}
	return rs.hasAllByName(timp.file, name...)
}

// compliance-only resolutionSetImp methods

// newResolutionSet constructs a new, empty instance of resolutionSetImp
func newResolutionSet(lg *licenseGraphImp) *resolutionSetImp {
	return &resolutionSetImp{lg, make(map[string]*licenseConditionSetImp)}
}

// add applies all of the license conditions in `cs` to `file` if not already applied.
func (rs *resolutionSetImp) add(file string, cs *licenseConditionSetImp) {
	if r, ok := rs.resolutions[file]; ok {
		rs.resolutions[file] = r.union(cs)
	} else {
		rs.resolutions[file] = cs
	}
}

// hasTarget returns true if `file` appears as one of the targets in the set.
func (rs *resolutionSetImp) hasTarget(file string) bool {
	_, isPresent := rs.resolutions[file]
	return isPresent
}

// hasAnyByName returns true if the target for `file` has at least 1 condition matching `name`.
func (rs *resolutionSetImp) hasAnyByName(file string, name... ConditionNames) bool {
	cs, isPresent := rs.resolutions[file]
	if !isPresent {
		return false
	}
	for _, cn := range name {
		for _, n := range cn {
			_, isPresent = cs.conditions[n]
			if isPresent {
				return true
			}
		}
	}
	return false
}

// hasAllByName returns true if the target for `file` has at least 1 condition for each element of `name`.
func (rs *resolutionSetImp) hasAllByName(file string, name... ConditionNames) bool {
	cs, isPresent := rs.resolutions[file]
	if !isPresent {
		return false
	}
	for _, cn := range name {
		found := false
		for _, n := range cn {
			_, isPresent = cs.conditions[n]
			if isPresent {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// targetNodeImp implements TargetNode
type targetNodeImp struct {
	lg   *licenseGraphImp
	file string
}

// Name returns the name identifying the target node. i.e. the path to the corresponding license metadata file
func (tn targetNodeImp) Name() string {
	return tn.lg.targets[tn.file].name
}

// IsContainer returns true if the target represents a container that merely aggregates other targets.
func (tn targetNodeImp) IsContainer() bool {
	return tn.lg.targets[tn.file].isContainer
}

// Projects returns the list of projects that define the target. (unordered)
func (tn targetNodeImp) Projects() []string {
	return append([]string{}, tn.lg.targets[tn.file].projects...)
}

// LicenseTexts returns the list of paths to license text for the target. (unordered)
func (tn targetNodeImp) LicenseTexts() []string {
	return append([]string{}, tn.lg.targets[tn.file].licenseTexts...)
}

// targetEdgeImp implements TargetEdge
type targetEdgeImp struct {
	lg *licenseGraphImp
	e *dependencyEdge
}

// Target returns the depending end of the edge.
func (e targetEdgeImp) Target() TargetNode {
	return targetNodeImp{e.lg, e.e.target}
}

// Dependency returns the depended-on end of the edge.
func (e targetEdgeImp) Dependency() TargetNode {
	return targetNodeImp{e.lg, e.e.dependency}
}

// Annotations describe the type of edge.
func (e targetEdgeImp) Annotations() TargetEdgeAnnotations {
	return e.e.annotations
}

// compliance-only targetEdgeImp methods

// isDynamicLink returns true if the edge represents a shared or dynamic link at runtime.
func (e targetEdgeImp) isDynamicLink() bool {
	return edgeIsDynamicLink(e.e)
}

// isDerivation returns true if the edge represents a dependency incorporated into the target as a derivative work.
func (e targetEdgeImp) isDerivation() bool {
	return edgeIsDerivation(e.e)
}


// targetEdgeAnnotationsImp implements TargetEdgeAnnotations
type targetEdgeAnnotationsImp map[string]interface{}

// HasAnnotation returns true if `ann` is attached to the edge.
func (ea targetEdgeAnnotationsImp) HasAnnotation(ann string) bool {
	_, ok := ea[ann]
	return ok
}

// ListAnnotations returns the list of annotation names attached to the edge. (unordered)
func (ea targetEdgeAnnotationsImp) ListAnnotations() []string {
	l := make([]string, 0, len(ea))
	for ann := range ea {
		l = append(l, ann)
	}
	return l
}


// licenseGraphImp implements the LicenseGraph interface.
type licenseGraphImp struct {
	// rootFiles identifies the original set of files to read (immutable)
	rootFiles []string
	// targets identifies the entire set of target node files (guarded by mu)
	targets map[string]*targetNode
	// edges lists the target edges from dependent to dependency (guarded by mu)
	edges []*dependencyEdge
	// index facilitates looking up edges from targets (creation guarded by my)
	index map[string][]*dependencyEdge
	// rs caches the results of a full graph resolution (creation guarded by mu)
	rs *resolutionSetImp
	// mu guards against concurrent update
	mu sync.Mutex
}

// indexForward guarantees the `index` map is populated to look up edges by `target`
func (lg *licenseGraphImp) indexForward() {
	lg.mu.Lock()
	defer func() {
		lg.mu.Unlock()
	}()

	if lg.index != nil {
		return
	}

	lg.index = make(map[string][]*dependencyEdge)
	for _, e := range lg.edges {
		if _, ok := lg.index[e.target]; ok {
			lg.index[e.target] = append(lg.index[e.target], e)
		} else {
			lg.index[e.target] = []*dependencyEdge{e}
		}
	}
}


// TargetNode returns the target node identified by `name`.
func (lg *licenseGraphImp) TargetNode(name string) TargetNode {
	if _, ok := lg.targets[name]; !ok {
		panic(fmt.Errorf("target node %q missing from graph", name))
	}
	return targetNodeImp{lg, name}
}

// HasTargetNode returns true if a target node identified by `name` appears in the graph.
func (lg *licenseGraphImp) HasTargetNode(name string) bool {
	_, isPresent := lg.targets[name]
	return isPresent
}

// Edges returns the list of edges in the graph. (unordered)
func (lg *licenseGraphImp) Edges() []TargetEdge {
	edges := make([]TargetEdge, 0, len(lg.edges))
	for _, e := range lg.edges {
		edges = append(edges, targetEdgeImp{lg, e})
	}
	return edges
}

// compliance-only licenseGraphImp methods

// newLicenseGraphImp constructs a new instance of licenseGraphImp.
func newLicenseGraphImp() *licenseGraphImp {
	return &licenseGraphImp{
		rootFiles: []string{},
		targets: make(map[string]*targetNode),
	}
}
