my_prefix := TARGET_
include $(BUILD_SYSTEM)/multilib.mk

my_sanitize_target := $(strip $(SANITIZE_TARGET))

ifndef my_module_multilib
# libraries default to building for both architecturess
my_module_multilib := both
endif

LOCAL_2ND_ARCH_VAR_PREFIX :=
include $(BUILD_SYSTEM)/module_arch_supported.mk

ifeq ($(my_module_arch_supported),true)
include $(BUILD_SYSTEM)/static_library_internal.mk
endif

ifdef TARGET_2ND_ARCH

LOCAL_2ND_ARCH_VAR_PREFIX := $(TARGET_2ND_ARCH_VAR_PREFIX)
include $(BUILD_SYSTEM)/module_arch_supported.mk

ifeq ($(my_module_arch_supported),true)
# Build for TARGET_2ND_ARCH
OVERRIDE_BUILT_MODULE_PATH :=
LOCAL_BUILT_MODULE :=
LOCAL_INSTALLED_MODULE :=
LOCAL_MODULE_STEM :=
LOCAL_BUILT_MODULE_STEM :=
LOCAL_INSTALLED_MODULE_STEM :=
LOCAL_INTERMEDIATE_TARGETS :=

include $(BUILD_SYSTEM)/static_library_internal.mk

endif

LOCAL_2ND_ARCH_VAR_PREFIX :=

endif # TARGET_2ND_ARCH


ifneq ($(filter $(my_sanitize_target),$(LOCAL_SUPPORTED_SANITIZERS)),)
  LOCAL_MODULE := $(LOCAL_MODULE)_$(my_sanitize_target)
  LOCAL_CLANG := true
  LOCAL_SANITIZE := $(my_sanitize_target)

  LOCAL_2ND_ARCH_VAR_PREFIX :=
  include $(BUILD_SYSTEM)/module_arch_supported.mk

  ifeq ($(my_module_arch_supported),true)
    OVERRIDE_BUILT_MODULE_PATH :=
    LOCAL_BUILT_MODULE :=
    LOCAL_INSTALLED_MODULE :=
    LOCAL_MODULE_STEM :=
    LOCAL_BUILT_MODULE_STEM :=
    LOCAL_INSTALLED_MODULE_STEM :=
    LOCAL_INTERMEDIATE_TARGETS :=
    include $(BUILD_SYSTEM)/static_library_internal.mk
  endif

  ifdef TARGET_2ND_ARCH
    LOCAL_2ND_ARCH_VAR_PREFIX := $(TARGET_2ND_ARCH_VAR_PREFIX)
    include $(BUILD_SYSTEM)/module_arch_supported.mk

    ifeq ($(my_module_arch_supported),true)
      # Build for TARGET_2ND_ARCH
      OVERRIDE_BUILT_MODULE_PATH :=
      LOCAL_BUILT_MODULE :=
      LOCAL_INSTALLED_MODULE :=
      LOCAL_MODULE_STEM :=
      LOCAL_BUILT_MODULE_STEM :=
      LOCAL_INSTALLED_MODULE_STEM :=
      LOCAL_INTERMEDIATE_TARGETS :=

      include $(BUILD_SYSTEM)/static_library_internal.mk
    endif
  LOCAL_2ND_ARCH_VAR_PREFIX :=
  endif # TARGET_2ND_ARCH
else ifneq ($(my_sanitize_target),)
  LOCAL_SRC_FILES :=
  LOCAL_STATIC_LIBRARIES :=
  LOCAL_WHOLE_STATIC_LIBRARIES := $(LOCAL_MODULE)
  LOCAL_SHARED_LIBRARIES :=
  LOCAL_MODULE := $(LOCAL_MODULE)_$(my_sanitize_target)

  LOCAL_2ND_ARCH_VAR_PREFIX :=
  include $(BUILD_SYSTEM)/module_arch_supported.mk

  ifeq ($(my_module_arch_supported),true)
    OVERRIDE_BUILT_MODULE_PATH :=
    LOCAL_BUILT_MODULE :=
    LOCAL_INSTALLED_MODULE :=
    LOCAL_MODULE_STEM :=
    LOCAL_BUILT_MODULE_STEM :=
    LOCAL_INSTALLED_MODULE_STEM :=
    LOCAL_INTERMEDIATE_TARGETS :=
    include $(BUILD_SYSTEM)/static_library_internal.mk
  endif

  ifdef TARGET_2ND_ARCH
    LOCAL_2ND_ARCH_VAR_PREFIX := $(TARGET_2ND_ARCH_VAR_PREFIX)
    include $(BUILD_SYSTEM)/module_arch_supported.mk

    ifeq ($(my_module_arch_supported),true)
      # Build for TARGET_2ND_ARCH
      OVERRIDE_BUILT_MODULE_PATH :=
      LOCAL_BUILT_MODULE :=
      LOCAL_INSTALLED_MODULE :=
      LOCAL_MODULE_STEM :=
      LOCAL_BUILT_MODULE_STEM :=
      LOCAL_INSTALLED_MODULE_STEM :=
      LOCAL_INTERMEDIATE_TARGETS :=

      include $(BUILD_SYSTEM)/static_library_internal.mk
    endif
  LOCAL_2ND_ARCH_VAR_PREFIX :=
  endif # TARGET_2ND_ARCH
endif

my_module_arch_supported :=

###########################################################
## Copy headers to the install tree
###########################################################
include $(BUILD_COPY_HEADERS)
