#
# Copyright (C) 2017 The Android Open-Source Project
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

$(call inherit-product, $(SRC_TARGET_DIR)/product/generic_no_telephony.mk)

PRODUCT_NAME := aosp_system

#fake_packages/art-runtime-timestamp
_my_whitelist := \
  fake_packages/art-tools-timestamp \
  fake_packages/bpf_kern.o-timestamp \
  fake_packages/javax.obex-timestamp \
  fake_packages/org.apache.http.legacy-timestamp \
  fake_packages/selinux_policy-timestamp \
  fake_packages/shell_and_utilities-timestamp \
  fake_packages/vndk_package-timestamp \
  fake_packages/vndk_snapshot_package-timestamp \
  recovery/root/etc/mke2fs.conf \
  recovery/root/sbin/adbd \
  recovery/root/sbin/e2fsdroid_static \
  recovery/root/sbin/mke2fs_static \
  recovery/root/sbin/recovery \
  root/init \
  root/init.environ.rc \
  root/init.rc \
  root/plat_file_contexts \
  root/plat_hwservice_contexts \
  root/plat_property_contexts \
  root/plat_seapp_contexts \
  root/plat_service_contexts \
  root/sbin/charger \
  root/sepolicy \
  root/vendor_file_contexts \
  root/vendor_hwservice_contexts \
  root/vendor_property_contexts \
  root/vendor_seapp_contexts \
  root/vendor_service_contexts \
  root/vndservice_contexts \

$(call make-isolation-claim, $(TARGET_COPY_OUT_SYSTEM) $(TARGET_COPY_OUT_VENDOR), $(_my_whitelist))
