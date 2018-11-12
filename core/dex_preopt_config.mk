include $(BUILD_SYSTEM)/json.mk

DEX_PREOPT_CONFIG := $(PRODUCT_OUT)/dexpreopt.config

$(call json_start)

$(call add_json_str,  SystemPartitionPath,                $(TARGET_OUT))
$(call add_json_str,  SystemOtherPartitionPath,           $(TARGET_OUT_SYSTEM_OTHER))
$(call add_json_bool, DefaultNoStripping,                 $(filter nostripping,$(DEX_PREOPT_DEFAULT)))
$(call add_json_list, DisablePreoptModules,               $(DEXPREOPT_DISABLED_MODULES))
$(call add_json_bool, OnlyPreoptBootImageAndSystemServer, $(filter true,$(WITH_DEXPREOPT_BOOT_IMG_AND_SYSTEM_SERVER_ONLY)))
$(call add_json_bool, DontUncompressPrivAppsDex,          $(filter true,$(DONT_UNCOMPRESS_PRIV_APPS_DEXS)))
$(call add_json_list, ModulesLoadedByPrivilegedModules,   $(PRODUCT_LOADED_BY_PRIVILEGED_MODULES))
$(call add_json_bool, HasSystemOther,                     $(BOARD_USES_SYSTEM_OTHER_ODEX))
$(call add_json_list, PatternsOnSystemOther,              $(SYSTEM_OTHER_ODEX_FILTER))
$(call add_json_bool, DisableGenerateProfile,             $(filter false,$(WITH_DEX_PREOPT_GENERATE_PROFILE)))
$(call add_json_list, BootJars,                           $(DEXPREOPT_BOOT_JARS_MODULES))
$(call add_json_list, SystemServerJars,                   $(PRODUCT_SYSTEM_SERVER_JARS))
$(call add_json_list, SystemServerApps,                   $(PRODUCT_SYSTEM_SERVER_APPS))
$(call add_json_list, SpeedApps,                          $(PRODUCT_DEXPREOPT_SPEED_APPS))
$(call add_json_list, PreoptFlags,                        $(PRODUCT_DEX_PREOPT_DEFAULT_FLAGS))
$(call add_json_str,  DefaultCompilerFilter,              $(PRODUCT_DEX_PREOPT_DEFAULT_COMPILER_FILTER))
$(call add_json_str,  SystemServerCompilerFilter,         $(PRODUCT_SYSTEM_SERVER_COMPILER_FILTER))
$(call add_json_bool, GenerateDmFiles,                    $(PRODUCT_DEX_PREOPT_GENERATE_DM_FILES))
$(call add_json_bool, NoDebugInfo,                        $(filter false,$(WITH_DEXPREOPT_DEBUG_INFO)))
$(call add_json_bool, AlwaysSystemServerDebugInfo,        $(filter true,$(PRODUCT_SYSTEM_SERVER_DEBUG_INFO)))
$(call add_json_bool, NeverSystemServerDebugInfo,         $(filter false,$(PRODUCT_SYSTEM_SERVER_DEBUG_INFO)))
$(call add_json_bool, AlwaysOtherDebugInfo,               $(filter true,$(PRODUCT_OTHER_JAVA_DEBUG_INFO)))
$(call add_json_bool, NeverOtherDebugInfo,                $(filter false,$(PRODUCT_OTHER_JAVA_DEBUG_INFO)))
$(call add_json_list, MissingUsesLibraries,               $(INTERNAL_PLATFORM_MISSING_USES_LIBRARIES))
$(call add_json_bool, IsEng,                              $(filter eng,$(TARGET_BUILD_VARIANT)))
$(call add_json_bool, SanitizeLite,                       $(SANITIZE_LITE))
$(call add_json_bool, DefaultAppImages,                   $(WITH_DEX_PREOPT_APP_IMAGE))
$(call add_json_str,  Dex2oatXmx,                         $(DEX2OAT_XMX))
$(call add_json_str,  Dex2oatXms,                         $(DEX2OAT_XMS))

$(call add_json_map,  DefaultDexPreoptImageLocation)
$(call add_json_str,  $(TARGET_ARCH), $(DEFAULT_DEX_PREOPT_BUILT_IMAGE_LOCATION))
ifdef TARGET_2ND_ARCH
  $(call add_json_str, $(TARGET_2ND_ARCH), $($(TARGET_2ND_ARCH_VAR_PREFIX)DEFAULT_DEX_PREOPT_BUILT_IMAGE_LOCATION))
endif
$(call end_json_map)

$(call add_json_map,  CpuVariant)
$(call add_json_str,  $(TARGET_ARCH), $(DEX2OAT_TARGET_CPU_VARIANT))
ifdef TARGET_2ND_ARCH
  $(call add_json_str, $(TARGET_2ND_ARCH), $($(TARGET_2ND_ARCH_VAR_PREFIX)DEX2OAT_TARGET_CPU_VARIANT))
endif
$(call end_json_map)

$(call add_json_map,  InstructionSetFeatures)
$(call add_json_str,  $(TARGET_ARCH), $(DEX2OAT_TARGET_INSTRUCTION_SET_FEATURES))
ifdef TARGET_2ND_ARCH
  $(call add_json_str, $(TARGET_2ND_ARCH), $($(TARGET_2ND_ARCH_VAR_PREFIX)DEX2OAT_TARGET_INSTRUCTION_SET_FEATURES))
endif
$(call end_json_map)

$(call add_json_map,  Tools)
$(call add_json_str,  Profman,                            $(PROFMAN))
$(call add_json_str,  Dex2oat,                            $(DEX2OAT))
$(call add_json_str,  Aapt,                               $(AAPT))
$(call add_json_str,  SoongZip,                           $(SOONG_ZIP))
$(call add_json_str,  VerifyUsesLibraries,                $(BUILD_SYSTEM)/verify_uses_libraries.sh)
$(call add_json_str,  ConstructContext,                   $(BUILD_SYSTEM)/construct_context.sh)
$(call end_json_map)

$(call json_end)

$(file >$(DEX_PREOPT_CONFIG).tmp,$(json_contents))

$(shell if ! cmp -s $(DEX_PREOPT_CONFIG).tmp $(DEX_PREOPT_CONFIG); then \
	  mv $(DEX_PREOPT_CONFIG).tmp $(DEX_PREOPT_CONFIG); \
	else \
	  rm $(DEX_PREOPT_CONFIG).tmp; \
	fi)

DEXPREOPT_DEPS := \
  $(PROFMAN) \
  $(DEX2OAT) \
  $(AAPT) \
  $(SOONG_ZIP) \
  $(BUILD_SYSTEM)/verify_uses_libraries.sh \
  $(BUILD_SYSTEM)/construct_context.sh

DEXPREOPT_DEPS += $(DEFAULT_DEX_PREOPT_BUILT_IMAGE_FILENAME)
ifdef TARGET_2ND_ARCH
  DEXPREOPT_DEPS += $($(TARGET_2ND_ARCH_VAR_PREFIX)DEFAULT_DEX_PREOPT_BUILT_IMAGE_FILENAME)
endif
