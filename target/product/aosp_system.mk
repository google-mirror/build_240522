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

_recovery_whitelist := \
  recovery/root/etc/mke2fs.conf \
  recovery/root/sbin/adbd \
  recovery/root/sbin/e2fsdroid_static \
  recovery/root/sbin/mke2fs_static \
  recovery/root/sbin/recovery

_shell_and_utilities_whitelist := \
  vendor/bin/acpi \
  vendor/bin/awk \
  vendor/bin/base64 \
  vendor/bin/basename \
  vendor/bin/blockdev \
  vendor/bin/cal \
  vendor/bin/cat \
  vendor/bin/chcon \
  vendor/bin/chgrp \
  vendor/bin/chmod \
  vendor/bin/chown \
  vendor/bin/chroot \
  vendor/bin/chrt \
  vendor/bin/cksum \
  vendor/bin/clear \
  vendor/bin/cmp \
  vendor/bin/comm \
  vendor/bin/cp \
  vendor/bin/cpio \
  vendor/bin/cut \
  vendor/bin/date \
  vendor/bin/dd \
  vendor/bin/df \
  vendor/bin/diff \
  vendor/bin/dirname \
  vendor/bin/dmesg \
  vendor/bin/dos2unix \
  vendor/bin/du \
  vendor/bin/echo \
  vendor/bin/env \
  vendor/bin/expand \
  vendor/bin/expr \
  vendor/bin/fallocate \
  vendor/bin/false \
  vendor/bin/file \
  vendor/bin/find \
  vendor/bin/flock \
  vendor/bin/fmt \
  vendor/bin/free \
  vendor/bin/getenforce \
  vendor/bin/getevent \
  vendor/bin/getprop \
  vendor/bin/grep \
  vendor/bin/groups \
  vendor/bin/gunzip \
  vendor/bin/gzip \
  vendor/bin/head \
  vendor/bin/hostname \
  vendor/bin/hwclock \
  vendor/bin/id \
  vendor/bin/ifconfig \
  vendor/bin/inotifyd \
  vendor/bin/insmod \
  vendor/bin/ionice \
  vendor/bin/iorenice \
  vendor/bin/kill \
  vendor/bin/killall \
  vendor/bin/ln \
  vendor/bin/load_policy \
  vendor/bin/log \
  vendor/bin/logname \
  vendor/bin/logwrapper \
  vendor/bin/losetup \
  vendor/bin/ls \
  vendor/bin/lsmod \
  vendor/bin/lsof \
  vendor/bin/lspci \
  vendor/bin/lsusb \
  vendor/bin/md5sum \
  vendor/bin/microcom \
  vendor/bin/mkdir \
  vendor/bin/mkfifo \
  vendor/bin/mknod \
  vendor/bin/mkswap \
  vendor/bin/mktemp \
  vendor/bin/modinfo \
  vendor/bin/modprobe \
  vendor/bin/more \
  vendor/bin/mount \
  vendor/bin/mountpoint \
  vendor/bin/mv \
  vendor/bin/nc \
  vendor/bin/netcat \
  vendor/bin/netstat \
  vendor/bin/newfs_msdos \
  vendor/bin/nice \
  vendor/bin/nl \
  vendor/bin/nohup \
  vendor/bin/nsenter \
  vendor/bin/od \
  vendor/bin/paste \
  vendor/bin/patch \
  vendor/bin/pgrep \
  vendor/bin/pidof \
  vendor/bin/pkill \
  vendor/bin/pmap \
  vendor/bin/printenv \
  vendor/bin/printf \
  vendor/bin/ps \
  vendor/bin/pwd \
  vendor/bin/readlink \
  vendor/bin/realpath \
  vendor/bin/renice \
  vendor/bin/restorecon \
  vendor/bin/rm \
  vendor/bin/rmdir \
  vendor/bin/rmmod \
  vendor/bin/runcon \
  vendor/bin/sed \
  vendor/bin/sendevent \
  vendor/bin/seq \
  vendor/bin/setenforce \
  vendor/bin/setprop \
  vendor/bin/setsid \
  vendor/bin/sh \
  vendor/bin/sha1sum \
  vendor/bin/sha224sum \
  vendor/bin/sha256sum \
  vendor/bin/sha384sum \
  vendor/bin/sha512sum \
  vendor/bin/sleep \
  vendor/bin/sort \
  vendor/bin/split \
  vendor/bin/start \
  vendor/bin/stat \
  vendor/bin/stop \
  vendor/bin/strings \
  vendor/bin/stty \
  vendor/bin/swapoff \
  vendor/bin/swapon \
  vendor/bin/sync \
  vendor/bin/sysctl \
  vendor/bin/tac \
  vendor/bin/tail \
  vendor/bin/tar \
  vendor/bin/taskset \
  vendor/bin/tee \
  vendor/bin/time \
  vendor/bin/timeout \
  vendor/bin/toolbox \
  vendor/bin/top \
  vendor/bin/touch \
  vendor/bin/toybox_vendor \
  vendor/bin/tr \
  vendor/bin/true \
  vendor/bin/truncate \
  vendor/bin/tty \
  vendor/bin/ulimit \
  vendor/bin/umount \
  vendor/bin/uname \
  vendor/bin/uniq \
  vendor/bin/unix2dos \
  vendor/bin/unshare \
  vendor/bin/uptime \
  vendor/bin/usleep \
  vendor/bin/uudecode \
  vendor/bin/uuencode \
  vendor/bin/uuidgen \
  vendor/bin/vmstat \
  vendor/bin/wc \
  vendor/bin/which \
  vendor/bin/whoami \
  vendor/bin/xargs \
  vendor/bin/xxd \
  vendor/bin/yes \
  vendor/bin/zcat \
  vendor/etc/mkshrc

_selinux_policy_whitelist := \
  root/plat_file_contexts \
  root/plat_hwservice_contexts \
  root/plat_property_contexts \
  root/plat_seapp_contexts \
  root/plat_service_contexts \
  root/sepolicy \
  root/vendor_file_contexts \
  root/vendor_hwservice_contexts \
  root/vendor_property_contexts \
  root/vendor_seapp_contexts \
  root/vendor_service_contexts \
  root/vndservice_contexts \
  vendor/etc/selinux/plat_pub_versioned.cil \
  vendor/etc/selinux/plat_sepolicy_vers.txt \
  vendor/etc/selinux/precompiled_sepolicy \
  vendor/etc/selinux/precompiled_sepolicy.plat_and_mapping.sha256 \
  vendor/etc/selinux/vendor_mac_permissions.xml \
  vendor/etc/selinux/vendor_sepolicy.cil \

# The files not broken out into their own variables have a 1:1 mapping between
# file and module.
_my_whitelist := \
  $(_recovery_whitelist) \
  $(_from_shell_and_utilities) \
  $(_selinux_policy_whitelist) \
  root/init \
  root/init.environ.rc \
  root/init.rc \
  root/sbin/charger \
  vendor/bin/hw/android.hardware.cas@1.0-service \
  vendor/bin/hw/android.hardware.configstore@1.0-service \
  vendor/bin/hw/android.hardware.media.omx@1.0-service \
  vendor/bin/vndservice \
  vendor/bin/vndservicemanager \
  vendor/etc/fs_config_dirs \
  vendor/etc/fs_config_files \
  vendor/etc/init/android.hardware.cas@1.0-service.rc \
  vendor/etc/init/android.hardware.configstore@1.0-service.rc \
  vendor/etc/init/android.hardware.media.omx@1.0-service.rc \
  vendor/etc/init/vndservicemanager.rc \
  vendor/etc/vintf/compatibility_matrix.xml \
  vendor/lib/hw/audio.primary.default.so \
  vendor/lib/hw/gralloc.default.so \
  vendor/lib/hw/local_time.default.so \
  vendor/lib/hw/power.default.so \
  vendor/lib/hw/vibrator.default.so \
  vendor/lib/libavservices_minijail_vendor.so \
  vendor/lib/libeffects.so \
  vendor/lib/libeffectsconfig.so \
  vendor/lib/libhwminijail.so \
  vendor/lib/libmediacodecservice.so \
  vendor/lib/libreference-ril.so \
  vendor/lib/libril.so \
  vendor/lib/librilutils.so \
  vendor/lib/libwebrtc_audio_preprocessing.so \
  vendor/lib/mediacas/libclearkeycasplugin.so \
  vendor/lib/mediadrm/libdrmclearkeyplugin.so \
  vendor/lib/soundfx/libaudiopreprocessing.so \
  vendor/lib/soundfx/libbundlewrapper.so \
  vendor/lib/soundfx/libdownmix.so \
  vendor/lib/soundfx/libeffectproxy.so \
  vendor/lib/soundfx/libldnhncr.so \
  vendor/lib/soundfx/libreverbwrapper.so \
  vendor/lib/soundfx/libvisualizer.so \
  vendor/overlay/SysuiDarkTheme/SysuiDarkThemeOverlay.apk \



$(call make-isolation-claim, $(TARGET_COPY_OUT_SYSTEM), $(_my_whitelist))
