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
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "apacheLib.meta_lic", "apacheLib.meta_lic", "notice"},
				{"apacheLib.meta_lic", "apacheLib.meta_lic", "apacheLib.meta_lic", "notice"},
			},
		},
		{
			name: "firstpartytool",
			roots: []string{"apacheBin.meta_lic"},
			edges: []annotated{
				{"apacheBin.meta_lic", "apacheLib.meta_lic", []string{"toolchain"}},
			},
			expectedResolutions: []res{
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheLib.meta_lic", "apacheLib.meta_lic", "apacheLib.meta_lic", "notice"},
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
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "apacheContainer.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheLib.meta_lic", "apacheLib.meta_lic", "notice"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "apacheLib.meta_lic", "apacheLib.meta_lic", "notice"},
				{"apacheLib.meta_lic", "apacheLib.meta_lic", "apacheLib.meta_lic", "notice"},
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
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "apacheContainer.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheLib.meta_lic", "apacheLib.meta_lic", "notice"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheLib.meta_lic", "apacheLib.meta_lic", "apacheLib.meta_lic", "notice"},
			},
		},
		{
			name: "firstpartydynamic",
			roots: []string{"apacheBin.meta_lic"},
			edges: []annotated{
				{"apacheBin.meta_lic", "apacheLib.meta_lic", []string{"dynamic"}},
			},
			expectedResolutions: []res{
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheLib.meta_lic", "apacheLib.meta_lic", "apacheLib.meta_lic", "notice"},
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
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "apacheContainer.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheLib.meta_lic", "apacheLib.meta_lic", "apacheLib.meta_lic", "notice"},
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
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "apacheContainer.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheLib.meta_lic", "apacheLib.meta_lic", "apacheLib.meta_lic", "notice"},
			},
		},
		{
			name: "restricted",
			roots: []string{"apacheBin.meta_lic"},
			edges: []annotated{
				{"apacheBin.meta_lic", "gplLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
				{"gplLib.meta_lic", "gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
			},
		},
		{
			name: "restrictedtool",
			roots: []string{"apacheBin.meta_lic"},
			edges: []annotated{
				{"apacheBin.meta_lic", "gplLib.meta_lic", []string{"toolchain"}},
			},
			expectedResolutions: []res{
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"gplLib.meta_lic", "gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
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
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "apacheContainer.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheContainer.meta_lic", "gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
				{"gplLib.meta_lic", "gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
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
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "apacheContainer.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"gplLib.meta_lic", "gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
			},
		},
		{
			name: "restricteddynamic",
			roots: []string{"apacheBin.meta_lic"},
			edges: []annotated{
				{"apacheBin.meta_lic", "gplLib.meta_lic", []string{"dynamic"}},
			},
			expectedResolutions: []res{
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
				{"gplLib.meta_lic", "gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
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
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "apacheContainer.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheContainer.meta_lic", "gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
				{"gplLib.meta_lic", "gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
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
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "apacheContainer.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"gplLib.meta_lic", "gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
			},
		},
		{
			name: "weakrestricted",
			roots: []string{"apacheBin.meta_lic"},
			edges: []annotated{
				{"apacheBin.meta_lic", "lgplLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "lgplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "lgplLib.meta_lic", "lgplLib.meta_lic", "restricted"},
				{"lgplLib.meta_lic", "lgplLib.meta_lic", "lgplLib.meta_lic", "restricted"},
			},
		},
		{
			name: "weakrestrictedtool",
			roots: []string{"apacheBin.meta_lic"},
			edges: []annotated{
				{"apacheBin.meta_lic", "lgplLib.meta_lic", []string{"toolchain"}},
			},
			expectedResolutions: []res{
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"lgplLib.meta_lic", "lgplLib.meta_lic", "lgplLib.meta_lic", "restricted"},
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
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "apacheContainer.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "lgplLib.meta_lic", "restricted"},
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", "lgplLib.meta_lic", "restricted"},
				{"apacheContainer.meta_lic", "lgplLib.meta_lic", "lgplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "lgplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "lgplLib.meta_lic", "lgplLib.meta_lic", "restricted"},
				{"lgplLib.meta_lic", "lgplLib.meta_lic", "lgplLib.meta_lic", "restricted"},
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
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "apacheContainer.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "lgplLib.meta_lic", "restricted"},
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "lgplLib.meta_lic", "lgplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"lgplLib.meta_lic", "lgplLib.meta_lic", "lgplLib.meta_lic", "restricted"},
			},
		},
		{
			name: "weakrestricteddynamic",
			roots: []string{"apacheBin.meta_lic"},
			edges: []annotated{
				{"apacheBin.meta_lic", "lgplLib.meta_lic", []string{"dynamic"}},
			},
			expectedResolutions: []res{
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"lgplLib.meta_lic", "lgplLib.meta_lic", "lgplLib.meta_lic", "restricted"},
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
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "apacheContainer.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"lgplLib.meta_lic", "lgplLib.meta_lic", "lgplLib.meta_lic", "restricted"},
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
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "apacheContainer.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"lgplLib.meta_lic", "lgplLib.meta_lic", "lgplLib.meta_lic", "restricted"},
			},
		},
		{
			name: "classpath",
			roots: []string{"apacheBin.meta_lic"},
			edges: []annotated{
				{"apacheBin.meta_lic", "gplWithClasspathException.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
				{"gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
			},
		},
		{
			name: "classpathdependent",
			roots: []string{"dependentModule.meta_lic"},
			edges: []annotated{
				{"dependentModule.meta_lic", "gplWithClasspathException.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"dependentModule.meta_lic", "dependentModule.meta_lic", "dependentModule.meta_lic", "notice"},
				{"dependentModule.meta_lic", "dependentModule.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
				{"dependentModule.meta_lic", "gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
				{"gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
			},
		},
		{
			name: "classpathdynamic",
			roots: []string{"apacheBin.meta_lic"},
			edges: []annotated{
				{"apacheBin.meta_lic", "gplWithClasspathException.meta_lic", []string{"dynamic"}},
			},
			expectedResolutions: []res{
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
			},
		},
		{
			name: "classpathdependentdynamic",
			roots: []string{"dependentModule.meta_lic"},
			edges: []annotated{
				{"dependentModule.meta_lic", "gplWithClasspathException.meta_lic", []string{"dynamic"}},
			},
			expectedResolutions: []res{
				{"dependentModule.meta_lic", "dependentModule.meta_lic", "dependentModule.meta_lic", "notice"},
				{"dependentModule.meta_lic", "dependentModule.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
				{"dependentModule.meta_lic", "gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
				{"gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stderr := &bytes.Buffer{}
			lg, err := toGraph(stderr, tt.roots, tt.edges)
			if err != nil {
				t.Errorf("unexpected test data error: got %w, want no error", err)
				return
			}
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
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "apacheLib.meta_lic", "apacheLib.meta_lic", "notice"},
				{"apacheLib.meta_lic", "apacheLib.meta_lic", "apacheLib.meta_lic", "notice"},
			},
		},
		{
			name: "firstpartydynamic",
			roots: []string{"apacheBin.meta_lic"},
			edges: []annotated{
				{"apacheBin.meta_lic", "apacheLib.meta_lic", []string{"dynamic"}},
			},
			expectedResolutions: []res{
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheLib.meta_lic", "apacheLib.meta_lic", "apacheLib.meta_lic", "notice"},
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
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "mitLib.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "mitLib.meta_lic", "mitLib.meta_lic", "notice"},
				{"gplLib.meta_lic", "gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
				{"mitLib.meta_lic", "mitLib.meta_lic", "mitLib.meta_lic", "notice"},
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
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "mitLib.meta_lic", "mitLib.meta_lic", "notice"},
				{"gplBin.meta_lic", "gplBin.meta_lic", "gplBin.meta_lic", "restricted"},
				{"mitLib.meta_lic", "mitLib.meta_lic", "mitLib.meta_lic", "notice"},
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
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "apacheContainer.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheContainer.meta_lic", "gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheContainer.meta_lic", "mitBin.meta_lic", "mitBin.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "mitLib.meta_lic", "mitLib.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "mplLib.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheContainer.meta_lic", "mplLib.meta_lic", "mplLib.meta_lic", "reciprocal"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "mplLib.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "mplLib.meta_lic", "mplLib.meta_lic", "reciprocal"},
				{"gplLib.meta_lic", "gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
				{"mitBin.meta_lic", "mitBin.meta_lic", "mitBin.meta_lic", "notice"},
				{"mitBin.meta_lic", "mitLib.meta_lic", "mitLib.meta_lic", "notice"},
				{"mitLib.meta_lic", "mitLib.meta_lic", "mitLib.meta_lic", "notice"},
				{"mplLib.meta_lic", "mplLib.meta_lic", "mplLib.meta_lic", "reciprocal"},
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
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "apacheContainer.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"gplLib.meta_lic", "gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
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
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "mitLib.meta_lic", "gplLib.meta_lic", "restricted"},
				{"gplLib.meta_lic", "gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
				{"mitLib.meta_lic", "mitLib.meta_lic", "mitLib.meta_lic", "notice"},
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
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "apacheContainer.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheContainer.meta_lic", "gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheContainer.meta_lic", "mitBin.meta_lic", "mitBin.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "mplLib.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "mplLib.meta_lic", "gplLib.meta_lic", "restricted"},
				{"gplLib.meta_lic", "gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
				{"mitBin.meta_lic", "mitBin.meta_lic", "mitBin.meta_lic", "notice"},
				{"mitLib.meta_lic", "mitLib.meta_lic", "mitLib.meta_lic", "notice"},
				{"mplLib.meta_lic", "mplLib.meta_lic", "mplLib.meta_lic", "reciprocal"},
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
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "apacheContainer.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"gplLib.meta_lic", "gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
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
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "lgplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "lgplLib.meta_lic", "lgplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "mitLib.meta_lic", "mitLib.meta_lic", "notice"},
				{"apacheBin.meta_lic", "mitLib.meta_lic", "lgplLib.meta_lic", "restricted"},
				{"lgplLib.meta_lic", "lgplLib.meta_lic", "lgplLib.meta_lic", "restricted"},
				{"mitLib.meta_lic", "mitLib.meta_lic", "mitLib.meta_lic", "notice"},
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
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "mitLib.meta_lic", "mitLib.meta_lic", "notice"},
				{"lgplBin.meta_lic", "lgplBin.meta_lic", "lgplBin.meta_lic", "restricted"},
				{"mitLib.meta_lic", "mitLib.meta_lic", "mitLib.meta_lic", "notice"},
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
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "apacheContainer.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "lgplLib.meta_lic", "restricted"},
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", "lgplLib.meta_lic", "restricted"},
				{"apacheContainer.meta_lic", "lgplLib.meta_lic", "lgplLib.meta_lic", "restricted"},
				{"apacheContainer.meta_lic", "mitLib.meta_lic", "mitLib.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "mitLib.meta_lic", "lgplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "lgplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "lgplLib.meta_lic", "lgplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "mitLib.meta_lic", "mitLib.meta_lic", "notice"},
				{"apacheBin.meta_lic", "mitLib.meta_lic", "lgplLib.meta_lic", "restricted"},
				{"lgplLib.meta_lic", "lgplLib.meta_lic", "lgplLib.meta_lic", "restricted"},
				{"mitLib.meta_lic", "mitLib.meta_lic", "mitLib.meta_lic", "notice"},
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
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "apacheContainer.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "lgplLib.meta_lic", "restricted"},
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "lgplLib.meta_lic", "lgplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"lgplLib.meta_lic", "lgplLib.meta_lic", "lgplLib.meta_lic", "restricted"},
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
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "mitLib.meta_lic", "mitLib.meta_lic", "notice"},
				{"lgplLib.meta_lic", "lgplLib.meta_lic", "lgplLib.meta_lic", "restricted"},
				{"mitLib.meta_lic", "mitLib.meta_lic", "mitLib.meta_lic", "notice"},
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
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "apacheContainer.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"lgplLib.meta_lic", "lgplLib.meta_lic", "lgplLib.meta_lic", "restricted"},
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
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "apacheContainer.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"lgplLib.meta_lic", "lgplLib.meta_lic", "lgplLib.meta_lic", "restricted"},
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
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "mitLib.meta_lic", "mitLib.meta_lic", "notice"},
				{"apacheBin.meta_lic", "mitLib.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
				{"gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
				{"mitLib.meta_lic", "mitLib.meta_lic", "mitLib.meta_lic", "notice"},
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
				{"dependentModule.meta_lic", "dependentModule.meta_lic", "dependentModule.meta_lic", "notice"},
				{"dependentModule.meta_lic", "dependentModule.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
				{"dependentModule.meta_lic", "gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
				{"dependentModule.meta_lic", "mitLib.meta_lic", "mitLib.meta_lic", "notice"},
				{"dependentModule.meta_lic", "mitLib.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
				{"gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
				{"mitLib.meta_lic", "mitLib.meta_lic", "mitLib.meta_lic", "notice"},
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
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "mitLib.meta_lic", "mitLib.meta_lic", "notice"},
				{"gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
				{"mitLib.meta_lic", "mitLib.meta_lic", "mitLib.meta_lic", "notice"},
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
				{"dependentModule.meta_lic", "dependentModule.meta_lic", "dependentModule.meta_lic", "notice"},
				{"dependentModule.meta_lic", "dependentModule.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
				{"dependentModule.meta_lic", "gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
				{"dependentModule.meta_lic", "mitLib.meta_lic", "mitLib.meta_lic", "notice"},
				{"dependentModule.meta_lic", "mitLib.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
				{"gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
				{"mitLib.meta_lic", "mitLib.meta_lic", "mitLib.meta_lic", "notice"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stderr := &bytes.Buffer{}
			lg, err := toGraph(stderr, tt.roots, tt.edges)
			if err != nil {
				t.Errorf("unexpected test data error: got %w, want no error", err)
				return
			}
			expectedRs := toResolutionSet(lg, tt.expectedResolutions)
			actualRs := ResolveTopDownConditions(lg)
			checkSame(actualRs, expectedRs, t)
		})
	}
}

func TestWalkResolutionsForCondition(t *testing.T) {
	tests := []struct {
		name                string
		condition           ConditionNames
		roots               []string
		edges               []annotated
		expectedResolutions []res
	}{
		{
			name: "firstparty",
			condition: ImpliesNotice,
			roots: []string{"apacheBin.meta_lic"},
			edges: []annotated{
				{"apacheBin.meta_lic", "apacheLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "apacheLib.meta_lic", "apacheLib.meta_lic", "notice"},
			},
		},
		{
			name: "notice",
			condition: ImpliesNotice,
			roots: []string{"mitBin.meta_lic"},
			edges: []annotated{
				{"mitBin.meta_lic", "mitLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"mitBin.meta_lic", "mitBin.meta_lic", "mitBin.meta_lic", "notice"},
				{"mitBin.meta_lic", "mitLib.meta_lic", "mitLib.meta_lic", "notice"},
			},
		},
		{
			name: "fponlgplnotice",
			condition: ImpliesNotice,
			roots: []string{"apacheBin.meta_lic"},
			edges: []annotated{
				{"apacheBin.meta_lic", "lgplLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "lgplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "lgplLib.meta_lic", "lgplLib.meta_lic", "restricted"},
			},
		},
		{
			name: "fponlgpldynamicnotice",
			condition: ImpliesNotice,
			roots: []string{"apacheBin.meta_lic"},
			edges: []annotated{
				{"apacheBin.meta_lic", "lgplLib.meta_lic", []string{"dynamic"}},
			},
			expectedResolutions: []res{
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
			},
		},
		{
			name: "independentmodulenotice",
			condition: ImpliesNotice,
			roots: []string{"apacheBin.meta_lic"},
			edges: []annotated{
				{"apacheBin.meta_lic", "gplWithClasspathException.meta_lic", []string{"dynamic"}},
			},
			expectedResolutions: []res{
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
			},
		},
		{
			name: "independentmodulerestricted",
			condition: ImpliesRestricted,
			roots: []string{"apacheBin.meta_lic"},
			edges: []annotated{
				{"apacheBin.meta_lic", "gplWithClasspathException.meta_lic", []string{"dynamic"}},
			},
			expectedResolutions: []res{},
		},
		{
			name: "independentmodulestaticnotice",
			condition: ImpliesNotice,
			roots: []string{"apacheBin.meta_lic"},
			edges: []annotated{
				{"apacheBin.meta_lic", "gplWithClasspathException.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
			},
		},
		{
			name: "independentmodulestaticrestricted",
			condition: ImpliesRestricted,
			roots: []string{"apacheBin.meta_lic"},
			edges: []annotated{
				{"apacheBin.meta_lic", "gplWithClasspathException.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
			},
		},
		{
			name: "dependentmodulenotice",
			condition: ImpliesNotice,
			roots: []string{"dependentModule.meta_lic"},
			edges: []annotated{
				{"dependentModule.meta_lic", "gplWithClasspathException.meta_lic", []string{"dynamic"}},
			},
			expectedResolutions: []res{
				{"dependentModule.meta_lic", "dependentModule.meta_lic", "dependentModule.meta_lic", "notice"},
				{"dependentModule.meta_lic", "dependentModule.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
			},
		},
		{
			name: "dependentmodulerestricted",
			condition: ImpliesRestricted,
			roots: []string{"dependentModule.meta_lic"},
			edges: []annotated{
				{"dependentModule.meta_lic", "gplWithClasspathException.meta_lic", []string{"dynamic"}},
			},
			expectedResolutions: []res{
				{"dependentModule.meta_lic", "dependentModule.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
			},
		},
		{
			name: "lgplonfpnotice",
			condition: ImpliesNotice,
			roots: []string{"lgplBin.meta_lic"},
			edges: []annotated{
				{"lgplBin.meta_lic", "apacheLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"lgplBin.meta_lic", "lgplBin.meta_lic", "lgplBin.meta_lic", "restricted"},
				{"lgplBin.meta_lic", "apacheLib.meta_lic", "apacheLib.meta_lic", "notice"},
				{"lgplBin.meta_lic", "apacheLib.meta_lic", "lgplBin.meta_lic", "restricted"},
			},
		},
		{
			name: "lgplonfprestricted",
			condition: ImpliesRestricted,
			roots: []string{"lgplBin.meta_lic"},
			edges: []annotated{
				{"lgplBin.meta_lic", "apacheLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"lgplBin.meta_lic", "lgplBin.meta_lic", "lgplBin.meta_lic", "restricted"},
				{"lgplBin.meta_lic", "apacheLib.meta_lic", "lgplBin.meta_lic", "restricted"},
			},
		},
		{
			name: "lgplonfpdynamicnotice",
			condition: ImpliesNotice,
			roots: []string{"lgplBin.meta_lic"},
			edges: []annotated{
				{"lgplBin.meta_lic", "apacheLib.meta_lic", []string{"dynamic"}},
			},
			expectedResolutions: []res{
				{"lgplBin.meta_lic", "lgplBin.meta_lic", "lgplBin.meta_lic", "restricted"},
			},
		},
		{
			name: "lgplonfpdynamicrestricted",
			condition: ImpliesRestricted,
			roots: []string{"lgplBin.meta_lic"},
			edges: []annotated{
				{"lgplBin.meta_lic", "apacheLib.meta_lic", []string{"dynamic"}},
			},
			expectedResolutions: []res{
				{"lgplBin.meta_lic", "lgplBin.meta_lic", "lgplBin.meta_lic", "restricted"},
			},
		},
		{
			name: "gplonfpnotice",
			condition: ImpliesNotice,
			roots: []string{"gplBin.meta_lic"},
			edges: []annotated{
				{"gplBin.meta_lic", "apacheLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"gplBin.meta_lic", "gplBin.meta_lic", "gplBin.meta_lic", "restricted"},
				{"gplBin.meta_lic", "apacheLib.meta_lic", "apacheLib.meta_lic", "notice"},
				{"gplBin.meta_lic", "apacheLib.meta_lic", "gplBin.meta_lic", "restricted"},
			},
		},
		{
			name: "gplonfprestricted",
			condition: ImpliesRestricted,
			roots: []string{"gplBin.meta_lic"},
			edges: []annotated{
				{"gplBin.meta_lic", "apacheLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"gplBin.meta_lic", "gplBin.meta_lic", "gplBin.meta_lic", "restricted"},
				{"gplBin.meta_lic", "apacheLib.meta_lic", "gplBin.meta_lic", "restricted"},
			},
		},
		{
			name: "gplcontainernotice",
			condition: ImpliesNotice,
			roots: []string{"gplContainer.meta_lic"},
			edges: []annotated{
				{"gplContainer.meta_lic", "apacheLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"gplContainer.meta_lic", "gplContainer.meta_lic", "gplContainer.meta_lic", "restricted"},
				{"gplContainer.meta_lic", "apacheLib.meta_lic", "apacheLib.meta_lic", "notice"},
				{"gplContainer.meta_lic", "apacheLib.meta_lic", "gplContainer.meta_lic", "restricted"},
				{"apacheLib.meta_lic", "apacheLib.meta_lic", "apacheLib.meta_lic", "notice"},
				{"apacheLib.meta_lic", "apacheLib.meta_lic", "gplContainer.meta_lic", "restricted"},
			},
		},
		{
			name: "gplcontainerrestricted",
			condition: ImpliesRestricted,
			roots: []string{"gplContainer.meta_lic"},
			edges: []annotated{
				{"gplContainer.meta_lic", "apacheLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"gplContainer.meta_lic", "gplContainer.meta_lic", "gplContainer.meta_lic", "restricted"},
				{"gplContainer.meta_lic", "apacheLib.meta_lic", "gplContainer.meta_lic", "restricted"},
				{"apacheLib.meta_lic", "apacheLib.meta_lic", "gplContainer.meta_lic", "restricted"},
			},
		},
		{
			name: "gploncontainernotice",
			condition: ImpliesNotice,
			roots: []string{"apacheContainer.meta_lic"},
			edges: []annotated{
				{"apacheContainer.meta_lic", "apacheLib.meta_lic", []string{"static"}},
				{"apacheContainer.meta_lic", "gplLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "apacheContainer.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheContainer.meta_lic", "apacheLib.meta_lic", "apacheLib.meta_lic", "notice"},
				{"apacheContainer.meta_lic", "gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheLib.meta_lic", "apacheLib.meta_lic", "apacheLib.meta_lic", "notice"},
				{"gplLib.meta_lic", "gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
			},
		},
		{
			name: "gploncontainerrestricted",
			condition: ImpliesRestricted,
			roots: []string{"apacheContainer.meta_lic"},
			edges: []annotated{
				{"apacheContainer.meta_lic", "apacheLib.meta_lic", []string{"static"}},
				{"apacheContainer.meta_lic", "gplLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"apacheContainer.meta_lic", "apacheContainer.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheContainer.meta_lic", "gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
				{"gplLib.meta_lic", "gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
			},
		},
		{
			name: "gplonbinnotice",
			condition: ImpliesNotice,
			roots: []string{"apacheBin.meta_lic"},
			edges: []annotated{
				{"apacheBin.meta_lic", "apacheLib.meta_lic", []string{"static"}},
				{"apacheBin.meta_lic", "gplLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "apacheLib.meta_lic", "apacheLib.meta_lic", "notice"},
				{"apacheBin.meta_lic", "apacheLib.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
			},
		},
		{
			name: "gplonbinrestricted",
			condition: ImpliesRestricted,
			roots: []string{"apacheBin.meta_lic"},
			edges: []annotated{
				{"apacheBin.meta_lic", "apacheLib.meta_lic", []string{"static"}},
				{"apacheBin.meta_lic", "gplLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"apacheBin.meta_lic", "apacheBin.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "apacheLib.meta_lic", "gplLib.meta_lic", "restricted"},
				{"apacheBin.meta_lic", "gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
			},
		},
		{
			name: "gplonfpdynamicnotice",
			condition: ImpliesNotice,
			roots: []string{"gplBin.meta_lic"},
			edges: []annotated{
				{"gplBin.meta_lic", "apacheLib.meta_lic", []string{"dynamic"}},
			},
			expectedResolutions: []res{
				{"gplBin.meta_lic", "gplBin.meta_lic", "gplBin.meta_lic", "restricted"},
			},
		},
		{
			name: "gplonfpdynamicrestricted",
			condition: ImpliesRestricted,
			roots: []string{"gplBin.meta_lic"},
			edges: []annotated{
				{"gplBin.meta_lic", "apacheLib.meta_lic", []string{"dynamic"}},
			},
			expectedResolutions: []res{
				{"gplBin.meta_lic", "gplBin.meta_lic", "gplBin.meta_lic", "restricted"},
			},
		},
		{
			name: "gplonfpdynamicrestrictedshipped",
			condition: ImpliesRestricted,
			roots: []string{"gplBin.meta_lic", "apacheLib.meta_lic"},
			edges: []annotated{
				{"gplBin.meta_lic", "apacheLib.meta_lic", []string{"dynamic"}},
			},
			expectedResolutions: []res{
				{"gplBin.meta_lic", "gplBin.meta_lic", "gplBin.meta_lic", "restricted"},
				{"gplBin.meta_lic", "apacheLib.meta_lic", "gplBin.meta_lic", "restricted"},
			},
		},
		{
			name: "independentmodulereversenotice",
			condition: ImpliesNotice,
			roots: []string{"gplWithClasspathException.meta_lic"},
			edges: []annotated{
				{"gplWithClasspathException.meta_lic", "apacheBin.meta_lic", []string{"dynamic"}},
			},
			expectedResolutions: []res{
				{"gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
			},
		},
		{
			name: "independentmodulereverserestricted",
			condition: ImpliesRestricted,
			roots: []string{"gplWithClasspathException.meta_lic"},
			edges: []annotated{
				{"gplWithClasspathException.meta_lic", "apacheBin.meta_lic", []string{"dynamic"}},
			},
			expectedResolutions: []res{
				{"gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
			},
		},
		{
			name: "independentmodulereverserestrictedshipped",
			condition: ImpliesRestricted,
			roots: []string{"gplWithClasspathException.meta_lic", "apacheBin.meta_lic"},
			edges: []annotated{
				{"gplWithClasspathException.meta_lic", "apacheBin.meta_lic", []string{"dynamic"}},
			},
			expectedResolutions: []res{
				{"gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
			},
		},
		{
			name: "independentmodulereversestaticnotice",
			condition: ImpliesNotice,
			roots: []string{"gplWithClasspathException.meta_lic"},
			edges: []annotated{
				{"gplWithClasspathException.meta_lic", "apacheBin.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
				{"gplWithClasspathException.meta_lic", "apacheBin.meta_lic", "apacheBin.meta_lic", "notice"},
				{"gplWithClasspathException.meta_lic", "apacheBin.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
			},
		},
		{
			name: "independentmodulereversestaticrestricted",
			condition: ImpliesRestricted,
			roots: []string{"gplWithClasspathException.meta_lic"},
			edges: []annotated{
				{"gplWithClasspathException.meta_lic", "apacheBin.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
				{"gplWithClasspathException.meta_lic", "apacheBin.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
			},
		},
		{
			name: "dependentmodulereversenotice",
			condition: ImpliesNotice,
			roots: []string{"gplWithClasspathException.meta_lic"},
			edges: []annotated{
				{"gplWithClasspathException.meta_lic", "dependentModule.meta_lic", []string{"dynamic"}},
			},
			expectedResolutions: []res{
				{"gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
			},
		},
		{
			name: "dependentmodulereverserestricted",
			condition: ImpliesRestricted,
			roots: []string{"gplWithClasspathException.meta_lic"},
			edges: []annotated{
				{"gplWithClasspathException.meta_lic", "dependentModule.meta_lic", []string{"dynamic"}},
			},
			expectedResolutions: []res{
				{"gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
			},
		},
		{
			name: "dependentmodulereverserestrictedshipped",
			condition: ImpliesRestricted,
			roots: []string{"gplWithClasspathException.meta_lic", "dependentModule.meta_lic"},
			edges: []annotated{
				{"gplWithClasspathException.meta_lic", "dependentModule.meta_lic", []string{"dynamic"}},
			},
			expectedResolutions: []res{
				{"gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
				{"gplWithClasspathException.meta_lic", "dependentModule.meta_lic", "gplWithClasspathException.meta_lic", "restricted"},
			},
		},
		{
			name: "ponrnotice",
			condition: ImpliesNotice,
			roots: []string{"proprietary.meta_lic"},
			edges: []annotated{
				{"proprietary.meta_lic", "gplLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"proprietary.meta_lic", "proprietary.meta_lic", "proprietary.meta_lic", "proprietary"},
				{"proprietary.meta_lic", "proprietary.meta_lic", "gplLib.meta_lic", "restricted"},
				{"proprietary.meta_lic", "gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
			},
		},
		{
			name: "ponrrestricted",
			condition: ImpliesRestricted,
			roots: []string{"proprietary.meta_lic"},
			edges: []annotated{
				{"proprietary.meta_lic", "gplLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"proprietary.meta_lic", "gplLib.meta_lic", "gplLib.meta_lic", "restricted"},
				{"proprietary.meta_lic", "proprietary.meta_lic", "gplLib.meta_lic", "restricted"},
			},
		},
		{
			name: "ponrproprietary",
			condition: ImpliesProprietary,
			roots: []string{"proprietary.meta_lic"},
			edges: []annotated{
				{"proprietary.meta_lic", "gplLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"proprietary.meta_lic", "proprietary.meta_lic", "proprietary.meta_lic", "proprietary"},
			},
		},
		{
			name: "ronpnotice",
			condition: ImpliesNotice,
			roots: []string{"gplBin.meta_lic"},
			edges: []annotated{
				{"gplBin.meta_lic", "proprietary.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"gplBin.meta_lic", "gplBin.meta_lic", "gplBin.meta_lic", "restricted"},
				{"gplBin.meta_lic", "proprietary.meta_lic", "proprietary.meta_lic", "proprietary"},
				{"gplBin.meta_lic", "proprietary.meta_lic", "gplBin.meta_lic", "restricted"},
			},
		},
		{
			name: "ronprestricted",
			condition: ImpliesRestricted,
			roots: []string{"gplBin.meta_lic"},
			edges: []annotated{
				{"gplBin.meta_lic", "proprietary.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"gplBin.meta_lic", "gplBin.meta_lic", "gplBin.meta_lic", "restricted"},
				{"gplBin.meta_lic", "proprietary.meta_lic", "gplBin.meta_lic", "restricted"},
			},
		},
		{
			name: "ronpproprietary",
			condition: ImpliesProprietary,
			roots: []string{"gplBin.meta_lic"},
			edges: []annotated{
				{"gplBin.meta_lic", "proprietary.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"gplBin.meta_lic", "proprietary.meta_lic", "proprietary.meta_lic", "proprietary"},
			},
		},
		{
			name: "noticeonb_e_onotice",
			condition: ImpliesNotice,
			roots: []string{"mitBin.meta_lic"},
			edges: []annotated{
				{"mitBin.meta_lic", "by_exception.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"mitBin.meta_lic", "mitBin.meta_lic", "mitBin.meta_lic", "notice"},
				{"mitBin.meta_lic", "by_exception.meta_lic", "by_exception.meta_lic", "by_exception_only"},
			},
		},
		{
			name: "noticeonb_e_orestricted",
			condition: ImpliesRestricted,
			roots: []string{"mitBin.meta_lic"},
			edges: []annotated{
				{"mitBin.meta_lic", "by_exception.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{},
		},
		{
			name: "noticeonb_e_ob_e_o",
			condition: ImpliesByExceptionOnly,
			roots: []string{"mitBin.meta_lic"},
			edges: []annotated{
				{"mitBin.meta_lic", "by_exception.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"mitBin.meta_lic", "by_exception.meta_lic", "by_exception.meta_lic", "by_exception_only"},
			},
		},
		{
			name: "b_e_oonnoticenotice",
			condition: ImpliesNotice,
			roots: []string{"by_exception.meta_lic"},
			edges: []annotated{
				{"by_exception.meta_lic", "mitLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"by_exception.meta_lic", "by_exception.meta_lic", "by_exception.meta_lic", "by_exception_only"},
				{"by_exception.meta_lic", "mitLib.meta_lic", "mitLib.meta_lic", "notice"},
			},
		},
		{
			name: "b_e_oonnoticerestricted",
			condition: ImpliesRestricted,
			roots: []string{"by_exception.meta_lic"},
			edges: []annotated{
				{"by_exception.meta_lic", "mitLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{},
		},
		{
			name: "b_e_oonnoticeb_e_o",
			condition: ImpliesByExceptionOnly,
			roots: []string{"by_exception.meta_lic"},
			edges: []annotated{
				{"by_exception.meta_lic", "mitLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"by_exception.meta_lic", "by_exception.meta_lic", "by_exception.meta_lic", "by_exception_only"},
			},
		},
		{
			name: "noticeonrecipnotice",
			condition: ImpliesNotice,
			roots: []string{"mitBin.meta_lic"},
			edges: []annotated{
				{"mitBin.meta_lic", "mplLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"mitBin.meta_lic", "mitBin.meta_lic", "mitBin.meta_lic", "notice"},
				{"mitBin.meta_lic", "mplLib.meta_lic", "mplLib.meta_lic", "reciprocal"},
			},
		},
		{
			name: "noticeonreciprecip",
			condition: ImpliesReciprocal,
			roots: []string{"mitBin.meta_lic"},
			edges: []annotated{
				{"mitBin.meta_lic", "mplLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"mitBin.meta_lic", "mplLib.meta_lic", "mplLib.meta_lic", "reciprocal"},
			},
		},
		{
			name: "reciponnoticenotice",
			condition: ImpliesNotice,
			roots: []string{"mplBin.meta_lic"},
			edges: []annotated{
				{"mplBin.meta_lic", "mitLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"mplBin.meta_lic", "mplBin.meta_lic", "mplBin.meta_lic", "reciprocal"},
				{"mplBin.meta_lic", "mitLib.meta_lic", "mitLib.meta_lic", "notice"},
			},
		},
		{
			name: "reciponnoticerecip",
			condition: ImpliesReciprocal,
			roots: []string{"mplBin.meta_lic"},
			edges: []annotated{
				{"mplBin.meta_lic", "mitLib.meta_lic", []string{"static"}},
			},
			expectedResolutions: []res{
				{"mplBin.meta_lic", "mplBin.meta_lic", "mplBin.meta_lic", "reciprocal"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stderr := &bytes.Buffer{}
			lg, err := toGraph(stderr, tt.roots, tt.edges)
			if err != nil {
				t.Errorf("unexpected test data error: got %w, want no error", err)
				return
			}
			expectedRs := toResolutionSet(lg, tt.expectedResolutions)
			actualRs := WalkResolutionsForCondition(lg, ResolveTopDownConditions(lg), tt.condition)
			checkSame(actualRs, expectedRs, t)
		})
	}
}
