<<<<<<< HEAD   (4be654 Merge "Merge empty history for sparse-7121469-L4290000080720)
###########################################
## A thin wrapper around BUILD_EXECUTABLE
## Common flags for native benchmarks are added.
###########################################
$(call record-module-type,NATIVE_BENCHMARK)

LOCAL_STATIC_LIBRARIES += libgoogle-benchmark

LOCAL_MODULE_PATH_64 := $(TARGET_OUT_DATA_METRIC_TESTS)/$(LOCAL_MODULE)
LOCAL_MODULE_PATH_32 := $($(TARGET_2ND_ARCH_VAR_PREFIX)TARGET_OUT_DATA_METRIC_TESTS)/$(LOCAL_MODULE)

ifndef LOCAL_MULTILIB
ifndef LOCAL_32_BIT_ONLY
LOCAL_MULTILIB := both
endif
endif

include $(BUILD_EXECUTABLE)
=======
>>>>>>> BRANCH (fe6ad7 Merge "Version bump to RBT1.210107.001.A1 [core/build_id.mk])
