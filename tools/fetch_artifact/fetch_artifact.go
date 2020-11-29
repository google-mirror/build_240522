// Copyright 2020 Google Inc. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var target = flag.String("target", "", "the target to fetch from")
var buildID = flag.String("build_id", "", "the build id to fetch from")
var artifact = flag.String("artifact", "", "the artifact to download")
var writeToStdout = flag.Bool("pipe", false, "if the output is written to stdout or not")

func errPrint(msg string) {
	fmt.Fprint(os.Stderr, msg)
	os.Exit(1)
}

func main() {
	flag.Parse()

	url := fmt.Sprintf("https://androidbuildinternal.googleapis.com/android/internal/build/v3/builds/%s/%s/attempts/latest/artifacts/%s/url", *buildID, *target, *artifact)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		errPrint(fmt.Sprintf("unable to build request %v", err))
	}
	req.Header.Set("Accept", "application/json")

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		errPrint(fmt.Sprintf("Unable to make request %s", err))
	}
	defer res.Body.Close()

	if res.Status != "200 OK" {
		body, _ := ioutil.ReadAll(res.Body)
		errPrint(fmt.Sprintf("Unable to download artifact: %s\n %s", res.Status, body))
	}

	if *writeToStdout {
		io.Copy(os.Stdout, res.Body)
		return
	}

	s := strings.Split(*artifact, "/")
	fileName := s[len(s)-1]
	f, err := os.Create(fileName)
	if err != nil {
		errPrint(fmt.Sprintf("Unable to create file %s", err))
	}
	defer f.Close()
	io.Copy(f, res.Body)
	fmt.Printf("File %s created.\n", fileName)
}
