package main

import (
	"bytes"
	"testing"
)

func Test_doit(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "One",
			in: `
[
  {
    "rule": "@//build/soong/licenses:Android-Apache-2.0",
    "license_kinds": [
      {
        "target": "@//build/soong/licenses:SPDX-license-identifier-Apache-2.0",
        "name": "SPDX-license-identifier-Apache-2.0",
        "conditions": ["notice"]
      }
    ],
    "copyright_notice": "Copyright (C) The Android Open Source Project",
    "package_name": "Discombobulator",
    "package_url": null,
    "package_version": null,
    "license_text": "../testdata/notice/NOTICE_LICENSE"
 },
 {
    "rule": "@//external/scudo:external_scudo_license",
    "license_kinds": [
      {
        "target": "@//build/soong/licenses:SPDX-license-identifier-Apache-2.0",
        "name": "SPDX-license-identifier-Apache-2.0",
        "conditions": ["notice"]
      }
    ],
    "copyright_notice": "",
    "package_name": "Scudo Standalone",
    "package_url": null,
    "package_version": null,
    "license_text": "external/scudo/LICENSE.TXT"
  }
]
`,
			want: `<!DOCTYPE html>
<html>
  <head>
    <style type="text/css">
      body { padding: 2px; margin: 0; }
      .license { background-color: seashell; margin: 1em;}
      pre { padding: 1em; }</style></head>
  <body>
    The following software has been included in this product and contains the license and notice as shown below.<p>
    <strong>Discombobulator</strong><br>Copyright Notice: Copyright (C) The Android Open Source Project<br><a href=#0e6553ab7221430a352fb7706ebc2aad>License</a><hr>
    <strong>Scudo Standalone</strong><hr>
    <div id="0e6553ab7221430a352fb7706ebc2aad" class="license"><pre>%%%Notice License%%%

    </pre></div>
  </body>
</html>
`,
		}, // TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := bytes.Buffer{}
			newGenerator(tt.in).generate(&buf)
			got := buf.String()
			if got != tt.want {
				t.Errorf("doit() = %v, want %v", got, tt.want)
			}
		})
	}
}
