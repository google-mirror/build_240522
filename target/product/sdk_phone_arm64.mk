#
# Copyright (C) 2009 The Android Open Source Project
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
<<<<<<< HEAD   (4be654 Merge "Merge empty history for sparse-7121469-L4290000080720)

PRODUCT_PROPERTY_OVERRIDES += \
	rild.libpath=/vendor/lib64/libreference-ril.so
=======
QEMU_USE_SYSTEM_EXT_PARTITIONS := true
PRODUCT_USE_DYNAMIC_PARTITIONS := true
>>>>>>> BRANCH (fe6ad7 Merge "Version bump to RBT1.210107.001.A1 [core/build_id.mk])

# This is a build configuration for a full-featured build of the
# Open-Source part of the tree. It's geared toward a US-centric
# build quite specifically for the emulator, and might not be
# entirely appropriate to inherit from for on-device configurations.

<<<<<<< HEAD   (4be654 Merge "Merge empty history for sparse-7121469-L4290000080720)
# Note: the following lines need to stay at the beginning so that it can
# take priority  and override the rules it inherit from other mk files
# see copy file rules in core/Makefile
PRODUCT_COPY_FILES += \
    development/sys-img/advancedFeatures.ini.arm:advancedFeatures.ini \
    prebuilts/qemu-kernel/arm64/3.18/kernel-qemu2:kernel-ranchu \
    device/generic/goldfish/fstab.ranchu.arm:$(TARGET_COPY_OUT_VENDOR)/etc/fstab.ranchu

=======
# Enable mainline checking for exact this product name
ifeq (sdk_phone_arm64,$(TARGET_PRODUCT))
PRODUCT_ENFORCE_ARTIFACT_PATH_REQUIREMENTS := relaxed
endif

#
# All components inherited here go to system image
#
>>>>>>> BRANCH (fe6ad7 Merge "Version bump to RBT1.210107.001.A1 [core/build_id.mk])
$(call inherit-product, $(SRC_TARGET_DIR)/product/core_64_bit.mk)
<<<<<<< HEAD   (4be654 Merge "Merge empty history for sparse-7121469-L4290000080720)
$(call inherit-product, $(SRC_TARGET_DIR)/product/sdk_base.mk)
$(call inherit-product, $(SRC_TARGET_DIR)/board/generic_arm64/device.mk)
=======
$(call inherit-product, $(SRC_TARGET_DIR)/product/generic_system.mk)
>>>>>>> BRANCH (fe6ad7 Merge "Version bump to RBT1.210107.001.A1 [core/build_id.mk])

# AOSP emulator images build the AOSP messaging app.
# Google API images override with the Google API app.
# See vendor/google/products/sdk_google_phone_*.mk
PRODUCT_PACKAGES += \
    messaging

# Overrides
PRODUCT_BRAND := Android
PRODUCT_NAME := sdk_phone_arm64
PRODUCT_DEVICE := generic_arm64
PRODUCT_MODEL := Android SDK built for arm64


