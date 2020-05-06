# App prebuilt coming from Soong.
# Extra inputs:

ifneq ($(LOCAL_MODULE_MAKEFILE),$(SOONG_ANDROID_MK))
  $(call pretty-error,soong_apk_set.mk may only be used from Soong)
endif

LOCAL_BUILT_MODULE_STEM := base-master.apk
LOCAL_INSTALLED_MODULE_STEM := base-master.apk
LOCAL_IS_RUNTIME_RESOURCE_OVERLAY := true

#######################################
include $(BUILD_SYSTEM)/base_rules.mk
#######################################

## Extract master APK from APK set into given directory
# $(1) APK set
# $(2) master APK entry (e.g., splits/base-master.apk
# $(3) target directory

define extract-master-from-apk-set
$(LOCAL_BUILT_MODULE): $(1)
	@echo "Extracting $$@"
	unzip -pq $$< $(2) >$$@

ALL_MODULES.$(my_register_name).INSTALLED += $(3)/$(notdir $(2))
ALL_MODULES.$(my_register_name).BUILT_INSTALLED += $(3)/$(notdir $(2))
endef

my_base_master := splits/base-master.apk
$(eval $(call extract-master-from-apk-set,$(LOCAL_PREBUILT_MODULE_FILE),$(my_base_master),$(call local-intermediates-dir)))
LOCAL_POST_INSTALL_CMD := unzip -qo -j -d $(dir $(LOCAL_INSTALLED_MODULE)) \
 $(LOCAL_PREBUILT_MODULE_FILE) -x $(my_base_master)
$(LOCAL_INSTALLED_MODULE): PRIVATE_POST_INSTALL_CMD := $(LOCAL_POST_INSTALL_CMD)
my_base_master :=

SOONG_ALREADY_CONV := $(SOONG_ALREADY_CONV) $(LOCAL_MODULE)
