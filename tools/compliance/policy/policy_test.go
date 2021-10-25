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
	"bytes"
	"fmt"
	"sort"
	"testing"
)

const (
	GPL = `` +
`package_name: "Free Software"
license_kinds: "SPDX-license-identifier-GPL-2.0"
license_conditions: "restricted"
`

	Classpath = `` +
`package_name: "Free Software"
license_kinds: "SPDX-license-identifier-GPL-2.0-with-classpath-exception"
license_conditions: "restricted"
`

	DependentModule = `` +
`package_name: "Free Software"
license_kinds: "SPDX-license-identifier-MIT"
license_conditions: "notice"
`

	LGPL = `` +
`package_name: "Free Library"
license_kinds: "SPDX-license-identifier-LGPL-2.0"
license_conditions: "restricted"
`

	MPL = `` +
`package_name: "Reciprocal"
license_kinds: "SPDX-license-identifier-MPL-2.0"
license_conditions: "reciprocal"
`

	MIT = `` +
`package_name: "Android"
license_kinds: "SPDX-license-identifier-MIT"
license_conditions: "notice"
`

	Proprietary = `` +
`package_name: "Android"
license_kinds: "legacy_proprietary"
license_conditions: "proprietary"
`

	ByException = `` +
`package_name: "Special"
license_kinds: "legacy_by_exception_only"
license_conditions: "by_exception_only"
`

)

var (
	meta = map[string]string{
		"apacheBin.meta_lic": AOSP,
		"apacheLib.meta_lic": AOSP,
		"apacheContainer.meta_lic": AOSP,
		"dependentModule.meta_lic": DependentModule,
		"gplWithClasspathException.meta_lic": Classpath,
		"gplBin.meta_lic": GPL,
		"gplLib.meta_lic": GPL,
		"lgplBin.meta_lic": LGPL,
		"lgplLib.meta_lic": LGPL,
		"mitBin.meta_lic": MIT,
		"mitLib.meta_lic": MIT,
		"mplBin.meta_lic": MPL,
		"mplLib.meta_lic": MPL,
		"proprietary.meta_lic": Proprietary,
		"by_exception.meta_lic": ByException,
	}
)

func TestPolicy(t *testing.T) {
	tests := []struct {
		name       string
		edge       annotated
		treatAsAggregate bool
		expectedDepConditions []string
		expectedTargetConditions []string
	}{
		{
			name: "firstparty",
			edge: annotated{"apacheBin.meta_lic", "apacheLib.meta_lic", []string{"static"}},
			expectedDepConditions: []string{"apacheLib.meta_lic:notice"},
			expectedTargetConditions: []string{},
		},
		{
			name: "notice",
			edge: annotated{"mitBin.meta_lic", "mitLib.meta_lic", []string{"static"}},
			expectedDepConditions: []string{"mitLib.meta_lic:notice"},
			expectedTargetConditions: []string{},
		},
		{
			name: "fponlgpl",
			edge: annotated{"apacheBin.meta_lic", "lgplLib.meta_lic", []string{"static"}},
			expectedDepConditions: []string{"lgplLib.meta_lic:restricted"},
			expectedTargetConditions: []string{},
		},
		{
			name: "fponlgpldynamic",
			edge: annotated{"apacheBin.meta_lic", "lgplLib.meta_lic", []string{"dynamic"}},
			expectedDepConditions: []string{},
			expectedTargetConditions: []string{},
		},
		{
			name: "fpongpl",
			edge: annotated{"apacheBin.meta_lic", "gplLib.meta_lic", []string{"static"}},
			expectedDepConditions: []string{"gplLib.meta_lic:restricted"},
			expectedTargetConditions: []string{},
		},
		{
			name: "fpongpldynamic",
			edge: annotated{"apacheBin.meta_lic", "gplLib.meta_lic", []string{"dynamic"}},
			expectedDepConditions: []string{"gplLib.meta_lic:restricted"},
			expectedTargetConditions: []string{},
		},
		{
			name: "independentmodule",
			edge: annotated{"apacheBin.meta_lic", "gplWithClasspathException.meta_lic", []string{"dynamic"}},
			expectedDepConditions: []string{},
			expectedTargetConditions: []string{},
		},
		{
			name: "independentmodulestatic",
			edge: annotated{"apacheBin.meta_lic", "gplWithClasspathException.meta_lic", []string{"static"}},
			expectedDepConditions: []string{"gplWithClasspathException.meta_lic:restricted"},
			expectedTargetConditions: []string{},
		},
		{
			name: "dependentmodule",
			edge: annotated{"dependentModule.meta_lic", "gplWithClasspathException.meta_lic", []string{"dynamic"}},
			expectedDepConditions: []string{"gplWithClasspathException.meta_lic:restricted"},
			expectedTargetConditions: []string{},
		},

		{
			name: "lgplonfp",
			edge: annotated{"lgplBin.meta_lic", "apacheLib.meta_lic", []string{"static"}},
			expectedDepConditions: []string{"apacheLib.meta_lic:notice"},
			expectedTargetConditions: []string{"lgplBin.meta_lic:restricted"},
		},
		{
			name: "lgplonfpdynamic",
			edge: annotated{"lgplBin.meta_lic", "apacheLib.meta_lic", []string{"dynamic"}},
			expectedDepConditions: []string{},
			expectedTargetConditions: []string{},
		},
		{
			name: "gplonfp",
			edge: annotated{"gplBin.meta_lic", "apacheLib.meta_lic", []string{"static"}},
			expectedDepConditions: []string{"apacheLib.meta_lic:notice"},
			expectedTargetConditions: []string{"gplBin.meta_lic:restricted"},
		},
		{
			name: "gplonfpdynamic",
			edge: annotated{"gplBin.meta_lic", "apacheLib.meta_lic", []string{"dynamic"}},
			expectedDepConditions: []string{},
			expectedTargetConditions: []string{"gplBin.meta_lic:restricted"},
		},
		{
			name: "independentmodulereverse",
			edge: annotated{"gplWithClasspathException.meta_lic", "apacheBin.meta_lic", []string{"dynamic"}},
			expectedDepConditions: []string{},
			expectedTargetConditions: []string{},
		},
		{
			name: "independentmodulereversestatic",
			edge: annotated{"gplWithClasspathException.meta_lic", "apacheBin.meta_lic", []string{"static"}},
			expectedDepConditions: []string{"apacheBin.meta_lic:notice"},
			expectedTargetConditions: []string{"gplWithClasspathException.meta_lic:restricted"},
		},
		{
			name: "dependentmodulereverse",
			edge: annotated{"gplWithClasspathException.meta_lic", "dependentModule.meta_lic", []string{"dynamic"}},
			expectedDepConditions: []string{},
			expectedTargetConditions: []string{"gplWithClasspathException.meta_lic:restricted"},
		},
		{
			name: "ponr",
			edge: annotated{"proprietary.meta_lic", "gplLib.meta_lic", []string{"static"}},
			expectedDepConditions: []string{"gplLib.meta_lic:restricted"},
			expectedTargetConditions: []string{},
		},
		{
			name: "ronp",
			edge: annotated{"gplBin.meta_lic", "proprietary.meta_lic", []string{"static"}},
			expectedDepConditions: []string{"proprietary.meta_lic:proprietary"},
			expectedTargetConditions: []string{"gplBin.meta_lic:restricted"},
		},
		{
			name: "noticeonb_e_o",
			edge: annotated{"mitBin.meta_lic", "by_exception.meta_lic", []string{"static"}},
			expectedDepConditions: []string{"by_exception.meta_lic:by_exception_only"},
			expectedTargetConditions: []string{},
		},
		{
			name: "b_e_oonnotice",
			edge: annotated{"by_exception.meta_lic", "mitLib.meta_lic", []string{"static"}},
			expectedDepConditions: []string{"mitLib.meta_lic:notice"},
			expectedTargetConditions: []string{},
		},
		{
			name: "noticeonrecip",
			edge: annotated{"mitBin.meta_lic", "mplLib.meta_lic", []string{"static"}},
			expectedDepConditions: []string{"mplLib.meta_lic:reciprocal"},
			expectedTargetConditions: []string{},
		},
		{
			name: "reciponnotice",
			edge: annotated{"mplBin.meta_lic", "mitLib.meta_lic", []string{"static"}},
			expectedDepConditions: []string{"mitLib.meta_lic:notice"},
			expectedTargetConditions: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := make(testFS)
			stderr := &bytes.Buffer{}
			target := meta[tt.edge.target] + fmt.Sprintf("deps: {\n  file: \"%s\"\n", tt.edge.dep)
			for _, ann := range tt.edge.annotations {
				target += fmt.Sprintf("  annotations: \"%s\"\n", ann)
			}
			fs[tt.edge.target] = []byte(target + "}\n")
			fs[tt.edge.dep] = []byte(meta[tt.edge.dep])
			lg, err := ReadLicenseGraph(&fs, stderr, []string{tt.edge.target})
			if err != nil {
				t.Errorf("unexpected error reading graph: %w", err)
				return
			}
			if tt.expectedDepConditions != nil {
				cs := depConditionsApplicableToTarget(
					lg.Edges()[0].(targetEdgeImp),
					lg.TargetNode(tt.edge.dep).LicenseConditions().(*licenseConditionSetImp),
					tt.treatAsAggregate)
				actual := make([]string, 0, cs.Count())
				for _, lc := range cs.AsList() {
					actual = append(actual, lc.(licenseConditionImp).asString(":"))
				}
				sort.Strings(actual)
				sort.Strings(tt.expectedDepConditions)
				if len(actual) != len(tt.expectedDepConditions) {
					t.Errorf("unexpected number of dependency conditions: got %v with %d conditions, want %v with %d conditions",
						actual, len(actual), tt.expectedDepConditions, len(tt.expectedDepConditions))
				} else {
					for i := 0; i < len(actual); i++ {
						if actual[i] != tt.expectedDepConditions[i] {
							t.Errorf("unexpected dependency condition at element %d: got %q, want %q",
								i, actual[i], tt.expectedDepConditions[i])
						}
					}
				}
			}
			if tt.expectedTargetConditions != nil {
				cs := targetConditionsApplicableToDep(
					lg.Edges()[0].(targetEdgeImp),
					lg.TargetNode(tt.edge.target).LicenseConditions().(*licenseConditionSetImp),
					tt.treatAsAggregate)
				actual := make([]string, 0, cs.Count())
				for _, lc := range cs.AsList() {
					actual = append(actual, lc.(licenseConditionImp).asString(":"))
				}
				sort.Strings(actual)
				sort.Strings(tt.expectedTargetConditions)
				if len(actual) != len(tt.expectedTargetConditions) {
					t.Errorf("unexpected number of target conditions: got %v with %d conditions, want %v with %d conditions",
						actual, len(actual), tt.expectedTargetConditions, len(tt.expectedTargetConditions))
				} else {
					for i := 0; i < len(actual); i++ {
						if actual[i] != tt.expectedTargetConditions[i] {
							t.Errorf("unexpected target condition at element %d: got %q, want %q",
								i, actual[i], tt.expectedTargetConditions[i])
						}
					}
				}
			}
		})
	}
}
