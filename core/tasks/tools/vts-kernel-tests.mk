<<<<<<< HEAD   (159713 Merge "Merge empty history for sparse-9081464-L8140000095646)
=======
# Copyright (C) 2022 The Android Open Source Project
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

-include external/linux-kselftest/android/kselftest_test_list.mk
-include external/ltp/android/ltp_package_list.mk

include $(BUILD_SYSTEM)/tasks/tools/vts_package_utils.mk

# Copy kernel test modules to testcases directories
kernel_ltp_host_out := $(HOST_OUT_TESTCASES)/vts_kernel_ltp_tests
kernel_ltp_vts_out := $(HOST_OUT)/$(test_suite_name)/android-$(test_suite_name)/testcases/vts_kernel_ltp_tests
kernel_ltp_modules := \
    ltp \
    $(ltp_packages)

kernel_kselftest_host_out := $(HOST_OUT_TESTCASES)/vts_kernel_kselftest_tests
kernel_kselftest_vts_out := $(HOST_OUT)/$(test_suite_name)/android-$(test_suite_name)/testcases/vts_kernel_kselftest_tests
kernel_kselftest_modules := $(kselftest_modules)
>>>>>>> BRANCH (5235f6 Merge "Version bump to TKB1.220921.001.A1 [core/build_id.mk])
