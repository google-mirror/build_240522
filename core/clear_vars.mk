###########################################################
## Clear out values of all variables used by rule templates.
###########################################################

LOCAL_MODULE:=
LOCAL_MODULE_PATH:=
LOCAL_MODULE_RELATIVE_PATH :=
LOCAL_MODULE_STEM:=
LOCAL_DONT_CHECK_MODULE:=
LOCAL_CHECKED_MODULE:=
LOCAL_BUILT_MODULE:=
LOCAL_BUILT_MODULE_STEM:=
OVERRIDE_BUILT_MODULE_PATH:=
LOCAL_INSTALLED_MODULE:=
LOCAL_INSTALLED_MODULE_STEM:=
LOCAL_UNINSTALLABLE_MODULE:=
LOCAL_INTERMEDIATE_TARGETS:=
LOCAL_UNSTRIPPED_PATH:=
LOCAL_MODULE_CLASS:=
LOCAL_MODULE_SUFFIX:=
LOCAL_PACKAGE_NAME:=
LOCAL_OVERRIDES_PACKAGES:=
LOCAL_EXPORT_PACKAGE_RESOURCES:=
LOCAL_MANIFEST_PACKAGE_NAME:=
LOCAL_REQUIRED_MODULES:=
LOCAL_ACP_UNAVAILABLE:=
LOCAL_MODULE_TAGS:=
LOCAL_SRC_FILES:=
LOCAL_PREBUILT_OBJ_FILES:=
LOCAL_STATIC_JAVA_LIBRARIES:=
LOCAL_STATIC_LIBRARIES:=
# Group static libraries with "-Wl,--start-group" and "-Wl,--end-group" when linking.
LOCAL_GROUP_STATIC_LIBRARIES:=
LOCAL_WHOLE_STATIC_LIBRARIES:=
LOCAL_SHARED_LIBRARIES:=
LOCAL_IS_HOST_MODULE:=
LOCAL_CC:=
LOCAL_CXX:=
LOCAL_CPP_EXTENSION:=
LOCAL_NO_DEFAULT_COMPILER_FLAGS:=
LOCAL_FDO_SUPPORT:=
LOCAL_ARM_MODE:=
LOCAL_YACCFLAGS:=
LOCAL_ASFLAGS:=
LOCAL_CFLAGS:=
LOCAL_CPPFLAGS:=
LOCAL_CONLYFLAGS:=
LOCAL_RTTI_FLAG:=
LOCAL_C_INCLUDES:=
LOCAL_EXPORT_C_INCLUDE_DIRS:=
LOCAL_LDFLAGS:=
LOCAL_LDLIBS:=
LOCAL_AAPT_FLAGS:=
LOCAL_AAPT_INCLUDE_ALL_RESOURCES:=
LOCAL_SYSTEM_SHARED_LIBRARIES:=none
LOCAL_PREBUILT_LIBS:=
LOCAL_PREBUILT_EXECUTABLES:=
LOCAL_PREBUILT_JAVA_LIBRARIES:=
LOCAL_PREBUILT_STATIC_JAVA_LIBRARIES:=
LOCAL_PREBUILT_STRIP_COMMENTS:=
LOCAL_INTERMEDIATE_SOURCES:=
LOCAL_INTERMEDIATE_SOURCE_DIR:=
LOCAL_JAVACFLAGS:=
LOCAL_JAVA_LIBRARIES:=
LOCAL_JAVA_LAYERS_FILE:=
LOCAL_NO_STANDARD_LIBRARIES:=
LOCAL_CLASSPATH:=
LOCAL_DROIDDOC_USE_STANDARD_DOCLET:=
LOCAL_DROIDDOC_SOURCE_PATH:=
LOCAL_DROIDDOC_TEMPLATE_DIR:=
LOCAL_DROIDDOC_CUSTOM_TEMPLATE_DIR:=
LOCAL_DROIDDOC_ASSET_DIR:=
LOCAL_DROIDDOC_CUSTOM_ASSET_DIR:=
LOCAL_DROIDDOC_OPTIONS:=
LOCAL_DROIDDOC_HTML_DIR:=
LOCAL_ADDITIONAL_HTML_DIR:=
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
LOCAL_COMPRESS_MODULE_SYMBOLS:=
LOCAL_STRIP_MODULE:=
LOCAL_JNI_SHARED_LIBRARIES:=
LOCAL_JNI_SHARED_LIBRARIES_ABI:=
LOCAL_JAR_MANIFEST:=
LOCAL_INSTRUMENTATION_FOR:=
LOCAL_APK_LIBRARIES:=
LOCAL_MANIFEST_INSTRUMENTATION_FOR:=
LOCAL_AIDL_INCLUDES:=
LOCAL_JARJAR_RULES:=
LOCAL_ADDITIONAL_JAVA_DIR:=
LOCAL_ALLOW_UNDEFINED_SYMBOLS:=
LOCAL_DX_FLAGS:=
LOCAL_CERTIFICATE:=
LOCAL_SDK_VERSION:=
LOCAL_SDK_RES_VERSION:=
LOCAL_NDK_STL_VARIANT:=
LOCAL_EMMA_INSTRUMENT:=
LOCAL_PROGUARD_ENABLED:= # '',full,custom,nosystem,disabled,obfuscation,optimization
LOCAL_PROGUARD_FLAGS:=
LOCAL_PROGUARD_FLAG_FILES:=
LOCAL_EMMA_COVERAGE_FILTER:=
LOCAL_WARNINGS_ENABLE:=
LOCAL_FULL_MANIFEST_FILE:=
LOCAL_MANIFEST_FILE:=
LOCAL_RENDERSCRIPT_INCLUDES:=
LOCAL_RENDERSCRIPT_INCLUDES_OVERRIDE:=
LOCAL_RENDERSCRIPT_CC:=
LOCAL_RENDERSCRIPT_COMPATIBILITY:=
LOCAL_RENDERSCRIPT_FLAGS:=
LOCAL_RENDERSCRIPT_SKIP_INSTALL:=
LOCAL_RENDERSCRIPT_TARGET_API:=
LOCAL_BUILD_HOST_DEX:=
LOCAL_DEX_PREOPT:= # '',true,false,nostripping
LOCAL_DEX_PREOPT_IMAGE:=
LOCAL_PROTOC_OPTIMIZE_TYPE:= # lite(default),micro,nano,full
LOCAL_PROTOC_FLAGS:=
LOCAL_PROTO_JAVA_OUTPUT_PARAMS:=
LOCAL_NO_CRT:=
LOCAL_PROPRIETARY_MODULE:=
LOCAL_PRIVILEGED_MODULE:=
LOCAL_MODULE_OWNER:=
LOCAL_CTS_TEST_PACKAGE:=
LOCAL_CTS_TEST_RUNNER:=
LOCAL_CLANG:=
LOCAL_ADDRESS_SANITIZER:=
LOCAL_JAR_EXCLUDE_FILES:=
LOCAL_JAR_PACKAGES:=
LOCAL_LINT_FLAGS:=
LOCAL_SOURCE_FILES_ALL_GENERATED:= # '',true
# Don't delete the META_INF dir when merging static Java libraries.
LOCAL_DONT_DELETE_JAR_META_INF:=
LOCAL_ADDITIONAL_CERTIFICATES:=
LOCAL_PREBUILT_MODULE_FILE:=
LOCAL_POST_INSTALL_CMD:=
LOCAL_DIST_BUNDLED_BINARIES:=
LOCAL_HAL_STATIC_LIBRARIES:=
LOCAL_NO_SYNTAX_CHECK:=
LOCAL_NO_STATIC_ANALYZER:=
LOCAL_32_BIT_ONLY:= # '',true
LOCAL_NO_2ND_ARCH:= # '',true

# arch specific variables
LOCAL_SRC_FILES_$(TARGET_ARCH):=
LOCAL_CFLAGS_$(TARGET_ARCH):=
LOCAL_C_INCLUDES_$(TARGET_ARCH):=
LOCAL_ASFLAGS_$(TARGET_ARCH):=
LOCAL_NO_CRT_$(TARGET_ARCH):=
LOCAL_LDFLAGS_$(TARGET_ARCH):=
LOCAL_SHARED_LIBRARIES_$(TARGET_ARCH):=
LOCAL_STATIC_LIBRARIES_$(TARGET_ARCH):=
LOCAL_WHOLE_STATIC_LIBRARIES_$(TARGET_ARCH):=
LOCAL_GENERATED_SOURCES_$(TARGET_ARCH):=
ifdef TARGET_2ND_ARCH
LOCAL_SRC_FILES_$(TARGET_2ND_ARCH):=
LOCAL_CFLAGS_$(TARGET_2ND_ARCH):=
LOCAL_C_INCLUDES_$(TARGET_2ND_ARCH):=
LOCAL_ASFLAGS_$(TARGET_2ND_ARCH):=
LOCAL_NO_CRT_$(TARGET_2ND_ARCH):=
LOCAL_LDFLAGS_$(TARGET_2ND_ARCH):=
LOCAL_SHARED_LIBRARIES_$(TARGET_2ND_ARCH):=
LOCAL_STATIC_LIBRARIES_$(TARGET_2ND_ARCH):=
LOCAL_WHOLE_STATIC_LIBRARIES_$(TARGET_2ND_ARCH):=
LOCAL_GENERATED_SOURCES_$(TARGET_2ND_ARCH):=
endif

LOCAL_CFLAGS_32:=
LOCAL_CFLAGS_64:=
LOCAL_LDFLAGS_32:=
LOCAL_LDFLAGS_64:=
LOCAL_ASFLAGS_32:=
LOCAL_ASFLAGS_64:=
LOCAL_C_INCLUDES_32:=
LOCAL_C_INCLUDES_64:=

# Trim MAKEFILE_LIST so that $(call my-dir) doesn't need to
# iterate over thousands of entries every time.
# Leave the current makefile to make sure we don't break anything
# that expects to be able to find the name of the current makefile.
MAKEFILE_LIST := $(lastword $(MAKEFILE_LIST))
