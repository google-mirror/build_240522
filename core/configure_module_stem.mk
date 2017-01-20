my_multilib_stem := $(LOCAL_MODULE_STEM_$(if $($(LOCAL_2ND_ARCH_VAR_PREFIX)$(my_prefix)IS_64_BIT),64,32))
ifdef my_multilib_stem
  my_module_stem := $(my_multilib_stem)
else ifdef LOCAL_MODULE_STEM
  my_module_stem := $(LOCAL_MODULE_STEM)
else ifneq ($(LOCAL_EXTENDS_MODULE),)
  # If this module extends other module, stem of this module becomes that of other module.
  # TODO(jiyong): to be correct, stem of this module should be step of the original module.
  # However, it is impossible(correct me if I am wrong) to know the stem of other module while
  # executing Android.mk.
  # TODO(jiyong): should ensure that this extended module and the original module are installed
  # in different directories. However, this also seems to be impossible. For now, we just rely
  # on make generating duplicate target error.
  my_module_stem := $(LOCAL_EXTENDS_MODULE)
else
  my_module_stem := $(LOCAL_MODULE)
endif

ifdef LOCAL_BUILT_MODULE_STEM
  my_built_module_stem := $(LOCAL_BUILT_MODULE_STEM)
else
  my_built_module_stem := $(my_module_stem)$(LOCAL_MODULE_SUFFIX)
endif

ifdef LOCAL_INSTALLED_MODULE_STEM
  my_installed_module_stem := $(LOCAL_INSTALLED_MODULE_STEM)
else
  my_installed_module_stem := $(my_module_stem)$(LOCAL_MODULE_SUFFIX)
endif
