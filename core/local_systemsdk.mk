#
# Copyright (C) 2018 The Android Open Source Project
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

ifdef BOARD_SYSTEMSDK_VERSIONS
  # Apps and jars in vendor or odm partition are forced to build against System SDK.
  _is_vendor_app :=
  ifneq (,$(filter true,$(LOCAL_VENDOR_MODULE) $(LOCAL_ODM_MODULE) $(LOCAL_PROPRIETARY_MODULE)))
    # Note: no need to check LOCAL_MODULE_PATH* since LOCAL_[VENDOR|ODM|OEM]_MODULE is already
    # set correctly before this is included.
    _is_vendor_app := true
  endif
  ifneq (,$(filter JAVA_LIBRARIES APPS,$(LOCAL_MODULE_CLASS)))
    ifndef LOCAL_SDK_VERSION
      ifeq ($(_is_vendor_app),true)
        LOCAL_SDK_VERSION := system_current
      endif
    endif
  endif

  # The selected System SDK version should be one of BOARD_SYSTEMSDK_VERSIONS, which is the list of
  # System SDK versions that the device is using.
  ifneq (,$(call has-system-sdk-version,$(LOCAL_SDK_VERSION)))
    _system_sdk_version := $(call get-numeric-sdk-version,$(LOCAL_SDK_VERSION))
    ifneq ($(_system_sdk_version),$(filter $(_system_sdk_version),$(BOARD_SYSTEMSDK_VERSIONS)))
      $(call pretty-error,Incompatible LOCAL_SDK_VERSION '$(LOCAL_SDK_VERSION)'. \
             System SDK version '$(_system_sdk_version)' is not in BOARD_SYSTEMSDK_VERSIONS: $(BOARD_SYSTEMSDK_VERSIONS))
    endif
    _system_sdk_version :=
  endif
endif
