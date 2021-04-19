# BoardConfigModuleCommon.mk
#
# Common compile-time settings for module builds.

# Required for all module devices.
TARGET_USES_64_BIT_BINDER := true

# Necessary to make modules able to use the VNDK (which media.swcodec does).
BOARD_VNDK_VERSION := current
