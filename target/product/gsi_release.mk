#
# Copyright (C) 2019 The Android Open-Source Project
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

#
# The makefile contains the special settings for GSI releasing.
# This makefile is used for the build targets which used for releasing GSI.
#
# For example:
# - Released GSI contains skip_mount.cfg to skip mounting prodcut paritition
# - Released GSI contains more VNDK packages to support old version vendors
# - etc.
#

# Exclude GSI specific files
PRODUCT_ARTIFACT_PATH_REQUIREMENT_WHITELIST += \
    system/etc/init/config/skip_mount.cfg \
    system/etc/init/init.gsi.rc \

# Some GSI builds enable dexpreopt, whitelist these preopt files
PRODUCT_ARTIFACT_PATH_REQUIREMENT_WHITELIST += %.odex %.vdex %.art

# Exclude all files under system/product and system/product_services
PRODUCT_ARTIFACT_PATH_REQUIREMENT_WHITELIST += \
    system/product/% \
    system/product_services/%


# GSI doesn't support apex for now.
# Properties set in product take precedence over those in vendor.
PRODUCT_PRODUCT_PROPERTIES += \
    ro.apex.updatable=false

# Split selinux policy
PRODUCT_FULL_TREBLE_OVERRIDE := true

# Enable dynamic partition size
PRODUCT_USE_DYNAMIC_PARTITION_SIZE := true

# Needed by Pi newly launched device to pass VtsTrebleSysProp on GSI
PRODUCT_COMPATIBLE_PROPERTY_OVERRIDE := true

# GSI specific tasks on boot
PRODUCT_COPY_FILES += \
    build/make/target/product/gsi/skip_mount.cfg:system/etc/init/config/skip_mount.cfg \
    build/make/target/product/gsi/init.gsi.rc:system/etc/init/init.gsi.rc \

# Support addtional P VNDK packages
PRODUCT_EXTRA_VNDK_VERSIONS := 28
