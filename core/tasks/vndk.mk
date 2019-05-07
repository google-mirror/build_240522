# Copyright (C) 2017 The Android Open Source Project
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

current_makefile := $(lastword $(MAKEFILE_LIST))

# BOARD_VNDK_VERSION must be set to 'current' in order to generate a VNDK snapshot.
ifeq ($(BOARD_VNDK_VERSION),current)

# PLATFORM_VNDK_VERSION must be set.
ifneq (,$(PLATFORM_VNDK_VERSION))

# BOARD_VNDK_RUNTIME_DISABLE must not be set to 'true'.
ifneq ($(BOARD_VNDK_RUNTIME_DISABLE),true)

# Returns list of src:dest paths of the intermediate objs
#
# Args:
#   $(1): list of module and filename pairs (e.g., ld.config.txt:ld.config.27.txt ...)
define paths-of-intermediates
$(strip \
  $(foreach pair,$(1), \
    $(eval module := $(call word-colon,1,$(pair))) \
    $(eval built := $(ALL_MODULES.$(module).BUILT_INSTALLED)) \
    $(eval filename := $(call word-colon,2,$(pair))) \
    $(if $(wordlist 2,100,$(built)), \
      $(error Unable to handle multiple built files ($(module)): $(built))) \
    $(if $(built),$(call word-colon,1,$(built)):$(filename)) \
  ) \
)
endef

vndk_prebuilt_txts := \
  ld.config.txt \
  vndksp.libraries.txt \
  llndk.libraries.txt

vndk_snapshot_top := $(call intermediates-dir-for,PACKAGING,vndk-snapshot)
vndk_snapshot_out := $(vndk_snapshot_top)/vndk-snapshot
vndk_snapshot_configs_out := $(vndk_snapshot_top)/configs
soong_vndk_snapshot_out := $(SOONG_OUT_DIR)/vndk-snapshot/vndk-snapshot

vndk_snapshot_configs := \
  $(vndkcore.libraries.txt) \
  $(vndkprivate.libraries.txt) \
  $(module_paths.txt)

#######################################
# vndk_snapshot_zip
vndk_snapshot_variant := $(vndk_snapshot_out)/$(TARGET_ARCH)
vndk_snapshot_zip := $(PRODUCT_OUT)/android-vndk-$(TARGET_PRODUCT).zip

$(vndk_snapshot_zip): PRIVATE_VNDK_SNAPSHOT_OUT := $(vndk_snapshot_out)
$(vndk_snapshot_zip): PRIVATE_SOONG_VNDK_SNAPSHOT_OUT := $(soong_vndk_snapshot_out)
$(vndk_snapshot_zip): $(soong_vndk_snapshot_out)

deps := $(call paths-of-intermediates,$(foreach txt,$(vndk_prebuilt_txts), \
          $(txt):$(patsubst %.txt,%.$(PLATFORM_VNDK_VERSION).txt,$(txt)))) \
        $(foreach config,$(vndk_snapshot_configs),$(config):$(notdir $(config)))
$(vndk_snapshot_zip): PRIVATE_CONFIGS_OUT := $(vndk_snapshot_variant)/configs
$(vndk_snapshot_zip): PRIVATE_CONFIGS_INTERMEDIATES := $(deps)
$(vndk_snapshot_zip): $(foreach d,$(deps),$(call word-colon,1,$(d)))
deps :=

# Args
#   $(1): destination directory
#   $(2): list of files (src:dest) to copy
$(vndk_snapshot_zip): private-copy-intermediates = \
  $(if $(2),$(strip \
    @mkdir -p $(1) && \
    $(foreach file,$(2), \
      cp $(call word-colon,1,$(file)) $(call append-path,$(1),$(call word-colon,2,$(file))) && \
    ) \
    true \
  ))


$(vndk_snapshot_zip): $(SOONG_ZIP)
	@echo 'Generating VNDK snapshot: $@'
	@rm -f $@
	@rm -rf $(PRIVATE_VNDK_SNAPSHOT_OUT)
	@mkdir -p $(PRIVATE_VNDK_SNAPSHOT_OUT)
	@cp -r $(PRIVATE_SOONG_VNDK_SNAPSHOT_OUT)/* $(PRIVATE_VNDK_SNAPSHOT_OUT)
	$(call private-copy-intermediates, \
		$(PRIVATE_CONFIGS_OUT),$(PRIVATE_CONFIGS_INTERMEDIATES))
	$(hide) $(SOONG_ZIP) -o $@ -C $(PRIVATE_VNDK_SNAPSHOT_OUT) -D $(PRIVATE_VNDK_SNAPSHOT_OUT)

.PHONY: vndk
vndk: $(vndk_snapshot_zip)

$(call dist-for-goals, vndk, $(vndk_snapshot_zip))

# clear global vars
clang-ubsan-vndk-core :=
paths-of-intermediates :=
paths-of-notice-files :=
vndk_core_libs :=
vndk_sp_libs :=
vndk_snapshot_libs :=
vndk_prebuilt_txts :=
vndk_snapshot_configs :=
vndk_snapshot_top :=
vndk_snapshot_out :=
vndk_snapshot_configs_out :=
vndk_snapshot_variant :=
binder :=
vndk_lib_dir :=
vndk_lib_dir_2nd :=

else # BOARD_VNDK_RUNTIME_DISABLE is set to 'true'
error_msg := "CANNOT generate VNDK snapshot. BOARD_VNDK_RUNTIME_DISABLE must not be set to 'true'."
endif # BOARD_VNDK_RUNTIME_DISABLE

else # PLATFORM_VNDK_VERSION is NOT set
error_msg := "CANNOT generate VNDK snapshot. PLATFORM_VNDK_VERSION must be set."
endif # PLATFORM_VNDK_VERSION

else # BOARD_VNDK_VERSION is NOT set to 'current'
error_msg := "CANNOT generate VNDK snapshot. BOARD_VNDK_VERSION must be set to 'current'."
endif # BOARD_VNDK_VERSION

ifneq (,$(error_msg))

.PHONY: vndk
vndk:
	$(call echo-error,$(current_makefile),$(error_msg))
	exit 1

endif
