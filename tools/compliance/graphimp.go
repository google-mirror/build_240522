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
	appliesTo targetNodeImp
}

func (c licenseConditionImp) Name() string {
	return c.name
}

func (c licenseConditionImp) AppliesTo() TargetNode {
	return c.appliesTo
}


// licenseConditionSetImp implements LicenseConditionSet.
type licenseConditionSetImp struct {
	lg *licenseGraphImp
	conditions map[string]map[string]interface{}
}

func (cs *licenseConditionSetImp) Add(conditions... LicenseCondition) {
	cs.addAll(conditions)
}

func (cs *licenseConditionSetImp) ByName(name string) []LicenseCondition {
	l := make([]LicenseCondition, 0)
	if appliesTo, ok := cs.conditions[name]; ok {
		for t := range appliesTo {
			l = append(l, licenseConditionImp{name, targetNodeImp{cs.lg, t}})
		}
	}
	return l
}

func (cs *licenseConditionSetImp) Conditions() []LicenseCondition {
	l := make([]LicenseCondition, 0, len(cs.conditions))
	for c, appliesTo := range cs.conditions {
		for t := range appliesTo {
			l = append(l, licenseConditionImp{c, targetNodeImp{cs.lg, t}})
		}
	}
	return l
}

func (cs *licenseConditionSetImp) Copy() LicenseConditionSet {
	other := licenseConditionSetImp{cs.lg, make(map[string]map[string]interface{})}
	for name, appliesTo := range cs.conditions {
		for t := range appliesTo {
			other.add(name, targetNodeImp{cs.lg, t})
		}
	}
	return &other
}

func (cs *licenseConditionSetImp) HasAnyByName(name string) bool {
	if appliesTo, ok := cs.conditions[name]; ok {
		if len(appliesTo) > 0 {
			return true
		}
	}
	return false
}

func (cs *licenseConditionSetImp) HasCondition(name string, appliesTo TargetNode) bool {
	timp := appliesTo.(targetNodeImp)
	if cs.lg == nil {
		return false
	} else if timp.lg == nil {
		timp.lg = cs.lg
	} else if cs.lg != timp.lg {
		panic(fmt.Errorf("attempt to query license conditions from different graph"))
	}
	if targets, ok := cs.conditions[name]; ok {
		_, isPresent := targets[timp.file]
		return isPresent
	}
	return false
}

func (cs *licenseConditionSetImp) IsEmpty() bool {
	isEmpty := true
	for _, appliesTo := range cs.conditions {
		if len(appliesTo) > 0 {
			isEmpty = false
			break
		}
	}
	return isEmpty
}

// compliance-only licenseConditionSetImp methods

// newLicenseConditionSet constructs a set of conditions
func newLicenseConditionSet(appliesTo *targetNodeImp, condition... string) *licenseConditionSetImp {
	cs := &licenseConditionSetImp{nil, make(map[string]map[string]interface{})}
	if appliesTo != nil {
		cs.lg = appliesTo.lg
	} else if len(condition) > 0 {
		panic(fmt.Errorf("attempt to add conditions to nil target"))
	}
	for _, c := range condition {
		cs.conditions[c][appliesTo.file] = nil
	}
	return cs
}

// add changes the set to include `condition` if it does not already
func (cs *licenseConditionSetImp) add(condition string, appliesTo... targetNodeImp) {
	if len(appliesTo) == 0 {
		return
	}
	if _, ok := cs.conditions[condition]; !ok {
		cs.conditions[condition] = make(map[string]interface{})
	}

	for _, t := range appliesTo {
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

func (cs *licenseConditionSetImp) addAll(conditions []LicenseCondition) {
	if len(conditions) == 0 {
		return
	}
	for _, c := range conditions {
		cimp := c.(licenseConditionImp)
		if cs.lg == nil {
			cs.lg = cimp.appliesTo.lg
		} else if cimp.appliesTo.lg != cs.lg {
			panic(fmt.Errorf("attempting to combine license conditions from different graphs"))
		}
		if _, ok := cs.conditions[cimp.name]; !ok {
			cs.conditions[cimp.name] = make(map[string]interface{})
		}
		cs.conditions[cimp.name][cimp.appliesTo.file] = nil
	}
}

// removeAllByName changes the set to delete all conditions matching `name`.
func (cs *licenseConditionSetImp) removeAllByName(name string) {
	delete(cs.conditions, name)
}

// removeAllByTarget changes the set to delete all conditions that apply to target `file`.
func (cs *licenseConditionSetImp) removeAllByTarget(file string) {
	for c := range cs.conditions {
		delete(cs.conditions[c], file)
	}
}

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
	for c, appliesTo := range cs.conditions {
		for t := range appliesTo {
			if _, ok := result.conditions[c]; !ok {
				result.conditions[c] = make(map[string]interface{})
			}
			result.conditions[c][t] = nil
		}
	}
	for c, appliesTo := range other.conditions {
		for t := range appliesTo {
			if _, ok := result.conditions[c]; !ok {
				result.conditions[c] = make(map[string]interface{})
			}
			result.conditions[c][t] = nil
		}
	}
	return result
}

func (cs *licenseConditionSetImp) filter(condition string) *licenseConditionSetImp {
	result := &licenseConditionSetImp{cs.lg, make(map[string]map[string]interface{})}
	for c, appliesTo := range cs.conditions {
		if c != condition {
			continue
		}
		for t := range appliesTo {
			if _, ok := result.conditions[c]; !ok {
				result.conditions[c] = make(map[string]interface{})
			}
			result.conditions[c][t] = nil
		}
	}
	return result
}

// targetResolutionImp implements TargetResolution.
type targetResolutionImp struct {
	target targetNodeImp
	conditions licenseConditionSetImp
}

func (tr targetResolutionImp) Target() TargetNode {
	return tr.target
}

func (tr targetResolutionImp) Conditions() LicenseConditionSet {
	return tr.conditions.Copy()
}


// resolutionSetImp implements ResolutionSet.
type resolutionSetImp struct {
	lg *licenseGraphImp
	resolutions map[string]*licenseConditionSetImp
}

func (rs *resolutionSetImp) Targets() []TargetNode {
	targets := make([]TargetNode, 0, len(rs.resolutions))
	for t := range rs.resolutions {
		targets = append(targets, targetNodeImp{rs.lg, t})
	}
	return targets
}

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

func (rs *resolutionSetImp) HasTarget(t TargetNode) bool {
	timp := t.(targetNodeImp)
	if rs.lg == nil {
		rs.lg = timp.lg
	} else if timp.lg == nil {
		timp.lg = rs.lg
	} else if rs.lg != timp.lg {
		panic(fmt.Errorf("attempt to query resolved targets for wrong graph"))
	}
	return rs.hasTarget(timp.file)
}

// compliance-only resolutionSetImp methods

func newResolutionSet(lg *licenseGraphImp) *resolutionSetImp {
	return &resolutionSetImp{lg, make(map[string]*licenseConditionSetImp)}
}

func (rs *resolutionSetImp) add(file string, cs *licenseConditionSetImp) {
	if r, ok := rs.resolutions[file]; ok {
		rs.resolutions[file] = r.union(cs)
	} else {
		rs.resolutions[file] = cs
	}
}

func (rs *resolutionSetImp) hasTarget(file string) bool {
	_, isPresent := rs.resolutions[file]
	return isPresent
}


// targetNodeImp implements TargetNode
type targetNodeImp struct {
	lg   *licenseGraphImp
	file string
}

func (tn targetNodeImp) Name() string {
	return tn.lg.targets[tn.file].name
}

func (tn targetNodeImp) Projects() []string {
	return append([]string{}, tn.lg.targets[tn.file].projects...)
}


// targetSetImp implements TargetSet
type targetSetImp struct {
	lg    *licenseGraphImp
	files map[string]interface{}
}

func (ts targetSetImp) Add(node... TargetNode) {
	for _, tn := range node {
		imp := tn.(targetNodeImp)
		if ts.lg != imp.lg {
			panic(fmt.Errorf("attempt to Add target node from different metadata"))
		}
		ts.files[imp.file] = nil
	}
}

func (ts targetSetImp) Targets() []TargetNode {
	var files []TargetNode
	for f := range ts.files {
		files = append(files, targetNodeImp{ts.lg, f})
	}
	return files
}

func (ts targetSetImp) Projects() []string {
	pset := make(map[string]interface{})
	for f := range ts.files {
		for _, p := range ts.lg.targets[f].projects {
			pset[p] = nil
		}
	}
	projects := make([]string, 0, len(pset))
	for p := range pset {
		projects = append(projects, p)
	}
	sort.Strings(projects)
	return projects
}


// targetEdgeImp implements TargetEdge
type targetEdgeImp struct {
	lg *licenseGraphImp
	e *dependencyEdge
}

func (e targetEdgeImp) Target() TargetNode {
	return targetNodeImp{e.lg, e.e.target}
}

func (e targetEdgeImp) Dependency() TargetNode {
	return targetNodeImp{e.lg, e.e.dependency}
}

func (e targetEdgeImp) Annotations() TargetEdgeAnnotations {
	return e.e.annotations
}

// compliance-only targetEdgeImp methods

// returns true if the edge represents a dynamic link at runtime
func (e targetEdgeImp) isDynamicLink() bool {
	return edgeIsDynamicLink(e.e)
}

// returns true if the edge represents a dependency that is incorporated into the target as a derivative work
func (e targetEdgeImp) isDerivativeOf() bool {
	return edgeIsDerivativeOf(e.e)
}


// targetEdgeAnnotationsImp implements TargetEdgeAnnotations
type targetEdgeAnnotationsImp map[string]interface{}

func (ea targetEdgeAnnotationsImp) HasAnnotation(ann string) bool {
	_, ok := ea[ann]
	return ok
}

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
	// reverse facilitates looking up edges from dependencies (creation guarded by mu)
	reverse map[string][]*dependencyEdge
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

// indexReverse guarantees the `reverse` map is populated to look up edges by `dependency`
func (lg *licenseGraphImp) indexReverse() {
	lg.mu.Lock()
	defer func() {
		lg.mu.Unlock()
	}()

	if lg.reverse != nil {
		return
	}

	lg.reverse = make(map[string][]*dependencyEdge)
	for _, e := range lg.edges {
		if _, ok := lg.reverse[e.dependency]; ok {
			lg.reverse[e.dependency] = append(lg.reverse[e.dependency], e)
		} else {
			lg.reverse[e.dependency] = []*dependencyEdge{e}
		}
	}
}

func (lg *licenseGraphImp) TargetNode(name string) TargetNode {
	if _, ok := lg.targets[name]; !ok {
		panic(fmt.Errorf("target node %q missing from graph", name))
	}
	return targetNodeImp{lg, name}
}

func (lg *licenseGraphImp) HasTargetNode(name string) bool {
	_, isPresent := lg.targets[name]
	return isPresent
}

// AnyPath returns a TargetPath from target to dependency or nil if none exist.
func (lg *licenseGraphImp) AnyPath(target, dependency TargetNode) TargetPath {
	lg.indexForward()

	timp := target.(targetNodeImp)
	dimp := dependency.(targetNodeImp)

	stack := make([]string, 0)
	stack = append(stack, timp.file)

	path := make(TargetPath, 0)

	index := []int{0}
	for {
		edges := lg.index[stack[len(stack)-1]]
		i := index[len(index)-1]
		if i >= len(edges) {
			stack = stack[:len(stack)-1]
			if len(stack) < 1 {
				break
			}
			path = path[:len(path)-1]
			index = index[:len(index)-1]
			index[len(index)-1]++
			continue
		}
		stack = append(stack, edges[i].dependency)
		path = append(path, targetEdgeImp{lg, edges[i]})
		if dimp.file == stack[len(stack)-1] {
			return path
		}
		index = append(index, 0)
	}
	return nil
}

// AllPaths returns a slice of TargetPath from target to dependency with an entry
// for each distinct path from target to dependency.
func (lg *licenseGraphImp) AllPaths(target, dependency TargetNode) []TargetPath {
	lg.indexForward()

	timp := target.(targetNodeImp)
	dimp := dependency.(targetNodeImp)

	paths := make([]TargetPath, 0)

	stack := make([]string, 0)
	stack = append(stack, timp.file)

	path := make(TargetPath, 0)

	index := []int{0}
	for {
		edges := lg.index[stack[len(stack)-1]]
		i := index[len(index)-1]
		if i >= len(edges) {
			stack = stack[:len(stack)-1]
			if len(stack) < 1 {
				break
			}
			path = path[:len(path)-1]
			index = index[:len(index)-1]
			index[len(index)-1]++
			continue
		}
		stack = append(stack, edges[i].dependency)
		path = append(path, targetEdgeImp{lg, edges[i]})
		if dimp.file == stack[len(stack)-1] {
			c := append(TargetPath{}, path...)
			paths = append(paths, c)
			stack = stack[:len(stack)-1]
			if len(stack) < 1 {
				break
			}
			path = path[:len(path)-1]
			index = index[:len(index)-1]
			index[len(index)-1]++
			continue
		}
		index = append(index, 0)
	}
	return paths
}

func (lg *licenseGraphImp) WalkRestricted(fullWalk bool) ResolutionSet {
	rs := ResolveGraphConditions(lg)

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
			if !edgeIsDerivativeOf(e) {
				continue
			}
			dcs, ok := rs.(*resolutionSetImp).resolutions[e.dependency]
			if !ok {
				dcs = rmap[f].Copy().(*licenseConditionSetImp)
			} else {
				dcs = dcs.filter("restricted").union(rmap[f])
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
			if !edgeIsDerivativeOf(e) {
				continue
			}
			dcs, ok := rs.(*resolutionSetImp).resolutions[e.dependency]
			if !ok {
				continue
			}
			dcs = dcs.filter("restricted")
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

	for _, r := range rs.Targets() {
		cs, _ := rs.Conditions(r)
		if !cs.HasAnyByName("restricted") {
			continue
		}
		timp := r.(targetNodeImp)
		rcs := newLicenseConditionSet(&timp)
		rcs.addAll(cs.ByName("restricted"))
		if lg.targets[timp.file].isContainer {
			walkContainer(timp.file, rcs)
		} else {
			walkNonContainer(timp.file, rcs)
		}
	}

	return &resolutionSetImp{lg, rmap}
}

func (lg *licenseGraphImp) WalkDepsForCondition(condition string) ResolutionSet {
	rs := ResolveGraphConditions(lg)

	rmap := make(map[string]*licenseConditionSetImp)

	var walk func(string)

	walk = func(f string) {
		if _, ok := rmap[f]; ok {
			return
		}
		for _, e := range lg.index[f] {
			if !edgeIsDerivativeOf(e) {
				continue
			}
			dcs, ok := rs.(*resolutionSetImp).resolutions[e.dependency]
			if !ok {
				continue
			} else if !dcs.HasAnyByName(condition) {
				continue
			}
			rcs := newLicenseConditionSet(nil)
 			rcs.addAll(dcs.ByName(condition))
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
		if !cs.HasAnyByName(condition) {
			continue
		}
		rcs := newLicenseConditionSet(&targetNodeImp{lg, r})
		rcs.addAll(cs.ByName(condition))
		rmap[r] = rcs
		if lg.targets[r].isContainer {
			walk(r)
		}
	}
	return &resolutionSetImp{lg, rmap}
}

// newLicenseGraphImp constructs a new instance of licenseGraphImp.
func newLicenseGraphImp() *licenseGraphImp {
	return &licenseGraphImp{
		rootFiles: []string{},
		targets: make(map[string]*targetNode),
	}
}
