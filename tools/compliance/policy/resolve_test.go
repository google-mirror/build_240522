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
	"testing"
)

func TestResolveBottomUpConditions(t *testing.T) {
	tests := []struct {
		name                string
		roots               []string
		edges               []annotated
		expectedResolutions []res
	}{
		{
			name: "firstparty",
			roots: []string{"apacheBin.meta_lic"},
			edges: []annotated{
				{"apacheBin.meta_lic", "apacheLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "apacheLib.meta_lic", "notice"},
				{"apacheLib.meta_lic", "apacheLib.meta_lic", "notice"},
			},
		},
		{
			name: "firstpartytool",
			roots: []string{"apacheBin.meta_lic"},
			edges: []annotated{
				{"apacheBin.meta_lic", "apacheLib.meta_lic", []string{"toolchain"}},
			},
			expectedResolutions: []res{
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheLib.meta_lic", "apacheLib.meta_lic", "notice"},
			},
		},
		{
			name: "firstpartydeep",
			roots: []string{"apacheContainer.meta_lic"},
			edges: []annotated{
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", []string{"static"}},
				{"apacheBin.meta_lic", "apacheLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheLib.meta_lic", "notice"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "apacheLib.meta_lic", "notice"},
				{"apacheLib.meta_lic", "apacheLib.meta_lic", "notice"},
			},
		},
		{
			name: "firstpartywide",
			roots: []string{"apacheContainer.meta_lic"},
			edges: []annotated{
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", []string{"static"}},
				{"apacheContainer.meta_lic", "apacheLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheLib.meta_lic", "notice"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheLib.meta_lic", "apacheLib.meta_lic", "notice"},
			},
		},
		{
			name: "firstpartydynamic",
			roots: []string{"apacheBin.meta_lic"},
			edges: []annotated{
				{"apacheBin.meta_lic", "apacheLib.meta_lic", []string{"dynamic"}},
			},
			expectedResolutions: []res{
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheLib.meta_lic", "apacheLib.meta_lic", "notice"},
			},
		},
		{
			name: "firstpartydynamicdeep",
			roots: []string{"apacheContainer.meta_lic"},
			edges: []annotated{
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", []string{"static"}},
				{"apacheBin.meta_lic", "apacheLib.meta_lic", []string{"dynamic"}},
			},
			expectedResolutions: []res{
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheLib.meta_lic", "apacheLib.meta_lic", "notice"},
			},
		},
		{
			name: "firstpartydynamicwide",
			roots: []string{"apacheContainer.meta_lic"},
			edges: []annotated{
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", []string{"static"}},
				{"apacheContainer.meta_lic", "apacheLib.meta_lic", []string{"dynamic"}},
			},
			expectedResolutions: []res{
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheLib.meta_lic", "apacheLib.meta_lic", "notice"},
			},
		},
		{
			name: "restricted",
			roots: []string{"apacheBin.meta_lic"},
			edges: []annotated{
				{"apacheBin.meta_lic", "gplLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "gplLib.meta_lic", "restricted"},
				{"gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
			},
		},
		{
			name: "restrictedtool",
			roots: []string{"apacheBin.meta_lic"},
			edges: []annotated{
				{"apacheBin.meta_lic", "gplLib.meta_lic", []string{"toolchain"}},
			},
			expectedResolutions: []res{
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
			},
		},
		{
			name: "restricteddeep",
			roots: []string{"apacheContainer.meta_lic"},
			edges: []annotated{
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", []string{"static"}},
				{"apacheBin.meta_lic", "gplLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "gplLib.meta_lic", "restricted"},
				{"gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
			},
		},
		{
			name: "restrictedwide",
			roots: []string{"apacheContainer.meta_lic"},
			edges: []annotated{
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", []string{"static"}},
				{"apacheContainer.meta_lic", "gplLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
			},
		},
		{
			name: "restricteddynamic",
			roots: []string{"apacheBin.meta_lic"},
			edges: []annotated{
				{"apacheBin.meta_lic", "gplLib.meta_lic", []string{"dynamic"}},
			},
			expectedResolutions: []res{
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "gplLib.meta_lic", "restricted"},
				{"gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
			},
		},
		{
			name: "restricteddynamicdeep",
			roots: []string{"apacheContainer.meta_lic"},
			edges: []annotated{
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", []string{"static"}},
				{"apacheBin.meta_lic", "gplLib.meta_lic", []string{"dynamic"}},
			},
			expectedResolutions: []res{
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "gplLib.meta_lic", "restricted"},
				{"gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
			},
		},
		{
			name: "restricteddynamicwide",
			roots: []string{"apacheContainer.meta_lic"},
			edges: []annotated{
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", []string{"static"}},
				{"apacheContainer.meta_lic", "gplLib.meta_lic", []string{"dynamic"}},
			},
			expectedResolutions: []res{
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
			},
		},
		{
			name: "weakrestricted",
			roots: []string{"apacheBin.meta_lic"},
			edges: []annotated{
				{"apacheBin.meta_lic", "lgplLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "lgplLib.meta_lic", "restricted"},
				{"lgplLib.meta_lic", "lgplLib.meta_lic", "restricted"},
			},
		},
		{
			name: "weakrestrictedtool",
			roots: []string{"apacheBin.meta_lic"},
			edges: []annotated{
				{"apacheBin.meta_lic", "lgplLib.meta_lic", []string{"toolchain"}},
			},
			expectedResolutions: []res{
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"lgplLib.meta_lic", "lgplLib.meta_lic", "restricted"},
			},
		},
		{
			name: "weakrestricteddeep",
			roots: []string{"apacheContainer.meta_lic"},
			edges: []annotated{
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", []string{"static"}},
				{"apacheBin.meta_lic", "lgplLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "lgplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "lgplLib.meta_lic", "restricted"},
				{"lgplLib.meta_lic", "lgplLib.meta_lic", "restricted"},
			},
		},
		{
			name: "weakrestrictedwide",
			roots: []string{"apacheContainer.meta_lic"},
			edges: []annotated{
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", []string{"static"}},
				{"apacheContainer.meta_lic", "lgplLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "lgplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"lgplLib.meta_lic", "lgplLib.meta_lic", "restricted"},
			},
		},
		{
			name: "weakrestricteddynamic",
			roots: []string{"apacheBin.meta_lic"},
			edges: []annotated{
				{"apacheBin.meta_lic", "lgplLib.meta_lic", []string{"dynamic"}},
			},
			expectedResolutions: []res{
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"lgplLib.meta_lic", "lgplLib.meta_lic", "restricted"},
			},
		},
		{
			name: "weakrestricteddynamicdeep",
			roots: []string{"apacheContainer.meta_lic"},
			edges: []annotated{
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", []string{"static"}},
				{"apacheBin.meta_lic", "lgplLib.meta_lic", []string{"dynamic"}},
			},
			expectedResolutions: []res{
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"lgplLib.meta_lic", "lgplLib.meta_lic", "restricted"},
			},
		},
		{
			name: "weakrestricteddynamicwide",
			roots: []string{"apacheContainer.meta_lic"},
			edges: []annotated{
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", []string{"static"}},
				{"apacheContainer.meta_lic", "lgplLib.meta_lic", []string{"dynamic"}},
			},
			expectedResolutions: []res{
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"lgplLib.meta_lic", "lgplLib.meta_lic", "restricted"},
			},
		},
		{
			name: "classpath",
			roots: []string{"apacheBin.meta_lic"},
			edges: []annotated{
				{"apacheBin.meta_lic", "gplWithClasspathException.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
				{"gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
			},
		},
		{
			name: "classpathdependent",
			roots: []string{"dependentModule.meta_lic"},
			edges: []annotated{
				{"dependentModule.meta_lic", "gplWithClasspathException.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"dependentModule.meta_lic", "dependentModule.meta_lic", "notice"},
				{"dependentModule.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
				{"gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
			},
		},
		{
			name: "classpathdynamic",
			roots: []string{"apacheBin.meta_lic"},
			edges: []annotated{
				{"apacheBin.meta_lic", "gplWithClasspathException.meta_lic", []string{"dynamic"}},
			},
			expectedResolutions: []res{
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
			},
		},
		{
			name: "classpathdependentdynamic",
			roots: []string{"dependentModule.meta_lic"},
			edges: []annotated{
				{"dependentModule.meta_lic", "gplWithClasspathException.meta_lic", []string{"dynamic"}},
			},
			expectedResolutions: []res{
				{"dependentModule.meta_lic", "dependentModule.meta_lic", "notice"},
				{"dependentModule.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
				{"gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stderr := &bytes.Buffer{}
			licenseGraph, err := toGraph(stderr, tt.roots, tt.edges)
			if err != nil {
				t.Errorf("unexpected test data error: got %w, want no error", err)
				return
			}
			lg := licenseGraph.(*licenseGraphImp)
			expectedRs := toResolutionSet(lg, tt.expectedResolutions)
			actualRs := ResolveBottomUpConditions(lg)
			checkSame(actualRs, expectedRs, t)
		})
	}
}

func TestResolveTopDownConditions(t *testing.T) {
	tests := []struct {
		name                string
		roots               []string
		edges               []annotated
		expectedResolutions []res
	}{
		{
			name: "firstparty",
			roots: []string{"apacheBin.meta_lic"},
			edges: []annotated{
				{"apacheBin.meta_lic", "apacheLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "apacheLib.meta_lic", "notice"},
				{"apacheLib.meta_lic", "apacheLib.meta_lic", "notice"},
			},
		},
		{
			name: "firstpartydynamic",
			roots: []string{"apacheBin.meta_lic"},
			edges: []annotated{
				{"apacheBin.meta_lic", "apacheLib.meta_lic", []string{"dynamic"}},
			},
			expectedResolutions: []res{
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheLib.meta_lic", "apacheLib.meta_lic", "notice"},
			},
		},
		{
			name: "restricted",
			roots: []string{"apacheBin.meta_lic"},
			edges: []annotated{
				{"apacheBin.meta_lic", "gplLib.meta_lic", []string{"static"}},
				{"apacheBin.meta_lic", "mitLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "mitLib.meta_lic", "notice"},
				{"gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
				{"mitLib.meta_lic", "mitLib.meta_lic", "notice"},
				{"mitLib.meta_lic", "gplLib.meta_lic", "restricted"},
			},
		},
		{
			name: "restrictedtool",
			roots: []string{"apacheBin.meta_lic"},
			edges: []annotated{
				{"apacheBin.meta_lic", "gplBin.meta_lic", []string{"toolchain"}},
				{"apacheBin.meta_lic", "mitLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "mitLib.meta_lic", "notice"},
				{"gplBin.meta_lic", "gplBin.meta_lic", "restricted"},
				{"mitLib.meta_lic", "mitLib.meta_lic", "notice"},
			},
		},
		{
			name: "restricteddeep",
			roots: []string{"apacheContainer.meta_lic"},
			edges: []annotated{
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", []string{"static"}},
				{"apacheContainer.meta_lic", "mitBin.meta_lic", []string{"static"}},
				{"apacheBin.meta_lic", "gplLib.meta_lic", []string{"static"}},
				{"apacheBin.meta_lic", "mplLib.meta_lic", []string{"static"}},
				{"mitBin.meta_lic", "mitLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheContainer.meta_lic", "mitBin.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "mitLib.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "mplLib.meta_lic", "reciprocal"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "mplLib.meta_lic", "reciprocal"},
				{"gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
				{"mitBin.meta_lic", "mitBin.meta_lic", "notice"},
				{"mitBin.meta_lic", "mitLib.meta_lic", "notice"},
				{"mitLib.meta_lic", "mitLib.meta_lic", "notice"},
				{"mplLib.meta_lic", "mplLib.meta_lic", "reciprocal"},
				{"mplLib.meta_lic", "gplLib.meta_lic", "restricted"},
			},
		},
		{
			name: "restrictedwide",
			roots: []string{"apacheContainer.meta_lic"},
			edges: []annotated{
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", []string{"static"}},
				{"apacheContainer.meta_lic", "gplLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
			},
		},
		{
			name: "restricteddynamic",
			roots: []string{"apacheBin.meta_lic"},
			edges: []annotated{
				{"apacheBin.meta_lic", "gplLib.meta_lic", []string{"dynamic"}},
				{"apacheBin.meta_lic", "mitLib.meta_lic", []string{"dynamic"}},
			},
			expectedResolutions: []res{
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "gplLib.meta_lic", "restricted"},
				{"gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
				{"mitLib.meta_lic", "gplLib.meta_lic", "restricted"},
				{"mitLib.meta_lic", "mitLib.meta_lic", "notice"},
			},
		},
		{
			name: "restricteddynamicdeep",
			roots: []string{"apacheContainer.meta_lic"},
			edges: []annotated{
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", []string{"static"}},
				{"apacheContainer.meta_lic", "mitBin.meta_lic", []string{"static"}},
				{"apacheBin.meta_lic", "gplLib.meta_lic", []string{"dynamic"}},
				{"apacheBin.meta_lic", "mplLib.meta_lic", []string{"dynamic"}},
				{"mitBin.meta_lic", "mitLib.meta_lic", []string{"dynamic"}},
			},
			expectedResolutions: []res{
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheContainer.meta_lic", "mitBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "gplLib.meta_lic", "restricted"},
				{"gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
				{"mitBin.meta_lic", "mitBin.meta_lic", "notice"},
				{"mitLib.meta_lic", "mitLib.meta_lic", "notice"},
				{"mplLib.meta_lic", "mplLib.meta_lic", "reciprocal"},
				{"mplLib.meta_lic", "gplLib.meta_lic", "restricted"},
			},
		},
		{
			name: "restricteddynamicwide",
			roots: []string{"apacheContainer.meta_lic"},
			edges: []annotated{
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", []string{"static"}},
				{"apacheContainer.meta_lic", "gplLib.meta_lic", []string{"dynamic"}},
			},
			expectedResolutions: []res{
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
			},
		},
		{
			name: "weakrestricted",
			roots: []string{"apacheBin.meta_lic"},
			edges: []annotated{
				{"apacheBin.meta_lic", "lgplLib.meta_lic", []string{"static"}},
				{"apacheBin.meta_lic", "mitLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "lgplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "mitLib.meta_lic", "notice"},
				{"lgplLib.meta_lic", "lgplLib.meta_lic", "restricted"},
				{"mitLib.meta_lic", "lgplLib.meta_lic", "restricted"},
				{"mitLib.meta_lic", "mitLib.meta_lic", "notice"},
			},
		},
		{
			name: "weakrestrictedtool",
			roots: []string{"apacheBin.meta_lic"},
			edges: []annotated{
				{"apacheBin.meta_lic", "lgplBin.meta_lic", []string{"toolchain"}},
				{"apacheBin.meta_lic", "mitLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "mitLib.meta_lic", "notice"},
				{"lgplBin.meta_lic", "lgplBin.meta_lic", "restricted"},
				{"mitLib.meta_lic", "mitLib.meta_lic", "notice"},
			},
		},
		{
			name: "weakrestricteddeep",
			roots: []string{"apacheContainer.meta_lic"},
			edges: []annotated{
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", []string{"static"}},
				{"apacheBin.meta_lic", "lgplLib.meta_lic", []string{"static"}},
				{"apacheBin.meta_lic", "mitLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "lgplLib.meta_lic", "restricted"},
				{"apacheContainer.meta_lic", "mitLib.meta_lic", "notice"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "lgplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "mitLib.meta_lic", "notice"},
				{"lgplLib.meta_lic", "lgplLib.meta_lic", "restricted"},
				{"mitLib.meta_lic", "mitLib.meta_lic", "notice"},
				{"mitLib.meta_lic", "lgplLib.meta_lic", "restricted"},
			},
		},
		{
			name: "weakrestrictedwide",
			roots: []string{"apacheContainer.meta_lic"},
			edges: []annotated{
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", []string{"static"}},
				{"apacheContainer.meta_lic", "lgplLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "lgplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"lgplLib.meta_lic", "lgplLib.meta_lic", "restricted"},
			},
		},
		{
			name: "weakrestricteddynamic",
			roots: []string{"apacheBin.meta_lic"},
			edges: []annotated{
				{"apacheBin.meta_lic", "lgplLib.meta_lic", []string{"dynamic"}},
				{"apacheBin.meta_lic", "mitLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "mitLib.meta_lic", "notice"},
				{"lgplLib.meta_lic", "lgplLib.meta_lic", "restricted"},
				{"mitLib.meta_lic", "mitLib.meta_lic", "notice"},
			},
		},
		{
			name: "weakrestricteddynamicdeep",
			roots: []string{"apacheContainer.meta_lic"},
			edges: []annotated{
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", []string{"static"}},
				{"apacheBin.meta_lic", "lgplLib.meta_lic", []string{"dynamic"}},
			},
			expectedResolutions: []res{
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"lgplLib.meta_lic", "lgplLib.meta_lic", "restricted"},
			},
		},
		{
			name: "weakrestricteddynamicwide",
			roots: []string{"apacheContainer.meta_lic"},
			edges: []annotated{
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", []string{"static"}},
				{"apacheContainer.meta_lic", "lgplLib.meta_lic", []string{"dynamic"}},
			},
			expectedResolutions: []res{
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"lgplLib.meta_lic", "lgplLib.meta_lic", "restricted"},
			},
		},
		{
			name: "classpath",
			roots: []string{"apacheBin.meta_lic"},
			edges: []annotated{
				{"apacheBin.meta_lic", "gplWithClasspathException.meta_lic", []string{"static"}},
				{"apacheBin.meta_lic", "mitLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "mitLib.meta_lic", "notice"},
				{"gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
				{"mitLib.meta_lic", "mitLib.meta_lic", "notice"},
				{"mitLib.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
			},
		},
		{
			name: "classpathdependent",
			roots: []string{"dependentModule.meta_lic"},
			edges: []annotated{
				{"dependentModule.meta_lic", "gplWithClasspathException.meta_lic", []string{"static"}},
				{"dependentModule.meta_lic", "mitLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"dependentModule.meta_lic", "dependentModule.meta_lic", "notice"},
				{"dependentModule.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
				{"dependentModule.meta_lic", "mitLib.meta_lic", "notice"},
				{"gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
				{"mitLib.meta_lic", "mitLib.meta_lic", "notice"},
				{"mitLib.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
			},
		},
		{
			name: "classpathdynamic",
			roots: []string{"apacheBin.meta_lic"},
			edges: []annotated{
				{"apacheBin.meta_lic", "gplWithClasspathException.meta_lic", []string{"dynamic"}},
				{"apacheBin.meta_lic", "mitLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "mitLib.meta_lic", "notice"},
				{"gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
				{"mitLib.meta_lic", "mitLib.meta_lic", "notice"},
			},
		},
		{
			name: "classpathdependentdynamic",
			roots: []string{"dependentModule.meta_lic"},
			edges: []annotated{
				{"dependentModule.meta_lic", "gplWithClasspathException.meta_lic", []string{"dynamic"}},
				{"dependentModule.meta_lic", "mitLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"dependentModule.meta_lic", "dependentModule.meta_lic", "notice"},
				{"dependentModule.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
				{"dependentModule.meta_lic", "mitLib.meta_lic", "notice"},
				{"gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
				{"mitLib.meta_lic", "mitLib.meta_lic", "notice"},
				{"mitLib.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stderr := &bytes.Buffer{}
			licenseGraph, err := toGraph(stderr, tt.roots, tt.edges)
			if err != nil {
				t.Errorf("unexpected test data error: got %w, want no error", err)
				return
			}
			lg := licenseGraph.(*licenseGraphImp)
			expectedRs := toResolutionSet(lg, tt.expectedResolutions)
			actualRs := ResolveTopDownConditions(lg)
			checkSame(actualRs, expectedRs, t)
		})
	}
}

func TestResolveTopDownForCondition(t *testing.T) {
	tests := []struct {
		name                string
		condition           string
		roots               []string
		edges               []annotated
		expectedResolutions []res
	}{
		{
			name: "firstparty",
			condition: "notice",
			roots: []string{"apacheBin.meta_lic"},
			edges: []annotated{
				{"apacheBin.meta_lic", "apacheLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "apacheLib.meta_lic", "notice"},
				{"apacheLib.meta_lic", "apacheLib.meta_lic", "notice"},
			},
		},
		{
			name: "notice",
			condition: "notice",
			roots: []string{"mitBin.meta_lic"},
			edges: []annotated{
				{"mitBin.meta_lic", "mitLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"mitBin.meta_lic", "mitBin.meta_lic", "notice"},
				{"mitBin.meta_lic", "mitLib.meta_lic", "notice"},
				{"mitLib.meta_lic", "mitLib.meta_lic", "notice"},
			},
		},
		{
			name: "fponlgplnotice",
			condition: "notice",
			roots: []string{"apacheBin.meta_lic"},
			edges: []annotated{
				{"apacheBin.meta_lic", "lgplLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "lgplLib.meta_lic", "notice"},
				{"lgplLib.meta_lic", "lgplLib.meta_lic", "notice"},
			},
		},
		{
			name: "fponlgpldynamicnotice",
			condition: "notice",
			roots: []string{"apacheBin.meta_lic"},
			edges: []annotated{
				{"apacheBin.meta_lic", "lgplLib.meta_lic", []string{"dynamic"}},
			},
			expectedResolutions: []res{
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"lgplLib.meta_lic", "lgplLib.meta_lic", "notice"},
			},
		},
		{
			name: "independentmodulenotice",
			condition: "notice",
			roots: []string{"apacheBin.meta_lic"},
			edges: []annotated{
				{"apacheBin.meta_lic", "gplWithClasspathException.meta_lic", []string{"dynamic"}},
			},
			expectedResolutions: []res{
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "notice"},
			},
		},
		{
			name: "independentmodulerestricted",
			condition: "restricted",
			roots: []string{"apacheBin.meta_lic"},
			edges: []annotated{
				{"apacheBin.meta_lic", "gplWithClasspathException.meta_lic", []string{"dynamic"}},
			},
			expectedResolutions: []res{},
		},
		{
			name: "independentmodulestaticnotice",
			condition: "notice",
			roots: []string{"apacheBin.meta_lic"},
			edges: []annotated{
				{"apacheBin.meta_lic", "gplWithClasspathException.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "gplWithClasspathException.meta_lic", "notice"},
				{"gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "notice"},
			},
		},
		{
			name: "independentmodulestaticrestricted",
			condition: "restricted",
			roots: []string{"apacheBin.meta_lic"},
			edges: []annotated{
				{"apacheBin.meta_lic", "gplWithClasspathException.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"apacheBin.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
				{"gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
			},
		},
		{
			name: "dependentmodulenotice",
			condition: "notice",
			roots: []string{"dependentModule.meta_lic"},
			edges: []annotated{
				{"dependentModule.meta_lic", "gplWithClasspathException.meta_lic", []string{"dynamic"}},
			},
			expectedResolutions: []res{
				{"dependentModule.meta_lic", "dependentModule.meta_lic", "notice"},
				{"dependentModule.meta_lic", "gplWithClasspathException.meta_lic", "notice"},
				{"gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "notice"},
			},
		},
		{
			name: "dependentmodulerestricted",
			condition: "restricted",
			roots: []string{"dependentModule.meta_lic"},
			edges: []annotated{
				{"dependentModule.meta_lic", "gplWithClasspathException.meta_lic", []string{"dynamic"}},
			},
			expectedResolutions: []res{
				{"dependentModule.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
				{"gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
			},
		},
		{
			name: "lgplonfpnotice",
			condition: "notice",
			roots: []string{"lgplBin.meta_lic"},
			edges: []annotated{
				{"lgplBin.meta_lic", "apacheLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"lgplBin.meta_lic", "lgplBin.meta_lic", "notice"},
				{"lgplBin.meta_lic", "apacheLib.meta_lic", "notice"},
				{"apacheLib.meta_lic", "apacheLib.meta_lic", "notice"},
			},
		},
		{
			name: "lgplonfprestricted",
			condition: "restricted",
			roots: []string{"lgplBin.meta_lic"},
			edges: []annotated{
				{"lgplBin.meta_lic", "apacheLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"lgplBin.meta_lic", "lgplBin.meta_lic", "restricted"},
				{"apacheLib.meta_lic", "lgplBin.meta_lic", "restricted"},
			},
		},
		{
			name: "lgplonfpdynamicnotice",
			condition: "notice",
			roots: []string{"lgplBin.meta_lic"},
			edges: []annotated{
				{"lgplBin.meta_lic", "apacheLib.meta_lic", []string{"dynamic"}},
			},
			expectedResolutions: []res{
				{"lgplBin.meta_lic", "lgplBin.meta_lic", "notice"},
				{"apacheLib.meta_lic", "apacheLib.meta_lic", "notice"},
			},
		},
		{
			name: "lgplonfpdynamicrestricted",
			condition: "restricted",
			roots: []string{"lgplBin.meta_lic"},
			edges: []annotated{
				{"lgplBin.meta_lic", "apacheLib.meta_lic", []string{"dynamic"}},
			},
			expectedResolutions: []res{
				{"lgplBin.meta_lic", "lgplBin.meta_lic", "restricted"},
			},
		},
		{
			name: "gplonfpnotice",
			condition: "notice",
			roots: []string{"gplBin.meta_lic"},
			edges: []annotated{
				{"gplBin.meta_lic", "apacheLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"gplBin.meta_lic", "gplBin.meta_lic", "notice"},
				{"gplBin.meta_lic", "apacheLib.meta_lic", "notice"},
				{"apacheLib.meta_lic", "apacheLib.meta_lic", "notice"},
			},
		},
		{
			name: "gplonfprestricted",
			condition: "restricted",
			roots: []string{"gplBin.meta_lic"},
			edges: []annotated{
				{"gplBin.meta_lic", "apacheLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"gplBin.meta_lic", "gplBin.meta_lic", "restricted"},
				{"apacheLib.meta_lic", "gplBin.meta_lic", "restricted"},
			},
		},
		{
			name: "gplcontainernotice",
			condition: "notice",
			roots: []string{"gplContainer.meta_lic"},
			edges: []annotated{
				{"gplContainer.meta_lic", "apacheLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"gplContainer.meta_lic", "gplContainer.meta_lic", "notice"},
				{"gplContainer.meta_lic", "apacheLib.meta_lic", "notice"},
				{"apacheLib.meta_lic", "apacheLib.meta_lic", "notice"},
			},
		},
		{
			name: "gplcontainerrestricted",
			condition: "restricted",
			roots: []string{"gplContainer.meta_lic"},
			edges: []annotated{
				{"gplContainer.meta_lic", "apacheLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"gplContainer.meta_lic", "gplContainer.meta_lic", "restricted"},
				{"apacheLib.meta_lic", "gplContainer.meta_lic", "restricted"},
			},
		},
		{
			name: "gploncontainernotice",
			condition: "notice",
			roots: []string{"apacheContainer.meta_lic"},
			edges: []annotated{
				{"apacheContainer.meta_lic", "apacheLib.meta_lic", []string{"static"}},
				{"apacheContainer.meta_lic", "gplLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheLib.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "gplLib.meta_lic", "notice"},
				{"apacheLib.meta_lic", "apacheLib.meta_lic", "notice"},
				{"gplLib.meta_lic", "gplLib.meta_lic", "notice"},
			},
		},
		{
			name: "gploncontainerrestricted",
			condition: "restricted",
			roots: []string{"apacheContainer.meta_lic"},
			edges: []annotated{
				{"apacheContainer.meta_lic", "apacheLib.meta_lic", []string{"static"}},
				{"apacheContainer.meta_lic", "gplLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"apacheContainer.meta_lic", "gplLib.meta_lic", "restricted"},
				{"gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
			},
		},
		{
			name: "gplonbinnotice",
			condition: "notice",
			roots: []string{"apacheBin.meta_lic"},
			edges: []annotated{
				{"apacheBin.meta_lic", "apacheLib.meta_lic", []string{"static"}},
				{"apacheBin.meta_lic", "gplLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "apacheLib.meta_lic", "notice"},
				{"apacheBin.meta_lic", "gplLib.meta_lic", "notice"},
				{"apacheLib.meta_lic", "apacheLib.meta_lic", "notice"},
				{"gplLib.meta_lic", "gplLib.meta_lic", "notice"},
			},
		},
		{
			name: "gplonbinrestricted",
			condition: "restricted",
			roots: []string{"apacheBin.meta_lic"},
			edges: []annotated{
				{"apacheBin.meta_lic", "apacheLib.meta_lic", []string{"static"}},
				{"apacheBin.meta_lic", "gplLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"apacheBin.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheLib.meta_lic", "gplLib.meta_lic", "restricted"},
				{"gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
			},
		},
		{
			name: "gplonfpdynamicnotice",
			condition: "notice",
			roots: []string{"gplBin.meta_lic"},
			edges: []annotated{
				{"gplBin.meta_lic", "apacheLib.meta_lic", []string{"dynamic"}},
			},
			expectedResolutions: []res{
				{"gplBin.meta_lic", "gplBin.meta_lic", "notice"},
				{"apacheLib.meta_lic", "apacheLib.meta_lic", "notice"},
			},
		},
		{
			name: "gplonfpdynamicrestricted",
			condition: "restricted",
			roots: []string{"gplBin.meta_lic"},
			edges: []annotated{
				{"gplBin.meta_lic", "apacheLib.meta_lic", []string{"dynamic"}},
			},
			expectedResolutions: []res{
				{"gplBin.meta_lic", "gplBin.meta_lic", "restricted"},
				{"apacheLib.meta_lic", "gplBin.meta_lic", "restricted"},
			},
		},
		{
			name: "independentmodulereversenotice",
			condition: "notice",
			roots: []string{"gplWithClasspathException.meta_lic"},
			edges: []annotated{
				{"gplWithClasspathException.meta_lic", "apacheBin.meta_lic", []string{"dynamic"}},
			},
			expectedResolutions: []res{
				{"gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "notice"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
			},
		},
		{
			name: "independentmodulereverserestricted",
			condition: "restricted",
			roots: []string{"gplWithClasspathException.meta_lic"},
			edges: []annotated{
				{"gplWithClasspathException.meta_lic", "apacheBin.meta_lic", []string{"dynamic"}},
			},
			expectedResolutions: []res{
				{"gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
			},
		},
		{
			name: "independentmodulereversestaticnotice",
			condition: "notice",
			roots: []string{"gplWithClasspathException.meta_lic"},
			edges: []annotated{
				{"gplWithClasspathException.meta_lic", "apacheBin.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "notice"},
				{"gplWithClasspathException.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
			},
		},
		{
			name: "independentmodulereversestaticrestricted",
			condition: "restricted",
			roots: []string{"gplWithClasspathException.meta_lic"},
			edges: []annotated{
				{"gplWithClasspathException.meta_lic", "apacheBin.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
			},
		},
		{
			name: "dependentmodulereversenotice",
			condition: "notice",
			roots: []string{"gplWithClasspathException.meta_lic"},
			edges: []annotated{
				{"gplWithClasspathException.meta_lic", "dependentModule.meta_lic", []string{"dynamic"}},
			},
			expectedResolutions: []res{
				{"gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "notice"},
				{"dependentModule.meta_lic", "dependentModule.meta_lic", "notice"},
			},
		},
		{
			name: "dependentmodulereverserestricted",
			condition: "restricted",
			roots: []string{"gplWithClasspathException.meta_lic"},
			edges: []annotated{
				{"gplWithClasspathException.meta_lic", "dependentModule.meta_lic", []string{"dynamic"}},
			},
			expectedResolutions: []res{
				{"gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
				{"dependentModule.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
			},
		},
		{
			name: "ponrnotice",
			condition: "notice",
			roots: []string{"proprietary.meta_lic"},
			edges: []annotated{
				{"proprietary.meta_lic", "gplLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"proprietary.meta_lic", "proprietary.meta_lic", "notice"},
				{"proprietary.meta_lic", "gplLib.meta_lic", "notice"},
				{"gplLib.meta_lic", "gplLib.meta_lic", "notice"},
			},
		},
		{
			name: "ponrrestricted",
			condition: "restricted",
			roots: []string{"proprietary.meta_lic"},
			edges: []annotated{
				{"proprietary.meta_lic", "gplLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"proprietary.meta_lic", "gplLib.meta_lic", "restricted"},
				{"gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
			},
		},
		{
			name: "ponrproprietary",
			condition: "proprietary",
			roots: []string{"proprietary.meta_lic"},
			edges: []annotated{
				{"proprietary.meta_lic", "gplLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"proprietary.meta_lic", "proprietary.meta_lic", "proprietary"},
			},
		},
		{
			name: "ronpnotice",
			condition: "notice",
			roots: []string{"gplBin.meta_lic"},
			edges: []annotated{
				{"gplBin.meta_lic", "proprietary.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"gplBin.meta_lic", "gplBin.meta_lic", "notice"},
				{"gplBin.meta_lic", "proprietary.meta_lic", "notice"},
				{"proprietary.meta_lic", "proprietary.meta_lic", "notice"},
				{"proprietary.meta_lic", "gplBin.meta_lic", "notice"},
			},
		},
		{
			name: "ronprestricted",
			condition: "restricted",
			roots: []string{"gplBin.meta_lic"},
			edges: []annotated{
				{"gplBin.meta_lic", "proprietary.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"gplBin.meta_lic", "gplBin.meta_lic", "restricted"},
				{"proprietary.meta_lic", "gplBin.meta_lic", "restricted"},
			},
		},
		{
			name: "ronpproprietary",
			condition: "proprietary",
			roots: []string{"gplBin.meta_lic"},
			edges: []annotated{
				{"gplBin.meta_lic", "proprietary.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"gplBin.meta_lic", "proprietary.meta_lic", "proprietary"},
				{"proprietary.meta_lic", "proprietary.meta_lic", "proprietary"},
			},
		},
		{
			name: "noticeonb_e_onotice",
			condition: "notice",
			roots: []string{"mitBin.meta_lic"},
			edges: []annotated{
				{"mitBin.meta_lic", "by_exception.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"mitBin.meta_lic", "mitBin.meta_lic", "notice"},
				{"mitBin.meta_lic", "by_exception.meta_lic", "notice"},
				{"by_exception.meta_lic", "by_exception.meta_lic", "notice"},
			},
		},
		{
			name: "noticeonb_e_orestricted",
			condition: "restricted",
			roots: []string{"mitBin.meta_lic"},
			edges: []annotated{
				{"mitBin.meta_lic", "by_exception.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{},
		},
		{
			name: "noticeonb_e_ob_e_o",
			condition: "by_exception_only",
			roots: []string{"mitBin.meta_lic"},
			edges: []annotated{
				{"mitBin.meta_lic", "by_exception.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"mitBin.meta_lic", "by_exception.meta_lic", "by_exception_only"},
				{"by_exception.meta_lic", "by_exception.meta_lic", "by_exception_only"},
			},
		},
		{
			name: "b_e_oonnoticenotice",
			condition: "notice",
			roots: []string{"by_exception.meta_lic"},
			edges: []annotated{
				{"by_exception.meta_lic", "mitLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"by_exception.meta_lic", "by_exception.meta_lic", "notice"},
				{"by_exception.meta_lic", "mitLib.meta_lic", "notice"},
				{"mitLib.meta_lic", "mitLib.meta_lic", "notice"},
			},
		},
		{
			name: "b_e_oonnoticerestricted",
			condition: "restricted",
			roots: []string{"by_exception.meta_lic"},
			edges: []annotated{
				{"by_exception.meta_lic", "mitLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{},
		},
		{
			name: "b_e_oonnoticeb_e_o",
			condition: "by_exception_only",
			roots: []string{"by_exception.meta_lic"},
			edges: []annotated{
				{"by_exception.meta_lic", "mitLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"by_exception.meta_lic", "by_exception.meta_lic", "by_exception_only"},
			},
		},
		{
			name: "noticeonrecipnotice",
			condition: "notice",
			roots: []string{"mitBin.meta_lic"},
			edges: []annotated{
				{"mitBin.meta_lic", "mplLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"mitBin.meta_lic", "mitBin.meta_lic", "notice"},
				{"mitBin.meta_lic", "mplLib.meta_lic", "notice"},
				{"mplLib.meta_lic", "mplLib.meta_lic", "notice"},
			},
		},
		{
			name: "noticeonreciprecip",
			condition: "reciprocal",
			roots: []string{"mitBin.meta_lic"},
			edges: []annotated{
				{"mitBin.meta_lic", "mplLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"mitBin.meta_lic", "mplLib.meta_lic", "reciprocal"},
				{"mplLib.meta_lic", "mplLib.meta_lic", "reciprocal"},
			},
		},
		{
			name: "reciponnoticenotice",
			condition: "notice",
			roots: []string{"mplBin.meta_lic"},
			edges: []annotated{
				{"mplBin.meta_lic", "mitLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"mplBin.meta_lic", "mplBin.meta_lic", "notice"},
				{"mplBin.meta_lic", "mitLib.meta_lic", "notice"},
				{"mitLib.meta_lic", "mitLib.meta_lic", "notice"},
			},
		},
		{
			name: "reciponnoticerecip",
			condition: "reciprocal",
			roots: []string{"mplBin.meta_lic"},
			edges: []annotated{
				{"mplBin.meta_lic", "mitLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"mplBin.meta_lic", "mplBin.meta_lic", "reciprocal"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
//			stderr := &bytes.Buffer{}
//			licenseGraph, err := toGraph(stderr, tt.roots, tt.edges)
//			if err != nil {
//				t.Errorf("unexpected test data error: got %w, want no error", err)
//				return
//			}
//			lg := licenseGraph.(*licenseGraphImp)
//			expectedRs := toResolutionSet(lg, tt.expectedResolutions)
//			actualRs := ResolveTopDownForCondition(lg, tt.condition)
//			checkSame(actualRs, expectedRs, t)
		})
	}
}
