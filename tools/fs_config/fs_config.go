// Copyright (C) 2019 The Android Open Source Project
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

package fs_config

import (
	"android/soong/android"
	"fmt"
)

var pctx = android.NewPackageContext("android/soong/fs_config")

func init() {
	android.RegisterModuleType("target_fs_config_gen_filegroup", targetFSConfigGenFactory)
	android.RegisterModuleType("fs_config", fsConfigFactory)
}

// target_fs_config_gen_filegroup is used to expose the files pointed to by TARGET_FS_CONFIG_GEN to
// genrules in Soong. If TARGET_FS_CONFIG_GEN is empty, it will export an empty file instead.
func targetFSConfigGenFactory() android.Module {
	module := &targetFSConfigGen{}
	android.InitAndroidModule(module)
	return module
}

var _ android.SourceFileProducer = (*targetFSConfigGen)(nil)

type targetFSConfigGen struct {
	android.ModuleBase
	paths android.Paths
}

func (targetFSConfigGen) DepsMutator(ctx android.BottomUpMutatorContext) {}

func (t *targetFSConfigGen) GenerateAndroidBuildActions(ctx android.ModuleContext) {
	if ret := ctx.DeviceConfig().TargetFSConfigGen(); len(ret) != 0 {
		t.paths = android.PathsForSource(ctx, ret)
	} else {
		path := android.PathForModuleGen(ctx, "empty")
		t.paths = android.Paths{path}

		rule := android.NewRuleBuilder(pctx, ctx)
		rule.Command().Text("rm -rf").Output(path)
		rule.Command().Text("touch").Output(path)
		rule.Build("fs_config_empty", "create empty file")
	}
}

func (t *targetFSConfigGen) Srcs() android.Paths {
	return t.paths
}

// fs_config is used to create fs_config_dirs or fs_config_files files for the partition where this
// module is installed to.
func fsConfigFactory() android.Module {
	module := &fsConfig{}
	module.AddProperties(&module.properties)
	android.InitAndroidArchModule(module, android.DeviceSupported, android.MultilibCommon)
	return module
}

type fsConfig struct {
	android.ModuleBase
	properties fsConfigProperties

	installDir android.InstallPath
	outputFile android.OutputPath
}

type fsConfigProperties struct {
	// dirs or files
	Type string
}

type configType string

const (
	configTypeDirs    configType = "dirs"
	configTypeFiles   configType = "files"
	configTypeInvalid configType = "invalid"
)

func (c configType) fileName() string {
	return "fs_config_" + string(c)
}

func (f *fsConfig) configType(ctx android.ModuleContext) configType {
	configType := f.properties.Type
	switch configType {
	case "dirs":
		return configTypeDirs
	case "files":
		return configTypeFiles
	default:
		ctx.PropertyErrorf("type", "should be either dirs or files, but got %s", configType)
		return configTypeInvalid
	}
}

func (f *fsConfig) DepsMutator(ctx android.BottomUpMutatorContext) {
}

func (f *fsConfig) GenerateAndroidBuildActions(ctx android.ModuleContext) {
	fsConfigHeader := android.PathForSource(ctx, "system/core/libcutils/include/private/android_filesystem_config.h")
	linuxCapabilityHeader := android.PathForSource(ctx, "bionic/libc/kernel/uapi/linux/capability.h")
	configType := f.configType(ctx)
	partition := f.PartitionTag(ctx.DeviceConfig())

	f.installDir = android.PathForModuleInstall(ctx, "etc")
	f.outputFile = android.PathForModuleOut(ctx, configType.fileName()).OutputPath

	builder := android.NewRuleBuilder(pctx, ctx)
	cmd := builder.Command()
	cmd.BuiltTool("fs_config_generator").
		Text("fsconfig").
		FlagWithInput("--aid-header ", fsConfigHeader).
		FlagWithInput("--capability-header ", linuxCapabilityHeader).
		FlagWithArg("--partition ", partition)
	if partition == "system" {
		cmd.FlagWithArg("--all-partitions ", "vendor,oem,odm,vendor_dlkm,odm_dlkm,system_dlkm")
	}

	if configType == configTypeFiles {
		cmd.Flag("--files")
	} else {
		cmd.Flag("--dirs")
	}

	cmd.FlagWithOutput("--out_file ", f.outputFile)

	if ret := ctx.DeviceConfig().TargetFSConfigGen(); len(ret) != 0 {
		cmd.Inputs(android.PathsForSource(ctx, ret))
	} else {
		cmd.Text("/dev/null")
	}

	builder.Build("fs_config_"+string(configType), fmt.Sprintf("Building %s", f.BaseModuleName()))
}

var _ android.AndroidMkEntriesProvider = (*fsConfig)(nil)

func (f *fsConfig) AndroidMkEntries() []android.AndroidMkEntries {
	return []android.AndroidMkEntries{
		android.AndroidMkEntries{
			Class:      "ETC",
			OutputFile: android.OptionalPathForPath(f.outputFile),
			ExtraEntries: []android.AndroidMkExtraEntriesFunc{
				func(ctx android.AndroidMkExtraEntriesContext, entries *android.AndroidMkEntries) {
					entries.SetString("LOCAL_MODULE_PATH", f.installDir.String())
					entries.SetString("LOCAL_INSTALLED_MODULE_STEM", f.outputFile.Base())
				},
			},
		},
	}
}
