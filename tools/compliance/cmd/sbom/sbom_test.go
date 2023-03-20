// Copyright 2022 Google LLC
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
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"

	"android/soong/tools/compliance"
	"github.com/spdx/tools-golang/builder/builder2v2"
	"github.com/spdx/tools-golang/spdx/common"
	spdx "github.com/spdx/tools-golang/spdx/v2_2"
)

func TestMain(m *testing.M) {
	// Change into the parent directory before running the tests
	// so they can find the testdata directory.
	if err := os.Chdir(".."); err != nil {
		fmt.Printf("failed to change to testdata directory: %s\n", err)
		os.Exit(1)
	}
	os.Exit(m.Run())
}

func Test(t *testing.T) {
	tests := []struct {
		condition    string
		name         string
		outDir       string
		roots        []string
		stripPrefix  string
		expectedOut  *spdx.Document
		expectedDeps []string
	}{
		{
			condition: "firstparty",
			name:      "apex",
			roots:     []string{"highest.apex.meta_lic"},
			expectedOut: &spdx.Document{
				SPDXIdentifier: "DOCUMENT",
				CreationInfo:   getCreationInfo(t),
				Packages: []*spdx.Package{
					{
						PackageName:             "testdata-firstparty-highest.apex.meta_lic",
						PackageVersion:          "NOASSERTION",
						PackageDownloadLocation: "NOASSERTION",
						PackageSPDXIdentifier:   common.ElementID("testdata-firstparty-highest.apex.meta_lic"),
						PackageLicenseConcluded: "LicenseRef-testdata-firstparty-FIRST_PARTY_LICENSE",
					},
					{
						PackageName:             "testdata-firstparty-bin-bin1.meta_lic",
						PackageVersion:          "NOASSERTION",
						PackageDownloadLocation: "NOASSERTION",
						PackageSPDXIdentifier:   common.ElementID("testdata-firstparty-bin-bin1.meta_lic"),
						PackageLicenseConcluded: "LicenseRef-testdata-firstparty-FIRST_PARTY_LICENSE",
					},
					{
						PackageName:             "testdata-firstparty-bin-bin2.meta_lic",
						PackageVersion:          "NOASSERTION",
						PackageDownloadLocation: "NOASSERTION",
						PackageSPDXIdentifier:   common.ElementID("testdata-firstparty-bin-bin2.meta_lic"),
						PackageLicenseConcluded: "LicenseRef-testdata-firstparty-FIRST_PARTY_LICENSE",
					},
					{
						PackageName:             "testdata-firstparty-lib-liba.so.meta_lic",
						PackageVersion:          "NOASSERTION",
						PackageDownloadLocation: "NOASSERTION",
						PackageSPDXIdentifier:   common.ElementID("testdata-firstparty-lib-liba.so.meta_lic"),
						PackageLicenseConcluded: "LicenseRef-testdata-firstparty-FIRST_PARTY_LICENSE",
					},
					{
						PackageName:             "testdata-firstparty-lib-libb.so.meta_lic",
						PackageVersion:          "NOASSERTION",
						PackageDownloadLocation: "NOASSERTION",
						PackageSPDXIdentifier:   common.ElementID("testdata-firstparty-lib-libb.so.meta_lic"),
						PackageLicenseConcluded: "LicenseRef-testdata-firstparty-FIRST_PARTY_LICENSE",
					},
					{
						PackageName:             "testdata-firstparty-lib-libc.a.meta_lic",
						PackageVersion:          "NOASSERTION",
						PackageDownloadLocation: "NOASSERTION",
						PackageSPDXIdentifier:   common.ElementID("testdata-firstparty-lib-libc.a.meta_lic"),
						PackageLicenseConcluded: "LicenseRef-testdata-firstparty-FIRST_PARTY_LICENSE",
					},
					{
						PackageName:             "testdata-firstparty-lib-libd.so.meta_lic",
						PackageVersion:          "NOASSERTION",
						PackageDownloadLocation: "NOASSERTION",
						PackageSPDXIdentifier:   common.ElementID("testdata-firstparty-lib-libd.so.meta_lic"),
						PackageLicenseConcluded: "LicenseRef-testdata-firstparty-FIRST_PARTY_LICENSE",
					},
				},
				Relationships: []*spdx.Relationship{
					{
						RefA:         common.MakeDocElementID("", "DOCUMENT"),
						RefB:         common.MakeDocElementID("", "testdata-firstparty-highest.apex.meta_lic"),
						Relationship: "DESCRIBES",
					},
					{
						RefA:         common.MakeDocElementID("", "testdata-firstparty-highest.apex.meta_lic"),
						RefB:         common.MakeDocElementID("", "testdata-firstparty-bin-bin1.meta_lic"),
						Relationship: "CONTAINS",
					},
					{
						RefA:         common.MakeDocElementID("", "testdata-firstparty-highest.apex.meta_lic"),
						RefB:         common.MakeDocElementID("", "testdata-firstparty-bin-bin2.meta_lic"),
						Relationship: "CONTAINS",
					},
					{
						RefA:         common.MakeDocElementID("", "testdata-firstparty-highest.apex.meta_lic"),
						RefB:         common.MakeDocElementID("", "testdata-firstparty-lib-liba.so.meta_lic"),
						Relationship: "CONTAINS",
					},
					{
						RefA:         common.MakeDocElementID("", "testdata-firstparty-highest.apex.meta_lic"),
						RefB:         common.MakeDocElementID("", "testdata-firstparty-lib-libb.so.meta_lic"),
						Relationship: "CONTAINS",
					},
					{
						RefA:         common.MakeDocElementID("", "testdata-firstparty-bin-bin1.meta_lic"),
						RefB:         common.MakeDocElementID("", "testdata-firstparty-lib-liba.so.meta_lic"),
						Relationship: "CONTAINS",
					},
					{
						RefA:         common.MakeDocElementID("", "testdata-firstparty-bin-bin1.meta_lic"),
						RefB:         common.MakeDocElementID("", "testdata-firstparty-lib-libc.a.meta_lic"),
						Relationship: "CONTAINS",
					},
					{
						RefA:         common.MakeDocElementID("", "testdata-firstparty-lib-libb.so.meta_lic"),
						RefB:         common.MakeDocElementID("", "testdata-firstparty-bin-bin2.meta_lic"),
						Relationship: "RUNTIME_DEPENDENCY_OF",
					},
					{
						RefA:         common.MakeDocElementID("", "testdata-firstparty-lib-libd.so.meta_lic"),
						RefB:         common.MakeDocElementID("", "testdata-firstparty-bin-bin2.meta_lic"),
						Relationship: "RUNTIME_DEPENDENCY_OF",
					},
				},
				OtherLicenses: []*spdx.OtherLicense{
					{
						LicenseIdentifier: "LicenseRef-testdata-firstparty-FIRST_PARTY_LICENSE",
						ExtractedText:     "&&&First Party License&&&",
						LicenseName:       "testdata-firstparty-FIRST_PARTY_LICENSE",
					},
				},
			},
			expectedDeps: []string{
				"testdata/firstparty/FIRST_PARTY_LICENSE",
				"testdata/firstparty/bin/bin1.meta_lic",
				"testdata/firstparty/bin/bin2.meta_lic",
				"testdata/firstparty/highest.apex.meta_lic",
				"testdata/firstparty/lib/liba.so.meta_lic",
				"testdata/firstparty/lib/libb.so.meta_lic",
				"testdata/firstparty/lib/libc.a.meta_lic",
				"testdata/firstparty/lib/libd.so.meta_lic",
			},
		},
		{
			condition: "notice",
			name:      "binary",
			roots:     []string{"bin/bin1.meta_lic"},
			expectedOut: &spdx.Document{
				SPDXIdentifier: "DOCUMENT",
				CreationInfo:   getCreationInfo(t),
				Packages: []*spdx.Package{
					{
						PackageName:             "testdata-notice-bin-bin1.meta_lic",
						PackageVersion:          "NOASSERTION",
						PackageDownloadLocation: "NOASSERTION",
						PackageSPDXIdentifier:   common.ElementID("testdata-notice-bin-bin1.meta_lic"),
						PackageLicenseConcluded: "LicenseRef-testdata-firstparty-FIRST_PARTY_LICENSE",
					},
					{
						PackageName:             "testdata-notice-lib-liba.so.meta_lic",
						PackageVersion:          "NOASSERTION",
						PackageDownloadLocation: "NOASSERTION",
						PackageSPDXIdentifier:   common.ElementID("testdata-notice-lib-liba.so.meta_lic"),
						PackageLicenseConcluded: "LicenseRef-testdata-notice-NOTICE_LICENSE",
					},
					{
						PackageName:             "testdata-notice-lib-libc.a.meta_lic",
						PackageVersion:          "NOASSERTION",
						PackageDownloadLocation: "NOASSERTION",
						PackageSPDXIdentifier:   common.ElementID("testdata-notice-lib-libc.a.meta_lic"),
						PackageLicenseConcluded: "LicenseRef-testdata-notice-NOTICE_LICENSE",
					},
				},
				Relationships: []*spdx.Relationship{
					{
						RefA:         common.MakeDocElementID("", "DOCUMENT"),
						RefB:         common.MakeDocElementID("", "testdata-notice-bin-bin1.meta_lic"),
						Relationship: "DESCRIBES",
					},
					{
						RefA:         common.MakeDocElementID("", "testdata-notice-bin-bin1.meta_lic"),
						RefB:         common.MakeDocElementID("", "testdata-notice-lib-liba.so.meta_lic"),
						Relationship: "CONTAINS",
					},
					{
						RefA:         common.MakeDocElementID("", "testdata-notice-bin-bin1.meta_lic"),
						RefB:         common.MakeDocElementID("", "testdata-notice-lib-libc.a.meta_lic"),
						Relationship: "CONTAINS",
					},
				},
				OtherLicenses: []*spdx.OtherLicense{
					{
						LicenseIdentifier: "LicenseRef-testdata-firstparty-FIRST_PARTY_LICENSE",
						ExtractedText:     "&&&First Party License&&&",
						LicenseName:       "testdata-firstparty-FIRST_PARTY_LICENSE",
					},
					{
						LicenseIdentifier: "LicenseRef-testdata-notice-NOTICE_LICENSE",
						ExtractedText:     "%%%Notice License%%%",
						LicenseName:       "testdata-notice-NOTICE_LICENSE",
					},
				},
			},
			expectedDeps: []string{
				"testdata/firstparty/FIRST_PARTY_LICENSE",
				"testdata/notice/NOTICE_LICENSE",
				"testdata/notice/bin/bin1.meta_lic",
				"testdata/notice/lib/liba.so.meta_lic",
				"testdata/notice/lib/libc.a.meta_lic",
			},
		},
		{
			condition: "reciprocal",
			name:      "application",
			roots:     []string{"application.meta_lic"},
			expectedOut: &spdx.Document{
				SPDXIdentifier: "DOCUMENT",
				CreationInfo:   getCreationInfo(t),
				Packages: []*spdx.Package{
					{
						PackageName:             "testdata-reciprocal-application.meta_lic",
						PackageVersion:          "NOASSERTION",
						PackageDownloadLocation: "NOASSERTION",
						PackageSPDXIdentifier:   common.ElementID("testdata-reciprocal-application.meta_lic"),
						PackageLicenseConcluded: "LicenseRef-testdata-firstparty-FIRST_PARTY_LICENSE",
					},
					{
						PackageName:             "testdata-reciprocal-bin-bin3.meta_lic",
						PackageVersion:          "NOASSERTION",
						PackageDownloadLocation: "NOASSERTION",
						PackageSPDXIdentifier:   common.ElementID("testdata-reciprocal-bin-bin3.meta_lic"),
						PackageLicenseConcluded: "LicenseRef-testdata-notice-NOTICE_LICENSE",
					},
					{
						PackageName:             "testdata-reciprocal-lib-liba.so.meta_lic",
						PackageVersion:          "NOASSERTION",
						PackageDownloadLocation: "NOASSERTION",
						PackageSPDXIdentifier:   common.ElementID("testdata-reciprocal-lib-liba.so.meta_lic"),
						PackageLicenseConcluded: "LicenseRef-testdata-reciprocal-RECIPROCAL_LICENSE",
					},
					{
						PackageName:             "testdata-reciprocal-lib-libb.so.meta_lic",
						PackageVersion:          "NOASSERTION",
						PackageDownloadLocation: "NOASSERTION",
						PackageSPDXIdentifier:   common.ElementID("testdata-reciprocal-lib-libb.so.meta_lic"),
						PackageLicenseConcluded: "LicenseRef-testdata-firstparty-FIRST_PARTY_LICENSE",
					},
				},
				Relationships: []*spdx.Relationship{
					{
						RefA:         common.MakeDocElementID("", "DOCUMENT"),
						RefB:         common.MakeDocElementID("", "testdata-reciprocal-application.meta_lic"),
						Relationship: "DESCRIBES",
					},
					{
						RefA:         common.MakeDocElementID("", "testdata-reciprocal-bin-bin3.meta_lic"),
						RefB:         common.MakeDocElementID("", "testdata-reciprocal-application.meta_lic"),
						Relationship: "BUILD_TOOL_OF",
					},
					{
						RefA:         common.MakeDocElementID("", "testdata-reciprocal-application.meta_lic"),
						RefB:         common.MakeDocElementID("", "testdata-reciprocal-lib-liba.so.meta_lic"),
						Relationship: "CONTAINS",
					},
					{
						RefA:         common.MakeDocElementID("", "testdata-reciprocal-lib-libb.so.meta_lic"),
						RefB:         common.MakeDocElementID("", "testdata-reciprocal-application.meta_lic"),
						Relationship: "RUNTIME_DEPENDENCY_OF",
					},
				},
				OtherLicenses: []*spdx.OtherLicense{
					{
						LicenseIdentifier: "LicenseRef-testdata-firstparty-FIRST_PARTY_LICENSE",
						ExtractedText:     "&&&First Party License&&&",
						LicenseName:       "testdata-firstparty-FIRST_PARTY_LICENSE",
					},
					{
						LicenseIdentifier: "LicenseRef-testdata-notice-NOTICE_LICENSE",
						ExtractedText:     "%%%Notice License%%%",
						LicenseName:       "testdata-notice-NOTICE_LICENSE",
					},
					{
						LicenseIdentifier: "LicenseRef-testdata-reciprocal-RECIPROCAL_LICENSE",
						ExtractedText:     "$$$Reciprocal License$$$",
						LicenseName:       "testdata-reciprocal-RECIPROCAL_LICENSE",
					},
				},
			},
			expectedDeps: []string{
				"testdata/firstparty/FIRST_PARTY_LICENSE",
				"testdata/notice/NOTICE_LICENSE",
				"testdata/reciprocal/RECIPROCAL_LICENSE",
				"testdata/reciprocal/application.meta_lic",
				"testdata/reciprocal/bin/bin3.meta_lic",
				"testdata/reciprocal/lib/liba.so.meta_lic",
				"testdata/reciprocal/lib/libb.so.meta_lic",
			},
		},
		{
			condition: "restricted",
			name:      "library",
			roots:     []string{"lib/libd.so.meta_lic"},
			expectedOut: &spdx.Document{
				SPDXIdentifier: "DOCUMENT",
				CreationInfo:   getCreationInfo(t),
				Packages: []*spdx.Package{
					{
						PackageName:             "testdata-restricted-lib-libd.so.meta_lic",
						PackageVersion:          "NOASSERTION",
						PackageDownloadLocation: "NOASSERTION",
						PackageSPDXIdentifier:   common.ElementID("testdata-restricted-lib-libd.so.meta_lic"),
						PackageLicenseConcluded: "LicenseRef-testdata-notice-NOTICE_LICENSE",
					},
				},
				Relationships: []*spdx.Relationship{
					{
						RefA:         common.MakeDocElementID("", "DOCUMENT"),
						RefB:         common.MakeDocElementID("", "testdata-restricted-lib-libd.so.meta_lic"),
						Relationship: "DESCRIBES",
					},
				},
				OtherLicenses: []*spdx.OtherLicense{
					{
						LicenseIdentifier: "LicenseRef-testdata-notice-NOTICE_LICENSE",
						ExtractedText:     "%%%Notice License%%%",
						LicenseName:       "testdata-notice-NOTICE_LICENSE",
					},
				},
			},
			expectedDeps: []string{
				"testdata/notice/NOTICE_LICENSE",
				"testdata/restricted/lib/libd.so.meta_lic",
			},
		},
		{
			condition: "proprietary",
			name:      "container",
			roots:     []string{"container.zip.meta_lic"},
			expectedOut: &spdx.Document{
				SPDXIdentifier: "DOCUMENT",
				CreationInfo:   getCreationInfo(t),
				Packages: []*spdx.Package{
					{
						PackageName:             "testdata-proprietary-container.zip.meta_lic",
						PackageVersion:          "NOASSERTION",
						PackageDownloadLocation: "NOASSERTION",
						PackageSPDXIdentifier:   common.ElementID("testdata-proprietary-container.zip.meta_lic"),
						PackageLicenseConcluded: "LicenseRef-testdata-firstparty-FIRST_PARTY_LICENSE",
					},
					{
						PackageName:             "testdata-proprietary-bin-bin1.meta_lic",
						PackageVersion:          "NOASSERTION",
						PackageDownloadLocation: "NOASSERTION",
						PackageSPDXIdentifier:   common.ElementID("testdata-proprietary-bin-bin1.meta_lic"),
						PackageLicenseConcluded: "LicenseRef-testdata-firstparty-FIRST_PARTY_LICENSE",
					},
					{
						PackageName:             "testdata-proprietary-bin-bin2.meta_lic",
						PackageVersion:          "NOASSERTION",
						PackageDownloadLocation: "NOASSERTION",
						PackageSPDXIdentifier:   common.ElementID("testdata-proprietary-bin-bin2.meta_lic"),
						PackageLicenseConcluded: "LicenseRef-testdata-proprietary-PROPRIETARY_LICENSE",
					},
					{
						PackageName:             "testdata-proprietary-lib-liba.so.meta_lic",
						PackageVersion:          "NOASSERTION",
						PackageDownloadLocation: "NOASSERTION",
						PackageSPDXIdentifier:   common.ElementID("testdata-proprietary-lib-liba.so.meta_lic"),
						PackageLicenseConcluded: "LicenseRef-testdata-proprietary-PROPRIETARY_LICENSE",
					},
					{
						PackageName:             "testdata-proprietary-lib-libb.so.meta_lic",
						PackageVersion:          "NOASSERTION",
						PackageDownloadLocation: "NOASSERTION",
						PackageSPDXIdentifier:   common.ElementID("testdata-proprietary-lib-libb.so.meta_lic"),
						PackageLicenseConcluded: "LicenseRef-testdata-restricted-RESTRICTED_LICENSE",
					},
					{
						PackageName:             "testdata-proprietary-lib-libc.a.meta_lic",
						PackageVersion:          "NOASSERTION",
						PackageDownloadLocation: "NOASSERTION",
						PackageSPDXIdentifier:   common.ElementID("testdata-proprietary-lib-libc.a.meta_lic"),
						PackageLicenseConcluded: "LicenseRef-testdata-proprietary-PROPRIETARY_LICENSE",
					},
					{
						PackageName:             "testdata-proprietary-lib-libd.so.meta_lic",
						PackageVersion:          "NOASSERTION",
						PackageDownloadLocation: "NOASSERTION",
						PackageSPDXIdentifier:   common.ElementID("testdata-proprietary-lib-libd.so.meta_lic"),
						PackageLicenseConcluded: "LicenseRef-testdata-notice-NOTICE_LICENSE",
					},
				},
				Relationships: []*spdx.Relationship{
					{
						RefA:         common.MakeDocElementID("", "DOCUMENT"),
						RefB:         common.MakeDocElementID("", "testdata-proprietary-container.zip.meta_lic"),
						Relationship: "DESCRIBES",
					},
					{
						RefA:         common.MakeDocElementID("", "testdata-proprietary-container.zip.meta_lic"),
						RefB:         common.MakeDocElementID("", "testdata-proprietary-bin-bin1.meta_lic"),
						Relationship: "CONTAINS",
					},
					{
						RefA:         common.MakeDocElementID("", "testdata-proprietary-container.zip.meta_lic"),
						RefB:         common.MakeDocElementID("", "testdata-proprietary-bin-bin2.meta_lic"),
						Relationship: "CONTAINS",
					},
					{
						RefA:         common.MakeDocElementID("", "testdata-proprietary-container.zip.meta_lic"),
						RefB:         common.MakeDocElementID("", "testdata-proprietary-lib-liba.so.meta_lic"),
						Relationship: "CONTAINS",
					},
					{
						RefA:         common.MakeDocElementID("", "testdata-proprietary-container.zip.meta_lic"),
						RefB:         common.MakeDocElementID("", "testdata-proprietary-lib-libb.so.meta_lic"),
						Relationship: "CONTAINS",
					},
					{
						RefA:         common.MakeDocElementID("", "testdata-proprietary-bin-bin1.meta_lic"),
						RefB:         common.MakeDocElementID("", "testdata-proprietary-lib-liba.so.meta_lic"),
						Relationship: "CONTAINS",
					},
					{
						RefA:         common.MakeDocElementID("", "testdata-proprietary-bin-bin1.meta_lic"),
						RefB:         common.MakeDocElementID("", "testdata-proprietary-lib-libc.a.meta_lic"),
						Relationship: "CONTAINS",
					},
					{
						RefA:         common.MakeDocElementID("", "testdata-proprietary-lib-libb.so.meta_lic"),
						RefB:         common.MakeDocElementID("", "testdata-proprietary-bin-bin2.meta_lic"),
						Relationship: "RUNTIME_DEPENDENCY_OF",
					},
					{
						RefA:         common.MakeDocElementID("", "testdata-proprietary-lib-libd.so.meta_lic"),
						RefB:         common.MakeDocElementID("", "testdata-proprietary-bin-bin2.meta_lic"),
						Relationship: "RUNTIME_DEPENDENCY_OF",
					},
				},
				OtherLicenses: []*spdx.OtherLicense{
					{
						LicenseIdentifier: "LicenseRef-testdata-firstparty-FIRST_PARTY_LICENSE",
						ExtractedText:     "&&&First Party License&&&",
						LicenseName:       "testdata-firstparty-FIRST_PARTY_LICENSE",
					},
					{
						LicenseIdentifier: "LicenseRef-testdata-notice-NOTICE_LICENSE",
						ExtractedText:     "%%%Notice License%%%",
						LicenseName:       "testdata-notice-NOTICE_LICENSE",
					},
					{
						LicenseIdentifier: "LicenseRef-testdata-proprietary-PROPRIETARY_LICENSE",
						ExtractedText:     "@@@Proprietary License@@@",
						LicenseName:       "testdata-proprietary-PROPRIETARY_LICENSE",
					},
					{
						LicenseIdentifier: "LicenseRef-testdata-restricted-RESTRICTED_LICENSE",
						ExtractedText:     "###Restricted License###",
						LicenseName:       "testdata-restricted-RESTRICTED_LICENSE",
					},
				},
			},
			expectedDeps: []string{
				"testdata/firstparty/FIRST_PARTY_LICENSE",
				"testdata/notice/NOTICE_LICENSE",
				"testdata/proprietary/PROPRIETARY_LICENSE",
				"testdata/proprietary/bin/bin1.meta_lic",
				"testdata/proprietary/bin/bin2.meta_lic",
				"testdata/proprietary/container.zip.meta_lic",
				"testdata/proprietary/lib/liba.so.meta_lic",
				"testdata/proprietary/lib/libb.so.meta_lic",
				"testdata/proprietary/lib/libc.a.meta_lic",
				"testdata/proprietary/lib/libd.so.meta_lic",
				"testdata/restricted/RESTRICTED_LICENSE",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.condition+" "+tt.name, func(t *testing.T) {
			stdout := &bytes.Buffer{}
			stderr := &bytes.Buffer{}

			rootFiles := make([]string, 0, len(tt.roots))
			for _, r := range tt.roots {
				rootFiles = append(rootFiles, "testdata/"+tt.condition+"/"+r)
			}

			ctx := context{stdout, stderr, compliance.GetFS(tt.outDir), "Android", []string{tt.stripPrefix}, fakeTime}

			spdxDoc, deps, err := sbomGenerator(&ctx, rootFiles...)
			if err != nil {
				t.Fatalf("sbom: error = %v, stderr = %v", err, stderr)
				return
			}
			if stderr.Len() > 0 {
				t.Errorf("sbom: gotStderr = %v, want none", stderr)
			}

			gotData, err := json.Marshal(spdxDoc)
			if err != nil {
				t.Fatalf("sbom: failed to marshal spdx doc: %v", err)
				return
			}

			t.Logf("Got SPDX Doc: %s", string(gotData))

			expectedData, err := json.Marshal(tt.expectedOut)
			if err != nil {
				t.Fatalf("sbom: failed to marshal spdx doc: %v", err)
				return
			}

			t.Logf("Want SPDX Doc: %s", string(expectedData))

			if !compareSpdxDocs(t, spdxDoc, tt.expectedOut) {
				t.Errorf("SBOM: Test failed!")
			}

			t.Logf("got deps: %q", deps)

			t.Logf("want deps: %q", tt.expectedDeps)

			if g, w := deps, tt.expectedDeps; !reflect.DeepEqual(g, w) {
				t.Errorf("unexpected deps, wanted:\n%s\ngot:\n%s\n",
					strings.Join(w, "\n"), strings.Join(g, "\n"))
			}
		})
	}
}

func getCreationInfo(t *testing.T) *spdx.CreationInfo {
	ci, err := builder2v2.BuildCreationInfoSection2_2("Organization", "Google LLC", nil)
	if err != nil {
		t.Error("Unable to get creation info: %v", err)
		return nil
	}
	return ci
}

// compareSpdxDocs deep-compares two spdx docs by going through the info section, packages, relationships and licenses
func compareSpdxDocs(t *testing.T, actual, expected *spdx.Document) bool {

	if actual == nil || expected == nil {
		return actual == expected
	}
	// compare creation info
	if !compareSpdxCreationInfo(t, actual.CreationInfo, expected.CreationInfo) {
		return false
	}

	// compare packages
	if len(actual.Packages) != len(expected.Packages) {
		t.Error("Error: Number of Packages in actual is different! Got %d: Expected %d", len(actual.Packages), len(expected.Packages))
		return false
	}

	for i, pkg := range actual.Packages {
		if !compareSpdxPackages(t, pkg, expected.Packages[i]) {
			return false
		}
	}

	// compare licenses
	if len(actual.OtherLicenses) != len(expected.OtherLicenses) {
		t.Error("Error: Number of Licenses in actual is different! Got %d: Expected %d", len(actual.OtherLicenses), len(expected.OtherLicenses))
		return false
	}
	for i, license := range actual.OtherLicenses {
		if !compareLicenses(t, license, expected.OtherLicenses[i]) {
			return false
		}
	}

	//compare Relationships
	if len(actual.Relationships) != len(expected.Relationships) {
		t.Error("Error: Number of Licenses in actual is different! Got %d: Expected %d", len(actual.Relationships), len(expected.Relationships))
		return false
	}
	for i, rl := range actual.Relationships {
		if !compareRelationShips(t, rl, expected.Relationships[i]) {
			t.Error("SBOM: Relationship Error! Got: %s, Want: %s", rl, expected.Relationships[i])
			return false
		}
	}

	return true
}

func compareSpdxCreationInfo(t *testing.T, actual, expected *spdx.CreationInfo) bool {
	if actual == nil || expected == nil {
		t.Error("SBOM: Creation info Error! Got %s: Expected %s", actual, expected)
		return actual == expected
	}

	if actual.LicenseListVersion != expected.LicenseListVersion {
		t.Error("SBOM: Creation info license version Error! Got %s: Expected %s", actual.LicenseListVersion, expected.LicenseListVersion)
		return false
	}

	if len(actual.Creators) != len(expected.Creators) {
		t.Error("SBOM: Creation info creators Error! Got %d: Expected %d", actual.Creators, expected.Creators)
		return false
	}

	for i, info := range actual.Creators {
		if info != expected.Creators[i] {
			t.Error("SBOM: Creation info creators Error! Got %s: Expected %s", info, expected.Creators[i])
			return false
		}
	}

	return true
}

func compareSpdxPackages(t *testing.T, actual, expected *spdx.Package) bool {
	if actual == nil || expected == nil {
		t.Error("SBOM: Packages Error! Got %s: Expected %s", actual, expected)
		return actual == expected
	}
	if actual.PackageName != expected.PackageName {
		t.Error("SBOM: Package name Error! Got %s: Expected %s", actual.PackageName, expected.PackageName)
		return false
	}

	if actual.PackageVersion != expected.PackageVersion {
		t.Error("SBOM: Package version Error! Got %s: Expected %s", actual.PackageVersion, expected.PackageVersion)
		return false
	}

	if actual.PackageSPDXIdentifier != expected.PackageSPDXIdentifier {
		t.Error("SBOM: Package identifier Error! Got %s: Expected %s", actual.PackageSPDXIdentifier, expected.PackageSPDXIdentifier)
		return false
	}

	if actual.PackageDownloadLocation != expected.PackageDownloadLocation {
		t.Error("SBOM: Package download location Error! Got %s: Expected %s", actual.PackageDownloadLocation, expected.PackageDownloadLocation)
		return false
	}

	if actual.PackageLicenseConcluded != expected.PackageLicenseConcluded {
		t.Error("SBOM: Package license concluded Error! Got %s: Expected %s", actual.PackageLicenseConcluded, expected.PackageLicenseConcluded)
		return false
	}

	return true
}

func compareRelationShips(t *testing.T, actual, expected *spdx.Relationship) bool {
	if actual == nil || expected == nil {
		t.Error("SBOM: Relationships Error! Got %s: Expected %s", actual, expected)
		return actual == expected
	}

	if actual.RefA != expected.RefA {
		t.Error("SBOM: Relationship RefA Error! Got %s: Expected %s", actual.RefA, expected.RefA)
		return false
	}

	if actual.RefB != expected.RefB {
		t.Error("SBOM: Relationship RefB Error! Got %s: Expected %s", actual.RefB, expected.RefB)
		return false
	}

	if actual.Relationship != expected.Relationship {
		t.Error("SBOM: Relationship type Error! Got %s: Expected %s", actual.Relationship, expected.Relationship)
		return false
	}

	return true
}

func compareLicenses(t *testing.T, actual, expected *spdx.OtherLicense) bool {
	if actual == nil || expected == nil {
		t.Error("SBOM: Licenses Error! Got %s: Expected %s", actual, expected)
		return actual == expected
	}

	if actual.LicenseName != expected.LicenseName {
		t.Error("SBOM: License Name Error! Got %s: Expected %s", actual.LicenseName, expected.LicenseName)
		return false
	}

	if actual.LicenseIdentifier != expected.LicenseIdentifier {
		t.Error("SBOM: License Identifier Error! Got %s: Expected %s", actual.LicenseIdentifier, expected.LicenseIdentifier)
		return false
	}

	if stripSpecialChars(actual.ExtractedText) != stripSpecialChars(expected.ExtractedText) {
		t.Error("SBOM: License Extracted Text Error! Got: %s", actual.ExtractedText, "want: %s", expected.ExtractedText)
		return false
	}

	return true
}

func stripSpecialChars(str string) string {
	// Remove all special characters except letters and digits
	re := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	stripped := re.ReplaceAllString(str, "")

	// Remove any non-letter characters before the first letter
	re = regexp.MustCompile(`^[^a-zA-Z]*([a-zA-Z])`)
	stripped = re.ReplaceAllString(stripped, "$1")

	// Remove any non-letter characters after the last letter
	re = regexp.MustCompile(`([a-zA-Z])[^a-zA-Z]*$`)
	stripped = re.ReplaceAllString(stripped, "$1")

	return stripped
}

func fakeTime() time.Time {
	return time.UnixMicro(0).UTC()
}
