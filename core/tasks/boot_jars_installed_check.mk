#
# Copyright (C) 2021 The Android Open Source Project
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

# Checks that the boot jars variables only refer to APEXes that are installed,
# i.e. present in product_MODULES. If they're not then it's likely that an APEX
# has been overridden without a corresponding entry in
# PRODUCT_BOOT_JAR_MODULE_OVERRIDES.

# Skip for unbundled builds that don't produce a platform image.
ifeq (,$(TARGET_BUILD_UNBUNDLED))

# $(1): Name of a boot jars variable with <apex>:<jar> pairs.
define check-boot-jar-modules-installed
  $(foreach pair,$($(1)), \
    $(if $(filter $(call word-colon,1,$(pair)),platform $(product_MODULES)),, \
      $(error $(pair) in $(1) refers to an APEX that is not installed)))
endef

$(call check-boot-jar-modules-installed,PRODUCT_BOOT_JARS)
$(call check-boot-jar-modules-installed,PRODUCT_UPDATABLE_BOOT_JARS)
$(call check-boot-jar-modules-installed,ART_APEX_JARS)

endif
