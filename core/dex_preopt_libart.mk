####################################
# ART boot image installation
# Input variables:
#   my_boot_image_name: the boot image to install
#   my_boot_image_arch: the architecture to install (e.g. TARGET_ARCH, not expanded)
#   my_boot_image_out:  the install directory (e.g. $(PRODUCT_OUT))
#   my_boot_image_syms: the symbols director (e.g. $(TARGET_OUT_UNSTRIPPED))
#   my_boot_image_root: make variable used to store installed image module name
#   my_boot_image_vdex_extra_symlink_arch: list of extra architectures to install vdex symlinks for
#
# Install the boot images compiled by Soong.
# Create a phony package named dexpreopt_bootjar.$(my_boot_image_name)_$($(my_boot_image_arch))
# that installs all of bootjars' dexpreopt images.
# The name of the phony module is saved in $(my_boot_image_root).
#
####################################

LOCAL_PATH := $(BUILD_SYSTEM)

my_suffix := $(my_boot_image_name)_$($(my_boot_image_arch))
is_host := $(strip $(filter HOST%,$(my_boot_image_arch)))
is_primary_arch := $(strip $(filter-out %_2ND_ARCH,$(my_boot_image_arch)))

# Define a module for bootjar dexpreopt file.
# $(1): module name
# $(2): prebuilt file path
# $(3): destination path
define dexpreopt-bootjar-module
include $(CLEAR_VARS)
LOCAL_MODULE := $(1)
LOCAL_PREBUILT_MODULE_FILE := $(2)
LOCAL_MODULE_PATH := $(dir $(3))
LOCAL_MODULE_STEM := $(notdir $(3))
ifneq (,$(is_host))
  LOCAL_IS_HOST_MODULE := true
endif
LOCAL_MODULE_CLASS := ETC
include $(BUILD_PREBUILT)
endef

my_bootjar_modules := \
  $(foreach p,$(DEXPREOPT_IMAGE_BUILT_INSTALLED_$(my_suffix)), \
    $(eval src := $(call word-colon,1,$(p))) \
    $(eval dest := $(my_boot_image_out)$(call word-colon,2,$(p))) \
    $(eval module_name := $(notdir $(dest)).dexpreopt_bootjar.$(my_suffix)) \
    $(eval $(call dexpreopt-bootjar-module,$(module_name),$(src),$(dest))) \
    $(module_name) \
  )

my_bootjar_modules += \
  $(foreach p,$(DEXPREOPT_IMAGE_VDEX_BUILT_INSTALLED_$(my_suffix)), \
    $(eval src := $(call word-colon,1,$(p))) \
    $(eval dest := $(my_boot_image_out)$(call word-colon,2,$(p))) \
    $(eval shared_vdex_dest := $(dir $(patsubst %/,%,$(dir $(dest))))$(notdir $(dest))) \
    $(eval shared_vdex_module_name := $(notdir $(dest)).dexpreopt_bootjar.$(my_boot_image_name)) \
    $(if $(is_primary_arch), \
      $(eval $(call dexpreopt-bootjar-module,$(shared_vdex_module_name),$(src),$(shared_vdex_dest),)) \
      $(eval ALL_MODULES.$(my_register_name).INSTALLED += $(dest)) \
      $(eval $(my_all_targets) : $(dest)) \
      $(call symlink-file,$(LOCAL_INSTALLED_MODULE),../$(notdir $(dest)),$(dest)) \
      $(foreach arch,$(my_boot_image_vdex_extra_symlink_arch), \
        $(eval dest := $(dir $(patsubst %/,%,$(dir $(dest))))$(arch)/$(notdir $(dest))) \
        $(eval ALL_MODULES.$(my_register_name).INSTALLED += $(dest)) \
        $(eval $(my_all_targets) : $(dest)) \
        $(call symlink-file,$(LOCAL_INSTALLED_MODULE),../$(notdir $(dest)),$(dest)) \
      )\
    ) \
    $(shared_vdex_module_name) \
  )

my_bootjar_unstripped := $(call copy-many-files,$(DEXPREOPT_IMAGE_UNSTRIPPED_BUILT_INSTALLED_$(my_suffix)),$(my_boot_image_syms))

include $(CLEAR_VARS)
LOCAL_MODULE := dexpreopt_bootjar.$(my_suffix)
LOCAL_REQUIRED_MODULES := $(my_bootjar_modules)
LOCAL_ADDITIONAL_DEPENDENCIES := $(my_bootjar_unstripped)
ifneq (,$(is_host))
  LOCAL_IS_HOST_MODULE := true
endif
include $(BUILD_PHONY_PACKAGE)

$(my_boot_image_root) += $(LOCAL_MODULE)

my_suffix :=
my_bootjar_modules :=
my_bootjar_unstripped :=
