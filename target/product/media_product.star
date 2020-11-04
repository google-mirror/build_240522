# This starlark config file contains the product partition contents for
# media-capable devices (non-wearables). Only add something here
# if it definitely doesn't belong on wearables. Otherwise, choose
# base_product.star.
load(":product_config.star", "prodconf")
load(":base_product.star", "base_product")

# /product packages
media_product = prodconf(
    "media_product",
    [base_product],
    PRODUCT_PACKAGES=[
        "webview",
    ],
)
