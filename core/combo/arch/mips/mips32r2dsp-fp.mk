# Configuration for Android on MIPS.
# Generating binaries for MIPS32R2/hard-float/little-endian/dsp

ARCH_MIPS_HAS_DSP  	:=true
ARCH_MIPS_DSP_REV	:=1
ARCH_MIPS_HAS_FPU       :=true
ARCH_HAVE_ALIGNED_DOUBLES :=true
arch_variant_cflags := \
    -mips32r2 \
    -mfp32 \
    -modd-spreg \
    -mdsp \
    -msynci

# Workaround until clang generates this itself:
arch_variant_cflags += -D__mips_isa_rev=2

arch_variant_ldflags := \
    -Wl,-melf32ltsmip
