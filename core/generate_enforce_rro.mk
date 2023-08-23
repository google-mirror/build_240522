include $(CLEAR_VARS)

enforce_rro_module := $(enforce_rro_source_module)__$(PRODUCT_NAME)__auto_generated_rro_$(enforce_rro_partition)
LOCAL_PACKAGE_NAME := $(enforce_rro_module)

intermediates := $(call intermediates-dir-for,APPS,$(LOCAL_PACKAGE_NAME),,COMMON)
rro_android_manifest_file := $(intermediates)/AndroidManifest.xml

ifneq ($(PRODUCT_AUTO_GENERATE_RRO_PRIORITY),)
  enforce_rro_android_priority := $(PRODUCT_AUTO_GENERATE_RRO_PRIORITY)
else
  enforce_rro_android_priority := 0
endif

ifeq (true,$(enforce_rro_source_is_manifest_package_name))
  use_package_name_arg := --use-package-name
else
  use_package_name_arg :=
$(rro_android_manifest_file): $(enforce_rro_source_manifest_package_info)
endif

$(rro_android_manifest_file): PRIVATE_PACKAGE_INFO := $(enforce_rro_source_manifest_package_info)
$(rro_android_manifest_file): PRIVATE_USE_PACKAGE_NAME := $(use_package_name_arg)
$(rro_android_manifest_file): PRIVATE_PARTITION := $(enforce_rro_partition)
# Set the default priority to 0 but make it possible to override with PRODUCT_AUTO_GENERATE_RRO_PRIORITY
$(rro_android_manifest_file): PRIVATE_PRIORITY := $(enforce_rro_android_priority)
$(rro_android_manifest_file): build/make/tools/generate-enforce-rro-android-manifest.py
	$(hide) build/make/tools/generate-enforce-rro-android-manifest.py \
	    --package-info $(PRIVATE_PACKAGE_INFO) \
	    $(PRIVATE_USE_PACKAGE_NAME) \
	    --partition $(PRIVATE_PARTITION) \
	    --priority $(PRIVATE_PRIORITY) \
	    -o $@

LOCAL_PATH:= $(intermediates)

# TODO(b/187404676): remove this condition when the prebuilt for packges exporting resource exists.
ifeq (,$(TARGET_BUILD_UNBUNDLED))
ifeq ($(enforce_rro_use_res_lib),true)
  LOCAL_RES_LIBRARIES := $(enforce_rro_source_module)
endif
endif

LOCAL_FULL_MANIFEST_FILE := $(rro_android_manifest_file)

LOCAL_AAPT_FLAGS += --auto-add-overlay --keep-raw-values
LOCAL_RESOURCE_DIR := $(enforce_rro_source_overlays)

ifeq (odm,$(enforce_rro_partition))
  LOCAL_ODM_MODULE := true
else ifeq (oem,$(enforce_rro_partition))
  LOCAL_OEM_MODULE := true
else ifeq (product,$(enforce_rro_partition))
  LOCAL_PRODUCT_MODULE := true
else ifeq (systemext,$(enforce_rro_partition))
  LOCAL_SYSTEM_EXT_MODULE := true
else ifeq (vendor,$(enforce_rro_partition))
  LOCAL_VENDOR_MODULE := true
else
  $(error Unsupported partition. Want: [vendor/systemext/product/oem/odm] Got: [$(enforce_rro_partition)])
endif
ifneq (,$(TARGET_BUILD_UNBUNDLED))
  LOCAL_SDK_VERSION := current
else ifneq (,$(LOCAL_RES_LIBRARIES))
  # Technically we are linking against the app (if only to grab its resources),
  # and because it's potentially not building against the SDK, we can't either.
  LOCAL_PRIVATE_PLATFORM_APIS := true
else ifeq (framework-res,$(enforce_rro_source_module))
  LOCAL_PRIVATE_PLATFORM_APIS := true
else
  LOCAL_SDK_VERSION := current
endif

include $(BUILD_RRO_PACKAGE)
