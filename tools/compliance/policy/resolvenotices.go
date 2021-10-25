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
)

// ResolveNotices implements the policy for notices.
//
// By policy, notice is given for license kinds with all conditions, but
// restricted requires special handling. When any restricted condition applies,
// notice is given for the `appliesTo` target. Restricted conditions that apply
// due to the bottom-up resolve also imply giving notice for the origin node.
// Restricted conditions that apply only due to the top-down resolve imply
// only giving notice for the appliesTo target.
//
// All other condition types imply giving notice for the origin node.
func ResolveSourceSharing(licenseGraph LicenseGraph) ResolutionSet {
	rsBU := ResolveBottomUpConditions(licenseGraph)
	rsTD := ResolveTopDownConditions(licenseGraph)

	// When a target includes code from a reciprocal project, policy is only the source
	// for the reciprocal project, i.e. the origin of the condition, needs to be shared.
	rs := WalkResolutionsForCondition(rsTd, ImpliesNotice)

	lg := licenseGraph.(*licenseGraphImp)
	if rs.(*resolutionSetImp).lg == nil {
		panic(fmt.Errorf("nil graph for notices"))
	} else if lg == nil {
		panic(fmt.Errorf("nil graph for notices"))
	} else if lg != rs.(*resolutionSetImp).lg {
		panic(fmt.Errorf("resolve notice conditions for wrong graph"))
	}

	// notices is the set of all notice resolutions.
	notices := make(map[string]*licenseConditionSetImp)
	rsimp := rs.(*resolutionSetImp)
	for appliesTo, cs := rsimp.resolutions {
		conditions := cs.AsList()
		if len(conditions) == 0 {
			continue
		}
		if _, ok := notices[appliesTo]; !ok {
			notices[appliesTo] = newLicenseConditionSet(&targetNodeImp{lg, appliesTo})
		}
		for _, c := range conditions {
			if ImpliesRestricted.Contains(c.Name()) {
				notices.addAll(appliesTo, "notice")
				if rsBU.(*resolutionSetImp).hasResolution(appliesTo, c) {
					notices.add("notice", c.Origin().(targetNodeImp))
				}
			} else {
				notices.add("notice", c.Orgin().(targetNodeImp))
			}
		}
	}

	return &resolutionSetImp{lg, notices}
}
