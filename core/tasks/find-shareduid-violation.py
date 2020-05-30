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

from collections import defaultdict
from glob import glob

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

    shareduid_app_dict = defaultdict(list)

    def make_aapt_cmds(filename):
        return [aapt + ' dump ' + filename + ' --file AndroidManifest.xml',
                aapt + ' dump xmltree ' + filename + ' --file AndroidManifest.xml']

    def extract_shared_uid(file):
        for cmd in make_aapt_cmds(file):
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

    for part, location in partitions:
        for f in glob(os.path.join(product_out, location, "*", "*", "*.apk")):
            apk_file = os.path.basename(f)
            shared_uid = extract_shared_uid(f)

            if shared_uid is None:
                continue
            shareduid_app_dict[shared_uid].append((part, apk_file))

    output = defaultdict(lambda: defaultdict(list))

    # android.uid.phone
    uid_phone_whitelist = (
      "Stk.apk",
      "MmsService.apk",
      "ONS.apk",
      "TeleService.apk",
      "TelephonyProvider.apk",
    )
    # android.uid.system
    uid_system_whitelist = (
      # system
      "KeyChain.apk",
      "WallpaperBackup.apk",
      "DynamicSystemInstallationService.apk",
      "FusedLocation.apk",
      "LocalTransport.apk",
      "InputDevices.apk",
      "SettingsProvider.apk",
      "Telecom.apk",
       # product
      "Settings.apk",
    )

    violators = []

    for uid, app_infos in shareduid_app_dict.items():
        _partitions = {p for p, _ in app_infos}
        if len(_partitions) <= 1:
            continue
        for part in _partitions:
            for p, a in app_infos:
                if p != part:
                    continue
                if uid == "android.uid.phone" and a not in uid_phone_whitelist:
                    violators.append([a, uid])
                if uid == "android.uid.system" and a not in uid_system_whitelist:
                    violators.append([a, uid])
                output[uid][part].extend([a])

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
