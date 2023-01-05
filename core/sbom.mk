# For SBOM generation

ifdef my_register_name
  ifneq (, $(strip $(ALL_MODULES.$(my_register_name).INSTALLED)))
    $(foreach installed_file,$(ALL_MODULES.$(my_register_name).INSTALLED),\
      $(eval ALL_INSTALLED_FILES.$(installed_file) := $(my_register_name))\
    )
  endif
endif