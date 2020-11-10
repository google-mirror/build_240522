#
# Copyright (C) 2020 The Android Open Source Project
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

host_init_verifier_output := $(PRODUCT_OUT)/host_init_verifier_output.txt

$(host_init_verifier_output): \
    $(INSTALLED_SYSTEMIMAGE_TARGET) \
    $(INSTALLED_SYSTEM_EXTIMAGE_TARGET) \
    $(INSTALLED_VENDORIMAGE_TARGET) \
    $(INSTALLED_ODMIMAGE_TARGET) \
    $(INSTALLED_PRODUCTIMAGE_TARGET)

# Run host_init_verifier on the partition staging directories.
$(host_init_verifier_output): $(HOST_INIT_VERIFIER)
	$(HOST_INIT_VERIFIER) \
		--out_system $(PRODUCT_OUT)/$(TARGET_COPY_OUT_SYSTEM) \
		--out_system_ext $(PRODUCT_OUT)/$(TARGET_COPY_OUT_SYSTEM_EXT) \
		--out_vendor $(PRODUCT_OUT)/$(TARGET_COPY_OUT_VENDOR) \
		--out_odm $(PRODUCT_OUT)/$(TARGET_COPY_OUT_ODM) \
		--out_product $(PRODUCT_OUT)/$(TARGET_COPY_OUT_PRODUCT) \
		> $@

$(call dist-for-goals,droidcore,$(host_init_verifier_output))
