# Java prebuilt coming from Soong.
# Extra inputs:
# LOCAL_SOONG_HEADER_JAR
# LOCAL_SOONG_DEX_JAR

ifneq ($(LOCAL_MODULE_MAKEFILE),$(SOONG_ANDROID_MK))
  $(call pretty-error,soong_java_prebuilt.mk may only be used from Soong)
endif

LOCAL_MODULE_SUFFIX := .jar
LOCAL_BUILT_MODULE_STEM := javalib.jar

#######################################
include $(BUILD_SYSTEM)/base_rules.mk
#######################################

full_classes_jar := $(intermediates.COMMON)/classes.jar
full_classes_header_jar := $(intermediates.COMMON)/classes-header.jar
common_javalib.jar := $(intermediates.COMMON)/javalib.jar

LOCAL_FULL_CLASSES_PRE_JACOCO_JAR := $(LOCAL_PREBUILT_MODULE_FILE)

#######################################
include $(BUILD_SYSTEM)/jacoco.mk
#######################################

$(eval $(call copy-one-file,$(LOCAL_FULL_CLASSES_JACOCO_JAR),$(full_classes_jar)))

ifdef LOCAL_SOONG_HEADER_JAR
$(eval $(call copy-one-file,$(LOCAL_SOONG_HEADER_JAR),$(full_classes_header_jar)))
else
$(eval $(call copy-one-file,$(full_classes_jar),$(full_classes_header_jar)))
endif

ifdef LOCAL_SOONG_DEX_JAR
$(eval $(call copy-one-file,$(LOCAL_SOONG_DEX_JAR),$(common_javalib.jar)))
$(eval $(call copy-one-file,$(LOCAL_SOONG_DEX_JAR),$(LOCAL_BUILT_MODULE)))
java-dex : $(LOCAL_BUILT_MODULE)
else
$(eval $(call copy-one-file,$(full_classes_jar),$(LOCAL_BUILT_MODULE)))
endif

javac-check : $(full_classes_jar)
javac-check-$(LOCAL_MODULE) : $(full_classes_jar)

# Built in equivalent to include $(CLEAR_VARS)
LOCAL_SOONG_HEADER_JAR :=
LOCAL_SOONG_DEX_JAR :=
