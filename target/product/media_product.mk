#
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
#

# This makefile contains the product partition contents for
# media-capable devices (non-wearables). Only add something here
# if it definitely doesn't belong on wearables. Otherwise, choose
# base_vendor.mk.
$(call inherit-product, $(SRC_TARGET_DIR)/product/base_product.mk)

# /product packages
PRODUCT_PACKAGES += \
    webview \
