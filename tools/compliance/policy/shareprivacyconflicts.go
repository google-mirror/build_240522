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

// SourceSharePrivacyConflict describes an individual conflict between a source-sharing
// condition and a source privacy condition
type SourceSharePrivacyConflict struct {
	AppliesTo        TargetNode
	ShareCondition   LicenseCondition
	PrivacyCondition LicenseCondition
}

// Error returns a string describing the conflict.
func (conflict SourceSharePrivacyConflict) Error() string {
	return fmt.Sprintf("%s %s from %s and must share from %s %s\n",
		conflict.AppliesTo.Name(),
		conflict.PrivacyCondition.Name(), conflict.PrivacyCondition.Origin().Name(),
		conflict.ShareCondition.Name(), conflict.ShareCondition.Origin().Name())
}

// IsEqualTo returns true when `conflict` and `other` describe the same conflict.
func (conflict SourceSharePrivacyConflict) IsEqualTo(other SourceSharePrivacyConflict) bool {
	return conflict.AppliesTo.Name() == other.AppliesTo.Name() &&
		conflict.ShareCondition.Name() == other.ShareCondition.Name() &&
		conflict.ShareCondition.Origin().Name() == other.ShareCondition.Origin().Name() &&
		conflict.PrivacyCondition.Name() == other.PrivacyCondition.Name() &&
		conflict.PrivacyCondition.Origin().Name() == other.PrivacyCondition.Origin().Name()
}


// ConflictingSharedPrivateSource lists all of the targets where conflicting conditions to
// share the source and to keep the source private apply to the target.
func ConflictingSharedPrivateSource(licenseGraph LicenseGraph) []SourceSharePrivacyConflict {
	// shareSource is the set of all source-sharing resolutions.
	shareSource := ResolveSourceSharing(licenseGraph)

	// privateSource is the set of all source privacy resolutions.
	privateSource := ResolveSourcePrivacy(licenseGraph)

	// combined is the combination of source-sharing and source privacy.
	combined := JoinResolutions(shareSource, privateSource)

	// size is the size of the result
	size := 0
	for _, appliesTo := range combined.AppliesTo() {
		cs := combined.Conditions(appliesTo)
		size += cs.CountByName(ImpliesShared) * cs.CountByName(ImpliesPrivate)
	}
	result := make([]SourceSharePrivacyConflict, 0, size)
	for _, appliesTo := range combined.AppliesTo() {
		cs := combined.Conditions(appliesTo)

		pconditions := cs.ByName(ImpliesPrivate)
		ssconditions := cs.ByName(ImpliesShared)

		// report all conflicting condition combinations
		for _, p := range pconditions {
			for _, ss := range ssconditions {
				result = append(result, SourceSharePrivacyConflict{appliesTo, ss, p})
			}
		}
	}
	return result
}
