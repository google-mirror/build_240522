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
func ResolveSourcePrivacy(lg *LicenseGraph) *ResolutionSet {
	rs := ResolveTopDownConditions(lg)


	if rs.lg == nil {
		panic(fmt.Errorf("nil graph for proprietary"))
	} else if lg == nil {
		panic(fmt.Errorf("nil graph for privacy"))
	} else if lg != rs.lg {
		panic(fmt.Errorf("resolve proprietary for wrong graph"))
	}

	return WalkResolutionsForCondition(rs, ImpliesPrivate)
}
