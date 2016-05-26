#
# Copyright (C) 2016 The Android Open Source Project
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

$(if $(filter %: :%,$(prebuilt_list)), \
  $(error $(LOCAL_PATH): Leading or trailing colons in "$(prebuilt_list)")) \
$(foreach t,$(prebuilt_list), \
  $(eval tw := $(subst :, ,$(strip $(t)))) \
  $(if $(word 3,$(tw)),$(error $(LOCAL_PATH): Bad prebuilt filename '$(t)')) \
  $(if $(word 2,$(tw)), \
    $(eval LOCAL_MODULE := $(word 1,$(tw))) \
    $(eval LOCAL_SRC_FILES := $(word 2,$(tw))) \
   , \
    $(eval LOCAL_MODULE := $(basename $(notdir $(t)))) \
    $(eval LOCAL_SRC_FILES := $(t)) \
   ) \
  $(if $(LOCAL_BUILT_MODULE_STEM),, \
    $(if $(word 2,$(tw)), \
      $(eval LOCAL_BUILT_MODULE_STEM := $(LOCAL_MODULE)$(suffix $(LOCAL_SRC_FILES))) \
     , \
      $(eval LOCAL_BUILT_MODULE_STEM := $(notdir $(LOCAL_SRC_FILES))) \
     ) \
   ) \
  $(eval LOCAL_MODULE_SUFFIX := $(suffix $(LOCAL_SRC_FILES))) \
  $(eval include $(BUILD_PREBUILT)) \
  $(eval LOCAL_BUILT_MODULE :=) \
  $(eval LOCAL_INSTALLED_MODULE :=) \
 )

