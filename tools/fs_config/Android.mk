# Copyright (C) 2008 The Android Open Source Project
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

LOCAL_PATH := $(call my-dir)

# One can override the default android_filesystem_config.h file in one of two ways:
#
# 1. The old way:
#   To Build the custom target binary for the host to generate the fs_config
#   override files. The executable is hard coded to include the
#   $(TARGET_ANDROID_FILESYSTEM_CONFIG_H) file if it exists.
#   Expectations:
#      device/<vendor>/<device>/android_filesystem_config.h
#          fills in struct fs_path_config android_device_dirs[] and
#                   struct fs_path_config android_device_files[]
#      device/<vendor>/<device>/device.mk
#          PRODUCT_PACKAGES += fs_config_dirs fs_config_files
#   If not specified, check if default one to be found
#
# 2. The new way:
#   set TARGET_FS_CONFIG_GEN to contain a list of intermediate format files
#   for generating the android_filesystem_config.h file.
#
# More information can be found in the README
ANDROID_FS_CONFIG_H := android_filesystem_config.h

ifneq ($(TARGET_ANDROID_FILESYSTEM_CONFIG_H),)
ifneq ($(TARGET_FS_CONFIG_GEN),)
$(error Cannot set TARGET_ANDROID_FILESYSTEM_CONFIG_H and TARGET_FS_CONFIG_GEN simultaneously)
endif

# One and only one file can be specified.
ifneq ($(words $(TARGET_ANDROID_FILESYSTEM_CONFIG_H)),1)
$(error Multiple fs_config files specified, \
 see "$(TARGET_ANDROID_FILESYSTEM_CONFIG_H)".)
endif

<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
ifeq ($(filter %/$(ANDROID_FS_CONFIG_H),$(TARGET_ANDROID_FILESYSTEM_CONFIG_H)),)
$(error TARGET_ANDROID_FILESYSTEM_CONFIG_H file name must be $(ANDROID_FS_CONFIG_H), \
 see "$(notdir $(TARGET_ANDROID_FILESYSTEM_CONFIG_H))".)
endif

my_fs_config_h := $(TARGET_ANDROID_FILESYSTEM_CONFIG_H)
else ifneq ($(wildcard $(TARGET_DEVICE_DIR)/$(ANDROID_FS_CONFIG_H)),)

ifneq ($(TARGET_FS_CONFIG_GEN),)
$(error Cannot provide $(TARGET_DEVICE_DIR)/$(ANDROID_FS_CONFIG_H) and set TARGET_FS_CONFIG_GEN simultaneously)
endif
my_fs_config_h := $(TARGET_DEVICE_DIR)/$(ANDROID_FS_CONFIG_H)

else
my_fs_config_h := $(LOCAL_PATH)/default/$(ANDROID_FS_CONFIG_H)
endif

##################################
include $(CLEAR_VARS)
LOCAL_SRC_FILES := fs_config_generate.c
LOCAL_MODULE := fs_config_generate_$(TARGET_DEVICE)
LOCAL_MODULE_CLASS := EXECUTABLES
LOCAL_SHARED_LIBRARIES := libcutils
LOCAL_CFLAGS := -Werror -Wno-error=\#warnings

ifneq ($(TARGET_FS_CONFIG_GEN),)
system_android_filesystem_config := system/core/include/private/android_filesystem_config.h

# Generate the "generated_oem_aid.h" file
oem := $(local-generated-sources-dir)/generated_oem_aid.h
$(oem): PRIVATE_LOCAL_PATH := $(LOCAL_PATH)
$(oem): PRIVATE_TARGET_FS_CONFIG_GEN := $(TARGET_FS_CONFIG_GEN)
$(oem): PRIVATE_ANDROID_FS_HDR := $(system_android_filesystem_config)
$(oem): PRIVATE_CUSTOM_TOOL = $(PRIVATE_LOCAL_PATH)/fs_config_generator.py oemaid --aid-header=$(PRIVATE_ANDROID_FS_HDR) $(PRIVATE_TARGET_FS_CONFIG_GEN) > $@
$(oem): $(TARGET_FS_CONFIG_GEN) $(LOCAL_PATH)/fs_config_generator.py
	$(transform-generated-source)

# Generate the fs_config header
gen := $(local-generated-sources-dir)/$(ANDROID_FS_CONFIG_H)
$(gen): PRIVATE_LOCAL_PATH := $(LOCAL_PATH)
$(gen): PRIVATE_TARGET_FS_CONFIG_GEN := $(TARGET_FS_CONFIG_GEN)
$(gen): PRIVATE_ANDROID_FS_HDR := $(system_android_filesystem_config)
$(gen): PRIVATE_CUSTOM_TOOL = $(PRIVATE_LOCAL_PATH)/fs_config_generator.py fsconfig --aid-header=$(PRIVATE_ANDROID_FS_HDR) $(PRIVATE_TARGET_FS_CONFIG_GEN) > $@
$(gen): $(TARGET_FS_CONFIG_GEN) $(system_android_filesystem_config) $(LOCAL_PATH)/fs_config_generator.py
	$(transform-generated-source)

LOCAL_GENERATED_SOURCES := $(oem) $(gen)

my_fs_config_h := $(gen)
my_gen_oem_aid := $(oem)
gen :=
oem :=
endif

LOCAL_C_INCLUDES := $(dir $(my_fs_config_h)) $(dir $(my_gen_oem_aid))

include $(BUILD_HOST_EXECUTABLE)
fs_config_generate_bin := $(LOCAL_INSTALLED_MODULE)
# List of all supported vendor, oem and odm Partitions
=======
# Use snapshots if exist
vendor_android_filesystem_config := $(strip \
  $(if $(filter-out current,$(BOARD_VNDK_VERSION)), \
    $(SOONG_VENDOR_$(BOARD_VNDK_VERSION)_SNAPSHOT_DIR)/include/$(system_android_filesystem_config)))
ifeq (,$(wildcard $(vendor_android_filesystem_config)))
vendor_android_filesystem_config := $(system_android_filesystem_config)
endif

vendor_capability_header := $(strip \
  $(if $(filter-out current,$(BOARD_VNDK_VERSION)), \
    $(SOONG_VENDOR_$(BOARD_VNDK_VERSION)_SNAPSHOT_DIR)/include/$(system_capability_header)))
ifeq (,$(wildcard $(vendor_capability_header)))
vendor_capability_header := $(system_capability_header)
endif

# List of supported vendor, oem, odm, vendor_dlkm and odm_dlkm Partitions
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
fs_config_generate_extra_partition_list := $(strip \
  $(if $(BOARD_USES_VENDORIMAGE)$(BOARD_VENDORIMAGE_FILE_SYSTEM_TYPE),vendor) \
  $(if $(BOARD_USES_OEMIMAGE)$(BOARD_OEMIMAGE_FILE_SYSTEM_TYPE),oem) \
<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
  $(if $(BOARD_USES_ODMIMAGE)$(BOARD_ODMIMAGE_FILE_SYSTEM_TYPE),odm))
=======
  $(if $(BOARD_USES_ODMIMAGE)$(BOARD_ODMIMAGE_FILE_SYSTEM_TYPE),odm) \
  $(if $(BOARD_USES_VENDOR_DLKMIMAGE)$(BOARD_VENDOR_DLKMIMAGE_FILE_SYSTEM_TYPE),vendor_dlkm) \
  $(if $(BOARD_USES_ODM_DLKMIMAGE)$(BOARD_ODM_DLKMIMAGE_FILE_SYSTEM_TYPE),odm_dlkm) \
)

##################################
# Generate the <p>/etc/fs_config_dirs binary files for each partition.
# Add fs_config_dirs to PRODUCT_PACKAGES in the device make file to enable.
include $(CLEAR_VARS)

LOCAL_MODULE := fs_config_dirs
LOCAL_LICENSE_KINDS := legacy_restricted
LOCAL_LICENSE_CONDITIONS := restricted
LOCAL_REQUIRED_MODULES := \
  fs_config_dirs_system \
  fs_config_dirs_system_ext \
  fs_config_dirs_product \
  fs_config_dirs_nonsystem
include $(BUILD_PHONY_PACKAGE)

##################################
# Generate the <p>/etc/fs_config_files binary files for each partition.
# Add fs_config_files to PRODUCT_PACKAGES in the device make file to enable.
include $(CLEAR_VARS)

LOCAL_MODULE := fs_config_files
LOCAL_LICENSE_KINDS := legacy_restricted
LOCAL_LICENSE_CONDITIONS := restricted
LOCAL_REQUIRED_MODULES := \
  fs_config_files_system \
  fs_config_files_system_ext \
  fs_config_files_product \
  fs_config_files_nonsystem
include $(BUILD_PHONY_PACKAGE)

##################################
# Generate the system_ext/etc/fs_config_dirs binary file for the target if the
# system_ext partition is generated. Add fs_config_dirs or fs_config_dirs_system_ext
# to PRODUCT_PACKAGES in the device make file to enable.
include $(CLEAR_VARS)

LOCAL_MODULE := fs_config_dirs_system_ext
LOCAL_LICENSE_KINDS := legacy_restricted
LOCAL_LICENSE_CONDITIONS := restricted
LOCAL_REQUIRED_MODULES := $(if $(BOARD_USES_SYSTEM_EXTIMAGE)$(BOARD_SYSTEM_EXTIMAGE_FILE_SYSTEM_TYPE),_fs_config_dirs_system_ext)
include $(BUILD_PHONY_PACKAGE)

##################################
# Generate the system_ext/etc/fs_config_files binary file for the target if the
# system_ext partition is generated. Add fs_config_files or fs_config_files_system_ext
# to PRODUCT_PACKAGES in the device make file to enable.
include $(CLEAR_VARS)

LOCAL_MODULE := fs_config_files_system_ext
LOCAL_LICENSE_KINDS := legacy_restricted
LOCAL_LICENSE_CONDITIONS := restricted
LOCAL_REQUIRED_MODULES := $(if $(BOARD_USES_SYSTEM_EXTIMAGE)$(BOARD_SYSTEM_EXTIMAGE_FILE_SYSTEM_TYPE),_fs_config_files_system_ext)
include $(BUILD_PHONY_PACKAGE)

##################################
# Generate the product/etc/fs_config_dirs binary file for the target if the
# product partition is generated. Add fs_config_dirs or fs_config_dirs_product
# to PRODUCT_PACKAGES in the device make file to enable.
include $(CLEAR_VARS)

LOCAL_MODULE := fs_config_dirs_product
LOCAL_LICENSE_KINDS := legacy_restricted
LOCAL_LICENSE_CONDITIONS := restricted
LOCAL_REQUIRED_MODULES := $(if $(BOARD_USES_PRODUCTIMAGE)$(BOARD_PRODUCTIMAGE_FILE_SYSTEM_TYPE),_fs_config_dirs_product)
include $(BUILD_PHONY_PACKAGE)

##################################
# Generate the product/etc/fs_config_files binary file for the target if the
# product partition is generated. Add fs_config_files or fs_config_files_product
# to PRODUCT_PACKAGES in the device make file to enable.
include $(CLEAR_VARS)

LOCAL_MODULE := fs_config_files_product
LOCAL_LICENSE_KINDS := legacy_restricted
LOCAL_LICENSE_CONDITIONS := restricted
LOCAL_REQUIRED_MODULES := $(if $(BOARD_USES_PRODUCTIMAGE)$(BOARD_PRODUCTIMAGE_FILE_SYSTEM_TYPE),_fs_config_files_product)
include $(BUILD_PHONY_PACKAGE)

##################################
# Generate the <p>/etc/fs_config_dirs binary files for all enabled partitions
# excluding /system, /system_ext and /product. Add fs_config_dirs_nonsystem to
# PRODUCT_PACKAGES in the device make file to enable.
include $(CLEAR_VARS)

LOCAL_MODULE := fs_config_dirs_nonsystem
LOCAL_LICENSE_KINDS := legacy_restricted
LOCAL_LICENSE_CONDITIONS := restricted
LOCAL_REQUIRED_MODULES := $(foreach t,$(fs_config_generate_extra_partition_list),_fs_config_dirs_$(t))
include $(BUILD_PHONY_PACKAGE)

##################################
# Generate the <p>/etc/fs_config_files binary files for all enabled partitions
# excluding /system, /system_ext and /product. Add fs_config_files_nonsystem to
# PRODUCT_PACKAGES in the device make file to enable.
include $(CLEAR_VARS)

LOCAL_MODULE := fs_config_files_nonsystem
LOCAL_LICENSE_KINDS := legacy_restricted
LOCAL_LICENSE_CONDITIONS := restricted
LOCAL_REQUIRED_MODULES := $(foreach t,$(fs_config_generate_extra_partition_list),_fs_config_files_$(t))
include $(BUILD_PHONY_PACKAGE)
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)

##################################
# Generate the system/etc/fs_config_dirs binary file for the target
# Add fs_config_dirs to PRODUCT_PACKAGES in the device make file to enable
include $(CLEAR_VARS)

<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
LOCAL_MODULE := fs_config_dirs
=======
LOCAL_MODULE := fs_config_dirs_system
LOCAL_LICENSE_KINDS := legacy_restricted
LOCAL_LICENSE_CONDITIONS := restricted
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
LOCAL_MODULE_CLASS := ETC
LOCAL_REQUIRED_MODULES := $(foreach t,$(fs_config_generate_extra_partition_list),$(LOCAL_MODULE)_$(t))
include $(BUILD_SYSTEM)/base_rules.mk
$(LOCAL_BUILT_MODULE): $(fs_config_generate_bin)
	@mkdir -p $(dir $@)
	$< -D $(if $(fs_config_generate_extra_partition_list), \
	   -P '$(subst $(space),$(comma),$(addprefix -,$(fs_config_generate_extra_partition_list)))') \
	   -o $@

##################################
# Generate the system/etc/fs_config_files binary file for the target
# Add fs_config_files to PRODUCT_PACKAGES in the device make file to enable
include $(CLEAR_VARS)

<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
LOCAL_MODULE := fs_config_files
=======
LOCAL_MODULE := fs_config_files_system
LOCAL_LICENSE_KINDS := legacy_restricted
LOCAL_LICENSE_CONDITIONS := restricted
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
LOCAL_MODULE_CLASS := ETC
LOCAL_REQUIRED_MODULES := $(foreach t,$(fs_config_generate_extra_partition_list),$(LOCAL_MODULE)_$(t))
include $(BUILD_SYSTEM)/base_rules.mk
$(LOCAL_BUILT_MODULE): $(fs_config_generate_bin)
	@mkdir -p $(dir $@)
	$< -F $(if $(fs_config_generate_extra_partition_list), \
	   -P '$(subst $(space),$(comma),$(addprefix -,$(fs_config_generate_extra_partition_list)))') \
	   -o $@

ifneq ($(filter vendor,$(fs_config_generate_extra_partition_list)),)
##################################
# Generate the vendor/etc/fs_config_dirs binary file for the target
# Add fs_config_dirs or fs_config_dirs_vendor to PRODUCT_PACKAGES in
# the device make file to enable.
include $(CLEAR_VARS)

<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
LOCAL_MODULE := fs_config_dirs_vendor
=======
LOCAL_MODULE := _fs_config_dirs_vendor
LOCAL_LICENSE_KINDS := legacy_restricted
LOCAL_LICENSE_CONDITIONS := restricted
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
LOCAL_MODULE_CLASS := ETC
LOCAL_INSTALLED_MODULE_STEM := fs_config_dirs
LOCAL_MODULE_PATH := $(TARGET_OUT_VENDOR)/etc
include $(BUILD_SYSTEM)/base_rules.mk
<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
$(LOCAL_BUILT_MODULE): $(fs_config_generate_bin)
=======
$(LOCAL_BUILT_MODULE): PRIVATE_ANDROID_FS_HDR := $(vendor_android_filesystem_config)
$(LOCAL_BUILT_MODULE): PRIVATE_ANDROID_CAP_HDR := $(vendor_capability_header)
$(LOCAL_BUILT_MODULE): PRIVATE_TARGET_FS_CONFIG_GEN := $(TARGET_FS_CONFIG_GEN)
$(LOCAL_BUILT_MODULE): $(LOCAL_PATH)/fs_config_generator.py $(TARGET_FS_CONFIG_GEN) $(vendor_android_filesystem_config) $(vendor_capability_header)
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
	@mkdir -p $(dir $@)
	$< -D -P vendor -o $@

##################################
# Generate the vendor/etc/fs_config_files binary file for the target
# Add fs_config_files or fs_config_files_vendor to PRODUCT_PACKAGES in
# the device make file to enable
include $(CLEAR_VARS)

<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
LOCAL_MODULE := fs_config_files_vendor
=======
LOCAL_MODULE := _fs_config_files_vendor
LOCAL_LICENSE_KINDS := legacy_restricted
LOCAL_LICENSE_CONDITIONS := restricted
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
LOCAL_MODULE_CLASS := ETC
LOCAL_INSTALLED_MODULE_STEM := fs_config_files
LOCAL_MODULE_PATH := $(TARGET_OUT_VENDOR)/etc
include $(BUILD_SYSTEM)/base_rules.mk
<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
$(LOCAL_BUILT_MODULE): $(fs_config_generate_bin)
=======
$(LOCAL_BUILT_MODULE): PRIVATE_ANDROID_FS_HDR := $(vendor_android_filesystem_config)
$(LOCAL_BUILT_MODULE): PRIVATE_ANDROID_CAP_HDR := $(vendor_capability_header)
$(LOCAL_BUILT_MODULE): PRIVATE_TARGET_FS_CONFIG_GEN := $(TARGET_FS_CONFIG_GEN)
$(LOCAL_BUILT_MODULE): $(LOCAL_PATH)/fs_config_generator.py $(TARGET_FS_CONFIG_GEN) $(vendor_android_filesystem_config) $(vendor_capability_header)
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
	@mkdir -p $(dir $@)
	$< -F -P vendor -o $@

endif

ifneq ($(filter oem,$(fs_config_generate_extra_partition_list)),)
##################################
# Generate the oem/etc/fs_config_dirs binary file for the target
# Add fs_config_dirs or fs_config_dirs_oem to PRODUCT_PACKAGES in
# the device make file to enable
include $(CLEAR_VARS)

<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
LOCAL_MODULE := fs_config_dirs_oem
=======
LOCAL_MODULE := _fs_config_dirs_oem
LOCAL_LICENSE_KINDS := legacy_restricted
LOCAL_LICENSE_CONDITIONS := restricted
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
LOCAL_MODULE_CLASS := ETC
LOCAL_INSTALLED_MODULE_STEM := fs_config_dirs
LOCAL_MODULE_PATH := $(TARGET_OUT_OEM)/etc
include $(BUILD_SYSTEM)/base_rules.mk
$(LOCAL_BUILT_MODULE): $(fs_config_generate_bin)
	@mkdir -p $(dir $@)
	$< -D -P oem -o $@

##################################
# Generate the oem/etc/fs_config_files binary file for the target
# Add fs_config_files or fs_config_files_oem to PRODUCT_PACKAGES in
# the device make file to enable
include $(CLEAR_VARS)

<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
LOCAL_MODULE := fs_config_files_oem
=======
LOCAL_MODULE := _fs_config_files_oem
LOCAL_LICENSE_KINDS := legacy_restricted
LOCAL_LICENSE_CONDITIONS := restricted
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
LOCAL_MODULE_CLASS := ETC
LOCAL_INSTALLED_MODULE_STEM := fs_config_files
LOCAL_MODULE_PATH := $(TARGET_OUT_OEM)/etc
include $(BUILD_SYSTEM)/base_rules.mk
$(LOCAL_BUILT_MODULE): $(fs_config_generate_bin)
	@mkdir -p $(dir $@)
	$< -F -P oem -o $@

endif

ifneq ($(filter odm,$(fs_config_generate_extra_partition_list)),)
##################################
# Generate the odm/etc/fs_config_dirs binary file for the target
# Add fs_config_dirs or fs_config_dirs_odm to PRODUCT_PACKAGES in
# the device make file to enable
include $(CLEAR_VARS)

<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
LOCAL_MODULE := fs_config_dirs_odm
=======
LOCAL_MODULE := _fs_config_dirs_odm
LOCAL_LICENSE_KINDS := legacy_restricted
LOCAL_LICENSE_CONDITIONS := restricted
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
LOCAL_MODULE_CLASS := ETC
LOCAL_INSTALLED_MODULE_STEM := fs_config_dirs
LOCAL_MODULE_PATH := $(TARGET_OUT_ODM)/etc
include $(BUILD_SYSTEM)/base_rules.mk
<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
$(LOCAL_BUILT_MODULE): $(fs_config_generate_bin)
=======
$(LOCAL_BUILT_MODULE): PRIVATE_ANDROID_FS_HDR := $(vendor_android_filesystem_config)
$(LOCAL_BUILT_MODULE): PRIVATE_ANDROID_CAP_HDR := $(vendor_capability_header)
$(LOCAL_BUILT_MODULE): PRIVATE_TARGET_FS_CONFIG_GEN := $(TARGET_FS_CONFIG_GEN)
$(LOCAL_BUILT_MODULE): $(LOCAL_PATH)/fs_config_generator.py $(TARGET_FS_CONFIG_GEN) $(vendor_android_filesystem_config) $(vendor_capability_header)
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
	@mkdir -p $(dir $@)
	$< -D -P odm -o $@

##################################
# Generate the odm/etc/fs_config_files binary file for the target
# Add fs_config_files of fs_config_files_odm to PRODUCT_PACKAGES in
# the device make file to enable
include $(CLEAR_VARS)

<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
LOCAL_MODULE := fs_config_files_odm
=======
LOCAL_MODULE := _fs_config_files_odm
LOCAL_LICENSE_KINDS := legacy_restricted
LOCAL_LICENSE_CONDITIONS := restricted
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
LOCAL_MODULE_CLASS := ETC
LOCAL_INSTALLED_MODULE_STEM := fs_config_files
LOCAL_MODULE_PATH := $(TARGET_OUT_ODM)/etc
include $(BUILD_SYSTEM)/base_rules.mk
<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
$(LOCAL_BUILT_MODULE): $(fs_config_generate_bin)
=======
$(LOCAL_BUILT_MODULE): PRIVATE_ANDROID_FS_HDR := $(vendor_android_filesystem_config)
$(LOCAL_BUILT_MODULE): PRIVATE_ANDROID_CAP_HDR := $(vendor_capability_header)
$(LOCAL_BUILT_MODULE): PRIVATE_TARGET_FS_CONFIG_GEN := $(TARGET_FS_CONFIG_GEN)
$(LOCAL_BUILT_MODULE): $(LOCAL_PATH)/fs_config_generator.py $(TARGET_FS_CONFIG_GEN) $(vendor_android_filesystem_config) $(vendor_capability_header)
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
	@mkdir -p $(dir $@)
	$< -F -P odm -o $@

endif

<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
# The newer passwd/group targets are only generated if you
# use the new TARGET_FS_CONFIG_GEN method.
ifneq ($(TARGET_FS_CONFIG_GEN),)
=======
ifneq ($(filter vendor_dlkm,$(fs_config_generate_extra_partition_list)),)
##################################
# Generate the vendor_dlkm/etc/fs_config_dirs binary file for the target
# Add fs_config_dirs or fs_config_dirs_nonsystem to PRODUCT_PACKAGES in
# the device make file to enable
include $(CLEAR_VARS)

LOCAL_MODULE := _fs_config_dirs_vendor_dlkm
LOCAL_LICENSE_KINDS := legacy_restricted
LOCAL_LICENSE_CONDITIONS := restricted
LOCAL_MODULE_CLASS := ETC
LOCAL_INSTALLED_MODULE_STEM := fs_config_dirs
LOCAL_MODULE_PATH := $(TARGET_OUT_VENDOR_DLKM)/etc
include $(BUILD_SYSTEM)/base_rules.mk
$(LOCAL_BUILT_MODULE): PRIVATE_ANDROID_FS_HDR := $(vendor_android_filesystem_config)
$(LOCAL_BUILT_MODULE): PRIVATE_ANDROID_CAP_HDR := $(vendor_capability_header)
$(LOCAL_BUILT_MODULE): PRIVATE_TARGET_FS_CONFIG_GEN := $(TARGET_FS_CONFIG_GEN)
$(LOCAL_BUILT_MODULE): $(LOCAL_PATH)/fs_config_generator.py $(TARGET_FS_CONFIG_GEN) $(vendor_android_filesystem_config) $(vendor_capability_header)
	@mkdir -p $(dir $@)
	$< fsconfig \
	   --aid-header $(PRIVATE_ANDROID_FS_HDR) \
	   --capability-header $(PRIVATE_ANDROID_CAP_HDR) \
	   --partition vendor_dlkm \
	   --dirs \
	   --out_file $@ \
	   $(or $(PRIVATE_TARGET_FS_CONFIG_GEN),/dev/null)
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)

##################################
# Build the oemaid header library when fs config files are present.
# Intentionally break build if you require generated AIDs
# header file, but are not using any fs config files.
include $(CLEAR_VARS)
<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
LOCAL_MODULE := oemaids_headers
LOCAL_EXPORT_C_INCLUDE_DIRS := $(dir $(my_gen_oem_aid))
LOCAL_EXPORT_C_INCLUDE_DEPS := $(my_gen_oem_aid)
include $(BUILD_HEADER_LIBRARY)
=======

LOCAL_MODULE := _fs_config_files_vendor_dlkm
LOCAL_LICENSE_KINDS := legacy_restricted
LOCAL_LICENSE_CONDITIONS := restricted
LOCAL_MODULE_CLASS := ETC
LOCAL_INSTALLED_MODULE_STEM := fs_config_files
LOCAL_MODULE_PATH := $(TARGET_OUT_VENDOR_DLKM)/etc
include $(BUILD_SYSTEM)/base_rules.mk
$(LOCAL_BUILT_MODULE): PRIVATE_ANDROID_FS_HDR := $(vendor_android_filesystem_config)
$(LOCAL_BUILT_MODULE): PRIVATE_ANDROID_CAP_HDR := $(vendor_capability_header)
$(LOCAL_BUILT_MODULE): PRIVATE_TARGET_FS_CONFIG_GEN := $(TARGET_FS_CONFIG_GEN)
$(LOCAL_BUILT_MODULE): $(LOCAL_PATH)/fs_config_generator.py $(TARGET_FS_CONFIG_GEN) $(vendor_android_filesystem_config) $(vendor_capability_header)
	@mkdir -p $(dir $@)
	$< fsconfig \
	   --aid-header $(PRIVATE_ANDROID_FS_HDR) \
	   --capability-header $(PRIVATE_ANDROID_CAP_HDR) \
	   --partition vendor_dlkm \
	   --files \
	   --out_file $@ \
	   $(or $(PRIVATE_TARGET_FS_CONFIG_GEN),/dev/null)

endif

ifneq ($(filter odm_dlkm,$(fs_config_generate_extra_partition_list)),)
##################################
# Generate the odm_dlkm/etc/fs_config_dirs binary file for the target
# Add fs_config_dirs or fs_config_dirs_nonsystem to PRODUCT_PACKAGES
# in the device make file to enable
include $(CLEAR_VARS)

LOCAL_MODULE := _fs_config_dirs_odm_dlkm
LOCAL_LICENSE_KINDS := legacy_restricted
LOCAL_LICENSE_CONDITIONS := restricted
LOCAL_MODULE_CLASS := ETC
LOCAL_INSTALLED_MODULE_STEM := fs_config_dirs
LOCAL_MODULE_PATH := $(TARGET_OUT_ODM_DLKM)/etc
include $(BUILD_SYSTEM)/base_rules.mk
$(LOCAL_BUILT_MODULE): PRIVATE_ANDROID_FS_HDR := $(vendor_android_filesystem_config)
$(LOCAL_BUILT_MODULE): PRIVATE_ANDROID_CAP_HDR := $(vendor_capability_header)
$(LOCAL_BUILT_MODULE): PRIVATE_TARGET_FS_CONFIG_GEN := $(TARGET_FS_CONFIG_GEN)
$(LOCAL_BUILT_MODULE): $(LOCAL_PATH)/fs_config_generator.py $(TARGET_FS_CONFIG_GEN) $(vendor_android_filesystem_config) $(vendor_capability_header)
	@mkdir -p $(dir $@)
	$< fsconfig \
	   --aid-header $(PRIVATE_ANDROID_FS_HDR) \
	   --capability-header $(PRIVATE_ANDROID_CAP_HDR) \
	   --partition odm_dlkm \
	   --dirs \
	   --out_file $@ \
	   $(or $(PRIVATE_TARGET_FS_CONFIG_GEN),/dev/null)
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)

##################################
# Generate the vendor/etc/passwd text file for the target
# This file may be empty if no AIDs are defined in
# TARGET_FS_CONFIG_GEN files.
include $(CLEAR_VARS)

<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
LOCAL_MODULE := passwd
=======
LOCAL_MODULE := _fs_config_files_odm_dlkm
LOCAL_LICENSE_KINDS := legacy_restricted
LOCAL_LICENSE_CONDITIONS := restricted
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
LOCAL_MODULE_CLASS := ETC
LOCAL_VENDOR_MODULE := true

include $(BUILD_SYSTEM)/base_rules.mk
<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)

$(LOCAL_BUILT_MODULE): PRIVATE_LOCAL_PATH := $(LOCAL_PATH)
=======
$(LOCAL_BUILT_MODULE): PRIVATE_ANDROID_FS_HDR := $(vendor_android_filesystem_config)
$(LOCAL_BUILT_MODULE): PRIVATE_ANDROID_CAP_HDR := $(vendor_capability_header)
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
$(LOCAL_BUILT_MODULE): PRIVATE_TARGET_FS_CONFIG_GEN := $(TARGET_FS_CONFIG_GEN)
<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
=======
$(LOCAL_BUILT_MODULE): $(LOCAL_PATH)/fs_config_generator.py $(TARGET_FS_CONFIG_GEN) $(vendor_android_filesystem_config) $(vendor_capability_header)
	@mkdir -p $(dir $@)
	$< fsconfig \
	   --aid-header $(PRIVATE_ANDROID_FS_HDR) \
	   --capability-header $(PRIVATE_ANDROID_CAP_HDR) \
	   --partition odm_dlkm \
	   --files \
	   --out_file $@ \
	   $(or $(PRIVATE_TARGET_FS_CONFIG_GEN),/dev/null)

endif

ifneq ($(BOARD_USES_PRODUCTIMAGE)$(BOARD_PRODUCTIMAGE_FILE_SYSTEM_TYPE),)
##################################
# Generate the product/etc/fs_config_dirs binary file for the target
# Add fs_config_dirs or fs_config_dirs_product to PRODUCT_PACKAGES in
# the device make file to enable
include $(CLEAR_VARS)

LOCAL_MODULE := _fs_config_dirs_product
LOCAL_LICENSE_KINDS := legacy_restricted
LOCAL_LICENSE_CONDITIONS := restricted
LOCAL_MODULE_CLASS := ETC
LOCAL_INSTALLED_MODULE_STEM := fs_config_dirs
LOCAL_MODULE_PATH := $(TARGET_OUT_PRODUCT)/etc
include $(BUILD_SYSTEM)/base_rules.mk
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
$(LOCAL_BUILT_MODULE): PRIVATE_ANDROID_FS_HDR := $(system_android_filesystem_config)
$(LOCAL_BUILT_MODULE): $(LOCAL_PATH)/fs_config_generator.py $(TARGET_FS_CONFIG_GEN) $(system_android_filesystem_config)
	@mkdir -p $(dir $@)
	$(hide) $< passwd --required-prefix=vendor_ --aid-header=$(PRIVATE_ANDROID_FS_HDR) $(PRIVATE_TARGET_FS_CONFIG_GEN) > $@

##################################
# Generate the vendor/etc/group text file for the target
# This file may be empty if no AIDs are defined in
# TARGET_FS_CONFIG_GEN files.
include $(CLEAR_VARS)

<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
LOCAL_MODULE := group
=======
LOCAL_MODULE := _fs_config_files_product
LOCAL_LICENSE_KINDS := legacy_restricted
LOCAL_LICENSE_CONDITIONS := restricted
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
LOCAL_MODULE_CLASS := ETC
LOCAL_VENDOR_MODULE := true

include $(BUILD_SYSTEM)/base_rules.mk

$(LOCAL_BUILT_MODULE): PRIVATE_LOCAL_PATH := $(LOCAL_PATH)
$(LOCAL_BUILT_MODULE): PRIVATE_TARGET_FS_CONFIG_GEN := $(TARGET_FS_CONFIG_GEN)
<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
=======
$(LOCAL_BUILT_MODULE): $(LOCAL_PATH)/fs_config_generator.py $(TARGET_FS_CONFIG_GEN) $(system_android_filesystem_config) $(system_capability_header)
	@mkdir -p $(dir $@)
	$< fsconfig \
	   --aid-header $(PRIVATE_ANDROID_FS_HDR) \
	   --capability-header $(PRIVATE_ANDROID_CAP_HDR) \
	   --partition product \
	   --files \
	   --out_file $@ \
	   $(or $(PRIVATE_TARGET_FS_CONFIG_GEN),/dev/null)
endif

ifneq ($(BOARD_USES_SYSTEM_EXTIMAGE)$(BOARD_SYSTEM_EXTIMAGE_FILE_SYSTEM_TYPE),)
##################################
# Generate the system_ext/etc/fs_config_dirs binary file for the target
# Add fs_config_dirs or fs_config_dirs_system_ext to PRODUCT_PACKAGES in
# the device make file to enable
include $(CLEAR_VARS)

LOCAL_MODULE := _fs_config_dirs_system_ext
LOCAL_LICENSE_KINDS := legacy_restricted
LOCAL_LICENSE_CONDITIONS := restricted
LOCAL_MODULE_CLASS := ETC
LOCAL_INSTALLED_MODULE_STEM := fs_config_dirs
LOCAL_MODULE_PATH := $(TARGET_OUT_SYSTEM_EXT)/etc
include $(BUILD_SYSTEM)/base_rules.mk
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
$(LOCAL_BUILT_MODULE): PRIVATE_ANDROID_FS_HDR := $(system_android_filesystem_config)
$(LOCAL_BUILT_MODULE): $(LOCAL_PATH)/fs_config_generator.py $(TARGET_FS_CONFIG_GEN) $(system_android_filesystem_config)
	@mkdir -p $(dir $@)
<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
	$(hide) $< group --required-prefix=vendor_ --aid-header=$(PRIVATE_ANDROID_FS_HDR) $(PRIVATE_TARGET_FS_CONFIG_GEN) > $@
=======
	$< fsconfig \
	   --aid-header $(PRIVATE_ANDROID_FS_HDR) \
	   --capability-header $(PRIVATE_ANDROID_CAP_HDR) \
	   --partition system_ext \
	   --dirs \
	   --out_file $@ \
	   $(or $(PRIVATE_TARGET_FS_CONFIG_GEN),/dev/null)

##################################
# Generate the system_ext/etc/fs_config_files binary file for the target
# Add fs_config_files or fs_config_files_system_ext to PRODUCT_PACKAGES in
# the device make file to enable
include $(CLEAR_VARS)

LOCAL_MODULE := _fs_config_files_system_ext
LOCAL_LICENSE_KINDS := legacy_restricted
LOCAL_LICENSE_CONDITIONS := restricted
LOCAL_MODULE_CLASS := ETC
LOCAL_INSTALLED_MODULE_STEM := fs_config_files
LOCAL_MODULE_PATH := $(TARGET_OUT_SYSTEM_EXT)/etc
include $(BUILD_SYSTEM)/base_rules.mk
$(LOCAL_BUILT_MODULE): PRIVATE_ANDROID_FS_HDR := $(system_android_filesystem_config)
$(LOCAL_BUILT_MODULE): PRIVATE_ANDROID_CAP_HDR := $(system_capability_header)
$(LOCAL_BUILT_MODULE): PRIVATE_TARGET_FS_CONFIG_GEN := $(TARGET_FS_CONFIG_GEN)
$(LOCAL_BUILT_MODULE): $(LOCAL_PATH)/fs_config_generator.py $(TARGET_FS_CONFIG_GEN) $(system_android_filesystem_config) $(system_capability_header)
	@mkdir -p $(dir $@)
	$< fsconfig \
	   --aid-header $(PRIVATE_ANDROID_FS_HDR) \
	   --capability-header $(PRIVATE_ANDROID_CAP_HDR) \
	   --partition system_ext \
	   --files \
	   --out_file $@ \
	   $(or $(PRIVATE_TARGET_FS_CONFIG_GEN),/dev/null)
endif
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)

system_android_filesystem_config :=
endif

ANDROID_FS_CONFIG_H :=
my_fs_config_h :=
fs_config_generate_bin :=
my_gen_oem_aid :=
fs_config_generate_extra_partition_list :=
