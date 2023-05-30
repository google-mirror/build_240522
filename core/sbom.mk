# For SBOM generation
# This is included by base_rules.mk and is not necessary to be included in other .mk files
# unless a .mk file changes its installed file after including base_rules.mk.

ifdef my_register_name
  ifneq (, $(strip $(ALL_MODULES.$(my_register_name).INSTALLED)))
    $(foreach installed_file,$(ALL_MODULES.$(my_register_name).INSTALLED),\
      $(eval ALL_INSTALLED_FILES.$(installed_file) := $(my_register_name))\
      $(eval ALL_INSTALLED_FILES.$(installed_file).STATIC_LIBRARIES := $(strip $(sort $(LOCAL_STATIC_LIBRARIES))))\
      $(eval ALL_INSTALLED_FILES.$(installed_file).WHOLE_STATIC_LIBRARIES := $(strip $(sort $(LOCAL_WHOLE_STATIC_LIBRARIES))))\
    )
  endif
endif