ifneq ($(BOARD_VNDK_VERSION),)
LOCAL_PATH:= $(call my-dir)

#####################################################################
# Create the list of vndk libraries from the source code.
INTERNAL_VNDK_LIB_LIST := $(call intermediates-dir-for,PACKAGING,vndk)/libs.txt
$(INTERNAL_VNDK_LIB_LIST):
	@echo "Generate: $@"
	@mkdir -p $(dir $@)
	$(hide) echo -n > $@
	$(hide) $(foreach lib, $(LLNDK_LIBRARIES), \
	  echo LLNDK: $(lib).so >> $@;)
	$(hide) $(foreach lib, $(VNDK_SAMEPROCESS_LIBRARIES), \
	  echo VNDK-SP: $(lib).so >> $@;)
	$(hide) $(foreach lib, $(VNDK_CORE_LIBRARIES), \
	  echo VNDK-core: $(lib).so >> $@;)
	$(hide) $(foreach lib, $(VNDK_PRIVATE_LIBRARIES), \
	  echo VNDK-private: $(lib).so >> $@;)

#####################################################################
# This is the up-to-date list of vndk libs.
# TODO(b/62012285): the lib list should be stored somewhere under
# /prebuilts/vndk
ifeq (REL,$(PLATFORM_VERSION_CODENAME))
LATEST_VNDK_LIB_LIST := $(LOCAL_PATH)/$(PLATFORM_VNDK_VERSION).txt
else
LATEST_VNDK_LIB_LIST := $(LOCAL_PATH)/current.txt
endif

#####################################################################
# Check the generate list against the latest list stored in the
# source tree
.PHONY: check-vndk-list

# Check if vndk list is changed
droidcore: check-vndk-list

check-vndk-list-timestamp := $(call intermediates-dir-for,PACKAGING,vndk)/check-list-timestamp
check-vndk-list: $(check-vndk-list-timestamp)

_vndk_check_failure_message := "VNDK library list has changed."
ifeq (REL,$(PLATFORM_VERSION_CODENAME))
_vndk_check_failure_message += "This isn't allowed in API locked branches."
else
_vndk_check_failure_message += "Run 'make update-vndk-list' to update the list."
endif

$(check-vndk-list-timestamp): $(INTERNAL_VNDK_LIB_LIST) $(LATEST_VNDK_LIB_LIST)
	$(hide) ( diff --old-line-format="Removed %L" \
	  --new-line-format="Added %L" \
	  --unchanged-line-format="" \
	  $(LATEST_VNDK_LIB_LIST) $(INTERNAL_VNDK_LIB_LIST) \
	  || ( echo $(_vndk_check_failure_message); exit 1 ))
	$(hide) mkdir -p $(dir $@)
	$(hide) touch $@

# Update the latest VNDK lib list
.PHONY: update-vndk-list

update-vndk-list: $(INTERNAL_VNDK_LIB_LIST)
	$(if $(filter-out REL,$(PLATFORM_VERSION_CODENAME)), \
		$(hide) echo "Generate: $(LATEST_VNDK_LIB_LIST)"; \
		mkdir -p $(dir $(LATEST_VNDK_LIB_LIST)); \
		rm -f $(LATEST_VNDK_LIB_LIST); \
		cp $(INTERNAL_VNDK_LIB_LIST) $(LATEST_VNDK_LIB_LIST); \
		echo "$(LATEST_VNDK_LIB_LIST) updated.", \
		@echo "Updating VNDK library list is NOT allowed in API locked branches.")

include $(CLEAR_VARS)
LOCAL_MODULE := vndk_package
LOCAL_REQUIRED_MODULES := \
    $(addsuffix .vendor,$(VNDK_CORE_LIBRARIES)) \
    $(addsuffix .vendor,$(VNDK_SAMEPROCESS_LIBRARIES)) \
    $(LLNDK_LIBRARIES) \
    llndk.libraries.txt \
    vndksp.libraries.txt
include $(BUILD_PHONY_PACKAGE)

include $(CLEAR_VARS)
LOCAL_MODULE := vndk_snapshot_package
LOCAL_REQUIRED_MODULES := \
    $(foreach vndk_ver,$(PRODUCT_EXTRA_VNDK_VERSIONS),vndk_v$(vndk_ver)_$(TARGET_ARCH))
include $(BUILD_PHONY_PACKAGE)

endif # BOARD_VNDK_VERSION is set
