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

# Package up a compatibility test suite in a zip file.
#
# Input variables:
#   test_suite_name: the name of this test suite eg. cts
#   test_suite_tradefed: the name of this test suite's tradefed wrapper
#   test_suite_dynamic_config: the path to this test suite's dynamic configuration file
#   test_suite_readme: the path to a README file for this test suite
#   test_suite_prebuilt_tools: the set of prebuilt tools to be included directly
#                         in the 'tools' subdirectory of the test suite.
#   test_suite_tools: the set of tools for this test suite
#
# Output variables:
#   compatibility_zip: the path to the output zip file.

out_dir := $(HOST_OUT)/$(test_suite_name)/android-$(test_suite_name)
test_artifacts := $(COMPATIBILITY.$(test_suite_name).FILES)
test_tools := $(HOST_OUT_JAVA_LIBRARIES)/hosttestlib.jar \
  $(HOST_OUT_JAVA_LIBRARIES)/tradefed.jar \
  $(HOST_OUT_JAVA_LIBRARIES)/tradefed-no-fwk.jar \
  $(HOST_OUT_JAVA_LIBRARIES)/tradefed-test-framework.jar \
  $(HOST_OUT_JAVA_LIBRARIES)/loganalysis.jar \
  $(HOST_OUT_JAVA_LIBRARIES)/compatibility-host-util.jar \
  $(HOST_OUT_JAVA_LIBRARIES)/compatibility-host-util-tests.jar \
  $(HOST_OUT_JAVA_LIBRARIES)/compatibility-common-util-tests.jar \
  $(HOST_OUT_JAVA_LIBRARIES)/compatibility-tradefed-tests.jar \
  $(HOST_OUT_JAVA_LIBRARIES)/$(test_suite_tradefed).jar \
  $(HOST_OUT_JAVA_LIBRARIES)/$(test_suite_tradefed)-tests.jar \
  $(HOST_OUT_EXECUTABLES)/$(test_suite_tradefed) \
  $(test_suite_readme)

test_tools += $(test_suite_tools)

# The JDK to package into the test suite zip file.  Always package the linux JDK.
test_suite_jdk := $(ANDROID_JAVA_HOME)/../linux-x86/jdk.zip

# Include host shared libraries
host_shared_libs := $(call copy-many-files, $(COMPATIBILITY.$(test_suite_name).HOST_SHARED_LIBRARY.FILES))

compatibility_zip := $(out_dir).zip
$(compatibility_zip): PRIVATE_NAME := android-$(test_suite_name)
$(compatibility_zip): PRIVATE_OUT_DIR := $(out_dir)
$(compatibility_zip): PRIVATE_TOOLS := $(test_tools) $(test_suite_prebuilt_tools)
$(compatibility_zip): PRIVATE_SUITE_NAME := $(test_suite_name)
$(compatibility_zip): PRIVATE_DYNAMIC_CONFIG := $(test_suite_dynamic_config)
$(compatibility_zip): PRIVATE_JDK := $(test_suite_jdk)
$(compatibility_zip): $(test_artifacts) $(test_tools) $(test_suite_prebuilt_tools) $(test_suite_dynamic_config) $(host_shared_libs) $(test_suite_jdk) $(MERGE_ZIPS) $(SOONG_ZIP) $(ZIP2ZIP) | $(ADB) $(ACP)
# Make dir structure
	mkdir -p $(PRIVATE_OUT_DIR)/tools $(PRIVATE_OUT_DIR)/testcases
	rm -f $@ $@.tmp $@.jdk
	echo $(BUILD_NUMBER_FROM_FILE) > $(PRIVATE_OUT_DIR)/tools/version.txt
# Copy tools
	cp $(PRIVATE_TOOLS) $(PRIVATE_OUT_DIR)/tools
	$(if $(PRIVATE_DYNAMIC_CONFIG),cp $(PRIVATE_DYNAMIC_CONFIG) $(PRIVATE_OUT_DIR)/testcases/$(PRIVATE_SUITE_NAME).dynamic)
	find $(PRIVATE_OUT_DIR)/tools $(PRIVATE_OUT_DIR)/testcases | sort >$@.list
	$(SOONG_ZIP) -d -o $@.tmp -C $(dir $@) -l $@.list
	$(ZIP2ZIP) -o $@.jdk -i $(PRIVATE_JDK) **/*:jdk/
	$(MERGE_ZIPS) $@ $@.tmp $@.jdk
	rm -f $@.tmp $@.jdk

# Reset all input variables
test_suite_name :=
test_suite_tradefed :=
test_suite_dynamic_config :=
test_suite_readme :=
test_suite_prebuilt_tools :=
test_suite_tools :=
test_suite_jdk :=
host_shared_libs :=
