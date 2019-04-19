LOCAL_PATH:= $(call my-dir)

#######################################
# verity_key
include $(CLEAR_VARS)

LOCAL_MODULE := verity_key
LOCAL_SRC_FILES := $(LOCAL_MODULE)
LOCAL_MODULE_CLASS := ETC
LOCAL_MODULE_PATH := $(TARGET_ROOT_OUT)

include $(BUILD_PREBUILT)

#######################################
# adb key, if configured via PRODUCT_ADB_KEYS
ifdef PRODUCT_ADB_KEYS
  ifneq ($(filter eng userdebug,$(TARGET_BUILD_VARIANT)),)
    include $(CLEAR_VARS)
    LOCAL_MODULE := adb_keys
    LOCAL_MODULE_CLASS := ETC
    LOCAL_MODULE_PATH := $(TARGET_ROOT_OUT)
    LOCAL_PREBUILT_MODULE_FILE := $(PRODUCT_ADB_KEYS)
    include $(BUILD_PREBUILT)
  endif
endif


#######################################
# otacerts: A keystore with the authorized keys in it, used to verify the authenticity of
# downloaded OTA packages.
include $(CLEAR_VARS)

LOCAL_MODULE := otacerts
LOCAL_MODULE_CLASS := ETC
LOCAL_MODULE_STEM := otacerts.zip
LOCAL_MODULE_PATH := $(TARGET_OUT_ETC)/security
include $(BUILD_SYSTEM)/base_rules.mk
$(LOCAL_BUILT_MODULE): PRIVATE_CERT := $(DEFAULT_SYSTEM_DEV_CERTIFICATE).x509.pem
$(LOCAL_BUILT_MODULE): $(SOONG_ZIP) $(DEFAULT_SYSTEM_DEV_CERTIFICATE).x509.pem
	$(SOONG_ZIP) -o $@ -C $(dir $(PRIVATE_CERT)) -f $(PRIVATE_CERT)


#######################################
# otacerts for recovery image.
include $(CLEAR_VARS)

LOCAL_MODULE := otacerts.recovery
LOCAL_MODULE_CLASS := ETC
LOCAL_MODULE_STEM := otacerts.zip
LOCAL_MODULE_PATH := $(TARGET_RECOVERY_ROOT_OUT)/system/etc/security
include $(BUILD_SYSTEM)/base_rules.mk

extra_keys := $(patsubst %,%.x509.pem,$(PRODUCT_EXTRA_RECOVERY_KEYS))

$(LOCAL_BUILT_MODULE): PRIVATE_CERT := $(DEFAULT_SYSTEM_DEV_CERTIFICATE).x509.pem
$(LOCAL_BUILT_MODULE): PRIVATE_EXTRA_RECOVERY_KEYS := $(extra_keys)
$(LOCAL_BUILT_MODULE): $(SOONG_ZIP) $(DEFAULT_SYSTEM_DEV_CERTIFICATE).x509.pem $(extra_keys)
	$(SOONG_ZIP) -o $@ \
	    $(foreach key_file, $(PRIVATE_CERT) $(extra_keys), -C $(dir $(key_file)) -f $(key_file))


#######################################
# update_engine_payload_key, used by update_engine. We use the same key as otacerts but in RSA
# public key format.
include $(CLEAR_VARS)

LOCAL_MODULE := update_engine_payload_key
LOCAL_MODULE_CLASS := ETC
LOCAL_MODULE_STEM := update-payload-key.pub.pem
LOCAL_MODULE_PATH := $(TARGET_OUT_ETC)/update_engine
include $(BUILD_SYSTEM)/base_rules.mk
$(LOCAL_BUILT_MODULE): $(DEFAULT_SYSTEM_DEV_CERTIFICATE).x509.pem
	rm -f $@
	mkdir -p $(dir $@)
	openssl x509 -pubkey -noout -in $< > $@

#######################################
# update_engine_payload_key for recovery image, used by update_engine.
include $(CLEAR_VARS)

LOCAL_MODULE := update_engine_payload_key.recovery
LOCAL_MODULE_CLASS := ETC
LOCAL_MODULE_STEM := update-payload-key.pub.pem
LOCAL_MODULE_PATH := $(TARGET_RECOVERY_ROOT_OUT)/system/etc/update_engine
include $(BUILD_SYSTEM)/base_rules.mk
$(LOCAL_BUILT_MODULE): $(DEFAULT_SYSTEM_DEV_CERTIFICATE).x509.pem
	rm -f $@
	mkdir -p $(dir $@)
	openssl x509 -pubkey -noout -in $< > $@
