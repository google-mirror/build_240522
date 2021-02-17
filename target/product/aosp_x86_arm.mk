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


<<<<<<< HEAD   (4be654 Merge "Merge empty history for sparse-7121469-L4290000080720)
# aosp_x86 with arm libraries needed by binary translation.
=======
#
# All components inherited here go to system image
#
$(call inherit-product, $(SRC_TARGET_DIR)/product/generic_system.mk)
>>>>>>> BRANCH (fe6ad7 Merge "Version bump to RBT1.210107.001.A1 [core/build_id.mk])

include $(SRC_TARGET_DIR)/product/full_x86.mk

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

PRODUCT_NAME := aosp_x86_arm
PRODUCT_DEVICE := generic_x86_arm
