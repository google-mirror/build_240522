# Inputs:
#   LOCAL_MODULE_CLASS, LOCAL_MODULE, LOCAL_MODULE_MAKEFILE, LOCAL_BUILT_MODULE
#   from base_rules.mk: my_kind, my_host_cross
#   my_common: empty or COMMON, like the argument to intermediates-dir-for
#   my_2nd_arch_prefix: usually LOCAL_2ND_ARCH_VAR_PREFIX, separate for JNI installation
#
#   my_link_type: the tags to apply to this module
#   my_warn_types: the tags to warn about in our dependencies
#   my_allowed_types: the tags to allow in our dependencies
#   my_link_deps: the dependencies, in the form of <MODULE_CLASS>:<name>
#

my_link_prefix := LINK_TYPE:$(call find-idf-prefix,$(my_kind),$(my_host_cross)):$(if $(my_common),$(my_common):_,_:$(if $(my_2nd_arch_prefix),$(my_2nd_arch_prefix),_))
link_type := $(my_link_prefix):$(LOCAL_MODULE_CLASS):$(LOCAL_MODULE)
ALL_LINK_TYPES := $(ALL_LINK_TYPES) $(link_type)
$(link_type).TYPE := $(my_link_type)
$(link_type).WARN := $(my_warn_types)
$(link_type).ALLOWED := $(my_allowed_types)
$(link_type).DEPS := $(addprefix $(my_link_prefix):,$(my_link_deps))
$(link_type).BUILT := $(LOCAL_BUILT_MODULE)

ifeq ($(LOCAL_MODULE_MAKEFILE),$(SOONG_ANDROID_MK))
  $(link_type).MAKEFILE := $(LOCAL_PATH)

  ifdef LOCAL_SOONG_VISIBILITY
    $(link_type).VISIBILITY := $(LOCAL_SOONG_VISIBILITY)
  endif
else
  $(link_type).MAKEFILE := $(LOCAL_MODULE_MAKEFILE)

  my_local_path := $(call clean-path,$(LOCAL_PATH))/
  $(foreach dep,$($(link_type).DEPS),\
    $(if $($(dep).VISIBILITY),\
      $(if $(filter $($(dep).VISIBILITY),$(my_local_path)),,\
        $(eval _other := //$$($$(dep).MAKEFILE):$$(word 6,$$(subst :,$$(space),$$(dep))))\
        $(call pretty-warning,depends on $(_other) which is not visible to this module))))
  my_local_path :=
endif

link_type :=
my_allowed_types :=
my_link_prefix :=
my_link_type :=
my_warn_types :=
