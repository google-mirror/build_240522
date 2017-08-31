$(call record-module-type,STATIC_LIBRARY)
my_prefix := TARGET_
include $(BUILD_SYSTEM)/multilib.mk

ifndef my_module_multilib
# libraries default to building for both architecturess
my_module_multilib := both
endif

LOCAL_2ND_ARCH_VAR_PREFIX :=
include $(BUILD_SYSTEM)/module_arch_supported.mk

ifeq ($(my_module_arch_supported),true)

# We need to build a pagerando version of the library in case any pagerando
# binaries need this static library
saved_local_pagerando := $(LOCAL_PAGERANDO)
include $(BUILD_SYSTEM)/pagerando.mk
ifeq ($(LOCAL_PAGERANDO),true)
  LOCAL_PAGERANDO_MODULE_SUFFIX := _pagerando
  LOCAL_PAGERANDO_INTERMEDIATES_SUFFIX := _pagerando
  LOCAL_PAGERANDO_STATIC_SUFFIX := _pagerando

  include $(BUILD_SYSTEM)/static_library_internal.mk

  # Clear build variables to set up for non-pagerando build
  OVERRIDE_BUILT_MODULE_PATH :=
  LOCAL_BUILT_MODULE :=
  LOCAL_INSTALLED_MODULE :=
  LOCAL_MODULE_STEM :=
  LOCAL_BUILT_MODULE_STEM :=
  LOCAL_INSTALLED_MODULE_STEM :=
  LOCAL_INTERMEDIATE_TARGETS :=
  LOCAL_PAGERANDO_MODULE_SUFFIX :=
  LOCAL_PAGERANDO_INTERMEDIATES_SUFFIX :=
  LOCAL_PAGERANDO_STATIC_SUFFIX :=
  LOCAL_PAGERANDO :=
endif

include $(BUILD_SYSTEM)/static_library_internal.mk

endif # my_module_arch_supported

ifdef TARGET_2ND_ARCH

LOCAL_2ND_ARCH_VAR_PREFIX := $(TARGET_2ND_ARCH_VAR_PREFIX)
include $(BUILD_SYSTEM)/module_arch_supported.mk

ifeq ($(my_module_arch_supported),true)
# Build for TARGET_2ND_ARCH
OVERRIDE_BUILT_MODULE_PATH :=
LOCAL_BUILT_MODULE :=
LOCAL_INSTALLED_MODULE :=
LOCAL_INTERMEDIATE_TARGETS :=

# We need to build a pagerando version of the library in case any pagerando
# binaries need this static library
LOCAL_PAGERANDO := $(saved_local_pagerando)
include $(BUILD_SYSTEM)/pagerando.mk
ifeq ($(LOCAL_PAGERANDO),true)
  LOCAL_PAGERANDO_MODULE_SUFFIX := _pagerando
  LOCAL_PAGERANDO_INTERMEDIATES_SUFFIX := _pagerando
  LOCAL_PAGERANDO_STATIC_SUFFIX := _pagerando

  include $(BUILD_SYSTEM)/static_library_internal.mk

  # Clear build variables to set up for non-pagerando build
  OVERRIDE_BUILT_MODULE_PATH :=
  LOCAL_BUILT_MODULE :=
  LOCAL_INSTALLED_MODULE :=
  LOCAL_MODULE_STEM :=
  LOCAL_BUILT_MODULE_STEM :=
  LOCAL_INSTALLED_MODULE_STEM :=
  LOCAL_INTERMEDIATE_TARGETS :=
  LOCAL_PAGERANDO_MODULE_SUFFIX :=
  LOCAL_PAGERANDO_INTERMEDIATES_SUFFIX :=
  LOCAL_PAGERANDO_STATIC_SUFFIX :=
  LOCAL_PAGERANDO :=
endif

include $(BUILD_SYSTEM)/static_library_internal.mk

endif # my_module_arch_supported

LOCAL_2ND_ARCH_VAR_PREFIX :=

endif # TARGET_2ND_ARCH

my_module_arch_supported :=

###########################################################
## Copy headers to the install tree
###########################################################
include $(BUILD_COPY_HEADERS)
