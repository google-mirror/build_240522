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

#
# Convert file file to the PRODUCT_COPY_FILES/PRODUCT_SDK_ADDON_COPY_FILES
# format: for each file F return $(F):$(PREFIX)/$(notdir $(F))
# $(1): files list
# $(2): prefix

define copy-files
$(foreach f,$(1),$(f):$(2)/$(notdir $(f)))
endef

#
# Convert the list of file names to the list of PRODUCT_COPY_FILES items
# $(1): from pattern
# $(2): to pattern
# $(3): file names
# E.g., calling product-copy-files-by-pattern with
#   (from/%, to/%, a b)
# returns
#   from/a:to/a from/b:to/b
define product-copy-files-by-pattern
$(join $(patsubst %,$(1),$(3)),$(patsubst %,:$(2),$(3)))
endef

# Return empty unless the board matches
define is-board-platform2
$(filter $(1), $(TARGET_BOARD_PLATFORM))
endef

# Return empty unless the board is in the list
define is-board-platform-in-list2
$(filter $(1),$(TARGET_BOARD_PLATFORM))
endef

# Return empty unless the board is QCOM
define is-vendor-board-qcom
$(if $(strip $(TARGET_BOARD_PLATFORM) $(QCOM_BOARD_PLATFORMS)),\
  $(filter $(TARGET_BOARD_PLATFORM),$(QCOM_BOARD_PLATFORMS)),\
  $(error both TARGET_BOARD_PLATFORM=$(TARGET_BOARD_PLATFORM) and QCOM_BOARD_PLATFORMS=$(QCOM_BOARD_PLATFORMS)))
endef

# ---------------------------------------------------------------

# These are the valid values of TARGET_BUILD_VARIANT.  Also, if anything else is passed
# as the variant in the PRODUCT-$TARGET_BUILD_PRODUCT-$TARGET_BUILD_VARIANT form,
# it will be treated as a goal, and the eng variant will be used.
INTERNAL_VALID_VARIANTS := user userdebug eng

# ---------------------------------------------------------------
# Provide "PRODUCT-<prodname>-<goal>" targets, which lets you build
# a particular configuration without needing to set up the environment.
#
ifeq ($(CALLED_FROM_SETUP),true)
product_goals := $(strip $(filter PRODUCT-%,$(MAKECMDGOALS)))
ifdef product_goals
  # Scrape the product and build names out of the goal,
  # which should be of the form PRODUCT-<productname>-<buildname>.
  #
  ifneq ($(words $(product_goals)),1)
    $(error Only one PRODUCT-* goal may be specified; saw "$(product_goals)")
  endif
  goal_name := $(product_goals)
  product_goals := $(patsubst PRODUCT-%,%,$(product_goals))
  product_goals := $(subst -, ,$(product_goals))
  ifneq ($(words $(product_goals)),2)
    $(error Bad PRODUCT-* goal "$(goal_name)")
  endif

  # The product they want
  TARGET_PRODUCT := $(word 1,$(product_goals))

  # The variant they want
  TARGET_BUILD_VARIANT := $(word 2,$(product_goals))

  ifeq ($(TARGET_BUILD_VARIANT),tests)
    $(error "tests" has been deprecated as a build variant. Use it as a build goal instead.)
  endif

  # The build server wants to do make PRODUCT-dream-sdk
  # which really means TARGET_PRODUCT=dream make sdk.
  ifneq ($(filter-out $(INTERNAL_VALID_VARIANTS),$(TARGET_BUILD_VARIANT)),)
    override MAKECMDGOALS := $(MAKECMDGOALS) $(TARGET_BUILD_VARIANT)
    TARGET_BUILD_VARIANT := userdebug
    default_goal_substitution :=
  else
    default_goal_substitution := droid
  endif

  # Replace the PRODUCT-* goal with the build goal that it refers to.
  # Note that this will ensure that it appears in the same relative
  # position, in case it matters.
  override MAKECMDGOALS := $(patsubst $(goal_name),$(default_goal_substitution),$(MAKECMDGOALS))
endif
endif # CALLED_FROM_SETUP
# else: Use the value set in the environment or buildspec.mk.

# ---------------------------------------------------------------
# Provide "APP-<appname>" targets, which lets you build
# an unbundled app.
#
ifeq ($(CALLED_FROM_SETUP),true)
unbundled_goals := $(strip $(filter APP-%,$(MAKECMDGOALS)))
ifdef unbundled_goals
  ifneq ($(words $(unbundled_goals)),1)
    $(error Only one APP-* goal may be specified; saw "$(unbundled_goals)")
  endif
  TARGET_BUILD_APPS := $(strip $(subst -, ,$(patsubst APP-%,%,$(unbundled_goals))))
  ifneq ($(filter droid,$(MAKECMDGOALS)),)
    override MAKECMDGOALS := $(patsubst $(unbundled_goals),,$(MAKECMDGOALS))
  else
    override MAKECMDGOALS := $(patsubst $(unbundled_goals),droid,$(MAKECMDGOALS))
  endif
endif # unbundled_goals
endif

# Now that we've parsed APP-* and PRODUCT-*, mark these as readonly
TARGET_BUILD_APPS ?=
.KATI_READONLY := \
  TARGET_PRODUCT \
  TARGET_BUILD_VARIANT \
  TARGET_BUILD_APPS

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

ifneq ($(ALLOW_RULES_IN_PRODUCT_CONFIG),)
_product_config_saved_KATI_ALLOW_RULES := $(.KATI_ALLOW_RULES)
.KATI_ALLOW_RULES := $(ALLOW_RULES_IN_PRODUCT_CONFIG)
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

ifndef RBC_PRODUCT_CONFIG
$(call import-products, $(current_product_makefile))
else
  $(shell build/soong/scripts/rbc-run $(current_product_makefile) \
      >$(OUT_DIR)/rbctemp.mk)
  ifneq ($(.SHELLSTATUS),0)
    $(error product configuration converter failed: $(.SHELLSTATUS))
  endif
  include $(OUT_DIR)/rbctemp.mk
  PRODUCTS += $(current_product_makefile)
endif
endif  # Import all or just the current product makefile

<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
# Sanity check
=======
ifndef RBC_PRODUCT_CONFIG
# Quick check
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
$(check-all-products)
endif

<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
=======
ifeq ($(SKIP_ARTIFACT_PATH_REQUIREMENT_PRODUCTS_CHECK),)
# Import all the products that have made artifact path requirements, so that we can verify
# the artifacts they produce.
# These are imported after check-all-products because some of them might not be real products.
$(foreach makefile,$(ARTIFACT_PATH_REQUIREMENT_PRODUCTS),\
  $(if $(filter-out $(makefile),$(PRODUCTS)),$(eval $(call import-products,$(makefile))))\
)
endif

ifneq ($(ALLOW_RULES_IN_PRODUCT_CONFIG),)
.KATI_ALLOW_RULES := $(_saved_KATI_ALLOW_RULES)
_product_config_saved_KATI_ALLOW_RULES :=
endif

>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
ifneq ($(filter dump-products, $(MAKECMDGOALS)),)
$(dump-products)
$(error done)
endif

ifndef RBC_PRODUCT_CONFIG
# Convert a short name like "sooner" into the path to the product
# file defining that product.
#
INTERNAL_PRODUCT := $(call resolve-short-product-name, $(TARGET_PRODUCT))
ifneq ($(current_product_makefile),$(INTERNAL_PRODUCT))
$(error PRODUCT_NAME inconsistent in $(current_product_makefile) and $(INTERNAL_PRODUCT))
endif


<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
=======
############################################################################
# Strip and assign the PRODUCT_ variables.
$(call strip-product-vars)
else
INTERNAL_PRODUCT := $(current_product_makefile)
endif

current_product_makefile :=
all_product_makefiles :=
all_product_configs :=
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)

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

<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
# Convert spaces to commas.
PRODUCT_AAPT_CONFIG := \
    $(subst $(space),$(comma),$(strip $(PRODUCT_AAPT_CONFIG)))
=======
###########################################################
## Add 'platform:' prefix to jars not in <apex>:<module> format.
##
## This makes sure that a jar corresponds to ConfigureJarList format of <apex> and <module> pairs
## where needed.
##
## $(1): a list of jars either in <module> or <apex>:<module> format
###########################################################

define qualify-platform-jars
  $(foreach jar,$(1),$(if $(findstring :,$(jar)),,platform:)$(jar))
endef

# Extra boot jars must be appended at the end after common boot jars.
PRODUCT_BOOT_JARS += $(PRODUCT_BOOT_JARS_EXTRA)
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)

<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
PRODUCT_BRAND := $(strip $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_BRAND))
=======
PRODUCT_BOOT_JARS := $(call qualify-platform-jars,$(PRODUCT_BOOT_JARS))

# b/191127295: force core-icu4j onto boot image. It comes from a non-updatable APEX jar, but has
# historically been part of the boot image; even though APEX jars are not meant to be part of the
# boot image.
# TODO(b/191686720): remove PRODUCT_APEX_BOOT_JARS to avoid a special handling of core-icu4j
# in make rules.
PRODUCT_APEX_BOOT_JARS := $(filter-out com.android.i18n:core-icu4j,$(PRODUCT_APEX_BOOT_JARS))
# All APEX jars come after /system and /system_ext jars, so adding core-icu4j at the end of the list
PRODUCT_BOOT_JARS += com.android.i18n:core-icu4j
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)

<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
PRODUCT_MODEL := $(strip $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_MODEL))
=======
# The extra system server jars must be appended at the end after common system server jars.
PRODUCT_SYSTEM_SERVER_JARS += $(PRODUCT_SYSTEM_SERVER_JARS_EXTRA)

PRODUCT_SYSTEM_SERVER_JARS := $(call qualify-platform-jars,$(PRODUCT_SYSTEM_SERVER_JARS))

ifndef PRODUCT_SYSTEM_NAME
  PRODUCT_SYSTEM_NAME := $(PRODUCT_NAME)
endif
ifndef PRODUCT_SYSTEM_DEVICE
  PRODUCT_SYSTEM_DEVICE := $(PRODUCT_DEVICE)
endif
ifndef PRODUCT_SYSTEM_BRAND
  PRODUCT_SYSTEM_BRAND := $(PRODUCT_BRAND)
endif
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
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

<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
# A list of words like <source path>:<destination path>[:<owner>].
# The file at the source path should be copied to the destination path
# when building  this product.  <destination path> is relative to
# $(PRODUCT_OUT), so it should look like, e.g., "system/etc/file.xml".
# The rules for these copy steps are defined in build/make/core/Makefile.
# The optional :<owner> is used to indicate the owner of a vendor file.
PRODUCT_COPY_FILES := \
    $(strip $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_COPY_FILES))
=======
$(foreach pair,$(PRODUCT_APEX_BOOT_JARS), \
  $(eval jar := $(call word-colon,2,$(pair))) \
  $(if $(findstring $(jar), $(PRODUCT_BOOT_JARS)), \
    $(error A jar in PRODUCT_APEX_BOOT_JARS must not be in PRODUCT_BOOT_JARS, but $(jar) is)))
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)

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

<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
# Whether any paths are excluded from sanitization when SANITIZE_TARGET=cfi
PRODUCT_CFI_EXCLUDE_PATHS := \
    $(strip $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_CFI_EXCLUDE_PATHS))
=======
ifdef PRODUCT_INSTALL_DEBUG_POLICY_TO_SYSTEM_EXT
  ifeq (,$(filter gsi_arm gsi_arm64 gsi_x86 gsi_x86_64,$(PRODUCT_NAME)))
    $(error Only GSI products are allowed to set PRODUCT_INSTALL_DEBUG_POLICY_TO_SYSTEM_EXT)
  endif
endif

ifndef PRODUCT_USE_DYNAMIC_PARTITIONS
  PRODUCT_USE_DYNAMIC_PARTITIONS := $(PRODUCT_RETROFIT_DYNAMIC_PARTITIONS)
endif
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)

# Whether any paths should have CFI enabled for components
PRODUCT_CFI_INCLUDE_PATHS := \
    $(strip $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_CFI_INCLUDE_PATHS))

# which Soong namespaces to export to Make
PRODUCT_SOONG_NAMESPACES := \
    $(strip $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_SOONG_NAMESPACES))

# A flag to override PRODUCT_COMPATIBLE_PROPERTY
PRODUCT_COMPATIBLE_PROPERTY_OVERRIDE := \
    $(strip $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_COMPATIBLE_PROPERTY_OVERRIDE))

<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
# Whether the whitelist of actionable compatible properties should be disabled or not
PRODUCT_ACTIONABLE_COMPATIBLE_PROPERTY_DISABLE := \
    $(strip $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_ACTIONABLE_COMPATIBLE_PROPERTY_DISABLE))
=======
ifeq ($(PRODUCT_SET_DEBUGFS_RESTRICTIONS),)
  ifdef PRODUCT_SHIPPING_API_LEVEL
    ifeq (true,$(call math_gt_or_eq,$(PRODUCT_SHIPPING_API_LEVEL),31))
      PRODUCT_SET_DEBUGFS_RESTRICTIONS := true
    endif
  endif
endif

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

###########################################
# APEXes are by default not compressed
#
# APEX compression can be forcibly enabled (resp. disabled) by
# setting OVERRIDE_PRODUCT_COMPRESSED_APEX to true (resp. false), e.g. by
# setting the OVERRIDE_PRODUCT_COMPRESSED_APEX environment variable.
ifdef OVERRIDE_PRODUCT_COMPRESSED_APEX
  PRODUCT_COMPRESSED_APEX := $(OVERRIDE_PRODUCT_COMPRESSED_APEX)
endif

$(KATI_obsolete_var OVERRIDE_PRODUCT_EXTRA_VNDK_VERSIONS \
    ,Use PRODUCT_EXTRA_VNDK_VERSIONS instead)

# If build command defines OVERRIDE_PRODUCT_ENFORCE_PRODUCT_PARTITION_INTERFACE,
# override PRODUCT_ENFORCE_PRODUCT_PARTITION_INTERFACE with it unless it is
# defined as `false`. If the value is `false` clear
# PRODUCT_ENFORCE_PRODUCT_PARTITION_INTERFACE
# OVERRIDE_PRODUCT_ENFORCE_PRODUCT_PARTITION_INTERFACE can be used for
# testing only.
ifdef OVERRIDE_PRODUCT_ENFORCE_PRODUCT_PARTITION_INTERFACE
  ifeq (false,$(OVERRIDE_PRODUCT_ENFORCE_PRODUCT_PARTITION_INTERFACE))
    PRODUCT_ENFORCE_PRODUCT_PARTITION_INTERFACE :=
  else
    PRODUCT_ENFORCE_PRODUCT_PARTITION_INTERFACE := $(OVERRIDE_PRODUCT_ENFORCE_PRODUCT_PARTITION_INTERFACE)
  endif
else ifeq ($(PRODUCT_SHIPPING_API_LEVEL),)
  # No shipping level defined
else ifeq ($(call math_gt,$(PRODUCT_SHIPPING_API_LEVEL),29),true)
  # Enforce product interface if PRODUCT_SHIPPING_API_LEVEL is greater than 29.
  PRODUCT_ENFORCE_PRODUCT_PARTITION_INTERFACE := true
endif

$(KATI_obsolete_var OVERRIDE_PRODUCT_ENFORCE_PRODUCT_PARTITION_INTERFACE,Use PRODUCT_ENFORCE_PRODUCT_PARTITION_INTERFACE instead)

# If build command defines PRODUCT_USE_PRODUCT_VNDK_OVERRIDE as `false`,
# PRODUCT_PRODUCT_VNDK_VERSION will not be defined automatically.
# PRODUCT_USE_PRODUCT_VNDK_OVERRIDE can be used for testing only.
PRODUCT_USE_PRODUCT_VNDK := false
ifneq ($(PRODUCT_USE_PRODUCT_VNDK_OVERRIDE),)
  PRODUCT_USE_PRODUCT_VNDK := $(PRODUCT_USE_PRODUCT_VNDK_OVERRIDE)
else ifeq ($(PRODUCT_SHIPPING_API_LEVEL),)
  # No shipping level defined
else ifeq ($(call math_gt,$(PRODUCT_SHIPPING_API_LEVEL),29),true)
  # Enforce product interface for VNDK if PRODUCT_SHIPPING_API_LEVEL is greater
  # than 29.
  PRODUCT_USE_PRODUCT_VNDK := true
endif

ifeq ($(PRODUCT_USE_PRODUCT_VNDK),true)
  ifndef PRODUCT_PRODUCT_VNDK_VERSION
    PRODUCT_PRODUCT_VNDK_VERSION := current
  endif
endif

$(KATI_obsolete_var PRODUCT_USE_PRODUCT_VNDK,Use PRODUCT_PRODUCT_VNDK_VERSION instead)
$(KATI_obsolete_var PRODUCT_USE_PRODUCT_VNDK_OVERRIDE,Use PRODUCT_PRODUCT_VNDK_VERSION instead)

ifdef PRODUCT_ENFORCE_RRO_EXEMPTED_TARGETS
    $(error PRODUCT_ENFORCE_RRO_EXEMPTED_TARGETS is deprecated, consider using RRO for \
      $(PRODUCT_ENFORCE_RRO_EXEMPTED_TARGETS))
endif

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
    PVMFW \
    SYSTEM \
    SYSTEM_OTHER \
    VENDOR \
    PRODUCT \
    SYSTEM_EXT \
    ODM \
    VENDOR_DLKM \
    ODM_DLKM \
    CACHE \
    RAMDISK \
    USERDATA \
    BOOT \
    RECOVERY, \
  $(eval $(call product-build-image-config,$(image))))

product-build-image-config :=

$(call readonly-product-vars)
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
