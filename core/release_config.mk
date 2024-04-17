# Copyright (C) 2023 The Android Open Source Project
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


# -----------------------------------------------------------------
# Choose the flag files
# -----------------------------------------------------------------
# Release configs are defined in reflease_config_map files, which map
# the short name (e.g. -next) used in lunch to the starlark files
# defining the build flag values.
#
# (If you're thinking about aconfig flags, there is one build flag,
# RELEASE_ACONFIG_VALUE_SETS, that sets which aconfig_value_set
# module to use to set the aconfig flag values.)
#
# The short release config names *can* appear multiple times, to allow
# for AOSP and vendor specific flags under the same name, but the
# individual flag values must appear in exactly one config.  Vendor
# does not override AOSP, or anything like that.  This is because
# vendor code usually includes prebuilts, and having vendor compile
# with different flags from AOSP increases the likelihood of flag
# mismatch.

# Do this first, because we're going to unset TARGET_RELEASE before
# including anyone, so they don't start making conditionals based on it.
# This logic is in make because starlark doesn't understand optional
# vendor files.

# If this is a google source tree, restrict it to only the one file
# which has OWNERS control.  If it isn't let others define their own.
# TODO: Remove wildcard for build/release one when all branch manifests
# have updated.

config_map_files := $(wildcard build/release/release_config_map.textproto) \
    $(wildcard vendor/google_shared/build/release/release_config_map.textproto) \
    $(if $(wildcard vendor/google/release/release_config_map.textproto), \
        vendor/google/release/release_config_map.textproto, \
        $(sort \
            $(wildcard device/*/release/release_config_map.textproto) \
            $(wildcard device/*/*/release/release_config_map.textproto) \
            $(wildcard vendor/*/release/release_config_map.textproto) \
            $(wildcard vendor/*/*/release/release_config_map.textproto) \
        ) \
    )

# PRODUCT_RELEASE_CONFIG_MAPS is set by Soong using an initial run of product
# config to capture only the list of config maps needed by the build.
# Keep them in the order provided, but remove duplicates.
$(foreach map,$(PRODUCT_RELEASE_CONFIG_MAPS), \
    $(if $(filter $(map),$(config_map_files)),,$(eval config_map_files += $(map))) \
)

_args := $(foreach map,$(config_map_files), --map $(map) )
$(KATI_shell_no_rerun $(OUT_DIR)/release-config $(_args) >$(OUT_DIR)/release-config.out)
$(if $(filter-out 0,$(.SHELLSTATUS)),$(error release-config failed to run))
# This will also set _all_release_configs for us.
$(eval include $(OUT_DIR)/soong/release-config/release_config-$(TARGET_PRODUCT)-$(TARGET_RELEASE).mk)
$(KATI_extra_file_deps $(config_map_files))

ifeq ($(TARGET_RELEASE),)
    # We allow some internal paths to explicitly set TARGET_RELEASE to the
    # empty string.  For the most part, 'make' treats unset and empty string as
    # the same.  But the following line differentiates, and will only assign
    # if the variable was completely unset.
    TARGET_RELEASE ?= was_unset
    ifeq ($(TARGET_RELEASE),was_unset)
        $(error No release config set for target; please set TARGET_RELEASE, or if building on the command line use 'lunch <target>-<release>-<build_type>', where release is one of: $(_all_release_configs))
    endif
    # Instead of leaving this string empty, we want to default to a valid
    # setting.  Full builds coming through this path is a bug, but in case
    # of such a bug, we want to at least get consistent, valid results.
    # This matches the default value used in release-config.
    TARGET_RELEASE = trunk_staging
endif

# During pass 1 of product config, using a non-existent release config is not an error.
# We can safely assume that we are doing pass 1 if DUMP_MANY_VARS=="PRODUCT_RELEASE_CONFIG_MAPS".
ifneq (PRODUCT_RELEASE_CONFIG_MAPS,$(DUMP_MANY_VARS))
    ifeq ($(filter $(_all_release_configs), $(TARGET_RELEASE)),)
        $(error No release config found for TARGET_RELEASE: $(TARGET_RELEASE). Available releases are: $(_all_release_configs))
    endif
endif

# TODO: Remove this check after enough people have sourced lunch that we don't
# need to worry about it trying to do get_build_vars TARGET_RELEASE. Maybe after ~9/2023
ifneq ($(CALLED_FROM_SETUP),true)
define TARGET_RELEASE
$(error TARGET_RELEASE may not be accessed directly. Use individual flags.)
endef
else
TARGET_RELEASE:=
endif
.KATI_READONLY := TARGET_RELEASE

_all_release_configs:=
config_map_files:=
