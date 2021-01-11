#
# Copyright (C) 2019 The Android Open Source Project
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

# Base modules and settings for the product partition.
PRODUCT_PACKAGES += \
    fs_config_dirs_product \
    fs_config_files_product \
    group_product \
    ModuleMetadata \
    passwd_product \
    product_compatibility_matrix.xml \
    product_manifest.xml \
    selinux_policy_product \

# System server should not contain compiled code.
PRODUCT_SYSTEM_SERVER_COMPILER_FILTER := verify

# Use the apex image for preopting.
DEXPREOPT_USE_ART_IMAGE := true

# Have the runtime pick up the JIT-zygote image.
PRODUCT_PROPERTY_OVERRIDES += \
    dalvik.vm.boot-image=boot.art:/nonx/boot-framework.art!/system/etc/boot-image.prof \
    dalvik.vm.image-dex2oat-filter=assume-verified
