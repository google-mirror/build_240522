# This prodcut config Starlark contains the non-system partition contents for
# media-capable devices (non-wearables). Only add something here
# if it definitely doesn't belong on wearables. Otherwise, choose
# base_vendor.mk.
load(":product_config.star", "prodconf")

load(":base_vendor.star", "base_vendor")

# /vendor packages
media_vendor = prodconf(
    "media_vendor",
    [base_vendor],
    PRODUCT_PACKAGES = [
        "libaudiopreprocessing",
        "libwebrtc_audio_preprocessing",
    ],
    )
