package main

import "testing"

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
    "package_name": "Android",
    "package_url": null,
    "package_version": null,
    "license_text": "build/soong/licenses/LICENSE"
  }
]
`,
			want: `
<!DOCTYPE html>
<html>
  <head>
    <style type="text/css"
	  body { padding: 2px; margin: 0; }
	  ul { list-style-type: none; margin: 0; padding: 0; }
	  li { padding-left: 1em; }
	  .file-list { margin-left: 1em; }
	</style>
	</head>
  <body>
  </body>
</html>
`,
		}, // TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := doit(tt.in); got != tt.want {
				t.Errorf("doit() = %v, want %v", got, tt.want)
			}
		})
	}
}
