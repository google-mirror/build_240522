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

# Provides a functioning ART environment without Android frameworks

ifeq ($(TARGET_CORE_JARS),)
$(error TARGET_CORE_JARS is empty; cannot update PRODUCT_PACKAGES variable)
endif

# Minimal boot classpath. This should be a subset of PRODUCT_BOOT_JARS, and equivalent to
# TARGET_CORE_JARS.
PRODUCT_PACKAGES += \
    $(TARGET_CORE_JARS)

# Additional mixins to the boot classpath.
PRODUCT_PACKAGES += \
    android.test.base \

# Why are we pulling in ext, which is frameworks/base, depending on tagsoup and nist-sip?
PRODUCT_PACKAGES += \
    ext \

# Why are we pulling in expat, which is used in frameworks, only, it seem?
PRODUCT_PACKAGES += \
    libexpat \

# Libcore.
PRODUCT_PACKAGES += \
    libjavacore \
    libopenjdk \

# Libcore ICU. TODO: Try to figure out if/why we need them explicitly.
PRODUCT_PACKAGES += \
    libicui18n \
    libicuuc \

<<<<<<< HEAD   (11d6ae Merge "Merge empty history for sparse-8121823-L3120000095288)
# ART.
PRODUCT_PACKAGES += art-runtime
# ART/dex helpers.
PRODUCT_PACKAGES += art-tools
=======
ifeq (true,$(art_target_include_debug_build))
  PRODUCT_PACKAGES += com.android.art.debug
  apex_test_module := art-check-debug-apex-gen-fakebin
else
  PRODUCT_PACKAGES += com.android.art
  apex_test_module := art-check-release-apex-gen-fakebin
endif

ifeq (true,$(call soong_config_get,art_module,source_build))
  PRODUCT_HOST_PACKAGES += $(apex_test_module)
endif

art_target_include_debug_build :=
apex_test_module :=
>>>>>>> BRANCH (244bfb Merge "Version bump to TKB1.220323.002.A1 [core/build_id.mk])

# Certificates.
PRODUCT_PACKAGES += \
    cacerts \

PRODUCT_PACKAGES += \
    hiddenapi-package-whitelist.xml \

<<<<<<< HEAD   (11d6ae Merge "Merge empty history for sparse-8121823-L3120000095288)
PRODUCT_SYSTEM_DEFAULT_PROPERTIES += \
    dalvik.vm.image-dex2oat-Xms=64m \
    dalvik.vm.image-dex2oat-Xmx=64m \
    dalvik.vm.dex2oat-Xms=64m \
    dalvik.vm.dex2oat-Xmx=512m \
=======
ifeq (,$(TARGET_BUILD_UNBUNDLED))
  # Don't depend on the framework boot image profile in unbundled builds where
  # frameworks/base may not be present.
  # TODO(b/179900989): We may not need this check once we stop using full
  # platform products on the thin ART manifest branch.
  PRODUCT_DEX_PREOPT_BOOT_IMAGE_PROFILE_LOCATION += frameworks/base/boot/boot-image-profile.txt
endif

# The dalvik.vm.dexopt.thermal-cutoff property must contain one of the values
# listed here:
#
# https://source.android.com/devices/architecture/hidl/thermal-mitigation#thermal-api
#
# If the thermal status of the device reaches or exceeds the value set here
# background dexopt will be terminated and rescheduled using an exponential
# backoff polcy.
#
# The thermal cutoff value is currently set to THERMAL_STATUS_MODERATE.
PRODUCT_SYSTEM_PROPERTIES += \
>>>>>>> BRANCH (244bfb Merge "Version bump to TKB1.220323.002.A1 [core/build_id.mk])
    dalvik.vm.usejit=true \
    dalvik.vm.usejitprofiles=true \
    dalvik.vm.dexopt.secondary=true \
    dalvik.vm.appimageformat=lz4

PRODUCT_PROPERTY_OVERRIDES += \
    ro.dalvik.vm.native.bridge=0

# Different dexopt types for different package update/install times.
# On eng builds, make "boot" reasons only extract for faster turnaround.
ifeq (eng,$(TARGET_BUILD_VARIANT))
    PRODUCT_SYSTEM_DEFAULT_PROPERTIES += \
        pm.dexopt.first-boot=extract \
        pm.dexopt.boot=extract
else
    PRODUCT_SYSTEM_DEFAULT_PROPERTIES += \
        pm.dexopt.first-boot=quicken \
        pm.dexopt.boot=verify
endif

# The install filter is speed-profile in order to enable the use of
# profiles from the dex metadata files. Note that if a profile is not provided
# or if it is empty speed-profile is equivalent to (quicken + empty app image).
PRODUCT_SYSTEM_DEFAULT_PROPERTIES += \
    pm.dexopt.install=speed-profile \
    pm.dexopt.bg-dexopt=speed-profile \
    pm.dexopt.ab-ota=speed-profile \
    pm.dexopt.inactive=verify \
    pm.dexopt.shared=speed

# Enable minidebuginfo generation unless overridden.
PRODUCT_SYSTEM_DEFAULT_PROPERTIES += \
    dalvik.vm.dex2oat-minidebuginfo=true
<<<<<<< HEAD   (11d6ae Merge "Merge empty history for sparse-8121823-L3120000095288)
=======

# Enable Madvising of the whole art, odex and vdex files to MADV_WILLNEED.
# The size specified here is the size limit of how much of the file
# (in bytes) is madvised.
# We madvise the whole .art file to MADV_WILLNEED with UINT_MAX limit.
# For odex and vdex files, we limit madvising to 100MB.
PRODUCT_SYSTEM_PROPERTIES += \
    dalvik.vm.madvise.vdexfile.size=104857600 \
    dalvik.vm.madvise.odexfile.size=104857600 \
    dalvik.vm.madvise.artfile.size=4294967295
>>>>>>> BRANCH (244bfb Merge "Version bump to TKB1.220323.002.A1 [core/build_id.mk])
