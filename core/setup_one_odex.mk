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

my_built_odex := $(call get-odex-file-path,$($(my_2nd_arch_prefix)DEX2OAT_TARGET_ARCH),$(LOCAL_BUILT_MODULE))
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

# By default pass the special class loader context to skip the classpath and collision check.
# TODO: We should modify the build system to pass used libraries properly. b/26880306
dex2oat_class_loader_context := \&

# For the system server, we know all modules. Detect these modules and construct
# an actual classloader context.
ifneq (,$(filter $(PRODUCT_SYSTEM_SERVER_JARS),$(LOCAL_MODULE)))
$(warning Detected system server jar: $(LOCAL_MODULE) in $(PRODUCT_SYSTEM_SERVER_JARS))
classpath_modules :=
classpath_modules_go_on := true
define check-and-save-module
  ifeq ($$(classpath_modules_go_on),true)
    ifneq ($(1),$(LOCAL_MODULE))
      classpath_modules := $$(classpath_modules) $(1)
    else
      classpath_modules_go_on := false
    endif
  endif

endef
$(eval $(foreach module,$(PRODUCT_SYSTEM_SERVER_JARS),$(call check-and-save-module,$(module))))
check-and-save-module :=

# Transform the dependencies to their classes.dex files. This is computing intermediates
# for other modules to find the unstripped dex files.
classpath_dex := $(foreach module,$(classpath_modules),$(call intermediates-dir-for,$(LOCAL_MODULE_CLASS),$(module),$(call def-host-aux-target),COMMON,,)/classes.dex)

# Add these as build dependencies.
define add_dex_dependency
$(1) : $(2)

endef
$(eval $(foreach dex,$(classpath_dex),$(call add_dex_dependency,$(my_built_odex),$(dex))))
add_dex_dependency :=

# Construct a classloader context.
#
# The paths need to be absolute.
classpath_list := $(call normalize-path-list,$(classpath_dex))
dex2oat_class_loader_context := "PCL[$(classpath_list)]"
endif

$(my_built_odex): PRIVATE_DEX2OAT_CLASS_LOADER_CONTEXT := $(dex2oat_class_loader_context)

my_installed_odex := $(call get-odex-installed-file-path,$($(my_2nd_arch_prefix)DEX2OAT_TARGET_ARCH),$(LOCAL_INSTALLED_MODULE))

my_built_vdex := $(patsubst %.odex,%.vdex,$(my_built_odex))
my_installed_vdex := $(patsubst %.odex,%.vdex,$(my_installed_odex))
my_installed_art := $(patsubst %.odex,%.art,$(my_installed_odex))

ifndef LOCAL_DEX_PREOPT_APP_IMAGE
# Local override not defined, use the global one.
ifeq (true,$(WITH_DEX_PREOPT_APP_IMAGE))
  LOCAL_DEX_PREOPT_APP_IMAGE := true
endif
endif

ifeq (true,$(LOCAL_DEX_PREOPT_APP_IMAGE))
my_built_art := $(patsubst %.odex,%.art,$(my_built_odex))
$(my_built_odex): PRIVATE_ART_FILE_PREOPT_FLAGS := --app-image-file=$(my_built_art) \
    --image-format=lz4
$(eval $(call copy-one-file,$(my_built_art),$(my_installed_art)))
built_art += $(my_built_art)
installed_art += $(my_installed_art)
built_installed_art += $(my_built_art):$(my_installed_art)
endif

$(eval $(call copy-one-file,$(my_built_odex),$(my_installed_odex)))
$(eval $(call copy-one-file,$(my_built_vdex),$(my_installed_vdex)))

built_odex += $(my_built_odex)
built_vdex += $(my_built_vdex)

installed_odex += $(my_installed_odex)
installed_vdex += $(my_installed_vdex)

built_installed_odex += $(my_built_odex):$(my_installed_odex)
built_installed_vdex += $(my_built_vdex):$(my_installed_vdex)
