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

# arm64 emulator specific definitions
TARGET_ARCH := arm64
TARGET_ARCH_VARIANT := armv8-a
TARGET_CPU_VARIANT := generic
TARGET_CPU_ABI := arm64-v8a

TARGET_2ND_ARCH := arm
TARGET_2ND_CPU_ABI := armeabi-v7a
TARGET_2ND_CPU_ABI2 := armeabi

ifneq ($(TARGET_BUILD_APPS)$(filter cts sdk,$(MAKECMDGOALS)),)
# DO NOT USE
# DO NOT USE
#
# This architecture / CPU variant must NOT be used for any 64 bit
# platform builds. It is the lowest common denominator required
# to build an unbundled application or cts for all supported 32 and 64 bit
# platforms.
#
# If you're building a 64 bit platform (and not an application) the
# ARM-v8 specification allows you to assume all the features available in an
# armv7-a-neon CPU. You should set the following as 2nd arch/cpu variant:
#
# TARGET_2ND_ARCH_VARIANT := armv8-a
# TARGET_2ND_CPU_VARIANT := generic
#
# DO NOT USE
# DO NOT USE
TARGET_2ND_ARCH_VARIANT := armv7-a-neon
# DO NOT USE
# DO NOT USE
TARGET_2ND_CPU_VARIANT := generic
# DO NOT USE
# DO NOT USE
else
TARGET_2ND_ARCH_VARIANT := armv8-a
TARGET_2ND_CPU_VARIANT := generic
endif

include build/make/target/board/BoardConfigGsiCommon.mk

BOARD_KERNEL-4.19-GZ_BOOTIMAGE_PARTITION_SIZE := 47185920
BOARD_KERNEL-5.4_BOOTIMAGE_PARTITION_SIZE := 67108864
BOARD_KERNEL-5.4-ALLSYMS_BOOTIMAGE_PARTITION_SIZE := 67108864
BOARD_KERNEL-5.4-GZ_BOOTIMAGE_PARTITION_SIZE := 47185920
BOARD_KERNEL-5.4-GZ-ALLSYMS_BOOTIMAGE_PARTITION_SIZE := 47185920
BOARD_KERNEL-5.4-LZ4_BOOTIMAGE_PARTITION_SIZE := 53477376
BOARD_KERNEL-5.4-LZ4-ALLSYMS_BOOTIMAGE_PARTITION_SIZE := 53477376
BOARD_KERNEL-5.10_BOOTIMAGE_PARTITION_SIZE := 67108864
BOARD_KERNEL-5.10-ALLSYMS_BOOTIMAGE_PARTITION_SIZE := 67108864
BOARD_KERNEL-5.10-GZ_BOOTIMAGE_PARTITION_SIZE := 47185920
BOARD_KERNEL-5.10-GZ-ALLSYMS_BOOTIMAGE_PARTITION_SIZE := 47185920
BOARD_KERNEL-5.10-LZ4_BOOTIMAGE_PARTITION_SIZE := 53477376
BOARD_KERNEL-5.10-LZ4-ALLSYMS_BOOTIMAGE_PARTITION_SIZE := 53477376
BOARD_KERNEL-MAINLINE_BOOTIMAGE_PARTITION_SIZE := 67108864
BOARD_KERNEL-MAINLINE-GZ_BOOTIMAGE_PARTITION_SIZE := 47185920
BOARD_KERNEL-MAINLINE-LZ4_BOOTIMAGE_PARTITION_SIZE := 53477376

BOARD_USERDATAIMAGE_PARTITION_SIZE := 576716800

BOARD_RAMDISK_USE_LZ4 := true
BOARD_BOOT_HEADER_VERSION := 4
BOARD_MKBOOTIMG_ARGS += --header_version $(BOARD_BOOT_HEADER_VERSION)

# Enable GKI 2.0 signing.
BOARD_GKI_SIGNING_KEY_PATH := build/make/target/product/gsi/testkey_rsa2048.pem
BOARD_GKI_SIGNING_ALGORITHM := SHA256_RSA2048

BOARD_KERNEL_BINARIES := \
    kernel-4.19-gz \
    kernel-5.4 kernel-5.4-gz kernel-5.4-lz4 \
    kernel-5.10 kernel-5.10-gz kernel-5.10-lz4 \
    kernel-mainline kernel-mainline-gz kernel-mainline-lz4 \

ifneq (,$(filter userdebug eng,$(TARGET_BUILD_VARIANT)))
BOARD_KERNEL_BINARIES += \
    kernel-5.4-allsyms kernel-5.4-gz-allsyms kernel-5.4-lz4-allsyms \
    kernel-5.10-allsyms kernel-5.10-gz-allsyms kernel-5.10-lz4-allsyms \

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

# Some vendors still haven't cleaned up all device specific directories under
# root!

# TODO(b/111434759, b/111287060) SoC specific hacks
BOARD_ROOT_EXTRA_SYMLINKS += /vendor/lib/dsp:/dsp
BOARD_ROOT_EXTRA_SYMLINKS += /mnt/vendor/persist:/persist
BOARD_ROOT_EXTRA_SYMLINKS += /vendor/firmware_mnt:/firmware

# TODO(b/36764215): remove this setting when the generic system image
# no longer has QCOM-specific directories under /.
BOARD_SEPOLICY_DIRS += build/make/target/board/generic_arm64/sepolicy
