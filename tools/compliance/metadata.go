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
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"sync"
)

var (
	// ConcurrentReaders is the size of the task pool for limiting resource usage e.g. open files.
	ConcurrentReaders = 5

	// isAllWS is a regular expression to match strings that have no non-whitespace characters.
	isAllWS = regexp.MustCompile("^\\s*$")
	// parseMetaLine is a regular expression to parse a license metadata file line into name and value.
	parseMetaLine = regexp.MustCompile("^([^:]*)[:]\\s\\s*\"?([^\"]*)\"?$")
	// colonPair is a regular expression to parse a colon-separated pair of strings.
	colonPair = regexp.MustCompile("^([^:]*)[:](.*)$")
)

// LicenseSet describes the metadata for a set of licenses.
type LicenseSet interface {
}

// licenseSetImp implements the LicenseSet interface.
type licenseSetImp struct {
	// rootFiles identifies the original set of files to read (immutable)
	rootFiles []string
	// metadata identifies the requested files (initially nil) and the parsed results (guarded by recv.mu)
	metadata map[string]*metadataFile
}

// newLicenseSetImp constructs a new instance of licenseSetImp.
func newLicenseSetImp() *licenseSetImp {
	return &licenseSetImp{
		[]string{},
		make(map[string]*metadataFile),
	}
}

// installMap describes a pair of strings where `prefix` at the start of a source path is substituted with `replacement`.
type installMap struct {
	prefix      string
	replacement string
}

// metadataFile describes the contents of a license metadata file.
type metadataFile struct {
	// packageName identifies the source package. License texts are named relative to the package name.
	packageName string
	// projects identifies the git project(s) containing the associated source code.
	projects []string
	// licenseKinds lists the kinds of licenses. e.g. SPDX-license-identifier-Apache-2.0 or legacy_notice
	licenseKinds []string
	// licenseConditions lists the conditions that apply to the license kinds. e.g. notice or restricted
	licenseConditions []string
	// licenseTexts lists the filenames of the associated license text(s).
	licenseTexts []string
	// isContainer is true for target types that merely aggregate. e.g. .img or .zip files
	isContainer bool
	// built lists the built targets
	built []string
	// installed lists the installed targets
	installed []string
	// installMap identifies the substitutions to make to path names when moving into installed location
	installMap []installMap
	// sources lists the targets depended on
	sources []string
	// effectiveLicenseConditions lists the conditions and the inherited condition(s). e.g. restricted
	effectiveLicenseConditions []string
	// effectiveLicenseTexts lists the transitive closure of license texts including dependencies.
	effectiveLicenseTexts []string
	// deps lists the license metadata files depended on
	deps []string
}

// Add sets or appends `value` to the field identified by `name`.
func (mf *metadataFile) Add(name, value string) error {
	switch name {
	case "license_package_name":
		if len(mf.packageName) > 0 {
			return fmt.Errorf("too many package names %q and %q", mf.packageName, value)
		}
		mf.packageName = value
	case "root":
		mf.projects = append(mf.projects, value)
	case "license_kind":
		mf.licenseKinds = append(mf.licenseKinds, value)
	case "license_condition":
		mf.licenseConditions = append(mf.licenseConditions, value)
	case "license_text":
		mf.licenseTexts = append(mf.licenseTexts, value)
	case "is_container":
		if value == "true" {
			mf.isContainer = true
		} else if value == "false" {
			mf.isContainer = false
		} else {
			return fmt.Errorf("invalid boolean is_container %q", value)
		}
	case "built":
		mf.built = append(mf.built, value)
	case "installed":
		mf.installed = append(mf.installed, value)
	case "install_map":
		if matches := colonPair.FindStringSubmatch(value); matches != nil {
			mf.installMap = append(mf.installMap, installMap{matches[1], matches[2]})
		} else {
			return fmt.Errorf("invalid install map %q", value)
		}
	case "source":
		mf.sources = append(mf.sources, value)
	case "dep":
		mf.deps = append(mf.deps, value)
	default:
		return fmt.Errorf("unknown metadata key %q for value %q", name, value)
	}
	return nil
}

// result describes the outcome of reading and parsing a single license metadata file.
type result struct {
	// file identifies the path to the license metadata file
	file string
	// meta contains the parsed metadata or nil if an error
	meta *metadataFile
	// err is nil unless an error occurs
	err error
}

// receiver coordinates the tasks for reading and parsing license metadata files.
type receiver struct {
	// ls accumulates the read metadata (guarded by mu)
	ls *licenseSetImp
	// task provides a fixed-size task pool to limit concurrent open files etc.
	task chan bool
	// results returns one metadata file result at a time
	results chan *result
	// wg detects when done
	wg sync.WaitGroup
	// mu guards ls against concurrent update
	mu sync.Mutex
}

// ReadLicenseMetadata reads and parses `files` and their dependencies into a LicenseSet.
func ReadLicenseMetadata(files []string) (LicenseSet, error) {
	if len(files) == 0 {
		return nil, fmt.Errorf("no license metadata to analyze")
	}

	ls := newLicenseSetImp()
	for _, f := range files {
		if strings.HasSuffix(f, ".meta_lic") {
			ls.rootFiles = append(ls.rootFiles, f)
		} else {
			ls.rootFiles = append(ls.rootFiles, f+".meta_lic")
		}
	}

	recv := &receiver{ls: ls, task: make(chan bool, ConcurrentReaders), results: make(chan *result, ConcurrentReaders), wg: sync.WaitGroup{}}
	for i := 0; i < ConcurrentReaders; i++ {
		recv.task <- true
	}

	readFiles := func() {
		// identify the metadata files to schedule reading tasks for
		for _, f := range ls.rootFiles {
			ls.metadata[f] = nil
		}

		// schedule tasks to read the files
		for _, f := range ls.rootFiles {
			recv.wg.Add(1)
			<-recv.task
			go readFile(recv, f)
		}

		// schedule a task to wait until finished and close the channels.
		go func() {
			recv.wg.Wait()
			close(recv.task)
			close(recv.results)
		}()
	}
	go readFiles()

	// tasks to read license metadata files are scheduled; read and process results from channel
	var err error
	for recv.results != nil {
		select {
		case r, ok := <-recv.results:
			if ok {
				// handle errors by nil'ing ls, setting err, and clobbering results channel
				if r.err != nil {
					err = r.err
					ls = nil
					recv.results = nil
					continue
				}
				// record the parsed metadata guarded by mutex
				recv.mu.Lock()
				recv.ls.metadata[r.file] = r.meta
				recv.mu.Unlock()
			} else {
				// finished -- nil the results channel
				recv.results = nil
			}
		}
	}

	// single task has reference to `ls` at this point
	if ls != nil {
		// reading and parsing succeeded -- calculate effective conditions and license texts
		for f := range ls.metadata {
			flattenConditions(ls.metadata[f], ls)
			flattenTexts(ls.metadata[f], ls)
		}
	}
	return ls, err
}

// readFile is a task to read and parse a single license metadata file, and to schedule
// additional tasks for reading and parsing dependencies as necessary.
func readFile(recv *receiver, file string) {
	mf := &metadataFile{}

	// read the file
	f, err := os.Open(file)
	if err != nil {
		recv.results <- &result{file, nil, err}
		fmt.Fprintf(os.Stderr, "cannot open file %q: %q", file, err)
		return
	}
	defer func() {
		f.Close()
	}()

	scanner := bufio.NewScanner(f)
	for line := 1; scanner.Scan(); line++ {
		l := scanner.Text()

		// skip blank lines
		if isAllWS.MatchString(l) {
			continue
		}

		// parse 'name: value' or 'name: "quoted value"'
		if matches := parseMetaLine.FindStringSubmatch(l); matches != nil {
			// add name/value pair to `metadataFile`
			err = mf.Add(matches[1], matches[2])
			if err != nil {
				// add line number context
				err = fmt.Errorf("%v line %d: %q", err, line, l)
			}
		} else {
			// matched nil means line was not in expected pattern
			err = fmt.Errorf("unparseable line %d: %q", line, l)
		}

		// handle error result if necessary
		if err != nil {
			recv.results <- &result{file, nil, err}
			fmt.Fprintf(os.Stderr, "cannot parse file %q: %q", file, err)
			return
		}
	}
	err = scanner.Err()
	if err != nil && err != io.EOF {
		// handle scanner error if necessary
		recv.results <- &result{file, nil, err}
		fmt.Fprintf(os.Stderr, "cannot read file %q: %q", file, err)
		return
	}

	// send result for this file and release task before scheduling dependencies,
	// but do not signal done to WaitGroup until dependencies are scheduled.
	recv.results <- &result{file, mf, nil}
	recv.task <- true

	// schedule tasks as necessary to read dependencies
	for _, d := range mf.deps {
		// decide and record whether to schedul task in critical section
		recv.mu.Lock()
		_, alreadyScheduled := recv.ls.metadata[d]
		if !alreadyScheduled {
			recv.ls.metadata[d] = nil
		}
		recv.mu.Unlock()
		// schedule task to read dependency file outside critical section
		if !alreadyScheduled {
			recv.wg.Add(1)
			<-recv.task
			go readFile(recv, d)
		}
	}

	// signal done after scheduling dependencies
	recv.wg.Done()
}

// flattenConditions calculates `effectiveLicenseConditions` for `mf` inheriting 'restricted' as necessary.
func flattenConditions(mf *metadataFile, ls *licenseSetImp) {
	// copy license conditions
	mf.effectiveLicenseConditions = append(mf.effectiveLicenseConditions, mf.licenseConditions...)

	// check whether already has "restricted"
	for _, c := range mf.effectiveLicenseConditions {
		if strings.HasPrefix(c, "restrict") {
			return
		}
	}

	// add "restricted" if any dependency has "restricted"
	for _, d := range mf.deps {
		for _, c := range ls.metadata[d].licenseConditions {
			if strings.HasPrefix(c, "restrict") {
				mf.effectiveLicenseConditions = append(mf.effectiveLicenseConditions, "restricted")
				return
			}
		}
	}
}

// flattenTexts calculates the transitive closure of license texts for `mf` from its dependencies.
func flattenTexts(mf *metadataFile, ls *licenseSetImp) {
	// copy license texts
	mf.effectiveLicenseTexts = append(mf.effectiveLicenseTexts, mf.licenseTexts...)

	// indicate license texts already recorded
	recorded := make(map[string]bool)
	for _, t := range mf.effectiveLicenseTexts {
		recorded[t] = true
	}

	// add unrecorded license texts from dependencies
	for _, d := range mf.deps {
		for _, t := range ls.metadata[d].licenseTexts {
			if _, alreadyRecorded := recorded[t]; !alreadyRecorded {
				mf.effectiveLicenseTexts = append(mf.effectiveLicenseTexts, t)
				recorded[t] = true
			}
		}
	}
}
