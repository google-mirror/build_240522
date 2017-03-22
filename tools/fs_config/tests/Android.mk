#
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
#

LOCAL_PATH := $(call my-dir)

test_module_prefix := fs_conf-
test_tags := tests

# -----------------------------------------------------------------------------
# Unit tests.
# -----------------------------------------------------------------------------

test_c_flags := \
    -fstack-protector-all \
    -g \
    -Wall -Wextra \
    -Werror \
    -fno-builtin \

test_src_files := \
    fs_conf_test.cpp

################################## test executable
include $(CLEAR_VARS)
LOCAL_SRC_FILES := ../fs_config_generate.c
LOCAL_MODULE := fs_config_generate_test
LOCAL_MODULE_TAGS := $(test_tags)
LOCAL_MODULE_CLASS := EXECUTABLES
LOCAL_SHARED_LIBRARIES := libcutils
LOCAL_CFLAGS := -Werror
include $(BUILD_HOST_EXECUTABLE)
fs_config_generate_test_bin := $(LOCAL_INSTALLED_MODULE)

################################## gTest tool
include $(CLEAR_VARS)
LOCAL_MODULE := $(test_module_prefix)unit-tests
LOCAL_MODULE_TAGS := $(test_tags)
LOCAL_MODULE_CLASS := EXECUTABLES
LOCAL_CFLAGS += $(test_c_flags)
LOCAL_CFLAGS += -DFS_CONFIG_GENERATE='"$(fs_config_generate_test_bin)"' -DHOST
LOCAL_SHARED_LIBRARIES := liblog libcutils libbase
LOCAL_STATIC_LIBRARIES := libgtest_host libgtest_main_host
LOCAL_ADDITIONAL_DEPENDENCIES += $(fs_config_generate_test_bin)
LOCAL_SRC_FILES := $(test_src_files)
include $(BUILD_HOST_EXECUTABLE)
