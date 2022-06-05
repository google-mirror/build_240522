$(call record-module-type,HOST_EXECUTABLE)
LOCAL_IS_HOST_MODULE := true
my_prefix := owner
LOCAL_HOST_PREFIX := owner
include $(BUILD_SYSTEM)/multilib.mk
