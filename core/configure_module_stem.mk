my_multilib_stem := $(LOCAL_MODULE_STEM_$(if $($(LOCAL_2ND_ARCH_VAR_PREFIX)TARGET_IS_64_BIT),64,32))
ifdef my_multilib_stem
LOCAL_MODULE_STEM := $(my_multilib_stem)
endif

ifndef LOCAL_MODULE_STEM
  LOCAL_MODULE_STEM := $(LOCAL_MODULE)
endif

ifndef LOCAL_BUILT_MODULE_STEM
  LOCAL_BUILT_MODULE_STEM := $(LOCAL_MODULE_STEM)$(LOCAL_MODULE_SUFFIX)
endif

ifndef LOCAL_INSTALLED_MODULE_STEM
  LOCAL_INSTALLED_MODULE_STEM := $(LOCAL_MODULE_STEM)$(LOCAL_MODULE_SUFFIX)
endif
