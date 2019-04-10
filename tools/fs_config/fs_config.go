package fs_config

import (
	"android/soong/android"
	"android/soong/genrule"
)

func init() {
	android.RegisterModuleType("fs_config_genrule", fsConfigGenruleFactory)
}

func targetFsConfigGenFunc(ctx android.ModuleContext) (val string, deps android.Paths) {
	if ret := ctx.DeviceConfig().TargetFSConfigGen(); ret != nil && *ret != "" {
		return *ret, []android.Path{android.PathForSource(ctx, *ret)}
	}
	return "/dev/null", nil
}

func fsConfigGenruleFactory() android.Module {
	module := genrule.NewGenRule()
	module.AddCustomFunc("target_fs_config_gen", targetFsConfigGenFunc)
	android.InitAndroidModule(module)
	android.InitDefaultableModule(module)
	return module
}
