#
# Copyright (C) 2017 The Android Open Source Project
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

# Sets Android Go recommended default values for propreties.

# Set lowram options
PRODUCT_PROPERTY_OVERRIDES += \
     ro.config.low_ram=true \
<<<<<<< HEAD   (5c8d84 Merge "Merge empty history for sparse-6676661-L8360000065797)
     ro.lmk.critical_upgrade=true \
     ro.lmk.upgrade_pressure=40 \
     ro.lmk.downgrade_pressure=60 \
     ro.lmk.kill_heaviest_task=false \
     ro.statsd.enable=false

# set threshold to filter unused apps
PRODUCT_PROPERTY_OVERRIDES += \
     pm.dexopt.downgrade_after_inactive_days=10

=======
>>>>>>> BRANCH (a10c18 Merge "Version bump to RT11.201014.001.A1 [core/build_id.mk])

# Speed profile services and wifi-service to reduce RAM and storage.
PRODUCT_SYSTEM_SERVER_COMPILER_FILTER := speed-profile

# Always preopt extracted APKs to prevent extracting out of the APK for gms
# modules.
PRODUCT_ALWAYS_PREOPT_EXTRACTED_APK := true

# Use a profile based boot image for this device. Note that this is currently a
# generic profile and not Android Go optimized.
PRODUCT_USE_PROFILE_FOR_BOOT_IMAGE := true
PRODUCT_DEX_PREOPT_BOOT_IMAGE_PROFILE_LOCATION := frameworks/base/config/boot-image-profile.txt

# set the compiler filter for shared apks to quicken.
# Rationale: speed has a lot of dex code expansion, it uses more ram and space
# compared to quicken. Using quicken for shared APKs on Go devices may save RAM.
# Note that this is a trade-off: here we trade clean pages for dirty pages,
# extra cpu and battery. That's because the quicken files will be jit-ed in all
# the processes that load of shared apk and the code cache is not shared.
# Some notable apps that will be affected by this are gms and chrome.
# b/65591595.
PRODUCT_PROPERTY_OVERRIDES += \
     pm.dexopt.shared=quicken

# Default heap sizes. Allow up to 256m for large heaps to make sure a single app
# doesn't take all of the RAM.
PRODUCT_PROPERTY_OVERRIDES += dalvik.vm.heapgrowthlimit=128m
PRODUCT_PROPERTY_OVERRIDES += dalvik.vm.heapsize=256m

# Do not generate libartd.
PRODUCT_ART_TARGET_INCLUDE_DEBUG_BUILD := false

<<<<<<< HEAD   (5c8d84 Merge "Merge empty history for sparse-6676661-L8360000065797)
=======
# Do not spin up a separate process for the network stack on go devices, use an in-process APK.
PRODUCT_PACKAGES += InProcessNetworkStack
PRODUCT_PACKAGES += CellBroadcastAppPlatform
PRODUCT_PACKAGES += CellBroadcastServiceModulePlatform
PRODUCT_PACKAGES += com.android.tethering.inprocess

>>>>>>> BRANCH (a10c18 Merge "Version bump to RT11.201014.001.A1 [core/build_id.mk])
# Strip the local variable table and the local variable type table to reduce
# the size of the system image. This has no bearing on stack traces, but will
# leave less information available via JDWP.
PRODUCT_MINIMIZE_JAVA_DEBUG_INFO := true
<<<<<<< HEAD   (5c8d84 Merge "Merge empty history for sparse-6676661-L8360000065797)
=======

# Disable Scudo outside of eng builds to save RAM.
ifneq (,$(filter eng, $(TARGET_BUILD_VARIANT)))
  PRODUCT_DISABLE_SCUDO := true
endif

# Add the system properties.
TARGET_SYSTEM_PROP += \
    build/make/target/board/go_defaults_common.prop

# use the go specific handheld_core_hardware.xml from frameworks
PRODUCT_COPY_FILES += \
    frameworks/native/data/etc/go_handheld_core_hardware.xml:$(TARGET_COPY_OUT_VENDOR)/etc/permissions/handheld_core_hardware.xml

# Dedupe VNDK libraries with identical core variants.
TARGET_VNDK_USE_CORE_VARIANT := true
>>>>>>> BRANCH (a10c18 Merge "Version bump to RT11.201014.001.A1 [core/build_id.mk])
