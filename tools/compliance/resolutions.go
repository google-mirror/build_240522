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

// Resolution describes an action to resolve one or more license conditions.
//
// `AttachesTo` identifies the target node that when distributed triggers the action.
// `ActsOn` identifies the target node that is the object of the action.
// `Resolves` identifies one or more license conditions that the action resolves.
//
// e.g. Suppose an MIT library is linked to a binary that also links to GPL code.
//
// A resolution would attach to the binary to share (act on) the MIT library to
// resolve the restricted condition originating from the GPL code.
type Resolution interface {
	AttachesTo() TargetNode
	ActsOn() TargetNode
	Resolves() LicenseConditionSet
}

// ResolutionSet describes an immutable set of targets and the license
// conditions each target must satisfy or "resolve" in a specific context.
//
// Ultimately, the purpose of recording the license metadata and building a
// license graph is to identify, describe, and verify the necessary actions or
// operations for compliance policy.
//
// i.e. What is the source-sharing policy? Has it been met? Meet it.
//
// i.e. Are there incompatible policy requirements? Such as a source-sharing
// policy applied to code that policy also says may not be shared? If so, stop
// and remove the dependencies that create the situation.
//
// The ResolutionSet is the base unit for mapping license conditions to the
// targets triggering some necessary action per policy. Different ResolutionSet
// values may be calculated for different contexts.
//
// e.g. Suppose an unencumbered binary links in a notice .a library.
//
// An "unencumbered" condition would originate from the binary, and a "notice"
// condition would originate from the .a library. A ResolutionSet for the
// context of the Notice policy might apply both conditions to the binary while
// preserving the origin of each condition. By applying the notice condition to
// the binary, the ResolutionSet stipulates the policy that the release of the
// unencumbered binary must provide suitable notice for the .a library.
//
// The resulting ResolutionSet could be used for building a notice file, for
// validating that a suitable notice has been built into the distribution, or
// for reporting what notices need to be given.
//
// Resolutions for different contexts may be combined in a new ResolutionSet
// using JoinResolutions(...).
//
// See: resolve.go for:
//  * ResolveBottomUpConditions(...)
//  * ResolveTopDownForCondition(...)
// See also: policy.go for:
//  * ResolveSourceSharing(...)
//  * ResolveSourcePrivacy(...)
type ResolutionSet interface {
	// AttachesTo identifies the list of targets triggering action to
	// resolve conditions. (unordered)
	AttachesTo() []TargetNode

	// Resolutions returns the list of resolutions that `attachedTo`
	// target must resolve. Returns empty list if no conditions apply.
	Resolutions(attachedTo TargetNode) []Resolution

	// ResolutionsByActsOn returns the list of resolutions that must
	// `actOn` to resolvee. Returns empty list if no conditions apply.
	ResolutionsByActsOn(actOn TargetNode) []Resolution

	// Origins identifies the list of targets originating conditions to
	// resolve. (unordered)
	Origins() []TargetNode

	// ActsOn identifies the list of targets to be acted on to resolve
	// conditions. (unordered)
	ActsOn() []TargetNode

	// AttachesToByOrigin identifies the list of targets requiring action to
	// resolve conditions originating at `origin`. (unordered)
	AttachesToByOrigin(origin TargetNode) []TargetNode

	// AttachesToTarget returns true if the set contains conditions that
	// apply to `attachedTo`.
	AttachesToTarget(attachedTo TargetNode) bool

	// AnyByNameAttachToTarget returns true if the set contains conditions
	// matching `names` that attach to `attachesTo`.
	AnyByNameApplyToTarget(attachesTo TargetNode, names ...ConditionNames) bool

	// AllByNameAttachTo returns true if the set contains at least one
	// condition matching each element of `names` for `attachesTo`.
	AllByNameApplyToTarget(attachesTo TargetNode, names ...ConditionNames) bool
}
