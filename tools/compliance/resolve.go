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

func ResolveDependencyConditions(e TargetEdge) LicenseConditionSet {
	eimp := e.(targetEdgeImp)
	cs := eimp.lg.targets[eimp.e.dependency].licenseConditions.Copy()
	result := cs.(licenseConditionSetImp)

	if !eimp.isDerivativeOf() {
		result.removeAllByName("unencumbered")
		result.removeAllByName("permissive")
		result.removeAllByName("notice")
		result.removeAllByName("proprietary")
	}
	return result
}

func ResolvePathConditions(path TargetPath) LicenseConditionSet {
	result := newLicenseConditionSet(nil)
	if len(path) == 0 {
		return result
	}
	for i := len(path)-1; i >= 0; i-- {
		eimp := path[i].(targetEdgeImp)
		result = result.union(&eimp.lg.targets[eimp.e.dependency].licenseConditions)
		if !eimp.isDerivativeOf() {
			result.removeAllByName("unencumbered")
			result.removeAllByName("permissive")
			result.removeAllByName("notice")
			result.removeAllByName("proprietary")
		}
	}
	return result
}

func ResolveGraphConditions(graph LicenseGraph) ResolutionSet {
}
