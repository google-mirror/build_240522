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
	"strings"
)

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
type Resolution struct {
	attachesTo, actsOn string
	cs                 *LicenseConditionSet
}

func (r Resolution) AttachesTo() TargetNode {
	return TargetNode{r.cs.lg, r.attachesTo}
}

func (r Resolution) ActsOn() TargetNode {
	return TargetNode{r.cs.lg, r.actsOn}
}

func (r Resolution) Resolves() *LicenseConditionSet {
	return r.cs.Copy()
}

func (r Resolution) asString() string {
	var sb strings.Builder
	cl := r.cs.AsList()
	sort.Sort(cl)
	fmt.Fprintf(&sb, "%s -> %s -> %s", r.attachesTo, r.actsOn, cl.String())
	return sb.String()
}

func (r Resolution) byName(names ConditionNames) Resolution {
	return Resolution{r.attachesTo, r.actsOn, r.cs.ByName(names)}
}

func (r Resolution) byOrigin(origin TargetNode) Resolution {
	return Resolution{r.attachesTo, r.actsOn, r.cs.ByOrigin(origin)}
}

func (r Resolution) copy() Resolution {
	return Resolution{r.attachesTo, r.actsOn, r.cs.Copy()}
}


type ResolutionList []Resolution

// ResolutionList partially orders Resolutions by AttachesTo() and ActsOn() leaving `Resolves()` unordered.
func (l ResolutionList) Len() int      { return len(l) }
func (l ResolutionList) Swap(i, j int) { l[i], l[j] = l[j], l[i] }
func (l ResolutionList) Less(i, j int) bool {
	if l[i].AttachesTo().Name() == l[j].AttachesTo().Name() {
		return l[i].ActsOn().Name() < l[j].ActsOn().Name()
	}
	return l[i].AttachesTo().Name() < l[j].AttachesTo().Name()
}

func (rl ResolutionList) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "[")
	sep := ""
	for _, r := range rl {
		fmt.Fprintf(&sb, "%s%s", sep, r.asString())
		sep = ", "
	}
	fmt.Fprintf(&sb, "]")
	return sb.String()
}

func (rl ResolutionList) AllConditions() *LicenseConditionSet {
	result := newLicenseConditionSet(nil)
	for _, r := range rl {
		result.AddSet(r.cs)
	}
	return result
}

func (rl ResolutionList) ByName(names ConditionNames) ResolutionList {
	result := make(ResolutionList, 0, rl.CountByName(names))
	for _, r := range rl {
		if r.Resolves().HasAnyByName(names) {
			result = append(result, r.byName(names))
		}
	}
	return result
}

func (rl ResolutionList) CountByName(names ConditionNames) int {
	c := 0
	for _, r := range rl {
		if r.Resolves().HasAnyByName(names) {
			c++
		}
	}
	return c
}

func (rl ResolutionList) CountConditionsByName(names ConditionNames) int {
	c := 0
	for _, r := range rl {
		c += r.Resolves().CountByName(names)
	}
	return c
}

func (rl ResolutionList) ByAttachesTo(attachesTo TargetNode) ResolutionList {
	result := make(ResolutionList, 0, rl.CountByActsOn(attachesTo))
	for _, r := range rl {
		if r.attachesTo == attachesTo.file {
			result = append(result, r.copy())
		}
	}
	return result
}

func (rl ResolutionList) CountByAttachesTo(attachesTo TargetNode) int {
	c := 0
	for _, r := range rl {
		if r.attachesTo == attachesTo.file {
			c++
		}
	}
	return c
}

func (rl ResolutionList) ByActsOn(actsOn TargetNode) ResolutionList {
	result := make(ResolutionList, 0, rl.CountByActsOn(actsOn))
	for _, r := range rl {
		if r.actsOn == actsOn.file {
			result = append(result, r.copy())
		}
	}
	return result
}

func (rl ResolutionList) CountByActsOn(actsOn TargetNode) int {
	c := 0
	for _, r := range rl {
		if r.actsOn == actsOn.file {
			c++
		}
	}
	return c
}

func (rl ResolutionList) ByOrigin(origin TargetNode) ResolutionList {
	result := make(ResolutionList, 0, rl.CountByOrigin(origin))
	for _, r := range rl {
		if r.Resolves().HasAnyByOrigin(origin) {
			result = append(result, r.byOrigin(origin))
		}
	}
	return result
}

func (rl ResolutionList) CountByOrigin(origin TargetNode) int {
	c := 0
	for _, r := range rl {
		if r.Resolves().HasAnyByOrigin(origin) {
			c++
		}
	}
	return c
}

