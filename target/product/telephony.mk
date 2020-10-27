#
# Copyright (C) 2007 The Android Open Source Project
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

<<<<<<< HEAD   (5c8d84 Merge "Merge empty history for sparse-6676661-L8360000065797)
# This is the list of product-level settings that are specific
# to products that have telephony hardware.

PRODUCT_PACKAGES := \
    CarrierConfig \
    CarrierDefaultApp \
    Dialer \
    CallLogBackup \
    CellBroadcastReceiver \
    EmergencyInfo \
    rild

PRODUCT_COPY_FILES := \
=======
# All modules for telephony
$(call inherit-product, $(SRC_TARGET_DIR)/product/telephony_system.mk)
$(call inherit-product, $(SRC_TARGET_DIR)/product/telephony_system_ext.mk)
$(call inherit-product, $(SRC_TARGET_DIR)/product/telephony_vendor.mk)
$(call inherit-product, $(SRC_TARGET_DIR)/product/telephony_product.mk)
>>>>>>> BRANCH (a10c18 Merge "Version bump to RT11.201014.001.A1 [core/build_id.mk])
