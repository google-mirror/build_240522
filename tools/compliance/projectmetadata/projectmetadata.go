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

package projectmetadata

import (
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"strings"
	"sync"

	"android/soong/compliance/project_metadata_proto"

	"google.golang.org/protobuf/encoding/prototext"
)

var (
	// ConcurrentReaders is the size of the task pool for limiting resource usage e.g. open files.
	ConcurrentReaders = 5
)

// ProjectMetadata contains the METADATA for a git project.
type ProjectMetadata struct {
	proto project_metadata_proto.Metadata

	// project is the path to the directory containing the METADATA file.
	project string
}

// String returns a string representation of the metadata for error messages.
func (pm *ProjectMetadata) String() string {
	return fmt.Sprintf("project: %q\n%s", pm.project, pm.proto.String())
}

// VersionedName returns the name of the project including the version if any.
func (pm *ProjectMetadata) VersionedName() string {
	name := pm.proto.GetName()
	if name != "" {
		tp := pm.proto.GetThirdParty()
		if tp != nil {
			version := tp.GetVersion()
			if version != "" {
				if version[0] == 'v' || version[0] == 'V' {
					return name + "_" + version
				} else {
					return name + "_v_" + version
				}
			}
		}
		return name
	}
	return pm.proto.GetDescription()
}

// Index reads and caches ProjectMetadata (thread safe)
type Index struct {
	// projectWaitGroup maps project name to a wait group if read has already started.
	// guarded by mu
	projectWaitGroup map[string]*sync.WaitGroup

	// projectMetadata maps project name to the metadata content or nil after read.
	// guarded by mu
	projectMetadata map[string]*ProjectMetadata

	// projectError maps project name to the error encountered reading the metadata file if any.
	// guarded by mu
	projectError map[string]error

	// task provides a fixed-size task pool to limit concurrent open files etc.
	// guarded by mu
	task chan bool

	// guards concurrent access of projectWaitGroup, projectMetadata, projectError and task channel above.
	mu sync.Mutex

	// rootFS locates the root of the file system from which to read the files.
	rootFS fs.FS
}

// NewIndex constructs a project metadata `Index` for the given file system.
func NewIndex(rootFS fs.FS) *Index {
	ix := &Index{
		projectWaitGroup: make(map[string]*sync.WaitGroup),
		projectMetadata:  make(map[string]*ProjectMetadata),
		projectError:     make(map[string]error),
		rootFS:           rootFS,
	}
	return ix
}

// MetadataForProjects returns 0..n ProjectMetadata for n `projects`, or an error.
// Each project that has a METADATA.android or a METADATA file in the root of the project will have
// a corresponding ProjectMetadata in the result. Projects with neither file get skipped. A nil
// result with no error indicates none of the given `projects` has a METADATA file.
// (thread safe -- can be called concurrently from multiple goroutines)
func (ix *Index) MetadataForProjects(projects ...string) ([]*ProjectMetadata, error) {
	if ConcurrentReaders < 1 {
		return nil, fmt.Errorf("need at least one task in project metadata pool")
	}
	if len(projects) == 0 {
		return nil, nil
	}
	// Identify the projects that have never been read
	projectsToRead := make([]string, 0, len(projects))
	ix.mu.Lock()
	if ix.task == nil {
		ix.task = make(chan bool, ConcurrentReaders)
		for i := 0; i < ConcurrentReaders; i++ {
			ix.task <- true
		}
	}
	for _, p := range projects {
		if _, ok := ix.projectWaitGroup[p]; !ok {
			projectsToRead = append(projectsToRead, p)
			wg := sync.WaitGroup{}
			wg.Add(1)
			ix.projectWaitGroup[p] = &wg
		}
	}
	ix.mu.Unlock()
	// read parses the METADATA file at `path` for `project` and records the result.
	read := func(project, path string) {
		pm, err := ix.readMetadataFile(project, path)
		ix.mu.Lock()
		if err == nil {
			ix.projectMetadata[project] = pm
		} else {
			ix.projectError[project] = err
		}
		ix.projectWaitGroup[project].Done()
		ix.mu.Unlock()
	}
	findMeta := func(project string) {
		<-ix.task
		defer func() { ix.task <- true }()

		// Support METADATA.android for projects that already have a different sort of METADATA file.
		path := filepath.Join(project, "METADATA.android")
		fi, err := fs.Stat(ix.rootFS, path)
		if err == nil {
			if fi.Mode().IsRegular() {
				read(project, path)
				return
			}
		}
		// No METADATA.android try METADATA file.
		path = filepath.Join(project, "METADATA")
		fi, err = fs.Stat(ix.rootFS, path)
		if err == nil {
			if fi.Mode().IsRegular() {
				read(project, path)
				return
			}
		}
		// no METADATA file exists -- use nil
		ix.mu.Lock()
		ix.projectMetadata[project] = nil
		ix.projectWaitGroup[project].Done()
		ix.mu.Unlock()
	}
	// Look for the METADATA files to read, and record any missing.
	for _, p := range projectsToRead {
		go findMeta(p)
	}
	// Wait until all of the projects have been read.
	for _, p := range projects {
		ix.mu.Lock()
		wg := ix.projectWaitGroup[p]
		ix.mu.Unlock()
		wg.Wait()
	}
	// Combine any errors into a single error.
	var msg strings.Builder
	result := make([]*ProjectMetadata, 0, len(projects))
	for _, p := range projects {
		ix.mu.Lock()
		err, hasError := ix.projectError[p]
		pm, ok := ix.projectMetadata[p]
		ix.mu.Unlock()
		if hasError && err != nil {
			fmt.Fprintf(&msg, "  %v\n", err)
		}
		if ok && pm != nil {
			result = append(result, pm)
		}
	}
	if msg.Len() > 0 {
		return nil, fmt.Errorf("error reading project(s):\n%s", msg.String())
	}
	if len(result) == 0 {
		return nil, nil
	}
	return result, nil
}

// readMetadataFile tries to read and parse a METADATA file at `path` for `project`.
func (ix *Index) readMetadataFile(project, path string) (*ProjectMetadata, error) {
	f, err := ix.rootFS.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error opening project %q metadata %q: %w", project, path, err)
	}

	// read the file
	data, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("error reading project %q metadata %q: %w", project, path, err)
	}
	f.Close()

	uo := prototext.UnmarshalOptions{DiscardUnknown: true}
	pm := &ProjectMetadata{project: project}
	err = uo.Unmarshal(data, &pm.proto)
	if err != nil {
		return nil, fmt.Errorf("error in project %q metadata %q: %w", project, path, err)
	}

	return pm, nil
}
