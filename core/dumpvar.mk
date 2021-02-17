# ---------------------------------------------------------------
# the setpath shell function in envsetup.sh uses this to figure out
# what to add to the path given the config we have chosen.
ifeq ($(CALLED_FROM_SETUP),true)

ifneq ($(filter /%,$(SOONG_HOST_OUT_EXECUTABLES)),)
ABP := $(SOONG_HOST_OUT_EXECUTABLES)
else
ABP := $(PWD)/$(SOONG_HOST_OUT_EXECUTABLES)
endif
ifneq ($(filter /%,$(HOST_OUT_EXECUTABLES)),)
ABP := $(ABP):$(HOST_OUT_EXECUTABLES)
else
ABP := $(ABP):$(PWD)/$(HOST_OUT_EXECUTABLES)
endif

ANDROID_BUILD_PATHS := $(ABP)
ANDROID_PREBUILTS := prebuilt/$(HOST_PREBUILT_TAG)
ANDROID_GCC_PREBUILTS := prebuilts/gcc/$(HOST_PREBUILT_TAG)

# Dump mulitple variables to "<var>=<value>" pairs, one per line.
# The output may be executed as bash script.
# Input variables:
#   DUMP_MANY_VARS: the list of variable names.
#   DUMP_VAR_PREFIX: an optional prefix of the variable name added to the output.
# The value is printed in parts because large variables like PRODUCT_PACKAGES
# can exceed the maximum linux command line size
.PHONY: dump-many-vars
dump-many-vars :
	@$(foreach v, $(DUMP_MANY_VARS),\
<<<<<<< HEAD   (4be654 Merge "Merge empty history for sparse-7121469-L4290000080720)
	  echo "$(DUMP_VAR_PREFIX)$(v)='$($(v))'";)
=======
	  printf "%s='%s" '$(DUMP_VAR_PREFIX)$(v)' '$(firstword $($(v)))'; \
	  $(foreach part, $(wordlist 2, $(words $($(v))), $($(v))),\
	    printf " %s" '$(part)'$(newline))\
	  printf "'\n";)
>>>>>>> BRANCH (fe6ad7 Merge "Version bump to RBT1.210107.001.A1 [core/build_id.mk])

endif # CALLED_FROM_SETUP
