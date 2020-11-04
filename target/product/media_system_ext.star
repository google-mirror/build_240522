load(":product_config.star", "prodconf")
load(":base_system_ext.star", "base_system_ext")

# /system_ext packages
media_system_ext = prodconf(
    "media_system_ext",
    [base_system_ext],
    PRODUCT_PACKAGES=[
        "vndk_apex_snapshot_package",
    ],
)
