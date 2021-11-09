# Put some miscellaneous rules here

# HACK: clear LOCAL_PATH from including last build target before calling
# intermedites-dir-for
LOCAL_PATH := $(BUILD_SYSTEM)

# -----------------------------------------------------------------
# Define rules to copy PRODUCT_COPY_FILES defined by the product.
# PRODUCT_COPY_FILES contains words like <source file>:<dest file>[:<owner>].
# <dest file> is relative to $(PRODUCT_OUT), so it should look like,
# e.g., "system/etc/file.xml".
# The filter part means "only eval the copy-one-file rule if this
# src:dest pair is the first one to match the same dest"
#$(1): the src:dest pair
#$(2): the dest
define check-product-copy-files
$(if $(filter-out $(TARGET_COPY_OUT_SYSTEM_OTHER)/%,$(2)), \
  $(if $(filter %.apk, $(2)),$(error \
     Prebuilt apk found in PRODUCT_COPY_FILES: $(1), use BUILD_PREBUILT instead!)))
endef
# filter out the duplicate <source file>:<dest file> pairs.
unique_product_copy_files_pairs :=
$(foreach cf,$(PRODUCT_COPY_FILES), \
    $(if $(filter $(unique_product_copy_files_pairs),$(cf)),,\
        $(eval unique_product_copy_files_pairs += $(cf))))
unique_product_copy_files_destinations :=
product_copy_files_ignored :=
$(foreach cf,$(unique_product_copy_files_pairs), \
    $(eval _src := $(call word-colon,1,$(cf))) \
    $(eval _dest := $(call word-colon,2,$(cf))) \
    $(call check-product-copy-files,$(cf),$(_dest)) \
    $(if $(filter $(unique_product_copy_files_destinations),$(_dest)), \
        $(eval product_copy_files_ignored += $(cf)), \
        $(eval _fulldest := $(call append-path,$(PRODUCT_OUT),$(_dest))) \
        $(if $(filter %.xml,$(_dest)),\
            $(eval $(call copy-xml-file-checked,$(_src),$(_fulldest))),\
            $(if $(and $(filter %.jar,$(_dest)),$(filter $(basename $(notdir $(_dest))),$(PRODUCT_LOADED_BY_PRIVILEGED_MODULES))),\
                $(eval $(call copy-and-uncompress-dexs,$(_src),$(_fulldest))), \
                $(eval $(call copy-one-file,$(_src),$(_fulldest))))) \
        $(eval ALL_DEFAULT_INSTALLED_MODULES += $(_fulldest)) \
        $(eval unique_product_copy_files_destinations += $(_dest))))

# Dump a list of overriden (and ignored PRODUCT_COPY_FILES entries)
$(file >$(PRODUCT_OUT)/product_copy_files_ignored.txt,$(subst $(space),$(newline),$(strip $(product_copy_files_ignored))))
ifdef dist_goal
$(file >$(DIST_DIR)/logs/product_copy_files_ignored.txt,$(subst $(space),$(newline),$(strip $(product_copy_files_ignored))))
endif

<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
=======
$(call dist-for-goals,droidcore-unbundled,$(pcf_ignored_file):logs/$(notdir $(pcf_ignored_file)))

pcf_ignored_file :=
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
product_copy_files_ignored :=
unique_product_copy_files_pairs :=
unique_product_copy_files_destinations :=

# -----------------------------------------------------------------
# Returns the max allowed size for an image suitable for hash verification
# (e.g., boot.img, recovery.img, etc).
# The value 69632 derives from MAX_VBMETA_SIZE + MAX_FOOTER_SIZE in $(AVBTOOL).
# $(1): partition size to flash the image
define get-hash-image-max-size
$(if $(1), \
  $(if $(filter true,$(BOARD_AVB_ENABLE)), \
    $(eval _hash_meta_size := 69632), \
    $(eval _hash_meta_size := 0)) \
  $(1)-$(_hash_meta_size))
endef

# -----------------------------------------------------------------
# Define rules to copy headers defined in copy_headers.mk
# If more than one makefile declared a header, print a warning,
# then copy the last one defined. This matches the previous make
# behavior.
$(foreach dest,$(ALL_COPIED_HEADERS), \
    $(eval _srcs := $(ALL_COPIED_HEADERS.$(dest).SRC)) \
    $(eval _src := $(word $(words $(_srcs)),$(_srcs))) \
    $(if $(call streq,$(_src),$(_srcs)),, \
        $(warning Duplicate header copy: $(dest)) \
	$(warning Defined in: $(ALL_COPIED_HEADERS.$(dest).MAKEFILE))) \
    $(eval $(call copy-one-header,$(_src),$(dest))))
all_copied_headers: $(ALL_COPIED_HEADERS)

# -----------------------------------------------------------------
# docs/index.html
ifeq (,$(TARGET_BUILD_APPS))
gen := $(OUT_DOCS)/index.html
ALL_DOCS += $(gen)
$(gen): frameworks/base/docs/docs-redirect-index.html
	@mkdir -p $(dir $@)
	@cp -f $< $@
endif

ndk_doxygen_out := $(OUT_NDK_DOCS)
ndk_headers := $(SOONG_OUT_DIR)/ndk/sysroot/usr/include
ndk_docs_src_dir := frameworks/native/docs
ndk_doxyfile := $(ndk_docs_src_dir)/Doxyfile
ifneq ($(wildcard $(ndk_docs_src_dir)),)
ndk_docs_srcs := $(addprefix $(ndk_docs_src_dir)/,\
    $(call find-files-in-subdirs,$(ndk_docs_src_dir),"*",.))
$(ndk_doxygen_out)/index.html: $(ndk_docs_srcs) $(SOONG_OUT_DIR)/ndk.timestamp
	@mkdir -p $(ndk_doxygen_out)
	@echo "Generating NDK docs to $(ndk_doxygen_out)"
	@( cat $(ndk_doxyfile); \
	    echo "INPUT=$(ndk_headers)"; \
	    echo "HTML_OUTPUT=$(ndk_doxygen_out)" \
	) | doxygen -

# Note: Not a part of the docs target because we don't have doxygen available.
# You can run this target locally if you have doxygen installed.
ndk-docs: $(ndk_doxygen_out)/index.html
.PHONY: ndk-docs
endif

<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
# -----------------------------------------------------------------
# property_overrides_split_enabled
property_overrides_split_enabled :=
ifeq ($(BOARD_PROPERTY_OVERRIDES_SPLIT_ENABLED), true)
  property_overrides_split_enabled := true
=======
ifeq ($(HOST_OS),linux)
$(call dist-for-goals,sdk,$(API_FINGERPRINT))
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
endif

# -----------------------------------------------------------------
# FINAL_VENDOR_DEFAULT_PROPERTIES will be installed in vendor/default.prop if
# property_overrides_split_enabled is true. Otherwise it will be installed in
# ROOT/default.prop.
ifdef BOARD_VNDK_VERSION
  ifeq ($(BOARD_VNDK_VERSION),current)
    FINAL_VENDOR_DEFAULT_PROPERTIES := ro.vndk.version=$(PLATFORM_VNDK_VERSION)
  else
    FINAL_VENDOR_DEFAULT_PROPERTIES := ro.vndk.version=$(BOARD_VNDK_VERSION)
  endif
  ifdef BOARD_VNDK_RUNTIME_DISABLE
    FINAL_VENDOR_DEFAULT_PROPERTIES += ro.vndk.lite=true
  endif
else
  FINAL_VENDOR_DEFAULT_PROPERTIES := ro.vndk.version=$(PLATFORM_VNDK_VERSION)
  FINAL_VENDOR_DEFAULT_PROPERTIES += ro.vndk.lite=true
endif
FINAL_VENDOR_DEFAULT_PROPERTIES += \
    $(call collapse-pairs, $(PRODUCT_DEFAULT_PROPERTY_OVERRIDES))
FINAL_VENDOR_DEFAULT_PROPERTIES := $(call uniq-pairs-by-first-component, \
    $(FINAL_VENDOR_DEFAULT_PROPERTIES),=)

# -----------------------------------------------------------------
# prop.default
ifdef property_overrides_split_enabled
INSTALLED_DEFAULT_PROP_TARGET := $(TARGET_OUT)/etc/prop.default
INSTALLED_DEFAULT_PROP_OLD_TARGET := $(TARGET_ROOT_OUT)/default.prop
ALL_DEFAULT_INSTALLED_MODULES += $(INSTALLED_DEFAULT_PROP_OLD_TARGET)
$(INSTALLED_DEFAULT_PROP_OLD_TARGET): $(INSTALLED_DEFAULT_PROP_TARGET)
else
# legacy path
INSTALLED_DEFAULT_PROP_TARGET := $(TARGET_ROOT_OUT)/default.prop
endif
ALL_DEFAULT_INSTALLED_MODULES += $(INSTALLED_DEFAULT_PROP_TARGET)
FINAL_DEFAULT_PROPERTIES := \
    $(call collapse-pairs, $(ADDITIONAL_DEFAULT_PROPERTIES)) \
    $(call collapse-pairs, $(PRODUCT_SYSTEM_DEFAULT_PROPERTIES))
ifndef property_overrides_split_enabled
  FINAL_DEFAULT_PROPERTIES += \
      $(call collapse-pairs, $(FINAL_VENDOR_DEFAULT_PROPERTIES))
endif
FINAL_DEFAULT_PROPERTIES := $(call uniq-pairs-by-first-component, \
    $(FINAL_DEFAULT_PROPERTIES),=)

intermediate_system_build_prop := $(call intermediates-dir-for,ETC,system_build_prop)/build.prop

$(INSTALLED_DEFAULT_PROP_TARGET): $(intermediate_system_build_prop)
	@echo Target buildinfo: $@
	@mkdir -p $(dir $@)
	@rm -f $@
	$(hide) echo "#" > $@; \
	        echo "# ADDITIONAL_DEFAULT_PROPERTIES" >> $@; \
	        echo "#" >> $@;
	$(hide) $(foreach line,$(FINAL_DEFAULT_PROPERTIES), \
		echo "$(line)" >> $@;)
	$(hide) echo "#" >> $@; \
	        echo "# BOOTIMAGE_BUILD_PROPERTIES" >> $@; \
	        echo "#" >> $@;
	$(hide) echo ro.bootimage.build.date=`$(DATE_FROM_FILE)`>>$@
	$(hide) echo ro.bootimage.build.date.utc=`$(DATE_FROM_FILE) +%s`>>$@
	$(hide) echo ro.bootimage.build.fingerprint="$(BUILD_FINGERPRINT_FROM_FILE)">>$@
	$(hide) build/make/tools/post_process_props.py $@
ifdef property_overrides_split_enabled
	$(hide) mkdir -p $(TARGET_ROOT_OUT)
	$(hide) ln -sf system/etc/prop.default $(INSTALLED_DEFAULT_PROP_OLD_TARGET)
endif

# -----------------------------------------------------------------
# vendor default.prop
INSTALLED_VENDOR_DEFAULT_PROP_TARGET :=
ifdef property_overrides_split_enabled
INSTALLED_VENDOR_DEFAULT_PROP_TARGET := $(TARGET_OUT_VENDOR)/default.prop
ALL_DEFAULT_INSTALLED_MODULES += $(INSTALLED_VENDOR_DEFAULT_PROP_TARGET)

$(INSTALLED_VENDOR_DEFAULT_PROP_TARGET): $(INSTALLED_DEFAULT_PROP_TARGET)
	@echo Target buildinfo: $@
	@mkdir -p $(dir $@)
	$(hide) echo "#" > $@; \
	        echo "# ADDITIONAL VENDOR DEFAULT PROPERTIES" >> $@; \
	        echo "#" >> $@;
	$(hide) $(foreach line,$(FINAL_VENDOR_DEFAULT_PROPERTIES), \
		echo "$(line)" >> $@;)
	$(hide) build/make/tools/post_process_props.py $@

endif  # property_overrides_split_enabled

# -----------------------------------------------------------------
# build.prop
INSTALLED_BUILD_PROP_TARGET := $(TARGET_OUT)/build.prop
ALL_DEFAULT_INSTALLED_MODULES += $(INSTALLED_BUILD_PROP_TARGET)
FINAL_BUILD_PROPERTIES := \
    $(call collapse-pairs, $(ADDITIONAL_BUILD_PROPERTIES))
FINAL_BUILD_PROPERTIES := $(call uniq-pairs-by-first-component, \
    $(FINAL_BUILD_PROPERTIES),=)

# A list of arbitrary tags describing the build configuration.
# Force ":=" so we can use +=
BUILD_VERSION_TAGS := $(BUILD_VERSION_TAGS)
ifeq ($(TARGET_BUILD_TYPE),debug)
  BUILD_VERSION_TAGS += debug
endif
# The "test-keys" tag marks builds signed with the old test keys,
# which are available in the SDK.  "dev-keys" marks builds signed with
# non-default dev keys (usually private keys from a vendor directory).
# Both of these tags will be removed and replaced with "release-keys"
# when the target-files is signed in a post-build step.
ifeq ($(DEFAULT_SYSTEM_DEV_CERTIFICATE),build/target/product/security/testkey)
BUILD_KEYS := test-keys
else
BUILD_KEYS := dev-keys
endif
BUILD_VERSION_TAGS += $(BUILD_KEYS)
BUILD_VERSION_TAGS := $(subst $(space),$(comma),$(sort $(BUILD_VERSION_TAGS)))

# A human-readable string that descibes this build in detail.
build_desc := $(TARGET_PRODUCT)-$(TARGET_BUILD_VARIANT) $(PLATFORM_VERSION) $(BUILD_ID) $(BUILD_NUMBER_FROM_FILE) $(BUILD_VERSION_TAGS)
$(intermediate_system_build_prop): PRIVATE_BUILD_DESC := $(build_desc)

# The string used to uniquely identify the combined build and product; used by the OTA server.
ifeq (,$(strip $(BUILD_FINGERPRINT)))
  ifeq ($(strip $(HAS_BUILD_NUMBER)),false)
    BF_BUILD_NUMBER := $(USER)$$($(DATE_FROM_FILE) +%m%d%H%M)
  else
    BF_BUILD_NUMBER := $(file <$(BUILD_NUMBER_FILE))
  endif
  BUILD_FINGERPRINT := $(PRODUCT_BRAND)/$(TARGET_PRODUCT)/$(TARGET_DEVICE):$(PLATFORM_VERSION)/$(BUILD_ID)/$(BF_BUILD_NUMBER):$(TARGET_BUILD_VARIANT)/$(BUILD_VERSION_TAGS)
endif
# unset it for safety.
BF_BUILD_NUMBER :=

BUILD_FINGERPRINT_FILE := $(PRODUCT_OUT)/build_fingerprint.txt
ifneq (,$(shell mkdir -p $(PRODUCT_OUT) && echo $(BUILD_FINGERPRINT) >$(BUILD_FINGERPRINT_FILE) && grep " " $(BUILD_FINGERPRINT_FILE)))
  $(error BUILD_FINGERPRINT cannot contain spaces: "$(file <$(BUILD_FINGERPRINT_FILE))")
endif
BUILD_FINGERPRINT_FROM_FILE := $$(cat $(BUILD_FINGERPRINT_FILE))
# unset it for safety.
BUILD_FINGERPRINT :=

# The string used to uniquely identify the system build; used by the OTA server.
# This purposefully excludes any product-specific variables.
ifeq (,$(strip $(BUILD_THUMBPRINT)))
  BUILD_THUMBPRINT := $(PLATFORM_VERSION)/$(BUILD_ID)/$(BUILD_NUMBER_FROM_FILE):$(TARGET_BUILD_VARIANT)/$(BUILD_VERSION_TAGS)
endif

BUILD_THUMBPRINT_FILE := $(PRODUCT_OUT)/build_thumbprint.txt
ifneq (,$(shell mkdir -p $(PRODUCT_OUT) && echo $(BUILD_THUMBPRINT) >$(BUILD_THUMBPRINT_FILE) && grep " " $(BUILD_THUMBPRINT_FILE)))
  $(error BUILD_THUMBPRINT cannot contain spaces: "$(file <$(BUILD_THUMBPRINT_FILE))")
endif
BUILD_THUMBPRINT_FROM_FILE := $$(cat $(BUILD_THUMBPRINT_FILE))
# unset it for safety.
BUILD_THUMBPRINT :=

KNOWN_OEM_THUMBPRINT_PROPERTIES := \
    ro.product.brand \
    ro.product.name \
    ro.product.device
OEM_THUMBPRINT_PROPERTIES := $(filter $(KNOWN_OEM_THUMBPRINT_PROPERTIES),\
    $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_OEM_PROPERTIES))

# Display parameters shown under Settings -> About Phone
ifeq ($(TARGET_BUILD_VARIANT),user)
  # User builds should show:
  # release build number or branch.buld_number non-release builds

  # Dev. branches should have DISPLAY_BUILD_NUMBER set
  ifeq (true,$(DISPLAY_BUILD_NUMBER))
    BUILD_DISPLAY_ID := $(BUILD_ID).$(BUILD_NUMBER_FROM_FILE) $(BUILD_KEYS)
  else
    BUILD_DISPLAY_ID := $(BUILD_ID) $(BUILD_KEYS)
  endif
else
  # Non-user builds should show detailed build information
  BUILD_DISPLAY_ID := $(build_desc)
endif

# Accepts a whitespace separated list of product locales such as
# (en_US en_AU en_GB...) and returns the first locale in the list with
# underscores replaced with hyphens. In the example above, this will
# return "en-US".
define get-default-product-locale
$(strip $(subst _,-, $(firstword $(1))))
endef

BUILDINFO_SH := build/make/tools/buildinfo.sh
VENDOR_BUILDINFO_SH := build/make/tools/vendor_buildinfo.sh

# TARGET_BUILD_FLAVOR and ro.build.flavor are used only by the test
# harness to distinguish builds. Only add _asan for a sanitized build
# if it isn't already a part of the flavor (via a dedicated lunch
# config for example).
TARGET_BUILD_FLAVOR := $(TARGET_PRODUCT)-$(TARGET_BUILD_VARIANT)
ifneq (, $(filter address, $(SANITIZE_TARGET)))
ifeq (,$(findstring _asan,$(TARGET_BUILD_FLAVOR)))
TARGET_BUILD_FLAVOR := $(TARGET_BUILD_FLAVOR)_asan
endif
endif

ifdef TARGET_SYSTEM_PROP
system_prop_file := $(TARGET_SYSTEM_PROP)
else
system_prop_file := $(wildcard $(TARGET_DEVICE_DIR)/system.prop)
endif
$(intermediate_system_build_prop): $(BUILDINFO_SH) $(INTERNAL_BUILD_ID_MAKEFILE) $(BUILD_SYSTEM)/version_defaults.mk $(system_prop_file) $(INSTALLED_ANDROID_INFO_TXT_TARGET)
	@echo Target buildinfo: $@
	@mkdir -p $(dir $@)
	$(hide) echo > $@
ifneq ($(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_OEM_PROPERTIES),)
	$(hide) echo "#" >> $@; \
	        echo "# PRODUCT_OEM_PROPERTIES" >> $@; \
	        echo "#" >> $@;
	$(hide) $(foreach prop,$(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_OEM_PROPERTIES), \
		echo "import /oem/oem.prop $(prop)" >> $@;)
endif
	$(hide) TARGET_BUILD_TYPE="$(TARGET_BUILD_VARIANT)" \
			TARGET_BUILD_FLAVOR="$(TARGET_BUILD_FLAVOR)" \
			TARGET_DEVICE="$(TARGET_DEVICE)" \
			PRODUCT_NAME="$(TARGET_PRODUCT)" \
			PRODUCT_BRAND="$(PRODUCT_BRAND)" \
			PRODUCT_DEFAULT_LOCALE="$(call get-default-product-locale,$(PRODUCT_LOCALES))" \
			PRODUCT_DEFAULT_WIFI_CHANNELS="$(PRODUCT_DEFAULT_WIFI_CHANNELS)" \
			PRODUCT_MODEL="$(PRODUCT_MODEL)" \
			PRODUCT_MANUFACTURER="$(PRODUCT_MANUFACTURER)" \
			PRIVATE_BUILD_DESC="$(PRIVATE_BUILD_DESC)" \
			BUILD_ID="$(BUILD_ID)" \
			BUILD_DISPLAY_ID="$(BUILD_DISPLAY_ID)" \
			DATE="$(DATE_FROM_FILE)" \
			BUILD_NUMBER="$(BUILD_NUMBER_FROM_FILE)" \
			BOARD_BUILD_SYSTEM_ROOT_IMAGE="$(BOARD_BUILD_SYSTEM_ROOT_IMAGE)" \
			AB_OTA_UPDATER="$(AB_OTA_UPDATER)" \
			PLATFORM_VERSION="$(PLATFORM_VERSION)" \
			PLATFORM_SECURITY_PATCH="$(PLATFORM_SECURITY_PATCH)" \
			PLATFORM_BASE_OS="$(PLATFORM_BASE_OS)" \
			PLATFORM_SDK_VERSION="$(PLATFORM_SDK_VERSION)" \
			PLATFORM_PREVIEW_SDK_VERSION="$(PLATFORM_PREVIEW_SDK_VERSION)" \
			PLATFORM_VERSION_CODENAME="$(PLATFORM_VERSION_CODENAME)" \
			PLATFORM_VERSION_ALL_CODENAMES="$(PLATFORM_VERSION_ALL_CODENAMES)" \
			PLATFORM_MIN_SUPPORTED_TARGET_SDK_VERSION="$(PLATFORM_MIN_SUPPORTED_TARGET_SDK_VERSION)" \
			BUILD_VERSION_TAGS="$(BUILD_VERSION_TAGS)" \
			BUILD_FINGERPRINT="$(BUILD_FINGERPRINT_FROM_FILE)" \
			$(if $(OEM_THUMBPRINT_PROPERTIES),BUILD_THUMBPRINT="$(BUILD_THUMBPRINT_FROM_FILE)") \
			TARGET_CPU_ABI_LIST="$(TARGET_CPU_ABI_LIST)" \
			TARGET_CPU_ABI_LIST_32_BIT="$(TARGET_CPU_ABI_LIST_32_BIT)" \
			TARGET_CPU_ABI_LIST_64_BIT="$(TARGET_CPU_ABI_LIST_64_BIT)" \
			TARGET_CPU_ABI="$(TARGET_CPU_ABI)" \
			TARGET_CPU_ABI2="$(TARGET_CPU_ABI2)" \
			TARGET_AAPT_CHARACTERISTICS="$(TARGET_AAPT_CHARACTERISTICS)" \
	        bash $(BUILDINFO_SH) >> $@
	$(hide) $(foreach file,$(system_prop_file), \
		if [ -f "$(file)" ]; then \
			echo "#" >> $@; \
			echo Target buildinfo from: "$(file)"; \
			echo "# from $(file)" >> $@; \
			echo "#" >> $@; \
			cat $(file) >> $@; \
		fi;)
	$(if $(FINAL_BUILD_PROPERTIES), \
		$(hide) echo >> $@; \
		        echo "#" >> $@; \
		        echo "# ADDITIONAL_BUILD_PROPERTIES" >> $@; \
		        echo "#" >> $@; )
	$(hide) $(foreach line,$(FINAL_BUILD_PROPERTIES), \
		echo "$(line)" >> $@;)
	$(hide) cat $(INSTALLED_ANDROID_INFO_TXT_TARGET) | grep 'require version-' | sed -e 's/require version-/ro.build.expect./g' >> $@
	$(hide) build/make/tools/post_process_props.py $@ $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_SYSTEM_PROPERTY_BLACKLIST)

build_desc :=

ifeq (,$(filter true, $(TARGET_NO_KERNEL) $(TARGET_NO_RECOVERY)))
INSTALLED_RECOVERYIMAGE_TARGET := $(PRODUCT_OUT)/recovery.img
else
INSTALLED_RECOVERYIMAGE_TARGET :=
endif

$(INSTALLED_BUILD_PROP_TARGET): $(intermediate_system_build_prop) $(INSTALLED_RECOVERYIMAGE_TARGET)
	@echo "Target build info: $@"
	$(hide) grep -v 'ro.product.first_api_level' $(intermediate_system_build_prop) > $@
ifdef INSTALLED_RECOVERYIMAGE_TARGET
	$(hide) echo ro.expect.recovery_id=`cat $(RECOVERYIMAGE_ID_FILE)` >> $@
endif

# -----------------------------------------------------------------
# vendor build.prop
#
# For verifying that the vendor build is what we think it is
INSTALLED_VENDOR_BUILD_PROP_TARGET := $(TARGET_OUT_VENDOR)/build.prop
ALL_DEFAULT_INSTALLED_MODULES += $(INSTALLED_VENDOR_BUILD_PROP_TARGET)

ifdef property_overrides_split_enabled
FINAL_VENDOR_BUILD_PROPERTIES += \
    $(call collapse-pairs, $(PRODUCT_PROPERTY_OVERRIDES))
FINAL_VENDOR_BUILD_PROPERTIES := $(call uniq-pairs-by-first-component, \
    $(FINAL_VENDOR_BUILD_PROPERTIES),=)
endif  # property_overrides_split_enabled

$(INSTALLED_VENDOR_BUILD_PROP_TARGET): $(VENDOR_BUILDINFO_SH) $(intermediate_system_build_prop)
	@echo Target vendor buildinfo: $@
	@mkdir -p $(dir $@)
	$(hide) echo > $@
	$(hide) grep 'ro.product.first_api_level' $(intermediate_system_build_prop) >> $@ || true
	$(hide) echo ro.vendor.build.date=`$(DATE_FROM_FILE)`>>$@
	$(hide) echo ro.vendor.build.date.utc=`$(DATE_FROM_FILE) +%s`>>$@
	$(hide) echo ro.vendor.build.fingerprint="$(BUILD_FINGERPRINT_FROM_FILE)">>$@
	$(hide) echo ro.vendor.build.security_patch="$(VENDOR_SECURITY_PATCH)">>$@
	$(hide) echo ro.vendor.product.cpu.abilist="$(TARGET_CPU_ABI_LIST)">>$@
	$(hide) echo ro.vendor.product.cpu.abilist32="$(TARGET_CPU_ABI_LIST_32_BIT)">>$@
	$(hide) echo ro.vendor.product.cpu.abilist64="$(TARGET_CPU_ABI_LIST_64_BIT)">>$@
	$(hide) TARGET_DEVICE="$(TARGET_DEVICE)" \
			PRODUCT_NAME="$(TARGET_PRODUCT)" \
			PRODUCT_BRAND="$(PRODUCT_BRAND)" \
			PRODUCT_MODEL="$(PRODUCT_MODEL)" \
			PRODUCT_MANUFACTURER="$(PRODUCT_MANUFACTURER)" \
			TARGET_BOOTLOADER_BOARD_NAME="$(TARGET_BOOTLOADER_BOARD_NAME)" \
			TARGET_BOARD_PLATFORM="$(TARGET_BOARD_PLATFORM)" \
	        bash $(VENDOR_BUILDINFO_SH) >> $@
ifdef property_overrides_split_enabled
	$(hide) echo "#" >> $@; \
	        echo "# ADDITIONAL VENDOR BUILD PROPERTIES" >> $@; \
	        echo "#" >> $@;
	$(hide) $(foreach line,$(FINAL_VENDOR_BUILD_PROPERTIES), \
		echo "$(line)" >> $@;)
endif  # property_overrides_split_enabled
	$(hide) build/make/tools/post_process_props.py $@ $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_VENDOR_PROPERTY_BLACKLIST)

# -----------------------------------------------------------------
# product build.prop
INSTALLED_PRODUCT_BUILD_PROP_TARGET := $(TARGET_OUT_PRODUCT)/build.prop
ALL_DEFAULT_INSTALLED_MODULES += $(INSTALLED_PRODUCT_BUILD_PROP_TARGET)

FINAL_PRODUCT_PROPERTIES += \
    $(call collapse-pairs, $(PRODUCT_PRODUCT_PROPERTIES))
FINAL_PRODUCT_PROPERTIES := $(call uniq-pairs-by-first-component, \
    $(FINAL_PRODUCT_PROPERTIES),=)

$(INSTALLED_PRODUCT_BUILD_PROP_TARGET):
	@echo Target product buildinfo: $@
	@mkdir -p $(dir $@)
	$(hide) echo > $@
ifdef BOARD_USES_PRODUCTIMAGE
	$(hide) echo ro.product.build.date=`$(DATE_FROM_FILE)`>>$@
	$(hide) echo ro.product.build.date.utc=`$(DATE_FROM_FILE) +%s`>>$@
	$(hide) echo ro.product.build.fingerprint="$(BUILD_FINGERPRINT_FROM_FILE)">>$@
endif  # BOARD_USES_PRODUCTIMAGE
	$(hide) echo "#" >> $@; \
	        echo "# ADDITIONAL PRODUCT PROPERTIES" >> $@; \
	        echo "#" >> $@;
	$(hide) $(foreach line,$(FINAL_PRODUCT_PROPERTIES), \
		echo "$(line)" >> $@;)
	$(hide) build/make/tools/post_process_props.py $@

# ----------------------------------------------------------------

# -----------------------------------------------------------------
# sdk-build.prop
#
# There are certain things in build.prop that we don't want to
# ship with the sdk; remove them.

# This must be a list of entire property keys followed by
# "=" characters, without any internal spaces.
sdk_build_prop_remove := \
	ro.build.user= \
	ro.build.host= \
	ro.product.brand= \
	ro.product.manufacturer= \
	ro.product.device=
# TODO: Remove this soon-to-be obsolete property
sdk_build_prop_remove += ro.build.product=
INSTALLED_SDK_BUILD_PROP_TARGET := $(PRODUCT_OUT)/sdk/sdk-build.prop
$(INSTALLED_SDK_BUILD_PROP_TARGET): $(INSTALLED_BUILD_PROP_TARGET)
	@echo SDK buildinfo: $@
	@mkdir -p $(dir $@)
	$(hide) grep -v "$(subst $(space),\|,$(strip \
<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
				$(sdk_build_prop_remove)))" $< > $@.tmp
	$(hide) for x in $(sdk_build_prop_remove); do \
				echo "$$x"generic >> $@.tmp; done
=======
	            $(sdk_build_prop_remove)))" $< > $@.tmp
	$(hide) for x in $(strip $(sdk_build_prop_remove)); do \
	            echo "$$x"generic >> $@.tmp; done
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
	$(hide) mv $@.tmp $@

# -----------------------------------------------------------------
<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
# package stats
PACKAGE_STATS_FILE := $(PRODUCT_OUT)/package-stats.txt
PACKAGES_TO_STAT := \
    $(sort $(filter $(TARGET_OUT)/% $(TARGET_OUT_DATA)/%, \
	$(filter %.jar %.apk, $(ALL_DEFAULT_INSTALLED_MODULES))))
$(PACKAGE_STATS_FILE): $(PACKAGES_TO_STAT)
	@echo Package stats: $@
	@mkdir -p $(dir $@)
	$(hide) rm -f $@
ifeq ($(PACKAGES_TO_STAT),)
# Create empty package stats file if target builds no jar(s) or apk(s).
	$(hide) touch $@
=======
# declare recovery ramdisk files
ifeq ($(BUILDING_RECOVERY_IMAGE),true)
INTERNAL_RECOVERY_RAMDISK_FILES_TIMESTAMP := $(call intermediates-dir-for,PACKAGING,recovery)/ramdisk_files-timestamp
endif

# -----------------------------------------------------------------
# Declare vendor ramdisk fragments
INTERNAL_VENDOR_RAMDISK_FRAGMENTS :=

ifeq (true,$(BOARD_INCLUDE_RECOVERY_RAMDISK_IN_VENDOR_BOOT))
  ifneq (,$(filter recovery,$(BOARD_VENDOR_RAMDISK_FRAGMENTS)))
    $(error BOARD_VENDOR_RAMDISK_FRAGMENTS must not contain "recovery" if \
      BOARD_INCLUDE_RECOVERY_RAMDISK_IN_VENDOR_BOOT is set)
  endif
  INTERNAL_VENDOR_RAMDISK_FRAGMENTS += recovery
  VENDOR_RAMDISK_FRAGMENT.recovery.STAGING_DIR := $(TARGET_RECOVERY_ROOT_OUT)
  VENDOR_RAMDISK_FRAGMENT.recovery.FILES := $(INTERNAL_RECOVERY_RAMDISK_FILES_TIMESTAMP)
  BOARD_VENDOR_RAMDISK_FRAGMENT.recovery.MKBOOTIMG_ARGS += --ramdisk_type RECOVERY
  .KATI_READONLY := VENDOR_RAMDISK_FRAGMENT.recovery.STAGING_DIR
endif

# Validation check and assign default --ramdisk_type.
$(foreach vendor_ramdisk_fragment,$(BOARD_VENDOR_RAMDISK_FRAGMENTS), \
  $(if $(and $(BOARD_VENDOR_RAMDISK_FRAGMENT.$(vendor_ramdisk_fragment).KERNEL_MODULE_DIRS), \
             $(BOARD_VENDOR_RAMDISK_FRAGMENT.$(vendor_ramdisk_fragment).PREBUILT)), \
    $(error Must not specify KERNEL_MODULE_DIRS for prebuilt vendor ramdisk fragment "$(vendor_ramdisk_fragment)": $(BOARD_VENDOR_RAMDISK_FRAGMENT.$(vendor_ramdisk_fragment).KERNEL_MODULE_DIRS))) \
  $(eval VENDOR_RAMDISK_FRAGMENT.$(vendor_ramdisk_fragment).STAGING_DIR := $(call intermediates-dir-for,PACKAGING,vendor_ramdisk_fragment-stage-$(vendor_ramdisk_fragment))) \
  $(eval VENDOR_RAMDISK_FRAGMENT.$(vendor_ramdisk_fragment).FILES :=) \
  $(if $(BOARD_VENDOR_RAMDISK_FRAGMENT.$(vendor_ramdisk_fragment).KERNEL_MODULE_DIRS), \
    $(if $(filter --ramdisk_type,$(BOARD_VENDOR_RAMDISK_FRAGMENT.$(vendor_ramdisk_fragment).MKBOOTIMG_ARGS)),, \
      $(eval BOARD_VENDOR_RAMDISK_FRAGMENT.$(vendor_ramdisk_fragment).MKBOOTIMG_ARGS += --ramdisk_type DLKM))) \
)

# Create the "kernel module directory" to "vendor ramdisk fragment" inverse mapping.
$(foreach vendor_ramdisk_fragment,$(BOARD_VENDOR_RAMDISK_FRAGMENTS), \
  $(foreach kmd,$(BOARD_VENDOR_RAMDISK_FRAGMENT.$(vendor_ramdisk_fragment).KERNEL_MODULE_DIRS), \
    $(eval kmd_vrf := KERNEL_MODULE_DIR_VENDOR_RAMDISK_FRAGMENT_$(kmd)) \
    $(if $($(kmd_vrf)),$(error Kernel module directory "$(kmd)" belongs to multiple vendor ramdisk fragments: "$($(kmd_vrf))" "$(vendor_ramdisk_fragment)", each kernel module directory should belong to exactly one or none vendor ramdisk fragment)) \
    $(eval $(kmd_vrf) := $(vendor_ramdisk_fragment)) \
  ) \
)
INTERNAL_VENDOR_RAMDISK_FRAGMENTS += $(BOARD_VENDOR_RAMDISK_FRAGMENTS)

# Strip the list in case of any whitespace.
INTERNAL_VENDOR_RAMDISK_FRAGMENTS := \
  $(strip $(INTERNAL_VENDOR_RAMDISK_FRAGMENTS))

# Assign --ramdisk_name for each vendor ramdisk fragment.
$(foreach vendor_ramdisk_fragment,$(INTERNAL_VENDOR_RAMDISK_FRAGMENTS), \
  $(if $(filter --ramdisk_name,$(BOARD_VENDOR_RAMDISK_FRAGMENT.$(vendor_ramdisk_fragment).MKBOOTIMG_ARGS)), \
    $(error Must not specify --ramdisk_name for vendor ramdisk fragment: $(vendor_ramdisk_fragment))) \
  $(eval BOARD_VENDOR_RAMDISK_FRAGMENT.$(vendor_ramdisk_fragment).MKBOOTIMG_ARGS += --ramdisk_name $(vendor_ramdisk_fragment)) \
  $(eval .KATI_READONLY := BOARD_VENDOR_RAMDISK_FRAGMENT.$(vendor_ramdisk_fragment).MKBOOTIMG_ARGS) \
)

# -----------------------------------------------------------------
# kernel modules

# Depmod requires a well-formed kernel version so 0.0 is used as a placeholder.
DEPMOD_STAGING_SUBDIR :=$= lib/modules/0.0

define copy-and-strip-kernel-module
$(2): $(1)
	$(LLVM_STRIP) -o $(2) --strip-debug $(1)
endef

# $(1): modules list
# $(2): output dir
# $(3): mount point
# $(4): staging dir
# $(5): module load list
# $(6): module load list filename
# $(7): module archive
# $(8): staging dir for stripped modules
# $(9): module directory name
# Returns a list of src:dest pairs to install the modules using copy-many-files.
define build-image-kernel-modules
  $(if $(9), \
    $(eval _dir := $(9)/), \
    $(eval _dir :=)) \
  $(foreach module,$(1), \
    $(eval _src := $(module)) \
    $(if $(8), \
      $(eval _src := $(8)/$(notdir $(module))) \
      $(eval $(call copy-and-strip-kernel-module,$(module),$(_src)))) \
    $(_src):$(2)/lib/modules/$(_dir)$(notdir $(module))) \
  $(eval $(call build-image-kernel-modules-depmod,$(1),$(3),$(4),$(5),$(6),$(7),$(2),$(9))) \
  $(4)/$(DEPMOD_STAGING_SUBDIR)/modules.dep:$(2)/lib/modules/$(_dir)modules.dep \
  $(4)/$(DEPMOD_STAGING_SUBDIR)/modules.alias:$(2)/lib/modules/$(_dir)modules.alias \
  $(4)/$(DEPMOD_STAGING_SUBDIR)/modules.softdep:$(2)/lib/modules/$(_dir)modules.softdep \
  $(4)/$(DEPMOD_STAGING_SUBDIR)/$(6):$(2)/lib/modules/$(_dir)$(6)
endef

# $(1): modules list
# $(2): mount point
# $(3): staging dir
# $(4): module load list
# $(5): module load list filename
# $(6): module archive
# $(7): output dir
# $(8): module directory name
# TODO(b/144844424): If a module archive is being used, this step (which
# generates obj/PACKAGING/.../modules.dep) also unzips the module archive into
# the output directory. This should be moved to a module with a
# LOCAL_POST_INSTALL_CMD so that if modules.dep is removed from the output dir,
# the archive modules are restored along with modules.dep.
define build-image-kernel-modules-depmod
$(3)/$(DEPMOD_STAGING_SUBDIR)/modules.dep: .KATI_IMPLICIT_OUTPUTS := $(3)/$(DEPMOD_STAGING_SUBDIR)/modules.alias $(3)/$(DEPMOD_STAGING_SUBDIR)/modules.softdep $(3)/$(DEPMOD_STAGING_SUBDIR)/$(5)
$(3)/$(DEPMOD_STAGING_SUBDIR)/modules.dep: $(DEPMOD)
$(3)/$(DEPMOD_STAGING_SUBDIR)/modules.dep: PRIVATE_MODULES := $(strip $(1))
$(3)/$(DEPMOD_STAGING_SUBDIR)/modules.dep: PRIVATE_MOUNT_POINT := $(2)
$(3)/$(DEPMOD_STAGING_SUBDIR)/modules.dep: PRIVATE_MODULE_DIR := $(3)/$(DEPMOD_STAGING_SUBDIR)/$(2)/lib/modules/$(8)
$(3)/$(DEPMOD_STAGING_SUBDIR)/modules.dep: PRIVATE_STAGING_DIR := $(3)
$(3)/$(DEPMOD_STAGING_SUBDIR)/modules.dep: PRIVATE_LOAD_MODULES := $(strip $(4))
$(3)/$(DEPMOD_STAGING_SUBDIR)/modules.dep: PRIVATE_LOAD_FILE := $(3)/$(DEPMOD_STAGING_SUBDIR)/$(5)
$(3)/$(DEPMOD_STAGING_SUBDIR)/modules.dep: PRIVATE_MODULE_ARCHIVE := $(6)
$(3)/$(DEPMOD_STAGING_SUBDIR)/modules.dep: PRIVATE_OUTPUT_DIR := $(7)
$(3)/$(DEPMOD_STAGING_SUBDIR)/modules.dep: $(1) $(6)
	@echo depmod $$(PRIVATE_STAGING_DIR)
	rm -rf $$(PRIVATE_STAGING_DIR)
	mkdir -p $$(PRIVATE_MODULE_DIR)
	$(if $(6),\
	  unzip -qoDD -d $$(PRIVATE_MODULE_DIR) $$(PRIVATE_MODULE_ARCHIVE); \
	  mkdir -p $$(PRIVATE_OUTPUT_DIR)/lib; \
	  cp -r  $(3)/$(DEPMOD_STAGING_SUBDIR)/$(2)/lib/modules $$(PRIVATE_OUTPUT_DIR)/lib/; \
	  find $$(PRIVATE_MODULE_DIR) -type f -name *.ko | xargs basename -a > $$(PRIVATE_LOAD_FILE); \
	)
	$(if $(1),\
	  cp $$(PRIVATE_MODULES) $$(PRIVATE_MODULE_DIR)/; \
	  for MODULE in $$(PRIVATE_LOAD_MODULES); do \
	    basename $$$$MODULE >> $$(PRIVATE_LOAD_FILE); \
	  done; \
	)
	$(DEPMOD) -b $$(PRIVATE_STAGING_DIR) 0.0
	# Turn paths in modules.dep into absolute paths
	sed -i.tmp -e 's|\([^: ]*lib/modules/[^: ]*\)|/\1|g' $$(PRIVATE_STAGING_DIR)/$$(DEPMOD_STAGING_SUBDIR)/modules.dep
	touch $$(PRIVATE_LOAD_FILE)
endef

# $(1): staging dir
# $(2): modules list
# $(3): module load list
# $(4): module load list filename
# $(5): output dir
define module-load-list-copy-paths
  $(eval $(call build-image-module-load-list,$(1),$(2),$(3),$(4))) \
  $(1)/$(DEPMOD_STAGING_SUBDIR)/$(4):$(5)/lib/modules/$(4)
endef

# $(1): staging dir
# $(2): modules list
# $(3): module load list
# $(4): module load list filename
define build-image-module-load-list
$(1)/$(DEPMOD_STAGING_SUBDIR)/$(4): PRIVATE_LOAD_MODULES := $(3)
$(1)/$(DEPMOD_STAGING_SUBDIR)/$(4): $(2)
	@echo load-list $$(@)
	@echo '$$(strip $$(notdir $$(PRIVATE_LOAD_MODULES)))' | tr ' ' '\n' > $$(@)
endef

# $(1): source options file
# $(2): destination pathname
# Returns a build rule that checks the syntax of and installs a kernel modules
# options file. Strip and squeeze any extra space and blank lines.
# For use via $(eval).
define build-image-kernel-modules-options-file
$(2): $(1)
	@echo "libmodprobe options $$(@)"
	$(hide) mkdir -p "$$(dir $$@)"
	$(hide) rm -f "$$@"
	$(hide) awk <"$$<" >"$$@" \
	  '/^#/ { print; next } \
	   NF == 0 { next } \
	   NF < 2 || $$$$1 != "options" \
	     { print "Invalid options line " FNR ": " $$$$0 >"/dev/stderr"; \
	       exit_status = 1; next } \
	   { $$$$1 = $$$$1; print } \
	   END { exit exit_status }'
endef

# $(1): source blocklist file
# $(2): destination pathname
# Returns a build rule that checks the syntax of and installs a kernel modules
# blocklist file. Strip and squeeze any extra space and blank lines.
# For use via $(eval).
define build-image-kernel-modules-blocklist-file
$(2): $(1)
	@echo "libmodprobe blocklist $$(@)"
	$(hide) mkdir -p "$$(dir $$@)"
	$(hide) rm -f "$$@"
	$(hide) awk <"$$<" >"$$@" \
	  '/^#/ { print; next } \
	   NF == 0 { next } \
	   NF != 2 || $$$$1 != "blocklist" \
	     { print "Invalid blocklist line " FNR ": " $$$$0 >"/dev/stderr"; \
	       exit_status = 1; next } \
	   { $$$$1 = $$$$1; print } \
	   END { exit exit_status }'
endef

# $(1): image name
# $(2): build output directory (TARGET_OUT_VENDOR, TARGET_RECOVERY_ROOT_OUT, etc)
# $(3): mount point
# $(4): module load filename
# $(5): stripped staging directory
# $(6): kernel module directory name (top is an out of band value for no directory)
define build-image-kernel-modules-dir
$(if $(filter top,$(6)),\
  $(eval _kver :=)$(eval _sep :=),\
  $(eval _kver := $(6))$(eval _sep :=_))\
$(if $(5),\
  $(eval _stripped_staging_dir := $(5)$(_sep)$(_kver)),\
  $(eval _stripped_staging_dir :=))\
$(if $(strip $(BOARD_$(1)_KERNEL_MODULES$(_sep)$(_kver))$(BOARD_$(1)_KERNEL_MODULES_ARCHIVE$(_sep)$(_kver))),\
  $(if $(BOARD_$(1)_KERNEL_MODULES_LOAD$(_sep)$(_kver)),,\
    $(eval BOARD_$(1)_KERNEL_MODULES_LOAD$(_sep)$(_kver) := $(BOARD_$(1)_KERNEL_MODULES$(_sep)$(_kver)))) \
  $(call copy-many-files,$(call build-image-kernel-modules,$(BOARD_$(1)_KERNEL_MODULES$(_sep)$(_kver)),$(2),$(3),$(call intermediates-dir-for,PACKAGING,depmod_$(1)$(_sep)$(_kver)),$(BOARD_$(1)_KERNEL_MODULES_LOAD$(_sep)$(_kver)),$(4),$(BOARD_$(1)_KERNEL_MODULES_ARCHIVE$(_sep)$(_kver)),$(_stripped_staging_dir),$(_kver)))) \
$(if $(_kver), \
  $(eval _dir := $(_kver)/), \
  $(eval _dir :=)) \
$(if $(BOARD_$(1)_KERNEL_MODULES_OPTIONS_FILE$(_sep)$(_kver)), \
  $(eval $(call build-image-kernel-modules-options-file, \
    $(BOARD_$(1)_KERNEL_MODULES_OPTIONS_FILE$(_sep)$(_kver)), \
    $(2)/lib/modules/$(_dir)modules.options)) \
  $(2)/lib/modules/$(_dir)modules.options) \
$(if $(BOARD_$(1)_KERNEL_MODULES_BLOCKLIST_FILE$(_sep)$(_kver)), \
  $(eval $(call build-image-kernel-modules-blocklist-file, \
    $(BOARD_$(1)_KERNEL_MODULES_BLOCKLIST_FILE$(_sep)$(_kver)), \
    $(2)/lib/modules/$(_dir)modules.blocklist)) \
  $(2)/lib/modules/$(_dir)modules.blocklist)
endef

# $(1): kernel module directory name (top is an out of band value for no directory)
define build-recovery-as-boot-load
$(if $(filter top,$(1)),\
  $(eval _kver :=)$(eval _sep :=),\
  $(eval _kver := $(1))$(eval _sep :=_))\
  $(if $(BOARD_GENERIC_RAMDISK_KERNEL_MODULES_LOAD$(_sep)$(_kver)),\
    $(call copy-many-files,$(call module-load-list-copy-paths,$(call intermediates-dir-for,PACKAGING,ramdisk_module_list$(_sep)$(_kver)),$(BOARD_GENERIC_RAMDISK_KERNEL_MODULES$(_sep)$(_kver)),$(BOARD_GENERIC_RAMDISK_KERNEL_MODULES_LOAD$(_sep)$(_kver)),modules.load,$(TARGET_RECOVERY_ROOT_OUT))))
endef

# $(1): kernel module directory name (top is an out of band value for no directory)
define build-vendor-ramdisk-recovery-load
$(if $(filter top,$(1)),\
  $(eval _kver :=)$(eval _sep :=),\
  $(eval _kver := $(1))$(eval _sep :=_))\
  $(if $(BOARD_VENDOR_RAMDISK_RECOVERY_KERNEL_MODULES_LOAD$(_sep)$(_kver)),\
    $(call copy-many-files,$(call module-load-list-copy-paths,$(call intermediates-dir-for,PACKAGING,vendor_ramdisk_recovery_module_list$(_sep)$(_kver)),$(BOARD_VENDOR_RAMDISK_KERNEL_MODULES$(_sep)$(_kver)),$(BOARD_VENDOR_RAMDISK_RECOVERY_KERNEL_MODULES_LOAD$(_sep)$(_kver)),modules.load.recovery,$(TARGET_VENDOR_RAMDISK_OUT))))
endef

# $(1): kernel module directory name (top is an out of band value for no directory)
define build-vendor-charger-load
$(if $(filter top,$(1)),\
  $(eval _kver :=)$(eval _sep :=),\
  $(eval _kver := $(1))$(eval _sep :=_))\
  $(if $(BOARD_VENDOR_CHARGER_KERNEL_MODULES_LOAD$(_sep)$(_kver)),\
    $(call copy-many-files,$(call module-load-list-copy-paths,$(call intermediates-dir-for,PACKAGING,vendor_charger_module_list$(_sep)$(_kver)),$(BOARD_VENDOR_CHARGER_KERNEL_MODULES$(_sep)$(_kver)),$(BOARD_VENDOR_CHARGER_KERNEL_MODULES_LOAD$(_sep)$(_kver)),modules.load.charger,$(TARGET_OUT_VENDOR))))
endef

ifneq ($(BUILDING_VENDOR_BOOT_IMAGE),true)
  # If there is no vendor boot partition, store vendor ramdisk kernel modules in the
  # boot ramdisk.
  BOARD_GENERIC_RAMDISK_KERNEL_MODULES += $(BOARD_VENDOR_RAMDISK_KERNEL_MODULES)
  BOARD_GENERIC_RAMDISK_KERNEL_MODULES_LOAD += $(BOARD_VENDOR_RAMDISK_KERNEL_MODULES_LOAD)
endif

ifeq ($(BOARD_GENERIC_RAMDISK_KERNEL_MODULES_LOAD),)
  BOARD_GENERIC_RAMDISK_KERNEL_MODULES_LOAD := $(BOARD_GENERIC_RAMDISK_KERNEL_MODULES)
endif

ifneq ($(strip $(BOARD_GENERIC_RAMDISK_KERNEL_MODULES)),)
  ifeq ($(BOARD_USES_RECOVERY_AS_BOOT), true)
    BOARD_RECOVERY_KERNEL_MODULES += $(BOARD_GENERIC_RAMDISK_KERNEL_MODULES)
  endif
endif

ifneq ($(BOARD_DO_NOT_STRIP_RECOVERY_MODULES),true)
	RECOVERY_STRIPPED_MODULE_STAGING_DIR := $(call intermediates-dir-for,PACKAGING,depmod_recovery_stripped)
else
	RECOVERY_STRIPPED_MODULE_STAGING_DIR :=
endif

ifneq ($(BOARD_DO_NOT_STRIP_VENDOR_MODULES),true)
	VENDOR_STRIPPED_MODULE_STAGING_DIR := $(call intermediates-dir-for,PACKAGING,depmod_vendor_stripped)
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
else
	$(hide) build/make/tools/dump-package-stats $^ > $@
endif

<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
.PHONY: package-stats
package-stats: $(PACKAGE_STATS_FILE)
=======
ifneq ($(BOARD_DO_NOT_STRIP_VENDOR_RAMDISK_MODULES),true)
  VENDOR_RAMDISK_STRIPPED_MODULE_STAGING_DIR := $(call intermediates-dir-for,PACKAGING,depmod_vendor_ramdisk_stripped)
else
  VENDOR_RAMDISK_STRIPPED_MODULE_STAGING_DIR :=
endif

BOARD_KERNEL_MODULE_DIRS += top
$(foreach kmd,$(BOARD_KERNEL_MODULE_DIRS), \
  $(eval ALL_DEFAULT_INSTALLED_MODULES += $(call build-image-kernel-modules-dir,RECOVERY,$(TARGET_RECOVERY_ROOT_OUT),,modules.load.recovery,$(RECOVERY_STRIPPED_MODULE_STAGING_DIR),$(kmd))) \
  $(eval vendor_ramdisk_fragment := $(KERNEL_MODULE_DIR_VENDOR_RAMDISK_FRAGMENT_$(kmd))) \
  $(if $(vendor_ramdisk_fragment), \
    $(eval output_dir := $(VENDOR_RAMDISK_FRAGMENT.$(vendor_ramdisk_fragment).STAGING_DIR)) \
    $(eval result_var := VENDOR_RAMDISK_FRAGMENT.$(vendor_ramdisk_fragment).FILES) \
    $(eval ### else ###), \
    $(eval output_dir := $(TARGET_VENDOR_RAMDISK_OUT)) \
    $(eval result_var := ALL_DEFAULT_INSTALLED_MODULES)) \
  $(eval $(result_var) += $(call build-image-kernel-modules-dir,VENDOR_RAMDISK,$(output_dir),,modules.load,$(VENDOR_RAMDISK_STRIPPED_MODULE_STAGING_DIR),$(kmd))) \
  $(eval ALL_DEFAULT_INSTALLED_MODULES += $(call build-vendor-ramdisk-recovery-load,$(kmd))) \
  $(eval ALL_DEFAULT_INSTALLED_MODULES += $(call build-image-kernel-modules-dir,VENDOR,$(if $(filter true,$(BOARD_USES_VENDOR_DLKMIMAGE)),$(TARGET_OUT_VENDOR_DLKM),$(TARGET_OUT_VENDOR)),vendor,modules.load,$(VENDOR_STRIPPED_MODULE_STAGING_DIR),$(kmd))) \
  $(eval ALL_DEFAULT_INSTALLED_MODULES += $(call build-vendor-charger-load,$(kmd))) \
  $(eval ALL_DEFAULT_INSTALLED_MODULES += $(call build-image-kernel-modules-dir,ODM,$(if $(filter true,$(BOARD_USES_ODM_DLKMIMAGE)),$(TARGET_OUT_ODM_DLKM),$(TARGET_OUT_ODM)),odm,modules.load,,$(kmd))) \
  $(if $(filter true,$(BOARD_USES_RECOVERY_AS_BOOT)),\
    $(eval ALL_DEFAULT_INSTALLED_MODULES += $(call build-recovery-as-boot-load,$(kmd))),\
    $(eval ALL_DEFAULT_INSTALLED_MODULES += $(call build-image-kernel-modules-dir,GENERIC_RAMDISK,$(TARGET_RAMDISK_OUT),,modules.load,,$(kmd)))))
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)

# -----------------------------------------------------------------
# Cert-to-package mapping.  Used by the post-build signing tools.
# Use a macro to add newline to each echo command
define _apkcerts_write_line
$(hide) echo -n 'name="$(1).apk" certificate="$2" private_key="$3"' >> $5
$(if $(4), $(hide) echo -n ' compressed="$4"' >> $5)
$(hide) echo '' >> $5

endef

name := $(TARGET_PRODUCT)
ifeq ($(TARGET_BUILD_TYPE),debug)
  name := $(name)_debug
endif
name := $(name)-apkcerts-$(FILE_NAME_TAG)
intermediates := \
	$(call intermediates-dir-for,PACKAGING,apkcerts)
APKCERTS_FILE := $(intermediates)/$(name).txt
# We don't need to really build all the modules.
# TODO: rebuild APKCERTS_FILE if any app change its cert.
$(APKCERTS_FILE):
	@echo APK certs list: $@
	@mkdir -p $(dir $@)
	@rm -f $@
	$(foreach p,$(PACKAGES),\
	  $(if $(PACKAGES.$(p).EXTERNAL_KEY),\
	    $(call _apkcerts_write_line,$(p),"EXTERNAL","",$(PACKAGES.$(p).COMPRESSED),$@),\
	    $(call _apkcerts_write_line,$(p),$(PACKAGES.$(p).CERTIFICATE),$(PACKAGES.$(p).PRIVATE_KEY),$(PACKAGES.$(p).COMPRESSED),$@)))
	# In case value of PACKAGES is empty.
	$(hide) touch $@

.PHONY: apkcerts-list
apkcerts-list: $(APKCERTS_FILE)

ifneq (,$(TARGET_BUILD_APPS))
  $(call dist-for-goals, apps_only, $(APKCERTS_FILE):apkcerts.txt)
endif

# -----------------------------------------------------------------
# build system stats
BUILD_SYSTEM_STATS := $(PRODUCT_OUT)/build_system_stats.txt
$(BUILD_SYSTEM_STATS):
	@rm -f $@
	@$(foreach s,$(STATS.MODULE_TYPE),echo "modules_type_make,$(s),$(words $(STATS.MODULE_TYPE.$(s)))" >>$@;)
	@$(foreach s,$(STATS.SOONG_MODULE_TYPE),echo "modules_type_soong,$(s),$(STATS.SOONG_MODULE_TYPE.$(s))" >>$@;)
$(call dist-for-goals,droidcore-unbundled,$(BUILD_SYSTEM_STATS))

# -----------------------------------------------------------------
# Modules ready to be converted to Soong, ordered by how many
# modules depend on them.
SOONG_CONV := $(sort $(SOONG_CONV))
SOONG_CONV_DATA := $(call intermediates-dir-for,PACKAGING,soong_conversion)/soong_conv_data
$(SOONG_CONV_DATA):
	@rm -f $@
	@$(foreach s,$(SOONG_CONV),echo "$(s),$(SOONG_CONV.$(s).TYPE),$(sort $(SOONG_CONV.$(s).PROBLEMS)),$(sort $(filter-out $(SOONG_ALREADY_CONV),$(SOONG_CONV.$(s).DEPS)))" >>$@;)

SOONG_TO_CONVERT_SCRIPT := build/make/tools/soong_to_convert.py
SOONG_TO_CONVERT := $(PRODUCT_OUT)/soong_to_convert.txt
$(SOONG_TO_CONVERT): $(SOONG_CONV_DATA) $(SOONG_TO_CONVERT_SCRIPT)
	@rm -f $@
	$(hide) $(SOONG_TO_CONVERT_SCRIPT) $< >$@
$(call dist-for-goals,droidcore-unbundled,$(SOONG_TO_CONVERT))

<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
=======
MK2BP_CATALOG_SCRIPT := build/make/tools/mk2bp_catalog.py
MK2BP_REMAINING_HTML := $(PRODUCT_OUT)/mk2bp_remaining.html
$(MK2BP_REMAINING_HTML): PRIVATE_CODE_SEARCH_BASE_URL := "https://cs.android.com/android/platform/superproject/+/master:"
$(MK2BP_REMAINING_HTML): $(SOONG_CONV_DATA) $(MK2BP_CATALOG_SCRIPT)
	@rm -f $@
	$(hide) $(MK2BP_CATALOG_SCRIPT) \
		--device=$(TARGET_DEVICE) \
		--title="Remaining Android.mk files for $(TARGET_DEVICE)-$(TARGET_BUILD_VARIANT)" \
		--codesearch=$(PRIVATE_CODE_SEARCH_BASE_URL) \
		--out_dir="$(OUT_DIR)" \
		--mode=html \
		> $@
$(call dist-for-goals,droidcore-unbundled,$(MK2BP_REMAINING_HTML))

MK2BP_REMAINING_CSV := $(PRODUCT_OUT)/mk2bp_remaining.csv
$(MK2BP_REMAINING_CSV): $(SOONG_CONV_DATA) $(MK2BP_CATALOG_SCRIPT)
	@rm -f $@
	$(hide) $(MK2BP_CATALOG_SCRIPT) \
		--device=$(TARGET_DEVICE) \
		--out_dir="$(OUT_DIR)" \
		--mode=csv \
		> $@
$(call dist-for-goals,droidcore-unbundled,$(MK2BP_REMAINING_CSV))

>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
# -----------------------------------------------------------------
# Modules use -Wno-error, or added default -Wall -Werror
WALL_WERROR := $(PRODUCT_OUT)/wall_werror.txt
$(WALL_WERROR):
	@rm -f $@
	echo "# Modules using -Wno-error" >> $@
	for m in $(sort $(SOONG_MODULES_USING_WNO_ERROR) $(MODULES_USING_WNO_ERROR)); do echo $$m >> $@; done
	echo "# Modules added default -Wall" >> $@
	for m in $(sort $(SOONG_MODULES_ADDED_WALL) $(MODULES_ADDED_WALL)); do echo $$m >> $@; done

$(call dist-for-goals,droidcore-unbundled,$(WALL_WERROR))

# -----------------------------------------------------------------
<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
=======
# C/C++ flag information for modules
$(call dist-for-goals,droidcore-unbundled,$(SOONG_MODULES_CFLAG_ARTIFACTS))

# -----------------------------------------------------------------
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
# Modules missing profile files
PGO_PROFILE_MISSING := $(PRODUCT_OUT)/pgo_profile_file_missing.txt
$(PGO_PROFILE_MISSING):
	@rm -f $@
	echo "# Modules missing PGO profile files" >> $@
	for m in $(SOONG_MODULES_MISSING_PGO_PROFILE_FILE); do echo $$m >> $@; done

$(call dist-for-goals,droidcore,$(PGO_PROFILE_MISSING))

# -----------------------------------------------------------------
# The dev key is used to sign this package, and as the key required
# for future OTA packages installed by this system.  Actual product
# deliverables will be re-signed by hand.  We expect this file to
# exist with the suffixes ".x509.pem" and ".pk8".
DEFAULT_KEY_CERT_PAIR := $(DEFAULT_SYSTEM_DEV_CERTIFICATE)


# Rules that need to be present for the all targets, even
# if they don't do anything.
.PHONY: systemimage
systemimage:

# -----------------------------------------------------------------

.PHONY: event-log-tags

# Produce an event logs tag file for everything we know about, in order
# to properly allocate numbers.  Then produce a file that's filtered
# for what's going to be installed.

all_event_log_tags_file := $(TARGET_OUT_COMMON_INTERMEDIATES)/all-event-log-tags.txt

event_log_tags_file := $(TARGET_OUT)/etc/event-log-tags

# Include tags from all packages that we know about
all_event_log_tags_src := \
    $(sort $(foreach m, $(ALL_MODULES), $(ALL_MODULES.$(m).EVENT_LOG_TAGS)))

# PDK builds will already have a full list of tags that needs to get merged
# in with the ones from source
pdk_fusion_log_tags_file := $(patsubst $(PRODUCT_OUT)/%,$(_pdk_fusion_intermediates)/%,$(filter $(event_log_tags_file),$(ALL_PDK_FUSION_FILES)))

$(all_event_log_tags_file): PRIVATE_SRC_FILES := $(all_event_log_tags_src) $(pdk_fusion_log_tags_file)
$(all_event_log_tags_file): $(all_event_log_tags_src) $(pdk_fusion_log_tags_file) $(MERGETAGS) build/make/tools/event_log_tags.py
	$(hide) mkdir -p $(dir $@)
	$(hide) $(MERGETAGS) -o $@ $(PRIVATE_SRC_FILES)

# Include tags from all packages included in this product, plus all
# tags that are part of the system (ie, not in a vendor/ or device/
# directory).
event_log_tags_src := \
    $(sort $(foreach m,\
      $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_PACKAGES) \
      $(call module-names-for-tag-list,user), \
      $(ALL_MODULES.$(m).EVENT_LOG_TAGS)) \
      $(filter-out vendor/% device/% out/%,$(all_event_log_tags_src)))

$(event_log_tags_file): PRIVATE_SRC_FILES := $(event_log_tags_src) $(pdk_fusion_log_tags_file)
$(event_log_tags_file): PRIVATE_MERGED_FILE := $(all_event_log_tags_file)
$(event_log_tags_file): $(event_log_tags_src) $(all_event_log_tags_file) $(pdk_fusion_log_tags_file) $(MERGETAGS) build/make/tools/event_log_tags.py
	$(hide) mkdir -p $(dir $@)
	$(hide) $(MERGETAGS) -o $@ -m $(PRIVATE_MERGED_FILE) $(PRIVATE_SRC_FILES)

event-log-tags: $(event_log_tags_file)

ALL_DEFAULT_INSTALLED_MODULES += $(event_log_tags_file)


# #################################################################
# Targets for boot/OS images
# #################################################################
ifneq ($(strip $(TARGET_NO_BOOTLOADER)),true)
  INSTALLED_BOOTLOADER_MODULE := $(PRODUCT_OUT)/bootloader
  ifeq ($(strip $(TARGET_BOOTLOADER_IS_2ND)),true)
    INSTALLED_2NDBOOTLOADER_TARGET := $(PRODUCT_OUT)/2ndbootloader
  else
    INSTALLED_2NDBOOTLOADER_TARGET :=
  endif
else
  INSTALLED_BOOTLOADER_MODULE :=
  INSTALLED_2NDBOOTLOADER_TARGET :=
endif # TARGET_NO_BOOTLOADER
ifneq ($(strip $(TARGET_NO_KERNEL)),true)
  INSTALLED_KERNEL_TARGET := $(PRODUCT_OUT)/kernel
else
  INSTALLED_KERNEL_TARGET :=
endif

# -----------------------------------------------------------------
<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
=======
# the root dir
INTERNAL_ROOT_FILES := $(filter $(TARGET_ROOT_OUT)/%, \
	$(ALL_GENERATED_SOURCES) \
	$(ALL_DEFAULT_INSTALLED_MODULES))

INSTALLED_FILES_FILE_ROOT := $(PRODUCT_OUT)/installed-files-root.txt
INSTALLED_FILES_JSON_ROOT := $(INSTALLED_FILES_FILE_ROOT:.txt=.json)
$(INSTALLED_FILES_FILE_ROOT): .KATI_IMPLICIT_OUTPUTS := $(INSTALLED_FILES_JSON_ROOT)
$(INSTALLED_FILES_FILE_ROOT) : $(INTERNAL_ROOT_FILES) $(FILESLIST) $(FILESLIST_UTIL)
	@echo Installed file list: $@
	mkdir -p $(dir $@)
	rm -f $@
	$(FILESLIST) $(TARGET_ROOT_OUT) > $(@:.txt=.json)
	$(FILESLIST_UTIL) -c $(@:.txt=.json) > $@

ifeq ($(HOST_OS),linux)
$(call dist-for-goals, sdk sdk_addon, $(INSTALLED_FILES_FILE_ROOT))
endif

#------------------------------------------------------------------
# dtb
ifdef BOARD_INCLUDE_DTB_IN_BOOTIMG
INSTALLED_DTBIMAGE_TARGET := $(PRODUCT_OUT)/dtb.img
ifdef BOARD_PREBUILT_DTBIMAGE_DIR
$(INSTALLED_DTBIMAGE_TARGET) : $(sort $(wildcard $(BOARD_PREBUILT_DTBIMAGE_DIR)/*.dtb))
	cat $^ > $@
endif
endif

# -----------------------------------------------------------------
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
# the ramdisk
INTERNAL_RAMDISK_FILES := $(filter $(TARGET_ROOT_OUT)/%, \
	$(ALL_GENERATED_SOURCES) \
	$(ALL_DEFAULT_INSTALLED_MODULES))

<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
=======
INSTALLED_FILES_FILE_RAMDISK := $(PRODUCT_OUT)/installed-files-ramdisk.txt
INSTALLED_FILES_JSON_RAMDISK := $(INSTALLED_FILES_FILE_RAMDISK:.txt=.json)
$(INSTALLED_FILES_FILE_RAMDISK): .KATI_IMPLICIT_OUTPUTS := $(INSTALLED_FILES_JSON_RAMDISK)
$(INSTALLED_FILES_FILE_RAMDISK) : $(INTERNAL_RAMDISK_FILES) $(FILESLIST) $(FILESLIST_UTIL)
	@echo Installed file list: $@
	mkdir -p $(TARGET_RAMDISK_OUT)
	mkdir -p $(dir $@)
	rm -f $@
	$(FILESLIST) $(TARGET_RAMDISK_OUT) > $(@:.txt=.json)
	$(FILESLIST_UTIL) -c $(@:.txt=.json) > $@

ifeq ($(HOST_OS),linux)
$(call dist-for-goals, sdk win_sdk sdk_addon, $(INSTALLED_FILES_FILE_RAMDISK))
endif
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
BUILT_RAMDISK_TARGET := $(PRODUCT_OUT)/ramdisk.img

# We just build this directly to the install location.
INSTALLED_RAMDISK_TARGET := $(BUILT_RAMDISK_TARGET)
<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
$(INSTALLED_RAMDISK_TARGET): $(MKBOOTFS) $(INTERNAL_RAMDISK_FILES) | $(MINIGZIP)
	$(call pretty,"Target ram disk: $@")
	$(hide) $(MKBOOTFS) -d $(TARGET_OUT) $(TARGET_ROOT_OUT) | $(MINIGZIP) > $@
=======
$(INSTALLED_RAMDISK_TARGET): PRIVATE_DIRS := debug_ramdisk dev metadata mnt proc second_stage_resources sys
$(INSTALLED_RAMDISK_TARGET): $(MKBOOTFS) $(INTERNAL_RAMDISK_FILES) $(INSTALLED_FILES_FILE_RAMDISK) | $(COMPRESSION_COMMAND_DEPS)
	$(call pretty,"Target ramdisk: $@")
	$(hide) mkdir -p $(addprefix $(TARGET_RAMDISK_OUT)/,$(PRIVATE_DIRS))
ifeq (true,$(BOARD_USES_GENERIC_KERNEL_IMAGE))
	$(hide) mkdir -p $(addprefix $(TARGET_RAMDISK_OUT)/first_stage_ramdisk/,$(PRIVATE_DIRS))
endif
	$(hide) $(MKBOOTFS) -d $(TARGET_OUT) $(TARGET_RAMDISK_OUT) | $(COMPRESSION_COMMAND) > $@
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)

.PHONY: ramdisk-nodeps
ramdisk-nodeps: $(MKBOOTFS) | $(MINIGZIP)
	@echo "make $@: ignoring dependencies"
	$(hide) $(MKBOOTFS) -d $(TARGET_OUT) $(TARGET_ROOT_OUT) | $(MINIGZIP) > $(INSTALLED_RAMDISK_TARGET)

INSTALLED_BOOTIMAGE_TARGET := $(PRODUCT_OUT)/boot.img

ifneq ($(strip $(TARGET_NO_KERNEL)),true)

# -----------------------------------------------------------------
# the boot image, which is a collection of other images.
<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
=======

# This is defined here since we may be building recovery as boot
# below and only want to define this once
ifneq ($(strip $(BOARD_KERNEL_BINARIES)),)
  BUILT_BOOTIMAGE_TARGET := $(foreach k,$(subst kernel,boot,$(BOARD_KERNEL_BINARIES)), $(PRODUCT_OUT)/$(k).img)
else
  BUILT_BOOTIMAGE_TARGET := $(PRODUCT_OUT)/boot.img
endif

INTERNAL_PREBUILT_BOOTIMAGE :=

my_installed_prebuilt_gki_apex := $(strip $(foreach package,$(PRODUCT_PACKAGES),$(if $(ALL_MODULES.$(package).EXTRACTED_BOOT_IMAGE),$(package))))
ifdef my_installed_prebuilt_gki_apex
  ifneq (1,$(words $(my_installed_prebuilt_gki_apex))) # len(my_installed_prebuilt_gki_apex) > 1
    $(error More than one prebuilt GKI APEXes are installed: $(my_installed_prebuilt_gki_apex))
  endif # len(my_installed_prebuilt_gki_apex) > 1

  ifdef BOARD_PREBUILT_BOOTIMAGE
    $(error Must not define BOARD_PREBUILT_BOOTIMAGE because a prebuilt GKI APEX is installed: $(my_installed_prebuilt_gki_apex))
  endif # BOARD_PREBUILT_BOOTIMAGE defined

  my_apex_extracted_boot_image := $(ALL_MODULES.$(my_installed_prebuilt_gki_apex).EXTRACTED_BOOT_IMAGE)
  INSTALLED_BOOTIMAGE_TARGET := $(PRODUCT_OUT)/boot.img
  $(eval $(call copy-one-file,$(my_apex_extracted_boot_image),$(INSTALLED_BOOTIMAGE_TARGET)))

  INTERNAL_PREBUILT_BOOTIMAGE := $(my_apex_extracted_boot_image)

else # my_installed_prebuilt_gki_apex not defined

# $1: boot image target
# returns the kernel used to make the bootimage
define bootimage-to-kernel
  $(if $(BOARD_KERNEL_BINARIES),\
    $(PRODUCT_OUT)/$(subst .img,,$(subst boot,kernel,$(notdir $(1)))),\
    $(INSTALLED_KERNEL_TARGET))
endef

ifdef BOARD_BOOTIMAGE_PARTITION_SIZE
  BOARD_KERNEL_BOOTIMAGE_PARTITION_SIZE := $(BOARD_BOOTIMAGE_PARTITION_SIZE)
endif

# $1: boot image file name
# $2: boot image variant (boot, boot-debug, boot-test-harness)
define get-bootimage-partition-size
  $(BOARD_$(call to-upper,$(subst .img,,$(subst $(2),kernel,$(notdir $(1)))))_BOOTIMAGE_PARTITION_SIZE)
endef

ifneq ($(strip $(TARGET_NO_KERNEL)),true)
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
INTERNAL_BOOTIMAGE_ARGS := \
	$(addprefix --second ,$(INSTALLED_2NDBOOTLOADER_TARGET)) \
	--kernel $(INSTALLED_KERNEL_TARGET)

ifneq ($(BOARD_BUILD_SYSTEM_ROOT_IMAGE),true)
INTERNAL_BOOTIMAGE_ARGS += --ramdisk $(INSTALLED_RAMDISK_TARGET)
endif

INTERNAL_BOOTIMAGE_FILES := $(filter-out --%,$(INTERNAL_BOOTIMAGE_ARGS))

ifdef BOARD_KERNEL_BASE
  INTERNAL_BOOTIMAGE_ARGS += --base $(BOARD_KERNEL_BASE)
endif

ifdef BOARD_KERNEL_PAGESIZE
  INTERNAL_BOOTIMAGE_ARGS += --pagesize $(BOARD_KERNEL_PAGESIZE)
endif

ifeq ($(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_SUPPORTS_VERITY),true)
ifeq ($(BOARD_BUILD_SYSTEM_ROOT_IMAGE),true)
VERITY_KEYID := veritykeyid=id:`openssl x509 -in $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_VERITY_SIGNING_KEY).x509.pem -text \
                | grep keyid | sed 's/://g' | tr -d '[:space:]' | tr '[:upper:]' '[:lower:]' | sed 's/keyid//g'`
endif
endif

INTERNAL_KERNEL_CMDLINE := $(strip $(BOARD_KERNEL_CMDLINE) buildvariant=$(TARGET_BUILD_VARIANT) $(VERITY_KEYID))
ifdef INTERNAL_KERNEL_CMDLINE
INTERNAL_BOOTIMAGE_ARGS += --cmdline "$(INTERNAL_KERNEL_CMDLINE)"
endif

INTERNAL_MKBOOTIMG_VERSION_ARGS := \
    --os_version $(PLATFORM_VERSION) \
    --os_patch_level $(PLATFORM_SECURITY_PATCH)

<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
# BOARD_USES_RECOVERY_AS_BOOT = true must have BOARD_BUILD_SYSTEM_ROOT_IMAGE = true.
ifeq ($(BOARD_USES_RECOVERY_AS_BOOT),true)
ifneq ($(BOARD_BUILD_SYSTEM_ROOT_IMAGE),true)
  $(error BOARD_BUILD_SYSTEM_ROOT_IMAGE must be enabled for BOARD_USES_RECOVERY_AS_BOOT.)
endif
endif
=======
ifdef BOARD_GKI_SIGNING_KEY_PATH
ifndef BOARD_GKI_SIGNING_ALGORITHM
$(error BOARD_GKI_SIGNING_ALGORITHM should be defined with BOARD_GKI_SIGNING_KEY_PATH)
endif
INTERNAL_MKBOOTIMG_GKI_SINGING_ARGS := \
    --gki_signing_key $(BOARD_GKI_SIGNING_KEY_PATH) \
    --gki_signing_algorithm $(BOARD_GKI_SIGNING_ALGORITHM) \
    --gki_signing_avbtool_path $(AVBTOOL)
endif

# Using double quote to pass BOARD_GKI_SIGNING_SIGNATURE_ARGS as a single string
# to MKBOOTIMG, although it may contain multiple args.
ifdef BOARD_GKI_SIGNING_SIGNATURE_ARGS
INTERNAL_MKBOOTIMG_GKI_SINGING_ARGS += \
    --gki_signing_signature_args "$(BOARD_GKI_SIGNING_SIGNATURE_ARGS)"
endif

# Define these only if we are building boot
ifdef BUILDING_BOOT_IMAGE
INSTALLED_BOOTIMAGE_TARGET := $(BUILT_BOOTIMAGE_TARGET)
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)

# We build recovery as boot image if BOARD_USES_RECOVERY_AS_BOOT is true.
ifneq ($(BOARD_USES_RECOVERY_AS_BOOT),true)
ifeq ($(TARGET_BOOTIMAGE_USE_EXT2),true)
$(error TARGET_BOOTIMAGE_USE_EXT2 is not supported anymore)

<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
else ifeq (true,$(BOARD_AVB_ENABLE)) # TARGET_BOOTIMAGE_USE_EXT2 != true
=======
$(foreach b,$(INSTALLED_BOOTIMAGE_TARGET), $(eval $(call add-dependency,$(b),$(call bootimage-to-kernel,$(b)))))

ifeq (true,$(BOARD_AVB_ENABLE))

# $1: boot image target
define build_boot_board_avb_enabled
  $(MKBOOTIMG) --kernel $(call bootimage-to-kernel,$(1)) $(INTERNAL_BOOTIMAGE_ARGS) $(INTERNAL_MKBOOTIMG_VERSION_ARGS) \
               $(INTERNAL_MKBOOTIMG_GKI_SINGING_ARGS) $(BOARD_MKBOOTIMG_ARGS) --output $(1)
  $(call assert-max-image-size,$(1),$(call get-hash-image-max-size,$(call get-bootimage-partition-size,$(1),boot)))
  $(AVBTOOL) add_hash_footer \
          --image $(1) \
          --partition_size $(call get-bootimage-partition-size,$(1),boot) \
          --partition_name boot $(INTERNAL_AVB_BOOT_SIGNING_ARGS) \
          $(BOARD_AVB_BOOT_ADD_HASH_FOOTER_ARGS)
endef
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)

$(INSTALLED_BOOTIMAGE_TARGET): $(MKBOOTIMG) $(AVBTOOL) $(INTERNAL_BOOTIMAGE_FILES) $(BOARD_AVB_BOOT_KEY_PATH) $(BOARD_GKI_SIGNING_KEY_PATH)
	$(call pretty,"Target boot image: $@")
	$(hide) $(MKBOOTIMG) $(INTERNAL_BOOTIMAGE_ARGS) $(INTERNAL_MKBOOTIMG_VERSION_ARGS) $(BOARD_MKBOOTIMG_ARGS) --output $@
	$(hide) $(call assert-max-image-size,$@,$(call get-hash-image-max-size,$(BOARD_BOOTIMAGE_PARTITION_SIZE)))
	$(hide) $(AVBTOOL) add_hash_footer \
	  --image $@ \
	  --partition_size $(BOARD_BOOTIMAGE_PARTITION_SIZE) \
	  --partition_name boot $(INTERNAL_AVB_BOOT_SIGNING_ARGS) \
	  $(BOARD_AVB_BOOT_ADD_HASH_FOOTER_ARGS)

.PHONY: bootimage-nodeps
bootimage-nodeps: $(MKBOOTIMG) $(AVBTOOL) $(BOARD_AVB_BOOT_KEY_PATH) $(BOARD_GKI_SIGNING_KEY_PATH)
	@echo "make $@: ignoring dependencies"
	$(hide) $(MKBOOTIMG) $(INTERNAL_BOOTIMAGE_ARGS) $(INTERNAL_MKBOOTIMG_VERSION_ARGS) $(BOARD_MKBOOTIMG_ARGS) --output $(INSTALLED_BOOTIMAGE_TARGET)
	$(hide) $(call assert-max-image-size,$(INSTALLED_BOOTIMAGE_TARGET),$(call get-hash-image-max-size,$(BOARD_BOOTIMAGE_PARTITION_SIZE)))
	$(hide) $(AVBTOOL) add_hash_footer \
	  --image $(INSTALLED_BOOTIMAGE_TARGET) \
	  --partition_size $(BOARD_BOOTIMAGE_PARTITION_SIZE) \
	  --partition_name boot $(INTERNAL_AVB_BOOT_SIGNING_ARGS) \
	  $(BOARD_AVB_BOOT_ADD_HASH_FOOTER_ARGS)

else ifeq (true,$(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_SUPPORTS_BOOT_SIGNER)) # BOARD_AVB_ENABLE != true

$(INSTALLED_BOOTIMAGE_TARGET): $(MKBOOTIMG) $(INTERNAL_BOOTIMAGE_FILES) $(BOOT_SIGNER)
	$(call pretty,"Target boot image: $@")
	$(hide) $(MKBOOTIMG) $(INTERNAL_BOOTIMAGE_ARGS) $(INTERNAL_MKBOOTIMG_VERSION_ARGS) $(BOARD_MKBOOTIMG_ARGS) --output $@
	$(BOOT_SIGNER) /boot $@ $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_VERITY_SIGNING_KEY).pk8 $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_VERITY_SIGNING_KEY).x509.pem $@
	$(hide) $(call assert-max-image-size,$@,$(BOARD_BOOTIMAGE_PARTITION_SIZE))

.PHONY: bootimage-nodeps
bootimage-nodeps: $(MKBOOTIMG) $(BOOT_SIGNER)
	@echo "make $@: ignoring dependencies"
	$(hide) $(MKBOOTIMG) $(INTERNAL_BOOTIMAGE_ARGS) $(INTERNAL_MKBOOTIMG_VERSION_ARGS) $(BOARD_MKBOOTIMG_ARGS) --output $(INSTALLED_BOOTIMAGE_TARGET)
	$(BOOT_SIGNER) /boot $(INSTALLED_BOOTIMAGE_TARGET) $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_VERITY_SIGNING_KEY).pk8 $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_VERITY_SIGNING_KEY).x509.pem $(INSTALLED_BOOTIMAGE_TARGET)
	$(hide) $(call assert-max-image-size,$(INSTALLED_BOOTIMAGE_TARGET),$(BOARD_BOOTIMAGE_PARTITION_SIZE))

else ifeq (true,$(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_SUPPORTS_VBOOT)) # PRODUCT_SUPPORTS_BOOT_SIGNER != true

$(INSTALLED_BOOTIMAGE_TARGET): $(MKBOOTIMG) $(INTERNAL_BOOTIMAGE_FILES) $(VBOOT_SIGNER) $(FUTILITY)
	$(call pretty,"Target boot image: $@")
	$(hide) $(MKBOOTIMG) $(INTERNAL_BOOTIMAGE_ARGS) $(INTERNAL_MKBOOTIMG_VERSION_ARGS) $(BOARD_MKBOOTIMG_ARGS) --output $@.unsigned
	$(VBOOT_SIGNER) $(FUTILITY) $@.unsigned $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_VBOOT_SIGNING_KEY).vbpubk $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_VBOOT_SIGNING_KEY).vbprivk $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_VBOOT_SIGNING_SUBKEY).vbprivk $@.keyblock $@
	$(hide) $(call assert-max-image-size,$@,$(BOARD_BOOTIMAGE_PARTITION_SIZE))

.PHONY: bootimage-nodeps
bootimage-nodeps: $(MKBOOTIMG) $(VBOOT_SIGNER) $(FUTILITY)
	@echo "make $@: ignoring dependencies"
	$(hide) $(MKBOOTIMG) $(INTERNAL_BOOTIMAGE_ARGS) $(INTERNAL_MKBOOTIMG_VERSION_ARGS) $(BOARD_MKBOOTIMG_ARGS) --output $(INSTALLED_BOOTIMAGE_TARGET).unsigned
	$(VBOOT_SIGNER) $(FUTILITY) $(INSTALLED_BOOTIMAGE_TARGET).unsigned $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_VBOOT_SIGNING_KEY).vbpubk $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_VBOOT_SIGNING_KEY).vbprivk $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_VBOOT_SIGNING_SUBKEY).vbprivk $(INSTALLED_BOOTIMAGE_TARGET).keyblock $(INSTALLED_BOOTIMAGE_TARGET)
	$(hide) $(call assert-max-image-size,$(INSTALLED_BOOTIMAGE_TARGET),$(BOARD_BOOTIMAGE_PARTITION_SIZE))

else # PRODUCT_SUPPORTS_VBOOT != true

$(INSTALLED_BOOTIMAGE_TARGET): $(MKBOOTIMG) $(INTERNAL_BOOTIMAGE_FILES)
	$(call pretty,"Target boot image: $@")
	$(hide) $(MKBOOTIMG) $(INTERNAL_BOOTIMAGE_ARGS) $(INTERNAL_MKBOOTIMG_VERSION_ARGS) $(BOARD_MKBOOTIMG_ARGS) --output $@
	$(hide) $(call assert-max-image-size,$@,$(BOARD_BOOTIMAGE_PARTITION_SIZE))

.PHONY: bootimage-nodeps
bootimage-nodeps: $(MKBOOTIMG)
	@echo "make $@: ignoring dependencies"
	$(hide) $(MKBOOTIMG) $(INTERNAL_BOOTIMAGE_ARGS) $(INTERNAL_MKBOOTIMG_VERSION_ARGS) $(BOARD_MKBOOTIMG_ARGS) --output $(INSTALLED_BOOTIMAGE_TARGET)
	$(hide) $(call assert-max-image-size,$(INSTALLED_BOOTIMAGE_TARGET),$(BOARD_BOOTIMAGE_PARTITION_SIZE))

endif # TARGET_BOOTIMAGE_USE_EXT2
endif # BOARD_USES_RECOVERY_AS_BOOT

else	# TARGET_NO_KERNEL
ifdef BOARD_PREBUILT_BOOTIMAGE
<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
ifneq ($(BOARD_BUILD_SYSTEM_ROOT_IMAGE),true)
# Remove when b/63676296 is resolved.
$(error Prebuilt bootimage is only supported for AB targets)
endif
$(eval $(call copy-one-file,$(BOARD_PREBUILT_BOOTIMAGE),$(INSTALLED_BOOTIMAGE_TARGET)))
=======
INTERNAL_PREBUILT_BOOTIMAGE := $(BOARD_PREBUILT_BOOTIMAGE)
INSTALLED_BOOTIMAGE_TARGET := $(PRODUCT_OUT)/boot.img

ifeq ($(BOARD_AVB_ENABLE),true)
$(INSTALLED_BOOTIMAGE_TARGET): $(INTERNAL_PREBUILT_BOOTIMAGE) $(AVBTOOL) $(BOARD_AVB_BOOT_KEY_PATH)
	cp $(INTERNAL_PREBUILT_BOOTIMAGE) $@
	$(AVBTOOL) add_hash_footer \
	    --image $@ \
	    --partition_size $(BOARD_BOOTIMAGE_PARTITION_SIZE) \
	    --partition_name boot $(INTERNAL_AVB_BOOT_SIGNING_ARGS) \
	    $(BOARD_AVB_BOOT_ADD_HASH_FOOTER_ARGS)
else
$(INSTALLED_BOOTIMAGE_TARGET): $(INTERNAL_PREBUILT_BOOTIMAGE)
	cp $(INTERNAL_PREBUILT_BOOTIMAGE) $@
endif # BOARD_AVB_ENABLE

else # BOARD_PREBUILT_BOOTIMAGE not defined
INSTALLED_BOOTIMAGE_TARGET :=
endif # BOARD_PREBUILT_BOOTIMAGE
endif # TARGET_NO_KERNEL
endif # my_installed_prebuilt_gki_apex not defined

my_apex_extracted_boot_image :=
my_installed_prebuilt_gki_apex :=

# -----------------------------------------------------------------
# vendor boot image
ifeq ($(BUILDING_VENDOR_BOOT_IMAGE),true)

ifeq ($(PRODUCT_SUPPORTS_VERITY),true)
  $(error vboot 1.0 does not support vendor_boot partition)
endif

INTERNAL_VENDOR_RAMDISK_FILES := $(filter $(TARGET_VENDOR_RAMDISK_OUT)/%, \
    $(ALL_GENERATED_SOURCES) \
    $(ALL_DEFAULT_INSTALLED_MODULES))

INTERNAL_VENDOR_RAMDISK_TARGET := $(call intermediates-dir-for,PACKAGING,vendor_boot)/vendor_ramdisk.cpio$(RAMDISK_EXT)

# Exclude recovery files in the default vendor ramdisk if including a standalone
# recovery ramdisk in vendor_boot.
ifeq (true,$(BOARD_MOVE_RECOVERY_RESOURCES_TO_VENDOR_BOOT))
ifneq (true,$(BOARD_INCLUDE_RECOVERY_RAMDISK_IN_VENDOR_BOOT))
$(INTERNAL_VENDOR_RAMDISK_TARGET): $(INTERNAL_RECOVERY_RAMDISK_FILES_TIMESTAMP)
$(INTERNAL_VENDOR_RAMDISK_TARGET): PRIVATE_ADDITIONAL_DIR := $(TARGET_RECOVERY_ROOT_OUT)
endif
endif

$(INTERNAL_VENDOR_RAMDISK_TARGET): $(MKBOOTFS) $(INTERNAL_VENDOR_RAMDISK_FILES) | $(COMPRESSION_COMMAND_DEPS)
	$(MKBOOTFS) -d $(TARGET_OUT) $(TARGET_VENDOR_RAMDISK_OUT) $(PRIVATE_ADDITIONAL_DIR) | $(COMPRESSION_COMMAND) > $@

INSTALLED_VENDOR_RAMDISK_TARGET := $(PRODUCT_OUT)/vendor_ramdisk.img
$(INSTALLED_VENDOR_RAMDISK_TARGET): $(INTERNAL_VENDOR_RAMDISK_TARGET)
	@echo "Target vendor ramdisk: $@"
	$(copy-file-to-target)

INSTALLED_FILES_FILE_VENDOR_RAMDISK := $(PRODUCT_OUT)/installed-files-vendor-ramdisk.txt
INSTALLED_FILES_JSON_VENDOR_RAMDISK := $(INSTALLED_FILES_FILE_VENDOR_RAMDISK:.txt=.json)
$(INSTALLED_FILES_FILE_VENDOR_RAMDISK): .KATI_IMPLICIT_OUTPUTS := $(INSTALLED_FILES_JSON_VENDOR_RAMDISK)
$(INSTALLED_FILES_FILE_VENDOR_RAMDISK): $(INTERNAL_VENDOR_RAMDISK_TARGET)
$(INSTALLED_FILES_FILE_VENDOR_RAMDISK): $(INTERNAL_VENDOR_RAMDISK_FILES) $(FILESLIST) $(FILESLIST_UTIL)
	@echo Installed file list: $@
	mkdir -p $(dir $@)
	rm -f $@
	$(FILESLIST) $(TARGET_VENDOR_RAMDISK_OUT) > $(@:.txt=.json)
	$(FILESLIST_UTIL) -c $(@:.txt=.json) > $@

ifdef BOARD_INCLUDE_DTB_IN_BOOTIMG
  INTERNAL_VENDOR_BOOTIMAGE_ARGS += --dtb $(INSTALLED_DTBIMAGE_TARGET)
endif
ifdef BOARD_KERNEL_BASE
  INTERNAL_VENDOR_BOOTIMAGE_ARGS += --base $(BOARD_KERNEL_BASE)
endif
ifdef BOARD_KERNEL_PAGESIZE
  INTERNAL_VENDOR_BOOTIMAGE_ARGS += --pagesize $(BOARD_KERNEL_PAGESIZE)
endif
ifdef INTERNAL_KERNEL_CMDLINE
  INTERNAL_VENDOR_BOOTIMAGE_ARGS += --vendor_cmdline "$(INTERNAL_KERNEL_CMDLINE)"
endif

ifdef INTERNAL_BOOTCONFIG
  INTERNAL_VENDOR_BOOTCONFIG_TARGET := $(PRODUCT_OUT)/vendor-bootconfig.img
  $(INTERNAL_VENDOR_BOOTCONFIG_TARGET):
	rm -f $@
	$(foreach param,$(INTERNAL_BOOTCONFIG), \
	 printf "%s\n" $(param) >> $@;)
  INTERNAL_VENDOR_BOOTIMAGE_ARGS += --vendor_bootconfig $(INTERNAL_VENDOR_BOOTCONFIG_TARGET)
endif

# $(1): Build target name
# $(2): Staging dir to be compressed
# $(3): Build dependencies
define build-vendor-ramdisk-fragment-target
$(1): $(3) $(MKBOOTFS) | $(COMPRESSION_COMMAND_DEPS)
	$(MKBOOTFS) -d $(TARGET_OUT) $(2) | $(COMPRESSION_COMMAND) > $$@
endef

# $(1): Ramdisk name
define build-vendor-ramdisk-fragment
$(strip \
  $(eval build_target := $(call intermediates-dir-for,PACKAGING,vendor_ramdisk_fragments)/$(1).cpio$(RAMDISK_EXT)) \
  $(eval $(call build-vendor-ramdisk-fragment-target,$(build_target),$(VENDOR_RAMDISK_FRAGMENT.$(1).STAGING_DIR),$(VENDOR_RAMDISK_FRAGMENT.$(1).FILES))) \
  $(build_target) \
)
endef

# $(1): Ramdisk name
# $(2): Prebuilt file path
define build-prebuilt-vendor-ramdisk-fragment
$(strip \
  $(eval build_target := $(call intermediates-dir-for,PACKAGING,prebuilt_vendor_ramdisk_fragments)/$(1)) \
  $(eval $(call copy-one-file,$(2),$(build_target))) \
  $(build_target) \
)
endef

INTERNAL_VENDOR_RAMDISK_FRAGMENT_TARGETS :=
INTERNAL_VENDOR_RAMDISK_FRAGMENT_ARGS :=

$(foreach vendor_ramdisk_fragment,$(INTERNAL_VENDOR_RAMDISK_FRAGMENTS), \
  $(eval prebuilt_vendor_ramdisk_fragment_file := $(BOARD_VENDOR_RAMDISK_FRAGMENT.$(vendor_ramdisk_fragment).PREBUILT)) \
  $(if $(prebuilt_vendor_ramdisk_fragment_file), \
    $(eval vendor_ramdisk_fragment_target := $(call build-prebuilt-vendor-ramdisk-fragment,$(vendor_ramdisk_fragment),$(prebuilt_vendor_ramdisk_fragment_file))) \
    $(eval ### else ###), \
    $(eval vendor_ramdisk_fragment_target := $(call build-vendor-ramdisk-fragment,$(vendor_ramdisk_fragment)))) \
  $(eval INTERNAL_VENDOR_RAMDISK_FRAGMENT_TARGETS += $(vendor_ramdisk_fragment_target)) \
  $(eval INTERNAL_VENDOR_RAMDISK_FRAGMENT_ARGS += $(BOARD_VENDOR_RAMDISK_FRAGMENT.$(vendor_ramdisk_fragment).MKBOOTIMG_ARGS) --vendor_ramdisk_fragment $(vendor_ramdisk_fragment_target)) \
)

INSTALLED_VENDOR_BOOTIMAGE_TARGET := $(PRODUCT_OUT)/vendor_boot.img
$(INSTALLED_VENDOR_BOOTIMAGE_TARGET): $(MKBOOTIMG) $(INTERNAL_VENDOR_RAMDISK_TARGET) $(INSTALLED_DTBIMAGE_TARGET)
$(INSTALLED_VENDOR_BOOTIMAGE_TARGET): $(INTERNAL_VENDOR_RAMDISK_FRAGMENT_TARGETS) $(INTERNAL_VENDOR_BOOTCONFIG_TARGET)
ifeq ($(BOARD_AVB_ENABLE),true)
$(INSTALLED_VENDOR_BOOTIMAGE_TARGET): $(AVBTOOL) $(BOARD_AVB_VENDOR_BOOTIMAGE_KEY_PATH)
	$(call pretty,"Target vendor_boot image: $@")
	$(MKBOOTIMG) $(INTERNAL_VENDOR_BOOTIMAGE_ARGS) $(BOARD_MKBOOTIMG_ARGS) --vendor_ramdisk $(INTERNAL_VENDOR_RAMDISK_TARGET) $(INTERNAL_VENDOR_RAMDISK_FRAGMENT_ARGS) --vendor_boot $@
	$(call assert-max-image-size,$@,$(BOARD_VENDOR_BOOTIMAGE_PARTITION_SIZE))
	$(AVBTOOL) add_hash_footer \
           --image $@ \
	   --partition_size $(BOARD_VENDOR_BOOTIMAGE_PARTITION_SIZE) \
	   --partition_name vendor_boot $(INTERNAL_AVB_VENDOR_BOOT_SIGNING_ARGS) \
	   $(BOARD_AVB_VENDOR_BOOT_ADD_HASH_FOOTER_ARGS)
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
else
INTERNAL_KERNEL_CMDLINE := $(strip $(BOARD_KERNEL_CMDLINE))
# HACK: The top-level targets depend on the bootimage.  Not all targets
# can produce a bootimage, though, and emulator targets need the ramdisk
# instead.  Fake it out by calling the ramdisk the bootimage.
# TODO: make the emulator use bootimages, and make mkbootimg accept
#       kernel-less inputs.
INSTALLED_BOOTIMAGE_TARGET := $(INSTALLED_RAMDISK_TARGET)
endif
endif

# -----------------------------------------------------------------
# NOTICE files
#
# We are required to publish the licenses for all code under BSD, GPL and
# Apache licenses (and possibly other more exotic ones as well). We err on the
# side of caution, so the licenses for other third-party code are included here
# too.
#
# This needs to be before the systemimage rules, because it adds to
# ALL_DEFAULT_INSTALLED_MODULES, which those use to pick which files
# go into the systemimage.

.PHONY: notice_files

# Create the rule to combine the files into text and html/xml forms
# $(1) - xml_excluded_vendor|xml_vendor|html
# $(2) - Plain text output file
# $(3) - HTML/XML output file
# $(4) - File title
# $(5) - Directory to use.  Notice files are all $(4)/src.  Other
#		 directories in there will be used for scratch
# $(6) - Dependencies for the output files
#
# The algorithm here is that we go collect a hash for each of the notice
# files and write the names of the files that match that hash.  Then
# to generate the real files, we go print out all of the files and their
# hashes.
#
# These rules are fairly complex, so they depend on this makefile so if
# it changes, they'll run again.
#
# TODO: We could clean this up so that we just record the locations of the
# original notice files instead of making rules to copy them somwehere.
# Then we could traverse that without quite as much bash drama.
define combine-notice-files
$(2) $(3): PRIVATE_MESSAGE := $(4)
$(2) $(3): PRIVATE_DIR := $(5)
$(2) : $(3)
$(3) : $(6) $(BUILD_SYSTEM)/Makefile build/make/tools/generate-notice-files.py
	build/make/tools/generate-notice-files.py --text-output $(2) \
		$(if $(filter $(1),xml_excluded_vendor),-e vendor --xml-output, \
		  $(if $(filter $(1),xml_vendor),-i vendor --xml-output, \
		    --html-output)) $(3) \
		-t $$(PRIVATE_MESSAGE) -s $$(PRIVATE_DIR)/src
notice_files: $(2) $(3)
endef

# TODO These intermediate NOTICE.txt/NOTICE.html files should go into
# TARGET_OUT_NOTICE_FILES now that the notice files are gathered from
# the src subdirectory.

target_notice_file_txt := $(TARGET_OUT_INTERMEDIATES)/NOTICE.txt
target_notice_file_html_or_xml := $(TARGET_OUT_INTERMEDIATES)/NOTICE.html
target_notice_file_html_or_xml_gz := $(TARGET_OUT_INTERMEDIATES)/NOTICE.html.gz
installed_notice_html_or_xml_gz := $(TARGET_OUT)/etc/NOTICE.html.gz
tools_notice_file_txt := $(HOST_OUT_INTERMEDIATES)/NOTICE.txt
tools_notice_file_html := $(HOST_OUT_INTERMEDIATES)/NOTICE.html

# TODO(b/69865032): Make PRODUCT_NOTICE_SPLIT the default behavior.
<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
ifeq ($(PRODUCT_NOTICE_SPLIT),true)
target_notice_file_html_or_xml := $(TARGET_OUT_INTERMEDIATES)/NOTICE.xml
target_notice_file_html_or_xml_gz := $(TARGET_OUT_INTERMEDIATES)/NOTICE.xml.gz
=======
ifneq ($(PRODUCT_NOTICE_SPLIT),true)
target_notice_file_html := $(TARGET_OUT_INTERMEDIATES)/NOTICE.html
target_notice_file_html_gz := $(TARGET_OUT_INTERMEDIATES)/NOTICE.html.gz
installed_notice_html_or_xml_gz := $(TARGET_OUT)/etc/NOTICE.html.gz
$(eval $(call combine-notice-files, html, \
	        $(target_notice_file_txt), \
	        $(target_notice_file_html), \
	        "Notices for files contained in the filesystem images in this directory:", \
	        $(TARGET_OUT_NOTICE_FILES), \
	        $(ALL_DEFAULT_INSTALLED_MODULES) $(kernel_notice_file), \
	        $(exclude_target_dirs)))
$(target_notice_file_html_gz): $(target_notice_file_html) | $(MINIGZIP)
	$(hide) $(MINIGZIP) -9 < $< > $@
$(installed_notice_html_or_xml_gz): $(target_notice_file_html_gz)
	$(copy-file-to-target)

$(call declare-0p-target,$(target_notice_file_html_gz))
$(call declare-0p-target,$(installed_notice_html_or_xml_gz))
else
target_notice_file_xml := $(TARGET_OUT_INTERMEDIATES)/NOTICE.xml
target_notice_file_xml_gz := $(TARGET_OUT_INTERMEDIATES)/NOTICE.xml.gz
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
installed_notice_html_or_xml_gz := $(TARGET_OUT)/etc/NOTICE.xml.gz

target_vendor_notice_file_txt := $(TARGET_OUT_INTERMEDIATES)/NOTICE_VENDOR.txt
target_vendor_notice_file_xml := $(TARGET_OUT_INTERMEDIATES)/NOTICE_VENDOR.xml
target_vendor_notice_file_xml_gz := $(TARGET_OUT_INTERMEDIATES)/NOTICE_VENDOR.xml.gz
installed_vendor_notice_xml_gz := $(TARGET_OUT_VENDOR)/etc/NOTICE.xml.gz
<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
=======

target_product_notice_file_txt := $(TARGET_OUT_INTERMEDIATES)/NOTICE_PRODUCT.txt
target_product_notice_file_xml := $(TARGET_OUT_INTERMEDIATES)/NOTICE_PRODUCT.xml
target_product_notice_file_xml_gz := $(TARGET_OUT_INTERMEDIATES)/NOTICE_PRODUCT.xml.gz
installed_product_notice_xml_gz := $(TARGET_OUT_PRODUCT)/etc/NOTICE.xml.gz

target_system_ext_notice_file_txt := $(TARGET_OUT_INTERMEDIATES)/NOTICE_SYSTEM_EXT.txt
target_system_ext_notice_file_xml := $(TARGET_OUT_INTERMEDIATES)/NOTICE_SYSTEM_EXT.xml
target_system_ext_notice_file_xml_gz := $(TARGET_OUT_INTERMEDIATES)/NOTICE_SYSTEM_EXT.xml.gz
installed_system_ext_notice_xml_gz := $(TARGET_OUT_SYSTEM_EXT)/etc/NOTICE.xml.gz

target_odm_notice_file_txt := $(TARGET_OUT_INTERMEDIATES)/NOTICE_ODM.txt
target_odm_notice_file_xml := $(TARGET_OUT_INTERMEDIATES)/NOTICE_ODM.xml
target_odm_notice_file_xml_gz := $(TARGET_OUT_INTERMEDIATES)/NOTICE_ODM.xml.gz
installed_odm_notice_xml_gz := $(TARGET_OUT_ODM)/etc/NOTICE.xml.gz

target_vendor_dlkm_notice_file_txt := $(TARGET_OUT_INTERMEDIATES)/NOTICE_VENDOR_DLKM.txt
target_vendor_dlkm_notice_file_xml := $(TARGET_OUT_INTERMEDIATES)/NOTICE_VENDOR_DLKM.xml
target_vendor_dlkm_notice_file_xml_gz := $(TARGET_OUT_INTERMEDIATES)/NOTICE_VENDOR_DLKM.xml.gz
installed_vendor_dlkm_notice_xml_gz := $(TARGET_OUT_VENDOR_DLKM)/etc/NOTICE.xml.gz

target_odm_dlkm_notice_file_txt := $(TARGET_OUT_INTERMEDIATES)/NOTICE_ODM_DLKM.txt
target_odm_dlkm_notice_file_xml := $(TARGET_OUT_INTERMEDIATES)/NOTICE_ODM_DLKM.xml
target_odm_dlkm_notice_file_xml_gz := $(TARGET_OUT_INTERMEDIATES)/NOTICE_ODM_DLKM.xml.gz
installed_odm_dlkm_notice_xml_gz := $(TARGET_OUT_ODM_DLKM)/etc/NOTICE.xml.gz

# Notice files are copied to TARGET_OUT_NOTICE_FILES as a side-effect of their module
# being built. A notice xml file must depend on all modules that could potentially
# install a license file relevant to it.
license_modules := $(ALL_DEFAULT_INSTALLED_MODULES) $(kernel_notice_file)
# Only files copied to a system image need system image notices.
license_modules := $(filter $(PRODUCT_OUT)/%,$(license_modules))
# Phonys/fakes don't have notice files (though their deps might)
license_modules := $(filter-out $(TARGET_OUT_FAKE)/%,$(license_modules))
# testcases are not relevant to the system image.
license_modules := $(filter-out $(TARGET_OUT_TESTCASES)/%,$(license_modules))
# filesystem images: system, vendor, product, system_ext, odm, vendor_dlkm, and odm_dlkm
license_modules_system := $(filter $(TARGET_OUT)/%,$(license_modules))
# system_other is relevant to system partition.
license_modules_system += $(filter $(TARGET_OUT_SYSTEM_OTHER)/%,$(license_modules))
license_modules_vendor := $(filter $(TARGET_OUT_VENDOR)/%,$(license_modules))
license_modules_product := $(filter $(TARGET_OUT_PRODUCT)/%,$(license_modules))
license_modules_system_ext := $(filter $(TARGET_OUT_SYSTEM_EXT)/%,$(license_modules))
license_modules_odm := $(filter $(TARGET_OUT_ODM)/%,$(license_modules))
license_modules_vendor_dlkm := $(filter $(TARGET_OUT_VENDOR_DLKM)/%,$(license_modules))
license_modules_odm_dlkm := $(filter $(TARGET_OUT_ODM_DLKM)/%,$(license_modules))
license_modules_agg := $(license_modules_system) \
                       $(license_modules_vendor) \
                       $(license_modules_product) \
                       $(license_modules_system_ext) \
                       $(license_modules_odm) \
                       $(license_modules_vendor_dlkm) \
                       $(license_modules_odm_dlkm)
# targets used for debug symbols only and do not get copied to the device
license_modules_symbols_only := $(filter $(PRODUCT_OUT)/apex/%,$(license_modules))

license_modules_rest := $(filter-out $(license_modules_agg),$(license_modules))
license_modules_rest := $(filter-out $(license_modules_symbols_only),$(license_modules_rest))

# Identify the other targets we expect to have notices for:
# targets copied to the device but are not readable by the UI (e.g. must boot
# into a different partition to read or don't have an associated /etc
# directory) must have their notices built somewhere readable.
license_modules_rehomed := $(filter-out $(PRODUCT_OUT)/%/%,$(license_modules_rest))  # files in root have no /etc
license_modules_rehomed += $(filter $(PRODUCT_OUT)/recovery/%,$(license_modules_rest))
license_modules_rehomed += $(filter $(PRODUCT_OUT)/root/%,$(license_modules_rest))
license_modules_rehomed += $(filter $(PRODUCT_OUT)/data/%,$(license_modules_rest))
license_modules_rehomed += $(filter $(PRODUCT_OUT)/ramdisk/%,$(license_modules_rest))
license_modules_rehomed += $(filter $(PRODUCT_OUT)/debug_ramdisk/%,$(license_modules_rest))
license_modules_rehomed += $(filter $(PRODUCT_OUT)/vendor_ramdisk/%,$(license_modules_rest))
license_modules_rehomed += $(filter $(PRODUCT_OUT)/persist/%,$(license_modules_rest))
license_modules_rehomed += $(filter $(PRODUCT_OUT)/persist.img,$(license_modules_rest))
license_modules_rehomed += $(filter $(PRODUCT_OUT)/system_other/%,$(license_modules_rest))
license_modules_rehomed += $(filter $(PRODUCT_OUT)/kernel%,$(license_modules_rest))
license_modules_rehomed += $(filter $(PRODUCT_OUT)/%.img,$(license_modules_rest))
license_modules_rehomed += $(filter $(PRODUCT_OUT)/%.bin,$(license_modules_rest))

# after removing targets in system images, targets reported in system images, and
# targets used for debug symbols that do not need notices, nothing must remain.
license_modules_rest := $(filter-out $(license_modules_rehomed),$(license_modules_rest))
$(call maybe-print-list-and-error, $(license_modules_rest), \
        "Targets added under $(PRODUCT_OUT)/ unaccounted for notice handling.")

# If we are building in a configuration that includes a prebuilt vendor.img, we can't
# update its notice file, so include those notices in the system partition instead
ifdef BOARD_PREBUILT_VENDORIMAGE
license_modules_system += $(license_modules_rehomed)
system_xml_directories := xml_excluded_vendor_product_odm_vendor_dlkm_odm_dlkm
system_notice_file_message := "Notices for files contained in all filesystem images except vendor/system_ext/product/odm/vendor_dlkm/odm_dlkm in this directory:"
else
license_modules_vendor += $(license_modules_rehomed)
system_xml_directories := xml_system
system_notice_file_message := "Notices for files contained in the system filesystem image in this directory:"
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
endif

ifndef TARGET_BUILD_APPS
kernel_notice_file := $(TARGET_OUT_NOTICE_FILES)/src/kernel.txt
winpthreads_notice_file := $(TARGET_OUT_NOTICE_FILES)/src/winpthreads.txt
pdk_fusion_notice_files := $(filter $(TARGET_OUT_NOTICE_FILES)/%, $(ALL_PDK_FUSION_FILES))

<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
ifdef target_vendor_notice_file_xml_gz
$(eval $(call combine-notice-files, xml_excluded_vendor, \
			$(target_notice_file_txt), \
			$(target_notice_file_html_or_xml), \
			"Notices for files contained in the filesystem images in this directory:", \
			$(TARGET_OUT_NOTICE_FILES), \
			$(ALL_DEFAULT_INSTALLED_MODULES) $(kernel_notice_file) $(pdk_fusion_notice_files)))
$(eval $(call combine-notice-files, xml_vendor, \
			$(target_vendor_notice_file_txt), \
			$(target_vendor_notice_file_xml), \
			"Notices for files contained in the vendor filesystem image in this directory:", \
			$(TARGET_OUT_NOTICE_FILES), \
			$(target_notice_file_html_or_xml)))
else
$(eval $(call combine-notice-files, html, \
			$(target_notice_file_txt), \
			$(target_notice_file_html_or_xml), \
			"Notices for files contained in the filesystem images in this directory:", \
			$(TARGET_OUT_NOTICE_FILES), \
			$(ALL_DEFAULT_INSTALLED_MODULES) $(kernel_notice_file) $(pdk_fusion_notice_files)))
endif
=======
$(target_notice_file_xml_gz): $(target_notice_file_xml) | $(MINIGZIP)
	$(hide) $(MINIGZIP) -9 < $< > $@
$(target_vendor_notice_file_xml_gz): $(target_vendor_notice_file_xml) | $(MINIGZIP)
	$(hide) $(MINIGZIP) -9 < $< > $@
$(target_product_notice_file_xml_gz): $(target_product_notice_file_xml) | $(MINIGZIP)
	$(hide) $(MINIGZIP) -9 < $< > $@
$(target_system_ext_notice_file_xml_gz): $(target_system_ext_notice_file_xml) | $(MINIGZIP)
	$(hide) $(MINIGZIP) -9 < $< > $@
$(target_odm_notice_file_xml_gz): $(target_odm_notice_file_xml) | $(MINIGZIP)
	$(hide) $(MINIGZIP) -9 < $< > $@
$(target_vendor_dlkm_notice_file_xml_gz): $(target_vendor_dlkm_notice_file_xml) | $(MINIGZIP)
	$(hide) $(MINIGZIP) -9 < $< > $@
$(target_odm_dlkm_notice_file_xml_gz): $(target_odm_dlkm_notice_file_xml) | $(MINIGZIP)
	$(hide) $(MINIGZIP) -9 < $< > $@
$(installed_notice_html_or_xml_gz): $(target_notice_file_xml_gz)
	$(copy-file-to-target)
$(installed_vendor_notice_xml_gz): $(target_vendor_notice_file_xml_gz)
	$(copy-file-to-target)
$(installed_product_notice_xml_gz): $(target_product_notice_file_xml_gz)
	$(copy-file-to-target)
$(installed_system_ext_notice_xml_gz): $(target_system_ext_notice_file_xml_gz)
	$(copy-file-to-target)
$(installed_odm_notice_xml_gz): $(target_odm_notice_file_xml_gz)
	$(copy-file-to-target)
$(installed_vendor_dlkm_notice_xml_gz): $(target_vendor_dlkm_notice_file_xml_gz)
	$(copy-file-to-target)
$(installed_odm_dlkm_notice_xml_gz): $(target_odm_dlkm_notice_file_xml_gz)
	$(copy-file-to-target)

$(call declare-0p-target,$(target_notice_file_xml))
$(call declare-0p-target,$(target_notice_file_xml_gz))
$(call declare-0p-target,$(target_vendor_notice_file_xml))
$(call declare-0p-target,$(target_vendor_notice_file_xml_gz))
$(call declare-0p-target,$(target_product_notice_file_xml))
$(call declare-0p-target,$(target_product_notice_file_xml_gz))
$(call declare-0p-target,$(target_system_ext_notice_file_xml))
$(call declare-0p-target,$(target_system_ext_notice_file_xml_gz))
$(call declare-0p-target,$(target_odm_notice_file_xml))
$(call declare-0p-target,$(target_odm_notice_file_xml_gz))
$(call declare-0p-target,$(target_vendor_dlkm_notice_file_xml))
$(call declare-0p-target,$(target_vendor_dlkm_notice_file_xml_gz))
$(call declare-0p-target,$(target_odm_dlkm_notice_file_xml))
$(call declare-0p-target,$(target_odm_dlkm_notice_file_xml_gz))
$(call declare-0p-target,$(installed_notice_html_or_xml_gz))
$(call declare-0p-target,$(installed_vendor_notice_xml_gz))
$(call declare-0p-target,$(installed_product_notice_xml_gz))
$(call declare-0p-target,$(installed_system_ext_notice_xml_gz))
$(call declare-0p-target,$(installed_odm_notice_xml_gz))
$(call declare-0p-target,$(installed_vendor_dlkm_notice_xml_gz))
$(call declare-0p-target,$(installed_odm_dlkm_notice_xml_gz))

ALL_DEFAULT_INSTALLED_MODULES += $(installed_notice_html_or_xml_gz)
ALL_DEFAULT_INSTALLED_MODULES += $(installed_vendor_notice_xml_gz)
ALL_DEFAULT_INSTALLED_MODULES += $(installed_product_notice_xml_gz)
ALL_DEFAULT_INSTALLED_MODULES += $(installed_system_ext_notice_xml_gz)
ALL_DEFAULT_INSTALLED_MODULES += $(installed_odm_notice_xml_gz)
ALL_DEFAULT_INSTALLED_MODULES += $(installed_vendor_dlkm_notice_xml_gz)
ALL_DEFAULT_INSTALLED_MODULES += $(installed_odm_dlkm_notice_xml_gz)
endif # PRODUCT_NOTICE_SPLIT

ALL_DEFAULT_INSTALLED_MODULES += $(installed_notice_html_or_xml_gz)
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)

$(eval $(call combine-notice-files, html, \
			$(tools_notice_file_txt), \
			$(tools_notice_file_html), \
			"Notices for files contained in the tools directory:", \
			$(HOST_OUT_NOTICE_FILES), \
			$(ALL_DEFAULT_INSTALLED_MODULES) \
			$(winpthreads_notice_file)))

# Install the html file at /system/etc/NOTICE.html.gz.
# This is not ideal, but this is very late in the game, after a lot of
# the module processing has already been done -- in fact, we used the
# fact that all that has been done to get the list of modules that we
# need notice files for.
$(target_notice_file_html_or_xml_gz): $(target_notice_file_html_or_xml) | $(MINIGZIP)
	$(hide) $(MINIGZIP) -9 < $< > $@
$(installed_notice_html_or_xml_gz): $(target_notice_file_html_or_xml_gz)
	$(copy-file-to-target)

ifdef target_vendor_notice_file_xml_gz
# Install the vendor html file at /vendor/etc/NOTICE.xml.gz.
$(target_vendor_notice_file_xml_gz): $(target_vendor_notice_file_xml) | $(MINIGZIP)
	$(hide) $(MINIGZIP) -9 < $< > $@
$(installed_vendor_notice_xml_gz): $(target_vendor_notice_file_xml_gz)
	$(copy-file-to-target)
endif

# if we've been run my mm, mmm, etc, don't reinstall this every time
ifeq ($(ONE_SHOT_MAKEFILE),)
  ALL_DEFAULT_INSTALLED_MODULES += $(installed_notice_html_or_xml_gz)
  ifdef target_vendor_notice_file_xml_gz
    ALL_DEFAULT_INSTALLED_MODULES += $(installed_vendor_notice_xml_gz)
  endif
endif
endif  # TARGET_BUILD_APPS

# The kernel isn't really a module, so to get its module file in there, we
# make the target NOTICE files depend on this particular file too, which will
# then be in the right directory for the find in combine-notice-files to work.
$(kernel_notice_file): \
	    $(BUILD_SYSTEM)/LINUX_KERNEL_COPYING \
	    | $(ACP)
	@echo Copying: $@
	$(hide) mkdir -p $(dir $@)
	$(hide) $(ACP) $< $@

$(winpthreads_notice_file): \
	    $(BUILD_SYSTEM)/WINPTHREADS_COPYING \
	    | $(ACP)
	@echo Copying: $@
	$(hide) mkdir -p $(dir $@)
	$(hide) $(ACP) $< $@

# -----------------------------------------------------------------
# Build a keystore with the authorized keys in it, used to verify the
# authenticity of downloaded OTA packages.
#
# This rule adds to ALL_DEFAULT_INSTALLED_MODULES, so it needs to come
# before the rules that use that variable to build the image.
ALL_DEFAULT_INSTALLED_MODULES += $(TARGET_OUT_ETC)/security/otacerts.zip
$(TARGET_OUT_ETC)/security/otacerts.zip: KEY_CERT_PAIR := $(DEFAULT_KEY_CERT_PAIR)
$(TARGET_OUT_ETC)/security/otacerts.zip: $(addsuffix .x509.pem,$(DEFAULT_KEY_CERT_PAIR)) | $(ZIPTIME)
	$(hide) rm -f $@
	$(hide) mkdir -p $(dir $@)
	$(hide) zip -qjX $@ $<
	$(remove-timestamps-from-package)

# Carry the public key for update_engine if it's a non-IoT target that
# uses the AB updater. We use the same key as otacerts but in RSA public key
# format.
ifeq ($(AB_OTA_UPDATER),true)
ifneq ($(PRODUCT_IOT),true)
ALL_DEFAULT_INSTALLED_MODULES += $(TARGET_OUT_ETC)/update_engine/update-payload-key.pub.pem
$(TARGET_OUT_ETC)/update_engine/update-payload-key.pub.pem: $(addsuffix .x509.pem,$(DEFAULT_KEY_CERT_PAIR))
	$(hide) rm -f $@
	$(hide) mkdir -p $(dir $@)
	$(hide) openssl x509 -pubkey -noout -in $< > $@

ALL_DEFAULT_INSTALLED_MODULES += $(TARGET_RECOVERY_ROOT_OUT)/etc/update_engine/update-payload-key.pub.pem
$(TARGET_RECOVERY_ROOT_OUT)/etc/update_engine/update-payload-key.pub.pem: $(TARGET_OUT_ETC)/update_engine/update-payload-key.pub.pem
	$(hide) cp -f $< $@
endif
endif

.PHONY: otacerts
otacerts: $(TARGET_OUT_ETC)/security/otacerts.zip


# #################################################################
# Targets for user images
# #################################################################

INTERNAL_USERIMAGES_EXT_VARIANT :=
ifeq ($(TARGET_USERIMAGES_USE_EXT2),true)
INTERNAL_USERIMAGES_USE_EXT := true
INTERNAL_USERIMAGES_EXT_VARIANT := ext2
else
ifeq ($(TARGET_USERIMAGES_USE_EXT3),true)
INTERNAL_USERIMAGES_USE_EXT := true
INTERNAL_USERIMAGES_EXT_VARIANT := ext3
else
ifeq ($(TARGET_USERIMAGES_USE_EXT4),true)
INTERNAL_USERIMAGES_USE_EXT := true
INTERNAL_USERIMAGES_EXT_VARIANT := ext4
endif
endif
endif

# These options tell the recovery updater/installer how to mount the partitions writebale.
# <fstype>=<fstype_opts>[|<fstype_opts>]...
# fstype_opts := <opt>[,<opt>]...
#         opt := <name>[=<value>]
# The following worked on Nexus devices with Kernel 3.1, 3.4, 3.10
DEFAULT_TARGET_RECOVERY_FSTYPE_MOUNT_OPTIONS := ext4=max_batch_time=0,commit=1,data=ordered,barrier=1,errors=panic,nodelalloc

ifneq (true,$(TARGET_USERIMAGES_SPARSE_EXT_DISABLED))
  INTERNAL_USERIMAGES_SPARSE_EXT_FLAG := -s
endif

ifeq ($(INTERNAL_USERIMAGES_USE_EXT),true)
INTERNAL_USERIMAGES_DEPS := $(SIMG2IMG)
INTERNAL_USERIMAGES_DEPS += $(MKEXTUSERIMG) $(MAKE_EXT4FS) $(E2FSCK)
ifeq ($(TARGET_USERIMAGES_USE_F2FS),true)
INTERNAL_USERIMAGES_DEPS += $(MKF2FSUSERIMG) $(MAKE_F2FS)
endif
endif

ifeq ($(BOARD_AVB_ENABLE),true)
INTERNAL_USERIMAGES_DEPS += $(AVBTOOL)
endif

ifneq (true,$(TARGET_USERIMAGES_SPARSE_SQUASHFS_DISABLED))
  INTERNAL_USERIMAGES_SPARSE_SQUASHFS_FLAG := -s
endif
ifneq ($(filter $(BOARD_PRODUCTIMAGE_FILE_SYSTEM_TYPE) $(BOARD_VENDORIMAGE_FILE_SYSTEM_TYPE) $(BOARD_SYSTEMIMAGE_FILE_SYSTEM_TYPE),squashfs),)
INTERNAL_USERIMAGES_DEPS += $(MAKE_SQUASHFS) $(MKSQUASHFSUSERIMG) $(IMG2SIMG)
endif

INTERNAL_USERIMAGES_BINARY_PATHS := $(sort $(dir $(INTERNAL_USERIMAGES_DEPS)))

ifeq (true,$(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_SUPPORTS_VERITY))
INTERNAL_USERIMAGES_DEPS += $(BUILD_VERITY_TREE) $(APPEND2SIMG) $(VERITY_SIGNER)
ifeq (true,$(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_SUPPORTS_VERITY_FEC))
INTERNAL_USERIMAGES_DEPS += $(FEC)
endif
endif

SELINUX_FC := $(call intermediates-dir-for,ETC,file_contexts.bin)/file_contexts.bin
INTERNAL_USERIMAGES_DEPS += $(SELINUX_FC)

INTERNAL_USERIMAGES_DEPS += $(BLK_ALLOC_TO_BASE_FS)

ifeq ($(INTERNAL_USERIMAGES_USE_EXT),true)
INTERNAL_USERIMAGES_DEPS += $(MKE2FS_CONF)
endif

# $(1) the partition name (eg system)
# $(2) the image prop file
define add-common-flags-to-image-props
$(eval _var := $(call to-upper,$(1)))
$(hide) echo "$(1)_selinux_fc=$(SELINUX_FC)" >> $(2)
$(hide) echo "building_$(1)_image=$(BUILDING_$(_var)_IMAGE)" >> $(2)
endef

# $(1) the partition name (eg system)
# $(2) the image prop file
define add-common-ro-flags-to-image-props
$(eval _var := $(call to-upper,$(1)))
$(if $(BOARD_$(_var)IMAGE_EROFS_COMPRESSOR),$(hide) echo "$(1)_erofs_compressor=$(BOARD_$(_var)IMAGE_EROFS_COMPRESSOR)" >> $(2))
$(if $(BOARD_$(_var)IMAGE_EROFS_PCLUSTER_SIZE),$(hide) echo "$(1)_erofs_pcluster_size=$(BOARD_$(_var)IMAGE_EROFS_PCLUSTER_SIZE)" >> $(2))
$(if $(BOARD_$(_var)IMAGE_EXTFS_INODE_COUNT),$(hide) echo "$(1)_extfs_inode_count=$(BOARD_$(_var)IMAGE_EXTFS_INODE_COUNT)" >> $(2))
$(if $(BOARD_$(_var)IMAGE_EXTFS_RSV_PCT),$(hide) echo "$(1)_extfs_rsv_pct=$(BOARD_$(_var)IMAGE_EXTFS_RSV_PCT)" >> $(2))
$(if $(BOARD_$(_var)IMAGE_F2FS_SLOAD_COMPRESS_FLAGS),$(hide) echo "$(1)_f2fs_sldc_flags=$(BOARD_$(_var)IMAGE_F2FS_SLOAD_COMPRESS_FLAGS)" >> $(2))
$(if $(BOARD_$(_var)IMAGE_FILE_SYSTEM_COMPRESS),$(hide) echo "$(1)_f2fs_compress=$(BOARD_$(_var)IMAGE_FILE_SYSTEM_COMPRESS)" >> $(2))
$(if $(BOARD_$(_var)IMAGE_FILE_SYSTEM_TYPE),$(hide) echo "$(1)_fs_type=$(BOARD_$(_var)IMAGE_FILE_SYSTEM_TYPE)" >> $(2))
$(if $(BOARD_$(_var)IMAGE_JOURNAL_SIZE),$(hide) echo "$(1)_journal_size=$(BOARD_$(_var)IMAGE_JOURNAL_SIZE)" >> $(2))
$(if $(BOARD_$(_var)IMAGE_PARTITION_RESERVED_SIZE),$(hide) echo "$(1)_reserved_size=$(BOARD_$(_var)IMAGE_PARTITION_RESERVED_SIZE)" >> $(2))
$(if $(BOARD_$(_var)IMAGE_PARTITION_SIZE),$(hide) echo "$(1)_size=$(BOARD_$(_var)IMAGE_PARTITION_SIZE)" >> $(2))
$(if $(BOARD_$(_var)IMAGE_SQUASHFS_BLOCK_SIZE),$(hide) echo "$(1)_squashfs_block_size=$(BOARD_$(_var)IMAGE_SQUASHFS_BLOCK_SIZE)" >> $(2))
$(if $(BOARD_$(_var)IMAGE_SQUASHFS_COMPRESSOR),$(hide) echo "$(1)_squashfs_compressor=$(BOARD_$(_var)IMAGE_SQUASHFS_COMPRESSOR)" >> $(2))
$(if $(BOARD_$(_var)IMAGE_SQUASHFS_COMPRESSOR_OPT),$(hide) echo "$(1)_squashfs_compressor_opt=$(BOARD_$(_var)IMAGE_SQUASHFS_COMPRESSOR_OPT)" >> $(2))
$(if $(BOARD_$(_var)IMAGE_SQUASHFS_DISABLE_4K_ALIGN),$(hide) echo "$(1)_squashfs_disable_4k_align=$(BOARD_$(_var)IMAGE_SQUASHFS_DISABLE_4K_ALIGN)" >> $(2))
$(if $(PRODUCT_$(_var)_BASE_FS_PATH),$(hide) echo "$(1)_base_fs_file=$(PRODUCT_$(_var)_BASE_FS_PATH)" >> $(2))
$(eval _size := $(BOARD_$(_var)IMAGE_PARTITION_SIZE))
$(eval _reserved := $(BOARD_$(_var)IMAGE_PARTITION_RESERVED_SIZE))
$(eval _headroom := $(PRODUCT_$(_var)_HEADROOM))
$(if $(or $(_size), $(_reserved), $(_headroom)),,
    $(hide) echo "$(1)_disable_sparse=true" >> $(2))
$(call add-common-flags-to-image-props,$(1),$(2))
endef

# $(1): the path of the output dictionary file
<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
# $(2): additional "key=value" pairs to append to the dictionary file.
define generate-userimage-prop-dictionary
=======
# $(2): a subset of "system vendor cache userdata product system_ext oem odm vendor_dlkm odm_dlkm"
# $(3): additional "key=value" pairs to append to the dictionary file.
define generate-image-prop-dictionary
$(if $(filter $(2),system),\
    $(if $(INTERNAL_SYSTEM_OTHER_PARTITION_SIZE),$(hide) echo "system_other_size=$(INTERNAL_SYSTEM_OTHER_PARTITION_SIZE)" >> $(1))
    $(if $(PRODUCT_SYSTEM_HEADROOM),$(hide) echo "system_headroom=$(PRODUCT_SYSTEM_HEADROOM)" >> $(1))
    $(call add-common-ro-flags-to-image-props,system,$(1))
)
$(if $(filter $(2),system_other),\
    $(hide) echo "building_system_other_image=$(BUILDING_SYSTEM_OTHER_IMAGE)" >> $(1)
    $(if $(INTERNAL_SYSTEM_OTHER_PARTITION_SIZE),,
        $(hide) echo "system_other_disable_sparse=true" >> $(1))
)
$(if $(filter $(2),userdata),\
    $(if $(BOARD_USERDATAIMAGE_FILE_SYSTEM_TYPE),$(hide) echo "userdata_fs_type=$(BOARD_USERDATAIMAGE_FILE_SYSTEM_TYPE)" >> $(1))
    $(if $(BOARD_USERDATAIMAGE_PARTITION_SIZE),$(hide) echo "userdata_size=$(BOARD_USERDATAIMAGE_PARTITION_SIZE)" >> $(1))
    $(if $(PRODUCT_FS_CASEFOLD),$(hide) echo "needs_casefold=$(PRODUCT_FS_CASEFOLD)" >> $(1))
    $(if $(PRODUCT_QUOTA_PROJID),$(hide) echo "needs_projid=$(PRODUCT_QUOTA_PROJID)" >> $(1))
    $(if $(PRODUCT_FS_COMPRESSION),$(hide) echo "needs_compress=$(PRODUCT_FS_COMPRESSION)" >> $(1))
    $(call add-common-flags-to-image-props,userdata,$(1))
)
$(if $(filter $(2),cache),\
    $(if $(BOARD_CACHEIMAGE_FILE_SYSTEM_TYPE),$(hide) echo "cache_fs_type=$(BOARD_CACHEIMAGE_FILE_SYSTEM_TYPE)" >> $(1))
    $(if $(BOARD_CACHEIMAGE_PARTITION_SIZE),$(hide) echo "cache_size=$(BOARD_CACHEIMAGE_PARTITION_SIZE)" >> $(1))
    $(call add-common-flags-to-image-props,cache,$(1))
)
$(if $(filter $(2),vendor),\
    $(call add-common-ro-flags-to-image-props,vendor,$(1))
)
$(if $(filter $(2),product),\
    $(call add-common-ro-flags-to-image-props,product,$(1))
)
$(if $(filter $(2),system_ext),\
    $(call add-common-ro-flags-to-image-props,system_ext,$(1))
)
$(if $(filter $(2),odm),\
    $(call add-common-ro-flags-to-image-props,odm,$(1))
)
$(if $(filter $(2),vendor_dlkm),\
    $(call add-common-ro-flags-to-image-props,vendor_dlkm,$(1))
)
$(if $(filter $(2),odm_dlkm),\
    $(call add-common-ro-flags-to-image-props,odm_dlkm,$(1))
)
$(if $(filter $(2),oem),\
    $(if $(BOARD_OEMIMAGE_PARTITION_SIZE),$(hide) echo "oem_size=$(BOARD_OEMIMAGE_PARTITION_SIZE)" >> $(1))
    $(if $(BOARD_OEMIMAGE_JOURNAL_SIZE),$(hide) echo "oem_journal_size=$(BOARD_OEMIMAGE_JOURNAL_SIZE)" >> $(1))
    $(if $(BOARD_OEMIMAGE_EXTFS_INODE_COUNT),$(hide) echo "oem_extfs_inode_count=$(BOARD_OEMIMAGE_EXTFS_INODE_COUNT)" >> $(1))
    $(if $(BOARD_OEMIMAGE_EXTFS_RSV_PCT),$(hide) echo "oem_extfs_rsv_pct=$(BOARD_OEMIMAGE_EXTFS_RSV_PCT)" >> $(1))
    $(call add-common-flags-to-image-props,oem,$(1))
)
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
$(hide) echo "ext_mkuserimg=$(notdir $(MKEXTUSERIMG))" >> $(1)
$(if $(INTERNAL_USERIMAGES_EXT_VARIANT),$(hide) echo "fs_type=$(INTERNAL_USERIMAGES_EXT_VARIANT)" >> $(1))
<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
$(if $(BOARD_SYSTEMIMAGE_PARTITION_SIZE),$(hide) echo "system_size=$(BOARD_SYSTEMIMAGE_PARTITION_SIZE)" >> $(1))
$(if $(BOARD_SYSTEMIMAGE_FILE_SYSTEM_TYPE),$(hide) echo "system_fs_type=$(BOARD_SYSTEMIMAGE_FILE_SYSTEM_TYPE)" >> $(1))
$(if $(BOARD_SYSTEMIMAGE_EXTFS_INODE_COUNT),$(hide) echo "system_extfs_inode_count=$(BOARD_SYSTEMIMAGE_EXTFS_INODE_COUNT)" >> $(1))
$(if $(BOARD_SYSTEMIMAGE_EXTFS_RSV_PCT),$(hide) echo "system_extfs_rsv_pct=$(BOARD_SYSTEMIMAGE_EXTFS_RSV_PCT)" >> $(1))
$(if $(BOARD_SYSTEMIMAGE_JOURNAL_SIZE),$(hide) echo "system_journal_size=$(BOARD_SYSTEMIMAGE_JOURNAL_SIZE)" >> $(1))
=======
$(if $(INTERNAL_USERIMAGES_SPARSE_EXT_FLAG),$(hide) echo "extfs_sparse_flag=$(INTERNAL_USERIMAGES_SPARSE_EXT_FLAG)" >> $(1))
$(if $(INTERNAL_USERIMAGES_SPARSE_EROFS_FLAG),$(hide) echo "erofs_sparse_flag=$(INTERNAL_USERIMAGES_SPARSE_EROFS_FLAG)" >> $(1))
$(if $(INTERNAL_USERIMAGES_SPARSE_SQUASHFS_FLAG),$(hide) echo "squashfs_sparse_flag=$(INTERNAL_USERIMAGES_SPARSE_SQUASHFS_FLAG)" >> $(1))
$(if $(INTERNAL_USERIMAGES_SPARSE_F2FS_FLAG),$(hide) echo "f2fs_sparse_flag=$(INTERNAL_USERIMAGES_SPARSE_F2FS_FLAG)" >> $(1))
$(if $(BOARD_EROFS_COMPRESSOR),$(hide) echo "erofs_default_compressor=$(BOARD_EROFS_COMPRESSOR)" >> $(1))
$(if $(BOARD_EROFS_PCLUSTER_SIZE),$(hide) echo "erofs_pcluster_size=$(BOARD_EROFS_PCLUSTER_SIZE)" >> $(1))
$(if $(BOARD_EROFS_SHARE_DUP_BLOCKS),$(hide) echo "erofs_share_dup_blocks=$(BOARD_EROFS_SHARE_DUP_BLOCKS)" >> $(1))
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
$(if $(BOARD_EXT4_SHARE_DUP_BLOCKS),$(hide) echo "ext4_share_dup_blocks=$(BOARD_EXT4_SHARE_DUP_BLOCKS)" >> $(1))
$(if $(BOARD_SYSTEMIMAGE_SQUASHFS_COMPRESSOR),$(hide) echo "system_squashfs_compressor=$(BOARD_SYSTEMIMAGE_SQUASHFS_COMPRESSOR)" >> $(1))
$(if $(BOARD_SYSTEMIMAGE_SQUASHFS_COMPRESSOR_OPT),$(hide) echo "system_squashfs_compressor_opt=$(BOARD_SYSTEMIMAGE_SQUASHFS_COMPRESSOR_OPT)" >> $(1))
$(if $(BOARD_SYSTEMIMAGE_SQUASHFS_BLOCK_SIZE),$(hide) echo "system_squashfs_block_size=$(BOARD_SYSTEMIMAGE_SQUASHFS_BLOCK_SIZE)" >> $(1))
$(if $(BOARD_SYSTEMIMAGE_SQUASHFS_DISABLE_4K_ALIGN),$(hide) echo "system_squashfs_disable_4k_align=$(BOARD_SYSTEMIMAGE_SQUASHFS_DISABLE_4K_ALIGN)" >> $(1))
$(if $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_SYSTEM_BASE_FS_PATH),$(hide) echo "system_base_fs_file=$(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_SYSTEM_BASE_FS_PATH)" >> $(1))
$(if $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_SYSTEM_HEADROOM),$(hide) echo "system_headroom=$(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_SYSTEM_HEADROOM)" >> $(1))
$(if $(BOARD_USERDATAIMAGE_FILE_SYSTEM_TYPE),$(hide) echo "userdata_fs_type=$(BOARD_USERDATAIMAGE_FILE_SYSTEM_TYPE)" >> $(1))
$(if $(BOARD_USERDATAIMAGE_PARTITION_SIZE),$(hide) echo "userdata_size=$(BOARD_USERDATAIMAGE_PARTITION_SIZE)" >> $(1))
$(if $(BOARD_FLASH_LOGICAL_BLOCK_SIZE), $(hide) echo "flash_logical_block_size=$(BOARD_FLASH_LOGICAL_BLOCK_SIZE)" >> $(1))
$(if $(BOARD_FLASH_ERASE_BLOCK_SIZE), $(hide) echo "flash_erase_block_size=$(BOARD_FLASH_ERASE_BLOCK_SIZE)" >> $(1))
$(if $(BOARD_CACHEIMAGE_FILE_SYSTEM_TYPE),$(hide) echo "cache_fs_type=$(BOARD_CACHEIMAGE_FILE_SYSTEM_TYPE)" >> $(1))
$(if $(BOARD_CACHEIMAGE_PARTITION_SIZE),$(hide) echo "cache_size=$(BOARD_CACHEIMAGE_PARTITION_SIZE)" >> $(1))
$(if $(BOARD_VENDORIMAGE_FILE_SYSTEM_TYPE),$(hide) echo "vendor_fs_type=$(BOARD_VENDORIMAGE_FILE_SYSTEM_TYPE)" >> $(1))
$(if $(BOARD_VENDORIMAGE_EXTFS_INODE_COUNT),$(hide) echo "vendor_extfs_inode_count=$(BOARD_VENDORIMAGE_EXTFS_INODE_COUNT)" >> $(1))
$(if $(BOARD_VENDORIMAGE_EXTFS_RSV_PCT),$(hide) echo "vendor_extfs_rsv_pct=$(BOARD_VENDORIMAGE_EXTFS_RSV_PCT)" >> $(1))
$(if $(BOARD_VENDORIMAGE_PARTITION_SIZE),$(hide) echo "vendor_size=$(BOARD_VENDORIMAGE_PARTITION_SIZE)" >> $(1))
$(if $(BOARD_VENDORIMAGE_JOURNAL_SIZE),$(hide) echo "vendor_journal_size=$(BOARD_VENDORIMAGE_JOURNAL_SIZE)" >> $(1))
$(if $(BOARD_VENDORIMAGE_SQUASHFS_COMPRESSOR),$(hide) echo "vendor_squashfs_compressor=$(BOARD_VENDORIMAGE_SQUASHFS_COMPRESSOR)" >> $(1))
$(if $(BOARD_VENDORIMAGE_SQUASHFS_COMPRESSOR_OPT),$(hide) echo "vendor_squashfs_compressor_opt=$(BOARD_VENDORIMAGE_SQUASHFS_COMPRESSOR_OPT)" >> $(1))
$(if $(BOARD_VENDORIMAGE_SQUASHFS_BLOCK_SIZE),$(hide) echo "vendor_squashfs_block_size=$(BOARD_VENDORIMAGE_SQUASHFS_BLOCK_SIZE)" >> $(1))
$(if $(BOARD_VENDORIMAGE_SQUASHFS_DISABLE_4K_ALIGN),$(hide) echo "vendor_squashfs_disable_4k_align=$(BOARD_VENDORIMAGE_SQUASHFS_DISABLE_4K_ALIGN)" >> $(1))
$(if $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_VENDOR_BASE_FS_PATH),$(hide) echo "vendor_base_fs_file=$(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_VENDOR_BASE_FS_PATH)" >> $(1))
$(if $(BOARD_PRODUCTIMAGE_FILE_SYSTEM_TYPE),$(hide) echo "product_fs_type=$(BOARD_PRODUCTIMAGE_FILE_SYSTEM_TYPE)" >> $(1))
$(if $(BOARD_PRODUCTIMAGE_EXTFS_INODE_COUNT),$(hide) echo "product_extfs_inode_count=$(BOARD_PRODUCTIMAGE_EXTFS_INODE_COUNT)" >> $(1))
$(if $(BOARD_PRODUCTIMAGE_EXTFS_RSV_PCT),$(hide) echo "product_extfs_rsv_pct=$(BOARD_PRODUCTIMAGE_EXTFS_RSV_PCT)" >> $(1))
$(if $(BOARD_PRODUCTIMAGE_PARTITION_SIZE),$(hide) echo "product_size=$(BOARD_PRODUCTIMAGE_PARTITION_SIZE)" >> $(1))
$(if $(BOARD_PRODUCTIMAGE_JOURNAL_SIZE),$(hide) echo "product_journal_size=$(BOARD_PRODUCTIMAGE_JOURNAL_SIZE)" >> $(1))
$(if $(BOARD_PRODUCTIMAGE_SQUASHFS_COMPRESSOR),$(hide) echo "product_squashfs_compressor=$(BOARD_PRODUCTIMAGE_SQUASHFS_COMPRESSOR)" >> $(1))
$(if $(BOARD_PRODUCTIMAGE_SQUASHFS_COMPRESSOR_OPT),$(hide) echo "product_squashfs_compressor_opt=$(BOARD_PRODUCTIMAGE_SQUASHFS_COMPRESSOR_OPT)" >> $(1))
$(if $(BOARD_PRODUCTIMAGE_SQUASHFS_BLOCK_SIZE),$(hide) echo "product_squashfs_block_size=$(BOARD_PRODUCTIMAGE_SQUASHFS_BLOCK_SIZE)" >> $(1))
$(if $(BOARD_PRODUCTIMAGE_SQUASHFS_DISABLE_4K_ALIGN),$(hide) echo "product_squashfs_disable_4k_align=$(BOARD_PRODUCTIMAGE_SQUASHFS_DISABLE_4K_ALIGN)" >> $(1))
$(if $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_PRODUCT_BASE_FS_PATH),$(hide) echo "product_base_fs_file=$(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_PRODUCT_BASE_FS_PATH)" >> $(1))
$(if $(BOARD_OEMIMAGE_PARTITION_SIZE),$(hide) echo "oem_size=$(BOARD_OEMIMAGE_PARTITION_SIZE)" >> $(1))
$(if $(BOARD_OEMIMAGE_JOURNAL_SIZE),$(hide) echo "oem_journal_size=$(BOARD_OEMIMAGE_JOURNAL_SIZE)" >> $(1))
$(if $(BOARD_OEMIMAGE_EXTFS_INODE_COUNT),$(hide) echo "oem_extfs_inode_count=$(BOARD_OEMIMAGE_EXTFS_INODE_COUNT)" >> $(1))
$(if $(BOARD_OEMIMAGE_EXTFS_RSV_PCT),$(hide) echo "oem_extfs_rsv_pct=$(BOARD_OEMIMAGE_EXTFS_RSV_PCT)" >> $(1))
$(if $(INTERNAL_USERIMAGES_SPARSE_EXT_FLAG),$(hide) echo "extfs_sparse_flag=$(INTERNAL_USERIMAGES_SPARSE_EXT_FLAG)" >> $(1))
$(if $(INTERNAL_USERIMAGES_SPARSE_SQUASHFS_FLAG),$(hide) echo "squashfs_sparse_flag=$(INTERNAL_USERIMAGES_SPARSE_SQUASHFS_FLAG)" >> $(1))
$(hide) echo "selinux_fc=$(SELINUX_FC)" >> $(1)
$(if $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_SUPPORTS_BOOT_SIGNER),$(hide) echo "boot_signer=$(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_SUPPORTS_BOOT_SIGNER)" >> $(1))
$(if $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_SUPPORTS_VERITY),$(hide) echo "verity=$(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_SUPPORTS_VERITY)" >> $(1))
$(if $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_SUPPORTS_VERITY),$(hide) echo "verity_key=$(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_VERITY_SIGNING_KEY)" >> $(1))
$(if $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_SUPPORTS_VERITY),$(hide) echo "verity_signer_cmd=$(notdir $(VERITY_SIGNER))" >> $(1))
$(if $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_SUPPORTS_VERITY_FEC),$(hide) echo "verity_fec=$(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_SUPPORTS_VERITY_FEC)" >> $(1))
$(if $(filter eng, $(TARGET_BUILD_VARIANT)),$(hide) echo "verity_disable=true" >> $(1))
$(if $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_SYSTEM_VERITY_PARTITION),$(hide) echo "system_verity_block_device=$(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_SYSTEM_VERITY_PARTITION)" >> $(1))
$(if $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_VENDOR_VERITY_PARTITION),$(hide) echo "vendor_verity_block_device=$(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_VENDOR_VERITY_PARTITION)" >> $(1))
$(if $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_PRODUCT_VERITY_PARTITION),$(hide) echo "product_verity_block_device=$(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_PRODUCT_VERITY_PARTITION)" >> $(1))
$(if $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_SUPPORTS_VBOOT),$(hide) echo "vboot=$(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_SUPPORTS_VBOOT)" >> $(1))
$(if $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_SUPPORTS_VBOOT),$(hide) echo "vboot_key=$(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_VBOOT_SIGNING_KEY)" >> $(1))
$(if $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_SUPPORTS_VBOOT),$(hide) echo "vboot_subkey=$(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_VBOOT_SIGNING_SUBKEY)" >> $(1))
$(if $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_SUPPORTS_VBOOT),$(hide) echo "futility=$(notdir $(FUTILITY))" >> $(1))
$(if $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_SUPPORTS_VBOOT),$(hide) echo "vboot_signer_cmd=$(VBOOT_SIGNER)" >> $(1))
$(if $(BOARD_AVB_ENABLE),$(hide) echo "avb_avbtool=$(notdir $(AVBTOOL))" >> $(1))
$(if $(BOARD_AVB_ENABLE),$(hide) echo "avb_system_hashtree_enable=$(BOARD_AVB_ENABLE)" >> $(1))
$(if $(BOARD_AVB_ENABLE),$(hide) echo "avb_system_add_hashtree_footer_args=$(BOARD_AVB_SYSTEM_ADD_HASHTREE_FOOTER_ARGS)" >> $(1))
$(if $(BOARD_AVB_ENABLE),\
    $(if $(BOARD_AVB_SYSTEM_KEY_PATH),\
        $(hide) echo "avb_system_key_path=$(BOARD_AVB_SYSTEM_KEY_PATH)" >> $(1)
        $(hide) echo "avb_system_algorithm=$(BOARD_AVB_SYSTEM_ALGORITHM)" >> $(1)
        $(hide) echo "avb_system_rollback_index_location=$(BOARD_AVB_SYSTEM_ROLLBACK_INDEX_LOCATION)" >> $(1)))
$(if $(BOARD_AVB_ENABLE),$(hide) echo "avb_vendor_hashtree_enable=$(BOARD_AVB_ENABLE)" >> $(1))
$(if $(BOARD_AVB_ENABLE),$(hide) echo "avb_vendor_add_hashtree_footer_args=$(BOARD_AVB_VENDOR_ADD_HASHTREE_FOOTER_ARGS)" >> $(1))
$(if $(BOARD_AVB_ENABLE),\
    $(if $(BOARD_AVB_VENDOR_KEY_PATH),\
        $(hide) echo "avb_vendor_key_path=$(BOARD_AVB_VENDOR_KEY_PATH)" >> $(1)
        $(hide) echo "avb_vendor_algorithm=$(BOARD_AVB_VENDOR_ALGORITHM)" >> $(1)
        $(hide) echo "avb_vendor_rollback_index_location=$(BOARD_AVB_VENDOR_ROLLBACK_INDEX_LOCATION)" >> $(1)))
$(if $(BOARD_AVB_ENABLE),$(hide) echo "avb_product_hashtree_enable=$(BOARD_AVB_ENABLE)" >> $(1))
$(if $(BOARD_AVB_ENABLE),$(hide) echo "avb_product_add_hashtree_footer_args=$(BOARD_AVB_PRODUCT_ADD_HASHTREE_FOOTER_ARGS)" >> $(1))
$(if $(BOARD_AVB_ENABLE),\
    $(if $(BOARD_AVB_PRODUCT_KEY_PATH),\
        $(hide) echo "avb_product_key_path=$(BOARD_AVB_PRODUCT_KEY_PATH)" >> $(1)
        $(hide) echo "avb_product_algorithm=$(BOARD_AVB_PRODUCT_ALGORITHM)" >> $(1)
        $(hide) echo "avb_product_rollback_index_location=$(BOARD_AVB_PRODUCT_ROLLBACK_INDEX_LOCATION)" >> $(1)))
$(if $(filter true,$(BOARD_USES_RECOVERY_AS_BOOT)),\
    $(hide) echo "recovery_as_boot=true" >> $(1))
$(if $(filter true,$(BOARD_BUILD_SYSTEM_ROOT_IMAGE)),\
    $(hide) echo "system_root_image=true" >> $(1);\
    echo "ramdisk_dir=$(TARGET_ROOT_OUT)" >> $(1))
$(if $(2),$(hide) $(foreach kv,$(2),echo "$(kv)" >> $(1);))
endef

# $(1): modules list
# $(2): output dir
# $(3): mount point
# $(4): staging dir
# Depmod requires a well-formed kernel version so 0.0 is used as a placeholder.
define build-image-kernel-modules
    $(hide) rm -rf $(2)/lib/modules
    $(hide) mkdir -p $(2)/lib/modules
    $(hide) cp $(1) $(2)/lib/modules/
    $(hide) rm -rf $(4)
    $(hide) mkdir -p $(4)/lib/modules/0.0/$(3)lib/modules
    $(hide) cp $(1) $(4)/lib/modules/0.0/$(3)lib/modules
    $(hide) $(DEPMOD) -b $(4) 0.0
    $(hide) sed -e 's/\(.*modules.*\):/\/\1:/g' -e 's/ \([^ ]*modules[^ ]*\)/ \/\1/g' $(4)/lib/modules/0.0/modules.dep > $(2)/lib/modules/modules.dep
    $(hide) cp $(4)/lib/modules/0.0/modules.alias $(2)/lib/modules
endef

# -----------------------------------------------------------------
# Recovery image

# Recovery image exists if we are building recovery, or building recovery as boot.
ifneq (,$(INSTALLED_RECOVERYIMAGE_TARGET)$(filter true,$(BOARD_USES_RECOVERY_AS_BOOT)))

INTERNAL_RECOVERYIMAGE_FILES := $(filter $(TARGET_RECOVERY_OUT)/%, \
    $(ALL_DEFAULT_INSTALLED_MODULES))

<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
recovery_initrc := $(call include-path-for, recovery)/etc/init.rc
=======
INSTALLED_FILES_FILE_RECOVERY := $(PRODUCT_OUT)/installed-files-recovery.txt
INSTALLED_FILES_JSON_RECOVERY := $(INSTALLED_FILES_FILE_RECOVERY:.txt=.json)

ifeq ($(BOARD_USES_RECOVERY_AS_BOOT),true)
INSTALLED_BOOTIMAGE_TARGET := $(BUILT_BOOTIMAGE_TARGET)
endif

# TODO(b/30414428): Can't depend on INTERNAL_RECOVERYIMAGE_FILES alone like other
# INSTALLED_FILES_FILE_* rules. Because currently there're cp/rsync/rm commands in
# build-recoveryimage-target, which would touch the files under TARGET_RECOVERY_OUT and race with
# the call to FILELIST.
$(INSTALLED_FILES_FILE_RECOVERY): $(INTERNAL_RECOVERY_RAMDISK_FILES_TIMESTAMP)

$(INSTALLED_FILES_FILE_RECOVERY): .KATI_IMPLICIT_OUTPUTS := $(INSTALLED_FILES_JSON_RECOVERY)
$(INSTALLED_FILES_FILE_RECOVERY): $(INTERNAL_RECOVERYIMAGE_FILES) $(FILESLIST) $(FILESLIST_UTIL)
	@echo Installed file list: $@
	mkdir -p $(dir $@)
	rm -f $@
	$(FILESLIST) $(TARGET_RECOVERY_ROOT_OUT) > $(@:.txt=.json)
	$(FILESLIST_UTIL) -c $(@:.txt=.json) > $@

>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
recovery_sepolicy := \
    $(TARGET_RECOVERY_ROOT_OUT)/sepolicy \
    $(TARGET_RECOVERY_ROOT_OUT)/plat_file_contexts \
    $(TARGET_RECOVERY_ROOT_OUT)/vendor_file_contexts \
    $(TARGET_RECOVERY_ROOT_OUT)/plat_property_contexts \
    $(TARGET_RECOVERY_ROOT_OUT)/vendor_property_contexts

ifdef BOARD_ODM_SEPOLICY_DIRS
recovery_sepolicy += \
    $(TARGET_RECOVERY_ROOT_OUT)/odm_file_contexts \
    $(TARGET_RECOVERY_ROOT_OUT)/odm_property_contexts
endif

# Passed into rsync from non-recovery root to recovery root, to avoid overwriting recovery-specific
# SELinux files
IGNORE_RECOVERY_SEPOLICY := $(patsubst $(TARGET_RECOVERY_OUT)/%,--exclude=/%,$(recovery_sepolicy))

recovery_kernel := $(INSTALLED_KERNEL_TARGET) # same as a non-recovery system
recovery_ramdisk := $(PRODUCT_OUT)/ramdisk-recovery.img
recovery_build_props := $(intermediate_system_build_prop)
ifdef property_overrides_split_enabled
recovery_build_props += $(INSTALLED_VENDOR_BUILD_PROP_TARGET)
endif
ifdef BOARD_USES_PRODUCTIMAGE
recovery_build_props += $(INSTALLED_PRODUCT_BUILD_PROP_TARGET)
endif
recovery_resources_common := $(call include-path-for, recovery)/res

# Set recovery_density to the density bucket of the device.
recovery_density := unknown
ifneq (,$(PRODUCT_AAPT_PREF_CONFIG))
# If PRODUCT_AAPT_PREF_CONFIG includes a dpi bucket, then use that value.
recovery_density := $(filter %dpi,$(PRODUCT_AAPT_PREF_CONFIG))
else
# Otherwise, use the default medium density.
recovery_densities := mdpi
endif

ifneq (,$(wildcard $(recovery_resources_common)-$(recovery_density)))
recovery_resources_common := $(recovery_resources_common)-$(recovery_density)
else
recovery_resources_common := $(recovery_resources_common)-xhdpi
endif

# Select the 18x32 font on high-density devices (xhdpi and up); and
# the 12x22 font on other devices.  Note that the font selected here
# can be overridden for a particular device by putting a font.png in
# its private recovery resources.

ifneq (,$(filter xxxhdpi 560dpi xxhdpi 400dpi xhdpi,$(recovery_density)))
recovery_font := $(call include-path-for, recovery)/fonts/18x32.png
else
recovery_font := $(call include-path-for, recovery)/fonts/12x22.png
endif

ifndef TARGET_PRIVATE_RES_DIRS
TARGET_PRIVATE_RES_DIRS := $(wildcard $(TARGET_DEVICE_DIR)/recovery/res)
endif
recovery_resource_deps := $(shell find $(recovery_resources_common) \
  $(TARGET_PRIVATE_RES_DIRS) -type f)
ifdef TARGET_RECOVERY_FSTAB
recovery_fstab := $(TARGET_RECOVERY_FSTAB)
else
recovery_fstab := $(strip $(wildcard $(TARGET_DEVICE_DIR)/recovery.fstab))
endif
ifdef TARGET_RECOVERY_WIPE
recovery_wipe := $(TARGET_RECOVERY_WIPE)
else
recovery_wipe :=
endif

# Traditionally with non-A/B OTA we have:
#   boot.img + recovery-from-boot.p + recovery-resource.dat = recovery.img.
# recovery-resource.dat is needed only if we carry an imgdiff patch of the boot and recovery images
# and invoke install-recovery.sh on the first boot post an OTA update.
#
# We no longer need that if one of the following conditions holds:
#   a) We carry a full copy of the recovery image - no patching needed
#      (BOARD_USES_FULL_RECOVERY_IMAGE = true);
#   b) We build a single image that contains boot and recovery both - no recovery image to install
#      (BOARD_USES_RECOVERY_AS_BOOT = true);
#   c) We build the root into system image - not needing the resource file as we do bsdiff
#      (BOARD_BUILD_SYSTEM_ROOT_IMAGE = true).
#   d) We include the recovery DTBO image within recovery - not needing the resource file as we
#      do bsdiff because boot and recovery will contain different number of entries
#      (BOARD_INCLUDE_RECOVERY_DTBO = true).
# Note that condition b) implies condition c), because of the earlier check in this file:
# "BOARD_USES_RECOVERY_AS_BOOT = true must have BOARD_BUILD_SYSTEM_ROOT_IMAGE = true" (not vice
# versa though).

ifeq (,$(filter true, $(BOARD_USES_FULL_RECOVERY_IMAGE) $(BOARD_BUILD_SYSTEM_ROOT_IMAGE) \
  $(BOARD_INCLUDE_RECOVERY_DTBO)))
# Named '.dat' so we don't attempt to use imgdiff for patching it.
RECOVERY_RESOURCE_ZIP := $(TARGET_OUT)/etc/recovery-resource.dat
else
RECOVERY_RESOURCE_ZIP :=
endif

INTERNAL_RECOVERYIMAGE_ARGS := \
	$(addprefix --second ,$(INSTALLED_2NDBOOTLOADER_TARGET)) \
	--kernel $(recovery_kernel) \
	--ramdisk $(recovery_ramdisk)

# Assumes this has already been stripped
ifdef INTERNAL_KERNEL_CMDLINE
  INTERNAL_RECOVERYIMAGE_ARGS += --cmdline "$(INTERNAL_KERNEL_CMDLINE)"
endif
ifdef BOARD_KERNEL_BASE
  INTERNAL_RECOVERYIMAGE_ARGS += --base $(BOARD_KERNEL_BASE)
endif
ifdef BOARD_KERNEL_PAGESIZE
  INTERNAL_RECOVERYIMAGE_ARGS += --pagesize $(BOARD_KERNEL_PAGESIZE)
endif
ifdef BOARD_INCLUDE_RECOVERY_DTBO
  INTERNAL_RECOVERYIMAGE_ARGS += --recovery_dtbo $(BOARD_PREBUILT_DTBOIMAGE)
endif

# Keys authorized to sign OTA packages this build will accept.  The
# build always uses dev-keys for this; release packaging tools will
# substitute other keys for this one.
OTA_PUBLIC_KEYS := $(DEFAULT_SYSTEM_DEV_CERTIFICATE).x509.pem

# Generate a file containing the keys that will be read by the
# recovery binary.
RECOVERY_INSTALL_OTA_KEYS := \
	$(call intermediates-dir-for,PACKAGING,ota_keys)/keys
DUMPKEY_JAR := $(HOST_OUT_JAVA_LIBRARIES)/dumpkey.jar
$(RECOVERY_INSTALL_OTA_KEYS): PRIVATE_OTA_PUBLIC_KEYS := $(OTA_PUBLIC_KEYS)
$(RECOVERY_INSTALL_OTA_KEYS): extra_keys := $(patsubst %,%.x509.pem,$(PRODUCT_EXTRA_RECOVERY_KEYS))
$(RECOVERY_INSTALL_OTA_KEYS): $(OTA_PUBLIC_KEYS) $(DUMPKEY_JAR) $(extra_keys)
	@echo "DumpPublicKey: $@ <= $(PRIVATE_OTA_PUBLIC_KEYS) $(extra_keys)"
	@rm -rf $@
	@mkdir -p $(dir $@)
	$(JAVA) -jar $(DUMPKEY_JAR) $(PRIVATE_OTA_PUBLIC_KEYS) $(extra_keys) > $@

RECOVERYIMAGE_ID_FILE := $(PRODUCT_OUT)/recovery.id

# $(1): output file
define build-recoveryimage-target
<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
  # Making recovery image
  $(hide) mkdir -p $(TARGET_RECOVERY_OUT)
  $(hide) mkdir -p $(TARGET_RECOVERY_ROOT_OUT)/etc $(TARGET_RECOVERY_ROOT_OUT)/sdcard $(TARGET_RECOVERY_ROOT_OUT)/tmp
  # Copying baseline ramdisk...
  # Use rsync because "cp -Rf" fails to overwrite broken symlinks on Mac.
  $(hide) rsync -a --exclude=etc --exclude=sdcard $(IGNORE_RECOVERY_SEPOLICY) $(IGNORE_CACHE_LINK) $(TARGET_ROOT_OUT) $(TARGET_RECOVERY_OUT)
  # Copy adbd from system/bin to recovery/root/sbin
  $(hide) cp -f $(TARGET_OUT_EXECUTABLES)/adbd $(TARGET_RECOVERY_ROOT_OUT)/sbin/adbd
  # Modifying ramdisk contents...
  $(if $(BOARD_RECOVERY_KERNEL_MODULES), \
    $(call build-image-kernel-modules,$(BOARD_RECOVERY_KERNEL_MODULES),$(TARGET_RECOVERY_ROOT_OUT),,$(call intermediates-dir-for,PACKAGING,depmod_recovery)))
  # Removes $(TARGET_RECOVERY_ROOT_OUT)/init*.rc EXCEPT init.recovery*.rc.
  $(hide) find $(TARGET_RECOVERY_ROOT_OUT) -maxdepth 1 -name 'init*.rc' -type f -not -name "init.recovery.*.rc" | xargs rm -f
  $(hide) cp -f $(recovery_initrc) $(TARGET_RECOVERY_ROOT_OUT)/
  $(hide) cp $(TARGET_ROOT_OUT)/init.recovery.*.rc $(TARGET_RECOVERY_ROOT_OUT)/ || true # Ignore error when the src file doesn't exist.
  $(hide) mkdir -p $(TARGET_RECOVERY_ROOT_OUT)/res
  $(hide) rm -rf $(TARGET_RECOVERY_ROOT_OUT)/res/*
  $(hide) cp -rf $(recovery_resources_common)/* $(TARGET_RECOVERY_ROOT_OUT)/res
  $(hide) cp -f $(recovery_font) $(TARGET_RECOVERY_ROOT_OUT)/res/images/font.png
  $(hide) $(foreach item,$(TARGET_PRIVATE_RES_DIRS), \
    cp -rf $(item) $(TARGET_RECOVERY_ROOT_OUT)/$(newline))
  $(hide) $(foreach item,$(recovery_fstab), \
    cp -f $(item) $(TARGET_RECOVERY_ROOT_OUT)/etc/recovery.fstab)
  $(if $(strip $(recovery_wipe)), \
    $(hide) cp -f $(recovery_wipe) $(TARGET_RECOVERY_ROOT_OUT)/etc/recovery.wipe)
  $(hide) cp $(RECOVERY_INSTALL_OTA_KEYS) $(TARGET_RECOVERY_ROOT_OUT)/res/keys
  $(hide) cat $(INSTALLED_DEFAULT_PROP_TARGET) \
          > $(TARGET_RECOVERY_ROOT_OUT)/prop.default
  $(if $(INSTALLED_VENDOR_DEFAULT_PROP_TARGET), \
    $(hide) cat $(INSTALLED_VENDOR_DEFAULT_PROP_TARGET) \
            >> $(TARGET_RECOVERY_ROOT_OUT)/prop.default)
  $(hide) cat $(recovery_build_props) \
          >> $(TARGET_RECOVERY_ROOT_OUT)/prop.default
  $(hide) ln -sf prop.default $(TARGET_RECOVERY_ROOT_OUT)/default.prop
  $(BOARD_RECOVERY_IMAGE_PREPARE)
  $(if $(filter true,$(BOARD_BUILD_SYSTEM_ROOT_IMAGE)), \
    $(hide) mkdir -p $(TARGET_RECOVERY_ROOT_OUT)/system_root; \
            rm -rf $(TARGET_RECOVERY_ROOT_OUT)/system; \
            ln -sf /system_root/system $(TARGET_RECOVERY_ROOT_OUT)/system) # Mount the system_root_image to /system_root and symlink /system.
  $(hide) $(MKBOOTFS) -d $(TARGET_OUT) $(TARGET_RECOVERY_ROOT_OUT) | $(MINIGZIP) > $(recovery_ramdisk)
  $(if $(filter true,$(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_SUPPORTS_VBOOT)), \
    $(hide) $(MKBOOTIMG) $(INTERNAL_RECOVERYIMAGE_ARGS) $(INTERNAL_MKBOOTIMG_VERSION_ARGS) $(BOARD_MKBOOTIMG_ARGS) --output $(1).unsigned, \
    $(hide) $(MKBOOTIMG) $(INTERNAL_RECOVERYIMAGE_ARGS) $(INTERNAL_MKBOOTIMG_VERSION_ARGS) $(BOARD_MKBOOTIMG_ARGS) --output $(1) --id > $(RECOVERYIMAGE_ID_FILE))
  $(if $(filter true,$(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_SUPPORTS_BOOT_SIGNER)),\
=======
  $(if $(filter true,$(PRODUCT_SUPPORTS_VBOOT)), \
    $(MKBOOTIMG) $(if $(strip $(2)),--kernel $(strip $(2))) $(INTERNAL_RECOVERYIMAGE_ARGS) \
                 $(INTERNAL_MKBOOTIMG_VERSION_ARGS) $(BOARD_RECOVERY_MKBOOTIMG_ARGS) \
                 --output $(1).unsigned, \
    $(MKBOOTIMG) $(if $(strip $(2)),--kernel $(strip $(2))) $(INTERNAL_RECOVERYIMAGE_ARGS) \
                 $(INTERNAL_MKBOOTIMG_VERSION_ARGS) $(INTERNAL_MKBOOTIMG_GKI_SINGING_ARGS) \
                 $(BOARD_RECOVERY_MKBOOTIMG_ARGS) --output $(1))
  $(if $(filter true,$(PRODUCT_SUPPORTS_BOOT_SIGNER)),\
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
    $(if $(filter true,$(BOARD_USES_RECOVERY_AS_BOOT)),\
      $(BOOT_SIGNER) /boot $(1) $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_VERITY_SIGNING_KEY).pk8 $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_VERITY_SIGNING_KEY).x509.pem $(1),\
      $(BOOT_SIGNER) /recovery $(1) $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_VERITY_SIGNING_KEY).pk8 $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_VERITY_SIGNING_KEY).x509.pem $(1)\
    )\
  )
  $(if $(filter true,$(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_SUPPORTS_VBOOT)), \
    $(VBOOT_SIGNER) $(FUTILITY) $(1).unsigned $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_VBOOT_SIGNING_KEY).vbpubk $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_VBOOT_SIGNING_KEY).vbprivk $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_VBOOT_SIGNING_SUBKEY).vbprivk $(1).keyblock $(1))
  $(if $(filter true,$(BOARD_USES_RECOVERY_AS_BOOT)), \
    $(hide) $(call assert-max-image-size,$(1),$(call get-hash-image-max-size,$(BOARD_BOOTIMAGE_PARTITION_SIZE))), \
    $(hide) $(call assert-max-image-size,$(1),$(call get-hash-image-max-size,$(BOARD_RECOVERYIMAGE_PARTITION_SIZE))))
  $(if $(filter true,$(BOARD_AVB_ENABLE)), \
    $(if $(filter true,$(BOARD_USES_RECOVERY_AS_BOOT)), \
      $(hide) $(AVBTOOL) add_hash_footer --image $(1) --partition_size $(BOARD_BOOTIMAGE_PARTITION_SIZE) --partition_name boot $(INTERNAL_AVB_BOOT_SIGNING_ARGS) $(BOARD_AVB_BOOT_ADD_HASH_FOOTER_ARGS),\
      $(hide) $(AVBTOOL) add_hash_footer --image $(1) --partition_size $(BOARD_RECOVERYIMAGE_PARTITION_SIZE) --partition_name recovery $(INTERNAL_AVB_RECOVERY_SIGNING_ARGS) $(BOARD_AVB_RECOVERY_ADD_HASH_FOOTER_ARGS)))
endef

<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
ADBD := $(TARGET_OUT_EXECUTABLES)/adbd
=======
recoveryimage-deps := $(MKBOOTIMG) $(recovery_ramdisk) $(recovery_kernel)
ifeq (true,$(PRODUCT_SUPPORTS_BOOT_SIGNER))
  recoveryimage-deps += $(BOOT_SIGNER)
endif
ifeq (true,$(PRODUCT_SUPPORTS_VBOOT))
  recoveryimage-deps += $(VBOOT_SIGNER)
endif
ifeq (true,$(BOARD_AVB_ENABLE))
  recoveryimage-deps += $(AVBTOOL) $(BOARD_AVB_BOOT_KEY_PATH)
endif
ifdef BOARD_GKI_SIGNING_KEY_PATH
  recoveryimage-deps += $(BOARD_GKI_SIGNING_KEY_PATH) $(AVBTOOL)
endif
ifdef BOARD_INCLUDE_RECOVERY_DTBO
  ifdef BOARD_PREBUILT_RECOVERY_DTBOIMAGE
    recoveryimage-deps += $(BOARD_PREBUILT_RECOVERY_DTBOIMAGE)
  else
    recoveryimage-deps += $(BOARD_PREBUILT_DTBOIMAGE)
  endif
endif
ifdef BOARD_INCLUDE_RECOVERY_ACPIO
  recoveryimage-deps += $(BOARD_RECOVERY_ACPIO)
endif
ifdef BOARD_INCLUDE_DTB_IN_BOOTIMG
  recoveryimage-deps += $(INSTALLED_DTBIMAGE_TARGET)
endif
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)

ifeq ($(BOARD_USES_RECOVERY_AS_BOOT),true)
ifeq (true,$(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_SUPPORTS_BOOT_SIGNER))
$(INSTALLED_BOOTIMAGE_TARGET) : $(BOOT_SIGNER)
endif
ifeq (true,$(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_SUPPORTS_VBOOT))
$(INSTALLED_BOOTIMAGE_TARGET) : $(VBOOT_SIGNER)
endif
ifeq (true,$(BOARD_AVB_ENABLE))
$(INSTALLED_BOOTIMAGE_TARGET) : $(AVBTOOL) $(BOARD_AVB_BOOT_KEY_PATH)
endif
$(INSTALLED_BOOTIMAGE_TARGET): $(MKBOOTFS) $(MKBOOTIMG) $(MINIGZIP) $(ADBD) \
		$(INSTALLED_RAMDISK_TARGET) \
		$(INTERNAL_RECOVERYIMAGE_FILES) \
		$(recovery_initrc) $(recovery_sepolicy) $(recovery_kernel) \
		$(INSTALLED_2NDBOOTLOADER_TARGET) \
		$(recovery_build_props) $(recovery_resource_deps) \
		$(recovery_fstab) \
		$(RECOVERY_INSTALL_OTA_KEYS) \
		$(INSTALLED_VENDOR_DEFAULT_PROP_TARGET) \
		$(BOARD_RECOVERY_KERNEL_MODULES) \
		$(DEPMOD)
		$(call pretty,"Target boot image from recovery: $@")
		$(call build-recoveryimage-target, $@)
endif

$(INSTALLED_RECOVERYIMAGE_TARGET): $(MKBOOTFS) $(MKBOOTIMG) $(MINIGZIP) $(ADBD) \
		$(INSTALLED_RAMDISK_TARGET) \
		$(INSTALLED_BOOTIMAGE_TARGET) \
		$(INTERNAL_RECOVERYIMAGE_FILES) \
		$(recovery_initrc) $(recovery_sepolicy) $(recovery_kernel) \
		$(INSTALLED_2NDBOOTLOADER_TARGET) \
		$(recovery_build_props) $(recovery_resource_deps) \
		$(recovery_fstab) \
		$(RECOVERY_INSTALL_OTA_KEYS) \
		$(INSTALLED_VENDOR_DEFAULT_PROP_TARGET) \
		$(BOARD_RECOVERY_KERNEL_MODULES) \
		$(DEPMOD)
		$(call build-recoveryimage-target, $@)

ifdef RECOVERY_RESOURCE_ZIP
$(RECOVERY_RESOURCE_ZIP): $(INSTALLED_RECOVERYIMAGE_TARGET) | $(ZIPTIME)
	$(hide) mkdir -p $(dir $@)
	$(hide) find $(TARGET_RECOVERY_ROOT_OUT)/res -type f | sort | zip -0qrjX $@ -@
	$(remove-timestamps-from-package)
endif

.PHONY: recoveryimage-nodeps
recoveryimage-nodeps:
	@echo "make $@: ignoring dependencies"
	$(call build-recoveryimage-target, $(INSTALLED_RECOVERYIMAGE_TARGET))

else # INSTALLED_RECOVERYIMAGE_TARGET not defined
RECOVERY_RESOURCE_ZIP :=
endif

.PHONY: recoveryimage
recoveryimage: $(INSTALLED_RECOVERYIMAGE_TARGET) $(RECOVERY_RESOURCE_ZIP)

ifneq ($(BOARD_NAND_PAGE_SIZE),)
$(error MTD device is no longer supported and thus BOARD_NAND_PAGE_SIZE is deprecated.)
endif

ifneq ($(BOARD_NAND_SPARE_SIZE),)
$(error MTD device is no longer supported and thus BOARD_NAND_SPARE_SIZE is deprecated.)
endif

<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
=======

# -----------------------------------------------------------------
# Build debug ramdisk and debug boot image.
ifneq ($(BUILDING_DEBUG_BOOT_IMAGE)$(BUILDING_DEBUG_VENDOR_BOOT_IMAGE),)

INTERNAL_DEBUG_RAMDISK_FILES := $(filter $(TARGET_DEBUG_RAMDISK_OUT)/%, \
    $(ALL_GENERATED_SOURCES) \
    $(ALL_DEFAULT_INSTALLED_MODULES))

# Directories to be picked into the debug ramdisk.
# As these directories are all merged into one single cpio archive, the order
# matters. If there are multiple files with the same pathname, then the last one
# wins.
#
# ramdisk-debug.img will merge the content from either ramdisk.img or
# ramdisk-recovery.img, depending on whether BOARD_USES_RECOVERY_AS_BOOT
# is set or not.
ifeq ($(BOARD_USES_RECOVERY_AS_BOOT),true)
  INTERNAL_DEBUG_RAMDISK_SRC_DIRS := $(TARGET_RECOVERY_ROOT_OUT)
  INTERNAL_DEBUG_RAMDISK_SRC_RAMDISK_TARGET := $(recovery_ramdisk)
else  # BOARD_USES_RECOVERY_AS_BOOT == true
  INTERNAL_DEBUG_RAMDISK_SRC_DIRS := $(TARGET_RAMDISK_OUT)
  INTERNAL_DEBUG_RAMDISK_SRC_RAMDISK_TARGET := $(INSTALLED_RAMDISK_TARGET)
endif # BOARD_USES_RECOVERY_AS_BOOT != true

INTERNAL_DEBUG_RAMDISK_SRC_DIRS += $(TARGET_DEBUG_RAMDISK_OUT)
INTERNAL_DEBUG_RAMDISK_SRC_DEPS := $(INTERNAL_DEBUG_RAMDISK_SRC_RAMDISK_TARGET) $(INTERNAL_DEBUG_RAMDISK_FILES)

# INSTALLED_FILES_FILE_DEBUG_RAMDISK would ensure TARGET_DEBUG_RAMDISK_OUT is created.
INSTALLED_FILES_FILE_DEBUG_RAMDISK := $(PRODUCT_OUT)/installed-files-ramdisk-debug.txt
INSTALLED_FILES_JSON_DEBUG_RAMDISK := $(INSTALLED_FILES_FILE_DEBUG_RAMDISK:.txt=.json)
$(INSTALLED_FILES_FILE_DEBUG_RAMDISK): .KATI_IMPLICIT_OUTPUTS := $(INSTALLED_FILES_JSON_DEBUG_RAMDISK)
$(INSTALLED_FILES_FILE_DEBUG_RAMDISK): $(INTERNAL_DEBUG_RAMDISK_SRC_DEPS)
$(INSTALLED_FILES_FILE_DEBUG_RAMDISK): $(FILESLIST) $(FILESLIST_UTIL)
	@echo "Installed file list: $@"
	$(hide) rm -f $@
	$(hide) mkdir -p $(dir $@) $(TARGET_DEBUG_RAMDISK_OUT)
	touch $(TARGET_DEBUG_RAMDISK_OUT)/force_debuggable
	$(FILESLIST) $(INTERNAL_DEBUG_RAMDISK_SRC_DIRS) > $(@:.txt=.json)
	$(FILESLIST_UTIL) -c $(@:.txt=.json) > $@

ifdef BUILDING_DEBUG_BOOT_IMAGE

# -----------------------------------------------------------------
# the debug ramdisk, which is the original ramdisk plus additional
# files: force_debuggable, adb_debug.prop and userdebug sepolicy.
# When /force_debuggable is present, /init will load userdebug sepolicy
# and property files to allow adb root, if the device is unlocked.
INSTALLED_DEBUG_RAMDISK_TARGET := $(PRODUCT_OUT)/ramdisk-debug.img

$(INSTALLED_DEBUG_RAMDISK_TARGET): $(INSTALLED_FILES_FILE_DEBUG_RAMDISK)
$(INSTALLED_DEBUG_RAMDISK_TARGET): $(MKBOOTFS) | $(COMPRESSION_COMMAND_DEPS)
	@echo "Target debug ramdisk: $@"
	$(hide) rm -f $@
	$(hide) mkdir -p $(dir $@)
	$(MKBOOTFS) -d $(TARGET_OUT) $(INTERNAL_DEBUG_RAMDISK_SRC_DIRS) | $(COMPRESSION_COMMAND) > $@

.PHONY: ramdisk_debug-nodeps
ramdisk_debug-nodeps: $(MKBOOTFS) | $(COMPRESSION_COMMAND_DEPS)
	@echo "make $@: ignoring dependencies"
	$(hide) rm -f $(INSTALLED_DEBUG_RAMDISK_TARGET)
	$(hide) mkdir -p $(dir $(INSTALLED_DEBUG_RAMDISK_TARGET)) $(INTERNAL_DEBUG_RAMDISK_SRC_DIRS)
	$(MKBOOTFS) -d $(TARGET_OUT) $(INTERNAL_DEBUG_RAMDISK_SRC_DIRS) | $(COMPRESSION_COMMAND) > $(INSTALLED_DEBUG_RAMDISK_TARGET)

# -----------------------------------------------------------------
# the boot-debug.img, which is the kernel plus ramdisk-debug.img
#
# Note: it's intentional to skip signing for boot-debug.img, because it
# can only be used if the device is unlocked with verification error.
ifneq ($(strip $(BOARD_KERNEL_BINARIES)),)
  INSTALLED_DEBUG_BOOTIMAGE_TARGET := $(foreach k,$(subst kernel,boot-debug,$(BOARD_KERNEL_BINARIES)), \
         $(PRODUCT_OUT)/$(k).img)
else
  INSTALLED_DEBUG_BOOTIMAGE_TARGET := $(PRODUCT_OUT)/boot-debug.img
endif

# Replace ramdisk.img in $(MKBOOTIMG) ARGS with ramdisk-debug.img to build boot-debug.img
$(INSTALLED_DEBUG_BOOTIMAGE_TARGET): $(INSTALLED_DEBUG_RAMDISK_TARGET)
ifeq ($(BOARD_USES_RECOVERY_AS_BOOT),true)
  INTERNAL_DEBUG_BOOTIMAGE_ARGS := $(subst $(INTERNAL_DEBUG_RAMDISK_SRC_RAMDISK_TARGET),$(INSTALLED_DEBUG_RAMDISK_TARGET),$(INTERNAL_RECOVERYIMAGE_ARGS))
else
  INTERNAL_DEBUG_BOOTIMAGE_ARGS := $(subst $(INTERNAL_DEBUG_RAMDISK_SRC_RAMDISK_TARGET),$(INSTALLED_DEBUG_RAMDISK_TARGET),$(INTERNAL_BOOTIMAGE_ARGS))
endif

# If boot.img is chained but boot-debug.img is not signed, libavb in bootloader
# will fail to find valid AVB metadata from the end of /boot, thus stop booting.
# Using a test key to sign boot-debug.img to continue booting with the mismatched
# public key, if the device is unlocked.
ifneq ($(BOARD_AVB_BOOT_KEY_PATH),)
$(INSTALLED_DEBUG_BOOTIMAGE_TARGET): $(AVBTOOL) $(BOARD_AVB_BOOT_TEST_KEY_PATH)
endif

BOARD_AVB_BOOT_TEST_KEY_PATH := external/avb/test/data/testkey_rsa2048.pem
INTERNAL_AVB_BOOT_TEST_SIGNING_ARGS := --algorithm SHA256_RSA2048 --key $(BOARD_AVB_BOOT_TEST_KEY_PATH)
# $(1): the bootimage to sign
# $(2): boot image variant (boot, boot-debug, boot-test-harness)
define test-key-sign-bootimage
$(call assert-max-image-size,$(1),$(call get-hash-image-max-size,$(call get-bootimage-partition-size,$(1),$(2))))
$(AVBTOOL) add_hash_footer \
  --image $(1) \
  --partition_size $(call get-bootimage-partition-size,$(1),$(2))\
  --partition_name boot $(INTERNAL_AVB_BOOT_TEST_SIGNING_ARGS) \
  $(BOARD_AVB_BOOT_ADD_HASH_FOOTER_ARGS)
$(call assert-max-image-size,$(1),$(call get-bootimage-partition-size,$(1),$(2)))
endef

# $(1): output file
define build-debug-bootimage-target
  $(MKBOOTIMG) --kernel $(PRODUCT_OUT)/$(subst .img,,$(subst boot-debug,kernel,$(notdir $(1)))) \
    $(INTERNAL_DEBUG_BOOTIMAGE_ARGS) $(INTERNAL_MKBOOTIMG_VERSION_ARGS) \
    $(INTERNAL_MKBOOTIMG_GKI_SINGING_ARGS) $(BOARD_MKBOOTIMG_ARGS) --output $1
  $(if $(BOARD_AVB_BOOT_KEY_PATH),$(call test-key-sign-bootimage,$1,boot-debug))
endef

# Depends on original boot.img and ramdisk-debug.img, to build the new boot-debug.img
$(INSTALLED_DEBUG_BOOTIMAGE_TARGET): $(MKBOOTIMG) $(INSTALLED_BOOTIMAGE_TARGET) $(BOARD_GKI_SIGNING_KEY_PATH) $(AVBTOOL)
	$(call pretty,"Target boot debug image: $@")
	$(call build-debug-bootimage-target, $@)

.PHONY: bootimage_debug-nodeps
bootimage_debug-nodeps: $(MKBOOTIMG) $(BOARD_GKI_SIGNING_KEY_PATH) $(AVBTOOL)
	echo "make $@: ignoring dependencies"
	$(foreach b,$(INSTALLED_DEBUG_BOOTIMAGE_TARGET),$(call build-debug-bootimage-target,$b))

endif # BUILDING_DEBUG_BOOT_IMAGE

# -----------------------------------------------------------------
# vendor debug ramdisk
# Combines vendor ramdisk files and debug ramdisk files to build the vendor debug ramdisk.
ifdef BUILDING_DEBUG_VENDOR_BOOT_IMAGE

INTERNAL_VENDOR_DEBUG_RAMDISK_FILES := $(filter $(TARGET_VENDOR_DEBUG_RAMDISK_OUT)/%, \
    $(ALL_GENERATED_SOURCES) \
    $(ALL_DEFAULT_INSTALLED_MODULES))

# The debug vendor ramdisk combines vendor ramdisk and debug ramdisk.
INTERNAL_DEBUG_VENDOR_RAMDISK_SRC_DIRS := $(TARGET_VENDOR_RAMDISK_OUT)
INTERNAL_DEBUG_VENDOR_RAMDISK_SRC_DEPS := $(INTERNAL_VENDOR_RAMDISK_TARGET)

# Exclude recovery files in the default vendor ramdisk if including a standalone
# recovery ramdisk in vendor_boot.
ifeq (true,$(BOARD_MOVE_RECOVERY_RESOURCES_TO_VENDOR_BOOT))
  ifneq (true,$(BOARD_INCLUDE_RECOVERY_RAMDISK_IN_VENDOR_BOOT))
    INTERNAL_DEBUG_VENDOR_RAMDISK_SRC_DIRS += $(TARGET_RECOVERY_ROOT_OUT)
    INTERNAL_DEBUG_VENDOR_RAMDISK_SRC_DEPS += $(INTERNAL_RECOVERY_RAMDISK_FILES_TIMESTAMP)
  endif # BOARD_INCLUDE_RECOVERY_RAMDISK_IN_VENDOR_BOOT != true
endif # BOARD_MOVE_RECOVERY_RESOURCES_TO_VENDOR_BOOT == true

INTERNAL_DEBUG_VENDOR_RAMDISK_SRC_DIRS += $(TARGET_VENDOR_DEBUG_RAMDISK_OUT) $(TARGET_DEBUG_RAMDISK_OUT)
INTERNAL_DEBUG_VENDOR_RAMDISK_SRC_DEPS += $(INTERNAL_VENDOR_DEBUG_RAMDISK_FILES) $(INSTALLED_FILES_FILE_DEBUG_RAMDISK)

# INSTALLED_FILES_FILE_VENDOR_DEBUG_RAMDISK would ensure TARGET_VENDOR_DEBUG_RAMDISK_OUT is created.
INSTALLED_FILES_FILE_VENDOR_DEBUG_RAMDISK := $(PRODUCT_OUT)/installed-files-vendor-ramdisk-debug.txt
INSTALLED_FILES_JSON_VENDOR_DEBUG_RAMDISK := $(INSTALLED_FILES_FILE_VENDOR_DEBUG_RAMDISK:.txt=.json)
$(INSTALLED_FILES_FILE_VENDOR_DEBUG_RAMDISK): .KATI_IMPLICIT_OUTPUTS := $(INSTALLED_FILES_JSON_VENDOR_DEBUG_RAMDISK)
$(INSTALLED_FILES_FILE_VENDOR_DEBUG_RAMDISK): $(INTERNAL_DEBUG_VENDOR_RAMDISK_SRC_DEPS)
$(INSTALLED_FILES_FILE_VENDOR_DEBUG_RAMDISK): $(FILESLIST) $(FILESLIST_UTIL)
	@echo "Installed file list: $@"
	$(hide) rm -f $@
	$(hide) mkdir -p $(dir $@) $(TARGET_VENDOR_DEBUG_RAMDISK_OUT)
	$(FILESLIST) $(INTERNAL_DEBUG_VENDOR_RAMDISK_SRC_DIRS) > $(@:.txt=.json)
	$(FILESLIST_UTIL) -c $(@:.txt=.json) > $@

INTERNAL_VENDOR_DEBUG_RAMDISK_TARGET := $(call intermediates-dir-for,PACKAGING,vendor_boot-debug)/vendor_ramdisk-debug.cpio$(RAMDISK_EXT)

$(INTERNAL_VENDOR_DEBUG_RAMDISK_TARGET): $(INSTALLED_FILES_FILE_VENDOR_DEBUG_RAMDISK)
$(INTERNAL_VENDOR_DEBUG_RAMDISK_TARGET): $(MKBOOTFS) | $(COMPRESSION_COMMAND_DEPS)
	$(hide) rm -f $@
	$(hide) mkdir -p $(dir $@)
	$(MKBOOTFS) -d $(TARGET_OUT) $(INTERNAL_DEBUG_VENDOR_RAMDISK_SRC_DIRS) | $(COMPRESSION_COMMAND) > $@

INSTALLED_VENDOR_DEBUG_RAMDISK_TARGET := $(PRODUCT_OUT)/vendor_ramdisk-debug.img
$(INSTALLED_VENDOR_DEBUG_RAMDISK_TARGET): $(INTERNAL_VENDOR_DEBUG_RAMDISK_TARGET)
	@echo "Target debug vendor ramdisk: $@"
	$(copy-file-to-target)

# -----------------------------------------------------------------
# vendor_boot-debug.img.
INSTALLED_VENDOR_DEBUG_BOOTIMAGE_TARGET := $(PRODUCT_OUT)/vendor_boot-debug.img

# The util to sign vendor_boot-debug.img with a test key.
BOARD_AVB_VENDOR_BOOT_TEST_KEY_PATH := external/avb/test/data/testkey_rsa2048.pem
INTERNAL_AVB_VENDOR_BOOT_TEST_SIGNING_ARGS := --algorithm SHA256_RSA2048 --key $(BOARD_AVB_VENDOR_BOOT_TEST_KEY_PATH)
# $(1): the vendor bootimage to sign
define test-key-sign-vendor-bootimage
$(call assert-max-image-size,$(1),$(call get-hash-image-max-size,$(BOARD_VENDOR_BOOTIMAGE_PARTITION_SIZE)))
$(AVBTOOL) add_hash_footer \
  --image $(1) \
  --partition_size $(BOARD_VENDOR_BOOTIMAGE_PARTITION_SIZE) \
  --partition_name vendor_boot $(INTERNAL_AVB_VENDOR_BOOT_TEST_SIGNING_ARGS) \
  $(BOARD_AVB_VENDOR_BOOT_ADD_HASH_FOOTER_ARGS)
$(call assert-max-image-size,$(1),$(BOARD_VENDOR_BOOTIMAGE_PARTITION_SIZE))
endef

ifneq ($(BOARD_AVB_VENDOR_BOOT_KEY_PATH),)
$(INSTALLED_VENDOR_DEBUG_BOOTIMAGE_TARGET): $(AVBTOOL) $(BOARD_AVB_VENDOR_BOOT_TEST_KEY_PATH)
endif

# Depends on vendor_boot.img and vendor-ramdisk-debug.cpio.gz to build the new vendor_boot-debug.img
$(INSTALLED_VENDOR_DEBUG_BOOTIMAGE_TARGET): $(MKBOOTIMG) $(INSTALLED_VENDOR_BOOTIMAGE_TARGET) $(INTERNAL_VENDOR_DEBUG_RAMDISK_TARGET)
$(INSTALLED_VENDOR_DEBUG_BOOTIMAGE_TARGET): $(INTERNAL_VENDOR_RAMDISK_FRAGMENT_TARGETS)
	$(call pretty,"Target vendor_boot debug image: $@")
	$(MKBOOTIMG) $(INTERNAL_VENDOR_BOOTIMAGE_ARGS) $(BOARD_MKBOOTIMG_ARGS) --vendor_ramdisk $(INTERNAL_VENDOR_DEBUG_RAMDISK_TARGET) $(INTERNAL_VENDOR_RAMDISK_FRAGMENT_ARGS) --vendor_boot $@
	$(call assert-max-image-size,$@,$(BOARD_VENDOR_BOOTIMAGE_PARTITION_SIZE))
	$(if $(BOARD_AVB_VENDOR_BOOT_KEY_PATH),$(call test-key-sign-vendor-bootimage,$@))

endif # BUILDING_DEBUG_VENDOR_BOOT_IMAGE

# Appends a few test harness specific properties into the adb_debug.prop.
ADDITIONAL_TEST_HARNESS_PROPERTIES := ro.audio.silent=1
ADDITIONAL_TEST_HARNESS_PROPERTIES += ro.test_harness=1

INTERNAL_DEBUG_RAMDISK_ADB_DEBUG_PROP_TARGET := $(strip $(filter $(TARGET_DEBUG_RAMDISK_OUT)/adb_debug.prop,$(INTERNAL_DEBUG_RAMDISK_FILES)))
INTERNAL_TEST_HARNESS_RAMDISK_ADB_DEBUG_PROP_TARGET := $(TARGET_TEST_HARNESS_RAMDISK_OUT)/adb_debug.prop
$(INTERNAL_TEST_HARNESS_RAMDISK_ADB_DEBUG_PROP_TARGET): $(INTERNAL_DEBUG_RAMDISK_ADB_DEBUG_PROP_TARGET)
	$(hide) rm -f $@
	$(hide) mkdir -p $(dir $@)
ifdef INTERNAL_DEBUG_RAMDISK_ADB_DEBUG_PROP_TARGET
	$(hide) cp $(INTERNAL_DEBUG_RAMDISK_ADB_DEBUG_PROP_TARGET) $@
endif
	$(hide) echo "" >> $@
	$(hide) echo "#" >> $@
	$(hide) echo "# ADDITIONAL TEST HARNESS PROPERTIES" >> $@
	$(hide) echo "#" >> $@
	$(hide) $(foreach line,$(ADDITIONAL_TEST_HARNESS_PROPERTIES), \
	          echo "$(line)" >> $@;)

INTERNAL_TEST_HARNESS_RAMDISK_FILES := $(filter $(TARGET_TEST_HARNESS_RAMDISK_OUT)/%, \
    $(INTERNAL_TEST_HARNESS_RAMDISK_ADB_DEBUG_PROP_TARGET) \
    $(ALL_GENERATED_SOURCES) \
    $(ALL_DEFAULT_INSTALLED_MODULES))

# The order is important here. The test harness ramdisk staging directory has to
# come last so that it can override the adb_debug.prop in the debug ramdisk
# staging directory.
INTERNAL_TEST_HARNESS_RAMDISK_SRC_DIRS := $(INTERNAL_DEBUG_RAMDISK_SRC_DIRS) $(TARGET_TEST_HARNESS_RAMDISK_OUT)
INTERNAL_TEST_HARNESS_RAMDISK_SRC_DEPS := $(INSTALLED_FILES_FILE_DEBUG_RAMDISK) $(INTERNAL_TEST_HARNESS_RAMDISK_FILES)

ifdef BUILDING_DEBUG_BOOT_IMAGE

# -----------------------------------------------------------------
# The test harness ramdisk, which is based off debug_ramdisk, plus a
# few additional test-harness-specific properties in adb_debug.prop.
INSTALLED_TEST_HARNESS_RAMDISK_TARGET := $(PRODUCT_OUT)/ramdisk-test-harness.img

$(INSTALLED_TEST_HARNESS_RAMDISK_TARGET): $(INTERNAL_TEST_HARNESS_RAMDISK_SRC_DEPS)
$(INSTALLED_TEST_HARNESS_RAMDISK_TARGET): $(MKBOOTFS) | $(COMPRESSION_COMMAND_DEPS)
	@echo "Target test harness ramdisk: $@"
	$(hide) rm -f $@
	$(hide) mkdir -p $(dir $@)
	$(MKBOOTFS) -d $(TARGET_OUT) $(INTERNAL_TEST_HARNESS_RAMDISK_SRC_DIRS) | $(COMPRESSION_COMMAND) > $@

.PHONY: ramdisk_test_harness-nodeps
ramdisk_test_harness-nodeps: $(MKBOOTFS) | $(COMPRESSION_COMMAND_DEPS)
	@echo "make $@: ignoring dependencies"
	$(hide) rm -f $(INSTALLED_TEST_HARNESS_RAMDISK_TARGET)
	$(hide) mkdir -p $(dir $(INSTALLED_TEST_HARNESS_RAMDISK_TARGET)) $(INTERNAL_TEST_HARNESS_RAMDISK_SRC_DIRS)
	$(MKBOOTFS) -d $(TARGET_OUT) $(INTERNAL_TEST_HARNESS_RAMDISK_SRC_DIRS) | $(COMPRESSION_COMMAND) > $(INSTALLED_TEST_HARNESS_RAMDISK_TARGET)

# -----------------------------------------------------------------
# the boot-test-harness.img, which is the kernel plus ramdisk-test-harness.img
#
# Note: it's intentional to skip signing for boot-test-harness.img, because it
# can only be used if the device is unlocked with verification error.
ifneq ($(strip $(BOARD_KERNEL_BINARIES)),)
  INSTALLED_TEST_HARNESS_BOOTIMAGE_TARGET := $(foreach k,$(subst kernel,boot-test-harness,$(BOARD_KERNEL_BINARIES)), \
    $(PRODUCT_OUT)/$(k).img)
else
  INSTALLED_TEST_HARNESS_BOOTIMAGE_TARGET := $(PRODUCT_OUT)/boot-test-harness.img
endif

# Replace ramdisk-debug.img in $(MKBOOTIMG) ARGS with ramdisk-test-harness.img to build boot-test-harness.img
$(INSTALLED_TEST_HARNESS_BOOTIMAGE_TARGET): $(INSTALLED_TEST_HARNESS_RAMDISK_TARGET)
INTERNAL_TEST_HARNESS_BOOTIMAGE_ARGS := $(subst $(INSTALLED_DEBUG_RAMDISK_TARGET),$(INSTALLED_TEST_HARNESS_RAMDISK_TARGET),$(INTERNAL_DEBUG_BOOTIMAGE_ARGS))

# If boot.img is chained but boot-test-harness.img is not signed, libavb in bootloader
# will fail to find valid AVB metadata from the end of /boot, thus stop booting.
# Using a test key to sign boot-test-harness.img to continue booting with the mismatched
# public key, if the device is unlocked.
ifneq ($(BOARD_AVB_BOOT_KEY_PATH),)
$(INSTALLED_TEST_HARNESS_BOOTIMAGE_TARGET): $(AVBTOOL) $(BOARD_AVB_BOOT_TEST_KEY_PATH)
endif

# $(1): output file
define build-boot-test-harness-target
  $(MKBOOTIMG) --kernel $(PRODUCT_OUT)/$(subst .img,,$(subst boot-test-harness,kernel,$(notdir $(1)))) \
    $(INTERNAL_TEST_HARNESS_BOOTIMAGE_ARGS) $(INTERNAL_MKBOOTIMG_VERSION_ARGS) \
    $(INTERNAL_MKBOOTIMG_GKI_SINGING_ARGS) $(BOARD_MKBOOTIMG_ARGS) --output $@
  $(if $(BOARD_AVB_BOOT_KEY_PATH),$(call test-key-sign-bootimage,$@,boot-test-harness))
endef

# Build the new boot-test-harness.img, based on boot-debug.img and ramdisk-test-harness.img.
$(INSTALLED_TEST_HARNESS_BOOTIMAGE_TARGET): $(MKBOOTIMG) $(INSTALLED_DEBUG_BOOTIMAGE_TARGET) $(BOARD_GKI_SIGNING_KEY_PATH) $(AVBTOOL)
	$(call pretty,"Target boot test harness image: $@")
	$(call build-boot-test-harness-target,$@)

.PHONY: bootimage_test_harness-nodeps
bootimage_test_harness-nodeps: $(MKBOOTIMG) $(BOARD_GKI_SIGNING_KEY_PATH) $(AVBTOOL)
	echo "make $@: ignoring dependencies"
	$(foreach b,$(INSTALLED_TEST_HARNESS_BOOTIMAGE_TARGET),$(call build-boot-test-harness-target,$b))

endif # BUILDING_DEBUG_BOOT_IMAGE

# -----------------------------------------------------------------
# vendor test harness ramdisk, which is a vendor ramdisk combined with
# a test harness ramdisk.
ifdef BUILDING_DEBUG_VENDOR_BOOT_IMAGE

INTERNAL_VENDOR_TEST_HARNESS_RAMDISK_TARGET := $(call intermediates-dir-for,PACKAGING,vendor_boot-test-harness)/vendor_ramdisk-test-harness.cpio$(RAMDISK_EXT)

# The order is important here. The test harness ramdisk staging directory has to
# come last so that it can override the adb_debug.prop in the debug ramdisk
# staging directory.
INTERNAL_TEST_HARNESS_VENDOR_RAMDISK_SRC_DIRS := $(INTERNAL_DEBUG_VENDOR_RAMDISK_SRC_DIRS) $(TARGET_TEST_HARNESS_RAMDISK_OUT)
INTERNAL_TEST_HARNESS_VENDOR_RAMDISK_SRC_DEPS := $(INSTALLED_FILES_FILE_VENDOR_DEBUG_RAMDISK) $(INTERNAL_TEST_HARNESS_RAMDISK_FILES)

$(INTERNAL_VENDOR_TEST_HARNESS_RAMDISK_TARGET): $(INTERNAL_TEST_HARNESS_VENDOR_RAMDISK_SRC_DEPS)
$(INTERNAL_VENDOR_TEST_HARNESS_RAMDISK_TARGET): $(MKBOOTFS) | $(COMPRESSION_COMMAND_DEPS)
	$(hide) rm -f $@
	$(hide) mkdir -p $(dir $@)
	$(MKBOOTFS) -d $(TARGET_OUT) $(INTERNAL_TEST_HARNESS_VENDOR_RAMDISK_SRC_DIRS) | $(COMPRESSION_COMMAND) > $@

INSTALLED_VENDOR_TEST_HARNESS_RAMDISK_TARGET := $(PRODUCT_OUT)/vendor_ramdisk-test-harness.img
$(INSTALLED_VENDOR_TEST_HARNESS_RAMDISK_TARGET): $(INTERNAL_VENDOR_TEST_HARNESS_RAMDISK_TARGET)
	@echo "Target test harness vendor ramdisk: $@"
	$(copy-file-to-target)

# -----------------------------------------------------------------
# vendor_boot-test-harness.img.
INSTALLED_VENDOR_TEST_HARNESS_BOOTIMAGE_TARGET := $(PRODUCT_OUT)/vendor_boot-test-harness.img

ifneq ($(BOARD_AVB_VENDOR_BOOT_KEY_PATH),)
$(INSTALLED_VENDOR_TEST_HARNESS_BOOTIMAGE_TARGET): $(AVBTOOL) $(BOARD_AVB_VENDOR_BOOT_TEST_KEY_PATH)
endif

# Depends on vendor_boot.img and vendor_ramdisk-test-harness.cpio$(RAMDISK_EXT) to build the new vendor_boot-test-harness.img
$(INSTALLED_VENDOR_TEST_HARNESS_BOOTIMAGE_TARGET): $(MKBOOTIMG) $(INSTALLED_VENDOR_BOOTIMAGE_TARGET)
$(INSTALLED_VENDOR_TEST_HARNESS_BOOTIMAGE_TARGET): $(INTERNAL_VENDOR_TEST_HARNESS_RAMDISK_TARGET)
$(INSTALLED_VENDOR_TEST_HARNESS_BOOTIMAGE_TARGET): $(INTERNAL_VENDOR_RAMDISK_FRAGMENT_TARGETS)
	$(call pretty,"Target vendor_boot test harness image: $@")
	$(MKBOOTIMG) $(INTERNAL_VENDOR_BOOTIMAGE_ARGS) $(BOARD_MKBOOTIMG_ARGS) --vendor_ramdisk $(INTERNAL_VENDOR_TEST_HARNESS_RAMDISK_TARGET) $(INTERNAL_VENDOR_RAMDISK_FRAGMENT_ARGS) --vendor_boot $@
	$(call assert-max-image-size,$@,$(BOARD_VENDOR_BOOTIMAGE_PARTITION_SIZE))
	$(if $(BOARD_AVB_VENDOR_BOOT_KEY_PATH),$(call test-key-sign-vendor-bootimage,$@))

endif # BUILDING_DEBUG_VENDOR_BOOT_IMAGE

endif # BUILDING_DEBUG_BOOT_IMAGE || BUILDING_DEBUG_VENDOR_BOOT_IMAGE


# Creates a compatibility symlink between two partitions, e.g. /system/vendor to /vendor
# $1: from location (e.g $(TARGET_OUT)/vendor)
# $2: destination location (e.g. /vendor)
# $3: partition image name (e.g. vendor.img)
define create-partition-compat-symlink
$(eval \
$1:
	@echo Symlink $(patsubst $(PRODUCT_OUT)/%,%,$1) to $2
	mkdir -p $(dir $1)
	if [ -d $1 ] && [ ! -h $1 ]; then \
	  echo 'Non-symlink $1 detected!' 1>&2; \
	  echo 'You cannot install files to $1 while building a separate $3!' 1>&2; \
	  exit 1; \
	fi
	ln -sfn $2 $1
$1: .KATI_SYMLINK_OUTPUTS := $1
)
$1
endef


>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
# -----------------------------------------------------------------
# system image
#
# Remove overridden packages from $(ALL_PDK_FUSION_FILES)
PDK_FUSION_SYSIMG_FILES := \
    $(filter-out $(foreach p,$(overridden_packages),$(p) %/$(p).apk), \
        $(ALL_PDK_FUSION_FILES))

INTERNAL_SYSTEMIMAGE_FILES := $(filter $(TARGET_OUT)/%, \
    $(ALL_GENERATED_SOURCES) \
    $(ALL_DEFAULT_INSTALLED_MODULES) \
    $(PDK_FUSION_SYSIMG_FILES) \
    $(RECOVERY_RESOURCE_ZIP)) \
    $(PDK_FUSION_SYMLINK_STAMP)

FULL_SYSTEMIMAGE_DEPS := $(INTERNAL_SYSTEMIMAGE_FILES) $(INTERNAL_USERIMAGES_DEPS)

# ASAN libraries in the system image - add dependency.
ASAN_IN_SYSTEM_INSTALLED := $(TARGET_OUT)/asan.tar.bz2
ifneq (,$(SANITIZE_TARGET))
  ifeq (true,$(SANITIZE_TARGET_SYSTEM))
    FULL_SYSTEMIMAGE_DEPS += $(ASAN_IN_SYSTEM_INSTALLED)
  endif
endif

# -----------------------------------------------------------------
# installed file list
# Depending on anything that $(BUILT_SYSTEMIMAGE) depends on.
# We put installed-files.txt ahead of image itself in the dependency graph
# so that we can get the size stat even if the build fails due to too large
# system image.
INSTALLED_FILES_FILE := $(PRODUCT_OUT)/installed-files.txt
$(INSTALLED_FILES_FILE): $(FULL_SYSTEMIMAGE_DEPS) $(FILESLIST)
	@echo Installed file list: $@
<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
	@mkdir -p $(dir $@)
	@rm -f $@
	$(hide) $(FILESLIST) $(TARGET_OUT) > $(@:.txt=.json)
	$(hide) build/make/tools/fileslist_util.py -c $(@:.txt=.json) > $@
=======
	mkdir -p $(dir $@)
	rm -f $@
	$(FILESLIST) $(TARGET_OUT) > $(@:.txt=.json)
	$(FILESLIST_UTIL) -c $(@:.txt=.json) > $@
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)

.PHONY: installed-file-list
installed-file-list: $(INSTALLED_FILES_FILE)

ifeq ($(HOST_OS),linux)
$(call dist-for-goals, sdk sdk_addon, $(INSTALLED_FILES_FILE))
endif

systemimage_intermediates := \
    $(call intermediates-dir-for,PACKAGING,systemimage)
BUILT_SYSTEMIMAGE := $(systemimage_intermediates)/system.img

# Create symlink /system/vendor to /vendor if necessary.
ifdef BOARD_USES_VENDORIMAGE
define create-system-vendor-symlink
$(hide) if [ -d $(TARGET_OUT)/vendor ] && [ ! -h $(TARGET_OUT)/vendor ]; then \
  echo 'Non-symlink $(TARGET_OUT)/vendor detected!' 1>&2; \
  echo 'You cannot install files to $(TARGET_OUT)/vendor while building a separate vendor.img!' 1>&2; \
  exit 1; \
fi
$(hide) ln -sf /vendor $(TARGET_OUT)/vendor
endef
else
define create-system-vendor-symlink
endef
endif

# Create symlink /system/product to /product if necessary.
ifdef BOARD_USES_PRODUCTIMAGE
define create-system-product-symlink
$(hide) if [ -d $(TARGET_OUT)/product ] && [ ! -h $(TARGET_OUT)/product ]; then \
  echo 'Non-symlink $(TARGET_OUT)/product detected!' 1>&2; \
  echo 'You cannot install files to $(TARGET_OUT)/product while building a separate product.img!' 1>&2; \
  exit 1; \
fi
$(hide) ln -sf /product $(TARGET_OUT)/product
endef
else
define create-system-product-symlink
endef
endif

# $(1): output file
define build-systemimage-target
  @echo "Target system fs image: $(1)"
  $(call create-system-vendor-symlink)
  $(call create-system-product-symlink)
  @mkdir -p $(dir $(1)) $(systemimage_intermediates) && rm -rf $(systemimage_intermediates)/system_image_info.txt
  $(call generate-userimage-prop-dictionary, $(systemimage_intermediates)/system_image_info.txt, \
      skip_fsck=true)
  $(hide) PATH=$(foreach p,$(INTERNAL_USERIMAGES_BINARY_PATHS),$(p):)$$PATH \
      build/make/tools/releasetools/build_image.py \
      $(TARGET_OUT) $(systemimage_intermediates)/system_image_info.txt $(1) $(TARGET_OUT) \
      || ( echo "Out of space? the tree size of $(TARGET_OUT) is (MB): " 1>&2 ;\
           du -sm $(TARGET_OUT) 1>&2;\
           if [ "$(INTERNAL_USERIMAGES_EXT_VARIANT)" == "ext4" ]; then \
               maxsize=$(BOARD_SYSTEMIMAGE_PARTITION_SIZE); \
               echo "The max is $$(( maxsize / 1048576 )) MB." 1>&2 ;\
           else \
               echo "The max is $$(( $(BOARD_SYSTEMIMAGE_PARTITION_SIZE) / 1048576 )) MB." 1>&2 ;\
           fi; \
           mkdir -p $(DIST_DIR); cp $(INSTALLED_FILES_FILE) $(DIST_DIR)/installed-files-rescued.txt; \
           exit 1 )
endef

$(BUILT_SYSTEMIMAGE): $(FULL_SYSTEMIMAGE_DEPS) $(INSTALLED_FILES_FILE) $(BUILD_IMAGE_SRCS)
	$(call build-systemimage-target,$@)

<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
INSTALLED_SYSTEMIMAGE := $(PRODUCT_OUT)/system.img
=======
$(call declare-1p-container,$(BUILT_SYSTEMIMAGE),system/extras)
$(call declare-container-license-deps,$(BUILT_SYSTEMIMAGE),$(FULL_SYSTEMIMAGE_DEPS),$(PRODUCT_OUT)/:)

INSTALLED_SYSTEMIMAGE_TARGET := $(PRODUCT_OUT)/system.img
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
SYSTEMIMAGE_SOURCE_DIR := $(TARGET_OUT)

# The system partition needs room for the recovery image as well.  We
# now store the recovery image as a binary patch using the boot image
# as the source (since they are very similar).  Generate the patch so
# we can see how big it's going to be, and include that in the system
# image size check calculation.
ifneq ($(INSTALLED_RECOVERYIMAGE_TARGET),)
ifneq ($(BOARD_USES_FULL_RECOVERY_IMAGE),true)
ifneq (,$(filter true, $(BOARD_BUILD_SYSTEM_ROOT_IMAGE) $(BOARD_INCLUDE_RECOVERY_DTBO)))
diff_tool := $(HOST_OUT_EXECUTABLES)/bsdiff
else
diff_tool := $(HOST_OUT_EXECUTABLES)/imgdiff
endif
intermediates := $(call intermediates-dir-for,PACKAGING,recovery_patch)
RECOVERY_FROM_BOOT_PATCH := $(intermediates)/recovery_from_boot.p
$(RECOVERY_FROM_BOOT_PATCH): PRIVATE_DIFF_TOOL := $(diff_tool)
$(RECOVERY_FROM_BOOT_PATCH): \
		$(INSTALLED_RECOVERYIMAGE_TARGET) \
		$(INSTALLED_BOOTIMAGE_TARGET) \
		$(diff_tool)
	@echo "Construct recovery from boot"
	mkdir -p $(dir $@)
	$(PRIVATE_DIFF_TOOL) $(INSTALLED_BOOTIMAGE_TARGET) $(INSTALLED_RECOVERYIMAGE_TARGET) $@
else # $(BOARD_USES_FULL_RECOVERY_IMAGE) == true
RECOVERY_FROM_BOOT_PATCH := $(INSTALLED_RECOVERYIMAGE_TARGET)
endif
endif


$(INSTALLED_SYSTEMIMAGE): $(BUILT_SYSTEMIMAGE) $(RECOVERY_FROM_BOOT_PATCH)
	@echo "Install system fs image: $@"
	$(copy-file-to-target)
	$(hide) $(call assert-max-image-size,$@ $(RECOVERY_FROM_BOOT_PATCH),$(BOARD_SYSTEMIMAGE_PARTITION_SIZE))

<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
systemimage: $(INSTALLED_SYSTEMIMAGE)
=======
$(call declare-1p-container,$(INSTALLED_SYSTEMIMAGE_TARGET),)
$(call declare-container-license-deps,$(INSTALLED_SYSTEMIMAGE_TARGET),$(BUILT_SYSTEMIMAGE),$(BUILT_SYSTEMIMAGE):/)

systemimage: $(INSTALLED_SYSTEMIMAGE_TARGET)
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)

.PHONY: systemlicense
systemlicense: $(call license-metadata-dir)/$(INSTALLED_SYSTEMIMAGE_TARGET).meta_lic reportmissinglicenses

.PHONY: systemimage-nodeps snod
systemimage-nodeps snod: $(filter-out systemimage-nodeps snod,$(MAKECMDGOALS)) \
	            | $(INTERNAL_USERIMAGES_DEPS)
	@echo "make $@: ignoring dependencies"
	$(call build-systemimage-target,$(INSTALLED_SYSTEMIMAGE))
	$(hide) $(call assert-max-image-size,$(INSTALLED_SYSTEMIMAGE),$(BOARD_SYSTEMIMAGE_PARTITION_SIZE))

ifneq (,$(filter systemimage-nodeps snod, $(MAKECMDGOALS)))
ifeq (true,$(WITH_DEXPREOPT))
$(warning Warning: with dexpreopt enabled, you may need a full rebuild.)
endif
endif

.PHONY: sync
sync: $(INTERNAL_SYSTEMIMAGE_FILES)

#######
## system tarball
define build-systemtarball-target
  $(call pretty,"Target system fs tarball: $(INSTALLED_SYSTEMTARBALL_TARGET)")
  $(call create-system-vendor-symlink)
  $(call create-system-product-symlink)
  $(MKTARBALL) $(FS_GET_STATS) \
    $(PRODUCT_OUT) system $(PRIVATE_SYSTEM_TAR) \
    $(INSTALLED_SYSTEMTARBALL_TARGET) $(TARGET_OUT)
endef

ifndef SYSTEM_TARBALL_FORMAT
    SYSTEM_TARBALL_FORMAT := bz2
endif

system_tar := $(PRODUCT_OUT)/system.tar
INSTALLED_SYSTEMTARBALL_TARGET := $(system_tar).$(SYSTEM_TARBALL_FORMAT)
$(INSTALLED_SYSTEMTARBALL_TARGET): PRIVATE_SYSTEM_TAR := $(system_tar)
$(INSTALLED_SYSTEMTARBALL_TARGET): $(FS_GET_STATS) $(INTERNAL_SYSTEMIMAGE_FILES)
	$(build-systemtarball-target)

.PHONY: systemtarball-nodeps
systemtarball-nodeps: $(FS_GET_STATS) \
                      $(filter-out systemtarball-nodeps stnod,$(MAKECMDGOALS))
	$(build-systemtarball-target)

.PHONY: stnod
stnod: systemtarball-nodeps

# -----------------------------------------------------------------
## platform.zip: system, plus other files to be used in PDK fusion build,
## in a zip file
##
## PDK_PLATFORM_ZIP_PRODUCT_BINARIES is used to store specified files to platform.zip.
## The variable will be typically set from BoardConfig.mk.
## Files under out dir will be rejected to prevent possible conflicts with other rules.
pdk_odex_javalibs := $(strip $(foreach m,$(DEXPREOPT.MODULES.JAVA_LIBRARIES),\
  $(if $(filter $(DEXPREOPT.$(m).INSTALLED),$(ALL_DEFAULT_INSTALLED_MODULES)),$(m))))
pdk_odex_apps := $(strip $(foreach m,$(DEXPREOPT.MODULES.APPS),\
  $(if $(filter $(DEXPREOPT.$(m).INSTALLED),$(ALL_DEFAULT_INSTALLED_MODULES)),$(m))))
pdk_classes_dex := $(strip \
  $(foreach m,$(pdk_odex_javalibs),$(call intermediates-dir-for,JAVA_LIBRARIES,$(m),,COMMON)/javalib.jar) \
  $(foreach m,$(pdk_odex_apps),$(call intermediates-dir-for,APPS,$(m))/package.dex.apk))

pdk_odex_config_mk := $(PRODUCT_OUT)/pdk_dexpreopt_config.mk
$(pdk_odex_config_mk): PRIVATE_JAVA_LIBRARIES := $(pdk_odex_javalibs)
$(pdk_odex_config_mk): PRIVATE_APPS := $(pdk_odex_apps)
$(pdk_odex_config_mk) :
	@echo "PDK odex config makefile: $@"
	$(hide) mkdir -p $(dir $@)
	$(hide) echo "# Auto-generated. Do not modify." > $@
	$(hide) echo "PDK.DEXPREOPT.JAVA_LIBRARIES:=$(PRIVATE_JAVA_LIBRARIES)" >> $@
	$(hide) echo "PDK.DEXPREOPT.APPS:=$(PRIVATE_APPS)" >> $@
	$(foreach m,$(PRIVATE_JAVA_LIBRARIES),\
	  $(hide) echo "PDK.DEXPREOPT.$(m).SRC:=$(patsubst $(OUT_DIR)/%,%,$(call intermediates-dir-for,JAVA_LIBRARIES,$(m),,COMMON)/javalib.jar)" >> $@$(newline)\
	  $(hide) echo "PDK.DEXPREOPT.$(m).DEX_PREOPT:=$(DEXPREOPT.$(m).DEX_PREOPT)" >> $@$(newline)\
	  $(hide) echo "PDK.DEXPREOPT.$(m).MULTILIB:=$(DEXPREOPT.$(m).MULTILIB)" >> $@$(newline)\
	  $(hide) echo "PDK.DEXPREOPT.$(m).DEX_PREOPT_FLAGS:=$(DEXPREOPT.$(m).DEX_PREOPT_FLAGS)" >> $@$(newline)\
	  )
	$(foreach m,$(PRIVATE_APPS),\
	  $(hide) echo "PDK.DEXPREOPT.$(m).SRC:=$(patsubst $(OUT_DIR)/%,%,$(call intermediates-dir-for,APPS,$(m))/package.dex.apk)" >> $@$(newline)\
	  $(hide) echo "PDK.DEXPREOPT.$(m).DEX_PREOPT:=$(DEXPREOPT.$(m).DEX_PREOPT)" >> $@$(newline)\
	  $(hide) echo "PDK.DEXPREOPT.$(m).MULTILIB:=$(DEXPREOPT.$(m).MULTILIB)" >> $@$(newline)\
	  $(hide) echo "PDK.DEXPREOPT.$(m).DEX_PREOPT_FLAGS:=$(DEXPREOPT.$(m).DEX_PREOPT_FLAGS)" >> $@$(newline)\
	  $(hide) echo "PDK.DEXPREOPT.$(m).PRIVILEGED_MODULE:=$(DEXPREOPT.$(m).PRIVILEGED_MODULE)" >> $@$(newline)\
	  $(hide) echo "PDK.DEXPREOPT.$(m).VENDOR_MODULE:=$(DEXPREOPT.$(m).VENDOR_MODULE)" >> $@$(newline)\
	  $(hide) echo "PDK.DEXPREOPT.$(m).TARGET_ARCH:=$(DEXPREOPT.$(m).TARGET_ARCH)" >> $@$(newline)\
	  $(hide) echo "PDK.DEXPREOPT.$(m).STRIPPED_SRC:=$(patsubst $(PRODUCT_OUT)/%,%,$(DEXPREOPT.$(m).INSTALLED_STRIPPED))" >> $@$(newline)\
	  )

PDK_PLATFORM_ZIP_PRODUCT_BINARIES := $(filter-out $(OUT_DIR)/%,$(PDK_PLATFORM_ZIP_PRODUCT_BINARIES))
INSTALLED_PLATFORM_ZIP := $(PRODUCT_OUT)/platform.zip

$(INSTALLED_PLATFORM_ZIP): PRIVATE_DEX_FILES := $(pdk_classes_dex)
$(INSTALLED_PLATFORM_ZIP): PRIVATE_ODEX_CONFIG := $(pdk_odex_config_mk)
$(INSTALLED_PLATFORM_ZIP) : $(INTERNAL_SYSTEMIMAGE_FILES) $(pdk_odex_config_mk)
	$(call pretty,"Platform zip package: $(INSTALLED_PLATFORM_ZIP)")
	$(hide) rm -f $@
	$(hide) cd $(dir $@) && zip -qryX $(notdir $@) \
		$(TARGET_COPY_OUT_SYSTEM) \
		$(patsubst $(PRODUCT_OUT)/%, %, $(TARGET_OUT_NOTICE_FILES)) \
		$(addprefix symbols/,$(PDK_SYMBOL_FILES_LIST))
ifdef BOARD_VENDORIMAGE_FILE_SYSTEM_TYPE
	$(hide) cd $(dir $@) && zip -qryX $(notdir $@) \
		$(TARGET_COPY_OUT_VENDOR)
endif
ifdef BOARD_PRODUCTIMAGE_FILE_SYSTEM_TYPE
	$(hide) cd $(dir $@) && zip -qryX $(notdir $@) \
		$(TARGET_COPY_OUT_PRODUCT)
endif
ifneq ($(PDK_PLATFORM_JAVA_ZIP_CONTENTS),)
	$(hide) cd $(OUT_DIR) && zip -qryX $(patsubst $(OUT_DIR)/%,%,$@) $(PDK_PLATFORM_JAVA_ZIP_CONTENTS)
endif
ifneq ($(PDK_PLATFORM_ZIP_PRODUCT_BINARIES),)
	$(hide) zip -qryX $@ $(PDK_PLATFORM_ZIP_PRODUCT_BINARIES)
endif
	@# Add dex-preopt files and config.
	$(if $(PRIVATE_DEX_FILES),$(hide) cd $(OUT_DIR) && zip -qryX $(patsubst $(OUT_DIR)/%,%,$@ $(PRIVATE_DEX_FILES)))
	$(hide) touch $(PRODUCT_OUT)/pdk.mk
	$(hide) zip -qryXj $@ $(PRIVATE_ODEX_CONFIG) $(PRODUCT_OUT)/pdk.mk

.PHONY: platform
platform: $(INSTALLED_PLATFORM_ZIP)

.PHONY: platform-java
platform-java: platform

# Dist the platform.zip
ifneq (,$(filter platform platform-java, $(MAKECMDGOALS)))
$(call dist-for-goals, platform platform-java, $(INSTALLED_PLATFORM_ZIP))
endif

# -----------------------------------------------------------------
## boot tarball
define build-boottarball-target
    $(hide) echo "Target boot fs tarball: $(INSTALLED_BOOTTARBALL_TARGET)"
    $(hide) mkdir -p $(PRODUCT_OUT)/boot
    $(hide) cp -f $(INTERNAL_BOOTIMAGE_FILES) $(PRODUCT_OUT)/boot/.
    $(hide) echo $(INTERNAL_KERNEL_CMDLINE) > $(PRODUCT_OUT)/boot/cmdline
    $(hide) $(MKTARBALL) $(FS_GET_STATS) \
                 $(PRODUCT_OUT) boot $(PRIVATE_BOOT_TAR) \
                 $(INSTALLED_BOOTTARBALL_TARGET) $(TARGET_OUT)
endef

ifndef BOOT_TARBALL_FORMAT
    BOOT_TARBALL_FORMAT := bz2
endif

boot_tar := $(PRODUCT_OUT)/boot.tar
INSTALLED_BOOTTARBALL_TARGET := $(boot_tar).$(BOOT_TARBALL_FORMAT)
$(INSTALLED_BOOTTARBALL_TARGET): PRIVATE_BOOT_TAR := $(boot_tar)
$(INSTALLED_BOOTTARBALL_TARGET): $(FS_GET_STATS) $(INTERNAL_BOOTIMAGE_FILES)
	$(build-boottarball-target)

.PHONY: boottarball-nodeps btnod
boottarball-nodeps btnod: $(FS_GET_STATS) \
                      $(filter-out boottarball-nodeps btnod,$(MAKECMDGOALS))
	$(build-boottarball-target)


# -----------------------------------------------------------------
# data partition image
INTERNAL_USERDATAIMAGE_FILES := \
    $(filter $(TARGET_OUT_DATA)/%,$(ALL_DEFAULT_INSTALLED_MODULES))

# Don't build userdata.img if it's extfs but no partition size
skip_userdata.img :=
ifdef INTERNAL_USERIMAGES_EXT_VARIANT
ifndef BOARD_USERDATAIMAGE_PARTITION_SIZE
skip_userdata.img := true
endif
endif

ifneq ($(skip_userdata.img),true)
userdataimage_intermediates := \
    $(call intermediates-dir-for,PACKAGING,userdata)
BUILT_USERDATAIMAGE_TARGET := $(PRODUCT_OUT)/userdata.img

define build-userdataimage-target
  $(call pretty,"Target userdata fs image: $(INSTALLED_USERDATAIMAGE_TARGET)")
  @mkdir -p $(TARGET_OUT_DATA)
  @mkdir -p $(userdataimage_intermediates) && rm -rf $(userdataimage_intermediates)/userdata_image_info.txt
  $(call generate-userimage-prop-dictionary, $(userdataimage_intermediates)/userdata_image_info.txt, skip_fsck=true)
  $(hide) PATH=$(foreach p,$(INTERNAL_USERIMAGES_BINARY_PATHS),$(p):)$$PATH \
      build/make/tools/releasetools/build_image.py \
      $(TARGET_OUT_DATA) $(userdataimage_intermediates)/userdata_image_info.txt $(INSTALLED_USERDATAIMAGE_TARGET) $(TARGET_OUT)
  $(hide) $(call assert-max-image-size,$(INSTALLED_USERDATAIMAGE_TARGET),$(BOARD_USERDATAIMAGE_PARTITION_SIZE))
endef

# We just build this directly to the install location.
INSTALLED_USERDATAIMAGE_TARGET := $(BUILT_USERDATAIMAGE_TARGET)
INSTALLED_USERDATAIMAGE_TARGET_DEPS := \
    $(INTERNAL_USERIMAGES_DEPS) \
    $(INTERNAL_USERDATAIMAGE_FILES) \
    $(BUILD_IMAGE_SRCS)
$(INSTALLED_USERDATAIMAGE_TARGET): $(INSTALLED_USERDATAIMAGE_TARGET_DEPS)
	$(build-userdataimage-target)

.PHONY: userdataimage-nodeps
userdataimage-nodeps: | $(INTERNAL_USERIMAGES_DEPS)
	$(build-userdataimage-target)

endif # not skip_userdata.img
skip_userdata.img :=

# ASAN libraries in the system image - build rule.
ASAN_OUT_DIRS_FOR_SYSTEM_INSTALL := $(sort $(patsubst $(PRODUCT_OUT)/%,%,\
  $(TARGET_OUT_SHARED_LIBRARIES) \
  $(2ND_TARGET_OUT_SHARED_LIBRARIES) \
  $(TARGET_OUT_VENDOR_SHARED_LIBRARIES) \
  $(2ND_TARGET_OUT_VENDOR_SHARED_LIBRARIES)))
# Extra options: Enforce the system user for the files to avoid having to change ownership.
ASAN_SYSTEM_INSTALL_OPTIONS := --owner=1000 --group=1000
# Note: experimentally, it seems not worth it to try to get "best" compression. We don't save
#       enough space.
$(ASAN_IN_SYSTEM_INSTALLED): $(INSTALLED_USERDATAIMAGE_TARGET_DEPS)
	tar cfj $(ASAN_IN_SYSTEM_INSTALLED) $(ASAN_SYSTEM_INSTALL_OPTIONS) -C $(TARGET_OUT_DATA)/.. $(ASAN_OUT_DIRS_FOR_SYSTEM_INSTALL) >/dev/null

#######
## data partition tarball
define build-userdatatarball-target
    $(call pretty,"Target userdata fs tarball: " \
                  "$(INSTALLED_USERDATATARBALL_TARGET)")
    $(MKTARBALL) $(FS_GET_STATS) \
		$(PRODUCT_OUT) data $(PRIVATE_USERDATA_TAR) \
		$(INSTALLED_USERDATATARBALL_TARGET) $(TARGET_OUT)
endef

userdata_tar := $(PRODUCT_OUT)/userdata.tar
INSTALLED_USERDATATARBALL_TARGET := $(userdata_tar).bz2
$(INSTALLED_USERDATATARBALL_TARGET): PRIVATE_USERDATA_TAR := $(userdata_tar)
$(INSTALLED_USERDATATARBALL_TARGET): $(FS_GET_STATS) $(INTERNAL_USERDATAIMAGE_FILES)
	$(build-userdatatarball-target)

$(call dist-for-goals,userdatatarball,$(INSTALLED_USERDATATARBALL_TARGET))

.PHONY: userdatatarball-nodeps
userdatatarball-nodeps: $(FS_GET_STATS)
	$(build-userdatatarball-target)


# -----------------------------------------------------------------
# partition table image
ifdef BOARD_BPT_INPUT_FILES

BUILT_BPTIMAGE_TARGET := $(PRODUCT_OUT)/partition-table.img
BUILT_BPTJSON_TARGET := $(PRODUCT_OUT)/partition-table.bpt

INTERNAL_BVBTOOL_MAKE_TABLE_ARGS := \
	--output_gpt $(BUILT_BPTIMAGE_TARGET) \
	--output_json $(BUILT_BPTJSON_TARGET) \
	$(foreach file, $(BOARD_BPT_INPUT_FILES), --input $(file))

ifdef BOARD_BPT_DISK_SIZE
INTERNAL_BVBTOOL_MAKE_TABLE_ARGS += --disk_size $(BOARD_BPT_DISK_SIZE)
endif

define build-bptimage-target
  $(call pretty,"Target partition table image: $(INSTALLED_BPTIMAGE_TARGET)")
  $(hide) $(BPTTOOL) make_table $(INTERNAL_BVBTOOL_MAKE_TABLE_ARGS) $(BOARD_BPT_MAKE_TABLE_ARGS)
endef

INSTALLED_BPTIMAGE_TARGET := $(BUILT_BPTIMAGE_TARGET)
$(BUILT_BPTJSON_TARGET): $(INSTALLED_BPTIMAGE_TARGET)
	$(hide) touch -c $(BUILT_BPTJSON_TARGET)

$(INSTALLED_BPTIMAGE_TARGET): $(BPTTOOL) $(BOARD_BPT_INPUT_FILES)
	$(build-bptimage-target)

.PHONY: bptimage-nodeps
bptimage-nodeps:
	$(build-bptimage-target)

endif # BOARD_BPT_INPUT_FILES

# -----------------------------------------------------------------
# cache partition image
ifdef BOARD_CACHEIMAGE_FILE_SYSTEM_TYPE
INTERNAL_CACHEIMAGE_FILES := \
    $(filter $(TARGET_OUT_CACHE)/%,$(ALL_DEFAULT_INSTALLED_MODULES))

cacheimage_intermediates := \
    $(call intermediates-dir-for,PACKAGING,cache)
BUILT_CACHEIMAGE_TARGET := $(PRODUCT_OUT)/cache.img

define build-cacheimage-target
  $(call pretty,"Target cache fs image: $(INSTALLED_CACHEIMAGE_TARGET)")
  @mkdir -p $(TARGET_OUT_CACHE)
  @mkdir -p $(cacheimage_intermediates) && rm -rf $(cacheimage_intermediates)/cache_image_info.txt
  $(call generate-userimage-prop-dictionary, $(cacheimage_intermediates)/cache_image_info.txt, skip_fsck=true)
  $(hide) PATH=$(foreach p,$(INTERNAL_USERIMAGES_BINARY_PATHS),$(p):)$$PATH \
      build/make/tools/releasetools/build_image.py \
      $(TARGET_OUT_CACHE) $(cacheimage_intermediates)/cache_image_info.txt $(INSTALLED_CACHEIMAGE_TARGET) $(TARGET_OUT)
  $(hide) $(call assert-max-image-size,$(INSTALLED_CACHEIMAGE_TARGET),$(BOARD_CACHEIMAGE_PARTITION_SIZE))
endef

# We just build this directly to the install location.
INSTALLED_CACHEIMAGE_TARGET := $(BUILT_CACHEIMAGE_TARGET)
$(INSTALLED_CACHEIMAGE_TARGET): $(INTERNAL_USERIMAGES_DEPS) $(INTERNAL_CACHEIMAGE_FILES) $(BUILD_IMAGE_SRCS)
	$(build-cacheimage-target)

.PHONY: cacheimage-nodeps
cacheimage-nodeps: | $(INTERNAL_USERIMAGES_DEPS)
	$(build-cacheimage-target)

else # BOARD_CACHEIMAGE_FILE_SYSTEM_TYPE
# we need to ignore the broken cache link when doing the rsync
IGNORE_CACHE_LINK := --exclude=cache
endif # BOARD_CACHEIMAGE_FILE_SYSTEM_TYPE

# -----------------------------------------------------------------
# system_other partition image
ifeq ($(BOARD_USES_SYSTEM_OTHER_ODEX),true)
BOARD_USES_SYSTEM_OTHER := true

# Marker file to identify that odex files are installed
INSTALLED_SYSTEM_OTHER_ODEX_MARKER := $(TARGET_OUT_SYSTEM_OTHER)/system-other-odex-marker
ALL_DEFAULT_INSTALLED_MODULES += $(INSTALLED_SYSTEM_OTHER_ODEX_MARKER)
$(INSTALLED_SYSTEM_OTHER_ODEX_MARKER):
	$(hide) touch $@
endif

ifdef BOARD_USES_SYSTEM_OTHER
INTERNAL_SYSTEMOTHERIMAGE_FILES := \
    $(filter $(TARGET_OUT_SYSTEM_OTHER)/%,\
      $(ALL_DEFAULT_INSTALLED_MODULES)\
      $(ALL_PDK_FUSION_FILES)) \
    $(PDK_FUSION_SYMLINK_STAMP)

INSTALLED_FILES_FILE_SYSTEMOTHER := $(PRODUCT_OUT)/installed-files-system-other.txt
$(INSTALLED_FILES_FILE_SYSTEMOTHER) : $(INTERNAL_SYSTEMOTHERIMAGE_FILES) $(FILESLIST)
	@echo Installed file list: $@
<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
	@mkdir -p $(dir $@)
	@rm -f $@
	$(hide) $(FILESLIST) $(TARGET_OUT_SYSTEM_OTHER) > $(@:.txt=.json)
	$(hide) build/make/tools/fileslist_util.py -c $(@:.txt=.json) > $@
=======
	mkdir -p $(dir $@)
	rm -f $@
	$(FILESLIST) $(TARGET_OUT_SYSTEM_OTHER) > $(@:.txt=.json)
	$(FILESLIST_UTIL) -c $(@:.txt=.json) > $@

# Determines partition size for system_other.img.
ifeq ($(PRODUCT_RETROFIT_DYNAMIC_PARTITIONS),true)
ifneq ($(filter system,$(BOARD_SUPER_PARTITION_BLOCK_DEVICES)),)
INTERNAL_SYSTEM_OTHER_PARTITION_SIZE := $(BOARD_SUPER_PARTITION_SYSTEM_DEVICE_SIZE)
endif
endif

ifndef INTERNAL_SYSTEM_OTHER_PARTITION_SIZE
INTERNAL_SYSTEM_OTHER_PARTITION_SIZE:= $(BOARD_SYSTEMIMAGE_PARTITION_SIZE)
endif
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)

systemotherimage_intermediates := \
    $(call intermediates-dir-for,PACKAGING,system_other)
BUILT_SYSTEMOTHERIMAGE_TARGET := $(PRODUCT_OUT)/system_other.img

# Note that we assert the size is SYSTEMIMAGE_PARTITION_SIZE since this is the 'b' system image.
define build-systemotherimage-target
  $(call pretty,"Target system_other fs image: $(INSTALLED_SYSTEMOTHERIMAGE_TARGET)")
  @mkdir -p $(TARGET_OUT_SYSTEM_OTHER)
  @mkdir -p $(systemotherimage_intermediates) && rm -rf $(systemotherimage_intermediates)/system_other_image_info.txt
  $(call generate-userimage-prop-dictionary, $(systemotherimage_intermediates)/system_other_image_info.txt, skip_fsck=true)
  $(hide) PATH=$(foreach p,$(INTERNAL_USERIMAGES_BINARY_PATHS),$(p):)$$PATH \
      build/make/tools/releasetools/build_image.py \
      $(TARGET_OUT_SYSTEM_OTHER) $(systemotherimage_intermediates)/system_other_image_info.txt $(INSTALLED_SYSTEMOTHERIMAGE_TARGET) $(TARGET_OUT)
  $(hide) $(call assert-max-image-size,$(INSTALLED_SYSTEMOTHERIMAGE_TARGET),$(BOARD_SYSTEMIMAGE_PARTITION_SIZE))
endef

# We just build this directly to the install location.
INSTALLED_SYSTEMOTHERIMAGE_TARGET := $(BUILT_SYSTEMOTHERIMAGE_TARGET)
ifneq (true,$(SANITIZE_LITE))
# Only create system_other when not building the second stage of a SANITIZE_LITE build.
$(INSTALLED_SYSTEMOTHERIMAGE_TARGET): $(INTERNAL_USERIMAGES_DEPS) $(INTERNAL_SYSTEMOTHERIMAGE_FILES) $(INSTALLED_FILES_FILE_SYSTEMOTHER)
	$(build-systemotherimage-target)
endif

.PHONY: systemotherimage-nodeps
systemotherimage-nodeps: | $(INTERNAL_USERIMAGES_DEPS)
	$(build-systemotherimage-target)

endif # BOARD_USES_SYSTEM_OTHER


# -----------------------------------------------------------------
# vendor partition image
ifdef BOARD_VENDORIMAGE_FILE_SYSTEM_TYPE
INTERNAL_VENDORIMAGE_FILES := \
    $(filter $(TARGET_OUT_VENDOR)/%,\
      $(ALL_DEFAULT_INSTALLED_MODULES)\
      $(ALL_PDK_FUSION_FILES)) \
    $(PDK_FUSION_SYMLINK_STAMP)

# platform.zip depends on $(INTERNAL_VENDORIMAGE_FILES).
$(INSTALLED_PLATFORM_ZIP) : $(INTERNAL_VENDORIMAGE_FILES)

INSTALLED_FILES_FILE_VENDOR := $(PRODUCT_OUT)/installed-files-vendor.txt
$(INSTALLED_FILES_FILE_VENDOR) : $(INTERNAL_VENDORIMAGE_FILES) $(FILESLIST)
	@echo Installed file list: $@
<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
	@mkdir -p $(dir $@)
	@rm -f $@
	$(hide) $(FILESLIST) $(TARGET_OUT_VENDOR) > $(@:.txt=.json)
	$(hide) build/make/tools/fileslist_util.py -c $(@:.txt=.json) > $@
=======
	mkdir -p $(dir $@)
	rm -f $@
	$(FILESLIST) $(TARGET_OUT_VENDOR) > $(@:.txt=.json)
	$(FILESLIST_UTIL) -c $(@:.txt=.json) > $@
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)

vendorimage_intermediates := \
    $(call intermediates-dir-for,PACKAGING,vendor)
BUILT_VENDORIMAGE_TARGET := $(PRODUCT_OUT)/vendor.img
define build-vendorimage-target
  $(call pretty,"Target vendor fs image: $(INSTALLED_VENDORIMAGE_TARGET)")
  @mkdir -p $(TARGET_OUT_VENDOR)
  @mkdir -p $(vendorimage_intermediates) && rm -rf $(vendorimage_intermediates)/vendor_image_info.txt
  $(call generate-userimage-prop-dictionary, $(vendorimage_intermediates)/vendor_image_info.txt, skip_fsck=true)
  $(if $(BOARD_VENDOR_KERNEL_MODULES), \
    $(call build-image-kernel-modules,$(BOARD_VENDOR_KERNEL_MODULES),$(TARGET_OUT_VENDOR),vendor/,$(call intermediates-dir-for,PACKAGING,depmod_vendor)))
  $(hide) PATH=$(foreach p,$(INTERNAL_USERIMAGES_BINARY_PATHS),$(p):)$$PATH \
      build/make/tools/releasetools/build_image.py \
      $(TARGET_OUT_VENDOR) $(vendorimage_intermediates)/vendor_image_info.txt $(INSTALLED_VENDORIMAGE_TARGET) $(TARGET_OUT)
  $(hide) $(call assert-max-image-size,$(INSTALLED_VENDORIMAGE_TARGET),$(BOARD_VENDORIMAGE_PARTITION_SIZE))
endef

# We just build this directly to the install location.
INSTALLED_VENDORIMAGE_TARGET := $(BUILT_VENDORIMAGE_TARGET)
$(INSTALLED_VENDORIMAGE_TARGET): $(INTERNAL_USERIMAGES_DEPS) $(INTERNAL_VENDORIMAGE_FILES) $(INSTALLED_FILES_FILE_VENDOR) $(BUILD_IMAGE_SRCS) $(DEPMOD) $(BOARD_VENDOR_KERNEL_MODULES)
	$(build-vendorimage-target)

.PHONY: vendorimage-nodeps vnod
vendorimage-nodeps vnod: | $(INTERNAL_USERIMAGES_DEPS) $(DEPMOD)
	$(build-vendorimage-target)

sync: $(INTERNAL_VENDORIMAGE_FILES)

else ifdef BOARD_PREBUILT_VENDORIMAGE
INSTALLED_VENDORIMAGE_TARGET := $(PRODUCT_OUT)/vendor.img
$(eval $(call copy-one-file,$(BOARD_PREBUILT_VENDORIMAGE),$(INSTALLED_VENDORIMAGE_TARGET)))
endif

# -----------------------------------------------------------------
# product partition image
ifdef BOARD_PRODUCTIMAGE_FILE_SYSTEM_TYPE
INTERNAL_PRODUCTIMAGE_FILES := \
    $(filter $(TARGET_OUT_PRODUCT)/%,\
      $(ALL_DEFAULT_INSTALLED_MODULES)\
      $(ALL_PDK_FUSION_FILES)) \
    $(PDK_FUSION_SYMLINK_STAMP)

# platform.zip depends on $(INTERNAL_PRODUCTIMAGE_FILES).
$(INSTALLED_PLATFORM_ZIP) : $(INTERNAL_PRODUCTIMAGE_FILES)

INSTALLED_FILES_FILE_PRODUCT := $(PRODUCT_OUT)/installed-files-product.txt
$(INSTALLED_FILES_FILE_PRODUCT) : $(INTERNAL_PRODUCTIMAGE_FILES) $(FILESLIST)
	@echo Installed file list: $@
<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
	@mkdir -p $(dir $@)
	@rm -f $@
	$(hide) $(FILESLIST) $(TARGET_OUT_PRODUCT) > $(@:.txt=.json)
	$(hide) build/tools/fileslist_util.py -c $(@:.txt=.json) > $@
=======
	mkdir -p $(dir $@)
	rm -f $@
	$(FILESLIST) $(TARGET_OUT_PRODUCT) > $(@:.txt=.json)
	$(FILESLIST_UTIL) -c $(@:.txt=.json) > $@
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)

productimage_intermediates := \
    $(call intermediates-dir-for,PACKAGING,product)
BUILT_PRODUCTIMAGE_TARGET := $(PRODUCT_OUT)/product.img
define build-productimage-target
  $(call pretty,"Target product fs image: $(INSTALLED_PRODUCTIMAGE_TARGET)")
  @mkdir -p $(TARGET_OUT_PRODUCT)
  @mkdir -p $(productimage_intermediates) && rm -rf $(productimage_intermediates)/product_image_info.txt
  $(call generate-userimage-prop-dictionary, $(productimage_intermediates)/product_image_info.txt, skip_fsck=true)
  $(hide) PATH=$(foreach p,$(INTERNAL_USERIMAGES_BINARY_PATHS),$(p):)$$PATH \
      ./build/tools/releasetools/build_image.py \
      $(TARGET_OUT_PRODUCT) $(productimage_intermediates)/product_image_info.txt $(INSTALLED_PRODUCTIMAGE_TARGET) $(TARGET_OUT)
  $(hide) $(call assert-max-image-size,$(INSTALLED_PRODUCTIMAGE_TARGET),$(BOARD_PRODUCTIMAGE_PARTITION_SIZE))
endef

# We just build this directly to the install location.
INSTALLED_PRODUCTIMAGE_TARGET := $(BUILT_PRODUCTIMAGE_TARGET)
$(INSTALLED_PRODUCTIMAGE_TARGET): $(INTERNAL_USERIMAGES_DEPS) $(INTERNAL_PRODUCTIMAGE_FILES) $(INSTALLED_FILES_FILE_PRODUCT) $(BUILD_IMAGE_SRCS)
	$(build-productimage-target)

.PHONY: productimage-nodeps pnod
productimage-nodeps pnod: | $(INTERNAL_USERIMAGES_DEPS)
	$(build-productimage-target)

sync: $(INTERNAL_PRODUCTIMAGE_FILES)

else ifdef BOARD_PREBUILT_PRODUCTIMAGE
INSTALLED_PRODUCTIMAGE_TARGET := $(PRODUCT_OUT)/product.img
$(eval $(call copy-one-file,$(BOARD_PREBUILT_PRODUCTIMAGE),$(INSTALLED_PRODUCTIMAGE_TARGET)))
endif

# -----------------------------------------------------------------
<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
=======
# system_ext partition image
ifdef BUILDING_SYSTEM_EXT_IMAGE
INTERNAL_SYSTEM_EXTIMAGE_FILES := \
    $(filter $(TARGET_OUT_SYSTEM_EXT)/%,\
      $(ALL_DEFAULT_INSTALLED_MODULES))

INSTALLED_FILES_FILE_SYSTEM_EXT := $(PRODUCT_OUT)/installed-files-system_ext.txt
INSTALLED_FILES_JSON_SYSTEM_EXT := $(INSTALLED_FILES_FILE_SYSTEM_EXT:.txt=.json)
$(INSTALLED_FILES_FILE_SYSTEM_EXT): .KATI_IMPLICIT_OUTPUTS := $(INSTALLED_FILES_JSON_SYSTEM_EXT)
$(INSTALLED_FILES_FILE_SYSTEM_EXT) : $(INTERNAL_SYSTEM_EXTIMAGE_FILES) $(FILESLIST) $(FILESLIST_UTIL)
	@echo Installed file list: $@
	mkdir -p $(dir $@)
	rm -f $@
	$(FILESLIST) $(TARGET_OUT_SYSTEM_EXT) > $(@:.txt=.json)
	$(FILESLIST_UTIL) -c $(@:.txt=.json) > $@

system_extimage_intermediates := \
    $(call intermediates-dir-for,PACKAGING,system_ext)
BUILT_SYSTEM_EXTIMAGE_TARGET := $(PRODUCT_OUT)/system_ext.img
define build-system_extimage-target
  $(call pretty,"Target system_ext fs image: $(INSTALLED_SYSTEM_EXTIMAGE_TARGET)")
  @mkdir -p $(TARGET_OUT_SYSTEM_EXT)
  @mkdir -p $(system_extimage_intermediates) && rm -rf $(system_extimage_intermediates)/system_ext_image_info.txt
  $(call generate-image-prop-dictionary, $(system_extimage_intermediates)/system_ext_image_info.txt,system_ext, skip_fsck=true)
  PATH=$(INTERNAL_USERIMAGES_BINARY_PATHS):$$PATH \
      $(BUILD_IMAGE) \
          $(TARGET_OUT_SYSTEM_EXT) \
          $(system_extimage_intermediates)/system_ext_image_info.txt \
          $(INSTALLED_SYSTEM_EXTIMAGE_TARGET) \
          $(TARGET_OUT)
  $(call assert-max-image-size,$(INSTALLED_PRODUCT_SERVICESIMAGE_TARGET),$(BOARD_PRODUCT_SERVICESIMAGE_PARTITION_SIZE))
endef

# We just build this directly to the install location.
INSTALLED_SYSTEM_EXTIMAGE_TARGET := $(BUILT_SYSTEM_EXTIMAGE_TARGET)
$(INSTALLED_SYSTEM_EXTIMAGE_TARGET): \
    $(INTERNAL_USERIMAGES_DEPS) \
    $(INTERNAL_SYSTEM_EXTIMAGE_FILES) \
    $(INSTALLED_FILES_FILE_SYSTEM_EXT)
	$(build-system_extimage-target)

.PHONY: systemextimage-nodeps senod
systemextimage-nodeps senod: | $(INTERNAL_USERIMAGES_DEPS)
	$(build-system_extimage-target)

sync: $(INTERNAL_SYSTEM_EXTIMAGE_FILES)

else ifdef BOARD_PREBUILT_SYSTEM_EXTIMAGE
INSTALLED_SYSTEM_EXTIMAGE_TARGET := $(PRODUCT_OUT)/system_ext.img
$(eval $(call copy-one-file,$(BOARD_PREBUILT_SYSTEM_EXTIMAGE),$(INSTALLED_SYSTEM_EXTIMAGE_TARGET)))
endif

# -----------------------------------------------------------------
# odm partition image
ifdef BUILDING_ODM_IMAGE
INTERNAL_ODMIMAGE_FILES := \
    $(filter $(TARGET_OUT_ODM)/%,\
      $(ALL_DEFAULT_INSTALLED_MODULES))

# Create symlinks for odm_dlkm on devices with a odm_dlkm partition:
# /odm/lib/modules -> /odm_dlkm/lib/modules
#
# On devices with a odm_dlkm partition,
# - /odm/lib/modules is a symlink to a directory that stores odm DLKMs.
# - /odm_dlkm/{etc,...} store other odm_dlkm files directly. The odm_dlkm partition is
#   mounted at /odm_dlkm at runtime and the symlinks created in system/core/rootdir/Android.mk
#   are hidden.
# On devices without a odm_dlkm partition,
# - /odm/lib/modules stores odm DLKMs directly.
# - /odm_dlkm/{etc,...} are symlinks to directories that store other odm_dlkm files.
#   See system/core/rootdir/Android.mk for a list of created symlinks.
# The odm DLKMs and other odm_dlkm files must not be accessed using other paths because they
# are not guaranteed to exist on all devices.
ifdef BOARD_USES_ODM_DLKMIMAGE
  INTERNAL_ODMIMAGE_FILES += $(call create-partition-compat-symlink,$(TARGET_OUT_ODM)/lib/modules,/odm_dlkm/lib/modules,odm_dlkm.img)
endif

INSTALLED_FILES_FILE_ODM := $(PRODUCT_OUT)/installed-files-odm.txt
INSTALLED_FILES_JSON_ODM := $(INSTALLED_FILES_FILE_ODM:.txt=.json)
$(INSTALLED_FILES_FILE_ODM): .KATI_IMPLICIT_OUTPUTS := $(INSTALLED_FILES_JSON_ODM)
$(INSTALLED_FILES_FILE_ODM) : $(INTERNAL_ODMIMAGE_FILES) $(FILESLIST) $(FILESLIST_UTIL)
	@echo Installed file list: $@
	mkdir -p $(dir $@)
	rm -f $@
	$(FILESLIST) $(TARGET_OUT_ODM) > $(@:.txt=.json)
	$(FILESLIST_UTIL) -c $(@:.txt=.json) > $@

odmimage_intermediates := \
    $(call intermediates-dir-for,PACKAGING,odm)
BUILT_ODMIMAGE_TARGET := $(PRODUCT_OUT)/odm.img
define build-odmimage-target
  $(call pretty,"Target odm fs image: $(INSTALLED_ODMIMAGE_TARGET)")
  @mkdir -p $(TARGET_OUT_ODM)
  @mkdir -p $(odmimage_intermediates) && rm -rf $(odmimage_intermediates)/odm_image_info.txt
  $(call generate-image-prop-dictionary, $(odmimage_intermediates)/odm_image_info.txt, odm, \
	  skip_fsck=true)
  PATH=$(INTERNAL_USERIMAGES_BINARY_PATHS):$$PATH \
      $(BUILD_IMAGE) \
          $(TARGET_OUT_ODM) $(odmimage_intermediates)/odm_image_info.txt \
          $(INSTALLED_ODMIMAGE_TARGET) $(TARGET_OUT)
  $(call assert-max-image-size,$(INSTALLED_ODMIMAGE_TARGET),$(BOARD_ODMIMAGE_PARTITION_SIZE))
endef

# We just build this directly to the install location.
INSTALLED_ODMIMAGE_TARGET := $(BUILT_ODMIMAGE_TARGET)
$(INSTALLED_ODMIMAGE_TARGET): \
    $(INTERNAL_USERIMAGES_DEPS) \
    $(INTERNAL_ODMIMAGE_FILES) \
    $(INSTALLED_FILES_FILE_ODM)
	$(build-odmimage-target)

.PHONY: odmimage-nodeps onod
odmimage-nodeps onod: | $(INTERNAL_USERIMAGES_DEPS)
	$(build-odmimage-target)

sync: $(INTERNAL_ODMIMAGE_FILES)

else ifdef BOARD_PREBUILT_ODMIMAGE
INSTALLED_ODMIMAGE_TARGET := $(PRODUCT_OUT)/odm.img
$(eval $(call copy-one-file,$(BOARD_PREBUILT_ODMIMAGE),$(INSTALLED_ODMIMAGE_TARGET)))
endif

# -----------------------------------------------------------------
# vendor_dlkm partition image
ifdef BUILDING_VENDOR_DLKM_IMAGE
INTERNAL_VENDOR_DLKMIMAGE_FILES := \
    $(filter $(TARGET_OUT_VENDOR_DLKM)/%,\
      $(ALL_DEFAULT_INSTALLED_MODULES))

INSTALLED_FILES_FILE_VENDOR_DLKM := $(PRODUCT_OUT)/installed-files-vendor_dlkm.txt
INSTALLED_FILES_JSON_VENDOR_DLKM := $(INSTALLED_FILES_FILE_VENDOR_DLKM:.txt=.json)
$(INSTALLED_FILES_FILE_VENDOR_DLKM): .KATI_IMPLICIT_OUTPUTS := $(INSTALLED_FILES_JSON_VENDOR_DLKM)
$(INSTALLED_FILES_FILE_VENDOR_DLKM) : $(INTERNAL_VENDOR_DLKMIMAGE_FILES) $(FILESLIST) $(FILESLIST_UTIL)
	@echo Installed file list: $@
	mkdir -p $(dir $@)
	rm -f $@
	$(FILESLIST) $(TARGET_OUT_VENDOR_DLKM) > $(@:.txt=.json)
	$(FILESLIST_UTIL) -c $(@:.txt=.json) > $@

vendor_dlkmimage_intermediates := \
    $(call intermediates-dir-for,PACKAGING,vendor_dlkm)
BUILT_VENDOR_DLKMIMAGE_TARGET := $(PRODUCT_OUT)/vendor_dlkm.img
define build-vendor_dlkmimage-target
  $(call pretty,"Target vendor_dlkm fs image: $(INSTALLED_VENDOR_DLKMIMAGE_TARGET)")
  @mkdir -p $(TARGET_OUT_VENDOR_DLKM)
  @mkdir -p $(vendor_dlkmimage_intermediates) && rm -rf $(vendor_dlkmimage_intermediates)/vendor_dlkm_image_info.txt
  $(call generate-image-prop-dictionary, $(vendor_dlkmimage_intermediates)/vendor_dlkm_image_info.txt, \
	  vendor_dlkm, skip_fsck=true)
  PATH=$(INTERNAL_USERIMAGES_BINARY_PATHS):$$PATH \
      $(BUILD_IMAGE) \
          $(TARGET_OUT_VENDOR_DLKM) $(vendor_dlkmimage_intermediates)/vendor_dlkm_image_info.txt \
          $(INSTALLED_VENDOR_DLKMIMAGE_TARGET) $(TARGET_OUT)
  $(call assert-max-image-size,$(INSTALLED_VENDOR_DLKMIMAGE_TARGET),$(BOARD_VENDOR_DLKMIMAGE_PARTITION_SIZE))
endef

# We just build this directly to the install location.
INSTALLED_VENDOR_DLKMIMAGE_TARGET := $(BUILT_VENDOR_DLKMIMAGE_TARGET)
$(INSTALLED_VENDOR_DLKMIMAGE_TARGET): \
    $(INTERNAL_USERIMAGES_DEPS) \
    $(INTERNAL_VENDOR_DLKMIMAGE_FILES) \
    $(INSTALLED_FILES_FILE_VENDOR_DLKM)
	$(build-vendor_dlkmimage-target)

.PHONY: vendor_dlkmimage-nodeps vdnod
vendor_dlkmimage-nodeps vdnod: | $(INTERNAL_USERIMAGES_DEPS)
	$(build-vendor_dlkmimage-target)

sync: $(INTERNAL_VENDOR_DLKMIMAGE_FILES)

else ifdef BOARD_PREBUILT_VENDOR_DLKMIMAGE
INSTALLED_VENDOR_DLKMIMAGE_TARGET := $(PRODUCT_OUT)/vendor_dlkm.img
$(eval $(call copy-one-file,$(BOARD_PREBUILT_VENDOR_DLKMIMAGE),$(INSTALLED_VENDOR_DLKMIMAGE_TARGET)))
endif

# -----------------------------------------------------------------
# odm_dlkm partition image
ifdef BUILDING_ODM_DLKM_IMAGE
INTERNAL_ODM_DLKMIMAGE_FILES := \
    $(filter $(TARGET_OUT_ODM_DLKM)/%,\
      $(ALL_DEFAULT_INSTALLED_MODULES))

INSTALLED_FILES_FILE_ODM_DLKM := $(PRODUCT_OUT)/installed-files-odm_dlkm.txt
INSTALLED_FILES_JSON_ODM_DLKM := $(INSTALLED_FILES_FILE_ODM_DLKM:.txt=.json)
$(INSTALLED_FILES_FILE_ODM_DLKM): .KATI_IMPLICIT_OUTPUTS := $(INSTALLED_FILES_JSON_ODM_DLKM)
$(INSTALLED_FILES_FILE_ODM_DLKM) : $(INTERNAL_ODM_DLKMIMAGE_FILES) $(FILESLIST) $(FILESLIST_UTIL)
	@echo Installed file list: $@
	mkdir -p $(dir $@)
	rm -f $@
	$(FILESLIST) $(TARGET_OUT_ODM_DLKM) > $(@:.txt=.json)
	$(FILESLIST_UTIL) -c $(@:.txt=.json) > $@

odm_dlkmimage_intermediates := \
    $(call intermediates-dir-for,PACKAGING,odm_dlkm)
BUILT_ODM_DLKMIMAGE_TARGET := $(PRODUCT_OUT)/odm_dlkm.img
define build-odm_dlkmimage-target
  $(call pretty,"Target odm_dlkm fs image: $(INSTALLED_ODM_DLKMIMAGE_TARGET)")
  @mkdir -p $(TARGET_OUT_ODM_DLKM)
  @mkdir -p $(odm_dlkmimage_intermediates) && rm -rf $(odm_dlkmimage_intermediates)/odm_dlkm_image_info.txt
  $(call generate-image-prop-dictionary, $(odm_dlkmimage_intermediates)/odm_dlkm_image_info.txt, \
	  odm_dlkm, skip_fsck=true)
  PATH=$(INTERNAL_USERIMAGES_BINARY_PATHS):$$PATH \
      $(BUILD_IMAGE) \
          $(TARGET_OUT_ODM_DLKM) $(odm_dlkmimage_intermediates)/odm_dlkm_image_info.txt \
          $(INSTALLED_ODM_DLKMIMAGE_TARGET) $(TARGET_OUT)
  $(call assert-max-image-size,$(INSTALLED_ODM_DLKMIMAGE_TARGET),$(BOARD_ODM_DLKMIMAGE_PARTITION_SIZE))
endef

# We just build this directly to the install location.
INSTALLED_ODM_DLKMIMAGE_TARGET := $(BUILT_ODM_DLKMIMAGE_TARGET)
$(INSTALLED_ODM_DLKMIMAGE_TARGET): \
    $(INTERNAL_USERIMAGES_DEPS) \
    $(INTERNAL_ODM_DLKMIMAGE_FILES) \
    $(INSTALLED_FILES_FILE_ODM_DLKM)
	$(build-odm_dlkmimage-target)

.PHONY: odm_dlkmimage-nodeps odnod
odm_dlkmimage-nodeps odnod: | $(INTERNAL_USERIMAGES_DEPS)
	$(build-odm_dlkmimage-target)

sync: $(INTERNAL_ODM_DLKMIMAGE_FILES)

else ifdef BOARD_PREBUILT_ODM_DLKMIMAGE
INSTALLED_ODM_DLKMIMAGE_TARGET := $(PRODUCT_OUT)/odm_dlkm.img
$(eval $(call copy-one-file,$(BOARD_PREBUILT_ODM_DLKMIMAGE),$(INSTALLED_ODM_DLKMIMAGE_TARGET)))
endif


# -----------------------------------------------------------------
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
# dtbo image
ifdef BOARD_PREBUILT_DTBOIMAGE
INSTALLED_DTBOIMAGE_TARGET := $(PRODUCT_OUT)/dtbo.img

ifeq ($(BOARD_AVB_ENABLE),true)
$(INSTALLED_DTBOIMAGE_TARGET): $(BOARD_PREBUILT_DTBOIMAGE) $(AVBTOOL) $(BOARD_AVB_DTBO_KEY_PATH)
	cp $(BOARD_PREBUILT_DTBOIMAGE) $@
	$(AVBTOOL) add_hash_footer \
		--image $@ \
		--partition_size $(BOARD_DTBOIMG_PARTITION_SIZE) \
		--partition_name dtbo $(INTERNAL_AVB_DTBO_SIGNING_ARGS) \
		$(BOARD_AVB_DTBO_ADD_HASH_FOOTER_ARGS)
else
$(INSTALLED_DTBOIMAGE_TARGET): $(BOARD_PREBUILT_DTBOIMAGE)
	cp $(BOARD_PREBUILT_DTBOIMAGE) $@
endif

<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
=======
endif # BOARD_PREBUILT_DTBOIMAGE

# -----------------------------------------------------------------
# Protected VM firmware image
ifeq ($(BOARD_USES_PVMFWIMAGE),true)
INSTALLED_PVMFWIMAGE_TARGET := $(PRODUCT_OUT)/pvmfw.img
INTERNAL_PREBUILT_PVMFWIMAGE := packages/modules/Virtualization/pvmfw/pvmfw.img

ifdef BOARD_PREBUILT_PVMFWIMAGE
BUILT_PVMFWIMAGE_TARGET := $(BOARD_PREBUILT_PVMFWIMAGE)
else ifeq ($(BUILDING_PVMFW_IMAGE),true)
BUILT_PVMFWIMAGE_TARGET := $(INTERNAL_PREBUILT_PVMFWIMAGE)
endif

ifeq ($(BOARD_AVB_ENABLE),true)
$(INSTALLED_PVMFWIMAGE_TARGET): $(BUILT_PVMFWIMAGE_TARGET) $(AVBTOOL) $(BOARD_AVB_PVMFW_KEY_PATH)
	cp $(BUILT_PVMFWIMAGE_TARGET) $@
	$(AVBTOOL) add_hash_footer \
	    --image $@ \
	    --partition_size $(BOARD_PVMFWIMAGE_PARTITION_SIZE) \
	    --partition_name pvmfw $(INTERNAL_AVB_PVMFW_SIGNING_ARGS) \
	    $(BOARD_AVB_PVMFW_ADD_HASH_FOOTER_ARGS)
else
$(eval $(call copy-one-file,$(BUILT_PVMFWIMAGE_TARGET),$(INSTALLED_PVMFWIMAGE_TARGET)))
endif

endif # BOARD_USES_PVMFWIMAGE

# Returns a list of image targets corresponding to the given list of partitions. For example, it
# returns "$(INSTALLED_PRODUCTIMAGE_TARGET)" for "product", or "$(INSTALLED_SYSTEMIMAGE_TARGET)
# $(INSTALLED_VENDORIMAGE_TARGET)" for "system vendor".
# (1): list of partitions like "system", "vendor" or "system product system_ext".
define images-for-partitions
$(strip $(foreach item,$(1),$(if $(filter $(item),system_other),$(INSTALLED_SYSTEMOTHERIMAGE_TARGET),$(INSTALLED_$(call to-upper,$(item))IMAGE_TARGET))))
endef

# -----------------------------------------------------------------
# custom images
INSTALLED_CUSTOMIMAGES_TARGET :=

ifneq ($(strip $(BOARD_CUSTOMIMAGES_PARTITION_LIST)),)
INTERNAL_AVB_CUSTOMIMAGES_SIGNING_ARGS :=

# Sign custom image.
# $(1): the prebuilt custom image.
# $(2): the mount point of the prebuilt custom image.
# $(3): the signed custom image target.
define sign_custom_image
$(3): $(1) $(INTERNAL_USERIMAGES_DEPS)
	@echo Target custom image: $(3)
	mkdir -p $(dir $(3))
	cp $(1) $(3)
ifeq ($(BOARD_AVB_ENABLE),true)
	PATH=$(INTERNAL_USERIMAGES_BINARY_PATHS):$$$$PATH \
          $(AVBTOOL) add_hashtree_footer \
          --image $(3) \
          --key $(BOARD_AVB_$(call to-upper,$(2))_KEY_PATH) \
          --algorithm $(BOARD_AVB_$(call to-upper,$(2))_ALGORITHM) \
          --partition_size $(BOARD_AVB_$(call to-upper,$(2))_PARTITION_SIZE) \
          --partition_name $(2) \
          $(INTERNAL_AVB_CUSTOMIMAGES_SIGNING_ARGS) \
          $(BOARD_AVB_$(call to-upper,$(2))_ADD_HASHTREE_FOOTER_ARGS)
endif
INSTALLED_CUSTOMIMAGES_TARGET += $(3)
endef

$(foreach partition,$(BOARD_CUSTOMIMAGES_PARTITION_LIST), \
  $(foreach image,$(BOARD_AVB_$(call to-upper,$(partition))_IMAGE_LIST), \
     $(eval $(call sign_custom_image,$(image),$(partition),$(PRODUCT_OUT)/$(notdir $(image))))))
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
endif

# -----------------------------------------------------------------
# vbmeta image
ifeq ($(BOARD_AVB_ENABLE),true)

BUILT_VBMETAIMAGE_TARGET := $(PRODUCT_OUT)/vbmeta.img
AVB_CHAIN_KEY_DIR := $(TARGET_OUT_INTERMEDIATES)/avb_chain_keys

ifdef BOARD_AVB_KEY_PATH
$(if $(BOARD_AVB_ALGORITHM),,$(error BOARD_AVB_ALGORITHM is not defined))
else
# If key path isn't specified, use the 4096-bit test key.
BOARD_AVB_ALGORITHM := SHA256_RSA4096
BOARD_AVB_KEY_PATH := external/avb/test/data/testkey_rsa4096.pem
endif

INTERNAL_AVB_SIGNING_ARGS := \
    --algorithm $(BOARD_AVB_ALGORITHM) --key $(BOARD_AVB_KEY_PATH)

BOOT_FOOTER_ARGS := BOARD_AVB_BOOT_ADD_HASH_FOOTER_ARGS
DTBO_FOOTER_ARGS := BOARD_AVB_DTBO_ADD_HASH_FOOTER_ARGS
SYSTEM_FOOTER_ARGS := BOARD_AVB_SYSTEM_ADD_HASHTREE_FOOTER_ARGS
VENDOR_FOOTER_ARGS := BOARD_AVB_VENDOR_ADD_HASHTREE_FOOTER_ARGS
RECOVERY_FOOTER_ARGS := BOARD_AVB_RECOVERY_ADD_HASH_FOOTER_ARGS
PRODUCT_FOOTER_ARGS := BOARD_AVB_PRODUCT_ADD_HASHTREE_FOOTER_ARGS

# Check and set required build variables for a chain partition.
# $(1): the partition to enable AVB chain, e.g., BOOT or SYSTEM.
define check-and-set-avb-chain-args
$(eval PART := $(1))
$(eval part=$(call to-lower,$(PART)))

$(eval _key_path := BOARD_AVB_$(PART)_KEY_PATH)
$(eval _signing_algorithm := BOARD_AVB_$(PART)_ALGORITHM)
$(eval _rollback_index := BOARD_AVB_$(PART)_ROLLBACK_INDEX)
$(eval _rollback_index_location := BOARD_AVB_$(PART)_ROLLBACK_INDEX_LOCATION)
$(if $($(_key_path)),,$(error $(_key_path) is not defined))
$(if $($(_signing_algorithm)),,$(error $(_signing_algorithm) is not defined))
$(if $($(_rollback_index)),,$(error $(_rollback_index) is not defined))
$(if $($(_rollback_index_location)),,$(error $(_rollback_index_location) is not defined))

# Set INTERNAL_AVB_(PART)_SIGNING_ARGS
$(eval _signing_args := INTERNAL_AVB_$(PART)_SIGNING_ARGS)
$(eval $(_signing_args) := \
    --algorithm $($(_signing_algorithm)) --key $($(_key_path)))

$(eval INTERNAL_AVB_MAKE_VBMETA_IMAGE_ARGS += \
    --chain_partition $(part):$($(_rollback_index_location)):$(AVB_CHAIN_KEY_DIR)/$(part).avbpubkey)

# Set rollback_index via footer args
$(eval _footer_args := $(PART)_FOOTER_ARGS)
$(eval $($(_footer_args)) += --rollback_index $($(_rollback_index)))
endef

ifdef BOARD_AVB_BOOT_KEY_PATH
$(eval $(call check-and-set-avb-chain-args,BOOT))
else
INTERNAL_AVB_MAKE_VBMETA_IMAGE_ARGS += \
    --include_descriptors_from_image $(INSTALLED_BOOTIMAGE_TARGET)
endif

ifdef BOARD_AVB_SYSTEM_KEY_PATH
$(eval $(call check-and-set-avb-chain-args,SYSTEM))
else
INTERNAL_AVB_MAKE_VBMETA_IMAGE_ARGS += \
    --include_descriptors_from_image $(INSTALLED_SYSTEMIMAGE)
endif

ifdef INSTALLED_VENDORIMAGE_TARGET
ifdef BOARD_AVB_VENDOR_KEY_PATH
$(eval $(call check-and-set-avb-chain-args,VENDOR))
else
INTERNAL_AVB_MAKE_VBMETA_IMAGE_ARGS += \
    --include_descriptors_from_image $(INSTALLED_VENDORIMAGE_TARGET)
endif
endif

ifdef INSTALLED_PRODUCTIMAGE_TARGET
ifdef BOARD_AVB_PRODUCT_KEY_PATH
$(eval $(call check-and-set-avb-chain-args,PRODUCT))
else
INTERNAL_AVB_MAKE_VBMETA_IMAGE_ARGS += \
    --include_descriptors_from_image $(INSTALLED_PRODUCTIMAGE_TARGET)
endif
endif

ifdef INSTALLED_DTBOIMAGE_TARGET
ifdef BOARD_AVB_DTBO_KEY_PATH
$(eval $(call check-and-set-avb-chain-args,DTBO))
else
INTERNAL_AVB_MAKE_VBMETA_IMAGE_ARGS += \
    --include_descriptors_from_image $(INSTALLED_DTBOIMAGE_TARGET)
endif
endif

ifdef INSTALLED_RECOVERYIMAGE_TARGET
ifdef BOARD_AVB_RECOVERY_KEY_PATH
$(eval $(call check-and-set-avb-chain-args,RECOVERY))
else
INTERNAL_AVB_MAKE_VBMETA_IMAGE_ARGS += \
    --include_descriptors_from_image $(INSTALLED_RECOVERYIMAGE_TARGET)
endif
endif

BOARD_AVB_MAKE_VBMETA_IMAGE_ARGS += --padding_size 4096

# Add kernel cmdline descriptor for kernel to mount system.img as root with
# dm-verity. This works when system.img is either chained or not-chained:
# - chained: The --setup_as_rootfs_from_kernel option will add dm-verity kernel
#   cmdline descriptor to system.img
# - not-chained: The --include_descriptors_from_image option for make_vbmeta_image
#   will include the kernel cmdline descriptor from system.img into vbmeta.img
ifeq ($(BOARD_BUILD_SYSTEM_ROOT_IMAGE),true)
BOARD_AVB_SYSTEM_ADD_HASHTREE_FOOTER_ARGS += --setup_as_rootfs_from_kernel
endif

ifdef BOARD_AVB_ROLLBACK_INDEX
BOARD_AVB_MAKE_VBMETA_IMAGE_ARGS += --rollback_index $(BOARD_AVB_ROLLBACK_INDEX)
endif

ifndef BOARD_BOOTIMAGE_PARTITION_SIZE
  $(error BOARD_BOOTIMAGE_PARTITION_SIZE must be set for BOARD_AVB_ENABLE)
endif

ifndef BOARD_SYSTEMIMAGE_PARTITION_SIZE
  $(error BOARD_SYSTEMIMAGE_PARTITION_SIZE must be set for BOARD_AVB_ENABLE)
endif

# $(1): the directory to extract public keys to
define extract-avb-chain-public-keys
  $(if $(BOARD_AVB_BOOT_KEY_PATH),\
    $(hide) $(AVBTOOL) extract_public_key --key $(BOARD_AVB_BOOT_KEY_PATH) \
      --output $(1)/boot.avbpubkey)
  $(if $(BOARD_AVB_SYSTEM_KEY_PATH),\
    $(hide) $(AVBTOOL) extract_public_key --key $(BOARD_AVB_SYSTEM_KEY_PATH) \
      --output $(1)/system.avbpubkey)
  $(if $(BOARD_AVB_VENDOR_KEY_PATH),\
    $(hide) $(AVBTOOL) extract_public_key --key $(BOARD_AVB_VENDOR_KEY_PATH) \
      --output $(1)/vendor.avbpubkey)
  $(if $(BOARD_AVB_PRODUCT_KEY_PATH),\
    $(hide) $(AVBTOOL) extract_public_key --key $(BOARD_AVB_PRODUCT_KEY_PATH) \
      --output $(1)/product.avbpubkey)
  $(if $(BOARD_AVB_DTBO_KEY_PATH),\
    $(hide) $(AVBTOOL) extract_public_key --key $(BOARD_AVB_DTBO_KEY_PATH) \
      --output $(1)/dtbo.avbpubkey)
  $(if $(BOARD_AVB_RECOVERY_KEY_PATH),\
    $(hide) $(AVBTOOL) extract_public_key --key $(BOARD_AVB_RECOVERY_KEY_PATH) \
      --output $(1)/recovery.avbpubkey)
endef

define build-vbmetaimage-target
  $(call pretty,"Target vbmeta image: $(INSTALLED_VBMETAIMAGE_TARGET)")
  $(hide) mkdir -p $(AVB_CHAIN_KEY_DIR)
  $(call extract-avb-chain-public-keys, $(AVB_CHAIN_KEY_DIR))
  $(hide) $(AVBTOOL) make_vbmeta_image \
    $(INTERNAL_AVB_MAKE_VBMETA_IMAGE_ARGS) \
    $(INTERNAL_AVB_SIGNING_ARGS) \
    $(BOARD_AVB_MAKE_VBMETA_IMAGE_ARGS) \
    --output $@
  $(hide) rm -rf $(AVB_CHAIN_KEY_DIR)
endef

INSTALLED_VBMETAIMAGE_TARGET := $(BUILT_VBMETAIMAGE_TARGET)
$(INSTALLED_VBMETAIMAGE_TARGET): \
		$(AVBTOOL) \
		$(INSTALLED_BOOTIMAGE_TARGET) \
		$(INSTALLED_SYSTEMIMAGE) \
		$(INSTALLED_VENDORIMAGE_TARGET) \
		$(INSTALLED_PRODUCTIMAGE_TARGET) \
		$(INSTALLED_DTBOIMAGE_TARGET) \
		$(INSTALLED_RECOVERYIMAGE_TARGET) \
		$(BOARD_AVB_KEY_PATH)
	$(build-vbmetaimage-target)

.PHONY: vbmetaimage-nodeps
vbmetaimage-nodeps:
	$(build-vbmetaimage-target)

else ifeq (true,$(BOARD_BUILD_DISABLED_VBMETAIMAGE))
BUILT_DISABLED_VBMETAIMAGE := $(PRODUCT_OUT)/vbmeta.img

INSTALLED_VBMETAIMAGE_TARGET := $(BUILT_DISABLED_VBMETAIMAGE)
$(INSTALLED_VBMETAIMAGE_TARGET): $(AVBTOOL)
	$(hide) $(AVBTOOL) make_vbmeta_image \
	  --flag 2 --padding_size 4096 --output $@

endif # BOARD_AVB_ENABLE

# List of files from all images
INTERNAL_ALLIMAGES_FILES := \
    $(FULL_SYSTEMIMAGE_DEPS) \
    $(INTERNAL_RAMDISK_FILES) \
    $(INTERNAL_USERDATAIMAGE_FILES) \
    $(INTERNAL_VENDORIMAGE_FILES) \
    $(INTERNAL_PRODUCTIMAGE_FILES) \
    $(INTERNAL_SYSTEM_EXTIMAGE_FILES) \
    $(INTERNAL_ODMIMAGE_FILES) \
    $(INTERNAL_VENDOR_DLKMIMAGE_FILES) \
    $(INTERNAL_ODM_DLKMIMAGE_FILES) \

# -----------------------------------------------------------------
<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
=======
# Check VINTF of build

# Note: vendor_dlkm and odm_dlkm does not have VINTF files.
ifeq (,$(TARGET_BUILD_UNBUNDLED))

intermediates := $(call intermediates-dir-for,PACKAGING,check_vintf_all)
check_vintf_all_deps :=

# The build system only writes VINTF metadata to */etc/vintf paths. Legacy paths aren't needed here
# because they are only used for prebuilt images.
check_vintf_common_srcs_patterns := \
  $(TARGET_OUT)/etc/vintf/% \
  $(TARGET_OUT_VENDOR)/etc/vintf/% \
  $(TARGET_OUT_ODM)/etc/vintf/% \
  $(TARGET_OUT_PRODUCT)/etc/vintf/% \
  $(TARGET_OUT_SYSTEM_EXT)/etc/vintf/% \

check_vintf_common_srcs := $(sort $(filter $(check_vintf_common_srcs_patterns),$(INTERNAL_ALLIMAGES_FILES)))
check_vintf_common_srcs_patterns :=

check_vintf_has_system :=
check_vintf_has_vendor :=

ifneq (,$(filter EMPTY_ODM_SKU_PLACEHOLDER,$(ODM_MANIFEST_SKUS)))
$(error EMPTY_ODM_SKU_PLACEHOLDER is an internal variable and cannot be used for ODM_MANIFEST_SKUS)
endif
ifneq (,$(filter EMPTY_VENDOR_SKU_PLACEHOLDER,$(DEVICE_MANIFEST_SKUS)))
$(error EMPTY_VENDOR_SKU_PLACEHOLDER is an internal variable and cannot be used for DEIVCE_MANIFEST_SKUS)
endif

# -- Check system manifest / matrix including fragments (excluding other framework manifests / matrices, e.g. product);
check_vintf_system_deps := $(filter $(TARGET_OUT)/etc/vintf/%, $(check_vintf_common_srcs))
ifneq ($(check_vintf_system_deps),)
check_vintf_has_system := true

check_vintf_system_log := $(intermediates)/check_vintf_system.log
check_vintf_all_deps += $(check_vintf_system_log)
$(check_vintf_system_log): $(HOST_OUT_EXECUTABLES)/checkvintf $(check_vintf_system_deps)
	@( $< --check-one --dirmap /system:$(TARGET_OUT) > $@ 2>&1 ) || ( cat $@ && exit 1 )
check_vintf_system_log :=

vintffm_log := $(intermediates)/vintffm.log
check_vintf_all_deps += $(vintffm_log)
$(vintffm_log): $(HOST_OUT_EXECUTABLES)/vintffm $(check_vintf_system_deps)
	@( $< --check --dirmap /system:$(TARGET_OUT) \
	  $(VINTF_FRAMEWORK_MANIFEST_FROZEN_DIR) > $@ 2>&1 ) || ( cat $@ && exit 1 )

endif # check_vintf_system_deps
check_vintf_system_deps :=

# -- Check vendor manifest / matrix including fragments (excluding other device manifests / matrices)
check_vintf_vendor_deps := $(filter $(TARGET_OUT_VENDOR)/etc/vintf/%, $(check_vintf_common_srcs))
ifneq ($(check_vintf_vendor_deps),)
check_vintf_has_vendor := true
check_vintf_vendor_log := $(intermediates)/check_vintf_vendor.log
check_vintf_all_deps += $(check_vintf_vendor_log)
# Check vendor SKU=(empty) case when:
# - DEVICE_MANIFEST_FILE is not empty; OR
# - DEVICE_MANIFEST_FILE is empty AND DEVICE_MANIFEST_SKUS is empty (only vendor manifest fragments are used)
$(check_vintf_vendor_log): PRIVATE_VENDOR_SKUS := \
  $(if $(DEVICE_MANIFEST_FILE),EMPTY_VENDOR_SKU_PLACEHOLDER,\
    $(if $(DEVICE_MANIFEST_SKUS),,EMPTY_VENDOR_SKU_PLACEHOLDER)) \
  $(DEVICE_MANIFEST_SKUS)
$(check_vintf_vendor_log): $(HOST_OUT_EXECUTABLES)/checkvintf $(check_vintf_vendor_deps)
	$(foreach vendor_sku,$(PRIVATE_VENDOR_SKUS), \
	  ( $< --check-one --dirmap /vendor:$(TARGET_OUT_VENDOR) \
	       --property ro.boot.product.vendor.sku=$(filter-out EMPTY_VENDOR_SKU_PLACEHOLDER,$(vendor_sku)) \
	       > $@ 2>&1 ) || ( cat $@ && exit 1 ); )
check_vintf_vendor_log :=
endif # check_vintf_vendor_deps
check_vintf_vendor_deps :=

# -- Kernel version and configurations.
ifeq ($(PRODUCT_OTA_ENFORCE_VINTF_KERNEL_REQUIREMENTS),true)

BUILT_KERNEL_CONFIGS_FILE := $(intermediates)/kernel_configs.txt
BUILT_KERNEL_VERSION_FILE := $(intermediates)/kernel_version.txt

my_board_extracted_kernel :=

# BOARD_KERNEL_CONFIG_FILE and BOARD_KERNEL_VERSION can be used to override the values extracted
# from INSTALLED_KERNEL_TARGET.
ifdef BOARD_KERNEL_CONFIG_FILE
ifdef BOARD_KERNEL_VERSION
$(BUILT_KERNEL_CONFIGS_FILE): $(BOARD_KERNEL_CONFIG_FILE)
	cp $< $@
$(BUILT_KERNEL_VERSION_FILE):
	echo $(BOARD_KERNEL_VERSION) > $@

my_board_extracted_kernel := true
endif # BOARD_KERNEL_VERSION
endif # BOARD_KERNEL_CONFIG_FILE

ifneq ($(my_board_extracted_kernel),true)
# Tools for decompression that is not in PATH.
# Check $(EXTRACT_KERNEL) for decompression algorithms supported by the script.
# Algorithms that are in the script but not in this list will be found in PATH.
my_decompress_tools := \
    lz4:$(HOST_OUT_EXECUTABLES)/lz4 \

endif # my_board_extracted_kernel

ifneq ($(my_board_extracted_kernel),true)
ifdef INSTALLED_KERNEL_TARGET
$(BUILT_KERNEL_CONFIGS_FILE): .KATI_IMPLICIT_OUTPUTS := $(BUILT_KERNEL_VERSION_FILE)
$(BUILT_KERNEL_CONFIGS_FILE): PRIVATE_DECOMPRESS_TOOLS := $(my_decompress_tools)
$(BUILT_KERNEL_CONFIGS_FILE): $(foreach pair,$(my_decompress_tools),$(call word-colon,2,$(pair)))
$(BUILT_KERNEL_CONFIGS_FILE): $(EXTRACT_KERNEL) $(firstword $(INSTALLED_KERNEL_TARGET))
	$< --tools $(PRIVATE_DECOMPRESS_TOOLS) --input $(firstword $(INSTALLED_KERNEL_TARGET)) \
	  --output-configs $@ \
	  --output-release $(BUILT_KERNEL_VERSION_FILE)

my_board_extracted_kernel := true
endif # INSTALLED_KERNEL_TARGET
endif # my_board_extracted_kernel

ifneq ($(my_board_extracted_kernel),true)
ifdef INSTALLED_BOOTIMAGE_TARGET
$(BUILT_KERNEL_CONFIGS_FILE): .KATI_IMPLICIT_OUTPUTS := $(BUILT_KERNEL_VERSION_FILE)
$(BUILT_KERNEL_CONFIGS_FILE): PRIVATE_DECOMPRESS_TOOLS := $(my_decompress_tools)
$(BUILT_KERNEL_CONFIGS_FILE): $(foreach pair,$(my_decompress_tools),$(call word-colon,2,$(pair)))
$(BUILT_KERNEL_CONFIGS_FILE): PRIVATE_UNPACKED_BOOTIMG := $(intermediates)/unpacked_bootimage
$(BUILT_KERNEL_CONFIGS_FILE): \
        $(HOST_OUT_EXECUTABLES)/unpack_bootimg \
        $(EXTRACT_KERNEL) \
        $(INSTALLED_BOOTIMAGE_TARGET)
	$(HOST_OUT_EXECUTABLES)/unpack_bootimg --boot_img $(INSTALLED_BOOTIMAGE_TARGET) --out $(PRIVATE_UNPACKED_BOOTIMG)
	$(EXTRACT_KERNEL) --tools $(PRIVATE_DECOMPRESS_TOOLS) --input $(PRIVATE_UNPACKED_BOOTIMG)/kernel \
	  --output-configs $@ \
	  --output-release $(BUILT_KERNEL_VERSION_FILE)

my_board_extracted_kernel := true
endif # INSTALLED_BOOTIMAGE_TARGET
endif # my_board_extracted_kernel

ifneq ($(my_board_extracted_kernel),true)
$(warning Neither INSTALLED_KERNEL_TARGET nor INSTALLED_BOOTIMAGE_TARGET is defined when \
    PRODUCT_OTA_ENFORCE_VINTF_KERNEL_REQUIREMENTS  is true. Information about the updated kernel \
    cannot be built into OTA update package. You can fix this by: \
    (1) setting TARGET_NO_KERNEL to false and installing the built kernel to $(PRODUCT_OUT)/kernel,\
        so that kernel information will be extracted from the built kernel; or \
    (2) Add a prebuilt boot image and specify it in BOARD_PREBUILT_BOOTIMAGE; or \
    (3) extracting kernel configuration and defining BOARD_KERNEL_CONFIG_FILE and \
        BOARD_KERNEL_VERSION manually; or \
    (4) unsetting PRODUCT_OTA_ENFORCE_VINTF_KERNEL_REQUIREMENTS manually.)
# Clear their values to indicate that these two files does not exist.
BUILT_KERNEL_CONFIGS_FILE :=
BUILT_KERNEL_VERSION_FILE :=
endif

my_decompress_tools :=
my_board_extracted_kernel :=

endif # PRODUCT_OTA_ENFORCE_VINTF_KERNEL_REQUIREMENTS

# -- Check VINTF compatibility of build.
# Skip partial builds; only check full builds. Only check if:
# - PRODUCT_ENFORCE_VINTF_MANIFEST is true
# - system / vendor VINTF metadata exists
# - Building product / system_ext / odm images if board has product / system_ext / odm images
ifeq ($(PRODUCT_ENFORCE_VINTF_MANIFEST),true)
ifeq ($(check_vintf_has_system),true)
ifeq ($(check_vintf_has_vendor),true)
ifeq ($(filter true,$(BUILDING_ODM_IMAGE)),$(filter true,$(BOARD_USES_ODMIMAGE)))
ifeq ($(filter true,$(BUILDING_PRODUCT_IMAGE)),$(filter true,$(BOARD_USES_PRODUCTIMAGE)))
ifeq ($(filter true,$(BUILDING_SYSTEM_EXT_IMAGE)),$(filter true,$(BOARD_USES_SYSTEM_EXTIMAGE)))

check_vintf_compatible_log := $(intermediates)/check_vintf_compatible.log
check_vintf_all_deps += $(check_vintf_compatible_log)

check_vintf_compatible_args :=
check_vintf_compatible_deps := $(check_vintf_common_srcs)

ifeq ($(PRODUCT_OTA_ENFORCE_VINTF_KERNEL_REQUIREMENTS),true)
ifneq (,$(BUILT_KERNEL_VERSION_FILE)$(BUILT_KERNEL_CONFIGS_FILE))
check_vintf_compatible_args += --kernel $(BUILT_KERNEL_VERSION_FILE):$(BUILT_KERNEL_CONFIGS_FILE)
check_vintf_compatible_deps += $(BUILT_KERNEL_CONFIGS_FILE) $(BUILT_KERNEL_VERSION_FILE)
endif # BUILT_KERNEL_VERSION_FILE != "" || BUILT_KERNEL_CONFIGS_FILE != ""
endif # PRODUCT_OTA_ENFORCE_VINTF_KERNEL_REQUIREMENTS

check_vintf_compatible_args += \
  --dirmap /system:$(TARGET_OUT) \
  --dirmap /vendor:$(TARGET_OUT_VENDOR) \
  --dirmap /odm:$(TARGET_OUT_ODM) \
  --dirmap /product:$(TARGET_OUT_PRODUCT) \
  --dirmap /system_ext:$(TARGET_OUT_SYSTEM_EXT) \

ifdef PRODUCT_SHIPPING_API_LEVEL
check_vintf_compatible_args += --property ro.product.first_api_level=$(PRODUCT_SHIPPING_API_LEVEL)
endif # PRODUCT_SHIPPING_API_LEVEL

$(check_vintf_compatible_log): PRIVATE_CHECK_VINTF_ARGS := $(check_vintf_compatible_args)
$(check_vintf_compatible_log): PRIVATE_CHECK_VINTF_DEPS := $(check_vintf_compatible_deps)
# Check ODM SKU=(empty) case when:
# - ODM_MANIFEST_FILES is not empty; OR
# - ODM_MANIFEST_FILES is empty AND ODM_MANIFEST_SKUS is empty (only ODM manifest fragments are used)
$(check_vintf_compatible_log): PRIVATE_ODM_SKUS := \
  $(if $(ODM_MANIFEST_FILES),EMPTY_ODM_SKU_PLACEHOLDER,\
    $(if $(ODM_MANIFEST_SKUS),,EMPTY_ODM_SKU_PLACEHOLDER)) \
  $(ODM_MANIFEST_SKUS)
# Check vendor SKU=(empty) case when:
# - DEVICE_MANIFEST_FILE is not empty; OR
# - DEVICE_MANIFEST_FILE is empty AND DEVICE_MANIFEST_SKUS is empty (only vendor manifest fragments are used)
$(check_vintf_compatible_log): PRIVATE_VENDOR_SKUS := \
  $(if $(DEVICE_MANIFEST_FILE),EMPTY_VENDOR_SKU_PLACEHOLDER,\
    $(if $(DEVICE_MANIFEST_SKUS),,EMPTY_VENDOR_SKU_PLACEHOLDER)) \
  $(DEVICE_MANIFEST_SKUS)
$(check_vintf_compatible_log): $(HOST_OUT_EXECUTABLES)/checkvintf $(check_vintf_compatible_deps)
	@echo "PRODUCT_OTA_ENFORCE_VINTF_KERNEL_REQUIREMENTS=$(PRODUCT_OTA_ENFORCE_VINTF_KERNEL_REQUIREMENTS)" > $@
	@echo -n -e 'Deps: \n  ' >> $@
	@sed 's/ /\n  /g' <<< "$(PRIVATE_CHECK_VINTF_DEPS)" >> $@
	@echo -n -e 'Args: \n  ' >> $@
	@cat <<< "$(PRIVATE_CHECK_VINTF_ARGS)" >> $@
	$(foreach odm_sku,$(PRIVATE_ODM_SKUS), $(foreach vendor_sku,$(PRIVATE_VENDOR_SKUS), \
	  echo "For ODM SKU = $(odm_sku), vendor SKU = $(vendor_sku)" >> $@; \
	  ( $< --check-compat $(PRIVATE_CHECK_VINTF_ARGS) \
	       --property ro.boot.product.hardware.sku=$(filter-out EMPTY_ODM_SKU_PLACEHOLDER,$(odm_sku)) \
	       --property ro.boot.product.vendor.sku=$(filter-out EMPTY_VENDOR_SKU_PLACEHOLDER,$(vendor_sku)) \
	       >> $@ 2>&1 ) || (cat $@ && exit 1); ))

check_vintf_compatible_log :=
check_vintf_compatible_args :=
check_vintf_compatible_deps :=

endif # BUILDING_SYSTEM_EXT_IMAGE equals BOARD_USES_SYSTEM_EXTIMAGE
endif # BUILDING_PRODUCT_IMAGE equals BOARD_USES_PRODUCTIMAGE
endif # BUILDING_ODM_IMAGE equals BOARD_USES_ODMIMAGE
endif # check_vintf_has_vendor
endif # check_vintf_has_system
endif # PRODUCT_ENFORCE_VINTF_MANIFEST

# Add all logs of VINTF checks to dist builds
droid_targets: $(check_vintf_all_deps)
$(call dist-for-goals, droid_targets, $(check_vintf_all_deps))

# Helper alias to check all VINTF of current build.
.PHONY: check-vintf-all
check-vintf-all: $(check_vintf_all_deps)
	$(foreach file,$^,echo "$(file)"; cat "$(file)"; echo;)

check_vintf_has_vendor :=
check_vintf_has_system :=
check_vintf_common_srcs :=
check_vintf_all_deps :=
intermediates :=
endif # !TARGET_BUILD_UNBUNDLED

# -----------------------------------------------------------------
# Check image sizes <= size of super partition

ifeq (,$(TARGET_BUILD_UNBUNDLED))

ifeq (true,$(PRODUCT_BUILD_SUPER_PARTITION))

PARTITIONS_AND_OTHER_IN_SUPER := $(BOARD_SUPER_PARTITION_PARTITION_LIST)

# Add the system other image to the misc_info. Because factory ota may install system_other to the super partition.
ifdef BUILDING_SYSTEM_OTHER_IMAGE
PARTITIONS_AND_OTHER_IN_SUPER += system_other
endif # BUILDING_SYSTEM_OTHER_IMAGE

# $(1): misc_info.txt
# #(2): optional log file
define check-all-partition-sizes-target
  mkdir -p $(dir $(1))
  rm -f $(1)
  $(call dump-super-image-info, $(1))
  $(foreach partition,$(PARTITIONS_AND_OTHER_IN_SUPER), \
    echo "$(partition)_image="$(call images-for-partitions,$(partition)) >> $(1);)
  $(CHECK_PARTITION_SIZES) $(if $(2),--logfile $(2),-v) $(1)
endef

check_all_partition_sizes_log := $(call intermediates-dir-for,PACKAGING,check-all-partition-sizes)/check_all_partition_sizes.log
droid_targets: $(check_all_partition_sizes_log)
$(call dist-for-goals, droid_targets, $(check_all_partition_sizes_log))

$(check_all_partition_sizes_log): \
    $(CHECK_PARTITION_SIZES) \
    $(call images-for-partitions,$(PARTITIONS_AND_OTHER_IN_SUPER))
	$(call check-all-partition-sizes-target, \
	  $(call intermediates-dir-for,PACKAGING,check-all-partition-sizes)/misc_info.txt, \
	  $@)

.PHONY: check-all-partition-sizes
check-all-partition-sizes: $(check_all_partition_sizes_log)

.PHONY: check-all-partition-sizes-nodeps
check-all-partition-sizes-nodeps:
	$(call check-all-partition-sizes-target, \
	  $(call intermediates-dir-for,PACKAGING,check-all-partition-sizes-nodeps)/misc_info.txt)

endif # PRODUCT_BUILD_SUPER_PARTITION

endif # !TARGET_BUILD_UNBUNDLED

# -----------------------------------------------------------------
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
# bring in the installer image generation defines if necessary
ifeq ($(TARGET_USE_DISKINSTALLER),true)
include bootable/diskinstaller/config.mk
endif

# -----------------------------------------------------------------
# host tools needed to build dist and OTA packages

build_ota_package := true
ifeq ($(TARGET_SKIP_OTA_PACKAGE),true)
build_ota_package := false
endif
ifeq ($(BUILD_OS),darwin)
<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
build_ota_package := false
endif
ifneq ($(strip $(SANITIZE_TARGET)),)
build_ota_package := false
endif
ifeq ($(TARGET_PRODUCT),sdk)
build_ota_package := false
endif
ifneq ($(filter generic%,$(TARGET_DEVICE)),)
build_ota_package := false
endif
ifeq ($(TARGET_NO_KERNEL),true)
build_ota_package := false
endif
ifeq ($(recovery_fstab),)
build_ota_package := false
endif
ifeq ($(TARGET_BUILD_PDK),true)
build_ota_package := false
=======
  build_ota_package := false
  build_otatools_package := false
else
  # Set build_ota_package, and allow opt-out below.
  build_ota_package := true
  ifeq ($(TARGET_SKIP_OTA_PACKAGE),true)
    build_ota_package := false
  endif
  ifneq (,$(filter address, $(SANITIZE_TARGET)))
    build_ota_package := false
  endif
  ifeq ($(TARGET_PRODUCT),sdk)
    build_ota_package := false
  endif
  # A target without a kernel may be one of the following:
  # - A generic target. In this case, the OTA package usually isn't built.
  #   PRODUCT_BUILD_GENERIC_OTA_PACKAGE may be set to true to force OTA package
  #   generation.
  # - A real device target, with TARGET_NO_KERNEL set to true and
  #   BOARD_PREBUILT_BOOTIMAGE set. In this case, it is valid to generate
  #   an OTA package.
  ifneq ($(PRODUCT_BUILD_GENERIC_OTA_PACKAGE),true)
    ifneq ($(filter generic%,$(TARGET_DEVICE)),)
      build_ota_package := false
    endif
    ifeq ($(INSTALLED_BOOTIMAGE_TARGET),)
      ifeq ($(TARGET_NO_KERNEL),true)
        build_ota_package := false
      endif
    endif # INSTALLED_BOOTIMAGE_TARGET == ""
    ifeq ($(recovery_fstab),)
      build_ota_package := false
    endif
  endif # PRODUCT_BUILD_GENERIC_OTA_PACKAGE

  # Set build_otatools_package, and allow opt-out below.
  build_otatools_package := true
  ifeq ($(TARGET_SKIP_OTATOOLS_PACKAGE),true)
    build_otatools_package := false
  endif
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
endif

ifeq ($(build_ota_package),true)
OTATOOLS :=  $(HOST_OUT_EXECUTABLES)/minigzip \
  $(HOST_OUT_EXECUTABLES)/aapt \
  $(HOST_OUT_EXECUTABLES)/checkvintf \
  $(HOST_OUT_EXECUTABLES)/mkbootfs \
  $(HOST_OUT_EXECUTABLES)/mkbootimg \
  $(HOST_OUT_EXECUTABLES)/fs_config \
  $(HOST_OUT_EXECUTABLES)/zipalign \
  $(HOST_OUT_EXECUTABLES)/bsdiff \
  $(HOST_OUT_EXECUTABLES)/imgdiff \
  $(HOST_OUT_JAVA_LIBRARIES)/dumpkey.jar \
  $(HOST_OUT_JAVA_LIBRARIES)/signapk.jar \
  $(HOST_OUT_JAVA_LIBRARIES)/BootSignature.jar \
  $(HOST_OUT_JAVA_LIBRARIES)/VeritySigner.jar \
  $(HOST_OUT_EXECUTABLES)/mke2fs \
  $(HOST_OUT_EXECUTABLES)/mkuserimg_mke2fs.sh \
  $(HOST_OUT_EXECUTABLES)/e2fsdroid \
  $(HOST_OUT_EXECUTABLES)/mksquashfsimage.sh \
  $(HOST_OUT_EXECUTABLES)/mksquashfs \
  $(HOST_OUT_EXECUTABLES)/mkf2fsuserimg.sh \
  $(HOST_OUT_EXECUTABLES)/make_f2fs \
  $(HOST_OUT_EXECUTABLES)/sload_f2fs \
  $(HOST_OUT_EXECUTABLES)/simg2img \
  $(HOST_OUT_EXECUTABLES)/e2fsck \
  $(HOST_OUT_EXECUTABLES)/build_verity_tree \
  $(HOST_OUT_EXECUTABLES)/generate_verity_key \
  $(HOST_OUT_EXECUTABLES)/verity_signer \
  $(HOST_OUT_EXECUTABLES)/verity_verifier \
  $(HOST_OUT_EXECUTABLES)/append2simg \
  $(HOST_OUT_EXECUTABLES)/img2simg \
  $(HOST_OUT_EXECUTABLES)/boot_signer \
  $(HOST_OUT_EXECUTABLES)/fec \
  $(HOST_OUT_EXECUTABLES)/brillo_update_payload \
  $(HOST_OUT_EXECUTABLES)/lib/shflags/shflags \
  $(HOST_OUT_EXECUTABLES)/delta_generator \
  $(AVBTOOL) \
  $(BLK_ALLOC_TO_BASE_FS) \
  $(BROTLI)

<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
ifeq (true,$(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_SUPPORTS_VBOOT))
OTATOOLS += \
  $(FUTILITY) \
  $(VBOOT_SIGNER)
=======
INTERNAL_OTATOOLS_MODULES := \
  aapt2 \
  add_img_to_target_files \
  apksigner \
  append2simg \
  avbtool \
  blk_alloc_to_base_fs \
  boot_signer \
  brillo_update_payload \
  brotli \
  bsdiff \
  build_image \
  build_super_image \
  build_verity_metadata \
  build_verity_tree \
  care_map_generator \
  check_ota_package_signature \
  check_target_files_signatures \
  check_target_files_vintf \
  checkvintf \
  delta_generator \
  e2fsck \
  e2fsdroid \
  fc_sort \
  fec \
  fsck.f2fs \
  fs_config \
  generate_verity_key \
  host_init_verifier \
  img2simg \
  img_from_target_files \
  imgdiff \
  libconscrypt_openjdk_jni \
  lpmake \
  lpunpack \
  lz4 \
  make_f2fs \
  make_f2fs_casefold \
  merge_target_files \
  minigzip \
  mk_combined_img \
  mkbootfs \
  mkbootimg \
  mke2fs \
  mke2fs.conf \
  mkfs.erofs \
  mkerofsimage.sh \
  mkf2fsuserimg.sh \
  mksquashfs \
  mksquashfsimage.sh \
  mkuserimg_mke2fs \
  ota_from_target_files \
  repack_bootimg \
  secilc \
  sefcontext_compile \
  sgdisk \
  shflags \
  sign_apex \
  sign_target_files_apks \
  sign_virt_apex \
  signapk \
  simg2img \
  sload_f2fs \
  toybox \
  tune2fs \
  unpack_bootimg \
  update_host_simulator \
  validate_target_files \
  verity_signer \
  verity_verifier \
  zipalign \
  zucchini \

# Additional tools to unpack and repack the apex file.
INTERNAL_OTATOOLS_MODULES += \
  apexer \
  apex_compression_tool \
  deapexer \
  debugfs_static \
  merge_zips \
  resize2fs \
  soong_zip \

ifeq (true,$(PRODUCT_SUPPORTS_VBOOT))
INTERNAL_OTATOOLS_MODULES += \
  futility \
  vboot_signer
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
endif

# Shared libraries.
OTATOOLS += \
  $(HOST_LIBRARY_PATH)/libc++$(HOST_SHLIB_SUFFIX) \
  $(HOST_LIBRARY_PATH)/liblog$(HOST_SHLIB_SUFFIX) \
  $(HOST_LIBRARY_PATH)/libcutils$(HOST_SHLIB_SUFFIX) \
  $(HOST_LIBRARY_PATH)/libselinux$(HOST_SHLIB_SUFFIX) \
  $(HOST_LIBRARY_PATH)/libcrypto_utils$(HOST_SHLIB_SUFFIX) \
  $(HOST_LIBRARY_PATH)/libcrypto-host$(HOST_SHLIB_SUFFIX) \
  $(HOST_LIBRARY_PATH)/libext2fs-host$(HOST_SHLIB_SUFFIX) \
  $(HOST_LIBRARY_PATH)/libext2_blkid-host$(HOST_SHLIB_SUFFIX) \
  $(HOST_LIBRARY_PATH)/libext2_com_err-host$(HOST_SHLIB_SUFFIX) \
  $(HOST_LIBRARY_PATH)/libext2_e2p-host$(HOST_SHLIB_SUFFIX) \
  $(HOST_LIBRARY_PATH)/libext2_misc$(HOST_SHLIB_SUFFIX) \
  $(HOST_LIBRARY_PATH)/libext2_profile-host$(HOST_SHLIB_SUFFIX) \
  $(HOST_LIBRARY_PATH)/libext2_quota-host$(HOST_SHLIB_SUFFIX) \
  $(HOST_LIBRARY_PATH)/libext2_uuid-host$(HOST_SHLIB_SUFFIX) \
  $(HOST_LIBRARY_PATH)/libconscrypt_openjdk_jni$(HOST_SHLIB_SUFFIX) \
  $(HOST_LIBRARY_PATH)/libbrillo$(HOST_SHLIB_SUFFIX) \
  $(HOST_LIBRARY_PATH)/libbrillo-stream$(HOST_SHLIB_SUFFIX) \
  $(HOST_LIBRARY_PATH)/libchrome$(HOST_SHLIB_SUFFIX) \
  $(HOST_LIBRARY_PATH)/libcurl-host$(HOST_SHLIB_SUFFIX) \
  $(HOST_LIBRARY_PATH)/libevent-host$(HOST_SHLIB_SUFFIX) \
  $(HOST_LIBRARY_PATH)/libprotobuf-cpp-lite$(HOST_SHLIB_SUFFIX) \
  $(HOST_LIBRARY_PATH)/libssl-host$(HOST_SHLIB_SUFFIX) \
  $(HOST_LIBRARY_PATH)/libz-host$(HOST_SHLIB_SUFFIX) \
  $(HOST_LIBRARY_PATH)/libsparse-host$(HOST_SHLIB_SUFFIX) \
  $(HOST_LIBRARY_PATH)/libbase$(HOST_SHLIB_SUFFIX) \
  $(HOST_LIBRARY_PATH)/libpcre2$(HOST_SHLIB_SUFFIX) \
  $(HOST_LIBRARY_PATH)/libbrotli$(HOST_SHLIB_SUFFIX)


.PHONY: otatools
otatools: $(OTATOOLS)

BUILT_OTATOOLS_PACKAGE := $(PRODUCT_OUT)/otatools.zip
$(BUILT_OTATOOLS_PACKAGE): zip_root := $(call intermediates-dir-for,PACKAGING,otatools)/otatools

OTATOOLS_DEPS := \
  system/extras/verity/build_verity_metadata.py \
  system/extras/ext4_utils/mke2fs.conf \
  $(sort $(shell find external/avb/test/data -type f -name "testkey_*.pem" -o \
      -name "atx_metadata.bin")) \
  $(sort $(shell find system/update_engine/scripts -name \*.pyc -prune -o -type f -print)) \
  $(sort $(shell find build/target/product/security -type f -name \*.x509.pem -o -name \*.pk8 -o \
      -name verity_key)) \
  $(sort $(shell find device $(wildcard vendor) -type f -name \*.pk8 -o -name verifiedboot\* -o \
      -name \*.x509.pem -o -name oem\*.prop))

OTATOOLS_RELEASETOOLS := \
  $(sort $(shell find build/make/tools/releasetools -name \*.pyc -prune -o -type f))

ifeq (true,$(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_SUPPORTS_VBOOT))
OTATOOLS_DEPS += \
  $(sort $(shell find external/vboot_reference/tests/devkeys -type f))
endif

$(BUILT_OTATOOLS_PACKAGE): $(OTATOOLS) $(OTATOOLS_DEPS) $(OTATOOLS_RELEASETOOLS) | $(ACP)
	@echo "Package OTA tools: $@"
	$(hide) rm -rf $@ $(zip_root)
	$(hide) mkdir -p $(dir $@) $(zip_root)/bin $(zip_root)/framework $(zip_root)/releasetools
	$(call copy-files-with-structure,$(OTATOOLS),$(HOST_OUT)/,$(zip_root))
	$(hide) $(ACP) -r -d -p build/make/tools/releasetools/* $(zip_root)/releasetools
	$(hide) rm -rf $@ $(zip_root)/releasetools/*.pyc
	$(hide) (cd $(zip_root) && zip -qryX $(abspath $@) *)
	$(hide) echo $(OTATOOLS_DEPS) | xargs zip -qryX $(abspath $@)>/dev/null || true

.PHONY: otatools-package
otatools-package: $(BUILT_OTATOOLS_PACKAGE)

<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
endif # build_ota_package
=======
endif # build_otatools_package

# -----------------------------------------------------------------
#  misc_info.txt

INSTALLED_MISC_INFO_TARGET := $(PRODUCT_OUT)/misc_info.txt

ifeq ($(TARGET_RELEASETOOLS_EXTENSIONS),)
# default to common dir for device vendor
tool_extensions := $(TARGET_DEVICE_DIR)/../common
else
tool_extensions := $(TARGET_RELEASETOOLS_EXTENSIONS)
endif
.KATI_READONLY := tool_extensions

# $1: boot image file name
define misc_boot_size
$(subst .img,_size,$(1))=$(BOARD_KERNEL$(call to-upper,$(subst boot,,$(subst .img,,$(1))))_BOOTIMAGE_PARTITION_SIZE)
endef

$(INSTALLED_MISC_INFO_TARGET):
	rm -f $@
	$(call pretty,"Target misc_info.txt: $@")
	$(hide) echo "recovery_api_version=$(RECOVERY_API_VERSION)" >> $@
	$(hide) echo "fstab_version=$(RECOVERY_FSTAB_VERSION)" >> $@
ifdef BOARD_FLASH_BLOCK_SIZE
	$(hide) echo "blocksize=$(BOARD_FLASH_BLOCK_SIZE)" >> $@
endif
ifneq ($(strip $(BOARD_BOOTIMAGE_PARTITION_SIZE))$(strip $(BOARD_KERNEL_BINARIES)),)
	$(foreach b,$(INSTALLED_BOOTIMAGE_TARGET),\
		echo "$(call misc_boot_size,$(notdir $(b)))" >> $@;)
endif
ifeq ($(INSTALLED_BOOTIMAGE_TARGET),)
	$(hide) echo "no_boot=true" >> $@
else
	echo "boot_images=$(foreach b,$(INSTALLED_BOOTIMAGE_TARGET),$(notdir $(b)))" >> $@
endif
ifeq ($(BOARD_RAMDISK_USE_LZ4),true)
	echo "lz4_ramdisks=true" >> $@
endif
ifneq ($(INSTALLED_VENDOR_BOOTIMAGE_TARGET),)
	echo "vendor_boot=true" >> $@
	echo "vendor_boot_size=$(BOARD_VENDOR_BOOTIMAGE_PARTITION_SIZE)" >> $@
endif
ifeq ($(INSTALLED_RECOVERYIMAGE_TARGET),)
	$(hide) echo "no_recovery=true" >> $@
endif
ifdef BOARD_INCLUDE_RECOVERY_DTBO
	$(hide) echo "include_recovery_dtbo=true" >> $@
endif
ifdef BOARD_INCLUDE_RECOVERY_ACPIO
	$(hide) echo "include_recovery_acpio=true" >> $@
endif
ifdef BOARD_RECOVERYIMAGE_PARTITION_SIZE
	$(hide) echo "recovery_size=$(BOARD_RECOVERYIMAGE_PARTITION_SIZE)" >> $@
endif
ifdef TARGET_RECOVERY_FSTYPE_MOUNT_OPTIONS
	@# TARGET_RECOVERY_FSTYPE_MOUNT_OPTIONS can be empty to indicate that nothing but defaults should be used.
	$(hide) echo "recovery_mount_options=$(TARGET_RECOVERY_FSTYPE_MOUNT_OPTIONS)" >> $@
else
	$(hide) echo "recovery_mount_options=$(DEFAULT_TARGET_RECOVERY_FSTYPE_MOUNT_OPTIONS)" >> $@
endif
	$(hide) echo "tool_extensions=$(tool_extensions)" >> $@
	$(hide) echo "default_system_dev_certificate=$(DEFAULT_SYSTEM_DEV_CERTIFICATE)" >> $@
ifdef PRODUCT_EXTRA_RECOVERY_KEYS
	$(hide) echo "extra_recovery_keys=$(PRODUCT_EXTRA_RECOVERY_KEYS)" >> $@
endif
	$(hide) echo 'mkbootimg_args=$(BOARD_MKBOOTIMG_ARGS)' >> $@
	$(hide) echo 'recovery_mkbootimg_args=$(BOARD_RECOVERY_MKBOOTIMG_ARGS)' >> $@
	$(hide) echo 'mkbootimg_version_args=$(INTERNAL_MKBOOTIMG_VERSION_ARGS)' >> $@
ifdef BOARD_GKI_SIGNING_KEY_PATH
	$(hide) echo 'gki_signing_key_path=$(BOARD_GKI_SIGNING_KEY_PATH)' >> $@
	$(hide) echo 'gki_signing_algorithm=$(BOARD_GKI_SIGNING_ALGORITHM)' >> $@
endif
ifdef BOARD_GKI_SIGNING_SIGNATURE_ARGS
	$(hide) echo 'gki_signing_signature_args=$(BOARD_GKI_SIGNING_SIGNATURE_ARGS)' >> $@
endif
	$(hide) echo "multistage_support=1" >> $@
	$(hide) echo "blockimgdiff_versions=3,4" >> $@
ifeq ($(PRODUCT_BUILD_GENERIC_OTA_PACKAGE),true)
	$(hide) echo "build_generic_ota_package=true" >> $@
endif
ifneq ($(OEM_THUMBPRINT_PROPERTIES),)
	# OTA scripts are only interested in fingerprint related properties
	$(hide) echo "oem_fingerprint_properties=$(OEM_THUMBPRINT_PROPERTIES)" >> $@
endif
ifneq (,$(filter address, $(SANITIZE_TARGET)))
	# We need to create userdata.img with real data because the instrumented libraries are in userdata.img.
	$(hide) echo "userdata_img_with_data=true" >> $@
endif
ifeq ($(BOARD_USES_FULL_RECOVERY_IMAGE),true)
	$(hide) echo "full_recovery_image=true" >> $@
endif
ifdef BOARD_USES_VENDORIMAGE
	$(hide) echo "board_uses_vendorimage=true" >> $@
endif
ifeq ($(BOARD_AVB_ENABLE),true)
ifeq ($(BUILDING_VBMETA_IMAGE),true)
	$(hide) echo "avb_building_vbmeta_image=true" >> $@
endif # BUILDING_VBMETA_IMAGE
	$(hide) echo "avb_enable=true" >> $@
	$(hide) echo "avb_vbmeta_key_path=$(BOARD_AVB_KEY_PATH)" >> $@
	$(hide) echo "avb_vbmeta_algorithm=$(BOARD_AVB_ALGORITHM)" >> $@
	$(hide) echo "avb_vbmeta_args=$(BOARD_AVB_MAKE_VBMETA_IMAGE_ARGS)" >> $@
	$(hide) echo "avb_boot_add_hash_footer_args=$(BOARD_AVB_BOOT_ADD_HASH_FOOTER_ARGS)" >> $@
ifdef BOARD_AVB_BOOT_KEY_PATH
	$(hide) echo "avb_boot_key_path=$(BOARD_AVB_BOOT_KEY_PATH)" >> $@
	$(hide) echo "avb_boot_algorithm=$(BOARD_AVB_BOOT_ALGORITHM)" >> $@
	$(hide) echo "avb_boot_rollback_index_location=$(BOARD_AVB_BOOT_ROLLBACK_INDEX_LOCATION)" >> $@
endif # BOARD_AVB_BOOT_KEY_PATH
	echo "avb_vendor_boot_add_hash_footer_args=$(BOARD_AVB_VENDOR_BOOT_ADD_HASH_FOOTER_ARGS)" >> $@
ifdef BOARD_AVB_VENDOR_BOOT_KEY_PATH
	echo "avb_vendor_boot_key_path=$(BOARD_AVB_VENDOR_BOOT_KEY_PATH)" >> $@
	echo "avb_vendor_boot_algorithm=$(BOARD_AVB_VENDOR_BOOT_ALGORITHM)" >> $@
	echo "avb_vendor_boot_rollback_index_location=$(BOARD_AVB_VENDOR_BOOT_ROLLBACK_INDEX_LOCATION)" >> $@
endif # BOARD_AVB_VENDOR_BOOT_KEY_PATH
	$(hide) echo "avb_recovery_add_hash_footer_args=$(BOARD_AVB_RECOVERY_ADD_HASH_FOOTER_ARGS)" >> $@
ifdef BOARD_AVB_RECOVERY_KEY_PATH
	$(hide) echo "avb_recovery_key_path=$(BOARD_AVB_RECOVERY_KEY_PATH)" >> $@
	$(hide) echo "avb_recovery_algorithm=$(BOARD_AVB_RECOVERY_ALGORITHM)" >> $@
	$(hide) echo "avb_recovery_rollback_index_location=$(BOARD_AVB_RECOVERY_ROLLBACK_INDEX_LOCATION)" >> $@
endif # BOARD_AVB_RECOVERY_KEY_PATH
ifneq (,$(strip $(BOARD_CUSTOMIMAGES_PARTITION_LIST)))
	$(hide) echo "avb_custom_images_partition_list=$(BOARD_CUSTOMIMAGES_PARTITION_LIST)" >> $@
	$(hide) $(foreach partition,$(BOARD_CUSTOMIMAGES_PARTITION_LIST), \
	    echo "avb_$(partition)_key_path=$(BOARD_AVB_$(call to-upper,$(partition))_KEY_PATH)"  >> $@; \
	    echo "avb_$(partition)_algorithm=$(BOARD_AVB_$(call to-upper,$(partition))_ALGORITHM)"  >> $@; \
	    echo "avb_$(partition)_add_hashtree_footer_args=$(BOARD_AVB_$(call to-upper,$(partition))_ADD_HASHTREE_FOOTER_ARGS)"  >> $@; \
	    echo "avb_$(partition)_rollback_index_location=$(BOARD_AVB_$(call to-upper,$(partition))_ROLLBACK_INDEX_LOCATION)"  >> $@; \
	    echo "avb_$(partition)_partition_size=$(BOARD_AVB_$(call to-upper,$(partition))_PARTITION_SIZE)"  >> $@; \
	    echo "avb_$(partition)_image_list=$(foreach image,$(BOARD_AVB_$(call to-upper,$(partition))_IMAGE_LIST),$(notdir $(image)))" >> $@;)
endif # BOARD_CUSTOMIMAGES_PARTITION_LIST
ifneq (,$(strip $(BOARD_AVB_VBMETA_SYSTEM)))
	$(hide) echo "avb_vbmeta_system=$(BOARD_AVB_VBMETA_SYSTEM)" >> $@
	$(hide) echo "avb_vbmeta_system_args=$(BOARD_AVB_MAKE_VBMETA_SYSTEM_IMAGE_ARGS)" >> $@
	$(hide) echo "avb_vbmeta_system_key_path=$(BOARD_AVB_VBMETA_SYSTEM_KEY_PATH)" >> $@
	$(hide) echo "avb_vbmeta_system_algorithm=$(BOARD_AVB_VBMETA_SYSTEM_ALGORITHM)" >> $@
	$(hide) echo "avb_vbmeta_system_rollback_index_location=$(BOARD_AVB_VBMETA_SYSTEM_ROLLBACK_INDEX_LOCATION)" >> $@
endif # BOARD_AVB_VBMETA_SYSTEM
ifneq (,$(strip $(BOARD_AVB_VBMETA_VENDOR)))
	$(hide) echo "avb_vbmeta_vendor=$(BOARD_AVB_VBMETA_VENDOR)" >> $@
	$(hide) echo "avb_vbmeta_vendor_args=$(BOARD_AVB_MAKE_VBMETA_SYSTEM_IMAGE_ARGS)" >> $@
	$(hide) echo "avb_vbmeta_vendor_key_path=$(BOARD_AVB_VBMETA_VENDOR_KEY_PATH)" >> $@
	$(hide) echo "avb_vbmeta_vendor_algorithm=$(BOARD_AVB_VBMETA_VENDOR_ALGORITHM)" >> $@
	$(hide) echo "avb_vbmeta_vendor_rollback_index_location=$(BOARD_AVB_VBMETA_VENDOR_ROLLBACK_INDEX_LOCATION)" >> $@
endif # BOARD_AVB_VBMETA_VENDOR_KEY_PATH
endif # BOARD_AVB_ENABLE
ifdef BOARD_BPT_INPUT_FILES
	$(hide) echo "board_bpt_enable=true" >> $@
	$(hide) echo "board_bpt_make_table_args=$(BOARD_BPT_MAKE_TABLE_ARGS)" >> $@
	$(hide) echo "board_bpt_input_files=$(BOARD_BPT_INPUT_FILES)" >> $@
endif
ifdef BOARD_BPT_DISK_SIZE
	$(hide) echo "board_bpt_disk_size=$(BOARD_BPT_DISK_SIZE)" >> $@
endif
	$(call generate-userimage-prop-dictionary, $@)
ifeq ($(AB_OTA_UPDATER),true)
	@# Include the build type in META/misc_info.txt so the server can easily differentiate production builds.
	$(hide) echo "build_type=$(TARGET_BUILD_VARIANT)" >> $@
	$(hide) echo "ab_update=true" >> $@
endif
ifeq ($(TARGET_OTA_ALLOW_NON_AB),true)
	$(hide) echo "allow_non_ab=true" >> $@
endif
ifdef BOARD_PREBUILT_DTBOIMAGE
	$(hide) echo "has_dtbo=true" >> $@
ifeq ($(BOARD_AVB_ENABLE),true)
	$(hide) echo "dtbo_size=$(BOARD_DTBOIMG_PARTITION_SIZE)" >> $@
	$(hide) echo "avb_dtbo_add_hash_footer_args=$(BOARD_AVB_DTBO_ADD_HASH_FOOTER_ARGS)" >> $@
ifdef BOARD_AVB_DTBO_KEY_PATH
	$(hide) echo "avb_dtbo_key_path=$(BOARD_AVB_DTBO_KEY_PATH)" >> $@
	$(hide) echo "avb_dtbo_algorithm=$(BOARD_AVB_DTBO_ALGORITHM)" >> $@
	$(hide) echo "avb_dtbo_rollback_index_location=$(BOARD_AVB_DTBO_ROLLBACK_INDEX_LOCATION)" >> $@
endif # BOARD_AVB_DTBO_KEY_PATH
endif # BOARD_AVB_ENABLE
endif # BOARD_PREBUILT_DTBOIMAGE
ifeq ($(BOARD_USES_PVMFWIMAGE),true)
	$(hide) echo "has_pvmfw=true" >> $@
ifeq ($(BOARD_AVB_ENABLE),true)
	$(hide) echo "pvmfw_size=$(BOARD_PVMFWIMAGE_PARTITION_SIZE)" >> $@
	$(hide) echo "avb_pvmfw_add_hash_footer_args=$(BOARD_AVB_PVMFW_ADD_HASH_FOOTER_ARGS)" >> $@
ifdef BOARD_AVB_PVMFW_KEY_PATH
	$(hide) echo "avb_pvmfw_key_path=$(BOARD_AVB_PVMFW_KEY_PATH)" >> $@
	$(hide) echo "avb_pvmfw_algorithm=$(BOARD_AVB_PVMFW_ALGORITHM)" >> $@
	$(hide) echo "avb_pvmfw_rollback_index_location=$(BOARD_AVB_PVMFW_ROLLBACK_INDEX_LOCATION)" >> $@
endif # BOARD_AVB_PVMFW_KEY_PATH
endif # BOARD_AVB_ENABLE
endif # BOARD_USES_PVMFWIMAGE
	$(call dump-dynamic-partitions-info,$@)
	@# VINTF checks
ifeq ($(PRODUCT_ENFORCE_VINTF_MANIFEST),true)
	$(hide) echo "vintf_enforce=true" >> $@
endif
ifdef ODM_MANIFEST_SKUS
	$(hide) echo "vintf_odm_manifest_skus=$(ODM_MANIFEST_SKUS)" >> $@
endif
ifdef ODM_MANIFEST_FILES
	$(hide) echo "vintf_include_empty_odm_sku=true" >> $@
endif
ifdef DEVICE_MANIFEST_SKUS
	$(hide) echo "vintf_vendor_manifest_skus=$(DEVICE_MANIFEST_SKUS)" >> $@
endif
ifdef DEVICE_MANIFEST_FILE
	$(hide) echo "vintf_include_empty_vendor_sku=true" >> $@
endif
ifeq ($(BOARD_BOOTLOADER_IN_UPDATE_PACKAGE),true)
	$(hide) echo "bootloader_in_update_package=true" >> $@
endif
ifeq ($(BOARD_EXCLUDE_KERNEL_FROM_RECOVERY_IMAGE),true)
	$(hide) echo "exclude_kernel_from_recovery_image=true" >> $@
endif
ifneq ($(BOARD_PARTIAL_OTA_UPDATE_PARTITIONS_LIST),)
	$(hide) echo "partial_ota_update_partitions_list=$(BOARD_PARTIAL_OTA_UPDATE_PARTITIONS_LIST)" >> $@
endif
ifeq ($(BUILDING_WITH_VSDK),true)
	$(hide) echo "building_with_vsdk=true" >> $@
endif
ifeq ($(TARGET_FLATTEN_APEX),false)
	$(hide) echo "target_flatten_apex=false" >> $@
endif

.PHONY: misc_info
misc_info: $(INSTALLED_MISC_INFO_TARGET)

droidcore-unbundled: $(INSTALLED_MISC_INFO_TARGET)
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)

# -----------------------------------------------------------------
# A zip of the directories that map to the target filesystem.
# This zip can be used to create an OTA package or filesystem image
# as a post-build step.
#
name := $(TARGET_PRODUCT)
ifeq ($(TARGET_BUILD_TYPE),debug)
  name := $(name)_debug
endif
name := $(name)-target_files-$(FILE_NAME_TAG)

intermediates := $(call intermediates-dir-for,PACKAGING,target_files)
BUILT_TARGET_FILES_PACKAGE := $(intermediates)/$(name).zip
$(BUILT_TARGET_FILES_PACKAGE): intermediates := $(intermediates)
$(BUILT_TARGET_FILES_PACKAGE): \
		zip_root := $(intermediates)/$(name)

# $(1): Directory to copy
# $(2): Location to copy it to
# The "ls -A" is to prevent "acp s/* d" from failing if s is empty.
define package_files-copy-root
  if [ -d "$(strip $(1))" -a "$$(ls -A $(1))" ]; then \
    mkdir -p $(2) && \
    $(ACP) -rd $(strip $(1))/* $(2); \
  fi
endef

built_ota_tools :=

# We can't build static executables when SANITIZE_TARGET=address
ifeq ($(strip $(SANITIZE_TARGET)),)
built_ota_tools += \
    $(call intermediates-dir-for,EXECUTABLES,updater,,,$(TARGET_PREFER_32_BIT))/updater
endif

$(BUILT_TARGET_FILES_PACKAGE): PRIVATE_OTA_TOOLS := $(built_ota_tools)

$(BUILT_TARGET_FILES_PACKAGE): PRIVATE_RECOVERY_API_VERSION := $(RECOVERY_API_VERSION)
$(BUILT_TARGET_FILES_PACKAGE): PRIVATE_RECOVERY_FSTAB_VERSION := $(RECOVERY_FSTAB_VERSION)

ifeq ($(TARGET_RELEASETOOLS_EXTENSIONS),)
# default to common dir for device vendor
tool_extensions := $(TARGET_DEVICE_DIR)/../common
else
tool_extensions := $(TARGET_RELEASETOOLS_EXTENSIONS)
endif
tool_extension := $(wildcard $(tool_extensions)/releasetools.py)
$(BUILT_TARGET_FILES_PACKAGE): PRIVATE_TOOL_EXTENSIONS := $(tool_extensions)
$(BUILT_TARGET_FILES_PACKAGE): PRIVATE_TOOL_EXTENSION := $(tool_extension)

ifeq ($(AB_OTA_UPDATER),true)
updater_dep := system/update_engine/update_engine.conf
else
# Build OTA tools if not using the AB Updater.
updater_dep := $(built_ota_tools)
endif
$(BUILT_TARGET_FILES_PACKAGE): $(updater_dep)

# If we are using recovery as boot, output recovery files to BOOT/.
ifeq ($(BOARD_USES_RECOVERY_AS_BOOT),true)
$(BUILT_TARGET_FILES_PACKAGE): PRIVATE_RECOVERY_OUT := BOOT
else
$(BUILT_TARGET_FILES_PACKAGE): PRIVATE_RECOVERY_OUT := RECOVERY
endif

ifeq ($(AB_OTA_UPDATER),true)
  ifdef BRILLO_VENDOR_PARTITIONS
    $(BUILT_TARGET_FILES_PACKAGE): $(foreach p,$(BRILLO_VENDOR_PARTITIONS),\
                                     $(call word-colon,1,$(p))/$(call word-colon,2,$(p)))
  endif
  ifdef OSRELEASED_DIRECTORY
    $(BUILT_TARGET_FILES_PACKAGE): $(TARGET_OUT_OEM)/$(OSRELEASED_DIRECTORY)/product_id
    $(BUILT_TARGET_FILES_PACKAGE): $(TARGET_OUT_OEM)/$(OSRELEASED_DIRECTORY)/product_version
    $(BUILT_TARGET_FILES_PACKAGE): $(TARGET_OUT_ETC)/$(OSRELEASED_DIRECTORY)/system_version
  endif
endif

# Run fs_config while creating the target files package
# $1: root directory
# $2: add prefix
define fs_config
(cd $(1); find . -type d | sed 's,$$,/,'; find . \! -type d) | cut -c 3- | sort | sed 's,^,$(2),' | $(HOST_OUT_EXECUTABLES)/fs_config -C -D $(TARGET_OUT) -S $(SELINUX_FC) -R "$(2)"
endef

<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
=======
define filter-out-missing-vendor
$(if $(INSTALLED_VENDORIMAGE_TARGET),$(1),$(filter-out vendor,$(1)))
endef
define filter-out-missing-vendor_dlkm
$(if $(INSTALLED_VENDOR_DLKMIMAGE_TARGET),$(1),$(filter-out vendor_dlkm,$(1)))
endef
define filter-out-missing-odm
$(if $(INSTALLED_ODMIMAGE_TARGET),$(1),$(filter-out odm,$(1)))
endef
define filter-out-missing-odm_dlkm
$(if $(INSTALLED_ODM_DLKMIMAGE_TARGET),$(1),$(filter-out odm_dlkm,$(1)))
endef
# Filter out vendor,vendor_dlkm,odm,odm_dlkm from the list for AOSP targets.
# $(1): list
define filter-out-missing-partitions
$(call filter-out-missing-vendor,\
  $(call filter-out-missing-vendor_dlkm,\
    $(call filter-out-missing-odm,\
      $(call filter-out-missing-odm_dlkm,$(1)))))
endef

# Information related to dynamic partitions and virtual A/B. This information
# is needed for building the super image (see dump-super-image-info) and
# building OTA packages.
# $(1): file
define dump-dynamic-partitions-info
  $(if $(filter true,$(PRODUCT_USE_DYNAMIC_PARTITIONS)), \
    echo "use_dynamic_partitions=true" >> $(1))
  $(if $(filter true,$(PRODUCT_RETROFIT_DYNAMIC_PARTITIONS)), \
    echo "dynamic_partition_retrofit=true" >> $(1))
  echo "lpmake=$(notdir $(LPMAKE))" >> $(1)
  $(if $(filter true,$(PRODUCT_BUILD_SUPER_PARTITION)), $(if $(BOARD_SUPER_PARTITION_SIZE), \
    echo "build_super_partition=true" >> $(1)))
  $(if $(BUILDING_SUPER_EMPTY_IMAGE), \
    echo "build_super_empty_partition=true" >> $(1))
  $(if $(filter true,$(BOARD_BUILD_RETROFIT_DYNAMIC_PARTITIONS_OTA_PACKAGE)), \
    echo "build_retrofit_dynamic_partitions_ota_package=true" >> $(1))
  echo "super_metadata_device=$(BOARD_SUPER_PARTITION_METADATA_DEVICE)" >> $(1)
  $(if $(BOARD_SUPER_PARTITION_BLOCK_DEVICES), \
    echo "super_block_devices=$(BOARD_SUPER_PARTITION_BLOCK_DEVICES)" >> $(1))
  $(foreach device,$(BOARD_SUPER_PARTITION_BLOCK_DEVICES), \
    echo "super_$(device)_device_size=$(BOARD_SUPER_PARTITION_$(call to-upper,$(device))_DEVICE_SIZE)" >> $(1);)
  $(if $(BOARD_SUPER_PARTITION_PARTITION_LIST), \
    echo "dynamic_partition_list=$(call filter-out-missing-partitions,$(BOARD_SUPER_PARTITION_PARTITION_LIST))" >> $(1))
  $(if $(BOARD_SUPER_PARTITION_GROUPS),
    echo "super_partition_groups=$(BOARD_SUPER_PARTITION_GROUPS)" >> $(1))
  $(foreach group,$(BOARD_SUPER_PARTITION_GROUPS), \
    echo "super_$(group)_group_size=$(BOARD_$(call to-upper,$(group))_SIZE)" >> $(1); \
    $(if $(BOARD_$(call to-upper,$(group))_PARTITION_LIST), \
      echo "super_$(group)_partition_list=$(call filter-out-missing-partitions,$(BOARD_$(call to-upper,$(group))_PARTITION_LIST))" >> $(1);))
  $(if $(filter true,$(TARGET_USERIMAGES_SPARSE_EXT_DISABLED)), \
    echo "build_non_sparse_super_partition=true" >> $(1))
  $(if $(filter true,$(TARGET_USERIMAGES_SPARSE_F2FS_DISABLED)), \
    echo "build_non_sparse_super_partition=true" >> $(1))
  $(if $(filter true,$(BOARD_SUPER_IMAGE_IN_UPDATE_PACKAGE)), \
    echo "super_image_in_update_package=true" >> $(1))
  $(if $(BOARD_SUPER_PARTITION_SIZE), \
    echo "super_partition_size=$(BOARD_SUPER_PARTITION_SIZE)" >> $(1))
  $(if $(BOARD_SUPER_PARTITION_ALIGNMENT), \
    echo "super_partition_alignment=$(BOARD_SUPER_PARTITION_ALIGNMENT)" >> $(1))
  $(if $(BOARD_SUPER_PARTITION_WARN_LIMIT), \
    echo "super_partition_warn_limit=$(BOARD_SUPER_PARTITION_WARN_LIMIT)" >> $(1))
  $(if $(BOARD_SUPER_PARTITION_ERROR_LIMIT), \
    echo "super_partition_error_limit=$(BOARD_SUPER_PARTITION_ERROR_LIMIT)" >> $(1))
  $(if $(filter true,$(PRODUCT_VIRTUAL_AB_OTA)), \
    echo "virtual_ab=true" >> $(1))
  $(if $(filter true,$(PRODUCT_VIRTUAL_AB_COMPRESSION)), \
    echo "virtual_ab_compression=true" >> $(1))
  $(if $(filter true,$(PRODUCT_VIRTUAL_AB_OTA_RETROFIT)), \
    echo "virtual_ab_retrofit=true" >> $(1))
endef

# By conditionally including the dependency of the target files package on the
# full system image deps, we speed up builds that do not build the system
# image.
ifdef BUILDING_SYSTEM_IMAGE
  $(BUILT_TARGET_FILES_PACKAGE): $(FULL_SYSTEMIMAGE_DEPS)
endif

ifdef BUILDING_USERDATA_IMAGE
  $(BUILT_TARGET_FILES_PACKAGE): $(INTERNAL_USERDATAIMAGE_FILES)
endif

ifdef BUILDING_SYSTEM_OTHER_IMAGE
  $(BUILT_TARGET_FILES_PACKAGE): $(INTERNAL_SYSTEMOTHERIMAGE_FILES)
endif

ifdef BUILDING_VENDOR_BOOT_IMAGE
  $(BUILT_TARGET_FILES_PACKAGE): $(INTERNAL_VENDOR_RAMDISK_FILES)
  $(BUILT_TARGET_FILES_PACKAGE): $(INTERNAL_VENDOR_RAMDISK_FRAGMENT_TARGETS)
  $(BUILT_TARGET_FILES_PACKAGE): $(INTERNAL_VENDOR_BOOTCONFIG_TARGET)
  # The vendor ramdisk may be built from the recovery ramdisk.
  ifeq (true,$(BOARD_MOVE_RECOVERY_RESOURCES_TO_VENDOR_BOOT))
    $(BUILT_TARGET_FILES_PACKAGE): $(INTERNAL_RECOVERY_RAMDISK_FILES_TIMESTAMP)
  endif
endif

ifdef BUILDING_RECOVERY_IMAGE
  # TODO(b/30414428): Can't depend on INTERNAL_RECOVERYIMAGE_FILES alone like other
  # BUILT_TARGET_FILES_PACKAGE dependencies because currently there're cp/rsync/rm
  # commands in build-recoveryimage-target, which would touch the files under
  # TARGET_RECOVERY_OUT and race with packaging target-files.zip.
  ifeq ($(BOARD_USES_RECOVERY_AS_BOOT),true)
    $(BUILT_TARGET_FILES_PACKAGE): $(INSTALLED_BOOTIMAGE_TARGET)
  else
    $(BUILT_TARGET_FILES_PACKAGE): $(INSTALLED_RECOVERYIMAGE_TARGET)
  endif
  $(BUILT_TARGET_FILES_PACKAGE): $(INTERNAL_RECOVERYIMAGE_FILES)
endif

# Conditionally depend on the image files if the image is being built so the
# target-files.zip rule doesn't wait on the image creation rule, or the image
# if it is coming from a prebuilt.

ifdef BUILDING_VENDOR_IMAGE
  $(BUILT_TARGET_FILES_PACKAGE): $(INTERNAL_VENDORIMAGE_FILES)
else ifdef BOARD_PREBUILT_VENDORIMAGE
  $(BUILT_TARGET_FILES_PACKAGE): $(INSTALLED_VENDORIMAGE_TARGET)
endif

ifdef BUILDING_PRODUCT_IMAGE
  $(BUILT_TARGET_FILES_PACKAGE): $(INTERNAL_PRODUCTIMAGE_FILES)
else ifdef BOARD_PREBUILT_PRODUCTIMAGE
  $(BUILT_TARGET_FILES_PACKAGE): $(INSTALLED_PRODUCTIMAGE_TARGET)
endif

ifdef BUILDING_SYSTEM_EXT_IMAGE
  $(BUILT_TARGET_FILES_PACKAGE): $(INTERNAL_SYSTEM_EXTIMAGE_FILES)
else ifdef BOARD_PREBUILT_SYSTEM_EXTIMAGE
  $(BUILT_TARGET_FILES_PACKAGE): $(INSTALLED_SYSTEM_EXTIMAGE_TARGET)
endif

ifdef BUILDING_BOOT_IMAGE
  $(BUILT_TARGET_FILES_PACKAGE): $(INTERNAL_RAMDISK_FILES)
endif
ifneq (,$(INTERNAL_PREBUILT_BOOTIMAGE) $(filter true,$(BOARD_COPY_BOOT_IMAGE_TO_TARGET_FILES)))
  $(BUILT_TARGET_FILES_PACKAGE): $(INSTALLED_BOOTIMAGE_TARGET)
endif

ifdef BUILDING_ODM_IMAGE
  $(BUILT_TARGET_FILES_PACKAGE): $(INTERNAL_ODMIMAGE_FILES)
else ifdef BOARD_PREBUILT_ODMIMAGE
  $(BUILT_TARGET_FILES_PACKAGE): $(INSTALLED_ODMIMAGE_TARGET)
endif

ifdef BUILDING_VENDOR_DLKM_IMAGE
  $(BUILT_TARGET_FILES_PACKAGE): $(INTERNAL_VENDOR_DLKMIMAGE_FILES)
else ifdef BOARD_PREBUILT_VENDOR_DLKMIMAGE
  $(BUILT_TARGET_FILES_PACKAGE): $(INSTALLED_VENDOR_DLKMIMAGE_TARGET)
endif

ifdef BUILDING_ODM_DLKM_IMAGE
  $(BUILT_TARGET_FILES_PACKAGE): $(INTERNAL_ODM_DLKMIMAGE_FILES)
else ifdef BOARD_PREBUILT_ODM_DLKMIMAGE
  $(BUILT_TARGET_FILES_PACKAGE): $(INSTALLED_ODM_DLKMIMAGE_TARGET)
endif

ifeq ($(BUILD_QEMU_IMAGES),true)
  MK_VBMETA_BOOT_KERNEL_CMDLINE_SH := device/generic/goldfish/tools/mk_vbmeta_boot_params.sh
  $(BUILT_TARGET_FILES_PACKAGE): $(MK_VBMETA_BOOT_KERNEL_CMDLINE_SH)
endif

ifdef BOARD_PREBUILT_BOOTLOADER
$(BUILT_TARGET_FILES_PACKAGE): $(INSTALLED_BOOTLOADER_MODULE)
droidcore-unbundled: $(INSTALLED_BOOTLOADER_MODULE)
endif

>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
# Depending on the various images guarantees that the underlying
# directories are up-to-date.
$(BUILT_TARGET_FILES_PACKAGE): \
		$(INSTALLED_BOOTIMAGE_TARGET) \
		$(INSTALLED_RADIOIMAGE_TARGET) \
		$(INSTALLED_RECOVERYIMAGE_TARGET) \
		$(FULL_SYSTEMIMAGE_DEPS) \
		$(INSTALLED_USERDATAIMAGE_TARGET) \
		$(INSTALLED_CACHEIMAGE_TARGET) \
		$(INSTALLED_VENDORIMAGE_TARGET) \
		$(INSTALLED_PRODUCTIMAGE_TARGET) \
		$(INSTALLED_VBMETAIMAGE_TARGET) \
		$(INSTALLED_DTBOIMAGE_TARGET) \
		$(INTERNAL_SYSTEMOTHERIMAGE_FILES) \
		$(INSTALLED_ANDROID_INFO_TXT_TARGET) \
		$(INSTALLED_KERNEL_TARGET) \
		$(INSTALLED_2NDBOOTLOADER_TARGET) \
		$(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_SYSTEM_BASE_FS_PATH) \
		$(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_VENDOR_BASE_FS_PATH) \
		$(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_PRODUCT_BASE_FS_PATH) \
		$(SELINUX_FC) \
		$(APKCERTS_FILE) \
		$(SOONG_ZIP) \
		$(HOST_OUT_EXECUTABLES)/fs_config \
		$(HOST_OUT_EXECUTABLES)/imgdiff \
		$(HOST_OUT_EXECUTABLES)/bsdiff \
		$(BUILD_IMAGE_SRCS) \
		$(BUILT_VENDOR_MANIFEST) \
		$(BUILT_VENDOR_MATRIX) \
		| $(ACP)
	@echo "Package target files: $@"
	$(call create-system-vendor-symlink)
	$(call create-system-product-symlink)
	$(hide) rm -rf $@ $@.list $(zip_root)
	$(hide) mkdir -p $(dir $@) $(zip_root)
ifneq (,$(INSTALLED_RECOVERYIMAGE_TARGET)$(filter true,$(BOARD_USES_RECOVERY_AS_BOOT)))
	@# Components of the recovery image
	$(hide) mkdir -p $(zip_root)/$(PRIVATE_RECOVERY_OUT)
# Exclude recovery files in the default vendor ramdisk if including a standalone
# recovery ramdisk in vendor_boot.
ifneq (true,$(BOARD_INCLUDE_RECOVERY_RAMDISK_IN_VENDOR_BOOT))
	$(hide) $(call package_files-copy-root, \
<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
		$(TARGET_RECOVERY_ROOT_OUT),$(zip_root)/$(PRIVATE_RECOVERY_OUT)/RAMDISK)
=======
	    $(TARGET_RECOVERY_ROOT_OUT),$(zip_root)/$(PRIVATE_RECOVERY_OUT)/RAMDISK)
endif
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
ifdef INSTALLED_KERNEL_TARGET
	$(hide) cp $(INSTALLED_KERNEL_TARGET) $(zip_root)/$(PRIVATE_RECOVERY_OUT)/kernel
endif
ifdef INSTALLED_2NDBOOTLOADER_TARGET
	$(hide) cp $(INSTALLED_2NDBOOTLOADER_TARGET) $(zip_root)/$(PRIVATE_RECOVERY_OUT)/second
endif
ifdef BOARD_INCLUDE_RECOVERY_DTBO
	$(hide) cp $(INSTALLED_DTBOIMAGE_TARGET) $(zip_root)/$(PRIVATE_RECOVERY_OUT)/recovery_dtbo
endif
ifdef INTERNAL_KERNEL_CMDLINE
	$(hide) echo "$(INTERNAL_KERNEL_CMDLINE)" > $(zip_root)/$(PRIVATE_RECOVERY_OUT)/cmdline
endif
ifdef BOARD_KERNEL_BASE
	$(hide) echo "$(BOARD_KERNEL_BASE)" > $(zip_root)/$(PRIVATE_RECOVERY_OUT)/base
endif
ifdef BOARD_KERNEL_PAGESIZE
	$(hide) echo "$(BOARD_KERNEL_PAGESIZE)" > $(zip_root)/$(PRIVATE_RECOVERY_OUT)/pagesize
endif
endif # INSTALLED_RECOVERYIMAGE_TARGET defined or BOARD_USES_RECOVERY_AS_BOOT is true
	@# Components of the boot image
	$(hide) mkdir -p $(zip_root)/BOOT
ifeq ($(BOARD_BUILD_SYSTEM_ROOT_IMAGE),true)
	$(hide) mkdir -p $(zip_root)/ROOT
	$(hide) $(call package_files-copy-root, \
		$(TARGET_ROOT_OUT),$(zip_root)/ROOT)
else
	$(hide) $(call package_files-copy-root, \
		$(TARGET_ROOT_OUT),$(zip_root)/BOOT/RAMDISK)
endif
	@# If we are using recovery as boot, this is already done when processing recovery.
ifneq ($(BOARD_USES_RECOVERY_AS_BOOT),true)
ifdef INSTALLED_KERNEL_TARGET
	$(hide) cp $(INSTALLED_KERNEL_TARGET) $(zip_root)/BOOT/kernel
endif
ifdef INSTALLED_2NDBOOTLOADER_TARGET
	$(hide) cp $(INSTALLED_2NDBOOTLOADER_TARGET) $(zip_root)/BOOT/second
endif
ifdef INTERNAL_KERNEL_CMDLINE
	$(hide) echo "$(INTERNAL_KERNEL_CMDLINE)" > $(zip_root)/BOOT/cmdline
endif
ifdef BOARD_KERNEL_BASE
	$(hide) echo "$(BOARD_KERNEL_BASE)" > $(zip_root)/BOOT/base
endif
ifdef BOARD_KERNEL_PAGESIZE
	$(hide) echo "$(BOARD_KERNEL_PAGESIZE)" > $(zip_root)/BOOT/pagesize
endif
endif # BOARD_USES_RECOVERY_AS_BOOT
	$(hide) $(foreach t,$(INSTALLED_RADIOIMAGE_TARGET),\
	            mkdir -p $(zip_root)/RADIO; \
	            cp $(t) $(zip_root)/RADIO/$(notdir $(t));)
<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
=======
ifdef INSTALLED_VENDOR_BOOTIMAGE_TARGET
	mkdir -p $(zip_root)/VENDOR_BOOT
	$(call package_files-copy-root, \
	    $(TARGET_VENDOR_RAMDISK_OUT),$(zip_root)/VENDOR_BOOT/RAMDISK)
ifdef INSTALLED_DTBIMAGE_TARGET
	cp $(INSTALLED_DTBIMAGE_TARGET) $(zip_root)/VENDOR_BOOT/dtb
endif
ifdef INTERNAL_VENDOR_BOOTCONFIG_TARGET
	cp $(INTERNAL_VENDOR_BOOTCONFIG_TARGET) $(zip_root)/VENDOR_BOOT/vendor_bootconfig
endif
ifdef BOARD_KERNEL_BASE
	echo "$(BOARD_KERNEL_BASE)" > $(zip_root)/VENDOR_BOOT/base
endif
ifdef BOARD_KERNEL_PAGESIZE
	echo "$(BOARD_KERNEL_PAGESIZE)" > $(zip_root)/VENDOR_BOOT/pagesize
endif
	echo "$(INTERNAL_KERNEL_CMDLINE)" > $(zip_root)/VENDOR_BOOT/vendor_cmdline
ifdef INTERNAL_VENDOR_RAMDISK_FRAGMENTS
	echo "$(INTERNAL_VENDOR_RAMDISK_FRAGMENTS)" > "$(zip_root)/VENDOR_BOOT/vendor_ramdisk_fragments"
	$(foreach vendor_ramdisk_fragment,$(INTERNAL_VENDOR_RAMDISK_FRAGMENTS), \
	  mkdir -p $(zip_root)/VENDOR_BOOT/RAMDISK_FRAGMENTS/$(vendor_ramdisk_fragment); \
	  echo "$(BOARD_VENDOR_RAMDISK_FRAGMENT.$(vendor_ramdisk_fragment).MKBOOTIMG_ARGS)" > "$(zip_root)/VENDOR_BOOT/RAMDISK_FRAGMENTS/$(vendor_ramdisk_fragment)/mkbootimg_args"; \
	  $(eval prebuilt_ramdisk := $(BOARD_VENDOR_RAMDISK_FRAGMENT.$(vendor_ramdisk_fragment).PREBUILT)) \
	  $(if $(prebuilt_ramdisk), \
	    cp "$(prebuilt_ramdisk)" "$(zip_root)/VENDOR_BOOT/RAMDISK_FRAGMENTS/$(vendor_ramdisk_fragment)/prebuilt_ramdisk";, \
	    $(call package_files-copy-root, \
	      $(VENDOR_RAMDISK_FRAGMENT.$(vendor_ramdisk_fragment).STAGING_DIR), \
	      $(zip_root)/VENDOR_BOOT/RAMDISK_FRAGMENTS/$(vendor_ramdisk_fragment)/RAMDISK); \
	  ))
endif # INTERNAL_VENDOR_RAMDISK_FRAGMENTS != ""
endif # INSTALLED_VENDOR_BOOTIMAGE_TARGET
ifdef BUILDING_SYSTEM_IMAGE
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
	@# Contents of the system image
	$(hide) $(call package_files-copy-root, \
		$(SYSTEMIMAGE_SOURCE_DIR),$(zip_root)/SYSTEM)
	@# Contents of the data image
	$(hide) $(call package_files-copy-root, \
		$(TARGET_OUT_DATA),$(zip_root)/DATA)
ifdef BOARD_VENDORIMAGE_FILE_SYSTEM_TYPE
	@# Contents of the vendor image
	$(hide) $(call package_files-copy-root, \
		$(TARGET_OUT_VENDOR),$(zip_root)/VENDOR)
endif
ifdef BOARD_PRODUCTIMAGE_FILE_SYSTEM_TYPE
	@# Contents of the product image
	$(hide) $(call package_files-copy-root, \
		$(TARGET_OUT_PRODUCT),$(zip_root)/PRODUCT)
endif
ifdef INSTALLED_SYSTEMOTHERIMAGE_TARGET
	@# Contents of the system_other image
	$(hide) $(call package_files-copy-root, \
		$(TARGET_OUT_SYSTEM_OTHER),$(zip_root)/SYSTEM_OTHER)
endif
	@# Extra contents of the OTA package
	$(hide) mkdir -p $(zip_root)/OTA
	$(hide) cp $(INSTALLED_ANDROID_INFO_TXT_TARGET) $(zip_root)/OTA/
ifneq ($(AB_OTA_UPDATER),true)
ifneq ($(built_ota_tools),)
	$(hide) mkdir -p $(zip_root)/OTA/bin
	$(hide) cp $(PRIVATE_OTA_TOOLS) $(zip_root)/OTA/bin/
endif
endif
	@# Files that do not end up in any images, but are necessary to
	@# build them.
	$(hide) mkdir -p $(zip_root)/META
	$(hide) cp $(APKCERTS_FILE) $(zip_root)/META/apkcerts.txt
ifneq ($(tool_extension),)
	$(hide) cp $(PRIVATE_TOOL_EXTENSION) $(zip_root)/META/
endif
	$(hide) echo "$(PRODUCT_OTA_PUBLIC_KEYS)" > $(zip_root)/META/otakeys.txt
	$(hide) cp $(SELINUX_FC) $(zip_root)/META/file_contexts.bin
	$(hide) echo "recovery_api_version=$(PRIVATE_RECOVERY_API_VERSION)" > $(zip_root)/META/misc_info.txt
	$(hide) echo "fstab_version=$(PRIVATE_RECOVERY_FSTAB_VERSION)" >> $(zip_root)/META/misc_info.txt
ifdef BOARD_FLASH_BLOCK_SIZE
	$(hide) echo "blocksize=$(BOARD_FLASH_BLOCK_SIZE)" >> $(zip_root)/META/misc_info.txt
endif
ifdef BOARD_BOOTIMAGE_PARTITION_SIZE
	$(hide) echo "boot_size=$(BOARD_BOOTIMAGE_PARTITION_SIZE)" >> $(zip_root)/META/misc_info.txt
endif
ifeq ($(INSTALLED_RECOVERYIMAGE_TARGET),)
	$(hide) echo "no_recovery=true" >> $(zip_root)/META/misc_info.txt
endif
ifdef BOARD_INCLUDE_RECOVERY_DTBO
	$(hide) echo "include_recovery_dtbo=true" >> $(zip_root)/META/misc_info.txt
endif
ifdef BOARD_RECOVERYIMAGE_PARTITION_SIZE
	$(hide) echo "recovery_size=$(BOARD_RECOVERYIMAGE_PARTITION_SIZE)" >> $(zip_root)/META/misc_info.txt
endif
ifdef TARGET_RECOVERY_FSTYPE_MOUNT_OPTIONS
	@# TARGET_RECOVERY_FSTYPE_MOUNT_OPTIONS can be empty to indicate that nothing but defaults should be used.
	$(hide) echo "recovery_mount_options=$(TARGET_RECOVERY_FSTYPE_MOUNT_OPTIONS)" >> $(zip_root)/META/misc_info.txt
else
	$(hide) echo "recovery_mount_options=$(DEFAULT_TARGET_RECOVERY_FSTYPE_MOUNT_OPTIONS)" >> $(zip_root)/META/misc_info.txt
endif
	$(hide) echo "tool_extensions=$(PRIVATE_TOOL_EXTENSIONS)" >> $(zip_root)/META/misc_info.txt
	$(hide) echo "default_system_dev_certificate=$(DEFAULT_SYSTEM_DEV_CERTIFICATE)" >> $(zip_root)/META/misc_info.txt
ifdef PRODUCT_EXTRA_RECOVERY_KEYS
	$(hide) echo "extra_recovery_keys=$(PRODUCT_EXTRA_RECOVERY_KEYS)" >> $(zip_root)/META/misc_info.txt
endif
	$(hide) echo 'mkbootimg_args=$(BOARD_MKBOOTIMG_ARGS)' >> $(zip_root)/META/misc_info.txt
	$(hide) echo 'mkbootimg_version_args=$(INTERNAL_MKBOOTIMG_VERSION_ARGS)' >> $(zip_root)/META/misc_info.txt
	$(hide) echo "multistage_support=1" >> $(zip_root)/META/misc_info.txt
	$(hide) echo "blockimgdiff_versions=3,4" >> $(zip_root)/META/misc_info.txt
ifneq ($(OEM_THUMBPRINT_PROPERTIES),)
	# OTA scripts are only interested in fingerprint related properties
	$(hide) echo "oem_fingerprint_properties=$(OEM_THUMBPRINT_PROPERTIES)" >> $(zip_root)/META/misc_info.txt
endif
ifneq ($(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_SYSTEM_BASE_FS_PATH),)
	$(hide) cp $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_SYSTEM_BASE_FS_PATH) \
	  $(zip_root)/META/$(notdir $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_SYSTEM_BASE_FS_PATH))
endif
ifneq ($(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_VENDOR_BASE_FS_PATH),)
	$(hide) cp $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_VENDOR_BASE_FS_PATH) \
	  $(zip_root)/META/$(notdir $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_VENDOR_BASE_FS_PATH))
endif
ifneq ($(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_PRODUCT_BASE_FS_PATH),)
	$(hide) cp $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_PRODUCT_BASE_FS_PATH) \
	  $(zip_root)/META/$(notdir $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_PRODUCT_BASE_FS_PATH))
endif
ifneq ($(strip $(SANITIZE_TARGET)),)
	# We need to create userdata.img with real data because the instrumented libraries are in userdata.img.
	$(hide) echo "userdata_img_with_data=true" >> $(zip_root)/META/misc_info.txt
endif
ifeq ($(BOARD_USES_FULL_RECOVERY_IMAGE),true)
	$(hide) echo "full_recovery_image=true" >> $(zip_root)/META/misc_info.txt
endif
ifeq ($(BOARD_AVB_ENABLE),true)
	$(hide) echo "avb_enable=true" >> $(zip_root)/META/misc_info.txt
	$(hide) echo "avb_vbmeta_key_path=$(BOARD_AVB_KEY_PATH)" >> $(zip_root)/META/misc_info.txt
	$(hide) echo "avb_vbmeta_algorithm=$(BOARD_AVB_ALGORITHM)" >> $(zip_root)/META/misc_info.txt
	$(hide) echo "avb_vbmeta_args=$(BOARD_AVB_MAKE_VBMETA_IMAGE_ARGS)" >> $(zip_root)/META/misc_info.txt
	$(hide) echo "avb_boot_add_hash_footer_args=$(BOARD_AVB_BOOT_ADD_HASH_FOOTER_ARGS)" >> $(zip_root)/META/misc_info.txt
ifdef BOARD_AVB_BOOT_KEY_PATH
	$(hide) echo "avb_boot_key_path=$(BOARD_AVB_BOOT_KEY_PATH)" >> $(zip_root)/META/misc_info.txt
	$(hide) echo "avb_boot_algorithm=$(BOARD_AVB_BOOT_ALGORITHM)" >> $(zip_root)/META/misc_info.txt
	$(hide) echo "avb_boot_rollback_index_location=$(BOARD_AVB_BOOT_ROLLBACK_INDEX_LOCATION)" >> $(zip_root)/META/misc_info.txt
endif # BOARD_AVB_BOOT_KEY_PATH
	$(hide) echo "avb_recovery_add_hash_footer_args=$(BOARD_AVB_RECOVERY_ADD_HASH_FOOTER_ARGS)" >> $(zip_root)/META/misc_info.txt
ifdef BOARD_AVB_RECOVERY_KEY_PATH
	$(hide) echo "avb_recovery_key_path=$(BOARD_AVB_RECOVERY_KEY_PATH)" >> $(zip_root)/META/misc_info.txt
	$(hide) echo "avb_recovery_algorithm=$(BOARD_AVB_RECOVERY_ALGORITHM)" >> $(zip_root)/META/misc_info.txt
	$(hide) echo "avb_recovery_rollback_index_location=$(BOARD_AVB_RECOVERY_ROLLBACK_INDEX_LOCATION)" >> $(zip_root)/META/misc_info.txt
endif # BOARD_AVB_RECOVERY_KEY_PATH
endif # BOARD_AVB_ENABLE
ifdef BOARD_BPT_INPUT_FILES
	$(hide) echo "board_bpt_enable=true" >> $(zip_root)/META/misc_info.txt
	$(hide) echo "board_bpt_make_table_args=$(BOARD_BPT_MAKE_TABLE_ARGS)" >> $(zip_root)/META/misc_info.txt
	$(hide) echo "board_bpt_input_files=$(BOARD_BPT_INPUT_FILES)" >> $(zip_root)/META/misc_info.txt
endif
ifdef BOARD_BPT_DISK_SIZE
	$(hide) echo "board_bpt_disk_size=$(BOARD_BPT_DISK_SIZE)" >> $(zip_root)/META/misc_info.txt
endif
	$(call generate-userimage-prop-dictionary, $(zip_root)/META/misc_info.txt)
ifneq ($(INSTALLED_RECOVERYIMAGE_TARGET),)
	$(hide) PATH=$(foreach p,$(INTERNAL_USERIMAGES_BINARY_PATHS),$(p):)$$PATH MKBOOTIMG=$(MKBOOTIMG) \
	    build/make/tools/releasetools/make_recovery_patch $(zip_root) $(zip_root)
endif
ifeq ($(AB_OTA_UPDATER),true)
	@# When using the A/B updater, include the updater config files in the zip.
	$(hide) cp $(TOPDIR)system/update_engine/update_engine.conf $(zip_root)/META/update_engine_config.txt
	$(hide) for part in $(strip $(AB_OTA_PARTITIONS)); do \
	  echo "$${part}" >> $(zip_root)/META/ab_partitions.txt; \
	done
	$(hide) for conf in $(strip $(AB_OTA_POSTINSTALL_CONFIG)); do \
	  echo "$${conf}" >> $(zip_root)/META/postinstall_config.txt; \
	done
	@# Include the build type in META/misc_info.txt so the server can easily differentiate production builds.
	$(hide) echo "build_type=$(TARGET_BUILD_VARIANT)" >> $(zip_root)/META/misc_info.txt
	$(hide) echo "ab_update=true" >> $(zip_root)/META/misc_info.txt
ifdef BRILLO_VENDOR_PARTITIONS
	$(hide) mkdir -p $(zip_root)/VENDOR_IMAGES
	$(hide) for f in $(BRILLO_VENDOR_PARTITIONS); do \
	  pair1="$$(echo $$f | awk -F':' '{print $$1}')"; \
	  pair2="$$(echo $$f | awk -F':' '{print $$2}')"; \
	  src=$${pair1}/$${pair2}; \
	  dest=$(zip_root)/VENDOR_IMAGES/$${pair2}; \
	  mkdir -p $$(dirname "$${dest}"); \
	  cp $${src} $${dest}; \
	done;
endif
ifdef OSRELEASED_DIRECTORY
	$(hide) cp $(TARGET_OUT_OEM)/$(OSRELEASED_DIRECTORY)/product_id $(zip_root)/META/product_id.txt
	$(hide) cp $(TARGET_OUT_OEM)/$(OSRELEASED_DIRECTORY)/product_version $(zip_root)/META/product_version.txt
	$(hide) cp $(TARGET_OUT_ETC)/$(OSRELEASED_DIRECTORY)/system_version $(zip_root)/META/system_version.txt
endif
endif
ifeq ($(BREAKPAD_GENERATE_SYMBOLS),true)
	@# If breakpad symbols have been generated, add them to the zip.
	$(hide) $(ACP) -r $(TARGET_OUT_BREAKPAD) $(zip_root)/BREAKPAD
endif
# BOARD_BUILD_DISABLED_VBMETAIMAGE is used to build a special vbmeta.img
# that disables AVB verification. The content is fixed and we can just copy
# it to $(zip_root)/IMAGES without passing some info into misc_info.txt for
# regeneration.
ifeq (true,$(BOARD_BUILD_DISABLED_VBMETAIMAGE))
	$(hide) mkdir -p $(zip_root)/IMAGES
	$(hide) cp $(INSTALLED_VBMETAIMAGE_TARGET) $(zip_root)/IMAGES/
endif
ifdef BOARD_PREBUILT_VENDORIMAGE
	$(hide) mkdir -p $(zip_root)/IMAGES
	$(hide) cp $(INSTALLED_VENDORIMAGE_TARGET) $(zip_root)/IMAGES/
endif
ifdef BOARD_PREBUILT_PRODUCTIMAGE
	$(hide) mkdir -p $(zip_root)/IMAGES
	$(hide) cp $(INSTALLED_PRODUCTIMAGE_TARGET) $(zip_root)/IMAGES/
endif
<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
ifdef BOARD_PREBUILT_BOOTIMAGE
=======
ifdef BOARD_PREBUILT_SYSTEM_EXTIMAGE
	$(hide) mkdir -p $(zip_root)/IMAGES
	$(hide) cp $(INSTALLED_SYSTEM_EXTIMAGE_TARGET) $(zip_root)/IMAGES/
endif
ifndef BOARD_PREBUILT_BOOTIMAGE
ifneq (,$(INTERNAL_PREBUILT_BOOTIMAGE) $(filter true,$(BOARD_COPY_BOOT_IMAGE_TO_TARGET_FILES)))
ifdef INSTALLED_BOOTIMAGE_TARGET
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
	$(hide) mkdir -p $(zip_root)/IMAGES
	$(hide) cp $(INSTALLED_BOOTIMAGE_TARGET) $(zip_root)/IMAGES/
<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
=======
endif # INSTALLED_BOOTIMAGE_TARGET
endif # INTERNAL_PREBUILT_BOOTIMAGE != "" || BOARD_COPY_BOOT_IMAGE_TO_TARGET_FILES == true
else # BOARD_PREBUILT_BOOTIMAGE is defined
	$(hide) mkdir -p $(zip_root)/PREBUILT_IMAGES
	$(hide) cp $(INSTALLED_BOOTIMAGE_TARGET) $(zip_root)/PREBUILT_IMAGES/
endif # BOARD_PREBUILT_BOOTIMAGE
ifdef BOARD_PREBUILT_ODMIMAGE
	$(hide) mkdir -p $(zip_root)/IMAGES
	$(hide) cp $(INSTALLED_ODMIMAGE_TARGET) $(zip_root)/IMAGES/
endif
ifdef BOARD_PREBUILT_VENDOR_DLKMIMAGE
	$(hide) mkdir -p $(zip_root)/IMAGES
	$(hide) cp $(INSTALLED_VENDOR_DLKMIMAGE_TARGET) $(zip_root)/IMAGES/
endif
ifdef BOARD_PREBUILT_ODM_DLKMIMAGE
	$(hide) mkdir -p $(zip_root)/IMAGES
	$(hide) cp $(INSTALLED_ODM_DLKMIMAGE_TARGET) $(zip_root)/IMAGES/
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
endif
ifdef BOARD_PREBUILT_DTBOIMAGE
	$(hide) mkdir -p $(zip_root)/PREBUILT_IMAGES
	$(hide) cp $(INSTALLED_DTBOIMAGE_TARGET) $(zip_root)/PREBUILT_IMAGES/
	$(hide) echo "has_dtbo=true" >> $(zip_root)/META/misc_info.txt
ifeq ($(BOARD_AVB_ENABLE),true)
	$(hide) echo "dtbo_size=$(BOARD_DTBOIMG_PARTITION_SIZE)" >> $(zip_root)/META/misc_info.txt
	$(hide) echo "avb_dtbo_add_hash_footer_args=$(BOARD_AVB_DTBO_ADD_HASH_FOOTER_ARGS)" >> $(zip_root)/META/misc_info.txt
ifdef BOARD_AVB_DTBO_KEY_PATH
	$(hide) echo "avb_dtbo_key_path=$(BOARD_AVB_DTBO_KEY_PATH)" >> $(zip_root)/META/misc_info.txt
	$(hide) echo "avb_dtbo_algorithm=$(BOARD_AVB_DTBO_ALGORITHM)" >> $(zip_root)/META/misc_info.txt
	$(hide) echo "avb_dtbo_rollback_index_location=$(BOARD_AVB_DTBO_ROLLBACK_INDEX_LOCATION)" \
	    >> $(zip_root)/META/misc_info.txt
endif # BOARD_AVB_DTBO_KEY_PATH
endif # BOARD_AVB_ENABLE
endif # BOARD_PREBUILT_DTBOIMAGE
<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
=======
ifdef BOARD_PREBUILT_PVMFWIMAGE
	$(hide) mkdir -p $(zip_root)/PREBUILT_IMAGES
	$(hide) cp $(INSTALLED_PVMFWIMAGE_TARGET) $(zip_root)/PREBUILT_IMAGES/
else ifeq ($(BUILDING_PVMFW_IMAGE),true)
	$(hide) mkdir -p $(zip_root)/IMAGES
	$(hide) cp $(INSTALLED_PVMFWIMAGE_TARGET) $(zip_root)/IMAGES/
endif
ifdef BOARD_PREBUILT_BOOTLOADER
	$(hide) mkdir -p $(zip_root)/IMAGES
	$(hide) cp $(INSTALLED_BOOTLOADER_MODULE) $(zip_root)/IMAGES/
endif
ifneq ($(strip $(BOARD_CUSTOMIMAGES_PARTITION_LIST)),)
	$(hide) mkdir -p $(zip_root)/PREBUILT_IMAGES
	$(hide) $(foreach partition,$(BOARD_CUSTOMIMAGES_PARTITION_LIST), \
	    $(foreach image,$(BOARD_AVB_$(call to-upper,$(partition))_IMAGE_LIST),cp $(image) $(zip_root)/PREBUILT_IMAGES/;))
endif # BOARD_CUSTOMIMAGES_PARTITION_LIST
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
	@# The radio images in BOARD_PACK_RADIOIMAGES will be additionally copied from RADIO/ into
	@# IMAGES/, which then will be added into <product>-img.zip. Such images must be listed in
	@# INSTALLED_RADIOIMAGE_TARGET.
	$(hide) $(foreach part,$(BOARD_PACK_RADIOIMAGES), \
	    echo $(part) >> $(zip_root)/META/pack_radioimages.txt;)
	@# Run fs_config on all the system, vendor, boot ramdisk,
	@# and recovery ramdisk files in the zip, and save the output
	$(hide) $(call fs_config,$(zip_root)/SYSTEM,system/) > $(zip_root)/META/filesystem_config.txt
ifdef BOARD_VENDORIMAGE_FILE_SYSTEM_TYPE
	$(hide) $(call fs_config,$(zip_root)/VENDOR,vendor/) > $(zip_root)/META/vendor_filesystem_config.txt
endif
ifdef BOARD_PRODUCTIMAGE_FILE_SYSTEM_TYPE
	$(hide) $(call fs_config,$(zip_root)/PRODUCT,product/) > $(zip_root)/META/product_filesystem_config.txt
endif
ifeq ($(BOARD_BUILD_SYSTEM_ROOT_IMAGE),true)
	@# When using BOARD_BUILD_SYSTEM_ROOT_IMAGE, ROOT always contains the files for the root under
	@# normal boot. BOOT/RAMDISK exists only if additionally using BOARD_USES_RECOVERY_AS_BOOT.
	$(hide) $(call fs_config,$(zip_root)/ROOT,) > $(zip_root)/META/root_filesystem_config.txt
ifeq ($(BOARD_USES_RECOVERY_AS_BOOT),true)
	$(hide) $(call fs_config,$(zip_root)/BOOT/RAMDISK,) > $(zip_root)/META/boot_filesystem_config.txt
endif
else # BOARD_BUILD_SYSTEM_ROOT_IMAGE != true
	$(hide) $(call fs_config,$(zip_root)/BOOT/RAMDISK,) > $(zip_root)/META/boot_filesystem_config.txt
endif
ifneq ($(INSTALLED_RECOVERYIMAGE_TARGET),)
	$(hide) $(call fs_config,$(zip_root)/RECOVERY/RAMDISK,) > $(zip_root)/META/recovery_filesystem_config.txt
endif
ifdef INSTALLED_SYSTEMOTHERIMAGE_TARGET
	$(hide) $(call fs_config,$(zip_root)/SYSTEM_OTHER,system/) > $(zip_root)/META/system_other_filesystem_config.txt
endif
	@# Metadata for compatibility verification.
	$(hide) cp $(BUILT_SYSTEM_MANIFEST) $(zip_root)/META/system_manifest.xml
	$(hide) cp $(BUILT_SYSTEM_COMPATIBILITY_MATRIX) $(zip_root)/META/system_matrix.xml
ifdef BUILT_VENDOR_MANIFEST
	$(hide) cp $(BUILT_VENDOR_MANIFEST) $(zip_root)/META/vendor_manifest.xml
endif
ifdef BUILT_VENDOR_MATRIX
	$(hide) cp $(BUILT_VENDOR_MATRIX) $(zip_root)/META/vendor_matrix.xml
endif

	$(hide) PATH=$(foreach p,$(INTERNAL_USERIMAGES_BINARY_PATHS),$(p):)$$PATH MKBOOTIMG=$(MKBOOTIMG) \
	    build/make/tools/releasetools/add_img_to_target_files -a -v -p $(HOST_OUT) $(zip_root)
	@# Zip everything up, preserving symlinks and placing META/ files first to
	@# help early validation of the .zip file while uploading it.
	$(hide) find $(zip_root)/META | sort >$@.list
	$(hide) find $(zip_root) -path $(zip_root)/META -prune -o -print | sort >>$@.list
	$(hide) $(SOONG_ZIP) -d -o $@ -C $(zip_root) -r $@.list

.PHONY: target-files-package
target-files-package: $(BUILT_TARGET_FILES_PACKAGE)

ifneq ($(filter $(MAKECMDGOALS),target-files-package),)
$(call dist-for-goals, target-files-package, $(BUILT_TARGET_FILES_PACKAGE))
endif

# -----------------------------------------------------------------
# NDK Sysroot Package
NDK_SYSROOT_TARGET := $(PRODUCT_OUT)/ndk_sysroot.tar.bz2
$(NDK_SYSROOT_TARGET): ndk
	@echo Package NDK sysroot...
	$(hide) tar cjf $@ -C $(SOONG_OUT_DIR) ndk

ifeq ($(HOST_OS),linux)
$(call dist-for-goals,sdk,$(NDK_SYSROOT_TARGET))
endif

ifeq ($(build_ota_package),true)
# -----------------------------------------------------------------
# OTA update package

name := $(TARGET_PRODUCT)
ifeq ($(TARGET_BUILD_TYPE),debug)
  name := $(name)_debug
endif
name := $(name)-ota-$(FILE_NAME_TAG)

INTERNAL_OTA_PACKAGE_TARGET := $(PRODUCT_OUT)/$(name).zip

$(INTERNAL_OTA_PACKAGE_TARGET): KEY_CERT_PAIR := $(DEFAULT_KEY_CERT_PAIR)

ifeq ($(AB_OTA_UPDATER),true)
$(INTERNAL_OTA_PACKAGE_TARGET): $(BRILLO_UPDATE_PAYLOAD)
else
$(INTERNAL_OTA_PACKAGE_TARGET): $(BROTLI)
endif

$(INTERNAL_OTA_PACKAGE_TARGET): $(BUILT_TARGET_FILES_PACKAGE) \
		build/make/tools/releasetools/ota_from_target_files
	@echo "Package OTA: $@"
	$(hide) PATH=$(foreach p,$(INTERNAL_USERIMAGES_BINARY_PATHS),$(p):)$$PATH MKBOOTIMG=$(MKBOOTIMG) \
	   build/make/tools/releasetools/ota_from_target_files -v \
	   --block \
	   --extracted_input_target_files $(patsubst %.zip,%,$(BUILT_TARGET_FILES_PACKAGE)) \
	   -p $(HOST_OUT) \
	   -k $(KEY_CERT_PAIR) \
	   $(if $(OEM_OTA_CONFIG), -o $(OEM_OTA_CONFIG)) \
	   $(BUILT_TARGET_FILES_PACKAGE) $@

.PHONY: otapackage
otapackage: $(INTERNAL_OTA_PACKAGE_TARGET)

endif    # build_ota_package

# -----------------------------------------------------------------
# The update package

name := $(TARGET_PRODUCT)
ifeq ($(TARGET_BUILD_TYPE),debug)
  name := $(name)_debug
endif
name := $(name)-img-$(FILE_NAME_TAG)

INTERNAL_UPDATE_PACKAGE_TARGET := $(PRODUCT_OUT)/$(name).zip

$(INTERNAL_UPDATE_PACKAGE_TARGET): $(BUILT_TARGET_FILES_PACKAGE) $(ZIP2ZIP)
	@echo "Package: $@"
	$(hide) $(ZIP2ZIP) -i $(BUILT_TARGET_FILES_PACKAGE) -o $@ \
	   OTA/android-info.txt:android-info.txt "IMAGES/*.img:."

.PHONY: updatepackage
updatepackage: $(INTERNAL_UPDATE_PACKAGE_TARGET)

# The mac build doesn't build dex2oat, so create the zip file only if the build OS is linux.
ifeq ($(BUILD_OS),linux)
ifneq ($(DEX2OAT),)
dexpreopt_tools_deps := $(DEXPREOPT_GEN_DEPS) $(DEXPREOPT_GEN) $(AAPT2)
DEXPREOPT_TOOLS_ZIP := $(PRODUCT_OUT)/dexpreopt_tools.zip
$(DEXPREOPT_TOOLS_ZIP): $(dexpreopt_tools_deps)
$(DEXPREOPT_TOOLS_ZIP): PRIVATE_DEXPREOPT_TOOLS_DEPS := $(dexpreopt_tools_deps)
$(DEXPREOPT_TOOLS_ZIP): $(SOONG_ZIP)
	$(hide) mkdir -p $(dir $@)
	$(hide) $(SOONG_ZIP) -d -o $@ -j $(addprefix -f ,$(PRIVATE_DEXPREOPT_TOOLS_DEPS)) -f $$(realpath $(DEX2OAT))
endif # DEX2OAT is set
endif # BUILD_OS == linux

DEXPREOPT_CONFIG_ZIP := $(PRODUCT_OUT)/dexpreopt_config.zip

$(DEXPREOPT_CONFIG_ZIP): $(INSTALLED_SYSTEMIMAGE_TARGET) \
    $(INSTALLED_VENDORIMAGE_TARGET) \
    $(INSTALLED_ODMIMAGE_TARGET) \
    $(INSTALLED_PRODUCTIMAGE_TARGET) \

ifeq (,$(TARGET_BUILD_UNBUNDLED))
$(DEXPREOPT_CONFIG_ZIP): $(DEX_PREOPT_CONFIG_FOR_MAKE) \
	  $(DEX_PREOPT_SOONG_CONFIG_FOR_MAKE) \

endif

$(DEXPREOPT_CONFIG_ZIP): $(SOONG_ZIP)
	$(hide) mkdir -p $(dir $@) $(PRODUCT_OUT)/dexpreopt_config

ifeq (,$(TARGET_BUILD_UNBUNDLED))
ifneq (,$(DEX_PREOPT_CONFIG_FOR_MAKE))
	$(hide) cp $(DEX_PREOPT_CONFIG_FOR_MAKE) $(PRODUCT_OUT)/dexpreopt_config
endif
ifneq (,$(DEX_PREOPT_SOONG_CONFIG_FOR_MAKE))
	$(hide) cp $(DEX_PREOPT_SOONG_CONFIG_FOR_MAKE) $(PRODUCT_OUT)/dexpreopt_config
endif
endif #!TARGET_BUILD_UNBUNDLED
	$(hide) $(SOONG_ZIP) -d -o $@ -C $(PRODUCT_OUT)/dexpreopt_config -D $(PRODUCT_OUT)/dexpreopt_config

.PHONY: dexpreopt_config_zip
dexpreopt_config_zip: $(DEXPREOPT_CONFIG_ZIP)

# -----------------------------------------------------------------
# A zip of the symbols directory.  Keep the full paths to make it
# more obvious where these files came from.
#
name := $(TARGET_PRODUCT)
ifeq ($(TARGET_BUILD_TYPE),debug)
  name := $(name)_debug
endif
name := $(name)-symbols-$(FILE_NAME_TAG)

SYMBOLS_ZIP := $(PRODUCT_OUT)/$(name).zip
# For apps_only build we'll establish the dependency later in build/make/core/main.mk.
<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
ifndef TARGET_BUILD_APPS
$(SYMBOLS_ZIP): $(INSTALLED_SYSTEMIMAGE) \
		$(INSTALLED_BOOTIMAGE_TARGET) \
		$(INSTALLED_USERDATAIMAGE_TARGET) \
		$(INSTALLED_VENDORIMAGE_TARGET) \
		$(INSTALLED_PRODUCTIMAGE_TARGET) \
		$(updater_dep)
=======
ifeq (,$(TARGET_BUILD_UNBUNDLED))
$(SYMBOLS_ZIP): $(INTERNAL_ALLIMAGES_FILES) $(updater_dep)
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
endif
$(SYMBOLS_ZIP): PRIVATE_LIST_FILE := $(call intermediates-dir-for,PACKAGING,symbols)/filelist
$(SYMBOLS_ZIP): $(SOONG_ZIP)
	@echo "Package symbols: $@"
	$(hide) rm -rf $@ $(PRIVATE_LIST_FILE)
	$(hide) mkdir -p $(dir $@) $(TARGET_OUT_UNSTRIPPED) $(dir $(PRIVATE_LIST_FILE))
	$(hide) find $(TARGET_OUT_UNSTRIPPED) | sort >$(PRIVATE_LIST_FILE)
	$(hide) $(SOONG_ZIP) -d -o $@ -C $(OUT_DIR)/.. -l $(PRIVATE_LIST_FILE)
# -----------------------------------------------------------------
# A zip of the coverage directory.
#
name := $(TARGET_PRODUCT)
ifeq ($(TARGET_BUILD_TYPE),debug)
name := $(name)_debug
endif
name := $(name)-coverage-$(FILE_NAME_TAG)
COVERAGE_ZIP := $(PRODUCT_OUT)/$(name).zip
<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
ifndef TARGET_BUILD_APPS
$(COVERAGE_ZIP): $(INSTALLED_SYSTEMIMAGE) \
		$(INSTALLED_BOOTIMAGE_TARGET) \
		$(INSTALLED_USERDATAIMAGE_TARGET) \
		$(INSTALLED_VENDORIMAGE_TARGET) \
		$(INSTALLED_PRODUCTIMAGE_TARGET)
=======
ifeq (,$(TARGET_BUILD_UNBUNDLED))
$(COVERAGE_ZIP): $(INTERNAL_ALLIMAGES_FILES)
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
endif
$(COVERAGE_ZIP): PRIVATE_LIST_FILE := $(call intermediates-dir-for,PACKAGING,coverage)/filelist
$(COVERAGE_ZIP): $(SOONG_ZIP)
	@echo "Package coverage: $@"
	$(hide) rm -rf $@ $(PRIVATE_LIST_FILE)
	$(hide) mkdir -p $(dir $@) $(TARGET_OUT_COVERAGE) $(dir $(PRIVATE_LIST_FILE))
	$(hide) find $(TARGET_OUT_COVERAGE) | sort >$(PRIVATE_LIST_FILE)
	$(hide) $(SOONG_ZIP) -d -o $@ -C $(TARGET_OUT_COVERAGE) -l $(PRIVATE_LIST_FILE)

<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
=======
#------------------------------------------------------------------
# Export the LLVM profile data tool and dependencies for Clang coverage processing
#
ifeq (true,$(CLANG_COVERAGE))
  LLVM_PROFDATA := $(LLVM_PREBUILTS_BASE)/linux-x86/$(LLVM_PREBUILTS_VERSION)/bin/llvm-profdata
  LIBCXX := $(LLVM_PREBUILTS_BASE)/linux-x86/$(LLVM_PREBUILTS_VERSION)/lib64/libc++.so.1
  PROFDATA_ZIP := $(PRODUCT_OUT)/llvm-profdata.zip
  $(PROFDATA_ZIP): $(SOONG_ZIP)
	$(hide) $(SOONG_ZIP) -d -o $@ -C $(LLVM_PREBUILTS_BASE)/linux-x86/$(LLVM_PREBUILTS_VERSION) -f $(LLVM_PROFDATA) -f $(LIBCXX)

  $(call dist-for-goals,droidcore-unbundled apps_only,$(PROFDATA_ZIP))
endif

>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
# -----------------------------------------------------------------
# A zip of the Android Apps. Not keeping full path so that we don't
# include product names when distributing
#
name := $(TARGET_PRODUCT)
ifeq ($(TARGET_BUILD_TYPE),debug)
  name := $(name)_debug
endif
name := $(name)-apps-$(FILE_NAME_TAG)

APPS_ZIP := $(PRODUCT_OUT)/$(name).zip
$(APPS_ZIP): $(INSTALLED_SYSTEMIMAGE)
	@echo "Package apps: $@"
	$(hide) rm -rf $@
	$(hide) mkdir -p $(dir $@)
	$(hide) apps_to_zip=`find $(TARGET_OUT_APPS) $(TARGET_OUT_APPS_PRIVILEGED) -mindepth 2 -maxdepth 3 -name "*.apk"`; \
	if [ -z "$$apps_to_zip" ]; then \
		echo "No apps to zip up. Generating empty apps archive." ; \
		a=$$(mktemp /tmp/XXXXXXX) && touch $$a && zip $@ $$a && zip -d $@ $$a; \
	else \
		zip -qjX $@ $$apps_to_zip; \
	fi

ifeq (true,$(EMMA_INSTRUMENT))
#------------------------------------------------------------------
# An archive of classes for use in generating code-coverage reports
# These are the uninstrumented versions of any classes that were
# to be instrumented.
# Any dependencies are set up later in build/make/core/main.mk.

JACOCO_REPORT_CLASSES_ALL := $(PRODUCT_OUT)/jacoco-report-classes-all.jar
$(JACOCO_REPORT_CLASSES_ALL): PRIVATE_TARGET_JACOCO_DIR := $(call intermediates-dir-for,PACKAGING,jacoco)
$(JACOCO_REPORT_CLASSES_ALL): PRIVATE_HOST_JACOCO_DIR := $(call intermediates-dir-for,PACKAGING,jacoco,HOST)
$(JACOCO_REPORT_CLASSES_ALL): PRIVATE_TARGET_PROGUARD_USAGE_DIR := $(call intermediates-dir-for,PACKAGING,proguard_usage)
$(JACOCO_REPORT_CLASSES_ALL): PRIVATE_HOST_PROGUARD_USAGE_DIR := $(call intermediates-dir-for,PACKAGING,proguard_usage,HOST)
$(JACOCO_REPORT_CLASSES_ALL) :
	@echo "Collecting uninstrumented classes"
<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
	$(hide) find $(TARGET_COMMON_OUT_ROOT) $(HOST_COMMON_OUT_ROOT) -name "jacoco-report-classes.jar" | \
		zip -@ -0 -q -X $@
# Meaning of these options:
# -@ scan stdin for file paths to add to the zip
# -0 don't do any compression
# -q supress most output
# -X skip storing extended file attributes
=======
	mkdir -p $(PRIVATE_TARGET_JACOCO_DIR) $(PRIVATE_HOST_JACOCO_DIR) $(PRIVATE_TARGET_PROGUARD_USAGE_DIR) $(PRIVATE_HOST_PROGUARD_USAGE_DIR)
	$(SOONG_ZIP) -o $@ -L 0 \
	  -C $(PRIVATE_TARGET_JACOCO_DIR) -P out/target/common/obj -D $(PRIVATE_TARGET_JACOCO_DIR) \
	  -C $(PRIVATE_HOST_JACOCO_DIR) -P out/target/common/obj -D $(PRIVATE_HOST_JACOCO_DIR) \
	  -C $(PRIVATE_TARGET_PROGUARD_USAGE_DIR) -P out/target/common/obj -D $(PRIVATE_TARGET_PROGUARD_USAGE_DIR) \
	  -C $(PRIVATE_HOST_PROGUARD_USAGE_DIR) -P out/target/common/obj -D $(PRIVATE_HOST_PROGUARD_USAGE_DIR)
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)

ifeq (,$(TARGET_BUILD_UNBUNDLED))
  $(JACOCO_REPORT_CLASSES_ALL): $(INTERNAL_ALLIMAGES_FILES)
endif
endif # EMMA_INSTRUMENT=true


#------------------------------------------------------------------
# A zip of Proguard obfuscation dictionary files.
# Only for apps_only build.
#
ifdef TARGET_BUILD_APPS
PROGUARD_DICT_ZIP := $(PRODUCT_OUT)/$(TARGET_PRODUCT)-proguard-dict-$(FILE_NAME_TAG).zip
<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
# the dependency will be set up later in build/make/core/main.mk.
$(PROGUARD_DICT_ZIP) :
=======
# For apps_only build we'll establish the dependency later in build/make/core/main.mk.
ifeq (,$(TARGET_BUILD_UNBUNDLED))
$(PROGUARD_DICT_ZIP): $(INTERNAL_ALLIMAGES_FILES) $(updater_dep)
endif
$(PROGUARD_DICT_ZIP): PRIVATE_PACKAGING_DIR := $(call intermediates-dir-for,PACKAGING,proguard_dictionary)
$(PROGUARD_DICT_ZIP): $(SOONG_ZIP)
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
	@echo "Packaging Proguard obfuscation dictionary files."
<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
	$(hide) dict_files=`find $(TARGET_OUT_COMMON_INTERMEDIATES)/APPS -name proguard_dictionary`; \
		if [ -n "$$dict_files" ]; then \
		  unobfuscated_jars=$${dict_files//proguard_dictionary/classes.jar}; \
		  zip -qX $@ $$dict_files $$unobfuscated_jars; \
		else \
		  touch $(dir $@)/zipdummy; \
		  (cd $(dir $@) && zip -q $(notdir $@) zipdummy); \
		  zip -qd $@ zipdummy; \
		  rm $(dir $@)/zipdummy; \
		fi
=======
	mkdir -p $(dir $@) $(PRIVATE_PACKAGING_DIR)
	$(SOONG_ZIP) --ignore_missing_files -d -o $@ -C $(PRIVATE_PACKAGING_DIR) -P out/target/common/obj -D $(PRIVATE_PACKAGING_DIR)
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)

<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
endif # TARGET_BUILD_APPS
=======
#------------------------------------------------------------------
# A zip of Proguard usage files.
#
PROGUARD_USAGE_ZIP := $(PRODUCT_OUT)/$(TARGET_PRODUCT)-proguard-usage-$(FILE_NAME_TAG).zip
# For apps_only build we'll establish the dependency later in build/make/core/main.mk.
ifeq (,$(TARGET_BUILD_UNBUNDLED))
$(PROGUARD_USAGE_ZIP): \
    $(INSTALLED_SYSTEMIMAGE_TARGET) \
    $(INSTALLED_RAMDISK_TARGET) \
    $(INSTALLED_BOOTIMAGE_TARGET) \
    $(INSTALLED_USERDATAIMAGE_TARGET) \
    $(INSTALLED_VENDORIMAGE_TARGET) \
    $(INSTALLED_PRODUCTIMAGE_TARGET) \
    $(INSTALLED_SYSTEM_EXTIMAGE_TARGET) \
    $(INSTALLED_ODMIMAGE_TARGET) \
    $(INSTALLED_VENDOR_DLKMIMAGE_TARGET) \
    $(INSTALLED_ODM_DLKMIMAGE_TARGET) \
    $(updater_dep)
endif
$(PROGUARD_USAGE_ZIP): PRIVATE_LIST_FILE := $(call intermediates-dir-for,PACKAGING,proguard_usage.zip)/filelist
$(PROGUARD_USAGE_ZIP): PRIVATE_PACKAGING_DIR := $(call intermediates-dir-for,PACKAGING,proguard_usage)
$(PROGUARD_USAGE_ZIP): $(MERGE_ZIPS)
	@echo "Packaging Proguard usage files."
	mkdir -p $(dir $@) $(PRIVATE_PACKAGING_DIR) $(dir $(PRIVATE_LIST_FILE))
	find $(PRIVATE_PACKAGING_DIR) -name proguard_usage.zip > $(PRIVATE_LIST_FILE)
	$(MERGE_ZIPS) $@ @$(PRIVATE_LIST_FILE)

ifeq (true,$(PRODUCT_USE_DYNAMIC_PARTITIONS))

# Dump variables used by build_super_image.py (for building super.img and super_empty.img).
# $(1): output file
define dump-super-image-info
  $(call dump-dynamic-partitions-info,$(1))
  $(if $(filter true,$(AB_OTA_UPDATER)), \
    echo "ab_update=true" >> $(1))
endef

endif # PRODUCT_USE_DYNAMIC_PARTITIONS

# -----------------------------------------------------------------
# super partition image (dist)

ifeq (true,$(PRODUCT_BUILD_SUPER_PARTITION))

# BOARD_SUPER_PARTITION_SIZE must be defined to build super image.
ifneq ($(BOARD_SUPER_PARTITION_SIZE),)

ifneq (true,$(PRODUCT_RETROFIT_DYNAMIC_PARTITIONS))

# For real devices and for dist builds, build super image from target files to an intermediate directory.
INTERNAL_SUPERIMAGE_DIST_TARGET := $(call intermediates-dir-for,PACKAGING,super.img)/super.img
$(INTERNAL_SUPERIMAGE_DIST_TARGET): extracted_input_target_files := $(patsubst %.zip,%,$(BUILT_TARGET_FILES_PACKAGE))
$(INTERNAL_SUPERIMAGE_DIST_TARGET): $(LPMAKE) $(BUILT_TARGET_FILES_PACKAGE) $(BUILD_SUPER_IMAGE)
	$(call pretty,"Target super fs image from target files: $@")
	PATH=$(dir $(LPMAKE)):$$PATH \
	    $(BUILD_SUPER_IMAGE) -v $(extracted_input_target_files) $@

# Skip packing it in dist package because it is in update package.
ifneq (true,$(BOARD_SUPER_IMAGE_IN_UPDATE_PACKAGE))
$(call dist-for-goals,dist_files,$(INTERNAL_SUPERIMAGE_DIST_TARGET))
endif

.PHONY: superimage_dist
superimage_dist: $(INTERNAL_SUPERIMAGE_DIST_TARGET)

endif # PRODUCT_RETROFIT_DYNAMIC_PARTITIONS != "true"
endif # BOARD_SUPER_PARTITION_SIZE != ""
endif # PRODUCT_BUILD_SUPER_PARTITION == "true"

# -----------------------------------------------------------------
# super partition image for development

ifeq (true,$(PRODUCT_BUILD_SUPER_PARTITION))
ifneq ($(BOARD_SUPER_PARTITION_SIZE),)
ifneq (true,$(PRODUCT_RETROFIT_DYNAMIC_PARTITIONS))

# Build super.img by using $(INSTALLED_*IMAGE_TARGET) to $(1)
# $(1): built image path
# $(2): misc_info.txt path; its contents should match expectation of build_super_image.py
define build-superimage-target
  mkdir -p $(dir $(2))
  rm -rf $(2)
  $(call dump-super-image-info,$(2))
  $(foreach p,$(BOARD_SUPER_PARTITION_PARTITION_LIST), \
    echo "$(p)_image=$(INSTALLED_$(call to-upper,$(p))IMAGE_TARGET)" >> $(2);)
  $(if $(BUILDING_SYSTEM_OTHER_IMAGE), $(if $(filter system,$(BOARD_SUPER_PARTITION_PARTITION_LIST)), \
    echo "system_other_image=$(INSTALLED_SYSTEMOTHERIMAGE_TARGET)" >> $(2);))
  mkdir -p $(dir $(1))
  PATH=$(dir $(LPMAKE)):$$PATH \
    $(BUILD_SUPER_IMAGE) -v $(2) $(1)
endef

INSTALLED_SUPERIMAGE_TARGET := $(PRODUCT_OUT)/super.img
INSTALLED_SUPERIMAGE_DEPENDENCIES := $(LPMAKE) $(BUILD_SUPER_IMAGE) \
    $(foreach p, $(BOARD_SUPER_PARTITION_PARTITION_LIST), $(INSTALLED_$(call to-upper,$(p))IMAGE_TARGET))

ifdef BUILDING_SYSTEM_OTHER_IMAGE
ifneq ($(filter system,$(BOARD_SUPER_PARTITION_PARTITION_LIST)),)
INSTALLED_SUPERIMAGE_DEPENDENCIES += $(INSTALLED_SYSTEMOTHERIMAGE_TARGET)
endif
endif

# If BOARD_BUILD_SUPER_IMAGE_BY_DEFAULT is set, super.img is built from images in the
# $(PRODUCT_OUT) directory, and is built to $(PRODUCT_OUT)/super.img. Also, it will
# be built for non-dist builds. This is useful for devices that uses super.img directly, e.g.
# virtual devices.
ifeq (true,$(BOARD_BUILD_SUPER_IMAGE_BY_DEFAULT))
$(INSTALLED_SUPERIMAGE_TARGET): $(INSTALLED_SUPERIMAGE_DEPENDENCIES)
	$(call pretty,"Target super fs image for debug: $@")
	$(call build-superimage-target,$(INSTALLED_SUPERIMAGE_TARGET),\
          $(call intermediates-dir-for,PACKAGING,superimage_debug)/misc_info.txt)

droidcore-unbundled: $(INSTALLED_SUPERIMAGE_TARGET)

# For devices that uses super image directly, the superimage target points to the file in $(PRODUCT_OUT).
.PHONY: superimage
superimage: $(INSTALLED_SUPERIMAGE_TARGET)

$(call dist-for-goals,dist_files,$(INSTALLED_MISC_INFO_TARGET):super_misc_info.txt)
endif # BOARD_BUILD_SUPER_IMAGE_BY_DEFAULT

# Build $(PRODUCT_OUT)/super.img without dependencies.
.PHONY: superimage-nodeps supernod
superimage-nodeps supernod: intermediates :=
superimage-nodeps supernod: | $(INSTALLED_SUPERIMAGE_DEPENDENCIES)
	$(call pretty,"make $(INSTALLED_SUPERIMAGE_TARGET): ignoring dependencies")
	$(call build-superimage-target,$(INSTALLED_SUPERIMAGE_TARGET),\
	  $(call intermediates-dir-for,PACKAGING,superimage-nodeps)/misc_info.txt)

endif # PRODUCT_RETROFIT_DYNAMIC_PARTITIONS != "true"
endif # BOARD_SUPER_PARTITION_SIZE != ""
endif # PRODUCT_BUILD_SUPER_PARTITION == "true"

# -----------------------------------------------------------------
# super empty image
ifdef BUILDING_SUPER_EMPTY_IMAGE

INSTALLED_SUPERIMAGE_EMPTY_TARGET := $(PRODUCT_OUT)/super_empty.img
$(INSTALLED_SUPERIMAGE_EMPTY_TARGET): intermediates := $(call intermediates-dir-for,PACKAGING,super_empty)
$(INSTALLED_SUPERIMAGE_EMPTY_TARGET): $(LPMAKE) $(BUILD_SUPER_IMAGE)
	$(call pretty,"Target empty super fs image: $@")
	mkdir -p $(intermediates)
	rm -rf $(intermediates)/misc_info.txt
	$(call dump-super-image-info,$(intermediates)/misc_info.txt)
	PATH=$(dir $(LPMAKE)):$$PATH \
	    $(BUILD_SUPER_IMAGE) -v $(intermediates)/misc_info.txt $@

$(call dist-for-goals,dist_files,$(INSTALLED_SUPERIMAGE_EMPTY_TARGET))

endif # BUILDING_SUPER_EMPTY_IMAGE


# -----------------------------------------------------------------
# The update package

name := $(TARGET_PRODUCT)
ifeq ($(TARGET_BUILD_TYPE),debug)
  name := $(name)_debug
endif
name := $(name)-img-$(FILE_NAME_TAG)

INTERNAL_UPDATE_PACKAGE_TARGET := $(PRODUCT_OUT)/$(name).zip

$(INTERNAL_UPDATE_PACKAGE_TARGET): $(BUILT_TARGET_FILES_PACKAGE) $(IMG_FROM_TARGET_FILES)
	$(call pretty,"Package: $@")
	PATH=$(INTERNAL_USERIMAGES_BINARY_PATHS):$(dir $(ZIP2ZIP)):$$PATH \
	    $(IMG_FROM_TARGET_FILES) \
	        --additional IMAGES/VerifiedBootParams.textproto:VerifiedBootParams.textproto \
	        $(BUILT_TARGET_FILES_PACKAGE) $@

.PHONY: updatepackage
updatepackage: $(INTERNAL_UPDATE_PACKAGE_TARGET)
$(call dist-for-goals,updatepackage,$(INTERNAL_UPDATE_PACKAGE_TARGET))

>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)

# -----------------------------------------------------------------
# dalvik something
.PHONY: dalvikfiles
dalvikfiles: $(INTERNAL_DALVIK_MODULES)

ifeq ($(BUILD_QEMU_IMAGES),true)
INSTALLED_QEMU_SYSTEMIMAGE := $(PRODUCT_OUT)/system-qemu.img
MK_QEMU_IMAGE_SH := device/generic/goldfish/tools/mk_qemu_image.sh
SGDISK_HOST := $(HOST_OUT_EXECUTABLES)/sgdisk
$(INSTALLED_QEMU_SYSTEMIMAGE): $(INSTALLED_SYSTEMIMAGE) $(MK_QEMU_IMAGE_SH) $(SGDISK_HOST)
	@echo Create system-qemu.img
	(export SGDISK=$(SGDISK_HOST); $(MK_QEMU_IMAGE_SH) ${PRODUCT_OUT}/system.img)

systemimage: $(INSTALLED_QEMU_SYSTEMIMAGE)
<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
droidcore: $(INSTALLED_QEMU_SYSTEMIMAGE)
ifeq ($(BOARD_USES_VENDORIMAGE),true)
=======
droidcore-unbundled: $(INSTALLED_QEMU_SYSTEMIMAGE)
endif
ifdef INSTALLED_VENDORIMAGE_TARGET
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
INSTALLED_QEMU_VENDORIMAGE := $(PRODUCT_OUT)/vendor-qemu.img
$(INSTALLED_QEMU_VENDORIMAGE): $(INSTALLED_VENDORIMAGE_TARGET) $(MK_QEMU_IMAGE_SH) $(SGDISK_HOST)
	@echo Create vendor-qemu.img
	(export SGDISK=$(SGDISK_HOST); $(MK_QEMU_IMAGE_SH) ${PRODUCT_OUT}/vendor.img)

vendorimage: $(INSTALLED_QEMU_VENDORIMAGE)
droidcore-unbundled: $(INSTALLED_QEMU_VENDORIMAGE)
endif
<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
ifeq ($(BOARD_USES_PRODUCTIMAGE),true)
=======

ifdef INSTALLED_RAMDISK_TARGET
ifdef INSTALLED_VENDOR_BOOTIMAGE_TARGET
ifdef INTERNAL_VENDOR_RAMDISK_TARGET
INSTALLED_QEMU_RAMDISKIMAGE := $(PRODUCT_OUT)/ramdisk-qemu.img
$(INSTALLED_QEMU_RAMDISKIMAGE): $(INTERNAL_VENDOR_RAMDISK_TARGET) $(INSTALLED_RAMDISK_TARGET)
	@echo Create ramdisk-qemu.img
	(cat $(INSTALLED_RAMDISK_TARGET) $(INTERNAL_VENDOR_RAMDISK_TARGET) > $(INSTALLED_QEMU_RAMDISKIMAGE))

droidcore-unbundled: $(INSTALLED_QEMU_RAMDISKIMAGE)
endif
endif
endif

ifdef INSTALLED_PRODUCTIMAGE_TARGET
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
INSTALLED_QEMU_PRODUCTIMAGE := $(PRODUCT_OUT)/product-qemu.img
$(INSTALLED_QEMU_PRODUCTIMAGE): $(INSTALLED_PRODUCTIMAGE_TARGET) $(MK_QEMU_IMAGE_SH) $(SGDISK_HOST)
	@echo Create product-qemu.img
	(export SGDISK=$(SGDISK_HOST); $(MK_QEMU_IMAGE_SH) ${PRODUCT_OUT}/product.img)

productimage: $(INSTALLED_QEMU_PRODUCTIMAGE)
droidcore-unbundled: $(INSTALLED_QEMU_PRODUCTIMAGE)
endif
<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
=======
ifdef INSTALLED_SYSTEM_EXTIMAGE_TARGET
INSTALLED_QEMU_SYSTEM_EXTIMAGE := $(PRODUCT_OUT)/system_ext-qemu.img
$(INSTALLED_QEMU_SYSTEM_EXTIMAGE): $(INSTALLED_SYSTEM_EXTIMAGE_TARGET) $(MK_QEMU_IMAGE_SH) $(SGDISK_HOST) $(SIMG2IMG)
	@echo Create system_ext-qemu.img
	(export SGDISK=$(SGDISK_HOST) SIMG2IMG=$(SIMG2IMG); $(MK_QEMU_IMAGE_SH) $(INSTALLED_SYSTEM_EXTIMAGE_TARGET))

systemextimage: $(INSTALLED_QEMU_SYSTEM_EXTIMAGE)
droidcore-unbundled: $(INSTALLED_QEMU_SYSTEM_EXTIMAGE)
endif
ifdef INSTALLED_ODMIMAGE_TARGET
INSTALLED_QEMU_ODMIMAGE := $(PRODUCT_OUT)/odm-qemu.img
$(INSTALLED_QEMU_ODMIMAGE): $(INSTALLED_ODMIMAGE_TARGET) $(MK_QEMU_IMAGE_SH) $(SGDISK_HOST)
	@echo Create odm-qemu.img
	(export SGDISK=$(SGDISK_HOST); $(MK_QEMU_IMAGE_SH) $(INSTALLED_ODMIMAGE_TARGET))

odmimage: $(INSTALLED_QEMU_ODMIMAGE)
droidcore-unbundled: $(INSTALLED_QEMU_ODMIMAGE)
endif

ifdef INSTALLED_VENDOR_DLKMIMAGE_TARGET
INSTALLED_QEMU_VENDOR_DLKMIMAGE := $(PRODUCT_OUT)/vendor_dlkm-qemu.img
$(INSTALLED_QEMU_VENDOR_DLKMIMAGE): $(INSTALLED_VENDOR_DLKMIMAGE_TARGET) $(MK_QEMU_IMAGE_SH) $(SGDISK_HOST)
	@echo Create vendor_dlkm-qemu.img
	(export SGDISK=$(SGDISK_HOST); $(MK_QEMU_IMAGE_SH) $(INSTALLED_VENDOR_DLKMIMAGE_TARGET))

vendor_dlkmimage: $(INSTALLED_QEMU_VENDOR_DLKMIMAGE)
droidcore-unbundled: $(INSTALLED_QEMU_VENDOR_DLKMIMAGE)
endif

ifdef INSTALLED_ODM_DLKMIMAGE_TARGET
INSTALLED_QEMU_ODM_DLKMIMAGE := $(PRODUCT_OUT)/odm_dlkm-qemu.img
$(INSTALLED_QEMU_ODM_DLKMIMAGE): $(INSTALLED_ODM_DLKMIMAGE_TARGET) $(MK_QEMU_IMAGE_SH) $(SGDISK_HOST)
	@echo Create odm_dlkm-qemu.img
	(export SGDISK=$(SGDISK_HOST); $(MK_QEMU_IMAGE_SH) $(INSTALLED_ODM_DLKMIMAGE_TARGET))

odm_dlkmimage: $(INSTALLED_QEMU_ODM_DLKMIMAGE)
droidcore-unbundled: $(INSTALLED_QEMU_ODM_DLKMIMAGE)
endif

QEMU_VERIFIED_BOOT_PARAMS := $(PRODUCT_OUT)/VerifiedBootParams.textproto
$(QEMU_VERIFIED_BOOT_PARAMS): $(INSTALLED_VBMETAIMAGE_TARGET) $(INSTALLED_SYSTEMIMAGE_TARGET) \
    $(MK_VBMETA_BOOT_KERNEL_CMDLINE_SH) $(AVBTOOL)
	@echo Creating $@
	(export AVBTOOL=$(AVBTOOL); $(MK_VBMETA_BOOT_KERNEL_CMDLINE_SH) $(INSTALLED_VBMETAIMAGE_TARGET) \
    $(INSTALLED_SYSTEMIMAGE_TARGET) $(QEMU_VERIFIED_BOOT_PARAMS))

systemimage: $(QEMU_VERIFIED_BOOT_PARAMS)
droidcore-unbundled: $(QEMU_VERIFIED_BOOT_PARAMS)

>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
endif
# -----------------------------------------------------------------
# The emulator package
ifeq ($(BUILD_EMULATOR),true)
INTERNAL_EMULATOR_PACKAGE_FILES += \
        $(HOST_OUT_EXECUTABLES)/emulator$(HOST_EXECUTABLE_SUFFIX) \
        prebuilts/qemu-kernel/$(TARGET_ARCH)/kernel-qemu \
        $(INSTALLED_RAMDISK_TARGET) \
		$(INSTALLED_SYSTEMIMAGE) \
		$(INSTALLED_USERDATAIMAGE_TARGET)

name := $(TARGET_PRODUCT)-emulator-$(FILE_NAME_TAG)

INTERNAL_EMULATOR_PACKAGE_TARGET := $(PRODUCT_OUT)/$(name).zip

$(INTERNAL_EMULATOR_PACKAGE_TARGET): $(INTERNAL_EMULATOR_PACKAGE_FILES)
	@echo "Package: $@"
	$(hide) zip -qjX $@ $(INTERNAL_EMULATOR_PACKAGE_FILES)

endif
# -----------------------------------------------------------------
# Old PDK stuffs, retired
# The pdk package (Platform Development Kit)

#ifneq (,$(filter pdk,$(MAKECMDGOALS)))
#  include development/pdk/Pdk.mk
#endif


# -----------------------------------------------------------------
# The SDK

ifneq ($(filter sdk,$(MAKECMDGOALS)),)

# The SDK includes host-specific components, so it belongs under HOST_OUT.
sdk_dir := $(HOST_OUT)/sdk/$(TARGET_PRODUCT)

# Build a name that looks like:
#
#     linux-x86   --> android-sdk_12345_linux-x86
#     darwin-x86  --> android-sdk_12345_mac-x86
#     windows-x86 --> android-sdk_12345_windows
#
ifneq ($(HOST_OS),linux)
  $(error Building the monolithic SDK is only supported on Linux)
endif
sdk_name := android-sdk_$(FILE_NAME_TAG)
INTERNAL_SDK_HOST_OS_NAME := linux-$(SDK_HOST_ARCH)
sdk_name := $(sdk_name)_$(INTERNAL_SDK_HOST_OS_NAME)

sdk_dep_file := $(sdk_dir)/sdk_deps.mk

ATREE_FILES :=
-include $(sdk_dep_file)

# if we don't have a real list, then use "everything"
ifeq ($(strip $(ATREE_FILES)),)
ATREE_FILES := \
	$(ALL_DEFAULT_INSTALLED_MODULES) \
	$(INSTALLED_RAMDISK_TARGET) \
	$(ALL_DOCS) \
	$(ALL_SDK_FILES)
endif

atree_dir := development/build


sdk_atree_files := $(atree_dir)/sdk.exclude.atree

# development/build/sdk-android-<abi>.atree is used to differentiate
# between architecture models (e.g. ARMv5TE versus ARMv7) when copying
# files like the kernel image. We use TARGET_CPU_ABI because we don't
# have a better way to distinguish between CPU models.
ifneq (,$(strip $(wildcard $(atree_dir)/sdk-android-$(TARGET_CPU_ABI).atree)))
  sdk_atree_files += $(atree_dir)/sdk-android-$(TARGET_CPU_ABI).atree
endif

ifneq ($(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_SDK_ATREE_FILES),)
sdk_atree_files += $(PRODUCTS.$(INTERNAL_PRODUCT).PRODUCT_SDK_ATREE_FILES)
else
sdk_atree_files += $(atree_dir)/sdk.atree
endif

include $(BUILD_SYSTEM)/sdk_font.mk

deps := \
	$(target_notice_file_txt) \
	$(tools_notice_file_txt) \
	$(OUT_DOCS)/offline-sdk-timestamp \
	$(SYMBOLS_ZIP) \
	$(COVERAGE_ZIP) \
	$(INSTALLED_SYSTEMIMAGE) \
	$(INSTALLED_QEMU_SYSTEMIMAGE) \
	$(INSTALLED_QEMU_VENDORIMAGE) \
	$(INSTALLED_USERDATAIMAGE_TARGET) \
	$(INSTALLED_RAMDISK_TARGET) \
	$(INSTALLED_SDK_BUILD_PROP_TARGET) \
	$(INSTALLED_BUILD_PROP_TARGET) \
	$(ATREE_FILES) \
	$(sdk_atree_files) \
	$(HOST_OUT_EXECUTABLES)/atree \
	$(HOST_OUT_EXECUTABLES)/line_endings \
	$(SDK_FONT_DEPS)

INTERNAL_SDK_TARGET := $(sdk_dir)/$(sdk_name).zip
$(INTERNAL_SDK_TARGET): PRIVATE_NAME := $(sdk_name)
$(INTERNAL_SDK_TARGET): PRIVATE_DIR := $(sdk_dir)/$(sdk_name)
$(INTERNAL_SDK_TARGET): PRIVATE_DEP_FILE := $(sdk_dep_file)
$(INTERNAL_SDK_TARGET): PRIVATE_INPUT_FILES := $(sdk_atree_files)

# Set SDK_GNU_ERROR to non-empty to fail when a GNU target is built.
#
#SDK_GNU_ERROR := true

$(INTERNAL_SDK_TARGET): $(deps)
	@echo "Package SDK: $@"
	$(hide) rm -rf $(PRIVATE_DIR) $@
	$(hide) for f in $(strip $(target_gnu_MODULES)); do \
	  if [ -f $$f ]; then \
	    echo SDK: $(if $(SDK_GNU_ERROR),ERROR:,warning:) \
	        including GNU target $$f >&2; \
	    FAIL=$(SDK_GNU_ERROR); \
	  fi; \
	done; \
	if [ $$FAIL ]; then exit 1; fi
	$(hide) echo $(notdir $(SDK_FONT_DEPS)) | tr " " "\n"  > $(SDK_FONT_TEMP)/fontsInSdk.txt
	$(hide) ( \
<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
		ATREE_STRIP="strip -x" \
		$(HOST_OUT_EXECUTABLES)/atree \
		$(addprefix -f ,$(PRIVATE_INPUT_FILES)) \
			-m $(PRIVATE_DEP_FILE) \
			-I . \
			-I $(PRODUCT_OUT) \
			-I $(HOST_OUT) \
			-I $(TARGET_COMMON_OUT_ROOT) \
			-v "PLATFORM_NAME=android-$(PLATFORM_VERSION)" \
			-v "OUT_DIR=$(OUT_DIR)" \
			-v "HOST_OUT=$(HOST_OUT)" \
			-v "TARGET_ARCH=$(TARGET_ARCH)" \
			-v "TARGET_CPU_ABI=$(TARGET_CPU_ABI)" \
			-v "DLL_EXTENSION=$(HOST_SHLIB_SUFFIX)" \
			-v "FONT_OUT=$(SDK_FONT_TEMP)" \
			-o $(PRIVATE_DIR) && \
		cp -f $(target_notice_file_txt) \
				$(PRIVATE_DIR)/system-images/android-$(PLATFORM_VERSION)/$(TARGET_CPU_ABI)/NOTICE.txt && \
		cp -f $(tools_notice_file_txt) $(PRIVATE_DIR)/platform-tools/NOTICE.txt && \
		HOST_OUT_EXECUTABLES=$(HOST_OUT_EXECUTABLES) HOST_OS=$(HOST_OS) \
			development/build/tools/sdk_clean.sh $(PRIVATE_DIR) && \
		chmod -R ug+rwX $(PRIVATE_DIR) && \
		cd $(dir $@) && zip -rqX $(notdir $@) $(PRIVATE_NAME) \
=======
	    ATREE_STRIP="$(HOST_STRIP) -x" \
	    $(HOST_OUT_EXECUTABLES)/atree \
	    $(addprefix -f ,$(PRIVATE_INPUT_FILES)) \
	        -m $(PRIVATE_DEP_FILE) \
	        -I . \
	        -I $(PRODUCT_OUT) \
	        -I $(HOST_OUT) \
	        -I $(TARGET_COMMON_OUT_ROOT) \
	        -v "PLATFORM_NAME=android-$(PLATFORM_VERSION)" \
	        -v "OUT_DIR=$(OUT_DIR)" \
	        -v "HOST_OUT=$(HOST_OUT)" \
	        -v "TARGET_ARCH=$(TARGET_ARCH)" \
	        -v "TARGET_CPU_ABI=$(TARGET_CPU_ABI)" \
	        -v "DLL_EXTENSION=$(HOST_SHLIB_SUFFIX)" \
	        -v "FONT_OUT=$(SDK_FONT_TEMP)" \
	        -o $(PRIVATE_DIR) && \
	    cp -f $(target_notice_file_txt) \
	            $(PRIVATE_DIR)/system-images/android-$(PLATFORM_VERSION)/$(TARGET_CPU_ABI)/NOTICE.txt && \
	    HOST_OUT_EXECUTABLES=$(HOST_OUT_EXECUTABLES) HOST_OS=$(HOST_OS) \
	        development/build/tools/sdk_clean.sh $(PRIVATE_DIR) && \
	    chmod -R ug+rwX $(PRIVATE_DIR) && \
	    cd $(dir $@) && zip -rqX $(notdir $@) $(PRIVATE_NAME) \
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
	) || ( rm -rf $(PRIVATE_DIR) $@ && exit 44 )

MAIN_SDK_DIR  := $(sdk_dir)
MAIN_SDK_ZIP  := $(INTERNAL_SDK_TARGET)

endif # sdk in MAKECMDGOALS

# -----------------------------------------------------------------
# Findbugs
INTERNAL_FINDBUGS_XML_TARGET := $(PRODUCT_OUT)/findbugs.xml
INTERNAL_FINDBUGS_HTML_TARGET := $(PRODUCT_OUT)/findbugs.html
$(INTERNAL_FINDBUGS_XML_TARGET): $(ALL_FINDBUGS_FILES)
	@echo UnionBugs: $@
	$(hide) $(FINDBUGS_DIR)/unionBugs $(ALL_FINDBUGS_FILES) \
	> $@
$(INTERNAL_FINDBUGS_HTML_TARGET): $(INTERNAL_FINDBUGS_XML_TARGET)
	@echo ConvertXmlToText: $@
	$(hide) $(FINDBUGS_DIR)/convertXmlToText -html:fancy.xsl \
	$(INTERNAL_FINDBUGS_XML_TARGET) > $@

# -----------------------------------------------------------------
# Findbugs

# -----------------------------------------------------------------
# These are some additional build tasks that need to be run.
ifneq ($(dont_bother),true)
include $(sort $(wildcard $(BUILD_SYSTEM)/tasks/*.mk))
-include $(sort $(wildcard vendor/*/build/tasks/*.mk))
-include $(sort $(wildcard device/*/build/tasks/*.mk))
-include $(sort $(wildcard product/*/build/tasks/*.mk))
# Also the project-specific tasks
-include $(sort $(wildcard vendor/*/*/build/tasks/*.mk))
-include $(sort $(wildcard device/*/*/build/tasks/*.mk))
-include $(sort $(wildcard product/*/*/build/tasks/*.mk))
# Also add test specifc tasks
include $(sort $(wildcard platform_testing/build/tasks/*.mk))
include $(sort $(wildcard test/vts/tools/build/tasks/*.mk))
endif

include $(BUILD_SYSTEM)/product-graph.mk

# -----------------------------------------------------------------
# Create SDK repository packages. Must be done after tasks/* since
# we need the addon rules defined.
ifneq ($(sdk_repo_goal),)
include $(TOPDIR)development/build/tools/sdk_repo.mk
endif

<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
#------------------------------------------------------------------
# Find lsdump paths
FIND_LSDUMPS_FILE := $(PRODUCT_OUT)/lsdump_paths.txt
$(FIND_LSDUMPS_FILE) : $(LSDUMP_PATHS)
	$(hide) rm -rf $@ && echo "$^" > $@
=======
# -----------------------------------------------------------------
# Soong generates the list of all shared libraries that are depended on by fuzz
# targets. It saves this list as a source:destination pair to
# FUZZ_TARGET_SHARED_DEPS_INSTALL_PAIRS, where the source is the path to the
# build of the unstripped shared library, and the destination is the
# /data/fuzz/$ARCH/lib (for device) or /fuzz/$ARCH/lib (for host) directory
# where fuzz target shared libraries are to be "reinstalled". The
# copy-many-files below generates the rules to copy the unstripped shared
# libraries to the device or host "reinstallation" directory. These rules are
# depended on by each module in soong_cc_prebuilt.mk, where the module will have
# a dependency on each shared library that it needs to be "reinstalled".
FUZZ_SHARED_DEPS := $(call copy-many-files,$(strip $(FUZZ_TARGET_SHARED_DEPS_INSTALL_PAIRS)))

# -----------------------------------------------------------------
# The rule to build all fuzz targets for C++ and Rust, and package them.
# Note: The packages are created in Soong, and in a perfect world,
# we'd be able to create the phony rule there. But, if we want to
# have dist goals for the fuzz target, we need to have the PHONY
# target defined in make. MakeVarsContext.DistForGoal doesn't take
# into account that a PHONY rule create by Soong won't be available
# during make, and such will fail with `writing to readonly
# directory`, because kati will see 'haiku' as being a file, not a
# phony target.
.PHONY: haiku
haiku: $(SOONG_FUZZ_PACKAGING_ARCH_MODULES) $(ALL_FUZZ_TARGETS)
$(call dist-for-goals,haiku,$(SOONG_FUZZ_PACKAGING_ARCH_MODULES))

.PHONY: haiku-rust
haiku-rust: $(SOONG_RUST_FUZZ_PACKAGING_ARCH_MODULES) $(ALL_RUST_FUZZ_TARGETS)
$(call dist-for-goals,haiku-rust,$(SOONG_RUST_FUZZ_PACKAGING_ARCH_MODULES))
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
