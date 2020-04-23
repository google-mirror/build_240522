#!/usr/bin/env python3

import os
import subprocess
import sys

if len(sys.argv) < 6:
    image_path = os.environ["OUT"] + "/system.img"
    mount_dir = "systemdir/"
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

def make_vbmeta_image_cmds(key):
    return build_top + '/external/avb/avbtool make_vbmeta_image --algorithm SHA256_RSA4096 --key ' + key + ' --padding_size 4096' \
        + ' --rollback_index 0 --include_descriptors_from_image ' + os.environ["OUT"] + "/system.img" + ' --include_descriptors_from_image ' \
        + os.environ["OUT"] + "/product.img " + ' --output ' + os.environ["OUT"] + "/vbmeta_system.img"

def build_super_image_cmds():
    return build_top + '/build/make/tools/releasetools/build_super_image.py ' + ' -v ' + os.environ["OUT"]  + '/obj/PACKAGING/superimage_debug_intermediates/misc_info.txt ' \
        + os.environ["OUT"] + '/super.img'

execute(build_top + '/out/host/linux-x86/bin/simg2img ' + image_path + ' ' + image_path + '.bak')
success, out, error_msg = execute('sudo ' + build_top + '/build/tools/image_tools/resign_apk.py ' \
                                  + image_path + ' ' + mount_dir + ' ' + build_top + ' ' + key_path_prefix + ' ' + image_key_path)

if success:
    #execute('cp ' + image_path + '.bak ' + image_path + '.bk')
    sign_image_result, _, sign_image_error = execute(make_sign_image_cmds(image_path + '.bak', str(os.path.getsize(image_path + '.bak')), build_top + image_key_path))
    execute(build_top + '/out/host/linux-x86/bin/img2simg ' + image_path + '.bak ' + os.environ["OUT"] + "/system.img")
    make_vbmeta_image_result, _, make_vbmeta_image_error = execute(make_vbmeta_image_cmds(build_top + image_key_path))
    build_image_result, _, build_image_error = execute(build_super_image_cmds())
    if build_image_result:
        print("successfully", build_image_result)
        sys.exit(0)
    else:
        print(build_image_error, file=sys.stderr)
        sys.exit(1)
else:
    print(error_msg, file=sys.stderr)
    sys.exit(1)
