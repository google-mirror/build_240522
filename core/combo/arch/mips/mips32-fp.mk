# Configuration for Android on MIPS.
# Generating binaries for MIPS32/hard-float/little-endian

ARCH_MIPS_HAS_FPU	:=true
ARCH_HAVE_ALIGNED_DOUBLES :=true
arch_variant_cflags := \
    -mips32 \
    -mfp32 \
    -modd-spreg \

# Workaround until clang generates this itself:
arch_variant_cflags += -D__mips_isa_rev=1

arch_variant_ldflags := \
    -Wl,-melf32ltsmip
