#############################################################################
## Standard rules for installing runtime resouce overlay APKs.
##
## Set LOCAL_RRO_THEME to the theme name if the package should apply only to
## a particular theme as set by ro.boot.vendor.overlay.theme system property.
##
## If LOCAL_RRO_THEME is not set, the package will apply always, independent
## of themes.
##
#############################################################################

LOCAL_IS_RUNTIME_RESOURCE_OVERLAY := true

ifneq ($(LOCAL_SRC_FILES),)
  $(error runtime resource overlay package should not contain sources)
endif

ifeq ($(LOCAL_VENDOR_MODULE),true)
  LOCAL_MODULE_PATH := $(TARGET_OUT_VENDOR)
else ifeq ($(LOCAL_PRODUCT_MODULE),true)
  LOCAL_MODULE_PATH := $(TARGET_OUT_PRODUCT)
else ifeq ($(LOCAL_PRODUCT_SERVICES_MODULE),true)
  LOCAL_MODULE_PATH := $(TARGET_OUT_PRODUCT_SERVICES)
else
  LOCAL_MODULE_PATH := $(TARGET_OUT_VENDOR)
endif

ifeq ($(LOCAL_RRO_THEME),)
  LOCAL_MODULE_PATH := $(LOCAL_MODULE_PATH)/overlay
else
  LOCAL_MODULE_PATH := $(LOCAL_MODULE_PATH)/overlay/$(LOCAL_RRO_THEME)
endif

include $(BUILD_SYSTEM)/package.mk

