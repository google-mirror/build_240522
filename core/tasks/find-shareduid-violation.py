#!/usr/bin/env python3
import os
import subprocess
from glob import glob
import re
from collections import defaultdict
import sys
import json

if len(sys.argv) < 3:
    product_out = os.environ["ANDROID_PRODUCT_OUT"]
    aapt = "aapt2"
else:
    product_out = sys.argv[1]
    aapt = sys.argv[2]


def extract_shared_uid(file):
    manifest = subprocess.check_output(
        aapt + ' dump --file AndroidManifest.xml ' + file, shell=True).decode().split('\n')
    for l in manifest:
        if "sharedUserId" in l:
            return l.split('"')[-2]
    return None


partitions = ["system", "vendor", "product"]

shareduid_app_dict = defaultdict(list)

for p in partitions:
    for f in glob(os.path.join(product_out, p, "**", "*.apk"), recursive=True):
        apk_file = os.path.basename(f)
        shared_uid = extract_shared_uid(f)

        if shared_uid is None:
            continue
        shareduid_app_dict[shared_uid].append((p, apk_file))


output = defaultdict(lambda: defaultdict(list))

for uid, app_infos in shareduid_app_dict.items():
    partitions = {p for p, _ in app_infos}
    if len(partitions) > 1:
        for part in partitions:
            output[uid][part].extend([a for p, a in app_infos if p == part])

print(json.dumps(output, indent=2, sort_keys=True))
