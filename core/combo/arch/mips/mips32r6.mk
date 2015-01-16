# Configuration for Android on MIPS.
# Generating binaries for MIPS32R6/hard-float/little-endian

ARCH_MIPS_REV6 := true
arch_variant_cflags := \
    -mips32r6 \
    -mfp64 \
    -mno-odd-spreg \
    -msynci

# Workaround until clang generates this itself:
arch_variant_cflags += -D__mips_isa_rev=6

arch_variant_ldflags := \
    -Wl,-melf32ltsmip
