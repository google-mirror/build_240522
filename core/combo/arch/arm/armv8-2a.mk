<<<<<<< HEAD   (5c8d84 Merge "Merge empty history for sparse-6676661-L8360000065797)
=======
# Configuration for Linux on ARM.
# Generating binaries for the ARMv8-2a architecture
#
# Many libraries are not aware of armv8-2a, and AArch32 is (almost) a superset
# of armv7-a-neon. So just let them think we are just like v7.
ARCH_ARM_HAVE_VFP               := true
ARCH_ARM_HAVE_VFP_D32           := true
ARCH_ARM_HAVE_NEON              := true
>>>>>>> BRANCH (a10c18 Merge "Version bump to RT11.201014.001.A1 [core/build_id.mk])
