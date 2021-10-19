#
# Copyright 2021 The Android Open Source Project
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

# Inherit this in order to support Isolated Compilation in a protected VM.

PRODUCT_PACKAGES += \
    com.android.compos \
    com.android.virt \

PRODUCT_APEX_SYSTEM_SERVER_JARS += \
    com.android.compos:service-compos

PRODUCT_ARTIFACT_PATH_REQUIREMENT_ALLOWED_LIST += \
    system/apex/com.android.compos.apex \
    system/apex/com.android.virt.apex \
    system/framework/oat/%@service-compos.jar@classes.odex \
    system/framework/oat/%@service-compos.jar@classes.vdex \
    system/lib64/%.dylib.so \
    system/lib64/libgfxstream_backend.so \

PRODUCT_SYSTEM_PROPERTIES += \
    ro.config.isolated_compilation_enabled=true \
