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
// multiple resolution sets. All sets must be derived from the same license
// graph.
//
// e.g. combine "restricted", "reciprocal", and "proprietary" resolutions.
func JoinResolutions(resolutions ...ResolutionSet) ResolutionSet {
	if len(resolutions) < 1 {
		panic(fmt.Errorf("attempt to join 0 resolution sets"))
	}
	rmap := make(map[string]*licenseConditionSetImp)
	lg := resolutions[0].(*resolutionSetImp).lg
	for _, r := range resolutions {
		if len(r.(*resolutionSetImp).resolutions) < 1 {
			continue
		}
		if lg == nil {
			lg = r.(*resolutionSetImp).lg
		} else if r.(*resolutionSetImp).lg == nil {
			panic(fmt.Errorf("attempt to join nil resolution set with %d resolutions", len(r.(*resolutionSetImp).resolutions)))
		} else if lg != r.(*resolutionSetImp).lg {
			panic(fmt.Errorf("attempt to join resolutions from multiple graphs"))
		}
		for appliesTo, cs := range r.(*resolutionSetImp).resolutions {
			if cs.Count() < 1 {
				continue
			}
			if _, ok := rmap[appliesTo]; !ok {
				rmap[appliesTo] = cs.Copy().(*licenseConditionSetImp)
				continue
			}
			rmap[appliesTo].addSet(cs)
		}
	}
	return &resolutionSetImp{lg, rmap}
}

// resolutionSetImp implements ResolutionSet.
type resolutionSetImp struct {
	// lg defines the scope of the resolutions to one LicenseGraph
	lg *licenseGraphImp

	// resolutions maps names of target with applicable conditions to the set of conditions that apply.
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

// Conditions returns the set of conditions applied to `appliesTo`.
//
// Panics if `appliesTo` does not appear in the set.
func (rs *resolutionSetImp) Conditions(appliesTo TargetNode) LicenseConditionSet {
	timp := appliesTo.(targetNodeImp)
	if rs.lg == nil {
		rs.lg = timp.lg
	} else if timp.lg == nil {
		timp.lg = rs.lg
	} else if rs.lg != timp.lg {
		panic(fmt.Errorf("attempt to query target resolutions for wrong graph"))
	}
	cs, ok := rs.resolutions[timp.file]
	if !ok {
		return newLicenseConditionSet(&timp)
	}
	return cs.Copy()
}

// Origins identifies the list of originating targets with conditions to resolve. (unordered)
func (rs *resolutionSetImp) Origins() []TargetNode {
	oset := make(map[string]interface{})
	for _, cs := range rs.resolutions {
		for _, origins := range cs.conditions {
			for origin := range origins {
				oset[origin] = nil
			}
		}
	}
	origins := make([]TargetNode, 0, len(oset))
	for origin := range oset {
		origins = append(origins, targetNodeImp{rs.lg, origin})
	}
	return origins
}

// AppliesToByOrigin identifies the list of targets requiring action to resolve conditions originating at `origin`. (unordered)
func (rs *resolutionSetImp) AppliesToByOrigin(origin TargetNode) []TargetNode {
	oimp := origin.(targetNodeImp)
	if rs.lg == nil {
		rs.lg = oimp.lg
	} else if oimp.lg == nil {
		oimp.lg = rs.lg
	} else if rs.lg != oimp.lg {
		panic(fmt.Errorf("attempt to query targets by origin for wrong graph"))
	}
	tset := make(map[string]interface{})
	for target, cs := range rs.resolutions {
		if cs.hasByOrigin(oimp.file) {
			tset[target] = nil
		}
	}
	targets := make([]TargetNode, 0, len(tset))
	for target := range tset {
		targets = append(targets, targetNodeImp{rs.lg, target})
	}
	return targets
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
func (rs resolutionSetImp) AnyByNameApplyToTarget(appliesTo TargetNode, names ...ConditionNames) bool {
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
func (rs resolutionSetImp) AllByNameApplyToTarget(appliesTo TargetNode, names ...ConditionNames) bool {
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

// newResolutionSet constructs a new, empty instance of resolutionSetImp for graph `lg`.
func newResolutionSet(lg *licenseGraphImp) *resolutionSetImp {
	return &resolutionSetImp{lg, make(map[string]*licenseConditionSetImp)}
}

// add applies all of the license conditions in `cs` to `file` if not already applied.
func (rs *resolutionSetImp) add(file string, cs *licenseConditionSetImp) {
	if r, ok := rs.resolutions[file]; ok {
		r.addSet(cs)
	} else {
		rs.resolutions[file] = cs
	}
}

// hasTarget returns true if `file` appears as one of the targets in the set.
func (rs *resolutionSetImp) hasTarget(file string) bool {
	_, isPresent := rs.resolutions[file]
	return isPresent
}

// hasAnyByName returns true if the target for `file` has at least 1 condition matching `names`.
func (rs *resolutionSetImp) hasAnyByName(file string, names ...ConditionNames) bool {
	cs, isPresent := rs.resolutions[file]
	if !isPresent {
		return false
	}
	for _, cn := range names {
		for _, name := range cn {
			_, isPresent = cs.conditions[name]
			if isPresent {
				return true
			}
		}
	}
	return false
}

// hasAllByName returns true if the target for `file` has at least 1 condition for each element of `names`.
func (rs *resolutionSetImp) hasAllByName(file string, names ...ConditionNames) bool {
	cs, isPresent := rs.resolutions[file]
	if !isPresent {
		return false
	}
	for _, cn := range names {
		found := false
		for _, name := range cn {
			_, isPresent = cs.conditions[name]
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
