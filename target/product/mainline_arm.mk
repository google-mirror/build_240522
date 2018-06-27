
$(call inherit-product, $(SRC_TARGET_DIR)/product/mainline_system.mk)
$(call inherit-product, $(SRC_TARGET_DIR)/product/base_vendor.mk)

PRODUCT_NAME := mainline_arm
PRODUCT_DEVICE := generic
PRODUCT_BRAND := generic
PRODUCT_SHIPPING_API_LEVEL := 28

PRODUCT_ENFORCE_ARTIFACT_PATH_REQUIREMENTS := true
PRODUCT_ARTIFACT_PATH_REQUIREMENT_WHITELIST := system/etc/seccomp_policy/mediacodec.policy
