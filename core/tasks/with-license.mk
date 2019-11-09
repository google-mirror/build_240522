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


# DO NOT PROCEED without a license file.
ifndef LICENSE_NOTICE_FILE
$(error with-license requires LICENSE_NOTICE_FILE to be set)
endif

.PHONY: with-license

name := $(TARGET_PRODUCT)
ifeq ($(TARGET_BUILD_TYPE),debug)
	name := $(name)_debug
endif

name := $(name)-img-$(FILE_NAME_TAG)-with-license

with_license_intermediates := \
	$(call intermediates-dir-for,PACKAGING,with_license)

# Create a with-license artifact target
license-image-input-zip := $(with_license_intermediates)/$(name).zip
$(license-image-input-zip) : $(BUILT_TARGET_FILES_PACKAGE) $(ZIP2ZIP)
	$(ZIP2ZIP) -i $(BUILT_TARGET_FILES_PACKAGE) -o $@ \
		RADIO/bootloader.img:bootloader.img RADIO/radio.img:radio.img \
		IMAGES/system.img:system.img IMAGES/vendor.img:vendor.img \
		IMAGES/boot.img:boot.img OTA/android-info.txt:android-info.txt
with-license-zip := $(PRODUCT_OUT)/$(name).sh
$(with-license-zip) : $(license-image-input-zip) $(LICENSE_NOTICE_FILE)
	# Args: <output> <input archive> <comment> <license file>
	build/make/tools/generate-self-extracting-archive.py $@ $(license-image-input-zip) \
		$(name) $(LICENSE_NOTICE_FILE)
	chmod ug+x $@
with-license : $(with-license-zip)
$(call dist-for-goals, with-license, $(with-license-zip))
