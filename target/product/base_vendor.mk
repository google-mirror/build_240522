#
# Copyright (C) 2018 The Android Open Source Project
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

# Base modules and settings for recovery.
PRODUCT_PACKAGES += \
    adbd.recovery \
    init_second_stage.recovery \
    ld.config.recovery.txt \
    linker.recovery \
    recovery \
    shell_and_utilities_recovery \
    watchdogd.recovery \

# Base modules and settings for the vendor partition.
PRODUCT_PACKAGES += \
    android.hardware.cas@1.0-service \
    android.hardware.configstore@1.1-service \
    android.hardware.media.omx@1.0-service \
    fs_config_files_nonsystem \
    fs_config_dirs_nonsystem \
    gralloc.default \
    group \
    libbundlewrapper \
    libclearkeycasplugin \
    libdownmix \
    libdrmclearkeyplugin \
    libdynproc \
    libeffectproxy \
    libeffects \
    libldnhncr \
    libreference-ril \
    libreverbwrapper \
    libril \
    libvisualizer \
    passwd \
    selinux_policy_nonsystem \
    shell_and_utilities_vendor \
    vndservice \
    vndservicemanager \

# VINTF data for vendor image
PRODUCT_PACKAGES += \
    device_manifest.xml \
    device_compatibility_matrix.xml \
