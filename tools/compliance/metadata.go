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
	"sort"
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

// LicenseGraph describes the license metadata for a set of root targets.
type LicenseGraph interface {
	// WalkRestricted returns the set of metadata files to treat as 'restricted'.
	WalkRestricted() TargetSet
	// WalkDepsForCondition returns the set of metadata with license `condition`.
	WalkDepsForCondition(condition string) TargetSet
}

// TargetNode describes a target corresponding to a specific license metadata file.
type TargetNode interface {
	Projects() []string
}

// TargetSet describes a set of TargetNode.
type TargetSet interface {
	Add(node... TargetNode)
	Targets() []TargetNode
	Projects() []string
}

// TargetEdge describes an edge from a target to a dependency.
type TargetEdge interface {
	Target() TargetNode
	Dependency() TargetNode
	Annotations() TargetEdgeAnnotations
}

// TargetPath describes a sequence of edges where the dependency of one is the target of the next
// describing a path from the target of the 1st element to the dependency of the last element.
type TargetPath []TargetEdge

// TargetEdgeAnnotations describes a set of annotations attached to an edge from a target to a dependency.
type TargetEdgeAnnotations interface {
	HasAnnotation(ann string) bool
	ListAnnotations() []string
}

// targetNodeImp implements TargetNode
type targetNodeImp struct {
	lm   *licenseMetadataImp
	file string
}

func (lf TargetNode) Projects() []string {
	return append([]string{}, lf.lm.targets[lf.file].projects...)
}

// targetSetImp implements TargetSet
type targetSetImp struct {
	lm    *licenseMetadataImp
	files map[string]interface{}
}

func (ls TargetSet) Add(node... TargetNode) {
	for tn := range node {
		if ls.lm != tn.lm {
			panic(fmt.Errorf("attempt to Add target node from different metadata"))
		}
		ls.files[tn.file] = nil
	}
}

func (ls TargetSet) Targets() []TargetNode {
	var files []TargetNode
	for f := range ls.files {
		files = append(files, TargetNode{ls.lm, f})
	}
	return files
}

func (ls TargetSet) Projects() []string {
	pset := make(map[string]interface{})
	for f := range ls.files {
		for _, p := range ls.lm.targets[f].projects {
			pset[p] = nil
		}
	}
	projects := make([]string, 0, len(pset))
	for p := range pset {
		projects = append(projects, p)
	}
	sort.Strings(projects)
	return projects
}

// targetEdgeImp implements TargetEdge
type targetEdgeImp struct {
	lm *licenseMetadataImp
	e *dependencyEdge
}

func (e targetEdgeImp) Target() TargetNode {
	return TargetNode{e.lm, e.d.target}
}

func (e targetEdgeImp) Dependency() TargetNode {
	return TargetNode{e.lm, e.d.dependency}
}

func (e targetEdgeImp) Annotations() TargetEdgeAnnotations {
	return e.annotations
}

// targetEdgeAnnotationsImp implements TargetEdgeAnnotations
type targetEdgeAnnotationsImp map[string]interface{}

func (ea targetEdgeAnnotationsImp) HasAnnotation(ann string) bool {
	_, ok := ea[ann]
	return ok
}

func (ea targetEdgeAnnotationsImp) ListAnnotations() []string {
	l := make([]string, 0, len(ea))
	for ann := range ea {
		l = append(l, ann)
	}
	return l
}


// licenseGraphImp implements the LicenseGraph interface.
type licenseGraphImp struct {
	// rootFiles identifies the original set of files to read (immutable)
	rootFiles []string
	// targets identifies the entire set of target node files (guarded by mu)
	targets map[string]*targetNode
	// edges lists the target edges from dependent to dependency (guarded by mu)
	edges []*dependencyEdge
	// index facilitates looking up edges from targets (creation guarded by my)
	index map[string][]*dependencyEdge
	// reverse facilitates looking up edges from dependencies (creation guarded by mu)
	reverse map[string[]*dependencyEdge
	// mu guards against concurrent update
	mu sync.Mutex
}

// indexForward guarantees the `index` map is populated to look up edges by `target`
func (lm *licenseGraphImp) indexForward() {
	mu.Lock()
	defer func() {
		mu.Unlock()
	}()

	if lm.index != nil {
		return
	}

	lm.index = make(map[string][]*dependencyEdge)
	for _, e := range lm.edges {
		if _, ok := lm.index[e.target]; ok {
			lm.index[e.target] = append(lm.index[e.target], e)
		} else {
			lm.index[e.target] = []*dependencyEdge{e}
		}
	}
}

// indexReverse guarantees the `reverse` map is populated to look up edges by `dependency`
func (lm *licenseGraphImp) indexReverse() {
	mu.Lock()
	defer func() {
		mu.Unlock()
	}()

	if lm.reverse != nil {
		return
	}

	lm.reverse = make(map[string][]*dependencyEdge)
	for _, e := range lm.edges {
		if _, ok := lm.reverse[e.dependency]; ok {
			lm.reverse[e.dependency] = append(lm.reverse[e.dependency], e)
		} else {
			lm.reverse[e.dependency] = []*dependencyEdge{e}
		}
	}
}

// AnyPath returns a TargetPath from target to dependency or nil if none exist.
func (lm *licenseGraphImp) AnyPath(target, dependency TargetNode) TargetPath {
	lm.indexForward()

	stack := make([]string)
	stack = append(stack, target.file)

	path := make(TargetPath)

	index := []int{0}
	for {
		edges := lm.index[stack[len(stack)-1]]
		i := index[len(index)-1]
		if i >= len(edges) {
			stack = stack[:len(stack)-1]
			if len(stack) < 1 {
				break
			}
			path = path[:len(path)-1]
			index = index[:len(index)-1]
			index[len(index)-1]++
			continue
		}
		stack = append(stack, edges[i].dependency)
		path = append(path, targetPathImp{lm, edges[i]})
		if dependency.file == stack[len(stack)-1] {
			return path
		}
		index = append(index, 0)
	}
	return nil
}

// AllPaths returns a slice of TargetPath from target to dependency with an entry
// for each distinct path from target to dependency.
func (lm *licenseGraphImp) AllPaths(target, dependency TargetNode) []TargetPath {
	lm.indexForward()

	paths := make([]TargetPath)

	stack := make([]string)
	stack = append(stack, target.file)

	path := make(TargetPath)

	index := []int{0}
	for {
		edges := lm.index[stack[len(stack)-1]]
		i := index[len(index)-1]
		if i >= len(edges) {
			stack = stack[:len(stack)-1]
			if len(stack) < 1 {
				break
			}
			path = path[:len(path)-1]
			index = index[:len(index)-1]
			index[len(index)-1]++
			continue
		}
		stack = append(stack, edges[i].dependency)
		path = append(path, targetPathImp{lm, edges[i]})
		if dependency.file == stack[len(stack)-1] {
			c := append(TargetPath{}, path...)
			paths = append(paths, c)
			stack = stack[:len(stack)-1]
			if len(stack) < 1 {
				break
			}
			path = path[:len(path)-1]
			index = index[:len(index)-1]
			index[len(index)-1]++
			continue
		}
		index = append(index, 0)
	}
	return paths
}


func (lm *licenseGraphImp) WalkRestricted() TargetSet {
	rmap := make(map[string]interface{})
	cmap := make(map[string]interface{})

	var walkContainer, walkNonContainer func(string)

	walkNonContainer = func(f string) {
		rmap[f] = nil
		for _, d := range lm.metadata[f].deps {
			if _, alreadyWalked := rmap[d]; alreadyWalked {
				if _, asContainer := cmap[d]; asContainer {
					delete(cmap, d)
					walkNonContainer(d)
				}
				continue
			}
			walkNonContainer(d)
		}
	}

	walkContainer = func(f string) {
		rmap[f] = nil
		cmap[f] = nil
		for _, d := range lm.metadata[f].deps {
			if !lm.metadata[d].isRestricted {
				continue
			}
			if _, alreadyWalked := rmap[d]; alreadyWalked {
				continue
			}
			if lm.metadata[d].isContainer {
				walkContainer(d)
			} else {
				walkNonContainer(d)
			}
		}
	}

	for _, r := range lm.rootFiles {
		if !lm.metadata[r].isRestricted {
			continue
		}
		if _, alreadyWalked := rmap[r]; alreadyWalked {
			continue
		}
		if lm.metadata[r].isContainer {
			walkContainer(r)
		} else {
			walkNonContainer(r)
		}
	}

	return TargetSet{lm, rmap}
}

func (lm *licenseGraphImp) WalkDepsForCondition(condition string) TargetSet {
	if strings.HasPrefix(condition, "restrict") {
		return lm.WalkRestricted()
	}

	matched := make(map[string]interface{})
	wmap := make(map[string]interface{})

	var walk func(string)

	walk = func(f string) {
		wmap[f] = nil
		for _, c := range lm.metadata[f].licenseConditions {
			if c == condition {
				matched[f] = nil
				break
			}
		}
		for _, d := range lm.metadata[f].deps {
			if _, alreadyWalked := wmap[d]; alreadyWalked {
				continue
			}
			walk(d)
		}
	}

	for _, r := range lm.rootFiles {
		if _, alreadyWalked := wmap[r]; alreadyWalked {
			continue
		}
		walk(r)
	}

	return TargetSet{lm, matched}
}

// newLicenseGraphImp constructs a new instance of licenseGraphImp.
func newLicenseGraphImp() *licenseGraphImp {
	return &licenseGraphImp{
		[]string{},
		make(map[string]*targetNode),
	}
}

// result describes the outcome of reading and parsing a single license metadata file.
type result struct {
	// file identifies the path to the license metadata file
	file string
	// target contains the parsed metadata or nil if an error
	target *targetNode
	// edges contains the parsed dependencies
	edges []dependencyEdge
	// err is nil unless an error occurs
	err error
}

// receiver coordinates the tasks for reading and parsing license metadata files.
type receiver struct {
	// ls accumulates the read metadata
	lm *licenseGraphImp
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

	lm := newLicenseGraphImp()
	for _, f := range files {
		if strings.HasSuffix(f, ".meta_lic") {
			lm.rootFiles = append(lm.rootFiles, f)
		} else {
			lm.rootFiles = append(lm.rootFiles, f+".meta_lic")
		}
	}

	recv := &receiver{lm: lm, task: make(chan bool, ConcurrentReaders), results: make(chan *result, ConcurrentReaders), wg: sync.WaitGroup{}}
	for i := 0; i < ConcurrentReaders; i++ {
		recv.task <- true
	}

	readFiles := func() {
		// identify the metadata files to schedule reading tasks for
		for _, f := range lm.rootFiles {
			lm.targets[f] = nil
		}

		// schedule tasks to read the files
		for _, f := range lm.rootFiles {
			recv.wg.Add(1)
			<-recv.task
			go readFile(recv, f)
		}

		// schedule a task to wait until finished and close the channelm.
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
					lm = nil
					recv.results = nil
					continue
				}

				// record the parsed metadata (guarded by mutex)
				recv.mu.Lock()
				recv.lm.targets[r.file] = r.target
				if len(r.edges) > 0 {
					recv.lm.edges = append(recv.lm.edges, r.edges...)
				}
				recv.mu.Unlock()
			} else {
				// finished -- nil the results channel
				recv.results = nil
			}
		}
	}

	// single task has reference to `ls` at this point
	if lm != nil {
		// reading and parsing succeeded -- calculate effective conditions and license texts
		for f := range lm.metadata {
			lm.metadata[f].isRestricted = isRestricted(lm.metadata[f], lm)
			lm.metadata[f].isProprietary = isProprietary(lm.metadata[f], lm)
			flattenTexts(lm.metadata[f], lm)
		}
	}
	return lm, err
}

// addDependency parses a dependency from `tn` and adds it to `edges`
func addDependency(edges *[]*dependencyEdge, target string, value string) error {
	fields := strings.Split(value, ":")
	dependency := fields[0]
	annotations := make(targetEdgeAnnotationsImp)
	for a := range fields[1:] {
		annotations[a] = nil
	}
	edges = append(edges, &dependencyEdge{target, dependency, annotations})
	return nil
}

// readFile is a task to read and parse a single license metadata file, and to schedule
// additional tasks for reading and parsing dependencies as necessary.
func readFile(recv *receiver, file string) {
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

	tn := &targetNode{file}
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
	recv.results <- &result{file, tn, edges, nil}
	recv.task <- true

	// schedule tasks as necessary to read dependencies
	for _, e := range edges {
		// decide and record whether to schedul task in critical section
		recv.mu.Lock()
		_, alreadyScheduled := recv.lm.targets[e.dependency]
		if !alreadyScheduled {
			recv.lm.targets[e.dependency] = nil
		}
		recv.mu.Unlock()
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

// isRestricted returns true if `mf` or any dependency has 'restricted' condition.
func isRestricted(mf *metadataFile, lm *licenseGraphImp) bool {
	// check whether already has "restricted"
	for _, c := range mf.licenseConditions {
		if strings.HasPrefix(c, "restrict") {
			return true
		}
	}

	// add "restricted" if any dependency has "restricted"
	for _, d := range mf.deps {
		if isRestricted(lm.metadata[d], lm) {
			return true
		}
	}

	return false
}

// isProprietary returns true if `mf` or any dependency has 'proprietary' condition.
func isProprietary(mf *metadataFile, lm *licenseGraphImp) bool {
	// check whether already has "proprietary"
	for _, c := range mf.licenseConditions {
		if c == "proprietary" {
			return true
		}
	}

	// add "proprietary" if any dependency has "proprietary"
	for _, d := range mf.deps {
		if isProprietary(lm.metadata[d], lm) {
			return true
		}
	}

	return false
}

// flattenTexts calculates the transitive closure of license texts for `mf` from its dependencies.
func flattenTexts(mf *metadataFile, lm *licenseGraphImp) {
	// copy license texts
	mf.effectiveLicenseTexts = append(mf.effectiveLicenseTexts, mf.licenseTexts...)

	// indicate license texts already recorded
	recorded := make(map[string]bool)
	for _, t := range mf.effectiveLicenseTexts {
		recorded[t] = true
	}

	// add unrecorded license texts from dependencies
	for _, d := range mf.deps {
		for _, t := range lm.metadata[d].licenseTexts {
			if _, alreadyRecorded := recorded[t]; !alreadyRecorded {
				mf.effectiveLicenseTexts = append(mf.effectiveLicenseTexts, t)
				recorded[t] = true
			}
		}
	}
}
