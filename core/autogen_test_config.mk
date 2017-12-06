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
my_test_config_file :=
my_android_test_xml :=

# Output test config files to testcases directory.
ifeq (,$(filter general-tests, $(LOCAL_COMPATIBILITY_SUITE)))
  LOCAL_COMPATIBILITY_SUITE += general-tests
endif

LOCAL_AUTOGEN_TEST_CONFIG_FILE := $(dir $(LOCAL_BUILT_MODULE))$(LOCAL_PACKAGE_NAME).config
ifeq (true,$(LOCAL_AUTOGEN_INSTRUMENTATION_TEST_CONFIG))
$(LOCAL_AUTOGEN_TEST_CONFIG_FILE) : PRIVATE_PACKAGE := com.android.notification.tests
ifeq (true,$(LOCAL_AUTOGEN_TEST_CONFIG_IS_ANDROIDJUNITTEST))
  $(LOCAL_AUTOGEN_TEST_CONFIG_FILE) : PRIVATE_test_type := AndroidJUnitTest
  $(LOCAL_AUTOGEN_TEST_CONFIG_FILE) : PRIVATE_test_runner := android.support.test.runner.AndroidJUnitRunner
else
  $(LOCAL_AUTOGEN_TEST_CONFIG_FILE) : PRIVATE_test_type := InstrumentationTest
  $(LOCAL_AUTOGEN_TEST_CONFIG_FILE) : PRIVATE_test_runner := android.test.InstrumentationTestRunner
endif
$(LOCAL_AUTOGEN_TEST_CONFIG_FILE) :
	@echo "Auto generating test config $(notdir $@)"
	$(hide) sed 's&{TEST_TYPE}&$(PRIVATE_test_type)&' $(INSTRUMENTATION_TEST_CONFIG_TEMPLATE) > $@.temp1
	$(hide) sed 's&{RUNNER}&$(PRIVATE_test_runner)&' $@.temp1 > $@.temp2
	$(hide) sed 's&{PACKAGE}&$(LOCAL_AUTOGEN_TEST_CONFIG_PACKAGE)&' $@.temp2 > $@.temp1
	$(hide) sed 's&{MODULE}&$(PRIVATE_MODULE)&' $@.temp1 > $@
	rm -f $@.temp1 $@.temp2
endif # ifeq (true,$(LOCAL_AUTOGEN_INSTRUMENTATION_TEST_CONFIG))

ifeq (true,$(LOCAL_AUTOGEN_NATIVE_TEST_CONFIG))
  $(error Auto generating test config file for native test is not supported yet.)
endif # ifeq (true,$(LOCAL_AUTOGEN_NATIVE_TEST_CONFIG))

LOCAL_INTERMEDIATE_TARGETS += $(LOCAL_AUTOGEN_TEST_CONFIG_FILE)
$(LOCAL_BUILT_MODULE): $(LOCAL_AUTOGEN_TEST_CONFIG_FILE)
