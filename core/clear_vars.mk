###########################################################
## Clear out values of all variables used by rule templates.
###########################################################

LOCAL_MODULE:=
LOCAL_MODULE_PATH:=
LOCAL_MODULE_STEM:=
LOCAL_BUILT_MODULE:=
LOCAL_BUILT_MODULE_STEM:=
OVERRIDE_BUILT_MODULE_PATH:=
LOCAL_INSTALLED_MODULE:=
LOCAL_UNINSTALLABLE_MODULE:=
LOCAL_INTERMEDIATE_TARGETS:=
LOCAL_UNSTRIPPED_PATH:=
LOCAL_MODULE_CLASS:=
LOCAL_MODULE_SUFFIX:=
LOCAL_PACKAGE_NAME:=
LOCAL_OVERRIDES_PACKAGES:=
LOCAL_EXPORT_PACKAGE_RESOURCES:=
LOCAL_REQUIRED_MODULES:=
LOCAL_ACP_UNAVAILABLE:=
LOCAL_MODULE_TAGS:=
LOCAL_SRC_FILES:=
LOCAL_PREBUILT_OBJ_FILES:=
LOCAL_STATIC_JAVA_LIBRARIES:=
LOCAL_STATIC_LIBRARIES:=
LOCAL_WHOLE_STATIC_LIBRARIES:=
LOCAL_SHARED_LIBRARIES:=
LOCAL_IS_HOST_MODULE:=
LOCAL_CC:=
LOCAL_CXX:=
LOCAL_CPP_EXTENSION:=
LOCAL_NO_DEFAULT_COMPILER_FLAGS:=
LOCAL_ARM_MODE:=
LOCAL_YACCFLAGS:=
LOCAL_ASFLAGS:=
LOCAL_CFLAGS:=
LOCAL_CPPFLAGS:=
LOCAL_C_INCLUDES:=
LOCAL_LDFLAGS:=
LOCAL_LDLIBS:=
LOCAL_AAPT_FLAGS:=
LOCAL_SYSTEM_SHARED_LIBRARIES:=none
LOCAL_PREBUILT_LIBS:=
LOCAL_PREBUILT_EXECUTABLES:=
LOCAL_PREBUILT_JAVA_LIBRARIES:=
LOCAL_PREBUILT_STATIC_JAVA_LIBRARIES:=
LOCAL_INTERMEDIATE_SOURCES:=
LOCAL_JAVA_LIBRARIES:=
LOCAL_NO_STANDARD_LIBRARIES:=
LOCAL_CLASSPATH:=
LOCAL_DROIDDOC_SOURCE_PATH:=
LOCAL_DROIDDOC_TEMPLATE_DIR:=
LOCAL_DROIDDOC_CUSTOM_TEMPLATE_DIR:=
LOCAL_DROIDDOC_ASSET_DIR:=
LOCAL_DROIDDOC_CUSTOM_ASSET_DIR:=
LOCAL_DROIDDOC_OPTIONS:=
LOCAL_DROIDDOC_HTML_DIR:=
LOCAL_ASSET_FILES:=
LOCAL_ASSET_DIR:=
LOCAL_RESOURCE_DIR:=
LOCAL_JAVA_RESOURCE_DIRS:=
LOCAL_JAVA_RESOURCE_FILES:=
LOCAL_GENERATED_SOURCES:=
LOCAL_COPY_HEADERS_TO:=
LOCAL_COPY_HEADERS:=
LOCAL_FORCE_STATIC_EXECUTABLE:=
LOCAL_ADDITIONAL_DEPENDENCIES:=
LOCAL_PRELINK_MODULE:=
LOCAL_COMPRESS_MODULE_SYMBOLS:=
LOCAL_STRIP_MODULE:=
LOCAL_POST_PROCESS_COMMAND:=true
LOCAL_JNI_SHARED_LIBRARIES:=
LOCAL_JAR_MANIFEST:=
LOCAL_INSTRUMENTATION_FOR:=
LOCAL_INSTRUMENTATION_FOR_PACKAGE_NAME:=
LOCAL_AIDL_INCLUDES:=
LOCAL_JARJAR_RULES:=
LOCAL_ADDITIONAL_JAVA_DIR:=
LOCAL_ALLOW_UNDEFINED_SYMBOLS:=
LOCAL_DX_FLAGS:=
LOCAL_CERTIFICATE:=
LOCAL_SDK_VERSION:=
LOCAL_NO_EMMA_INSTRUMENT:=
LOCAL_NO_EMMA_COMPILE:=

# Trim MAKEFILE_LIST so that $(call my-dir) doesn't need to
# iterate over thousands of entries every time.
# Leave the current makefile to make sure we don't break anything
# that expects to be able to find the name of the current makefile.
MAKEFILE_LIST := $(lastword $(MAKEFILE_LIST))
