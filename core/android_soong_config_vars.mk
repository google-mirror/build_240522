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

# This file defines the Soong Config Variable namespace ANDROID, and also any
# variables in that namespace.

# The expectation is that no vendor should be using the ANDROID namespace. This
# check ensures that we don't collide with any existing vendor usage.

ifdef SOONG_CONFIG_ANDROID
$(error The Soong config namespace ANDROID is reserved.)
endif

$(call add_soong_config_namespace,ANDROID)

# Add variables to the namespace below:

$(call add_soong_config_var,ANDROID,TARGET_DYNAMIC_64_32_MEDIASERVER)
$(call add_soong_config_var,ANDROID,TARGET_DYNAMIC_64_32_DRMSERVER)
$(call add_soong_config_var,ANDROID,TARGET_ENABLE_MEDIADRM_64)
$(call add_soong_config_var,ANDROID,BOARD_USES_ODMIMAGE)
$(call add_soong_config_var,ANDROID,BOARD_USES_RECOVERY_AS_BOOT)
$(call add_soong_config_var,ANDROID,CHECK_DEV_TYPE_VIOLATIONS)
$(call add_soong_config_var,ANDROID,PRODUCT_INSTALL_DEBUG_POLICY_TO_SYSTEM_EXT)

# Default behavior for the tree wrt building modules or using prebuilts. This
# can always be overridden by setting the environment variable
# MODULE_BUILD_FROM_SOURCE.
BRANCH_DEFAULT_MODULE_BUILD_FROM_SOURCE := $(RELEASE_DEFAULT_MODULE_BUILD_FROM_SOURCE)
# TODO(b/301454934): The value from build flag is set to empty when use `False`
# The condition below can be removed after the issue get sorted.
ifeq (,$(BRANCH_DEFAULT_MODULE_BUILD_FROM_SOURCE))
  BRANCH_DEFAULT_MODULE_BUILD_FROM_SOURCE := false
endif

ifneq (,$(MODULE_BUILD_FROM_SOURCE))
  # Keep an explicit setting.
else ifeq (,$(filter docs sdk win_sdk sdk_addon,$(MAKECMDGOALS)))
  MODULE_BUILD_FROM_SOURCE := true
else ifneq (,$(PRODUCT_MODULE_BUILD_FROM_SOURCE))
  # Let products override the branch default.
  MODULE_BUILD_FROM_SOURCE := $(PRODUCT_MODULE_BUILD_FROM_SOURCE)
else
  MODULE_BUILD_FROM_SOURCE := $(BRANCH_DEFAULT_MODULE_BUILD_FROM_SOURCE)
endif

ifneq (,$(ART_MODULE_BUILD_FROM_SOURCE))
  # Keep an explicit setting.
else ifneq (,$(findstring .android.art,$(TARGET_BUILD_APPS)))
  # Build ART modules from source if they are listed in TARGET_BUILD_APPS.
  ART_MODULE_BUILD_FROM_SOURCE := true
else
  # Do the same as other modules by default.
  ART_MODULE_BUILD_FROM_SOURCE := $(MODULE_BUILD_FROM_SOURCE)
endif

$(call soong_config_set,art_module,source_build,$(ART_MODULE_BUILD_FROM_SOURCE))
ifdef ART_DEBUG_OPT_FLAG
$(call soong_config_set,art_module,art_debug_opt_flag,$(ART_DEBUG_OPT_FLAG))
endif

ifdef TARGET_BOARD_AUTO
  $(call add_soong_config_var_value, ANDROID, target_board_auto, $(TARGET_BOARD_AUTO))
endif

# Apex build mode variables
ifdef APEX_BUILD_FOR_PRE_S_DEVICES
$(call add_soong_config_var_value,ANDROID,library_linking_strategy,prefer_static)
else
ifdef KEEP_APEX_INHERIT
$(call add_soong_config_var_value,ANDROID,library_linking_strategy,prefer_static)
endif
endif

ifeq (true,$(MODULE_BUILD_FROM_SOURCE))
$(call add_soong_config_var_value,ANDROID,module_build_from_source,true)
endif

# Messaging app vars
ifeq (eng,$(TARGET_BUILD_VARIANT))
$(call soong_config_set,messaging,build_variant_eng,true)
endif

# Enable SystemUI optimizations by default unless explicitly set.
SYSTEMUI_OPTIMIZE_JAVA ?= true
$(call add_soong_config_var,ANDROID,SYSTEMUI_OPTIMIZE_JAVA)

# Enable Compose in SystemUI by default.
SYSTEMUI_USE_COMPOSE ?= true
$(call add_soong_config_var,ANDROID,SYSTEMUI_USE_COMPOSE)

ifdef PRODUCT_AVF_ENABLED
$(call add_soong_config_var_value,ANDROID,avf_enabled,$(PRODUCT_AVF_ENABLED))
endif

ifdef PRODUCT_AVF_MICRODROID_GUEST_GKI_VERSION
$(call add_soong_config_var_value,ANDROID,avf_microdroid_guest_gki_version,$(PRODUCT_AVF_MICRODROID_GUEST_GKI_VERSION))
endif

ifdef PRODUCT_MEMCG_V2_FORCE_ENABLED
$(call add_soong_config_var_value,ANDROID,memcg_v2_force_enabled,$(PRODUCT_MEMCG_V2_FORCE_ENABLED))
endif

ifdef PRODUCT_CGROUP_V2_SYS_APP_ISOLATION_ENABLED
$(call add_soong_config_var_value,ANDROID,cgroup_v2_sys_app_isolation,$(PRODUCT_CGROUP_V2_SYS_APP_ISOLATION_ENABLED))
endif

$(call add_soong_config_var_value,ANDROID,release_avf_allow_preinstalled_apps,$(RELEASE_AVF_ALLOW_PREINSTALLED_APPS))
$(call add_soong_config_var_value,ANDROID,release_avf_enable_device_assignment,$(RELEASE_AVF_ENABLE_DEVICE_ASSIGNMENT))
$(call add_soong_config_var_value,ANDROID,release_avf_enable_dice_changes,$(RELEASE_AVF_ENABLE_DICE_CHANGES))
$(call add_soong_config_var_value,ANDROID,release_avf_enable_llpvm_changes,$(RELEASE_AVF_ENABLE_LLPVM_CHANGES))
$(call add_soong_config_var_value,ANDROID,release_avf_enable_multi_tenant_microdroid_vm,$(RELEASE_AVF_ENABLE_MULTI_TENANT_MICRODROID_VM))
$(call add_soong_config_var_value,ANDROID,release_avf_enable_remote_attestation,$(RELEASE_AVF_ENABLE_REMOTE_ATTESTATION))
$(call add_soong_config_var_value,ANDROID,release_avf_enable_vendor_modules,$(RELEASE_AVF_ENABLE_VENDOR_MODULES))
$(call add_soong_config_var_value,ANDROID,release_avf_enable_virt_cpufreq,$(RELEASE_AVF_ENABLE_VIRT_CPUFREQ))
$(call add_soong_config_var_value,ANDROID,release_avf_microdroid_kernel_version,$(RELEASE_AVF_MICRODROID_KERNEL_VERSION))
$(call add_soong_config_var_value,ANDROID,release_avf_support_custom_vm_with_paravirtualized_devices,$(RELEASE_AVF_SUPPORT_CUSTOM_VM_WITH_PARAVIRTUALIZED_DEVICES))

$(call add_soong_config_var_value,ANDROID,release_binder_death_recipient_weak_from_jni,$(RELEASE_BINDER_DEATH_RECIPIENT_WEAK_FROM_JNI))

$(call add_soong_config_var_value,ANDROID,release_package_libandroid_runtime_punch_holes,$(RELEASE_PACKAGE_LIBANDROID_RUNTIME_PUNCH_HOLES))

$(call add_soong_config_var_value,ANDROID,release_selinux_data_data_ignore,$(RELEASE_SELINUX_DATA_DATA_IGNORE))

# Enable system_server optimizations by default unless explicitly set or if
# there may be dependent runtime jars.
# TODO(b/240588226): Remove the off-by-default exceptions after handling
# system_server jars automatically w/ R8.
ifeq (true,$(PRODUCT_BROKEN_SUBOPTIMAL_ORDER_OF_SYSTEM_SERVER_JARS))
  # If system_server jar ordering is broken, don't assume services.jar can be
  # safely optimized in isolation, as there may be dependent jars.
  SYSTEM_OPTIMIZE_JAVA ?= false
else ifneq (platform:services,$(lastword $(PRODUCT_SYSTEM_SERVER_JARS)))
  # If services is not the final jar in the dependency ordering, don't assume
  # it can be safely optimized in isolation, as there may be dependent jars.
  SYSTEM_OPTIMIZE_JAVA ?= false
else
  SYSTEM_OPTIMIZE_JAVA ?= true
endif

ifeq (true,$(FULL_SYSTEM_OPTIMIZE_JAVA))
  SYSTEM_OPTIMIZE_JAVA := true
endif

$(call add_soong_config_var,ANDROID,SYSTEM_OPTIMIZE_JAVA)
$(call add_soong_config_var,ANDROID,FULL_SYSTEM_OPTIMIZE_JAVA)

# TODO(b/319697968): Remove this build flag support when metalava fully supports flagged api
$(call soong_config_set,ANDROID,release_hidden_api_exportable_stubs,$(RELEASE_HIDDEN_API_EXPORTABLE_STUBS))

# Check for SupplementalApi module.
ifeq ($(wildcard packages/modules/SupplementalApi),)
$(call add_soong_config_var_value,ANDROID,include_nonpublic_framework_api,false)
else
$(call add_soong_config_var_value,ANDROID,include_nonpublic_framework_api,true)
endif

# Add crashrecovery build flag to soong
$(call soong_config_set,ANDROID,release_crashrecovery_module,$(RELEASE_CRASHRECOVERY_MODULE))
ifeq (true,$(RELEASE_CRASHRECOVERY_FILE_MOVE))
  $(call soong_config_set,ANDROID,crashrecovery_files_in_module,true)
  $(call soong_config_set,ANDROID,crashrecovery_files_in_platform,false)
else
  $(call soong_config_set,ANDROID,crashrecovery_files_in_module,false)
  $(call soong_config_set,ANDROID,crashrecovery_files_in_platform,true)
endif
# Weirdly required because platform_bootclasspath is using AUTO namespace
$(call soong_config_set,AUTO,release_crashrecovery_module,$(RELEASE_CRASHRECOVERY_MODULE))
