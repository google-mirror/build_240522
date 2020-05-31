#!/usr/bin/env python3
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
import argparse
import json
import os
import subprocess
import sys

from multiprocessing import Pool, cpu_count
from multiprocessing.dummy import Pool as ThreadPool

from collections import defaultdict
from glob import glob

# android.uid.phone
UID_PHONE_WHITELIST = (
  # Regular apps
  "MmsService.apk",
  "ONS.apk",
  "PresencePolling.apk",
  "RcsService.apk",
  "Stk.apk",
  "TeleService.apk",
  "TelephonyProvider.apk",

  # Test-only apps
  "CellBroadcastReceiverTests.apk",
  "PresencePollingTestHelper.apk",
)

# android.uid.system
UID_SYSTEM_WHITELIST = (
  # /system
  "DynamicSystemInstallationService.apk",
  "FusedLocation.apk",
  "InProcessNetworkStack.apk",
  "InputDevices.apk",
  "KeyChain.apk",
  "LocalTransport.apk",
  "SettingsProvider.apk",
  "Telecom.apk",
  "TvQuickSettings.apk",
  "TvSettings.apk",
  "WallpaperBackup.apk",
   # /product
  "Settings.apk",

  # Test-only apps
  "com.android.car.obd2.test.apk",
  "sl4a.apk",
  "AndroidCarApiTest.apk",
  "AppLaunch.apk",
  "AppLaunchWear.apk",
  "AppSmoke.apk",
  "AppCompatibilityTest.apk",
  "CarDeveloperOptions.apk",
  "CarService.apk",
  "CarServiceTest.apk",
  "CarServiceUnitTest.apk",
  "CarServicesTest.apk",
  "DataIdleTest.apk",
  "Development.apk",
  "DownloadManagerTestApp.apk",
  "FrameworksCoreFeatureFlagTests.apk",
  "EmbeddedKitchenSinkApp.apk",
  "FrameworksCorePackageManagerTests.apk",
  "FrameworksCoreSystemPropertiesTests.apk",
  "HdmiCecTests.apk",
  "StatsdLoadtest.apk",
  "StatsdDogfood.apk",
  "GarageModeTestApp.apk",
  "KeyChainTestsSupport.apk",
  "NetworkSecurityConfigTests.apk",
  "SystemUITests.apk",
  "TelecomUnitTests.apk",
  "TestablesTests.apk",
  "UsageReportingTest.apk",
  "UsageStatsTest.apk",
  "UxRestrictionsSample.apk",
  "VehicleHALTest.apk",
)

def parse_args():
    """Parse commandline arguments."""
    parser = argparse.ArgumentParser(description='Find sharedUserId violators')
    parser.add_argument('--product_out', help='PRODUCT_OUT directory',
                        default=os.environ.get("PRODUCT_OUT"))
    parser.add_argument('--aapt', help='Path to aapt or aapt2',
                        default="aapt2")
    parser.add_argument('--shipping_api', help='PRODUCT_SHIPPING_API_LEVEL',
                        default=0)
    parser.add_argument('--copy_out_system', help='TARGET_COPY_OUT_SYSTEM',
                        default="system")
    parser.add_argument('--copy_out_vendor', help='TARGET_COPY_OUT_VENDOR',
                        default="vendor")
    parser.add_argument('--copy_out_product', help='TARGET_COPY_OUT_PRODUCT',
                        default="product")
    parser.add_argument('--copy_out_system_ext', help='TARGET_COPY_OUT_SYSTEM_EXT',
                        default="system_ext")
    return parser.parse_args()

def execute(cmd):
    p = subprocess.Popen(cmd, shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    out, err = map(lambda b: b.decode('utf-8'), p.communicate())
    return p.returncode == 0, out, err

def main():
    # Set up multiprocessing
    threads = (cpu_count() // 2) or 1
    pool = ThreadPool(threads)

    args = parse_args()

    shipping_api = int(args.shipping_api)
    product_out = args.product_out
    aapt = args.aapt

    partitions = (
            ("system", args.copy_out_system),
            ("vendor", args.copy_out_vendor),
            ("product", args.copy_out_product),
            ("system_ext", args.copy_out_system_ext),
    )

    apks = []
    for part, location in partitions:
        # Match only app and priv-app
        for f in glob(os.path.join(product_out, location, "*app", "**", "*.apk")):
            apks.append((part, f))

    shareduid_app_dict = defaultdict(list)

    def make_aapt_cmds(filename):
        return [aapt + ' dump ' + filename + ' --file AndroidManifest.xml',
                aapt + ' dump xmltree ' + filename + ' --file AndroidManifest.xml']

    def extract_shared_uid(apk_file):
        for cmd in make_aapt_cmds(apk_file):
            success, manifest, error_msg = execute(cmd)
            if success:
                break
        else:
            print(error_msg, file=sys.stderr)
            sys.exit()

        for l in manifest.split('\n'):
            if "sharedUserId" in l:
                return l.split('"')[-2]
        return None

    def process_apk(appinfo):
        apk_file = os.path.basename(appinfo[1])
        shared_uid = extract_shared_uid(appinfo[1])
        if shared_uid is None:
            return
        shareduid_app_dict[shared_uid].append((appinfo[0], apk_file))

    pool.map(process_apk, apks)
    pool.close()
    pool.join()

    output = defaultdict(lambda: defaultdict(list))

    violators = []

    for uid, app_infos in shareduid_app_dict.items():
        _partitions = {p for p, _ in app_infos}
        if len(_partitions) <= 1:
            continue
        for part in _partitions:
            # Sort by apk name
            for p, a in sorted(app_infos, key=lambda x: x[1]):
                if p != part:
                    continue
                output[uid][part].extend([a])
                # Collect violators of banned UIDs
                if uid == "android.uid.phone" and a not in UID_PHONE_WHITELIST:
                    violators.append([a, uid])
                if uid == "android.uid.system" and a not in UID_SYSTEM_WHITELIST:
                    violators.append([a, uid])

    def errorfmt(err):
        return '\033[31m%s\033[0m\n' % err

    if len(violators) > 0 and shipping_api >= 31:
        sys.stderr.write(errorfmt("Error: The following applications are"
            " violating the ban on using system-reserved sharedUserIds:"))
        for _app, _uid in violators:
            sys.stderr.write(f"{_app} for sharedUserId {_uid}\n")
        sys.stderr.write("\nPlease remove the sharedUserId attribute from the"
            " offending application's AndroidManifest.xml\n")
        sys.exit(1)

    print(json.dumps(output, indent=2, sort_keys=True))

if __name__ == '__main__':
    main()
