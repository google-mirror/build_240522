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

# One can override the default android_filesystem_config.h file by using TARGET_FS_CONFIG_GEN.
#   Set TARGET_FS_CONFIG_GEN to contain a list of intermediate format files
#   for generating the android_filesystem_config.h file.
#
# More information can be found in the README

ifneq ($(wildcard $(TARGET_DEVICE_DIR)/android_filesystem_config.h),)
$(error Using $(TARGET_DEVICE_DIR)/android_filesystem_config.h is deprecated, please use TARGET_FS_CONFIG_GEN instead)
endif

system_android_filesystem_config := system/core/libcutils/include/private/android_filesystem_config.h
system_capability_header := bionic/libc/kernel/uapi/linux/capability.h

# List of supported vendor, oem, odm, vendor_dlkm, odm_dlkm, product and system_ext Partitions
fs_config_generate_extra_partition_list := $(strip \
  $(if $(BOARD_USES_VENDORIMAGE)$(BOARD_VENDORIMAGE_FILE_SYSTEM_TYPE),vendor) \
  $(if $(BOARD_USES_OEMIMAGE)$(BOARD_OEMIMAGE_FILE_SYSTEM_TYPE),oem) \
  $(if $(BOARD_USES_ODMIMAGE)$(BOARD_ODMIMAGE_FILE_SYSTEM_TYPE),odm) \
  $(if $(BOARD_USES_VENDOR_DLKMIMAGE)$(BOARD_VENDOR_DLKMIMAGE_FILE_SYSTEM_TYPE),vendor_dlkm) \
  $(if $(BOARD_USES_ODM_DLKMIMAGE)$(BOARD_ODM_DLKMIMAGE_FILE_SYSTEM_TYPE),odm_dlkm) \
  $(if $(BOARD_USES_PRODUCTIMAGE)$(BOARD_PRODUCTIMAGE_FILE_SYSTEM_TYPE),product) \
  $(if $(BOARD_USES_SYSTEM_EXTIMAGE)$(BOARD_SYSTEM_EXTIMAGE_FILE_SYSTEM_TYPE),system_ext) \
)

##################################
# Generate the <p>/etc/fs_config_dirs binary files for each partition.
# Add fs_config_dirs to PRODUCT_PACKAGES in the device make file to enable.
include $(CLEAR_VARS)

LOCAL_MODULE := fs_config_dirs
LOCAL_REQUIRED_MODULES := \
	fs_config_dirs_system \
	$(foreach t,$(fs_config_generate_extra_partition_list),_$(LOCAL_MODULE)_$(t))
include $(BUILD_PHONY_PACKAGE)


##################################
# Generate the <p>/etc/fs_config_files binary files for each partition.
# Add fs_config_files to PRODUCT_PACKAGES in the device make file to enable.
include $(CLEAR_VARS)

LOCAL_MODULE := fs_config_files
LOCAL_REQUIRED_MODULES := \
  fs_config_files_system \
  $(foreach t,$(fs_config_generate_extra_partition_list),_$(LOCAL_MODULE)_$(t))
include $(BUILD_PHONY_PACKAGE)

##################################
# Generate the <p>/etc/fs_config_dirs binary files for all enabled partitions
# excluding /system. Add fs_config_dirs_nonsystem to PRODUCT_PACKAGES in the
# device make file to enable.
include $(CLEAR_VARS)

LOCAL_MODULE := fs_config_dirs_nonsystem
LOCAL_REQUIRED_MODULES := $(foreach t,$(fs_config_generate_extra_partition_list),_fs_config_dirs_$(t))
include $(BUILD_PHONY_PACKAGE)

##################################
# Generate the <p>/etc/fs_config_files binary files for all enabled partitions
# excluding /system. Add fs_config_files_nonsystem to PRODUCT_PACKAGES in the
# device make file to enable.
include $(CLEAR_VARS)

LOCAL_MODULE := fs_config_files_nonsystem
LOCAL_REQUIRED_MODULES := $(foreach t,$(fs_config_generate_extra_partition_list),_fs_config_files_$(t))
include $(BUILD_PHONY_PACKAGE)

##################################
# Generate the system/etc/fs_config_dirs binary file for the target
# Add fs_config_dirs or fs_config_dirs_system to PRODUCT_PACKAGES in
# the device make file to enable
include $(CLEAR_VARS)

LOCAL_MODULE := fs_config_dirs_system
LOCAL_MODULE_CLASS := ETC
LOCAL_INSTALLED_MODULE_STEM := fs_config_dirs
include $(BUILD_SYSTEM)/base_rules.mk
$(LOCAL_BUILT_MODULE): PRIVATE_ANDROID_FS_HDR := $(system_android_filesystem_config)
$(LOCAL_BUILT_MODULE): PRIVATE_ANDROID_CAP_HDR := $(system_capability_header)
$(LOCAL_BUILT_MODULE): PRIVATE_PARTITION_LIST := $(fs_config_generate_extra_partition_list)
$(LOCAL_BUILT_MODULE): PRIVATE_TARGET_FS_CONFIG_GEN := $(TARGET_FS_CONFIG_GEN)
$(LOCAL_BUILT_MODULE): $(LOCAL_PATH)/fs_config_generator.py $(TARGET_FS_CONFIG_GEN) $(system_android_filesystem_config) $(system_capability_header)
	@mkdir -p $(dir $@)
	$< fsconfig \
	   --aid-header $(PRIVATE_ANDROID_FS_HDR) \
	   --capability-header $(PRIVATE_ANDROID_CAP_HDR) \
	   --partition system \
	   --all-partitions "$(subst $(space),$(comma),$(PRIVATE_PARTITION_LIST))" \
	   --dirs \
	   --out_file $@ \
	   $(or $(PRIVATE_TARGET_FS_CONFIG_GEN),/dev/null)

##################################
# Generate the system/etc/fs_config_files binary file for the target
# Add fs_config_files or fs_config_files_system to PRODUCT_PACKAGES in
# the device make file to enable
include $(CLEAR_VARS)

LOCAL_MODULE := fs_config_files_system
LOCAL_MODULE_CLASS := ETC
LOCAL_INSTALLED_MODULE_STEM := fs_config_files
include $(BUILD_SYSTEM)/base_rules.mk
$(LOCAL_BUILT_MODULE): PRIVATE_ANDROID_FS_HDR := $(system_android_filesystem_config)
$(LOCAL_BUILT_MODULE): PRIVATE_ANDROID_CAP_HDR := $(system_capability_header)
$(LOCAL_BUILT_MODULE): PRIVATE_PARTITION_LIST := $(fs_config_generate_extra_partition_list)
$(LOCAL_BUILT_MODULE): PRIVATE_TARGET_FS_CONFIG_GEN := $(TARGET_FS_CONFIG_GEN)
$(LOCAL_BUILT_MODULE): $(LOCAL_PATH)/fs_config_generator.py $(TARGET_FS_CONFIG_GEN) $(system_android_filesystem_config) $(system_capability_header)
	@mkdir -p $(dir $@)
	$< fsconfig \
	   --aid-header $(PRIVATE_ANDROID_FS_HDR) \
	   --capability-header $(PRIVATE_ANDROID_CAP_HDR) \
	   --partition system \
	   --all-partitions "$(subst $(space),$(comma),$(PRIVATE_PARTITION_LIST))" \
	   --files \
	   --out_file $@ \
	   $(or $(PRIVATE_TARGET_FS_CONFIG_GEN),/dev/null)

##################################
# Generate the vendor/etc/fs_config_dirs binary file for the target if the
# vendor partition is generated. Add fs_config_dirs or fs_config_dirs_vendor to
# PRODUCT_PACKAGES in the device make file to enable.
include $(CLEAR_VARS)

LOCAL_MODULE := fs_config_dirs_vendor
LOCAL_REQUIRED_MODULES := $(if $(filter vendor,$(fs_config_generate_extra_partition_list)),_fs_config_dirs_vendor)
include $(BUILD_PHONY_PACKAGE)

##################################
# Generate the vendor/etc/fs_config_files binary file for the target if the
# vendor partition is generated. Add fs_config_files or fs_config_files_vendor
# to PRODUCT_PACKAGES in the device make file to enable.
include $(CLEAR_VARS)

LOCAL_MODULE := fs_config_files_vendor
LOCAL_REQUIRED_MODULES := $(if $(filter vendor,$(fs_config_generate_extra_partition_list)),_fs_config_files_vendor)
include $(BUILD_PHONY_PACKAGE)

ifneq ($(filter vendor,$(fs_config_generate_extra_partition_list)),)
##################################
# Generate the vendor/etc/fs_config_dirs binary file for the target
# Do not add _fs_config_dirs_vendor to PRODUCT_PACKAGES, instead use
# fs_config_dirs_vendor
include $(CLEAR_VARS)

LOCAL_MODULE := _fs_config_dirs_vendor
LOCAL_MODULE_CLASS := ETC
LOCAL_INSTALLED_MODULE_STEM := fs_config_dirs
LOCAL_MODULE_PATH := $(TARGET_OUT_VENDOR)/etc
include $(BUILD_SYSTEM)/base_rules.mk
$(LOCAL_BUILT_MODULE): PRIVATE_ANDROID_FS_HDR := $(system_android_filesystem_config)
$(LOCAL_BUILT_MODULE): PRIVATE_ANDROID_CAP_HDR := $(system_capability_header)
$(LOCAL_BUILT_MODULE): PRIVATE_TARGET_FS_CONFIG_GEN := $(TARGET_FS_CONFIG_GEN)
$(LOCAL_BUILT_MODULE): $(LOCAL_PATH)/fs_config_generator.py $(TARGET_FS_CONFIG_GEN) $(system_android_filesystem_config) $(system_capability_header)
	@mkdir -p $(dir $@)
	$< fsconfig \
	   --aid-header $(PRIVATE_ANDROID_FS_HDR) \
	   --capability-header $(PRIVATE_ANDROID_CAP_HDR) \
	   --partition vendor \
	   --dirs \
	   --out_file $@ \
	   $(or $(PRIVATE_TARGET_FS_CONFIG_GEN),/dev/null)

##################################
# Generate the vendor/etc/fs_config_files binary file for the target
# Do not add _fs_config_files_vendor to PRODUCT_PACKAGES, instead use
# fs_config_files_vendor
include $(CLEAR_VARS)

LOCAL_MODULE := _fs_config_files_vendor
LOCAL_MODULE_CLASS := ETC
LOCAL_INSTALLED_MODULE_STEM := fs_config_files
LOCAL_MODULE_PATH := $(TARGET_OUT_VENDOR)/etc
include $(BUILD_SYSTEM)/base_rules.mk
$(LOCAL_BUILT_MODULE): PRIVATE_ANDROID_FS_HDR := $(system_android_filesystem_config)
$(LOCAL_BUILT_MODULE): PRIVATE_ANDROID_CAP_HDR := $(system_capability_header)
$(LOCAL_BUILT_MODULE): PRIVATE_TARGET_FS_CONFIG_GEN := $(TARGET_FS_CONFIG_GEN)
$(LOCAL_BUILT_MODULE): $(LOCAL_PATH)/fs_config_generator.py $(TARGET_FS_CONFIG_GEN) $(system_android_filesystem_config) $(system_capability_header)
	@mkdir -p $(dir $@)
	$< fsconfig \
	   --aid-header $(PRIVATE_ANDROID_FS_HDR) \
	   --capability-header $(PRIVATE_ANDROID_CAP_HDR) \
	   --partition vendor \
	   --files \
	   --out_file $@ \
	   $(or $(PRIVATE_TARGET_FS_CONFIG_GEN),/dev/null)

endif

##################################
# Generate the oem/etc/fs_config_dirs binary file for the target if the oem
# partition is generated. Add fs_config_dirs or fs_config_dirs_oem to
# PRODUCT_PACKAGES in the device make file to enable.
include $(CLEAR_VARS)

LOCAL_MODULE := fs_config_dirs_oem
LOCAL_REQUIRED_MODULES := $(if $(filter oem,$(fs_config_generate_extra_partition_list)),_fs_config_dirs_oem)
include $(BUILD_PHONY_PACKAGE)

##################################
# Generate the oem/etc/fs_config_files binary file for the target if the
# oem partition is generated. Add fs_config_files or fs_config_files_oem
# to PRODUCT_PACKAGES in the device make file to enable.
include $(CLEAR_VARS)

LOCAL_MODULE := fs_config_files_oem
LOCAL_REQUIRED_MODULES := $(if $(filter oem,$(fs_config_generate_extra_partition_list)),_fs_config_files_oem)
include $(BUILD_PHONY_PACKAGE)

ifneq ($(filter oem,$(fs_config_generate_extra_partition_list)),)
##################################
# Generate the oem/etc/fs_config_dirs binary file for the target
# Do not add _fs_config_dirs_oem to PRODUCT_PACKAGES, instead use
# fs_config_dirs_oem
include $(CLEAR_VARS)

LOCAL_MODULE := _fs_config_dirs_oem
LOCAL_MODULE_CLASS := ETC
LOCAL_INSTALLED_MODULE_STEM := fs_config_dirs
LOCAL_MODULE_PATH := $(TARGET_OUT_OEM)/etc
include $(BUILD_SYSTEM)/base_rules.mk
$(LOCAL_BUILT_MODULE): PRIVATE_ANDROID_FS_HDR := $(system_android_filesystem_config)
$(LOCAL_BUILT_MODULE): PRIVATE_ANDROID_CAP_HDR := $(system_capability_header)
$(LOCAL_BUILT_MODULE): PRIVATE_TARGET_FS_CONFIG_GEN := $(TARGET_FS_CONFIG_GEN)
$(LOCAL_BUILT_MODULE): $(LOCAL_PATH)/fs_config_generator.py $(TARGET_FS_CONFIG_GEN) $(system_android_filesystem_config) $(system_capability_header)
	@mkdir -p $(dir $@)
	$< fsconfig \
	   --aid-header $(PRIVATE_ANDROID_FS_HDR) \
	   --capability-header $(PRIVATE_ANDROID_CAP_HDR) \
	   --partition oem \
	   --dirs \
	   --out_file $@ \
	   $(or $(PRIVATE_TARGET_FS_CONFIG_GEN),/dev/null)

##################################
# Generate the oem/etc/fs_config_files binary file for the target
# Do not add _fs_config_files_oem to PRODUCT_PACKAGES, instead use
# fs_config_files_oem
include $(CLEAR_VARS)

LOCAL_MODULE := _fs_config_files_oem
LOCAL_MODULE_CLASS := ETC
LOCAL_INSTALLED_MODULE_STEM := fs_config_files
LOCAL_MODULE_PATH := $(TARGET_OUT_OEM)/etc
include $(BUILD_SYSTEM)/base_rules.mk
$(LOCAL_BUILT_MODULE): PRIVATE_ANDROID_FS_HDR := $(system_android_filesystem_config)
$(LOCAL_BUILT_MODULE): PRIVATE_ANDROID_CAP_HDR := $(system_capability_header)
$(LOCAL_BUILT_MODULE): PRIVATE_TARGET_FS_CONFIG_GEN := $(TARGET_FS_CONFIG_GEN)
$(LOCAL_BUILT_MODULE): $(LOCAL_PATH)/fs_config_generator.py $(TARGET_FS_CONFIG_GEN) $(system_android_filesystem_config) $(system_capability_header)
	@mkdir -p $(dir $@)
	$< fsconfig \
	   --aid-header $(PRIVATE_ANDROID_FS_HDR) \
	   --capability-header $(PRIVATE_ANDROID_CAP_HDR) \
	   --partition oem \
	   --files \
	   --out_file $@ \
	   $(or $(PRIVATE_TARGET_FS_CONFIG_GEN),/dev/null)

endif

##################################
# Generate the odm/etc/fs_config_dirs binary file for the target if the odm
# partition is generated. Add fs_config_dirs or fs_config_dirs_odm to
# PRODUCT_PACKAGES in the device make file to enable.
include $(CLEAR_VARS)

LOCAL_MODULE := fs_config_dirs_odm
LOCAL_REQUIRED_MODULES := $(if $(filter odm,$(fs_config_generate_extra_partition_list)),_fs_config_dirs_odm)
include $(BUILD_PHONY_PACKAGE)

##################################
# Generate the odm/etc/fs_config_files binary file for the target if the
# odm partition is generated. Add fs_config_files or fs_config_files_odm
# to PRODUCT_PACKAGES in the device make file to enable.
include $(CLEAR_VARS)

LOCAL_MODULE := fs_config_files_odm
LOCAL_REQUIRED_MODULES := $(if $(filter odm,$(fs_config_generate_extra_partition_list)),_fs_config_files_odm)
include $(BUILD_PHONY_PACKAGE)

ifneq ($(filter odm,$(fs_config_generate_extra_partition_list)),)
##################################
# Generate the odm/etc/fs_config_dirs binary file for the target
# Do not add _fs_config_dirs_odm to PRODUCT_PACKAGES, instead use
# fs_config_dirs_odm
include $(CLEAR_VARS)

LOCAL_MODULE := _fs_config_dirs_odm
LOCAL_MODULE_CLASS := ETC
LOCAL_INSTALLED_MODULE_STEM := fs_config_dirs
LOCAL_MODULE_PATH := $(TARGET_OUT_ODM)/etc
include $(BUILD_SYSTEM)/base_rules.mk
$(LOCAL_BUILT_MODULE): PRIVATE_ANDROID_FS_HDR := $(system_android_filesystem_config)
$(LOCAL_BUILT_MODULE): PRIVATE_ANDROID_CAP_HDR := $(system_capability_header)
$(LOCAL_BUILT_MODULE): PRIVATE_TARGET_FS_CONFIG_GEN := $(TARGET_FS_CONFIG_GEN)
$(LOCAL_BUILT_MODULE): $(LOCAL_PATH)/fs_config_generator.py $(TARGET_FS_CONFIG_GEN) $(system_android_filesystem_config) $(system_capability_header)
	@mkdir -p $(dir $@)
	$< fsconfig \
	   --aid-header $(PRIVATE_ANDROID_FS_HDR) \
	   --capability-header $(PRIVATE_ANDROID_CAP_HDR) \
	   --partition odm \
	   --dirs \
	   --out_file $@ \
	   $(or $(PRIVATE_TARGET_FS_CONFIG_GEN),/dev/null)

##################################
# Generate the odm/etc/fs_config_files binary file for the target
# Do not add _fs_config_files_odm to PRODUCT_PACKAGES, instead use
# fs_config_files_odm
include $(CLEAR_VARS)

LOCAL_MODULE := _fs_config_files_odm
LOCAL_MODULE_CLASS := ETC
LOCAL_INSTALLED_MODULE_STEM := fs_config_files
LOCAL_MODULE_PATH := $(TARGET_OUT_ODM)/etc
include $(BUILD_SYSTEM)/base_rules.mk
$(LOCAL_BUILT_MODULE): PRIVATE_ANDROID_FS_HDR := $(system_android_filesystem_config)
$(LOCAL_BUILT_MODULE): PRIVATE_ANDROID_CAP_HDR := $(system_capability_header)
$(LOCAL_BUILT_MODULE): PRIVATE_TARGET_FS_CONFIG_GEN := $(TARGET_FS_CONFIG_GEN)
$(LOCAL_BUILT_MODULE): $(LOCAL_PATH)/fs_config_generator.py $(TARGET_FS_CONFIG_GEN) $(system_android_filesystem_config) $(system_capability_header)
	@mkdir -p $(dir $@)
	$< fsconfig \
	   --aid-header $(PRIVATE_ANDROID_FS_HDR) \
	   --capability-header $(PRIVATE_ANDROID_CAP_HDR) \
	   --partition odm \
	   --files \
	   --out_file $@ \
	   $(or $(PRIVATE_TARGET_FS_CONFIG_GEN),/dev/null)

endif

##################################
# Generate the vendor_dlkm/etc/fs_config_dirs binary file for the target if the
# vendor_dlkm partition is generated. Add fs_config_dirs or fs_config_dirs_vendor_dlkm
# to PRODUCT_PACKAGES in the device make file to enable.
include $(CLEAR_VARS)

LOCAL_MODULE := fs_config_dirs_vendor_dlkm
LOCAL_REQUIRED_MODULES := $(if $(filter vendor_dlkm,$(fs_config_generate_extra_partition_list)),_fs_config_dirs_vendor_dlkm)
include $(BUILD_PHONY_PACKAGE)

##################################
# Generate the vendor_dlkm/etc/fs_config_files binary file for the target if the
# vendor_dlkm partition is generated. Add fs_config_files or fs_config_files_vendor_dlkm
# to PRODUCT_PACKAGES in the device make file to enable.
include $(CLEAR_VARS)

LOCAL_MODULE := fs_config_files_vendor_dlkm
LOCAL_REQUIRED_MODULES := $(if $(filter vendor_dlkm,$(fs_config_generate_extra_partition_list)),_fs_config_files_vendor_dlkm)
include $(BUILD_PHONY_PACKAGE)

ifneq ($(filter vendor_dlkm,$(fs_config_generate_extra_partition_list)),)
##################################
# Generate the vendor_dlkm/etc/fs_config_dirs binary file for the target
# Do not add _fs_config_dirs_vendor_dlkm to PRODUCT_PACKAGES, instead use
# fs_config_dirs_vendor_dlkm
include $(CLEAR_VARS)

LOCAL_MODULE := _fs_config_dirs_vendor_dlkm
LOCAL_MODULE_CLASS := ETC
LOCAL_INSTALLED_MODULE_STEM := fs_config_dirs
LOCAL_MODULE_PATH := $(TARGET_OUT_VENDOR_DLKM)/etc
include $(BUILD_SYSTEM)/base_rules.mk
$(LOCAL_BUILT_MODULE): PRIVATE_ANDROID_FS_HDR := $(system_android_filesystem_config)
$(LOCAL_BUILT_MODULE): PRIVATE_ANDROID_CAP_HDR := $(system_capability_header)
$(LOCAL_BUILT_MODULE): PRIVATE_TARGET_FS_CONFIG_GEN := $(TARGET_FS_CONFIG_GEN)
$(LOCAL_BUILT_MODULE): $(LOCAL_PATH)/fs_config_generator.py $(TARGET_FS_CONFIG_GEN) $(system_android_filesystem_config) $(system_capability_header)
	@mkdir -p $(dir $@)
	$< fsconfig \
	   --aid-header $(PRIVATE_ANDROID_FS_HDR) \
	   --capability-header $(PRIVATE_ANDROID_CAP_HDR) \
	   --partition vendor_dlkm \
	   --dirs \
	   --out_file $@ \
	   $(or $(PRIVATE_TARGET_FS_CONFIG_GEN),/dev/null)

##################################
# Generate the vendor_dlkm/etc/fs_config_files binary file for the target
# Do not add _fs_config_files_vendor_dlkm to PRODUCT_PACKAGES, instead use
# fs_config_files_vendor_dlkm
include $(CLEAR_VARS)

LOCAL_MODULE := _fs_config_files_vendor_dlkm
LOCAL_MODULE_CLASS := ETC
LOCAL_INSTALLED_MODULE_STEM := fs_config_files
LOCAL_MODULE_PATH := $(TARGET_OUT_VENDOR_DLKM)/etc
include $(BUILD_SYSTEM)/base_rules.mk
$(LOCAL_BUILT_MODULE): PRIVATE_ANDROID_FS_HDR := $(system_android_filesystem_config)
$(LOCAL_BUILT_MODULE): PRIVATE_ANDROID_CAP_HDR := $(system_capability_header)
$(LOCAL_BUILT_MODULE): PRIVATE_TARGET_FS_CONFIG_GEN := $(TARGET_FS_CONFIG_GEN)
$(LOCAL_BUILT_MODULE): $(LOCAL_PATH)/fs_config_generator.py $(TARGET_FS_CONFIG_GEN) $(system_android_filesystem_config) $(system_capability_header)
	@mkdir -p $(dir $@)
	$< fsconfig \
	   --aid-header $(PRIVATE_ANDROID_FS_HDR) \
	   --capability-header $(PRIVATE_ANDROID_CAP_HDR) \
	   --partition vendor_dlkm \
	   --files \
	   --out_file $@ \
	   $(or $(PRIVATE_TARGET_FS_CONFIG_GEN),/dev/null)

endif

##################################
# Generate the odm_dlkm/etc/fs_config_dirs binary file for the target if the
# odm_dlkm partition is generated. Add fs_config_dirs or fs_config_dirs_odm_dlkm
# to PRODUCT_PACKAGES in the device make file to enable.
include $(CLEAR_VARS)

LOCAL_MODULE := fs_config_dirs_odm_dlkm
LOCAL_REQUIRED_MODULES := $(if $(filter odm_dlkm,$(fs_config_generate_extra_partition_list)),_fs_config_dirs_odm_dlkm)
include $(BUILD_PHONY_PACKAGE)

##################################
# Generate the odm_dlkm/etc/fs_config_files binary file for the target if the
# odm_dlkm partition is generated. Add fs_config_files or fs_config_files_odm_dlkm
# to PRODUCT_PACKAGES in the device make file to enable.
include $(CLEAR_VARS)

LOCAL_MODULE := fs_config_files_odm_dlkm
LOCAL_REQUIRED_MODULES := $(if $(filter odm_dlkm,$(fs_config_generate_extra_partition_list)),_fs_config_files_odm_dlkm)
include $(BUILD_PHONY_PACKAGE)

ifneq ($(filter odm_dlkm,$(fs_config_generate_extra_partition_list)),)
##################################
# Generate the odm_dlkm/etc/fs_config_dirs binary file for the target
# Do not add _fs_config_dirs_odm_dlkm to PRODUCT_PACKAGES, instead use
# fs_config_dirs_odm_dlkm
include $(CLEAR_VARS)

LOCAL_MODULE := _fs_config_dirs_odm_dlkm
LOCAL_MODULE_CLASS := ETC
LOCAL_INSTALLED_MODULE_STEM := fs_config_dirs
LOCAL_MODULE_PATH := $(TARGET_OUT_ODM_DLKM)/etc
include $(BUILD_SYSTEM)/base_rules.mk
$(LOCAL_BUILT_MODULE): PRIVATE_ANDROID_FS_HDR := $(system_android_filesystem_config)
$(LOCAL_BUILT_MODULE): PRIVATE_ANDROID_CAP_HDR := $(system_capability_header)
$(LOCAL_BUILT_MODULE): PRIVATE_TARGET_FS_CONFIG_GEN := $(TARGET_FS_CONFIG_GEN)
$(LOCAL_BUILT_MODULE): $(LOCAL_PATH)/fs_config_generator.py $(TARGET_FS_CONFIG_GEN) $(system_android_filesystem_config) $(system_capability_header)
	@mkdir -p $(dir $@)
	$< fsconfig \
	   --aid-header $(PRIVATE_ANDROID_FS_HDR) \
	   --capability-header $(PRIVATE_ANDROID_CAP_HDR) \
	   --partition odm_dlkm \
	   --dirs \
	   --out_file $@ \
	   $(or $(PRIVATE_TARGET_FS_CONFIG_GEN),/dev/null)

##################################
# Generate the odm_dlkm/etc/fs_config_files binary file for the target
# Do not add _fs_config_files_odm_dlkm to PRODUCT_PACKAGES, instead use
# fs_config_files_odm_dlkm
include $(CLEAR_VARS)

LOCAL_MODULE := _fs_config_files_odm_dlkm
LOCAL_MODULE_CLASS := ETC
LOCAL_INSTALLED_MODULE_STEM := fs_config_files
LOCAL_MODULE_PATH := $(TARGET_OUT_ODM_DLKM)/etc
include $(BUILD_SYSTEM)/base_rules.mk
$(LOCAL_BUILT_MODULE): PRIVATE_ANDROID_FS_HDR := $(system_android_filesystem_config)
$(LOCAL_BUILT_MODULE): PRIVATE_ANDROID_CAP_HDR := $(system_capability_header)
$(LOCAL_BUILT_MODULE): PRIVATE_TARGET_FS_CONFIG_GEN := $(TARGET_FS_CONFIG_GEN)
$(LOCAL_BUILT_MODULE): $(LOCAL_PATH)/fs_config_generator.py $(TARGET_FS_CONFIG_GEN) $(system_android_filesystem_config) $(system_capability_header)
	@mkdir -p $(dir $@)
	$< fsconfig \
	   --aid-header $(PRIVATE_ANDROID_FS_HDR) \
	   --capability-header $(PRIVATE_ANDROID_CAP_HDR) \
	   --partition odm_dlkm \
	   --files \
	   --out_file $@ \
	   $(or $(PRIVATE_TARGET_FS_CONFIG_GEN),/dev/null)

endif

##################################
# Generate the product/etc/fs_config_dirs binary file for the target if the
# product partition is generated. Add fs_config_dirs or fs_config_dirs_product
# to PRODUCT_PACKAGES in the device make file to enable.
include $(CLEAR_VARS)

LOCAL_MODULE := fs_config_dirs_product
LOCAL_REQUIRED_MODULES := $(if $(filter product,$(fs_config_generate_extra_partition_list)),_fs_config_dirs_product)
include $(BUILD_PHONY_PACKAGE)

##################################
# Generate the product/etc/fs_config_files binary file for the target if the
# product partition is generated. Add fs_config_files or fs_config_files_product
# to PRODUCT_PACKAGES in the device make file to enable.
include $(CLEAR_VARS)

LOCAL_MODULE := fs_config_files_product
LOCAL_REQUIRED_MODULES := $(if $(filter product,$(fs_config_generate_extra_partition_list)),_fs_config_files_product)
include $(BUILD_PHONY_PACKAGE)

ifneq ($(filter product,$(fs_config_generate_extra_partition_list)),)
##################################
# Generate the product/etc/fs_config_dirs binary file for the target
# Do not add _fs_config_dirs_product to PRODUCT_PACKAGES, instead use
# fs_config_dirs_product
include $(CLEAR_VARS)

LOCAL_MODULE := _fs_config_dirs_product
LOCAL_MODULE_CLASS := ETC
LOCAL_INSTALLED_MODULE_STEM := fs_config_dirs
LOCAL_MODULE_PATH := $(TARGET_OUT_PRODUCT)/etc
include $(BUILD_SYSTEM)/base_rules.mk
$(LOCAL_BUILT_MODULE): PRIVATE_ANDROID_FS_HDR := $(system_android_filesystem_config)
$(LOCAL_BUILT_MODULE): PRIVATE_ANDROID_CAP_HDR := $(system_capability_header)
$(LOCAL_BUILT_MODULE): PRIVATE_TARGET_FS_CONFIG_GEN := $(TARGET_FS_CONFIG_GEN)
$(LOCAL_BUILT_MODULE): $(LOCAL_PATH)/fs_config_generator.py $(TARGET_FS_CONFIG_GEN) $(system_android_filesystem_config) $(system_capability_header)
	@mkdir -p $(dir $@)
	$< fsconfig \
	   --aid-header $(PRIVATE_ANDROID_FS_HDR) \
	   --capability-header $(PRIVATE_ANDROID_CAP_HDR) \
	   --partition product \
	   --dirs \
	   --out_file $@ \
	   $(or $(PRIVATE_TARGET_FS_CONFIG_GEN),/dev/null)

##################################
# Generate the product/etc/fs_config_files binary file for the target
# Do not add _fs_config_files_product to PRODUCT_PACKAGES, instead use
# fs_config_files_product
include $(CLEAR_VARS)

LOCAL_MODULE := _fs_config_files_product
LOCAL_MODULE_CLASS := ETC
LOCAL_INSTALLED_MODULE_STEM := fs_config_files
LOCAL_MODULE_PATH := $(TARGET_OUT_PRODUCT)/etc
include $(BUILD_SYSTEM)/base_rules.mk
$(LOCAL_BUILT_MODULE): PRIVATE_ANDROID_FS_HDR := $(system_android_filesystem_config)
$(LOCAL_BUILT_MODULE): PRIVATE_ANDROID_CAP_HDR := $(system_capability_header)
$(LOCAL_BUILT_MODULE): PRIVATE_TARGET_FS_CONFIG_GEN := $(TARGET_FS_CONFIG_GEN)
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

##################################
# Generate the system_ext/etc/fs_config_dirs binary file for the target if the
# system_ext partition is generated. Add fs_config_dirs or fs_config_dirs_system_ext
# to PRODUCT_PACKAGES in the device make file to enable.
include $(CLEAR_VARS)

LOCAL_MODULE := fs_config_dirs_system_ext
LOCAL_REQUIRED_MODULES := $(if $(filter system_ext,$(fs_config_generate_extra_partition_list)),_fs_config_dirs_system_ext)
include $(BUILD_PHONY_PACKAGE)

##################################
# Generate the system_ext/etc/fs_config_files binary file for the target if the
# system_ext partition is generated. Add fs_config_files or fs_config_files_system_ext
# to PRODUCT_PACKAGES in the device make file to enable.
include $(CLEAR_VARS)

LOCAL_MODULE := fs_config_files_system_ext
LOCAL_REQUIRED_MODULES := $(if $(filter system_ext,$(fs_config_generate_extra_partition_list)),_fs_config_files_system_ext)
include $(BUILD_PHONY_PACKAGE)

ifneq ($(filter system_ext,$(fs_config_generate_extra_partition_list)),)
##################################
# Generate the system_ext/etc/fs_config_dirs binary file for the target
# Do not add _fs_config_dirs_system_ext to PRODUCT_PACKAGES, instead use
# fs_config_dirs_system_ext
include $(CLEAR_VARS)

LOCAL_MODULE := _fs_config_dirs_system_ext
LOCAL_MODULE_CLASS := ETC
LOCAL_INSTALLED_MODULE_STEM := fs_config_dirs
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
	   --dirs \
	   --out_file $@ \
	   $(or $(PRIVATE_TARGET_FS_CONFIG_GEN),/dev/null)

##################################
# Generate the system_ext/etc/fs_config_files binary file for the target
# Do not add _fs_config_files_system_ext to PRODUCT_PACKAGES, instead use
# fs_config_files_system_ext
include $(CLEAR_VARS)

LOCAL_MODULE := _fs_config_files_system_ext
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

system_android_filesystem_config :=
system_capability_header :=
fs_config_generate_extra_partition_list :=
