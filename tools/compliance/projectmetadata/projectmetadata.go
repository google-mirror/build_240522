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

// ProjectMetadata contains the METADATA for a git project.
type ProjectMetadata struct {
	proto project_metadata_proto.Metadata

	// project is the path to the directory containing the METADATA file.
	project string
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
	projectError map[string]error

	// guards concurrent access of projectWaitGroup and projectMetadata above.
	mu sync.Mutex

	// rootFS locates the root of the file system from which to read the files.
	rootFS fs.FS
}

// NewIndex constructs a project metadata `Index` for the given file system.
func NewIndex(rootFS fs.FS) *Index {
	return &Index{make(map[string]*sync.WaitGroup), make(map[string]*ProjectMetadata), make(map[string]error), sync.Mutex{}, rootFS}
}

// MetadataForProjects returns 0..n ProjectMetadata for n `projects`, or an error.
// Each project that has a METADATA.android or a METADATA file in the root of the project will have
// a corresponding ProjectMetadata in the result. Projects with neither file get skipped. A nil
// result with no error indicates none of the given `projects` has a METADATA file.
func (ix *Index) MetadataForProjects(projects ...string) ([]*ProjectMetadata, error) {
	if len(projects) == 0 {
		return nil, nil
	}
	// Identify the projects that have never been read
	projectsToRead := make([]string, 0, len(projects))
	ix.mu.Lock()
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
	// Look for the METADATA files to read, and record any missing.
	for _, p := range projectsToRead {
		// Support METADATA.android for projects that already have a different sort of METADATA file.
		path := filepath.Join(p, "METADATA.android")
		fi, err := fs.Stat(ix.rootFS, path)
		if err == nil {
			if fi.Mode().IsRegular() {
				go read(p, path)
				continue
			}
		}
		// No METADATA.android try METADATA file.
		path = filepath.Join(p, "METADATA")
		fi, err = fs.Stat(ix.rootFS, path)
		if err == nil {
			if fi.Mode().IsRegular() {
				go read(p, path)
				continue
			}
		}
		// no METADATA file exists -- use nil
		ix.mu.Lock()
		ix.projectMetadata[p] = nil
		ix.projectWaitGroup[p].Done()
		ix.mu.Unlock()
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
			fmt.Fprintf(&msg, "  %w\n", err)
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
