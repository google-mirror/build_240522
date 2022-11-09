# Copyright (C) 2022 The Android Open Source Project
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

.PHONY: art-host-cross-tests

intermediates_dir := $(call intermediates-dir-for,PACKAGING,art-host-cross-tests)
art_host_cross_tests_zip := $(PRODUCT_OUT)/art-host-cross-tests.zip
# Get the hostside libraries to be packaged in the test zip. Unlike
# device-tests.mk or general-tests.mk, the files are not copied to the
# testcases directory.
my_host_shared_lib_for_art_host_cross_tests := $(foreach f,$(COMPATIBILITY.art-host-cross-tests.HOST_CROSS_SHARED_LIBRARY.FILES),$(strip \
    $(eval _cmf_tuple := $(subst :, ,$(f))) \
    $(eval _cmf_src := $(word 1,$(_cmf_tuple))) \
    $(_cmf_src)))

$(art_host_cross_tests_zip) : PRIVATE_HOST_SHARED_LIBS := $(my_host_shared_lib_for_art_host_cross_tests)

$(art_host_cross_tests_zip) : $(COMPATIBILITY.art-host-cross-tests.FILES) $(my_host_shared_lib_for_art_host_cross_tests) $(SOONG_ZIP)
	echo $(sort $(COMPATIBILITY.art-host-cross-tests.FILES)) | tr " " "\n" > $@.list
	grep $(HOST_CROSS_OUT_TESTCASES) $@.list > $@-host.list || true
	$(hide) touch $@-host-cross-libs.list
	$(hide) for shared_lib in $(PRIVATE_HOST_SHARED_LIBS); do \
	  echo $$shared_lib >> $@-host-cross-libs.list; \
	done
	grep $(TARGET_OUT_TESTCASES) $@.list > $@-target.list || true
	$(hide) $(SOONG_ZIP) -d -o $@ -P host -C $(HOST_CROSS_OUT) -l $@-host.list \
	  -P target -C $(PRODUCT_OUT) -l $@-target.list \
	  -P host/testcases -C $(HOST_CROSS_OUT) -l $@-host-cross-libs.list
	rm -f $@.list $@-host.list $@-target.list $@-host-cross-libs.list

art-host-cross-tests: $(art_host_cross_tests_zip)
$(call dist-for-goals, art-host-cross-tests, $(art_host_cross_tests_zip))

$(call declare-1p-container,$(art_host_cross_tests_zip),)
$(call declare-container-license-deps,$(art_host_cross_tests_zip),$(COMPATIBILITY.art-host-cross-tests.FILES) $(my_host_shared_lib_for_art_host_cross_tests),$(PRODUCT_OUT)/:/)

tests: art-host-cross-tests

# ======================================================================
# FIXME: DEBUG.
.PHONY: my_host_shared_lib_for_art_host_cross_tests_list

my_host_shared_lib_for_art_host_cross_tests_list := $(PRODUCT_OUT)/my_host_shared_lib_for_art_host_cross_tests_list.txt

$(my_host_shared_lib_for_art_host_cross_tests_list) : PRIVATE_HOST_SHARED_LIBS := $(my_host_shared_lib_for_art_host_cross_tests)

$(my_host_shared_lib_for_art_host_cross_tests_list) :
	echo $(PRIVATE_HOST_SHARED_LIBS) >$@

my_host_shared_lib_for_art_host_cross_tests_list: $(my_host_shared_lib_for_art_host_cross_tests_list)
# ======================================================================
