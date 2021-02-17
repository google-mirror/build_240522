<<<<<<< HEAD   (4be654 Merge "Merge empty history for sparse-7121469-L4290000080720)
=======
#
# Copyright (C) 2019 The Android Open Source Project
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

# Inherit this when the target needs to support updating APEXes

ifneq ($(OVERRIDE_TARGET_FLATTEN_APEX),true)
  # com.android.apex.cts.shim.v1_prebuilt overrides CtsShimPrebuilt
  # and CtsShimPrivPrebuilt since they are packaged inside the APEX.
  PRODUCT_PACKAGES += com.android.apex.cts.shim.v1_prebuilt
  PRODUCT_VENDOR_PROPERTIES := ro.apex.updatable=true
  TARGET_FLATTEN_APEX := false
endif
>>>>>>> BRANCH (fe6ad7 Merge "Version bump to RBT1.210107.001.A1 [core/build_id.mk])
