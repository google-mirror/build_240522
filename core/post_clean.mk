# Copyright (C) 2012 The Android Open Source Project
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

# Clean steps that need global knowledge of individual modules.
# This file must be included after all Android.mks have been loaded.

#######################################################
# Check if we need to delete obsolete generated java files.
# When an proto/etc file gets deleted (or renamed), the generated java
# file is obsolete.
#
# We solve this by writing out the list of modules withtheir source files that
# will create generated source files. Then build/make/tools/clean_gen_java.py
# can compare that map against the previous, and if the file list has changed,
# remove the intermediate directories causing them to rebuild.
#
# We need to store off a list of all possible modules separately in case we're
# building with 'mm' and only know about a subset of the tree. The tool can
# ignore any map entries that aren't in this possible module list.
gen_java_modules :=
gen_java_files :=
$(foreach p,$(sort $(ALL_MODULES)), \
  $(foreach class,$(sort $(filter APPS JAVA_LIBRARIES,$(ALL_MODULES.$(p).CLASS))), \
    $(eval gen_java_modules += $$(class)/$$(p)_intermediates$$(newline)) \
    $(eval gs := $$(strip $$(ALL_MODULES.$$(p).PROTO_FILES) \
                          $$(ALL_MODULES.$$(p).RS_FILES))) \
    $(if $(gs), \
      $(eval gen_java_files += $$(class)/$$(p)_intermediates $$(gs)$$(newline)))))

current_gen_java_config := $(TARGET_OUT_COMMON_INTERMEDIATES)/current_gen_java.txt
current_all_modules_config := $(TARGET_OUT_COMMON_INTERMEDIATES)/current_modules.txt

$(shell mkdir -p $(dir $(current_gen_java_config)))

$(file >$(current_gen_java_config),$(gen_java_files))
$(file >$(current_all_modules_config),$(gen_java_modules))

current_gen_java_config :=
current_all_modules_config :=
gen_java_modules :=
gen_java_files :=
