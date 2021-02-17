<<<<<<< HEAD   (4be654 Merge "Merge empty history for sparse-7121469-L4290000080720)
=======
# Copyright (C) 2020 The Android Open Source Project
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

current_makefile := $(lastword $(MAKEFILE_LIST))

# BOARD_VNDK_VERSION must be set to 'current' in order to generate a vendor snapshot.
ifeq ($(BOARD_VNDK_VERSION),current)

.PHONY: vendor-snapshot
vendor-snapshot: $(SOONG_VENDOR_SNAPSHOT_ZIP)

$(call dist-for-goals, vendor-snapshot, $(SOONG_VENDOR_SNAPSHOT_ZIP))

.PHONY: vendor-fake-snapshot
vendor-fake-snapshot: $(SOONG_VENDOR_FAKE_SNAPSHOT_ZIP)

$(call dist-for-goals, vendor-fake-snapshot, $(SOONG_VENDOR_FAKE_SNAPSHOT_ZIP):fake/$(notdir $(SOONG_VENDOR_FAKE_SNAPSHOT_ZIP)))

else # BOARD_VNDK_VERSION is NOT set to 'current'

.PHONY: vendor-snapshot
vendor-snapshot: PRIVATE_MAKEFILE := $(current_makefile)
vendor-snapshot:
	$(call echo-error,$(PRIVATE_MAKEFILE),\
		"CANNOT generate Vendor snapshot. BOARD_VNDK_VERSION must be set to 'current'.")
	exit 1

.PHONY: vendor-fake-snapshot
vendor-fake-snapshot: PRIVATE_MAKEFILE := $(current_makefile)
vendor-fake-snapshot:
	$(call echo-error,$(PRIVATE_MAKEFILE),\
		"CANNOT generate Vendor snapshot. BOARD_VNDK_VERSION must be set to 'current'.")
	exit 1

endif # BOARD_VNDK_VERSION
>>>>>>> BRANCH (fe6ad7 Merge "Version bump to RBT1.210107.001.A1 [core/build_id.mk])
