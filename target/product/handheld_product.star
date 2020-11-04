# This Starlark config file contains the product partition contents for
# a generic phone or tablet device. Only add something here if
# it definitely doesn't belong on other types of devices (if it
# does, use base_product.mk).
load(":product_config.star", "prodconf")

load(":media_product.star", "media_product")

# /product packages
handheld_product = prodconf(
    "handheld_product",
    [media_product],
    PRODUCT_PACKAGES=[
        "Browser2",
        "Calendar",
        "Camera2",
        "Contacts",
        "DeskClock",
        "Gallery2",
        "LatinIME",
        "Music",
        "OneTimeInitializer",
        "preinstalled-packages-platform-handheld-product.xml",
        "QuickSearchBox",
        "SettingsIntelligence",
        "frameworks-base-overlays",
    ],
    PRODUCT_PACKAGES_DEBUG=[
        "frameworks-base-overlays-debug",
    ],
)
