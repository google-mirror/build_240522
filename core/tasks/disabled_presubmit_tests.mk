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

# This build rule creates a build artifact that contains a list of disabled
# presubmit tests specified in ALL_DISABLED_PRESUBMIT_TESTS, which are
# configuredfrom 2 sources:
# 1. Tests added to ALL_DISABLED_PRESUBMIT_TESTS variable directly, which are
#    usually integration tests that don't have build target.
# 2. Tests with module configured with LOCAL_PRESUBMIT_DISABLED set to true.

.PHONY: disabled_presubmit_tests

disabled-presubmit-tests-zip := $(PRODUCT_OUT)/disabled-presubmit-tests.zip
$(disabled-presubmit-tests-zip) : PRIVATE_all_disabled_presubmit_tests := $(ALL_DISABLED_PRESUBMIT_TESTS)
$(disabled-presubmit-tests-zip) : PRIVATE_disabled_presubmit_tests_file := $(PRODUCT_OUT)/disabled-presubmit-tests
$(disabled-presubmit-tests-zip) :
	@echo Compiling a list of disabled presubmit tests: $@
	$(hide) rm -f $(PRIVATE_disabled_presubmit_tests_file)
	$(hide) echo $(sort $(PRIVATE_all_disabled_presubmit_tests)) | tr " " "\n" >> $(PRIVATE_disabled_presubmit_tests_file)
	$(hide) $(SOONG_ZIP) -d -o $(disabled-presubmit-tests-zip) -C $(dir $@) -f $(PRIVATE_disabled_presubmit_tests_file)
	$(hide) rm -f $(PRIVATE_disabled_presubmit_tests_file)

$(disabled_presubmit_tests) : $(disabled-presubmit-tests-zip)
$(call dist-for-goals, disabled_presubmit_tests, $(disabled-presubmit-tests-zip))

tests: disabled_presubmit_tests
