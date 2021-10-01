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
	// WalkRestricted returns the set of metadata files to treat as 'restricted'.
	WalkRestricted(fullWalk bool) ResolutionSet
	// WalkDepsForCondition returns the set of metadata with license `condition`.
	WalkDepsForCondition(condition string) ResolutionSet
	AnyPath(target, dependency TargetNode) TargetPath
	AllPaths(target, dependency TargetNode) []TargetPath
}

// LicenseCondition describes an individual license condition or requirement.
type LicenseCondition interface {
	Name() string
	AppliesTo() TargetNode
}

// LicenseConditionSet describes the set of conditions applicable to a target or edge
// either before or after resolution.
type LicenseConditionSet interface {
	Add(condition... LicenseCondition)
	ByName(name string) []LicenseCondition
	Conditions() []LicenseCondition
	Copy() LicenseConditionSet
	HasAnyByName(name string) bool
	IsEmpty() bool
}

// TargetNode describes a target corresponding to a specific license metadata file.
type TargetNode interface {
	Name() string
	Projects() []string
}

// TargetResolution pairs a target with its resolved license conditions.
type TargetResolution interface {
	Target() TargetNode
	Conditions() LicenseConditionSet
}

// ResolutionSet describes a set of TargetResolution
type ResolutionSet interface {
	Targets() []TargetNode
	Conditions(target TargetNode) (LicenseConditionSet, bool)
	HasTarget(t TargetNode) bool
}

// TargetSet describes a set of TargetNode.
type TargetSet interface {
	Add(node... TargetNode)
	Targets() []TargetNode
	Projects() []string
}

// TargetEdge describes an edge from a target to a dependency.
type TargetEdge interface {
	Target() TargetNode
	Dependency() TargetNode
	Annotations() TargetEdgeAnnotations
}

// TargetPath describes a sequence of edges where the dependency of one is the target of the next
// describing a path from the target of the 1st element to the dependency of the last element.
type TargetPath []TargetEdge

// TargetEdgeAnnotations describes a set of annotations attached to an edge from a target to a dependency.
type TargetEdgeAnnotations interface {
	HasAnnotation(ann string) bool
	ListAnnotations() []string
}
