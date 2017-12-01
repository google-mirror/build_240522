
ifdef BOARD_ENFORCE_SYSTEM_SDK
# Set LOCAL_SDK_VERSION to system_current, If LOCAL_SDK_VERSION is not defined and LOCAL_VENDOR_MODULE is true
  my_vendor_module :=
  ifneq (,$(filter true,$(LOCAL_VENDOR_MODULE) $(LOCAL_ODM_MODULE) $(LOCAL_PROPRIETARY_MODULE)))
    my_vendor_module := true
  endif
  ifneq (,$(filter JAVA_LIBRARIES APPS,$(LOCAL_MODULE_CLASS)))
    ifndef LOCAL_SDK_VERSION
      ifeq ($(my_vendor_module),true)
        LOCAL_SDK_VERSION := system_current
      endif
    endif
  endif
endif
