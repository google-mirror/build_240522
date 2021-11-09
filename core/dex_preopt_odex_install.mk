# dexpreopt_odex_install.mk is used to define odex creation rules for JARs and APKs
# This file depends on variables set in base_rules.mk
<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
# Output variables: LOCAL_DEX_PREOPT, LOCAL_UNCOMPRESS_DEX, built_odex,
#                   dexpreopt_boot_jar_module
=======
# Input variables: my_manifest_or_apk
# Output variables: LOCAL_DEX_PREOPT, LOCAL_UNCOMPRESS_DEX

ifeq (true,$(LOCAL_USE_EMBEDDED_DEX))
  LOCAL_UNCOMPRESS_DEX := true
else
  LOCAL_UNCOMPRESS_DEX :=
endif
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)

# We explicitly uncompress APKs of privileged apps, and used by
# privileged apps
LOCAL_UNCOMPRESS_DEX := false
ifneq (true,$(DONT_UNCOMPRESS_PRIV_APPS_DEXS))
ifeq (true,$(LOCAL_PRIVILEGED_MODULE))
  LOCAL_UNCOMPRESS_DEX := true
else
  ifneq (,$(filter $(PRODUCT_LOADED_BY_PRIVILEGED_MODULES), $(LOCAL_MODULE)))
    LOCAL_UNCOMPRESS_DEX := true
  endif  # PRODUCT_LOADED_BY_PRIVILEGED_MODULES
endif  # LOCAL_PRIVILEGED_MODULE
endif  # DONT_UNCOMPRESS_PRIV_APPS_DEXS

# Setting LOCAL_DEX_PREOPT based on WITH_DEXPREOPT, LOCAL_DEX_PREOPT, etc
LOCAL_DEX_PREOPT := $(strip $(LOCAL_DEX_PREOPT))
ifneq (true,$(WITH_DEXPREOPT))
  LOCAL_DEX_PREOPT :=
else # WITH_DEXPREOPT=true
  ifeq (,$(TARGET_BUILD_APPS)) # TARGET_BUILD_APPS empty
    ifndef LOCAL_DEX_PREOPT # LOCAL_DEX_PREOPT undefined
      ifneq ($(filter $(TARGET_OUT)/%,$(my_module_path)),) # Installed to system.img.
        ifeq (,$(LOCAL_APK_LIBRARIES)) # LOCAL_APK_LIBRARIES empty
          # If we have product-specific config for this module?
          ifeq (disable,$(DEXPREOPT.$(TARGET_PRODUCT).$(LOCAL_MODULE).CONFIG))
            LOCAL_DEX_PREOPT := false
          else
            LOCAL_DEX_PREOPT := $(DEX_PREOPT_DEFAULT)
          endif
        else # LOCAL_APK_LIBRARIES not empty
          LOCAL_DEX_PREOPT := nostripping
        endif # LOCAL_APK_LIBRARIES not empty
      endif # Installed to system.img.
    endif # LOCAL_DEX_PREOPT undefined
  endif # TARGET_BUILD_APPS empty
endif # WITH_DEXPREOPT=true
ifeq (false,$(LOCAL_DEX_PREOPT))
  LOCAL_DEX_PREOPT :=
endif
<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
=======

# Disable preopt for tests.
ifneq (,$(filter $(LOCAL_MODULE_TAGS),tests))
  LOCAL_DEX_PREOPT :=
endif

# If we have product-specific config for this module?
ifneq (,$(filter $(LOCAL_MODULE),$(DEXPREOPT_DISABLED_MODULES)))
  LOCAL_DEX_PREOPT :=
endif

# Disable preopt for DISABLE_PREOPT
ifeq (true,$(DISABLE_PREOPT))
  LOCAL_DEX_PREOPT :=
endif

# Disable preopt if not WITH_DEXPREOPT
ifneq (true,$(WITH_DEXPREOPT))
  LOCAL_DEX_PREOPT :=
endif

>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
ifdef LOCAL_UNINSTALLABLE_MODULE
LOCAL_DEX_PREOPT :=
endif
<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
ifeq (,$(strip $(built_dex)$(my_prebuilt_src_file)$(LOCAL_SOONG_DEX_JAR))) # contains no java code
LOCAL_DEX_PREOPT :=
endif
=======

# Disable preopt if the app contains no java code.
ifeq (,$(strip $(built_dex)$(my_prebuilt_src_file)$(LOCAL_SOONG_DEX_JAR)))
  LOCAL_DEX_PREOPT :=
endif

>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
# if WITH_DEXPREOPT_BOOT_IMG_AND_SYSTEM_SERVER_ONLY=true and module is not in boot class path skip
# Also preopt system server jars since selinux prevents system server from loading anything from
# /data. If we don't do this they will need to be extracted which is not favorable for RAM usage
# or performance. If my_preopt_for_extracted_apk is true, we ignore the only preopt boot image
# options.
system_server_jars := $(foreach m,$(PRODUCT_SYSTEM_SERVER_JARS),$(call word-colon,2,$(m)))
ifneq (true,$(my_preopt_for_extracted_apk))
<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
ifeq (true,$(WITH_DEXPREOPT_BOOT_IMG_AND_SYSTEM_SERVER_ONLY))
ifeq ($(filter $(PRODUCT_SYSTEM_SERVER_JARS) $(DEXPREOPT_BOOT_JARS_MODULES),$(LOCAL_MODULE)),)
LOCAL_DEX_PREOPT :=
endif
endif
=======
  ifeq (true,$(WITH_DEXPREOPT_BOOT_IMG_AND_SYSTEM_SERVER_ONLY))
    ifeq ($(filter $(system_server_jars) $(DEXPREOPT_BOOT_JARS_MODULES),$(LOCAL_MODULE)),)
      LOCAL_DEX_PREOPT :=
    endif
  endif
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
endif

ifeq ($(LOCAL_DEX_PREOPT),true)

# Don't strip with dexes we explicitly uncompress (dexopt will not store the dex code).
ifeq ($(LOCAL_UNCOMPRESS_DEX),true)
LOCAL_DEX_PREOPT := nostripping
endif  # LOCAL_UNCOMPRESS_DEX

# system_other isn't there for an OTA, so don't strip
# if module is on system, and odex is on system_other.
ifeq ($(BOARD_USES_SYSTEM_OTHER_ODEX),true)
ifneq ($(call install-on-system-other, $(my_module_path)),)
LOCAL_DEX_PREOPT := nostripping
endif  # install-on-system-other
endif  # BOARD_USES_SYSTEM_OTHER_ODEX

# We also don't strip if all dexs are uncompressed (dexopt will not store the dex code),
# but that requires to inspect the source file, which is too early at this point (as we
# don't know if the source file will actually be used).
# See dexpreopt-remove-classes.dex.

endif  # LOCAL_DEX_PREOPT

built_odex :=
built_vdex :=
built_art :=
installed_odex :=
installed_vdex :=
installed_art :=
built_installed_odex :=
built_installed_vdex :=
built_installed_art :=
my_process_profile :=
my_profile_is_text_listing :=

ifeq (false,$(WITH_DEX_PREOPT_GENERATE_PROFILE))
LOCAL_DEX_PREOPT_GENERATE_PROFILE := false
endif

ifndef LOCAL_DEX_PREOPT_GENERATE_PROFILE


# If LOCAL_DEX_PREOPT_GENERATE_PROFILE is not defined, default it based on the existence of the
# profile class listing. TODO: Use product specific directory here.
my_classes_directory := $(PRODUCT_DEX_PREOPT_PROFILE_DIR)
LOCAL_DEX_PREOPT_PROFILE := $(my_classes_directory)/$(LOCAL_MODULE).prof

ifneq (,$(wildcard $(LOCAL_DEX_PREOPT_PROFILE)))
my_process_profile := true
my_profile_is_text_listing := false
endif
else
my_process_profile := $(LOCAL_DEX_PREOPT_GENERATE_PROFILE)
my_profile_is_text_listing := true
LOCAL_DEX_PREOPT_PROFILE := $(LOCAL_DEX_PREOPT_PROFILE_CLASS_LISTING)
endif

ifeq (true,$(my_process_profile))

ifeq (,$(LOCAL_DEX_PREOPT_APP_IMAGE))
LOCAL_DEX_PREOPT_APP_IMAGE := true
endif

<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
ifndef LOCAL_DEX_PREOPT_PROFILE
$(call pretty-error,Must have specified class listing (LOCAL_DEX_PREOPT_PROFILE))
endif
ifeq (,$(dex_preopt_profile_src_file))
$(call pretty-error, Internal error: dex_preopt_profile_src_file must be set)
endif
my_built_profile := $(dir $(LOCAL_BUILT_MODULE))/profile.prof
my_dex_location := $(patsubst $(PRODUCT_OUT)%,%,$(LOCAL_INSTALLED_MODULE))
# Remove compressed APK extension.
my_dex_location := $(patsubst %.gz,%,$(my_dex_location))
$(my_built_profile): PRIVATE_BUILT_MODULE := $(dex_preopt_profile_src_file)
$(my_built_profile): PRIVATE_DEX_LOCATION := $(my_dex_location)
$(my_built_profile): PRIVATE_SOURCE_CLASSES := $(LOCAL_DEX_PREOPT_PROFILE)
$(my_built_profile): $(LOCAL_DEX_PREOPT_PROFILE)
$(my_built_profile): $(PROFMAN)
$(my_built_profile): $(dex_preopt_profile_src_file)
ifeq (true,$(my_profile_is_text_listing))
# The profile is a test listing of classes (used for framework jars).
# We need to generate the actual binary profile before being able to compile.
	$(hide) mkdir -p $(dir $@)
	ANDROID_LOG_TAGS="*:e" $(PROFMAN) \
		--create-profile-from=$(PRIVATE_SOURCE_CLASSES) \
		--apk=$(PRIVATE_BUILT_MODULE) \
		--dex-location=$(PRIVATE_DEX_LOCATION) \
		--reference-profile-file=$@
else
# The profile is binary profile (used for apps). Run it through profman to
# ensure the profile keys match the apk.
$(my_built_profile):
	$(hide) mkdir -p $(dir $@)
	touch $@
	ANDROID_LOG_TAGS="*:i" $(PROFMAN) \
	  --copy-and-update-profile-key \
		--profile-file=$(PRIVATE_SOURCE_CLASSES) \
		--apk=$(PRIVATE_BUILT_MODULE) \
		--dex-location=$(PRIVATE_DEX_LOCATION) \
		--reference-profile-file=$@ \
	|| echo "Profile out of date for $(PRIVATE_BUILT_MODULE)"
endif

my_profile_is_text_listing :=
dex_preopt_profile_src_file :=

# Remove compressed APK extension.
my_installed_profile := $(patsubst %.gz,%,$(LOCAL_INSTALLED_MODULE)).prof

# my_installed_profile := $(LOCAL_INSTALLED_MODULE).prof
$(eval $(call copy-one-file,$(my_built_profile),$(my_installed_profile)))
build_installed_profile:=$(my_built_profile):$(my_installed_profile)
else
build_installed_profile:=
my_installed_profile :=
=======
################################################################################
# Local module variables and functions used in dexpreopt and manifest_check.
################################################################################

my_filtered_optional_uses_libraries := $(filter-out $(INTERNAL_PLATFORM_MISSING_USES_LIBRARIES), \
  $(LOCAL_OPTIONAL_USES_LIBRARIES))

# TODO(b/132357300): This may filter out too much, as PRODUCT_PACKAGES doesn't
# include all packages (the full list is unknown until reading all Android.mk
# makefiles). As a consequence, a library may be present but not included in
# dexpreopt, which will result in class loader context mismatch and a failure
# to load dexpreopt code on device. We should fix this, either by deferring
# dependency computation until the full list of product packages is known, or
# by adding product-specific lists of missing libraries.
my_filtered_optional_uses_libraries := $(filter $(PRODUCT_PACKAGES), \
  $(my_filtered_optional_uses_libraries))

ifeq ($(LOCAL_MODULE_CLASS),APPS)
  # compatibility libraries are added to class loader context of an app only if
  # targetSdkVersion in the app's manifest is lower than the given SDK version

  my_dexpreopt_libs_compat_28 := \
    org.apache.http.legacy

  my_dexpreopt_libs_compat_29 := \
    android.hidl.manager-V1.0-java \
    android.hidl.base-V1.0-java

  my_dexpreopt_libs_compat_30 := \
    android.test.base \
    android.test.mock

  my_dexpreopt_libs_compat := \
    $(my_dexpreopt_libs_compat_28) \
    $(my_dexpreopt_libs_compat_29) \
    $(my_dexpreopt_libs_compat_30)
else
  my_dexpreopt_libs_compat :=
endif

my_dexpreopt_libs := \
  $(LOCAL_USES_LIBRARIES) \
  $(my_filtered_optional_uses_libraries)

# Module dexpreopt.config depends on dexpreopt.config files of each
# <uses-library> dependency, because these libraries may be processed after
# the current module by Make (there's no topological order), so the dependency
# information (paths, class loader context) may not be ready yet by the time
# this dexpreopt.config is generated. So it's necessary to add file-level
# dependencies between dexpreopt.config files.
my_dexpreopt_dep_configs := $(foreach lib, \
  $(filter-out $(my_dexpreopt_libs_compat),$(LOCAL_USES_LIBRARIES) $(my_filtered_optional_uses_libraries)), \
  $(call intermediates-dir-for,JAVA_LIBRARIES,$(lib),,)/dexpreopt.config)

# 1: SDK version
# 2: list of libraries
#
# Make does not process modules in topological order wrt. <uses-library>
# dependencies, therefore we cannot rely on variables to get the information
# about dependencies (in particular, their on-device path and class loader
# context). This information is communicated via dexpreopt.config files: each
# config depends on configs for <uses-library> dependencies of this module,
# and the dex_preopt_config_merger.py script reads all configs and inserts the
# missing bits from dependency configs into the module config.
#
# By default on-device path is /system/framework/*.jar, and class loader
# subcontext is empty. These values are correct for compatibility libraries,
# which are special and not handled by dex_preopt_config_merger.py.
#
add_json_class_loader_context = \
  $(call add_json_array, $(1)) \
  $(foreach lib, $(2),\
    $(call add_json_map_anon) \
    $(call add_json_str, Name, $(lib)) \
    $(call add_json_str, Host, $(call intermediates-dir-for,JAVA_LIBRARIES,$(lib),,COMMON)/javalib.jar) \
    $(call add_json_str, Device, /system/framework/$(lib).jar) \
    $(call add_json_val, Subcontexts, null) \
    $(call end_json_map)) \
  $(call end_json_array)

################################################################################
# Verify <uses-library> coherence between the build system and the manifest.
################################################################################

# Some libraries do not have a manifest, so there is nothing to check against.
# Handle it as if the manifest had zero <uses-library> tags: it is ok unless the
# module has non-empty LOCAL_USES_LIBRARIES or LOCAL_OPTIONAL_USES_LIBRARIES.
ifndef my_manifest_or_apk
  ifneq (,$(strip $(LOCAL_USES_LIBRARIES)$(LOCAL_OPTIONAL_USES_LIBRARIES)))
    $(error $(LOCAL_MODULE) has non-empty <uses-library> list but no manifest)
  else
    LOCAL_ENFORCE_USES_LIBRARIES := false
  endif
endif

# Disable the check for tests.
ifneq (,$(filter $(LOCAL_MODULE_TAGS),tests))
  LOCAL_ENFORCE_USES_LIBRARIES := false
endif
ifneq (,$(LOCAL_COMPATIBILITY_SUITE))
  LOCAL_ENFORCE_USES_LIBRARIES := false
endif

# Disable the check if the app contains no java code.
ifeq (,$(strip $(built_dex)$(my_prebuilt_src_file)$(LOCAL_SOONG_DEX_JAR)))
  LOCAL_ENFORCE_USES_LIBRARIES := false
endif

# Disable <uses-library> checks if dexpreopt is globally disabled.
# Without dexpreopt the check is not necessary, and although it is good to have,
# it is difficult to maintain on non-linux build platforms where dexpreopt is
# generally disabled (the check may fail due to various unrelated reasons, such
# as a failure to get manifest from an APK).
ifneq (true,$(WITH_DEXPREOPT))
  LOCAL_ENFORCE_USES_LIBRARIES := false
else ifeq (true,$(WITH_DEXPREOPT_BOOT_IMG_AND_SYSTEM_SERVER_ONLY))
  LOCAL_ENFORCE_USES_LIBRARIES := false
endif

# Verify LOCAL_USES_LIBRARIES/LOCAL_OPTIONAL_USES_LIBRARIES against the manifest.
ifndef LOCAL_ENFORCE_USES_LIBRARIES
  LOCAL_ENFORCE_USES_LIBRARIES := true
endif

my_enforced_uses_libraries :=
ifeq (true,$(LOCAL_ENFORCE_USES_LIBRARIES))
  my_verify_script := build/soong/scripts/manifest_check.py
  my_uses_libs_args := $(patsubst %,--uses-library %,$(LOCAL_USES_LIBRARIES))
  my_optional_uses_libs_args := $(patsubst %,--optional-uses-library %, \
    $(LOCAL_OPTIONAL_USES_LIBRARIES))
  my_relax_check_arg := $(if $(filter true,$(RELAX_USES_LIBRARY_CHECK)), \
    --enforce-uses-libraries-relax,)
  my_dexpreopt_config_args := $(patsubst %,--dexpreopt-config %,$(my_dexpreopt_dep_configs))

  my_enforced_uses_libraries := $(intermediates.COMMON)/enforce_uses_libraries.status
  $(my_enforced_uses_libraries): PRIVATE_USES_LIBRARIES := $(my_uses_libs_args)
  $(my_enforced_uses_libraries): PRIVATE_OPTIONAL_USES_LIBRARIES := $(my_optional_uses_libs_args)
  $(my_enforced_uses_libraries): PRIVATE_DEXPREOPT_CONFIGS := $(my_dexpreopt_config_args)
  $(my_enforced_uses_libraries): PRIVATE_RELAX_CHECK := $(my_relax_check_arg)
  $(my_enforced_uses_libraries): $(AAPT)
  $(my_enforced_uses_libraries): $(my_verify_script)
  $(my_enforced_uses_libraries): $(my_dexpreopt_dep_configs)
  $(my_enforced_uses_libraries): $(my_manifest_or_apk)
	@echo Verifying uses-libraries: $<
	rm -f $@
	$(my_verify_script) \
	  --enforce-uses-libraries \
	  --enforce-uses-libraries-status $@ \
	  --aapt $(AAPT) \
	  $(PRIVATE_USES_LIBRARIES) \
	  $(PRIVATE_OPTIONAL_USES_LIBRARIES) \
	  $(PRIVATE_DEXPREOPT_CONFIGS) \
	  $(PRIVATE_RELAX_CHECK) \
	  $<
  $(LOCAL_BUILT_MODULE) : $(my_enforced_uses_libraries)
endif

################################################################################
# Dexpreopt command.
################################################################################

my_dexpreopt_archs :=
my_dexpreopt_images :=
my_dexpreopt_images_deps :=
my_dexpreopt_image_locations_on_host :=
my_dexpreopt_image_locations_on_device :=
my_dexpreopt_infix := boot
my_create_dexpreopt_config :=
ifeq (true, $(DEXPREOPT_USE_ART_IMAGE))
  my_dexpreopt_infix := art
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
endif

ifdef LOCAL_DEX_PREOPT
<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
=======
  ifeq (,$(filter PRESIGNED,$(LOCAL_CERTIFICATE)))
    # Store uncompressed dex files preopted in /system
    ifeq ($(BOARD_USES_SYSTEM_OTHER_ODEX),true)
      ifeq ($(call install-on-system-other, $(my_module_path)),)
        LOCAL_UNCOMPRESS_DEX := true
      endif  # install-on-system-other
    else  # BOARD_USES_SYSTEM_OTHER_ODEX
      LOCAL_UNCOMPRESS_DEX := true
    endif
  endif
  my_create_dexpreopt_config := true
endif
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)

<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
dexpreopt_boot_jar_module := $(filter $(DEXPREOPT_BOOT_JARS_MODULES),$(LOCAL_MODULE))

# Filter org.apache.http.legacy.boot.
ifeq ($(dexpreopt_boot_jar_module),org.apache.http.legacy.boot)
dexpreopt_boot_jar_module :=
endif

ifdef dexpreopt_boot_jar_module
# For libart, the boot jars' odex files are replaced by $(DEFAULT_DEX_PREOPT_INSTALLED_IMAGE).
# We use this installed_odex trick to get boot.art installed.
installed_odex := $(DEFAULT_DEX_PREOPT_INSTALLED_IMAGE)
# Append the odex for the 2nd arch if we have one.
installed_odex += $($(TARGET_2ND_ARCH_VAR_PREFIX)DEFAULT_DEX_PREOPT_INSTALLED_IMAGE)
else  # boot jar
ifeq ($(LOCAL_MODULE_CLASS),JAVA_LIBRARIES)
# For a Java library, by default we build odex for both 1st arch and 2nd arch.
# But it can be overridden with "LOCAL_MULTILIB := first".
ifneq (,$(filter $(PRODUCT_SYSTEM_SERVER_JARS),$(LOCAL_MODULE)))
# For system server jars, we build for only "first".
my_module_multilib := first
else
my_module_multilib := $(LOCAL_MULTILIB)
endif
# #################################################
# Odex for the 1st arch
my_2nd_arch_prefix :=
include $(BUILD_SYSTEM)/setup_one_odex.mk
# #################################################
# Odex for the 2nd arch
ifdef TARGET_2ND_ARCH
ifneq ($(TARGET_TRANSLATE_2ND_ARCH),true)
ifneq (first,$(my_module_multilib))
my_2nd_arch_prefix := $(TARGET_2ND_ARCH_VAR_PREFIX)
include $(BUILD_SYSTEM)/setup_one_odex.mk
endif  # my_module_multilib is not first.
endif  # TARGET_TRANSLATE_2ND_ARCH not true
endif  # TARGET_2ND_ARCH
# #################################################
else  # must be APPS
# The preferred arch
my_2nd_arch_prefix := $(LOCAL_2ND_ARCH_VAR_PREFIX)
include $(BUILD_SYSTEM)/setup_one_odex.mk
ifdef TARGET_2ND_ARCH
ifeq ($(LOCAL_MULTILIB),both)
# The non-preferred arch
my_2nd_arch_prefix := $(if $(LOCAL_2ND_ARCH_VAR_PREFIX),,$(TARGET_2ND_ARCH_VAR_PREFIX))
include $(BUILD_SYSTEM)/setup_one_odex.mk
endif  # LOCAL_MULTILIB is both
endif  # TARGET_2ND_ARCH
endif  # LOCAL_MODULE_CLASS
endif  # boot jar

built_odex := $(strip $(built_odex))
built_vdex := $(strip $(built_vdex))
built_art := $(strip $(built_art))
installed_odex := $(strip $(installed_odex))
installed_vdex := $(strip $(installed_vdex))
installed_art := $(strip $(installed_art))

ifdef built_odex
ifeq (true,$(my_process_profile))
$(built_odex): $(my_built_profile)
$(built_odex): PRIVATE_PROFILE_PREOPT_FLAGS := --profile-file=$(my_built_profile)
else
$(built_odex): PRIVATE_PROFILE_PREOPT_FLAGS :=
endif

ifndef LOCAL_DEX_PREOPT_FLAGS
LOCAL_DEX_PREOPT_FLAGS := $(DEXPREOPT.$(TARGET_PRODUCT).$(LOCAL_MODULE).CONFIG)
ifndef LOCAL_DEX_PREOPT_FLAGS
LOCAL_DEX_PREOPT_FLAGS := $(PRODUCT_DEX_PREOPT_DEFAULT_FLAGS)
endif
endif

my_system_server_compiler_filter := $(PRODUCT_SYSTEM_SERVER_COMPILER_FILTER)
ifeq (,$(my_system_server_compiler_filter))
my_system_server_compiler_filter := speed
endif

my_default_compiler_filter := $(PRODUCT_DEX_PREOPT_DEFAULT_COMPILER_FILTER)
ifeq (,$(my_default_compiler_filter))
# If no default compiler filter is specified, default to 'quicken' to save on storage.
my_default_compiler_filter := quicken
endif

ifeq (,$(filter --compiler-filter=%, $(LOCAL_DEX_PREOPT_FLAGS)))
  ifneq (,$(filter $(PRODUCT_SYSTEM_SERVER_JARS),$(LOCAL_MODULE)))
    # Jars of system server, use the product option if it is set, speed otherwise.
    LOCAL_DEX_PREOPT_FLAGS += --compiler-filter=$(my_system_server_compiler_filter)
  else
    ifneq (,$(filter $(PRODUCT_DEXPREOPT_SPEED_APPS) $(PRODUCT_SYSTEM_SERVER_APPS),$(LOCAL_MODULE)))
      # Apps loaded into system server, and apps the product default to being compiled with the
      # 'speed' compiler filter.
      LOCAL_DEX_PREOPT_FLAGS += --compiler-filter=speed
    else
      ifeq (true,$(my_process_profile))
        # For non system server jars, use speed-profile when we have a profile.
        LOCAL_DEX_PREOPT_FLAGS += --compiler-filter=speed-profile
      else
        LOCAL_DEX_PREOPT_FLAGS += --compiler-filter=$(my_default_compiler_filter)
=======
# dexpreopt is disabled when TARGET_BUILD_UNBUNDLED_IMAGE is true,
# but dexpreopt config files are required to dexpreopt in post-processing.
ifeq ($(TARGET_BUILD_UNBUNDLED_IMAGE),true)
  my_create_dexpreopt_config := true
endif

ifeq ($(my_create_dexpreopt_config), true)
  ifeq ($(LOCAL_MODULE_CLASS),JAVA_LIBRARIES)
    my_module_multilib := $(LOCAL_MULTILIB)
    # If the module is not an SDK library and it's a system server jar, only preopt the primary arch.
    ifeq (,$(filter $(JAVA_SDK_LIBRARIES),$(LOCAL_MODULE)))
      # For a Java library, by default we build odex for both 1st arch and 2nd arch.
      # But it can be overridden with "LOCAL_MULTILIB := first".
      ifneq (,$(filter $(PRODUCT_SYSTEM_SERVER_JARS),$(LOCAL_MODULE)))
        # For system server jars, we build for only "first".
        my_module_multilib := first
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
      endif
    endif
<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
=======

    # Only preopt primary arch for translated arch since there is only an image there.
    ifeq ($(TARGET_TRANSLATE_2ND_ARCH),true)
      my_module_multilib := first
    endif

    # #################################################
    # Odex for the 1st arch
    my_dexpreopt_archs += $(TARGET_ARCH)
    my_dexpreopt_images += $(DEXPREOPT_IMAGE_$(my_dexpreopt_infix)_$(TARGET_ARCH))
    my_dexpreopt_images_deps += $(DEXPREOPT_IMAGE_DEPS_$(my_dexpreopt_infix)_$(TARGET_ARCH))
    # Odex for the 2nd arch
    ifdef TARGET_2ND_ARCH
      ifneq ($(TARGET_TRANSLATE_2ND_ARCH),true)
        ifneq (first,$(my_module_multilib))
          my_dexpreopt_archs += $(TARGET_2ND_ARCH)
          my_dexpreopt_images += $(DEXPREOPT_IMAGE_$(my_dexpreopt_infix)_$(TARGET_2ND_ARCH))
          my_dexpreopt_images_deps += $(DEXPREOPT_IMAGE_DEPS_$(my_dexpreopt_infix)_$(TARGET_2ND_ARCH))
        endif  # my_module_multilib is not first.
      endif  # TARGET_TRANSLATE_2ND_ARCH not true
    endif  # TARGET_2ND_ARCH
    # #################################################
  else  # must be APPS
    # The preferred arch
    # Save the module multilib since setup_one_odex modifies it.
    my_2nd_arch_prefix := $(LOCAL_2ND_ARCH_VAR_PREFIX)
    my_dexpreopt_archs += $(TARGET_$(my_2nd_arch_prefix)ARCH)
    my_dexpreopt_images += \
        $(DEXPREOPT_IMAGE_$(my_dexpreopt_infix)_$(TARGET_$(my_2nd_arch_prefix)ARCH))
    my_dexpreopt_images_deps += \
        $(DEXPREOPT_IMAGE_DEPS_$(my_dexpreopt_infix)_$(TARGET_$(my_2nd_arch_prefix)ARCH))
    ifdef TARGET_2ND_ARCH
      ifeq ($(my_module_multilib),both)
        # The non-preferred arch
        my_2nd_arch_prefix := $(if $(LOCAL_2ND_ARCH_VAR_PREFIX),,$(TARGET_2ND_ARCH_VAR_PREFIX))
        my_dexpreopt_archs += $(TARGET_$(my_2nd_arch_prefix)ARCH)
        my_dexpreopt_images += \
            $(DEXPREOPT_IMAGE_$(my_dexpreopt_infix)_$(TARGET_$(my_2nd_arch_prefix)ARCH))
        my_dexpreopt_images_deps += \
            $(DEXPREOPT_IMAGE_DEPS_$(my_dexpreopt_infix)_$(TARGET_$(my_2nd_arch_prefix)ARCH))
      endif  # LOCAL_MULTILIB is both
    endif  # TARGET_2ND_ARCH
  endif  # LOCAL_MODULE_CLASS

  my_dexpreopt_image_locations_on_host += $(DEXPREOPT_IMAGE_LOCATIONS_ON_HOST$(my_dexpreopt_infix))
  my_dexpreopt_image_locations_on_device += $(DEXPREOPT_IMAGE_LOCATIONS_ON_DEVICE$(my_dexpreopt_infix))

  # Record dex-preopt config.
  DEXPREOPT.$(LOCAL_MODULE).DEX_PREOPT := $(LOCAL_DEX_PREOPT)
  DEXPREOPT.$(LOCAL_MODULE).MULTILIB := $(LOCAL_MULTILIB)
  DEXPREOPT.$(LOCAL_MODULE).DEX_PREOPT_FLAGS := $(LOCAL_DEX_PREOPT_FLAGS)
  DEXPREOPT.$(LOCAL_MODULE).PRIVILEGED_MODULE := $(LOCAL_PRIVILEGED_MODULE)
  DEXPREOPT.$(LOCAL_MODULE).VENDOR_MODULE := $(LOCAL_VENDOR_MODULE)
  DEXPREOPT.$(LOCAL_MODULE).TARGET_ARCH := $(LOCAL_MODULE_TARGET_ARCH)
  DEXPREOPT.$(LOCAL_MODULE).INSTALLED_STRIPPED := $(LOCAL_INSTALLED_MODULE)
  DEXPREOPT.MODULES.$(LOCAL_MODULE_CLASS) := $(sort \
    $(DEXPREOPT.MODULES.$(LOCAL_MODULE_CLASS)) $(LOCAL_MODULE))

  $(call json_start)

  # DexPath is not set: it will be filled in by dexpreopt_gen.

  $(call add_json_str,  Name,                           $(LOCAL_MODULE))
  $(call add_json_str,  DexLocation,                    $(patsubst $(PRODUCT_OUT)%,%,$(LOCAL_INSTALLED_MODULE)))
  $(call add_json_str,  BuildPath,                      $(LOCAL_BUILT_MODULE))
  $(call add_json_str,  ManifestPath,                   $(full_android_manifest))
  $(call add_json_str,  ExtrasOutputPath,               $$2)
  $(call add_json_bool, Privileged,                     $(filter true,$(LOCAL_PRIVILEGED_MODULE)))
  $(call add_json_bool, UncompressedDex,                $(filter true,$(LOCAL_UNCOMPRESS_DEX)))
  $(call add_json_bool, HasApkLibraries,                $(LOCAL_APK_LIBRARIES))
  $(call add_json_list, PreoptFlags,                    $(LOCAL_DEX_PREOPT_FLAGS))
  $(call add_json_str,  ProfileClassListing,            $(if $(my_process_profile),$(LOCAL_DEX_PREOPT_PROFILE)))
  $(call add_json_bool, ProfileIsTextListing,           $(my_profile_is_text_listing))
  $(call add_json_str,  EnforceUsesLibrariesStatusFile, $(my_enforced_uses_libraries))
  $(call add_json_bool, EnforceUsesLibraries,           $(filter true,$(LOCAL_ENFORCE_USES_LIBRARIES)))
  $(call add_json_str,  ProvidesUsesLibrary,            $(firstword $(LOCAL_PROVIDES_USES_LIBRARY) $(LOCAL_MODULE)))
  $(call add_json_map,  ClassLoaderContexts)
  $(call add_json_class_loader_context, any, $(my_dexpreopt_libs))
  $(call add_json_class_loader_context,  28, $(my_dexpreopt_libs_compat_28))
  $(call add_json_class_loader_context,  29, $(my_dexpreopt_libs_compat_29))
  $(call add_json_class_loader_context,  30, $(my_dexpreopt_libs_compat_30))
  $(call end_json_map)
  $(call add_json_list, Archs,                          $(my_dexpreopt_archs))
  $(call add_json_list, DexPreoptImages,                $(my_dexpreopt_images))
  $(call add_json_list, DexPreoptImageLocationsOnHost,  $(my_dexpreopt_image_locations_on_host))
  $(call add_json_list, DexPreoptImageLocationsOnDevice,$(my_dexpreopt_image_locations_on_device))
  $(call add_json_list, PreoptBootClassPathDexFiles,    $(DEXPREOPT_BOOTCLASSPATH_DEX_FILES))
  $(call add_json_list, PreoptBootClassPathDexLocations,$(DEXPREOPT_BOOTCLASSPATH_DEX_LOCATIONS))
  $(call add_json_bool, PreoptExtractedApk,             $(my_preopt_for_extracted_apk))
  $(call add_json_bool, NoCreateAppImage,               $(filter false,$(LOCAL_DEX_PREOPT_APP_IMAGE)))
  $(call add_json_bool, ForceCreateAppImage,            $(filter true,$(LOCAL_DEX_PREOPT_APP_IMAGE)))
  $(call add_json_bool, PresignedPrebuilt,              $(filter PRESIGNED,$(LOCAL_CERTIFICATE)))

  $(call json_end)

  my_dexpreopt_config := $(intermediates)/dexpreopt.config
  my_dexpreopt_config_for_postprocessing := $(PRODUCT_OUT)/dexpreopt_config/$(LOCAL_MODULE)_dexpreopt.config
  my_dexpreopt_config_merger := $(BUILD_SYSTEM)/dex_preopt_config_merger.py

  $(my_dexpreopt_config): $(my_dexpreopt_dep_configs) $(my_dexpreopt_config_merger)
  $(my_dexpreopt_config): PRIVATE_MODULE := $(LOCAL_MODULE)
  $(my_dexpreopt_config): PRIVATE_CONTENTS := $(json_contents)
  $(my_dexpreopt_config): PRIVATE_DEP_CONFIGS := $(my_dexpreopt_dep_configs)
  $(my_dexpreopt_config): PRIVATE_CONFIG_MERGER := $(my_dexpreopt_config_merger)
  $(my_dexpreopt_config):
	@echo "$(PRIVATE_MODULE) dexpreopt.config"
	echo -e -n '$(subst $(newline),\n,$(subst ','\'',$(subst \,\\,$(PRIVATE_CONTENTS))))' > $@
	$(PRIVATE_CONFIG_MERGER) $@ $(PRIVATE_DEP_CONFIGS)

$(eval $(call copy-one-file,$(my_dexpreopt_config),$(my_dexpreopt_config_for_postprocessing)))

$(LOCAL_INSTALLED_MODULE): $(my_dexpreopt_config_for_postprocessing)

# System server jars defined in Android.mk are deprecated.
ifneq (true, $(PRODUCT_BROKEN_DEPRECATED_MK_SYSTEM_SERVER_JARS))
  ifneq (,$(filter %:$(LOCAL_MODULE), $(PRODUCT_SYSTEM_SERVER_JARS) $(PRODUCT_APEX_SYSTEM_SERVER_JARS)))
    $(error System server jars defined in Android.mk are deprecated. \
      Convert $(LOCAL_MODULE) to Android.bp or temporarily disable the error \
      with 'PRODUCT_BROKEN_DEPRECATED_MK_SYSTEM_SERVER_JARS := true')
  endif
endif

ifdef LOCAL_DEX_PREOPT
  # System server jars must be copied into predefined locations expected by
  # dexpreopt. Copy rule must be exposed to Ninja (as it uses these files as
  # inputs), so it cannot go in dexpreopt.sh.
  ifneq (,$(filter %:$(LOCAL_MODULE), $(PRODUCT_SYSTEM_SERVER_JARS)))
    my_dexpreopt_jar_copy := $(OUT_DIR)/soong/system_server_dexjars/$(LOCAL_MODULE).jar
    $(my_dexpreopt_jar_copy): PRIVATE_BUILT_MODULE := $(LOCAL_BUILT_MODULE)
    $(my_dexpreopt_jar_copy): $(LOCAL_BUILT_MODULE)
	  @cp $(PRIVATE_BUILT_MODULE) $@
  endif

  my_dexpreopt_script := $(intermediates)/dexpreopt.sh
  my_dexpreopt_zip := $(intermediates)/dexpreopt.zip
  .KATI_RESTAT: $(my_dexpreopt_script)
  $(my_dexpreopt_script): PRIVATE_MODULE := $(LOCAL_MODULE)
  $(my_dexpreopt_script): PRIVATE_GLOBAL_SOONG_CONFIG := $(DEX_PREOPT_SOONG_CONFIG_FOR_MAKE)
  $(my_dexpreopt_script): PRIVATE_GLOBAL_CONFIG := $(DEX_PREOPT_CONFIG_FOR_MAKE)
  $(my_dexpreopt_script): PRIVATE_MODULE_CONFIG := $(my_dexpreopt_config)
  $(my_dexpreopt_script): $(DEXPREOPT_GEN)
  $(my_dexpreopt_script): $(my_dexpreopt_jar_copy)
  $(my_dexpreopt_script): $(my_dexpreopt_config) $(DEX_PREOPT_SOONG_CONFIG_FOR_MAKE) $(DEX_PREOPT_CONFIG_FOR_MAKE)
	@echo "$(PRIVATE_MODULE) dexpreopt gen"
	$(DEXPREOPT_GEN) \
	-global_soong $(PRIVATE_GLOBAL_SOONG_CONFIG) \
	-global $(PRIVATE_GLOBAL_CONFIG) \
	-module $(PRIVATE_MODULE_CONFIG) \
	-dexpreopt_script $@ \
	-out_dir $(OUT_DIR)

  my_dexpreopt_deps := $(my_dex_jar)
  my_dexpreopt_deps += $(if $(my_process_profile),$(LOCAL_DEX_PREOPT_PROFILE))
  my_dexpreopt_deps += \
    $(foreach lib, $(my_dexpreopt_libs) $(my_dexpreopt_libs_compat), \
      $(call intermediates-dir-for,JAVA_LIBRARIES,$(lib),,COMMON)/javalib.jar)
  my_dexpreopt_deps += $(my_dexpreopt_images_deps)
  my_dexpreopt_deps += $(DEXPREOPT_BOOTCLASSPATH_DEX_FILES)
  ifeq ($(LOCAL_ENFORCE_USES_LIBRARIES),true)
    my_dexpreopt_deps += $(intermediates.COMMON)/enforce_uses_libraries.status
  endif

  $(my_dexpreopt_zip): PRIVATE_MODULE := $(LOCAL_MODULE)
  $(my_dexpreopt_zip): $(my_dexpreopt_deps)
  $(my_dexpreopt_zip): | $(DEXPREOPT_GEN_DEPS)
  $(my_dexpreopt_zip): .KATI_DEPFILE := $(my_dexpreopt_zip).d
  $(my_dexpreopt_zip): PRIVATE_DEX := $(my_dex_jar)
  $(my_dexpreopt_zip): PRIVATE_SCRIPT := $(my_dexpreopt_script)
  $(my_dexpreopt_zip): $(my_dexpreopt_script)
	@echo "$(PRIVATE_MODULE) dexpreopt"
	bash $(PRIVATE_SCRIPT) $(PRIVATE_DEX) $@

  ifdef LOCAL_POST_INSTALL_CMD
    # Add a shell command separator
    LOCAL_POST_INSTALL_CMD += &&
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
  endif
endif

my_generate_dm := $(PRODUCT_DEX_PREOPT_GENERATE_DM_FILES)
ifeq (,$(filter $(LOCAL_DEX_PREOPT_FLAGS),--compiler-filter=verify))
# Generating DM files only makes sense for verify, avoid doing for non verify compiler filter APKs.
my_generate_dm := false
endif

# No reason to use a dm file if the dex is already uncompressed.
ifeq ($(LOCAL_UNCOMPRESS_DEX),true)
my_generate_dm := false
endif

ifeq (true,$(my_generate_dm))
LOCAL_DEX_PREOPT_FLAGS += --copy-dex-files=false
LOCAL_DEX_PREOPT := nostripping
my_built_dm := $(dir $(LOCAL_BUILT_MODULE))generated.dm
my_installed_dm := $(patsubst %.apk,%,$(LOCAL_INSTALLED_MODULE)).dm
my_copied_vdex := $(dir $(LOCAL_BUILT_MODULE))primary.vdex
$(eval $(call copy-one-file,$(built_vdex),$(my_copied_vdex)))
$(my_built_dm): PRIVATE_INPUT_VDEX := $(my_copied_vdex)
$(my_built_dm): $(my_copied_vdex) $(ZIPTIME)
	$(hide) mkdir -p $(dir $@)
	$(hide) rm -f $@
	$(hide) zip -qD -j -X -9 $@ $(PRIVATE_INPUT_VDEX)
	$(ZIPTIME) $@
$(eval $(call copy-one-file,$(my_built_dm),$(my_installed_dm)))
endif

<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
# By default, emit debug info.
my_dexpreopt_debug_info := true
# If the global setting suppresses mini-debug-info, disable it.
ifeq (false,$(WITH_DEXPREOPT_DEBUG_INFO))
  my_dexpreopt_debug_info := false
endif

# PRODUCT_SYSTEM_SERVER_DEBUG_INFO overrides WITH_DEXPREOPT_DEBUG_INFO.
# PRODUCT_OTHER_JAVA_DEBUG_INFO overrides WITH_DEXPREOPT_DEBUG_INFO.
ifneq (,$(filter $(PRODUCT_SYSTEM_SERVER_JARS),$(LOCAL_MODULE)))
  ifeq (true,$(PRODUCT_SYSTEM_SERVER_DEBUG_INFO))
    my_dexpreopt_debug_info := true
  else ifeq (false,$(PRODUCT_SYSTEM_SERVER_DEBUG_INFO))
    my_dexpreopt_debug_info := false
  endif
else
  ifeq (true,$(PRODUCT_OTHER_JAVA_DEBUG_INFO))
    my_dexpreopt_debug_info := true
  else ifeq (false,$(PRODUCT_OTHER_JAVA_DEBUG_INFO))
    my_dexpreopt_debug_info := false
  endif
endif

# Never enable on eng.
ifeq (eng,$(filter eng, $(TARGET_BUILD_VARIANT)))
my_dexpreopt_debug_info := false
endif

# Add dex2oat flag for debug-info/no-debug-info.
ifeq (true,$(my_dexpreopt_debug_info))
  LOCAL_DEX_PREOPT_FLAGS += --generate-mini-debug-info
else ifeq (false,$(my_dexpreopt_debug_info))
  LOCAL_DEX_PREOPT_FLAGS += --no-generate-mini-debug-info
endif

# Set the compiler reason to 'prebuilt' to identify the oat files produced
# during the build, as opposed to compiled on the device.
LOCAL_DEX_PREOPT_FLAGS += --compilation-reason=prebuilt

$(built_odex): PRIVATE_DEX_PREOPT_FLAGS := $(LOCAL_DEX_PREOPT_FLAGS)
$(built_vdex): $(built_odex)
$(built_art): $(built_odex)
endif

ifneq (true,$(my_generate_dm))
  # Add the installed_odex to the list of installed files for this module if we aren't generating a
  # dm file.
  ALL_MODULES.$(my_register_name).INSTALLED += $(installed_odex)
  ALL_MODULES.$(my_register_name).INSTALLED += $(installed_vdex)
  ALL_MODULES.$(my_register_name).INSTALLED += $(installed_art)

  ALL_MODULES.$(my_register_name).BUILT_INSTALLED += $(built_installed_odex)
  ALL_MODULES.$(my_register_name).BUILT_INSTALLED += $(built_installed_vdex)
  ALL_MODULES.$(my_register_name).BUILT_INSTALLED += $(built_installed_art)

  # Make sure to install the .odex and .vdex when you run "make <module_name>"
  $(my_all_targets): $(installed_odex) $(installed_vdex) $(installed_art)
else
  ALL_MODULES.$(my_register_name).INSTALLED += $(my_installed_dm)
  ALL_MODULES.$(my_register_name).BUILT_INSTALLED += $(my_built_dm) $(my_installed_dm)

  # Make sure to install the .dm when you run "make <module_name>"
  $(my_all_targets): $(installed_dm)
endif

# Record dex-preopt config.
DEXPREOPT.$(LOCAL_MODULE).DEX_PREOPT := $(LOCAL_DEX_PREOPT)
DEXPREOPT.$(LOCAL_MODULE).MULTILIB := $(LOCAL_MULTILIB)
DEXPREOPT.$(LOCAL_MODULE).DEX_PREOPT_FLAGS := $(LOCAL_DEX_PREOPT_FLAGS)
DEXPREOPT.$(LOCAL_MODULE).PRIVILEGED_MODULE := $(LOCAL_PRIVILEGED_MODULE)
DEXPREOPT.$(LOCAL_MODULE).VENDOR_MODULE := $(LOCAL_VENDOR_MODULE)
DEXPREOPT.$(LOCAL_MODULE).TARGET_ARCH := $(LOCAL_MODULE_TARGET_ARCH)
DEXPREOPT.$(LOCAL_MODULE).INSTALLED := $(installed_odex)
DEXPREOPT.$(LOCAL_MODULE).INSTALLED_STRIPPED := $(LOCAL_INSTALLED_MODULE)
DEXPREOPT.MODULES.$(LOCAL_MODULE_CLASS) := $(sort \
  $(DEXPREOPT.MODULES.$(LOCAL_MODULE_CLASS)) $(LOCAL_MODULE))

=======
  my_dexpreopt_config :=
  my_dexpreopt_script :=
  my_dexpreopt_zip :=
  my_dexpreopt_config_for_postprocessing :=
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
endif # LOCAL_DEX_PREOPT
<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)

# Profile doesn't depend on LOCAL_DEX_PREOPT.
ALL_MODULES.$(my_register_name).INSTALLED += $(my_installed_profile)
ALL_MODULES.$(my_register_name).BUILT_INSTALLED += $(build_installed_profile)

my_process_profile :=

$(my_all_targets): $(my_installed_profile)
=======
endif # my_create_dexpreopt_config
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
