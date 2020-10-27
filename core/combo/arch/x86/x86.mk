# This file contains feature macro definitions specific to the
# base 'x86' platform ABI.
#
# It is also used to build full_x86-eng / sdk_x86-eng platform images that
# are run in the emulator under KVM emulation (i.e. running directly on
# the host development machine's CPU).

# These features are optional and shall not be included in the base platform
# Otherwise, sdk_x86-eng system images might fail to run on some
# developer machines.
<<<<<<< HEAD   (5c8d84 Merge "Merge empty history for sparse-6676661-L8360000065797)
ARCH_X86_HAVE_SSSE3 := false
ARCH_X86_HAVE_MOVBE := false
ARCH_X86_HAVE_POPCNT := false
=======
>>>>>>> BRANCH (a10c18 Merge "Version bump to RT11.201014.001.A1 [core/build_id.mk])
