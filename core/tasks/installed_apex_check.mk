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

# Check that the entries on PRODUCT_INSTALL_APEXES are correct, and that all
# APEX modules getting installed are listed in it.

# Skip for unbundled builds that don't produce a platform image.
ifeq (,$(TARGET_BUILD_UNBUNDLED))

_not_apexes := \
  $(foreach m,$(PRODUCT_INSTALL_APEXES),\
    $(if $(ALL_MODULES.$(m).IS_APEX),,$(m)))

$(call maybe-print-list-and-error,$(_not_apexes),\
The following entries on PRODUCT_INSTALL_APEXES are not recognized as APEX modules)
_not_apexes :=

_not_installed := \
  $(foreach m,$(PRODUCT_INSTALL_APEXES),\
    $(if $(ALL_MODULES.$(m).INSTALLED),,$(m)))

$(call maybe-print-list-and-error,$(_not_installed),\
The following modules on PRODUCT_INSTALL_APEXES are not installed)
_not_installed :=

_product_apexes := \
  $(foreach m,$(PRODUCT_PACKAGES),\
    $(and $(ALL_MODULES.$(m).IS_APEX),$(ALL_MODULES.$(m).INSTALLED),$(m)))

$(call maybe-print-list-and-error,$(_product_apexes),\
The following APEX modules are installed and should be listed in PRODUCT_INSTALL_APEXES rather than PRODUCT_PACKAGES)
_product_apexes :=

_other_installed_apexes := \
  $(filter-out $(PRODUCT_INSTALL_APEXES),\
    $(foreach m,$(product_MODULES),\
      $(and $(ALL_MODULES.$(m).IS_APEX),$(ALL_MODULES.$(m).INSTALLED),$(m))))

# Strip the special VNDK APEXes, which are installable by default and don't need
# to be listed on PRODUCT_INSTALL_APEXES. See comment in the
# apexBundle.installable() function in build/soong/apex/apex.go for further details.
_other_installed_apexes := $(filter-out com.android.vndk.%,$(_other_installed_apexes))

$(call maybe-print-list-and-error,$(_other_installed_apexes),\
The following APEX modules are installed but not listed in PRODUCT_INSTALL_APEXES)
_other_installed_apexes :=

endif
