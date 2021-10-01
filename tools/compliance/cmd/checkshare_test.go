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
	"testing"
)

func Test(t *testing.T) {
	tests := []struct {
		condition string
		name      string
		roots     []string
	} {
	    {
		condition: "firstparty",
		name: "apex",
		roots: []string{"highest.apex.meta_lic"},
	    },
	    {
		condition: "firstparty",
		name: "container",
		roots: []string{"container.zip.meta_lic"},
	    },
	    {
		condition: "firstparty",
		name: "application",
		roots: []string{"application.meta_lic"},
	    },
	    {
		condition: "firstparty",
		name: "binary",
		roots: []string{"bin/bin1.meta_lic"},
	    },
	    {
		condition: "firstparty",
		name: "library",
		roots: []string{"lib/libd.so.meta_lic"},
	    },
	    {
		condition: "notice",
		name: "apex",
		roots: []string{"highest.apex.meta_lic"},
	    },
	    {
		condition: "notice",
		name: "container",
		roots: []string{"container.zip.meta_lic"},
	    },
	    {
		condition: "notice",
		name: "application",
		roots: []string{"application.meta_lic"},
	    },
	    {
		condition: "notice",
		name: "binary",
		roots: []string{"bin/bin1.meta_lic"},
	    },
	    {
		condition: "notice",
		name: "library",
		roots: []string{"lib/libd.so.meta_lic"},
	    },
	    {
		condition: "reciprocal",
		name: "apex",
		roots: []string{"highest.apex.meta_lic"},
	    },
	    {
		condition: "reciprocal",
		name: "container",
		roots: []string{"container.zip.meta_lic"},
	    },
	    {
		condition: "reciprocal",
		name: "application",
		roots: []string{"application.meta_lic"},
	    },
	    {
		condition: "reciprocal",
		name: "binary",
		roots: []string{"bin/bin1.meta_lic"},
	    },
	    {
		condition: "reciprocal",
		name: "library",
		roots: []string{"lib/libd.so.meta_lic"},
	    },
	    {
		condition: "restricted",
		name: "apex",
		roots: []string{"highest.apex.meta_lic"},
	    },
	    {
		condition: "restricted",
		name: "container",
		roots: []string{"container.zip.meta_lic"},
	    },
	    {
		condition: "restricted",
		name: "application",
		roots: []string{"application.meta_lic"},
	    },
	    {
		condition: "restricted",
		name: "binary",
		roots: []string{"bin/bin1.meta_lic"},
	    },
	    {
		condition: "restricted",
		name: "library",
		roots: []string{"lib/libd.so.meta_lic"},
	    },
	    {
		condition: "proprietary",
		name: "apex",
		roots: []string{"highest.apex.meta_lic"},
	    },
	    {
		condition: "proprietary",
		name: "container",
		roots: []string{"container.zip.meta_lic"},
	    },
	    {
		condition: "proprietary",
		name: "application",
		roots: []string{"application.meta_lic"},
	    },
	    {
		condition: "proprietary",
		name: "binary",
		roots: []string{"bin/bin1.meta_lic"},
	    },
	    {
		condition: "proprietary",
		name: "library",
		roots: []string{"lib/libd.so.meta_lic"},
	    },
	}
	for _, tt := range tests {
		t.Run(tt.condition + " " + tt.name, func(t *testing.T) {
			stdout := &bytes.Buffer{}
			stderr := &bytes.Buffer{}

			rootFiles := make([]string, 0, len(tt.roots))
			for _, r := range tt.roots {
				rootFiles = append(rootFiles, "testdata/" + tt.condition + "/" + r)
			}
			err := checkShare(stdout, stderr, rootFiles...)
			fmt.Fprintf(os.Stderr, "\nstderr:\n%s\n", stderr)
			fmt.Fprintf(os.Stderr, "\nstdout:\n%s\n", stdout)
			if err == nil {
				fmt.Fprintf(os.Stderr, "%s %s: passed\n", tt.condition, tt.name)
				return
			}
			fmt.Fprintf(os.Stderr, "%s %s: failed, err=%q\n", tt.condition, tt.name, err.Error())
		})
	}
	t.Errorf("completed %d tests", len(tests))
}
