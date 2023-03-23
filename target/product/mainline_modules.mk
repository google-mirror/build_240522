# Copyright (C) 2023 The Android Open Source Project
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

# The mainline_modules architecture agnostic product is used exclusively by
# Bazel. Bazel operates in unbundled mode by default, and the target
# architecture for the build is decided by bazel itself.

$(call inherit-product, $(SRC_TARGET_DIR)/product/module_common.mk)

PRODUCT_NAME := mainline_modules
PRODUCT_BRAND := Android
PRODUCT_DEVICE := mainline_modules
