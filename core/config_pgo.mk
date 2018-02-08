################################################################################
# Global constants and utilities
################################################################################

# utility to split a variable and split by commas
strip_and_split = $(subst $(comma),$(space),$(strip $(1)))

my_global_pgo_profile_projects:=toolchain/pgo-profiles vendor/google_data/pgo-profiles

my_global_pgo_instrument_flag:=-fprofile-generate=/data/local/tmp
my_global_pgo_sampling_flag:=-gline-tables-only
my_global_pgo_profile_use_flags:=-fprofile-use=
my_global_pgo_profile_use_other_flags:=-Wno-backend-plugin

# Process and store PGO-specific variables
my_pgo_benchmarks := $(call strip_and_split,$(LOCAL_PGO_BENCHMARKS))
my_pgo_profile_file := $(strip $(LOCAL_PGO_PROFILE_FILE))
my_pgo_instrumentation := $(strip $(LOCAL_PGO_INSTRUMENTATION))
my_pgo_sampling := $(strip $(LOCAL_PGO_SAMPLING))
my_pgo_enable_profile_use := $(strip $(LOCAL_PGO_ENABLE_PROFILE_USE))
my_pgo_cflags := $(strip $(LOCAL_PGO_CFLAGS))

################################################################################
# Sanitize variables that are flags to be either 'true' or the empty string.
################################################################################

# Sanitize profiling flags by making them empty if not explicitly set to true.
ifneq ($(my_pgo_instrumentation),true)
  my_pgo_instrumentation :=
endif
ifneq ($(my_pgo_sampling),true)
  my_pgo_sampling:=
endif

# enable_profile_use is true by default.  Set it to true either if it is empty
# or explicitly set to true.
ifeq (,$(my_pgo_enable_profile_use))
  my_pgo_enable_profile_use := true
endif
ifneq (true,$(my_pgo_enable_profile_use))
  my_pgo_enable_profile_use :=
endif

################################################################################
# Validate PGO properties:
# - instrumentation and sampling are not simultaneously enabled
# - sampling is not supported yet
# - if any of benchmarks, profile_file, and profile_kind (instrumentation or
#   sampling) is specified, then all of them need to be specified.
################################################################################

my_pgo_any_profile_kinds_set := $(or $(my_pgo_instrumentation),$(my_pgo_sampling))
my_pgo_all_profile_kinds_set := $(and $(my_pgo_instrumentation),$(my_pgo_sampling))

my_pgo_any_flags_set := $(or $(my_pgo_any_profile_kinds_set),$(my_pgo_benchmarks),$(my_pgo_profile_file))
my_pgo_all_flags_set := $(and $(my_pgo_any_profile_kinds_set),$(my_pgo_benchmarks),$(my_pgo_profile_file))

ifneq (,$(my_pgo_all_profile_kinds_set))
  $(error "Both LOCAL_PGO_INSTRUMENTATION and LOCAL_PGO_SAMPLING are set to true")
endif

ifneq (,$(my_pgo_sampling))
  $(error "LOCAL_PGO_SAMPLING is not supported yet")
endif

ifneq (,$(my_pgo_any_flags_set))
  ifeq (,$(my_pgo_all_flags_set))
    $(error "Missing PGO properties.  One of more of LOCAL_PGO_BENCHMARKS, LOCAL_PGO_PROFILE_FILE, and profile kind (LOCAL_PGO_INSTRUMENTATION or LOCAL_PGO_SAMPLING) is missing")
  endif
endif

################################################################################
# Process PGO properties
################################################################################

# Extract current module's benchmarks that are specified in $ANDROID_PGO_INSTRUMENT
my_global_pgo_enabled_benchmarks := $(call strip_and_split,$(ANDROID_PGO_INSTRUMENT))
my_pgo_enabled_benchmarks := $(filter $(my_pgo_benchmarks), $(my_global_pgo_enabled_benchmarks))

# Check if profile-use is disabled for this module.
my_pgo_use_profiles := true

# Profile use is disabled for a module if
# if no PGO-specific variable is set (check $(my_pgo_benchmarks) as proxy),
ifeq (,$(my_pgo_benchmarks))
  my_pgo_use_profiles := false

# if this module is being instrumented
else ifneq (,$(my_pgo_enabled_benchmarks))
  my_pgo_use_profiles := false

# ANDROID_PGO_NO_PROFILE_USE is set
else ifeq ($(ANDROID_PGO_NO_PROFILE_USE),true)
  my_pgo_use_profiles := false

# LOCAL_PGO_ENABLE_PROFILE_USE is set to false
else ifeq (,$(my_pgo_enable_profile_use))
  my_pgo_use_profiles := false
endif

# Search and find the first valid profile_file relative to the PGO profile projects.
define set-pgo-profiles
$(firstword $(strip \
    $(foreach file,
              $(foreach dir,$(my_global_pgo_profile_projects),$(dir)/$(1)),
              $(if $(wildcard $(file)), $(file))
    )
))
endef

# Instrument if any of the benchmarks is enabled.
ifneq ($(my_pgo_enabled_benchmarks),)
  my_cflags += -fprofile-generate=/data/local/tmp $(my_pgo_cflags)
  my_ldflags += -fprofile-generate=/data/local/tmp

# Otherwise, use PGO profiles if profile use is not disabled and a profile file
# exists
else ifeq (true,$(my_pgo_use_profiles)
  my_pgo_full_profile_file := $(call set-pgo-profiles,$(my_pgo_profile_file))
  ifneq (,$(my_pgo_full_profile_file))
    my_cflags += $(my_global_pgo_profile_use_flags)$(my_pgo_full_profile_file) \
                 $(my_global_pgo_profile_use_other_flags)
    my_ldflags += $(my_global_pgo_profile_use_flags)$(my_pgo_full_profile_file)
    my_additional_dependencies += $(my_pgo_full_profile_file)
  endif
endif
