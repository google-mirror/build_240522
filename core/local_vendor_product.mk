
#Set LOCAL_IN_VENDOR_OR_PRODUCT for modules going into product, vendor or odm partition, except for host modules
#If LOCAL_SDK_VERSION is set, thats a more restrictive set, so they dont need LOCAL_IN_VENDOR_OR_PRODUCT
ifndef LOCAL_IS_HOST_MODULE
ifndef LOCAL_SDK_VERSION
  ifneq (,$(filter true,$(LOCAL_VENDOR_MODULE) $(LOCAL_ODM_MODULE) $(LOCAL_OEM_MODULE) $(LOCAL_PROPRIETARY_MODULE)))
    LOCAL_IN_VENDOR_OR_PRODUCT:=true
    LOCAL_IN_VENDOR:=true
    # Note: no need to check LOCAL_MODULE_PATH* since LOCAL_[VENDOR|ODM|OEM]_MODULE is already
    # set correctly before this is included.
  endif
  ifdef PRODUCT_PRODUCT_VNDK_VERSION
    # Product modules also use VNDK when PRODUCT_PRODUCT_VNDK_VERSION is defined.
    ifeq (true,$(LOCAL_PRODUCT_MODULE))
      LOCAL_IN_VENDOR_OR_PRODUCT:=true
      LOCAL_IN_PRODUCT:=true
    endif
  endif
endif
endif

# Verify LOCAL_IN_VENDOR_OR_PRODUCT usage, and set LOCAL_SDK_VERSION if necessary

ifdef LOCAL_IS_HOST_MODULE
  ifdef LOCAL_IN_VENDOR_OR_PRODUCT
    $(shell echo $(LOCAL_MODULE_MAKEFILE): $(LOCAL_MODULE): Do not use LOCAL_IN_VENDOR_OR_PRODUCT with host modules >&2)
    $(error done)
  endif
endif
ifdef LOCAL_IN_VENDOR_OR_PRODUCT
  ifneq ($(LOCAL_IN_VENDOR_OR_PRODUCT),true)
    $(shell echo '$(LOCAL_MODULE_MAKEFILE): $(LOCAL_MODULE): LOCAL_IN_VENDOR_OR_PRODUCT must be "true" or empty, not "$(LOCAL_IN_VENDOR_OR_PRODUCT)"' >&2)
    $(error done)
  endif

  ifdef LOCAL_SDK_VERSION
    $(shell echo $(LOCAL_MODULE_MAKEFILE): $(LOCAL_MODULE): LOCAL_IN_VENDOR_OR_PRODUCT must not be used with LOCAL_SDK_VERSION >&2)
    $(error done)
  endif
endif

