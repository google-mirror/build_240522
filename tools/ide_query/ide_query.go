package main

import (
	"container/list"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"slices"
	"strings"

	"google.golang.org/protobuf/proto"
	pb "ide_query/ide_query_proto"
)

// Env contains information about the current environment.
type Env struct {
	LunchTarget LunchTarget
	RepoDir     string
	OutDir      string
}

// LunchTarget is a parsed Android lunch target.
// Input format: <product_name>-<release_type>-<build_variant>
type LunchTarget struct {
	Product string
	Release string
	Variant string
}

var _ flag.Value = (*LunchTarget)(nil)

// Get implements flag.Value.
func (l *LunchTarget) Get() any {
	return l
}

// Set implements flag.Value.
func (l *LunchTarget) Set(s string) error {
	parts := strings.Split(s, "-")
	if len(parts) != 3 {
		return fmt.Errorf("invalid lunch target: %q, must have form <product_name>-<release_type>-<build_variant>", s)
	}
	*l = LunchTarget{
		Product: parts[0],
		Release: parts[1],
		Variant: parts[2],
	}
	return nil
}

// String implements flag.Value.
func (l *LunchTarget) String() string {
	return fmt.Sprintf("%s-%s-%s", l.Product, l.Release, l.Variant)
}

func main() {
	var env Env
	env.OutDir = os.Getenv("OUT_DIR")
	env.RepoDir = os.Getenv("TOP")
	flag.Var(&env.LunchTarget, "lunch_target", "The lunch target to query")
	flag.Parse()
	if flag.NArg() == 0 {
		fmt.Println("No files provided.")
		os.Exit(1)
		return
	}

	files := flag.Args()

	ctx := context.Background()
	javaDepsPath := path.Join(env.RepoDir, env.OutDir, "soong/module_bp_java_deps.json")
	// TODO(michaelmerg): Figure out if module_bp_java_deps.json is outdated.
	if _, err := os.Stat(javaDepsPath); os.IsNotExist(err) {
		runMake(ctx, env, "nothing")
	}

	javaModules, err := loadJavaModules(javaDepsPath)
	if err != nil {
		log.Fatalf("Failed to load java modules: %v", err)
	}

	modules := make(map[string]*javaModule) // file path -> module
	for _, f := range files {
		for _, m := range javaModules {
			if !slices.Contains(m.Srcs, f) {
				continue
			}
			if prev, ok := modules[f]; ok && !strings.HasSuffix(prev.Name, ".impl") {
				log.Printf("File %q found in module %q but is already part of module %q", f, m.Name, prev.Name)
				continue
			}
			modules[f] = m
		}
	}

	var toMake []string
	for _, m := range modules {
		toMake = append(toMake, m.Name)
	}
	fmt.Printf("Running make for modules: %v\n", strings.Join(toMake, ", "))
	if err := runMake(ctx, env, toMake...); err != nil {
		log.Fatalf("Failed to run make: %v", err)
	}

	var sources []*pb.SourceFile
	for _, f := range files {
		file := &pb.SourceFile{
			Path:       f,
			WorkingDir: env.RepoDir,
		}
		sources = append(sources, file)

		m := modules[f]
		if m == nil {
			file.Status = &pb.Status{
				Code:    pb.Status_FAILURE,
				Message: proto.String("File not found in any module."),
			}
			continue
		}

		deps := transitiveDeps(m, javaModules)
		var genFiles []*pb.GeneratedFile
		outPrefix := env.OutDir + "/"
		for _, d := range deps {
			if relPath, ok := strings.CutPrefix(d, outPrefix); ok {
				// Note: Contents will be filled below.
				genFiles = append(genFiles, &pb.GeneratedFile{Path: relPath})
			}
		}

		file.Generated = genFiles
		file.Deps = deps
		file.Status = &pb.Status{Code: pb.Status_OK}
	}

	for _, s := range sources {
		for _, g := range s.GetGenerated() {
			contents, err := os.ReadFile(path.Join(env.OutDir, g.GetPath()))
			if err != nil {
				fmt.Printf("Failed to read generated file %q: %v. File contents will be missing.\n", g.GetPath(), err)
				continue
			}
			g.Contents = contents
		}
	}

	res := &pb.IdeAnalysis{
		BuildArtifactRoot: env.OutDir,
		Sources:           sources,
		Status:            &pb.Status{Code: pb.Status_OK},
	}
	data, err := proto.Marshal(res)
	if err != nil {
		log.Fatalf("Failed to marshal result proto: %v", err)
	}

	err = os.WriteFile(path.Join(env.OutDir, "ide_query.pb"), data, 0644)
	if err != nil {
		log.Fatalf("Failed to write result proto: %v", err)
	}

	for _, s := range sources {
		fmt.Printf("%s: %v (Deps: %d, Generated: %d)\n", s.GetPath(), s.GetStatus(), len(s.GetDeps()), len(s.GetGenerated()))
	}
}

// runMake runs Soong build for the given modules.
func runMake(ctx context.Context, env Env, modules ...string) error {
	args := []string{
		"--make-mode",
		"ANDROID_BUILD_ENVIRONMENT_CONFIG=googler-cog",
		"TARGET_PRODUCT=" + env.LunchTarget.Product,
		"TARGET_RELEASE=" + env.LunchTarget.Release,
		"TARGET_BUILD_VARIANT=" + env.LunchTarget.Variant,
	}
	args = append(args, modules...)
	cmd := exec.CommandContext(ctx, "build/soong/soong_ui.bash", args...)
	cmd.Dir = env.RepoDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

type javaModule struct {
	Name    string
	Path    []string `json:"path,omitempty"`
	Deps    []string `json:"dependencies,omitempty"`
	Srcs    []string `json:"srcs,omitempty"`
	Jars    []string `json:"jars,omitempty"`
	SrcJars []string `json:"srcjars,omitempty"`
}

func loadJavaModules(path string) (map[string]*javaModule, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var ret map[string]*javaModule // module name -> module
	if err = json.Unmarshal(data, &ret); err != nil {
		return nil, err
	}

	for name, module := range ret {
		if strings.HasSuffix(name, "-jarjar") {
			delete(ret, name)
			continue
		}

		module.Name = name
	}
	return ret, nil
}

func transitiveDeps(m *javaModule, modules map[string]*javaModule) []string {
	var ret []string
	q := list.New()
	q.PushBack(m.Name)
	seen := make(map[string]bool) // module names -> true
	for q.Len() > 0 {
		name := q.Remove(q.Front()).(string)
		mod := modules[name]
		if mod == nil {
			continue
		}

		ret = append(ret, mod.Srcs...)
		ret = append(ret, mod.SrcJars...)
		ret = append(ret, mod.Jars...)
		for _, d := range mod.Deps {
			if seen[d] {
				continue
			}
			seen[d] = true
			q.PushBack(d)
		}
	}
	slices.Sort(ret)
	ret = slices.Compact(ret)
	return ret
}
