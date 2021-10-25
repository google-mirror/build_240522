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

// LicenseGraph describes the immutable license metadata for a set of root
// targets and the transitive closure of their dependencies.
//
// Alternatively, a graph is a set of edges. In this case directed, annotated
// edges from targets to dependencies.
//
// A LicenseGraph provides the frame of reference for all of the other types
// defined here. It is possible to have multiple graphs, and to have targets,
// edges, and resolutions from multiple graphs. But it is an error to try to
// mix items from different graphs in the same operation.
// May panic if attempted.
//
// The compliance package assumes specific private implementations of each of
// these interfaces. May panic if attempts are made to combine different
// implementations of some interfaces with expected implementations of other
// interfaces here.
type LicenseGraph interface {
	// Edges returns a list of the edges comprising the graph.
	Edges() []TargetEdge

	// TargetNode returns the named target. Panics if named node does not
	// appear in graph.
	TargetNode(name string) TargetNode

	// HasTargetNode returns true if the named target exists in the graph.
	HasTargetNode(name string) bool
}

// TargetEdge describes a directed, annotated edge from a target to a
// dependency. (immutable)
//
// A LicenseGraph, above, is a set of TargetEdges.
//
// i.e. `Target` depends on `Dependency` in the manner described by
// `Annotations`.
type TargetEdge interface {
	// Target identifies the target that depends on the dependency.
	Target() TargetNode

	// Dependency identifies the target depended on by the target.
	Dependency() TargetNode

	// Annotations describes the set of annotations attached to the edge.
	Annotations() TargetEdgeAnnotations
}


// TargetNode describes a module or target identified by the name of a specific
// metadata file. (immutable)
//
// Each metadata files corresponds to a Soong module or to a Make target.
type TargetNode interface {
	// Name returns the string that identifies the target node.
	// i.e. path to license metadata file
	Name() string

	// PackageName returns the string that identifes the package for the
	// target.
	PackageName() string

	// ModuleTypes returns the list of module types implementing the
	// target. (unordered)
	//
	// In an ideal world, only 1 module type would implement each target,
	// but the interactions between Soong and Make for host versus product
	// and for a variety of architectures sometimes causes multiple module
	// types per target. Often a regular build target and a prebuilt.
	ModuleTypes() []string

	// ModuleClasses returns the list of module classes implementing the
	// target. (unordered)
	ModuleClasses() []string

	// Projects returns the projects defining the target node. (unordered)
	//
	// In an ideal world, only 1 project defines a target, but the
	// interaction between Soong and Make for a variety of architectures
	// and for host versus product means a module is sometimes defined more
	// than once.
	Projects() []string

	// LicenseKinds returns the list of license kind names for the module
	// or target.
	LicenseKinds() []string

	// LicenseConditions returns a copy of the set of license conditions
	// originating at the target. The values that appear and how each is
	// resolved is a matter of policy.
	LicenseConditions() LicenseConditionSet

	// LicenseTexts returns the paths to the files containing the license
	// texts for the target.
	LicenseTexts() []string

	// IsContainer returns true if the target represents a container that
	// merely aggregates other targets.
	IsContainer() bool

	// Built returns the list of file names built by the module or target.
	Built() []string

	// Installed returns the list of file names installed by the module or
	// target.
	Installed() []string

	// InstallMap returns the list of transformations to make to path names
	// when moving into their installed locations.
	InstallMap() []InstallMap

	// Sources returns the list of file names depended on by the target,
	// which may be a proper subset of those made available by dependency
	// modules.
	Sources() []string
}

// InstallMap describes the mapping from an input filesystem file to file in a
// container.
type InstallMap struct {
	// FromPath is the input path on the filesystem.
	FromPath string

	// ContainerPath is the path to the same file inside the container or
	// installed location.
	ContainerPath string
}


// TargetEdgeAnnotations describes a set of annotations attached to an edge
// from a target to a dependency.
type TargetEdgeAnnotations interface {
	// HasAnnotation returns true if an annotation `ann` is in the set.
	HasAnnotation(ann string) bool

	// ListAnnotations returns the list of string annotations in the set.
	ListAnnotations() []string
}
