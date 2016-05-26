#
# Copyright (C) 2008 The Android Open Source Project
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

ifneq ($(LOCAL_MODULE)$(LOCAL_MODULE_CLASS),)
$(error $(LOCAL_PATH): LOCAL_MODULE or LOCAL_MODULE_CLASS not needed by \
  BUILD_MULTI_PREBUILT, use BUILD_PREBUILT instead!)
endif

ifdef LOCAL_PREBUILT_LIBS
tmp_LOCAL_STRIP_MODULE := $(LOCAL_STRIP_MODULE)
tmp_LOCAL_PREBUILT_LIBS := $(LOCAL_PREBUILT_LIBS)
LOCAL_PREBUILT_LIBS :=

# static libs
LOCAL_MODULE_CLASS := STATIC_LIBRARIES
prebuilt_list := $(filter %.a,$(tmp_LOCAL_PREBUILT_LIBS))
OVERRIDE_BUILT_MODULE_PATH :=
LOCAL_UNINSTALLABLE_MODULE := true
LOCAL_STRIP_MODULE :=
include $(BUILD_SYSTEM)/multi_prebuilt_internal.mk

# shared libs
LOCAL_MODULE_CLASS := SHARED_LIBRARIES
prebuilt_list := $(filter-out %.a,$(tmp_LOCAL_PREBUILT_LIBS))
OVERRIDE_BUILT_MODULE_PATH := $($(if $(prebuilt_is_host),HOST,TARGET)_OUT_INTERMEDIATE_LIBRARIES)
LOCAL_UNINSTALLABLE_MODULE :=
LOCAL_STRIP_MODULE := $(tmp_LOCAL_STRIP_MODULE)
include $(BUILD_SYSTEM)/multi_prebuilt_internal.mk

tmp_LOCAL_STRIP_MODULE :=
tmp_LOCAL_PREBUILT_LIBS :=

else ifdef LOCAL_PREBUILT_EXECUTABLES
LOCAL_MODULE_CLASS := EXECUTABLES
prebuilt_list := $(LOCAL_PREBUILT_EXECUTABLES)
LOCAL_PREBUILT_EXECUTABLES :=
include $(BUILD_SYSTEM)/multi_prebuilt_internal.mk

else ifdef LOCAL_PREBUILT_JAVA_LIBRARIES
LOCAL_MODULE_CLASS := JAVA_LIBRARIES
prebuilt_list := $(LOCAL_PREBUILT_JAVA_LIBRARIES)
LOCAL_PREBUILT_JAVA_LIBRARIES :=
LOCAL_BUILT_MODULE_STEM := javalib.jar
include $(BUILD_SYSTEM)/multi_prebuilt_internal.mk

else ifdef LOCAL_PREBUILT_STATIC_JAVA_LIBRARIES
LOCAL_MODULE_CLASS := JAVA_LIBRARIES
prebuilt_list := $(LOCAL_PREBUILT_STATIC_JAVA_LIBRARIES)
LOCAL_PREBUILT_STATIC_JAVA_LIBRARIES :=
LOCAL_BUILT_MODULE_STEM := javalib.jar
LOCAL_UNINSTALLABLE_MODULE := true
include $(BUILD_SYSTEM)/multi_prebuilt_internal.mk

endif

prebuilt_list :=

