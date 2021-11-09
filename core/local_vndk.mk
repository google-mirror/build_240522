
#Set LOCAL_USE_VNDK for modules going into vendor partition, except for host modules
#If LOCAL_SDK_VERSION is set, thats a more restrictive set, so they dont need LOCAL_USE_VNDK
ifndef LOCAL_IS_HOST_MODULE
ifndef LOCAL_SDK_VERSION
  ifneq (,$(filter true,$(LOCAL_VENDOR_MODULE) $(LOCAL_ODM_MODULE) $(LOCAL_OEM_MODULE) $(LOCAL_PROPRIETARY_MODULE)))
    LOCAL_USE_VNDK:=true
    LOCAL_USE_VNDK_VENDOR:=true
    # Note: no need to check LOCAL_MODULE_PATH* since LOCAL_[VENDOR|ODM|OEM]_MODULE is already
    # set correctly before this is included.
  endif
endif
endif

# Verify LOCAL_USE_VNDK usage, and set LOCAL_SDK_VERSION if necessary

ifdef LOCAL_IS_HOST_MODULE
  ifdef LOCAL_USE_VNDK
    $(shell echo $(LOCAL_MODULE_MAKEFILE): $(LOCAL_MODULE): Do not use LOCAL_USE_VNDK with host modules >&2)
    $(error done)
  endif
endif
ifdef LOCAL_USE_VNDK
  ifneq ($(LOCAL_USE_VNDK),true)
    $(shell echo '$(LOCAL_MODULE_MAKEFILE): $(LOCAL_MODULE): LOCAL_USE_VNDK must be "true" or empty, not "$(LOCAL_USE_VNDK)"' >&2)
    $(error done)
  endif

  ifdef LOCAL_SDK_VERSION
    $(shell echo $(LOCAL_MODULE_MAKEFILE): $(LOCAL_MODULE): LOCAL_USE_VNDK must not be used with LOCAL_SDK_VERSION >&2)
    $(error done)
  endif

  # If we're not using the VNDK, drop all restrictions
  ifndef BOARD_VNDK_VERSION
    LOCAL_USE_VNDK:=
<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
=======
    LOCAL_USE_VNDK_VENDOR:=
    LOCAL_USE_VNDK_PRODUCT:=
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
  endif
endif

