# Copyright (C) 2016 The Android Open Source Project
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

# Configuration for generic_x86 + arm libraries needed by binary translation.

# The generic product target doesn't have any hardware-specific pieces.
TARGET_NO_BOOTLOADER := true
TARGET_NO_KERNEL := true
TARGET_CPU_ABI := x86
TARGET_ARCH := x86
TARGET_ARCH_VARIANT := x86

<<<<<<< HEAD   (5c8d84 Merge "Merge empty history for sparse-6676661-L8360000065797)
TARGET_2ND_ARCH := arm
TARGET_2ND_CPU_ABI := armeabi-v7a
TARGET_2ND_CPU_ABI2 := armeabi
TARGET_2ND_ARCH_VARIANT := armv7-a
TARGET_2ND_CPU_VARIANT := generic

# Tell the build system this isn't a typical 64bit+32bit multilib configuration.
TARGET_TRANSLATE_2ND_ARCH := true
=======
TARGET_NATIVE_BRIDGE_ARCH := arm
TARGET_NATIVE_BRIDGE_ARCH_VARIANT := armv7-a-neon
TARGET_NATIVE_BRIDGE_CPU_VARIANT := generic
TARGET_NATIVE_BRIDGE_ABI := armeabi-v7a armeabi
>>>>>>> BRANCH (a10c18 Merge "Version bump to RT11.201014.001.A1 [core/build_id.mk])

# no hardware camera
USE_CAMERA_STUB := true

<<<<<<< HEAD   (5c8d84 Merge "Merge empty history for sparse-6676661-L8360000065797)
# Enable dex-preoptimization to speed up the first boot sequence
# of an SDK AVD. Note that this operation only works on Linux for now
ifeq ($(HOST_OS),linux)
  ifeq ($(WITH_DEXPREOPT),)
    WITH_DEXPREOPT := true
    WITH_DEXPREOPT_BOOT_IMG_AND_SYSTEM_SERVER_ONLY := false
  endif
endif

TARGET_USES_HWC2 := true
NUM_FRAMEBUFFER_SURFACE_BUFFERS := 3
=======
#
# The inclusion order below is important.
# The settings in latter makefiles overwrite those in the former.
#
include build/make/target/board/BoardConfigMainlineCommon.mk
include build/make/target/board/BoardConfigEmuCommon.mk
>>>>>>> BRANCH (a10c18 Merge "Version bump to RT11.201014.001.A1 [core/build_id.mk])

<<<<<<< HEAD   (5c8d84 Merge "Merge empty history for sparse-6676661-L8360000065797)
# Build OpenGLES emulation host and guest libraries
BUILD_EMULATOR_OPENGL := true
=======
# the settings differ from BoardConfigMainlineCommon.mk
BOARD_USES_SYSTEM_OTHER_ODEX :=

# Resize to 4G to accomodate ASAN and CTS
BOARD_USERDATAIMAGE_PARTITION_SIZE := 4294967296
>>>>>>> BRANCH (a10c18 Merge "Version bump to RT11.201014.001.A1 [core/build_id.mk])

# Build and enable the OpenGL ES View renderer. When running on the emulator,
# the GLES renderer disables itself if host GL acceleration isn't available.
USE_OPENGL_RENDERER := true

TARGET_USERIMAGES_USE_EXT4 := true
BOARD_SYSTEMIMAGE_PARTITION_SIZE := 1879048192  # 1.75 GB
BOARD_USERDATAIMAGE_PARTITION_SIZE := 576716800
BOARD_CACHEIMAGE_PARTITION_SIZE := 69206016
BOARD_CACHEIMAGE_FILE_SYSTEM_TYPE := ext4
BOARD_FLASH_BLOCK_SIZE := 512
TARGET_USERIMAGES_SPARSE_EXT_DISABLED := true

BOARD_SEPOLICY_DIRS += build/target/board/generic/sepolicy
