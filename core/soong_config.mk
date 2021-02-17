SOONG_MAKEVARS_MK := $(SOONG_OUT_DIR)/make_vars-$(TARGET_PRODUCT).mk
SOONG_VARIABLES := $(SOONG_OUT_DIR)/soong.variables
SOONG_ANDROID_MK := $(SOONG_OUT_DIR)/Android-$(TARGET_PRODUCT).mk

BINDER32BIT :=
ifneq ($(TARGET_USES_64_BIT_BINDER),true)
ifneq ($(TARGET_IS_64_BIT),true)
BINDER32BIT := true
endif
endif

ifeq ($(WRITE_SOONG_VARIABLES),true)
# Converts a list to a JSON list.
# $1: List separator.
# $2: List.
_json_list = [$(if $(2),"$(subst $(1),"$(comma)",$(2))")]

# Converts a space-separated list to a JSON list.
json_list = $(call _json_list,$(space),$(1))

# Converts a comma-separated list to a JSON list.
csv_to_json_list = $(call _json_list,$(comma),$(1))

# 1: Key name
# 2: Value
add_json_val = $(eval _contents := $$(_contents)    "$$(strip $$(1))":$$(space)$$(strip $$(2))$$(comma)$$(newline))
add_json_str = $(call add_json_val,$(1),"$(strip $(2))")
add_json_list = $(call add_json_val,$(1),$(call json_list,$(patsubst %,%,$(2))))
add_json_csv = $(call add_json_val,$(1),$(call csv_to_json_list,$(strip $(2))))
add_json_bool = $(call add_json_val,$(1),$(if $(strip $(2)),true,false))

invert_bool = $(if $(strip $(1)),,true)

# Create soong.variables with copies of makefile settings.  Runs every build,
# but only updates soong.variables if it changes
$(shell mkdir -p $(dir $(SOONG_VARIABLES)))
_contents := {$(newline)

$(call add_json_str,  Make_suffix, -$(TARGET_PRODUCT))

$(call add_json_str,  BuildId,                           $(BUILD_ID))
$(call add_json_str,  BuildNumberFromFile,               $$$(BUILD_NUMBER_FROM_FILE))

$(call add_json_str,  Platform_version_name,             $(PLATFORM_VERSION))
$(call add_json_val,  Platform_sdk_version,              $(PLATFORM_SDK_VERSION))
$(call add_json_str,  Platform_sdk_codename,             $(PLATFORM_VERSION_CODENAME))
$(call add_json_bool, Platform_sdk_final,                $(filter REL,$(PLATFORM_VERSION_CODENAME)))
$(call add_json_csv,  Platform_version_active_codenames, $(PLATFORM_VERSION_ALL_CODENAMES))
$(call add_json_csv,  Platform_version_future_codenames, $(PLATFORM_VERSION_FUTURE_CODENAMES))

$(call add_json_bool, Allow_missing_dependencies,        $(ALLOW_MISSING_DEPENDENCIES))
<<<<<<< HEAD   (4be654 Merge "Merge empty history for sparse-7121469-L4290000080720)
$(call add_json_bool, Unbundled_build,                   $(TARGET_BUILD_APPS))
$(call add_json_bool, Pdk,                               $(filter true,$(TARGET_BUILD_PDK)))
=======
$(call add_json_bool, Unbundled_build,                   $(TARGET_BUILD_UNBUNDLED))
$(call add_json_bool, Unbundled_build_apps,              $(TARGET_BUILD_APPS))
$(call add_json_bool, Always_use_prebuilt_sdks,          $(TARGET_BUILD_USE_PREBUILT_SDKS))
$(call add_json_bool, Skip_boot_jars_check,              $(SKIP_BOOT_JARS_CHECK))
>>>>>>> BRANCH (fe6ad7 Merge "Version bump to RBT1.210107.001.A1 [core/build_id.mk])

$(call add_json_bool, Debuggable,                        $(filter userdebug eng,$(TARGET_BUILD_VARIANT)))
$(call add_json_bool, Eng,                               $(filter eng,$(TARGET_BUILD_VARIANT)))

$(call add_json_str,  DeviceName,                        $(TARGET_DEVICE))
$(call add_json_str,  DeviceArch,                        $(TARGET_ARCH))
$(call add_json_str,  DeviceArchVariant,                 $(TARGET_ARCH_VARIANT))
$(call add_json_str,  DeviceCpuVariant,                  $(TARGET_CPU_VARIANT))
$(call add_json_list, DeviceAbi,                         $(TARGET_CPU_ABI) $(TARGET_CPU_ABI2))

$(call add_json_str,  DeviceSecondaryArch,               $(TARGET_2ND_ARCH))
$(call add_json_str,  DeviceSecondaryArchVariant,        $(TARGET_2ND_ARCH_VARIANT))
$(call add_json_str,  DeviceSecondaryCpuVariant,         $(TARGET_2ND_CPU_VARIANT))
$(call add_json_list, DeviceSecondaryAbi,                $(TARGET_2ND_CPU_ABI) $(TARGET_2ND_CPU_ABI2))

<<<<<<< HEAD   (4be654 Merge "Merge empty history for sparse-7121469-L4290000080720)
=======
$(call add_json_bool, Aml_abis,                          $(if $(filter mainline_sdk,$(TARGET_ARCH_SUITE)),true))
$(call add_json_bool, Ndk_abis,                          $(if $(filter ndk,         $(TARGET_ARCH_SUITE)),true))

$(call add_json_str,  NativeBridgeArch,                  $(TARGET_NATIVE_BRIDGE_ARCH))
$(call add_json_str,  NativeBridgeArchVariant,           $(TARGET_NATIVE_BRIDGE_ARCH_VARIANT))
$(call add_json_str,  NativeBridgeCpuVariant,            $(TARGET_NATIVE_BRIDGE_CPU_VARIANT))
$(call add_json_list, NativeBridgeAbi,                   $(TARGET_NATIVE_BRIDGE_ABI))
$(call add_json_str,  NativeBridgeRelativePath,          $(TARGET_NATIVE_BRIDGE_RELATIVE_PATH))

$(call add_json_str,  NativeBridgeSecondaryArch,         $(TARGET_NATIVE_BRIDGE_2ND_ARCH))
$(call add_json_str,  NativeBridgeSecondaryArchVariant,  $(TARGET_NATIVE_BRIDGE_2ND_ARCH_VARIANT))
$(call add_json_str,  NativeBridgeSecondaryCpuVariant,   $(TARGET_NATIVE_BRIDGE_2ND_CPU_VARIANT))
$(call add_json_list, NativeBridgeSecondaryAbi,          $(TARGET_NATIVE_BRIDGE_2ND_ABI))
$(call add_json_str,  NativeBridgeSecondaryRelativePath, $(TARGET_NATIVE_BRIDGE_2ND_RELATIVE_PATH))

>>>>>>> BRANCH (fe6ad7 Merge "Version bump to RBT1.210107.001.A1 [core/build_id.mk])
$(call add_json_str,  HostArch,                          $(HOST_ARCH))
$(call add_json_str,  HostSecondaryArch,                 $(HOST_2ND_ARCH))
$(call add_json_bool, HostStaticBinaries,                $(BUILD_HOST_static))

$(call add_json_str,  CrossHost,                         $(HOST_CROSS_OS))
$(call add_json_str,  CrossHostArch,                     $(HOST_CROSS_ARCH))
$(call add_json_str,  CrossHostSecondaryArch,            $(HOST_CROSS_2ND_ARCH))

$(call add_json_list, ResourceOverlays,                  $(PRODUCT_PACKAGE_OVERLAYS) $(DEVICE_PACKAGE_OVERLAYS))
$(call add_json_list, EnforceRROTargets,                 $(PRODUCT_ENFORCE_RRO_TARGETS))
$(call add_json_list, EnforceRROExcludedOverlays,        $(PRODUCT_ENFORCE_RRO_EXCLUDED_OVERLAYS))

$(call add_json_str,  AAPTCharacteristics,               $(TARGET_AAPT_CHARACTERISTICS))
$(call add_json_list, AAPTConfig,                        $(PRODUCT_AAPT_CONFIG))
$(call add_json_str,  AAPTPreferredConfig,               $(PRODUCT_AAPT_PREF_CONFIG))
$(call add_json_list, AAPTPrebuiltDPI,                   $(PRODUCT_AAPT_PREBUILT_DPI))

$(call add_json_str,  DefaultAppCertificate,             $(PRODUCT_DEFAULT_DEV_CERTIFICATE))

$(call add_json_str,  AppsDefaultVersionName,            $(APPS_DEFAULT_VERSION_NAME))

$(call add_json_list, SanitizeHost,                      $(SANITIZE_HOST))
$(call add_json_list, SanitizeDevice,                    $(SANITIZE_TARGET))
$(call add_json_list, SanitizeDeviceDiag,                $(SANITIZE_TARGET_DIAG))
$(call add_json_list, SanitizeDeviceArch,                $(SANITIZE_TARGET_ARCH))

$(call add_json_bool, Safestack,                         $(filter true,$(USE_SAFESTACK)))
$(call add_json_bool, EnableCFI,                         $(call invert_bool,$(filter false,$(ENABLE_CFI))))
$(call add_json_list, CFIExcludePaths,                   $(CFI_EXCLUDE_PATHS) $(PRODUCT_CFI_EXCLUDE_PATHS))
$(call add_json_list, CFIIncludePaths,                   $(CFI_INCLUDE_PATHS) $(PRODUCT_CFI_INCLUDE_PATHS))
$(call add_json_list, IntegerOverflowExcludePaths,       $(INTEGER_OVERFLOW_EXCLUDE_PATHS) $(PRODUCT_INTEGER_OVERFLOW_EXCLUDE_PATHS))

$(call add_json_bool, ClangTidy,                         $(filter 1 true,$(WITH_TIDY)))
$(call add_json_str,  TidyChecks,                        $(WITH_TIDY_CHECKS))

<<<<<<< HEAD   (4be654 Merge "Merge empty history for sparse-7121469-L4290000080720)
$(call add_json_bool, NativeCoverage,                    $(filter true,$(NATIVE_COVERAGE)))
$(call add_json_list, CoveragePaths,                     $(COVERAGE_PATHS))
$(call add_json_list, CoverageExcludePaths,              $(COVERAGE_EXCLUDE_PATHS))
=======
$(call add_json_list, JavaCoveragePaths,                 $(JAVA_COVERAGE_PATHS))
$(call add_json_list, JavaCoverageExcludePaths,          $(JAVA_COVERAGE_EXCLUDE_PATHS))

$(call add_json_bool, GcovCoverage,                      $(filter true,$(NATIVE_COVERAGE)))
$(call add_json_bool, ClangCoverage,                     $(filter true,$(CLANG_COVERAGE)))
$(call add_json_list, NativeCoveragePaths,               $(NATIVE_COVERAGE_PATHS))
$(call add_json_list, NativeCoverageExcludePaths,        $(NATIVE_COVERAGE_EXCLUDE_PATHS))

$(call add_json_bool, SamplingPGO,                       $(filter true,$(SAMPLING_PGO)))
>>>>>>> BRANCH (fe6ad7 Merge "Version bump to RBT1.210107.001.A1 [core/build_id.mk])

$(call add_json_bool, ArtUseReadBarrier,                 $(call invert_bool,$(filter false,$(PRODUCT_ART_USE_READ_BARRIER))))
$(call add_json_bool, Binder32bit,                       $(BINDER32BIT))
$(call add_json_bool, Brillo,                            $(BRILLO))
$(call add_json_str,  BtConfigIncludeDir,                $(BOARD_BLUETOOTH_BDROID_BUILDCFG_INCLUDE_DIR))
<<<<<<< HEAD   (4be654 Merge "Merge empty history for sparse-7121469-L4290000080720)
$(call add_json_bool, Device_uses_hwc2,                  $(filter true,$(TARGET_USES_HWC2)))
$(call add_json_list, DeviceKernelHeaders,               $(TARGET_PROJECT_SYSTEM_INCLUDES))
$(call add_json_bool, DevicePrefer32BitExecutables,      $(filter true,$(TARGET_PREFER_32_BIT_EXECUTABLES)))
=======
$(call add_json_list, DeviceKernelHeaders,               $(TARGET_DEVICE_KERNEL_HEADERS) $(TARGET_BOARD_KERNEL_HEADERS) $(TARGET_PRODUCT_KERNEL_HEADERS))
>>>>>>> BRANCH (fe6ad7 Merge "Version bump to RBT1.210107.001.A1 [core/build_id.mk])
$(call add_json_str,  DeviceVndkVersion,                 $(BOARD_VNDK_VERSION))
$(call add_json_str,  Platform_vndk_version,             $(PLATFORM_VNDK_VERSION))
$(call add_json_list, ExtraVndkVersions,                 $(PRODUCT_EXTRA_VNDK_VERSIONS))
$(call add_json_list, DeviceSystemSdkVersions,           $(BOARD_SYSTEMSDK_VERSIONS))
$(call add_json_str,  RecoverySnapshotVersion,           $(RECOVERY_SNAPSHOT_VERSION))
$(call add_json_list, Platform_systemsdk_versions,       $(PLATFORM_SYSTEMSDK_VERSIONS))
$(call add_json_bool, Malloc_not_svelte,                 $(call invert_bool,$(filter true,$(MALLOC_SVELTE))))
$(call add_json_bool, Malloc_zero_contents,              $(call invert_bool,$(filter false,$(MALLOC_ZERO_CONTENTS))))
$(call add_json_bool, Malloc_pattern_fill_contents,      $(MALLOC_PATTERN_FILL_CONTENTS))
$(call add_json_str,  Override_rs_driver,                $(OVERRIDE_RS_DRIVER))

$(call add_json_bool, Treble_linker_namespaces,          $(filter true,$(PRODUCT_TREBLE_LINKER_NAMESPACES)))
$(call add_json_bool, Enforce_vintf_manifest,            $(filter true,$(PRODUCT_ENFORCE_VINTF_MANIFEST)))

$(call add_json_bool, Uml,                               $(filter true,$(TARGET_USER_MODE_LINUX)))
$(call add_json_bool, Use_lmkd_stats_log,                $(filter true,$(TARGET_LMKD_STATS_LOG)))
$(call add_json_str,  VendorPath,                        $(TARGET_COPY_OUT_VENDOR))
$(call add_json_str,  OdmPath,                           $(TARGET_COPY_OUT_ODM))
$(call add_json_str,  VendorDlkmPath,                    $(TARGET_COPY_OUT_VENDOR_DLKM))
$(call add_json_str,  OdmDlkmPath,                       $(TARGET_COPY_OUT_ODM_DLKM))
$(call add_json_str,  ProductPath,                       $(TARGET_COPY_OUT_PRODUCT))
$(call add_json_bool, MinimizeJavaDebugInfo,             $(filter true,$(PRODUCT_MINIMIZE_JAVA_DEBUG_INFO)))

$(call add_json_bool, UseGoma,                           $(filter-out false,$(USE_GOMA)))
$(call add_json_bool, Arc,                               $(filter true,$(TARGET_ARC)))

$(call add_json_str,  DistDir,                           $(if $(dist_goal), $(DIST_DIR)))

$(call add_json_list, NamespacesToExport,                $(PRODUCT_SOONG_NAMESPACES))

$(call add_json_list, PgoAdditionalProfileDirs,          $(PGO_ADDITIONAL_PROFILE_DIRS))

<<<<<<< HEAD   (4be654 Merge "Merge empty history for sparse-7121469-L4290000080720)
_contents := $(_contents)    "VendorVars": {$(newline)
=======
$(call add_json_list, BoardReqdMaskPolicy,               $(BOARD_REQD_MASK_POLICY))
$(call add_json_list, BoardVendorSepolicyDirs,           $(BOARD_VENDOR_SEPOLICY_DIRS) $(BOARD_SEPOLICY_DIRS))
$(call add_json_list, BoardOdmSepolicyDirs,              $(BOARD_ODM_SEPOLICY_DIRS))
$(call add_json_list, BoardVendorDlkmSepolicyDirs,       $(BOARD_VENDOR_DLKM_SEPOLICY_DIRS))
$(call add_json_list, BoardOdmDlkmSepolicyDirs,          $(BOARD_ODM_DLKM_SEPOLICY_DIRS))
# TODO: BOARD_PLAT_* dirs only kept for compatibility reasons. Will be a hard error on API level 31
$(call add_json_list, SystemExtPublicSepolicyDirs,       $(SYSTEM_EXT_PUBLIC_SEPOLICY_DIRS) $(BOARD_PLAT_PUBLIC_SEPOLICY_DIR))
$(call add_json_list, SystemExtPrivateSepolicyDirs,      $(SYSTEM_EXT_PRIVATE_SEPOLICY_DIRS) $(BOARD_PLAT_PRIVATE_SEPOLICY_DIR))
$(call add_json_list, BoardSepolicyM4Defs,               $(BOARD_SEPOLICY_M4DEFS))
$(call add_json_str,  BoardSepolicyVers,                 $(BOARD_SEPOLICY_VERS))

$(call add_json_bool, Flatten_apex,                      $(filter true,$(TARGET_FLATTEN_APEX)))
$(call add_json_bool, ForceApexSymlinkOptimization,      $(filter true,$(TARGET_FORCE_APEX_SYMLINK_OPTIMIZATION)))

$(call add_json_str,  DexpreoptGlobalConfig,             $(DEX_PREOPT_CONFIG))

$(call add_json_list, ManifestPackageNameOverrides,      $(PRODUCT_MANIFEST_PACKAGE_NAME_OVERRIDES))
$(call add_json_list, PackageNameOverrides,              $(PRODUCT_PACKAGE_NAME_OVERRIDES))
$(call add_json_list, CertificateOverrides,              $(PRODUCT_CERTIFICATE_OVERRIDES))

$(call add_json_bool, EnforceSystemCertificate,          $(ENFORCE_SYSTEM_CERTIFICATE))
$(call add_json_list, EnforceSystemCertificateAllowList, $(ENFORCE_SYSTEM_CERTIFICATE_ALLOW_LIST))

$(call add_json_list, ProductHiddenAPIStubs,             $(PRODUCT_HIDDENAPI_STUBS))
$(call add_json_list, ProductHiddenAPIStubsSystem,       $(PRODUCT_HIDDENAPI_STUBS_SYSTEM))
$(call add_json_list, ProductHiddenAPIStubsTest,         $(PRODUCT_HIDDENAPI_STUBS_TEST))

$(call add_json_list, ProductPublicSepolicyDirs,         $(PRODUCT_PUBLIC_SEPOLICY_DIRS))
$(call add_json_list, ProductPrivateSepolicyDirs,        $(PRODUCT_PRIVATE_SEPOLICY_DIRS))

$(call add_json_list, TargetFSConfigGen,                 $(TARGET_FS_CONFIG_GEN))

$(call add_json_list, MissingUsesLibraries,              $(INTERNAL_PLATFORM_MISSING_USES_LIBRARIES))

$(call add_json_map, VendorVars)
>>>>>>> BRANCH (fe6ad7 Merge "Version bump to RBT1.210107.001.A1 [core/build_id.mk])
$(foreach namespace,$(SOONG_CONFIG_NAMESPACES),\
  $(eval _contents := $$(_contents)        "$(namespace)": {$$(newline)) \
  $(foreach key,$(SOONG_CONFIG_$(namespace)),\
    $(eval _contents := $$(_contents)            "$(key)": "$(SOONG_CONFIG_$(namespace)_$(key))",$$(newline)))\
  $(eval _contents := $$(_contents)$(if $(strip $(SOONG_CONFIG_$(namespace))),__SV_END)        },$$(newline)))
_contents := $(_contents)$(if $(strip $(SOONG_CONFIG_NAMESPACES)),__SV_END)    },$(newline)

<<<<<<< HEAD   (4be654 Merge "Merge empty history for sparse-7121469-L4290000080720)
_contents := $(subst $(comma)$(newline)__SV_END,$(newline),$(_contents)__SV_END}$(newline))
=======
$(call add_json_bool, EnforceProductPartitionInterface,  $(PRODUCT_ENFORCE_PRODUCT_PARTITION_INTERFACE))
$(call add_json_str,  DeviceCurrentApiLevelForVendorModules,  $(BOARD_CURRENT_API_LEVEL_FOR_VENDOR_MODULES))

$(call add_json_bool, EnforceInterPartitionJavaSdkLibrary, $(PRODUCT_ENFORCE_INTER_PARTITION_JAVA_SDK_LIBRARY))
$(call add_json_list, InterPartitionJavaLibraryAllowList, $(PRODUCT_INTER_PARTITION_JAVA_LIBRARY_ALLOWLIST))
>>>>>>> BRANCH (fe6ad7 Merge "Version bump to RBT1.210107.001.A1 [core/build_id.mk])

<<<<<<< HEAD   (4be654 Merge "Merge empty history for sparse-7121469-L4290000080720)
$(file >$(SOONG_VARIABLES).tmp,$(_contents))
=======
$(call add_json_bool, InstallExtraFlattenedApexes, $(PRODUCT_INSTALL_EXTRA_FLATTENED_APEXES))

$(call add_json_bool, CompressedApex, $(PRODUCT_COMPRESSED_APEX))

$(call add_json_bool, BoardUsesRecoveryAsBoot, $(BOARD_USES_RECOVERY_AS_BOOT))

$(call add_json_list, BoardKernelBinaries, $(BOARD_KERNEL_BINARIES))
$(call add_json_list, BoardKernelModuleInterfaceVersions, $(BOARD_KERNEL_MODULE_INTERFACE_VERSIONS))

$(call add_json_bool, BoardMoveRecoveryResourcesToVendorBoot, $(BOARD_MOVE_RECOVERY_RESOURCES_TO_VENDOR_BOOT))

$(call json_end)

$(file >$(SOONG_VARIABLES).tmp,$(json_contents))
>>>>>>> BRANCH (fe6ad7 Merge "Version bump to RBT1.210107.001.A1 [core/build_id.mk])

$(shell if ! cmp -s $(SOONG_VARIABLES).tmp $(SOONG_VARIABLES); then \
	  mv $(SOONG_VARIABLES).tmp $(SOONG_VARIABLES); \
	else \
	  rm $(SOONG_VARIABLES).tmp; \
	fi)

_json_list :=
json_list :=
csv_to_json_list :=
add_json_val :=
add_json_str :=
add_json_list :=
add_json_csv :=
add_json_bool :=
invert_bool :=
_contents :=

endif # CONFIGURE_SOONG
