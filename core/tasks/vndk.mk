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

# Returns arch-specific libclang_rt.ubsan* library name.
# Because VNDK_CORE_LIBRARIES includes all arch variants for libclang_rt.ubsan*
# libs, the arch-specific libs are selected separately.
#
# Args:
#   $(1): if not empty, evaluates for TARGET_2ND_ARCH
define clang-ubsan-vndk-core
$(eval prefix := $(if $(1),2ND_,))
$(addsuffix .vendor,$($(addprefix $(prefix),UBSAN_RUNTIME_LIBRARY)))
endef

# Returns list of file paths of the intermediate objs
#
# Args:
#   $(1): list of obj names (e.g., libfoo.vendor, ld.config.txt, ...)
#   $(2): target class (e.g., SHARED_LIBRARIES, STATIC_LIBRARIES, ETC)
#   $(3): if not empty, evaluates for TARGET_2ND_ARCH
define paths-of-intermediates
$(strip \
  $(foreach obj,$(1), \
    $(eval file_name := $(if $(filter SHARED_LIBRARIES,$(2)),$(patsubst %.so,%,$(obj)).so,$(obj))) \
    $(eval dir := $(call intermediates-dir-for,$(2),$(obj),,,$(3))) \
    $(call append-path,$(dir),$(file_name)) \
  ) \
)
endef

vndk_core_libs := $(addsuffix .vendor,$(filter-out libclang_rt.ubsan%,$(VNDK_CORE_LIBRARIES)))
vndk_sp_libs := $(addsuffix .vendor,$(VNDK_SAMEPROCESS_LIBRARIES))
vndk_snapshot_libs := \
  $(vndk_core_libs) \
  $(vndk_sp_libs)
vndk_snapshot_txts := \
  ld.config.txt \
  vndksp.libraries.txt \
  llndk.libraries.txt

# If in the future libclang_rt.ubsan* is removed from the VNDK-core list,
# need to update the related logic in this file.
ifeq (,$(filter libclang_rt.ubsan%,$(VNDK_CORE_LIBRARIES)))
  $(error libclang_rt.ubsan* is no longer a VNDK-core library)
endif

# for TARGET_ARCH
clang_ubsan_vndk_core := $(call clang-ubsan-vndk-core)
vndk_snapshot_libs += \
  $(clang_ubsan_vndk_core)

# ifdef TARGET_2ND_ARCH
# clang_ubsan_vndk_core_2ND := $(call clang-ubsan-vndk-core,true)
# vndk_snapshot_libs += \
#   $(clang_ubsan_vndk_core_2ND)
# endif


#######################################
# vndk_snapshot_zip
vndk_snapshot_top := $(call intermediates-dir-for,PACKAGING,vndk-snapshot)
vndk_snapshot_arch := $(vndk_snapshot_top)/arch-$(TARGET_ARCH)-$(TARGET_ARCH_VARIANT)
vndk_snapshot_zip := $(PRODUCT_OUT)/android-vndk-$(TARGET_ARCH).zip
$(vndk_snapshot_zip): PRIVATE_VNDK_SNAPSHOT_TOP := $(vndk_snapshot_top)
$(vndk_snapshot_zip): PRIVATE_VNDK_CORE_OUT := $(vndk_snapshot_arch)/shared/vndk-core
$(vndk_snapshot_zip): PRIVATE_VNDK_CORE_INTERMEDIATES := \
  $(call paths-of-intermediates,$(vndk_core_libs) $(clang_ubsan_vndk_core),SHARED_LIBRARIES)
$(vndk_snapshot_zip): PRIVATE_VNDK_SP_OUT := $(vndk_snapshot_arch)/shared/vndk-sp
$(vndk_snapshot_zip): PRIVATE_VNDK_SP_INTERMEDIATES := \
  $(call paths-of-intermediates,$(vndk_sp_libs),SHARED_LIBRARIES)
$(vndk_snapshot_zip): PRIVATE_TEXT_FILES_OUT := $(vndk_snapshot_arch)
$(vndk_snapshot_zip): PRIVATE_TEXT_FILES_INTERMEDIATES := \
  $(call paths-of-intermediates,$(vndk_snapshot_txts),ETC)
$(vndk_snapshot_zip): PRIVATE_CONFIGS_OUT := $(vndk_snapshot_arch)/configs
$(vndk_snapshot_zip): PRIVATE_CONFIGS_INTERMEDIATES := $(export_includes)

# TODO(jaeshin): Package additional arch variants such as (arch, variant)=(arm, armv8-a)
# ifdef TARGET_2ND_ARCH
# vndk_snapshot_arch_2ND := $(vndk_snapshot_top)/arch-$(TARGET_2ND_ARCH)-$(TARGET_2ND_ARCH_VARIANT)
# $(vndk_snapshot_zip): PRIVATE_VNDK_CORE_OUT_2ND := $(vndk_snapshot_arch_2ND)/shared/vndk-core
# $(vndk_snapshot_zip): PRIVATE_VNDK_CORE_INTERMEDIATES_2ND := \
#   $(call paths-of-intermediates,$(vndk_core_libs) $(clang_ubsan_vndk_core_2ND),SHARED_LIBRARIES,true)
# $(vndk_snapshot_zip): PRIVATE_VNDK_SP_OUT_2ND := $(vndk_snapshot_arch_2ND)/shared/vndk-sp
# $(vndk_snapshot_zip): PRIVATE_VNDK_SP_INTERMEDIATES_2ND := \
#   $(call paths-of-intermediates,$(vndk_sp_libs),SHARED_LIBRARIES,true)
# $(vndk_snapshot_zip): PRIVATE_TEXT_FILES_OUT_2ND := $(vndk_snapshot_arch_2ND)
# $(vndk_snapshot_zip): PRIVATE_TEXT_FILES_INTERMEDIATES_2ND := \
#   $(call paths-of-intermediates,$(vndk_snapshot_txts),ETC,true)
# endif

# Args
#   $(1): destination directory
#   $(2): list of files to copy
$(vndk_snapshot_zip): private-copy-vndk-intermediates = \
	$(if $(2),$(strip \
	  @mkdir -p $(1); \
	  $(foreach file,$(2),$(if $(wildcard $(file)),cp -p $(file) $(call append-path,$(1),$(subst .vendor,,$(notdir $(file))));)) \
	))

$(vndk_snapshot_zip): $(vndk_snapshot_libs) $(vndk_snapshot_txts) $(SOONG_ZIP)
	@echo 'Generating VNDK snapshot: $@'
	@rm -f $@
	@rm -rf $(PRIVATE_VNDK_SNAPSHOT_TOP)
	@mkdir -p $(PRIVATE_VNDK_SNAPSHOT_TOP)
	$(call private-copy-vndk-intermediates, \
		$(PRIVATE_VNDK_CORE_OUT),$(PRIVATE_VNDK_CORE_INTERMEDIATES))
	$(call private-copy-vndk-intermediates, \
		$(PRIVATE_VNDK_SP_OUT),$(PRIVATE_VNDK_SP_INTERMEDIATES))
	$(call private-copy-vndk-intermediates, \
		$(PRIVATE_TEXT_FILES_OUT),$(PRIVATE_TEXT_FILES_INTERMEDIATES))
# ifdef TARGET_2ND_ARCH
# 	$(call private-copy-vndk-intermediates, \
# 		$(PRIVATE_VNDK_CORE_OUT_2ND),$(PRIVATE_VNDK_CORE_INTERMEDIATES_2ND))
# 	$(call private-copy-vndk-intermediates, \
# 		$(PRIVATE_VNDK_SP_OUT_2ND),$(PRIVATE_VNDK_SP_INTERMEDIATES_2ND))
# 	$(call private-copy-vndk-intermediates, \
# 		$(PRIVATE_TEXT_FILES_OUT_2ND),$(PRIVATE_TEXT_FILES_INTERMEDIATES_2ND))
# endif
	$(hide) $(SOONG_ZIP) -o $@ -P android-vndk-snapshot -C $(PRIVATE_VNDK_SNAPSHOT_TOP) \
	-D $(PRIVATE_VNDK_SNAPSHOT_TOP)

.PHONY: vndk
vndk: $(vndk_snapshot_zip)

$(call dist-for-goals, vndk, $(vndk_snapshot_zip))

# clear global vars
clang-ubsan-vndk-core :=
paths-of-intermediates :=
vndk_core_libs :=
vndk_sp_libs :=
vndk_snapshot_libs :=
vndk_snapshot_txts :=
vndk_snapshot_top :=
clang_ubsan_vndk_core :=
vndk_snapshot_arch :=
# ifdef TARGET_2ND_ARCH
# clang_ubsan_vndk_core_2ND :=
# vndk_snapshot_arch_2ND :=
# endif

else # BOARD_VNDK_VERSION is NOT set to 'current'

.PHONY: vndk
vndk:
	$(call echo-error,$(current_makefile),CANNOT generate VNDK snapshot. BOARD_VNDK_VERSION must be set to 'current'.)
	exit 1

endif # BOARD_VNDK_VERSION
