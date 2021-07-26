$(call inherit-product, $(SRC_TARGET_DIR)/product/virtual_ab_ota/compression.mk)


PRODUCT_VENDOR_PROPERTIES += ro.virtual_ab.compression.xor.enabled=true
