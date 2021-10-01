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
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

var (
	ConcurrentReaders = 5
)

type LicenseSet interface {
}

type licenseSetImp struct {
	rootFiles []string
	metadata  map[string]*metadataFile
}

func newLicenseSetImp() *licenseSetImp {
	return &licenseSetImp{
		[]string{},
		make(map[string]*metadataFile),
	}
}

type installMap struct {
	prefix string
	replacement string
}

type metadataFile struct {
	packageName string
	licenseKinds []string
	licenseConditions []string
	licenseTexts []string
	isContainer bool
	built []string
	installed []string
	installMap []installMap
	sources []string
	effectiveConditions []string
	effectiveLicenseText []string
	deps []string
	subDeps []string
}

type result struct {
	file    string
	meta    *metadataFile
	err     error
	elapsed time.Duration
}

type receiver struct {
	task    chan bool
	results chan *result
	wg      sync.WaitGroup
}

var recv *receiver

func ReadLicenseMetadata(files []string) (LicenseSet, error) {
	globalStart := time.Now()
	var ls *licenseSetImp

	if len(files) == 0 {
		return ls, fmt.Errorf("no license metadata to analyze")
	}

	ls = newLicenseSetImp()
	for f := range files {
		if strings.HasSuffix(f, ".meta_lic") {
			ls.rootFiles = append(ls.rootFiles, f)
		} else {
			ls.rootFiles = append(ls.rootFiles, f + ".meta_lic")
		}
	}
	for f := range ls.rootFiles {
		ls.metadata[f] = nil
	}

	recv = &receiver{make(chan bool, ConcurrentReaders), make(chan *result, ConcurrentReaders), sync.WaitGroup{}}
	for i := 0; i < ConcurrentReaders; i++ {
		recv.task <- true
	}

