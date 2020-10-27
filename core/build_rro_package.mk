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

<<<<<<< HEAD   (5c8d84 Merge "Merge empty history for sparse-6676661-L8360000065797)
=======
partition :=
ifeq ($(strip $(LOCAL_ODM_MODULE)),true)
  partition := $(TARGET_OUT_ODM)
else ifeq ($(strip $(LOCAL_VENDOR_MODULE)),true)
  partition := $(TARGET_OUT_VENDOR)
else ifeq ($(strip $(LOCAL_SYSTEM_EXT_MODULE)),true)
  partition := $(TARGET_OUT_SYSTEM_EXT)
else
  partition := $(TARGET_OUT_PRODUCT)
endif

>>>>>>> BRANCH (a10c18 Merge "Version bump to RT11.201014.001.A1 [core/build_id.mk])
ifeq ($(LOCAL_RRO_THEME),)
  LOCAL_MODULE_PATH := $(TARGET_OUT_VENDOR)/overlay
else
  LOCAL_MODULE_PATH := $(TARGET_OUT_VENDOR)/overlay/$(LOCAL_RRO_THEME)
endif

<<<<<<< HEAD   (5c8d84 Merge "Merge empty history for sparse-6676661-L8360000065797)
=======
# Do not remove resources without default values nor dedupe resource
# configurations with the same value
LOCAL_AAPT_FLAGS += \
    --no-resource-deduping \
    --no-resource-removal

partition :=

>>>>>>> BRANCH (a10c18 Merge "Version bump to RT11.201014.001.A1 [core/build_id.mk])
include $(BUILD_SYSTEM)/package.mk

