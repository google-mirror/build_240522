load(":product_config.star", "prodconf")

telephony_system = prodconf(
    "telephony_system",
    [],
    PRODUCT_PACKAGES=[
        "ONS",
        "CarrierDefaultApp",
        "CallLogBackup",
        "com.android.cellbroadcast",
        "CellBroadcastLegacyApp",
    ],
)
