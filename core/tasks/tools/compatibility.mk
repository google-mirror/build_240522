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
  $(HOST_OUT_JAVA_LIBRARIES)/loganalysis.jar \
  $(HOST_OUT_JAVA_LIBRARIES)/compatibility-host-util.jar \
  $(HOST_OUT_JAVA_LIBRARIES)/compatibility-host-util-tests.jar \
  $(HOST_OUT_JAVA_LIBRARIES)/compatibility-common-util-tests.jar \
  $(HOST_OUT_JAVA_LIBRARIES)/compatibility-tradefed-tests.jar \
  $(HOST_OUT_JAVA_LIBRARIES)/host-libprotobuf-java-full.jar \
  $(HOST_OUT_JAVA_LIBRARIES)/$(test_suite_tradefed).jar \
  $(HOST_OUT_JAVA_LIBRARIES)/$(test_suite_tradefed)-tests.jar \
  $(HOST_OUT_EXECUTABLES)/$(test_suite_tradefed) \
  $(test_suite_readme)

test_tools += $(test_suite_tools)

<<<<<<< HEAD   (123cec Merge "Merge empty history for sparse-8898769-L7910000095637)
=======
# The JDK to package into the test suite zip file.  Always package the linux JDK.
test_suite_jdk_dir := $(ANDROID_JAVA_HOME)/../linux-x86
test_suite_jdk := $(call intermediates-dir-for,PACKAGING,$(test_suite_name)_jdk,HOST)/jdk.zip
$(test_suite_jdk): PRIVATE_JDK_DIR := $(test_suite_jdk_dir)
$(test_suite_jdk): PRIVATE_SUBDIR := $(test_suite_subdir)
$(test_suite_jdk): $(shell find $(test_suite_jdk_dir) -type f | sort)
$(test_suite_jdk): $(SOONG_ZIP)
	$(SOONG_ZIP) -o $@ -P $(PRIVATE_SUBDIR)/jdk -C $(PRIVATE_JDK_DIR) -D $(PRIVATE_JDK_DIR)

$(call declare-license-metadata,$(test_suite_jdk),SPDX-license-identifier-GPL-2.0-with-classpath-exception,restricted,\
  $(test_suite_jdk_dir)/legal/java.base/LICENSE,JDK,prebuilts/jdk/$(notdir $(patsubst %/,%,$(dir $(test_suite_jdk_dir)))))

# Copy license metadata
$(call declare-copy-target-license-metadata,$(out_dir)/$(notdir $(test_suite_jdk)),$(test_suite_jdk))
$(foreach t,$(test_tools) $(test_suite_prebuilt_tools),\
  $(eval _dst := $(out_dir)/tools/$(notdir $(t)))\
  $(if $(strip $(ALL_TARGETS.$(t).META_LIC)),\
    $(call declare-copy-target-license-metadata,$(_dst),$(t)),\
    $(warning $(t) has no license metadata)\
  )\
)
test_copied_tools := $(foreach t,$(test_tools) $(test_suite_prebuilt_tools), $(out_dir)/tools/$(notdir $(t))) $(out_dir)/$(notdir $(test_suite_jdk))


# Include host shared libraries
host_shared_libs := $(call copy-many-files, $(COMPATIBILITY.$(test_suite_name).HOST_SHARED_LIBRARY.FILES))

$(if $(strip $(host_shared_libs)),\
  $(foreach p,$(COMPATIBILITY.$(test_suite_name).HOST_SHARED_LIBRARY.FILES),\
    $(eval _src := $(call word-colon,1,$(p)))\
    $(eval _dst := $(call word-colon,2,$(p)))\
    $(if $(strip $(ALL_TARGETS.$(_src).META_LIC)),\
      $(call declare-copy-target-license-metadata,$(_dst),$(_src)),\
      $(warning $(_src) has no license metadata for $(_dst))\
    )\
  )\
)

compatibility_zip_deps := \
  $(test_artifacts) \
  $(test_tools) \
  $(test_suite_prebuilt_tools) \
  $(test_suite_dynamic_config) \
  $(test_suite_jdk) \
  $(MERGE_ZIPS) \
  $(SOONG_ZIP) \
  $(host_shared_libs) \
  $(test_suite_extra_deps) \

compatibility_zip_resources := $(out_dir)/tools $(out_dir)/testcases $(out_dir)/lib $(out_dir)/lib64

# Test Suite NOTICE files
test_suite_notice_txt := $(out_dir)/NOTICE.txt
test_suite_notice_html := $(out_dir)/NOTICE.html

compatibility_zip_deps += $(test_suite_notice_txt)
compatibility_zip_resources += $(test_suite_notice_txt)

compatibility_tests_list_zip := $(out_dir)-tests_list.zip

>>>>>>> BRANCH (b4676c Merge "Version bump to TKB1.220911.001.A1 [core/build_id.mk])
compatibility_zip := $(out_dir).zip
$(compatibility_zip): PRIVATE_NAME := android-$(test_suite_name)
$(compatibility_zip): PRIVATE_OUT_DIR := $(out_dir)
$(compatibility_zip): PRIVATE_TOOLS := $(test_tools) $(test_suite_prebuilt_tools)
$(compatibility_zip): PRIVATE_SUITE_NAME := $(test_suite_name)
$(compatibility_zip): PRIVATE_DYNAMIC_CONFIG := $(test_suite_dynamic_config)
$(compatibility_zip): $(test_artifacts) $(test_tools) $(test_suite_prebuilt_tools) $(test_suite_dynamic_config) $(SOONG_ZIP) | $(ADB) $(ACP)
# Make dir structure
	$(hide) mkdir -p $(PRIVATE_OUT_DIR)/tools $(PRIVATE_OUT_DIR)/testcases
# Copy tools
<<<<<<< HEAD   (123cec Merge "Merge empty history for sparse-8898769-L7910000095637)
	$(hide) $(ACP) -fp $(PRIVATE_TOOLS) $(PRIVATE_OUT_DIR)/tools
	$(if $(PRIVATE_DYNAMIC_CONFIG),$(hide) $(ACP) -fp $(PRIVATE_DYNAMIC_CONFIG) $(PRIVATE_OUT_DIR)/testcases/$(PRIVATE_SUITE_NAME).dynamic)
	$(hide) find $(dir $@)/$(PRIVATE_NAME) | sort >$@.list
	$(hide) $(SOONG_ZIP) -d -o $@ -C $(dir $@) -l $@.list
=======
	cp $(PRIVATE_TOOLS) $(PRIVATE_OUT_DIR)/tools
	$(if $(PRIVATE_DYNAMIC_CONFIG),$(hide) cp $(PRIVATE_DYNAMIC_CONFIG) $(PRIVATE_OUT_DIR)/testcases/$(PRIVATE_SUITE_NAME).dynamic)
	find $(PRIVATE_RESOURCES) | sort >$@.list
	$(SOONG_ZIP) -d -o $@.tmp -C $(dir $@) -l $@.list
	$(MERGE_ZIPS) $@ $@.tmp $(PRIVATE_JDK)
	rm -f $@.tmp
# Build a list of tests
	rm -f $(PRIVATE_tests_list)
	$(hide) grep -e .*\\.config$$ $@.list | sed s%$(PRIVATE_OUT_DIR)/testcases/%%g > $(PRIVATE_tests_list)
	$(SOONG_ZIP) -d -o $(PRIVATE_tests_list_zip) -j -f $(PRIVATE_tests_list)
	rm -f $(PRIVATE_tests_list)

$(call declare-0p-target,$(compatibility_tests_list_zip),)

$(call declare-1p-container,$(compatibility_zip),)
$(call declare-container-license-deps,$(compatibility_zip),$(compatibility_zip_deps) $(test_copied_tools), $(out_dir)/:/)

$(eval $(call html-notice-rule,$(test_suite_notice_html),"Test suites","Notices for files contained in the test suites filesystem image:",$(compatibility_zip),$(compatibility_zip)))
$(eval $(call text-notice-rule,$(test_suite_notice_txt),"Test suites","Notices for files contained in the test suites filesystem image:",$(compatibility_zip),$(compatibility_zip)))

$(call declare-0p-target,$(test_suite_notice_html))
$(call declare-0p-target,$(test_suite_notice_txt))

$(call declare-1p-copy-files,$(test_suite_dynamic_config),)
$(call declare-1p-copy-files,$(test_suite_prebuilt_tools),)
>>>>>>> BRANCH (b4676c Merge "Version bump to TKB1.220911.001.A1 [core/build_id.mk])

# Reset all input variables
test_suite_name :=
test_suite_tradefed :=
test_suite_dynamic_config :=
test_suite_readme :=
test_suite_prebuilt_tools :=
test_suite_tools :=
