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

# This file sets up Java code coverage via Jacoco
# This file is only intended to be included internally by the build system
# (at the time of authorship, it is included by java.mk and
# java_host_library.mk)

my_jacoco_include_filter :=
my_jacoco_exclude_filter :=

ifeq ($(LOCAL_EMMA_INSTRUMENT),true)
  ifeq ($(ANDROID_COMPILE_WITH_JACK),false)
    # determine Jacoco include/exclude filters
    DEFAULT_JACOCO_EXCLUDE_FILTER := org/junit/*,org/jacoco/*,org/mockito/*
    # copy filters from Jack but also skip some known java packages
    my_jacoco_include_filter := $(strip $(LOCAL_JACK_COVERAGE_INCLUDE_FILTER))
    my_jacoco_exclude_filter := $(strip $(DEFAULT_JACOCO_EXCLUDE_FILTER),$(LOCAL_JACK_COVERAGE_EXCLUDE_FILTER))

    # replace "." with "/" and "," with " ", and quote each arg
    ifneq ($(strip $(my_jacoco_include_filter)),)
      my_jacoco_include_args := $(strip $(my_jacoco_include_filter))

      my_jacoco_include_args := $(subst .,/,$(my_jacoco_include_args))
      my_jacoco_include_args := "$(subst $(comma),"$(space)",$(my_jacoco_include_args))"
    else
      my_jacoco_include_args :=
    endif

    # replace "." with "/" and "," with " ", and quote each arg
    ifneq ($(strip $(my_jacoco_exclude_filter)),)
      my_jacoco_exclude_args := $(my_jacoco_exclude_filter)

      my_jacoco_exclude_args := $(subst .,/,$(my_jacoco_exclude_args))
      my_jacoco_exclude_args := $(subst $(comma)$(comma),$(comma),$(my_jacoco_exclude_args))
      my_jacoco_exclude_args := "$(subst $(comma)," ", $(my_jacoco_exclude_args))"
    else
      my_jacoco_exclude_args :=
    endif
  endif #ANDROID_COMPILE_WITH_JACK==false
endif #LOCAL_EMMA_INSTRUMENT == true

# determine whether to run the instrumenter based on whether there is any work
# for it to do
ifneq ($(my_jacoco_include_filter),)

  my_jacoco_files := $(intermediates.COMMON)/jacoco

  # make a task that unzips the classes that we want to instrument from the
  # input jar
  my_jacoco_unzipped_path := $(my_jacoco_files)/work/classes-to-instrument/classes
  my_jacoco_unzipped_timestamp_path := $(my_jacoco_files)/work/classes-to-instrument/updated.stamp
$(my_jacoco_unzipped_timestamp_path): PRIVATE_JACOCO_UNZIPPED_PATH := $(my_jacoco_unzipped_path)
$(my_jacoco_unzipped_timestamp_path): PRIVATE_JACOCO_UNZIPPED_TIMESTAMP_PATH := $(my_jacoco_unzipped_timestamp_path)
$(my_jacoco_unzipped_timestamp_path): PRIVATE_JACOCO_INCLUDE_ARGS := $(my_jacoco_include_args)
$(my_jacoco_unzipped_timestamp_path): PRIVATE_JACOCO_EXCLUDE_ARGS := $(my_jacoco_exclude_args)
$(my_jacoco_unzipped_timestamp_path): PRIVATE_FULL_CLASSES_PRE_JACOCO_JAR := $(LOCAL_FULL_CLASSES_PRE_JACOCO_JAR)
$(my_jacoco_unzipped_timestamp_path): $(LOCAL_FULL_CLASSES_PRE_JACOCO_JAR)
	rm -rf $(PRIVATE_JACOCO_UNZIPPED_PATH) $@
	mkdir -p $(PRIVATE_JACOCO_UNZIPPED_PATH)
	unzip -q $(PRIVATE_FULL_CLASSES_PRE_JACOCO_JAR) \
	  -d $(PRIVATE_JACOCO_UNZIPPED_PATH) \
	  $(PRIVATE_JACOCO_INCLUDE_ARGS)
	rm -rf $(PRIVATE_JACOCO_EXCLUDE_ARGS)
	touch $(PRIVATE_JACOCO_UNZIPPED_TIMESTAMP_PATH)
# Unfortunately in the previous task above,
# 'rm -rf $(PRIVATE_JACOCO_EXCLUDE_ARGS)' needs to be a separate
# shell command after 'unzip'.
# We can't just use the '-x' (exclude) option of 'unzip' because if both
# inclusions and exclusions are specified and an exclusion matches no
# inclusions, then 'unzip' exits with an error (error 11).
# We could ignore the error, but that would make the process less reliable


  # make a task that zips only the classes that will be instrumented
  # (for passing in to the report generator later)
  my_jacoco_classes_to_report_on_path := $(my_jacoco_files)/report-resources/jacoco-report-classes.jar
$(my_jacoco_classes_to_report_on_path): PRIVATE_JACOCO_UNZIPPED_PATH := $(my_jacoco_unzipped_path)
$(my_jacoco_classes_to_report_on_path): $(my_jacoco_unzipped_timestamp_path)
	rm -f $@
	zip -q $@ \
	  -r $(PRIVATE_JACOCO_UNZIPPED_PATH)



  # make a task that invokes instrumentation
  my_jacoco_instrumented_path := $(my_jacoco_files)/work/instrumented/classes
  my_jacoco_instrumented_timestamp_path := $(my_jacoco_files)/work/instrumented/updated.stamp
$(my_jacoco_instrumented_timestamp_path): PRIVATE_JACOCO_INSTRUMENTED_PATH := $(my_jacoco_instrumented_path)
$(my_jacoco_instrumented_timestamp_path): PRIVATE_JACOCO_INSTRUMENTED_TIMESTAMP_PATH := $(my_jacoco_instrumented_timestamp_path)
$(my_jacoco_instrumented_timestamp_path): PRIVATE_JACOCO_UNZIPPED_PATH := $(my_jacoco_unzipped_path)
$(my_jacoco_instrumented_timestamp_path): $(my_jacoco_unzipped_timestamp_path) $(JACOCO_CLI_JAR)
	rm -rf $(PRIVATE_JACOCO_INSTRUMENTED_PATH)
	mkdir -p $(PRIVATE_JACOCO_INSTRUMENTED_PATH)
	java -jar $(JACOCO_CLI_JAR) \
	  instrument \
	  -quiet \
	  -dest '$(PRIVATE_JACOCO_INSTRUMENTED_PATH)' \
	  $(PRIVATE_JACOCO_UNZIPPED_PATH)
	touch $(PRIVATE_JACOCO_INSTRUMENTED_TIMESTAMP_PATH)


  # make a task that zips both the instrumented classes and the uninstrumented
  # classes (this jar is the instrumented application to execute)
  my_jacoco_temp_jar_path := $(my_jacoco_files)/work/usable.jar
  LOCAL_FULL_CLASSES_JACOCO_JAR := $(intermediates.COMMON)/classes-jacoco.jar
$(LOCAL_FULL_CLASSES_JACOCO_JAR): PRIVATE_TEMP_JAR_PATH := $(my_jacoco_temp_jar_path)
$(LOCAL_FULL_CLASSES_JACOCO_JAR): PRIVATE_JACOCO_INSTRUMENTED_PATH := $(my_jacoco_instrumented_path)
$(LOCAL_FULL_CLASSES_JACOCO_JAR): PRIVATE_JACOCO_INSTRUMENTED_TIMESTAMP_PATH := $(my_jacoco_instrumented_timestamp_path)
$(LOCAL_FULL_CLASSES_JACOCO_JAR): PRIVATE_FULL_CLASSES_PRE_JACOCO_JAR := $(LOCAL_FULL_CLASSES_PRE_JACOCO_JAR)
$(LOCAL_FULL_CLASSES_JACOCO_JAR): $(my_jacoco_instrumented_timestamp_path) $(LOCAL_FULL_CLASSES_PRE_JACOCO_JAR)
	rm -f $@ $(PRIVATE_TEMP_JAR_PATH)
	cp $(PRIVATE_FULL_CLASSES_PRE_JACOCO_JAR) $(PRIVATE_TEMP_JAR_PATH)
	JarPath=`readlink -f $(PRIVATE_TEMP_JAR_PATH)`
	cd $(PRIVATE_JACOCO_INSTRUMENTED_PATH)
	zip -q -r $$JarPath .
	cd -
	mv $(PRIVATE_TEMP_JAR_PATH) $@


  # this is used to trigger $(my_jacoco_classes_to_report_on_path) to build
  # when $(LOCAL_FULL_CLASSES_JACOCO_JAR) builds, but it isn't truly a
  # dependency.
$(LOCAL_FULL_CLASSES_JACOCO_JAR): $(my_jacoco_classes_to_report_on_path)

else #my_jacoco_include_filter == ''
  LOCAL_FULL_CLASSES_JACOCO_JAR := $(LOCAL_FULL_CLASSES_PRE_JACOCO_JAR)
endif #my_jacoco_include_filter != ''

LOCAL_INTERMEDIATE_TARGETS += $(LOCAL_FULL_CLASSES_JACOCO_JAR)
