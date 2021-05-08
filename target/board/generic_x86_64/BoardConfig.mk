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

# x86_64 emulator specific definitions
TARGET_CPU_ABI := x86_64
TARGET_ARCH := x86_64
TARGET_ARCH_VARIANT := x86_64

TARGET_2ND_CPU_ABI := x86
TARGET_2ND_ARCH := x86
TARGET_2ND_ARCH_VARIANT := x86_64

include build/make/target/board/BoardConfigGsiCommon.mk

BOARD_KERNEL-5.4_BOOTIMAGE_PARTITION_SIZE := 67108864
BOARD_KERNEL-5.4-ALLSYMS_BOOTIMAGE_PARTITION_SIZE := 67108864
BOARD_KERNEL-5.10_BOOTIMAGE_PARTITION_SIZE := 67108864
BOARD_KERNEL-5.10-ALLSYMS_BOOTIMAGE_PARTITION_SIZE := 67108864

BOARD_USERDATAIMAGE_PARTITION_SIZE := 576716800

BOARD_RAMDISK_USE_LZ4 := true
BOARD_BOOT_HEADER_VERSION := 4
BOARD_MKBOOTIMG_ARGS += --header_version $(BOARD_BOOT_HEADER_VERSION)

# Enable GKI 2.0 signing.
BOARD_GKI_SIGNING_KEY_PATH := build/make/target/product/gsi/testkey_rsa2048.pem
BOARD_GKI_SIGNING_ALGORITHM := SHA256_RSA2048

BOARD_KERNEL_BINARIES := \
    kernel-5.4 \
    kernel-5.10 \

ifneq (,$(filter userdebug eng,$(TARGET_BUILD_VARIANT)))
BOARD_KERNEL_BINARIES += \
    kernel-5.4-allsyms \
    kernel-5.10-allsyms \

endif

# Boot image
BOARD_USES_RECOVERY_AS_BOOT :=
TARGET_NO_KERNEL := false
BOARD_USES_GENERIC_KERNEL_IMAGE := true
# TODO(b/187432172): Add 5.10-android12-unstable
BOARD_KERNEL_MODULE_INTERFACE_VERSIONS := \
    5.4-android12-0 \

# Copy boot image in $OUT to target files. This is defined for targets where
# the installed GKI APEXes are built from source.
BOARD_COPY_BOOT_IMAGE_TO_TARGET_FILES := true

# No vendor_boot
BOARD_MOVE_RECOVERY_RESOURCES_TO_VENDOR_BOOT :=

# No recovery
BOARD_EXCLUDE_KERNEL_FROM_RECOVERY_IMAGE :=
