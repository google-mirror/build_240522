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

const (
	defaultEdgesSize = 1000
)

// licenseGraphImp implements the publicly immutable LicenseGraph interface.
//
// Internally, the implementation uses lazy evaluation and caching in critical sections.
type licenseGraphImp struct {
	// rootFiles identifies the original set of files to read. (immutable)
	//
	// Defines the starting "top" for top-down walks.
	//
	// Alternatively, an instance of licenseGraphImp conceptually defines a scope within
	// the universe of build graphs as a sub-graph rooted at rootFiles where all edges
	// and targets for the instance are defined relative to and within that scope. For
	// most analyses, the correct scope is to root the graph at all of the distributed
	// artifacts.
	rootFiles []string

	// edges lists the directed edges in the graph from target to dependency. (guarded by mu)
	//
	// Alternatively, the graph is the set of `edges`.
	edges []*dependencyEdge

	// targets identifies, indexes by name, and describes the entire set of target node files. (guarded by mu)
	targets map[string]*targetNode

	// index facilitates looking up edges from targets. (creation guarded by my)
	//
	// This is a forward index from target to dependencies. i.e. "top-down"
	index map[string][]*dependencyEdge

	// rs caches the results of a full bottom-up resolve. (creation guarded by mu)
	//
	// A bottom-up resolve is a prerequisite for all of the top-down resolves so caching
	// the result is a performance win.
	rs *resolutionSetImp

	// mu guards against concurrent update.
	mu sync.Mutex
}

// newLicenseGraphImp constructs a new, empty instance of licenseGraphImp.
func newLicenseGraphImp() *licenseGraphImp {
	return &licenseGraphImp{
		rootFiles: []string{},
		edges:     make([]*dependencyEdge, 0, 1000),
		targets:   make(map[string]*targetNode),
	}
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

// implementations of the public methods

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

// Targets returns the list of target nodes in the graph. (unordered)
func (lg *licenseGraphImp) Targets() []TargetNode {
	targets := make([]TargetNode, 0, len(lg.targets))
	for target := range lg.targets {
		targets = append(targets, targetNodeImp{lg, target})
	}
	return targets
}

// targetEdgeImp implements TargetEdge
type targetEdgeImp struct {
	// lg identifies the scope, i.e. license graph, in which the edge appears.
	lg *licenseGraphImp

	// e identifies describes the target, dependency, and annotations of the edge.
	e *dependencyEdge
}

// Target returns the depending end of the edge.
//
// Target needs Dependency to build.
func (e targetEdgeImp) Target() TargetNode {
	return targetNodeImp{e.lg, e.e.target}
}

// Dependency returns the depended-on end of the edge.
//
// Dependency builds without Target, but Target needs Dependency to build.
func (e targetEdgeImp) Dependency() TargetNode {
	return targetNodeImp{e.lg, e.e.dependency}
}

// Annotations describe the type of edge.
//
// Only the annotations prescribed by policy have any meaning for licensing, and
// the meaning for licensing is likewise prescribed by policy. Other annotations
// are preserved and ignored by policy.
func (e targetEdgeImp) Annotations() TargetEdgeAnnotations {
	return e.e.annotations
}

// compliance-only targetEdgeImp methods

// isDynamicLink returns true if the edge represents a shared or dynamic link at runtime.
func (e targetEdgeImp) isDynamicLink() bool {
	return edgeIsDynamicLink(e.e) // defined by policy.go
}

// isDerivation returns true if the edge represents a dependency incorporated into the target as a derivative work.
func (e targetEdgeImp) isDerivation() bool {
	return edgeIsDerivation(e.e) // defined by policy.go
}

// areIndependentModules returns true of the ends of the edges are in different packages.
func (e targetEdgeImp) areIndependentModules() bool {
	return edgeNodesAreIndependentModules(e) // defined by policy.go
}

// targetNodeImp implements TargetNode
//
// The target can appear as the target in edges or as the dependency in edges.
// Most targets appear as both target in one edge and dependency in other edges.
type targetNodeImp struct {
	// lg identifies the scope, i.e. license graph, in which the target appears.
	lg *licenseGraphImp

	// file identifies the target node by the path to the license metadata file.
	file string
}

// Name returns the name identifying the target node. i.e. the path to the corresponding license metadata file
func (tn targetNodeImp) Name() string {
	return tn.lg.targets[tn.file].name
}

// PackageName returns the string identifying what package the target is part of.
func (tn targetNodeImp) PackageName() string {
	return tn.lg.targets[tn.file].name
}

// ModuleTypes returns the list of module types that define the target. (unordered)
func (tn targetNodeImp) ModuleTypes() []string {
	return append([]string{}, tn.lg.targets[tn.file].ModuleTypes...)
}

// ModuleClasses returns the list of module classes associated with the target. (unordered)
func (tn targetNodeImp) ModuleClasses() []string {
	return append([]string{}, tn.lg.targets[tn.file].ModuleClasses...)
}

// Projects returns the list of projects that define the target. (unordered)
func (tn targetNodeImp) Projects() []string {
	return append([]string{}, tn.lg.targets[tn.file].Projects...)
}

// LicenseKinds returns the list of license kind names associated with the target. (unordered)
// e.g. SPDX-license-identifier-Apache-2.0 or legacy_notice
func (tn targetNodeImp) LicenseKinds() []string {
	return append([]string{}, tn.lg.targets[tn.file].LicenseKinds...)
}

// LicenseConditions returns the set of license condition names originating at the target. (unordered)
// e.g. restricted or by_exception_only
func (tn targetNodeImp) LicenseConditions() LicenseConditionSet {
	return newLicenseConditionSet(&tn, tn.lg.targets[tn.file].LicenseConditions...)
}

// LicenseTexts returns the list of paths to license text for the target. (unordered)
func (tn targetNodeImp) LicenseTexts() []string {
	return append([]string{}, tn.lg.targets[tn.file].LicenseTexts...)
}

// IsContainer returns true if the target represents a container that merely aggregates other targets.
func (tn targetNodeImp) IsContainer() bool {
	return tn.lg.targets[tn.file].GetIsContainer()
}

// Built returns the list of files built by the module or target. (unordered)
func (tn targetNodeImp) Built() []string {
	return append([]string{}, tn.lg.targets[tn.file].Built...)
}

// Installed returns the list of files installed by the module or target. (unordered)
func (tn targetNodeImp) Installed() []string {
	return append([]string{}, tn.lg.targets[tn.file].Installed...)
}

// InstallMap returns the list of path name transformations to make to move files from their original
// location in the file system to their destination inside a container. (unordered)
func (tn targetNodeImp) InstallMap() []InstallMap {
	result := make([]InstallMap, 0, len(tn.lg.targets[tn.file].InstallMap))
	for _, im := range tn.lg.targets[tn.file].InstallMap {
		result = append(result, InstallMap{im.GetFromPath(), im.GetContainerPath()})
	}
	return result
}

// Sources returns the list of files actually depended on by the target, which may be a proper subset
// of the files available from dependencies. (unordered)
func (tn targetNodeImp) Sources() []string {
	return append([]string{}, tn.lg.targets[tn.file].Sources...)
}

// targetEdgeAnnotationsImp implements TargetEdgeAnnotations
type targetEdgeAnnotationsImp map[string]interface{}

// HasAnnotation returns true if `ann` is attached to the edge.
func (ea targetEdgeAnnotationsImp) HasAnnotation(ann string) bool {
	_, ok := ea[ann]
	return ok
}

// AsList returns the list of annotation names attached to the edge. (unordered)
func (ea targetEdgeAnnotationsImp) AsList() []string {
	l := make([]string, 0, len(ea))
	for ann := range ea {
		l = append(l, ann)
	}
	return l
}
