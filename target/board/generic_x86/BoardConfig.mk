# config.mk
#
# Product-specific compile-time definitions.
#

include build/make/target/board/treble_common_32.mk

# The emulator doesn't support sparsed images, yet.
# Overwrite the setting in treble_common.mk.
TARGET_USERIMAGES_SPARSE_EXT_DISABLED := true

# Remove the following once it's not set to true in treble_common.mk.
BOARD_VNDK_RUNTIME_DISABLE :=

TARGET_CPU_ABI := x86
TARGET_ARCH := x86
TARGET_ARCH_VARIANT := x86
TARGET_PRELINK_MODULE := false

# The IA emulator (qemu) uses the Goldfish devices
HAVE_HTC_AUDIO_DRIVER := true
BOARD_USES_GENERIC_AUDIO := true

# no hardware camera
USE_CAMERA_STUB := true

# Build OpenGLES emulation host and guest libraries
BUILD_EMULATOR_OPENGL := true

# Build partitioned system.img and vendor.img (if applicable)
# for qemu, otherwise, init cannot find PART_NAME
BUILD_QEMU_IMAGES := true

# Build and enable the OpenGL ES View renderer. When running on the emulator,
# the GLES renderer disables itself if host GL acceleration isn't available.
USE_OPENGL_RENDERER := true

BOARD_USERDATAIMAGE_PARTITION_SIZE := 576716800

# 128 MB vendor image. May need to be adjusted if not enough.
BOARD_VENDORIMAGE_PARTITION_SIZE := 134217728
BOARD_VENDORIMAGE_FILE_SYSTEM_TYPE := ext4

DEVICE_MATRIX_FILE   := device/generic/goldfish/compatibility_matrix.xml

BOARD_SEPOLICY_DIRS += \
        build/target/board/generic/sepolicy \
        build/target/board/generic_x86/sepolicy

BOARD_VNDK_VERSION := current

# Enable A/B update
TARGET_NO_RECOVERY := true
BOARD_BUILD_SYSTEM_ROOT_IMAGE := true
