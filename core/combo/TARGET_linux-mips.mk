#
# Copyright (C) 2010 The Android Open Source Project
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

# Configuration for Linux on MIPS.
# Included by combo/select.mk

# You can set TARGET_ARCH_VARIANT to use an arch version other
# than mips32r2-fp. Each value should correspond to a file named
# $(BUILD_COMBOS)/arch/<name>.mk which must contain
# makefile variable definitions similar to the preprocessor
# defines in build/core/combo/include/arch/<combo>/AndroidConfig.h. Their
# purpose is to allow module Android.mk files to selectively compile
# different versions of code based upon the funtionality and
# instructions available in a given architecture version.
#
# The blocks also define specific arch_variant_cflags, which
# include defines, and compiler settings for the given architecture
# version.
#
ifeq ($(strip $(TARGET_$(combo_2nd_arch_prefix)ARCH_VARIANT)),)
TARGET_$(combo_2nd_arch_prefix)ARCH_VARIANT := mips32r2-fp
endif

# Decouple NDK library selection with platform compiler version
$(combo_2nd_arch_prefix)TARGET_NDK_GCC_VERSION := 4.8

ifeq ($(strip $(TARGET_GCC_VERSION_EXP)),)
$(combo_2nd_arch_prefix)TARGET_GCC_VERSION := 4.9
else
$(combo_2nd_arch_prefix)TARGET_GCC_VERSION := $(TARGET_GCC_VERSION_EXP)
endif

TARGET_ARCH_SPECIFIC_MAKEFILE := $(BUILD_COMBOS)/arch/$(TARGET_$(combo_2nd_arch_prefix)ARCH)/$(TARGET_$(combo_2nd_arch_prefix)ARCH_VARIANT).mk
ifeq ($(strip $(wildcard $(TARGET_ARCH_SPECIFIC_MAKEFILE))),)
$(error Unknown MIPS architecture variant: $(TARGET_$(combo_2nd_arch_prefix)ARCH_VARIANT))
endif

include $(TARGET_ARCH_SPECIFIC_MAKEFILE)
include $(BUILD_SYSTEM)/combo/fdo.mk

# You can set TARGET_TOOLS_PREFIX to get gcc from somewhere else
ifeq ($(strip $($(combo_2nd_arch_prefix)TARGET_TOOLS_PREFIX)),)
$(combo_2nd_arch_prefix)TARGET_TOOLCHAIN_ROOT := prebuilts/gcc/$(HOST_PREBUILT_TAG)/mips/mips64el-linux-android-$($(combo_2nd_arch_prefix)TARGET_GCC_VERSION)
$(combo_2nd_arch_prefix)TARGET_TOOLS_PREFIX := $($(combo_2nd_arch_prefix)TARGET_TOOLCHAIN_ROOT)/bin/mips64el-linux-android-
endif

$(combo_2nd_arch_prefix)TARGET_CC := $($(combo_2nd_arch_prefix)TARGET_TOOLS_PREFIX)gcc$(HOST_EXECUTABLE_SUFFIX)
$(combo_2nd_arch_prefix)TARGET_CXX := $($(combo_2nd_arch_prefix)TARGET_TOOLS_PREFIX)g++$(HOST_EXECUTABLE_SUFFIX)
$(combo_2nd_arch_prefix)TARGET_AR := $($(combo_2nd_arch_prefix)TARGET_TOOLS_PREFIX)ar$(HOST_EXECUTABLE_SUFFIX)
$(combo_2nd_arch_prefix)TARGET_OBJCOPY := $($(combo_2nd_arch_prefix)TARGET_TOOLS_PREFIX)objcopy$(HOST_EXECUTABLE_SUFFIX)
$(combo_2nd_arch_prefix)TARGET_LD := $($(combo_2nd_arch_prefix)TARGET_TOOLS_PREFIX)ld$(HOST_EXECUTABLE_SUFFIX)
$(combo_2nd_arch_prefix)TARGET_READELF := $($(combo_2nd_arch_prefix)TARGET_TOOLS_PREFIX)readelf$(HOST_EXECUTABLE_SUFFIX)
$(combo_2nd_arch_prefix)TARGET_STRIP := $($(combo_2nd_arch_prefix)TARGET_TOOLS_PREFIX)strip$(HOST_EXECUTABLE_SUFFIX)

$(combo_2nd_arch_prefix)TARGET_NO_UNDEFINED_LDFLAGS := -Wl,--no-undefined

TARGET_mips_CFLAGS :=	-O2 \
			-fomit-frame-pointer \
			-fno-strict-aliasing    \
			-funswitch-loops

# Set FORCE_MIPS_DEBUGGING to "true" in your buildspec.mk
# or in your environment to gdb debugging easier.
# Don't forget to do a clean build.
ifeq ($(FORCE_MIPS_DEBUGGING),true)
  TARGET_mips_CFLAGS += -fno-omit-frame-pointer
endif

android_config_h := $(call select-android-config-h,linux-mips)

$(combo_2nd_arch_prefix)TARGET_GLOBAL_CFLAGS += \
			$(TARGET_mips_CFLAGS) \
			-U__unix -U__unix__ -Umips \
			-ffunction-sections \
			-fdata-sections \
			-funwind-tables \
			-Wa,--noexecstack \
			-Werror=format-security \
			-D_FORTIFY_SOURCE=2 \
			-no-canonical-prefixes \
			-fno-canonical-system-headers \
			$(arch_variant_cflags) \
			-include $(android_config_h) \
			-I $(dir $(android_config_h))

ifneq ($(ARCH_MIPS_PAGE_SHIFT),)
$(combo_2nd_arch_prefix)TARGET_GLOBAL_CFLAGS += -DPAGE_SHIFT=$(ARCH_MIPS_PAGE_SHIFT)
endif

$(combo_2nd_arch_prefix)TARGET_GLOBAL_LDFLAGS += \
			-Wl,-z,noexecstack \
			-Wl,-z,relro \
			-Wl,-z,now \
			-Wl,--warn-shared-textrel \
			-Wl,--fatal-warnings \
			$(arch_variant_ldflags)

$(combo_2nd_arch_prefix)TARGET_GLOBAL_CPPFLAGS += -fvisibility-inlines-hidden

# More flags/options can be added here
$(combo_2nd_arch_prefix)TARGET_RELEASE_CFLAGS := \
			-DNDEBUG \
			-g \
			-Wstrict-aliasing=2 \
			-fgcse-after-reload \
			-frerun-cse-after-loop \
			-frename-registers

libc_root := bionic/libc
libm_root := bionic/libm


## on some hosts, the target cross-compiler is not available so do not run this command
ifneq ($(wildcard $($(combo_2nd_arch_prefix)TARGET_CC)),)
# We compile with the global cflags to ensure that
# any flags which affect libgcc are correctly taken
# into account.
$(combo_2nd_arch_prefix)TARGET_LIBGCC := \
  $(shell $($(combo_2nd_arch_prefix)TARGET_CC) $($(combo_2nd_arch_prefix)TARGET_GLOBAL_CFLAGS) -print-file-name=libgcc.a)
$(combo_2nd_arch_prefix)TARGET_LIBATOMIC := \
  $(shell $($(combo_2nd_arch_prefix)TARGET_CC) $($(combo_2nd_arch_prefix)TARGET_GLOBAL_CFLAGS) -print-file-name=libatomic.a)
LIBGCC_EH := $(shell $($(combo_2nd_arch_prefix)TARGET_CC) $($(combo_2nd_arch_prefix)TARGET_GLOBAL_CFLAGS) -print-file-name=libgcc_eh.a)
ifneq ($(LIBGCC_EH),libgcc_eh.a)
  $(combo_2nd_arch_prefix)TARGET_LIBGCC += $(LIBGCC_EH)
endif
$(combo_2nd_arch_prefix)TARGET_LIBGCOV := $(shell $($(combo_2nd_arch_prefix)TARGET_CC) $($(combo_2nd_arch_prefix)TARGET_GLOBAL_CFLAGS) \
        --print-file-name=libgcov.a)
endif

KERNEL_HEADERS_COMMON := $(libc_root)/kernel/uapi
KERNEL_HEADERS_ARCH   := $(libc_root)/kernel/uapi/asm-mips # mips covers both mips and mips64.
KERNEL_HEADERS := $(KERNEL_HEADERS_COMMON) $(KERNEL_HEADERS_ARCH)

$(combo_2nd_arch_prefix)TARGET_C_INCLUDES := \
	$(libc_root)/arch-mips/include \
	$(libc_root)/include \
	$(KERNEL_HEADERS) \
	$(libm_root)/include \
	$(libm_root)/include/mips \

$(combo_2nd_arch_prefix)TARGET_CRTBEGIN_STATIC_O := $($(combo_2nd_arch_prefix)TARGET_OUT_INTERMEDIATE_LIBRARIES)/crtbegin_static.o
$(combo_2nd_arch_prefix)TARGET_CRTBEGIN_DYNAMIC_O := $($(combo_2nd_arch_prefix)TARGET_OUT_INTERMEDIATE_LIBRARIES)/crtbegin_dynamic.o
$(combo_2nd_arch_prefix)TARGET_CRTEND_O := $($(combo_2nd_arch_prefix)TARGET_OUT_INTERMEDIATE_LIBRARIES)/crtend_android.o

$(combo_2nd_arch_prefix)TARGET_CRTBEGIN_SO_O := $($(combo_2nd_arch_prefix)TARGET_OUT_INTERMEDIATE_LIBRARIES)/crtbegin_so.o
$(combo_2nd_arch_prefix)TARGET_CRTEND_SO_O := $($(combo_2nd_arch_prefix)TARGET_OUT_INTERMEDIATE_LIBRARIES)/crtend_so.o

$(combo_2nd_arch_prefix)TARGET_STRIP_MODULE:=true

$(combo_2nd_arch_prefix)TARGET_DEFAULT_SYSTEM_SHARED_LIBRARIES := libc libm

$(combo_2nd_arch_prefix)TARGET_CUSTOM_LD_COMMAND := true

define $(combo_2nd_arch_prefix)transform-o-to-shared-lib-inner
$(hide) $(PRIVATE_CXX) \
	-nostdlib -Wl,-soname,$(notdir $@) \
	-Wl,--gc-sections \
	$(if $(filter true,$(PRIVATE_CLANG)),-shared,-Wl,-shared) \
	$(PRIVATE_TARGET_GLOBAL_LD_DIRS) \
	$(if $(filter true,$(PRIVATE_NO_CRT)),,$(PRIVATE_TARGET_CRTBEGIN_SO_O)) \
	$(PRIVATE_ALL_OBJECTS) \
	-Wl,--whole-archive \
	$(call normalize-target-libraries,$(PRIVATE_ALL_WHOLE_STATIC_LIBRARIES)) \
	-Wl,--no-whole-archive \
	$(if $(PRIVATE_GROUP_STATIC_LIBRARIES),-Wl$(comma)--start-group) \
	$(call normalize-target-libraries,$(PRIVATE_ALL_STATIC_LIBRARIES)) \
	$(if $(PRIVATE_GROUP_STATIC_LIBRARIES),-Wl$(comma)--end-group) \
	$(if $(TARGET_BUILD_APPS),$(PRIVATE_TARGET_LIBGCC)) \
	$(PRIVATE_TARGET_FDO_LIB) \
	$(call normalize-target-libraries,$(PRIVATE_ALL_SHARED_LIBRARIES)) \
	-o $@ \
	$(PRIVATE_TARGET_GLOBAL_LDFLAGS) \
	$(PRIVATE_LDFLAGS) \
	$(PRIVATE_TARGET_LIBATOMIC) \
	$(if $(filter true,$(NATIVE_COVERAGE)),$(PRIVATE_TARGET_LIBGCOV)) \
	$(if $(PRIVATE_LIBCXX),,$(PRIVATE_TARGET_LIBGCC)) \
	$(if $(filter true,$(PRIVATE_NO_CRT)),,$(PRIVATE_TARGET_CRTEND_SO_O)) \
	$(PRIVATE_LDLIBS)
endef

define $(combo_2nd_arch_prefix)transform-o-to-executable-inner
$(hide) $(PRIVATE_CXX) -nostdlib -Bdynamic -pie \
	-Wl,-dynamic-linker,/system/bin/linker \
	-Wl,--gc-sections \
	-Wl,-z,nocopyreloc \
	$(PRIVATE_TARGET_GLOBAL_LD_DIRS) \
	-Wl,-rpath-link=$(PRIVATE_TARGET_OUT_INTERMEDIATE_LIBRARIES) \
	$(if $(filter true,$(PRIVATE_NO_CRT)),,$(PRIVATE_TARGET_CRTBEGIN_DYNAMIC_O)) \
	$(PRIVATE_ALL_OBJECTS) \
	-Wl,--whole-archive \
	$(call normalize-target-libraries,$(PRIVATE_ALL_WHOLE_STATIC_LIBRARIES)) \
	-Wl,--no-whole-archive \
	$(if $(PRIVATE_GROUP_STATIC_LIBRARIES),-Wl$(comma)--start-group) \
	$(call normalize-target-libraries,$(PRIVATE_ALL_STATIC_LIBRARIES)) \
	$(if $(PRIVATE_GROUP_STATIC_LIBRARIES),-Wl$(comma)--end-group) \
	$(if $(TARGET_BUILD_APPS),$(PRIVATE_TARGET_LIBGCC)) \
	$(PRIVATE_TARGET_FDO_LIB) \
	$(call normalize-target-libraries,$(PRIVATE_ALL_SHARED_LIBRARIES)) \
	-o $@ \
	$(PRIVATE_TARGET_GLOBAL_LDFLAGS) \
	$(PRIVATE_LDFLAGS) \
	$(PRIVATE_TARGET_LIBATOMIC) \
	$(if $(filter true,$(NATIVE_COVERAGE)),$(PRIVATE_TARGET_LIBGCOV)) \
	$(if $(PRIVATE_LIBCXX),,$(PRIVATE_TARGET_LIBGCC)) \
	$(if $(filter true,$(PRIVATE_NO_CRT)),,$(PRIVATE_TARGET_CRTEND_O)) \
	$(PRIVATE_LDLIBS)
endef

define $(combo_2nd_arch_prefix)transform-o-to-static-executable-inner
$(hide) $(PRIVATE_CXX) -nostdlib -Bstatic \
	-Wl,--gc-sections \
	-o $@ \
	$(PRIVATE_TARGET_GLOBAL_LD_DIRS) \
	$(if $(filter true,$(PRIVATE_NO_CRT)),,$(PRIVATE_TARGET_CRTBEGIN_STATIC_O)) \
	$(PRIVATE_TARGET_GLOBAL_LDFLAGS) \
	$(PRIVATE_LDFLAGS) \
	$(PRIVATE_ALL_OBJECTS) \
	-Wl,--whole-archive \
	$(call normalize-target-libraries,$(PRIVATE_ALL_WHOLE_STATIC_LIBRARIES)) \
	-Wl,--no-whole-archive \
	$(call normalize-target-libraries,$(filter-out %libc_nomalloc.a,$(filter-out %libc.a,$(PRIVATE_ALL_STATIC_LIBRARIES)))) \
	-Wl,--start-group \
	$(call normalize-target-libraries,$(filter %libc.a,$(PRIVATE_ALL_STATIC_LIBRARIES))) \
	$(call normalize-target-libraries,$(filter %libc_nomalloc.a,$(PRIVATE_ALL_STATIC_LIBRARIES))) \
	$(PRIVATE_TARGET_FDO_LIB) \
	$(PRIVATE_TARGET_LIBATOMIC) \
	$(if $(filter true,$(NATIVE_COVERAGE)),$(PRIVATE_TARGET_LIBGCOV)) \
	$(if $(PRIVATE_LIBCXX),,$(PRIVATE_TARGET_LIBGCC)) \
	-Wl,--end-group \
	$(if $(filter true,$(PRIVATE_NO_CRT)),,$(PRIVATE_TARGET_CRTEND_O))
endef
