load(":product_config.star", "prodconf")

# Includes all AOSP product packages
load(":handheld_product.star", "handheld_product")
load(":telephony_product.star", "telephony_product")

# Default AOSP sounds.
# TODO(asmundak): why there is inherit-product-if-exists?
#$(call inherit-product-if-exists, frameworks/base/data/sounds/AllAudio.mk)
#load("//frameworks/base/data/sounds/AllAudio.star", "AllAudio")

aosp_product = prodconf(
    "aosp_product",
    [handheld_product, telephony_product],
    # Additional settings used in all AOSP builds
    PRODUCT_PRODUCT_PROPERTIES = [
        # TODO(asmundak): what are those settings and what does '?=' mean
        "ro.config.ringtone?=Ring_Synth_04.ogg",
        "ro.config.notification_sound?=pixiedust.ogg",
        "ro.com.android.dataroaming?=true",
    ],

    # More AOSP packages
    PRODUCT_PACKAGES = [
        "messaging",
        "PhotoTable",
        "preinstalled-packages-platform-aosp-product.xml",
        "WallpaperPicker",
    ],

    # Telephony:
    #   Provide a APN configuration to GSI product
    PRODUCT_COPY_FILES = [
        "device/sample/etc/apns-full-conf.xml:$(TARGET_COPY_OUT_PRODUCT)/etc/apns-conf.xml",
    ],
)
