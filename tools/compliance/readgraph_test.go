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
	"io"
	"io/fs"
	"sort"
	"strings"
	"testing"
)

const (
	AOSP = `` +
`package_name: "Android"
license_kinds: "SPDX-license-identifier-Apache-2.0"
license_conditions: "notice"
`
)

type testFS map[string][]byte

func (fs *testFS) Open(name string) (fs.File, error) {
	if _, ok := (*fs)[name]; !ok {
		return nil, fmt.Errorf("unknown file %q", name)
	}
	return &testFile{fs, name, 0}, nil
}

type testFile struct {
	fs   *testFS
	name string
	posn int
}

func (f *testFile) Stat() (fs.FileInfo, error) {
	return nil, fmt.Errorf("unimplemented")
}

func (f *testFile) Read(b []byte) (int, error) {
	if f.posn < 0 {
		return 0, fmt.Errorf("file not open: %q", f.name)
	}
	if f.posn >= len((*f.fs)[f.name]) {
		return 0, io.EOF
	}
	n := copy(b, (*f.fs)[f.name][f.posn:])
	f.posn += n
	return n, nil
}

func (f *testFile) Close() error {
	if f.posn < 0 {
		return fmt.Errorf("file aready closed: %q", f.name)
	}
	f.posn = -1
	return nil
}

type edge struct {
	target, dep string
}

func (e edge) String() string {
	return e.target + " -> " + e.dep
}

// byEdge orders edges by target then dep name.
type byEdge []edge

func (l byEdge) Len() int      { return len(l) }
func (l byEdge) Swap(i, j int) { l[i], l[j] = l[j], l[i] }
func (l byEdge) Less(i, j int) bool {
	return l[i].target < l[j].target || (l[i].target == l[j].target && l[i].dep < l[j].dep)
}

func TestReadLicenseGraph(t *testing.T) {
	tests := []struct {
		name         string
		fs           *testFS
		roots        []string
		expectedError string
		expectedEdges []edge
		expectedTargets []string
	}{
		{
			name: "trivial",
			fs: &testFS{
				"app.meta_lic": []byte("package_name: \"Android\"\n"),
			},
			roots: []string{"app.meta_lic"},
			expectedEdges: []edge{},
			expectedTargets: []string{"app.meta_lic"},
		},
		{
			name: "unterminated",
			fs: &testFS{
				"app.meta_lic": []byte("package_name: \"Android\n"),
			},
			roots:         []string{"app.meta_lic"},
			expectedError: `invalid character '\n' in string`,
		},
		{
			name: "danglingref",
			fs: &testFS{
				"app.meta_lic": []byte(AOSP + "deps: {\n  file: \"lib.meta_lic\"\n}\n"),
			},
			roots:         []string{"app.meta_lic"},
			expectedError: `unknown file "lib.meta_lic"`,
		},
		{
			name: "singleedge",
			fs: &testFS{
				"app.meta_lic": []byte(AOSP + "deps: {\n  file: \"lib.meta_lic\"\n}\n"),
				"lib.meta_lic": []byte(AOSP),
			},
			roots:         []string{"app.meta_lic"},
			expectedEdges: []edge{{"app.meta_lic", "lib.meta_lic"}},
			expectedTargets: []string{"app.meta_lic", "lib.meta_lic"},
		},
		{
			name: "fullgraph",
			fs: &testFS{
				"apex.meta_lic": []byte(AOSP + "deps: {\n  file: \"app.meta_lic\"\n}\ndeps: {\n  file: \"bin.meta_lic\"\n}\n"),
				"app.meta_lic": []byte(AOSP),
				"bin.meta_lic": []byte(AOSP + "deps: {\n  file: \"lib.meta_lic\"\n}\n"),
				"lib.meta_lic": []byte(AOSP),
			},
			roots:         []string{"apex.meta_lic"},
			expectedEdges: []edge{
				{"apex.meta_lic", "app.meta_lic"},
				{"apex.meta_lic", "bin.meta_lic"},
				{"bin.meta_lic", "lib.meta_lic"},
			},
			expectedTargets: []string{"apex.meta_lic", "app.meta_lic", "bin.meta_lic", "lib.meta_lic"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stderr := &bytes.Buffer{}
			lg, err := ReadLicenseGraph(tt.fs, stderr, tt.roots)
			if err != nil {
				if len(tt.expectedError) == 0 {
					t.Errorf("unexpected error: got %w, want no error", err)
				} else if !strings.Contains(err.Error(), tt.expectedError) {
					t.Errorf("unexpected error: got %w, want %q", err, tt.expectedError)
				}
				return
			}
			if 0 < len(tt.expectedError) {
				t.Errorf("unexpected success: got no error, want %q err", tt.expectedError)
				return
			}
			if lg == nil {
				t.Errorf("missing license graph: got nil, want license graph")
				return
			}
			actualEdges := make([]edge, 0)
			for _, e := range lg.Edges() {
				actualEdges = append(actualEdges, edge{e.Target().Name(), e.Dependency().Name()})
			}
			sort.Sort(byEdge(tt.expectedEdges))
			sort.Sort(byEdge(actualEdges))
			if len(tt.expectedEdges) != len(actualEdges) {
				t.Errorf("unexpected number of edges: got %v with %d elements, want %v with %d elements",
					actualEdges, len(actualEdges), tt.expectedEdges, len(tt.expectedEdges))
			} else {
				for i := 0; i < len(actualEdges); i++ {
					if tt.expectedEdges[i] != actualEdges[i] {
						t.Errorf("unexpected edge at element %d: got %s, want %s", i, actualEdges[i], tt.expectedEdges[i])
					}
				}
			}
			actualTargets := make([]string, 0)
			for _, t := range lg.Targets() {
				actualTargets = append(actualTargets, t.Name())
			}
			sort.Strings(tt.expectedTargets)
			sort.Strings(actualTargets)
			if len(tt.expectedTargets) != len(actualTargets) {
				t.Errorf("unexpected number of targets: got %v with %d elements, want %v with %d elements",
					actualTargets, len(actualTargets), tt.expectedTargets, len(tt.expectedTargets))
			} else {
				for i := 0; i < len(actualTargets); i++ {
					if tt.expectedTargets[i] != actualTargets[i] {
						t.Errorf("unexpected target at element %d: got %s, want %s", i, actualTargets[i], tt.expectedTargets[i])
					}
				}
			}
		})
	}
}
