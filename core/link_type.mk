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

my_link_prefix := LINK_TYPE:$(call find-idf-prefix,$(my_kind),$(my_host_cross)):$(if $(my_common),$(my_common),_):$(if $(my_2nd_arch_prefix),$(my_2nd_arch_prefix),_)
link_type := $(my_link_prefix):$(LOCAL_MODULE_CLASS):$(LOCAL_MODULE)
ALL_LINK_TYPES := $(ALL_LINK_TYPES) $(link_type)
$(link_type).TYPE := $(my_link_type)
$(link_type).MAKEFILE := $(LOCAL_MODULE_MAKEFILE)
$(link_type).WARN := $(my_warn_types)
$(link_type).ALLOWED := $(my_allowed_types)
$(link_type).DEPS := $(addprefix $(my_link_prefix):,$(my_link_deps))

# Write out the file-based link_type rules for the ALLOW_MISSING_DEPENDENCIES
# case. We always need to write the file for mm to work.
my_link_type_deps := $(strip $(foreach l,$(my_link_deps),\
  $(eval _class := $(word 1,$(subst :,$(space),$(l))))\
  $(eval _name := $(word 2,$(subst :,$(space),$(l))))\
  $(call intermediates-dir-for,$(_class),$(_name),$(my_kind),$(my_common),$(my_2nd_arch_prefix),$(my_host_cross))/$(my_2nd_arch_prefix)link_type))

# We need $(my_2nd_arch_prefix) in the filename too, for multilib apps who
# don't have separate intermediate directories, but need to link against both
# architectures of shared libraries.
#
# We could combine them into a single file with dependencies on both arches,
# since the types should be the same, but that adds complexity for a very small
# speedup.
my_link_type_file := $(call intermediates-dir-for,$(LOCAL_MODULE_CLASS),$(LOCAL_MODULE),$(my_kind),$(my_common),$(my_2nd_arch_prefix),$(my_host_cross))/$(my_2nd_arch_prefix)link_type
$(my_link_type_file): PRIVATE_LINK_TYPE := $(my_link_type)
$(my_link_type_file): PRIVATE_WARN_TYPES := $(my_warn_types)
$(my_link_type_file): PRIVATE_ALLOWED_TYPES := $(my_allowed_types)
$(my_link_type_file): PRIVATE_DEPS := $(my_link_type_deps)
$(my_link_type_file): PRIVATE_MODULE := $(LOCAL_MODULE)
$(my_link_type_file): PRIVATE_MAKEFILE := $(LOCAL_MODULE_MAKEFILE)
$(my_link_type_file): $(my_link_type_deps) $(CHECK_LINK_TYPE)
	@echo Check module type: $@
	$(hide) mkdir -p $(dir $@) && rm -f $@
ifeq ($(ALLOW_MISSING_DEPENDENCIES),true)
	$(hide) $(CHECK_LINK_TYPE) --makefile $(PRIVATE_MAKEFILE) --module $(PRIVATE_MODULE) \
	  --type "$(PRIVATE_LINK_TYPE)" $(addprefix --allowed ,$(PRIVATE_ALLOWED_TYPES)) \
	  $(addprefix --warn ,$(PRIVATE_WARN_TYPES)) $(PRIVATE_DEPS)
endif
	$(hide) echo "$(PRIVATE_LINK_TYPE)" >$@

$(LOCAL_BUILT_MODULE): | $(my_link_type_file)

link_type :=
my_allowed_types :=
my_link_prefix :=
my_link_type :=
my_link_type_deps :=
my_link_type_file :=
my_warn_types :=
