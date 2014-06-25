# Configuration for Android on mips64r6.

ARCH_MIPS_HAS_FPU	:=true
ARCH_HAVE_ALIGNED_DOUBLES :=true
arch_variant_cflags := \
    -EL \
    -march=mips64r6 \
    -mtune=mips64r6 \
    -mips64r6 \
    -mhard-float \
    -msynci

arch_variant_ldflags := \
    -EL
