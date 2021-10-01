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

// JoinResolutions returns a new ResolutionSet combining the resolutions from
// multiple resolution sets derived from the same license graph.
//
// e.g. combine "restricted", "reciprocal", and "proprietary" resolutions.
func JoinResolutions(resolutions... ResolutionSet) ResolutionSet {
        rmap := make(map[string]*licenseConditionSetImp)
        var lg *licenseGraphImp
        for _, r := range resolutions {
                if lg == nil {
                        lg = r.(*resolutionSetImp).lg
                } else if lg != r.(*resolutionSetImp).lg {
                        panic(fmt.Errorf("attempt to join resolutions from multiple graphs"))
                }
                for t, cs := range r.(*resolutionSetImp).resolutions {
                        if _, ok := rmap[t]; !ok {
                                rmap[t] = cs.Copy().(*licenseConditionSetImp)
                                continue
                        }
                        rmap[t].addAll(cs.Conditions())
                }
        }
        return &resolutionSetImp{lg, rmap}
}


// resolutionSetImp implements ResolutionSet.
type resolutionSetImp struct {
	lg *licenseGraphImp
	resolutions map[string]*licenseConditionSetImp
}

// AppliesTo returns the list of targets with applicable license conditions attached.
func (rs *resolutionSetImp) AppliesTo() []TargetNode {
	targets := make([]TargetNode, 0, len(rs.resolutions))
	for t := range rs.resolutions {
		targets = append(targets, targetNodeImp{rs.lg, t})
	}
	return targets
}

// Conditions returns the set of conditions applied to `appliesTo`, and whether `appliesTo` is in set.
func (rs *resolutionSetImp) Conditions(appliesTo TargetNode) (LicenseConditionSet, bool) {
	timp := appliesTo.(targetNodeImp)
	if rs.lg == nil {
		rs.lg = timp.lg
	} else if timp.lg == nil {
		timp.lg = rs.lg
	} else if rs.lg != timp.lg {
		panic(fmt.Errorf("attempt to query target resolutions for wrong graph"))
	}
	cs, ok := rs.resolutions[timp.file]
	return cs.Copy(), ok
}

// AppliesToTarget returns true if `appliesTo` appears in the set.
func (rs *resolutionSetImp) AppliesToTarget(appliesTo TargetNode) bool {
	timp := appliesTo.(targetNodeImp)
	if rs.lg == nil {
		rs.lg = timp.lg
	} else if timp.lg == nil {
		timp.lg = rs.lg
	} else if rs.lg != timp.lg {
		panic(fmt.Errorf("attempt to query resolved targets for wrong graph"))
	}
	return rs.hasTarget(timp.file)
}

// AnyByNameApplyToTarget returns true if `appliesTo` appears in the set with any conditions matching `name`.
func (rs resolutionSetImp) AnyByNameApplyToTarget(appliesTo TargetNode, names... ConditionNames) bool {
	timp := appliesTo.(targetNodeImp)
	if rs.lg == nil {
		rs.lg = timp.lg
	} else if timp.lg == nil {
		timp.lg = rs.lg
	} else if rs.lg != timp.lg {
		panic(fmt.Errorf("attempt to query target resolutions for wrong graph"))
	}
	return rs.hasAnyByName(timp.file, names...)
}

// HasAllByName returns true if `target` appears in the set with conditions matching every element of `name`.
func (rs resolutionSetImp) AllByNameApplyToTarget(appliesTo TargetNode, names... ConditionNames) bool {
	timp := appliesTo.(targetNodeImp)
	if rs.lg == nil {
		rs.lg = timp.lg
	} else if timp.lg == nil {
		timp.lg = rs.lg
	} else if rs.lg != timp.lg {
		panic(fmt.Errorf("attempt to query target resolutions for wrong graph"))
	}
	return rs.hasAllByName(timp.file, names...)
}

// compliance-only resolutionSetImp methods

// newResolutionSet constructs a new, empty instance of resolutionSetImp
func newResolutionSet(lg *licenseGraphImp) *resolutionSetImp {
	return &resolutionSetImp{lg, make(map[string]*licenseConditionSetImp)}
}

// add applies all of the license conditions in `cs` to `file` if not already applied.
func (rs *resolutionSetImp) add(file string, cs *licenseConditionSetImp) {
	if r, ok := rs.resolutions[file]; ok {
		rs.resolutions[file] = r.union(cs)
	} else {
		rs.resolutions[file] = cs
	}
}

// hasTarget returns true if `file` appears as one of the targets in the set.
func (rs *resolutionSetImp) hasTarget(file string) bool {
	_, isPresent := rs.resolutions[file]
	return isPresent
}

// hasAnyByName returns true if the target for `file` has at least 1 condition matching `name`.
func (rs *resolutionSetImp) hasAnyByName(file string, name... ConditionNames) bool {
	cs, isPresent := rs.resolutions[file]
	if !isPresent {
		return false
	}
	for _, cn := range name {
		for _, n := range cn {
			_, isPresent = cs.conditions[n]
			if isPresent {
				return true
			}
		}
	}
	return false
}

// hasAllByName returns true if the target for `file` has at least 1 condition for each element of `name`.
func (rs *resolutionSetImp) hasAllByName(file string, name... ConditionNames) bool {
	cs, isPresent := rs.resolutions[file]
	if !isPresent {
		return false
	}
	for _, cn := range name {
		found := false
		for _, n := range cn {
			_, isPresent = cs.conditions[n]
			if isPresent {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
