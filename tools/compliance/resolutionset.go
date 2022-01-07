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
	"strings"
)

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
type ResolutionSet map[*TargetNode]ActionSet

func (rs ResolutionSet) AttachesTo() TargetNodeList {
	result := make(TargetNodeList, 0, len(rs))
	for attachesTo := range rs {
		result = append(result, attachesTo)
	}
	return result
}

func (rs ResolutionSet) AttachesToTarget(target *TargetNode) bool {
	_, isPresent := rs[target]
	return isPresent
}

func (rs ResolutionSet) CountResolutions(attachesTo *TargetNode) int {
	as, ok := rs[attachesTo]
	if !ok {
		return 0
	}
	return len(as)
}

func (rs ResolutionSet) Resolutions(attachesTo *TargetNode) ResolutionList {
	as, ok := rs[attachesTo]
	if !ok {
		return ResolutionList{}
	}
	result := make(ResolutionList, 0, len(as))
	for actsOn, cs := range as {
		result = append(result, Resolution{attachesTo, actsOn, cs})
	}
	return result
}

func (rs ResolutionSet) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "{")
	sep := ""
	for attachesTo, as := range rs {
		fmt.Fprintf(&sb, "%s%s -> %s", sep, attachesTo.Name(), as.String())
		sep = ", "
	}
	fmt.Fprintf(&sb, "}")
	return sb.String()
}

type ActionSet map[*TargetNode]LicenseConditionSet

func (as ActionSet) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "{")
	sep := ""
	for actsOn, cs := range as {
		fmt.Fprintf(&sb, "%s%s%s", sep, actsOn.Name(), cs.String())
		sep = ", "
	}
	fmt.Fprintf(&sb, "}")
	return sb.String()
}
