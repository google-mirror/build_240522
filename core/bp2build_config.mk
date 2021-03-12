# All modules in these directories and their subdirectories are automatically
# converted by bp2build.
BP2BUILD_CONVERT_SUBTREE ?= \
    bionic \
    system/core/libcutils \
    system/logging/liblog

# ..except these directories, where modules are opted-in with:
#
#     bazel_module: { bp2build_available: true }
#
# This list should be empty. If not, attach a bug number with an explanation for
# why it's blocked.
#
# bionic/libc: b/182339414
# bionic/linker: b/182338959
BP2BUILD_CONVERT_OPT_IN_MODULES ?= \
    bionic/libc \
    bionic/linker
