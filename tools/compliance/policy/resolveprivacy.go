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

// ResolveSourcePrivacy implements the policy for source privacy.
//
// Based on the "proprietary" walk, selectively prunes inherited conditions
// from pure aggregates.
func ResolveSourcePrivacy(licenseGraph LicenseGraph) ResolutionSet {
	// For pure aggregates, a container may contain both binaries with private source
	// and binaries with shared source so the pure aggregate does not need to inherit
	// the "proprietary" condition. Policy requires non-aggregates and
	// not-pure-aggregates with source sharing requirements to share all of the source
	// affecting any private source contained therein.
	proprietary := ResolveTopDownForCondition(licenseGraph, "proprietary")

	lg := licenseGraph.(*licenseGraphImp)
	if proprietary.(*resolutionSetImp).lg == nil {
		panic(fmt.Errorf("nil graph for proprietary"))
	} else if lg == nil {
		panic(fmt.Errorf("nil graph for privacy"))
	} else if lg != proprietary.(*resolutionSetImp).lg {
		panic(fmt.Errorf("resolve proprietary for wrong graph"))
	}

	// privateSource is the set of all source privacy resolutions pruned to account for
	// pure aggregations.
	privateSource := make(map[string]*licenseConditionSetImp)

	rsimp := proprietary.(*resolutionSetImp)
	for appliesTo, cs := range rsimp.resolutions {
		found := false
		for _, origins := range cs.conditions {
			// selfConditionsApplicableForConditionName(...) above causes all
			// conditions in a ResolutionSet to record whether they apply to
			// pure aggregates.
			for origin, treatAsAggregate := range origins {
				if treatAsAggregate == nil || treatAsAggregate == false || origin == appliesTo {
					found = true
				}
			}
		}
		if !found {
			// all conditions for `appliesTo` pruned so skip rest
			continue
		}
		for name, origins := range cs.conditions {
			for origin, treatAsAggregate := range origins {
				if treatAsAggregate == nil || treatAsAggregate == false || origin == appliesTo {
					if _, ok := privateSource[appliesTo]; !ok {
						privateSource[appliesTo] = newLicenseConditionSet(&targetNodeImp{lg, appliesTo})
					}
					if _, ok := privateSource[appliesTo].conditions[name]; !ok {
						privateSource[appliesTo].conditions[name] = make(map[string]interface{})
					}
					privateSource[appliesTo].conditions[name][origin] = rsimp.resolutions[appliesTo].conditions[name][origin]
				}
			}
		}
	}

	return &resolutionSetImp{lg, privateSource}
}
