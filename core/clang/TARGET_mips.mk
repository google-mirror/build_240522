$(clang_2nd_arch_prefix)RS_TRIPLE := armv7-linux-androideabi
$(clang_2nd_arch_prefix)RS_TRIPLE_CFLAGS :=
RS_COMPAT_TRIPLE := mipsel-linux-android

$(clang_2nd_arch_prefix)TARGET_LIBPROFILE_RT := $(LLVM_RTLIB_PATH)/libclang_rt.profile-mipsel-android.a

# Address sanitizer clang config
$(clang_2nd_arch_prefix)ADDRESS_SANITIZER_LINKER := /system/bin/linker_asan
