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

LOCAL_MODULE := ziptime
LOCAL_MODULE_HOST_OS := darwin linux windows

include $(BUILD_HOST_EXECUTABLE)
