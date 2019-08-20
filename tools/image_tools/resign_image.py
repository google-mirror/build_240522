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
import os
import subprocess
import sys

if len(sys.argv) < 6:
    image_path = os.environ["OUT"] + "/system.img"
    mount_dir = "/mnt/systemdir/"
    build_top = os.environ["ANDROID_BUILD_TOP"]
    key_path_prefix = "/vendor/sprd/partner/CusKeys.bak/"
    image_key_path = "/vendor/sprd/proprietories-source/packimage_scripts/signimage/sprd/config/rsa4096_system.pem"
else:
    image_path = sys.argv[1]
    mount_dir = sys.argv[2]
    build_top = sys.argv[3]
    key_path_prefix = sys.argv[4]
    image_key_path = sys.argv[5]

def execute(cmd):
    print(cmd)
    p = subprocess.Popen(cmd, shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    out, err = map(lambda b: b.decode('utf-8'), p.communicate())
    print(p.returncode, out, err)
    return p.returncode == 0, out, err

def make_sign_image_cmds(file, partition_size, key):
    return build_top + '/external/avb/avbtool add_hashtree_footer --partition_size ' + partition_size + ' --partition_name system' \
	+ ' --image ' + file + ' --key ' + key + ' --algorithm SHA256_RSA4096 --rollback_index 0 --setup_as_rootfs_from_kernel'


success, out, error_msg = execute('sudo ' + build_top + '/build/tools/image_tools/resign_apk.py ' \
                                  + image_path + ' ' + mount_dir + ' ' + build_top + ' ' + key_path_prefix + ' ' + image_key_path)
if success:
    sign_image_result, _, sign_image_error = execute(make_sign_image_cmds(image_path, str(os.path.getsize(image_path)), build_top + image_key_path))
    if sign_image_result:
        print("successfully", sign_image_result)
        sys.exit(0)
    else:
        print(sign_image_error, file=sys.stderr)
        sys.exit(1)
else:
    print(error_msg, file=sys.stderr)
    sys.exit(1)
