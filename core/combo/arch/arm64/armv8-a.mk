ARCH_ARM_HAVE_ARMV8A    := true
ARCH_ARM_HAVE_NEON              := true

SUPPORTED_ARMV8_CPU_VARIANTS = \
  generic    \
  cortex-a53 \
  cortex-a57 \
  cortex-a57.cortex-a53 \
  denver64   \

arch_variant_cflags :=

ifeq (,$(filter $(TARGET_CPU_VARIANT),$(SUPPORTED_ARMV8_CPU_VARIANTS)))
$(info "Supported ARMv8 CPU variants: $(SUPPORTED_ARMV8_CPU_VARIANTS)")
$(error "Unsupported ARMv8 CPU VARIANT")
endif

# CPU Variant.
ifneq (,$(filter generic cortex-a57.cortex-a53 denver64,$(TARGET_$(combo_2nd_arch_prefix)CPU_VARIANT)))
  # FIXME (Nvidia): Add -mcpu=denver support to GCC and Clang.
  # FIXME (ARM): Add -mcpu=cortex-a57.cortex-a53 to Clang
  arch_variant_cflags := -march=armv8-a
else
  arch_variant_cflags := -mcpu=$(strip $(TARGET_$(combo_2nd_arch_prefix)CPU_VARIANT))
endif
