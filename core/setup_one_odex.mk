#
# Copyright (C) 2014 The Android Open Source Project
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
#

# Set up variables and dependency for one odex file
# Input variables: my_2nd_arch_prefix
# Output(modified) variables: built_odex, installed_odex, built_installed_odex
# Output(modified) variables for DEX_PREOPT_IN_DATA: data_odex

my_local_odex_in_data :=

ifneq ($(filter %.apk,$(LOCAL_INSTALLED_MODULE)),)
ifneq ($(filter $(LOCAL_MODULE),$(DEX_PREOPT_IN_DATA_LIST)),)
my_local_odex_in_data := true
endif
endif

ifeq ($(my_local_odex_in_data),true)
my_built_odex := $(call get-odex-data-path,$($(my_2nd_arch_prefix)DEX2OAT_TARGET_ARCH),$(LOCAL_BUILT_MODULE),$(patsubst $(PRODUCT_OUT)/%,%,$(LOCAL_INSTALLED_MODULE)))
else
my_built_odex := $(call get-odex-file-path,$($(my_2nd_arch_prefix)DEX2OAT_TARGET_ARCH),$(LOCAL_BUILT_MODULE))
endif

ifdef LOCAL_DEX_PREOPT_IMAGE_LOCATION
my_dex_preopt_image_location := $(LOCAL_DEX_PREOPT_IMAGE_LOCATION)
else
my_dex_preopt_image_location := $($(my_2nd_arch_prefix)DEFAULT_DEX_PREOPT_BUILT_IMAGE_LOCATION)
endif
my_dex_preopt_image_filename := $(call get-image-file-path,$($(my_2nd_arch_prefix)DEX2OAT_TARGET_ARCH),$(my_dex_preopt_image_location))
$(my_built_odex): PRIVATE_2ND_ARCH_VAR_PREFIX := $(my_2nd_arch_prefix)
$(my_built_odex): PRIVATE_DEX_LOCATION := $(patsubst $(PRODUCT_OUT)%,%,$(LOCAL_INSTALLED_MODULE))
$(my_built_odex): PRIVATE_DEX_PREOPT_IMAGE_LOCATION := $(my_dex_preopt_image_location)
$(my_built_odex) : $($(my_2nd_arch_prefix)DEXPREOPT_ONE_FILE_DEPENDENCY_BUILT_BOOT_PREOPT) \
    $(DEXPREOPT_ONE_FILE_DEPENDENCY_TOOLS) \
    $(my_dex_preopt_image_filename)

built_odex += $(my_built_odex)

ifeq ($(my_local_odex_in_data),true)
DALVIK_CACHE_PREFIX := data/dalvik-cache
PRODUCT_OUT_DALVIK_CACHE := $(PRODUCT_OUT)/$(DALVIK_CACHE_PREFIX)

my_data_odex := $(call get-odex-data-path,$($(my_2nd_arch_prefix)DEX2OAT_TARGET_ARCH),$(PRODUCT_OUT_DALVIK_CACHE)/,$(patsubst $(PRODUCT_OUT)/%,%,$(LOCAL_INSTALLED_MODULE)))

data_odex += $(my_data_odex)
built_installed_odex += $(my_built_odex):$(my_data_odex)
else
my_installed_odex := $(call get-odex-file-path,$($(my_2nd_arch_prefix)DEX2OAT_TARGET_ARCH),$(LOCAL_INSTALLED_MODULE))

installed_odex += $(my_installed_odex)
built_installed_odex += $(my_built_odex):$(my_installed_odex)
endif
