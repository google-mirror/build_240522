# App prebuilt coming from Soong.
# Extra inputs:
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

$(eval $(call copy-one-file,$(LOCAL_SOONG_CLASSES_JAR),$(full_classes_jar)))
$(eval $(call copy-one-file,$(LOCAL_SOONG_CLASSES_JAR),$(full_classes_pre_proguard_jar)))

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

ifneq ($(TURBINE_ENABLED),false)
ifdef LOCAL_SOONG_HEADER_JAR
$(eval $(call copy-one-file,$(LOCAL_SOONG_HEADER_JAR),$(full_classes_header_jar)))
else
$(eval $(call copy-one-file,$(full_classes_jar),$(full_classes_header_jar)))
endif
endif # TURBINE_ENABLED != false


$(eval $(call copy-one-file,$(LOCAL_PREBUILT_MODULE_FILE),$(LOCAL_BUILT_MODULE)))

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

<<<<<<< HEAD   (240c89 Merge "Merge empty history for sparse-9140227-L7340000095670)
=======
# install symbol files of JNI libraries
my_jni_lib_symbols_copy_files := $(foreach f,$(LOCAL_SOONG_JNI_LIBS_SYMBOLS),\
  $(call word-colon,1,$(f)):$(patsubst $(PRODUCT_OUT)/%,$(TARGET_OUT_UNSTRIPPED)/%,$(call word-colon,2,$(f))))
$(LOCAL_BUILT_MODULE): | $(call copy-many-files, $(my_jni_lib_symbols_copy_files))

# embedded JNI will already have been handled by soong
my_embed_jni :=
my_prebuilt_jni_libs :=
ifdef LOCAL_SOONG_JNI_LIBS_$(TARGET_ARCH)
  my_2nd_arch_prefix :=
  LOCAL_JNI_SHARED_LIBRARIES := $(LOCAL_SOONG_JNI_LIBS_$(TARGET_ARCH))
  partition_lib_pairs :=  $(LOCAL_SOONG_JNI_LIBS_PARTITION_$(TARGET_ARCH))
  include $(BUILD_SYSTEM)/install_jni_libs_internal.mk
endif
ifdef TARGET_2ND_ARCH
  ifdef LOCAL_SOONG_JNI_LIBS_$(TARGET_2ND_ARCH)
    my_2nd_arch_prefix := $(TARGET_2ND_ARCH_VAR_PREFIX)
    LOCAL_JNI_SHARED_LIBRARIES := $(LOCAL_SOONG_JNI_LIBS_$(TARGET_2ND_ARCH))
    partition_lib_pairs :=  $(LOCAL_SOONG_JNI_LIBS_PARTITION_$(TARGET_2ND_ARCH))
    include $(BUILD_SYSTEM)/install_jni_libs_internal.mk
  endif
endif
LOCAL_SHARED_JNI_LIBRARIES :=
my_embed_jni :=
my_prebuilt_jni_libs :=
my_2nd_arch_prefix :=
partition_lib_pairs :=

>>>>>>> BRANCH (0a00ba Merge "Version bump to TKB1.221005.001.A1 [core/build_id.mk])
PACKAGES := $(PACKAGES) $(LOCAL_MODULE)
ifdef LOCAL_CERTIFICATE
  PACKAGES.$(LOCAL_MODULE).CERTIFICATE := $(LOCAL_CERTIFICATE)
  PACKAGES.$(LOCAL_MODULE).PRIVATE_KEY := $(patsubst %.x509.pem,%.pk8,$(LOCAL_CERTIFICATE))
endif

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
