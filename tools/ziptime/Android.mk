#
# Copyright 2008 The Android Open Source Project
#
# Zip alignment tool
#

LOCAL_PATH:= $(call my-dir)
include $(CLEAR_VARS)

LOCAL_SRC_FILES := \
	ZipTime.cpp \
	ZipEntry.cpp \
	ZipFile.cpp

LOCAL_STATIC_LIBRARIES := \
	libutils \
	liblog \

LOCAL_LDLIBS_linux += -lrt

ifneq ($(strip $(BUILD_HOST_static)),)
LOCAL_LDLIBS += -lpthread
endif # BUILD_HOST_static

LOCAL_MODULE := ziptime
LOCAL_MODULE_HOST_OS := darwin linux windows

include $(BUILD_HOST_EXECUTABLE)
