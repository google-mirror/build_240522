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

ifndef BOARD_VNDK_VERSION
snapshot_file_name := android-vndk-$(TARGET_ARCH).zip
else
snapshot_file_name := android-vndk-$(BOARD_VNDK_VERSION)-$(TARGET_ARCH).zip
endif

vndk_snapshot_zip := $(PRODUCT_OUT)/$(snapshot_file_name)
vndk_snapshot_out := $(PRODUCT_OUT)/vndk-snapshot
vndk_core_libs := $(addsuffix .vendor,$(filter-out libclang_rt.%,$(VNDK_CORE_LIBRARIES)))
vndk_sp_libs := $(addsuffix .vendor,$(VNDK_SAMEPROCESS_LIBRARIES))
vndk_snapshot_dependencies := \
  $(vndk_core_libs) \
  $(vndk_sp_libs)

# $(1): if not empty, evaluates for TARGET_2ND_ARCH
define vndk_core_libs_additional
  $(eval prefix := $(if $(1),2ND_,)) \
  $(foreach lib,ADDRESS_SANITIZER_RUNTIME_LIBRARY UBSAN_RUNTIME_LIBRARY, \
    $(addsuffix .vendor,$($(addprefix $(prefix),$(lib)))))
endef

# $(1): list of lib names without '.so' suffix (e.g., libX.vendor)
# $(2): if not empty, evaluates for TARGET_2ND_ARCH
define vndk_intermediates
  $(strip \
    $(foreach lib,$(1), \
      $(call append-path,$(call intermediates-dir-for,SHARED_LIBRARIES,$(lib),,,$(2)),$(lib).so)))
endef

vndk_core_out_$(TARGET_ARCH) := $(vndk_snapshot_out)/arch-$(TARGET_ARCH)/lib/vndk-core
vndk_sp_out_$(TARGET_ARCH) := $(vndk_snapshot_out)/arch-$(TARGET_ARCH)/lib/vndk-sp
vndk_core_libs_additional_$(TARGET_ARCH) := $(call vndk_core_libs_additional)
vndk_core_intermediates_$(TARGET_ARCH) := $(call vndk_intermediates,$(vndk_core_libs) $(vndk_core_libs_additional_$(TARGET_ARCH)))
vndk_sp_intermediates_$(TARGET_ARCH) := $(call vndk_intermediates,$(vndk_sp_libs))
vndk_snapshot_dependencies += \
  $(vndk_core_libs_additional_$(TARGET_ARCH))

ifeq (true,$(TARGET_IS_64_BIT))
vndk_core_out_$(TARGET_2ND_ARCH) := $(vndk_snapshot_out)/arch-$(TARGET_2ND_ARCH)/lib/vndk-core
vndk_sp_out_$(TARGET_2ND_ARCH) := $(vndk_snapshot_out)/arch-$(TARGET_2ND_ARCH)/lib/vndk-sp
vndk_core_libs_additional_$(TARGET_2ND_ARCH) := $(call vndk_core_libs_additional,true)
vndk_core_intermediates_$(TARGET_2ND_ARCH) := $(call vndk_intermediates,$(vndk_core_libs) $(vndk_core_libs_additional_$(TARGET_2ND_ARCH)),true)
vndk_sp_intermediates_$(TARGET_2ND_ARCH) := $(call vndk_intermediates,$(vndk_sp_libs),true)
vndk_snapshot_dependencies += \
  $(vndk_core_libs_additional_$(TARGET_2ND_ARCH))
endif

# $(1): $(TARGET_ARCH) or $(TARGET_2ND_ARCH)
define copy_vndk_libs
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
	@mkdir -p $(vndk_snapshot_out)
	$(call copy_vndk_libs,$(TARGET_ARCH))
ifeq (true,$(TARGET_IS_64_BIT))
	$(call copy_vndk_libs,$(TARGET_2ND_ARCH))
endif
	$(hide) $(SOONG_ZIP) -o $@ -P $(notdir $(vndk_snapshot_out)) -C $(vndk_snapshot_out) \
		-D $(vndk_snapshot_out)
	@rm -rf $(vndk_snapshot_out)

.PHONY: vndk
vndk: $(vndk_snapshot_zip)

$(call dist-for-goals, vndk, $(vndk_snapshot_zip))
