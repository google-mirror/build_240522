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
ifeq (true,$(is_native))
# Auto generating test config file for native test
$(LOCAL_AUTOGEN_TEST_CONFIG_FILE) :
	@echo "Auto generating test config $(notdir $@)"
	$(hide) sed 's&{MODULE}&$(PRIVATE_MODULE)&g' $(NATIVE_TEST_CONFIG_TEMPLATE) > $@
my_auto_generate_config := true
else
# Auto generating test config file for instrumentation test
ifdef LOCAL_MANIFEST_FILE
  android_manifest := $(LOCAL_PATH)/$(LOCAL_MANIFEST_FILE)
else
  android_manifest := $(LOCAL_PATH)/AndroidManifest.xml
endif

ifneq (,$(wildcard $(android_manifest)))
$(LOCAL_AUTOGEN_TEST_CONFIG_FILE) : $(android_manifest)
	@echo "Auto generating test config $(notdir $@)"
	# The rule will parse AndroidManifest.xml file to compile the test config file.
	# If the AndroidManifest.xml file doesn't have instrumentation element, a dummy
	# test config file will be created.
	# The grep command to parse instrumentation element:
	#   -P: activate perl-regexp for grep
	#   -z: suppress newline at the end of line, subtituting it for null character,
	#       so grep can search across multiple lines.
	#		-o: print only matching.
	# The tr command merges the multiple lines of instrumentation element into a
	# single line.
	$(shell) instrumentation=$(shell grep -Pzo '<instrumentation([^<]*\n*)+' $^ | tr '\n' ' '); \
		if [[ -z $$instrumentation ]]; then \
			cat $(EMPTY_TEST_CONFIG) > $@; \
		else \
			package=$(shell echo $$instrumentation | grep -oe 'android:targetPackage="[^"]*' | awk '{split($$0,a,"=\""); print a[2]}'); \
			label=$(shell echo $$instrumentation | grep -oe 'android:label="[^"]*' | awk '{split($$0,a,"=\""); print a[2]}'); \
			runner=$(shell echo $$instrumentation | grep -oe 'android:name="[^"]*' | awk '{split($$0,a,"=\""); print a[2]}'); \
			if [[ $$runner == *"AndroidJUnitRunner" ]]; then \
				test_type=AndroidJUnitTest; \
			else \
				test_type=InstrumentationTest; \
			fi; \
			sed "s&{TEST_TYPE}&$$test_type&" $(INSTRUMENTATION_TEST_CONFIG_TEMPLATE) > $@.temp1; \
			sed "s&{RUNNER}&$$runner&" $@.temp1 > $@.temp2; \
			sed "s&{LABEL}&$$label&" $@.temp2 > $@.temp1; \
			sed "s&{PACKAGE}&$$package&" $@.temp1 > $@.temp2; \
			sed "s&{MODULE}&$(PRIVATE_MODULE)&" $@.temp2 > $@; \
			rm -f $@.temp1 $@.temp2; \
		fi
my_auto_generate_config := true
endif # ifeq (,$(wildcard $(android_manifest)))
endif # ifneq (true,$(is_native))

ifeq (true,$(my_auto_generate_config))
  LOCAL_INTERMEDIATE_TARGETS += $(LOCAL_AUTOGEN_TEST_CONFIG_FILE)
  $(LOCAL_BUILT_MODULE): $(LOCAL_AUTOGEN_TEST_CONFIG_FILE)
  ALL_MODULES.$(my_register_name).auto_test_config := true
else
  LOCAL_AUTOGEN_TEST_CONFIG_FILE :=
endif

my_auto_generate_config :=
