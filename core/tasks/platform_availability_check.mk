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

# Check whether there is any module that isn't available for platform
# is installed to the platform.

_modules_not_available_for_platform := \
$(strip $(foreach m,$(product_MODULES),\
  $(if $(filter-out FAKE,$(ALL_MODULES.$(m).CLASS)),\
    $(if $(filter true,$(ALL_MODULES.$(m).NOT_AVAILABLE_FOR_PLATFORM)),\
      $(m)))))

ifneq ($(_modules_not_available_for_platform),)
$(info Fllowing modules are unavailable for platform but are requested to be installed:)
$(foreach m,$(sort $(_modules_not_available_for_platform)),$(info $(m):$(ALL_MODULES.$(m).PATH)))
$(error Build failed)
endif
