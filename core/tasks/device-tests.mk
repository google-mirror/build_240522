# Copyright (C) 2017 The Android Open Source Project
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agrls eed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


.PHONY: device-tests
device-tests-zip := device-tests.zip
$(device-tests-zip): $(COMPATIBILITY.device-tests.FILES) $(SOONG_ZIP)
	$(hide) find $(HOST_OUT_TESTCASES) $(TARGET_OUT_TESTCASES) | sort >$(OUT_DIR)/$@.list
	$(hide) $(SOONG_ZIP) -d -o $(OUT_DIR)/$@ -C $(dir $@) -l $(OUT_DIR)/$@.list

device-tests: $(device-tests-zip)
$(call dist-for-goals, device-tests, $(device-tests-zip))
