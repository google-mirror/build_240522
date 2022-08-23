load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

# Skylib provides common utilities for writing bazel rules and functions.
# For docs see https://github.com/bazelbuild/bazel-skylib/blob/main/README.md
http_archive(
    name = "bazel_skylib",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-skylib/releases/download/1.2.1/bazel-skylib-1.2.1.tar.gz",
        "https://github.com/bazelbuild/bazel-skylib/releases/download/1.2.1/bazel-skylib-1.2.1.tar.gz",
    ],
    sha256 = "f7be3474d42aae265405a592bb7da8e171919d74c16f082a5457840f06054728",
)
load("@bazel_skylib//:workspace.bzl", "bazel_skylib_workspace")
bazel_skylib_workspace()

# Repository that provides the clang compilers
# Bind the directory as a local repo so we can use our own BUILD file without
# touching the one added by go/roboleaf.
new_local_repository(
    name = "clang",
    path = "prebuilts/clang/host",
    build_file = "build/bazel/toolchains/cc/BUILD.clang",
)

# Repository that provides include / libs from GCC
# Bind the directory as a local repo so we can use our own BUILD file without
# touching the one added by go/roboleaf.
new_local_repository(
    name = "gcc_lib",
    path = "prebuilts/gcc/linux-x86/host",
    build_file = "build/bazel/toolchains/cc/BUILD.gcc_lib",
)
