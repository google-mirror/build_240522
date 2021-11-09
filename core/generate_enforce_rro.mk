include $(CLEAR_VARS)

LOCAL_PACKAGE_NAME := $(enforce_rro_module)

intermediates := $(call intermediates-dir-for,APPS,$(LOCAL_PACKAGE_NAME),,COMMON)
rro_android_manifest_file := $(intermediates)/AndroidManifest.xml

ifeq (true,$(enforce_rro_source_is_manifest_package_name))
$(rro_android_manifest_file): PRIVATE_PACKAGE_NAME := $(enforce_rro_source_manifest_package_info)
$(rro_android_manifest_file): build/make/tools/generate-enforce-rro-android-manifest.py
	$(hide) build/make/tools/generate-enforce-rro-android-manifest.py -u -p $(PRIVATE_PACKAGE_NAME) -o $@
else
$(rro_android_manifest_file): PRIVATE_SOURCE_MANIFEST_FILE := $(enforce_rro_source_manifest_package_info)
$(rro_android_manifest_file): $(enforce_rro_source_manifest_package_info) build/make/tools/generate-enforce-rro-android-manifest.py
	$(hide) build/make/tools/generate-enforce-rro-android-manifest.py -p $(PRIVATE_SOURCE_MANIFEST_FILE) -o $@
endif

LOCAL_PATH:= $(intermediates)

# TODO(b/187404676): remove this condition when the prebuilt for packges exporting resource exists.
ifeq (,$(TARGET_BUILD_UNBUNDLED))
ifeq ($(enforce_rro_use_res_lib),true)
LOCAL_RES_LIBRARIES := $(enforce_rro_source_module)
endif
endif

LOCAL_FULL_MANIFEST_FILE := $(rro_android_manifest_file)
LOCAL_CERTIFICATE := platform

LOCAL_AAPT_FLAGS += --auto-add-overlay
LOCAL_RESOURCE_DIR := $(enforce_rro_source_overlays)

<<<<<<< HEAD   (3619c8 Merge "Merge empty history for sparse-7625297-L4670000095071)
=======
ifeq (product,$(enforce_rro_partition))
  LOCAL_PRODUCT_MODULE := true
else ifeq (vendor,$(enforce_rro_partition))
  LOCAL_VENDOR_MODULE := true
else
  $(error Unsupported partition. Want: [vendor/product] Got: [$(enforce_rro_partition)])
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

>>>>>>> BRANCH (77b382 Merge "Version bump to AAQ4.211109.001 [core/build_id.mk]" i)
include $(BUILD_RRO_PACKAGE)
