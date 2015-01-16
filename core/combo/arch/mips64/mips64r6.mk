# Configuration for Android on mips64r6.

ARCH_MIPS_REV6 := true
arch_variant_cflags := \
    -mips64r6 \
    -msynci

# Workaround until clang generates this itself:
arch_variant_cflags += -D__mips_isa_rev=6