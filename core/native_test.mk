###########################################
## A thin wrapper around BUILD_EXECUTABLE
## Common flags for native tests are added.
###########################################
$(call record-module-type,NATIVE_TEST)

ifdef LOCAL_MODULE_CLASS
ifneq ($(LOCAL_MODULE_CLASS),NATIVE_TESTS)
$(error $(LOCAL_PATH): LOCAL_MODULE_CLASS must be NATIVE_TESTS with BUILD_HOST_NATIVE_TEST)
endif
endif

# Implicitly run this test under MTE SYNC for aarch64 binaries. This is a no-op
# on non-MTE hardware.
my_arch := $(TARGET_$(LOCAL_2ND_ARCH_VAR_PREFIX)ARCH)
ifneq (,$(filter arm64,$(my_arch)))
	LOCAL_WHOLE_STATIC_LIBRARIES += note_memtag_heap_sync
endif

LOCAL_MODULE_CLASS := NATIVE_TESTS

include $(BUILD_SYSTEM)/target_test_internal.mk

ifndef LOCAL_MULTILIB
ifndef LOCAL_32_BIT_ONLY
LOCAL_MULTILIB := both
endif
endif

include $(BUILD_EXECUTABLE)
