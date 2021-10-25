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
			expectedResolutions: []res{},
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
			actualRs := ResolveBottomUpConditions(licenseGraph)
			checkSame(actualRs, expectedRs, t)
		})
	}
}
