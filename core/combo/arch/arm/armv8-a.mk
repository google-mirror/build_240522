# Configuration for Linux on ARM.
# Generating binaries for the ARMv8-a architecture and higher.
#
ARCH_ARM_HAVE_ARMV7A            := true
ARCH_ARM_HAVE_ARMV8A            := true
ARCH_ARM_HAVE_VFP               := true
ARCH_ARM_HAVE_VFP_D32           := true
ARCH_ARM_HAVE_NEON              := true

SUPPORTED_ARMV8_CPU_VARIANTS = \
  generic    \
  cortex-a53 \
  cortex-a57 \
  cortex-a57.cortex-a53 \
  denver     \

ifeq (,$(filter $(TARGET_$(combo_2nd_arch_prefix)CPU_VARIANT),$(SUPPORTED_ARMV8_CPU_VARIANTS)))
$(info "Supported ARMv8 CPU variants: $(SUPPORTED_ARMV8_CPU_VARIANTS)")
$(error "Unsupported ARMv8 CPU VARIANT")
endif


# CPU Variant.
# FIXME: This cannot work until we fix:
#  - arm-linux-androideabi-ld.gold support for AArch32 ARMv8A
#   (see https://bugs.launchpad.net/binutils-linaro/+bug/1154165)
#
#  - libc/arch-arm/include/machine/cpu-features.h
#   (add support for __ARM_ARCH_8A__, define __ARM_ARCH__ 8)
#
#  - NDK and all the other projects that do not know how to deal with __ARM_ARCH_8A__
#   (eg: external/chromium_org/base/atomicops_internals_arm_gcc.h
# ifneq (,$(filter generic denver,$(TARGET_$(combo_2nd_arch_prefix)CPU_VARIANT)))
# # FIXME (Nvidia): Add -mcpu=denver support to GCC and Clang.
# arch_variant_cflags := -march=armv8-a
# else
# arch_variant_cflags := -mcpu=$(strip $(TARGET_$(combo_2nd_arch_prefix)CPU_VARIANT))
# endif
#
# FIXME: For now we use the closest ARMv7 equivalent.
arch_variant_cflags += -mcpu=cortex-a15

# FPU Variant.
arch_variant_cflags += -mfpu=neon-fp-armv8

# Soft Float ABI.
arch_variant_cflags += -mfloat-abi=softfp

# LPAE Support.
# Fake an ARM compiler flag as these processors support LPAE which GCC/clang
# don't advertise.
arch_variant_cflags += -D__ARM_FEATURE_LPAE=1

# Errata fixes.
arch_variant_ldflags := -Wl,--no-fix-cortex-a8
