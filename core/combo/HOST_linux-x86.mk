#
# Copyright (C) 2006 The Android Open Source Project
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

# Configuration for builds hosted on linux-x86.
# Included by combo/select.mk

# $(1): The file to check
define get-file-size
stat --format "%s" "$(1)" | tr -d '\n'
endef

ifeq ($(strip $(HOST_TOOLCHAIN_PREFIX)),)
HOST_TOOLCHAIN_PREFIX=prebuilts/gcc/linux-x86/host/x86_64-linux-glibc2.11-4.6/bin/x86_64-linux-
endif
# Don't do anything if the toolchain is not there
ifneq (,$(strip $(wildcard $(HOST_TOOLCHAIN_PREFIX)gcc)))
HOST_CC  := $(HOST_TOOLCHAIN_PREFIX)gcc
HOST_CXX := $(HOST_TOOLCHAIN_PREFIX)g++
HOST_AR  := $(HOST_TOOLCHAIN_PREFIX)ar
endif # $(HOST_TOOLCHAIN_PREFIX)gcc exists

HOST_TOOLCHAIN_FOR_CLANG := prebuilts/gcc/linux-x86/host/x86_64-linux-glibc2.11-4.6/
ifneq ($(strip $(BUILD_HOST_64bit)),)
# By default we build everything in 32-bit, because it gives us
# more consistency between the host tools and the target.
# BUILD_HOST_64bit=1 overrides it for tool like emulator
# which can benefit from 64-bit host arch.
HOST_GLOBAL_CFLAGS += -m64 -Wa,--noexecstack
HOST_GLOBAL_LDFLAGS += -m64 -Wl,-z,noexecstack
else
# We expect SSE3 floating point math.
HOST_GLOBAL_CFLAGS += -mstackrealign -msse3 -mfpmath=sse -m32 -Wa,--noexecstack
HOST_GLOBAL_LDFLAGS += -m32 -Wl,-z,noexecstack
endif # BUILD_HOST_64bit

ifneq ($(strip $(BUILD_HOST_static)),)
# Statically-linked binaries are desirable for sandboxed environment
HOST_GLOBAL_LDFLAGS += -static
endif # BUILD_HOST_static

HOST_GLOBAL_CFLAGS += -fPIC \
    -include $(call select-android-config-h,linux-x86)

# Disable new longjmp in glibc 2.11 and later. See bug 2967937.
HOST_GLOBAL_CFLAGS += -U_FORTIFY_SOURCE -D_FORTIFY_SOURCE=0

# Workaround differences in inttypes.h between host and target.
# See bug 12708004.
HOST_GLOBAL_CFLAGS += -D__STDC_FORMAT_MACROS -D__STDC_CONSTANT_MACROS

HOST_NO_UNDEFINED_LDFLAGS := -Wl,--no-undefined
