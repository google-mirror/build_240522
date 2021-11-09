#######################################################
## Shared definitions for all target test compilations.
#######################################################

ifeq ($(LOCAL_GTEST),true)
  LOCAL_CFLAGS += -DGTEST_OS_LINUX_ANDROID -DGTEST_HAS_STD_STRING

  ifndef LOCAL_SDK_VERSION
    LOCAL_STATIC_LIBRARIES += libgtest_main libgtest
  else
    # TODO(danalbert): Remove the suffix from the module since we only need the
    # one variant now.
    my_ndk_gtest_suffix := _c++
    LOCAL_STATIC_LIBRARIES += \
        libgtest_main_ndk$(my_ndk_gtest_suffix) \
        libgtest_ndk$(my_ndk_gtest_suffix)
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
<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
=======
endif

# Implicitly run this test under MTE SYNC for aarch64 binaries. This is a no-op
# on non-MTE hardware.
ifneq (,$(filter arm64,$(TARGET_$(LOCAL_2ND_ARCH_VAR_PREFIX)ARCH)))
  LOCAL_WHOLE_STATIC_LIBRARIES += note_memtag_heap_sync
endif
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
