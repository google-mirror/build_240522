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
QEMU_USE_SYSTEM_EXT_PARTITIONS := true
PRODUCT_USE_DYNAMIC_PARTITIONS := true

#
# All components inherited here go to system image
#
$(call inherit-product, $(SRC_TARGET_DIR)/product/mainline_system.mk)

# Enable mainline checking for excat this product name
ifeq (sdk_phone_x86,$(TARGET_PRODUCT))
PRODUCT_ENFORCE_ARTIFACT_PATH_REQUIREMENTS := relaxed
endif

#
# All components inherited here go to product image
#
$(call inherit-product, $(SRC_TARGET_DIR)/product/aosp_product.mk)

#
# All components inherited here go to vendor image
#
$(call inherit-product-if-exists, device/generic/goldfish/x86-vendor.mk)
$(call inherit-product, $(SRC_TARGET_DIR)/product/emulator_vendor.mk)
$(call inherit-product, $(SRC_TARGET_DIR)/product/emulator_system.mk)
$(call inherit-product, $(SRC_TARGET_DIR)/board/generic_x86/device.mk)

$(call inherit-product, $(SRC_TARGET_DIR)/product/aosp_x86.mk)

# Define the host tools and libs that are parts of the SDK.
-include sdk/build/product_sdk.mk
-include development/build/product_sdk.mk

# Overrides
PRODUCT_BRAND := Android
PRODUCT_NAME := sdk_phone_x86
PRODUCT_DEVICE := generic_x86
PRODUCT_MODEL := Android SDK built for x86
