# Copyright 2005 The Android Open Source Project
#
# Android.mk for retouch
#

ifneq ($(TARGET_SIMULATOR),true)

LOCAL_PATH:= $(call my-dir)

# First part: build the host executable, "retouch-prepare".
#
# On the host, we scan relocation lists produced by Apriori,
# and output file offset+value pairs, ready for retouching.

include $(CLEAR_VARS)

LOCAL_LDLIBS += -ldl
LOCAL_CFLAGS += -O2 -g
LOCAL_CFLAGS += -fno-function-sections -fno-data-sections -fno-inline
LOCAL_CFLAGS += -Wall -Wno-unused-function #-Werror
LOCAL_CFLAGS += -DDEBUG

ifeq ($(TARGET_ARCH),arm)
LOCAL_CFLAGS += -DARM_SPECIFIC_HACKS
LOCAL_CFLAGS += -DBIG_ENDIAN=1
endif

ifeq ($(HOST_OS),darwin)
LOCAL_CFLAGS += -DFSCANF_IS_BROKEN
endif
ifeq ($(HOST_OS),windows)
LOCAL_CFLAGS += -DFSCANF_IS_BROKEN
LOCAL_LDLIBS += -lintl
endif

LOCAL_SRC_FILES := \
	retouch-prepare.c

LOCAL_C_INCLUDES:= \
	$(LOCAL_PATH)/ \
	external/elfutils/lib/ \
	external/elfutils/libelf/ \
	external/elfutils/libebl/ \
	external/elfcopy/

LOCAL_STATIC_LIBRARIES := libelfcopy libelf libebl #dl

ifeq ($(TARGET_ARCH),arm)
LOCAL_STATIC_LIBRARIES += libebl_arm
endif

LOCAL_MODULE := retouch-prepare

include $(BUILD_HOST_EXECUTABLE)

# Second part: build the encoding test tool, "retouch-deltas".
# This can be used outside the build process to check the contents
# of encoded retouch files.

include $(CLEAR_VARS)

LOCAL_LDLIBS += -ldl
LOCAL_CFLAGS += -O2 -g
LOCAL_CFLAGS += -fno-function-sections -fno-data-sections -fno-inline
LOCAL_CFLAGS += -Wall -Wno-unused-function #-Werror
LOCAL_CFLAGS += -DDEBUG

ifeq ($(TARGET_ARCH),arm)
LOCAL_CFLAGS += -DARM_SPECIFIC_HACKS
LOCAL_CFLAGS += -DBIG_ENDIAN=1
endif

ifeq ($(HOST_OS),darwin)
LOCAL_CFLAGS += -DFSCANF_IS_BROKEN
endif
ifeq ($(HOST_OS),windows)
LOCAL_CFLAGS += -DFSCANF_IS_BROKEN
LOCAL_LDLIBS += -lintl
endif

LOCAL_SRC_FILES := \
	retouch-deltas.c

LOCAL_C_INCLUDES:= \
	$(LOCAL_PATH)/

LOCAL_STATIC_LIBRARIES := libelfcopy libelf libebl #dl

ifeq ($(TARGET_ARCH),arm)
LOCAL_STATIC_LIBRARIES += libebl_arm
endif

LOCAL_MODULE := retouch-deltas

include $(BUILD_HOST_EXECUTABLE)

endif # TARGET_SIMULATOR != true
