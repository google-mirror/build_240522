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

# This is a generic phone product that isn't specialized for a specific device.
# It includes the base Android platform.

$(call inherit-product, $(SRC_TARGET_DIR)/product/aosp_system.mk)
$(call inherit-product, $(SRC_TARGET_DIR)/product/telephony.mk)

# Overrides
PRODUCT_BRAND := generic
PRODUCT_DEVICE := generic
PRODUCT_NAME := generic

PRODUCT_ENFORCE_ISOLATION_CLAIMS := true
PRODUCT_ISOLATION_CLAIM_WHITELIST := \
	system/app/CarrierDefaultApp/CarrierDefaultApp.apk \
	system/priv-app/CallLogBackup/CallLogBackup.apk \
	system/priv-app/CarrierConfig/CarrierConfig.apk \
	system/priv-app/CellBroadcastReceiver/CellBroadcastReceiver.apk \
	system/priv-app/Dialer/Dialer.apk \
	system/priv-app/EmergencyInfo/EmergencyInfo.apk \
	vendor/bin/hw/rild \
	vendor/etc/init/rild.rc \
