SOONG := $(OUT_DIR)/soong
SOONG_BUILD_NINJA := $(OUT_DIR)/build.ninja
SOONG_ANDROID_MK := $(OUT_DIR)/Android.mk
SOONG_ENV := . $(OUT_DIR)/.soong.bootstrap &&
SOONG_VARIABLES := $(OUT_DIR)/soong.variables

# Bootstrap soong.  Run only the first time for clean builds
$(SOONG):
	$(hide) mkdir -p $(dir $@)
	$(hide) cd $(dir $@) && ../bootstrap.bash # TODO: relpath from $OUT_DIR

# Create soong.variables with copies of makefile settings.  Runs every build,
# but only updates soong.variables if it changes
SOONG_VARIABLES_TMP := $(SOONG_VARIABLES).$$$$
$(SOONG_VARIABLES): FORCE
	$(hide) (\
	echo '{'; \
	echo '    "Device_uses_jemalloc": $(if $(filter dlmalloc,$(MALLOC_IMPL)),false,true),'; \
	echo '    "Device_uses_dlmalloc": $(if $(filter dlmalloc,$(MALLOC_IMPL)),true,false),'; \
	echo '    $(if $(BOARD_MALLOC_ALIGNMENT),"Dlmalloc_alignment": $(BOARD_MALLOC_ALIGNMENT)$(comma),)'; \
	echo '    "Platform_sdk_version": $(PLATFORM_SDK_VERSION),'; \
	echo ''; \
	echo '    "DeviceName": "$(TARGET_DEVICE)",'; \
	echo '    "DeviceArch": "$(TARGET_ARCH)",'; \
	echo '    "DeviceArchVariant": "$(TARGET_ARCH_VARIANT)",'; \
	echo '    "DeviceCpuVariant": "$(TARGET_CPU_VARIANT)",'; \
	echo '    "DeviceAbi": ["$(TARGET_CPU_ABI)", "$(TARGET_CPU_ABI2)"],'; \
	echo '    "DeviceUsesClang": $(if $(USE_CLANG_PLATFORM_BUILD),$(USE_CLANG_PLATFORM_BUILD),false),'; \
	echo ''; \
	echo '    "DeviceSecondaryArch": "$(TARGET_2ND_ARCH)",'; \
	echo '    "DeviceSecondaryArchVariant": "$(TARGET_2ND_ARCH_VARIANT)",'; \
	echo '    "DeviceSecondaryCpuVariant": "$(TARGET_2ND_CPU_VARIANT)",'; \
	echo '    "DeviceSecondaryAbi": ["$(TARGET_2ND_CPU_ABI)", "$(TARGET_2ND_CPU_ABI2)"],'; \
	echo ''; \
	echo '    "HostArch": "$(HOST_ARCH)",'; \
	echo '    "HostSecondaryArch": "$(HOST_2ND_ARCH)"'; \
	echo '}') > $(SOONG_VARIABLES_TMP); \
	if ! cmp -s $(SOONG_VARIABLES_TMP) $(SOONG_VARIABLES); then \
	  mv $(SOONG_VARIABLES_TMP) $(SOONG_VARIABLES); \
	else \
	  rm $(SOONG_VARIABLES_TMP); \
	fi

# Build an Android.mk listing all soong outputs as prebuilts
$(SOONG_ANDROID_MK): $(SOONG) $(SOONG_VARIABLES) FORCE
	$(hide) $(SOONG) $(TOP)/$@ $(NINJA_ARGS)
