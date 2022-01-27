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

package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	// Change into the testdata directory before running the tests.
	if err := os.Chdir("testdata"); err != nil {
		fmt.Printf("failed to change to testdata directory: %s\n", err)
		os.Exit(1)
	}
	os.Exit(m.Run())
}

func Test_plaintext(t *testing.T) {
	tests := []struct {
		condition   string
		name        string
		roots       []string
		ctx         context
		expectedOut []string
	}{
		{
			condition: "firstparty",
			name:      "apex",
			roots:     []string{"highest.apex.meta_lic"},
			expectedOut: []string{
				"firstparty/bin/bin1.meta_lic firstparty/lib/liba.so.meta_lic static",
				"firstparty/bin/bin1.meta_lic firstparty/lib/libc.a.meta_lic static",
				"firstparty/bin/bin2.meta_lic firstparty/lib/libb.so.meta_lic dynamic",
				"firstparty/bin/bin2.meta_lic firstparty/lib/libd.so.meta_lic dynamic",
				"firstparty/highest.apex.meta_lic firstparty/bin/bin1.meta_lic static",
				"firstparty/highest.apex.meta_lic firstparty/bin/bin2.meta_lic static",
				"firstparty/highest.apex.meta_lic firstparty/lib/liba.so.meta_lic static",
				"firstparty/highest.apex.meta_lic firstparty/lib/libb.so.meta_lic static",
			},
		},
		{
			condition: "firstparty",
			name:      "apex_trimmed",
			roots:     []string{"highest.apex.meta_lic"},
			ctx:       context{stripPrefix: "firstparty/"},
			expectedOut: []string{
				"bin/bin1.meta_lic lib/liba.so.meta_lic static",
				"bin/bin1.meta_lic lib/libc.a.meta_lic static",
				"bin/bin2.meta_lic lib/libb.so.meta_lic dynamic",
				"bin/bin2.meta_lic lib/libd.so.meta_lic dynamic",
				"highest.apex.meta_lic bin/bin1.meta_lic static",
				"highest.apex.meta_lic bin/bin2.meta_lic static",
				"highest.apex.meta_lic lib/liba.so.meta_lic static",
				"highest.apex.meta_lic lib/libb.so.meta_lic static",
			},
		},
		{
			condition: "firstparty",
			name:      "apex_trimmed_labelled",
			roots:     []string{"highest.apex.meta_lic"},
			ctx:       context{stripPrefix: "firstparty/", labelConditions: true},
			expectedOut: []string{
				"bin/bin1.meta_lic:notice lib/liba.so.meta_lic:notice static",
				"bin/bin1.meta_lic:notice lib/libc.a.meta_lic:notice static",
				"bin/bin2.meta_lic:notice lib/libb.so.meta_lic:notice dynamic",
				"bin/bin2.meta_lic:notice lib/libd.so.meta_lic:notice dynamic",
				"highest.apex.meta_lic:notice bin/bin1.meta_lic:notice static",
				"highest.apex.meta_lic:notice bin/bin2.meta_lic:notice static",
				"highest.apex.meta_lic:notice lib/liba.so.meta_lic:notice static",
				"highest.apex.meta_lic:notice lib/libb.so.meta_lic:notice static",
			},
		},
		{
			condition: "firstparty",
			name:      "container",
			roots:     []string{"container.zip.meta_lic"},
			expectedOut: []string{
				"firstparty/bin/bin1.meta_lic firstparty/lib/liba.so.meta_lic static",
				"firstparty/bin/bin1.meta_lic firstparty/lib/libc.a.meta_lic static",
				"firstparty/bin/bin2.meta_lic firstparty/lib/libb.so.meta_lic dynamic",
				"firstparty/bin/bin2.meta_lic firstparty/lib/libd.so.meta_lic dynamic",
				"firstparty/container.zip.meta_lic firstparty/bin/bin1.meta_lic static",
				"firstparty/container.zip.meta_lic firstparty/bin/bin2.meta_lic static",
				"firstparty/container.zip.meta_lic firstparty/lib/liba.so.meta_lic static",
				"firstparty/container.zip.meta_lic firstparty/lib/libb.so.meta_lic static",
			},
		},
		{
			condition: "firstparty",
			name:      "application",
			roots:     []string{"application.meta_lic"},
			expectedOut: []string{
				"firstparty/application.meta_lic firstparty/bin/bin3.meta_lic toolchain",
				"firstparty/application.meta_lic firstparty/lib/liba.so.meta_lic static",
				"firstparty/application.meta_lic firstparty/lib/libb.so.meta_lic dynamic",
			},
		},
		{
			condition: "firstparty",
			name:      "binary",
			roots:     []string{"bin/bin1.meta_lic"},
			expectedOut: []string{
				"firstparty/bin/bin1.meta_lic firstparty/lib/liba.so.meta_lic static",
				"firstparty/bin/bin1.meta_lic firstparty/lib/libc.a.meta_lic static",
			},
		},
		{
			condition:   "firstparty",
			name:        "library",
			roots:       []string{"lib/libd.so.meta_lic"},
			expectedOut: []string{},
		},
		{
			condition: "notice",
			name:      "apex",
			roots:     []string{"highest.apex.meta_lic"},
			expectedOut: []string{
				"notice/bin/bin1.meta_lic notice/lib/liba.so.meta_lic static",
				"notice/bin/bin1.meta_lic notice/lib/libc.a.meta_lic static",
				"notice/bin/bin2.meta_lic notice/lib/libb.so.meta_lic dynamic",
				"notice/bin/bin2.meta_lic notice/lib/libd.so.meta_lic dynamic",
				"notice/highest.apex.meta_lic notice/bin/bin1.meta_lic static",
				"notice/highest.apex.meta_lic notice/bin/bin2.meta_lic static",
				"notice/highest.apex.meta_lic notice/lib/liba.so.meta_lic static",
				"notice/highest.apex.meta_lic notice/lib/libb.so.meta_lic static",
			},
		},
		{
			condition: "notice",
			name:      "apex_trimmed",
			roots:     []string{"highest.apex.meta_lic"},
			ctx:       context{stripPrefix: "notice/"},
			expectedOut: []string{
				"bin/bin1.meta_lic lib/liba.so.meta_lic static",
				"bin/bin1.meta_lic lib/libc.a.meta_lic static",
				"bin/bin2.meta_lic lib/libb.so.meta_lic dynamic",
				"bin/bin2.meta_lic lib/libd.so.meta_lic dynamic",
				"highest.apex.meta_lic bin/bin1.meta_lic static",
				"highest.apex.meta_lic bin/bin2.meta_lic static",
				"highest.apex.meta_lic lib/liba.so.meta_lic static",
				"highest.apex.meta_lic lib/libb.so.meta_lic static",
			},
		},
		{
			condition: "notice",
			name:      "apex_trimmed_labelled",
			roots:     []string{"highest.apex.meta_lic"},
			ctx:       context{stripPrefix: "notice/", labelConditions: true},
			expectedOut: []string{
				"bin/bin1.meta_lic:notice lib/liba.so.meta_lic:notice static",
				"bin/bin1.meta_lic:notice lib/libc.a.meta_lic:notice static",
				"bin/bin2.meta_lic:notice lib/libb.so.meta_lic:notice dynamic",
				"bin/bin2.meta_lic:notice lib/libd.so.meta_lic:notice dynamic",
				"highest.apex.meta_lic:notice bin/bin1.meta_lic:notice static",
				"highest.apex.meta_lic:notice bin/bin2.meta_lic:notice static",
				"highest.apex.meta_lic:notice lib/liba.so.meta_lic:notice static",
				"highest.apex.meta_lic:notice lib/libb.so.meta_lic:notice static",
			},
		},
		{
			condition: "notice",
			name:      "container",
			roots:     []string{"container.zip.meta_lic"},
			expectedOut: []string{
				"notice/bin/bin1.meta_lic notice/lib/liba.so.meta_lic static",
				"notice/bin/bin1.meta_lic notice/lib/libc.a.meta_lic static",
				"notice/bin/bin2.meta_lic notice/lib/libb.so.meta_lic dynamic",
				"notice/bin/bin2.meta_lic notice/lib/libd.so.meta_lic dynamic",
				"notice/container.zip.meta_lic notice/bin/bin1.meta_lic static",
				"notice/container.zip.meta_lic notice/bin/bin2.meta_lic static",
				"notice/container.zip.meta_lic notice/lib/liba.so.meta_lic static",
				"notice/container.zip.meta_lic notice/lib/libb.so.meta_lic static",
			},
		},
		{
			condition: "notice",
			name:      "application",
			roots:     []string{"application.meta_lic"},
			expectedOut: []string{
				"notice/application.meta_lic notice/bin/bin3.meta_lic toolchain",
				"notice/application.meta_lic notice/lib/liba.so.meta_lic static",
				"notice/application.meta_lic notice/lib/libb.so.meta_lic dynamic",
			},
		},
		{
			condition: "notice",
			name:      "binary",
			roots:     []string{"bin/bin1.meta_lic"},
			expectedOut: []string{
				"notice/bin/bin1.meta_lic notice/lib/liba.so.meta_lic static",
				"notice/bin/bin1.meta_lic notice/lib/libc.a.meta_lic static",
			},
		},
		{
			condition:   "notice",
			name:        "library",
			roots:       []string{"lib/libd.so.meta_lic"},
			expectedOut: []string{},
		},
		{
			condition: "reciprocal",
			name:      "apex",
			roots:     []string{"highest.apex.meta_lic"},
			expectedOut: []string{
				"reciprocal/bin/bin1.meta_lic reciprocal/lib/liba.so.meta_lic static",
				"reciprocal/bin/bin1.meta_lic reciprocal/lib/libc.a.meta_lic static",
				"reciprocal/bin/bin2.meta_lic reciprocal/lib/libb.so.meta_lic dynamic",
				"reciprocal/bin/bin2.meta_lic reciprocal/lib/libd.so.meta_lic dynamic",
				"reciprocal/highest.apex.meta_lic reciprocal/bin/bin1.meta_lic static",
				"reciprocal/highest.apex.meta_lic reciprocal/bin/bin2.meta_lic static",
				"reciprocal/highest.apex.meta_lic reciprocal/lib/liba.so.meta_lic static",
				"reciprocal/highest.apex.meta_lic reciprocal/lib/libb.so.meta_lic static",
			},
		},
		{
			condition: "reciprocal",
			name:      "apex_trimmed",
			roots:     []string{"highest.apex.meta_lic"},
			ctx:       context{stripPrefix: "reciprocal/"},
			expectedOut: []string{
				"bin/bin1.meta_lic lib/liba.so.meta_lic static",
				"bin/bin1.meta_lic lib/libc.a.meta_lic static",
				"bin/bin2.meta_lic lib/libb.so.meta_lic dynamic",
				"bin/bin2.meta_lic lib/libd.so.meta_lic dynamic",
				"highest.apex.meta_lic bin/bin1.meta_lic static",
				"highest.apex.meta_lic bin/bin2.meta_lic static",
				"highest.apex.meta_lic lib/liba.so.meta_lic static",
				"highest.apex.meta_lic lib/libb.so.meta_lic static",
			},
		},
		{
			condition: "reciprocal",
			name:      "apex_trimmed_labelled",
			roots:     []string{"highest.apex.meta_lic"},
			ctx:       context{stripPrefix: "reciprocal/", labelConditions: true},
			expectedOut: []string{
				"bin/bin1.meta_lic:notice lib/liba.so.meta_lic:reciprocal static",
				"bin/bin1.meta_lic:notice lib/libc.a.meta_lic:reciprocal static",
				"bin/bin2.meta_lic:notice lib/libb.so.meta_lic:notice dynamic",
				"bin/bin2.meta_lic:notice lib/libd.so.meta_lic:notice dynamic",
				"highest.apex.meta_lic:notice bin/bin1.meta_lic:notice static",
				"highest.apex.meta_lic:notice bin/bin2.meta_lic:notice static",
				"highest.apex.meta_lic:notice lib/liba.so.meta_lic:reciprocal static",
				"highest.apex.meta_lic:notice lib/libb.so.meta_lic:notice static",
			},
		},
		{
			condition: "reciprocal",
			name:      "container",
			roots:     []string{"container.zip.meta_lic"},
			expectedOut: []string{
				"reciprocal/bin/bin1.meta_lic reciprocal/lib/liba.so.meta_lic static",
				"reciprocal/bin/bin1.meta_lic reciprocal/lib/libc.a.meta_lic static",
				"reciprocal/bin/bin2.meta_lic reciprocal/lib/libb.so.meta_lic dynamic",
				"reciprocal/bin/bin2.meta_lic reciprocal/lib/libd.so.meta_lic dynamic",
				"reciprocal/container.zip.meta_lic reciprocal/bin/bin1.meta_lic static",
				"reciprocal/container.zip.meta_lic reciprocal/bin/bin2.meta_lic static",
				"reciprocal/container.zip.meta_lic reciprocal/lib/liba.so.meta_lic static",
				"reciprocal/container.zip.meta_lic reciprocal/lib/libb.so.meta_lic static",
			},
		},
		{
			condition: "reciprocal",
			name:      "application",
			roots:     []string{"application.meta_lic"},
			expectedOut: []string{
				"reciprocal/application.meta_lic reciprocal/bin/bin3.meta_lic toolchain",
				"reciprocal/application.meta_lic reciprocal/lib/liba.so.meta_lic static",
				"reciprocal/application.meta_lic reciprocal/lib/libb.so.meta_lic dynamic",
			},
		},
		{
			condition: "reciprocal",
			name:      "binary",
			roots:     []string{"bin/bin1.meta_lic"},
			expectedOut: []string{
				"reciprocal/bin/bin1.meta_lic reciprocal/lib/liba.so.meta_lic static",
				"reciprocal/bin/bin1.meta_lic reciprocal/lib/libc.a.meta_lic static",
			},
		},
		{
			condition:   "reciprocal",
			name:        "library",
			roots:       []string{"lib/libd.so.meta_lic"},
			expectedOut: []string{},
		},
		{
			condition: "restricted",
			name:      "apex",
			roots:     []string{"highest.apex.meta_lic"},
			expectedOut: []string{
				"restricted/bin/bin1.meta_lic restricted/lib/liba.so.meta_lic static",
				"restricted/bin/bin1.meta_lic restricted/lib/libc.a.meta_lic static",
				"restricted/bin/bin2.meta_lic restricted/lib/libb.so.meta_lic dynamic",
				"restricted/bin/bin2.meta_lic restricted/lib/libd.so.meta_lic dynamic",
				"restricted/highest.apex.meta_lic restricted/bin/bin1.meta_lic static",
				"restricted/highest.apex.meta_lic restricted/bin/bin2.meta_lic static",
				"restricted/highest.apex.meta_lic restricted/lib/liba.so.meta_lic static",
				"restricted/highest.apex.meta_lic restricted/lib/libb.so.meta_lic static",
			},
		},
		{
			condition: "restricted",
			name:      "apex_trimmed",
			roots:     []string{"highest.apex.meta_lic"},
			ctx:       context{stripPrefix: "restricted/"},
			expectedOut: []string{
				"bin/bin1.meta_lic lib/liba.so.meta_lic static",
				"bin/bin1.meta_lic lib/libc.a.meta_lic static",
				"bin/bin2.meta_lic lib/libb.so.meta_lic dynamic",
				"bin/bin2.meta_lic lib/libd.so.meta_lic dynamic",
				"highest.apex.meta_lic bin/bin1.meta_lic static",
				"highest.apex.meta_lic bin/bin2.meta_lic static",
				"highest.apex.meta_lic lib/liba.so.meta_lic static",
				"highest.apex.meta_lic lib/libb.so.meta_lic static",
			},
		},
		{
			condition: "restricted",
			name:      "apex_trimmed_labelled",
			roots:     []string{"highest.apex.meta_lic"},
			ctx:       context{stripPrefix: "restricted/", labelConditions: true},
			expectedOut: []string{
				"bin/bin1.meta_lic:notice lib/liba.so.meta_lic:restricted_allows_dynamic_linking static",
				"bin/bin1.meta_lic:notice lib/libc.a.meta_lic:reciprocal static",
				"bin/bin2.meta_lic:notice lib/libb.so.meta_lic:restricted dynamic",
				"bin/bin2.meta_lic:notice lib/libd.so.meta_lic:notice dynamic",
				"highest.apex.meta_lic:notice bin/bin1.meta_lic:notice static",
				"highest.apex.meta_lic:notice bin/bin2.meta_lic:notice static",
				"highest.apex.meta_lic:notice lib/liba.so.meta_lic:restricted_allows_dynamic_linking static",
				"highest.apex.meta_lic:notice lib/libb.so.meta_lic:restricted static",
			},
		},
		{
			condition: "restricted",
			name:      "container",
			roots:     []string{"container.zip.meta_lic"},
			expectedOut: []string{
				"restricted/bin/bin1.meta_lic restricted/lib/liba.so.meta_lic static",
				"restricted/bin/bin1.meta_lic restricted/lib/libc.a.meta_lic static",
				"restricted/bin/bin2.meta_lic restricted/lib/libb.so.meta_lic dynamic",
				"restricted/bin/bin2.meta_lic restricted/lib/libd.so.meta_lic dynamic",
				"restricted/container.zip.meta_lic restricted/bin/bin1.meta_lic static",
				"restricted/container.zip.meta_lic restricted/bin/bin2.meta_lic static",
				"restricted/container.zip.meta_lic restricted/lib/liba.so.meta_lic static",
				"restricted/container.zip.meta_lic restricted/lib/libb.so.meta_lic static",
			},
		},
		{
			condition: "restricted",
			name:      "application",
			roots:     []string{"application.meta_lic"},
			expectedOut: []string{
				"restricted/application.meta_lic restricted/bin/bin3.meta_lic toolchain",
				"restricted/application.meta_lic restricted/lib/liba.so.meta_lic static",
				"restricted/application.meta_lic restricted/lib/libb.so.meta_lic dynamic",
			},
		},
		{
			condition: "restricted",
			name:      "binary",
			roots:     []string{"bin/bin1.meta_lic"},
			expectedOut: []string{
				"restricted/bin/bin1.meta_lic restricted/lib/liba.so.meta_lic static",
				"restricted/bin/bin1.meta_lic restricted/lib/libc.a.meta_lic static",
			},
		},
		{
			condition:   "restricted",
			name:        "library",
			roots:       []string{"lib/libd.so.meta_lic"},
			expectedOut: []string{},
		},
		{
			condition: "proprietary",
			name:      "apex",
			roots:     []string{"highest.apex.meta_lic"},
			expectedOut: []string{
				"proprietary/bin/bin1.meta_lic proprietary/lib/liba.so.meta_lic static",
				"proprietary/bin/bin1.meta_lic proprietary/lib/libc.a.meta_lic static",
				"proprietary/bin/bin2.meta_lic proprietary/lib/libb.so.meta_lic dynamic",
				"proprietary/bin/bin2.meta_lic proprietary/lib/libd.so.meta_lic dynamic",
				"proprietary/highest.apex.meta_lic proprietary/bin/bin1.meta_lic static",
				"proprietary/highest.apex.meta_lic proprietary/bin/bin2.meta_lic static",
				"proprietary/highest.apex.meta_lic proprietary/lib/liba.so.meta_lic static",
				"proprietary/highest.apex.meta_lic proprietary/lib/libb.so.meta_lic static",
			},
		},
		{
			condition: "proprietary",
			name:      "apex_trimmed",
			roots:     []string{"highest.apex.meta_lic"},
			ctx:       context{stripPrefix: "proprietary/"},
			expectedOut: []string{
				"bin/bin1.meta_lic lib/liba.so.meta_lic static",
				"bin/bin1.meta_lic lib/libc.a.meta_lic static",
				"bin/bin2.meta_lic lib/libb.so.meta_lic dynamic",
				"bin/bin2.meta_lic lib/libd.so.meta_lic dynamic",
				"highest.apex.meta_lic bin/bin1.meta_lic static",
				"highest.apex.meta_lic bin/bin2.meta_lic static",
				"highest.apex.meta_lic lib/liba.so.meta_lic static",
				"highest.apex.meta_lic lib/libb.so.meta_lic static",
			},
		},
		{
			condition: "proprietary",
			name:      "apex_trimmed_labelled",
			roots:     []string{"highest.apex.meta_lic"},
			ctx:       context{stripPrefix: "proprietary/", labelConditions: true},
			expectedOut: []string{
				"bin/bin1.meta_lic:notice lib/liba.so.meta_lic:by_exception_only:proprietary static",
				"bin/bin1.meta_lic:notice lib/libc.a.meta_lic:by_exception_only:proprietary static",
				"bin/bin2.meta_lic:by_exception_only:proprietary lib/libb.so.meta_lic:restricted dynamic",
				"bin/bin2.meta_lic:by_exception_only:proprietary lib/libd.so.meta_lic:notice dynamic",
				"highest.apex.meta_lic:notice bin/bin1.meta_lic:notice static",
				"highest.apex.meta_lic:notice bin/bin2.meta_lic:by_exception_only:proprietary static",
				"highest.apex.meta_lic:notice lib/liba.so.meta_lic:by_exception_only:proprietary static",
				"highest.apex.meta_lic:notice lib/libb.so.meta_lic:restricted static",
			},
		},
		{
			condition: "proprietary",
			name:      "container",
			roots:     []string{"container.zip.meta_lic"},
			expectedOut: []string{
				"proprietary/bin/bin1.meta_lic proprietary/lib/liba.so.meta_lic static",
				"proprietary/bin/bin1.meta_lic proprietary/lib/libc.a.meta_lic static",
				"proprietary/bin/bin2.meta_lic proprietary/lib/libb.so.meta_lic dynamic",
				"proprietary/bin/bin2.meta_lic proprietary/lib/libd.so.meta_lic dynamic",
				"proprietary/container.zip.meta_lic proprietary/bin/bin1.meta_lic static",
				"proprietary/container.zip.meta_lic proprietary/bin/bin2.meta_lic static",
				"proprietary/container.zip.meta_lic proprietary/lib/liba.so.meta_lic static",
				"proprietary/container.zip.meta_lic proprietary/lib/libb.so.meta_lic static",
			},
		},
		{
			condition: "proprietary",
			name:      "application",
			roots:     []string{"application.meta_lic"},
			expectedOut: []string{
				"proprietary/application.meta_lic proprietary/bin/bin3.meta_lic toolchain",
				"proprietary/application.meta_lic proprietary/lib/liba.so.meta_lic static",
				"proprietary/application.meta_lic proprietary/lib/libb.so.meta_lic dynamic",
			},
		},
		{
			condition: "proprietary",
			name:      "binary",
			roots:     []string{"bin/bin1.meta_lic"},
			expectedOut: []string{
				"proprietary/bin/bin1.meta_lic proprietary/lib/liba.so.meta_lic static",
				"proprietary/bin/bin1.meta_lic proprietary/lib/libc.a.meta_lic static",
			},
		},
		{
			condition:   "proprietary",
			name:        "library",
			roots:       []string{"lib/libd.so.meta_lic"},
			expectedOut: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.condition+" "+tt.name, func(t *testing.T) {
			expectedOut := &bytes.Buffer{}
			for _, eo := range tt.expectedOut {
				expectedOut.WriteString(eo)
				expectedOut.WriteString("\n")
			}

			stdout := &bytes.Buffer{}
			stderr := &bytes.Buffer{}

			rootFiles := make([]string, 0, len(tt.roots))
			for _, r := range tt.roots {
				rootFiles = append(rootFiles, tt.condition+"/"+r)
			}
			err := dumpGraph(&tt.ctx, stdout, stderr, rootFiles...)
			if err != nil {
				t.Fatalf("dumpgraph: error = %v, stderr = %v", err, stderr)
				return
			}
			if stderr.Len() > 0 {
				t.Errorf("dumpgraph: gotStderr = %v, want none", stderr)
			}
			out := stdout.String()
			expected := expectedOut.String()
			if out != expected {
				outList := strings.Split(out, "\n")
				expectedList := strings.Split(expected, "\n")
				startLine := 0
				for len(outList) > startLine && len(expectedList) > startLine && outList[startLine] == expectedList[startLine] {
					startLine++
				}
				t.Errorf("listshare: gotStdout = %v, want %v, somewhere near line %d Stdout = %v, want %v",
					out, expected, startLine+1, outList[startLine], expectedList[startLine])
			}
		})
	}
}

type testContext struct {
	nextNode int
	nodes    map[string]string
}

type matcher interface {
	matchString(*testContext) string
	typeString() string
}

type targetMatcher struct {
	target     string
	conditions []string
}

func (tm *targetMatcher) matchString(ctx *testContext) string {
	m := tm.target
	if len(tm.conditions) > 0 {
		m += "\\n" + strings.Join(tm.conditions, "\\n")
	}
	m = ctx.nodes[tm.target] + " [label=\"" + m + "\"];"
	return m
}

func (tm *targetMatcher) typeString() string {
	return "target"
}

type edgeMatcher struct {
	target      string
	dep         string
	annotations []string
}

func (em *edgeMatcher) matchString(ctx *testContext) string {
	return ctx.nodes[em.dep] + " -> " + ctx.nodes[em.target] + " [label=\"" + strings.Join(em.annotations, "\\n") + "\"];"
}

func (tm *edgeMatcher) typeString() string {
	return "edge"
}

type getMatcher func(*testContext) matcher

func matchTarget(target string, conditions ...string) getMatcher {
	return func(ctx *testContext) matcher {
		ctx.nodes[target] = fmt.Sprintf("n%d", ctx.nextNode)
		ctx.nextNode++
		return &targetMatcher{target, append([]string{}, conditions...)}
	}
}

func matchEdge(target, dep string, annotations ...string) getMatcher {
	return func(ctx *testContext) matcher {
		if _, ok := ctx.nodes[target]; !ok {
			panic(fmt.Errorf("no node for target %v in %v -> %v [label=\"%s\"];", target, dep, target, strings.Join(annotations, "\\n")))
		}
		if _, ok := ctx.nodes[dep]; !ok {
			panic(fmt.Errorf("no node for dep %v in %v -> %v [label=\"%s\"];", target, dep, target, strings.Join(annotations, "\\n")))
		}
		return &edgeMatcher{target, dep, append([]string{}, annotations...)}
	}
}

func Test_graphviz(t *testing.T) {
	tests := []struct {
		condition   string
		name        string
		roots       []string
		ctx         context
		expectedOut []getMatcher
	}{
		{
			condition: "firstparty",
			name:      "apex",
			roots:     []string{"highest.apex.meta_lic"},
			expectedOut: []getMatcher{
				matchTarget("firstparty/bin/bin1.meta_lic"),
				matchTarget("firstparty/bin/bin2.meta_lic"),
				matchTarget("firstparty/highest.apex.meta_lic"),
				matchTarget("firstparty/lib/liba.so.meta_lic"),
				matchTarget("firstparty/lib/libb.so.meta_lic"),
				matchTarget("firstparty/lib/libc.a.meta_lic"),
				matchTarget("firstparty/lib/libd.so.meta_lic"),
				matchEdge("firstparty/bin/bin1.meta_lic", "firstparty/lib/liba.so.meta_lic", "static"),
				matchEdge("firstparty/bin/bin1.meta_lic", "firstparty/lib/libc.a.meta_lic", "static"),
				matchEdge("firstparty/bin/bin2.meta_lic", "firstparty/lib/libb.so.meta_lic", "dynamic"),
				matchEdge("firstparty/bin/bin2.meta_lic", "firstparty/lib/libd.so.meta_lic", "dynamic"),
				matchEdge("firstparty/highest.apex.meta_lic", "firstparty/bin/bin1.meta_lic", "static"),
				matchEdge("firstparty/highest.apex.meta_lic", "firstparty/bin/bin2.meta_lic", "static"),
				matchEdge("firstparty/highest.apex.meta_lic", "firstparty/lib/liba.so.meta_lic", "static"),
				matchEdge("firstparty/highest.apex.meta_lic", "firstparty/lib/libb.so.meta_lic", "static"),
			},
		},
		{
			condition: "firstparty",
			name:      "apex_trimmed",
			roots:     []string{"highest.apex.meta_lic"},
			ctx:       context{stripPrefix: "firstparty/"},
			expectedOut: []getMatcher{
				matchTarget("bin/bin1.meta_lic"),
				matchTarget("bin/bin2.meta_lic"),
				matchTarget("highest.apex.meta_lic"),
				matchTarget("lib/liba.so.meta_lic"),
				matchTarget("lib/libb.so.meta_lic"),
				matchTarget("lib/libc.a.meta_lic"),
				matchTarget("lib/libd.so.meta_lic"),
				matchEdge("bin/bin1.meta_lic", "lib/liba.so.meta_lic", "static"),
				matchEdge("bin/bin1.meta_lic", "lib/libc.a.meta_lic", "static"),
				matchEdge("bin/bin2.meta_lic", "lib/libb.so.meta_lic", "dynamic"),
				matchEdge("bin/bin2.meta_lic", "lib/libd.so.meta_lic", "dynamic"),
				matchEdge("highest.apex.meta_lic", "bin/bin1.meta_lic", "static"),
				matchEdge("highest.apex.meta_lic", "bin/bin2.meta_lic", "static"),
				matchEdge("highest.apex.meta_lic", "lib/liba.so.meta_lic", "static"),
				matchEdge("highest.apex.meta_lic", "lib/libb.so.meta_lic", "static"),
			},
		},
		{
			condition: "firstparty",
			name:      "apex_trimmed_labelled",
			roots:     []string{"highest.apex.meta_lic"},
			ctx:       context{stripPrefix: "firstparty/", labelConditions: true},
			expectedOut: []getMatcher{
				matchTarget("bin/bin1.meta_lic", "notice"),
				matchTarget("bin/bin2.meta_lic", "notice"),
				matchTarget("highest.apex.meta_lic", "notice"),
				matchTarget("lib/liba.so.meta_lic", "notice"),
				matchTarget("lib/libb.so.meta_lic", "notice"),
				matchTarget("lib/libc.a.meta_lic", "notice"),
				matchTarget("lib/libd.so.meta_lic", "notice"),
				matchEdge("bin/bin1.meta_lic", "lib/liba.so.meta_lic", "static"),
				matchEdge("bin/bin1.meta_lic", "lib/libc.a.meta_lic", "static"),
				matchEdge("bin/bin2.meta_lic", "lib/libb.so.meta_lic", "dynamic"),
				matchEdge("bin/bin2.meta_lic", "lib/libd.so.meta_lic", "dynamic"),
				matchEdge("highest.apex.meta_lic", "bin/bin1.meta_lic", "static"),
				matchEdge("highest.apex.meta_lic", "bin/bin2.meta_lic", "static"),
				matchEdge("highest.apex.meta_lic", "lib/liba.so.meta_lic", "static"),
				matchEdge("highest.apex.meta_lic", "lib/libb.so.meta_lic", "static"),
			},
		},
		{
			condition: "firstparty",
			name:      "container",
			roots:     []string{"container.zip.meta_lic"},
			expectedOut: []getMatcher{
				matchTarget("firstparty/bin/bin1.meta_lic"),
				matchTarget("firstparty/bin/bin2.meta_lic"),
				matchTarget("firstparty/container.zip.meta_lic"),
				matchTarget("firstparty/lib/liba.so.meta_lic"),
				matchTarget("firstparty/lib/libb.so.meta_lic"),
				matchTarget("firstparty/lib/libc.a.meta_lic"),
				matchTarget("firstparty/lib/libd.so.meta_lic"),
				matchEdge("firstparty/bin/bin1.meta_lic", "firstparty/lib/liba.so.meta_lic", "static"),
				matchEdge("firstparty/bin/bin1.meta_lic", "firstparty/lib/libc.a.meta_lic", "static"),
				matchEdge("firstparty/bin/bin2.meta_lic", "firstparty/lib/libb.so.meta_lic", "dynamic"),
				matchEdge("firstparty/bin/bin2.meta_lic", "firstparty/lib/libd.so.meta_lic", "dynamic"),
				matchEdge("firstparty/container.zip.meta_lic", "firstparty/bin/bin1.meta_lic", "static"),
				matchEdge("firstparty/container.zip.meta_lic", "firstparty/bin/bin2.meta_lic", "static"),
				matchEdge("firstparty/container.zip.meta_lic", "firstparty/lib/liba.so.meta_lic", "static"),
				matchEdge("firstparty/container.zip.meta_lic", "firstparty/lib/libb.so.meta_lic", "static"),
			},
		},
		{
			condition: "firstparty",
			name:      "application",
			roots:     []string{"application.meta_lic"},
			expectedOut: []getMatcher{
				matchTarget("firstparty/application.meta_lic"),
				matchTarget("firstparty/bin/bin3.meta_lic"),
				matchTarget("firstparty/lib/liba.so.meta_lic"),
				matchTarget("firstparty/lib/libb.so.meta_lic"),
				matchEdge("firstparty/application.meta_lic", "firstparty/bin/bin3.meta_lic", "toolchain"),
				matchEdge("firstparty/application.meta_lic", "firstparty/lib/liba.so.meta_lic", "static"),
				matchEdge("firstparty/application.meta_lic", "firstparty/lib/libb.so.meta_lic", "dynamic"),
			},
		},
		{
			condition: "firstparty",
			name:      "binary",
			roots:     []string{"bin/bin1.meta_lic"},
			expectedOut: []getMatcher{
				matchTarget("firstparty/bin/bin1.meta_lic"),
				matchTarget("firstparty/lib/liba.so.meta_lic"),
				matchTarget("firstparty/lib/libc.a.meta_lic"),
				matchEdge("firstparty/bin/bin1.meta_lic", "firstparty/lib/liba.so.meta_lic", "static"),
				matchEdge("firstparty/bin/bin1.meta_lic", "firstparty/lib/libc.a.meta_lic", "static"),
			},
		},
		{
			condition:   "firstparty",
			name:        "library",
			roots:       []string{"lib/libd.so.meta_lic"},
			expectedOut: []getMatcher{matchTarget("firstparty/lib/libd.so.meta_lic")},
		},
		{
			condition: "notice",
			name:      "apex",
			roots:     []string{"highest.apex.meta_lic"},
			expectedOut: []getMatcher{
				matchTarget("notice/bin/bin1.meta_lic"),
				matchTarget("notice/bin/bin2.meta_lic"),
				matchTarget("notice/highest.apex.meta_lic"),
				matchTarget("notice/lib/liba.so.meta_lic"),
				matchTarget("notice/lib/libb.so.meta_lic"),
				matchTarget("notice/lib/libc.a.meta_lic"),
				matchTarget("notice/lib/libd.so.meta_lic"),
				matchEdge("notice/bin/bin1.meta_lic", "notice/lib/liba.so.meta_lic", "static"),
				matchEdge("notice/bin/bin1.meta_lic", "notice/lib/libc.a.meta_lic", "static"),
				matchEdge("notice/bin/bin2.meta_lic", "notice/lib/libb.so.meta_lic", "dynamic"),
				matchEdge("notice/bin/bin2.meta_lic", "notice/lib/libd.so.meta_lic", "dynamic"),
				matchEdge("notice/highest.apex.meta_lic", "notice/bin/bin1.meta_lic", "static"),
				matchEdge("notice/highest.apex.meta_lic", "notice/bin/bin2.meta_lic", "static"),
				matchEdge("notice/highest.apex.meta_lic", "notice/lib/liba.so.meta_lic", "static"),
				matchEdge("notice/highest.apex.meta_lic", "notice/lib/libb.so.meta_lic", "static"),
			},
		},
		{
			condition: "notice",
			name:      "apex_trimmed",
			roots:     []string{"highest.apex.meta_lic"},
			ctx:       context{stripPrefix: "notice/"},
			expectedOut: []getMatcher{
				matchTarget("bin/bin1.meta_lic"),
				matchTarget("bin/bin2.meta_lic"),
				matchTarget("highest.apex.meta_lic"),
				matchTarget("lib/liba.so.meta_lic"),
				matchTarget("lib/libb.so.meta_lic"),
				matchTarget("lib/libc.a.meta_lic"),
				matchTarget("lib/libd.so.meta_lic"),
				matchEdge("bin/bin1.meta_lic", "lib/liba.so.meta_lic", "static"),
				matchEdge("bin/bin1.meta_lic", "lib/libc.a.meta_lic", "static"),
				matchEdge("bin/bin2.meta_lic", "lib/libb.so.meta_lic", "dynamic"),
				matchEdge("bin/bin2.meta_lic", "lib/libd.so.meta_lic", "dynamic"),
				matchEdge("highest.apex.meta_lic", "bin/bin1.meta_lic", "static"),
				matchEdge("highest.apex.meta_lic", "bin/bin2.meta_lic", "static"),
				matchEdge("highest.apex.meta_lic", "lib/liba.so.meta_lic", "static"),
				matchEdge("highest.apex.meta_lic", "lib/libb.so.meta_lic", "static"),
			},
		},
		{
			condition: "notice",
			name:      "apex_trimmed_labelled",
			roots:     []string{"highest.apex.meta_lic"},
			ctx:       context{stripPrefix: "notice/", labelConditions: true},
			expectedOut: []getMatcher{
				matchTarget("bin/bin1.meta_lic", "notice"),
				matchTarget("bin/bin2.meta_lic", "notice"),
				matchTarget("highest.apex.meta_lic", "notice"),
				matchTarget("lib/liba.so.meta_lic", "notice"),
				matchTarget("lib/libb.so.meta_lic", "notice"),
				matchTarget("lib/libc.a.meta_lic", "notice"),
				matchTarget("lib/libd.so.meta_lic", "notice"),
				matchEdge("bin/bin1.meta_lic", "lib/liba.so.meta_lic", "static"),
				matchEdge("bin/bin1.meta_lic", "lib/libc.a.meta_lic", "static"),
				matchEdge("bin/bin2.meta_lic", "lib/libb.so.meta_lic", "dynamic"),
				matchEdge("bin/bin2.meta_lic", "lib/libd.so.meta_lic", "dynamic"),
				matchEdge("highest.apex.meta_lic", "bin/bin1.meta_lic", "static"),
				matchEdge("highest.apex.meta_lic", "bin/bin2.meta_lic", "static"),
				matchEdge("highest.apex.meta_lic", "lib/liba.so.meta_lic", "static"),
				matchEdge("highest.apex.meta_lic", "lib/libb.so.meta_lic", "static"),
			},
		},
		{
			condition: "notice",
			name:      "container",
			roots:     []string{"container.zip.meta_lic"},
			expectedOut: []getMatcher{
				matchTarget("notice/bin/bin1.meta_lic"),
				matchTarget("notice/bin/bin2.meta_lic"),
				matchTarget("notice/container.zip.meta_lic"),
				matchTarget("notice/lib/liba.so.meta_lic"),
				matchTarget("notice/lib/libb.so.meta_lic"),
				matchTarget("notice/lib/libc.a.meta_lic"),
				matchTarget("notice/lib/libd.so.meta_lic"),
				matchEdge("notice/bin/bin1.meta_lic", "notice/lib/liba.so.meta_lic", "static"),
				matchEdge("notice/bin/bin1.meta_lic", "notice/lib/libc.a.meta_lic", "static"),
				matchEdge("notice/bin/bin2.meta_lic", "notice/lib/libb.so.meta_lic", "dynamic"),
				matchEdge("notice/bin/bin2.meta_lic", "notice/lib/libd.so.meta_lic", "dynamic"),
				matchEdge("notice/container.zip.meta_lic", "notice/bin/bin1.meta_lic", "static"),
				matchEdge("notice/container.zip.meta_lic", "notice/bin/bin2.meta_lic", "static"),
				matchEdge("notice/container.zip.meta_lic", "notice/lib/liba.so.meta_lic", "static"),
				matchEdge("notice/container.zip.meta_lic", "notice/lib/libb.so.meta_lic", "static"),
			},
		},
		{
			condition: "notice",
			name:      "application",
			roots:     []string{"application.meta_lic"},
			expectedOut: []getMatcher{
				matchTarget("notice/application.meta_lic"),
				matchTarget("notice/bin/bin3.meta_lic"),
				matchTarget("notice/lib/liba.so.meta_lic"),
				matchTarget("notice/lib/libb.so.meta_lic"),
				matchEdge("notice/application.meta_lic", "notice/bin/bin3.meta_lic", "toolchain"),
				matchEdge("notice/application.meta_lic", "notice/lib/liba.so.meta_lic", "static"),
				matchEdge("notice/application.meta_lic", "notice/lib/libb.so.meta_lic", "dynamic"),
			},
		},
		{
			condition: "notice",
			name:      "binary",
			roots:     []string{"bin/bin1.meta_lic"},
			expectedOut: []getMatcher{
				matchTarget("notice/bin/bin1.meta_lic"),
				matchTarget("notice/lib/liba.so.meta_lic"),
				matchTarget("notice/lib/libc.a.meta_lic"),
				matchEdge("notice/bin/bin1.meta_lic", "notice/lib/liba.so.meta_lic", "static"),
				matchEdge("notice/bin/bin1.meta_lic", "notice/lib/libc.a.meta_lic", "static"),
			},
		},
		{
			condition:   "notice",
			name:        "library",
			roots:       []string{"lib/libd.so.meta_lic"},
			expectedOut: []getMatcher{matchTarget("notice/lib/libd.so.meta_lic")},
		},
		{
			condition: "reciprocal",
			name:      "apex",
			roots:     []string{"highest.apex.meta_lic"},
			expectedOut: []getMatcher{
				matchTarget("reciprocal/bin/bin1.meta_lic"),
				matchTarget("reciprocal/bin/bin2.meta_lic"),
				matchTarget("reciprocal/highest.apex.meta_lic"),
				matchTarget("reciprocal/lib/liba.so.meta_lic"),
				matchTarget("reciprocal/lib/libb.so.meta_lic"),
				matchTarget("reciprocal/lib/libc.a.meta_lic"),
				matchTarget("reciprocal/lib/libd.so.meta_lic"),
				matchEdge("reciprocal/bin/bin1.meta_lic", "reciprocal/lib/liba.so.meta_lic", "static"),
				matchEdge("reciprocal/bin/bin1.meta_lic", "reciprocal/lib/libc.a.meta_lic", "static"),
				matchEdge("reciprocal/bin/bin2.meta_lic", "reciprocal/lib/libb.so.meta_lic", "dynamic"),
				matchEdge("reciprocal/bin/bin2.meta_lic", "reciprocal/lib/libd.so.meta_lic", "dynamic"),
				matchEdge("reciprocal/highest.apex.meta_lic", "reciprocal/bin/bin1.meta_lic", "static"),
				matchEdge("reciprocal/highest.apex.meta_lic", "reciprocal/bin/bin2.meta_lic", "static"),
				matchEdge("reciprocal/highest.apex.meta_lic", "reciprocal/lib/liba.so.meta_lic", "static"),
				matchEdge("reciprocal/highest.apex.meta_lic", "reciprocal/lib/libb.so.meta_lic", "static"),
			},
		},
		{
			condition: "reciprocal",
			name:      "apex_trimmed",
			roots:     []string{"highest.apex.meta_lic"},
			ctx:       context{stripPrefix: "reciprocal/"},
			expectedOut: []getMatcher{
				matchTarget("bin/bin1.meta_lic"),
				matchTarget("bin/bin2.meta_lic"),
				matchTarget("highest.apex.meta_lic"),
				matchTarget("lib/liba.so.meta_lic"),
				matchTarget("lib/libb.so.meta_lic"),
				matchTarget("lib/libc.a.meta_lic"),
				matchTarget("lib/libd.so.meta_lic"),
				matchEdge("bin/bin1.meta_lic", "lib/liba.so.meta_lic", "static"),
				matchEdge("bin/bin1.meta_lic", "lib/libc.a.meta_lic", "static"),
				matchEdge("bin/bin2.meta_lic", "lib/libb.so.meta_lic", "dynamic"),
				matchEdge("bin/bin2.meta_lic", "lib/libd.so.meta_lic", "dynamic"),
				matchEdge("highest.apex.meta_lic", "bin/bin1.meta_lic", "static"),
				matchEdge("highest.apex.meta_lic", "bin/bin2.meta_lic", "static"),
				matchEdge("highest.apex.meta_lic", "lib/liba.so.meta_lic", "static"),
				matchEdge("highest.apex.meta_lic", "lib/libb.so.meta_lic", "static"),
			},
		},
		{
			condition: "reciprocal",
			name:      "apex_trimmed_labelled",
			roots:     []string{"highest.apex.meta_lic"},
			ctx:       context{stripPrefix: "reciprocal/", labelConditions: true},
			expectedOut: []getMatcher{
				matchTarget("bin/bin1.meta_lic", "notice"),
				matchTarget("bin/bin2.meta_lic", "notice"),
				matchTarget("highest.apex.meta_lic", "notice"),
				matchTarget("lib/liba.so.meta_lic", "reciprocal"),
				matchTarget("lib/libb.so.meta_lic", "notice"),
				matchTarget("lib/libc.a.meta_lic", "reciprocal"),
				matchTarget("lib/libd.so.meta_lic", "notice"),
				matchEdge("bin/bin1.meta_lic", "lib/liba.so.meta_lic", "static"),
				matchEdge("bin/bin1.meta_lic", "lib/libc.a.meta_lic", "static"),
				matchEdge("bin/bin2.meta_lic", "lib/libb.so.meta_lic", "dynamic"),
				matchEdge("bin/bin2.meta_lic", "lib/libd.so.meta_lic", "dynamic"),
				matchEdge("highest.apex.meta_lic", "bin/bin1.meta_lic", "static"),
				matchEdge("highest.apex.meta_lic", "bin/bin2.meta_lic", "static"),
				matchEdge("highest.apex.meta_lic", "lib/liba.so.meta_lic", "static"),
				matchEdge("highest.apex.meta_lic", "lib/libb.so.meta_lic", "static"),
			},
		},
		{
			condition: "reciprocal",
			name:      "container",
			roots:     []string{"container.zip.meta_lic"},
			expectedOut: []getMatcher{
				matchTarget("reciprocal/bin/bin1.meta_lic"),
				matchTarget("reciprocal/bin/bin2.meta_lic"),
				matchTarget("reciprocal/container.zip.meta_lic"),
				matchTarget("reciprocal/lib/liba.so.meta_lic"),
				matchTarget("reciprocal/lib/libb.so.meta_lic"),
				matchTarget("reciprocal/lib/libc.a.meta_lic"),
				matchTarget("reciprocal/lib/libd.so.meta_lic"),
				matchEdge("reciprocal/bin/bin1.meta_lic", "reciprocal/lib/liba.so.meta_lic", "static"),
				matchEdge("reciprocal/bin/bin1.meta_lic", "reciprocal/lib/libc.a.meta_lic", "static"),
				matchEdge("reciprocal/bin/bin2.meta_lic", "reciprocal/lib/libb.so.meta_lic", "dynamic"),
				matchEdge("reciprocal/bin/bin2.meta_lic", "reciprocal/lib/libd.so.meta_lic", "dynamic"),
				matchEdge("reciprocal/container.zip.meta_lic", "reciprocal/bin/bin1.meta_lic", "static"),
				matchEdge("reciprocal/container.zip.meta_lic", "reciprocal/bin/bin2.meta_lic", "static"),
				matchEdge("reciprocal/container.zip.meta_lic", "reciprocal/lib/liba.so.meta_lic", "static"),
				matchEdge("reciprocal/container.zip.meta_lic", "reciprocal/lib/libb.so.meta_lic", "static"),
			},
		},
		{
			condition: "reciprocal",
			name:      "application",
			roots:     []string{"application.meta_lic"},
			expectedOut: []getMatcher{
				matchTarget("reciprocal/application.meta_lic"),
				matchTarget("reciprocal/bin/bin3.meta_lic"),
				matchTarget("reciprocal/lib/liba.so.meta_lic"),
				matchTarget("reciprocal/lib/libb.so.meta_lic"),
				matchEdge("reciprocal/application.meta_lic", "reciprocal/bin/bin3.meta_lic", "toolchain"),
				matchEdge("reciprocal/application.meta_lic", "reciprocal/lib/liba.so.meta_lic", "static"),
				matchEdge("reciprocal/application.meta_lic", "reciprocal/lib/libb.so.meta_lic", "dynamic"),
			},
		},
		{
			condition: "reciprocal",
			name:      "binary",
			roots:     []string{"bin/bin1.meta_lic"},
			expectedOut: []getMatcher{
				matchTarget("reciprocal/bin/bin1.meta_lic"),
				matchTarget("reciprocal/lib/liba.so.meta_lic"),
				matchTarget("reciprocal/lib/libc.a.meta_lic"),
				matchEdge("reciprocal/bin/bin1.meta_lic", "reciprocal/lib/liba.so.meta_lic", "static"),
				matchEdge("reciprocal/bin/bin1.meta_lic", "reciprocal/lib/libc.a.meta_lic", "static"),
			},
		},
		{
			condition:   "reciprocal",
			name:        "library",
			roots:       []string{"lib/libd.so.meta_lic"},
			expectedOut: []getMatcher{matchTarget("reciprocal/lib/libd.so.meta_lic")},
		},
		{
			condition: "restricted",
			name:      "apex",
			roots:     []string{"highest.apex.meta_lic"},
			expectedOut: []getMatcher{
				matchTarget("restricted/bin/bin1.meta_lic"),
				matchTarget("restricted/bin/bin2.meta_lic"),
				matchTarget("restricted/highest.apex.meta_lic"),
				matchTarget("restricted/lib/liba.so.meta_lic"),
				matchTarget("restricted/lib/libb.so.meta_lic"),
				matchTarget("restricted/lib/libc.a.meta_lic"),
				matchTarget("restricted/lib/libd.so.meta_lic"),
				matchEdge("restricted/bin/bin1.meta_lic", "restricted/lib/liba.so.meta_lic", "static"),
				matchEdge("restricted/bin/bin1.meta_lic", "restricted/lib/libc.a.meta_lic", "static"),
				matchEdge("restricted/bin/bin2.meta_lic", "restricted/lib/libb.so.meta_lic", "dynamic"),
				matchEdge("restricted/bin/bin2.meta_lic", "restricted/lib/libd.so.meta_lic", "dynamic"),
				matchEdge("restricted/highest.apex.meta_lic", "restricted/bin/bin1.meta_lic", "static"),
				matchEdge("restricted/highest.apex.meta_lic", "restricted/bin/bin2.meta_lic", "static"),
				matchEdge("restricted/highest.apex.meta_lic", "restricted/lib/liba.so.meta_lic", "static"),
				matchEdge("restricted/highest.apex.meta_lic", "restricted/lib/libb.so.meta_lic", "static"),
			},
		},
		{
			condition: "restricted",
			name:      "apex_trimmed",
			roots:     []string{"highest.apex.meta_lic"},
			ctx:       context{stripPrefix: "restricted/"},
			expectedOut: []getMatcher{
				matchTarget("bin/bin1.meta_lic"),
				matchTarget("bin/bin2.meta_lic"),
				matchTarget("highest.apex.meta_lic"),
				matchTarget("lib/liba.so.meta_lic"),
				matchTarget("lib/libb.so.meta_lic"),
				matchTarget("lib/libc.a.meta_lic"),
				matchTarget("lib/libd.so.meta_lic"),
				matchEdge("bin/bin1.meta_lic", "lib/liba.so.meta_lic", "static"),
				matchEdge("bin/bin1.meta_lic", "lib/libc.a.meta_lic", "static"),
				matchEdge("bin/bin2.meta_lic", "lib/libb.so.meta_lic", "dynamic"),
				matchEdge("bin/bin2.meta_lic", "lib/libd.so.meta_lic", "dynamic"),
				matchEdge("highest.apex.meta_lic", "bin/bin1.meta_lic", "static"),
				matchEdge("highest.apex.meta_lic", "bin/bin2.meta_lic", "static"),
				matchEdge("highest.apex.meta_lic", "lib/liba.so.meta_lic", "static"),
				matchEdge("highest.apex.meta_lic", "lib/libb.so.meta_lic", "static"),
			},
		},
		{
			condition: "restricted",
			name:      "apex_trimmed_labelled",
			roots:     []string{"highest.apex.meta_lic"},
			ctx:       context{stripPrefix: "restricted/", labelConditions: true},
			expectedOut: []getMatcher{
				matchTarget("bin/bin1.meta_lic", "notice"),
				matchTarget("bin/bin2.meta_lic", "notice"),
				matchTarget("highest.apex.meta_lic", "notice"),
				matchTarget("lib/liba.so.meta_lic", "restricted_allows_dynamic_linking"),
				matchTarget("lib/libb.so.meta_lic", "restricted"),
				matchTarget("lib/libc.a.meta_lic", "reciprocal"),
				matchTarget("lib/libd.so.meta_lic", "notice"),
				matchEdge("bin/bin1.meta_lic", "lib/liba.so.meta_lic", "static"),
				matchEdge("bin/bin1.meta_lic", "lib/libc.a.meta_lic", "static"),
				matchEdge("bin/bin2.meta_lic", "lib/libb.so.meta_lic", "dynamic"),
				matchEdge("bin/bin2.meta_lic", "lib/libd.so.meta_lic", "dynamic"),
				matchEdge("highest.apex.meta_lic", "bin/bin1.meta_lic", "static"),
				matchEdge("highest.apex.meta_lic", "bin/bin2.meta_lic", "static"),
				matchEdge("highest.apex.meta_lic", "lib/liba.so.meta_lic", "static"),
				matchEdge("highest.apex.meta_lic", "lib/libb.so.meta_lic", "static"),
			},
		},
		{
			condition: "restricted",
			name:      "container",
			roots:     []string{"container.zip.meta_lic"},
			expectedOut: []getMatcher{
				matchTarget("restricted/bin/bin1.meta_lic"),
				matchTarget("restricted/bin/bin2.meta_lic"),
				matchTarget("restricted/container.zip.meta_lic"),
				matchTarget("restricted/lib/liba.so.meta_lic"),
				matchTarget("restricted/lib/libb.so.meta_lic"),
				matchTarget("restricted/lib/libc.a.meta_lic"),
				matchTarget("restricted/lib/libd.so.meta_lic"),
				matchEdge("restricted/bin/bin1.meta_lic", "restricted/lib/liba.so.meta_lic", "static"),
				matchEdge("restricted/bin/bin1.meta_lic", "restricted/lib/libc.a.meta_lic", "static"),
				matchEdge("restricted/bin/bin2.meta_lic", "restricted/lib/libb.so.meta_lic", "dynamic"),
				matchEdge("restricted/bin/bin2.meta_lic", "restricted/lib/libd.so.meta_lic", "dynamic"),
				matchEdge("restricted/container.zip.meta_lic", "restricted/bin/bin1.meta_lic", "static"),
				matchEdge("restricted/container.zip.meta_lic", "restricted/bin/bin2.meta_lic", "static"),
				matchEdge("restricted/container.zip.meta_lic", "restricted/lib/liba.so.meta_lic", "static"),
				matchEdge("restricted/container.zip.meta_lic", "restricted/lib/libb.so.meta_lic", "static"),
			},
		},
		{
			condition: "restricted",
			name:      "application",
			roots:     []string{"application.meta_lic"},
			expectedOut: []getMatcher{
				matchTarget("restricted/application.meta_lic"),
				matchTarget("restricted/bin/bin3.meta_lic"),
				matchTarget("restricted/lib/liba.so.meta_lic"),
				matchTarget("restricted/lib/libb.so.meta_lic"),
				matchEdge("restricted/application.meta_lic", "restricted/bin/bin3.meta_lic", "toolchain"),
				matchEdge("restricted/application.meta_lic", "restricted/lib/liba.so.meta_lic", "static"),
				matchEdge("restricted/application.meta_lic", "restricted/lib/libb.so.meta_lic", "dynamic"),
			},
		},
		{
			condition: "restricted",
			name:      "binary",
			roots:     []string{"bin/bin1.meta_lic"},
			expectedOut: []getMatcher{
				matchTarget("restricted/bin/bin1.meta_lic"),
				matchTarget("restricted/lib/liba.so.meta_lic"),
				matchTarget("restricted/lib/libc.a.meta_lic"),
				matchEdge("restricted/bin/bin1.meta_lic", "restricted/lib/liba.so.meta_lic", "static"),
				matchEdge("restricted/bin/bin1.meta_lic", "restricted/lib/libc.a.meta_lic", "static"),
			},
		},
		{
			condition:   "restricted",
			name:        "library",
			roots:       []string{"lib/libd.so.meta_lic"},
			expectedOut: []getMatcher{matchTarget("restricted/lib/libd.so.meta_lic")},
		},
		{
			condition: "proprietary",
			name:      "apex",
			roots:     []string{"highest.apex.meta_lic"},
			expectedOut: []getMatcher{
				matchTarget("proprietary/bin/bin1.meta_lic"),
				matchTarget("proprietary/bin/bin2.meta_lic"),
				matchTarget("proprietary/highest.apex.meta_lic"),
				matchTarget("proprietary/lib/liba.so.meta_lic"),
				matchTarget("proprietary/lib/libb.so.meta_lic"),
				matchTarget("proprietary/lib/libc.a.meta_lic"),
				matchTarget("proprietary/lib/libd.so.meta_lic"),
				matchEdge("proprietary/bin/bin1.meta_lic", "proprietary/lib/liba.so.meta_lic", "static"),
				matchEdge("proprietary/bin/bin1.meta_lic", "proprietary/lib/libc.a.meta_lic", "static"),
				matchEdge("proprietary/bin/bin2.meta_lic", "proprietary/lib/libb.so.meta_lic", "dynamic"),
				matchEdge("proprietary/bin/bin2.meta_lic", "proprietary/lib/libd.so.meta_lic", "dynamic"),
				matchEdge("proprietary/highest.apex.meta_lic", "proprietary/bin/bin1.meta_lic", "static"),
				matchEdge("proprietary/highest.apex.meta_lic", "proprietary/bin/bin2.meta_lic", "static"),
				matchEdge("proprietary/highest.apex.meta_lic", "proprietary/lib/liba.so.meta_lic", "static"),
				matchEdge("proprietary/highest.apex.meta_lic", "proprietary/lib/libb.so.meta_lic", "static"),
			},
		},
		{
			condition: "proprietary",
			name:      "apex_trimmed",
			roots:     []string{"highest.apex.meta_lic"},
			ctx:       context{stripPrefix: "proprietary/"},
			expectedOut: []getMatcher{
				matchTarget("bin/bin1.meta_lic"),
				matchTarget("bin/bin2.meta_lic"),
				matchTarget("highest.apex.meta_lic"),
				matchTarget("lib/liba.so.meta_lic"),
				matchTarget("lib/libb.so.meta_lic"),
				matchTarget("lib/libc.a.meta_lic"),
				matchTarget("lib/libd.so.meta_lic"),
				matchEdge("bin/bin1.meta_lic", "lib/liba.so.meta_lic", "static"),
				matchEdge("bin/bin1.meta_lic", "lib/libc.a.meta_lic", "static"),
				matchEdge("bin/bin2.meta_lic", "lib/libb.so.meta_lic", "dynamic"),
				matchEdge("bin/bin2.meta_lic", "lib/libd.so.meta_lic", "dynamic"),
				matchEdge("highest.apex.meta_lic", "bin/bin1.meta_lic", "static"),
				matchEdge("highest.apex.meta_lic", "bin/bin2.meta_lic", "static"),
				matchEdge("highest.apex.meta_lic", "lib/liba.so.meta_lic", "static"),
				matchEdge("highest.apex.meta_lic", "lib/libb.so.meta_lic", "static"),
			},
		},
		{
			condition: "proprietary",
			name:      "apex_trimmed_labelled",
			roots:     []string{"highest.apex.meta_lic"},
			ctx:       context{stripPrefix: "proprietary/", labelConditions: true},
			expectedOut: []getMatcher{
				matchTarget("bin/bin1.meta_lic", "notice"),
				matchTarget("bin/bin2.meta_lic", "by_exception_only", "proprietary"),
				matchTarget("highest.apex.meta_lic", "notice"),
				matchTarget("lib/liba.so.meta_lic", "by_exception_only", "proprietary"),
				matchTarget("lib/libb.so.meta_lic", "restricted"),
				matchTarget("lib/libc.a.meta_lic", "by_exception_only", "proprietary"),
				matchTarget("lib/libd.so.meta_lic", "notice"),
				matchEdge("bin/bin1.meta_lic", "lib/liba.so.meta_lic", "static"),
				matchEdge("bin/bin1.meta_lic", "lib/libc.a.meta_lic", "static"),
				matchEdge("bin/bin2.meta_lic", "lib/libb.so.meta_lic", "dynamic"),
				matchEdge("bin/bin2.meta_lic", "lib/libd.so.meta_lic", "dynamic"),
				matchEdge("highest.apex.meta_lic", "bin/bin1.meta_lic", "static"),
				matchEdge("highest.apex.meta_lic", "bin/bin2.meta_lic", "static"),
				matchEdge("highest.apex.meta_lic", "lib/liba.so.meta_lic", "static"),
				matchEdge("highest.apex.meta_lic", "lib/libb.so.meta_lic", "static"),
			},
		},
		{
			condition: "proprietary",
			name:      "container",
			roots:     []string{"container.zip.meta_lic"},
			expectedOut: []getMatcher{
				matchTarget("proprietary/bin/bin1.meta_lic"),
				matchTarget("proprietary/bin/bin2.meta_lic"),
				matchTarget("proprietary/container.zip.meta_lic"),
				matchTarget("proprietary/lib/liba.so.meta_lic"),
				matchTarget("proprietary/lib/libb.so.meta_lic"),
				matchTarget("proprietary/lib/libc.a.meta_lic"),
				matchTarget("proprietary/lib/libd.so.meta_lic"),
				matchEdge("proprietary/bin/bin1.meta_lic", "proprietary/lib/liba.so.meta_lic", "static"),
				matchEdge("proprietary/bin/bin1.meta_lic", "proprietary/lib/libc.a.meta_lic", "static"),
				matchEdge("proprietary/bin/bin2.meta_lic", "proprietary/lib/libb.so.meta_lic", "dynamic"),
				matchEdge("proprietary/bin/bin2.meta_lic", "proprietary/lib/libd.so.meta_lic", "dynamic"),
				matchEdge("proprietary/container.zip.meta_lic", "proprietary/bin/bin1.meta_lic", "static"),
				matchEdge("proprietary/container.zip.meta_lic", "proprietary/bin/bin2.meta_lic", "static"),
				matchEdge("proprietary/container.zip.meta_lic", "proprietary/lib/liba.so.meta_lic", "static"),
				matchEdge("proprietary/container.zip.meta_lic", "proprietary/lib/libb.so.meta_lic", "static"),
			},
		},
		{
			condition: "proprietary",
			name:      "application",
			roots:     []string{"application.meta_lic"},
			expectedOut: []getMatcher{
				matchTarget("proprietary/application.meta_lic"),
				matchTarget("proprietary/bin/bin3.meta_lic"),
				matchTarget("proprietary/lib/liba.so.meta_lic"),
				matchTarget("proprietary/lib/libb.so.meta_lic"),
				matchEdge("proprietary/application.meta_lic", "proprietary/bin/bin3.meta_lic", "toolchain"),
				matchEdge("proprietary/application.meta_lic", "proprietary/lib/liba.so.meta_lic", "static"),
				matchEdge("proprietary/application.meta_lic", "proprietary/lib/libb.so.meta_lic", "dynamic"),
			},
		},
		{
			condition: "proprietary",
			name:      "binary",
			roots:     []string{"bin/bin1.meta_lic"},
			expectedOut: []getMatcher{
				matchTarget("proprietary/bin/bin1.meta_lic"),
				matchTarget("proprietary/lib/liba.so.meta_lic"),
				matchTarget("proprietary/lib/libc.a.meta_lic"),
				matchEdge("proprietary/bin/bin1.meta_lic", "proprietary/lib/liba.so.meta_lic", "static"),
				matchEdge("proprietary/bin/bin1.meta_lic", "proprietary/lib/libc.a.meta_lic", "static"),
			},
		},
		{
			condition:   "proprietary",
			name:        "library",
			roots:       []string{"lib/libd.so.meta_lic"},
			expectedOut: []getMatcher{matchTarget("proprietary/lib/libd.so.meta_lic")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.condition+" "+tt.name, func(t *testing.T) {
			ctx := &testContext{0, make(map[string]string)}

			expectedOut := &bytes.Buffer{}
			for _, eo := range tt.expectedOut {
				m := eo(ctx)
				expectedOut.WriteString(m.matchString(ctx))
				expectedOut.WriteString("\n")
			}

			stdout := &bytes.Buffer{}
			stderr := &bytes.Buffer{}

			rootFiles := make([]string, 0, len(tt.roots))
			for _, r := range tt.roots {
				rootFiles = append(rootFiles, tt.condition+"/"+r)
			}
			tt.ctx.graphViz = true
			err := dumpGraph(&tt.ctx, stdout, stderr, rootFiles...)
			if err != nil {
				t.Fatalf("dumpgraph: error = %v, stderr = %v", err, stderr)
				return
			}
			if stderr.Len() > 0 {
				t.Errorf("dumpgraph: gotStderr = %v, want none", stderr)
			}
			outList := strings.Split(stdout.String(), "\n")
			outLine := 0
			if outList[outLine] != "strict digraph {" {
				t.Errorf("dumpgraph: got 1st line %v, want strict digraph {", outList[outLine])
			}
			outLine++
			if strings.HasPrefix(strings.TrimLeft(outList[outLine], " \t"), "rankdir") {
				outLine++
			}
			endOut := len(outList)
			for endOut > 0 && strings.TrimLeft(outList[endOut-1], " \t") == "" {
				endOut--
			}
			if outList[endOut-1] != "}" {
				t.Errorf("dumpgraph: got last line %v, want }", outList[endOut-1])
			}
			endOut--
			if strings.HasPrefix(strings.TrimLeft(outList[endOut-1], " \t"), "{rank=same") {
				endOut--
			}
			expectedList := strings.Split(expectedOut.String(), "\n")
			for len(expectedList) > 0 && expectedList[len(expectedList)-1] == "" {
				expectedList = expectedList[0 : len(expectedList)-1]
			}
			matchLine := 0

			for outLine < endOut && matchLine < len(expectedList) && strings.TrimLeft(outList[outLine], " \t") == expectedList[matchLine] {
				outLine++
				matchLine++
			}
			if outLine < endOut || matchLine < len(expectedList) {
				if outLine >= endOut {
					t.Errorf("dumpgraph: missing lines at end of graph, want %d lines %v", len(expectedList)-matchLine, strings.Join(expectedList[matchLine:], "\n"))
				} else if matchLine >= len(expectedList) {
					t.Errorf("dumpgraph: unexpected lines at end of graph starting line %d, got %v, want nothing", outLine+1, strings.Join(outList[outLine:], "\n"))
				} else {
					t.Errorf("dumpgraph: at line %d, got %v, want %v", outLine+1, strings.Join(outList[outLine:], "\n"), strings.Join(expectedList[matchLine:], "\n"))
				}
			}
		})
	}
}
