# Copyright 2005 The Android Open Source Project
#

LOCAL_PATH:= $(call my-dir)

include $(CLEAR_VARS)
LOCAL_CPP_EXTENSION := .cc
LOCAL_SRC_FILES := dextoc.cc
LOCAL_SHARED_LIBRARIES += libart liblog
LOCAL_C_INCLUDES := art/runtime
LOCAL_CFLAGS += -Wall
LOCAL_MODULE := dextoc
include $(BUILD_HOST_EXECUTABLE)
