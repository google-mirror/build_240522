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

// WalkTopDown does a top-down walk of `lg` calling `descend` and descending
// into depenencies when `descend` returns true.
func WalkTopDown(lg *LicenseGraph, descend func(*LicenseGraph, *TargetNode, []TargetEdge) bool) {
	path := make([]TargetEdge, 0, 32)

	var walk func(f string)
	walk = func(f string) {
		if !descend(lg, lg.targets[f], path) {
			return
		}
		for _, edge := range lg.index[f] {
			path = append(path, TargetEdge{lg, edge})
			walk(edge.dependency)
			path = path[:len(path)-1]
		}
	}

	for _, r := range lg.rootFiles {
		walk(r)
	}
}
