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

func resolve(e targetEdgeImp, cs LicenseConditionSet) *licenseConditionSetImp {
	result := cs.(*licenseConditionSetImp)
	if !e.isDerivativeOf() {
		result.removeAllByName("unencumbered")
		result.removeAllByName("permissive")
		result.removeAllByName("notice")
		result.removeAllByName("reciprocal")
		result.removeAllByName("proprietary")
	}
	return result
}

func edgeIsDynamicLink(e *dependencyEdge) bool {
	_, isPresent := e.annotations["dynamic"]
	return isPresent
}

func edgeIsDerivativeOf(e *dependencyEdge) bool {
	_, isDynamic := e.annotations["dynamic"]
	_, isToolchain := e.annotations["toolchain"]
	return !isDynamic && !isToolchain
}


func ResolveDependencyConditions(e TargetEdge) LicenseConditionSet {
	eimp := e.(targetEdgeImp)
	return resolve(eimp, eimp.lg.targets[eimp.e.dependency].licenseConditions)
}

func ResolvePathConditions(path TargetPath) LicenseConditionSet {
	result := newLicenseConditionSet(nil)
	if len(path) == 0 {
		return result
	}
	for i := len(path)-1; i >= 0; i-- {
		eimp := path[i].(targetEdgeImp)
		result = eimp.lg.targets[eimp.e.dependency].licenseConditions.union(result)
		result = resolve(eimp, result)
	}
	eimp := path[0].(targetEdgeImp)
	result = eimp.lg.targets[eimp.e.target].licenseConditions.union(result)
	return result
}

func ResolveGraphConditions(graph LicenseGraph) ResolutionSet {
	lg := graph.(*licenseGraphImp)

	lg.mu.Lock()
	rs := lg.rs
	lg.mu.Unlock()

	if rs != nil {
		return rs
	}

	lg.indexForward()

	rs = newResolutionSet(lg)

	var walk func(f string) *licenseConditionSetImp

	walk = func(f string) *licenseConditionSetImp {
		result := newLicenseConditionSet(&targetNodeImp{lg, f})
		if preresolved, ok := rs.resolutions[f]; ok {
			return result.union(preresolved)
		}
		for _, e := range lg.index[f] {
			cs := walk(e.dependency)
			cs = resolve(targetEdgeImp{lg, e}, cs)
			result = result.union(cs)
		}
		result = lg.targets[f].licenseConditions.union(result)
		rs.add(f, result)
		return result
	}

	for _, r := range lg.rootFiles {
		cs := walk(r)
		rs.add(r, lg.targets[r].licenseConditions.union(cs))
	}

	lg.mu.Lock()
	if lg.rs == nil{
		lg.rs = rs
	}
	lg.mu.Unlock()

        return rs
}
