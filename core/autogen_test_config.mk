#
# Copyright (C) 2017 The Android Open Source Project
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

$(warning *** Auto generating test config for: $(LOCAL_PACKAGE_NAME))

# If there is already a test config file named as $(LOCAL_MODULE).xml or AndroidTest.xml,
# auto generating test config should fail.
my_test_config_file := $(wildcard $(LOCAL_PATH)/$(LOCAL_MODULE).xml)
ifdef my_test_config_file
  $(error Test config file already exists: $(my_test_config_file))
endif
my_android_test_xml := $(wildcard $(LOCAL_PATH)/AndroidTest.xml)
ifdef my_android_test_xml
  $(error Test config file already exists: $(my_android_test_xml))
endif

# Output test config files to testcases directory.
ifeq (,$(filter general-tests, $(LOCAL_COMPATIBILITY_SUITE)))
  LOCAL_COMPATIBILITY_SUITE += general-tests
endif

test_config_file := $(dir $(LOCAL_BUILT_MODULE))$(LOCAL_PACKAGE_NAME).config
ifeq ($(LOCAL_IS_INSTRUMENTATION_TEST),true)
$(test_config_file) :
	@echo "Auto generating test config $(notdir $@)"
	$(hide) sed 's&{LOCAL_PACKAGE_NAME}&$(LOCAL_PACKAGE_NAME)&' $(INSTRUMENTATION_TEST_CONFIG_TEMPLATE) > $@
LOCAL_INTERMEDIATE_TARGETS += $(test_config_file)
endif

my_test_config_file :=
my_android_test_xml :=
test_config_file :=
