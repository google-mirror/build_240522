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

// ResolveSourceSharing implements the policy for source-sharing.
//
// Reciprocal and Restricted top-down walks are joined except for reciprocal
// only the conditions where the AppliesTo and Origin are the same are preserved.
//
// Per policy, the source for reciprocal triggering target does not need to be
// shared unless it is also the originating target.
func ResolveSourceSharing(licenseGraph LicenseGraph) ResolutionSet {
	rs := ResolveTopDownConditions(licenseGraph)

	// When a target includes code from a reciprocal project, policy is only the source
	// for the reciprocal project, i.e. the origin of the condition, needs to be shared.
	recip := WalkResolutionsForCondition(rs, ImpliesReciprocal)

	// When a target includes code from a restricted project, the policy
	// requires sharing the restricted project, the including project(s), and the
	// transitive closure of the including project(s) derivation dependencies.
	restrict := WalkResolutionsForCondition(rs, ImpliesRestricted)

	lg := licenseGraph.(*licenseGraphImp)
	if recip.(*resolutionSetImp).lg == nil {
		panic(fmt.Errorf("nil graph for reciprocal"))
	} else if restrict.(*resolutionSetImp).lg == nil {
		panic(fmt.Errorf("nil graph for restricted"))
	} else if lg == nil {
		panic(fmt.Errorf("nil graph for sharing"))
	} else if lg != recip.(*resolutionSetImp).lg {
		panic(fmt.Errorf("resolve reciprocal for wrong graph"))
	} else if lg != restrict.(*resolutionSetImp).lg {
		panic(fmt.Errorf("resolve restricted for wrong graph"))
	}

	// shareSource is the set of all source-sharing resolutions.
	shareSource := make(map[string]*licenseConditionSetImp)

	// For reciprocal, policy is to share the originating code only.
	rsimp := recip.(*resolutionSetImp)
	for appliesTo, cs := range rsimp.resolutions {
		conditions := cs.byOrigin(appliesTo)
		if len(conditions) == 0 {
			continue
		}
		if _, ok := shareSource[appliesTo]; !ok {
			shareSource[appliesTo] = newLicenseConditionSet(&targetNodeImp{lg, appliesTo})
		}
		shareSource[appliesTo].addList(conditions)
	}

	// For restricted, policy is to share the originating code and any code linked to it.
	rsimp = restrict.(*resolutionSetImp)
	for appliesTo, cs := range rsimp.resolutions {
		conditions := cs.AsList()
		if len(conditions) == 0 {
			continue
		}
		if _, ok := shareSource[appliesTo]; !ok {
			shareSource[appliesTo] = newLicenseConditionSet(&targetNodeImp{lg, appliesTo})
		}
		shareSource[appliesTo].addSet(rsimp.resolutions[appliesTo])
	}

	return &resolutionSetImp{lg, shareSource}
}
