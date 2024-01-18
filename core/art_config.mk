# ART configuration that has to be determined after product config is resolved.
#
# Inputs:
# PRODUCT_ENABLE_UFFD_GC: See comments in build/make/core/product.mk.
# OVERRIDE_ENABLE_UFFD_GC: Overrides PRODUCT_ENABLE_UFFD_GC. Can be passed from the commandline for
# debugging purposes.
# BOARD_API_LEVEL: See comments in build/make/core/main.mk.
# BOARD_SHIPPING_API_LEVEL: See comments in build/make/core/main.mk.
# PRODUCT_SHIPPING_API_LEVEL: See comments in build/make/core/product.mk.
#
# Outputs:
# ENABLE_UFFD_GC: Whether to use userfaultfd GC.

config_enable_uffd_gc := \
  $(firstword $(OVERRIDE_ENABLE_UFFD_GC) $(PRODUCT_ENABLE_UFFD_GC) default)

ifeq (,$(filter default true false,$(config_enable_uffd_gc)))
  $(error Unknown PRODUCT_ENABLE_UFFD_GC value: $(config_enable_uffd_gc))
endif

ENABLE_UFFD_GC := $(config_enable_uffd_gc)
# If the value is "default", it will be mangled by post_process_props.py.
ADDITIONAL_PRODUCT_PROPERTIES += ro.dalvik.vm.enable_uffd_gc=$(config_enable_uffd_gc)
