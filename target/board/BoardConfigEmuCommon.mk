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

# Emulator doesn't support sparse image format.
TARGET_USERIMAGES_SPARSE_EXT_DISABLED := true

ifeq ($(PRODUCT_USE_DYNAMIC_PARTITIONS),true)
# emulator is Non-A/B device
AB_OTA_UPDATER := false

# emulator needs super.img
BOARD_BUILD_SUPER_IMAGE_BY_DEFAULT := true

BOARD_EXT4_SHARE_DUP_BLOCKS := true

# 3G + header
BOARD_SUPER_PARTITION_SIZE := 3229614080
BOARD_SUPER_PARTITION_GROUPS := emulator_dynamic_partitions
BOARD_EMULATOR_DYNAMIC_PARTITIONS_PARTITION_LIST := \
    system \
    vendor

# 3G
BOARD_EMULATOR_DYNAMIC_PARTITIONS_SIZE := 3221225472
else
# ~140 MB vendor image. Please adjust system image / vendor image sizes
# when finalizing them. The partition size needs to be a multiple of image
# block size: 4096.
BOARD_VENDORIMAGE_PARTITION_SIZE := 140963840
# for older devices that do not support dynamic partition
BOARD_BUILD_SYSTEM_ROOT_IMAGE := true
endif

BOARD_VENDORIMAGE_FILE_SYSTEM_TYPE := ext4
BOARD_FLASH_BLOCK_SIZE := 512
DEVICE_MATRIX_FILE   := device/generic/goldfish/compatibility_matrix.xml

BOARD_SEPOLICY_DIRS += device/generic/goldfish/sepolicy/common
