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

type LicenseCondition uint16

const LicenseConditionMask = LicenseCondition(0x3ff)

const (
	UnencumberedCondition = LicenseCondition(0x0001)
	PermissiveCondition = LicenseCondition(0x0002)
	NoticeCondition = LicenseCondition(0x0004)
	ReciprocalCondition = LicenseCondition(0x0008)
	RestrictedCondition = LicenseCondition(0x0010)
	RestrictedClasspathExceptionCondition = LicenseCondition(0x0020)
	WeaklyRestrictedCondition = LicenseCondition(0x0040)
	ProprietaryCondition = LicenseCondition(0x0080)
	ByExceptionOnlyCondition = LicenseCondition(0x0100)
	NotAllowedCondition = LicenseCondition(0x0200)
)

var (
	RecognizedConditionNames = map[string]LicenseCondition{
		"unencumbered": UnencumberedCondition,
		"permissive": PermissiveCondition,
		"notice": NoticeCondition,
		"reciprocal": ReciprocalCondition,
		"restricted": RestrictedCondition,
		"restricted_with_classpath_exception": RestrictedClasspathExceptionCondition,
		"restricted_allows_dynamic_linking": WeaklyRestrictedCondition,
		"proprietary": ProprietaryCondition,
		"by_exception_only": ByExceptionOnlyCondition,
		"not_allowed": NotAllowedCondition,
	}
)

func (lc LicenseCondition) Name() string {
	switch lc {
	case UnencumberedCondition:
		return "unencumbered"
	case PermissiveCondition:
		return "permissive"
	case NoticeCondition:
		return "notice"
	case ReciprocalCondition:
		return "reciprocal"
	case RestrictedCondition:
		return "restricted"
	case RestrictedClasspathExceptionCondition:
		return "restricted_with_classpath_exception"
	case WeaklyRestrictedCondition:
		return "restricted_allows_dynamic_linking"
	case ProprietaryCondition:
		return "proprietary"
	case ByExceptionOnlyCondition:
		return "by_exception_only"
	case NotAllowedCondition:
		return "not_allowed"
	}
	panic(fmt.Errorf("unrecognized license condition: %04x", lc))
}
