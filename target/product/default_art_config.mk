#
# Copyright (C) 2020 The Android Open Source Project
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

ifeq ($(ART_APEX_JARS),)
  $(error ART_APEX_JARS is empty; cannot initialize PRODUCT_BOOT_JARS variable)
endif

# The order matters for runtime class lookup performance.
PRODUCT_BOOT_JARS := \
    $(ART_APEX_JARS) \
    framework-minus-apex \
    ext \
    com.android.i18n:core-icu4j \
    telephony-common \
    voip-common \
    ims-common

PRODUCT_UPDATABLE_BOOT_JARS := \
    com.android.conscrypt:conscrypt \
    com.android.media:updatable-media \
    com.android.mediaprovider:framework-mediaprovider \
    com.android.os.statsd:framework-statsd \
    com.android.permission:framework-permission \
    com.android.sdkext:framework-sdkextensions \
    com.android.wifi:framework-wifi \
    com.android.tethering:framework-tethering \
    com.android.ipsec:android.net.ipsec.ike

# Add the compatibility library that is needed when android.test.base
# is removed from the bootclasspath.
# Default to excluding android.test.base from the bootclasspath.
ifneq ($(REMOVE_ATB_FROM_BCP),false)
  PRODUCT_PACKAGES += framework-atb-backward-compatibility
  PRODUCT_BOOT_JARS += framework-atb-backward-compatibility
else
  PRODUCT_BOOT_JARS += android.test.base
endif

# Minimal configuration for running dex2oat (default argument values).
# PRODUCT_USES_DEFAULT_ART_CONFIG must be true to enable boot image compilation.
PRODUCT_USES_DEFAULT_ART_CONFIG := true
PRODUCT_SYSTEM_PROPERTIES += \
    dalvik.vm.image-dex2oat-Xms=64m \
    dalvik.vm.image-dex2oat-Xmx=64m \
    dalvik.vm.dex2oat-Xms=64m \
    dalvik.vm.dex2oat-Xmx=512m \

# TODO(b/172480615): Remove when platform uses ART Module prebuilts by default.
ifeq (,$(filter art_module,$(SOONG_CONFIG_NAMESPACES)))
  $(call add_soong_config_namespace,art_module)
  SOONG_CONFIG_art_module += source_build
endif
ifneq (,$(findstring .android.art,$(TARGET_BUILD_APPS)))
  # Build ART modules from source if they are listed in TARGET_BUILD_APPS.
  SOONG_CONFIG_art_module_source_build := true
else ifneq (,$(filter true,$(NATIVE_COVERAGE) $(CLANG_COVERAGE)))
  # Always build ART APEXes from source in coverage builds since the prebuilts
  # aren't built with instrumentation.
  # TODO(b/172480617): Find another solution for this.
  SOONG_CONFIG_art_module_source_build := true
else ifneq (,$(SANITIZE_TARGET)$(SANITIZE_HOST))
  # Prebuilts aren't built with sanitizers either.
  SOONG_CONFIG_art_module_source_build := true
else
  # This sets the default for building ART APEXes from source rather than
  # prebuilts (in packages/modules/ArtPrebuilt and prebuilt/module_sdk/art) in
  # all other platform builds.
  SOONG_CONFIG_art_module_source_build ?= true
endif
