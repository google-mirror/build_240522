# Copyright (C) 2019 The Android Open Source Project
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

.PHONY: vts-core

# =====================================================================
# Package vts related artifacts - BEGIN
# TODO(b/149249068): Clean up after all VTS tests are converted.
# =====================================================================

-include external/linux-kselftest/android/kselftest_test_list.mk
-include external/ltp/android/ltp_package_list.mk

VTS_CORE_OUT_ROOT := $(HOST_OUT)/vts-core
VTS_CORE_TESTCASES_OUT := $(VTS_CORE_OUT_ROOT)/android-vts-core/testcases

# Package vts-tradefed jars
test_suite_tools += $(HOST_OUT_JAVA_LIBRARIES)/vts-tradefed.jar \
    $(HOST_OUT_JAVA_LIBRARIES)/vts-tradefed-tests.jar

# Packaging rule for host-side Python logic, configs, and data files
host_framework_files := \
  $(call find-files-in-subdirs,test/vts,"*.py" -and -type f,.) \
  $(call find-files-in-subdirs,test/vts,"*.runner_conf" -and -type f,.) \
  $(call find-files-in-subdirs,test/vts,"*.push" -and -type f,.)
host_framework_copy_pairs := \
  $(foreach f,$(host_framework_files),\
    test/vts/$(f):$(VTS_CORE_TESTCASES_OUT)/vts/$(f))

# Packaging rule for android-vts.zip's testcases dir (DATA subdir).
# TODO(b/149249068): this should be fixed by packaging kernel tests as a standalone module like
# testcases/ltp and testcases/ltp64. Once such tests are no longer run through the python wrapper,
# we can stop packaging the kernel tests under testcases/DATA/nativetests(64)
target_native_modules := \
    $(kselftest_modules) \
    ltp \
    $(ltp_packages) \

target_native_copy_pairs := \
  $(call target-native-copy-pairs,$(target_native_modules),$(VTS_CORE_TESTCASES_OUT))

vts_copy_pairs := \
  $(call copy-many-files,$(host_framework_copy_pairs)) \
  $(call copy-many-files,$(target_native_copy_pairs))

# =====================================================================
# Package vts related artifacts - END
# =====================================================================

vts-core-zip := $(PRODUCT_OUT)/vts-core-tests.zip
# Create an artifact to include a list of test config files in vts-core.
vts-core-list-zip := $(PRODUCT_OUT)/vts-core_list.zip
# Create an artifact to include all test config files in vts-core.
vts-core-configs-zip := $(PRODUCT_OUT)/vts-core_configs.zip
my_host_shared_lib_for_vts_core := $(call copy-many-files,$(COMPATIBILITY.vts-core.HOST_SHARED_LIBRARY.FILES))
$(vts-core-zip) : .KATI_IMPLICIT_OUTPUTS := $(vts-core-list-zip) $(vts-core-configs-zip)
$(vts-core-zip) : PRIVATE_vts_core_list := $(PRODUCT_OUT)/vts-core_list
$(vts-core-zip) : PRIVATE_HOST_SHARED_LIBS := $(my_host_shared_lib_for_vts_core)
$(vts-core-zip) : $(COMPATIBILITY.vts-core.FILES) $(my_host_shared_lib_for_vts_core) $(SOONG_ZIP) $(vts_copy_pairs)
	echo $(sort $(COMPATIBILITY.vts-core.FILES)) | tr " " "\n" > $@.list
	grep $(HOST_OUT_TESTCASES) $@.list > $@-host.list || true
	grep -e .*\\.config$$ $@-host.list > $@-host-test-configs.list || true
	$(hide) for shared_lib in $(PRIVATE_HOST_SHARED_LIBS); do \
	  echo $$shared_lib >> $@-host.list; \
	done
	grep $(TARGET_OUT_TESTCASES) $@.list > $@-target.list || true
	grep -e .*\\.config$$ $@-target.list > $@-target-test-configs.list || true
	$(hide) $(SOONG_ZIP) -d -o $@ -P host -C $(HOST_OUT) -l $@-host.list -P target -C $(PRODUCT_OUT) -l $@-target.list
	$(hide) $(SOONG_ZIP) -d -o $(vts-core-configs-zip) \
	  -P host -C $(HOST_OUT) -l $@-host-test-configs.list \
	  -P target -C $(PRODUCT_OUT) -l $@-target-test-configs.list
	rm -f $(PRIVATE_vts_core_list)
	$(hide) grep -e .*\\.config$$ $@-host.list | sed s%$(HOST_OUT)%host%g > $(PRIVATE_vts_core_list)
	$(hide) grep -e .*\\.config$$ $@-target.list | sed s%$(PRODUCT_OUT)%target%g >> $(PRIVATE_vts_core_list)
	$(hide) $(SOONG_ZIP) -d -o $(vts-core-list-zip) -C $(dir $@) -f $(PRIVATE_vts_core_list)
	rm -f $@.list $@-host.list $@-target.list $@-host-test-configs.list $@-target-test-configs.list \
	  $(PRIVATE_vts_core_list)

vts-core: $(vts-core-zip)

test_suite_name := vts-core
test_suite_tradefed := vts-core-tradefed
test_suite_readme := test/vts/tools/vts-core-tradefed/README
include $(BUILD_SYSTEM)/tasks/tools/compatibility.mk
vts-core: $(compatibility_zip)

$(call dist-for-goals, vts-core, $(vts-core-zip) $(vts-core-list-zip) $(vts-core-configs-zip) $(compatibility_zip))

tests: vts-core
