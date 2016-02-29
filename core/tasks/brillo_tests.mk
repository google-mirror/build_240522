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

# Generates a zip file that includes brillo test modules
ifneq ($(filter $(MAKECMDGOALS),brillo_tests),)

include $(CLEAR_VARS)
LOCAL_MODULE := brillo_whitelist
LOCAL_PICKUP_FILES := device/generic/brillo/tests.txt
include $(BUILD_PHONY_PACKAGE)

# Include the whitelist itself and all files from it.
my_modules := brillo_whitelist \
  $(shell sed -n -e 's/\(^[^#].*\),.*/\1/p' device/generic/brillo/tests.txt)

my_package_name := brillo_tests

include $(BUILD_SYSTEM)/tasks/tools/package-modules.mk

.PHONY: brillo_tests
brillo_tests : $(my_package_zip)

name := $(TARGET_PRODUCT)-brillo-tests-$(FILE_NAME_TAG)
$(call dist-for-goals, brillo_tests, $(my_package_zip):$(name).zip)
endif
