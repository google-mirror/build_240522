###########################################################
## Track NOTICE files
###########################################################
$(call record-module-type,NOTICE_FILE)

ifneq ($(LOCAL_NOTICE_FILE),)
notice_file:=$(strip $(LOCAL_NOTICE_FILE))
else
notice_file:=$(strip $(wildcard $(LOCAL_PATH)/NOTICE))
endif

ifeq ($(LOCAL_MODULE_CLASS),GYP)
  # We ignore NOTICE files for modules of type GYP.
  notice_file :=
endif

# Soong generates stub libraries that don't need NOTICE files
ifdef LOCAL_NO_NOTICE_FILE
  ifneq ($(LOCAL_MODULE_MAKEFILE),$(SOONG_ANDROID_MK))
    $(call pretty-error,LOCAL_NO_NOTICE_FILE should not be used by Android.mk files)
  endif
  notice_file :=
endif

ifeq ($(LOCAL_MODULE_CLASS),NOTICE_FILES)
# If this is a NOTICE-only module, we don't include base_rule.mk,
# so my_prefix is not set at this point.
ifeq ($(LOCAL_IS_HOST_MODULE),true)
  my_prefix := HOST_
  LOCAL_HOST_PREFIX :=
else
  my_prefix := TARGET_
endif
endif

<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
=======
installed_notice_file :=

is_container:=$(strip $(LOCAL_MODULE_IS_CONTAINER))
ifeq (,$(is_container))
ifneq (,$(strip $(filter %.zip %.tar %.tgz %.tar.gz %.apk %.img %.srcszip %.apex, $(LOCAL_BUILT_MODULE))))
is_container:=true
else
is_container:=false
endif
else ifneq (,$(strip $(filter-out true false,$(is_container))))
$(error Unrecognized value '$(is_container)' for LOCAL_MODULE_IS_CONTAINER)
endif

ifeq (true,$(is_container))
# Include shared libraries' notices for "container" types, but not for binaries etc.
notice_deps := \
    $(strip \
        $(foreach d, \
            $(LOCAL_REQUIRED_MODULES) \
            $(LOCAL_STATIC_LIBRARIES) \
            $(LOCAL_WHOLE_STATIC_LIBRARIES) \
            $(LOCAL_SHARED_LIBRARIES) \
            $(LOCAL_DYLIB_LIBRARIES) \
            $(LOCAL_RLIB_LIBRARIES) \
            $(LOCAL_PROC_MACRO_LIBRARIES) \
            $(LOCAL_HEADER_LIBRARIES) \
            $(LOCAL_STATIC_JAVA_LIBRARIES) \
            $(LOCAL_JAVA_LIBRARIES) \
            $(LOCAL_JNI_SHARED_LIBRARIES) \
            ,$(subst :,_,$(d)):static \
        ) \
    )
else
notice_deps := \
    $(strip \
        $(foreach d, \
            $(LOCAL_REQUIRED_MODULES) \
            $(LOCAL_STATIC_LIBRARIES) \
            $(LOCAL_WHOLE_STATIC_LIBRARIES) \
            $(LOCAL_RLIB_LIBRARIES) \
            $(LOCAL_PROC_MACRO_LIBRARIES) \
            $(LOCAL_HEADER_LIBRARIES) \
            $(LOCAL_STATIC_JAVA_LIBRARIES) \
            ,$(subst :,_,$(d)):static \
        )$(foreach d, \
            $(LOCAL_SHARED_LIBRARIES) \
            $(LOCAL_DYLIB_LIBRARIES) \
            $(LOCAL_JAVA_LIBRARIES) \
            $(LOCAL_JNI_SHARED_LIBRARIES) \
            ,$(subst :,_,$(d)):dynamic \
        ) \
    )
endif
ifeq ($(LOCAL_IS_HOST_MODULE),true)
notice_deps := $(strip $(notice_deps) $(foreach d,$(LOCAL_HOST_REQUIRED_MODULES),$(subst :,_,$(d)):static))
else
notice_deps := $(strip $(notice_deps) $(foreach d,$(LOCAL_TARGET_REQUIRED_MODULES),$(subst :,_,$(d)):static))
endif

local_path := $(LOCAL_PATH)

ifdef my_register_name
ALL_MODULES.$(my_register_name).LICENSE_PACKAGE_NAME := $(strip $(license_package_name))
ALL_MODULES.$(my_register_name).MODULE_TYPE := $(strip $(ALL_MODULES.$(my_register_name).MODULE_TYPE) $(LOCAL_MODULE_TYPE))
ALL_MODULES.$(my_register_name).MODULE_CLASS := $(strip $(ALL_MODULES.$(my_register_name).MODULE_CLASS) $(LOCAL_MODULE_CLASS))
ALL_MODULES.$(my_register_name).LICENSE_KINDS := $(ALL_MODULES.$(my_register_name).LICENSE_KINDS) $(license_kinds)
ALL_MODULES.$(my_register_name).LICENSE_CONDITIONS := $(ALL_MODULES.$(my_register_name).LICENSE_CONDITIONS) $(license_conditions)
ALL_MODULES.$(my_register_name).LICENSE_INSTALL_MAP := $(ALL_MODULES.$(my_register_name).LICENSE_INSTALL_MAP) $(install_map)
ALL_MODULES.$(my_register_name).NOTICE_DEPS := $(ALL_MODULES.$(my_register_name).NOTICE_DEPS) $(notice_deps)
ALL_MODULES.$(my_register_name).IS_CONTAINER := $(strip $(filter-out false,$(ALL_MODULES.$(my_register_name).IS_CONTAINER) $(is_container)))
ALL_MODULES.$(my_register_name).PATH := $(strip $(ALL_MODULES.$(my_register_name).PATH) $(local_path))
endif

>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
ifdef notice_file

<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
=======
ifdef my_register_name
ALL_MODULES.$(my_register_name).NOTICES := $(ALL_MODULES.$(my_register_name).NOTICES) $(notice_file)
endif

>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
# This relies on the name of the directory in PRODUCT_OUT matching where
# it's installed on the target - i.e. system, data, etc.  This does
# not work for root and isn't exact, but it's probably good enough for
# compliance.
# Includes the leading slash
ifdef LOCAL_INSTALLED_MODULE
  module_installed_filename := $(patsubst $(PRODUCT_OUT)/%,%,$(LOCAL_INSTALLED_MODULE))
else
  # This module isn't installable
  ifneq ($(filter STATIC_LIBRARIES HEADER_LIBRARIES,$(LOCAL_MODULE_CLASS)),)
    # Stick the static libraries with the dynamic libraries.
    # We can't use xxx_OUT_STATIC_LIBRARIES because it points into
    # device-obj or host-obj.
    module_installed_filename := \
        $(patsubst $(PRODUCT_OUT)/%,%,$($(LOCAL_2ND_ARCH_VAR_PREFIX)$(my_prefix)OUT_SHARED_LIBRARIES))/$(notdir $(LOCAL_BUILT_MODULE))
  else
    ifeq ($(LOCAL_MODULE_CLASS),JAVA_LIBRARIES)
      # Stick the static java libraries with the regular java libraries.
      module_leaf := $(notdir $(LOCAL_BUILT_MODULE))
      # javalib.jar is the default name for the build module (and isn't meaningful)
      # If that's what we have, substitute the module name instead.  These files
      # aren't included on the device, so this name is synthetic anyway.
      ifneq ($(filter javalib.jar,$(module_leaf)),)
        module_leaf := $(LOCAL_MODULE).jar
      endif
      module_installed_filename := \
          $(patsubst $(PRODUCT_OUT)/%,%,$($(my_prefix)OUT_JAVA_LIBRARIES))/$(module_leaf)
    else
      $(error Cannot determine where to install NOTICE file for $(LOCAL_MODULE))
    endif # JAVA_LIBRARIES
  endif # STATIC_LIBRARIES
endif

# In case it's actually a host file
module_installed_filename := $(patsubst $(HOST_OUT)/%,%,$(module_installed_filename))
module_installed_filename := $(patsubst $(HOST_CROSS_OUT)/%,%,$(module_installed_filename))

installed_notice_file := $($(my_prefix)OUT_NOTICE_FILES)/src/$(module_installed_filename).txt

<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
=======
ifdef my_register_name
ALL_MODULES.$(my_register_name).INSTALLED_NOTICE_FILE := $(ALL_MODULES.$(my_register_name).INSTALLED_NOTICE_FILE) $(installed_notice_file)
ALL_MODULES.$(my_register_name).MODULE_INSTALLED_FILENAMES := $(ALL_MODULES.$(my_register_name).MODULE_INSTALLED_FILENAMES) $(module_installed_filename)
INSTALLED_NOTICE_FILES.$(installed_notice_file).MODULE := $(my_register_name)
else
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
$(installed_notice_file): PRIVATE_INSTALLED_MODULE := $(module_installed_filename)
<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
=======
$(installed_notice_file) : PRIVATE_NOTICES := $(sort $(foreach n,$(notice_file),$(if $(filter %:%,$(n)), $(call word-colon,1,$(n)), $(n))))
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)

$(installed_notice_file): $(foreach n,$(notice_file),$(if $(filter %:%,$(n)), $(call word-colon,1,$(n)), $(n)))
	@echo Notice file: $< -- $@
	$(hide) mkdir -p $(dir $@)
<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
	$(hide) cat $< > $@
=======
	$(hide) awk 'FNR==1 && NR > 1 {print "\n"} {print}' $(PRIVATE_NOTICES) > $@
endif
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)

ifdef LOCAL_INSTALLED_MODULE
# Make LOCAL_INSTALLED_MODULE depend on NOTICE files if they exist
# libraries so they get installed along with it.  Make it an order-only
# dependency so we don't re-install a module when the NOTICE changes.
$(LOCAL_INSTALLED_MODULE): | $(installed_notice_file)
endif

# To facilitate collecting NOTICE files for apps_only build,
# we install the NOTICE file even if a module gets built but not installed,
# because shared jni libraries won't be installed to the system image.
ifdef TARGET_BUILD_APPS
# for static Java libraries, we don't need to even build LOCAL_BUILT_MODULE,
# but just javalib.jar in the common intermediate dir.
ifeq ($(LOCAL_MODULE_CLASS),JAVA_LIBRARIES)
$(intermediates.COMMON)/javalib.jar : | $(installed_notice_file)
else
$(LOCAL_BUILT_MODULE): | $(installed_notice_file)
endif  # JAVA_LIBRARIES
endif  # TARGET_BUILD_APPS

else
# NOTICE file does not exist
installed_notice_file :=
endif

# Create a predictable, phony target to build this notice file.
# Define it even if the notice file doesn't exist so that other
# modules can depend on it.
notice_target := NOTICE-$(if \
    $(LOCAL_IS_HOST_MODULE),HOST,TARGET)-$(LOCAL_MODULE_CLASS)-$(LOCAL_MODULE)
.PHONY: $(notice_target)
$(notice_target): $(installed_notice_file)
