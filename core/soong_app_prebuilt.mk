# App prebuilt coming from Soong.
# Extra inputs:
<<<<<<< HEAD   (4be654 Merge "Merge empty history for sparse-7121469-L4290000080720)
=======
# LOCAL_SOONG_BUILT_INSTALLED
# LOCAL_SOONG_BUNDLE
# LOCAL_SOONG_CLASSES_JAR
# LOCAL_SOONG_DEX_JAR
# LOCAL_SOONG_HEADER_JAR
# LOCAL_SOONG_JACOCO_REPORT_CLASSES_JAR
# LOCAL_SOONG_PROGUARD_DICT
# LOCAL_SOONG_PROGUARD_USAGE
>>>>>>> BRANCH (fe6ad7 Merge "Version bump to RBT1.210107.001.A1 [core/build_id.mk])
# LOCAL_SOONG_RESOURCE_EXPORT_PACKAGE

ifneq ($(LOCAL_MODULE_MAKEFILE),$(SOONG_ANDROID_MK))
  $(call pretty-error,soong_app_prebuilt.mk may only be used from Soong)
endif

LOCAL_MODULE_SUFFIX := .apk
LOCAL_BUILT_MODULE_STEM := package.apk

#######################################
include $(BUILD_SYSTEM)/base_rules.mk
#######################################

full_classes_jar := $(intermediates.COMMON)/classes.jar
full_classes_pre_proguard_jar := $(intermediates.COMMON)/classes-pre-proguard.jar
full_classes_header_jar := $(intermediates.COMMON)/classes-header.jar

<<<<<<< HEAD   (4be654 Merge "Merge empty history for sparse-7121469-L4290000080720)
$(eval $(call copy-one-file,$(LOCAL_SOONG_CLASSES_JAR),$(full_classes_jar)))
$(eval $(call copy-one-file,$(LOCAL_SOONG_CLASSES_JAR),$(full_classes_pre_proguard_jar)))
=======
#######################################
include $(BUILD_SYSTEM)/base_rules.mk
#######################################

ifdef LOCAL_SOONG_CLASSES_JAR
  $(eval $(call copy-one-file,$(LOCAL_SOONG_CLASSES_JAR),$(full_classes_jar)))
  $(eval $(call copy-one-file,$(LOCAL_SOONG_CLASSES_JAR),$(full_classes_pre_proguard_jar)))
  $(eval $(call add-dependency,$(LOCAL_BUILT_MODULE),$(full_classes_jar)))

  ifneq ($(TURBINE_ENABLED),false)
    ifdef LOCAL_SOONG_HEADER_JAR
      $(eval $(call copy-one-file,$(LOCAL_SOONG_HEADER_JAR),$(full_classes_header_jar)))
    else
      $(eval $(call copy-one-file,$(full_classes_jar),$(full_classes_header_jar)))
    endif
  endif # TURBINE_ENABLED != false

  javac-check : $(full_classes_jar)
  javac-check-$(LOCAL_MODULE) : $(full_classes_jar)
  .PHONY: javac-check-$(LOCAL_MODULE)
endif

# Run veridex on product, system_ext and vendor modules.
# We skip it for unbundled app builds where we cannot build veridex.
module_run_appcompat :=
ifeq (true,$(non_system_module))
ifeq (,$(TARGET_BUILD_APPS))  # ! unbundled app build
ifneq ($(UNSAFE_DISABLE_HIDDENAPI_FLAGS),true)
  module_run_appcompat := true
endif
endif
endif

ifeq ($(module_run_appcompat),true)
  $(LOCAL_BUILT_MODULE): $(appcompat-files)
  $(LOCAL_BUILT_MODULE): PRIVATE_INSTALLED_MODULE := $(LOCAL_INSTALLED_MODULE)
  $(LOCAL_BUILT_MODULE): $(LOCAL_PREBUILT_MODULE_FILE)
	@echo "Copy: $@"
	$(copy-file-to-target)
	$(appcompat-header)
	$(run-appcompat)
else
  $(eval $(call copy-one-file,$(LOCAL_PREBUILT_MODULE_FILE),$(LOCAL_BUILT_MODULE)))
endif
>>>>>>> BRANCH (fe6ad7 Merge "Version bump to RBT1.210107.001.A1 [core/build_id.mk])

ifdef LOCAL_SOONG_JACOCO_REPORT_CLASSES_JAR
  $(eval $(call copy-one-file,$(LOCAL_SOONG_JACOCO_REPORT_CLASSES_JAR),\
    $(intermediates.COMMON)/jacoco-report-classes.jar))
  $(call add-dependency,$(LOCAL_BUILT_MODULE),\
    $(intermediates.COMMON)/jacoco-report-classes.jar)
endif

ifdef LOCAL_SOONG_PROGUARD_DICT
  $(eval $(call copy-one-file,$(LOCAL_SOONG_PROGUARD_DICT),\
    $(intermediates.COMMON)/proguard_dictionary))
  $(call add-dependency,$(LOCAL_BUILT_MODULE),\
    $(intermediates.COMMON)/proguard_dictionary)
endif

<<<<<<< HEAD   (4be654 Merge "Merge empty history for sparse-7121469-L4290000080720)
ifneq ($(TURBINE_ENABLED),false)
ifdef LOCAL_SOONG_HEADER_JAR
$(eval $(call copy-one-file,$(LOCAL_SOONG_HEADER_JAR),$(full_classes_header_jar)))
else
$(eval $(call copy-one-file,$(full_classes_jar),$(full_classes_header_jar)))
endif
endif # TURBINE_ENABLED != false


$(eval $(call copy-one-file,$(LOCAL_PREBUILT_MODULE_FILE),$(LOCAL_BUILT_MODULE)))
=======
ifdef LOCAL_SOONG_PROGUARD_USAGE_ZIP
  $(eval $(call copy-one-file,$(LOCAL_SOONG_PROGUARD_USAGE_ZIP),\
    $(intermediates.COMMON)/proguard_usage.zip))
  $(call add-dependency,$(LOCAL_BUILT_MODULE),\
    $(intermediates.COMMON)/proguard_usage.zip)
endif
>>>>>>> BRANCH (fe6ad7 Merge "Version bump to RBT1.210107.001.A1 [core/build_id.mk])

ifdef LOCAL_SOONG_RESOURCE_EXPORT_PACKAGE
resource_export_package := $(intermediates.COMMON)/package-export.apk
resource_export_stamp := $(intermediates.COMMON)/src/R.stamp

$(resource_export_package): PRIVATE_STAMP := $(resource_export_stamp)
$(resource_export_package): .KATI_IMPLICIT_OUTPUTS := $(resource_export_stamp)
$(resource_export_package): $(LOCAL_SOONG_RESOURCE_EXPORT_PACKAGE)
	@echo "Copy: $$@"
	$(copy-file-to-target)
	touch $(PRIVATE_STAMP)
$(call add-dependency,$(LOCAL_BUILT_MODULE),$(resource_export_package))

endif # LOCAL_SOONG_RESOURCE_EXPORT_PACKAGE

java-dex: $(LOCAL_SOONG_DEX_JAR)

ifdef LOCAL_DEX_PREOPT
# defines built_odex along with rule to install odex
include $(BUILD_SYSTEM)/dex_preopt_odex_install.mk

$(built_odex): $(LOCAL_SOONG_DEX_JAR)
	$(call dexpreopt-one-file,$<,$@)
endif

PACKAGES := $(PACKAGES) $(LOCAL_MODULE)
ifdef LOCAL_CERTIFICATE
  PACKAGES.$(LOCAL_MODULE).CERTIFICATE := $(LOCAL_CERTIFICATE)
  PACKAGES.$(LOCAL_MODULE).PRIVATE_KEY := $(patsubst %.x509.pem,%.pk8,$(LOCAL_CERTIFICATE))
endif
<<<<<<< HEAD   (4be654 Merge "Merge empty history for sparse-7121469-L4290000080720)
=======
include $(BUILD_SYSTEM)/app_certificate_validate.mk
PACKAGES.$(LOCAL_MODULE).OVERRIDES := $(strip $(LOCAL_OVERRIDES_PACKAGES))

ifneq ($(LOCAL_MODULE_STEM),)
  PACKAGES.$(LOCAL_MODULE).STEM := $(LOCAL_MODULE_STEM)
else
  PACKAGES.$(LOCAL_MODULE).STEM := $(LOCAL_MODULE)
endif

# Set a actual_partition_tag (calculated in base_rules.mk) for the package.
PACKAGES.$(LOCAL_MODULE).PARTITION := $(actual_partition_tag)

ifdef LOCAL_SOONG_BUNDLE
  ALL_MODULES.$(my_register_name).BUNDLE := $(LOCAL_SOONG_BUNDLE)
endif

ifdef LOCAL_SOONG_LINT_REPORTS
  ALL_MODULES.$(my_register_name).LINT_REPORTS := $(LOCAL_SOONG_LINT_REPORTS)
endif
>>>>>>> BRANCH (fe6ad7 Merge "Version bump to RBT1.210107.001.A1 [core/build_id.mk])

ifndef LOCAL_IS_HOST_MODULE
ifeq ($(LOCAL_SDK_VERSION),system_current)
my_link_type := java:system
else ifneq ($(LOCAL_SDK_VERSION),)
my_link_type := java:sdk
else
my_link_type := java:platform
endif
# warn/allowed types are both empty because Soong modules can't depend on
# make-defined modules.
my_warn_types :=
my_allowed_types :=

my_link_deps :=
my_2nd_arch_prefix := $(LOCAL_2ND_ARCH_VAR_PREFIX)
my_common := COMMON
include $(BUILD_SYSTEM)/link_type.mk
endif # !LOCAL_IS_HOST_MODULE

ifdef LOCAL_SOONG_RRO_DIRS
  $(call append_enforce_rro_sources, \
      $(my_register_name), \
      false, \
      $(LOCAL_FULL_MANIFEST_FILE), \
      $(LOCAL_EXPORT_PACKAGE_RESOURCES), \
      $(LOCAL_SOONG_RRO_DIRS))
endif
