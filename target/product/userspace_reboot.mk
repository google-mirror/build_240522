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

# Inherit this when the target supports userspace reboot

PRODUCT_PROPERTY_OVERRIDES += \
	ro.init.userspace_reboot.is_supported=true \
	ro.init.userspace_reboot.sigkill_timeout=10000 \ # 10 seconds
	ro.init.userspace_reboot.sigterm_timeout=5000 \ # 5 seconds
	ro.init.userspace_reboot.watchdog_timeout=300000 \ # 5 minutes

# TODO(b/135984674): configure userspace reboot related properties
