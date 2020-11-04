# This is the list of modules that are specific to products that have telephony
# hardware, and install to the product partition.
load(":product_config.star", "prodconf")

# /product packages
telephony_product = prodconf(
    "telephony_product",
    [],
    PRODUCT_PACKAGES=[
        "Dialer",
    ]
)
