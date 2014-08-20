#
# Copyright (C) 2014 The Android Open Source Project
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

# These definitions are related to the ARM64-specific virtual board named
# 'ranchu' which is required to run under emulator with the binaries built
# from https://qemu-android.googlesource.com/qemu-android
#
# In a nutshell, these are based on a recent version of QEMU, while the
# Android emulator sources are based on much older release.
#

PRODUCT_COPY_FILES += \
    device/linaro/ranchu/fstab.ranchu:root/fstab.ranchu \
    device/linaro/ranchu/init.ranchu.rc:root/init.ranchu.rc \
    device/linaro/ranchu/ueventd.ranchu.rc:root/ueventd.ranchu.rc
