#
# Copyright (C) 2018 The Android Open Source Project
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

# The following PRODUCT_* variables seem to be necessary if
# require-artifacts-in-path is invoked in this make file. Why??
PRODUCT_NAME := cf_vendor
PRODUCT_BRAND := generic

#PRODUCT_COPY_FILES += device/google/cuttlefish_kernel/4.4-x86_64/kernel:kernel

PRODUCT_PACKAGES := libvsoc-ril rild
PRODUCT_PACKAGES += Telecom	# This should be caught by the build system!

_vendor_whitelist :=
_vendor_paths := $(TARGET_COPY_OUT_VENDOT)/
$(call require-artifacts-in-path, $(_vendor_paths), $(_vendor_whitelist))
