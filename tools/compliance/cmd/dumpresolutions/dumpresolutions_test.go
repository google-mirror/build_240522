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

	"android/soong/tools/compliance"
)

func TestMain(m *testing.M) {
	// Change into the testdata directory before running the tests.
	if err := os.Chdir("../testdata"); err != nil {
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
				"firstparty/bin/bin1.meta_lic firstparty/bin/bin1.meta_lic notice",
				"firstparty/bin/bin1.meta_lic firstparty/lib/liba.so.meta_lic notice",
				"firstparty/bin/bin1.meta_lic firstparty/lib/libc.a.meta_lic notice",
				"firstparty/bin/bin2.meta_lic firstparty/bin/bin2.meta_lic notice",
				"firstparty/highest.apex.meta_lic firstparty/bin/bin1.meta_lic notice",
				"firstparty/highest.apex.meta_lic firstparty/bin/bin2.meta_lic notice",
				"firstparty/highest.apex.meta_lic firstparty/highest.apex.meta_lic notice",
				"firstparty/highest.apex.meta_lic firstparty/lib/liba.so.meta_lic notice",
				"firstparty/highest.apex.meta_lic firstparty/lib/libb.so.meta_lic notice",
				"firstparty/highest.apex.meta_lic firstparty/lib/libc.a.meta_lic notice",
				"firstparty/lib/liba.so.meta_lic firstparty/lib/liba.so.meta_lic notice",
				"firstparty/lib/libb.so.meta_lic firstparty/lib/libb.so.meta_lic notice",
			},
		},
		{
			condition: "firstparty",
			name:      "apex_trimmed",
			roots:     []string{"highest.apex.meta_lic"},
			ctx:       context{stripPrefix: "firstparty/"},
			expectedOut: []string{
				"bin/bin1.meta_lic bin/bin1.meta_lic notice",
				"bin/bin1.meta_lic lib/liba.so.meta_lic notice",
				"bin/bin1.meta_lic lib/libc.a.meta_lic notice",
				"bin/bin2.meta_lic bin/bin2.meta_lic notice",
				"highest.apex.meta_lic bin/bin1.meta_lic notice",
				"highest.apex.meta_lic bin/bin2.meta_lic notice",
				"highest.apex.meta_lic highest.apex.meta_lic notice",
				"highest.apex.meta_lic lib/liba.so.meta_lic notice",
				"highest.apex.meta_lic lib/libb.so.meta_lic notice",
				"highest.apex.meta_lic lib/libc.a.meta_lic notice",
				"lib/liba.so.meta_lic lib/liba.so.meta_lic notice",
				"lib/libb.so.meta_lic lib/libb.so.meta_lic notice",
			},
		},
		{
			condition: "firstparty",
			name:      "apex_trimmed_notice",
			roots:     []string{"highest.apex.meta_lic"},
			ctx: context{
				conditions:  []compliance.LicenseCondition{compliance.NoticeCondition},
				stripPrefix: "firstparty/",
			},
			expectedOut: []string{
				"bin/bin1.meta_lic bin/bin1.meta_lic notice",
				"bin/bin1.meta_lic lib/liba.so.meta_lic notice",
				"bin/bin1.meta_lic lib/libc.a.meta_lic notice",
				"bin/bin2.meta_lic bin/bin2.meta_lic notice",
				"highest.apex.meta_lic bin/bin1.meta_lic notice",
				"highest.apex.meta_lic bin/bin2.meta_lic notice",
				"highest.apex.meta_lic highest.apex.meta_lic notice",
				"highest.apex.meta_lic lib/liba.so.meta_lic notice",
				"highest.apex.meta_lic lib/libb.so.meta_lic notice",
				"highest.apex.meta_lic lib/libc.a.meta_lic notice",
				"lib/liba.so.meta_lic lib/liba.so.meta_lic notice",
				"lib/libb.so.meta_lic lib/libb.so.meta_lic notice",
			},
		},
		{
			condition: "firstparty",
			name:      "apex_trimmed_share",
			roots:     []string{"highest.apex.meta_lic"},
			ctx: context{
				conditions:  compliance.ImpliesShared.AsList(),
				stripPrefix: "firstparty/",
			},
			expectedOut: []string{},
		},
		{
			condition: "firstparty",
			name:      "apex_trimmed_private",
			roots:     []string{"highest.apex.meta_lic"},
			ctx: context{
				conditions:  compliance.ImpliesPrivate.AsList(),
				stripPrefix: "firstparty/",
			},
			expectedOut: []string{},
		},
		{
			condition: "firstparty",
			name:      "apex_trimmed_share_private",
			roots:     []string{"highest.apex.meta_lic"},
			ctx: context{
				conditions:  append(compliance.ImpliesPrivate.AsList(), compliance.ImpliesShared.AsList()...),
				stripPrefix: "firstparty/",
			},
			expectedOut: []string{},
		},
		{
			condition: "firstparty",
			name:      "apex_trimmed_labelled",
			roots:     []string{"highest.apex.meta_lic"},
			ctx:       context{stripPrefix: "firstparty/", labelConditions: true},
			expectedOut: []string{
				"bin/bin1.meta_lic:notice bin/bin1.meta_lic:notice notice",
				"bin/bin1.meta_lic:notice lib/liba.so.meta_lic:notice notice",
				"bin/bin1.meta_lic:notice lib/libc.a.meta_lic:notice notice",
				"bin/bin2.meta_lic:notice bin/bin2.meta_lic:notice notice",
				"highest.apex.meta_lic:notice bin/bin1.meta_lic:notice notice",
				"highest.apex.meta_lic:notice bin/bin2.meta_lic:notice notice",
				"highest.apex.meta_lic:notice highest.apex.meta_lic:notice notice",
				"highest.apex.meta_lic:notice lib/liba.so.meta_lic:notice notice",
				"highest.apex.meta_lic:notice lib/libb.so.meta_lic:notice notice",
				"highest.apex.meta_lic:notice lib/libc.a.meta_lic:notice notice",
				"lib/liba.so.meta_lic:notice lib/liba.so.meta_lic:notice notice",
				"lib/libb.so.meta_lic:notice lib/libb.so.meta_lic:notice notice",
			},
		},
		{
			condition: "firstparty",
			name:      "container",
			roots:     []string{"container.zip.meta_lic"},
			expectedOut: []string{
				"firstparty/bin/bin1.meta_lic firstparty/bin/bin1.meta_lic notice",
				"firstparty/bin/bin1.meta_lic firstparty/lib/liba.so.meta_lic notice",
				"firstparty/bin/bin1.meta_lic firstparty/lib/libc.a.meta_lic notice",
				"firstparty/bin/bin2.meta_lic firstparty/bin/bin2.meta_lic notice",
				"firstparty/container.zip.meta_lic firstparty/bin/bin1.meta_lic notice",
				"firstparty/container.zip.meta_lic firstparty/bin/bin2.meta_lic notice",
				"firstparty/container.zip.meta_lic firstparty/container.zip.meta_lic notice",
				"firstparty/container.zip.meta_lic firstparty/lib/liba.so.meta_lic notice",
				"firstparty/container.zip.meta_lic firstparty/lib/libb.so.meta_lic notice",
				"firstparty/container.zip.meta_lic firstparty/lib/libc.a.meta_lic notice",
				"firstparty/lib/liba.so.meta_lic firstparty/lib/liba.so.meta_lic notice",
				"firstparty/lib/libb.so.meta_lic firstparty/lib/libb.so.meta_lic notice",
			},
		},
		{
			condition: "firstparty",
			name:      "application",
			roots:     []string{"application.meta_lic"},
			expectedOut: []string{
				"firstparty/application.meta_lic firstparty/application.meta_lic notice",
				"firstparty/application.meta_lic firstparty/lib/liba.so.meta_lic notice",
			},
		},
		{
			condition: "firstparty",
			name:      "binary",
			roots:     []string{"bin/bin1.meta_lic"},
			expectedOut: []string{
				"firstparty/bin/bin1.meta_lic firstparty/bin/bin1.meta_lic notice",
				"firstparty/bin/bin1.meta_lic firstparty/lib/liba.so.meta_lic notice",
				"firstparty/bin/bin1.meta_lic firstparty/lib/libc.a.meta_lic notice",
			},
		},
		{
			condition: "firstparty",
			name:      "library",
			roots:     []string{"lib/libd.so.meta_lic"},
			expectedOut: []string{
				"firstparty/lib/libd.so.meta_lic firstparty/lib/libd.so.meta_lic notice",
			},
		},
		{
			condition: "notice",
			name:      "apex",
			roots:     []string{"highest.apex.meta_lic"},
			expectedOut: []string{
				"notice/bin/bin1.meta_lic notice/bin/bin1.meta_lic notice",
				"notice/bin/bin1.meta_lic notice/lib/liba.so.meta_lic notice",
				"notice/bin/bin1.meta_lic notice/lib/libc.a.meta_lic notice",
				"notice/bin/bin2.meta_lic notice/bin/bin2.meta_lic notice",
				"notice/highest.apex.meta_lic notice/bin/bin1.meta_lic notice",
				"notice/highest.apex.meta_lic notice/bin/bin2.meta_lic notice",
				"notice/highest.apex.meta_lic notice/highest.apex.meta_lic notice",
				"notice/highest.apex.meta_lic notice/lib/liba.so.meta_lic notice",
				"notice/highest.apex.meta_lic notice/lib/libb.so.meta_lic notice",
				"notice/highest.apex.meta_lic notice/lib/libc.a.meta_lic notice",
				"notice/lib/liba.so.meta_lic notice/lib/liba.so.meta_lic notice",
				"notice/lib/libb.so.meta_lic notice/lib/libb.so.meta_lic notice",
			},
		},
		{
			condition: "notice",
			name:      "apex_trimmed",
			roots:     []string{"highest.apex.meta_lic"},
			ctx:       context{stripPrefix: "notice/"},
			expectedOut: []string{
				"bin/bin1.meta_lic bin/bin1.meta_lic notice",
				"bin/bin1.meta_lic lib/liba.so.meta_lic notice",
				"bin/bin1.meta_lic lib/libc.a.meta_lic notice",
				"bin/bin2.meta_lic bin/bin2.meta_lic notice",
				"highest.apex.meta_lic bin/bin1.meta_lic notice",
				"highest.apex.meta_lic bin/bin2.meta_lic notice",
				"highest.apex.meta_lic highest.apex.meta_lic notice",
				"highest.apex.meta_lic lib/liba.so.meta_lic notice",
				"highest.apex.meta_lic lib/libb.so.meta_lic notice",
				"highest.apex.meta_lic lib/libc.a.meta_lic notice",
				"lib/liba.so.meta_lic lib/liba.so.meta_lic notice",
				"lib/libb.so.meta_lic lib/libb.so.meta_lic notice",
			},
		},
		{
			condition: "notice",
			name:      "apex_trimmed_notice",
			roots:     []string{"highest.apex.meta_lic"},
			ctx: context{
				conditions:  []compliance.LicenseCondition{compliance.NoticeCondition},
				stripPrefix: "notice/",
			},
			expectedOut: []string{
				"bin/bin1.meta_lic bin/bin1.meta_lic notice",
				"bin/bin1.meta_lic lib/liba.so.meta_lic notice",
				"bin/bin1.meta_lic lib/libc.a.meta_lic notice",
				"bin/bin2.meta_lic bin/bin2.meta_lic notice",
				"highest.apex.meta_lic bin/bin1.meta_lic notice",
				"highest.apex.meta_lic bin/bin2.meta_lic notice",
				"highest.apex.meta_lic highest.apex.meta_lic notice",
				"highest.apex.meta_lic lib/liba.so.meta_lic notice",
				"highest.apex.meta_lic lib/libb.so.meta_lic notice",
				"highest.apex.meta_lic lib/libc.a.meta_lic notice",
				"lib/liba.so.meta_lic lib/liba.so.meta_lic notice",
				"lib/libb.so.meta_lic lib/libb.so.meta_lic notice",
			},
		},
		{
			condition: "notice",
			name:      "apex_trimmed_share",
			roots:     []string{"highest.apex.meta_lic"},
			ctx: context{
				conditions:  compliance.ImpliesShared.AsList(),
				stripPrefix: "notice/",
			},
			expectedOut: []string{},
		},
		{
			condition: "notice",
			name:      "apex_trimmed_private",
			roots:     []string{"highest.apex.meta_lic"},
			ctx: context{
				conditions:  compliance.ImpliesPrivate.AsList(),
				stripPrefix: "notice/",
			},
			expectedOut: []string{},
		},
		{
			condition: "notice",
			name:      "apex_trimmed_share_private",
			roots:     []string{"highest.apex.meta_lic"},
			ctx: context{
				conditions:  append(compliance.ImpliesShared.AsList(), compliance.ImpliesPrivate.AsList()...),
				stripPrefix: "notice/",
			},
			expectedOut: []string{},
		},
		{
			condition: "notice",
			name:      "apex_trimmed_labelled",
			roots:     []string{"highest.apex.meta_lic"},
			ctx:       context{stripPrefix: "notice/", labelConditions: true},
			expectedOut: []string{
				"bin/bin1.meta_lic:notice bin/bin1.meta_lic:notice notice",
				"bin/bin1.meta_lic:notice lib/liba.so.meta_lic:notice notice",
				"bin/bin1.meta_lic:notice lib/libc.a.meta_lic:notice notice",
				"bin/bin2.meta_lic:notice bin/bin2.meta_lic:notice notice",
				"highest.apex.meta_lic:notice bin/bin1.meta_lic:notice notice",
				"highest.apex.meta_lic:notice bin/bin2.meta_lic:notice notice",
				"highest.apex.meta_lic:notice highest.apex.meta_lic:notice notice",
				"highest.apex.meta_lic:notice lib/liba.so.meta_lic:notice notice",
				"highest.apex.meta_lic:notice lib/libb.so.meta_lic:notice notice",
				"highest.apex.meta_lic:notice lib/libc.a.meta_lic:notice notice",
				"lib/liba.so.meta_lic:notice lib/liba.so.meta_lic:notice notice",
				"lib/libb.so.meta_lic:notice lib/libb.so.meta_lic:notice notice",
			},
		},
		{
			condition: "notice",
			name:      "container",
			roots:     []string{"container.zip.meta_lic"},
			expectedOut: []string{
				"notice/bin/bin1.meta_lic notice/bin/bin1.meta_lic notice",
				"notice/bin/bin1.meta_lic notice/lib/liba.so.meta_lic notice",
				"notice/bin/bin1.meta_lic notice/lib/libc.a.meta_lic notice",
				"notice/bin/bin2.meta_lic notice/bin/bin2.meta_lic notice",
				"notice/container.zip.meta_lic notice/bin/bin1.meta_lic notice",
				"notice/container.zip.meta_lic notice/bin/bin2.meta_lic notice",
				"notice/container.zip.meta_lic notice/container.zip.meta_lic notice",
				"notice/container.zip.meta_lic notice/lib/liba.so.meta_lic notice",
				"notice/container.zip.meta_lic notice/lib/libb.so.meta_lic notice",
				"notice/container.zip.meta_lic notice/lib/libc.a.meta_lic notice",
				"notice/lib/liba.so.meta_lic notice/lib/liba.so.meta_lic notice",
				"notice/lib/libb.so.meta_lic notice/lib/libb.so.meta_lic notice",
			},
		},
		{
			condition: "notice",
			name:      "application",
			roots:     []string{"application.meta_lic"},
			expectedOut: []string{
				"notice/application.meta_lic notice/application.meta_lic notice",
				"notice/application.meta_lic notice/lib/liba.so.meta_lic notice",
			},
		},
		{
			condition: "notice",
			name:      "binary",
			roots:     []string{"bin/bin1.meta_lic"},
			expectedOut: []string{
				"notice/bin/bin1.meta_lic notice/bin/bin1.meta_lic notice",
				"notice/bin/bin1.meta_lic notice/lib/liba.so.meta_lic notice",
				"notice/bin/bin1.meta_lic notice/lib/libc.a.meta_lic notice",
			},
		},
		{
			condition: "notice",
			name:      "library",
			roots:     []string{"lib/libd.so.meta_lic"},
			expectedOut: []string{
				"notice/lib/libd.so.meta_lic notice/lib/libd.so.meta_lic notice",
			},
		},
		{
			condition: "reciprocal",
			name:      "apex",
			roots:     []string{"highest.apex.meta_lic"},
			expectedOut: []string{
				"reciprocal/bin/bin1.meta_lic reciprocal/bin/bin1.meta_lic notice",
				"reciprocal/bin/bin1.meta_lic reciprocal/lib/liba.so.meta_lic reciprocal",
				"reciprocal/bin/bin1.meta_lic reciprocal/lib/libc.a.meta_lic reciprocal",
				"reciprocal/bin/bin2.meta_lic reciprocal/bin/bin2.meta_lic notice",
				"reciprocal/highest.apex.meta_lic reciprocal/bin/bin1.meta_lic notice",
				"reciprocal/highest.apex.meta_lic reciprocal/bin/bin2.meta_lic notice",
				"reciprocal/highest.apex.meta_lic reciprocal/highest.apex.meta_lic notice",
				"reciprocal/highest.apex.meta_lic reciprocal/lib/liba.so.meta_lic reciprocal",
				"reciprocal/highest.apex.meta_lic reciprocal/lib/libb.so.meta_lic notice",
				"reciprocal/highest.apex.meta_lic reciprocal/lib/libc.a.meta_lic reciprocal",
				"reciprocal/lib/liba.so.meta_lic reciprocal/lib/liba.so.meta_lic reciprocal",
				"reciprocal/lib/libb.so.meta_lic reciprocal/lib/libb.so.meta_lic notice",
			},
		},
		{
			condition: "reciprocal",
			name:      "apex_trimmed",
			roots:     []string{"highest.apex.meta_lic"},
			ctx:       context{stripPrefix: "reciprocal/"},
			expectedOut: []string{
				"bin/bin1.meta_lic bin/bin1.meta_lic notice",
				"bin/bin1.meta_lic lib/liba.so.meta_lic reciprocal",
				"bin/bin1.meta_lic lib/libc.a.meta_lic reciprocal",
				"bin/bin2.meta_lic bin/bin2.meta_lic notice",
				"highest.apex.meta_lic bin/bin1.meta_lic notice",
				"highest.apex.meta_lic bin/bin2.meta_lic notice",
				"highest.apex.meta_lic highest.apex.meta_lic notice",
				"highest.apex.meta_lic lib/liba.so.meta_lic reciprocal",
				"highest.apex.meta_lic lib/libb.so.meta_lic notice",
				"highest.apex.meta_lic lib/libc.a.meta_lic reciprocal",
				"lib/liba.so.meta_lic lib/liba.so.meta_lic reciprocal",
				"lib/libb.so.meta_lic lib/libb.so.meta_lic notice",
			},
		},
		{
			condition: "reciprocal",
			name:      "apex_trimmed_notice",
			roots:     []string{"highest.apex.meta_lic"},
			ctx: context{
				conditions:  []compliance.LicenseCondition{compliance.NoticeCondition},
				stripPrefix: "reciprocal/",
			},
			expectedOut: []string{
				"bin/bin1.meta_lic bin/bin1.meta_lic notice",
				"bin/bin2.meta_lic bin/bin2.meta_lic notice",
				"highest.apex.meta_lic bin/bin1.meta_lic notice",
				"highest.apex.meta_lic bin/bin2.meta_lic notice",
				"highest.apex.meta_lic highest.apex.meta_lic notice",
				"highest.apex.meta_lic lib/libb.so.meta_lic notice",
				"lib/libb.so.meta_lic lib/libb.so.meta_lic notice",
			},
		},
		{
			condition: "reciprocal",
			name:      "apex_trimmed_share",
			roots:     []string{"highest.apex.meta_lic"},
			ctx: context{
				conditions:  compliance.ImpliesShared.AsList(),
				stripPrefix: "reciprocal/",
			},
			expectedOut: []string{
				"bin/bin1.meta_lic lib/liba.so.meta_lic reciprocal",
				"bin/bin1.meta_lic lib/libc.a.meta_lic reciprocal",
				"highest.apex.meta_lic lib/liba.so.meta_lic reciprocal",
				"highest.apex.meta_lic lib/libc.a.meta_lic reciprocal",
				"lib/liba.so.meta_lic lib/liba.so.meta_lic reciprocal",
			},
		},
		{
			condition: "reciprocal",
			name:      "apex_trimmed_private",
			roots:     []string{"highest.apex.meta_lic"},
			ctx: context{
				conditions:  compliance.ImpliesPrivate.AsList(),
				stripPrefix: "reciprocal/",
			},
			expectedOut: []string{},
		},
		{
			condition: "reciprocal",
			name:      "apex_trimmed_share_private",
			roots:     []string{"highest.apex.meta_lic"},
			ctx: context{
				conditions:  append(compliance.ImpliesShared.AsList(), compliance.ImpliesPrivate.AsList()...),
				stripPrefix: "reciprocal/",
			},
			expectedOut: []string{
				"bin/bin1.meta_lic lib/liba.so.meta_lic reciprocal",
				"bin/bin1.meta_lic lib/libc.a.meta_lic reciprocal",
				"highest.apex.meta_lic lib/liba.so.meta_lic reciprocal",
				"highest.apex.meta_lic lib/libc.a.meta_lic reciprocal",
				"lib/liba.so.meta_lic lib/liba.so.meta_lic reciprocal",
			},
		},
		{
			condition: "reciprocal",
			name:      "apex_trimmed_labelled",
			roots:     []string{"highest.apex.meta_lic"},
			ctx:       context{stripPrefix: "reciprocal/", labelConditions: true},
			expectedOut: []string{
				"bin/bin1.meta_lic:notice bin/bin1.meta_lic:notice notice",
				"bin/bin1.meta_lic:notice lib/liba.so.meta_lic:reciprocal reciprocal",
				"bin/bin1.meta_lic:notice lib/libc.a.meta_lic:reciprocal reciprocal",
				"bin/bin2.meta_lic:notice bin/bin2.meta_lic:notice notice",
				"highest.apex.meta_lic:notice bin/bin1.meta_lic:notice notice",
				"highest.apex.meta_lic:notice bin/bin2.meta_lic:notice notice",
				"highest.apex.meta_lic:notice highest.apex.meta_lic:notice notice",
				"highest.apex.meta_lic:notice lib/liba.so.meta_lic:reciprocal reciprocal",
				"highest.apex.meta_lic:notice lib/libb.so.meta_lic:notice notice",
				"highest.apex.meta_lic:notice lib/libc.a.meta_lic:reciprocal reciprocal",
				"lib/liba.so.meta_lic:reciprocal lib/liba.so.meta_lic:reciprocal reciprocal",
				"lib/libb.so.meta_lic:notice lib/libb.so.meta_lic:notice notice",
			},
		},
		{
			condition: "reciprocal",
			name:      "container",
			roots:     []string{"container.zip.meta_lic"},
			expectedOut: []string{
				"reciprocal/bin/bin1.meta_lic reciprocal/bin/bin1.meta_lic notice",
				"reciprocal/bin/bin1.meta_lic reciprocal/lib/liba.so.meta_lic reciprocal",
				"reciprocal/bin/bin1.meta_lic reciprocal/lib/libc.a.meta_lic reciprocal",
				"reciprocal/bin/bin2.meta_lic reciprocal/bin/bin2.meta_lic notice",
				"reciprocal/container.zip.meta_lic reciprocal/bin/bin1.meta_lic notice",
				"reciprocal/container.zip.meta_lic reciprocal/bin/bin2.meta_lic notice",
				"reciprocal/container.zip.meta_lic reciprocal/container.zip.meta_lic notice",
				"reciprocal/container.zip.meta_lic reciprocal/lib/liba.so.meta_lic reciprocal",
				"reciprocal/container.zip.meta_lic reciprocal/lib/libb.so.meta_lic notice",
				"reciprocal/container.zip.meta_lic reciprocal/lib/libc.a.meta_lic reciprocal",
				"reciprocal/lib/liba.so.meta_lic reciprocal/lib/liba.so.meta_lic reciprocal",
				"reciprocal/lib/libb.so.meta_lic reciprocal/lib/libb.so.meta_lic notice",
			},
		},
		{
			condition: "reciprocal",
			name:      "application",
			roots:     []string{"application.meta_lic"},
			expectedOut: []string{
				"reciprocal/application.meta_lic reciprocal/application.meta_lic notice",
				"reciprocal/application.meta_lic reciprocal/lib/liba.so.meta_lic reciprocal",
			},
		},
		{
			condition: "reciprocal",
			name:      "binary",
			roots:     []string{"bin/bin1.meta_lic"},
			expectedOut: []string{
				"reciprocal/bin/bin1.meta_lic reciprocal/bin/bin1.meta_lic notice",
				"reciprocal/bin/bin1.meta_lic reciprocal/lib/liba.so.meta_lic reciprocal",
				"reciprocal/bin/bin1.meta_lic reciprocal/lib/libc.a.meta_lic reciprocal",
			},
		},
		{
			condition: "reciprocal",
			name:      "library",
			roots:     []string{"lib/libd.so.meta_lic"},
			expectedOut: []string{
				"reciprocal/lib/libd.so.meta_lic reciprocal/lib/libd.so.meta_lic notice",
			},
		},
		{
			condition: "restricted",
			name:      "apex",
			roots:     []string{"highest.apex.meta_lic"},
			expectedOut: []string{
				"restricted/bin/bin1.meta_lic restricted/bin/bin1.meta_lic notice:restricted_allows_dynamic_linking",
				"restricted/bin/bin1.meta_lic restricted/lib/liba.so.meta_lic restricted_allows_dynamic_linking",
				"restricted/bin/bin1.meta_lic restricted/lib/libc.a.meta_lic reciprocal:restricted_allows_dynamic_linking",
				"restricted/bin/bin2.meta_lic restricted/bin/bin2.meta_lic notice:restricted",
				"restricted/bin/bin2.meta_lic restricted/lib/libb.so.meta_lic restricted",
				"restricted/highest.apex.meta_lic restricted/bin/bin1.meta_lic notice:restricted_allows_dynamic_linking",
				"restricted/highest.apex.meta_lic restricted/bin/bin2.meta_lic notice:restricted",
				"restricted/highest.apex.meta_lic restricted/highest.apex.meta_lic notice:restricted:restricted_allows_dynamic_linking",
				"restricted/highest.apex.meta_lic restricted/lib/liba.so.meta_lic restricted_allows_dynamic_linking",
				"restricted/highest.apex.meta_lic restricted/lib/libb.so.meta_lic restricted",
				"restricted/highest.apex.meta_lic restricted/lib/libc.a.meta_lic reciprocal:restricted_allows_dynamic_linking",
				"restricted/lib/liba.so.meta_lic restricted/lib/liba.so.meta_lic restricted_allows_dynamic_linking",
				"restricted/lib/libb.so.meta_lic restricted/lib/libb.so.meta_lic restricted",
			},
		},
		{
			condition: "restricted",
			name:      "apex_trimmed",
			roots:     []string{"highest.apex.meta_lic"},
			ctx:       context{stripPrefix: "restricted/"},
			expectedOut: []string{
				"bin/bin1.meta_lic bin/bin1.meta_lic notice:restricted_allows_dynamic_linking",
				"bin/bin1.meta_lic lib/liba.so.meta_lic restricted_allows_dynamic_linking",
				"bin/bin1.meta_lic lib/libc.a.meta_lic reciprocal:restricted_allows_dynamic_linking",
				"bin/bin2.meta_lic bin/bin2.meta_lic notice:restricted",
				"bin/bin2.meta_lic lib/libb.so.meta_lic restricted",
				"highest.apex.meta_lic bin/bin1.meta_lic notice:restricted_allows_dynamic_linking",
				"highest.apex.meta_lic bin/bin2.meta_lic notice:restricted",
				"highest.apex.meta_lic highest.apex.meta_lic notice:restricted:restricted_allows_dynamic_linking",
				"highest.apex.meta_lic lib/liba.so.meta_lic restricted_allows_dynamic_linking",
				"highest.apex.meta_lic lib/libb.so.meta_lic restricted",
				"highest.apex.meta_lic lib/libc.a.meta_lic reciprocal:restricted_allows_dynamic_linking",
				"lib/liba.so.meta_lic lib/liba.so.meta_lic restricted_allows_dynamic_linking",
				"lib/libb.so.meta_lic lib/libb.so.meta_lic restricted",
			},
		},
		{
			condition: "restricted",
			name:      "apex_trimmed_notice",
			roots:     []string{"highest.apex.meta_lic"},
			ctx: context{
				conditions:  []compliance.LicenseCondition{compliance.NoticeCondition},
				stripPrefix: "restricted/",
			},
			expectedOut: []string{
				"bin/bin1.meta_lic bin/bin1.meta_lic notice",
				"bin/bin2.meta_lic bin/bin2.meta_lic notice",
				"highest.apex.meta_lic bin/bin1.meta_lic notice",
				"highest.apex.meta_lic bin/bin2.meta_lic notice",
				"highest.apex.meta_lic highest.apex.meta_lic notice",
			},
		},
		{
			condition: "restricted",
			name:      "apex_trimmed_share",
			roots:     []string{"highest.apex.meta_lic"},
			ctx: context{
				conditions:  compliance.ImpliesShared.AsList(),
				stripPrefix: "restricted/",
			},
			expectedOut: []string{
				"bin/bin1.meta_lic bin/bin1.meta_lic restricted_allows_dynamic_linking",
				"bin/bin1.meta_lic lib/liba.so.meta_lic restricted_allows_dynamic_linking",
				"bin/bin1.meta_lic lib/libc.a.meta_lic reciprocal:restricted_allows_dynamic_linking",
				"bin/bin2.meta_lic bin/bin2.meta_lic restricted",
				"bin/bin2.meta_lic lib/libb.so.meta_lic restricted",
				"highest.apex.meta_lic bin/bin1.meta_lic restricted_allows_dynamic_linking",
				"highest.apex.meta_lic bin/bin2.meta_lic restricted",
				"highest.apex.meta_lic highest.apex.meta_lic restricted:restricted_allows_dynamic_linking",
				"highest.apex.meta_lic lib/liba.so.meta_lic restricted_allows_dynamic_linking",
				"highest.apex.meta_lic lib/libb.so.meta_lic restricted",
				"highest.apex.meta_lic lib/libc.a.meta_lic reciprocal:restricted_allows_dynamic_linking",
				"lib/liba.so.meta_lic lib/liba.so.meta_lic restricted_allows_dynamic_linking",
				"lib/libb.so.meta_lic lib/libb.so.meta_lic restricted",
			},
		},
		{
			condition: "restricted",
			name:      "apex_trimmed_private",
			roots:     []string{"highest.apex.meta_lic"},
			ctx: context{
				conditions:  compliance.ImpliesPrivate.AsList(),
				stripPrefix: "restricted/",
			},
			expectedOut: []string{},
		},
		{
			condition: "restricted",
			name:      "apex_trimmed_share_private",
			roots:     []string{"highest.apex.meta_lic"},
			ctx: context{
				conditions:  append(compliance.ImpliesShared.AsList(), compliance.ImpliesPrivate.AsList()...),
				stripPrefix: "restricted/",
			},
			expectedOut: []string{
				"bin/bin1.meta_lic bin/bin1.meta_lic restricted_allows_dynamic_linking",
				"bin/bin1.meta_lic lib/liba.so.meta_lic restricted_allows_dynamic_linking",
				"bin/bin1.meta_lic lib/libc.a.meta_lic reciprocal:restricted_allows_dynamic_linking",
				"bin/bin2.meta_lic bin/bin2.meta_lic restricted",
				"bin/bin2.meta_lic lib/libb.so.meta_lic restricted",
				"highest.apex.meta_lic bin/bin1.meta_lic restricted_allows_dynamic_linking",
				"highest.apex.meta_lic bin/bin2.meta_lic restricted",
				"highest.apex.meta_lic highest.apex.meta_lic restricted:restricted_allows_dynamic_linking",
				"highest.apex.meta_lic lib/liba.so.meta_lic restricted_allows_dynamic_linking",
				"highest.apex.meta_lic lib/libb.so.meta_lic restricted",
				"highest.apex.meta_lic lib/libc.a.meta_lic reciprocal:restricted_allows_dynamic_linking",
				"lib/liba.so.meta_lic lib/liba.so.meta_lic restricted_allows_dynamic_linking",
				"lib/libb.so.meta_lic lib/libb.so.meta_lic restricted",
			},
		},
		{
			condition: "restricted",
			name:      "apex_trimmed_labelled",
			roots:     []string{"highest.apex.meta_lic"},
			ctx:       context{stripPrefix: "restricted/", labelConditions: true},
			expectedOut: []string{
				"bin/bin1.meta_lic:notice bin/bin1.meta_lic:notice notice:restricted_allows_dynamic_linking",
				"bin/bin1.meta_lic:notice lib/liba.so.meta_lic:restricted_allows_dynamic_linking restricted_allows_dynamic_linking",
				"bin/bin1.meta_lic:notice lib/libc.a.meta_lic:reciprocal reciprocal:restricted_allows_dynamic_linking",
				"bin/bin2.meta_lic:notice bin/bin2.meta_lic:notice notice:restricted",
				"bin/bin2.meta_lic:notice lib/libb.so.meta_lic:restricted restricted",
				"highest.apex.meta_lic:notice bin/bin1.meta_lic:notice notice:restricted_allows_dynamic_linking",
				"highest.apex.meta_lic:notice bin/bin2.meta_lic:notice notice:restricted",
				"highest.apex.meta_lic:notice highest.apex.meta_lic:notice notice:restricted:restricted_allows_dynamic_linking",
				"highest.apex.meta_lic:notice lib/liba.so.meta_lic:restricted_allows_dynamic_linking restricted_allows_dynamic_linking",
				"highest.apex.meta_lic:notice lib/libb.so.meta_lic:restricted restricted",
				"highest.apex.meta_lic:notice lib/libc.a.meta_lic:reciprocal reciprocal:restricted_allows_dynamic_linking",
				"lib/liba.so.meta_lic:restricted_allows_dynamic_linking lib/liba.so.meta_lic:restricted_allows_dynamic_linking restricted_allows_dynamic_linking",
				"lib/libb.so.meta_lic:restricted lib/libb.so.meta_lic:restricted restricted",
			},
		},
		{
			condition: "restricted",
			name:      "container",
			roots:     []string{"container.zip.meta_lic"},
			expectedOut: []string{
				"restricted/bin/bin1.meta_lic restricted/bin/bin1.meta_lic notice:restricted_allows_dynamic_linking",
				"restricted/bin/bin1.meta_lic restricted/lib/liba.so.meta_lic restricted_allows_dynamic_linking",
				"restricted/bin/bin1.meta_lic restricted/lib/libc.a.meta_lic reciprocal:restricted_allows_dynamic_linking",
				"restricted/bin/bin2.meta_lic restricted/bin/bin2.meta_lic notice:restricted",
				"restricted/bin/bin2.meta_lic restricted/lib/libb.so.meta_lic restricted",
				"restricted/container.zip.meta_lic restricted/bin/bin1.meta_lic notice:restricted_allows_dynamic_linking",
				"restricted/container.zip.meta_lic restricted/bin/bin2.meta_lic notice:restricted",
				"restricted/container.zip.meta_lic restricted/container.zip.meta_lic notice:restricted:restricted_allows_dynamic_linking",
				"restricted/container.zip.meta_lic restricted/lib/liba.so.meta_lic restricted_allows_dynamic_linking",
				"restricted/container.zip.meta_lic restricted/lib/libb.so.meta_lic restricted",
				"restricted/container.zip.meta_lic restricted/lib/libc.a.meta_lic reciprocal:restricted_allows_dynamic_linking",
				"restricted/lib/liba.so.meta_lic restricted/lib/liba.so.meta_lic restricted_allows_dynamic_linking",
				"restricted/lib/libb.so.meta_lic restricted/lib/libb.so.meta_lic restricted",
			},
		},
		{
			condition: "restricted",
			name:      "application",
			roots:     []string{"application.meta_lic"},
			expectedOut: []string{
				"restricted/application.meta_lic restricted/application.meta_lic notice:restricted:restricted_allows_dynamic_linking",
				"restricted/application.meta_lic restricted/lib/liba.so.meta_lic restricted:restricted_allows_dynamic_linking",
			},
		},
		{
			condition: "restricted",
			name:      "binary",
			roots:     []string{"bin/bin1.meta_lic"},
			expectedOut: []string{
				"restricted/bin/bin1.meta_lic restricted/bin/bin1.meta_lic notice:restricted_allows_dynamic_linking",
				"restricted/bin/bin1.meta_lic restricted/lib/liba.so.meta_lic restricted_allows_dynamic_linking",
				"restricted/bin/bin1.meta_lic restricted/lib/libc.a.meta_lic reciprocal:restricted_allows_dynamic_linking",
			},
		},
		{
			condition: "restricted",
			name:      "library",
			roots:     []string{"lib/libd.so.meta_lic"},
			expectedOut: []string{
				"restricted/lib/libd.so.meta_lic restricted/lib/libd.so.meta_lic notice",
			},
		},
		{
			condition: "proprietary",
			name:      "apex",
			roots:     []string{"highest.apex.meta_lic"},
			expectedOut: []string{
				"proprietary/bin/bin1.meta_lic proprietary/bin/bin1.meta_lic notice",
				"proprietary/bin/bin1.meta_lic proprietary/lib/liba.so.meta_lic proprietary:by_exception_only",
				"proprietary/bin/bin1.meta_lic proprietary/lib/libc.a.meta_lic proprietary:by_exception_only",
				"proprietary/bin/bin2.meta_lic proprietary/bin/bin2.meta_lic restricted:proprietary:by_exception_only",
				"proprietary/bin/bin2.meta_lic proprietary/lib/libb.so.meta_lic restricted",
				"proprietary/highest.apex.meta_lic proprietary/bin/bin1.meta_lic notice",
				"proprietary/highest.apex.meta_lic proprietary/bin/bin2.meta_lic restricted:proprietary:by_exception_only",
				"proprietary/highest.apex.meta_lic proprietary/highest.apex.meta_lic notice:restricted",
				"proprietary/highest.apex.meta_lic proprietary/lib/liba.so.meta_lic proprietary:by_exception_only",
				"proprietary/highest.apex.meta_lic proprietary/lib/libb.so.meta_lic restricted",
				"proprietary/highest.apex.meta_lic proprietary/lib/libc.a.meta_lic proprietary:by_exception_only",
				"proprietary/lib/liba.so.meta_lic proprietary/lib/liba.so.meta_lic proprietary:by_exception_only",
				"proprietary/lib/libb.so.meta_lic proprietary/lib/libb.so.meta_lic restricted",
			},
		},
		{
			condition: "proprietary",
			name:      "apex_trimmed",
			roots:     []string{"highest.apex.meta_lic"},
			ctx:       context{stripPrefix: "proprietary/"},
			expectedOut: []string{
				"bin/bin1.meta_lic bin/bin1.meta_lic notice",
				"bin/bin1.meta_lic lib/liba.so.meta_lic proprietary:by_exception_only",
				"bin/bin1.meta_lic lib/libc.a.meta_lic proprietary:by_exception_only",
				"bin/bin2.meta_lic bin/bin2.meta_lic restricted:proprietary:by_exception_only",
				"bin/bin2.meta_lic lib/libb.so.meta_lic restricted",
				"highest.apex.meta_lic bin/bin1.meta_lic notice",
				"highest.apex.meta_lic bin/bin2.meta_lic restricted:proprietary:by_exception_only",
				"highest.apex.meta_lic highest.apex.meta_lic notice:restricted",
				"highest.apex.meta_lic lib/liba.so.meta_lic proprietary:by_exception_only",
				"highest.apex.meta_lic lib/libb.so.meta_lic restricted",
				"highest.apex.meta_lic lib/libc.a.meta_lic proprietary:by_exception_only",
				"lib/liba.so.meta_lic lib/liba.so.meta_lic proprietary:by_exception_only",
				"lib/libb.so.meta_lic lib/libb.so.meta_lic restricted",
			},
		},
		{
			condition: "proprietary",
			name:      "apex_trimmed_notice",
			roots:     []string{"highest.apex.meta_lic"},
			ctx: context{
				conditions:  []compliance.LicenseCondition{compliance.NoticeCondition},
				stripPrefix: "proprietary/",
			},
			expectedOut: []string{
				"bin/bin1.meta_lic bin/bin1.meta_lic notice",
				"highest.apex.meta_lic bin/bin1.meta_lic notice",
				"highest.apex.meta_lic highest.apex.meta_lic notice",
			},
		},
		{
			condition: "proprietary",
			name:      "apex_trimmed_share",
			roots:     []string{"highest.apex.meta_lic"},
			ctx: context{
				conditions:  compliance.ImpliesShared.AsList(),
				stripPrefix: "proprietary/",
			},
			expectedOut: []string{
				"bin/bin2.meta_lic bin/bin2.meta_lic restricted",
				"bin/bin2.meta_lic lib/libb.so.meta_lic restricted",
				"highest.apex.meta_lic bin/bin2.meta_lic restricted",
				"highest.apex.meta_lic highest.apex.meta_lic restricted",
				"highest.apex.meta_lic lib/libb.so.meta_lic restricted",
				"lib/libb.so.meta_lic lib/libb.so.meta_lic restricted",
			},
		},
		{
			condition: "proprietary",
			name:      "apex_trimmed_private",
			roots:     []string{"highest.apex.meta_lic"},
			ctx: context{
				conditions:  compliance.ImpliesPrivate.AsList(),
				stripPrefix: "proprietary/",
			},
			expectedOut: []string{
				"bin/bin1.meta_lic lib/liba.so.meta_lic proprietary",
				"bin/bin1.meta_lic lib/libc.a.meta_lic proprietary",
				"bin/bin2.meta_lic bin/bin2.meta_lic proprietary",
				"highest.apex.meta_lic bin/bin2.meta_lic proprietary",
				"highest.apex.meta_lic lib/liba.so.meta_lic proprietary",
				"highest.apex.meta_lic lib/libc.a.meta_lic proprietary",
				"lib/liba.so.meta_lic lib/liba.so.meta_lic proprietary",
			},
		},
		{
			condition: "proprietary",
			name:      "apex_trimmed_share_private",
			roots:     []string{"highest.apex.meta_lic"},
			ctx: context{
				conditions:  append(compliance.ImpliesShared.AsList(), compliance.ImpliesPrivate.AsList()...),
				stripPrefix: "proprietary/",
			},
			expectedOut: []string{
				"bin/bin1.meta_lic lib/liba.so.meta_lic proprietary",
				"bin/bin1.meta_lic lib/libc.a.meta_lic proprietary",
				"bin/bin2.meta_lic bin/bin2.meta_lic restricted:proprietary",
				"bin/bin2.meta_lic lib/libb.so.meta_lic restricted",
				"highest.apex.meta_lic bin/bin2.meta_lic restricted:proprietary",
				"highest.apex.meta_lic highest.apex.meta_lic restricted",
				"highest.apex.meta_lic lib/liba.so.meta_lic proprietary",
				"highest.apex.meta_lic lib/libb.so.meta_lic restricted",
				"highest.apex.meta_lic lib/libc.a.meta_lic proprietary",
				"lib/liba.so.meta_lic lib/liba.so.meta_lic proprietary",
				"lib/libb.so.meta_lic lib/libb.so.meta_lic restricted",
			},
		},
		{
			condition: "proprietary",
			name:      "apex_trimmed_labelled",
			roots:     []string{"highest.apex.meta_lic"},
			ctx:       context{stripPrefix: "proprietary/", labelConditions: true},
			expectedOut: []string{
				"bin/bin1.meta_lic:notice bin/bin1.meta_lic:notice notice",
				"bin/bin1.meta_lic:notice lib/liba.so.meta_lic:proprietary:by_exception_only proprietary:by_exception_only",
				"bin/bin1.meta_lic:notice lib/libc.a.meta_lic:proprietary:by_exception_only proprietary:by_exception_only",
				"bin/bin2.meta_lic:proprietary:by_exception_only bin/bin2.meta_lic:proprietary:by_exception_only restricted:proprietary:by_exception_only",
				"bin/bin2.meta_lic:proprietary:by_exception_only lib/libb.so.meta_lic:restricted restricted",
				"highest.apex.meta_lic:notice bin/bin1.meta_lic:notice notice",
				"highest.apex.meta_lic:notice bin/bin2.meta_lic:proprietary:by_exception_only restricted:proprietary:by_exception_only",
				"highest.apex.meta_lic:notice highest.apex.meta_lic:notice notice:restricted",
				"highest.apex.meta_lic:notice lib/liba.so.meta_lic:proprietary:by_exception_only proprietary:by_exception_only",
				"highest.apex.meta_lic:notice lib/libb.so.meta_lic:restricted restricted",
				"highest.apex.meta_lic:notice lib/libc.a.meta_lic:proprietary:by_exception_only proprietary:by_exception_only",
				"lib/liba.so.meta_lic:proprietary:by_exception_only lib/liba.so.meta_lic:proprietary:by_exception_only proprietary:by_exception_only",
				"lib/libb.so.meta_lic:restricted lib/libb.so.meta_lic:restricted restricted",
			},
		},
		{
			condition: "proprietary",
			name:      "container",
			roots:     []string{"container.zip.meta_lic"},
			expectedOut: []string{
				"proprietary/bin/bin1.meta_lic proprietary/bin/bin1.meta_lic notice",
				"proprietary/bin/bin1.meta_lic proprietary/lib/liba.so.meta_lic proprietary:by_exception_only",
				"proprietary/bin/bin1.meta_lic proprietary/lib/libc.a.meta_lic proprietary:by_exception_only",
				"proprietary/bin/bin2.meta_lic proprietary/bin/bin2.meta_lic restricted:proprietary:by_exception_only",
				"proprietary/bin/bin2.meta_lic proprietary/lib/libb.so.meta_lic restricted",
				"proprietary/container.zip.meta_lic proprietary/bin/bin1.meta_lic notice",
				"proprietary/container.zip.meta_lic proprietary/bin/bin2.meta_lic restricted:proprietary:by_exception_only",
				"proprietary/container.zip.meta_lic proprietary/container.zip.meta_lic notice:restricted",
				"proprietary/container.zip.meta_lic proprietary/lib/liba.so.meta_lic proprietary:by_exception_only",
				"proprietary/container.zip.meta_lic proprietary/lib/libb.so.meta_lic restricted",
				"proprietary/container.zip.meta_lic proprietary/lib/libc.a.meta_lic proprietary:by_exception_only",
				"proprietary/lib/liba.so.meta_lic proprietary/lib/liba.so.meta_lic proprietary:by_exception_only",
				"proprietary/lib/libb.so.meta_lic proprietary/lib/libb.so.meta_lic restricted",
			},
		},
		{
			condition: "proprietary",
			name:      "application",
			roots:     []string{"application.meta_lic"},
			expectedOut: []string{
				"proprietary/application.meta_lic proprietary/application.meta_lic notice:restricted",
				"proprietary/application.meta_lic proprietary/lib/liba.so.meta_lic restricted:proprietary:by_exception_only",
			},
		},
		{
			condition: "proprietary",
			name:      "binary",
			roots:     []string{"bin/bin1.meta_lic"},
			expectedOut: []string{
				"proprietary/bin/bin1.meta_lic proprietary/bin/bin1.meta_lic notice",
				"proprietary/bin/bin1.meta_lic proprietary/lib/liba.so.meta_lic proprietary:by_exception_only",
				"proprietary/bin/bin1.meta_lic proprietary/lib/libc.a.meta_lic proprietary:by_exception_only",
			},
		},
		{
			condition: "proprietary",
			name:      "library",
			roots:     []string{"lib/libd.so.meta_lic"},
			expectedOut: []string{
				"proprietary/lib/libd.so.meta_lic proprietary/lib/libd.so.meta_lic notice",
			},
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
			_, err := dumpResolutions(&tt.ctx, stdout, stderr, rootFiles...)
			if err != nil {
				t.Fatalf("dumpresolutions: error = %v, stderr = %v", err, stderr)
				return
			}
			if stderr.Len() > 0 {
				t.Errorf("dumpresolutions: gotStderr = %v, want none", stderr)
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
	matchString(*testContext, *compliance.LicenseGraph) string
	typeString() string
}

type targetMatcher struct {
	target     string
	conditions []string
}

// newTestCondition constructs a test license condition in the license graph.
func newTestCondition(lg *compliance.LicenseGraph, conditionName ...string) compliance.LicenseConditionSet {
	cs := compliance.NewLicenseConditionSet()
	for _, name := range conditionName {
		cs = cs.Plus(compliance.RecognizedConditionNames[name])
	}
	if cs.IsEmpty() && len(conditionName) != 0 {
		panic(fmt.Errorf("attempt to create unrecognized condition: %q", conditionName))
	}
	return cs
}

func (tm *targetMatcher) matchString(ctx *testContext, lg *compliance.LicenseGraph) string {
	cs := newTestCondition(lg, tm.conditions...)
	m := tm.target
	if !cs.IsEmpty() {
		m += "\\n" + strings.Join(cs.Names(), "\\n")
	}
	m = ctx.nodes[tm.target] + " [label=\"" + m + "\"];"
	return m
}

func (tm *targetMatcher) typeString() string {
	return "target"
}

type resolutionMatcher struct {
	appliesTo  string
	actsOn     string
	conditions []string
}

func (rm *resolutionMatcher) matchString(ctx *testContext, lg *compliance.LicenseGraph) string {
	cs := newTestCondition(lg, rm.conditions...)
	return ctx.nodes[rm.appliesTo] + " -> " + ctx.nodes[rm.actsOn] +
		" [label=\"" + strings.Join(cs.Names(), "\\n") + "\"];"
}

func (rm *resolutionMatcher) typeString() string {
	return "resolution"
}

type getMatcher func(*testContext) matcher

func matchTarget(target string, conditions ...string) getMatcher {
	return func(ctx *testContext) matcher {
		ctx.nodes[target] = fmt.Sprintf("n%d", ctx.nextNode)
		ctx.nextNode++
		return &targetMatcher{target, append([]string{}, conditions...)}
	}
}

func matchResolution(appliesTo, actsOn string, conditions ...string) getMatcher {
	return func(ctx *testContext) matcher {
		if _, ok := ctx.nodes[appliesTo]; !ok {
			ctx.nodes[appliesTo] = fmt.Sprintf("unknown%d", ctx.nextNode)
			ctx.nextNode++
		}
		if _, ok := ctx.nodes[actsOn]; !ok {
			ctx.nodes[actsOn] = fmt.Sprintf("unknown%d", ctx.nextNode)
			ctx.nextNode++
		}
		return &resolutionMatcher{appliesTo, actsOn, append([]string{}, conditions...)}
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
				matchTarget("firstparty/lib/liba.so.meta_lic"),
				matchTarget("firstparty/lib/libc.a.meta_lic"),
				matchTarget("firstparty/bin/bin2.meta_lic"),
				matchTarget("firstparty/highest.apex.meta_lic"),
				matchTarget("firstparty/lib/libb.so.meta_lic"),
				matchResolution(
					"firstparty/bin/bin1.meta_lic",
					"firstparty/bin/bin1.meta_lic",
					"notice"),
				matchResolution(
					"firstparty/bin/bin1.meta_lic",
					"firstparty/lib/liba.so.meta_lic",
					"notice"),
				matchResolution(
					"firstparty/bin/bin1.meta_lic",
					"firstparty/lib/libc.a.meta_lic",
					"notice"),
				matchResolution(
					"firstparty/bin/bin2.meta_lic",
					"firstparty/bin/bin2.meta_lic",
					"notice"),
				matchResolution(
					"firstparty/highest.apex.meta_lic",
					"firstparty/bin/bin1.meta_lic",
					"notice"),
				matchResolution(
					"firstparty/highest.apex.meta_lic",
					"firstparty/bin/bin2.meta_lic",
					"notice"),
				matchResolution(
					"firstparty/highest.apex.meta_lic",
					"firstparty/highest.apex.meta_lic",
					"notice"),
				matchResolution(
					"firstparty/highest.apex.meta_lic",
					"firstparty/lib/liba.so.meta_lic",
					"notice"),
				matchResolution(
					"firstparty/highest.apex.meta_lic",
					"firstparty/lib/libb.so.meta_lic",
					"notice"),
				matchResolution(
					"firstparty/highest.apex.meta_lic",
					"firstparty/lib/libc.a.meta_lic",
					"notice"),
				matchResolution(
					"firstparty/lib/liba.so.meta_lic",
					"firstparty/lib/liba.so.meta_lic",
					"notice"),
				matchResolution(
					"firstparty/lib/libb.so.meta_lic",
					"firstparty/lib/libb.so.meta_lic",
					"notice"),
			},
		},
		{
			condition: "firstparty",
			name:      "apex_trimmed",
			roots:     []string{"highest.apex.meta_lic"},
			ctx:       context{stripPrefix: "firstparty/"},
			expectedOut: []getMatcher{
				matchTarget("bin/bin1.meta_lic"),
				matchTarget("lib/liba.so.meta_lic"),
				matchTarget("lib/libc.a.meta_lic"),
				matchTarget("bin/bin2.meta_lic"),
				matchTarget("highest.apex.meta_lic"),
				matchTarget("lib/libb.so.meta_lic"),
				matchResolution(
					"bin/bin1.meta_lic",
					"bin/bin1.meta_lic",
					"notice"),
				matchResolution(
					"bin/bin1.meta_lic",
					"lib/liba.so.meta_lic",
					"notice"),
				matchResolution(
					"bin/bin1.meta_lic",
					"lib/libc.a.meta_lic",
					"notice"),
				matchResolution(
					"bin/bin2.meta_lic",
					"bin/bin2.meta_lic",
					"notice"),
				matchResolution(
					"highest.apex.meta_lic",
					"bin/bin1.meta_lic",
					"notice"),
				matchResolution(
					"highest.apex.meta_lic",
					"bin/bin2.meta_lic",
					"notice"),
				matchResolution(
					"highest.apex.meta_lic",
					"highest.apex.meta_lic",
					"notice"),
				matchResolution(
					"highest.apex.meta_lic",
					"lib/liba.so.meta_lic",
					"notice"),
				matchResolution(
					"highest.apex.meta_lic",
					"lib/libb.so.meta_lic",
					"notice"),
				matchResolution(
					"highest.apex.meta_lic",
					"lib/libc.a.meta_lic",
					"notice"),
				matchResolution(
					"lib/liba.so.meta_lic",
					"lib/liba.so.meta_lic",
					"notice"),
				matchResolution(
					"lib/libb.so.meta_lic",
					"lib/libb.so.meta_lic",
					"notice"),
			},
		},
		{
			condition: "firstparty",
			name:      "apex_trimmed_notice",
			roots:     []string{"highest.apex.meta_lic"},
			ctx: context{
				conditions:  []compliance.LicenseCondition{compliance.NoticeCondition},
				stripPrefix: "firstparty/",
			},
			expectedOut: []getMatcher{
				matchTarget("bin/bin1.meta_lic"),
				matchTarget("lib/liba.so.meta_lic"),
				matchTarget("lib/libc.a.meta_lic"),
				matchTarget("bin/bin2.meta_lic"),
				matchTarget("highest.apex.meta_lic"),
				matchTarget("lib/libb.so.meta_lic"),
				matchResolution(
					"bin/bin1.meta_lic",
					"bin/bin1.meta_lic",
					"notice"),
				matchResolution(
					"bin/bin1.meta_lic",
					"lib/liba.so.meta_lic",
					"notice"),
				matchResolution(
					"bin/bin1.meta_lic",
					"lib/libc.a.meta_lic",
					"notice"),
				matchResolution(
					"bin/bin2.meta_lic",
					"bin/bin2.meta_lic",
					"notice"),
				matchResolution(
					"highest.apex.meta_lic",
					"bin/bin1.meta_lic",
					"notice"),
				matchResolution(
					"highest.apex.meta_lic",
					"bin/bin2.meta_lic",
					"notice"),
				matchResolution(
					"highest.apex.meta_lic",
					"highest.apex.meta_lic",
					"notice"),
				matchResolution(
					"highest.apex.meta_lic",
					"lib/liba.so.meta_lic",
					"notice"),
				matchResolution(
					"highest.apex.meta_lic",
					"lib/libb.so.meta_lic",
					"notice"),
				matchResolution(
					"highest.apex.meta_lic",
					"lib/libc.a.meta_lic",
					"notice"),
				matchResolution(
					"lib/liba.so.meta_lic",
					"lib/liba.so.meta_lic",
					"notice"),
				matchResolution(
					"lib/libb.so.meta_lic",
					"lib/libb.so.meta_lic",
					"notice"),
			},
		},
		{
			condition: "firstparty",
			name:      "apex_trimmed_share",
			roots:     []string{"highest.apex.meta_lic"},
			ctx: context{
				conditions:  compliance.ImpliesShared.AsList(),
				stripPrefix: "firstparty/",
			},
			expectedOut: []getMatcher{},
		},
		{
			condition: "firstparty",
			name:      "apex_trimmed_private",
			roots:     []string{"highest.apex.meta_lic"},
			ctx: context{
				conditions:  compliance.ImpliesPrivate.AsList(),
				stripPrefix: "firstparty/",
			},
			expectedOut: []getMatcher{},
		},
		{
			condition: "firstparty",
			name:      "apex_trimmed_share_private",
			roots:     []string{"highest.apex.meta_lic"},
			ctx: context{
				conditions:  compliance.ImpliesShared.Union(compliance.ImpliesPrivate).AsList(),
				stripPrefix: "firstparty/",
			},
			expectedOut: []getMatcher{},
		},
		{
			condition: "firstparty",
			name:      "apex_trimmed_labelled",
			roots:     []string{"highest.apex.meta_lic"},
			ctx:       context{stripPrefix: "firstparty/", labelConditions: true},
			expectedOut: []getMatcher{
				matchTarget("bin/bin1.meta_lic", "notice"),
				matchTarget("lib/liba.so.meta_lic", "notice"),
				matchTarget("lib/libc.a.meta_lic", "notice"),
				matchTarget("bin/bin2.meta_lic", "notice"),
				matchTarget("highest.apex.meta_lic", "notice"),
				matchTarget("lib/libb.so.meta_lic", "notice"),
				matchResolution(
					"bin/bin1.meta_lic",
					"bin/bin1.meta_lic",
					"notice"),
				matchResolution(
					"bin/bin1.meta_lic",
					"lib/liba.so.meta_lic",
					"notice"),
				matchResolution(
					"bin/bin1.meta_lic",
					"lib/libc.a.meta_lic",
					"notice"),
				matchResolution(
					"bin/bin2.meta_lic",
					"bin/bin2.meta_lic",
					"notice"),
				matchResolution(
					"highest.apex.meta_lic",
					"bin/bin1.meta_lic",
					"notice"),
				matchResolution(
					"highest.apex.meta_lic",
					"bin/bin2.meta_lic",
					"notice"),
				matchResolution(
					"highest.apex.meta_lic",
					"highest.apex.meta_lic",
					"notice"),
				matchResolution(
					"highest.apex.meta_lic",
					"lib/liba.so.meta_lic",
					"notice"),
				matchResolution(
					"highest.apex.meta_lic",
					"lib/libb.so.meta_lic",
					"notice"),
				matchResolution(
					"highest.apex.meta_lic",
					"lib/libc.a.meta_lic",
					"notice"),
				matchResolution(
					"lib/liba.so.meta_lic",
					"lib/liba.so.meta_lic",
					"notice"),
				matchResolution(
					"lib/libb.so.meta_lic",
					"lib/libb.so.meta_lic",
					"notice"),
			},
		},
		{
			condition: "firstparty",
			name:      "container",
			roots:     []string{"container.zip.meta_lic"},
			expectedOut: []getMatcher{
				matchTarget("firstparty/bin/bin1.meta_lic"),
				matchTarget("firstparty/lib/liba.so.meta_lic"),
				matchTarget("firstparty/lib/libc.a.meta_lic"),
				matchTarget("firstparty/bin/bin2.meta_lic"),
				matchTarget("firstparty/container.zip.meta_lic"),
				matchTarget("firstparty/lib/libb.so.meta_lic"),
				matchResolution(
					"firstparty/bin/bin1.meta_lic",
					"firstparty/bin/bin1.meta_lic",
					"notice"),
				matchResolution(
					"firstparty/bin/bin1.meta_lic",
					"firstparty/lib/liba.so.meta_lic",
					"notice"),
				matchResolution(
					"firstparty/bin/bin1.meta_lic",
					"firstparty/lib/libc.a.meta_lic",
					"notice"),
				matchResolution(
					"firstparty/bin/bin2.meta_lic",
					"firstparty/bin/bin2.meta_lic",
					"notice"),
				matchResolution(
					"firstparty/container.zip.meta_lic",
					"firstparty/bin/bin1.meta_lic",
					"notice"),
				matchResolution(
					"firstparty/container.zip.meta_lic",
					"firstparty/bin/bin2.meta_lic",
					"notice"),
				matchResolution(
					"firstparty/container.zip.meta_lic",
					"firstparty/container.zip.meta_lic",
					"notice"),
				matchResolution(
					"firstparty/container.zip.meta_lic",
					"firstparty/lib/liba.so.meta_lic",
					"notice"),
				matchResolution(
					"firstparty/container.zip.meta_lic",
					"firstparty/lib/libb.so.meta_lic",
					"notice"),
				matchResolution(
					"firstparty/container.zip.meta_lic",
					"firstparty/lib/libc.a.meta_lic",
					"notice"),
				matchResolution(
					"firstparty/lib/liba.so.meta_lic",
					"firstparty/lib/liba.so.meta_lic",
					"notice"),
				matchResolution(
					"firstparty/lib/libb.so.meta_lic",
					"firstparty/lib/libb.so.meta_lic",
					"notice"),
			},
		},
		{
			condition: "firstparty",
			name:      "application",
			roots:     []string{"application.meta_lic"},
			expectedOut: []getMatcher{
				matchTarget("firstparty/application.meta_lic"),
				matchTarget("firstparty/lib/liba.so.meta_lic"),
				matchResolution(
					"firstparty/application.meta_lic",
					"firstparty/application.meta_lic",
					"notice"),
				matchResolution(
					"firstparty/application.meta_lic",
					"firstparty/lib/liba.so.meta_lic",
					"notice"),
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
				matchResolution(
					"firstparty/bin/bin1.meta_lic",
					"firstparty/bin/bin1.meta_lic",
					"notice"),
				matchResolution(
					"firstparty/bin/bin1.meta_lic",
					"firstparty/lib/liba.so.meta_lic",
					"notice"),
				matchResolution(
					"firstparty/bin/bin1.meta_lic",
					"firstparty/lib/libc.a.meta_lic",
					"notice"),
			},
		},
		{
			condition: "firstparty",
			name:      "library",
			roots:     []string{"lib/libd.so.meta_lic"},
			expectedOut: []getMatcher{
				matchTarget("firstparty/lib/libd.so.meta_lic"),
				matchResolution(
					"firstparty/lib/libd.so.meta_lic",
					"firstparty/lib/libd.so.meta_lic",
					"notice"),
			},
		},
		{
			condition: "notice",
			name:      "apex",
			roots:     []string{"highest.apex.meta_lic"},
			expectedOut: []getMatcher{
				matchTarget("notice/bin/bin1.meta_lic"),
				matchTarget("notice/lib/liba.so.meta_lic"),
				matchTarget("notice/lib/libc.a.meta_lic"),
				matchTarget("notice/bin/bin2.meta_lic"),
				matchTarget("notice/highest.apex.meta_lic"),
				matchTarget("notice/lib/libb.so.meta_lic"),
				matchResolution(
					"notice/bin/bin1.meta_lic",
					"notice/bin/bin1.meta_lic",
					"notice"),
				matchResolution(
					"notice/bin/bin1.meta_lic",
					"notice/lib/liba.so.meta_lic",
					"notice"),
				matchResolution(
					"notice/bin/bin1.meta_lic",
					"notice/lib/libc.a.meta_lic",
					"notice"),
				matchResolution(
					"notice/bin/bin2.meta_lic",
					"notice/bin/bin2.meta_lic",
					"notice"),
				matchResolution(
					"notice/highest.apex.meta_lic",
					"notice/bin/bin1.meta_lic",
					"notice"),
				matchResolution(
					"notice/highest.apex.meta_lic",
					"notice/bin/bin2.meta_lic",
					"notice"),
				matchResolution(
					"notice/highest.apex.meta_lic",
					"notice/highest.apex.meta_lic",
					"notice"),
				matchResolution(
					"notice/highest.apex.meta_lic",
					"notice/lib/liba.so.meta_lic",
					"notice"),
				matchResolution(
					"notice/highest.apex.meta_lic",
					"notice/lib/libb.so.meta_lic",
					"notice"),
				matchResolution(
					"notice/highest.apex.meta_lic",
					"notice/lib/libc.a.meta_lic",
					"notice"),
				matchResolution(
					"notice/lib/liba.so.meta_lic",
					"notice/lib/liba.so.meta_lic",
					"notice"),
				matchResolution(
					"notice/lib/libb.so.meta_lic",
					"notice/lib/libb.so.meta_lic",
					"notice"),
			},
		},
		{
			condition: "notice",
			name:      "apex_trimmed",
			roots:     []string{"highest.apex.meta_lic"},
			ctx:       context{stripPrefix: "notice/"},
			expectedOut: []getMatcher{
				matchTarget("bin/bin1.meta_lic"),
				matchTarget("lib/liba.so.meta_lic"),
				matchTarget("lib/libc.a.meta_lic"),
				matchTarget("bin/bin2.meta_lic"),
				matchTarget("highest.apex.meta_lic"),
				matchTarget("lib/libb.so.meta_lic"),
				matchResolution(
					"bin/bin1.meta_lic",
					"bin/bin1.meta_lic",
					"notice"),
				matchResolution(
					"bin/bin1.meta_lic",
					"lib/liba.so.meta_lic",
					"notice"),
				matchResolution(
					"bin/bin1.meta_lic",
					"lib/libc.a.meta_lic",
					"notice"),
				matchResolution(
					"bin/bin2.meta_lic",
					"bin/bin2.meta_lic",
					"notice"),
				matchResolution(
					"highest.apex.meta_lic",
					"bin/bin1.meta_lic",
					"notice"),
				matchResolution(
					"highest.apex.meta_lic",
					"bin/bin2.meta_lic",
					"notice"),
				matchResolution(
					"highest.apex.meta_lic",
					"highest.apex.meta_lic",
					"notice"),
				matchResolution(
					"highest.apex.meta_lic",
					"lib/liba.so.meta_lic",
					"notice"),
				matchResolution(
					"highest.apex.meta_lic",
					"lib/libb.so.meta_lic",
					"notice"),
				matchResolution(
					"highest.apex.meta_lic",
					"lib/libc.a.meta_lic",
					"notice"),
				matchResolution(
					"lib/liba.so.meta_lic",
					"lib/liba.so.meta_lic",
					"notice"),
				matchResolution(
					"lib/libb.so.meta_lic",
					"lib/libb.so.meta_lic",
					"notice"),
			},
		},
		{
			condition: "notice",
			name:      "apex_trimmed_notice",
			roots:     []string{"highest.apex.meta_lic"},
			ctx: context{
				conditions:  []compliance.LicenseCondition{compliance.NoticeCondition},
				stripPrefix: "notice/",
			},
			expectedOut: []getMatcher{
				matchTarget("bin/bin1.meta_lic"),
				matchTarget("lib/liba.so.meta_lic"),
				matchTarget("lib/libc.a.meta_lic"),
				matchTarget("bin/bin2.meta_lic"),
				matchTarget("highest.apex.meta_lic"),
				matchTarget("lib/libb.so.meta_lic"),
				matchResolution(
					"bin/bin1.meta_lic",
					"bin/bin1.meta_lic",
					"notice"),
				matchResolution(
					"bin/bin1.meta_lic",
					"lib/liba.so.meta_lic",
					"notice"),
				matchResolution(
					"bin/bin1.meta_lic",
					"lib/libc.a.meta_lic",
					"notice"),
				matchResolution(
					"bin/bin2.meta_lic",
					"bin/bin2.meta_lic",
					"notice"),
				matchResolution(
					"highest.apex.meta_lic",
					"bin/bin1.meta_lic",
					"notice"),
				matchResolution(
					"highest.apex.meta_lic",
					"bin/bin2.meta_lic",
					"notice"),
				matchResolution(
					"highest.apex.meta_lic",
					"highest.apex.meta_lic",
					"notice"),
				matchResolution(
					"highest.apex.meta_lic",
					"lib/liba.so.meta_lic",
					"notice"),
				matchResolution(
					"highest.apex.meta_lic",
					"lib/libb.so.meta_lic",
					"notice"),
				matchResolution(
					"highest.apex.meta_lic",
					"lib/libc.a.meta_lic",
					"notice"),
				matchResolution(
					"lib/liba.so.meta_lic",
					"lib/liba.so.meta_lic",
					"notice"),
				matchResolution(
					"lib/libb.so.meta_lic",
					"lib/libb.so.meta_lic",
					"notice"),
			},
		},
		{
			condition: "notice",
			name:      "apex_trimmed_share",
			roots:     []string{"highest.apex.meta_lic"},
			ctx: context{
				conditions:  compliance.ImpliesShared.AsList(),
				stripPrefix: "notice/",
			},
			expectedOut: []getMatcher{},
		},
		{
			condition: "notice",
			name:      "apex_trimmed_private",
			roots:     []string{"highest.apex.meta_lic"},
			ctx: context{
				conditions:  compliance.ImpliesPrivate.AsList(),
				stripPrefix: "notice/",
			},
			expectedOut: []getMatcher{},
		},
		{
			condition: "notice",
			name:      "apex_trimmed_share_private",
			roots:     []string{"highest.apex.meta_lic"},
			ctx: context{
				conditions:  compliance.ImpliesShared.Union(compliance.ImpliesPrivate).AsList(),
				stripPrefix: "notice/",
			},
			expectedOut: []getMatcher{},
		},
		{
			condition: "notice",
			name:      "apex_trimmed_labelled",
			roots:     []string{"highest.apex.meta_lic"},
			ctx:       context{stripPrefix: "notice/", labelConditions: true},
			expectedOut: []getMatcher{
				matchTarget("bin/bin1.meta_lic", "notice"),
				matchTarget("lib/liba.so.meta_lic", "notice"),
				matchTarget("lib/libc.a.meta_lic", "notice"),
				matchTarget("bin/bin2.meta_lic", "notice"),
				matchTarget("highest.apex.meta_lic", "notice"),
				matchTarget("lib/libb.so.meta_lic", "notice"),
				matchResolution(
					"bin/bin1.meta_lic",
					"bin/bin1.meta_lic",
					"notice"),
				matchResolution(
					"bin/bin1.meta_lic",
					"lib/liba.so.meta_lic",
					"notice"),
				matchResolution(
					"bin/bin1.meta_lic",
					"lib/libc.a.meta_lic",
					"notice"),
				matchResolution(
					"bin/bin2.meta_lic",
					"bin/bin2.meta_lic",
					"notice"),
				matchResolution(
					"highest.apex.meta_lic",
					"bin/bin1.meta_lic",
					"notice"),
				matchResolution(
					"highest.apex.meta_lic",
					"bin/bin2.meta_lic",
					"notice"),
				matchResolution(
					"highest.apex.meta_lic",
					"highest.apex.meta_lic",
					"notice"),
				matchResolution(
					"highest.apex.meta_lic",
					"lib/liba.so.meta_lic",
					"notice"),
				matchResolution(
					"highest.apex.meta_lic",
					"lib/libb.so.meta_lic",
					"notice"),
				matchResolution(
					"highest.apex.meta_lic",
					"lib/libc.a.meta_lic",
					"notice"),
				matchResolution(
					"lib/liba.so.meta_lic",
					"lib/liba.so.meta_lic",
					"notice"),
				matchResolution(
					"lib/libb.so.meta_lic",
					"lib/libb.so.meta_lic",
					"notice"),
			},
		},
		{
			condition: "notice",
			name:      "container",
			roots:     []string{"container.zip.meta_lic"},
			expectedOut: []getMatcher{
				matchTarget("notice/bin/bin1.meta_lic"),
				matchTarget("notice/lib/liba.so.meta_lic"),
				matchTarget("notice/lib/libc.a.meta_lic"),
				matchTarget("notice/bin/bin2.meta_lic"),
				matchTarget("notice/container.zip.meta_lic"),
				matchTarget("notice/lib/libb.so.meta_lic"),
				matchResolution(
					"notice/bin/bin1.meta_lic",
					"notice/bin/bin1.meta_lic",
					"notice"),
				matchResolution(
					"notice/bin/bin1.meta_lic",
					"notice/lib/liba.so.meta_lic",
					"notice"),
				matchResolution(
					"notice/bin/bin1.meta_lic",
					"notice/lib/libc.a.meta_lic",
					"notice"),
				matchResolution(
					"notice/bin/bin2.meta_lic",
					"notice/bin/bin2.meta_lic",
					"notice"),
				matchResolution(
					"notice/container.zip.meta_lic",
					"notice/bin/bin1.meta_lic",
					"notice"),
				matchResolution(
					"notice/container.zip.meta_lic",
					"notice/bin/bin2.meta_lic",
					"notice"),
				matchResolution(
					"notice/container.zip.meta_lic",
					"notice/container.zip.meta_lic",
					"notice"),
				matchResolution(
					"notice/container.zip.meta_lic",
					"notice/lib/liba.so.meta_lic",
					"notice"),
				matchResolution(
					"notice/container.zip.meta_lic",
					"notice/lib/libb.so.meta_lic",
					"notice"),
				matchResolution(
					"notice/container.zip.meta_lic",
					"notice/lib/libc.a.meta_lic",
					"notice"),
				matchResolution(
					"notice/lib/liba.so.meta_lic",
					"notice/lib/liba.so.meta_lic",
					"notice"),
				matchResolution(
					"notice/lib/libb.so.meta_lic",
					"notice/lib/libb.so.meta_lic",
					"notice"),
			},
		},
		{
			condition: "notice",
			name:      "application",
			roots:     []string{"application.meta_lic"},
			expectedOut: []getMatcher{
				matchTarget("notice/application.meta_lic"),
				matchTarget("notice/lib/liba.so.meta_lic"),
				matchResolution(
					"notice/application.meta_lic",
					"notice/application.meta_lic",
					"notice"),
				matchResolution(
					"notice/application.meta_lic",
					"notice/lib/liba.so.meta_lic",
					"notice"),
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
				matchResolution(
					"notice/bin/bin1.meta_lic",
					"notice/bin/bin1.meta_lic",
					"notice"),
				matchResolution(
					"notice/bin/bin1.meta_lic",
					"notice/lib/liba.so.meta_lic",
					"notice"),
				matchResolution(
					"notice/bin/bin1.meta_lic",
					"notice/lib/libc.a.meta_lic",
					"notice"),
			},
		},
		{
			condition: "notice",
			name:      "library",
			roots:     []string{"lib/libd.so.meta_lic"},
			expectedOut: []getMatcher{
				matchTarget("notice/lib/libd.so.meta_lic"),
				matchResolution(
					"notice/lib/libd.so.meta_lic",
					"notice/lib/libd.so.meta_lic",
					"notice"),
			},
		},
		{
			condition: "reciprocal",
			name:      "apex",
			roots:     []string{"highest.apex.meta_lic"},
			expectedOut: []getMatcher{
				matchTarget("reciprocal/bin/bin1.meta_lic"),
				matchTarget("reciprocal/lib/liba.so.meta_lic"),
				matchTarget("reciprocal/lib/libc.a.meta_lic"),
				matchTarget("reciprocal/bin/bin2.meta_lic"),
				matchTarget("reciprocal/highest.apex.meta_lic"),
				matchTarget("reciprocal/lib/libb.so.meta_lic"),
				matchResolution(
					"reciprocal/bin/bin1.meta_lic",
					"reciprocal/bin/bin1.meta_lic",
					"notice"),
				matchResolution(
					"reciprocal/bin/bin1.meta_lic",
					"reciprocal/lib/liba.so.meta_lic",
					"reciprocal"),
				matchResolution(
					"reciprocal/bin/bin1.meta_lic",
					"reciprocal/lib/libc.a.meta_lic",
					"reciprocal"),
				matchResolution(
					"reciprocal/bin/bin2.meta_lic",
					"reciprocal/bin/bin2.meta_lic",
					"notice"),
				matchResolution(
					"reciprocal/highest.apex.meta_lic",
					"reciprocal/bin/bin1.meta_lic",
					"notice"),
				matchResolution(
					"reciprocal/highest.apex.meta_lic",
					"reciprocal/bin/bin2.meta_lic",
					"notice"),
				matchResolution(
					"reciprocal/highest.apex.meta_lic",
					"reciprocal/highest.apex.meta_lic",
					"notice"),
				matchResolution(
					"reciprocal/highest.apex.meta_lic",
					"reciprocal/lib/liba.so.meta_lic",
					"reciprocal"),
				matchResolution(
					"reciprocal/highest.apex.meta_lic",
					"reciprocal/lib/libb.so.meta_lic",
					"notice"),
				matchResolution(
					"reciprocal/highest.apex.meta_lic",
					"reciprocal/lib/libc.a.meta_lic",
					"reciprocal"),
				matchResolution(
					"reciprocal/lib/liba.so.meta_lic",
					"reciprocal/lib/liba.so.meta_lic",
					"reciprocal"),
				matchResolution(
					"reciprocal/lib/libb.so.meta_lic",
					"reciprocal/lib/libb.so.meta_lic",
					"notice"),
			},
		},
		{
			condition: "reciprocal",
			name:      "apex_trimmed",
			roots:     []string{"highest.apex.meta_lic"},
			ctx:       context{stripPrefix: "reciprocal/"},
			expectedOut: []getMatcher{
				matchTarget("bin/bin1.meta_lic"),
				matchTarget("lib/liba.so.meta_lic"),
				matchTarget("lib/libc.a.meta_lic"),
				matchTarget("bin/bin2.meta_lic"),
				matchTarget("highest.apex.meta_lic"),
				matchTarget("lib/libb.so.meta_lic"),
				matchResolution(
					"bin/bin1.meta_lic",
					"bin/bin1.meta_lic",
					"notice"),
				matchResolution(
					"bin/bin1.meta_lic",
					"lib/liba.so.meta_lic",
					"reciprocal"),
				matchResolution(
					"bin/bin1.meta_lic",
					"lib/libc.a.meta_lic",
					"reciprocal"),
				matchResolution(
					"bin/bin2.meta_lic",
					"bin/bin2.meta_lic",
					"notice"),
				matchResolution(
					"highest.apex.meta_lic",
					"bin/bin1.meta_lic",
					"notice"),
				matchResolution(
					"highest.apex.meta_lic",
					"bin/bin2.meta_lic",
					"notice"),
				matchResolution(
					"highest.apex.meta_lic",
					"highest.apex.meta_lic",
					"notice"),
				matchResolution(
					"highest.apex.meta_lic",
					"lib/liba.so.meta_lic",
					"reciprocal"),
				matchResolution(
					"highest.apex.meta_lic",
					"lib/libb.so.meta_lic",
					"notice"),
				matchResolution(
					"highest.apex.meta_lic",
					"lib/libc.a.meta_lic",
					"reciprocal"),
				matchResolution(
					"lib/liba.so.meta_lic",
					"lib/liba.so.meta_lic",
					"reciprocal"),
				matchResolution(
					"lib/libb.so.meta_lic",
					"lib/libb.so.meta_lic",
					"notice"),
			},
		},
		{
			condition: "reciprocal",
			name:      "apex_trimmed_notice",
			roots:     []string{"highest.apex.meta_lic"},
			ctx: context{
				conditions:  []compliance.LicenseCondition{compliance.NoticeCondition},
				stripPrefix: "reciprocal/",
			},
			expectedOut: []getMatcher{
				matchTarget("bin/bin1.meta_lic"),
				matchTarget("bin/bin2.meta_lic"),
				matchTarget("highest.apex.meta_lic"),
				matchTarget("lib/libb.so.meta_lic"),
				matchResolution(
					"bin/bin1.meta_lic",
					"bin/bin1.meta_lic",
					"notice"),
				matchResolution(
					"bin/bin2.meta_lic",
					"bin/bin2.meta_lic",
					"notice"),
				matchResolution(
					"highest.apex.meta_lic",
					"bin/bin1.meta_lic",
					"notice"),
				matchResolution(
					"highest.apex.meta_lic",
					"bin/bin2.meta_lic",
					"notice"),
				matchResolution(
					"highest.apex.meta_lic",
					"highest.apex.meta_lic",
					"notice"),
				matchResolution(
					"highest.apex.meta_lic",
					"lib/libb.so.meta_lic",
					"notice"),
				matchResolution(
					"lib/libb.so.meta_lic",
					"lib/libb.so.meta_lic",
					"notice"),
			},
		},
		{
			condition: "reciprocal",
			name:      "apex_trimmed_share",
			roots:     []string{"highest.apex.meta_lic"},
			ctx: context{
				conditions:  compliance.ImpliesShared.AsList(),
				stripPrefix: "reciprocal/",
			},
			expectedOut: []getMatcher{
				matchTarget("bin/bin1.meta_lic"),
				matchTarget("lib/liba.so.meta_lic"),
				matchTarget("lib/libc.a.meta_lic"),
				matchTarget("highest.apex.meta_lic"),
				matchResolution(
					"bin/bin1.meta_lic",
					"lib/liba.so.meta_lic",
					"reciprocal"),
				matchResolution(
					"bin/bin1.meta_lic",
					"lib/libc.a.meta_lic",
					"reciprocal"),
				matchResolution(
					"highest.apex.meta_lic",
					"lib/liba.so.meta_lic",
					"reciprocal"),
				matchResolution(
					"highest.apex.meta_lic",
					"lib/libc.a.meta_lic",
					"reciprocal"),
				matchResolution(
					"lib/liba.so.meta_lic",
					"lib/liba.so.meta_lic",
					"reciprocal"),
			},
		},
		{
			condition: "reciprocal",
			name:      "apex_trimmed_private",
			roots:     []string{"highest.apex.meta_lic"},
			ctx: context{
				conditions:  compliance.ImpliesPrivate.AsList(),
				stripPrefix: "reciprocal/",
			},
			expectedOut: []getMatcher{},
		},
		{
			condition: "reciprocal",
			name:      "apex_trimmed_share_private",
			roots:     []string{"highest.apex.meta_lic"},
			ctx: context{
				conditions:  compliance.ImpliesShared.Union(compliance.ImpliesPrivate).AsList(),
				stripPrefix: "reciprocal/",
			},
			expectedOut: []getMatcher{
				matchTarget("bin/bin1.meta_lic"),
				matchTarget("lib/liba.so.meta_lic"),
				matchTarget("lib/libc.a.meta_lic"),
				matchTarget("highest.apex.meta_lic"),
				matchResolution(
					"bin/bin1.meta_lic",
					"lib/liba.so.meta_lic",
					"reciprocal"),
				matchResolution(
					"bin/bin1.meta_lic",
					"lib/libc.a.meta_lic",
					"reciprocal"),
				matchResolution(
					"highest.apex.meta_lic",
					"lib/liba.so.meta_lic",
					"reciprocal"),
				matchResolution(
					"highest.apex.meta_lic",
					"lib/libc.a.meta_lic",
					"reciprocal"),
				matchResolution(
					"lib/liba.so.meta_lic",
					"lib/liba.so.meta_lic",
					"reciprocal"),
			},
		},
		{
			condition: "reciprocal",
			name:      "apex_trimmed_labelled",
			roots:     []string{"highest.apex.meta_lic"},
			ctx:       context{stripPrefix: "reciprocal/", labelConditions: true},
			expectedOut: []getMatcher{
				matchTarget("bin/bin1.meta_lic", "notice"),
				matchTarget("lib/liba.so.meta_lic", "reciprocal"),
				matchTarget("lib/libc.a.meta_lic", "reciprocal"),
				matchTarget("bin/bin2.meta_lic", "notice"),
				matchTarget("highest.apex.meta_lic", "notice"),
				matchTarget("lib/libb.so.meta_lic", "notice"),
				matchResolution(
					"bin/bin1.meta_lic",
					"bin/bin1.meta_lic",
					"notice"),
				matchResolution(
					"bin/bin1.meta_lic",
					"lib/liba.so.meta_lic",
					"reciprocal"),
				matchResolution(
					"bin/bin1.meta_lic",
					"lib/libc.a.meta_lic",
					"reciprocal"),
				matchResolution(
					"bin/bin2.meta_lic",
					"bin/bin2.meta_lic",
					"notice"),
				matchResolution(
					"highest.apex.meta_lic",
					"bin/bin1.meta_lic",
					"notice"),
				matchResolution(
					"highest.apex.meta_lic",
					"bin/bin2.meta_lic",
					"notice"),
				matchResolution(
					"highest.apex.meta_lic",
					"highest.apex.meta_lic",
					"notice"),
				matchResolution(
					"highest.apex.meta_lic",
					"lib/liba.so.meta_lic",
					"reciprocal"),
				matchResolution(
					"highest.apex.meta_lic",
					"lib/libb.so.meta_lic",
					"notice"),
				matchResolution(
					"highest.apex.meta_lic",
					"lib/libc.a.meta_lic",
					"reciprocal"),
				matchResolution(
					"lib/liba.so.meta_lic",
					"lib/liba.so.meta_lic",
					"reciprocal"),
				matchResolution(
					"lib/libb.so.meta_lic",
					"lib/libb.so.meta_lic",
					"notice"),
			},
		},
		{
			condition: "reciprocal",
			name:      "container",
			roots:     []string{"container.zip.meta_lic"},
			expectedOut: []getMatcher{
				matchTarget("reciprocal/bin/bin1.meta_lic"),
				matchTarget("reciprocal/lib/liba.so.meta_lic"),
				matchTarget("reciprocal/lib/libc.a.meta_lic"),
				matchTarget("reciprocal/bin/bin2.meta_lic"),
				matchTarget("reciprocal/container.zip.meta_lic"),
				matchTarget("reciprocal/lib/libb.so.meta_lic"),
				matchResolution(
					"reciprocal/bin/bin1.meta_lic",
					"reciprocal/bin/bin1.meta_lic",
					"notice"),
				matchResolution(
					"reciprocal/bin/bin1.meta_lic",
					"reciprocal/lib/liba.so.meta_lic",
					"reciprocal"),
				matchResolution(
					"reciprocal/bin/bin1.meta_lic",
					"reciprocal/lib/libc.a.meta_lic",
					"reciprocal"),
				matchResolution(
					"reciprocal/bin/bin2.meta_lic",
					"reciprocal/bin/bin2.meta_lic",
					"notice"),
				matchResolution(
					"reciprocal/container.zip.meta_lic",
					"reciprocal/bin/bin1.meta_lic",
					"notice"),
				matchResolution(
					"reciprocal/container.zip.meta_lic",
					"reciprocal/bin/bin2.meta_lic",
					"notice"),
				matchResolution(
					"reciprocal/container.zip.meta_lic",
					"reciprocal/container.zip.meta_lic",
					"notice"),
				matchResolution(
					"reciprocal/container.zip.meta_lic",
					"reciprocal/lib/liba.so.meta_lic",
					"reciprocal"),
				matchResolution(
					"reciprocal/container.zip.meta_lic",
					"reciprocal/lib/libb.so.meta_lic",
					"notice"),
				matchResolution(
					"reciprocal/container.zip.meta_lic",
					"reciprocal/lib/libc.a.meta_lic",
					"reciprocal"),
				matchResolution(
					"reciprocal/lib/liba.so.meta_lic",
					"reciprocal/lib/liba.so.meta_lic",
					"reciprocal"),
				matchResolution(
					"reciprocal/lib/libb.so.meta_lic",
					"reciprocal/lib/libb.so.meta_lic",
					"notice"),
			},
		},
		{
			condition: "reciprocal",
			name:      "application",
			roots:     []string{"application.meta_lic"},
			expectedOut: []getMatcher{
				matchTarget("reciprocal/application.meta_lic"),
				matchTarget("reciprocal/lib/liba.so.meta_lic"),
				matchResolution(
					"reciprocal/application.meta_lic",
					"reciprocal/application.meta_lic",
					"notice"),
				matchResolution(
					"reciprocal/application.meta_lic",
					"reciprocal/lib/liba.so.meta_lic",
					"reciprocal"),
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
				matchResolution(
					"reciprocal/bin/bin1.meta_lic",
					"reciprocal/bin/bin1.meta_lic",
					"notice"),
				matchResolution(
					"reciprocal/bin/bin1.meta_lic",
					"reciprocal/lib/liba.so.meta_lic",
					"reciprocal"),
				matchResolution(
					"reciprocal/bin/bin1.meta_lic",
					"reciprocal/lib/libc.a.meta_lic",
					"reciprocal"),
			},
		},
		{
			condition: "reciprocal",
			name:      "library",
			roots:     []string{"lib/libd.so.meta_lic"},
			expectedOut: []getMatcher{
				matchTarget("reciprocal/lib/libd.so.meta_lic"),
				matchResolution(
					"reciprocal/lib/libd.so.meta_lic",
					"reciprocal/lib/libd.so.meta_lic",
					"notice"),
			},
		},
		{
			condition: "restricted",
			name:      "apex",
			roots:     []string{"highest.apex.meta_lic"},
			expectedOut: []getMatcher{
				matchTarget("restricted/bin/bin1.meta_lic"),
				matchTarget("restricted/lib/liba.so.meta_lic"),
				matchTarget("restricted/lib/libc.a.meta_lic"),
				matchTarget("restricted/bin/bin2.meta_lic"),
				matchTarget("restricted/lib/libb.so.meta_lic"),
				matchTarget("restricted/highest.apex.meta_lic"),
				matchResolution(
					"restricted/bin/bin1.meta_lic",
					"restricted/bin/bin1.meta_lic",
					"restricted_allows_dynamic_linking",
					"notice"),
				matchResolution(
					"restricted/bin/bin1.meta_lic",
					"restricted/lib/liba.so.meta_lic",
					"restricted_allows_dynamic_linking"),
				matchResolution(
					"restricted/bin/bin1.meta_lic",
					"restricted/lib/libc.a.meta_lic",
					"reciprocal",
					"restricted_allows_dynamic_linking"),
				matchResolution(
					"restricted/bin/bin2.meta_lic",
					"restricted/bin/bin2.meta_lic",
					"restricted",
					"notice"),
				matchResolution(
					"restricted/bin/bin2.meta_lic",
					"restricted/lib/libb.so.meta_lic",
					"restricted"),
				matchResolution(
					"restricted/highest.apex.meta_lic",
					"restricted/bin/bin1.meta_lic",
					"restricted_allows_dynamic_linking",
					"notice"),
				matchResolution(
					"restricted/highest.apex.meta_lic",
					"restricted/bin/bin2.meta_lic",
					"restricted",
					"notice"),
				matchResolution(
					"restricted/highest.apex.meta_lic",
					"restricted/highest.apex.meta_lic",
					"restricted",
					"restricted_allows_dynamic_linking",
					"notice"),
				matchResolution(
					"restricted/highest.apex.meta_lic",
					"restricted/lib/liba.so.meta_lic",
					"restricted_allows_dynamic_linking"),
				matchResolution(
					"restricted/highest.apex.meta_lic",
					"restricted/lib/libb.so.meta_lic",
					"restricted"),
				matchResolution(
					"restricted/highest.apex.meta_lic",
					"restricted/lib/libc.a.meta_lic",
					"reciprocal",
					"restricted_allows_dynamic_linking"),
				matchResolution(
					"restricted/lib/liba.so.meta_lic",
					"restricted/lib/liba.so.meta_lic",
					"restricted_allows_dynamic_linking"),
				matchResolution(
					"restricted/lib/libb.so.meta_lic",
					"restricted/lib/libb.so.meta_lic",
					"restricted"),
			},
		},
		{
			condition: "restricted",
			name:      "apex_trimmed",
			roots:     []string{"highest.apex.meta_lic"},
			ctx:       context{stripPrefix: "restricted/"},
			expectedOut: []getMatcher{
				matchTarget("bin/bin1.meta_lic"),
				matchTarget("lib/liba.so.meta_lic"),
				matchTarget("lib/libc.a.meta_lic"),
				matchTarget("bin/bin2.meta_lic"),
				matchTarget("lib/libb.so.meta_lic"),
				matchTarget("highest.apex.meta_lic"),
				matchResolution(
					"bin/bin1.meta_lic",
					"bin/bin1.meta_lic",
					"restricted_allows_dynamic_linking",
					"notice"),
				matchResolution(
					"bin/bin1.meta_lic",
					"lib/liba.so.meta_lic",
					"restricted_allows_dynamic_linking"),
				matchResolution(
					"bin/bin1.meta_lic",
					"lib/libc.a.meta_lic",
					"reciprocal",
					"restricted_allows_dynamic_linking"),
				matchResolution(
					"bin/bin2.meta_lic",
					"bin/bin2.meta_lic",
					"restricted",
					"notice"),
				matchResolution(
					"bin/bin2.meta_lic",
					"lib/libb.so.meta_lic",
					"restricted"),
				matchResolution(
					"highest.apex.meta_lic",
					"bin/bin1.meta_lic",
					"restricted_allows_dynamic_linking",
					"notice"),
				matchResolution(
					"highest.apex.meta_lic",
					"bin/bin2.meta_lic",
					"restricted",
					"notice"),
				matchResolution(
					"highest.apex.meta_lic",
					"highest.apex.meta_lic",
					"restricted",
					"restricted_allows_dynamic_linking",
					"notice"),
				matchResolution(
					"highest.apex.meta_lic",
					"lib/liba.so.meta_lic",
					"restricted_allows_dynamic_linking"),
				matchResolution(
					"highest.apex.meta_lic",
					"lib/libb.so.meta_lic",
					"restricted"),
				matchResolution(
					"highest.apex.meta_lic",
					"lib/libc.a.meta_lic",
					"reciprocal",
					"restricted_allows_dynamic_linking"),
				matchResolution(
					"lib/liba.so.meta_lic",
					"lib/liba.so.meta_lic",
					"restricted_allows_dynamic_linking"),
				matchResolution(
					"lib/libb.so.meta_lic",
					"lib/libb.so.meta_lic",
					"restricted"),
			},
		},
		{
			condition: "restricted",
			name:      "apex_trimmed_notice",
			roots:     []string{"highest.apex.meta_lic"},
			ctx: context{
				conditions:  []compliance.LicenseCondition{compliance.NoticeCondition},
				stripPrefix: "restricted/",
			},
			expectedOut: []getMatcher{
				matchTarget("bin/bin1.meta_lic"),
				matchTarget("bin/bin2.meta_lic"),
				matchTarget("highest.apex.meta_lic"),
				matchResolution(
					"bin/bin1.meta_lic",
					"bin/bin1.meta_lic",
					"notice"),
				matchResolution(
					"bin/bin2.meta_lic",
					"bin/bin2.meta_lic",
					"notice"),
				matchResolution(
					"highest.apex.meta_lic",
					"bin/bin1.meta_lic",
					"notice"),
				matchResolution(
					"highest.apex.meta_lic",
					"bin/bin2.meta_lic",
					"notice"),
				matchResolution(
					"highest.apex.meta_lic",
					"highest.apex.meta_lic",
					"notice"),
			},
		},
		{
			condition: "restricted",
			name:      "apex_trimmed_share",
			roots:     []string{"highest.apex.meta_lic"},
			ctx: context{
				conditions:  compliance.ImpliesShared.AsList(),
				stripPrefix: "restricted/",
			},
			expectedOut: []getMatcher{
				matchTarget("bin/bin1.meta_lic"),
				matchTarget("lib/liba.so.meta_lic"),
				matchTarget("lib/libc.a.meta_lic"),
				matchTarget("bin/bin2.meta_lic"),
				matchTarget("lib/libb.so.meta_lic"),
				matchTarget("highest.apex.meta_lic"),
				matchResolution(
					"bin/bin1.meta_lic",
					"bin/bin1.meta_lic",
					"restricted_allows_dynamic_linking"),
				matchResolution(
					"bin/bin1.meta_lic",
					"lib/liba.so.meta_lic",
					"restricted_allows_dynamic_linking"),
				matchResolution(
					"bin/bin1.meta_lic",
					"lib/libc.a.meta_lic",
					"reciprocal",
					"restricted_allows_dynamic_linking"),
				matchResolution(
					"bin/bin2.meta_lic",
					"bin/bin2.meta_lic",
					"restricted"),
				matchResolution(
					"bin/bin2.meta_lic",
					"lib/libb.so.meta_lic",
					"restricted"),
				matchResolution(
					"highest.apex.meta_lic",
					"bin/bin1.meta_lic",
					"restricted_allows_dynamic_linking"),
				matchResolution(
					"highest.apex.meta_lic",
					"bin/bin2.meta_lic",
					"restricted"),
				matchResolution(
					"highest.apex.meta_lic",
					"highest.apex.meta_lic",
					"restricted",
					"restricted_allows_dynamic_linking"),
				matchResolution(
					"highest.apex.meta_lic",
					"lib/liba.so.meta_lic",
					"restricted_allows_dynamic_linking"),
				matchResolution(
					"highest.apex.meta_lic",
					"lib/libb.so.meta_lic",
					"restricted"),
				matchResolution(
					"highest.apex.meta_lic",
					"lib/libc.a.meta_lic",
					"reciprocal",
					"restricted_allows_dynamic_linking"),
				matchResolution(
					"lib/liba.so.meta_lic",
					"lib/liba.so.meta_lic",
					"restricted_allows_dynamic_linking"),
				matchResolution(
					"lib/libb.so.meta_lic",
					"lib/libb.so.meta_lic",
					"restricted"),
			},
		},
		{
			condition: "restricted",
			name:      "apex_trimmed_private",
			roots:     []string{"highest.apex.meta_lic"},
			ctx: context{
				conditions:  compliance.ImpliesPrivate.AsList(),
				stripPrefix: "restricted/",
			},
			expectedOut: []getMatcher{},
		},
		{
			condition: "restricted",
			name:      "apex_trimmed_share_private",
			roots:     []string{"highest.apex.meta_lic"},
			ctx: context{
				conditions:  compliance.ImpliesShared.Union(compliance.ImpliesPrivate).AsList(),
				stripPrefix: "restricted/",
			},
			expectedOut: []getMatcher{
				matchTarget("bin/bin1.meta_lic"),
				matchTarget("lib/liba.so.meta_lic"),
				matchTarget("lib/libc.a.meta_lic"),
				matchTarget("bin/bin2.meta_lic"),
				matchTarget("lib/libb.so.meta_lic"),
				matchTarget("highest.apex.meta_lic"),
				matchResolution(
					"bin/bin1.meta_lic",
					"bin/bin1.meta_lic",
					"restricted_allows_dynamic_linking"),
				matchResolution(
					"bin/bin1.meta_lic",
					"lib/liba.so.meta_lic",
					"restricted_allows_dynamic_linking"),
				matchResolution(
					"bin/bin1.meta_lic",
					"lib/libc.a.meta_lic",
					"reciprocal",
					"restricted_allows_dynamic_linking"),
				matchResolution(
					"bin/bin2.meta_lic",
					"bin/bin2.meta_lic",
					"restricted"),
				matchResolution(
					"bin/bin2.meta_lic",
					"lib/libb.so.meta_lic",
					"restricted"),
				matchResolution(
					"highest.apex.meta_lic",
					"bin/bin1.meta_lic",
					"restricted_allows_dynamic_linking"),
				matchResolution(
					"highest.apex.meta_lic",
					"bin/bin2.meta_lic",
					"restricted"),
				matchResolution(
					"highest.apex.meta_lic",
					"highest.apex.meta_lic",
					"restricted",
					"restricted_allows_dynamic_linking"),
				matchResolution(
					"highest.apex.meta_lic",
					"lib/liba.so.meta_lic",
					"restricted_allows_dynamic_linking"),
				matchResolution(
					"highest.apex.meta_lic",
					"lib/libb.so.meta_lic",
					"restricted"),
				matchResolution(
					"highest.apex.meta_lic",
					"lib/libc.a.meta_lic",
					"reciprocal",
					"restricted_allows_dynamic_linking"),
				matchResolution(
					"lib/liba.so.meta_lic",
					"lib/liba.so.meta_lic",
					"restricted_allows_dynamic_linking"),
				matchResolution(
					"lib/libb.so.meta_lic",
					"lib/libb.so.meta_lic",
					"restricted"),
			},
		},
		{
			condition: "restricted",
			name:      "apex_trimmed_labelled",
			roots:     []string{"highest.apex.meta_lic"},
			ctx:       context{stripPrefix: "restricted/", labelConditions: true},
			expectedOut: []getMatcher{
				matchTarget("bin/bin1.meta_lic", "notice"),
				matchTarget("lib/liba.so.meta_lic", "restricted_allows_dynamic_linking"),
				matchTarget("lib/libc.a.meta_lic", "reciprocal"),
				matchTarget("bin/bin2.meta_lic", "notice"),
				matchTarget("lib/libb.so.meta_lic", "restricted"),
				matchTarget("highest.apex.meta_lic", "notice"),
				matchResolution(
					"bin/bin1.meta_lic",
					"bin/bin1.meta_lic",
					"restricted_allows_dynamic_linking",
					"notice"),
				matchResolution(
					"bin/bin1.meta_lic",
					"lib/liba.so.meta_lic",
					"restricted_allows_dynamic_linking"),
				matchResolution(
					"bin/bin1.meta_lic",
					"lib/libc.a.meta_lic",
					"reciprocal",
					"restricted_allows_dynamic_linking"),
				matchResolution(
					"bin/bin2.meta_lic",
					"bin/bin2.meta_lic",
					"restricted",
					"notice"),
				matchResolution(
					"bin/bin2.meta_lic",
					"lib/libb.so.meta_lic",
					"restricted"),
				matchResolution(
					"highest.apex.meta_lic",
					"bin/bin1.meta_lic",
					"restricted_allows_dynamic_linking",
					"notice"),
				matchResolution(
					"highest.apex.meta_lic",
					"bin/bin2.meta_lic",
					"restricted",
					"notice"),
				matchResolution(
					"highest.apex.meta_lic",
					"highest.apex.meta_lic",
					"restricted",
					"restricted_allows_dynamic_linking",
					"notice"),
				matchResolution(
					"highest.apex.meta_lic",
					"lib/liba.so.meta_lic",
					"restricted_allows_dynamic_linking"),
				matchResolution(
					"highest.apex.meta_lic",
					"lib/libb.so.meta_lic",
					"restricted"),
				matchResolution(
					"highest.apex.meta_lic",
					"lib/libc.a.meta_lic",
					"reciprocal",
					"restricted_allows_dynamic_linking"),
				matchResolution(
					"lib/liba.so.meta_lic",
					"lib/liba.so.meta_lic",
					"restricted_allows_dynamic_linking"),
				matchResolution(
					"lib/libb.so.meta_lic",
					"lib/libb.so.meta_lic",
					"restricted"),
			},
		},
		{
			condition: "restricted",
			name:      "container",
			roots:     []string{"container.zip.meta_lic"},
			expectedOut: []getMatcher{
				matchTarget("restricted/bin/bin1.meta_lic"),
				matchTarget("restricted/lib/liba.so.meta_lic"),
				matchTarget("restricted/lib/libc.a.meta_lic"),
				matchTarget("restricted/bin/bin2.meta_lic"),
				matchTarget("restricted/lib/libb.so.meta_lic"),
				matchTarget("restricted/container.zip.meta_lic"),
				matchResolution(
					"restricted/bin/bin1.meta_lic",
					"restricted/bin/bin1.meta_lic",
					"restricted_allows_dynamic_linking",
					"notice"),
				matchResolution(
					"restricted/bin/bin1.meta_lic",
					"restricted/lib/liba.so.meta_lic",
					"restricted_allows_dynamic_linking"),
				matchResolution(
					"restricted/bin/bin1.meta_lic",
					"restricted/lib/libc.a.meta_lic",
					"reciprocal",
					"restricted_allows_dynamic_linking"),
				matchResolution(
					"restricted/bin/bin2.meta_lic",
					"restricted/bin/bin2.meta_lic",
					"restricted",
					"notice"),
				matchResolution(
					"restricted/bin/bin2.meta_lic",
					"restricted/lib/libb.so.meta_lic",
					"restricted"),
				matchResolution(
					"restricted/container.zip.meta_lic",
					"restricted/bin/bin1.meta_lic",
					"restricted_allows_dynamic_linking",
					"notice"),
				matchResolution(
					"restricted/container.zip.meta_lic",
					"restricted/bin/bin2.meta_lic",
					"restricted",
					"notice"),
				matchResolution(
					"restricted/container.zip.meta_lic",
					"restricted/container.zip.meta_lic",
					"restricted",
					"restricted_allows_dynamic_linking",
					"notice"),
				matchResolution(
					"restricted/container.zip.meta_lic",
					"restricted/lib/liba.so.meta_lic",
					"restricted_allows_dynamic_linking"),
				matchResolution(
					"restricted/container.zip.meta_lic",
					"restricted/lib/libb.so.meta_lic",
					"restricted"),
				matchResolution(
					"restricted/container.zip.meta_lic",
					"restricted/lib/libc.a.meta_lic",
					"reciprocal",
					"restricted_allows_dynamic_linking"),
				matchResolution(
					"restricted/lib/liba.so.meta_lic",
					"restricted/lib/liba.so.meta_lic",
					"restricted_allows_dynamic_linking"),
				matchResolution(
					"restricted/lib/libb.so.meta_lic",
					"restricted/lib/libb.so.meta_lic",
					"restricted"),
			},
		},
		{
			condition: "restricted",
			name:      "application",
			roots:     []string{"application.meta_lic"},
			expectedOut: []getMatcher{
				matchTarget("restricted/application.meta_lic"),
				matchTarget("restricted/lib/liba.so.meta_lic"),
				matchResolution(
					"restricted/application.meta_lic",
					"restricted/application.meta_lic",
					"restricted",
					"restricted_allows_dynamic_linking",
					"notice"),
				matchResolution(
					"restricted/application.meta_lic",
					"restricted/lib/liba.so.meta_lic",
					"restricted_allows_dynamic_linking",
					"restricted"),
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
				matchResolution(
					"restricted/bin/bin1.meta_lic",
					"restricted/bin/bin1.meta_lic",
					"restricted_allows_dynamic_linking",
					"notice"),
				matchResolution(
					"restricted/bin/bin1.meta_lic",
					"restricted/lib/liba.so.meta_lic",
					"restricted_allows_dynamic_linking"),
				matchResolution(
					"restricted/bin/bin1.meta_lic",
					"restricted/lib/libc.a.meta_lic",
					"restricted_allows_dynamic_linking",
					"reciprocal"),
			},
		},
		{
			condition: "restricted",
			name:      "library",
			roots:     []string{"lib/libd.so.meta_lic"},
			expectedOut: []getMatcher{
				matchTarget("restricted/lib/libd.so.meta_lic"),
				matchResolution(
					"restricted/lib/libd.so.meta_lic",
					"restricted/lib/libd.so.meta_lic",
					"notice"),
			},
		},
		{
			condition: "proprietary",
			name:      "apex",
			roots:     []string{"highest.apex.meta_lic"},
			expectedOut: []getMatcher{
				matchTarget("proprietary/bin/bin1.meta_lic"),
				matchTarget("proprietary/lib/liba.so.meta_lic"),
				matchTarget("proprietary/lib/libc.a.meta_lic"),
				matchTarget("proprietary/bin/bin2.meta_lic"),
				matchTarget("proprietary/lib/libb.so.meta_lic"),
				matchTarget("proprietary/highest.apex.meta_lic"),
				matchResolution(
					"proprietary/bin/bin1.meta_lic",
					"proprietary/bin/bin1.meta_lic",
					"notice"),
				matchResolution(
					"proprietary/bin/bin1.meta_lic",
					"proprietary/lib/liba.so.meta_lic",
					"by_exception_only",
					"proprietary"),
				matchResolution(
					"proprietary/bin/bin1.meta_lic",
					"proprietary/lib/libc.a.meta_lic",
					"by_exception_only",
					"proprietary"),
				matchResolution(
					"proprietary/bin/bin2.meta_lic",
					"proprietary/bin/bin2.meta_lic",
					"restricted",
					"by_exception_only",
					"proprietary"),
				matchResolution(
					"proprietary/bin/bin2.meta_lic",
					"proprietary/lib/libb.so.meta_lic",
					"restricted"),
				matchResolution(
					"proprietary/highest.apex.meta_lic",
					"proprietary/bin/bin1.meta_lic",
					"notice"),
				matchResolution(
					"proprietary/highest.apex.meta_lic",
					"proprietary/bin/bin2.meta_lic",
					"restricted",
					"by_exception_only",
					"proprietary"),
				matchResolution(
					"proprietary/highest.apex.meta_lic",
					"proprietary/highest.apex.meta_lic",
					"restricted",
					"notice"),
				matchResolution(
					"proprietary/highest.apex.meta_lic",
					"proprietary/lib/liba.so.meta_lic",
					"by_exception_only",
					"proprietary"),
				matchResolution(
					"proprietary/highest.apex.meta_lic",
					"proprietary/lib/libb.so.meta_lic",
					"restricted"),
				matchResolution(
					"proprietary/highest.apex.meta_lic",
					"proprietary/lib/libc.a.meta_lic",
					"by_exception_only",
					"proprietary"),
				matchResolution(
					"proprietary/lib/liba.so.meta_lic",
					"proprietary/lib/liba.so.meta_lic",
					"by_exception_only",
					"proprietary"),
				matchResolution(
					"proprietary/lib/libb.so.meta_lic",
					"proprietary/lib/libb.so.meta_lic",
					"restricted"),
			},
		},
		{
			condition: "proprietary",
			name:      "apex_trimmed",
			roots:     []string{"highest.apex.meta_lic"},
			ctx:       context{stripPrefix: "proprietary/"},
			expectedOut: []getMatcher{
				matchTarget("bin/bin1.meta_lic"),
				matchTarget("lib/liba.so.meta_lic"),
				matchTarget("lib/libc.a.meta_lic"),
				matchTarget("bin/bin2.meta_lic"),
				matchTarget("lib/libb.so.meta_lic"),
				matchTarget("highest.apex.meta_lic"),
				matchResolution(
					"bin/bin1.meta_lic",
					"bin/bin1.meta_lic",
					"notice"),
				matchResolution(
					"bin/bin1.meta_lic",
					"lib/liba.so.meta_lic",
					"by_exception_only",
					"proprietary"),
				matchResolution(
					"bin/bin1.meta_lic",
					"lib/libc.a.meta_lic",
					"by_exception_only",
					"proprietary"),
				matchResolution(
					"bin/bin2.meta_lic",
					"bin/bin2.meta_lic",
					"by_exception_only",
					"restricted",
					"proprietary"),
				matchResolution(
					"bin/bin2.meta_lic",
					"lib/libb.so.meta_lic",
					"restricted"),
				matchResolution(
					"highest.apex.meta_lic",
					"bin/bin1.meta_lic",
					"notice"),
				matchResolution(
					"highest.apex.meta_lic",
					"bin/bin2.meta_lic",
					"restricted",
					"by_exception_only",
					"proprietary"),
				matchResolution(
					"highest.apex.meta_lic",
					"highest.apex.meta_lic",
					"restricted",
					"notice"),
				matchResolution(
					"highest.apex.meta_lic",
					"lib/liba.so.meta_lic",
					"by_exception_only",
					"proprietary"),
				matchResolution(
					"highest.apex.meta_lic",
					"lib/libb.so.meta_lic",
					"restricted"),
				matchResolution(
					"highest.apex.meta_lic",
					"lib/libc.a.meta_lic",
					"by_exception_only",
					"proprietary"),
				matchResolution(
					"lib/liba.so.meta_lic",
					"lib/liba.so.meta_lic",
					"by_exception_only",
					"proprietary"),
				matchResolution(
					"lib/libb.so.meta_lic",
					"lib/libb.so.meta_lic",
					"restricted"),
			},
		},
		{
			condition: "proprietary",
			name:      "apex_trimmed_notice",
			roots:     []string{"highest.apex.meta_lic"},
			ctx: context{
				conditions:  []compliance.LicenseCondition{compliance.NoticeCondition},
				stripPrefix: "proprietary/",
			},
			expectedOut: []getMatcher{
				matchTarget("bin/bin1.meta_lic"),
				matchTarget("highest.apex.meta_lic"),
				matchResolution(
					"bin/bin1.meta_lic",
					"bin/bin1.meta_lic",
					"notice"),
				matchResolution(
					"highest.apex.meta_lic",
					"bin/bin1.meta_lic",
					"notice"),
				matchResolution(
					"highest.apex.meta_lic",
					"highest.apex.meta_lic",
					"notice"),
			},
		},
		{
			condition: "proprietary",
			name:      "apex_trimmed_share",
			roots:     []string{"highest.apex.meta_lic"},
			ctx: context{
				conditions:  compliance.ImpliesShared.AsList(),
				stripPrefix: "proprietary/",
			},
			expectedOut: []getMatcher{
				matchTarget("bin/bin2.meta_lic"),
				matchTarget("lib/libb.so.meta_lic"),
				matchTarget("highest.apex.meta_lic"),
				matchResolution(
					"bin/bin2.meta_lic",
					"bin/bin2.meta_lic",
					"restricted"),
				matchResolution(
					"bin/bin2.meta_lic",
					"lib/libb.so.meta_lic",
					"restricted"),
				matchResolution(
					"highest.apex.meta_lic",
					"bin/bin2.meta_lic",
					"restricted"),
				matchResolution(
					"highest.apex.meta_lic",
					"highest.apex.meta_lic",
					"restricted"),
				matchResolution(
					"highest.apex.meta_lic",
					"lib/libb.so.meta_lic",
					"restricted"),
				matchResolution(
					"lib/libb.so.meta_lic",
					"lib/libb.so.meta_lic",
					"restricted"),
			},
		},
		{
			condition: "proprietary",
			name:      "apex_trimmed_private",
			roots:     []string{"highest.apex.meta_lic"},
			ctx: context{
				conditions:  compliance.ImpliesPrivate.AsList(),
				stripPrefix: "proprietary/",
			},
			expectedOut: []getMatcher{
				matchTarget("bin/bin1.meta_lic"),
				matchTarget("lib/liba.so.meta_lic"),
				matchTarget("lib/libc.a.meta_lic"),
				matchTarget("bin/bin2.meta_lic"),
				matchTarget("highest.apex.meta_lic"),
				matchResolution(
					"bin/bin1.meta_lic",
					"lib/liba.so.meta_lic",
					"proprietary"),
				matchResolution(
					"bin/bin1.meta_lic",
					"lib/libc.a.meta_lic",
					"proprietary"),
				matchResolution(
					"bin/bin2.meta_lic",
					"bin/bin2.meta_lic",
					"proprietary"),
				matchResolution(
					"highest.apex.meta_lic",
					"bin/bin2.meta_lic",
					"proprietary"),
				matchResolution(
					"highest.apex.meta_lic",
					"lib/liba.so.meta_lic",
					"proprietary"),
				matchResolution(
					"highest.apex.meta_lic",
					"lib/libc.a.meta_lic",
					"proprietary"),
				matchResolution(
					"lib/liba.so.meta_lic",
					"lib/liba.so.meta_lic",
					"proprietary"),
			},
		},
		{
			condition: "proprietary",
			name:      "apex_trimmed_share_private",
			roots:     []string{"highest.apex.meta_lic"},
			ctx: context{
				conditions:  compliance.ImpliesShared.Union(compliance.ImpliesPrivate).AsList(),
				stripPrefix: "proprietary/",
			},
			expectedOut: []getMatcher{
				matchTarget("bin/bin1.meta_lic"),
				matchTarget("lib/liba.so.meta_lic"),
				matchTarget("lib/libc.a.meta_lic"),
				matchTarget("bin/bin2.meta_lic"),
				matchTarget("lib/libb.so.meta_lic"),
				matchTarget("highest.apex.meta_lic"),
				matchResolution(
					"bin/bin1.meta_lic",
					"lib/liba.so.meta_lic",
					"proprietary"),
				matchResolution(
					"bin/bin1.meta_lic",
					"lib/libc.a.meta_lic",
					"proprietary"),
				matchResolution(
					"bin/bin2.meta_lic",
					"bin/bin2.meta_lic",
					"restricted",
					"proprietary"),
				matchResolution(
					"bin/bin2.meta_lic",
					"lib/libb.so.meta_lic",
					"restricted"),
				matchResolution(
					"highest.apex.meta_lic",
					"bin/bin2.meta_lic",
					"restricted",
					"proprietary"),
				matchResolution(
					"highest.apex.meta_lic",
					"highest.apex.meta_lic",
					"restricted"),
				matchResolution(
					"highest.apex.meta_lic",
					"lib/liba.so.meta_lic",
					"proprietary"),
				matchResolution(
					"highest.apex.meta_lic",
					"lib/libb.so.meta_lic",
					"restricted"),
				matchResolution(
					"highest.apex.meta_lic",
					"lib/libc.a.meta_lic",
					"proprietary"),
				matchResolution(
					"lib/liba.so.meta_lic",
					"lib/liba.so.meta_lic",
					"proprietary"),
				matchResolution(
					"lib/libb.so.meta_lic",
					"lib/libb.so.meta_lic",
					"restricted"),
			},
		},
		{
			condition: "proprietary",
			name:      "apex_trimmed_labelled",
			roots:     []string{"highest.apex.meta_lic"},
			ctx:       context{stripPrefix: "proprietary/", labelConditions: true},
			expectedOut: []getMatcher{
				matchTarget("bin/bin1.meta_lic", "notice"),
				matchTarget("lib/liba.so.meta_lic", "by_exception_only", "proprietary"),
				matchTarget("lib/libc.a.meta_lic", "by_exception_only", "proprietary"),
				matchTarget("bin/bin2.meta_lic", "by_exception_only", "proprietary"),
				matchTarget("lib/libb.so.meta_lic", "restricted"),
				matchTarget("highest.apex.meta_lic", "notice"),
				matchResolution(
					"bin/bin1.meta_lic",
					"bin/bin1.meta_lic",
					"notice"),
				matchResolution(
					"bin/bin1.meta_lic",
					"lib/liba.so.meta_lic",
					"by_exception_only",
					"proprietary"),
				matchResolution(
					"bin/bin1.meta_lic",
					"lib/libc.a.meta_lic",
					"by_exception_only",
					"proprietary"),
				matchResolution(
					"bin/bin2.meta_lic",
					"bin/bin2.meta_lic",
					"restricted",
					"by_exception_only",
					"proprietary"),
				matchResolution(
					"bin/bin2.meta_lic",
					"lib/libb.so.meta_lic",
					"restricted"),
				matchResolution(
					"highest.apex.meta_lic",
					"bin/bin1.meta_lic",
					"notice"),
				matchResolution(
					"highest.apex.meta_lic",
					"bin/bin2.meta_lic",
					"restricted",
					"by_exception_only",
					"proprietary"),
				matchResolution(
					"highest.apex.meta_lic",
					"highest.apex.meta_lic",
					"restricted",
					"notice"),
				matchResolution(
					"highest.apex.meta_lic",
					"lib/liba.so.meta_lic",
					"by_exception_only",
					"proprietary"),
				matchResolution(
					"highest.apex.meta_lic",
					"lib/libb.so.meta_lic",
					"restricted"),
				matchResolution(
					"highest.apex.meta_lic",
					"lib/libc.a.meta_lic",
					"by_exception_only",
					"proprietary"),
				matchResolution(
					"lib/liba.so.meta_lic",
					"lib/liba.so.meta_lic",
					"by_exception_only",
					"proprietary"),
				matchResolution(
					"lib/libb.so.meta_lic",
					"lib/libb.so.meta_lic",
					"restricted"),
			},
		},
		{
			condition: "proprietary",
			name:      "container",
			roots:     []string{"container.zip.meta_lic"},
			expectedOut: []getMatcher{
				matchTarget("proprietary/bin/bin1.meta_lic"),
				matchTarget("proprietary/lib/liba.so.meta_lic"),
				matchTarget("proprietary/lib/libc.a.meta_lic"),
				matchTarget("proprietary/bin/bin2.meta_lic"),
				matchTarget("proprietary/lib/libb.so.meta_lic"),
				matchTarget("proprietary/container.zip.meta_lic"),
				matchResolution(
					"proprietary/bin/bin1.meta_lic",
					"proprietary/bin/bin1.meta_lic",
					"notice"),
				matchResolution(
					"proprietary/bin/bin1.meta_lic",
					"proprietary/lib/liba.so.meta_lic",
					"by_exception_only",
					"proprietary"),
				matchResolution(
					"proprietary/bin/bin1.meta_lic",
					"proprietary/lib/libc.a.meta_lic",
					"by_exception_only",
					"proprietary"),
				matchResolution(
					"proprietary/bin/bin2.meta_lic",
					"proprietary/bin/bin2.meta_lic",
					"restricted",
					"by_exception_only",
					"proprietary"),
				matchResolution(
					"proprietary/bin/bin2.meta_lic",
					"proprietary/lib/libb.so.meta_lic",
					"restricted"),
				matchResolution(
					"proprietary/container.zip.meta_lic",
					"proprietary/bin/bin1.meta_lic",
					"notice"),
				matchResolution(
					"proprietary/container.zip.meta_lic",
					"proprietary/bin/bin2.meta_lic",
					"restricted",
					"by_exception_only",
					"proprietary"),
				matchResolution(
					"proprietary/container.zip.meta_lic",
					"proprietary/container.zip.meta_lic",
					"restricted",
					"notice"),
				matchResolution(
					"proprietary/container.zip.meta_lic",
					"proprietary/lib/liba.so.meta_lic",
					"by_exception_only",
					"proprietary"),
				matchResolution(
					"proprietary/container.zip.meta_lic",
					"proprietary/lib/libb.so.meta_lic",
					"restricted"),
				matchResolution(
					"proprietary/container.zip.meta_lic",
					"proprietary/lib/libc.a.meta_lic",
					"by_exception_only",
					"proprietary"),
				matchResolution(
					"proprietary/lib/liba.so.meta_lic",
					"proprietary/lib/liba.so.meta_lic",
					"by_exception_only",
					"proprietary"),
				matchResolution(
					"proprietary/lib/libb.so.meta_lic",
					"proprietary/lib/libb.so.meta_lic",
					"restricted"),
			},
		},
		{
			condition: "proprietary",
			name:      "application",
			roots:     []string{"application.meta_lic"},
			expectedOut: []getMatcher{
				matchTarget("proprietary/application.meta_lic"),
				matchTarget("proprietary/lib/liba.so.meta_lic"),
				matchResolution(
					"proprietary/application.meta_lic",
					"proprietary/application.meta_lic",
					"notice",
					"restricted"),
				matchResolution(
					"proprietary/application.meta_lic",
					"proprietary/lib/liba.so.meta_lic",
					"restricted",
					"by_exception_only",
					"proprietary"),
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
				matchResolution(
					"proprietary/bin/bin1.meta_lic",
					"proprietary/bin/bin1.meta_lic",
					"notice"),
				matchResolution(
					"proprietary/bin/bin1.meta_lic",
					"proprietary/lib/liba.so.meta_lic",
					"by_exception_only",
					"proprietary"),
				matchResolution(
					"proprietary/bin/bin1.meta_lic",
					"proprietary/lib/libc.a.meta_lic",
					"by_exception_only",
					"proprietary"),
			},
		},
		{
			condition: "proprietary",
			name:      "library",
			roots:     []string{"lib/libd.so.meta_lic"},
			expectedOut: []getMatcher{
				matchTarget("proprietary/lib/libd.so.meta_lic"),
				matchResolution(
					"proprietary/lib/libd.so.meta_lic",
					"proprietary/lib/libd.so.meta_lic",
					"notice"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.condition+" "+tt.name, func(t *testing.T) {
			ctx := &testContext{0, make(map[string]string)}

			stdout := &bytes.Buffer{}
			stderr := &bytes.Buffer{}

			rootFiles := make([]string, 0, len(tt.roots))
			for _, r := range tt.roots {
				rootFiles = append(rootFiles, tt.condition+"/"+r)
			}
			tt.ctx.graphViz = true
			lg, err := dumpResolutions(&tt.ctx, stdout, stderr, rootFiles...)
			if err != nil {
				t.Fatalf("dumpresolutions: error = %v, stderr = %v", err, stderr)
				return
			}
			if stderr.Len() > 0 {
				t.Errorf("dumpresolutions: gotStderr = %v, want none", stderr)
			}

			expectedOut := &bytes.Buffer{}
			for _, eo := range tt.expectedOut {
				m := eo(ctx)
				expectedOut.WriteString(m.matchString(ctx, lg))
				expectedOut.WriteString("\n")
			}

			outList := strings.Split(stdout.String(), "\n")
			outLine := 0
			if outList[outLine] != "strict digraph {" {
				t.Errorf("dumpresolutions: got 1st line %v, want strict digraph {", outList[outLine])
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
				t.Errorf("dumpresolutions: got last line %v, want }", outList[endOut-1])
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
					t.Errorf("dumpresolutions: missing lines at end of graph, want %d lines %v", len(expectedList)-matchLine, strings.Join(expectedList[matchLine:], "\n"))
				} else if matchLine >= len(expectedList) {
					t.Errorf("dumpresolutions: unexpected lines at end of graph starting line %d, got %v, want nothing", outLine+1, strings.Join(outList[outLine:], "\n"))
				} else {
					t.Errorf("dumpresolutions: at line %d, got %v, want %v", outLine+1, strings.Join(outList[outLine:], "\n"), strings.Join(expectedList[matchLine:], "\n"))
				}
			}
		})
	}
}
