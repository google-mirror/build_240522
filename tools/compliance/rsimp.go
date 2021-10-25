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

// JoinResolutions returns a new ResolutionSet combining the resolutions from
// multiple resolution sets. All sets must be derived from the same license
// graph.
//
// e.g. combine "restricted", "reciprocal", and "proprietary" resolutions.
func JoinResolutions(resolutions ...ResolutionSet) ResolutionSet {
	if len(resolutions) < 1 {
		panic(fmt.Errorf("attempt to join 0 resolution sets"))
	}
	rmap := make(map[string]actionSet)
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
		for attachesTo, as := range r.(*resolutionSetImp).resolutions {
			if as.isEmpty() {
				continue
			}
			if _, ok := rmap[attachesTo]; !ok {
				rmap[attachesTo] = as.copy()
				continue
			}
			rmap[attachesTo].addSet(as)
		}
	}
	return &resolutionSetImp{lg, rmap}
}


// resolutionImp implements Resolution.
type resolutionImp struct {
	attachesTo, actsOn string
	cs                 *licenseConditionSetImp
}

func (r resolutionImp) AttachesTo() TargetNode {
	return targetNodeImp{r.cs.lg, r.attachesTo}
}

func (r resolutionImp) ActsOn() TargetNode {
	return targetNodeImp{r.cs.lg, r.actsOn}
}

func (r resolutionImp) Resolves() LicenseConditionSet {
	return r.cs.Copy()
}

func (r resolutionImp) asString() string {
	var sb strings.Builder
	cl := r.cs.AsList()
	sort.Sort(cl)
	fmt.Fprintf(&sb, "%s -> %s -> %s", r.attachesTo, r.actsOn, cl.String())
	return sb.String()
}

func (r resolutionImp) byName(names ConditionNames) resolutionImp {
	return resolutionImp{r.attachesTo, r.actsOn, r.cs.ByName(names).(*licenseConditionSetImp)}
}

func (r resolutionImp) byOrigin(origin TargetNode) resolutionImp {
	return resolutionImp{r.attachesTo, r.actsOn, r.cs.ByOrigin(origin).(*licenseConditionSetImp)}
}

func (r resolutionImp) copy() resolutionImp {
	return resolutionImp{r.attachesTo, r.actsOn, r.cs.Copy().(*licenseConditionSetImp)}
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
		fmt.Fprintf(&sb, "%s%s", sep, r.(resolutionImp).asString())
		sep = ", "
	}
	fmt.Fprintf(&sb, "]")
	return sb.String()
}

func (rl ResolutionList) AllConditions() LicenseConditionSet {
	result := newLicenseConditionSet(nil)
	for _, r := range rl {
		result.addSet(r.(resolutionImp).cs)
	}
	return result
}

func (rl ResolutionList) ByName(names ConditionNames) ResolutionList {
	result := make(ResolutionList, 0, rl.CountByName(names))
	for _, r := range rl {
		if r.Resolves().HasAnyByName(names) {
			result = append(result, r.(resolutionImp).byName(names))
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
		if r.(resolutionImp).attachesTo == attachesTo.Name() {
			result = append(result, r.(resolutionImp).copy())
		}
	}
	return result
}

func (rl ResolutionList) CountByAttachesTo(attachesTo TargetNode) int {
	c := 0
	for _, r := range rl {
		if r.(resolutionImp).attachesTo == attachesTo.Name() {
			c++
		}
	}
	return c
}

func (rl ResolutionList) ByActsOn(actsOn TargetNode) ResolutionList {
	result := make(ResolutionList, 0, rl.CountByActsOn(actsOn))
	for _, r := range rl {
		if r.(resolutionImp).actsOn == actsOn.Name() {
			result = append(result, r.(resolutionImp).copy())
		}
	}
	return result
}

func (rl ResolutionList) CountByActsOn(actsOn TargetNode) int {
	c := 0
	for _, r := range rl {
		if r.(resolutionImp).actsOn == actsOn.Name() {
			c++
		}
	}
	return c
}

func (rl ResolutionList) ByOrigin(origin TargetNode) ResolutionList {
	result := make(ResolutionList, 0, rl.CountByOrigin(origin))
	for _, r := range rl {
		if r.Resolves().HasAnyByOrigin(origin) {
			result = append(result, r.(resolutionImp).byOrigin(origin))
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



type actionSet map[string]*licenseConditionSetImp

func (as actionSet) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "{")
	osep := ""
	for actsOn, cs := range as {
		cl := cs.AsList()
		sort.Sort(cl)
		fmt.Fprintf(&sb, "%s%s -> [%s]", osep, actsOn, cl.String())
		osep = ", "
	}
	fmt.Fprintf(&sb, "}")
	return sb.String()
}

func (as actionSet) byName(names ConditionNames) actionSet {
	result := make(actionSet)
	for actsOn, cs := range as {
		bn := cs.ByName(names).(*licenseConditionSetImp)
		if bn.IsEmpty() {
			continue
		}
		result[actsOn] = bn
	}
	return result
}

func (as actionSet) filterActsOn(reachable *targetNodeSet) actionSet {
	result := make(actionSet)
	for actsOn, cs := range as {
		if !reachable.contains(actsOn) || cs.IsEmpty() {
			continue
		}
		result[actsOn] = cs.Copy().(*licenseConditionSetImp)
	}
	return result
}

func (as actionSet) copy() actionSet {
	result := make(actionSet)
	for actsOn, cs := range as {
		if cs.IsEmpty() {
			continue
		}
		result[actsOn] = cs.Copy().(*licenseConditionSetImp)
	}
	return result
}

func (as actionSet) addSet(other actionSet) {
	for actsOn, cs := range other {
		as.add(actsOn, cs)
	}
}

func (as actionSet) add(actsOn string, cs *licenseConditionSetImp) {
	if acs, ok := as[actsOn]; ok {
		acs.addSet(cs)
	} else {
		as[actsOn] = cs.Copy().(*licenseConditionSetImp)
	}
}

func (as actionSet) addCondition(actsOn string, lc licenseConditionImp) {
	if _, ok := as[actsOn]; !ok {
		as[actsOn] = newLicenseConditionSet(&lc.origin)
	}
	as[actsOn].Add(lc)
}

func (as actionSet) isEmpty() bool {
	for _, cs := range as {
		if !cs.IsEmpty() {
			return false
		}
	}
	return true
}

// resolutionSetImp implements ResolutionSet.
type resolutionSetImp struct {
	// lg defines the scope of the resolutions to one LicenseGraph
	lg *licenseGraphImp

	// resolutions maps names of target with applicable conditions to the set of conditions that apply.
	resolutions map[string]actionSet
}

func (rs *resolutionSetImp) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "{")
	sep := ""
	for attachesTo, as := range rs.resolutions {
		fmt.Fprintf(&sb, "%s%s -> %s", sep, attachesTo, as.String())
		sep = ", "
	}
	fmt.Fprintf(&sb, "}")
	return sb.String()
}

// AppliesTo returns the list of targets with applicable license conditions attached.
func (rs *resolutionSetImp) AttachesTo() []TargetNode {
	targets := make([]TargetNode, 0, len(rs.resolutions))
	for target := range rs.resolutions {
		targets = append(targets, targetNodeImp{rs.lg, target})
	}
	return targets
}

// ActsOn identifies the list of targets to act on (share, give notice etc.) to resolve conditions. (unordered)
func (rs *resolutionSetImp) ActsOn() []TargetNode {
	tset := make(map[string]interface{})
	for _, as := range rs.resolutions {
		for actsOn := range as {
			tset[actsOn] = nil
		}
	}
	targets := make([]TargetNode, 0, len(tset))
	for target := range tset {
		targets = append(targets, targetNodeImp{rs.lg, target})
	}
	return targets
}

// Origins identifies the list of originating targets with conditions to resolve. (unordered)
func (rs *resolutionSetImp) Origins() []TargetNode {
	tset := make(map[string]interface{})
	for _, as := range rs.resolutions {
		for _, cs := range as {
			for _, origins := range cs.conditions {
				for origin := range origins {
					tset[origin] = nil
				}
			}
		}
	}
	targets := make([]TargetNode, 0, len(tset))
	for target := range tset {
		targets = append(targets, targetNodeImp{rs.lg, target})
	}
	return targets
}

// Resolutions returns the set of resolutions attached to `attachedTo`.
//
// Panics if `attachedTo` does not appear in the set.
func (rs *resolutionSetImp) Resolutions(attachedTo TargetNode) []Resolution {
	timp := attachedTo.(targetNodeImp)
	if rs.lg == nil {
		rs.lg = timp.lg
	} else if timp.lg == nil {
		timp.lg = rs.lg
	} else if rs.lg != timp.lg {
		panic(fmt.Errorf("attempt to query target resolutions for wrong graph"))
	}
	as, ok := rs.resolutions[timp.file]
	if !ok {
		return []Resolution{}
	}
	result := make([]Resolution, 0, len(as))
	for actsOn, cs := range as {
		result = append(result, resolutionImp{timp.file, actsOn, cs.Copy().(*licenseConditionSetImp)})
	}
	return result
}

// ResolutionsByActsOn returns the set of resolutions that act on `actOn`.
//
// Panics if `attachedTo` does not appear in the set.
func (rs *resolutionSetImp) ResolutionsByActsOn(actOn TargetNode) []Resolution {
	timp := actOn.(targetNodeImp)
	if rs.lg == nil {
		rs.lg = timp.lg
	} else if timp.lg == nil {
		timp.lg = rs.lg
	} else if rs.lg != timp.lg {
		panic(fmt.Errorf("attempt to query target resolutions for wrong graph"))
	}
	c := 0
	for _, as := range rs.resolutions {
		if _, ok := as[actOn.Name()]; ok {
			c++
		}
	}
	result := make([]Resolution, 0, c)
	for attachedTo, as := range rs.resolutions {
		if cs, ok := as[actOn.Name()]; ok {
			result = append(result, resolutionImp{attachedTo, actOn.Name(), cs.Copy().(*licenseConditionSetImp)})
		}
	}
	return result
}

// AttachesToByOrigin identifies the list of targets requiring action to resolve conditions originating at `origin`. (unordered)
func (rs *resolutionSetImp) AttachesToByOrigin(origin TargetNode) []TargetNode {
	oimp := origin.(targetNodeImp)
	if rs.lg == nil {
		rs.lg = oimp.lg
	} else if oimp.lg == nil {
		oimp.lg = rs.lg
	} else if rs.lg != oimp.lg {
		panic(fmt.Errorf("attempt to query targets by origin for wrong graph"))
	}
	tset := make(map[string]interface{})
	for target, as := range rs.resolutions {
		for _, cs := range as {
			if cs.hasByOrigin(oimp.file) {
				tset[target] = nil
				break
			}
		}
	}
	targets := make([]TargetNode, 0, len(tset))
	for target := range tset {
		targets = append(targets, targetNodeImp{rs.lg, target})
	}
	return targets
}

// AttachesToTarget returns true if `attacheTo` appears in the set.
func (rs *resolutionSetImp) AttachesToTarget(attachedTo TargetNode) bool {
	timp := attachedTo.(targetNodeImp)
	if rs.lg == nil {
		rs.lg = timp.lg
	} else if timp.lg == nil {
		timp.lg = rs.lg
	} else if rs.lg != timp.lg {
		panic(fmt.Errorf("attempt to query resolved targets for wrong graph"))
	}
	return rs.hasTarget(timp.file)
}

// AnyByNameAttachToTarget returns true if `attachedTo` appears in the set with any conditions matching `name`.
func (rs resolutionSetImp) AnyByNameAttachToTarget(attachedTo TargetNode, names ...ConditionNames) bool {
	timp := attachedTo.(targetNodeImp)
	if rs.lg == nil {
		rs.lg = timp.lg
	} else if timp.lg == nil {
		timp.lg = rs.lg
	} else if rs.lg != timp.lg {
		panic(fmt.Errorf("attempt to query target resolutions for wrong graph"))
	}
	return rs.hasAnyByName(timp.file, names...)
}

// AllByNameAttachToTarget returns true if `attachedTo` appears in the set with conditions matching every element of `name`.
func (rs resolutionSetImp) AllByNameAttachToTarget(attachedTo TargetNode, names ...ConditionNames) bool {
	timp := attachedTo.(targetNodeImp)
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
	return &resolutionSetImp{lg, make(map[string]actionSet)}
}

// add attaches all of the license conditions in `cs` to `file` to act on the originating node if not already applied.
func (rs *resolutionSetImp) add(file string, cs *licenseConditionSetImp) {
	if cs.IsEmpty() {
		return
	}
	_, ok := rs.resolutions[file]
	if !ok {
		rs.resolutions[file] = make(actionSet)
	}

	for name, origins := range cs.conditions {
		for origin := range origins {
			if _, ok := rs.resolutions[file][origin]; !ok {
				rs.resolutions[file][origin] = newLicenseConditionSet(&targetNodeImp{rs.lg, origin}, name)
			} else {
				rs.resolutions[file][origin].addAll(origin, name)
			}
		}
	}
}

// add attaches all of the license conditions in `as` to `file` to act on the originating node if not already applied.
func (rs *resolutionSetImp) addConditions(file string, as actionSet) {
	_, ok := rs.resolutions[file]
	if !ok {
		rs.resolutions[file] = as.copy()
		return
	}
	rs.resolutions[file].addSet(as)
}

// add attaches all of the license conditions in `as` to `file` to act on `file` if not already applied.
func (rs *resolutionSetImp) addSelf(file string, as actionSet) {
	for _, cs := range as {
		if cs.IsEmpty() {
			return
		}
		_, ok := rs.resolutions[file]
		if !ok {
			rs.resolutions[file] = make(actionSet)
		}
		_, ok = rs.resolutions[file][file]
		if !ok {
			rs.resolutions[file][file] = newLicenseConditionSet(&targetNodeImp{rs.lg, file})
		}
		rs.resolutions[file][file].addSet(cs)
	}
}

// hasTarget returns true if `file` appears as one of the targets in the set.
func (rs *resolutionSetImp) hasTarget(file string) bool {
	_, isPresent := rs.resolutions[file]
	return isPresent
}

// hasAnyByName returns true if the target for `file` has at least 1 condition matching `names`.
func (rs *resolutionSetImp) hasAnyByName(file string, names ...ConditionNames) bool {
	objects, isPresent := rs.resolutions[file]
	if !isPresent {
		return false
	}
	for _, cs := range objects {
		for _, cn := range names {
			for _, name := range cn {
				_, isPresent = cs.conditions[name]
				if isPresent {
					return true
				}
			}
		}
	}
	return false
}

// hasAllByName returns true if the target for `file` has at least 1 condition for each element of `names`.
func (rs *resolutionSetImp) hasAllByName(file string, names ...ConditionNames) bool {
	as, isPresent := rs.resolutions[file]
	if !isPresent {
		return false
	}
	for _, cn := range names {
		found := false
	asloop:
		for _, cs := range as {
			for _, name := range cn {
				_, isPresent = cs.conditions[name]
				if isPresent {
					found = true
					break asloop
				}
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// hasResolution returns true if the exact `condition` applies to `file` acting on `actsOn`.
func (rs *resolutionSetImp) hasResolution(file, actsOn string, condition LicenseCondition) bool {
	as, isPresent := rs.resolutions[file]
	if !isPresent {
		return false
	}
	cs, isPresent := as[actsOn]
	if !isPresent {
		return false
	}
	origins, isPresent := cs.conditions[condition.Name()]
	if !isPresent {
		return false
	}
	_, isPresent = origins[condition.Origin().Name()]
	return isPresent
}
