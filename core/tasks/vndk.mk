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

include $(CLEAR_VARS)

ifneq ($(BOARD_VNDK_VERSION),current)
$(error BOARD_VNDK_VERSION must be set to 'current' in order to generate a VNDK snapshot. \
  Please add 'BOARD_VNDK_VERSION=current' to the make command)
endif

snapshot_file_name := android-vndk-$(TARGET_ARCH).zip
vndk_snapshot_zip := $(PRODUCT_OUT)/$(snapshot_file_name)
vndk_snapshot_out := $(call intermediates-dir-for,PACKAGING,vndk-snapshot)
vndk_core_libs := $(addsuffix .vendor,$(filter-out libclang_rt.%,$(VNDK_CORE_LIBRARIES)))
vndk_sp_libs := $(addsuffix .vendor,$(VNDK_SAMEPROCESS_LIBRARIES))
vndk_snapshot_dependencies := \
  $(vndk_core_libs) \
  $(vndk_sp_libs)

# Since VNDK_CORE_LIBRARIES includes all arch variants for libclang_rt.* libs,
# the arch-specific libclang libs are selected separately.
#
# Args:
# 	$(1): if not empty, evaluates for TARGET_2ND_ARCH
define libclang-vndk-core-libs
  $(eval libclang_soong_var_names := \
           ADDRESS_SANITIZER_RUNTIME_LIBRARY \
           UBSAN_RUNTIME_LIBRARY) \
  $(eval prefix := $(if $(1),2ND_,)) \
  $(foreach lib,$(libclang_soong_var_names), \
    $(addsuffix .vendor,$($(addprefix $(prefix),$(lib)))))
endef

# $(1): list of lib names without '.so' suffix (e.g., libX.vendor)
# $(2): if not empty, evaluates for TARGET_2ND_ARCH
define paths-for-vndk-intermediates
  $(strip \
    $(foreach lib,$(1), \
      $(call append-path,$(call intermediates-dir-for,SHARED_LIBRARIES,$(lib),,,$(2)),$(lib).so)))
endef

vndk_core_out_$(TARGET_ARCH) := $(vndk_snapshot_out)/arch-$(TARGET_ARCH)/lib/vndk-core
vndk_sp_out_$(TARGET_ARCH) := $(vndk_snapshot_out)/arch-$(TARGET_ARCH)/lib/vndk-sp
libclang-vndk-core-libs_$(TARGET_ARCH) := $(call libclang-vndk-core-libs)
vndk_core_intermediates_$(TARGET_ARCH) := $(call paths-for-vndk-intermediates,$(vndk_core_libs) $(libclang-vndk-core-libs_$(TARGET_ARCH)))
vndk_sp_intermediates_$(TARGET_ARCH) := $(call paths-for-vndk-intermediates,$(vndk_sp_libs))
vndk_snapshot_dependencies += \
  $(libclang-vndk-core-libs_$(TARGET_ARCH))

ifdef TARGET_2ND_ARCH
vndk_core_out_$(TARGET_2ND_ARCH) := $(vndk_snapshot_out)/arch-$(TARGET_2ND_ARCH)/lib/vndk-core
vndk_sp_out_$(TARGET_2ND_ARCH) := $(vndk_snapshot_out)/arch-$(TARGET_2ND_ARCH)/lib/vndk-sp
libclang-vndk-core-libs_$(TARGET_2ND_ARCH) := $(call libclang-vndk-core-libs,true)
vndk_core_intermediates_$(TARGET_2ND_ARCH) := $(call paths-for-vndk-intermediates,$(vndk_core_libs) $(libclang-vndk-core-libs_$(TARGET_2ND_ARCH)),true)
vndk_sp_intermediates_$(TARGET_2ND_ARCH) := $(call paths-for-vndk-intermediates,$(vndk_sp_libs),true)
vndk_snapshot_dependencies += \
  $(libclang-vndk-core-libs_$(TARGET_2ND_ARCH))
endif

# $(1): $(TARGET_ARCH) or $(TARGET_2ND_ARCH)
define copy-vndk-intermediates
@mkdir -p $(vndk_core_out_$(1))
$(hide) $(foreach lib, $(vndk_core_intermediates_$(1)), \
	cp -p $(lib) $(addprefix $(vndk_core_out_$(1))/,$(subst .vendor,,$(notdir $(lib))));)
@mkdir -p $(vndk_sp_out_$(1))
$(hide) $(foreach lib, $(vndk_sp_intermediates_$(1)), \
	cp -p $(lib) $(addprefix $(vndk_sp_out_$(1))/,$(subst .vendor,,$(notdir $(lib))));)
endef

$(vndk_snapshot_zip): $(vndk_snapshot_dependencies) $(SOONG_ZIP)
	@echo 'Generating VNDK snapshot...'
	@rm -f $@
	@rm -rf $(vndk_snapshot_out)
	@mkdir -p $(vndk_snapshot_out)
	$(call copy-vndk-intermediates,$(TARGET_ARCH))
ifdef TARGET_2ND_ARCH
	$(call copy-vndk-intermediates,$(TARGET_2ND_ARCH))
endif
	$(hide) $(SOONG_ZIP) -o $@ -P vndk-snapshot -C $(vndk_snapshot_out) \
		-D $(vndk_snapshot_out)

.PHONY: vndk
vndk: $(vndk_snapshot_zip)

$(call dist-for-goals, vndk, $(vndk_snapshot_zip))
