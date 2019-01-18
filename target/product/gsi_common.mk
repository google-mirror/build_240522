#
# Copyright (C) 2019 The Android Open-Source Project
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

$(call inherit-product, $(SRC_TARGET_DIR)/product/mainline_system.mk)

# GSI includes all AOSP product packages and placed under /system/product
$(call inherit-product, $(SRC_TARGET_DIR)/product/handheld_product.mk)
$(call inherit-product, $(SRC_TARGET_DIR)/product/telephony_product.mk)

# Default AOSP sounds (original source: full_base.mk)
$(call inherit-product-if-exists, frameworks/base/data/sounds/AllAudio.mk)

# Default languages (original source: full_base.mk)
$(call inherit-product, $(SRC_TARGET_DIR)/product/languages_full.mk)


# Enable mainline checking and the whitelist for GSI
PRODUCT_ENFORCE_ARTIFACT_PATH_REQUIREMENTS := true
PRODUCT_ARTIFACT_PATH_REQUIREMENT_WHITELIST := \
    system/app/messaging/messaging.apk \
    system/app/PhotoTable/oat/arm64/PhotoTable.odex \
    system/app/PhotoTable/oat/arm64/PhotoTable.vdex \
    system/app/PhotoTable/PhotoTable.apk \
    system/app/PrintRecommendationService/PrintRecommendationService.apk \
    system/app/WAPPushManager/WAPPushManager.apk \
    system/app/webview/webview.apk \
    system/bin/healthd \
    system/etc/vintf/manifest/manifest_healthd.xml \
    system/lib/libframesequence.so \
    system/lib/libgiftranscode.so \
    system/lib64/libframesequence.so \
    system/lib64/libgiftranscode.so \
    system/priv-app/Dialer/Dialer.apk \

# Exclude GSI specific files
PRODUCT_ARTIFACT_PATH_REQUIREMENT_WHITELIST += \
    system/etc/apns-conf.xml \
    system/etc/init/config/skip_mount.cfg \
    system/etc/init/healthd.rc \
    system/etc/init/init.gsi.rc \
    system/etc/libnfc-nci.conf \

# Exclude all files under system/product
PRODUCT_ARTIFACT_PATH_REQUIREMENT_WHITELIST += \
    system/product/%

# Exclude all default AOSP sounds
PRODUCT_ARTIFACT_PATH_REQUIREMENT_WHITELIST += \
    system/product/media/audio/% \


# Split selinux policy
PRODUCT_FULL_TREBLE_OVERRIDE := true

# Enable dynamic partition size
PRODUCT_USE_DYNAMIC_PARTITION_SIZE := true

# Enable A/B update
AB_OTA_UPDATER := true
AB_OTA_PARTITIONS := system
# PRODUCT_PACKAGES += \
#     update_engine \
#     update_verifier

# Needed by Pi newly launched device to pass VtsTrebleSysProp on GSI
PRODUCT_COMPATIBLE_PROPERTY_OVERRIDE := true

# GSI specific tasks on boot
PRODUCT_COPY_FILES += \
    build/make/target/product/gsi/skip_mount.cfg:system/etc/init/config/skip_mount.cfg \
    build/make/target/product/gsi/init.gsi.rc:system/etc/init/init.gsi.rc \

# Support addtional P vendor interface
PRODUCT_EXTRA_VNDK_VERSIONS := 28

# Default AOSP packages (original source: aosp_base_telephony.mk)
PRODUCT_PACKAGES += \
    messaging \

# Default AOSP packages (original source: full_base.mk)
PRODUCT_PACKAGES += \
    PhotoTable \
    WAPPushManager \

# Telephony:
#   Provide a default APN configuration
PRODUCT_COPY_FILES += \
    device/sample/etc/apns-full-conf.xml:system/etc/apns-conf.xml

# NFC:
#   Provide default libnfc-nci.conf file for devices that does not have one in
#   vendor/etc
PRODUCT_COPY_FILES += \
    device/generic/common/nfc/libnfc-nci.conf:system/etc/libnfc-nci.conf
