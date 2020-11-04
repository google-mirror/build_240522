load(":product_config.star", "prodconf")

# load(":handheld_system.star","handheld_system")
load(":telephony_system.star", "telephony_system")
# load(":languages_default.star", "languages_default")
# Add adb keys to debuggable AOSP builds (if they exist)
# TODO(asmundak):
# $(call inherit-product-if-exists, vendor/google/security/adb/vendor_key.mk)

# Enable updating of APEXes
# load(":updatable_apex.star", "updatable_apex")

generic_system = prodconf(
    "generic_system",
    [telephony_system],
    PRODUCT_PACKAGES=[
        # Shared java libs
        "com.android.nfc_extras",
        # Applications
        "LiveWallpapersPicker",
        "PartnerBookmarksProvider",
        "PresencePolling",
        "RcsService",
        "Stk",
        "Tag",
        "TimeZoneUpdater",

        # OTA support
        "recovery-refresh",
        "update_engine",
        "update_verifier",

        # Wrapped net utils for /vendor access.
        "netutils-wrapper-1.0",

        # Charger images
        "charger_res_images",

        # system_other support
        "cppreopts.sh",
        "otapreopt_script",

        # Bluetooth libraries
        "audio.a2dp.default",
        "audio.hearing_aid.default",

        # For ringtones that rely on forward lock encryption
        "libfwdlockengine",

        # System libraries commonly depended on by things on the system_ext or product partitions.
        # These lists will be pruned periodically.
        "android.hardware.biometrics.fingerprint@2.1",
        "android.hardware.radio@1.0",
        "android.hardware.radio@1.1",
        "android.hardware.radio@1.2",
        "android.hardware.radio@1.3",
        "android.hardware.radio@1.4",
        "android.hardware.radio.config@1.0",
        "android.hardware.radio.deprecated@1.0",
        "android.hardware.secure_element@1.0",
        "android.hardware.wifi@1.0",
        "libaudio-resampler",
        "libaudiohal",
        "libdrm",
        "liblogwrap",
        "liblz4",
        "libminui",
        "libnl",
        "libprotobuf-cpp-full",

        # These libraries are empty and have been combined into libhidlbase, but are still depended
        # on by things off /system.
        # TODO(b/135686713): remove these
        "libhidltransport",
        "libhwbinder",
        # Enable configurable audio policy
        "libaudiopolicyengineconfigurable",
        "libpolicy-subsystem",
    ],

    PRODUCT_PACKAGES_DEBUG=[
        "avbctl",
        "bootctl",
        "tinycap",
        "tinyhostless",
        "tinymix",
        "tinypcminfo",
        "tinyplay",
        "update_engine_client",
    ],

    PRODUCT_HOST_PACKAGES=[
        "tinyplay",
    ],

    # Include all zygote init scripts. "ro.zygote" will select one of them.
    PRODUCT_COPY_FILES=[
        "system/core/rootdir/init.zygote32.rc:system/etc/init/hw/init.zygote32.rc",
        "system/core/rootdir/init.zygote64.rc:system/etc/init/hw/init.zygote64.rc",
        "system/core/rootdir/init.zygote64_32.rc:system/etc/init/hw/init.zygote64_32.rc",
    ],

    # Enable dynamic partition size
    PRODUCT_USE_DYNAMIC_PARTITION_SIZE=True,
    PRODUCT_ENFORCE_RRO_TARGETS='*',

    # TODO(b/150820813) Settings depends on static overlay, remove this after eliminating the dependency.
    PRODUCT_ENFORCE_RRO_EXEMPTED_TARGETS="Settings",

    PRODUCT_NAME='generic_system',
    PRODUCT_BRAND='generic',

    # Define /system partition-specific product properties to identify that /system
    # partition is generic_system.
    PRODUCT_SYSTEM_NAME='mainline',
    PRODUCT_SYSTEM_BRAND='Android',
    PRODUCT_SYSTEM_MANUFACTURER='Android',
    PRODUCT_SYSTEM_MODEL='mainline',
    PRODUCT_SYSTEM_DEVICE='generic',
)

## TODO(Asmundak):
##_base_mk_allowed_list :=
##_my_allowed_list := $(_base_mk_allowed_list)

# For mainline, system.img should be mounted at /, so we include ROOT here.
##_my_paths := $(TARGET_COPY_OUT_ROOT)/ $(TARGET_COPY_OUT_SYSTEM)/
## $(call require-artifacts-in-path, $(_my_paths), $(_my_allowed_list))
