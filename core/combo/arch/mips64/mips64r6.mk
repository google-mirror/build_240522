# Configuration for Android on mips64r6.

ARCH_MIPS_HAS_FPU	:=true
ARCH_HAVE_ALIGNED_DOUBLES :=true
ARCH_MIPS_REV6 := true
arch_variant_cflags := \
    -mips64r6 \
    -msynci

