# VR application inclusions for certified devices.

BOARD_SEPOLICY_DIRS += device/google/vrservices/vrcore/sepolicy

# Include all certified-device-only Google VR applications here.

# VR Services.
PRODUCT_PACKAGES += \
    bufferhubd \
    performanced \
    virtual_touchpad \
    vr_hwc \

# Diagnosis tool to debug pdx, b/64377040.
PRODUCT_PACKAGES_DEBUG += pdx_tool

# Used by surface flinger's ConfigStore to decide whether or not to start vr
# flinger services at system boot.
USE_VR_FLINGER := true

# Platform Library Configurations for VR Services.
PRODUCT_PACKAGES += \
    com.google.vr.platform \
    com.google.vr.platform.xml
PRODUCT_BOOT_JARS += com.google.vr.platform
