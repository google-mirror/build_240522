# Copyright (C) 2015 The Android Open Source Project
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

test_suite_name := cts
test_suite_tradefed := cts-tradefed
test_suite_dynamic_config := cts/tools/cts-tradefed/DynamicConfig.xml
test_suite_readme := cts/tools/cts-tradefed/README

include $(BUILD_SYSTEM)/tasks/tools/compatibility.mk

.PHONY: cts
cts: $(compatibility_zip)
$(call dist-for-goals, cts, $(compatibility_zip))

.PHONY: cts_v2
cts_v2: cts

# Creates a "cts-verifier" directory that will contain:
#
# 1. Out directory with a "android-cts-verifier" containing the CTS Verifier
#    and other binaries it needs.
#
# 2. Zipped version of the android-cts-verifier directory to be included with
#    the build distribution.
##
cts-dir := $(HOST_OUT)/cts-verifier
verifier-dir-name := android-cts-verifier
verifier-dir := $(cts-dir)/$(verifier-dir-name)
verifier-zip-name := $(verifier-dir-name).zip
verifier-zip := $(cts-dir)/$(verifier-zip-name)

cts : $(verifier-zip)
$(verifier-zip): PRIVATE_DIR := $(cts-dir)
$(verifier-zip): $(SOONG_ANDROID_CTS_VERIFIER_ZIP)
	rm -rf $(PRIVATE_DIR)
	mkdir -p $(PRIVATE_DIR)
	unzip -q -d $(PRIVATE_DIR) $<
	$(copy-file-to-target)
