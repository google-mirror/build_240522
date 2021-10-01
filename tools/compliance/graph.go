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

// LicenseGraph describes the license metadata for a set of root targets.
type LicenseGraph interface {
	// TargetNode returns the named target.
	TargetNode(name string) TargetNode

	// HasTargetNode returns true if the named target exists in the graph.
	HasTargetNode(name string) bool

	// WalkRestricted returns the set of 'restricted' license condition resolutions by target.
	// When `fullWalk` is true, will revisit nodes to make sure all resolutions are reported.
	WalkRestricted(fullWalk bool) ResolutionSet

	// WalkDepsForCondition returns the set of license `condition` resolutions by target.
	WalkDepsForCondition(condition string) ResolutionSet

	// AnyPath returns a list of dependencies starting at `target` and ending at `dependency` or nil.
	AnyPath(target, dependency TargetNode) TargetPath

	// AllPaths returns a list of the paths starting at `target` and ending at `dependency`.
	AllPaths(target, dependency TargetNode) []TargetPath
}

// LicenseCondition describes an individual license condition or requirement.
type LicenseCondition interface {
	// Name identifies the type of condition. e.g. "notice" or "restricted".
	Name() string
	// AppliesTo identifies the target node that originated the condition.
	AppliesTo() TargetNode
}

// LicenseConditionSet describes the set of conditions applicable to a target or edge
// either before or after resolution.
type LicenseConditionSet interface {
	// Add makes `condition` a member of the set if it was not previously.
	Add(condition... LicenseCondition)

	// ByName returns all of the conditions in the set with `name`.
	ByName(name string) []LicenseCondition
	// HasAnyByName returns true if the set contains any contions with `name`.
	HasAnyByName(name string) bool

	// Conditions returns all of the license conditions in the set.
	Conditions() []LicenseCondition

	// Copy makes a duplicate of the set that can be modified independenly.
	Copy() LicenseConditionSet

	// IsEmpty returns true if the set contains zero elements.
	IsEmpty() bool
}

// TargetNode describes a target corresponding to a specific license metadata file.
type TargetNode interface {
	// Name returns the string that identifies the target note. i.e. path to license metadata for node
	Name() string

	// Projects returns the projects defining the target node.
	Projects() []string
}

// TargetResolution pairs a target with its resolved license conditions.
type TargetResolution interface {
	// Target identifies the target node requiring action to resolve the conditions.
	Target() TargetNode

	// Conditions identifies the license conditions or requirements that must be met.
	Conditions() LicenseConditionSet
}

// ResolutionSet describes a set of TargetResolution
type ResolutionSet interface {
	// Targets identifies the list of targets requiring action to resolve conditions.
	Targets() []TargetNode

	// Conditions identifies the license conditions or requirements that `target` must meet.
	Conditions(target TargetNode) (LicenseConditionSet, bool)

	// HasTarget returns true if the set contains conditions for `target`.
	HasTarget(target TargetNode) bool
}

// TargetSet describes a distinct set of target nodes.
type TargetSet interface {
	// Add makes each `node` a member of the set if it was not previously.
	Add(node... TargetNode)

	// Targets returns an iterable list of targets in the set.
	Targets() []TargetNode

	// Projects returns the distinct set of projects containing the targets in the set.
	Projects() []string
}

// TargetEdge describes a directed edge from a target to a dependency.
type TargetEdge interface {
	// Target identifies the target that depends on the dependency.
	Target() TargetNode

	// Dependency identifies the target depended on by the target.
	Dependency() TargetNode

	// Annotations describes the set of annotations attached to the edge.
	Annotations() TargetEdgeAnnotations
}

// TargetPath describes a sequence of edges where the dependency of one is the target of the next
// describing a path from the target of the 1st element to the dependency of the last element.
// If a -> b -> c, target path would have 2 annotated edges representing the dependency arrows.
type TargetPath []TargetEdge


// TargetEdgeAnnotations describes a set of annotations attached to an edge from a target to a dependency.
type TargetEdgeAnnotations interface {
	// HasAnnotation returns true if an annotation `ann` is in the set.
	HasAnnotation(ann string) bool

	// ListAnnotations returns the list of string annotations in the set.
	ListAnnotations() []string
}
