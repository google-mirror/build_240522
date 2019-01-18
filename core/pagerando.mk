
# Clear the internal variables, to make multilib builds work
ifndef LOCAL_IS_HOST_MODULE
  LOCAL_PAGERANDO := $(strip $(LOCAL_PAGERANDO))
  ifdef LOCAL_PAGERANDO_$(my_32_64_bit_suffix)
  LOCAL_PAGERANDO := $(strip $(LOCAL_PAGERANDO_$(my_32_64_bit_suffix)))
  endif
  ifdef LOCAL_PAGERANDO_$(TARGET_$(LOCAL_2ND_ARCH_VAR_PREFIX)ARCH)
  LOCAL_PAGERANDO := $(strip $(LOCAL_PAGERANDO_$(TARGET_$(LOCAL_2ND_ARCH_VAR_PREFIX)ARCH)))
  endif
  ifeq ($(PAGERANDO)|$(LOCAL_PAGERANDO),true|)
    LOCAL_PAGERANDO := true
  endif
endif

ifneq ($($(LOCAL_2ND_ARCH_VAR_PREFIX)TARGET_SUPPORTS_PAGERANDO),true)
  LOCAL_PAGERANDO := false
endif
