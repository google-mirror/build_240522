ARCH_ARM_HAVE_ARMV8A    := true

SUPPORTED_ARMV8_CPU_VARIANTS := \
  generic                     \
  cortex-a53                  \
  cortex-a57                  \
  cortex-a72                  \
  denver64                    \

ifeq (,$(filter $(TARGET_CPU_VARIANT),$(SUPPORTED_ARMV8_CPU_VARIANTS)))
$(error "Unsupported ARMv8 CPU VARIANT. Supported CPU_VARIANTs : $(SUPPORTED_ARMV8_CPU_VARIANTS)")
endif

# CPU Variant.
ifneq (,$(filter generic denver64,$(TARGET_$(combo_2nd_arch_prefix)CPU_VARIANT)))
  # TODO(NVIDIA): Update GCC/Clang to a version that has support for denver/denver64.
  arch_variant_cflags := -march=armv8-a
else
ifneq (,$(filter cortex-a72,$(TARGET_$(combo_2nd_arch_prefix)CPU_VARIANT)))
  # TODO(ARM): Update GCC/Clang to a version that has support for these CPUs.
  arch_variant_cflags := -mcpu=cortex-a57
else
  arch_variant_cflags := -mcpu=$(strip $(TARGET_$(combo_2nd_arch_prefix)CPU_VARIANT))
endif
endif
