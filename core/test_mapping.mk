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

# Create an artifact to include TEST_MAPPING files in source tree.

.PHONY: test_mapping

intermediates := $(call intermediates-dir-for,PACKAGING,test_mapping)
test_mappings_tarball := $(intermediates)/test_mappings.tar.bz2
test_mapping_list := $(file <$(OUT_DIR)/.module_paths/TEST_MAPPING.list)

$(test_mappings_tarball) : $(test_mapping_list)
	@echo "Building artifact to include TEST_MAPPING files."
	$(eval intermediates := $(dir $@))
	$(eval test_mapping_dir := $(intermediates)/test_mappings)
	$(hide) rm -rf $(test_mapping_dir)
	$(hide) mkdir -p $(test_mapping_dir)
	$(foreach f, $^,\
	  $(eval target_dir := $(test_mapping_dir)/$(dir $(f))) \
		$(shell mkdir -p $(target_dir)) \
		$(shell cp $(f) $(target_dir)))
	$(hide) cd $(intermediates) && tar -jcvf $(notdir $@) -C test_mappings .

test_mapping : $(test_mappings_tarball)

$(call dist-for-goals, dist_files test_mapping,$(test_mappings_tarball))

