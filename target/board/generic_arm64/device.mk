#
# Copyright (C) 2013 The Android Open Source Project
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

# This is a build configuration for the product aspects that
# are specific to the emulator.

PRODUCT_COPY_FILES := \
    device/generic/goldfish/camera/media_profiles.xml:$(TARGET_COPY_OUT_VENDOR)/etc/media_profiles.xml \
    frameworks/av/media/libstagefright/data/media_codecs_google_audio.xml:$(TARGET_COPY_OUT_VENDOR)/etc/media_codecs_google_audio.xml \
    frameworks/av/media/libstagefright/data/media_codecs_google_telephony.xml:$(TARGET_COPY_OUT_VENDOR)/etc/media_codecs_google_telephony.xml \
    frameworks/av/media/libstagefright/data/media_codecs_google_video.xml:$(TARGET_COPY_OUT_VENDOR)/etc/media_codecs_google_video.xml \
    device/generic/goldfish/camera/media_codecs.xml:$(TARGET_COPY_OUT_VENDOR)/etc/media_codecs.xml

# minimal configuration for audio policy.
PRODUCT_COPY_FILES += \
<<<<<<< HEAD   (4be654 Merge "Merge empty history for sparse-7121469-L4290000080720)
    frameworks/av/services/audiopolicy/config/audio_policy_configuration_generic.xml:system/etc/audio_policy_configuration.xml \
    frameworks/av/services/audiopolicy/config/primary_audio_policy_configuration.xml:system/etc/primary_audio_policy_configuration.xml \
    frameworks/av/services/audiopolicy/config/r_submix_audio_policy_configuration.xml:system/etc/r_submix_audio_policy_configuration.xml \
    frameworks/av/services/audiopolicy/config/audio_policy_volumes.xml:system/etc/audio_policy_volumes.xml \
    frameworks/av/services/audiopolicy/config/default_volume_tables.xml:system/etc/default_volume_tables.xml \

# NFC:
#   Provide default libnfc-nci.conf file for devices that does not have one in
#   vendor/etc because aosp system image (of aosp_$arch products) is going to
#   be used as GSI.
#   May need to remove the following for newly launched devices in P since this
#   NFC configuration file should be in vendor/etc, instead of system/etc
PRODUCT_COPY_FILES += \
    device/generic/common/nfc/libnfc-nci.conf:system/etc/libnfc-nci.conf
=======
    kernel/prebuilts/4.19/arm64/kernel-4.19-gz:kernel-4.19-gz \
    kernel/prebuilts/5.4/arm64/kernel-5.4:kernel-5.4 \
    kernel/prebuilts/5.4/arm64/kernel-5.4-gz:kernel-5.4-gz \
    kernel/prebuilts/5.4/arm64/kernel-5.4-lz4:kernel-5.4-lz4 \
    kernel/prebuilts/5.10/arm64/kernel-5.10:kernel-5.10 \
    kernel/prebuilts/5.10/arm64/kernel-5.10-gz:kernel-5.10-gz \
    kernel/prebuilts/5.10/arm64/kernel-5.10-lz4:kernel-5.10-lz4 \
    kernel/prebuilts/mainline/arm64/kernel-mainline-allsyms:kernel-mainline \
    kernel/prebuilts/mainline/arm64/kernel-mainline-gz-allsyms:kernel-mainline-gz \
    kernel/prebuilts/mainline/arm64/kernel-mainline-lz4-allsyms:kernel-mainline-lz4
>>>>>>> BRANCH (fe6ad7 Merge "Version bump to RBT1.210107.001.A1 [core/build_id.mk])

ifneq (,$(filter userdebug eng,$(TARGET_BUILD_VARIANT)))
PRODUCT_COPY_FILES += \
    kernel/prebuilts/5.4/arm64/kernel-5.4:kernel-5.4-allsyms \
    kernel/prebuilts/5.4/arm64/kernel-5.4-gz:kernel-5.4-gz-allsyms \
    kernel/prebuilts/5.4/arm64/kernel-5.4-lz4:kernel-5.4-lz4-allsyms \
    kernel/prebuilts/5.10/arm64/kernel-5.10:kernel-5.10-allsyms \
    kernel/prebuilts/5.10/arm64/kernel-5.10-gz:kernel-5.10-gz-allsyms \
    kernel/prebuilts/5.10/arm64/kernel-5.10-lz4:kernel-5.10-lz4-allsyms \

endif

PRODUCT_BUILD_VENDOR_BOOT_IMAGE := false
PRODUCT_BUILD_RECOVERY_IMAGE := false

$(call inherit-product, $(SRC_TARGET_DIR)/product/generic_ramdisk.mk)
