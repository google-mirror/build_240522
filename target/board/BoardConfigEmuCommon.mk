# BoardConfigEmuCommon.mk
#
# Common compile-time definitions for emulator
#

HAVE_HTC_AUDIO_DRIVER := true
BOARD_USES_GENERIC_AUDIO := true
TARGET_BOOTLOADER_BOARD_NAME := goldfish_$(TARGET_ARCH)

# no hardware camera
USE_CAMERA_STUB := true

NUM_FRAMEBUFFER_SURFACE_BUFFERS := 3

# Build OpenGLES emulation guest and host libraries
BUILD_EMULATOR_OPENGL := true
BUILD_QEMU_IMAGES := true

# Build and enable the OpenGL ES View renderer. When running on the emulator,
# the GLES renderer disables itself if host GL acceleration isn't available.
USE_OPENGL_RENDERER := true

ifeq ($(PRODUCT_USE_DYNAMIC_PARTITIONS),true)
BOARD_EXT4_SHARE_DUP_BLOCKS := true
# The super partition size is more than 2x of dynamic partition size for
# 1. lpmake requirement: size-of-super > 2 * size-of-dynamic-partitions
# 2. provide scratch spaces for adb remount to update system/vendor
BOARD_SUPER_PARTITION_SIZE := 8145338368
BOARD_SUPER_PARTITION_GROUPS := emulator_dynamic_partitions
BOARD_EMULATOR_DYNAMIC_PARTITIONS_PARTITION_LIST := \
    system \
    vendor

BOARD_EMULATOR_DYNAMIC_PARTITIONS_SIZE := 4068474880
else
# ~140 MB vendor image. Please adjust system image / vendor image sizes
# when finalizing them. The partition size needs to be a multiple of image
# block size: 4096.
BOARD_VENDORIMAGE_PARTITION_SIZE := 140963840
# TODO(b/125540538): Remove when emulator uses dynamic partitions
BOARD_BUILD_SYSTEM_ROOT_IMAGE := true
endif

BOARD_VENDORIMAGE_FILE_SYSTEM_TYPE := ext4
BOARD_FLASH_BLOCK_SIZE := 512
DEVICE_MATRIX_FILE   := device/generic/goldfish/compatibility_matrix.xml

BOARD_SEPOLICY_DIRS += device/generic/goldfish/sepolicy/common
