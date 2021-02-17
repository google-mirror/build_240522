<<<<<<< HEAD   (4be654 Merge "Merge empty history for sparse-7121469-L4290000080720)
=======
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

ifneq ($(wildcard test/mts/README.md),)

mts_test_suites :=
mts_test_suites += mts

$(foreach module, $(mts_modules), $(eval mts_test_suites += mts-$(module)))

$(foreach suite, $(mts_test_suites), \
	$(eval test_suite_name := $(suite)) \
	$(eval test_suite_tradefed := mts-tradefed) \
	$(eval test_suite_readme := test/mts/README.md) \
	$(eval include $(BUILD_SYSTEM)/tasks/tools/compatibility.mk) \
	$(eval .PHONY: $(suite)) \
	$(eval $(suite): $(compatibility_zip)) \
	$(eval $(call dist-for-goals, $(suite), $(compatibility_zip))) \
)

endif
>>>>>>> BRANCH (fe6ad7 Merge "Version bump to RBT1.210107.001.A1 [core/build_id.mk])
