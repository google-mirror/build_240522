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
	"strings"
	"testing"
)

func TestPolicy_edgeConditions(t *testing.T) {
	tests := []struct {
		name                     string
		edge                     annotated
		treatAsAggregate         bool
		otherCondition           string
		expectedDepConditions    []string
		expectedTargetConditions []string
	}{
		{
			name:                     "firstparty",
			edge:                     annotated{"apacheBin.meta_lic", "apacheLib.meta_lic", []string{"static"}},
			expectedDepConditions:    []string{"apacheLib.meta_lic:notice"},
			expectedTargetConditions: []string{},
		},
		{
			name:                     "notice",
			edge:                     annotated{"mitBin.meta_lic", "mitLib.meta_lic", []string{"static"}},
			expectedDepConditions:    []string{"mitLib.meta_lic:notice"},
			expectedTargetConditions: []string{},
		},
		{
			name:                     "fponlgpl",
			edge:                     annotated{"apacheBin.meta_lic", "lgplLib.meta_lic", []string{"static"}},
			expectedDepConditions:    []string{"lgplLib.meta_lic:restricted"},
			expectedTargetConditions: []string{},
		},
		{
			name:                     "fponlgpldynamic",
			edge:                     annotated{"apacheBin.meta_lic", "lgplLib.meta_lic", []string{"dynamic"}},
			expectedDepConditions:    []string{},
			expectedTargetConditions: []string{},
		},
		{
			name:                     "fpongpl",
			edge:                     annotated{"apacheBin.meta_lic", "gplLib.meta_lic", []string{"static"}},
			expectedDepConditions:    []string{"gplLib.meta_lic:restricted"},
			expectedTargetConditions: []string{},
		},
		{
			name:                     "fpongpldynamic",
			edge:                     annotated{"apacheBin.meta_lic", "gplLib.meta_lic", []string{"dynamic"}},
			expectedDepConditions:    []string{"gplLib.meta_lic:restricted"},
			expectedTargetConditions: []string{},
		},
		{
			name:                     "independentmodule",
			edge:                     annotated{"apacheBin.meta_lic", "gplWithClasspathException.meta_lic", []string{"dynamic"}},
			expectedDepConditions:    []string{},
			expectedTargetConditions: []string{},
		},
		{
			name:                     "independentmodulestatic",
			edge:                     annotated{"apacheBin.meta_lic", "gplWithClasspathException.meta_lic", []string{"static"}},
			expectedDepConditions:    []string{"gplWithClasspathException.meta_lic:restricted"},
			expectedTargetConditions: []string{},
		},
		{
			name:                     "dependentmodule",
			edge:                     annotated{"dependentModule.meta_lic", "gplWithClasspathException.meta_lic", []string{"dynamic"}},
			expectedDepConditions:    []string{"gplWithClasspathException.meta_lic:restricted"},
			expectedTargetConditions: []string{},
		},

		{
			name:                     "lgplonfp",
			edge:                     annotated{"lgplBin.meta_lic", "apacheLib.meta_lic", []string{"static"}},
			expectedDepConditions:    []string{"apacheLib.meta_lic:notice"},
			expectedTargetConditions: []string{"lgplBin.meta_lic:restricted"},
		},
		{
			name:                     "lgplonfpdynamic",
			edge:                     annotated{"lgplBin.meta_lic", "apacheLib.meta_lic", []string{"dynamic"}},
			expectedDepConditions:    []string{},
			expectedTargetConditions: []string{},
		},
		{
			name:                     "gplonfp",
			edge:                     annotated{"gplBin.meta_lic", "apacheLib.meta_lic", []string{"static"}},
			expectedDepConditions:    []string{"apacheLib.meta_lic:notice"},
			expectedTargetConditions: []string{"gplBin.meta_lic:restricted"},
		},
		{
			name:                     "gplcontainer",
			edge:                     annotated{"gplContainer.meta_lic", "apacheLib.meta_lic", []string{"static"}},
			treatAsAggregate:         true,
			expectedDepConditions:    []string{"apacheLib.meta_lic:notice"},
			expectedTargetConditions: []string{"gplContainer.meta_lic:restricted"},
		},
		{
			name:                     "gploncontainer",
			edge:                     annotated{"apacheContainer.meta_lic", "apacheLib.meta_lic", []string{"static"}},
			treatAsAggregate:         true,
			otherCondition:           "gplLib.meta_lic:restricted",
			expectedDepConditions:    []string{"apacheLib.meta_lic:notice", "gplLib.meta_lic:restricted"},
			expectedTargetConditions: []string{},
		},
		{
			name:                     "gplonbin",
			edge:                     annotated{"apacheBin.meta_lic", "apacheLib.meta_lic", []string{"static"}},
			treatAsAggregate:         false,
			otherCondition:           "gplLib.meta_lic:restricted",
			expectedDepConditions:    []string{"apacheLib.meta_lic:notice", "gplLib.meta_lic:restricted"},
			expectedTargetConditions: []string{"gplLib.meta_lic:restricted"},
		},
		{
			name:                     "gplonfpdynamic",
			edge:                     annotated{"gplBin.meta_lic", "apacheLib.meta_lic", []string{"dynamic"}},
			expectedDepConditions:    []string{},
			expectedTargetConditions: []string{"gplBin.meta_lic:restricted"},
		},
		{
			name:                     "independentmodulereverse",
			edge:                     annotated{"gplWithClasspathException.meta_lic", "apacheBin.meta_lic", []string{"dynamic"}},
			expectedDepConditions:    []string{},
			expectedTargetConditions: []string{},
		},
		{
			name:                     "independentmodulereversestatic",
			edge:                     annotated{"gplWithClasspathException.meta_lic", "apacheBin.meta_lic", []string{"static"}},
			expectedDepConditions:    []string{"apacheBin.meta_lic:notice"},
			expectedTargetConditions: []string{"gplWithClasspathException.meta_lic:restricted"},
		},
		{
			name:                     "dependentmodulereverse",
			edge:                     annotated{"gplWithClasspathException.meta_lic", "dependentModule.meta_lic", []string{"dynamic"}},
			expectedDepConditions:    []string{},
			expectedTargetConditions: []string{"gplWithClasspathException.meta_lic:restricted"},
		},
		{
			name:                     "ponr",
			edge:                     annotated{"proprietary.meta_lic", "gplLib.meta_lic", []string{"static"}},
			expectedDepConditions:    []string{"gplLib.meta_lic:restricted"},
			expectedTargetConditions: []string{},
		},
		{
			name:                     "ronp",
			edge:                     annotated{"gplBin.meta_lic", "proprietary.meta_lic", []string{"static"}},
			expectedDepConditions:    []string{"proprietary.meta_lic:proprietary"},
			expectedTargetConditions: []string{"gplBin.meta_lic:restricted"},
		},
		{
			name:                     "noticeonb_e_o",
			edge:                     annotated{"mitBin.meta_lic", "by_exception.meta_lic", []string{"static"}},
			expectedDepConditions:    []string{"by_exception.meta_lic:by_exception_only"},
			expectedTargetConditions: []string{},
		},
		{
			name:                     "b_e_oonnotice",
			edge:                     annotated{"by_exception.meta_lic", "mitLib.meta_lic", []string{"static"}},
			expectedDepConditions:    []string{"mitLib.meta_lic:notice"},
			expectedTargetConditions: []string{},
		},
		{
			name:                     "noticeonrecip",
			edge:                     annotated{"mitBin.meta_lic", "mplLib.meta_lic", []string{"static"}},
			expectedDepConditions:    []string{"mplLib.meta_lic:reciprocal"},
			expectedTargetConditions: []string{},
		},
		{
			name:                     "reciponnotice",
			edge:                     annotated{"mplBin.meta_lic", "mitLib.meta_lic", []string{"static"}},
			expectedDepConditions:    []string{"mitLib.meta_lic:notice"},
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
			licenseGraph, err := ReadLicenseGraph(&fs, stderr, []string{tt.edge.target})
			if err != nil {
				t.Errorf("unexpected error reading graph: %w", err)
				return
			}
			lg := licenseGraph.(*licenseGraphImp)
			// simulate a condition inherited from another edge/dependency.
			otherTarget := ""
			otherCondition := ""
			if len(tt.otherCondition) > 0 {
				fields := strings.Split(tt.otherCondition, ":")
				otherTarget = fields[0]
				otherCondition = fields[1]
				// other target must exist in graph
				lg.targets[otherTarget] = &targetNode{name: otherTarget}
				lg.targets[otherTarget].LicenseConditions = append(lg.targets[otherTarget].LicenseConditions, otherCondition)
			}
			if tt.expectedDepConditions != nil {
				depConditions := lg.TargetNode(tt.edge.dep).LicenseConditions().(*licenseConditionSetImp)
				if otherTarget != "" {
					depConditions.addAll(otherTarget, otherCondition)
				}
				cs := depConditionsApplicableToTarget(lg.Edges()[0].(targetEdgeImp), depConditions, tt.treatAsAggregate)
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
				targetConditions := lg.TargetNode(tt.edge.target).LicenseConditions().(*licenseConditionSetImp)
				if otherTarget != "" {
					targetConditions.addAll(otherTarget, otherCondition)
				}
				cs := targetConditionsApplicableToDep(lg.Edges()[0].(targetEdgeImp), targetConditions, tt.treatAsAggregate)
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
