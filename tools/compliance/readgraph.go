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

// result describes the outcome of reading and parsing a single license metadata file.
type result struct {
	// file identifies the path to the license metadata file
	file string
	// target contains the parsed metadata or nil if an error
	target *targetNode
	// edges contains the parsed dependencies
	edges []*dependencyEdge
	// err is nil unless an error occurs
	err error
}

// receiver coordinates the tasks for reading and parsing license metadata files.
type receiver struct {
	// ls accumulates the read metadata
	lg *licenseGraphImp
	// task provides a fixed-size task pool to limit concurrent open files etc.
	task chan bool
	// results returns one metadata file result at a time
	results chan *result
	// wg detects when done
	wg sync.WaitGroup
}

// ReadLicenseGraph reads and parses `files` and their dependencies into a LicenseGraph.
func ReadLicenseGraph(files []string) (LicenseGraph, error) {
	if len(files) == 0 {
		return nil, fmt.Errorf("no license metadata to analyze")
	}

	lg := newLicenseGraphImp()
	for _, f := range files {
		if strings.HasSuffix(f, ".meta_lic") {
			lg.rootFiles = append(lg.rootFiles, f)
		} else {
			lg.rootFiles = append(lg.rootFiles, f+".meta_lic")
		}
	}

	recv := &receiver{lg: lg, task: make(chan bool, ConcurrentReaders), results: make(chan *result, ConcurrentReaders), wg: sync.WaitGroup{}}
	for i := 0; i < ConcurrentReaders; i++ {
		recv.task <- true
	}

	readFiles := func() {
		// identify the metadata files to schedule reading tasks for
		for _, f := range lg.rootFiles {
			lg.targets[f] = nil
		}

		// schedule tasks to read the files
		for _, f := range lg.rootFiles {
			recv.wg.Add(1)
			<-recv.task
			go readFile(recv, f)
		}

		// schedule a task to wait until finished and close the channel.
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
					lg = nil
					recv.results = nil
					continue
				}

				// record the parsed metadata (guarded by mutex)
				recv.lg.mu.Lock()
				recv.lg.targets[r.file] = r.target
				if len(r.edges) > 0 {
					recv.lg.edges = append(recv.lg.edges, r.edges...)
				}
				recv.lg.mu.Unlock()
			} else {
				// finished -- nil the results channel
				recv.results = nil
			}
		}
	}

	return lg, err
}

// installMap describes a pair of strings where `prefix` at the start of a source path is substituted with `replacement`.
type installMap struct {
	prefix      string
	replacement string
}

type targetNode struct {
	// name is the path to the metadata file
	name string
	// packageName identifies the source package. License texts are named relative to the package name.
	packageName string
	// moduleType identifies the module type, if given.
	moduleType []string
	// moduleClass identifies the module class, if given.
	moduleClass []string
	// projects identifies the git project(s) containing the associated source code.
	projects []string
	// licenseKinds lists the kinds of licenses. e.g. SPDX-license-identifier-Apache-2.0 or legacy_notice
	licenseKinds []string
	// licenseConditions lists the conditions that apply to the license kinds. e.g. notice or restricted
	licenseConditions *licenseConditionSetImp
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
}

// Add sets or appends `value` to the field identified by `name`.
func (tn *targetNode) Add(name, value string) error {
	switch name {
	case "license_package_name":
		if len(tn.packageName) > 0 {
			return fmt.Errorf("too many package names %q and %q", tn.packageName, value)
		}
		tn.packageName = value
	case "module_type":
		tn.moduleType = append(tn.moduleType, value)
	case "module_class":
		tn.moduleClass = append(tn.moduleClass, value)
	case "root":
		tn.projects = append(tn.projects, value)
	case "license_kind":
		tn.licenseKinds = append(tn.licenseKinds, value)
	case "license_condition":
		tn.licenseConditions.add(value, targetNodeImp{nil, tn.name})
	case "license_text":
		tn.licenseTexts = append(tn.licenseTexts, value)
	case "is_container":
		if value == "true" {
			tn.isContainer = true
		} else if value == "false" {
			tn.isContainer = false
		} else {
			return fmt.Errorf("invalid boolean is_container %q", value)
		}
	case "built":
		tn.built = append(tn.built, value)
	case "installed":
		tn.installed = append(tn.installed, value)
	case "install_map":
		if matches := colonPair.FindStringSubmatch(value); matches != nil {
			tn.installMap = append(tn.installMap, installMap{matches[1], matches[2]})
		} else {
			return fmt.Errorf("invalid install map %q", value)
		}
	case "source":
		tn.sources = append(tn.sources, value)
	default:
		return fmt.Errorf("unknown metadata key %q for value %q", name, value)
	}
	return nil
}

type dependencyEdge struct {
	// target identifies the target node being built and/or installed.
	target string
	// dependency identifies the target node being depended on.
	dependency string
	// annotations are a set of text attributes attached to the edge.
	annotations targetEdgeAnnotationsImp
}

// addDependency parses a dependency from `tn` and adds it to `edges`
func addDependency(edges *[]*dependencyEdge, target string, value string) error {
	fields := strings.Split(value, ":")
	dependency := fields[0]
	if len(dependency) == 0 {
		return fmt.Errorf("invalid deps %q", value)
	}
	annotations := make(targetEdgeAnnotationsImp)
	for _, a := range fields[1:] {
		if len(a) == 0 {
			continue
		}
		annotations[a] = nil
	}
	*edges = append(*edges, &dependencyEdge{target, dependency, annotations})
	return nil
}

// readFile is a task to read and parse a single license metadata file, and to schedule
// additional tasks for reading and parsing dependencies as necessary.
func readFile(recv *receiver, file string) {
	// read the file
	f, err := os.Open(file)
	if err != nil {
		recv.results <- &result{file, nil, nil, err}
		fmt.Fprintf(os.Stderr, "cannot open file %q: %q", file, err)
		return
	}
	defer func() {
		f.Close()
	}()

	tn := &targetNode{name: file, licenseConditions: newLicenseConditionSet(&targetNodeImp{recv.lg, file})}
	edges := []*dependencyEdge{}

	scanner := bufio.NewScanner(f)
	for line := 1; scanner.Scan(); line++ {
		l := scanner.Text()

		// skip blank lines
		if isAllWS.MatchString(l) {
			continue
		}

		// parse 'name: value' or 'name: "quoted value"'
		if matches := parseMetaLine.FindStringSubmatch(l); matches != nil {
			if matches[1] == "dep" {
				// add dependency to `edges`
				err = addDependency(&edges, file, matches[2])
			} else {
				// add name/value pair to `targetNode`
				err = tn.Add(matches[1], matches[2])
			}
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
			recv.results <- &result{file, nil, nil, err}
			fmt.Fprintf(os.Stderr, "cannot parse file %q: %q", file, err)
			return
		}
	}
	err = scanner.Err()
	if err != nil && err != io.EOF {
		// handle scanner error if necessary
		recv.results <- &result{file, nil, nil, err}
		fmt.Fprintf(os.Stderr, "cannot read file %q: %q", file, err)
		return
	}

	// send result for this file and release task before scheduling dependencies,
	// but do not signal done to WaitGroup until dependencies are scheduled.
	recv.results <- &result{file, tn, edges, nil}
	recv.task <- true

	// schedule tasks as necessary to read dependencies
	for _, e := range edges {
		// decide and record whether to schedul task in critical section
		recv.lg.mu.Lock()
		_, alreadyScheduled := recv.lg.targets[e.dependency]
		if !alreadyScheduled {
			recv.lg.targets[e.dependency] = nil
		}
		recv.lg.mu.Unlock()
		// schedule task to read dependency file outside critical section
		if !alreadyScheduled {
			recv.wg.Add(1)
			<-recv.task
			go readFile(recv, e.dependency)
		}
	}

	// signal done after scheduling dependencies
	recv.wg.Done()
}
