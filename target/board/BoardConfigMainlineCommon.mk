# BoardConfigGsiCommon.mk
#
# Common compile-time definitions for mainline images.

# The generic product target doesn't have any hardware-specific pieces.
TARGET_NO_BOOTLOADER := true
TARGET_NO_KERNEL := true

TARGET_USERIMAGES_USE_EXT4 := true

# system-as-root is mandatory from Android P
TARGET_NO_RECOVERY := true
BOARD_BUILD_SYSTEM_ROOT_IMAGE := true

# Puts odex files on system_other, as well as causing dex files not to get
# stripped from APKs.
BOARD_USES_SYSTEM_OTHER_ODEX := true

# Audio: must using XML format for Treblized devices
USE_XML_AUDIO_POLICY_CONF := 1

# Enable stats logging in LMKD.
TARGET_LMKD_LOG_STATS := true

# Controls some sync timings in libhwui.
TARGET_USES_HWC2 := true
