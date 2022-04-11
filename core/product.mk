#
# Copyright (C) 2007 The Android Open Source Project
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

#
# Functions for including AndroidProducts.mk files
# PRODUCT_MAKEFILES is set up in AndroidProducts.mks.
# Format of PRODUCT_MAKEFILES:
# <product_name>:<path_to_the_product_makefile>
# If the <product_name> is the same as the base file name (without dir
# and the .mk suffix) of the product makefile, "<product_name>:" can be
# omitted.

# Search for AndroidProducts.mks in the given dir.
# $(1): the path to the dir
define _search-android-products-files-in-dir
$(sort $(shell test -d $(1) && find -L $(1) \
  -maxdepth 6 \
  -name .git -prune \
  -o -name AndroidProducts.mk -print))
endef

#
# Returns the list of all AndroidProducts.mk files.
# $(call ) isn't necessary.
#
define _find-android-products-files
$(foreach d, device vendor product,$(call _search-android-products-files-in-dir,$(d))) \
  $(SRC_TARGET_DIR)/product/AndroidProducts.mk
endef

#
# Returns the sorted concatenation of PRODUCT_MAKEFILES
# variables set in the given AndroidProducts.mk files.
# $(1): the list of AndroidProducts.mk files.
#
define get-product-makefiles
$(sort \
  $(foreach f,$(1), \
    $(eval PRODUCT_MAKEFILES :=) \
    $(eval LOCAL_DIR := $(patsubst %/,%,$(dir $(f)))) \
    $(eval include $(f)) \
    $(PRODUCT_MAKEFILES) \
   ) \
  $(eval PRODUCT_MAKEFILES :=) \
  $(eval LOCAL_DIR :=) \
 )
endef

#
# Returns the sorted concatenation of all PRODUCT_MAKEFILES
# variables set in all AndroidProducts.mk files.
# $(call ) isn't necessary.
#
define get-all-product-makefiles
$(call get-product-makefiles,$(_find-android-products-files))
endef

#
# Functions for including product makefiles
#

_product_var_list := \
    PRODUCT_NAME \
    PRODUCT_MODEL \
    PRODUCT_LOCALES \
    PRODUCT_AAPT_CONFIG \
    PRODUCT_AAPT_PREF_CONFIG \
    PRODUCT_AAPT_PREBUILT_DPI \
    PRODUCT_PACKAGES \
    PRODUCT_PACKAGES_DEBUG \
    PRODUCT_PACKAGES_ENG \
    PRODUCT_PACKAGES_TESTS \
    PRODUCT_DEVICE \
    PRODUCT_MANUFACTURER \
    PRODUCT_BRAND \
    PRODUCT_PROPERTY_OVERRIDES \
    PRODUCT_DEFAULT_PROPERTY_OVERRIDES \
    PRODUCT_PRODUCT_PROPERTIES \
    PRODUCT_CHARACTERISTICS \
    PRODUCT_COPY_FILES \
    PRODUCT_OTA_PUBLIC_KEYS \
    PRODUCT_EXTRA_RECOVERY_KEYS \
    PRODUCT_PACKAGE_OVERLAYS \
    DEVICE_PACKAGE_OVERLAYS \
    PRODUCT_ENFORCE_RRO_EXCLUDED_OVERLAYS \
    PRODUCT_ENFORCE_RRO_TARGETS \
    PRODUCT_SDK_ATREE_FILES \
    PRODUCT_SDK_ADDON_NAME \
    PRODUCT_SDK_ADDON_COPY_FILES \
    PRODUCT_SDK_ADDON_COPY_MODULES \
    PRODUCT_SDK_ADDON_DOC_MODULES \
    PRODUCT_SDK_ADDON_SYS_IMG_SOURCE_PROP \
    PRODUCT_SOONG_NAMESPACES \
    PRODUCT_DEFAULT_WIFI_CHANNELS \
    PRODUCT_DEFAULT_DEV_CERTIFICATE \
    PRODUCT_RESTRICT_VENDOR_FILES \
    PRODUCT_VENDOR_KERNEL_HEADERS \
    PRODUCT_BOOT_JARS \
    PRODUCT_SUPPORTS_BOOT_SIGNER \
    PRODUCT_SUPPORTS_VBOOT \
    PRODUCT_SUPPORTS_VERITY \
    PRODUCT_SUPPORTS_VERITY_FEC \
    PRODUCT_OEM_PROPERTIES \
    PRODUCT_SYSTEM_DEFAULT_PROPERTIES \
    PRODUCT_SYSTEM_PROPERTY_BLACKLIST \
    PRODUCT_VENDOR_PROPERTY_BLACKLIST \
    PRODUCT_SYSTEM_SERVER_APPS \
    PRODUCT_SYSTEM_SERVER_JARS \
    PRODUCT_ALWAYS_PREOPT_EXTRACTED_APK \
    PRODUCT_DEXPREOPT_SPEED_APPS \
    PRODUCT_LOADED_BY_PRIVILEGED_MODULES \
    PRODUCT_VBOOT_SIGNING_KEY \
    PRODUCT_VBOOT_SIGNING_SUBKEY \
    PRODUCT_VERITY_SIGNING_KEY \
    PRODUCT_SYSTEM_VERITY_PARTITION \
    PRODUCT_VENDOR_VERITY_PARTITION \
    PRODUCT_PRODUCT_VERITY_PARTITION \
    PRODUCT_SYSTEM_SERVER_DEBUG_INFO \
    PRODUCT_OTHER_JAVA_DEBUG_INFO \
    PRODUCT_DEX_PREOPT_MODULE_CONFIGS \
    PRODUCT_DEX_PREOPT_DEFAULT_COMPILER_FILTER \
    PRODUCT_DEX_PREOPT_DEFAULT_FLAGS \
    PRODUCT_DEX_PREOPT_BOOT_FLAGS \
    PRODUCT_DEX_PREOPT_PROFILE_DIR \
    PRODUCT_DEX_PREOPT_BOOT_IMAGE_PROFILE_LOCATION \
    PRODUCT_DEX_PREOPT_GENERATE_DM_FILES \
    PRODUCT_USE_PROFILE_FOR_BOOT_IMAGE \
<<<<<<< HEAD   (c2b35d Merge "Merge empty history for sparse-8348651-L2230000095368)
    PRODUCT_SYSTEM_SERVER_COMPILER_FILTER \
    PRODUCT_SANITIZER_MODULE_CONFIGS \
    PRODUCT_SYSTEM_BASE_FS_PATH \
    PRODUCT_VENDOR_BASE_FS_PATH \
    PRODUCT_PRODUCT_BASE_FS_PATH \
    PRODUCT_SHIPPING_API_LEVEL \
    VENDOR_PRODUCT_RESTRICT_VENDOR_FILES \
    VENDOR_EXCEPTION_MODULES \
    VENDOR_EXCEPTION_PATHS \
    PRODUCT_ART_TARGET_INCLUDE_DEBUG_BUILD \
    PRODUCT_ART_USE_READ_BARRIER \
    PRODUCT_IOT \
    PRODUCT_SYSTEM_HEADROOM \
    PRODUCT_MINIMIZE_JAVA_DEBUG_INFO \
    PRODUCT_INTEGER_OVERFLOW_EXCLUDE_PATHS \
    PRODUCT_ADB_KEYS \
    PRODUCT_CFI_INCLUDE_PATHS \
    PRODUCT_CFI_EXCLUDE_PATHS \
    PRODUCT_COMPATIBLE_PROPERTY_OVERRIDE \
    PRODUCT_ACTIONABLE_COMPATIBLE_PROPERTY_DISABLE \
=======
    PRODUCT_USES_DEFAULT_ART_CONFIG \

_product_single_value_vars += PRODUCT_SYSTEM_SERVER_COMPILER_FILTER
# Per-module sanitizer configs
_product_list_vars += PRODUCT_SANITIZER_MODULE_CONFIGS
_product_single_value_vars += PRODUCT_SYSTEM_BASE_FS_PATH
_product_single_value_vars += PRODUCT_VENDOR_BASE_FS_PATH
_product_single_value_vars += PRODUCT_PRODUCT_BASE_FS_PATH
_product_single_value_vars += PRODUCT_SYSTEM_EXT_BASE_FS_PATH
_product_single_value_vars += PRODUCT_ODM_BASE_FS_PATH
_product_single_value_vars += PRODUCT_VENDOR_DLKM_BASE_FS_PATH
_product_single_value_vars += PRODUCT_ODM_DLKM_BASE_FS_PATH
_product_single_value_vars += PRODUCT_SYSTEM_DLKM_BASE_FS_PATH

# The first API level this product shipped with
_product_single_value_vars += PRODUCT_SHIPPING_API_LEVEL

_product_list_vars += VENDOR_PRODUCT_RESTRICT_VENDOR_FILES
_product_list_vars += VENDOR_EXCEPTION_MODULES
_product_list_vars += VENDOR_EXCEPTION_PATHS
# Whether the product wants to ship libartd. For rules and meaning, see art/Android.mk.
_product_single_value_vars += PRODUCT_ART_TARGET_INCLUDE_DEBUG_BUILD

# Make this art variable visible to soong_config.mk.
_product_single_value_vars += PRODUCT_ART_USE_READ_BARRIER

# Add reserved headroom to a system image.
_product_single_value_vars += PRODUCT_SYSTEM_HEADROOM

# Whether to save disk space by minimizing java debug info
_product_single_value_vars += PRODUCT_MINIMIZE_JAVA_DEBUG_INFO

# Whether any paths are excluded from sanitization when SANITIZE_TARGET=integer_overflow
_product_list_vars += PRODUCT_INTEGER_OVERFLOW_EXCLUDE_PATHS

_product_single_value_vars += PRODUCT_ADB_KEYS

# Whether any paths should have CFI enabled for components
_product_list_vars += PRODUCT_CFI_INCLUDE_PATHS

# Whether any paths are excluded from sanitization when SANITIZE_TARGET=cfi
_product_list_vars += PRODUCT_CFI_EXCLUDE_PATHS

# Whether the Scudo hardened allocator is disabled platform-wide
_product_single_value_vars += PRODUCT_DISABLE_SCUDO

# List of extra VNDK versions to be included
_product_list_vars += PRODUCT_EXTRA_VNDK_VERSIONS

# Whether APEX should be compressed or not
_product_single_value_vars += PRODUCT_COMPRESSED_APEX

# VNDK version of product partition. It can be 'current' if the product
# partitions uses PLATFORM_VNDK_VERSION.
_product_single_value_vars += PRODUCT_PRODUCT_VNDK_VERSION

_product_single_value_vars += PRODUCT_ENFORCE_ARTIFACT_PATH_REQUIREMENTS
_product_single_value_vars += PRODUCT_ENFORCE_ARTIFACT_SYSTEM_CERTIFICATE_REQUIREMENT
_product_list_vars += PRODUCT_ARTIFACT_SYSTEM_CERTIFICATE_REQUIREMENT_ALLOW_LIST
_product_list_vars += PRODUCT_ARTIFACT_PATH_REQUIREMENT_HINT
_product_list_vars += PRODUCT_ARTIFACT_PATH_REQUIREMENT_ALLOWED_LIST

# List of modules that should be forcefully unmarked from being LOCAL_PRODUCT_MODULE, and hence
# installed on /system directory by default.
_product_list_vars += PRODUCT_FORCE_PRODUCT_MODULES_TO_SYSTEM_PARTITION

# When this is true, dynamic partitions is retrofitted on a device that has
# already been launched without dynamic partitions. Otherwise, the device
# is launched with dynamic partitions.
# This flag implies PRODUCT_USE_DYNAMIC_PARTITIONS.
_product_single_value_vars += PRODUCT_RETROFIT_DYNAMIC_PARTITIONS

# When this is true, various build time as well as runtime debugfs restrictions are enabled.
_product_single_value_vars += PRODUCT_SET_DEBUGFS_RESTRICTIONS

# Other dynamic partition feature flags.PRODUCT_USE_DYNAMIC_PARTITION_SIZE and
# PRODUCT_BUILD_SUPER_PARTITION default to the value of PRODUCT_USE_DYNAMIC_PARTITIONS.
_product_single_value_vars += \
    PRODUCT_USE_DYNAMIC_PARTITIONS \
    PRODUCT_USE_DYNAMIC_PARTITION_SIZE \
    PRODUCT_BUILD_SUPER_PARTITION \

# If set, kernel configuration requirements are present in OTA package (and will be enforced
# during OTA). Otherwise, kernel configuration requirements are enforced in VTS.
# Devices that checks the running kernel (instead of the kernel in OTA package) should not
# set this variable to prevent OTA failures.
_product_list_vars += PRODUCT_OTA_ENFORCE_VINTF_KERNEL_REQUIREMENTS

# If set to true, this product builds a generic OTA package, which installs generic system images
# onto matching devices. The product may only build a subset of system images (e.g. only
# system.img), so devices need to install the package in a system-only OTA manner.
_product_single_value_vars += PRODUCT_BUILD_GENERIC_OTA_PACKAGE

_product_list_vars += PRODUCT_MANIFEST_PACKAGE_NAME_OVERRIDES
_product_list_vars += PRODUCT_PACKAGE_NAME_OVERRIDES
_product_list_vars += PRODUCT_CERTIFICATE_OVERRIDES

# Controls for whether different partitions are built for the current product.
_product_single_value_vars += PRODUCT_BUILD_SYSTEM_IMAGE
_product_single_value_vars += PRODUCT_BUILD_SYSTEM_OTHER_IMAGE
_product_single_value_vars += PRODUCT_BUILD_VENDOR_IMAGE
_product_single_value_vars += PRODUCT_BUILD_PRODUCT_IMAGE
_product_single_value_vars += PRODUCT_BUILD_SYSTEM_EXT_IMAGE
_product_single_value_vars += PRODUCT_BUILD_ODM_IMAGE
_product_single_value_vars += PRODUCT_BUILD_VENDOR_DLKM_IMAGE
_product_single_value_vars += PRODUCT_BUILD_ODM_DLKM_IMAGE
_product_single_value_vars += PRODUCT_BUILD_SYSTEM_DLKM_IMAGE
_product_single_value_vars += PRODUCT_BUILD_CACHE_IMAGE
_product_single_value_vars += PRODUCT_BUILD_RAMDISK_IMAGE
_product_single_value_vars += PRODUCT_BUILD_USERDATA_IMAGE
_product_single_value_vars += PRODUCT_BUILD_RECOVERY_IMAGE
_product_single_value_vars += PRODUCT_BUILD_BOOT_IMAGE
_product_single_value_vars += PRODUCT_BUILD_INIT_BOOT_IMAGE
_product_single_value_vars += PRODUCT_BUILD_DEBUG_BOOT_IMAGE
_product_single_value_vars += PRODUCT_BUILD_VENDOR_BOOT_IMAGE
_product_single_value_vars += PRODUCT_BUILD_VENDOR_KERNEL_BOOT_IMAGE
_product_single_value_vars += PRODUCT_BUILD_DEBUG_VENDOR_BOOT_IMAGE
_product_single_value_vars += PRODUCT_BUILD_VBMETA_IMAGE
_product_single_value_vars += PRODUCT_BUILD_SUPER_EMPTY_IMAGE
_product_single_value_vars += PRODUCT_BUILD_PVMFW_IMAGE

# List of boot jars delivered via updatable APEXes, following the same format as
# PRODUCT_BOOT_JARS.
_product_list_vars += PRODUCT_APEX_BOOT_JARS

# If set, device uses virtual A/B.
_product_single_value_vars += PRODUCT_VIRTUAL_AB_OTA

# If set, device uses virtual A/B Compression.
_product_single_value_vars += PRODUCT_VIRTUAL_AB_COMPRESSION

# If set, device retrofits virtual A/B.
_product_single_value_vars += PRODUCT_VIRTUAL_AB_OTA_RETROFIT

# If set, forcefully generate a non-A/B update package.
# Note: A device configuration should inherit from virtual_ab_ota_plus_non_ab.mk
# instead of setting this variable directly.
# Note: Use TARGET_OTA_ALLOW_NON_AB in the build system because
# TARGET_OTA_ALLOW_NON_AB takes the value of AB_OTA_UPDATER into account.
_product_single_value_vars += PRODUCT_OTA_FORCE_NON_AB_PACKAGE

# If set, Java module in product partition cannot use hidden APIs.
_product_single_value_vars += PRODUCT_ENFORCE_PRODUCT_PARTITION_INTERFACE

# If set, only java_sdk_library can be used at inter-partition dependency.
# Note: Build error if BOARD_VNDK_VERSION is not set while
#       PRODUCT_ENFORCE_INTER_PARTITION_JAVA_SDK_LIBRARY is true, because
#       PRODUCT_ENFORCE_INTER_PARTITION_JAVA_SDK_LIBRARY has no meaning if
#       BOARD_VNDK_VERSION is not set.
# Note: When PRODUCT_ENFORCE_PRODUCT_PARTITION_INTERFACE is not set, there are
#       no restrictions at dependency between system and product partition.
_product_single_value_vars += PRODUCT_ENFORCE_INTER_PARTITION_JAVA_SDK_LIBRARY

# Allowlist for PRODUCT_ENFORCE_INTER_PARTITION_JAVA_SDK_LIBRARY option.
# Listed modules are allowed at inter-partition dependency even if it isn't
# a java_sdk_library module.
_product_list_vars += PRODUCT_INTER_PARTITION_JAVA_LIBRARY_ALLOWLIST

_product_single_value_vars += PRODUCT_INSTALL_EXTRA_FLATTENED_APEXES

# Install a copy of the debug policy to the system_ext partition, and allow
# init-second-stage to load debug policy from system_ext.
# This option is only meant to be set by compliance GSI targets.
_product_single_value_vars += PRODUCT_INSTALL_DEBUG_POLICY_TO_SYSTEM_EXT

# If set, metadata files for the following artifacts will be generated.
# - system/framework/*.jar
# - system/framework/oat/<arch>/*.{oat,vdex,art}
# - system/etc/boot-image.prof
# - system/etc/dirty-image-objects
# One fsverity metadata container file per one input file will be generated in
# system.img, with a suffix ".fsv_meta". e.g. a container file for
# "/system/framework/foo.jar" will be "system/framework/foo.jar.fsv_meta".
_product_single_value_vars += PRODUCT_SYSTEM_FSVERITY_GENERATE_METADATA

# If true, sets the default for MODULE_BUILD_FROM_SOURCE. This overrides
# BRANCH_DEFAULT_MODULE_BUILD_FROM_SOURCE but not an explicitly set value.
_product_single_value_vars += PRODUCT_MODULE_BUILD_FROM_SOURCE

.KATI_READONLY := _product_single_value_vars _product_list_vars
_product_var_list :=$= $(_product_single_value_vars) $(_product_list_vars)
>>>>>>> BRANCH (697279 Merge "Version bump to TKB1.220411.001.A1 [core/build_id.mk])

define dump-product
$(info ==== $(1) ====)\
$(foreach v,$(_product_var_list),\
$(info PRODUCTS.$(1).$(v) := $(PRODUCTS.$(1).$(v))))\
$(info --------)
endef

define dump-products
$(foreach p,$(PRODUCTS),$(call dump-product,$(p)))
endef

#
# $(1): product to inherit
#
<<<<<<< HEAD   (c2b35d Merge "Merge empty history for sparse-8348651-L2230000095368)
# Does three things:
=======
# To be called from product makefiles, and is later evaluated during the import-nodes
# call below. It does the following:
>>>>>>> BRANCH (697279 Merge "Version bump to TKB1.220411.001.A1 [core/build_id.mk])
#  1. Inherits all of the variables from $1.
#  2. Records the inheritance in the .INHERITS_FROM variable
<<<<<<< HEAD   (c2b35d Merge "Merge empty history for sparse-8348651-L2230000095368)
#  3. Records that we've visited this node, in ALL_PRODUCTS
=======
#
# (2) and the PRODUCTS variable can be used together to reconstruct the include hierarchy
# See e.g. product-graph.mk for an example of this.
>>>>>>> BRANCH (697279 Merge "Version bump to TKB1.220411.001.A1 [core/build_id.mk])
#
define inherit-product
<<<<<<< HEAD   (c2b35d Merge "Merge empty history for sparse-8348651-L2230000095368)
  $(if $(findstring ../,$(1)),\
    $(eval np := $(call normalize-paths,$(1))),\
    $(eval np := $(strip $(1))))\
  $(foreach v,$(_product_var_list), \
      $(eval $(v) := $($(v)) $(INHERIT_TAG)$(np))) \
  $(eval inherit_var := \
      PRODUCTS.$(strip $(word 1,$(_include_stack))).INHERITS_FROM) \
  $(eval $(inherit_var) := $(sort $($(inherit_var)) $(np))) \
  $(eval inherit_var:=) \
  $(eval ALL_PRODUCTS := $(sort $(ALL_PRODUCTS) $(word 1,$(_include_stack))))
=======
  $(eval _inherit_product_wildcard := $(wildcard $(1)))\
  $(if $(_inherit_product_wildcard),,$(error $(1) does not exist.))\
  $(foreach part,$(_inherit_product_wildcard),\
    $(if $(findstring ../,$(part)),\
      $(eval np := $(call normalize-paths,$(part))),\
      $(eval np := $(strip $(part))))\
    $(foreach v,$(_product_var_list), \
        $(eval $(v) := $($(v)) $(INHERIT_TAG)$(np))) \
    $(eval current_mk := $(strip $(word 1,$(_include_stack)))) \
    $(eval inherit_var := PRODUCTS.$(current_mk).INHERITS_FROM) \
    $(eval $(inherit_var) := $(sort $($(inherit_var)) $(np))) \
    $(call dump-inherit,$(strip $(word 1,$(_include_stack))),$(1)) \
    $(call dump-config-vals,$(current_mk),inherit))
>>>>>>> BRANCH (697279 Merge "Version bump to TKB1.220411.001.A1 [core/build_id.mk])
endef


#
# Do inherit-product only if $(1) exists
#
define inherit-product-if-exists
  $(if $(wildcard $(1)),$(call inherit-product,$(1)),)
endef

#
# $(1): product makefile list
#
#TODO: check to make sure that products have all the necessary vars defined
define import-products
$(call import-nodes,PRODUCTS,$(1),$(_product_var_list))
endef


#
# Does various consistency checks on all of the known products.
# Takes no parameters, so $(call ) is not necessary.
#
define check-all-products
$(if ,, \
  $(eval _cap_names :=) \
  $(foreach p,$(PRODUCTS), \
    $(eval pn := $(strip $(PRODUCTS.$(p).PRODUCT_NAME))) \
    $(if $(pn),,$(error $(p): PRODUCT_NAME must be defined.)) \
    $(if $(filter $(pn),$(_cap_names)), \
      $(error $(p): PRODUCT_NAME must be unique; "$(pn)" already used by $(strip \
          $(foreach \
            pp,$(PRODUCTS),
              $(if $(filter $(pn),$(PRODUCTS.$(pp).PRODUCT_NAME)), \
                $(pp) \
               ))) \
       ) \
     ) \
    $(eval _cap_names += $(pn)) \
    $(if $(call is-c-identifier,$(pn)),, \
      $(error $(p): PRODUCT_NAME must be a valid C identifier, not "$(pn)") \
     ) \
    $(eval pb := $(strip $(PRODUCTS.$(p).PRODUCT_BRAND))) \
    $(if $(pb),,$(error $(p): PRODUCT_BRAND must be defined.)) \
    $(foreach cf,$(strip $(PRODUCTS.$(p).PRODUCT_COPY_FILES)), \
      $(if $(filter 2 3,$(words $(subst :,$(space),$(cf)))),, \
        $(error $(p): malformed COPY_FILE "$(cf)") \
       ) \
     ) \
   ) \
)
endef


#
# Returns the product makefile path for the product with the provided name
#
# $(1): short product name like "generic"
#
define _resolve-short-product-name
  $(eval pn := $(strip $(1)))
  $(eval p := \
      $(foreach p,$(PRODUCTS), \
          $(if $(filter $(pn),$(PRODUCTS.$(p).PRODUCT_NAME)), \
            $(p) \
       )) \
   )
  $(eval p := $(sort $(p)))
  $(if $(filter 1,$(words $(p))), \
    $(p), \
    $(if $(filter 0,$(words $(p))), \
      $(error No matches for product "$(pn)"), \
      $(error Product "$(pn)" ambiguous: matches $(p)) \
    ) \
  )
endef
define resolve-short-product-name
$(strip $(call _resolve-short-product-name,$(1)))
endef


_product_stash_var_list := $(_product_var_list) \
	PRODUCT_BOOTCLASSPATH \
	PRODUCT_SYSTEM_SERVER_CLASSPATH \
	TARGET_ARCH \
	TARGET_ARCH_VARIANT \
	TARGET_CPU_VARIANT \
	TARGET_BOARD_PLATFORM \
	TARGET_BOARD_PLATFORM_GPU \
	TARGET_BOARD_KERNEL_HEADERS \
	TARGET_DEVICE_KERNEL_HEADERS \
	TARGET_PRODUCT_KERNEL_HEADERS \
	TARGET_BOOTLOADER_BOARD_NAME \
	TARGET_NO_BOOTLOADER \
	TARGET_NO_KERNEL \
	TARGET_NO_RECOVERY \
	TARGET_NO_RADIOIMAGE \
	TARGET_HARDWARE_3D \
	TARGET_CPU_ABI \
	TARGET_CPU_ABI2 \


_product_stash_var_list += \
	BOARD_WPA_SUPPLICANT_DRIVER \
	BOARD_WLAN_DEVICE \
	BOARD_USES_GENERIC_AUDIO \
	BOARD_KERNEL_CMDLINE \
	BOARD_KERNEL_BASE \
	BOARD_HAVE_BLUETOOTH \
	BOARD_VENDOR_USE_AKMD \
	BOARD_EGL_CFG \
	BOARD_BOOTIMAGE_PARTITION_SIZE \
	BOARD_RECOVERYIMAGE_PARTITION_SIZE \
	BOARD_SYSTEMIMAGE_PARTITION_SIZE \
	BOARD_SYSTEMIMAGE_FILE_SYSTEM_TYPE \
	BOARD_USERDATAIMAGE_FILE_SYSTEM_TYPE \
	BOARD_USERDATAIMAGE_PARTITION_SIZE \
	BOARD_CACHEIMAGE_FILE_SYSTEM_TYPE \
	BOARD_CACHEIMAGE_PARTITION_SIZE \
	BOARD_FLASH_BLOCK_SIZE \
	BOARD_VENDORIMAGE_PARTITION_SIZE \
	BOARD_VENDORIMAGE_FILE_SYSTEM_TYPE \
	BOARD_PRODUCTIMAGE_PARTITION_SIZE \
	BOARD_PRODUCTIMAGE_FILE_SYSTEM_TYPE \
	BOARD_INSTALLER_CMDLINE \


_product_stash_var_list += \
	DEFAULT_SYSTEM_DEV_CERTIFICATE \
	WITH_DEXPREOPT \
	WITH_DEXPREOPT_BOOT_IMG_AND_SYSTEM_SERVER_ONLY

#
# Mark the variables in _product_stash_var_list as readonly
#
define readonly-product-vars
$(foreach v,$(_product_stash_var_list), \
	$(eval $(v) ?=) \
	$(eval .KATI_READONLY := $(v)) \
 )
endef

define add-to-product-copy-files-if-exists
$(if $(wildcard $(word 1,$(subst :, ,$(1)))),$(1))
endef

# whitespace placeholder when we record module's dex-preopt config.
_PDPMC_SP_PLACE_HOLDER := |@SP@|
# Set up dex-preopt config for a module.
# $(1) list of module names
# $(2) the modules' dex-preopt config
define add-product-dex-preopt-module-config
$(eval _c := $(subst $(space),$(_PDPMC_SP_PLACE_HOLDER),$(strip $(2))))\
$(eval PRODUCT_DEX_PREOPT_MODULE_CONFIGS += \
  $(foreach m,$(1),$(m)=$(_c)))
endef

# whitespace placeholder when we record module's sanitizer config.
_PSMC_SP_PLACE_HOLDER := |@SP@|
# Set up sanitizer config for a module.
# $(1) list of module names
# $(2) the modules' sanitizer config
define add-product-sanitizer-module-config
$(eval _c := $(subst $(space),$(_PSMC_SP_PLACE_HOLDER),$(strip $(2))))\
$(eval PRODUCT_SANITIZER_MODULE_CONFIGS += \
  $(foreach m,$(1),$(m)=$(_c)))
endef
