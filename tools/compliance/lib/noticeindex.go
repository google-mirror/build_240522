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
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const (
	noProjectName = "\u2205"
)

var (
	nameRegexp = regexp.MustCompile(`^\s*name\s*:\s*"(.*)"\s*$`)
	descRegexp = regexp.MustCompile(`^\s*description\s*:\s*"(.*)"\s*$`)
	versionRegexp = regexp.MustCompile(`^\s*version\s*:\s*"(.*)"\s*$`)
	licensesPathRegexp = regexp.MustCompile(`licen[cs]es?/`)
)

/*
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

type NoticeResolutionList []*NoticeResolution
*/

func getInstallPaths(attachesTo *TargetNode, path TargetEdgePath) []string {
	if len(path) == 0 {
		return attachesTo.TargetFiles()
	}

	var getInstalls func(path TargetEdgePath) []string

	getInstalls = func(path TargetEdgePath) []string {
		// deps contains the output targets from the dependencies in the path
		var deps []string
		if len(path) > 1 {
			// recursively get the targets from the sub-path skipping 1 path segment
			deps = getInstalls(path[1:])
		} else {
			// stop recursion at 1 path segment
			deps = path[0].Dependency().TargetFiles()
		}

		size := 0
		prefixes := path[0].Target().TargetFiles()
		installMap := path[0].Target().InstallMap()
		sources := path[0].Target().Sources()
		for _, dep := range deps {
			found := false
			for _, source := range sources {
				if strings.HasPrefix(dep, source) {
					found = true
					break
				}
			}
			if !found {
				continue
			}
			for _, im := range installMap {
				if strings.HasPrefix(dep, im.FromPath) {
					size += len(prefixes)
					break
				}
			}
		}

		installs := make([]string, 0, size)
		for _, dep := range deps {
			found := false
			for _, source := range sources {
				if strings.HasPrefix(dep, source) {
					found = true
					break
				}
			}
			if !found {
				continue
			}
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
	return getInstalls(path)
}

type hash struct {
	key string
}

func (h hash) String() string {
	return h.key
}

type targetPair struct {
	target, dependency *TargetNode
}

type NoticeIndex struct {
	lg *LicenseGraph
	rs ResolutionSet
	shipped *TargetNodeSet
	// rootFS locates the root of the file system from which to read the files.
	rootFS fs.FS
	// hash maps license text filenames to content hashes
	hash map[string]hash
	// text maps content hashes to content
	text map[hash][]byte
	// hashLibInstall maps hashes to libraries to install paths.
	hashLibInstall map[hash]map[string]map[string]struct{}
	// installLibHash maps install paths to libraries to hashes.
	installLibHash map[string]map[string]map[hash]struct{}
	// libHash maps libraries to hashes.
	libHash map[string]map[hash]struct{}
	// projectName maps project directory names to project name text.
	projectName map[string]string
}

func (ni *NoticeIndex) Hashes() chan hash {
	c := make(chan hash)
	libs := make([]string, 0, len(ni.libHash))
	for libName := range ni.libHash {
		libs = append(libs, libName)
	}
	go func() {
		sort.Strings(libs)
		hashes := make(map[hash]struct{})
		for _, libName := range libs {
			for h := range ni.libHash[libName] {
				if _, ok := hashes[h]; ok {
					continue
				}
				hashes[h] = struct{}{}
				c <- h
			}
		}
		close(c)
	}()
	return c
}

func (ni *NoticeIndex) HashLibs(h hash) []string {
	libs := make([]string, 0, len(ni.hashLibInstall[h]))
	for libName := range ni.hashLibInstall[h] {
		libs = append(libs, libName)
	}
	sort.Strings(libs)
	return libs
}

func (ni *NoticeIndex) HashLibInstalls(h hash, libName string) []string {
	installs := make([]string, 0, len(ni.hashLibInstall[h][libName]))
	for installPath := range ni.hashLibInstall[h][libName] {
		installs = append(installs, installPath)
	}
	sort.Strings(installs)
	return installs
}

func (ni *NoticeIndex) HashText(h hash) []byte {
	return ni.text[h]
}

func IndexLicenseTexts(rootFS fs.FS, lg *LicenseGraph, rs ResolutionSet) (*NoticeIndex, error) {
	if rs == nil {
		rs = ResolveNotices(lg)
	}
	ni := &NoticeIndex{
		lg, rs, ShippedNodes(lg), rootFS,
		make(map[string]hash),
		make(map[hash][]byte),
		make(map[hash]map[string]map[string]struct{}),
		make(map[string]map[string]map[hash]struct{}),
		make(map[string]map[hash]struct{}),
		make(map[string]string),
	}

	fmt.Fprintf(os.Stderr, "%d shipped\n", len(ni.shipped.nodes))

	index := func(tn *TargetNode) (map[hash]struct{}, error) {
		hashes := make(map[hash]struct{})
		for _, text := range tn.LicenseTexts() {
			if _, ok := ni.hash[text]; !ok {
				err := ni.addText(text)
				if err != nil {
					return nil, err
				}
			}
			hash := ni.hash[text]
			if _, ok := hashes[hash]; !ok {
				hashes[hash] = struct{}{}
			}
		}
		return hashes, nil
	}

	var err error

	progress := 0
	WalkTopDown(NoEdgeContext{}, lg, func(lg *LicenseGraph, tn *TargetNode, path TargetEdgePath) bool {
		if err != nil {
			return false
		}
		if !ni.shipped.Contains(tn) {
			return false
		}
		installPaths := getInstallPaths(tn, path)
		var hashes map[hash]struct{}
		hashes, err = index(tn)
		if err != nil {
			return false
		}
		libNames := make(map[string]struct{})
		if tn.IsContainer() {
			libNames[ni.getLibName(tn)] = struct{}{}
			progress++
			fmt.Fprintf(os.Stderr, "\r%d container       \r", progress)
		}
		for _, r := range rs.Resolutions(tn) {
			libNames[ni.getLibName(r.actsOn)] = struct{}{}
			progress++
			fmt.Fprintf(os.Stderr, "\r%d file            \r", progress)
		}
		for libName := range libNames {
			ni.libHash[libName] = make(map[hash]struct{})
			for h := range hashes {
				if _, ok := ni.hashLibInstall[h]; !ok {
					ni.hashLibInstall[h] = make(map[string]map[string]struct{})
				}
				ni.libHash[libName][h] = struct{}{}
				for _, installPath := range installPaths {
					if _, ok := ni.installLibHash[installPath]; !ok {
						ni.installLibHash[installPath] = make(map[string]map[hash]struct{})
						ni.installLibHash[installPath][libName] = make(map[hash]struct{})
						ni.installLibHash[installPath][libName][h] = struct{}{}
					} else if _, ok = ni.installLibHash[installPath][libName]; !ok {
						ni.installLibHash[installPath][libName] = make(map[hash]struct{})
						ni.installLibHash[installPath][libName][h] = struct{}{}
					} else if _, ok = ni.installLibHash[installPath][libName][h]; !ok {
						ni.installLibHash[installPath][libName][h] = struct{}{}
					}
					if _, ok := ni.hashLibInstall[h]; !ok {
						ni.hashLibInstall[h] = make(map[string]map[string]struct{})
						ni.hashLibInstall[h][libName] = make(map[string]struct{})
						ni.hashLibInstall[h][libName][installPath] = struct{}{}
					} else if _, ok = ni.hashLibInstall[h][libName]; !ok {
						ni.hashLibInstall[h][libName] = make(map[string]struct{})
						ni.hashLibInstall[h][libName][installPath] = struct{}{}
					} else if _, ok = ni.hashLibInstall[h][libName][installPath]; !ok {
						ni.hashLibInstall[h][libName][installPath] = struct{}{}
					}
				}
			}
		}
		return tn.IsContainer()
	})

	if err != nil {
		return nil, err
	}

	return ni, nil
}

/*
func (ni *NoticeIndex) Resolutions() NoticeResolutionList {
	rl := make(NoticeResolutionList, 0, len(ni.rl))
	rl = append(rl, ni.rl...)
	return rl
}
*/

func (ni *NoticeIndex) Hash(fname string) hash {
	return ni.hash[fname]
}

// getLibName returns the name of the library associated with `noticeFor`.
func (ni *NoticeIndex) getLibName(noticeFor *TargetNode) string {
	// use name from METADATA if available
	ln := ni.checkMetadata(noticeFor)
	if len(ln) > 0 {
		return ln
	}
	// use package_name: from license{} module if available
	pn := noticeFor.PackageName()
	if len(pn) > 0 {
		return pn
	}
	for _, p := range noticeFor.Projects() {
		if strings.HasPrefix(p, "prebuilts/") {
			for _, licenseText := range noticeFor.LicenseTexts() {
				if !strings.HasPrefix(licenseText, "prebuilts/") {
					continue
				}
				for r, prefix := range SafePrebuiltPrefixes {
					match := r.FindString(licenseText)
					if len(match) == 0 {
						continue
					}
					strip := SafePathPrefixes[prefix]
					if strip {
						// strip entire prefix
						match = licenseText[len(match):]
					} else {
						// strip from prebuilts/ until safe prefix
						match = licenseText[len(match)-len(prefix):]
					}
					// remove LICENSE or NOTICE or other filename
					li := strings.LastIndex(match, "/")
					if 0 < li {
						match = match[:li]
					}
					// remove *licenses/ path segment and subdirectory if in path
					if offsets := licensesPathRegexp.FindAllStringIndex(match, -1); offsets != nil && 0 < offsets[len(offsets)-1][0] {
						match = match[:offsets[len(offsets)-1][0]]
						li = strings.LastIndex(match, "/")
						if 0 < li {
							match = match[:li]
						}
					}
					return match
				}
				break
			}
		}
		for prefix, strip := range SafePathPrefixes {
			if strings.HasPrefix(p, prefix) {
				if strip {
					return p[len(prefix):]
				} else {
					return p
				}
			}
		}
	}
	// strip off [./]meta_lic from license metadata path and extract base name
	n := noticeFor.name[:len(noticeFor.name)-9]
	li := strings.LastIndex(n, "/")
	if 0 < li {
		n = n[li+1:]
	}
	return n
}

// checkMetadata tries to look up a library name from a METADATA file associated with `tn`.
func (ni *NoticeIndex) checkMetadata(noticeFor *TargetNode) string {
	for _, p := range noticeFor.Projects() {
		if name, ok := ni.projectName[p]; ok {
			if name == noProjectName {
				continue
			}
			return name
		}
		f, err := ni.rootFS.Open(filepath.Join(p, "METADATA"))
		if err != nil {
			ni.projectName[p] = noProjectName
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
				if version[0] == 'v' || version[0] == 'V' {
					ni.projectName[p] = name + "_" + version
				} else {
					ni.projectName[p] = name + "_v_" + version
				}
			} else {
				ni.projectName[p] = name
			}
			return ni.projectName[p]
		}
		if description != "" {
			ni.projectName[p] = description
			return ni.projectName[p]
		}
		ni.projectName[p] = noProjectName
	}
	return ""
}

// addText reads and indexes the content of a license text file.
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

	hash := hash{fmt.Sprintf("%x", md5.Sum(text))}
	ni.hash[file] = hash
	if _, alreadyPresent := ni.text[hash]; !alreadyPresent {
		ni.text[hash] = text
	}

	return nil
}
