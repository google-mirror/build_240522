# This makefile contains the system_ext partition contents for
# a generic phone or tablet device. Only add something here if
# it definitely doesn't belong on other types of devices (if it
# does, use base_system_ext.mk).
load(":product_config.star", "prodconf")

load(":media_system_ext.star", "media_system_ext")

# /system_ext packages
handheld_system_ext = prodconf(
    "handheld_system_ext",
    [media_system_ext],
    PRODUCT_PACKAGES=[
        "Launcher3QuickStep",
        "Provision",
        "Settings",
        "StorageManager",
        "SystemUI",
        "WalpaperCropper",
    ],
)
