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

// LicenseGraph describes the license metadata for a set of root targets and the transitive
// closure of their dependencies.
//
// A LicenseGraph provides the frame of reference for all of the other types defined here.
// It is possible to have multiple graphs, and to have targets, edges, and resolutions from
// multiple graphs. But it is an error to try to mix items from different graphs in the same
// operation. May panic if attempted.
//
// Compliance assumes specific private implementations of each of these interfaces. May panic
// if attempts are made to combine different implementations of some interfaces with expected
// implementations of other interfaces here.
type LicenseGraph interface {
	// TargetNode returns the named target.
	TargetNode(name string) TargetNode

	Edges() []TargetEdge

	// HasTargetNode returns true if the named target exists in the graph.
	HasTargetNode(name string) bool
}

// ResolutionSet describes a set of targets and the license conditions each target must satisfy.
//
// Ultimately, the purpose of recording the license metadata and building a license graph is to
// identify, describe, and verify the necessary actions or operations for compliance.
//
// e.g. What are the source-sharing requirements? Have they been met? Meet them.
//
// e.g. Are there incompatible requirements? Such as a source-sharing requirement on code that must
// not be shared? If so, stop and remove the dependencies that create the situation.
//
// The ResolutionSet is the base unit for mapping license conditions to the targets where action
// may or must be taken in a specific context
//
// e.g. See: ResolveBottomUpConditions(...), ResolveTopDownCondition(...) and ResolveTopDownRestricted
type ResolutionSet interface {
	// Targets identifies the list of targets requiring action to resolve conditions (unordered).
	Targets() []TargetNode

	// Conditions identifies the license conditions or requirements that `target` must meet, and
	// whether the target appears in the set.
	Conditions(target TargetNode) (LicenseConditionSet, bool)

	// HasTarget returns true if the set contains conditions for `target`.
	HasTarget(target TargetNode) bool

	// HasAnyByName returns true if the set contains conditions with `name` for `target`.
	HasAnyByName(target TargetNode, name... ConditionNames) bool

	// HasAllByName returns true if the set contains all conditions with `name` for `target`.
	HasAllByName(target TargetNode, name... ConditionNames) bool
}


// TargetNode describes a target corresponding to a specific license metadata file.
//
// These correspond to Soong modules or to Make targets.
type TargetNode interface {
	// Name returns the string that identifies the target node. i.e. path to license metadata for node
	Name() string

	// IsContainer returns true if the target represents a container that merely aggregates other targets.
	IsContainer() bool

	// Projects returns the projects defining the target node.
	// In an ideal world, only 1 project defines a target, but the interaction between Soong and Make
	// for a variety of architectures and for host versus product means a module is sometimes defined
	// more than once.
	Projects() []string

	// LicenseTexts returns the paths to the files containing the license texts for the target.
	LicenseTexts() []string
}

// TargetEdge describes a directed, annotated edge from a target to a dependency.
//
// i.e. Target depends on Dependency in the manner described by Annotations.
type TargetEdge interface {
	// Target identifies the target that depends on the dependency.
	Target() TargetNode

	// Dependency identifies the target depended on by the target.
	Dependency() TargetNode

	// Annotations describes the set of annotations attached to the edge.
	Annotations() TargetEdgeAnnotations
}


// LicenseCondition describes an individual license condition or requirement.
//
// Origin identifies the origin of the requirement.
//
// e.g. Suppose an unencumbered binary links in a notice .a library. An "unencumbered" condition
// would originate from the binary, and a "notice" condition would originate from the .a library.
// A ResolutionSet might then apply both conditions to the binary by association preserving the
// origin of each condition.
type LicenseCondition interface {
	// Name identifies the type of condition. e.g. "notice" or "restricted".
	Name() string
	// Origin identifies the target node that originated the condition.
	Origin() TargetNode
}

// LicenseConditionSet describes a set of conditions.
//
// Use func NewLicenseConditionSet(conditions... LicenseCondition) LicenseConditionSet to construct
// LicenseConditionSet variables for local use.
type LicenseConditionSet interface {
	// Add makes each `condition` a member of the set if it was not previously.
	Add(condition... LicenseCondition)

	// ByName returns all of the conditions in the set with `name`.
	ByName(name... ConditionNames) []LicenseCondition
	// HasAnyByName returns true if the set contains any conditions with `name`.
	HasAnyByName(name... ConditionNames) bool

	// ByOrigin returns all of the conditions that originate at `target`.
	ByOrigin(target TargetNode) []LicenseCondition
	// HasAnyByOrigin returns true if the set contains any conditions originating at `target`.
	HasAnyByOrigin(target TargetNode) bool

	// Conditions returns all of the license conditions in the set (unordered).
	Conditions() []LicenseCondition

	// Copy makes a duplicate of the set that can be modified independently.
	Copy() LicenseConditionSet

	// IsEmpty returns true if the set contains zero elements.
	IsEmpty() bool
}

// TargetEdgeAnnotations describes a set of annotations attached to an edge from a target to a dependency.
type TargetEdgeAnnotations interface {
	// HasAnnotation returns true if an annotation `ann` is in the set.
	HasAnnotation(ann string) bool

	// ListAnnotations returns the list of string annotations in the set.
	ListAnnotations() []string
}
