#
# Copyright 2013 The Android Open-Source Project
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

PRODUCT_PROPERTY_OVERRIDES += \
	vendor.rild.libpath=/vendor/lib64/libreference-ril.so

# This is a build configuration for a full-featured build of the
# Open-Source part of the tree. It's geared toward a US-centric
# build quite specifically for the emulator, and might not be
# entirely appropriate to inherit from for on-device configurations.

<<<<<<< HEAD   (5c8d84 Merge "Merge empty history for sparse-6676661-L8360000065797)
PRODUCT_COPY_FILES += \
    development/sys-img/advancedFeatures.ini:advancedFeatures.ini \
    device/generic/goldfish/data/etc/encryptionkey.img:encryptionkey.img \
    prebuilts/qemu-kernel/x86_64/4.9/kernel-qemu2:kernel-ranchu
=======
#
# All components inherited here go to system image
#
$(call inherit-product, $(SRC_TARGET_DIR)/product/core_64_bit.mk)
$(call inherit-product, $(SRC_TARGET_DIR)/product/mainline_system.mk)

# Enable mainline checking for excat this product name
ifeq (aosp_x86_64,$(TARGET_PRODUCT))
PRODUCT_ENFORCE_ARTIFACT_PATH_REQUIREMENTS := relaxed
endif

#
# All components inherited here go to system_ext image
#
$(call inherit-product, $(SRC_TARGET_DIR)/product/handheld_system_ext.mk)
$(call inherit-product, $(SRC_TARGET_DIR)/product/telephony_system_ext.mk)
>>>>>>> BRANCH (a10c18 Merge "Version bump to RT11.201014.001.A1 [core/build_id.mk])

#
# All components inherited here go to product image
#
$(call inherit-product, $(SRC_TARGET_DIR)/product/aosp_product.mk)

#
# All components inherited here go to vendor image
#
$(call inherit-product-if-exists, device/generic/goldfish/x86_64-vendor.mk)
$(call inherit-product, $(SRC_TARGET_DIR)/product/emulator_vendor.mk)
$(call inherit-product, $(SRC_TARGET_DIR)/board/generic_x86_64/device.mk)

#
# Special settings for GSI releasing
#
ifeq (aosp_x86_64,$(TARGET_PRODUCT))
$(call inherit-product, $(SRC_TARGET_DIR)/product/gsi_release.mk)
endif


# TODO(b/78308559): includes vr_hwc into GSI before vr_hwc move to vendor
PRODUCT_PACKAGES += \
    vr_hwc

$(call inherit-product, $(SRC_TARGET_DIR)/product/emulator.mk)
$(call inherit-product, $(SRC_TARGET_DIR)/product/core_64_bit.mk)
$(call inherit-product, $(SRC_TARGET_DIR)/product/aosp_base_telephony.mk)
$(call inherit-product, $(SRC_TARGET_DIR)/board/generic_x86_64/device.mk)

# Needed by Pi newly launched device to pass VtsTrebleSysProp on GSI
PRODUCT_COMPATIBLE_PROPERTY_OVERRIDE := true

ifdef NET_ETH0_STARTONBOOT
  PRODUCT_PROPERTY_OVERRIDES += net.eth0.startonboot=1
endif

# Ensure we package the BIOS files too.
PRODUCT_PACKAGES += \
	bios.bin \
	vgabios-cirrus.bin \

# Overrides
PRODUCT_NAME := aosp_x86_64
PRODUCT_DEVICE := generic_x86_64
PRODUCT_BRAND := Android
PRODUCT_MODEL := AOSP on IA x86_64 Emulator
