###############################################################
## Renderscript support for java
## Adds rules to convert .rs files to .java and .bc files
###############################################################

renderscript_sources := $(filter %.rs,$(LOCAL_SRC_FILES))
LOCAL_SRC_FILES := $(filter-out %.rs,$(LOCAL_SRC_FILES))

rs_generated_res_zip :=
rs_generated_src_jar :=
rs_compatibility_jni_libs :=
ifneq ($(renderscript_sources),)
renderscript_sources_fullpath := $(addprefix $(LOCAL_PATH)/, $(renderscript_sources))
renderscript_intermediate.COMMON := $(intermediates.COMMON)/renderscript
rs_generated_res_zip := $(renderscript_intermediate.COMMON)/res.zip
rs_generated_src_jar := $(renderscript_intermediate.COMMON)/rs.srcjar

LOCAL_SRCJARS += $(rs_generated_src_jar)

# Defaulting to an empty string uses the latest available platform SDK.
renderscript_target_api :=

ifneq (,$(LOCAL_RENDERSCRIPT_TARGET_API))
  renderscript_target_api := $(LOCAL_RENDERSCRIPT_TARGET_API)
else
  ifneq (,$(LOCAL_SDK_VERSION))
    # Set target-api for LOCAL_SDK_VERSIONs other than current.
    ifneq (,$(filter-out current system_current test_current core_current, $(LOCAL_SDK_VERSION)))
      renderscript_target_api := $(call get-numeric-sdk-version,$(LOCAL_SDK_VERSION))
    endif
  endif  # LOCAL_SDK_VERSION is set
endif  # LOCAL_RENDERSCRIPT_TARGET_API is set

# For 64-bit, we always have to upgrade to at least 21 for compat build.
ifneq ($(LOCAL_RENDERSCRIPT_COMPATIBILITY),)
  ifeq ($(TARGET_IS_64_BIT),true)
    ifneq ($(filter $(RSCOMPAT_32BIT_ONLY_API_LEVELS),$(renderscript_target_api)),)
      renderscript_target_api := 21
    endif
  endif
endif

ifeq ($(LOCAL_RENDERSCRIPT_CC),)
LOCAL_RENDERSCRIPT_CC := $(LLVM_RS_CC)
endif

# Turn on all warnings and warnings as errors for RS compiles.
# This can be disabled with LOCAL_RENDERSCRIPT_FLAGS := -Wno-error
renderscript_flags := -Wall -Werror
renderscript_flags += $(LOCAL_RENDERSCRIPT_FLAGS)

# prepend the RenderScript system include path
ifneq ($(filter-out current system_current test_current core_current,$(LOCAL_SDK_VERSION))$(if $(TARGET_BUILD_APPS),$(filter current system_current test_current,$(LOCAL_SDK_VERSION))),)
# if a numeric LOCAL_SDK_VERSION, or current LOCAL_SDK_VERSION with TARGET_BUILD_APPS
LOCAL_RENDERSCRIPT_INCLUDES := \
    $(HISTORICAL_SDK_VERSIONS_ROOT)/renderscript/clang-include \
    $(HISTORICAL_SDK_VERSIONS_ROOT)/renderscript/include \
    $(LOCAL_RENDERSCRIPT_INCLUDES)
else
LOCAL_RENDERSCRIPT_INCLUDES := \
    $(TOPDIR)external/clang/lib/Headers \
    $(TOPDIR)frameworks/rs/script_api/include \
    $(LOCAL_RENDERSCRIPT_INCLUDES)
endif

ifneq ($(LOCAL_RENDERSCRIPT_INCLUDES_OVERRIDE),)
LOCAL_RENDERSCRIPT_INCLUDES := $(LOCAL_RENDERSCRIPT_INCLUDES_OVERRIDE)
endif

bc_files := $(patsubst %.rs,%.bc, $(notdir $(renderscript_sources)))
bc_dep_files := $(addprefix $(renderscript_intermediate.COMMON)/,$(patsubst %.bc,%.d,$(bc_files)))

$(rs_generated_src_jar): PRIVATE_RS_INCLUDES := $(LOCAL_RENDERSCRIPT_INCLUDES)
$(rs_generated_src_jar): PRIVATE_RS_CC := $(LOCAL_RENDERSCRIPT_CC)
$(rs_generated_src_jar): PRIVATE_RS_FLAGS := $(renderscript_flags)
$(rs_generated_src_jar): PRIVATE_RS_SOURCE_FILES := $(renderscript_sources_fullpath)
$(rs_generated_src_jar): PRIVATE_RS_OUTPUT_DIR := $(renderscript_intermediate.COMMON)
$(rs_generated_src_jar): PRIVATE_RS_TARGET_API := $(renderscript_target_api)
$(rs_generated_src_jar): PRIVATE_DEP_FILES := $(bc_dep_files)
$(rs_generated_src_jar): PRIVATE_RS_OUTPUT_RES_ZIP := $(rs_generated_res_zip)
$(rs_generated_src_jar): .KATI_IMPLICIT_OUTPUTS := $(rs_generated_res_zip)
$(rs_generated_src_jar): $(renderscript_sources_fullpath) $(LOCAL_RENDERSCRIPT_CC) $(SOONG_ZIP)
	$(transform-renderscripts-to-java-and-bc)

# include the dependency files (.d/.P) generated by llvm-rs-cc.
$(call include-depfile,$(rs_generated_src_jar).P,$(rs_generated_src_jar))

ifneq ($(LOCAL_RENDERSCRIPT_COMPATIBILITY),)


ifeq ($(filter $(RSCOMPAT_32BIT_ONLY_API_LEVELS),$(renderscript_target_api)),)
ifeq ($(TARGET_IS_64_BIT),true)
renderscript_intermediate.bc_folder := $(renderscript_intermediate.COMMON)/res/raw/bc64/
else
renderscript_intermediate.bc_folder := $(renderscript_intermediate.COMMON)/res/raw/bc32/
endif
else
renderscript_intermediate.bc_folder := $(renderscript_intermediate.COMMON)/res/raw/
endif

rs_generated_bc := $(addprefix \
    $(renderscript_intermediate.bc_folder), $(bc_files))

renderscript_intermediate := $(intermediates)/renderscript

# We don't need the .so files in bundled branches
# Prevent these from showing up on the device
# One exception is librsjni.so, which is needed for
# both native path and compat path.
rs_jni_lib := $(TARGET_OUT_INTERMEDIATE_LIBRARIES)/librsjni.so
LOCAL_JNI_SHARED_LIBRARIES += librsjni

ifneq (,$(TARGET_BUILD_APPS)$(FORCE_BUILD_RS_COMPAT))

rs_compatibility_jni_libs := $(addprefix \
    $(renderscript_intermediate)/librs., \
    $(patsubst %.bc,%.so, $(bc_files)))

$(rs_generated_src_jar): .KATI_IMPLICIT_OUTPUTS += $(rs_generated_bc)

rs_support_lib := $(TARGET_OUT_INTERMEDIATE_LIBRARIES)/libRSSupport.so
LOCAL_JNI_SHARED_LIBRARIES += libRSSupport

rs_support_io_lib :=
# check if the target api level support USAGE_IO
ifeq ($(filter $(RSCOMPAT_NO_USAGEIO_API_LEVELS),$(renderscript_target_api)),)
rs_support_io_lib := $(TARGET_OUT_INTERMEDIATE_LIBRARIES)/libRSSupportIO.so
LOCAL_JNI_SHARED_LIBRARIES += libRSSupportIO
endif

my_arch := $(TARGET_$(LOCAL_2ND_ARCH_VAR_PREFIX)ARCH)
ifneq (,$(filter arm64 mips64 x86_64,$(my_arch)))
  my_min_sdk_version := 21
else
  my_min_sdk_version := $(MIN_SUPPORTED_SDK_VERSION)
endif

$(rs_compatibility_jni_libs): $(RS_PREBUILT_CLCORE) \
    $(rs_support_lib) $(rs_support_io_lib) $(rs_jni_lib) $(rs_compiler_rt)
$(rs_compatibility_jni_libs): $(BCC_COMPAT)
$(rs_compatibility_jni_libs): PRIVATE_CXX := $(CXX_WRAPPER) $(TARGET_CXX)
$(rs_compatibility_jni_libs): PRIVATE_SDK_VERSION := $(my_min_sdk_version)
$(rs_compatibility_jni_libs): $(renderscript_intermediate)/librs.%.so: \
    $(renderscript_intermediate.bc_folder)%.bc \
    $(SOONG_OUT_DIR)/ndk.timestamp
	$(transform-bc-to-so)

endif

endif

LOCAL_INTERMEDIATE_TARGETS += $(rs_generated_src_jar)
# Make sure the generated resource will be added to the apk.
LOCAL_RESOURCE_DIR := $(renderscript_intermediate.COMMON)/res $(LOCAL_RESOURCE_DIR)
endif
