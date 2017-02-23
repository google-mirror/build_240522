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

PRODUCT_PROPERTY_OVERRIDES := \
    ro.ril.hsxpa=1 \
    ro.ril.gprsclass=10 \
#    ro.adb.qemud=1

PRODUCT_COPY_FILES := \
    device/generic/goldfish/data/etc/apns-conf.xml:system/etc/apns-conf.xml \
    device/generic/goldfish/camera/media_profiles.xml:system/etc/media_profiles.xml \
    frameworks/av/media/libstagefright/data/media_codecs_google_audio.xml:system/etc/media_codecs_google_audio.xml \
    frameworks/av/media/libstagefright/data/media_codecs_google_telephony.xml:system/etc/media_codecs_google_telephony.xml \
    frameworks/av/media/libstagefright/data/media_codecs_google_video.xml:system/etc/media_codecs_google_video.xml \
    device/generic/goldfish/camera/media_codecs.xml:system/etc/media_codecs.xml

# display
 PRODUCT_PROPERTY_OVERRIDES += \
     ro.sf.lcd_density=420

# Audio
PRODUCT_PROPERTY_OVERRIDES += \
     ro.config.vc_call_vol_steps=7

# HWUI cache sizes
PRODUCT_PROPERTY_OVERRIDES += \
     ro.hwui.texture_cache_size=56 \
     ro.hwui.layer_cache_size=32 \
     ro.hwui.path_cache_size=16

# Use Sdcardfs
PRODUCT_PROPERTY_OVERRIDES += \
     ro.sys.sdcardfs=1

PRODUCT_COPY_FILES += \
    device/google/marlin/init.common.rc:root/init.marlin.rc \
    device/google/marlin/init.common.usb.rc:root/init.marlin.usb.rc \
    device/google/marlin/fstab.common:root/fstab.marlin \
    device/google/marlin/ueventd.common.rc:root/ueventd.marlin.rc \
    device/google/marlin/init.recovery.common.rc:root/init.recovery.marlin.rc

# Sensor hub init script
PRODUCT_COPY_FILES += \
    device/google/marlin/init.common.nanohub.rc:root/init.marlin.nanohub.rc


AB_OTA_PARTITIONS := boot system
AB_OTA_UPDATER := true
PRODUCT_PACKAGES := update_engine update_verifier


# Adjust the Dalvik heap to be appropriate for a tablet.
$(call inherit-product-if-exists, frameworks/base/build/tablet-dalvik-heap.mk)
$(call inherit-product-if-exists, frameworks/native/build/tablet-dalvik-heap.mk)

PRODUCT_SYSTEM_VERITY_PARTITION := /dev/block/bootdevice/by-name/system
PRODUCT_VENDOR_VERITY_PARTITION := /dev/block/bootdevice/by-name/vendor
$(call inherit-product, build/target/product/verity.mk)
