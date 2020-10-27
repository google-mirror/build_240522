#
# Copyright 2016 The Android Open-Source Project
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


#
# All components inherited here go to system image
#
$(call inherit-product, $(SRC_TARGET_DIR)/product/mainline_system.mk)

<<<<<<< HEAD   (5c8d84 Merge "Merge empty history for sparse-6676661-L8360000065797)
include $(SRC_TARGET_DIR)/product/full_x86.mk
=======
# Enable mainline checking
ifeq (aosp_x86_arm,$(TARGET_PRODUCT))
PRODUCT_ENFORCE_ARTIFACT_PATH_REQUIREMENTS := relaxed
endif

# TODO (b/138382074): remove following setting after enable product/system_ext
PRODUCT_ARTIFACT_PATH_REQUIREMENT_ALLOWED_LIST += \
    system/product/% \
    system/system_ext/%

#
# All components inherited here go to system_ext image
#
$(call inherit-product, $(SRC_TARGET_DIR)/product/handheld_system_ext.mk)
$(call inherit-product, $(SRC_TARGET_DIR)/product/telephony_system_ext.mk)
>>>>>>> BRANCH (a10c18 Merge "Version bump to RT11.201014.001.A1 [core/build_id.mk])

<<<<<<< HEAD   (5c8d84 Merge "Merge empty history for sparse-6676661-L8360000065797)
# arm libraries. This is the list of shared libraries included in the NDK.
# Their dependency libraries will be automatically pulled in.
PRODUCT_PACKAGES += \
  libandroid_arm \
  libaaudio_arm \
  libc_arm \
  libdl_arm \
  libEGL_arm \
  libGLESv1_CM_arm \
  libGLESv2_arm \
  libGLESv3_arm \
  libjnigraphics_arm \
  liblog_arm \
  libm_arm \
  libmediandk_arm \
  libOpenMAXAL_arm \
  libstdc++_arm \
  libOpenSLES_arm \
  libz_arm \
=======
#
# All components inherited here go to product image
#
$(call inherit-product, $(SRC_TARGET_DIR)/product/aosp_product.mk)

#
# All components inherited here go to vendor image
#
$(call inherit-product-if-exists, device/generic/goldfish/x86-vendor.mk)
$(call inherit-product, $(SRC_TARGET_DIR)/product/emulator_vendor.mk)
$(call inherit-product, $(SRC_TARGET_DIR)/board/generic_x86_arm/device.mk)

>>>>>>> BRANCH (a10c18 Merge "Version bump to RT11.201014.001.A1 [core/build_id.mk])

PRODUCT_NAME := aosp_x86_arm
PRODUCT_DEVICE := generic_x86_arm
PRODUCT_BRAND := Android
PRODUCT_MODEL := AOSP on IA Emulator
