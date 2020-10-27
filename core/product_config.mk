#
# Copyright (C) 2008 The Android Open Source Project
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
#

# ---------------------------------------------------------------
# Generic functions
# TODO: Move these to definitions.make once we're able to include
# definitions.make before config.make.

###########################################################
## Return non-empty if $(1) is a C identifier; i.e., if it
## matches /^[a-zA-Z_][a-zA-Z0-9_]*$/.  We do this by first
## making sure that it isn't empty and doesn't start with
## a digit, then by removing each valid character.  If the
## final result is empty, then it was a valid C identifier.
##
## $(1): word to check
###########################################################

_ici_digits := 0 1 2 3 4 5 6 7 8 9
_ici_alphaunderscore := \
    a b c d e f g h i j k l m n o p q r s t u v w x y z \
    A B C D E F G H I J K L M N O P Q R S T U V W X Y Z _
define is-c-identifier
$(strip \
  $(if $(1), \
    $(if $(filter $(addsuffix %,$(_ici_digits)),$(1)), \
     , \
      $(eval w := $(1)) \
      $(foreach c,$(_ici_digits) $(_ici_alphaunderscore), \
        $(eval w := $(subst $(c),,$(w))) \
       ) \
      $(if $(w),,TRUE) \
      $(eval w :=) \
     ) \
   ) \
 )
endef

# TODO: push this into the combo files; unfortunately, we don't even
# know HOST_OS at this point.
trysed := $(shell echo a | sed -E -e 's/a/b/' 2>/dev/null)
ifeq ($(trysed),b)
  SED_EXTENDED := sed -E
else
  trysed := $(shell echo c | sed -r -e 's/c/d/' 2>/dev/null)
  ifeq ($(trysed),d)
    SED_EXTENDED := sed -r
  else
    $(error Unknown sed version)
  endif
endif

###########################################################
## List all of the files in a subdirectory in a format
## suitable for PRODUCT_COPY_FILES and
## PRODUCT_SDK_ADDON_COPY_FILES
##
## $(1): Glob to match file name
## $(2): Source directory
## $(3): Target base directory
###########################################################

define find-copy-subdir-files
$(sort $(shell find $(2) -name "$(1)" -type f | $(SED_EXTENDED) "s:($(2)/?(.*)):\\1\\:$(3)/\\2:" | sed "s://:/:g"))
endef

# ---------------------------------------------------------------
# Check for obsolete PRODUCT- and APP- goals
ifeq ($(CALLED_FROM_SETUP),true)
product_goals := $(strip $(filter PRODUCT-%,$(MAKECMDGOALS)))
ifdef product_goals
  $(error The PRODUCT-* goal is no longer supported. Use `TARGET_PRODUCT=<product> m droid` instead)
endif
unbundled_goals := $(strip $(filter APP-%,$(MAKECMDGOALS)))
ifdef unbundled_goals
  $(error The APP-* goal is no longer supported. Use `TARGET_BUILD_APPS="<app>" m droid` instead)
endif # unbundled_goals
endif

# Default to building dalvikvm on hosts that support it...
ifeq ($(HOST_OS),linux)
# ... or if the if the option is already set
ifeq ($(WITH_HOST_DALVIK),)
  WITH_HOST_DALVIK := true
endif
endif

# ---------------------------------------------------------------
# Include the product definitions.
# We need to do this to translate TARGET_PRODUCT into its
# underlying TARGET_DEVICE before we start defining any rules.
#
include $(BUILD_SYSTEM)/node_fns.mk
include $(BUILD_SYSTEM)/product.mk
include $(BUILD_SYSTEM)/device.mk

ifneq ($(strip $(TARGET_BUILD_APPS)),)
# An unbundled app build needs only the core product makefiles.
all_product_configs := $(call get-product-makefiles,\
    $(SRC_TARGET_DIR)/product/AndroidProducts.mk)
else
# Read in all of the product definitions specified by the AndroidProducts.mk
# files in the tree.
all_product_configs := $(get-all-product-makefiles)
endif

all_named_products :=

# Find the product config makefile for the current product.
# all_product_configs consists items like:
# <product_name>:<path_to_the_product_makefile>
# or just <path_to_the_product_makefile> in case the product name is the
# same as the base filename of the product config makefile.
current_product_makefile :=
all_product_makefiles :=
$(foreach f, $(all_product_configs),\
    $(eval _cpm_words := $(subst :,$(space),$(f)))\
    $(eval _cpm_word1 := $(word 1,$(_cpm_words)))\
    $(eval _cpm_word2 := $(word 2,$(_cpm_words)))\
    $(if $(_cpm_word2),\
        $(eval all_product_makefiles += $(_cpm_word2))\
        $(eval all_named_products += $(_cpm_word1))\
        $(if $(filter $(TARGET_PRODUCT),$(_cpm_word1)),\
            $(eval current_product_makefile += $(_cpm_word2)),),\
        $(eval all_product_makefiles += $(f))\
        $(eval all_named_products += $(basename $(notdir $(f))))\
        $(if $(filter $(TARGET_PRODUCT),$(basename $(notdir $(f)))),\
            $(eval current_product_makefile += $(f)),)))
_cpm_words :=
_cpm_word1 :=
_cpm_word2 :=
current_product_makefile := $(strip $(current_product_makefile))
all_product_makefiles := $(strip $(all_product_makefiles))

load_all_product_makefiles :=
ifneq (,$(filter product-graph, $(MAKECMDGOALS)))
ifeq ($(ANDROID_PRODUCT_GRAPH),--all)
load_all_product_makefiles := true
endif
endif
ifneq (,$(filter dump-products,$(MAKECMDGOALS)))
ifeq ($(ANDROID_DUMP_PRODUCTS),all)
load_all_product_makefiles := true
endif
endif

ifeq ($(load_all_product_makefiles),true)
# Import all product makefiles.
$(call import-products, $(all_product_makefiles))
else
# Import just the current product.
ifndef current_product_makefile
$(error Can not locate config makefile for product "$(TARGET_PRODUCT)")
endif
ifneq (1,$(words $(current_product_makefile)))
$(error Product "$(TARGET_PRODUCT)" ambiguous: matches $(current_product_makefile))
endif
$(call import-products, $(current_product_makefile))
endif  # Import all or just the current product makefile

# Sanity check
$(check-all-products)

ifneq ($(filter dump-products, $(MAKECMDGOALS)),)
$(dump-products)
$(error done)
endif

# Convert a short name like "sooner" into the path to the product
# file defining that product.
#
INTERNAL_PRODUCT := $(call resolve-short-product-name, $(TARGET_PRODUCT))
ifneq ($(current_product_makefile),$(INTERNAL_PRODUCT))
$(error PRODUCT_NAME inconsistent in $(current_product_makefile) and $(INTERNAL_PRODUCT))
endif
current_product_makefile :=
all_product_makefiles :=
all_product_configs :=

<<<<<<< HEAD   (5c8d84 Merge "Merge empty history for sparse-6676661-L8360000065797)
=======
# Jacoco agent JARS to be built and installed, if any.
ifeq ($(EMMA_INSTRUMENT),true)
  ifneq ($(EMMA_INSTRUMENT_STATIC),true)
    # For instrumented build, if Jacoco is not being included statically
    # in instrumented packages then include Jacoco classes into the
    # bootclasspath.
    $(foreach product,$(PRODUCTS),\
      $(eval PRODUCTS.$(product).PRODUCT_PACKAGES += jacocoagent)\
      $(eval PRODUCTS.$(product).PRODUCT_BOOT_JARS += jacocoagent))
  endif # EMMA_INSTRUMENT_STATIC
endif # EMMA_INSTRUMENT

############################################################################
# Strip and assign the PRODUCT_ variables.
$(call strip-product-vars)
>>>>>>> BRANCH (a10c18 Merge "Version bump to RT11.201014.001.A1 [core/build_id.mk])

#############################################################################

# A list of module names of BOOTCLASSPATH (jar files)
PRODUCT_BOOT_JARS := $(strip $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_BOOT_JARS))
PRODUCT_SYSTEM_SERVER_JARS := $(strip $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_SYSTEM_SERVER_JARS))
PRODUCT_SYSTEM_SERVER_APPS := $(strip $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_SYSTEM_SERVER_APPS))
PRODUCT_DEXPREOPT_SPEED_APPS := $(strip $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_DEXPREOPT_SPEED_APPS))
PRODUCT_LOADED_BY_PRIVILEGED_MODULES := $(strip $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_LOADED_BY_PRIVILEGED_MODULES))

# All of the apps that we force preopt, this overrides WITH_DEXPREOPT.
PRODUCT_ALWAYS_PREOPT_EXTRACTED_APK := $(strip $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_ALWAYS_PREOPT_EXTRACTED_APK))

# Find the device that this product maps to.
TARGET_DEVICE := $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_DEVICE)

# Figure out which resoure configuration options to use for this
# product.
PRODUCT_LOCALES := $(strip $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_LOCALES))
# TODO: also keep track of things like "port", "land" in product files.

# If CUSTOM_LOCALES contains any locales not already included
# in PRODUCT_LOCALES, add them to PRODUCT_LOCALES.
extra_locales := $(filter-out $(PRODUCT_LOCALES),$(CUSTOM_LOCALES))
ifneq (,$(extra_locales))
  ifneq ($(CALLED_FROM_SETUP),true)
    # Don't spam stdout, because envsetup.sh may be scraping values from it.
    $(info Adding CUSTOM_LOCALES [$(extra_locales)] to PRODUCT_LOCALES [$(PRODUCT_LOCALES)])
  endif
  PRODUCT_LOCALES += $(extra_locales)
  extra_locales :=
endif

# Add PRODUCT_LOCALES to PRODUCT_AAPT_CONFIG
PRODUCT_AAPT_CONFIG := $(strip $(PRODUCT_LOCALES) $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_AAPT_CONFIG))
PRODUCT_AAPT_PREF_CONFIG := $(strip $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_AAPT_PREF_CONFIG))
PRODUCT_AAPT_PREBUILT_DPI := $(strip $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_AAPT_PREBUILT_DPI))

# Keep a copy of the space-separated config
PRODUCT_AAPT_CONFIG_SP := $(PRODUCT_AAPT_CONFIG)

<<<<<<< HEAD   (5c8d84 Merge "Merge empty history for sparse-6676661-L8360000065797)
# Convert spaces to commas.
PRODUCT_AAPT_CONFIG := \
    $(subst $(space),$(comma),$(strip $(PRODUCT_AAPT_CONFIG)))

PRODUCT_BRAND := $(strip $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_BRAND))

PRODUCT_MODEL := $(strip $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_MODEL))
=======
# Extra boot jars must be appended at the end after common boot jars.
PRODUCT_BOOT_JARS += $(PRODUCT_BOOT_JARS_EXTRA)

# The extra system server jars must be appended at the end after common system server jars.
PRODUCT_SYSTEM_SERVER_JARS += $(PRODUCT_SYSTEM_SERVER_JARS_EXTRA)

ifndef PRODUCT_SYSTEM_NAME
  PRODUCT_SYSTEM_NAME := $(PRODUCT_NAME)
endif
ifndef PRODUCT_SYSTEM_DEVICE
  PRODUCT_SYSTEM_DEVICE := $(PRODUCT_DEVICE)
endif
ifndef PRODUCT_SYSTEM_BRAND
  PRODUCT_SYSTEM_BRAND := $(PRODUCT_BRAND)
endif
>>>>>>> BRANCH (a10c18 Merge "Version bump to RT11.201014.001.A1 [core/build_id.mk])
ifndef PRODUCT_MODEL
  PRODUCT_MODEL := $(strip $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_NAME))
endif

PRODUCT_MANUFACTURER := \
    $(strip $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_MANUFACTURER))
ifndef PRODUCT_MANUFACTURER
  PRODUCT_MANUFACTURER := unknown
endif

ifeq ($(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_CHARACTERISTICS),)
  TARGET_AAPT_CHARACTERISTICS := default
else
  TARGET_AAPT_CHARACTERISTICS := $(strip $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_CHARACTERISTICS))
endif

PRODUCT_DEFAULT_WIFI_CHANNELS := \
    $(strip $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_DEFAULT_WIFI_CHANNELS))

PRODUCT_DEFAULT_DEV_CERTIFICATE := \
    $(strip $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_DEFAULT_DEV_CERTIFICATE))
ifdef PRODUCT_DEFAULT_DEV_CERTIFICATE
ifneq (1,$(words $(PRODUCT_DEFAULT_DEV_CERTIFICATE)))
    $(error PRODUCT_DEFAULT_DEV_CERTIFICATE='$(PRODUCT_DEFAULT_DEV_CERTIFICATE)', \
      only 1 certificate is allowed.)
endif
endif

<<<<<<< HEAD   (5c8d84 Merge "Merge empty history for sparse-6676661-L8360000065797)
# A list of words like <source path>:<destination path>[:<owner>].
# The file at the source path should be copied to the destination path
# when building  this product.  <destination path> is relative to
# $(PRODUCT_OUT), so it should look like, e.g., "system/etc/file.xml".
# The rules for these copy steps are defined in build/make/core/Makefile.
# The optional :<owner> is used to indicate the owner of a vendor file.
PRODUCT_COPY_FILES := \
    $(strip $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_COPY_FILES))
=======
$(foreach pair,$(PRODUCT_UPDATABLE_BOOT_JARS), \
  $(if $(findstring $(call word-colon,2,$(pair)),$(PRODUCT_BOOT_JARS)), \
    $(error A jar in PRODUCT_UPDATABLE_BOOT_JARS must not be in PRODUCT_BOOT_JARS, \
      but $(call word-colon,2,$(pair)) is) \
  ) \
)

ENFORCE_SYSTEM_CERTIFICATE := $(PRODUCT_ENFORCE_ARTIFACT_SYSTEM_CERTIFICATE_REQUIREMENT)
ENFORCE_SYSTEM_CERTIFICATE_ALLOW_LIST := $(PRODUCT_ARTIFACT_SYSTEM_CERTIFICATE_REQUIREMENT_ALLOW_LIST)
>>>>>>> BRANCH (a10c18 Merge "Version bump to RT11.201014.001.A1 [core/build_id.mk])

# A list of property assignments, like "key = value", with zero or more
# whitespace characters on either side of the '='.
PRODUCT_PROPERTY_OVERRIDES := \
    $(strip $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_PROPERTY_OVERRIDES))
.KATI_READONLY := PRODUCT_PROPERTY_OVERRIDES

PRODUCT_SHIPPING_API_LEVEL := $(strip $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_SHIPPING_API_LEVEL))

# A list of property assignments, like "key = value", with zero or more
# whitespace characters on either side of the '='.
# used for adding properties to default.prop
PRODUCT_DEFAULT_PROPERTY_OVERRIDES := \
    $(strip $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_DEFAULT_PROPERTY_OVERRIDES))
.KATI_READONLY := PRODUCT_DEFAULT_PROPERTY_OVERRIDES

# A list of property assignments, like "key = value", with zero or more
# whitespace characters on either side of the '='.
# used for adding properties to default.prop of system partition
PRODUCT_SYSTEM_DEFAULT_PROPERTIES := \
    $(strip $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_SYSTEM_DEFAULT_PROPERTIES))
.KATI_READONLY := PRODUCT_SYSTEM_DEFAULT_PROPERTIES

# A list of property assignments, like "key = value", with zero or more
# whitespace characters on either side of the '='.
# used for adding properties to build.prop of product partition
PRODUCT_PRODUCT_PROPERTIES := \
    $(strip $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_PRODUCT_PROPERTIES))
.KATI_READONLY := PRODUCT_PRODUCT_PROPERTIES

# Should we use the default resources or add any product specific overlays
PRODUCT_PACKAGE_OVERLAYS := \
    $(strip $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_PACKAGE_OVERLAYS))
DEVICE_PACKAGE_OVERLAYS := \
        $(strip $(PRODUCTS.$(INTERNAL_PRODUCT).DEVICE_PACKAGE_OVERLAYS))

# The list of product-specific kernel header dirs
PRODUCT_VENDOR_KERNEL_HEADERS := \
    $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_VENDOR_KERNEL_HEADERS)

# The OTA key(s) specified by the product config, if any.  The names
# of these keys are stored in the target-files zip so that post-build
# signing tools can substitute them for the test key embedded by
# default.
PRODUCT_OTA_PUBLIC_KEYS := $(sort \
    $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_OTA_PUBLIC_KEYS))

PRODUCT_EXTRA_RECOVERY_KEYS := $(sort \
    $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_EXTRA_RECOVERY_KEYS))

PRODUCT_DEX_PREOPT_DEFAULT_COMPILER_FILTER := \
    $(strip $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_DEX_PREOPT_DEFAULT_COMPILER_FILTER))
PRODUCT_DEX_PREOPT_DEFAULT_FLAGS := \
    $(strip $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_DEX_PREOPT_DEFAULT_FLAGS))
PRODUCT_DEX_PREOPT_GENERATE_DM_FILES := \
    $(strip $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_DEX_PREOPT_GENERATE_DM_FILES))
PRODUCT_DEX_PREOPT_BOOT_FLAGS := \
    $(strip $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_DEX_PREOPT_BOOT_FLAGS))
PRODUCT_DEX_PREOPT_PROFILE_DIR := \
    $(strip $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_DEX_PREOPT_PROFILE_DIR))

# Boot image options.
PRODUCT_USE_PROFILE_FOR_BOOT_IMAGE := \
    $(strip $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_USE_PROFILE_FOR_BOOT_IMAGE))
PRODUCT_DEX_PREOPT_BOOT_IMAGE_PROFILE_LOCATION := \
    $(strip $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_DEX_PREOPT_BOOT_IMAGE_PROFILE_LOCATION))

PRODUCT_SYSTEM_SERVER_COMPILER_FILTER := \
    $(strip $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_SYSTEM_SERVER_COMPILER_FILTER))
PRODUCT_SYSTEM_SERVER_DEBUG_INFO := \
    $(strip $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_SYSTEM_SERVER_DEBUG_INFO))
PRODUCT_OTHER_JAVA_DEBUG_INFO := \
    $(strip $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_OTHER_JAVA_DEBUG_INFO))

# Resolve and setup per-module dex-preopt configs.
PRODUCT_DEX_PREOPT_MODULE_CONFIGS := \
    $(strip $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_DEX_PREOPT_MODULE_CONFIGS))
# If a module has multiple setups, the first takes precedence.
_pdpmc_modules :=
$(foreach c,$(PRODUCT_DEX_PREOPT_MODULE_CONFIGS),\
  $(eval m := $(firstword $(subst =,$(space),$(c))))\
  $(if $(filter $(_pdpmc_modules),$(m)),,\
    $(eval _pdpmc_modules += $(m))\
    $(eval cf := $(patsubst $(m)=%,%,$(c)))\
    $(eval cf := $(subst $(_PDPMC_SP_PLACE_HOLDER),$(space),$(cf)))\
    $(eval DEXPREOPT.$(TARGET_PRODUCT).$(m).CONFIG := $(cf))))
_pdpmc_modules :=

# Resolve and setup per-module sanitizer configs.
PRODUCT_SANITIZER_MODULE_CONFIGS := \
    $(strip $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_SANITIZER_MODULE_CONFIGS))
# If a module has multiple setups, the first takes precedence.
_psmc_modules :=
$(foreach c,$(PRODUCT_SANITIZER_MODULE_CONFIGS),\
  $(eval m := $(firstword $(subst =,$(space),$(c))))\
  $(if $(filter $(_psmc_modules),$(m)),,\
    $(eval _psmc_modules += $(m))\
    $(eval cf := $(patsubst $(m)=%,%,$(c)))\
    $(eval cf := $(subst $(_PSMC_SP_PLACE_HOLDER),$(space),$(cf)))\
    $(eval SANITIZER.$(TARGET_PRODUCT).$(m).CONFIG := $(cf))))
_psmc_modules :=

# Whether the product wants to ship libartd. For rules and meaning, see art/Android.mk.
PRODUCT_ART_TARGET_INCLUDE_DEBUG_BUILD := \
    $(strip $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_ART_TARGET_INCLUDE_DEBUG_BUILD))

# Make this art variable visible to soong_config.mk.
PRODUCT_ART_USE_READ_BARRIER := \
    $(strip $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_ART_USE_READ_BARRIER))

# Whether the product is an Android Things variant.
PRODUCT_IOT := \
    $(strip $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_IOT))

# Resource overlay list which must be excluded from enforcing RRO.
PRODUCT_ENFORCE_RRO_EXCLUDED_OVERLAYS := \
    $(strip $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_ENFORCE_RRO_EXCLUDED_OVERLAYS))

# Package list to apply enforcing RRO.
PRODUCT_ENFORCE_RRO_TARGETS := \
    $(strip $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_ENFORCE_RRO_TARGETS))

# Add reserved headroom to a system image.
PRODUCT_SYSTEM_HEADROOM := \
    $(strip $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_SYSTEM_HEADROOM))

# Whether to save disk space by minimizing java debug info
PRODUCT_MINIMIZE_JAVA_DEBUG_INFO := \
    $(strip $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_MINIMIZE_JAVA_DEBUG_INFO))

# Whether any paths are excluded from sanitization when SANITIZE_TARGET=integer_overflow
PRODUCT_INTEGER_OVERFLOW_EXCLUDE_PATHS := \
    $(strip $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_INTEGER_OVERFLOW_EXCLUDE_PATHS))

# ADB keys for debuggable builds
PRODUCT_ADB_KEYS :=
ifneq ($(filter eng userdebug,$(TARGET_BUILD_VARIANT)),)
  PRODUCT_ADB_KEYS := $(strip $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_ADB_KEYS))
endif
ifneq ($(filter-out 0 1,$(words $(PRODUCT_ADB_KEYS))),)
  $(error Only one file may be in PRODUCT_ADB_KEYS: $(PRODUCT_ADB_KEYS))
endif
.KATI_READONLY := PRODUCT_ADB_KEYS

# Whether any paths are excluded from sanitization when SANITIZE_TARGET=cfi
PRODUCT_CFI_EXCLUDE_PATHS := \
    $(strip $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_CFI_EXCLUDE_PATHS))

# Whether any paths should have CFI enabled for components
PRODUCT_CFI_INCLUDE_PATHS := \
    $(strip $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_CFI_INCLUDE_PATHS))

# which Soong namespaces to export to Make
PRODUCT_SOONG_NAMESPACES := \
    $(strip $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_SOONG_NAMESPACES))

# A flag to override PRODUCT_COMPATIBLE_PROPERTY
PRODUCT_COMPATIBLE_PROPERTY_OVERRIDE := \
    $(strip $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_COMPATIBLE_PROPERTY_OVERRIDE))

<<<<<<< HEAD   (5c8d84 Merge "Merge empty history for sparse-6676661-L8360000065797)
# Whether the whitelist of actionable compatible properties should be disabled or not
PRODUCT_ACTIONABLE_COMPATIBLE_PROPERTY_DISABLE := \
    $(strip $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_ACTIONABLE_COMPATIBLE_PROPERTY_DISABLE))
=======
ifdef PRODUCT_SHIPPING_API_LEVEL
  ifneq (,$(call math_gt_or_eq,29,$(PRODUCT_SHIPPING_API_LEVEL)))
    PRODUCT_PACKAGES += $(PRODUCT_PACKAGES_SHIPPING_API_LEVEL_29)
  endif
endif

# If build command defines OVERRIDE_PRODUCT_EXTRA_VNDK_VERSIONS,
# override PRODUCT_EXTRA_VNDK_VERSIONS with it.
ifdef OVERRIDE_PRODUCT_EXTRA_VNDK_VERSIONS
  PRODUCT_EXTRA_VNDK_VERSIONS := $(OVERRIDE_PRODUCT_EXTRA_VNDK_VERSIONS)
endif

$(KATI_obsolete_var OVERRIDE_PRODUCT_EXTRA_VNDK_VERSIONS \
    ,Use PRODUCT_EXTRA_VNDK_VERSIONS instead)

ifdef OVERRIDE_PRODUCT_ENFORCE_PRODUCT_PARTITION_INTERFACE
  ifeq (false,$(OVERRIDE_PRODUCT_ENFORCE_PRODUCT_PARTITION_INTERFACE))
    PRODUCT_ENFORCE_PRODUCT_PARTITION_INTERFACE :=
  else
    PRODUCT_ENFORCE_PRODUCT_PARTITION_INTERFACE := $(OVERRIDE_PRODUCT_ENFORCE_PRODUCT_PARTITION_INTERFACE)
  endif
else ifeq ($(PRODUCT_SHIPPING_API_LEVEL),)
  # No shipping level defined
else ifeq ($(call math_gt,$(PRODUCT_SHIPPING_API_LEVEL),29),true)
  PRODUCT_ENFORCE_PRODUCT_PARTITION_INTERFACE := true
endif

$(KATI_obsolete_var OVERRIDE_PRODUCT_ENFORCE_PRODUCT_PARTITION_INTERFACE,Use PRODUCT_ENFORCE_PRODUCT_PARTITION_INTERFACE instead)

define product-overrides-config
$$(foreach rule,$$(PRODUCT_$(1)_OVERRIDES),\
    $$(if $$(filter 2,$$(words $$(subst :,$$(space),$$(rule)))),,\
        $$(error Rule "$$(rule)" in PRODUCT_$(1)_OVERRIDE is not <module_name>:<new_value>)))
endef

$(foreach var, \
    MANIFEST_PACKAGE_NAME \
    PACKAGE_NAME \
    CERTIFICATE, \
  $(eval $(call product-overrides-config,$(var))))

# Macro to use below. $(1) is the name of the partition
define product-build-image-config
ifneq ($$(filter-out true false,$$(PRODUCT_BUILD_$(1)_IMAGE)),)
    $$(error Invalid PRODUCT_BUILD_$(1)_IMAGE: $$(PRODUCT_BUILD_$(1)_IMAGE) -- true false and empty are supported)
endif
endef

# Copy and check the value of each PRODUCT_BUILD_*_IMAGE variable
$(foreach image, \
    SYSTEM \
    SYSTEM_OTHER \
    VENDOR \
    PRODUCT \
    SYSTEM_EXT \
    ODM \
    CACHE \
    RAMDISK \
    USERDATA \
    BOOT \
    RECOVERY, \
  $(eval $(call product-build-image-config,$(image))))

product-build-image-config :=

$(call readonly-product-vars)
>>>>>>> BRANCH (a10c18 Merge "Version bump to RT11.201014.001.A1 [core/build_id.mk])
