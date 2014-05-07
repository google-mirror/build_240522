ifeq ($(HOST_OS),darwin)
# nothing required here yet
endif

ifeq ($(HOST_OS),linux)
CLANG_CONFIG_x86_LINUX_HOST_EXTRA_ASFLAGS := \
  --gcc-toolchain=$(HOST_TOOLCHAIN_FOR_CLANG) \
  --sysroot=$(HOST_TOOLCHAIN_FOR_CLANG)/sysroot \
  -no-integrated-as

CLANG_CONFIG_x86_LINUX_HOST_EXTRA_CFLAGS := \
  --gcc-toolchain=$(HOST_TOOLCHAIN_FOR_CLANG) \
  -no-integrated-as

ifneq ($(strip $(BUILD_HOST_64bit)),)
CLANG_CONFIG_x86_LINUX_HOST_EXTRA_CPPFLAGS :=   \
  --gcc-toolchain=$(HOST_TOOLCHAIN_FOR_CLANG) \
  --sysroot=$(HOST_TOOLCHAIN_FOR_CLANG)/sysroot \
  -isystem $(HOST_TOOLCHAIN_FOR_CLANG)/x86_64-linux/include/c++/4.6 \
  -isystem $(HOST_TOOLCHAIN_FOR_CLANG)/x86_64-linux/include/c++/4.6/x86_64-linux \
  -isystem $(HOST_TOOLCHAIN_FOR_CLANG)/x86_64-linux/include/c++/4.6/backward \
  -no-integrated-as

CLANG_CONFIG_x86_LINUX_HOST_EXTRA_LDFLAGS := \
  --gcc-toolchain=$(HOST_TOOLCHAIN_FOR_CLANG) \
  --sysroot=$(HOST_TOOLCHAIN_FOR_CLANG)/sysroot \
  -B$(HOST_TOOLCHAIN_FOR_CLANG)/x86_64-linux/bin \
  -B$(HOST_TOOLCHAIN_FOR_CLANG)/lib/gcc/x86_64-linux/4.6 \
  -L$(HOST_TOOLCHAIN_FOR_CLANG)/lib/gcc/x86_64-linux/4.6 \
  -L$(HOST_TOOLCHAIN_FOR_CLANG)/x86_64-linux/lib64/ \
  -no-integrated-as
else
CLANG_CONFIG_x86_LINUX_HOST_EXTRA_CPPFLAGS :=   \
  --gcc-toolchain=$(HOST_TOOLCHAIN_FOR_CLANG) \
  --sysroot=$(HOST_TOOLCHAIN_FOR_CLANG)/sysroot \
  -isystem $(HOST_TOOLCHAIN_FOR_CLANG)/x86_64-linux/include/c++/4.6 \
  -isystem $(HOST_TOOLCHAIN_FOR_CLANG)/x86_64-linux/include/c++/4.6/x86_64-linux/32 \
  -isystem $(HOST_TOOLCHAIN_FOR_CLANG)/x86_64-linux/include/c++/4.6/backward \
  -no-integrated-as

CLANG_CONFIG_x86_LINUX_HOST_EXTRA_LDFLAGS := \
  --gcc-toolchain=$(HOST_TOOLCHAIN_FOR_CLANG) \
  --sysroot=$(HOST_TOOLCHAIN_FOR_CLANG)/sysroot \
  -B$(HOST_TOOLCHAIN_FOR_CLANG)/x86_64-linux/bin \
  -B$(HOST_TOOLCHAIN_FOR_CLANG)/lib/gcc/x86_64-linux/4.6/32 \
  -L$(HOST_TOOLCHAIN_FOR_CLANG)/lib/gcc/x86_64-linux/4.6/32 \
  -L$(HOST_TOOLCHAIN_FOR_CLANG)/x86_64-linux/lib32/ \
  -no-integrated-as
endif
endif

ifeq ($(HOST_OS),windows)
CLANG_CONFIG_x86_LINUX_HOST_EXTRA_ASFLAGS :=

CLANG_CONFIG_x86_LINUX_HOST_EXTRA_CFLAGS := \
  -isystem $(HOST_TOOLCHAIN_FOR_CLANG)/x86_64-w64-mingw32/include

ifneq ($(strip $(BUILD_HOST_64bit)),)
CLANG_CONFIG_x86_LINUX_HOST_EXTRA_CPPFLAGS := \
  -isystem $(HOST_TOOLCHAIN_FOR_CLANG)/x86_64-w64-mingw32/include/c++/4.8.2 \
  -isystem $(HOST_TOOLCHAIN_FOR_CLANG)/x86_64-w64-mingw32/include/c++/4.8.2/x86_64-w64-mingw32 \
  -isystem $(HOST_TOOLCHAIN_FOR_CLANG)/x86_64-w64-mingw32/include/c++/4.8.2/backward

CLANG_CONFIG_x86_LINUX_HOST_EXTRA_LDFLAGS := \
  -B$(HOST_TOOLCHAIN_FOR_CLANG)/x86_64-w64-mingw32/bin \
  -B$(HOST_TOOLCHAIN_FOR_CLANG)/lib/gcc/x86_64-w64-mingw32/4.8.2 \
  -L$(HOST_TOOLCHAIN_FOR_CLANG)/lib/gcc/x86_64-w64-mingw32/4.8.2 \
  -L$(HOST_TOOLCHAIN_FOR_CLANG)/x86_64-w64-mingw32/lib64
else
CLANG_CONFIG_x86_LINUX_HOST_EXTRA_CPPFLAGS := \
  -isystem $(HOST_TOOLCHAIN_FOR_CLANG)/x86_64-w64-mingw32/include/c++/4.8.2 \
  -isystem $(HOST_TOOLCHAIN_FOR_CLANG)/x86_64-w64-mingw32/include/c++/4.8.2/x86_64-w64-mingw32/32 \
  -isystem $(HOST_TOOLCHAIN_FOR_CLANG)/x86_64-w64-mingw32/include/c++/4.8.2/backward

CLANG_CONFIG_x86_LINUX_HOST_EXTRA_LDFLAGS := \
  -B$(HOST_TOOLCHAIN_FOR_CLANG)/x86_64-w64-mingw32/bin \
  -B$(HOST_TOOLCHAIN_FOR_CLANG)/lib/gcc/x86_64-w64-mingw32/4.8.2/32 \
  -L$(HOST_TOOLCHAIN_FOR_CLANG)/lib/gcc/x86_64-w64-mingw32/4.8.2/32 \
  -L$(HOST_TOOLCHAIN_FOR_CLANG)/x86_64-w64-mingw32/lib32
endif
endif