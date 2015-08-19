# Convenient template for installing init.rc file as a prebuilt module.
# Input variables:
#   LOCAL_INIT_RC, the .rc file name without path, also used the module name.
#   LOCAL_INIT_RC_SRC, the source file path relative to LOCAL_PATH.
#   LOCAL_INIT_RC_INSTALL_PATH, the install path of the .rc file.
#

ifndef LOCAL_INIT_RC_SRC
LOCAL_INIT_RC_SRC := $(LOCAL_INIT_RC)
endif
ifndef LOCAL_INIT_RC_INSTALL_PATH
LOCAL_INIT_RC_INSTALL_PATH := $(TARGET_OUT)/init
endif

include $(CLEAR_VARS)
LOCAL_MODULE := $(LOCAL_INIT_RC)
LOCAL_SRC_FILES := $(LOCAL_INIT_RC_SRC)
LOCAL_MODULE_CLASS := ETC
LOCAL_MODULE_PATH := $(LOCAL_INIT_RC_INSTALL_PATH)
include $(BUILD_PREBUILT)

LOCAL_INIT_RC :=
LOCAL_INIT_RC_SRC :=
LOCAL_INIT_RC_INSTALL_PATH :=
