LOCAL_PATH:= $(call my-dir)

#######################################
# verity_key
include $(CLEAR_VARS)

LOCAL_MODULE := verity_key
LOCAL_LICENSE_KINDS := legacy_restricted
LOCAL_LICENSE_CONDITIONS := restricted
LOCAL_SRC_FILES := $(LOCAL_MODULE)
LOCAL_MODULE_CLASS := ETC
LOCAL_MODULE_PATH := $(TARGET_ROOT_OUT)

include $(BUILD_PREBUILT)

#######################################
<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
=======
# verity_key (installed to ramdisk)
#
# Enabling the target when using system-as-root would cause build failure, as TARGET_RAMDISK_OUT
# points to the same location as TARGET_ROOT_OUT.
ifneq ($(BOARD_BUILD_SYSTEM_ROOT_IMAGE),true)
  include $(CLEAR_VARS)
  LOCAL_MODULE := verity_key_ramdisk
  LOCAL_LICENSE_KINDS := legacy_restricted
  LOCAL_LICENSE_CONDITIONS := restricted
  LOCAL_MODULE_CLASS := ETC
  LOCAL_SRC_FILES := verity_key
  LOCAL_MODULE_STEM := verity_key
  LOCAL_MODULE_PATH := $(TARGET_RAMDISK_OUT)
  include $(BUILD_PREBUILT)
endif

#######################################
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
# adb key, if configured via PRODUCT_ADB_KEYS
ifdef PRODUCT_ADB_KEYS
  ifneq ($(filter eng userdebug,$(TARGET_BUILD_VARIANT)),)
    include $(CLEAR_VARS)
    LOCAL_MODULE := adb_keys
    LOCAL_LICENSE_KINDS := legacy_restricted
    LOCAL_LICENSE_CONDITIONS := restricted
    LOCAL_MODULE_CLASS := ETC
    LOCAL_MODULE_PATH := $(TARGET_ROOT_OUT)
    LOCAL_PREBUILT_MODULE_FILE := $(PRODUCT_ADB_KEYS)
    include $(BUILD_PREBUILT)
  endif
endif
<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
=======


#######################################
# otacerts: A keystore with the authorized keys in it, which is used to verify the authenticity of
# downloaded OTA packages.
include $(CLEAR_VARS)

LOCAL_MODULE := otacerts
LOCAL_LICENSE_KINDS := legacy_restricted
LOCAL_LICENSE_CONDITIONS := restricted
LOCAL_MODULE_CLASS := ETC
LOCAL_MODULE_STEM := otacerts.zip
LOCAL_MODULE_PATH := $(TARGET_OUT_ETC)/security
include $(BUILD_SYSTEM)/base_rules.mk
$(LOCAL_BUILT_MODULE): PRIVATE_CERT := $(DEFAULT_SYSTEM_DEV_CERTIFICATE).x509.pem
$(LOCAL_BUILT_MODULE): $(SOONG_ZIP) $(DEFAULT_SYSTEM_DEV_CERTIFICATE).x509.pem
	$(SOONG_ZIP) -o $@ -j -symlinks=false -f $(PRIVATE_CERT)


#######################################
# otacerts for recovery image.
include $(CLEAR_VARS)

LOCAL_MODULE := otacerts.recovery
LOCAL_LICENSE_KINDS := legacy_restricted
LOCAL_LICENSE_CONDITIONS := restricted
LOCAL_MODULE_CLASS := ETC
LOCAL_MODULE_STEM := otacerts.zip
LOCAL_MODULE_PATH := $(TARGET_RECOVERY_ROOT_OUT)/system/etc/security
include $(BUILD_SYSTEM)/base_rules.mk

extra_recovery_keys := $(patsubst %,%.x509.pem,$(PRODUCT_EXTRA_RECOVERY_KEYS))

$(LOCAL_BUILT_MODULE): PRIVATE_CERT := $(DEFAULT_SYSTEM_DEV_CERTIFICATE).x509.pem
$(LOCAL_BUILT_MODULE): PRIVATE_EXTRA_RECOVERY_KEYS := $(extra_recovery_keys)
$(LOCAL_BUILT_MODULE): \
	    $(SOONG_ZIP) \
	    $(DEFAULT_SYSTEM_DEV_CERTIFICATE).x509.pem \
	    $(extra_recovery_keys)
	$(SOONG_ZIP) -o $@ -j -symlinks=false \
	    $(foreach key_file, $(PRIVATE_CERT) $(PRIVATE_EXTRA_RECOVERY_KEYS), -f $(key_file))
>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
