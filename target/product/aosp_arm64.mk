#
# Copyright (C) 2017 The Android Open-Source Project
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

# The system image of aosp_arm64-userdebug is a GSI for the devices with:
# - ARM 64 bits user space
# - 64 bits binder interface
# - system-as-root
# - VNDK enforcement
# - compatible property override enabled

# This is a build configuration for a full-featured build of the
# Open-Source part of the tree. It's geared toward a US-centric
# build quite specifically for the emulator, and might not be
# entirely appropriate to inherit from for on-device configurations.

-include device/generic/goldfish/arm64-vendor.mk

# Copy different zygote settings for vendor.img to select by setting property
# ro.zygote=zygote64_32 or ro.zygote=zygote32_64:
#   1. 64-bit primary, 32-bit secondary OR
#   2. 32-bit primary, 64-bit secondary
# init.zygote64_32.rc is in the core_64_bit.mk below
PRODUCT_COPY_FILES += \
    system/core/rootdir/init.zygote32_64.rc:root/init.zygote32_64.rc

$(call inherit-product, $(SRC_TARGET_DIR)/product/emulator.mk)
$(call inherit-product, $(SRC_TARGET_DIR)/product/core_64_bit.mk)
$(call inherit-product, $(SRC_TARGET_DIR)/product/aosp_base_telephony.mk)
$(call inherit-product, $(SRC_TARGET_DIR)/board/generic_arm64/device.mk)

# Enable dynamic partition size
PRODUCT_USE_DYNAMIC_PARTITION_SIZE := true

# Enable A/B update
AB_OTA_UPDATER := true
AB_OTA_PARTITIONS := system
PRODUCT_PACKAGES += \
    update_engine \
    update_verifier

# Enable fake network location provider for GSI.
PRODUCT_PACKAGES += \
    StubNetworkLocation \

# Needed by Pi newly launched device to pass VtsTrebleSysProp on GSI
PRODUCT_COMPATIBLE_PROPERTY_OVERRIDE := true

# GSI specific tasks on boot
PRODUCT_COPY_FILES += \
    build/make/target/product/gsi/skip_mount.cfg:system/etc/init/config/skip_mount.cfg \
    build/make/target/product/gsi/init.gsi.rc:system/etc/init/init.gsi.rc \

# Support addtional P vendor interface
PRODUCT_EXTRA_VNDK_VERSIONS := 28

PRODUCT_NAME := aosp_arm64
PRODUCT_DEVICE := generic_arm64
PRODUCT_BRAND := Android
PRODUCT_MODEL := AOSP on ARM arm64 Emulator
