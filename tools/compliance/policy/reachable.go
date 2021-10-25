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

// ReachableNodes returns the set of reachable nodes in a license graph. (caches result)
func ReachableNodes(lg *LicenseGraph) *TargetNodeSet {
	lg.mu.Lock()
	reachable := lg.reachableNodes
	lg.mu.Unlock()
	if reachable != nil {
		return reachable
	}

	rmap := make(map[string]interface{})

	WalkTopDown(lg, func(lg *LicenseGraph, tn TargetNode, path []TargetEdge) bool {
		if _, alreadyWalked := rmap[tn.file]; alreadyWalked {
			return false
		}
		if len(path) > 0 {
			if !edgeIsDerivation(path[len(path)-1].e) {
				return false
			}
		}
		rmap[tn.file] = nil
		return true
	})

	reachable = &TargetNodeSet{lg, rmap}

	lg.mu.Lock()
	if lg.reachableNodes == nil {
		lg.reachableNodes = reachable
	}
	lg.mu.Unlock()
	return reachable
}
