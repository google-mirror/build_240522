#######################################################
## Shared definitions for all target test compilations.
#######################################################

ifeq ($(LOCAL_GTEST),true)
  LOCAL_CFLAGS += -DGTEST_OS_LINUX_ANDROID -DGTEST_HAS_STD_STRING

  ifndef LOCAL_SDK_VERSION
    LOCAL_STATIC_LIBRARIES += libgtest_main libgtest
  else
    LOCAL_STATIC_LIBRARIES += libgtest_main_ndk libgtest_ndk
  endif
endif

ifdef LOCAL_MODULE_PATH
$(error $(LOCAL_PATH): Do not set LOCAL_MODULE_PATH when building test $(LOCAL_MODULE))
endif

ifdef LOCAL_MODULE_PATH_32
$(error $(LOCAL_PATH): Do not set LOCAL_MODULE_PATH_32 when building test $(LOCAL_MODULE))
endif

ifdef LOCAL_MODULE_PATH_64
$(error $(LOCAL_PATH): Do not set LOCAL_MODULE_PATH_64 when building test $(LOCAL_MODULE))
endif

ifndef LOCAL_MODULE_RELATIVE_PATH
LOCAL_MODULE_RELATIVE_PATH := $(LOCAL_MODULE)
endif
