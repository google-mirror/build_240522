my_prefix := TARGET_
include $(BUILD_SYSTEM)/multilib.mk

ifndef my_module_multilib
# libraries default to building for both architecturess
my_module_multilib := both
endif

ifeq ($(my_module_multilib),both)
ifneq ($(LOCAL_MODULE_PATH),)
$(warning $(LOCAL_MODULE): LOCAL_MODULE_PATH for static and shared libraries is unsupported, use LOCAL_MODULE_RELATIVE_PATH instead)
endif

ifneq ($(LOCAL_UNSTRIPPED_PATH),)
$(warning $(LOCAL_MODULE): LOCAL_UNSTRIPPED_PATH for static and shared libraries is unsupported)
endif
endif # my_module_multilib == both


LOCAL_2ND_ARCH_VAR_PREFIX :=
include $(BUILD_SYSTEM)/module_arch_supported.mk

ifeq ($(my_module_arch_supported),true)
include $(BUILD_SYSTEM)/static_and_shared_library_internal.mk
endif

ifdef TARGET_2ND_ARCH

LOCAL_2ND_ARCH_VAR_PREFIX := $(TARGET_2ND_ARCH_VAR_PREFIX)
include $(BUILD_SYSTEM)/module_arch_supported.mk

ifeq ($(my_module_arch_supported),true)
# Build for TARGET_2ND_ARCH
include $(BUILD_SYSTEM)/static_and_shared_library_internal.mk

endif

LOCAL_2ND_ARCH_VAR_PREFIX :=

endif # TARGET_2ND_ARCH

my_module_arch_supported :=

###########################################################
## Copy headers to the install tree
###########################################################
include $(BUILD_COPY_HEADERS)
