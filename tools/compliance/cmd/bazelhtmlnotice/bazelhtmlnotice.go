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
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	inputFile  string
	outputFile = flag.String("o", "", "output file")
)

func quit(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func doit(in string) string {
	decoder := json.NewDecoder(strings.NewReader(in))
	decoder.DisallowUnknownFields() //useful to detect typos, e.g. in unit tests
	var licenses []License
	if err := decoder.Decode(&licenses); err != nil {
		quit(err)
	}
	return fmt.Sprintf("%#v", licenses)

}
func processArgs() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, `usage: bazelhtmlnotice -o <output> <input>`)
		flag.PrintDefaults()
		os.Exit(2)
	}
	flag.Parse()
	if *outputFile == "" || len(flag.Args()) != 1 {
		flag.Usage()
	}
	inputFile = flag.Arg(0)
}

type LicenseKind struct {
	Target     string   `json:"target"`
	Name       string   `json:"name"`
	Conditions []string `json:"conditions"`
}
type License struct {
	Rule            string        `json:"rule"`
	CopyrightNotice string        `json:"copyright_notice"`
	PackageName     string        `json:"package_name"`
	PackageUrl      string        `json:"package_url"`
	PackageVersion  string        `json:"package_version"`
	LicenseFile     string        `json:"license_text"`
	LicenseKinds    []LicenseKind `json:"license_kinds"`
}

func main() {
	ll := []License{
		{
			Rule:            "r",
			CopyrightNotice: "c",
			PackageName:     "p",
			PackageUrl:      "",
			PackageVersion:  "",
			LicenseFile:     "",
			LicenseKinds:    nil,
		},
	}
	json.NewEncoder(os.Stdout).Encode(ll)
	processArgs()
	data, err := os.ReadFile(inputFile)
	if err != nil {
		quit(err)
	}
	doit(string(data))
}
