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

LOCAL_AUTOGEN_TEST_CONFIG_FILE := $(dir $(LOCAL_BUILT_MODULE))$(LOCAL_MODULE).config
ifneq (true,$(is_native))
android_manifest := $(LOCAL_PATH)/AndroidManifest.xml
ifeq (,$(wildcard $(android_manifest)))
  $(warning No AndroidManifest.xml found at $(android_manifest), can't generate test config file.)
else

is_instrumentation :=  $(shell grep -o -e '<instrumentation' $(android_manifest))
ifeq (,$(is_instrumentation))
  $(warning $(LOCAL_MODULE) is not an instrumentation test, should not have LOCAL_COMPATIBILITY_SUITE setting.)
else
# Get package, test type and runner from AndroidManifest.xml
$(LOCAL_AUTOGEN_TEST_CONFIG_FILE) : PRIVATE_PACKAGE := $(shell \
  grep -Pazo '<instrumentation(.*\n*)+' $(android_manifest) | \
  grep -oe 'android:targetPackage="[^"]*' | head -n 1 | awk '{split($$0,a,"=\""); print a[2]}')
my_label := $(shell grep -Pazo '<instrumentation(.*\n*)+' $(android_manifest) | \
  grep -oe 'android:label="[^"]*' | head -n 1 | awk '{split($$0,a,"=\""); print a[2]}')
my_runner := $(shell grep -Pazo '<instrumentation(.*\n*)+' $(android_manifest) | \
  grep -oe 'android:name="[^"]*' | head -n 1 | awk '{split($$0,a,"=\""); print a[2]}')

ifneq (,$(findstring AndroidJUnitRunner,$(my_runner)))
  $(LOCAL_AUTOGEN_TEST_CONFIG_FILE) : PRIVATE_TEST_TYPE := AndroidJUnitTest
else
  $(LOCAL_AUTOGEN_TEST_CONFIG_FILE) : PRIVATE_TEST_TYPE := InstrumentationTest
endif
ifeq (,$(my_label)
  my_label := $(LOCAL_PACKAGE_NAME)
endif

$(LOCAL_AUTOGEN_TEST_CONFIG_FILE) : PRIVATE_LABEL := $(my_label)
$(LOCAL_AUTOGEN_TEST_CONFIG_FILE) : PRIVATE_TEST_RUNNER := $(my_runner)

$(LOCAL_AUTOGEN_TEST_CONFIG_FILE) :
	@echo "Auto generating test config $(notdir $@)"
	$(hide) sed 's&{TEST_TYPE}&$(PRIVATE_TEST_TYPE)&' $(INSTRUMENTATION_TEST_CONFIG_TEMPLATE) > $@.temp1
	$(hide) sed 's&{RUNNER}&$(PRIVATE_TEST_RUNNER)&' $@.temp1 > $@.temp2
	$(hide) sed 's&{LABEL}&$(PRIVATE_LABEL)&' $@.temp2 > $@.temp1
	$(hide) sed 's&{PACKAGE}&$(PRIVATE_PACKAGE)&' $@.temp1 > $@.temp2
	$(hide) sed 's&{MODULE}&$(PRIVATE_MODULE)&' $@.temp2 > $@
	rm -f $@.temp1 $@.temp2
my_auto_generate_config := true

endif # ifeq (,$(is_instrumentation))
endif # ifeq (,$(wildcard $(android_manifest)))
else
# Auto generating test config file for native test
$(LOCAL_AUTOGEN_TEST_CONFIG_FILE) :
	@echo "Auto generating test config $(notdir $@)"
	$(hide) sed 's&{MODULE}&$(PRIVATE_MODULE)&' $(NATIVE_TEST_CONFIG_TEMPLATE) > $@
my_auto_generate_config := true
endif # ifneq (true,$(is_native))

ifeq (true,$(my_auto_generate_config))
  LOCAL_INTERMEDIATE_TARGETS += $(LOCAL_AUTOGEN_TEST_CONFIG_FILE)
  $(LOCAL_BUILT_MODULE): $(LOCAL_AUTOGEN_TEST_CONFIG_FILE)
else
  LOCAL_AUTOGEN_TEST_CONFIG_FILE :=
endif

is_instrumentation :=
my_auto_generate_config :=
my_label :=
my_runner :=
