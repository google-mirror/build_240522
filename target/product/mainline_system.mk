#
# Copyright (C) 2017 The Android Open-Source Project
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

$(call inherit-product, $(SRC_TARGET_DIR)/product/generic_no_telephony.mk)

PRODUCT_NAME := aosp_system
PRODUCT_SHIPPING_API_LEVEL := 28

_shell_and_utilities_files := \
  bin/acpi \
  bin/base64 \
  bin/basename \
  bin/blockdev \
  bin/cal \
  bin/cat \
  bin/chcon \
  bin/chgrp \
  bin/chmod \
  bin/chown \
  bin/chroot \
  bin/chrt \
  bin/cksum \
  bin/clear \
  bin/cmp \
  bin/comm \
  bin/cp \
  bin/cpio \
  bin/cut \
  bin/date \
  bin/dd \
  bin/df \
  bin/diff \
  bin/dirname \
  bin/dmesg \
  bin/dos2unix \
  bin/du \
  bin/echo \
  bin/env \
  bin/expand \
  bin/expr \
  bin/fallocate \
  bin/false \
  bin/file \
  bin/find \
  bin/flock \
  bin/fmt \
  bin/free \
  bin/getenforce \
  bin/getevent \
  bin/getprop \
  bin/groups \
  bin/gunzip \
  bin/gzip \
  bin/head \
  bin/hostname \
  bin/hwclock \
  bin/id \
  bin/ifconfig \
  bin/inotifyd \
  bin/insmod \
  bin/ionice \
  bin/iorenice \
  bin/kill \
  bin/killall \
  bin/ln \
  bin/load_policy \
  bin/log \
  bin/logname \
  bin/losetup \
  bin/ls \
  bin/lsmod \
  bin/lsof \
  bin/lspci \
  bin/lsusb \
  bin/md5sum \
  bin/microcom \
  bin/mkdir \
  bin/mkfifo \
  bin/mknod \
  bin/mkswap \
  bin/mktemp \
  bin/modinfo \
  bin/modprobe \
  bin/more \
  bin/mount \
  bin/mountpoint \
  bin/mv \
  bin/nc \
  bin/netcat \
  bin/netstat \
  bin/newfs_msdos \
  bin/nice \
  bin/nl \
  bin/nohup \
  bin/nsenter \
  bin/od \
  bin/paste \
  bin/patch \
  bin/pgrep \
  bin/pidof \
  bin/pkill \
  bin/pmap \
  bin/printenv \
  bin/printf \
  bin/ps \
  bin/pwd \
  bin/readlink \
  bin/realpath \
  bin/renice \
  bin/restorecon \
  bin/rm \
  bin/rmdir \
  bin/rmmod \
  bin/runcon \
  bin/sed \
  bin/sendevent \
  bin/seq \
  bin/setenforce \
  bin/setprop \
  bin/setsid \
  bin/sh \
  bin/sha1sum \
  bin/sha224sum \
  bin/sha256sum \
  bin/sha384sum \
  bin/sha512sum \
  bin/sleep \
  bin/sort \
  bin/split \
  bin/start \
  bin/stat \
  bin/stop \
  bin/strings \
  bin/stty \
  bin/swapoff \
  bin/swapon \
  bin/sync \
  bin/sysctl \
  bin/tac \
  bin/tail \
  bin/tar \
  bin/taskset \
  bin/tee \
  bin/time \
  bin/timeout \
  bin/toolbox \
  bin/top \
  bin/touch \
  bin/tr \
  bin/true \
  bin/truncate \
  bin/tty \
  bin/ulimit \
  bin/umount \
  bin/uname \
  bin/uniq \
  bin/unix2dos \
  bin/unshare \
  bin/uptime \
  bin/usleep \
  bin/uudecode \
  bin/uuencode \
  bin/uuidgen \
  bin/vmstat \
  bin/wc \
  bin/which \
  bin/whoami \
  bin/xargs \
  bin/xxd \
  bin/yes \
  bin/zcat \

_shell_and_utilities_whitelist := \
  $(foreach f,$(_shell_and_utilities_files),recovery/root/system/$(f) vendor/$(f))

_shell_and_utilities_whitelist += \
  vendor/bin/awk \
  vendor/bin/grep \
  vendor/bin/logwrapper \
  vendor/etc/mkshrc


_unknown_source := \
  recovery/root/etc/ld.config.txt \
  recovery/root/sbin/e2fsdroid_static \
  recovery/root/sbin/mke2fs_static \
  recovery/root/sbin/recovery \
  recovery/root/system/bin/adbd \
  recovery/root/system/bin/linker \
  recovery/root/system/bin/linker_asan \
  recovery/root/system/lib/ld-android.so \
  recovery/root/system/lib/libbase.so \
  recovery/root/system/lib/libc++.so \
  recovery/root/system/lib/libc.so \
  recovery/root/system/lib/libcrypto.so \
  recovery/root/system/lib/libcutils.so \
  recovery/root/system/lib/libdl.so \
  recovery/root/system/lib/liblog.so \
  recovery/root/system/lib/libm.so \
  recovery/root/system/lib/libpackagelistparser.so \
  recovery/root/system/lib/libpcre2.so \
  recovery/root/system/lib/libselinux.so \
  recovery/root/system/lib/libz.so \
  vendor/etc/group \
  vendor/etc/passwd \



_shell_and_utilities_whitelist +=  recovery/root/system/bin/toybox vendor/bin/toybox_vendor

_selinux_policy_whitelist := \
  vendor/etc/selinux/vndservice_contexts \
  vendor/etc/selinux/plat_pub_versioned.cil \
  vendor/etc/selinux/plat_sepolicy_vers.txt \
  vendor/etc/selinux/precompiled_sepolicy \
  vendor/etc/selinux/precompiled_sepolicy.plat_and_mapping.sha256 \
  vendor/etc/selinux/vendor_file_contexts \
  vendor/etc/selinux/vendor_hwservice_contexts \
  vendor/etc/selinux/vendor_mac_permissions.xml \
  vendor/etc/selinux/vendor_property_contexts \
  vendor/etc/selinux/vendor_seapp_contexts \
  vendor/etc/selinux/vendor_sepolicy.cil \

_base_mk_whitelist := \
  $(_shell_and_utilities_whitelist) \
  $(_selinux_policy_whitelist) \
  $(_unknown_source) \
  root/init \
  root/init.environ.rc \
  root/init.rc \
  root/sbin/charger \
  vendor/bin/hw/android.hardware.configstore@1.0-service \
  vendor/bin/vndservice \
  vendor/bin/vndservicemanager \
  vendor/etc/fs_config_dirs \
  vendor/etc/fs_config_files \
  vendor/etc/init/android.hardware.configstore@1.0-service.rc \
  vendor/etc/init/vndservicemanager.rc \
  vendor/etc/vintf/compatibility_matrix.xml \
  vendor/lib/hw/gralloc.default.so \
  vendor/lib/libhwminijail.so \
  \
  vendor/bin/hw/android.hardware.cas@1.0-service \
  vendor/bin/hw/android.hardware.media.omx@1.0-service \
  vendor/etc/init/android.hardware.cas@1.0-service.rc \
  vendor/etc/init/android.hardware.media.omx@1.0-service.rc \
  vendor/lib/libavservices_minijail_vendor.so \
  vendor/lib/libeffects.so \
  vendor/lib/libeffectsconfig.so \
  vendor/lib/libmediacodecservice.so \
  vendor/lib/libreference-ril.so \
  vendor/lib/libril.so \
  vendor/lib/librilutils.so \
  vendor/lib/mediacas/libclearkeycasplugin.so \
  vendor/lib/mediadrm/libdrmclearkeyplugin.so \
  vendor/lib/soundfx/libbundlewrapper.so \
  vendor/lib/soundfx/libeffectproxy.so \
  vendor/lib/soundfx/libldnhncr.so \
  vendor/lib/soundfx/libreverbwrapper.so \
  vendor/lib/soundfx/libvisualizer.so \

_core_minimal_whitelist := vendor/lib/soundfx/libdownmix.so

_core_base_whitelist := \
  vendor/lib/libwebrtc_audio_preprocessing.so \
  vendor/lib/soundfx/libaudiopreprocessing.so \

_generic_no_telephony_whitelist := \
  vendor/lib/hw/audio.primary.default.so \
  vendor/lib/hw/local_time.default.so \
  vendor/lib/hw/power.default.so \
  vendor/lib/hw/vibrator.default.so \
  vendor/overlay/SysuiDarkTheme/SysuiDarkThemeOverlay.apk \


_my_whitelist := \
  $(_base_mk_whitelist) \
  $(_core_minimal_whitelist) \
  $(_generic_no_telephony_whitelist) \
  $(_core_base_whitelist) \
  recovery/root/etc/mke2fs.conf


$(call require-artifacts-in-path, $(TARGET_COPY_OUT_SYSTEM), $(_my_whitelist))
