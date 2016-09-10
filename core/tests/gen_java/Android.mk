# Can't use my-dir due to sanity check
LOCAL_PATH := build/core/tests/gen_java

include $(CLEAR_VARS)
LOCAL_MODULE := build_system_test_gen_java
ifeq ($(TARGET_PRODUCT),aosp_arm64)
LOCAL_SRC_FILES := test.proto
LOCAL_PROTOC_OPTIMIZE_TYPE := nano
endif
LOCAL_SRC_FILES += android/build/make/tests/Test.java
LOCAL_NO_STANDARD_LIBRARIES := true
LOCAL_JAVA_LIBRARIES := core-oj core-libart
LOCAL_DEX_PREOPT := false
include $(BUILD_STATIC_JAVA_LIBRARY)
