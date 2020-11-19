#
# Copyright (C) 2020 The Android Open Source Project
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

# This makefile installs contents of the generic ramdisk.
# Inherit from this makefile to declare that this product uses generic ramdisk.
# This makefile checks that other makefiles must not install things to the
# ramdisk.

PRODUCT_NAME := generic_ramdisk
PRODUCT_BRAND := generic

# Ramdisk
PRODUCT_PACKAGES += \
    init_first_stage \
    e2fsck_ramdisk \

# Debug ramdisk
PRODUCT_PACKAGES += \
    userdebug_plat_sepolicy.cil \

_my_paths := \
    $(TARGET_COPY_OUT_RAMDISK) \
    $(TARGET_COPY_OUT_DEBUG_RAMDISK) \

$(call require-artifacts-in-path, $(_my_paths),$(_my_allowed_list))
