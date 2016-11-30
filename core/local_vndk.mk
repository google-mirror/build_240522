
#If LOCAL_USE_VNDK is not set, set it if proprietary/odm/oem flag is true or module is going into vendor partition
ifndef LOCAL_USE_VNDK
  ifneq (,$(filter true,$(LOCAL_PROPRIETARY_MODULE)$(LOCAL_ODM_MODULE)$(LOCAL_OEM_MODULE)))
    LOCAL_USE_VNDK:=true
  else
    ifneq (,$(filter $(TARGET_OUT_VENDOR),$(LOCAL_MODULE_PATH)))
      LOCAL_USE_VNDK:=true
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
  endif
endif

