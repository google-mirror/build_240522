# Configuration for Linux on ARM.
# Generating binaries for the ARMv7-a architecture and higher with NEON.
#
ARCH_ARM_HAVE_ARMV7A            := true
ARCH_ARM_HAVE_VFP               := true
ARCH_ARM_HAVE_VFP_D32           := true
ARCH_ARM_HAVE_NEON              := true

SUPPORTED_ARMV7_CPU_VARIANTS := \
  generic                   \
  cortex-a5                 \
  cortex-a7                 \
  cortex-a8                 \
  cortex-a9                 \
  cortex-a12                \
  cortex-a15                \
  cortex-a17                \
  krait                     \

SUPPORTED_ARMV8_CPU_VARIANTS := \
  generic                   \
  cortex-a53                \
  cortex-a57                \
  cortex-a72                \
  denver                    \

ARMV7_CPU_VARIANTS_WITH_LPAE := \
  cortex-a7                 \
  cortex-a12                \
  cortex-a15                \
  cortex-a17                \
  krait                     \

ARMV7_CPU_VARIANTS_WITH_FPU_FP16 := \
  cortex-a5                 \
  cortex-a9                 \

ARMV7_CPU_VARIANTS_WITH_FPU_VFPV4 := \
  cortex-a7                 \
  cortex-a12                \
  cortex-a15                \
  cortex-a17                \
  krait                     \

ifneq (,$(filter $(TARGET_ARCH) $(TARGET_2ND_ARCH), arm64))
  ifeq (,$(filter $(TARGET_$(combo_2nd_arch_prefix)CPU_VARIANT),$(SUPPORTED_ARMV8_CPU_VARIANTS)))
    $(error "Unsupported ARMv8 CPU_VARIANT. Supported CPU_VARIANTs : $(SUPPORTED_ARMV8_CPU_VARIANTS)")
  endif
else
  ifeq (,$(filter $(TARGET_$(combo_2nd_arch_prefix)CPU_VARIANT),$(SUPPORTED_ARMV7_CPU_VARIANTS)))
    $(error "Unsupported ARMv7 CPU_VARIANT. Supported CPU_VARIANTs : $(SUPPORTED_ARMV7_CPU_VARIANTS)")
  endif
endif

# CPU Variant
ifneq (,$(filter $(TARGET_ARCH) $(TARGET_2ND_ARCH), arm64))
  # Assume the closest microarchitecture to a ARMv8.
  arch_variant_cflags := -mcpu=cortex-a15
else
ifeq ($(strip $(TARGET_$(combo_2nd_arch_prefix)CPU_VARIANT)), generic)
  arch_variant_cflags := -march=armv7-a
else
ifneq (,$(filter cortex-a17 krait,$(TARGET_$(combo_2nd_arch_prefix)CPU_VARIANT)))
  # TODO(ARM): Update GCC/Clang to a version that has support for cortex-a17.
  # TODO(QUALCOMM): Update GCC/Clang to a version that has support for krait.
  # TODO: krait is not a cortex-a15, we set the variant to cortex-a15 so that
  #       hardware divide operations are generated. This should be removed and a
  #       krait CPU variant added to GCC. For clang we specify -mcpu for krait in
  #       core/clang/arm.mk.
  arch_variant_cflags := -mcpu=cortex-a15
else
  arch_variant_cflags := -mcpu=$(strip $(TARGET_$(combo_2nd_arch_prefix)CPU_VARIANT))
endif
endif
endif

# FPU Variant
ifneq (,$(filter $(TARGET_ARCH) $(TARGET_2ND_ARCH), arm64))
  arch_variant_cflags += -mfpu=neon-fp-armv8
else
ifneq (,$(filter $(TARGET_$(combo_2nd_arch_prefix)CPU_VARIANT),$(ARMV7_CPU_VARIANTS_WITH_FPU_FP16)))
  arch_variant_cflags += -mfpu=neon-fp16
else
ifneq (,$(filter $(TARGET_$(combo_2nd_arch_prefix)CPU_VARIANT),$(ARMV7_CPU_VARIANTS_WITH_FPU_VFPV4)))
  arch_variant_cflags += -mfpu=neon-vfpv4
else
  arch_variant_cflags += -mfpu=neon
endif
endif
endif

# Soft Float ABI
arch_variant_cflags += -mfloat-abi=softfp

# LPAE Support
# Fake an ARM compiler flag as these processors support LPAE which GCC/clang don't advertise.
ifneq (,$(filter $(TARGET_ARCH) $(TARGET_2ND_ARCH), arm64))
  arch_variant_cflags += -D__ARM_FEATURE_LPAE=1
else
ifneq (,$(filter $(TARGET_$(combo_2nd_arch_prefix)CPU_VARIANT),$(ARMV7_CPU_VARIANTS_WITH_LPAE)))
  arch_variant_cflags += -D__ARM_FEATURE_LPAE=1
endif
endif


# Errata fixes
ifneq (,$(filter $(TARGET_ARCH) $(TARGET_2ND_ARCH), arm64))
  arch_variant_ldflags := -Wl,--no-fix-cortex-a8
else
ifneq (,$(filter generic cortex-a8,$(TARGET_$(combo_2nd_arch_prefix)CPU_VARIANT)))
  # Generic ARM might be a Cortex A8 -- better safe than sorry
  arch_variant_ldflags := -Wl,--fix-cortex-a8
else
  arch_variant_ldflags := -Wl,--no-fix-cortex-a8
endif
endif
