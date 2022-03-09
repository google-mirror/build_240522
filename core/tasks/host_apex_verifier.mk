#
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
#

HOST_APEX_VERIFIER := $(HOST_OUT_EXECUTABLES)/host_apex_verifier
DEAPEXER := $(HOST_OUT_EXECUTABLES)/deapexer
DEBUGFS := $(HOST_OUT_EXECUTABLES)/debugfs_static

host_apex_verifier_output := $(PRODUCT_OUT)/host_apex_verifier_output.txt

$(host_apex_verifier_output): \
	$(HOST_APEX_VERIFIER) \
	$(DEAPEXER) \
	$(DEBUGFS) \
	$(INSTALLED_SYSTEMIMAGE_TARGET) \
	$(INSTALLED_SYSTEM_EXTIMAGE_TARGET) \
	$(INSTALLED_PRODUCTIMAGE_TARGET) \
	$(INSTALLED_VENDORIMAGE_TARGET) \
	$(INSTALLED_ODMIMAGE_TARGET) \

# Run host_apex_verifier on the partition staging directories.
$(host_apex_verifier_output): $(HOST_APEX_VERIFIER)
	$(HOST_APEX_VERIFIER) \
		--deapexer $(DEAPEXER) \
		--debugfs $(DEBUGFS) \
		--out_system $(PRODUCT_OUT)/$(TARGET_COPY_OUT_SYSTEM) \
		--out_system_ext $(PRODUCT_OUT)/$(TARGET_COPY_OUT_SYSTEM_EXT) \
		--out_product $(PRODUCT_OUT)/$(TARGET_COPY_OUT_PRODUCT) \
		--out_vendor $(PRODUCT_OUT)/$(TARGET_COPY_OUT_VENDOR) \
		--out_odm $(PRODUCT_OUT)/$(TARGET_COPY_OUT_ODM) \
		> $@

$(call dist-for-goals,droidcore-unbundled,$(host_apex_verifier_output))

.PHONY: host-apex-verifier
host-apex-verifier: $(host_apex_verifier_output)
