# Copyright (C) 2013 The Android Open Source Project
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

# The generic product target doesn't have any hardware-specific pieces.
TARGET_NO_BOOTLOADER := true
TARGET_NO_KERNEL := true
TARGET_ARCH := arm64
TARGET_ARCH_VARIANT := armv8-a
TARGET_CPU_VARIANT := generic
TARGET_CPU_ABI := arm64-v8a
TARGET_BOOTLOADER_BOARD_NAME := goldfish_$(TARGET_ARCH)

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
TARGET_2ND_ARCH_VARIANT := armv7-a
# DO NOT USE
# DO NOT USE
TARGET_2ND_CPU_VARIANT := generic
# DO NOT USE
# DO NOT USE
else
TARGET_2ND_ARCH_VARIANT := armv8-a
TARGET_2ND_CPU_VARIANT := generic
endif


<<<<<<< HEAD   (4be654 Merge "Merge empty history for sparse-7121469-L4290000080720)
TARGET_USES_64_BIT_BINDER := true

# no hardware camera
USE_CAMERA_STUB := true

TARGET_USES_HWC2 := true
NUM_FRAMEBUFFER_SURFACE_BUFFERS := 3

# Build OpenGLES emulation host and guest libraries
BUILD_EMULATOR_OPENGL := true
BUILD_QEMU_IMAGES := true

# Build and enable the OpenGL ES View renderer. When running on the emulator,
# the GLES renderer disables itself if host GL acceleration isn't available.
USE_OPENGL_RENDERER := true

TARGET_USERIMAGES_USE_EXT4 := true
# Partition size is default 1.5GB (1536MB) for 64 bits projects
BOARD_SYSTEMIMAGE_PARTITION_SIZE := 1610612736
=======
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

>>>>>>> BRANCH (fe6ad7 Merge "Version bump to RBT1.210107.001.A1 [core/build_id.mk])
BOARD_USERDATAIMAGE_PARTITION_SIZE := 576716800
TARGET_COPY_OUT_VENDOR := vendor
# ~100 MB vendor image. Please adjust system image / vendor image sizes
# when finalizing them.
BOARD_VENDORIMAGE_PARTITION_SIZE := 100000000
BOARD_VENDORIMAGE_FILE_SYSTEM_TYPE := ext4
BOARD_FLASH_BLOCK_SIZE := 512
TARGET_USERIMAGES_SPARSE_EXT_DISABLED := true
DEVICE_MATRIX_FILE   := device/generic/goldfish/compatibility_matrix.xml

# Android generic system image always create metadata partition
BOARD_USES_METADATA_PARTITION := true

<<<<<<< HEAD   (4be654 Merge "Merge empty history for sparse-7121469-L4290000080720)
# Set this to create /cache mount point for non-A/B devices that mounts /cache.
# The partition size doesn't matter, just to make build pass.
BOARD_CACHEIMAGE_FILE_SYSTEM_TYPE := ext4
BOARD_CACHEIMAGE_PARTITION_SIZE := 16777216
=======
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
BOARD_KERNEL_MODULE_INTERFACE_VERSIONS := \
    5.4-android12-0 \
    5.10-android12-0 \

# Copy boot image in $OUT to target files. This is defined for targets where
# the installed GKI APEXes are built from source.
BOARD_COPY_BOOT_IMAGE_TO_TARGET_FILES := true

# No vendor_boot
BOARD_MOVE_RECOVERY_RESOURCES_TO_VENDOR_BOOT :=

# No recovery
BOARD_EXCLUDE_KERNEL_FROM_RECOVERY_IMAGE :=
>>>>>>> BRANCH (fe6ad7 Merge "Version bump to RBT1.210107.001.A1 [core/build_id.mk])

BOARD_PROPERTY_OVERRIDES_SPLIT_ENABLED := true
BOARD_SEPOLICY_DIRS += build/target/board/generic/sepolicy

# Android Verified Boot (AVB):
#   Builds a special vbmeta.img that disables AVB verification.
#   Otherwise, AVB will prevent the device from booting the generic system.img.
#   Also checks that BOARD_AVB_ENABLE is not set, to prevent adding verity
#   metadata into system.img.
ifeq ($(BOARD_AVB_ENABLE),true)
$(error BOARD_AVB_ENABLE cannot be set for GSI)
endif
BOARD_BUILD_DISABLED_VBMETAIMAGE := true

ifneq (,$(filter userdebug eng,$(TARGET_BUILD_VARIANT)))
# GSI is always userdebug and needs a couple of properties taking precedence
# over those set by the vendor.
TARGET_SYSTEM_PROP := build/make/target/board/gsi_system.prop
endif
BOARD_VNDK_VERSION := current

# Emulator system image is going to be used as GSI and some vendor still hasn't
# cleaned up all device specific directories under root!

# TODO(jiyong) These might be SoC specific.
BOARD_ROOT_EXTRA_FOLDERS += firmware firmware/radio persist
BOARD_ROOT_EXTRA_SYMLINKS := /vendor/lib/dsp:/dsp

# TODO(b/36764215): remove this setting when the generic system image
# no longer has QCOM-specific directories under /.
BOARD_SEPOLICY_DIRS += build/target/board/generic_arm64_ab/sepolicy

# Wifi.
BOARD_WLAN_DEVICE           := emulator
BOARD_HOSTAPD_DRIVER        := NL80211
BOARD_WPA_SUPPLICANT_DRIVER := NL80211
BOARD_HOSTAPD_PRIVATE_LIB   := lib_driver_cmd_simulated
BOARD_WPA_SUPPLICANT_PRIVATE_LIB := lib_driver_cmd_simulated
WPA_SUPPLICANT_VERSION      := VER_0_8_X
WIFI_DRIVER_FW_PATH_PARAM   := "/dev/null"
WIFI_DRIVER_FW_PATH_STA     := "/dev/null"
WIFI_DRIVER_FW_PATH_AP      := "/dev/null"

# Enable A/B update
TARGET_NO_RECOVERY := true
BOARD_BUILD_SYSTEM_ROOT_IMAGE := true
