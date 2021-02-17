<<<<<<< HEAD   (4be654 Merge "Merge empty history for sparse-7121469-L4290000080720)
################################################
## A thin wrapper around BUILD_HOST_EXECUTABLE
## Common flags for host fuzz tests are added.
################################################
$(call record-module-type,HOST_FUZZ_TEST)

LOCAL_CFLAGS += -fsanitize-coverage=trace-pc-guard,indirect-calls,trace-cmp
LOCAL_STATIC_LIBRARIES += libLLVMFuzzer

include $(BUILD_HOST_EXECUTABLE)
=======
>>>>>>> BRANCH (fe6ad7 Merge "Version bump to RBT1.210107.001.A1 [core/build_id.mk])
