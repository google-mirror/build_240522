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
include $(CLEAR_VARS)

LOCAL_SRC_FILES := fs_config.c
LOCAL_MODULE := fs_config
LOCAL_STATIC_LIBRARIES := libcutils libselinux
LOCAL_FORCE_STATIC_EXECUTABLE := true
LOCAL_CFLAGS := -Werror

include $(BUILD_HOST_EXECUTABLE)

# To Build the custom target binary for the host to generate the fs_config
# override files. The executable is hard coded to include the
# $(TARGET_ANDROID_FILESYSTEM_CONFIG_H) file if it exists.
# Expectations:
#    device/<vendor>/<device>/android_filesystem_config.h
#        fills in struct fs_path_config android_device_dirs[] and
#                 struct fs_path_config android_device_files[]
#    device/<vendor>/<device>/device.mk
#        PRODUCT_PACKAGES += fs_config_dirs fs_config_files

TARGET_ANDROID_FILESYSTEM_CONFIG_H := $(wildcard $(TARGET_DEVICE_DIR)/android_filesystem_config.h)

include $(CLEAR_VARS)

LOCAL_SRC_FILES := fs_config_generate.c
LOCAL_MODULE := fs_config_generate_$(TARGET_DEVICE)
LOCAL_MODULE_CLASS := EXECUTABLES
LOCAL_IS_HOST_MODULE := true
LOCAL_SHARED_LIBRARIES := libcutils

LOCAL_ADDITIONAL_DEPENDENCIES := $(TARGET_ANDROID_FILESYSTEM_CONFIG_H)
LOCAL_CFLAGS := $(TARGET_DEVICE_DIR:%=-I%) -Werror
LOCAL_CFLAGS += $(LOCAL_ADDITIONAL_DEPENDENCIES:%=-DANDROID_FILESYSTEM_CONFIG=%)

include $(BUILD_HOST_EXECUTABLE)

# Generate the system/etc/fs_config_dirs binary file for the target
# Add fs_config_dirs to PRODUCT_PACKAGES in the device make file to enable
include $(CLEAR_VARS)

LOCAL_MODULE := fs_config_dirs
LOCAL_MODULE_CLASS := ETC
include $(BUILD_SYSTEM)/base_rules.mk
TARGET_fs_config_dirs := $(LOCAL_BUILT_MODULE)
$(TARGET_fs_config_dirs): $(HOST_OUT_EXECUTABLES)/fs_config_generate_$(TARGET_DEVICE)
	-mkdir -p $(TARGET_fs_config_dirs:%/fs_config_dirs=%)
	$(HOST_OUT_EXECUTABLES)/fs_config_generate_$(TARGET_DEVICE) -D -o $(TARGET_fs_config_dirs)

# Generate the system/etc/fs_config_files binary file for the target
# Add fs_config_files to PRODUCT_PACKAGES in the device make file to enable
include $(CLEAR_VARS)

LOCAL_MODULE := fs_config_files
LOCAL_MODULE_CLASS := ETC
include $(BUILD_SYSTEM)/base_rules.mk
TARGET_fs_config_files := $(LOCAL_BUILT_MODULE)
$(TARGET_fs_config_files): $(HOST_OUT_EXECUTABLES)/fs_config_generate_$(TARGET_DEVICE)
	-mkdir -p $(TARGET_fs_config_files:%/fs_config_files=%)
	$(HOST_OUT_EXECUTABLES)/fs_config_generate_$(TARGET_DEVICE) -F -o $(TARGET_fs_config_files)
