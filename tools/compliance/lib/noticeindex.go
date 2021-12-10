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
	"crypto/md5"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	nameRegexp = regexp.MustCompile(`^\s*name\s*:\s*"(.*)"\s*$`)
	descRegexp = regexp.MustCompile(`^\s*description\s*:\s*"(.*)"\s*$`)
	versionRegexp = regexp.MustCompile(`^\s*version\s*:\s*"(.*)"\s*$`)
)

type NoticeResolution struct {
	attachesTo *TargetNode
	path TargetEdgePath
	noticeFor *TargetNode
}

func (nr *NoticeResolution) AttachesTo() *TargetNode {
	return nr.attachesTo
}

func (nr *NoticeResolution) NoticeFor() *TargetNode {
	return nr.noticeFor
}

func (nr *NoticeResolution) Path() *TargetEdgePath {
	return (&nr.path).Copy()
}

func (nr *NoticeResolution) InstallPaths() []string {
	if len(nr.path) == 0 {
		return nr.attachesTo.TargetFiles()
	}

	var getInstalls func(path TargetEdgePath) []string

	getInstalls = func(path TargetEdgePath) []string {
		var deps []string
		if len(path) > 1 {
			deps = getInstalls(path[1:])
		} else {
			deps = path[0].Dependency().TargetFiles()
		}

		size := 0
		prefixes := path[0].Target().TargetFiles()
		installMap := path[0].Target().InstallMap()
		for _, dep := range deps {
			for _, im := range installMap {
				if strings.HasPrefix(dep, im.FromPath) {
					size += len(prefixes)
					break
				}
			}
		}

		installs := make([]string, 0, size)
		for _, dep := range deps {
			for _, im := range installMap {
				if strings.HasPrefix(dep, im.FromPath) {
					for _, prefix := range prefixes {
						installs = append(installs, prefix + im.ContainerPath + dep[len(im.FromPath):])
					}
					break
				}
			}
		}
		return installs
	}
	return getInstalls(nr.path)
}

type NoticeResolutionList []*NoticeResolution

type hash struct {
	key string
}

type targetPair struct {
	target, dependency *TargetNode
}

type NoticeIndex struct {
	lg *LicenseGraph
	rs *ResolutionSet
	shipped *TargetNodeSet
	// rootFS locates the root of the file system from which to read the files.
	rootFS fs.FS
	hash map[string]hash
	text map[hash][]byte
	reverse map[hash]NoticeResolutionList
	rl NoticeResolutionList
	libName map[*TargetNode]string
	projectName map[string]string
}

func IndexLicenseTexts(rootFS fs.FS, lg *LicenseGraph, rs *ResolutionSet) (*NoticeIndex, error) {
	if rs == nil {
		rs = ResolveNotices(lg)
	}
	ni := &NoticeIndex{
		lg, rs, ShippedNodes(lg), rootFS,
		make(map[string]hash),
		make(map[hash][]byte),
		make(map[hash]NoticeResolutionList),
		make(NoticeResolutionList, 0, 0),
		make(map[*TargetNode]string),
		make(map[string]string),
	}

	size := 0
	WalkTopDown(lg, func(lg *LicenseGraph, tn *TargetNode, path TargetEdgePath) bool {
		if !ni.shipped.Contains(tn) {
			return false
		}
		ni.libName[tn] = ni.getLibName(tn)
		size++
		return tn.IsContainer()
	})

	ni.rl = make(NoticeResolutionList, 0, size)
	var err error

	WalkTopDown(lg, func(lg *LicenseGraph, tn *TargetNode, path TargetEdgePath) bool {
		if err != nil {
			return false
		}
		if !ni.shipped.Contains(tn) {
			return false
		}
		attachesTo := tn
		if len(path) > 0 {
			attachesTo = path[0].Target()
		}
		nr := &NoticeResolution{attachesTo, path, tn}
		ni.rl = append(ni.rl, nr)
		for _, text := range tn.LicenseTexts() {
			if _, ok := ni.hash[text]; !ok {
				err = ni.addText(text)
				if err != nil {
					return false
				}
			}
			hash := ni.hash[text]
			if _, ok := ni.reverse[hash]; !ok {
				ni.reverse[hash] = make(NoticeResolutionList, 0, 1)
			}
			ni.reverse[hash] = append(ni.reverse[hash], nr)
		}
		return false
	})

	if err != nil {
		return nil, err
	}

	return ni, nil
}

func (ni *NoticeIndex) LibraryName(nr *NoticeResolution) string {
	if libName, ok := ni.libName[nr.noticeFor]; ok {
		return libName
	}
	return ni.getLibName(nr.noticeFor)
}

func (ni *NoticeIndex) Resolutions() NoticeResolutionList {
	rl := make(NoticeResolutionList, 0, len(ni.rl))
	rl = append(rl, ni.rl...)
	return rl
}

func (ni *NoticeIndex) getLibName(tn *TargetNode) string {
	ln := ni.checkMetadata(tn)
	if len(ln) > 0 {
		return ln
	}
	pn := tn.PackageName()
	if len(pn) > 0 {
		return pn
	}
	n := tn.name[:len(tn.name)-9]
	li := strings.LastIndex(n, "/")
	if 0 < li {
		n = n[li+1:]
	}
	return n
}

func (ni *NoticeIndex) checkMetadata(tn *TargetNode) string {
	for _, p := range tn.Projects() {
		if name, ok := ni.projectName[p]; ok {
			return name
		}
		f, err := ni.rootFS.Open(filepath.Join(p, "METADATA"))
		if err != nil {
			continue
		}
		name := ""
		description := ""
		version := ""
		s := bufio.NewScanner(f)
		for s.Scan() {
			line := s.Text()
			m := nameRegexp.FindStringSubmatch(line)
			if m != nil {
				if 1 < len(m) && m[1] != "" {
					name = m[1]
				}
				if version != "" {
					break
				}
				continue
			}
			m = versionRegexp.FindStringSubmatch(line)
			if m != nil {
				if 1 < len(m) && m[1] != "" {
					version = m[1]
				}
				if name != "" {
					break
				}
				continue
			}
			m = descRegexp.FindStringSubmatch(line)
			if m != nil {
				if 1 < len(m) && m[1] != "" {
					description = m[1]
				}
			}
		}
		_ = s.Err()
		_ = f.Close()
		if name != "" {
			if version != "" {
				ni.projectName[p] = name + " " + version
			} else {
				ni.projectName[p] = name
			}
			return ni.projectName[p]
		}
		if description != "" {
			ni.projectName[p] = description
			return ni.projectName[p]
		}
	}
	return ""
}

func (ni *NoticeIndex) addText(file string) error {
	f, err := ni.rootFS.Open(filepath.Clean(file))
	if err != nil {
		return fmt.Errorf("error opening license text file %q: %w", file, err)
	}

	// read the file
	text, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("error reading license text file %q: %w", file, err)
	}

	hash := hash{fmt.Sprintf("%x", md5.New().Sum(text))}
	ni.hash[file] = hash
	if _, alreadyPresent := ni.text[hash]; !alreadyPresent {
		ni.text[hash] = text
	}

	return nil
}
